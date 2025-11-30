package service

import (
	"context"
	"errors"

	"blog/internal/dao"
	"blog/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type sUser struct{}

var User = sUser{}

func (s *sUser) Login(ctx context.Context, username, password string) (user *entity.User, err error) {
	err = dao.User.Ctx(ctx).Where(dao.User.Columns().Username, username).Where(dao.User.Columns().Password, password).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户名或密码错误")
	}
	return user, nil
}

func (s *sUser) Create(ctx context.Context, in *entity.User) (id int, err error) {
	// 构建数据
	data := g.Map{
		"username": in.Username,
		"password": in.Password,
		"email":    in.Email, // 邮箱必填
		"status":   in.Status,
	}
	if in.Nickname != "" {
		data["nickname"] = in.Nickname
	}
	if in.Avatar != "" {
		data["avatar"] = in.Avatar
	}
	result, err := dao.User.Ctx(ctx).Data(data).Insert()
	if err != nil {
		return 0, err
	}
	id64, _ := result.LastInsertId()
	return int(id64), nil
}

func (s *sUser) Update(ctx context.Context, in *entity.User) (err error) {
	data := g.Map{}
	if in.Username != "" {
		data["username"] = in.Username
	}
	if in.Password != "" {
		data["password"] = in.Password
	}
	if in.Nickname != "" {
		data["nickname"] = in.Nickname
	}
	if in.Email != "" {
		data["email"] = in.Email
	}
	if in.Avatar != "" {
		data["avatar"] = in.Avatar
	}
	if in.Status >= 0 {
		data["status"] = in.Status
	}
	_, err = dao.User.Ctx(ctx).Data(data).Where(dao.User.Columns().Id, in.Id).Update()
	return
}

func (s *sUser) GetList(ctx context.Context, username string, page, size int) (list []*entity.User, total int, err error) {
	m := dao.User.Ctx(ctx)

	if username != "" {
		m = m.WhereLike(dao.User.Columns().Username, "%"+username+"%")
	}

	total, err = m.Count()
	if err != nil {
		return
	}

	err = m.OrderDesc(dao.User.Columns().Id).Page(page, size).Scan(&list)
	return
}

func (s *sUser) GetOne(ctx context.Context, id int) (user *entity.User, err error) {
	err = dao.User.Ctx(ctx).Where(dao.User.Columns().Id, id).Scan(&user)
	return
}

func (s *sUser) Delete(ctx context.Context, id int) (err error) {
	_, err = dao.User.Ctx(ctx).Where(dao.User.Columns().Id, id).Delete()
	return
}

// 设置用户角色
func (s *sUser) SetRoles(ctx context.Context, userId int, roleIds []int) (err error) {
	// 先删除原有角色
	_, err = dao.UserRole.Ctx(ctx).Where(dao.UserRole.Columns().UserId, userId).Delete()
	if err != nil {
		return
	}
	// 插入新角色
	if len(roleIds) > 0 {
		for _, rid := range roleIds {
			_, err = dao.UserRole.Ctx(ctx).Data(g.Map{
				"user_id": userId,
				"role_id": rid,
			}).Insert()
			if err != nil {
				return
			}
		}
	}
	return
}

// 分配默认user角色
func (s *sUser) AssignDefaultUserRole(ctx context.Context, userId int) (err error) {
	// 查找user角色
	userRole, err := Role.GetByKey(ctx, "user")
	if err != nil || userRole == nil {
		// 如果user角色不存在，创建它
		userRole = &entity.Role{
			Name:   "用户",
			Key:    "user",
			Desc:   "普通用户角色",
			Status: 1,
		}
		roleId, err := Role.Create(ctx, userRole)
		if err != nil {
			return err
		}
		userRole.Id = roleId
	}

	// 分配角色给用户
	return s.SetRoles(ctx, userId, []int{userRole.Id})
}

// 获取用户角色
func (s *sUser) GetRoles(ctx context.Context, userId int) (roleIds []int, err error) {
	var list []struct {
		RoleId int `json:"roleId"`
	}
	err = dao.UserRole.Ctx(ctx).Fields(dao.UserRole.Columns().RoleId).Where(dao.UserRole.Columns().UserId, userId).Scan(&list)
	if err != nil {
		return
	}
	for _, item := range list {
		roleIds = append(roleIds, item.RoleId)
	}
	return
}

// 判断用户是否是admin角色
func (s *sUser) IsAdmin(ctx context.Context, userId int) (isAdmin bool, err error) {
	// 获取admin角色
	adminRole, err := Role.GetByKey(ctx, "admin")
	if err != nil || adminRole == nil {
		return false, nil
	}

	// 获取用户角色
	roleIds, err := s.GetRoles(ctx, userId)
	if err != nil {
		return false, err
	}

	// 检查用户是否拥有admin角色
	for _, roleId := range roleIds {
		if roleId == adminRole.Id {
			return true, nil
		}
	}

	return false, nil
}

// 获取用户权限（通过角色）
func (s *sUser) GetPermissions(ctx context.Context, userId int) (permissionCodes []string, err error) {
	// 获取用户角色
	roleIds, err := s.GetRoles(ctx, userId)
	if err != nil {
		return
	}
	if len(roleIds) == 0 {
		return []string{}, nil
	}

	// 获取所有角色的权限ID
	var permissionIds []int
	for _, roleId := range roleIds {
		rolePermissions, err := Role.GetPermissions(ctx, roleId)
		if err != nil {
			continue
		}
		permissionIds = append(permissionIds, rolePermissions...)
	}

	// 去重
	permissionMap := make(map[int]bool)
	var uniquePermissionIds []int
	for _, pid := range permissionIds {
		if !permissionMap[pid] {
			permissionMap[pid] = true
			uniquePermissionIds = append(uniquePermissionIds, pid)
		}
	}

	// 获取权限代码
	if len(uniquePermissionIds) > 0 {
		var permissions []*entity.Permission
		err = dao.Permission.Ctx(ctx).WhereIn(dao.Permission.Columns().Id, uniquePermissionIds).Scan(&permissions)
		if err != nil {
			return
		}
		for _, p := range permissions {
			permissionCodes = append(permissionCodes, p.Code)
		}
	}

	return
}
