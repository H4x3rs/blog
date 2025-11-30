// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Topic is the golang structure of table topic for DAO operations like Where/Data.
type Topic struct {
	g.Meta      `orm:"table:topic, do:true"`
	Id          any         //
	Name        any         // 专题名称
	Slug        any         // 别名
	Description any         // 专题描述
	CoverImage  any         // 封面图片
	Intro       any         // 专题介绍
	IsFeatured  any         // 是否置顶:0否,1是
	SortOrder   any         // 排序
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
}
