package topic

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
	g.Meta      `path:"/topic/create" method:"post" tags:"Topic" summary:"Create topic"`
	Name        string `json:"name" v:"required"`
	Slug        string `json:"slug" v:"required"`
	Description string `json:"description"`
	CoverImage  string `json:"coverImage"`
	Intro       string `json:"intro"`
	IsFeatured  int    `json:"isFeatured" d:"0"`
}
type CreateRes struct {
	ID int `json:"id"`
}

type UpdateReq struct {
	g.Meta      `path:"/topic/update" method:"post" tags:"Topic" summary:"Update topic"`
	ID          int    `json:"id" v:"required"`
	Name        string `json:"name" v:"required"`
	Slug        string `json:"slug" v:"required"`
	Description string `json:"description"`
	CoverImage  string `json:"coverImage"`
	Intro       string `json:"intro"`
	IsFeatured  int    `json:"isFeatured"`
}
type UpdateRes struct{}

type GetListReq struct {
	g.Meta `path:"/topic/getList" method:"post" tags:"Topic" summary:"Get topic list"`
	Name   string `json:"name" v:""`
	Page   int    `json:"page" v:"min:1#页码必须大于0" d:"1"`
	Size   int    `json:"size" v:"min:1|max:100#每页数量必须在1-100之间" d:"10"`
}
type GetListRes struct {
	List  []*entity.Topic `json:"list"`
	Total int             `json:"total"`
}

type GetOneReq struct {
	g.Meta `path:"/topic/getOne" method:"post" tags:"Topic" summary:"Get topic by id"`
	ID     int `json:"id" v:"required"`
}
type GetOneRes struct {
	*entity.Topic
}

type DeleteReq struct {
	g.Meta `path:"/topic/delete" method:"post" tags:"Topic" summary:"Delete topic"`
	ID     int `json:"id" v:"required"`
}
type DeleteRes struct{}

func (c *ControllerV1) Create(ctx context.Context, req *CreateReq) (res *CreateRes, err error) {
	id, err := service.Topic.Create(ctx, &entity.Topic{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Intro:       req.Intro,
		IsFeatured:  req.IsFeatured,
	})
	if err != nil {
		return nil, err
	}
	return &CreateRes{ID: id}, nil
}

func (c *ControllerV1) Update(ctx context.Context, req *UpdateReq) (res *UpdateRes, err error) {
	err = service.Topic.Update(ctx, &entity.Topic{
		Id:          req.ID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Intro:       req.Intro,
		IsFeatured:  req.IsFeatured,
	})
	return
}

func (c *ControllerV1) GetList(ctx context.Context, req *GetListReq) (res *GetListRes, err error) {
	list, total, err := service.Topic.GetList(ctx, req.Name, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &GetListRes{List: list, Total: total}, nil
}

func (c *ControllerV1) GetOne(ctx context.Context, req *GetOneReq) (res *GetOneRes, err error) {
	topic, err := service.Topic.GetOne(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &GetOneRes{Topic: topic}, nil
}

func (c *ControllerV1) Delete(ctx context.Context, req *DeleteReq) (res *DeleteRes, err error) {
	err = service.Topic.Delete(ctx, req.ID)
	return
}

// 专题文章管理
type GetTopicArticlesReq struct {
	g.Meta `path:"/topic/getTopicArticles" method:"post" tags:"Topic" summary:"Get topic articles"`
	ID     int `json:"id" v:"required"`
	Page   int `json:"page" v:"min:1#页码必须大于0" d:"1"`
	Size   int `json:"size" v:"min:1|max:100#每页数量必须在1-100之间" d:"10"`
}
type GetTopicArticlesRes struct {
	List  []map[string]interface{} `json:"list"`
	Total int                      `json:"total"`
}

type AddArticleToTopicReq struct {
	g.Meta    `path:"/topic/addArticleToTopic" method:"post" tags:"Topic" summary:"Add article to topic"`
	ID        int `json:"id" v:"required"`
	ArticleId int `json:"articleId" v:"required"`
	SortOrder int `json:"sortOrder" d:"0"`
}
type AddArticleToTopicRes struct{}

type RemoveArticleFromTopicReq struct {
	g.Meta    `path:"/topic/removeArticleFromTopic" method:"post" tags:"Topic" summary:"Remove article from topic"`
	ID        int `json:"id" v:"required"`
	ArticleId int `json:"articleId" v:"required"`
}
type RemoveArticleFromTopicRes struct{}

type UpdateTopicArticleSortReq struct {
	g.Meta    `path:"/topic/updateTopicArticleSort" method:"post" tags:"Topic" summary:"Update topic article sort"`
	ID        int `json:"id" v:"required"`
	SortOrder int `json:"sortOrder" v:"required"`
}
type UpdateTopicArticleSortRes struct{}

func (c *ControllerV1) GetTopicArticles(ctx context.Context, req *GetTopicArticlesReq) (res *GetTopicArticlesRes, err error) {
	list, total, err := service.Topic.GetTopicArticles(ctx, req.ID, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &GetTopicArticlesRes{List: list, Total: total}, nil
}

func (c *ControllerV1) AddArticleToTopic(ctx context.Context, req *AddArticleToTopicReq) (res *AddArticleToTopicRes, err error) {
	err = service.Topic.AddArticleToTopic(ctx, req.ID, req.ArticleId, req.SortOrder)
	return
}

func (c *ControllerV1) RemoveArticleFromTopic(ctx context.Context, req *RemoveArticleFromTopicReq) (res *RemoveArticleFromTopicRes, err error) {
	err = service.Topic.RemoveArticleFromTopic(ctx, req.ID, req.ArticleId)
	return
}

func (c *ControllerV1) UpdateTopicArticleSort(ctx context.Context, req *UpdateTopicArticleSortReq) (res *UpdateTopicArticleSortRes, err error) {
	err = service.Topic.UpdateTopicArticleSort(ctx, req.ID, req.SortOrder)
	return
}
