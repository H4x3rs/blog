// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ArticleTag is the golang structure for table article_tag.
type ArticleTag struct {
	Id        int         `json:"id"        orm:"id"         ` //
	ArticleId int         `json:"articleId" orm:"article_id" ` // 文章ID
	TagId     int         `json:"tagId"     orm:"tag_id"     ` // 标签ID
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
}

