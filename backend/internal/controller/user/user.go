package user

import (
	"context"
	"errors"

	"blog/internal/dao"
	"blog/internal/model/entity"
	"blog/internal/service"
	"blog/internal/utils"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

type LoginReq struct {
	g.Meta   `path:"/user/login" method:"post" tags:"User" summary:"User Login"`
	Username string `json:"username" v:"required"`
	Password string `json:"password" v:"required"`
}
type LoginRes struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type GetListReq struct {
	g.Meta   `path:"/user/getList" method:"post" tags:"User" summary:"Get user list"`
	Username string `json:"username" v:""`
	Page     int    `json:"page" v:"min:1#页码必须大于0" d:"1"`
	Size     int    `json:"size" v:"min:1|max:100#每页数量必须在1-100之间" d:"10"`
}
type GetListRes struct {
	List  []*entity.User `json:"list"`
	Total int            `json:"total"`
}

type GetOneReq struct {
	g.Meta `path:"/user/getOne" method:"post" tags:"User" summary:"Get user by id"`
	ID     int `json:"id" v:"required"`
}
type GetOneRes struct {
	*entity.User
}

type CreateReq struct {
	g.Meta   `path:"/user/create" method:"post" tags:"User" summary:"Create user"`
	Username string `json:"username" v:"required"`
	Password string `json:"password" v:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}
type CreateRes struct {
	ID int `json:"id"`
}

// 注册请求
type RegisterReq struct {
	g.Meta   `path:"/user/register" method:"post" tags:"User" summary:"User register"`
	Username string `json:"username" v:"required|length:3,20#用户名必填|用户名长度必须在3-20之间"`
	Password string `json:"password" v:"required|length:6,20#密码必填|密码长度必须在6-20之间"`
	Nickname string `json:"nickname" v:"length:1,50#昵称长度必须在1-50之间"`
	Email    string `json:"email" v:"required|email#邮箱必填|邮箱格式不正确"`
}
type RegisterRes struct {
	ID int `json:"id"`
}

type UpdateReq struct {
	g.Meta   `path:"/user/update" method:"post" tags:"User" summary:"Update user"`
	ID       int    `json:"id" v:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"status"`
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/user/delete" method:"post" tags:"User" summary:"Delete user"`
	ID     int `json:"id" v:"required"`
}
type DeleteRes struct{}

func (c *ControllerV1) Login(ctx context.Context, req *LoginReq) (res *LoginRes, err error) {
	user, err := service.User.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.Id, user.Username)
	if err != nil {
		return nil, err
	}

	return &LoginRes{
		Token:    token,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}

func (c *ControllerV1) GetList(ctx context.Context, req *GetListReq) (res *GetListRes, err error) {
	list, total, err := service.User.GetList(ctx, req.Username, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &GetListRes{List: list, Total: total}, nil
}

func (c *ControllerV1) GetOne(ctx context.Context, req *GetOneReq) (res *GetOneRes, err error) {
	user, err := service.User.GetOne(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &GetOneRes{User: user}, nil
}

func (c *ControllerV1) Create(ctx context.Context, req *CreateReq) (res *CreateRes, err error) {
	id, err := service.User.Create(ctx, &entity.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Avatar:   req.Avatar,
		Status:   1,
	})
	if err != nil {
		return nil, err
	}

	// 为新注册用户分配默认user角色
	err = service.User.AssignDefaultUserRole(ctx, id)
	if err != nil {
		// 记录错误但不影响用户创建
		glog.Error(ctx, "分配默认角色失败:", err)
	}

	return &CreateRes{ID: id}, nil
}

func (c *ControllerV1) Update(ctx context.Context, req *UpdateReq) (res *UpdateRes, err error) {
	err = service.User.Update(ctx, &entity.User{
		Id:       req.ID,
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Avatar:   req.Avatar,
		Status:   req.Status,
	})
	return
}

func (c *ControllerV1) Delete(ctx context.Context, req *DeleteReq) (res *DeleteRes, err error) {
	err = service.User.Delete(ctx, req.ID)
	return
}

// 个人中心相关
type UpdateProfileReq struct {
	g.Meta   `path:"/user/updateProfile" method:"post" tags:"User" summary:"Update user profile"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}
type UpdateProfileRes struct{}

type ChangePasswordReq struct {
	g.Meta      `path:"/user/changePassword" method:"post" tags:"User" summary:"Change password"`
	OldPassword string `json:"oldPassword" v:"required"`
	NewPassword string `json:"newPassword" v:"required"`
}
type ChangePasswordRes struct{}

type GetCurrentUserReq struct {
	g.Meta `path:"/user/getCurrentUser" method:"post" tags:"User" summary:"Get current user"`
}
type GetCurrentUserRes struct {
	*entity.User
}

func (c *ControllerV1) UpdateProfile(ctx context.Context, req *UpdateProfileReq) (res *UpdateProfileRes, err error) {
	// 从请求上下文中获取用户ID（通过中间件设置）
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id", 0).Int()

	if userID == 0 {
		return nil, errors.New("未登录或token无效")
	}

	err = service.User.Update(ctx, &entity.User{
		Id:       userID,
		Nickname: req.Nickname,
		Email:    req.Email,
		Avatar:   req.Avatar,
	})
	return
}

func (c *ControllerV1) ChangePassword(ctx context.Context, req *ChangePasswordReq) (res *ChangePasswordRes, err error) {
	// TODO: 从 token 中获取用户ID
	userID := 1 // Mock user ID
	// 验证旧密码
	user, err := service.User.GetOne(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.Password != req.OldPassword {
		return nil, errors.New("旧密码错误")
	}
	// 更新新密码
	err = service.User.Update(ctx, &entity.User{
		Id:       userID,
		Password: req.NewPassword,
	})
	return
}

func (c *ControllerV1) GetCurrentUser(ctx context.Context, req *GetCurrentUserReq) (res *GetCurrentUserRes, err error) {
	// 从请求上下文中获取用户ID（通过中间件设置）
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id", 0).Int()

	if userID == 0 {
		return nil, errors.New("未登录或token无效")
	}

	user, err := service.User.GetOne(ctx, userID)
	if err != nil {
		return nil, err
	}
	// 不返回密码
	user.Password = ""
	return &GetCurrentUserRes{User: user}, nil
}

// 设置用户角色
type SetRolesReq struct {
	g.Meta  `path:"/user/setRoles" method:"post" tags:"User" summary:"Set user roles"`
	ID      int   `json:"id" v:"required"`
	RoleIds []int `json:"roleIds"`
}
type SetRolesRes struct{}

func (c *ControllerV1) SetRoles(ctx context.Context, req *SetRolesReq) (res *SetRolesRes, err error) {
	err = service.User.SetRoles(ctx, req.ID, req.RoleIds)
	return
}

// 获取用户角色
type GetRolesReq struct {
	g.Meta `path:"/user/getRoles" method:"post" tags:"User" summary:"Get user roles"`
	ID     int `json:"id" v:"required"`
}
type GetRolesRes struct {
	RoleIds []int `json:"roleIds"`
}

func (c *ControllerV1) GetRoles(ctx context.Context, req *GetRolesReq) (res *GetRolesRes, err error) {
	roleIds, err := service.User.GetRoles(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &GetRolesRes{RoleIds: roleIds}, nil
}

// 获取当前用户权限
type GetCurrentUserPermissionsReq struct {
	g.Meta `path:"/user/getCurrentUserPermissions" method:"post" tags:"User" summary:"Get current user permissions"`
}
type GetCurrentUserPermissionsRes struct {
	PermissionCodes []string `json:"permissionCodes"`
}

func (c *ControllerV1) GetCurrentUserPermissions(ctx context.Context, req *GetCurrentUserPermissionsReq) (res *GetCurrentUserPermissionsRes, err error) {
	// 从请求上下文中获取用户ID（通过中间件设置）
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id", 0).Int()

	if userID == 0 {
		return nil, errors.New("未登录或token无效")
	}

	permissionCodes, err := service.User.GetPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &GetCurrentUserPermissionsRes{PermissionCodes: permissionCodes}, nil
}

// 获取当前用户角色
type GetCurrentUserRolesReq struct {
	g.Meta `path:"/user/getCurrentUserRoles" method:"post" tags:"User" summary:"Get current user roles"`
}
type GetCurrentUserRolesRes struct {
	RoleIds []int `json:"roleIds"`
}

func (c *ControllerV1) GetCurrentUserRoles(ctx context.Context, req *GetCurrentUserRolesReq) (res *GetCurrentUserRolesRes, err error) {
	// 从请求上下文中获取用户ID（通过中间件设置）
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id", 0).Int()

	if userID == 0 {
		return nil, errors.New("未登录或token无效")
	}

	roleIds, err := service.User.GetRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &GetCurrentUserRolesRes{RoleIds: roleIds}, nil
}

// 用户注册
func (c *ControllerV1) Register(ctx context.Context, req *RegisterReq) (res *RegisterRes, err error) {
	// 检查用户名是否已存在
	var existingUser *entity.User
	err = dao.User.Ctx(ctx).Where(dao.User.Columns().Username, req.Username).Scan(&existingUser)
	if err == nil && existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 如果提供了邮箱，检查邮箱是否已存在
	if req.Email != "" {
		var emailUser *entity.User
		err = dao.User.Ctx(ctx).Where(dao.User.Columns().Email, req.Email).Scan(&emailUser)
		if err == nil && emailUser != nil {
			return nil, errors.New("邮箱已被注册")
		}
	}

	// 创建用户
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	// 邮箱格式验证（邮箱已通过v标签验证为必填和格式）

	id, err := service.User.Create(ctx, &entity.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: nickname,
		Email:    req.Email,
		Status:   1,
	})
	if err != nil {
		return nil, err
	}

	// 为新注册用户分配默认user角色
	err = service.User.AssignDefaultUserRole(ctx, id)
	if err != nil {
		// 记录错误但不影响用户创建
		glog.Error(ctx, "分配默认角色失败:", err)
	}

	return &RegisterRes{ID: id}, nil
}
