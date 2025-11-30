// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RolePermission is the golang structure for table role_permission.
type RolePermission struct {
	Id           int         `json:"id"           orm:"id"            ` //
	RoleId       int         `json:"roleId"       orm:"role_id"       ` // 角色ID
	PermissionId int         `json:"permissionId" orm:"permission_id" ` // 权限ID
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` //
}
