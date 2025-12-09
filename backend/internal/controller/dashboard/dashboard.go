package dashboard

import (
	"context"

	"blog/internal/dao"
	"blog/internal/model/entity"
	"blog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

// 获取统计数据请求
type GetStatsReq struct {
	g.Meta `path:"/dashboard/stats" method:"post" tags:"Dashboard" summary:"Get dashboard statistics"`
}
type GetStatsRes struct {
	TotalViews    int64 `json:"totalViews"`    // 总访问量（文章阅读数总和）
	TotalArticles int   `json:"totalArticles"` // 文章总数
	TotalComments int   `json:"totalComments"` // 评论总数
	TotalUsers    int   `json:"totalUsers"`   // 注册用户数
}

// 获取访问趋势请求
type GetTrendReq struct {
	g.Meta `path:"/dashboard/trend" method:"post" tags:"Dashboard" summary:"Get visit trend"`
	Period string `json:"period" v:"in:week,month#时间段必须是week或month" d:"week"` // week:本周, month:本月
}
type TrendData struct {
	Date      string `json:"date"`
	Views     int64  `json:"views"`     // 访问量
	Reads     int64  `json:"reads"`     // 阅读量
}
type GetTrendRes struct {
	Data []TrendData `json:"data"`
}

// 获取分类占比请求
type GetCategoryRatioReq struct {
	g.Meta `path:"/dashboard/categoryRatio" method:"post" tags:"Dashboard" summary:"Get category ratio"`
}
type CategoryRatioData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}
type GetCategoryRatioRes struct {
	Data []CategoryRatioData `json:"data"`
}

// 获取最新评论请求
type GetLatestCommentsReq struct {
	g.Meta `path:"/dashboard/latestComments" method:"post" tags:"Dashboard" summary:"Get latest comments"`
	Limit  int `json:"limit" v:"min:1|max:20#数量必须在1-20之间" d:"3"`
}
type GetLatestCommentsRes struct {
	List []*LatestCommentItem `json:"list"`
}
type LatestCommentItem struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	UserName  string `json:"userName"`
	UserAvatar string `json:"userAvatar"`
	ArticleTitle string `json:"articleTitle"`
	Time      string `json:"time"`
}

// GetStats 获取统计数据
func (c *ControllerV1) GetStats(ctx context.Context, req *GetStatsReq) (res *GetStatsRes, err error) {
	// 总访问量（文章阅读数总和）
	var totalViews int64
	value, err := dao.Article.Ctx(ctx).
		Fields("SUM(views)").
		Where(dao.Article.Columns().Status, "published").
		Value()
	if err != nil {
		totalViews = 0
	} else {
		totalViews = value.Int64()
	}

	// 文章总数（已发布）
	totalArticles, err := dao.Article.Ctx(ctx).
		Where(dao.Article.Columns().Status, "published").
		Count()
	if err != nil {
		totalArticles = 0
	}

	// 评论总数（已审核）
	totalComments, err := dao.Comment.Ctx(ctx).
		Where(dao.Comment.Columns().Status, "approved").
		Count()
	if err != nil {
		totalComments = 0
	}

	// 注册用户数
	totalUsers, err := dao.User.Ctx(ctx).Count()
	if err != nil {
		totalUsers = 0
	}

	return &GetStatsRes{
		TotalViews:    totalViews,
		TotalArticles: totalArticles,
		TotalComments: totalComments,
		TotalUsers:    totalUsers,
	}, nil
}

