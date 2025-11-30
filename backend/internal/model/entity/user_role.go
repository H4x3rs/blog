// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserRole is the golang structure for table user_role.
type UserRole struct {
	Id        int         `json:"id"        orm:"id"         ` //
	UserId    int         `json:"userId"    orm:"user_id"    ` // 用户ID
	RoleId    int         `json:"roleId"    orm:"role_id"    ` // 角色ID
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
}
