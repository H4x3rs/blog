// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Article is the golang structure for table article.
type Article struct {
	Id         int         `json:"id"         orm:"id"          ` //
	Title      string      `json:"title"      orm:"title"       ` // 文章标题
	Desc       string      `json:"desc"       orm:"desc"        ` // 文章摘要
	Content    string      `json:"content"    orm:"content"     ` // 文章内容
	CoverImage string      `json:"coverImage" orm:"cover_image" ` // 封面图片
	CategoryId int         `json:"categoryId" orm:"category_id" ` // 分类ID
	Status      string      `json:"status"      orm:"status"       ` // 状态:published,draft
	Views       int         `json:"views"       orm:"views"        ` // 阅读数
	CreatedBy   int         `json:"createdBy"   orm:"created_by"   ` // 创建人员ID
	UpdatedBy   int         `json:"updatedBy"   orm:"updated_by"   ` // 更新人员ID
	PublishedBy int         `json:"publishedBy" orm:"published_by" ` // 发布人员ID
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   ` //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   ` //
}
