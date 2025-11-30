// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure of table user for DAO operations like Where/Data.
type User struct {
	g.Meta    `orm:"table:user, do:true"`
	Id        any         //
	Username  any         // 用户名
	Password  any         // 密码
	Nickname  any         // 昵称
	Email     any         // 邮箱
	Avatar    any         // 头像
	Status    any         // 状态:0禁用,1启用
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
