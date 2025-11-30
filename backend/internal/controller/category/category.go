package category

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
	g.Meta      `path:"/category/create" method:"post" tags:"Category" summary:"Create category"`
	Name        string `json:"name" v:"required"`
	Slug        string `json:"slug" v:"required"`
	Icon        string `json:"icon" v:""`
	Description string `json:"description" v:""`
}
type CreateRes struct {
	ID int `json:"id"`
}

type UpdateReq struct {
	g.Meta      `path:"/category/update" method:"post" tags:"Category" summary:"Update category"`
	ID          int    `json:"id" v:"required"`
	Name        string `json:"name" v:"required"`
	Slug        string `json:"slug" v:"required"`
	Icon        string `json:"icon" v:""`
	Description string `json:"description" v:""`
}
type UpdateRes struct{}

type GetListReq struct {
	g.Meta `path:"/category/getList" method:"post" tags:"Category" summary:"Get category list"`
	Name   string `json:"name" v:""`
	Slug   string `json:"slug" v:""`
	Page   int    `json:"page" v:"min:1#页码必须大于0" d:"1"`
	Size   int    `json:"size" v:"min:1|max:100#每页数量必须在1-100之间" d:"10"`
}
type GetListRes struct {
	List  []*entity.Category `json:"list"`
	Total int                `json:"total"`
}

type GetOneReq struct {
	g.Meta `path:"/category/getOne" method:"post" tags:"Category" summary:"Get category by id"`
	ID     int `json:"id" v:"required"`
}
type GetOneRes struct {
	*entity.Category
}

type GetBySlugReq struct {
	g.Meta `path:"/category/getBySlug" method:"post" tags:"Category" summary:"Get category by slug"`
	Slug   string `json:"slug" v:"required"`
}
type GetBySlugRes struct {
	*entity.Category
}

type DeleteReq struct {
	g.Meta `path:"/category/delete" method:"post" tags:"Category" summary:"Delete category"`
	ID     int `json:"id" v:"required"`
}
type DeleteRes struct{}

func (c *ControllerV1) Create(ctx context.Context, req *CreateReq) (res *CreateRes, err error) {
	id, err := service.Category.Create(ctx, &entity.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Icon:        req.Icon,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &CreateRes{ID: id}, nil
}

func (c *ControllerV1) Update(ctx context.Context, req *UpdateReq) (res *UpdateRes, err error) {
	err = service.Category.Update(ctx, &entity.Category{
		Id:          req.ID,
		Name:        req.Name,
		Slug:        req.Slug,
		Icon:        req.Icon,
		Description: req.Description,
	})
	return
}

func (c *ControllerV1) GetList(ctx context.Context, req *GetListReq) (res *GetListRes, err error) {
	list, total, err := service.Category.GetList(ctx, req.Name, req.Slug, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &GetListRes{List: list, Total: total}, nil
}

func (c *ControllerV1) GetOne(ctx context.Context, req *GetOneReq) (res *GetOneRes, err error) {
	category, err := service.Category.GetOne(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &GetOneRes{Category: category}, nil
}

func (c *ControllerV1) GetBySlug(ctx context.Context, req *GetBySlugReq) (res *GetBySlugRes, err error) {
	category, err := service.Category.GetBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	return &GetBySlugRes{Category: category}, nil
}

func (c *ControllerV1) Delete(ctx context.Context, req *DeleteReq) (res *DeleteRes, err error) {
	err = service.Category.Delete(ctx, req.ID)
	return
}
