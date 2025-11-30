// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Permission is the golang structure of table permission for DAO operations like Where/Data.
type Permission struct {
	g.Meta    `orm:"table:permission, do:true"`
	Id        any         //
	ParentId  any         // 父权限ID
	Name      any         // 权限名称
	Code      any         // 权限标识
	Type      any         // 类型:menu菜单,btn按钮
	Path      any         // 路径
	Icon      any         // 图标
	SortOrder any         // 排序
	Status    any         // 状态:0禁用,1启用
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
