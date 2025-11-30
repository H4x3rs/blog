// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Tag is the golang structure of table tag for DAO operations like Where/Data.
type Tag struct {
	g.Meta    `orm:"table:tag, do:true"`
	Id        any         //
	Name      any         // 标签名称
	Slug      any         // 别名
	Color     any         // 颜色
	Count     any         // 文章数
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
