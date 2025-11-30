package service

import (
	"context"
	"errors"
	"sort"

	"blog/internal/dao"
	"blog/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type sMenu struct{}

var Menu = sMenu{}

func (s *sMenu) Create(ctx context.Context, in *entity.Menu) (id int, err error) {
	result, err := dao.Menu.Ctx(ctx).Data(in).Insert()
	if err != nil {
		return 0, err
	}
	id64, _ := result.LastInsertId()
	return int(id64), nil
}

func (s *sMenu) Update(ctx context.Context, in *entity.Menu) (err error) {
	data := g.Map{}
	if in.Title != "" {
		data["title"] = in.Title
	}
	if in.Type != "" {
		data["type"] = in.Type
	}
	if in.Icon != "" {
		data["icon"] = in.Icon
	}
	if in.Path != "" {
		data["path"] = in.Path
	}
	if in.Component != "" {
		data["component"] = in.Component
	}
	if in.Permission != "" {
		data["permission"] = in.Permission
	}
	if in.ParentId >= 0 {
		data["parent_id"] = in.ParentId
	}
	if in.Order >= 0 {
		data["order"] = in.Order
	}
	if in.Status != "" {
		data["status"] = in.Status
	}
	_, err = dao.Menu.Ctx(ctx).Data(data).Where(dao.Menu.Columns().Id, in.Id).Update()
	return
}

func (s *sMenu) GetList(ctx context.Context) (list []*entity.Menu, err error) {
	err = dao.Menu.Ctx(ctx).OrderAsc(dao.Menu.Columns().Order).OrderAsc(dao.Menu.Columns().Id).Scan(&list)
	return
}

func (s *sMenu) GetTree(ctx context.Context) (list []map[string]interface{}, err error) {
	var all []*entity.Menu
	err = dao.Menu.Ctx(ctx).Where(dao.Menu.Columns().Status, "active").OrderAsc(dao.Menu.Columns().Order).OrderAsc(dao.Menu.Columns().Id).Scan(&all)
	if err != nil {
		return
	}

	// 构建树形结构
	menuMap := make(map[int]*entity.Menu)
	for _, m := range all {
		menuMap[m.Id] = m
	}

	// 按order排序的根菜单
	var rootMenus []*entity.Menu
	for _, m := range all {
		if m.ParentId == 0 {
			rootMenus = append(rootMenus, m)
		}
	}

	// 按order排序，order相同时按id排序
	sort.Slice(rootMenus, func(i, j int) bool {
		if rootMenus[i].Order != rootMenus[j].Order {
			return rootMenus[i].Order < rootMenus[j].Order
		}
		return rootMenus[i].Id < rootMenus[j].Id
	})

	var roots []map[string]interface{}
	for _, m := range rootMenus {
		item := s.menuToMap(m)
		item["children"] = s.buildMenuTree(m.Id, menuMap)
		roots = append(roots, item)
	}

	return roots, nil
}

func (s *sMenu) buildMenuTree(parentId int, menuMap map[int]*entity.Menu) []map[string]interface{} {
	// 先收集所有子菜单
	var childMenus []*entity.Menu
	for _, m := range menuMap {
		if m.ParentId == parentId {
			childMenus = append(childMenus, m)
		}
	}

	// 按order排序，order相同时按id排序
	sort.Slice(childMenus, func(i, j int) bool {
		if childMenus[i].Order != childMenus[j].Order {
			return childMenus[i].Order < childMenus[j].Order
		}
		return childMenus[i].Id < childMenus[j].Id
	})

	// 构建树
	var children []map[string]interface{}
	for _, m := range childMenus {
		item := s.menuToMap(m)
		item["children"] = s.buildMenuTree(m.Id, menuMap)
		children = append(children, item)
	}
	return children
}

func (s *sMenu) menuToMap(m *entity.Menu) map[string]interface{} {
	return map[string]interface{}{
		"id":         m.Id,
		"parentId":   m.ParentId,
		"type":       m.Type,
		"title":      m.Title,
		"icon":       m.Icon,
		"path":       m.Path,
		"component":  m.Component,
		"permission": m.Permission,
		"order":      m.Order,
		"status":     m.Status,
	}
}

func (s *sMenu) GetOne(ctx context.Context, id int) (menu *entity.Menu, err error) {
	err = dao.Menu.Ctx(ctx).Where(dao.Menu.Columns().Id, id).Scan(&menu)
	return
}

func (s *sMenu) Delete(ctx context.Context, id int) (err error) {
	// 检查是否有子菜单
	var count int
	count, err = dao.Menu.Ctx(ctx).Where(dao.Menu.Columns().ParentId, id).Count()
	if err != nil {
		return
	}
	if count > 0 {
		return errors.New("存在子菜单，无法删除")
	}
	_, err = dao.Menu.Ctx(ctx).Where(dao.Menu.Columns().Id, id).Delete()
	return
}

// 根据权限代码过滤菜单树
func (s *sMenu) GetTreeByPermissions(ctx context.Context, permissionCodes []string) (list []map[string]interface{}, err error) {
	// 获取所有菜单
	allMenus, err := s.GetTree(ctx)
	if err != nil {
		return
	}

	// 如果没有权限，返回空列表
	if len(permissionCodes) == 0 {
		return []map[string]interface{}{}, nil
	}

	// 创建权限代码集合（用于快速查找）
	permissionSet := make(map[string]bool)
	for _, code := range permissionCodes {
		permissionSet[code] = true
	}

	// 递归过滤菜单
	var filterMenu func(items []map[string]interface{}) []map[string]interface{}
	filterMenu = func(items []map[string]interface{}) []map[string]interface{} {
		var result []map[string]interface{}
		for _, item := range items {
			// 检查菜单类型和权限要求
			menuType, _ := item["type"].(string)
			permission, ok := item["permission"].(string)
			hasPermission := ok && permission != "" && permissionSet[permission]
			noPermissionRequired := !ok || permission == ""

			// 如果有子菜单，先递归过滤子菜单
			if children, ok := item["children"].([]map[string]interface{}); ok && len(children) > 0 {
				filteredChildren := filterMenu(children)
				// 只有当过滤后有子菜单时才保留父菜单
				if len(filteredChildren) > 0 {
					// 创建新的item副本，避免修改原始数据，保持顺序
					newItem := make(map[string]interface{})
					for k, v := range item {
						newItem[k] = v
					}
					newItem["children"] = filteredChildren

					// 对于目录类型（M），如果没有权限要求，直接显示；如果有权限要求，必须有权限
					if menuType == "M" {
						if noPermissionRequired || hasPermission {
							result = append(result, newItem)
						}
					} else {
						// 非目录类型，只要有权限或没有权限要求就显示
						if hasPermission || noPermissionRequired {
							result = append(result, newItem)
						}
					}
				}
			} else {
				// 没有子菜单，检查权限
				// 目录类型（M）必须有权限或没有权限要求
				// 菜单类型（C）或按钮类型（F），必须有权限或没有权限要求
				if menuType == "M" {
					// 目录类型如果没有子菜单，不应该显示
					continue
				} else if hasPermission || noPermissionRequired {
					// 创建新对象，保持顺序
					newItem := make(map[string]interface{})
					for k, v := range item {
						newItem[k] = v
					}
					result = append(result, newItem)
				}
			}
		}
		return result
	}

	list = filterMenu(allMenus)
	return
}
