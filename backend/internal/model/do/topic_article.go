// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TopicArticle is the golang structure of table topic_article for DAO operations like Where/Data.
type TopicArticle struct {
	g.Meta    `orm:"table:topic_article, do:true"`
	Id        any         //
	TopicId   any         // 专题ID
	ArticleId any         // 文章ID
	SortOrder any         // 排序
	CreatedAt *gtime.Time //
}
