// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TopicArticle is the golang structure for table topic_article.
type TopicArticle struct {
	Id        int         `json:"id"        orm:"id"         ` //
	TopicId   int         `json:"topicId"   orm:"topic_id"   ` // 专题ID
	ArticleId int         `json:"articleId" orm:"article_id" ` // 文章ID
	SortOrder int         `json:"sortOrder" orm:"sort_order" ` // 排序
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
}
