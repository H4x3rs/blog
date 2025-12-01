// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ArticleTag is the golang structure of table article_tag for DAO operations like Where/Data.
type ArticleTag struct {
	g.Meta    `orm:"table:article_tag, do:true"`
	Id        any         //
	ArticleId any         // 文章ID
	TagId     any         // 标签ID
	CreatedAt *gtime.Time //
}
