// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OperationLog is the golang structure of table operation_log for DAO operations like Where/Data.
type OperationLog struct {
	g.Meta        `orm:"table:operation_log, do:true"`
	Id            any         //
	UserId        any         // 操作用户ID
	Username      any         // 操作用户名
	OperationType any         // 操作类型：login, create, update, delete等
	Module        any         // 操作模块：user, article, category等
	OperationDesc any         // 操作描述
	RequestMethod any         // 请求方法：GET, POST等
	RequestPath   any         // 请求路径
	RequestParams any         // 请求参数（JSON格式）
	IpAddress     any         // IP地址
	UserAgent     any         // 用户代理
	Status        any         // 操作状态：0失败，1成功
	ErrorMessage  any         // 错误信息
	CreatedAt     *gtime.Time // 创建时间
}


