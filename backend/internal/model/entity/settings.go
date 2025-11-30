// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Settings is the golang structure for table settings.
type Settings struct {
	Id        int         `json:"id"        orm:"id"         ` //
	Key       string      `json:"key"       orm:"key"        ` // 设置键
	Value     string      `json:"value"     orm:"value"      ` // 设置值
	Type      string      `json:"type"      orm:"type"       ` // 类型:string,number,boolean,json
	Desc      string      `json:"desc"      orm:"desc"       ` // 描述
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` //
}
