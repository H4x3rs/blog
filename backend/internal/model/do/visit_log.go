// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// VisitLog is the golang structure of table visit_log for DAO operations like Where/Data.
type VisitLog struct {
	g.Meta    `orm:"table:visit_log, do:true"`
	Id        any         //
	Date      *gtime.Time // 访问日期
	Views     any         // 访问量
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}


