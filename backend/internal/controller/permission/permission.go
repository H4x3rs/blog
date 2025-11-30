package permission

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
	g.Meta   `path:"/permission/create" method:"post" tags:"Permission" summary:"Create permission"`
	Name     string `json:"name" v:"required"`
	Code     string `json:"code" v:"required"`
	Type     string `json:"type" v:"required"`
	Path     string `json:"path"`
	ParentId int    `json:"parentId" d:"0"`
}
type CreateRes struct {
	ID int `json:"id"`
}

type UpdateReq struct {
	g.Meta   `path:"/permission/update" method:"post" tags:"Permission" summary:"Update permission"`
	ID       int    `json:"id" v:"required"`
	Name     string `json:"name" v:"required"`
	Code     string `json:"code" v:"required"`
	Type     string `json:"type" v:"required"`
	Path     string `json:"path"`
	ParentId int    `json:"parentId"`
}
type UpdateRes struct{}

type GetListReq struct {
	g.Meta `path:"/permission/getList" method:"post" tags:"Permission" summary:"Get permission list"`
}
type GetListRes struct {
	List []map[string]interface{} `json:"list"`
}

type GetOneReq struct {
	g.Meta `path:"/permission/getOne" method:"post" tags:"Permission" summary:"Get permission by id"`
	ID     int `json:"id" v:"required"`
}
type GetOneRes struct {
	*entity.Permission
}

type DeleteReq struct {
	g.Meta `path:"/permission/delete" method:"post" tags:"Permission" summary:"Delete permission"`
	ID     int `json:"id" v:"required"`
}
type DeleteRes struct{}

func (c *ControllerV1) Create(ctx context.Context, req *CreateReq) (res *CreateRes, err error) {
	id, err := service.Permission.Create(ctx, &entity.Permission{
		ParentId:  req.ParentId,
		Name:      req.Name,
		Code:      req.Code,
		Type:      req.Type,
		Path:      req.Path,
		SortOrder: 0,
		Status:    1,
	})
	if err != nil {
		return nil, err
	}
	return &CreateRes{ID: id}, nil
}

func (c *ControllerV1) Update(ctx context.Context, req *UpdateReq) (res *UpdateRes, err error) {
	err = service.Permission.Update(ctx, &entity.Permission{
		Id:       req.ID,
		Name:     req.Name,
		Code:     req.Code,
		Type:     req.Type,
		Path:     req.Path,
		ParentId: req.ParentId,
	})
	return
}

func (c *ControllerV1) GetList(ctx context.Context, req *GetListReq) (res *GetListRes, err error) {
	list, err := service.Permission.GetTree(ctx)
	if err != nil {
		return nil, err
	}
	return &GetListRes{List: list}, nil
}

func (c *ControllerV1) GetOne(ctx context.Context, req *GetOneReq) (res *GetOneRes, err error) {
	permission, err := service.Permission.GetOne(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &GetOneRes{Permission: permission}, nil
}

func (c *ControllerV1) Delete(ctx context.Context, req *DeleteReq) (res *DeleteRes, err error) {
	err = service.Permission.Delete(ctx, req.ID)
	return
}