// GetTrend 获取访问趋势
func (c *ControllerV1) GetTrend(ctx context.Context, req *GetTrendReq) (res *GetTrendRes, err error) {
	var startTime *gtime.Time
	var days int

	now := gtime.Now()
	if req.Period == "week" {
		// 本周：从周一开始
		weekday := now.Weekday()
		daysAgo := int(weekday) - 1
		if daysAgo < 0 {
			daysAgo = 6
		}
		startTime = now.AddDate(0, 0, -daysAgo)
		startTime = gtime.New(startTime.Format("Y-m-d 00:00:00"))
		days = 7
	} else {
		// 本月：从1号开始
		startTime = gtime.New(now.Format("Y-m-01 00:00:00"))
		days = now.Day()
	}

	// 从 visit_log 表查询访问日志数据
	var visitLogs []*entity.VisitLog
	err = dao.VisitLog.Ctx(ctx).
		WhereGTE(dao.VisitLog.Columns().Date, startTime).
		WhereLT(dao.VisitLog.Columns().Date, startTime.AddDate(0, 0, days)).
		OrderAsc(dao.VisitLog.Columns().Date).
		Scan(&visitLogs)
	if err != nil {
		visitLogs = []*entity.VisitLog{}
	}

	// 创建日期到访问量的映射
	visitLogMap := make(map[string]int64)
	for _, log := range visitLogs {
		if log.Date != nil {
			dateStr := log.Date.Format("Y-m-d")
			visitLogMap[dateStr] = int64(log.Views)
		}
	}

	// 获取所有已发布文章的views总和（用于计算阅读量）
	var totalViews int64
	totalValue, err := dao.Article.Ctx(ctx).
		Fields("SUM(views)").
		Where(dao.Article.Columns().Status, "published").
		Value()
	if err == nil {
		totalViews = totalValue.Int64()
	}

	// 获取所有已发布文章数
	totalArticles, _ := dao.Article.Ctx(ctx).
		Where(dao.Article.Columns().Status, "published").
		Count()

	data := make([]TrendData, 0, days)
	for i := 0; i < days; i++ {
		date := startTime.AddDate(0, 0, i)
		dateStr := date.Format("Y-m-d")
		nextDate := date.AddDate(0, 0, 1)

		// 访问量：从 visit_log 表获取该日期的访问量
		var views int64
		if v, ok := visitLogMap[dateStr]; ok {
			views = v
		}

		// 阅读量：该日期新增的文章数 * 平均每篇文章的阅读量
		articleCount, err := dao.Article.Ctx(ctx).
			Where(dao.Article.Columns().Status, "published").
			WhereGTE(dao.Article.Columns().CreatedAt, date).
			WhereLT(dao.Article.Columns().CreatedAt, nextDate).
			Count()
		if err != nil {
			articleCount = 0
		}

		var reads int64
		if totalArticles > 0 && totalViews > 0 {
			// 使用平均阅读量来估算
			avgViews := totalViews / int64(totalArticles)
			reads = int64(articleCount) * avgViews
		} else {
			reads = views
		}

		data = append(data, TrendData{
			Date:  dateStr,
			Views: views, // 访问量：从 visit_log 表获取的每日访问量
			Reads: reads, // 阅读量：估算值（新增文章数 * 平均阅读量）
		})
	}

	return &GetTrendRes{Data: data}, nil
}

// GetCategoryRatio 获取分类占比
func (c *ControllerV1) GetCategoryRatio(ctx context.Context, req *GetCategoryRatioReq) (res *GetCategoryRatioRes, err error) {
	// 查询每个分类的文章数
	type CategoryCount struct {
		CategoryId int
		Count      int
	}

	var categoryCounts []CategoryCount
	err = dao.Article.Ctx(ctx).
		Fields("category_id, COUNT(*) as count").
		Where(dao.Article.Columns().Status, "published").
		WhereGT(dao.Article.Columns().CategoryId, 0).
		Group("category_id").
		Scan(&categoryCounts)
	if err != nil {
		return &GetCategoryRatioRes{Data: []CategoryRatioData{}}, nil
	}

	data := make([]CategoryRatioData, 0, len(categoryCounts))
	for _, cc := range categoryCounts {
		category, err := service.Category.GetOne(ctx, cc.CategoryId)
		if err == nil && category != nil {
			data = append(data, CategoryRatioData{
				Name:  category.Name,
				Value: cc.Count,
			})
		}
	}

	return &GetCategoryRatioRes{Data: data}, nil
}

// GetLatestComments 获取最新评论
func (c *ControllerV1) GetLatestComments(ctx context.Context, req *GetLatestCommentsReq) (res *GetLatestCommentsRes, err error) {
	comments, err := service.Comment.GetLatestComments(ctx, req.Limit)
	if err != nil {
		return &GetLatestCommentsRes{List: []*LatestCommentItem{}}, nil
	}

	list := make([]*LatestCommentItem, 0, len(comments))
	for _, comment := range comments {
		item := &LatestCommentItem{
			ID:      comment.Id,
			Content: comment.Content,
			Time:    comment.CreatedAt.Format("c"), // 使用 ISO8601 格式，便于前端解析
		}

		// 用户信息
		if comment.User != nil {
			item.UserName = comment.User.Nickname
			if item.UserName == "" {
				item.UserName = comment.User.Username
			}
			item.UserAvatar = comment.User.Avatar
		} else {
			item.UserName = "匿名用户"
		}

		// 文章标题
		article, err := service.Article.GetOne(ctx, comment.ArticleId, 0)
		if err == nil && article != nil {
			item.ArticleTitle = article.Title
		}

		list = append(list, item)
	}

	return &GetLatestCommentsRes{List: list}, nil
}

