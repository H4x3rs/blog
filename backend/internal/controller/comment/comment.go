package comment

import (
	"context"

	"blog/internal/model/entity"
	"blog/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

// 创建评论请求
type CreateReq struct {
	g.Meta    `path:"/comment/create" method:"post" tags:"Comment" summary:"Create a new comment"`
	ArticleId int    `json:"articleId" v:"required#文章ID必填"`
	ParentId  int    `json:"parentId" d:"0"` // 父评论ID，默认为0（顶级评论）
	Content   string `json:"content" v:"required|length:1,1000#评论内容必填|评论内容长度必须在1-1000之间"`
}
type CreateRes struct {
	ID int `json:"id"`
}

// 获取评论列表请求
type GetListReq struct {
	g.Meta    `path:"/comment/getList" method:"post" tags:"Comment" summary:"Get comment list"`
	ArticleId int `json:"articleId" v:"required#文章ID必填"`
	Page      int `json:"page" v:"min:1#页码必须大于0" d:"1"`
	Size      int `json:"size" v:"min:1|max:100#每页数量必须在1-100之间" d:"10"`
}
type GetListRes struct {
	List  []*entity.CommentWithUser `json:"list"`
	Total int                       `json:"total"`
}

// 更新评论请求
type UpdateReq struct {
	g.Meta  `path:"/comment/update" method:"post" tags:"Comment" summary:"Update comment"`
	ID      int    `json:"id" v:"required"`
	Content string `json:"content" v:"length:1,1000#评论内容长度必须在1-1000之间"`
	Status  string `json:"status" v:"in:approved,pending,rejected#状态值不正确"`
}
type UpdateRes struct{}

// 删除评论请求
type DeleteReq struct {
	g.Meta `path:"/comment/delete" method:"post" tags:"Comment" summary:"Delete comment"`
	ID     int `json:"id" v:"required"`
}
type DeleteRes struct{}

// Create 创建评论
func (c *ControllerV1) Create(ctx context.Context, req *CreateReq) (res *CreateRes, err error) {
	// 从请求上下文中获取用户ID（通过中间件设置，可能为0表示未登录）
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id", 0).Int()

	// 检查用户是否登录，匿名用户不允许评论
	if userID == 0 {
		return nil, gerror.New("请先登录后再发表评论")
	}

	comment := &entity.Comment{
		ArticleId: req.ArticleId,
		ParentId: req.ParentId,
		Content:  req.Content,
		Status:   "approved", // 默认已审核，可以根据需要改为pending
		UserId:   &userID,     // 必须设置用户ID
	}

	id, err := service.Comment.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return &CreateRes{ID: id}, nil
}

// GetList 获取评论列表
func (c *ControllerV1) GetList(ctx context.Context, req *GetListReq) (res *GetListRes, err error) {
	list, total, err := service.Comment.GetList(ctx, req.ArticleId, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	return &GetListRes{List: list, Total: total}, nil
}

// Update 更新评论
func (c *ControllerV1) Update(ctx context.Context, req *UpdateReq) (res *UpdateRes, err error) {
	comment := &entity.Comment{
		Content: req.Content,
		Status:  req.Status,
	}

	err = service.Comment.Update(ctx, req.ID, comment)
	return
}

// Delete 删除评论
func (c *ControllerV1) Delete(ctx context.Context, req *DeleteReq) (res *DeleteRes, err error) {
	err = service.Comment.Delete(ctx, req.ID)
	return
}

// 获取管理评论列表请求（支持筛选）
type GetManageListReq struct {
	g.Meta    `path:"/comment/getManageList" method:"post" tags:"Comment" summary:"Get comment list for management"`
	Status    string `json:"status"`     // 状态筛选：approved, pending, rejected, 空字符串表示全部
	ArticleId int    `json:"articleId" d:"0"` // 文章ID筛选，0表示全部
	Page      int    `json:"page" v:"min:1#页码必须大于0" d:"1"`
	Size      int    `json:"size" v:"min:1|max:100#每页数量必须在1-100之间" d:"20"`
}
type ManageCommentItem struct {
	ID          int    `json:"id"`
	ArticleId   int    `json:"articleId"`
	ArticleTitle string `json:"articleTitle"`
	ParentId    int    `json:"parentId"`
	Content     string `json:"content"`
	Status      string `json:"status"`
	UserName    string `json:"userName"`
	UserAvatar  string `json:"userAvatar"`
	CreatedAt   string `json:"createdAt"`
}
type GetManageListRes struct {
	List  []*ManageCommentItem `json:"list"`
	Total int                  `json:"total"`
}

// GetManageList 获取管理评论列表
func (c *ControllerV1) GetManageList(ctx context.Context, req *GetManageListReq) (res *GetManageListRes, err error) {
	list, total, err := service.Comment.GetManageList(ctx, req.Status, req.ArticleId, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	// 转换为管理列表格式
	manageList := make([]*ManageCommentItem, 0, len(list))
	for _, item := range list {
		manageItem := &ManageCommentItem{
			ID:        item.Id,
			ArticleId: item.ArticleId,
			ParentId:  item.ParentId,
			Content:   item.Content,
			Status:    item.Status,
			CreatedAt: item.CreatedAt.Format("c"),
		}

		// 用户信息
		if item.User != nil {
			manageItem.UserName = item.User.Nickname
			if manageItem.UserName == "" {
				manageItem.UserName = item.User.Username
			}
			manageItem.UserAvatar = item.User.Avatar
		}

		// 文章标题
		article, err := service.Article.GetOne(ctx, item.ArticleId, 0)
		if err == nil && article != nil {
			manageItem.ArticleTitle = article.Title
		}

		manageList = append(manageList, manageItem)
	}

	return &GetManageListRes{List: manageList, Total: total}, nil
}

