// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Topic is the golang structure for table topic.
type Topic struct {
	Id          int         `json:"id"          orm:"id"          ` //
	Name        string      `json:"name"        orm:"name"        ` // 专题名称
	Slug        string      `json:"slug"        orm:"slug"        ` // 别名
	Description string      `json:"description" orm:"description" ` // 专题描述
	CoverImage  string      `json:"coverImage"  orm:"cover_image" ` // 封面图片
	Intro       string      `json:"intro"       orm:"intro"       ` // 专题介绍
	IsFeatured  int         `json:"isFeatured"  orm:"is_featured" ` // 是否置顶:0否,1是
	SortOrder   int         `json:"sortOrder"   orm:"sort_order"  ` // 排序
	Views       int         `json:"views"       orm:"views"        ` // 阅读数
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  ` //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  ` //
}
