package service

import (
	"context"

	"blog/internal/dao"
	"blog/internal/model/do"
	"blog/internal/model/entity"
)

type sComment struct{}

var Comment = sComment{}

// Create 创建评论
func (s *sComment) Create(ctx context.Context, in *entity.Comment) (id int, err error) {
	result, err := dao.Comment.Ctx(ctx).Data(in).Insert()
	if err != nil {
		return 0, err
	}
	id64, _ := result.LastInsertId()
	return int(id64), nil
}

// GetList 获取评论列表（按文章ID）
func (s *sComment) GetList(ctx context.Context, articleId int, page, size int) (list []*entity.CommentWithUser, total int, err error) {
	// 先查询顶级评论（parent_id = 0）
	m := dao.Comment.Ctx(ctx).
		Where(dao.Comment.Columns().ArticleId, articleId).
		Where(dao.Comment.Columns().ParentId, 0).
		Where(dao.Comment.Columns().Status, "approved") // 只显示已审核的评论

	total, err = m.Count()
	if err != nil {
		return
	}

	var comments []*entity.Comment
	err = m.OrderDesc(dao.Comment.Columns().Id).Page(page, size).Scan(&comments)
	if err != nil {
		return
	}

	// 转换为带用户信息的评论
	list = make([]*entity.CommentWithUser, 0, len(comments))
	for _, comment := range comments {
		commentWithUser := &entity.CommentWithUser{
			Comment:    comment,
			Replies:    []*entity.CommentWithUser{},
			ReplyCount: 0,
		}

		// 获取用户信息
		if comment.UserId != nil && *comment.UserId > 0 {
			user, err := User.GetOne(ctx, *comment.UserId)
			if err == nil && user != nil {
				user.Password = "" // 不返回密码
				commentWithUser.User = user
			}
		}

		// 获取回复列表
		replies, replyCount, err := s.getReplies(ctx, comment.Id)
		if err == nil {
			commentWithUser.Replies = replies
			commentWithUser.ReplyCount = replyCount
		}

		list = append(list, commentWithUser)
	}

	return
}

// getReplies 获取评论的回复列表
func (s *sComment) getReplies(ctx context.Context, parentId int) (replies []*entity.CommentWithUser, count int, err error) {
	var commentList []*entity.Comment
	err = dao.Comment.Ctx(ctx).
		Where(dao.Comment.Columns().ParentId, parentId).
		Where(dao.Comment.Columns().Status, "approved").
		OrderAsc(dao.Comment.Columns().Id).
		Scan(&commentList)
	if err != nil {
		return
	}

	count = len(commentList)
	replies = make([]*entity.CommentWithUser, 0, count)

	for _, comment := range commentList {
		commentWithUser := &entity.CommentWithUser{
			Comment:    comment,
			Replies:    []*entity.CommentWithUser{},
			ReplyCount: 0,
		}

		// 获取用户信息
		if comment.UserId != nil && *comment.UserId > 0 {
			user, err := User.GetOne(ctx, *comment.UserId)
			if err == nil && user != nil {
				user.Password = "" // 不返回密码
				commentWithUser.User = user
			}
		}

		replies = append(replies, commentWithUser)
	}

	return
}

// GetOne 获取单条评论
func (s *sComment) GetOne(ctx context.Context, id int) (out *entity.Comment, err error) {
	err = dao.Comment.Ctx(ctx).Where(dao.Comment.Columns().Id, id).Scan(&out)
	return
}

// Update 更新评论
func (s *sComment) Update(ctx context.Context, id int, in *entity.Comment) (err error) {
	_, err = dao.Comment.Ctx(ctx).
		Data(do.Comment{
			Content: in.Content,
			Status:  in.Status,
		}).
		OmitEmpty().
		FieldsEx(dao.Comment.Columns().Id, dao.Comment.Columns().ArticleId, dao.Comment.Columns().UserId, dao.Comment.Columns().ParentId, dao.Comment.Columns().CreatedAt).
		Where(dao.Comment.Columns().Id, id).
		Update()
	return
}

// Delete 删除评论
func (s *sComment) Delete(ctx context.Context, id int) (err error) {
	// 删除评论及其所有回复
	_, err = dao.Comment.Ctx(ctx).
		WhereOr(dao.Comment.Columns().Id, id).
		WhereOr(dao.Comment.Columns().ParentId, id).
		Delete()
	return
}

// GetCountByArticle 获取文章评论数
func (s *sComment) GetCountByArticle(ctx context.Context, articleId int) (count int, err error) {
	count, err = dao.Comment.Ctx(ctx).
		Where(dao.Comment.Columns().ArticleId, articleId).
		Where(dao.Comment.Columns().Status, "approved").
		Count()
	return
}

// GetLatestComments 获取最新评论（用于Dashboard）
func (s *sComment) GetLatestComments(ctx context.Context, limit int) (list []*entity.CommentWithUser, err error) {
	var comments []*entity.Comment
	err = dao.Comment.Ctx(ctx).
		Where(dao.Comment.Columns().Status, "approved").
		OrderDesc(dao.Comment.Columns().Id).
		Limit(limit).
		Scan(&comments)
	if err != nil {
		return
	}

	list = make([]*entity.CommentWithUser, 0, len(comments))
	for _, comment := range comments {
		commentWithUser := &entity.CommentWithUser{
			Comment: comment,
		}

		// 获取用户信息
		if comment.UserId != nil && *comment.UserId > 0 {
			user, err := User.GetOne(ctx, *comment.UserId)
			if err == nil && user != nil {
				user.Password = ""
				commentWithUser.User = user
			}
		}

		// 获取文章标题
		article, err := Article.GetOne(ctx, comment.ArticleId, 0)
		if err == nil && article != nil {
			// 可以在这里添加文章信息到commentWithUser
			_ = article
		}

		list = append(list, commentWithUser)
	}

	return
}

// GetManageList 获取管理评论列表（支持筛选）
func (s *sComment) GetManageList(ctx context.Context, status string, articleId int, page, size int) (list []*entity.CommentWithUser, total int, err error) {
	m := dao.Comment.Ctx(ctx)

	// 状态筛选
	if status != "" {
		m = m.Where(dao.Comment.Columns().Status, status)
	}

	// 文章ID筛选
	if articleId > 0 {
		m = m.Where(dao.Comment.Columns().ArticleId, articleId)
	}

	total, err = m.Count()
	if err != nil {
		return
	}

	var comments []*entity.Comment
	err = m.OrderDesc(dao.Comment.Columns().Id).Page(page, size).Scan(&comments)
	if err != nil {
		return
	}

	// 转换为带用户信息的评论
	list = make([]*entity.CommentWithUser, 0, len(comments))
	for _, comment := range comments {
		commentWithUser := &entity.CommentWithUser{
			Comment: comment,
		}

		// 获取用户信息
		if comment.UserId != nil && *comment.UserId > 0 {
			user, err := User.GetOne(ctx, *comment.UserId)
			if err == nil && user != nil {
				user.Password = ""
				commentWithUser.User = user
			}
		}

		list = append(list, commentWithUser)
	}

	return
}
