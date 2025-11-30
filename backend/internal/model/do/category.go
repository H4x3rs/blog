// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Category is the golang structure of table category for DAO operations like Where/Data.
type Category struct {
	g.Meta      `orm:"table:category, do:true"`
	Id          any         //
	Name        any         // 分类名称
	Slug        any         // 别名
	Icon        any         // 图标名称（Element Plus图标组件名）
	Description any         // 分类描述
	Count       any         // 文章数
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
