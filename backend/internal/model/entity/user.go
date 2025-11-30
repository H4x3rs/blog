// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure for table user.
type User struct {
	Id        int         `json:"id"        orm:"id"         ` //
	Username  string      `json:"username"  orm:"username"   ` // 用户名
	Password  string      `json:"password"  orm:"password"   ` // 密码
	Nickname  string      `json:"nickname"  orm:"nickname"   ` // 昵称
	Email     string      `json:"email"     orm:"email"      ` // 邮箱
	Avatar    string      `json:"avatar"    orm:"avatar"     ` // 头像
	Status    int         `json:"status"    orm:"status"     ` // 状态:0禁用,1启用
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` //
}
