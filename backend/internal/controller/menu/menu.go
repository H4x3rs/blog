package menu

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
	g.Meta     `path:"/menu/create" method:"post" tags:"Menu" summary:"Create menu"`
	ParentId   int    `json:"parentId" d:"0"`
	Type       string `json:"type" v:"required"`
	Title      string `json:"title" v:"required"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Permission string `json:"permission"`
	Order      int    `json:"order" d:"0"`
	Status     string `json:"status" d:"active"`
}
type CreateRes struct {
	ID int `json:"id"`
}

type UpdateReq struct {
	g.Meta     `path:"/menu/update" method:"post" tags:"Menu" summary:"Update menu"`
	ID         int    `json:"id" v:"required"`
	ParentId   int    `json:"parentId"`
	Type       string `json:"type" v:"required"`
	Title      string `json:"title" v:"required"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Permission string `json:"permission"`
	Order      int    `json:"order"`
	Status     string `json:"status"`
}
type UpdateRes struct{}

type GetListReq struct {
	g.Meta `path:"/menu/getList" method:"post" tags:"Menu" summary:"Get menu list"`
}
type GetListRes struct {
	List []map[string]interface{} `json:"list"`
}

type GetOneReq struct {
	g.Meta `path:"/menu/getOne" method:"post" tags:"Menu" summary:"Get menu by id"`
	ID     int `json:"id" v:"required"`
}
type GetOneRes struct {
	*entity.Menu
}

type DeleteReq struct {
	g.Meta `path:"/menu/delete" method:"post" tags:"Menu" summary:"Delete menu"`
	ID     int `json:"id" v:"required"`
}
type DeleteRes struct{}

func (c *ControllerV1) Create(ctx context.Context, req *CreateReq) (res *CreateRes, err error) {
	id, err := service.Menu.Create(ctx, &entity.Menu{
		ParentId:   req.ParentId,
		Type:       req.Type,
		Title:      req.Title,
		Icon:       req.Icon,
		Path:       req.Path,
		Component:  req.Component,
		Permission: req.Permission,
		Order:      req.Order,
		Status:     req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &CreateRes{ID: id}, nil
}

func (c *ControllerV1) Update(ctx context.Context, req *UpdateReq) (res *UpdateRes, err error) {
	err = service.Menu.Update(ctx, &entity.Menu{
		Id:         req.ID,
		ParentId:   req.ParentId,
		Type:       req.Type,
		Title:      req.Title,
		Icon:       req.Icon,
		Path:       req.Path,
		Component:  req.Component,
		Permission: req.Permission,
		Order:      req.Order,
		Status:     req.Status,
	})
	return
}

func (c *ControllerV1) GetList(ctx context.Context, req *GetListReq) (res *GetListRes, err error) {
	list, err := service.Menu.GetTree(ctx)
	if err != nil {
		return nil, err
	}
	return &GetListRes{List: list}, nil
}

func (c *ControllerV1) GetOne(ctx context.Context, req *GetOneReq) (res *GetOneRes, err error) {
	menu, err := service.Menu.GetOne(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &GetOneRes{Menu: menu}, nil
}

func (c *ControllerV1) Delete(ctx context.Context, req *DeleteReq) (res *DeleteRes, err error) {
	err = service.Menu.Delete(ctx, req.ID)
	return
}

// 获取当前用户的菜单（根据权限过滤）
type GetCurrentUserMenusReq struct {
	g.Meta `path:"/menu/getCurrentUserMenus" method:"post" tags:"Menu" summary:"Get current user menus"`
}
type GetCurrentUserMenusRes struct {
	List []map[string]interface{} `json:"list"`
}

func (c *ControllerV1) GetCurrentUserMenus(ctx context.Context, req *GetCurrentUserMenusReq) (res *GetCurrentUserMenusRes, err error) {
	// 从请求上下文中获取用户ID（通过中间件设置）
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id", 0).Int()

	if userID == 0 {
		// 未登录，返回空菜单
		return &GetCurrentUserMenusRes{List: []map[string]interface{}{}}, nil
	}

	// 获取用户权限
	permissionCodes, err := service.User.GetPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 根据权限获取菜单
	list, err := service.Menu.GetTreeByPermissions(ctx, permissionCodes)
	if err != nil {
		return nil, err
	}
	return &GetCurrentUserMenusRes{List: list}, nil
}
