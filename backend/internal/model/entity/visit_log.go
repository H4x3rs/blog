// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// VisitLog is the golang structure for table visit_log.
type VisitLog struct {
	Id        int         `json:"id"         orm:"id"          ` //
	Date      *gtime.Time `json:"date"       orm:"date"        ` // 访问日期
	Views     int         `json:"views"      orm:"views"       ` // 访问量
	CreatedAt *gtime.Time `json:"createdAt"  orm:"created_at"  ` //
	UpdatedAt *gtime.Time `json:"updatedAt"  orm:"updated_at"  ` //
}

