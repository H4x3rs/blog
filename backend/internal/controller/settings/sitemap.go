package settings

import (
	"context"
	"fmt"
	"time"

	"blog/internal/dao"
	"blog/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

// 生成sitemap.xml
type SitemapReq struct {
	g.Meta `path:"/sitemap.xml" method:"get" tags:"SEO" summary:"Generate sitemap.xml"`
}
type SitemapRes struct{}

// 生成robots.txt
type RobotsReq struct {
	g.Meta `path:"/robots.txt" method:"get" tags:"SEO" summary:"Generate robots.txt"`
}
type RobotsRes struct{}

func (c *ControllerV1) Sitemap(ctx context.Context, req *SitemapReq) (res *SitemapRes, err error) {
	r := g.RequestFromCtx(ctx)

	// 获取网站URL
	scheme := "https"
	if r.Header.Get("X-Forwarded-Proto") == "http" || r.URL.Scheme == "http" {
		scheme = "http"
	}
	siteUrl := fmt.Sprintf("%s://%s", scheme, r.Host)

	// 获取所有已发布的文章
	var articles []*entity.Article
	err = dao.Article.Ctx(ctx).
		Where(dao.Article.Columns().Status, "published").
		OrderDesc(dao.Article.Columns().Id).
		Scan(&articles)
	if err != nil {
		return nil, err
	}

	// 获取所有分类
	var categories []*entity.Category
	err = dao.Category.Ctx(ctx).Scan(&categories)
	if err != nil {
		return nil, err
	}

	// 获取所有标签
	var tags []*entity.Tag
	err = dao.Tag.Ctx(ctx).Scan(&tags)
	if err != nil {
		return nil, err
	}

	// 获取所有专题
	var topics []*entity.Topic
	err = dao.Topic.Ctx(ctx).Scan(&topics)
	if err != nil {
		return nil, err
	}

	// 生成sitemap XML
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
`

	// 首页
	xml += fmt.Sprintf(`  <url>
    <loc>%s/</loc>
    <lastmod>%s</lastmod>
    <changefreq>daily</changefreq>
    <priority>1.0</priority>
  </url>
`, siteUrl, time.Now().Format("2006-01-02"))

	// 分类列表页
	xml += fmt.Sprintf(`  <url>
    <loc>%s/category</loc>
    <lastmod>%s</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
`, siteUrl, time.Now().Format("2006-01-02"))

	// 标签列表页
	xml += fmt.Sprintf(`  <url>
    <loc>%s/tag</loc>
    <lastmod>%s</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
`, siteUrl, time.Now().Format("2006-01-02"))

	// 专题列表页
	xml += fmt.Sprintf(`  <url>
    <loc>%s/topic</loc>
    <lastmod>%s</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
`, siteUrl, time.Now().Format("2006-01-02"))

	// 关于页面
	xml += fmt.Sprintf(`  <url>
    <loc>%s/about</loc>
    <lastmod>%s</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.5</priority>
  </url>
`, siteUrl, time.Now().Format("2006-01-02"))

	// 文章详情页
	for _, article := range articles {
		updatedAt := article.UpdatedAt
		if updatedAt == nil {
			updatedAt = article.CreatedAt
		}
		lastmod := updatedAt.Format("2006-01-02")
		xml += fmt.Sprintf(`  <url>
    <loc>%s/article/%d</loc>
    <lastmod>%s</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
`, siteUrl, article.Id, lastmod)
	}

	// 分类详情页
	for _, category := range categories {
		if category.Slug != "" {
			xml += fmt.Sprintf(`  <url>
    <loc>%s/category/%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.7</priority>
  </url>
`, siteUrl, category.Slug, time.Now().Format("2006-01-02"))
		}
	}

	// 标签详情页
	for _, tag := range tags {
		if tag.Slug != "" {
			xml += fmt.Sprintf(`  <url>
    <loc>%s/tag/%s</loc>
    <lastmod>%s</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.7</priority>
  </url>
`, siteUrl, tag.Slug, time.Now().Format("2006-01-02"))
		}
	}

	// 专题详情页
	for _, topic := range topics {
		xml += fmt.Sprintf(`  <url>
    <loc>%s/topic/%d</loc>
    <lastmod>%s</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.7</priority>
  </url>
`, siteUrl, topic.Id, time.Now().Format("2006-01-02"))
	}

	xml += `</urlset>`

	// 设置响应头
	r.Response.Header().Set("Content-Type", "application/xml; charset=utf-8")
	r.Response.Write(xml)

	return &SitemapRes{}, nil
}

func (c *ControllerV1) Robots(ctx context.Context, req *RobotsReq) (res *RobotsRes, err error) {
	r := g.RequestFromCtx(ctx)

	// 获取网站URL
	scheme := "https"
	if r.Header.Get("X-Forwarded-Proto") == "http" || r.URL.Scheme == "http" {
		scheme = "http"
	}
	siteUrl := fmt.Sprintf("%s://%s", scheme, r.Host)

	robots := fmt.Sprintf(`User-agent: *
Allow: /

# Sitemap
Sitemap: %s/sitemap.xml
`, siteUrl)

	r.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	r.Response.Write(robots)

	return &RobotsRes{}, nil
}
