// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Tag is the golang structure for table tag.
type Tag struct {
	Id        int         `json:"id"        orm:"id"         ` //
	Name      string      `json:"name"      orm:"name"       ` // 标签名称
	Slug      string      `json:"slug"      orm:"slug"       ` // 别名
	Color     string      `json:"color"     orm:"color"      ` // 颜色
	Count     int         `json:"count"     orm:"count"      ` // 文章数
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` //
}
