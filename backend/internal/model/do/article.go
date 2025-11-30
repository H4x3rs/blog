// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Article is the golang structure of table article for DAO operations like Where/Data.
type Article struct {
	g.Meta     `orm:"table:article, do:true"`
	Id         any         //
	Title      any         // 文章标题
	Desc       any         // 文章摘要
	Content    any         // 文章内容
	CoverImage any         // 封面图片
	CategoryId any         // 分类ID
	Status      any         // 状态:published,draft
	Views       any         // 阅读数
	CreatedBy   any         // 创建人员ID
	UpdatedBy   any         // 更新人员ID
	PublishedBy any         // 发布人员ID
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
