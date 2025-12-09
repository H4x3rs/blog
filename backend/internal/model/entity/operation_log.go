// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OperationLog is the golang structure for table operation_log.
type OperationLog struct {
	Id            int         `json:"id"             orm:"id"              ` //
	UserId        int         `json:"userId"         orm:"user_id"         ` // 操作用户ID
	Username      string      `json:"username"       orm:"username"        ` // 操作用户名
	OperationType string      `json:"operationType"  orm:"operation_type"  ` // 操作类型：login, create, update, delete等
	Module        string      `json:"module"         orm:"module"          ` // 操作模块：user, article, category等
	OperationDesc string      `json:"operationDesc"  orm:"operation_desc"  ` // 操作描述
	RequestMethod string      `json:"requestMethod"  orm:"request_method"  ` // 请求方法：GET, POST等
	RequestPath   string      `json:"requestPath"    orm:"request_path"    ` // 请求路径
	RequestParams string      `json:"requestParams"  orm:"request_params"  ` // 请求参数（JSON格式）
	IpAddress     string      `json:"ipAddress"      orm:"ip_address"      ` // IP地址
	UserAgent     string      `json:"userAgent"      orm:"user_agent"      ` // 用户代理
	Status        int         `json:"status"         orm:"status"          ` // 操作状态：0失败，1成功
	ErrorMessage  string      `json:"errorMessage"   orm:"error_message"   ` // 错误信息
	CreatedAt     *gtime.Time `json:"createdAt"      orm:"created_at"      ` // 创建时间
}

