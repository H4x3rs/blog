package service

import (
	"context"
	"errors"

	"blog/internal/dao"
	"blog/internal/model/entity"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type sTopic struct{}

var Topic = sTopic{}

func (s *sTopic) Create(ctx context.Context, in *entity.Topic) (id int, err error) {
	result, err := dao.Topic.Ctx(ctx).Data(in).Insert()
	if err != nil {
		return 0, err
	}
	id64, _ := result.LastInsertId()
	return int(id64), nil
}

func (s *sTopic) Update(ctx context.Context, in *entity.Topic) (err error) {
	_, err = dao.Topic.Ctx(ctx).Data(in).Where(dao.Topic.Columns().Id, in.Id).Update()
	return
}

func (s *sTopic) GetList(ctx context.Context, name string, page, size int) (list []*entity.Topic, total int, err error) {
	m := dao.Topic.Ctx(ctx)

	if name != "" {
		m = m.WhereLike(dao.Topic.Columns().Name, "%"+name+"%")
	}

	total, err = m.Count()
	if err != nil {
		return
	}

	err = m.OrderDesc(dao.Topic.Columns().IsFeatured).OrderAsc(dao.Topic.Columns().SortOrder).OrderDesc(dao.Topic.Columns().Id).Page(page, size).Scan(&list)
	return
}

func (s *sTopic) GetOne(ctx context.Context, id int) (topic *entity.Topic, err error) {
	err = dao.Topic.Ctx(ctx).Where(dao.Topic.Columns().Id, id).Scan(&topic)
	if err != nil {
		return nil, err
	}
	if topic == nil {
		return nil, nil
	}

	// 增加阅读量（使用原子操作，避免并发问题）
	_, err = dao.Topic.Ctx(ctx).
		Where(dao.Topic.Columns().Id, id).
		Data(gdb.Map{
			dao.Topic.Columns().Views: gdb.Raw("`views` + 1"),
		}).
		Update()
	if err != nil {
		// 如果增加阅读量失败，记录错误但不影响返回专题
		g.Log().Errorf(ctx, "增加专题阅读量失败: %v", err)
	} else {
		// 更新返回的专题对象的阅读量
		topic.Views++
	}

	return topic, nil
}

func (s *sTopic) Delete(ctx context.Context, id int) (err error) {
	_, err = dao.Topic.Ctx(ctx).Where(dao.Topic.Columns().Id, id).Delete()
	return
}

// 获取专题的文章列表（包含文章详情）
func (s *sTopic) GetTopicArticles(ctx context.Context, topicId int, page, size int) (list []map[string]interface{}, total int, err error) {
	// 先获取总数
	total, err = dao.TopicArticle.Ctx(ctx).Where(dao.TopicArticle.Columns().TopicId, topicId).Count()
	if err != nil {
		return
	}

	// 分页查询专题文章关联
	var topicArticles []*entity.TopicArticle
	err = dao.TopicArticle.Ctx(ctx).
		Where(dao.TopicArticle.Columns().TopicId, topicId).
		OrderAsc(dao.TopicArticle.Columns().SortOrder).
		OrderAsc(dao.TopicArticle.Columns().Id).
		Page(page, size).
		Scan(&topicArticles)
	if err != nil {
		return
	}

	if len(topicArticles) == 0 {
		return []map[string]interface{}{}, total, nil
	}

	articleIds := make([]int, 0, len(topicArticles))
	for _, ta := range topicArticles {
		articleIds = append(articleIds, ta.ArticleId)
	}

	var articles []*entity.Article
	err = dao.Article.Ctx(ctx).WhereIn(dao.Article.Columns().Id, articleIds).Scan(&articles)
	if err != nil {
		return
	}

	articleMap := make(map[int]*entity.Article)
	for _, article := range articles {
		articleMap[article.Id] = article
	}

	list = make([]map[string]interface{}, 0, len(topicArticles))
	for _, ta := range topicArticles {
		if article, ok := articleMap[ta.ArticleId]; ok {
			list = append(list, map[string]interface{}{
				"id":         ta.Id,
				"articleId":  article.Id,
				"title":      article.Title,
				"desc":       article.Desc,
				"coverImage": article.CoverImage,
				"views":      article.Views,
				"status":     article.Status,
				"categoryId": article.CategoryId,
				"sortOrder":  ta.SortOrder,
				"createdAt":  article.CreatedAt,
			})
		}
	}

	return
}

// 添加文章到专题
func (s *sTopic) AddArticleToTopic(ctx context.Context, topicId int, articleId int, sortOrder int) (err error) {
	// 检查文章是否已在专题中 - 使用文章ID来判断是否重复
	count, err := dao.TopicArticle.Ctx(ctx).
		Where(dao.TopicArticle.Columns().TopicId, topicId).
		Where(dao.TopicArticle.Columns().ArticleId, articleId).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该文章已在专题中，无法重复添加")
	}

	// 添加文章到专题
	_, err = dao.TopicArticle.Ctx(ctx).Data(map[string]interface{}{
		"topic_id":   topicId,
		"article_id": articleId,
		"sort_order": sortOrder,
	}).Insert()
	return
}

// 从专题移除文章
func (s *sTopic) RemoveArticleFromTopic(ctx context.Context, topicId int, articleId int) (err error) {
	_, err = dao.TopicArticle.Ctx(ctx).Where(dao.TopicArticle.Columns().TopicId, topicId).Where(dao.TopicArticle.Columns().ArticleId, articleId).Delete()
	return
}

// 更新专题文章排序
func (s *sTopic) UpdateTopicArticleSort(ctx context.Context, id int, sortOrder int) (err error) {
	_, err = dao.TopicArticle.Ctx(ctx).Data(map[string]interface{}{
		"sort_order": sortOrder,
	}).Where(dao.TopicArticle.Columns().Id, id).Update()
	return
}

// 获取专题的文章数量
func (s *sTopic) GetTopicArticleCount(ctx context.Context, topicId int) (count int, err error) {
	count, err = dao.TopicArticle.Ctx(ctx).Where(dao.TopicArticle.Columns().TopicId, topicId).Count()
	return
}
