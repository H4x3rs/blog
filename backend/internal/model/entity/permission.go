// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Permission is the golang structure for table permission.
type Permission struct {
	Id        int         `json:"id"        orm:"id"         ` //
	ParentId  int         `json:"parentId"  orm:"parent_id"  ` // 父权限ID
	Name      string      `json:"name"      orm:"name"       ` // 权限名称
	Code      string      `json:"code"      orm:"code"       ` // 权限标识
	Type      string      `json:"type"      orm:"type"       ` // 类型:menu菜单,btn按钮
	Path      string      `json:"path"      orm:"path"       ` // 路径
	Icon      string      `json:"icon"      orm:"icon"       ` // 图标
	SortOrder int         `json:"sortOrder" orm:"sort_order" ` // 排序
	Status    int         `json:"status"    orm:"status"     ` // 状态:0禁用,1启用
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` //
}
