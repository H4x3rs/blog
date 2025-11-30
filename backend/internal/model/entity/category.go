// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Category is the golang structure for table category.
type Category struct {
	Id          int         `json:"id"          orm:"id"           ` //
	Name        string      `json:"name"        orm:"name"         ` // 分类名称
	Slug        string      `json:"slug"        orm:"slug"         ` // 别名
	Icon        string      `json:"icon"        orm:"icon"         ` // 图标名称（Element Plus图标组件名）
	Description string      `json:"description" orm:"description" `  // 分类描述
	Count       int         `json:"count"       orm:"count"        ` // 文章数
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   ` //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   ` //
}
