package service

import (
	"context"
	"errors"

	"blog/internal/dao"
	"blog/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type sPermission struct{}

var Permission = sPermission{}

func (s *sPermission) Create(ctx context.Context, in *entity.Permission) (id int, err error) {
	result, err := dao.Permission.Ctx(ctx).Data(in).Insert()
	if err != nil {
		return 0, err
	}
	id64, _ := result.LastInsertId()
	return int(id64), nil
}

func (s *sPermission) Update(ctx context.Context, in *entity.Permission) (err error) {
	data := g.Map{}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.Code != "" {
		data["code"] = in.Code
	}
	if in.Type != "" {
		data["type"] = in.Type
	}
	if in.Path != "" {
		data["path"] = in.Path
	}
	if in.Icon != "" {
		data["icon"] = in.Icon
	}
	if in.ParentId >= 0 {
		data["parent_id"] = in.ParentId
	}
	if in.SortOrder >= 0 {
		data["sort_order"] = in.SortOrder
	}
	if in.Status >= 0 {
		data["status"] = in.Status
	}
	_, err = dao.Permission.Ctx(ctx).Data(data).Where(dao.Permission.Columns().Id, in.Id).Update()
	return
}

func (s *sPermission) GetList(ctx context.Context, parentId int) (list []*entity.Permission, err error) {
	m := dao.Permission.Ctx(ctx).Where(dao.Permission.Columns().Status, 1)
	if parentId >= 0 {
		m = m.Where(dao.Permission.Columns().ParentId, parentId)
	}
	err = m.OrderAsc(dao.Permission.Columns().SortOrder).OrderAsc(dao.Permission.Columns().Id).Scan(&list)
	return
}

func (s *sPermission) GetTree(ctx context.Context) (list []map[string]interface{}, err error) {
	var all []*entity.Permission
	err = dao.Permission.Ctx(ctx).Where(dao.Permission.Columns().Status, 1).OrderAsc(dao.Permission.Columns().SortOrder).OrderAsc(dao.Permission.Columns().Id).Scan(&all)
	if err != nil {
		return
	}

	// 构建树形结构
	permissionMap := make(map[int]*entity.Permission)
	for _, p := range all {
		permissionMap[p.Id] = p
	}

	var roots []map[string]interface{}
	for _, p := range all {
		if p.ParentId == 0 {
			item := s.permissionToMap(p)
			item["children"] = s.buildPermissionTree(p.Id, permissionMap)
			roots = append(roots, item)
		}
	}

	return roots, nil
}

func (s *sPermission) buildPermissionTree(parentId int, permissionMap map[int]*entity.Permission) []map[string]interface{} {
	var children []map[string]interface{}
	for _, p := range permissionMap {
		if p.ParentId == parentId {
			item := s.permissionToMap(p)
			item["children"] = s.buildPermissionTree(p.Id, permissionMap)
			children = append(children, item)
		}
	}
	return children
}

func (s *sPermission) permissionToMap(p *entity.Permission) map[string]interface{} {
	return map[string]interface{}{
		"id":        p.Id,
		"parentId":  p.ParentId,
		"name":      p.Name,
		"code":      p.Code,
		"type":      p.Type,
		"path":      p.Path,
		"icon":      p.Icon,
		"sortOrder": p.SortOrder,
		"status":    p.Status,
	}
}

func (s *sPermission) GetOne(ctx context.Context, id int) (permission *entity.Permission, err error) {
	err = dao.Permission.Ctx(ctx).Where(dao.Permission.Columns().Id, id).Scan(&permission)
	return
}

func (s *sPermission) Delete(ctx context.Context, id int) (err error) {
	// 检查是否有子权限
	var count int
	count, err = dao.Permission.Ctx(ctx).Where(dao.Permission.Columns().ParentId, id).Count()
	if err != nil {
		return
	}
	if count > 0 {
		return errors.New("存在子权限，无法删除")
	}
	// 删除角色权限关联
	_, err = dao.RolePermission.Ctx(ctx).Where(dao.RolePermission.Columns().PermissionId, id).Delete()
	if err != nil {
		return
	}
	// 删除权限
	_, err = dao.Permission.Ctx(ctx).Where(dao.Permission.Columns().Id, id).Delete()
	return
}

