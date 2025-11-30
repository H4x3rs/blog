// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Role is the golang structure for table role.
type Role struct {
	Id        int         `json:"id"        orm:"id"         ` //
	Name      string      `json:"name"      orm:"name"       ` // 角色名称
	Key       string      `json:"key"       orm:"key"        ` // 角色标识
	Desc      string      `json:"desc"      orm:"desc"       ` // 角色描述
	Status    int         `json:"status"    orm:"status"     ` // 状态:0禁用,1启用
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` //
}
