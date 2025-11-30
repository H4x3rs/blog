// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Menu is the golang structure of table menu for DAO operations like Where/Data.
type Menu struct {
	g.Meta     `orm:"table:menu, do:true"`
	Id         any         //
	ParentId   any         // 父菜单ID
	Type       any         // 类型:M目录,C菜单,F按钮
	Title      any         // 菜单标题
	Icon       any         // 图标
	Path       any         // 路由路径
	Component  any         // 组件路径
	Permission any         // 权限标识
	Order      any         // 排序
	Status     any         // 状态:active激活,inactive禁用
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
}
