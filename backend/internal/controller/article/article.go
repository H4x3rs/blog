package article

import (
	"context"
	"errors"

	"blog/internal/model/entity"
	"blog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

// RPC style POST interfaces
type CreateReq struct {
	g.Meta     `path:"/article/create" method:"post" tags:"Article" summary:"Create a new article"`
	Title      string `json:"title" v:"required"`
	Content    string `json:"content" v:"required"`
	Desc       string `json:"desc"`
	CoverImage string `json:"coverImage"`
	CategoryId int    `json:"categoryId"`
	Status     string `json:"status" d:"draft"`
	TagIds     []int  `json:"tagIds"` // 标签ID列表
}
type CreateRes struct {
	ID int `json:"id"`
}

type GetOneReq struct {
	g.Meta `path:"/article/getOne" method:"post" tags:"Article" summary:"Get one article"`
	ID     int `json:"id" v:"required"`
}
type GetOneRes struct {
	*entity.Article
	PublishedByUser *entity.User  `json:"publishedByUser,omitempty"`
	Tags            []*entity.Tag `json:"tags,omitempty"` // 标签列表
}

type GetListReq struct {
	g.Meta     `path:"/article/getList" method:"post" tags:"Article" summary:"Get article list"`
	Title      string `json:"title" v:""`
	CategoryId int    `json:"categoryId" v:""`
	Status     string `json:"status" v:""`
	OnlyMine   bool   `json:"onlyMine" v:""` // 是否只查询当前用户的文章（我的文章页面）
	AuthorId   int    `json:"authorId" v:""` // 根据作者ID查询（作者文章列表页面）
	Page       int    `json:"page" v:"min:1#页码必须大于0" d:"1"`
	Size       int    `json:"size" v:"min:1|max:100#每页数量必须在1-100之间" d:"10"`
}
type ArticleWithUser struct {
	*entity.Article
	PublishedByUser *entity.User     `json:"publishedByUser,omitempty"`
	CategoryName    string           `json:"categoryName,omitempty"`
	Category        *entity.Category `json:"category,omitempty"`
}

type GetListRes struct {
	List  []*ArticleWithUser `json:"list"`
	Total int                `json:"total"`
}

type UpdateReq struct {
	g.Meta     `path:"/article/update" method:"post" tags:"Article" summary:"Update article"`
	ID         int    `json:"id" v:"required"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Desc       string `json:"desc"`
	CoverImage string `json:"coverImage"`
	CategoryId int    `json:"categoryId"`
	Status     string `json:"status"`
	TagIds     []int  `json:"tagIds"` // 标签ID列表
}
type UpdateRes struct{}

type DeleteReq struct {
	g.Meta `path:"/article/delete" method:"post" tags:"Article" summary:"Delete article"`
	ID     int `json:"id" v:"required"`
}
type DeleteRes struct{}

func (c *ControllerV1) Create(ctx context.Context, req *CreateReq) (res *CreateRes, err error) {
	// 从请求上下文中获取用户ID（通过中间件设置）
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id", 0).Int()

	article := &entity.Article{
		Title:      req.Title,
		Content:    req.Content,
		Desc:       req.Desc,
		CoverImage: req.CoverImage,
		CategoryId: req.CategoryId,
		Status:     req.Status,
		CreatedBy:  userID,
		UpdatedBy:  userID, // 创建时也设置更新人员
	}

	// 如果状态是已发布，同时设置发布人员
	if req.Status == "published" {
		article.PublishedBy = userID
	}

	id, err := service.Article.Create(ctx, article)
	if err != nil {
		return nil, err
	}

	// 设置文章标签
	if len(req.TagIds) > 0 {
		err = service.Article.SetArticleTags(ctx, id, req.TagIds)
		if err != nil {
			// 如果设置标签失败，记录错误但不影响文章创建
			g.Log().Errorf(ctx, "设置文章标签失败: %v", err)
		}
	}

	return &CreateRes{ID: id}, nil
}

func (c *ControllerV1) GetOne(ctx context.Context, req *GetOneReq) (res *GetOneRes, err error) {
	// 从请求上下文中获取用户ID（通过中间件设置）
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id", 0).Int()

	article, err := service.Article.GetOne(ctx, req.ID, userID)
	if err != nil {
		return nil, err
	}

	// 如果文章不存在（可能是权限问题），返回错误
	if article == nil {
		return nil, errors.New("文章不存在或无权访问")
	}

	res = &GetOneRes{Article: article}

	// 如果有关联发布人员，获取用户信息
	if article.PublishedBy > 0 {
		user, err := service.User.GetOne(ctx, article.PublishedBy)
		if err == nil && user != nil {
			// 不返回密码
			user.Password = ""
			res.PublishedByUser = user
		}
	}

	// 获取文章标签
	tags, err := service.Article.GetArticleTags(ctx, req.ID)
	if err == nil {
		res.Tags = tags
	}

	return res, nil
}

func (c *ControllerV1) GetList(ctx context.Context, req *GetListReq) (res *GetListRes, err error) {
	// 从请求上下文中获取用户ID（通过中间件设置）
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id", 0).Int()

	// 确保 onlyMine 参数正确传递
	onlyMine := req.OnlyMine

	list, total, err := service.Article.GetList(ctx, req.Title, req.CategoryId, req.Status, req.Page, req.Size, userID, onlyMine, req.AuthorId)
	if err != nil {
		return nil, err
	}

	// 关联发布人员信息和分类信息
	articleList := make([]*ArticleWithUser, 0, len(list))
	for _, article := range list {
		articleWithUser := &ArticleWithUser{Article: article}

		// 如果有关联发布人员，获取用户信息
		if article.PublishedBy > 0 {
			user, err := service.User.GetOne(ctx, article.PublishedBy)
			if err == nil && user != nil {
				// 不返回密码
				user.Password = ""
				articleWithUser.PublishedByUser = user
			}
		}

		// 如果有关联分类，获取分类信息
		if article.CategoryId > 0 {
			category, err := service.Category.GetOne(ctx, article.CategoryId)
			if err == nil && category != nil {
				articleWithUser.CategoryName = category.Name
				articleWithUser.Category = category
			}
		}

		articleList = append(articleList, articleWithUser)
	}

	return &GetListRes{List: articleList, Total: total}, nil
}

func (c *ControllerV1) Update(ctx context.Context, req *UpdateReq) (res *UpdateRes, err error) {
	// 从请求上下文中获取用户ID（通过中间件设置）
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id", 0).Int()

	// 获取原有文章信息，判断状态是否改变
	oldArticle, err := service.Article.GetOne(ctx, req.ID, userID)
	if err != nil {
		return nil, err
	}

	// 如果文章不存在（可能是权限问题），返回错误
	if oldArticle == nil {
		return nil, errors.New("文章不存在或无权访问")
	}

	article := &entity.Article{
		Title:      req.Title,
		Content:    req.Content,
		Desc:       req.Desc,
		CoverImage: req.CoverImage,
		CategoryId: req.CategoryId,
		Status:     req.Status,
		UpdatedBy:  userID,
	}

	// 如果状态从草稿变为已发布，设置发布人员
	if oldArticle.Status != "published" && req.Status == "published" {
		article.PublishedBy = userID
	}

	err = service.Article.Update(ctx, req.ID, article, userID)
	if err != nil {
		return nil, err
	}

	// 更新文章标签
	tagIds := req.TagIds
	if tagIds == nil {
		tagIds = []int{} // 如果未传递，则清空标签
	}
	err = service.Article.SetArticleTags(ctx, req.ID, tagIds)
	if err != nil {
		// 如果设置标签失败，记录错误但不影响文章更新
		g.Log().Errorf(ctx, "更新文章标签失败: %v", err)
	}

	return
}

func (c *ControllerV1) Delete(ctx context.Context, req *DeleteReq) (res *DeleteRes, err error) {
	// 从请求上下文中获取用户ID（通过中间件设置）
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id", 0).Int()

	err = service.Article.Delete(ctx, req.ID, userID)
	return
}
