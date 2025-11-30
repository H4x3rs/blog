package service

import (
	"context"
	"errors"

	"blog/internal/dao"
	"blog/internal/model/do"
	"blog/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type sArticle struct{}

var Article = sArticle{}

func (s *sArticle) Create(ctx context.Context, in *entity.Article) (id int, err error) {
	result, err := dao.Article.Ctx(ctx).Data(in).Insert()
	if err != nil {
		return 0, err
	}
	id64, _ := result.LastInsertId()
	return int(id64), nil
}

func (s *sArticle) GetOne(ctx context.Context, id int, userId int) (out *entity.Article, err error) {
	m := dao.Article.Ctx(ctx).Where(dao.Article.Columns().Id, id)

	// 先查询文章，检查状态
	var tempArticle *entity.Article
	err = m.Scan(&tempArticle)
	if err != nil {
		return nil, err
	}
	if tempArticle == nil {
		return nil, nil
	}

	// 如果是已发布的文章，所有用户都可以查看（前台显示）
	if tempArticle.Status == "published" {
		return tempArticle, nil
	}

	// 如果是草稿或其他状态，非admin只能查看自己创建的文章（后台管理）
	if userId > 0 {
		isAdmin, err := User.IsAdmin(ctx, userId)
		if err == nil && !isAdmin {
			if tempArticle.CreatedBy != userId {
				return nil, nil // 返回nil表示无权访问
			}
		}
	} else {
		// 未登录用户不能查看草稿
		return nil, nil
	}

	return tempArticle, nil
}

func (s *sArticle) GetList(ctx context.Context, title string, categoryId int, status string, page, size int, userId int, onlyMine bool, authorId int) (list []*entity.Article, total int, err error) {
	m := dao.Article.Ctx(ctx)

	if title != "" {
		m = m.WhereLike(dao.Article.Columns().Title, "%"+title+"%")
	}
	if categoryId > 0 {
		m = m.Where(dao.Article.Columns().CategoryId, categoryId)
	}
	if status != "" {
		m = m.Where(dao.Article.Columns().Status, status)
	}

	// 如果指定了authorId，根据作者ID查询（作者文章列表页面）
	if authorId > 0 {
		// 查询该作者发布的所有已发布文章
		m = m.Where(dao.Article.Columns().PublishedBy, authorId)
		// 如果前端没有指定status，强制只显示已发布的文章
		if status == "" {
			m = m.Where(dao.Article.Columns().Status, "published")
		}
		// 注意：这里不return，继续执行后面的Count和Scan
	} else if onlyMine {
		// 如果onlyMine为true，只查询当前用户创建的文章（我的文章页面）
		// 无论用户是否是admin，都只返回当前用户的文章
		if userId > 0 {
			// 只查询当前用户创建的文章
			m = m.Where(dao.Article.Columns().CreatedBy, userId)
		} else {
			// 未登录用户，onlyMine为true时不应该有数据
			m = m.Where("1 = 0") // 返回空结果
		}
		// 注意：这里不return，继续执行后面的Count和Scan
	} else {
		// 如果查询的是已发布的文章，显示所有用户的文章（前台显示，与登录用户无关）
		// 如果查询的是草稿或其他状态，非admin只能查看自己创建的文章（后台管理）
		// 如果status为空字符串，表示查询所有状态，此时非admin用户只能查看自己创建的文章（我的文章页面）
		if status == "" {
			// status为空时，查询所有状态的文章
			if userId > 0 {
				// 已登录用户，非admin只能查看自己创建的文章
				isAdmin, err := User.IsAdmin(ctx, userId)
				if err == nil && !isAdmin {
					m = m.Where(dao.Article.Columns().CreatedBy, userId)
				}
			} else {
				// 未登录用户不能查看非已发布状态的文章，强制只显示已发布的
				m = m.Where(dao.Article.Columns().Status, "published")
			}
		} else if status != "published" {
			// 非已发布状态的文章，需要用户权限检查
			if userId > 0 {
				// 已登录用户，非admin只能查看自己创建的文章
				isAdmin, err := User.IsAdmin(ctx, userId)
				if err == nil && !isAdmin {
					m = m.Where(dao.Article.Columns().CreatedBy, userId)
				}
			} else {
				// 未登录用户不能查看非已发布状态的文章，强制只显示已发布的
				m = m.Where(dao.Article.Columns().Status, "published")
			}
		}
		// status == "published" 且 onlyMine == false 时，不进行用户过滤，显示所有已发布的文章
	}

	total, err = m.Count()
	if err != nil {
		return
	}

	err = m.OrderDesc(dao.Article.Columns().Id).Page(page, size).Scan(&list)
	return
}

func (s *sArticle) Update(ctx context.Context, id int, in *entity.Article, userId int) (err error) {
	// 如果不是admin，只能更新自己创建的文章
	if userId > 0 {
		isAdmin, err := User.IsAdmin(ctx, userId)
		if err == nil && !isAdmin {
			// 检查文章是否是当前用户创建的
			article, err := dao.Article.Ctx(ctx).Where(dao.Article.Columns().Id, id).One()
			if err != nil {
				return err
			}
			if article.IsEmpty() {
				return errors.New("文章不存在")
			}
			createdBy := article.Map()["created_by"]
			if createdBy == nil || g.NewVar(createdBy).Int() != userId {
				return errors.New("无权编辑此文章")
			}
		}
	}

	// 使用 do.Article 进行更新，OmitEmpty 排除零值字段，FieldsEx 排除不应更新的字段
	_, err = dao.Article.Ctx(ctx).
		Data(do.Article{
			Title:       in.Title,
			Content:     in.Content,
			Desc:        in.Desc,
			CoverImage:  in.CoverImage,
			CategoryId:  in.CategoryId,
			Status:      in.Status,
			UpdatedBy:   in.UpdatedBy,
			PublishedBy: in.PublishedBy,
		}).
		OmitEmpty().
		FieldsEx(dao.Article.Columns().Id, dao.Article.Columns().Views, dao.Article.Columns().CreatedAt, dao.Article.Columns().CreatedBy).
		Where(dao.Article.Columns().Id, id).
		Update()
	return
}

func (s *sArticle) Delete(ctx context.Context, id int, userId int) (err error) {
	// 如果不是admin，只能删除自己创建的文章
	if userId > 0 {
		isAdmin, err := User.IsAdmin(ctx, userId)
		if err == nil && !isAdmin {
			// 检查文章是否是当前用户创建的
			article, err := dao.Article.Ctx(ctx).Where(dao.Article.Columns().Id, id).One()
			if err != nil {
				return err
			}
			if article.IsEmpty() {
				return errors.New("文章不存在")
			}
			createdBy := article.Map()["created_by"]
			if createdBy == nil || g.NewVar(createdBy).Int() != userId {
				return errors.New("无权删除此文章")
			}
		}
	}

	_, err = dao.Article.Ctx(ctx).Where(dao.Article.Columns().Id, id).Delete()
	return
}
