// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Settings is the golang structure of table settings for DAO operations like Where/Data.
type Settings struct {
	g.Meta    `orm:"table:settings, do:true"`
	Id        any         //
	Key       any         // 设置键
	Value     any         // 设置值
	Type      any         // 类型:string,number,boolean,json
	Desc      any         // 描述
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
