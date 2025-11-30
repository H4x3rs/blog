// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Menu is the golang structure for table menu.
type Menu struct {
	Id         int         `json:"id"         orm:"id"         ` //
	ParentId   int         `json:"parentId"   orm:"parent_id"  ` // 父菜单ID
	Type       string      `json:"type"       orm:"type"       ` // 类型:M目录,C菜单,F按钮
	Title      string      `json:"title"      orm:"title"      ` // 菜单标题
	Icon       string      `json:"icon"       orm:"icon"       ` // 图标
	Path       string      `json:"path"       orm:"path"       ` // 路由路径
	Component  string      `json:"component"  orm:"component"  ` // 组件路径
	Permission string      `json:"permission" orm:"permission" ` // 权限标识
	Order      int         `json:"order"      orm:"order"      ` // 排序
	Status     string      `json:"status"     orm:"status"     ` // 状态:active激活,inactive禁用
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at" ` //
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at" ` //
}
