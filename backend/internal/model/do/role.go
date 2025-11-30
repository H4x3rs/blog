// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Role is the golang structure of table role for DAO operations like Where/Data.
type Role struct {
	g.Meta    `orm:"table:role, do:true"`
	Id        any         //
	Name      any         // 角色名称
	Key       any         // 角色标识
	Desc      any         // 角色描述
	Status    any         // 状态:0禁用,1启用
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
