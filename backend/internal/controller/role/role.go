package role

import (
	"context"

	"blog/internal/model/entity"
	"blog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

type CreateReq struct {
	g.Meta `path:"/role/create" method:"post" tags:"Role" summary:"Create role"`
	Name   string `json:"name" v:"required"`
	Key    string `json:"key" v:"required"`
	Desc   string `json:"desc"`
}
type CreateRes struct {
	ID int `json:"id"`
}

type UpdateReq struct {
	g.Meta `path:"/role/update" method:"post" tags:"Role" summary:"Update role"`
	ID     int    `json:"id" v:"required"`
	Name   string `json:"name" v:"required"`
	Key    string `json:"key" v:"required"`
	Desc   string `json:"desc"`
}
type UpdateRes struct{}

type GetListReq struct {
	g.Meta `path:"/role/getList" method:"post" tags:"Role" summary:"Get role list"`
	Page   int `json:"page" v:"min:1#页码必须大于0" d:"1"`
	Size   int `json:"size" v:"min:1|max:100#每页数量必须在1-100之间" d:"10"`
}
type GetListRes struct {
	List  []*entity.Role `json:"list"`
	Total int            `json:"total"`
}

type GetOneReq struct {
	g.Meta `path:"/role/getOne" method:"post" tags:"Role" summary:"Get role by id"`
	ID     int `json:"id" v:"required"`
}
type GetOneRes struct {
	*entity.Role
}

type DeleteReq struct {
	g.Meta `path:"/role/delete" method:"post" tags:"Role" summary:"Delete role"`
	ID     int `json:"id" v:"required"`
}
type DeleteRes struct{}

func (c *ControllerV1) Create(ctx context.Context, req *CreateReq) (res *CreateRes, err error) {
	id, err := service.Role.Create(ctx, &entity.Role{
		Name:   req.Name,
		Key:    req.Key,
		Desc:   req.Desc,
		Status: 1,
	})
	if err != nil {
		return nil, err
	}
	return &CreateRes{ID: id}, nil
}

func (c *ControllerV1) Update(ctx context.Context, req *UpdateReq) (res *UpdateRes, err error) {
	err = service.Role.Update(ctx, &entity.Role{
		Id:   req.ID,
		Name: req.Name,
		Key:  req.Key,
		Desc: req.Desc,
	})
	return
}

func (c *ControllerV1) GetList(ctx context.Context, req *GetListReq) (res *GetListRes, err error) {
	list, total, err := service.Role.GetList(ctx, "", req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &GetListRes{List: list, Total: total}, nil
}

func (c *ControllerV1) GetOne(ctx context.Context, req *GetOneReq) (res *GetOneRes, err error) {
	role, err := service.Role.GetOne(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &GetOneRes{Role: role}, nil
}

func (c *ControllerV1) Delete(ctx context.Context, req *DeleteReq) (res *DeleteRes, err error) {
	err = service.Role.Delete(ctx, req.ID)
	return
}

// 设置角色权限
type SetPermissionsReq struct {
	g.Meta        `path:"/role/setPermissions" method:"post" tags:"Role" summary:"Set role permissions"`
	ID            int   `json:"id" v:"required"`
	PermissionIds []int `json:"permissionIds"`
}
type SetPermissionsRes struct{}

func (c *ControllerV1) SetPermissions(ctx context.Context, req *SetPermissionsReq) (res *SetPermissionsRes, err error) {
	err = service.Role.SetPermissions(ctx, req.ID, req.PermissionIds)
	return
}

// 获取角色权限
type GetPermissionsReq struct {
	g.Meta `path:"/role/getPermissions" method:"post" tags:"Role" summary:"Get role permissions"`
	ID     int `json:"id" v:"required"`
}
type GetPermissionsRes struct {
	PermissionIds []int `json:"permissionIds"`
}

func (c *ControllerV1) GetPermissions(ctx context.Context, req *GetPermissionsReq) (res *GetPermissionsRes, err error) {
	permissionIds, err := service.Role.GetPermissions(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &GetPermissionsRes{PermissionIds: permissionIds}, nil
}
