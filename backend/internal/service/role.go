package service

import (
	"context"

	"blog/internal/dao"
	"blog/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type sRole struct{}

var Role = sRole{}

func (s *sRole) Create(ctx context.Context, in *entity.Role) (id int, err error) {
	result, err := dao.Role.Ctx(ctx).Data(in).Insert()
	if err != nil {
		return 0, err
	}
	id64, _ := result.LastInsertId()
	return int(id64), nil
}

func (s *sRole) Update(ctx context.Context, in *entity.Role) (err error) {
	data := g.Map{}
	if in.Name != "" {
		data["name"] = in.Name
	}
	if in.Key != "" {
		data["key"] = in.Key
	}
	if in.Desc != "" {
		data["desc"] = in.Desc
	}
	if in.Status >= 0 {
		data["status"] = in.Status
	}
	_, err = dao.Role.Ctx(ctx).Data(data).Where(dao.Role.Columns().Id, in.Id).Update()
	return
}

func (s *sRole) GetList(ctx context.Context, name string, page, size int) (list []*entity.Role, total int, err error) {
	m := dao.Role.Ctx(ctx)

	if name != "" {
		m = m.WhereLike(dao.Role.Columns().Name, "%"+name+"%")
	}

	total, err = m.Count()
	if err != nil {
		return
	}

	err = m.OrderDesc(dao.Role.Columns().Id).Page(page, size).Scan(&list)
	return
}

func (s *sRole) GetOne(ctx context.Context, id int) (role *entity.Role, err error) {
	err = dao.Role.Ctx(ctx).Where(dao.Role.Columns().Id, id).Scan(&role)
	return
}

// 通过key查找角色
func (s *sRole) GetByKey(ctx context.Context, key string) (role *entity.Role, err error) {
	err = dao.Role.Ctx(ctx).Where(dao.Role.Columns().Key, key).Scan(&role)
	return
}

func (s *sRole) Delete(ctx context.Context, id int) (err error) {
	// 删除角色权限关联
	_, err = dao.RolePermission.Ctx(ctx).Where(dao.RolePermission.Columns().RoleId, id).Delete()
	if err != nil {
		return
	}
	// 删除用户角色关联
	_, err = dao.UserRole.Ctx(ctx).Where(dao.UserRole.Columns().RoleId, id).Delete()
	if err != nil {
		return
	}
	// 删除角色
	_, err = dao.Role.Ctx(ctx).Where(dao.Role.Columns().Id, id).Delete()
	return
}

// 设置角色权限
func (s *sRole) SetPermissions(ctx context.Context, roleId int, permissionIds []int) (err error) {
	// 先删除原有权限
	_, err = dao.RolePermission.Ctx(ctx).Where(dao.RolePermission.Columns().RoleId, roleId).Delete()
	if err != nil {
		return
	}
	// 插入新权限
	if len(permissionIds) > 0 {
		for _, pid := range permissionIds {
			_, err = dao.RolePermission.Ctx(ctx).Data(g.Map{
				"role_id":       roleId,
				"permission_id": pid,
			}).Insert()
			if err != nil {
				return
			}
		}
	}
	return
}

// 获取角色权限
func (s *sRole) GetPermissions(ctx context.Context, roleId int) (permissionIds []int, err error) {
	var list []struct {
		PermissionId int `json:"permissionId"`
	}
	err = dao.RolePermission.Ctx(ctx).Fields(dao.RolePermission.Columns().PermissionId).Where(dao.RolePermission.Columns().RoleId, roleId).Scan(&list)
	if err != nil {
		return
	}
	for _, item := range list {
		permissionIds = append(permissionIds, item.PermissionId)
	}
	return
}
