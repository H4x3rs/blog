package tag

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
	g.Meta `path:"/tag/create" method:"post" tags:"Tag" summary:"Create tag"`
	Name   string `json:"name" v:"required"`
	Slug   string `json:"slug" v:"required"`
	Color  string `json:"color" d:"#409EFF"`
}
type CreateRes struct {
	ID int `json:"id"`
}

type UpdateReq struct {
	g.Meta `path:"/tag/update" method:"post" tags:"Tag" summary:"Update tag"`
	ID     int    `json:"id" v:"required"`
	Name   string `json:"name" v:"required"`
	Slug   string `json:"slug" v:"required"`
	Color  string `json:"color"`
}
type UpdateRes struct{}

type GetListReq struct {
	g.Meta `path:"/tag/getList" method:"post" tags:"Tag" summary:"Get tag list"`
	Name   string `json:"name" v:""`
	Slug   string `json:"slug" v:""`
	Page   int    `json:"page" v:"min:1#页码必须大于0" d:"1"`
	Size   int    `json:"size" v:"min:1|max:100#每页数量必须在1-100之间" d:"10"`
}
type GetListRes struct {
	List  []*entity.Tag `json:"list"`
	Total int           `json:"total"`
}

type GetOneReq struct {
	g.Meta `path:"/tag/getOne" method:"post" tags:"Tag" summary:"Get tag by id"`
	ID     int `json:"id" v:"required"`
}
type GetOneRes struct {
	*entity.Tag
}

type GetBySlugReq struct {
	g.Meta `path:"/tag/getBySlug" method:"post" tags:"Tag" summary:"Get tag by slug"`
	Slug   string `json:"slug" v:"required"`
}
type GetBySlugRes struct {
	*entity.Tag
}

type DeleteReq struct {
	g.Meta `path:"/tag/delete" method:"post" tags:"Tag" summary:"Delete tag"`
	ID     int `json:"id" v:"required"`
}
type DeleteRes struct{}

func (c *ControllerV1) Create(ctx context.Context, req *CreateReq) (res *CreateRes, err error) {
	id, err := service.Tag.Create(ctx, &entity.Tag{
		Name:  req.Name,
		Slug:  req.Slug,
		Color: req.Color,
	})
	if err != nil {
		return nil, err
	}
	return &CreateRes{ID: id}, nil
}

func (c *ControllerV1) Update(ctx context.Context, req *UpdateReq) (res *UpdateRes, err error) {
	err = service.Tag.Update(ctx, &entity.Tag{
		Id:    req.ID,
		Name:  req.Name,
		Slug:  req.Slug,
		Color: req.Color,
	})
	return
}

func (c *ControllerV1) GetList(ctx context.Context, req *GetListReq) (res *GetListRes, err error) {
	list, total, err := service.Tag.GetList(ctx, req.Name, req.Slug, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &GetListRes{List: list, Total: total}, nil
}

func (c *ControllerV1) GetOne(ctx context.Context, req *GetOneReq) (res *GetOneRes, err error) {
	tag, err := service.Tag.GetOne(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &GetOneRes{Tag: tag}, nil
}

func (c *ControllerV1) GetBySlug(ctx context.Context, req *GetBySlugReq) (res *GetBySlugRes, err error) {
	tag, err := service.Tag.GetBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	return &GetBySlugRes{Tag: tag}, nil
}

func (c *ControllerV1) Delete(ctx context.Context, req *DeleteReq) (res *DeleteRes, err error) {
	err = service.Tag.Delete(ctx, req.ID)
	return
}

