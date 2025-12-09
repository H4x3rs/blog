package service

import (
	"context"
	"encoding/json"

	"blog/internal/dao"
	"blog/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type sOperationLog struct{}

var OperationLog = sOperationLog{}

// Create 创建操作日志
func (s *sOperationLog) Create(ctx context.Context, in *entity.OperationLog) (id int, err error) {
	data := g.Map{
		"user_id":        in.UserId,
		"username":       in.Username,
		"operation_type": in.OperationType,
		"module":         in.Module,
		"operation_desc": in.OperationDesc,
		"request_method": in.RequestMethod,
		"request_path":   in.RequestPath,
		"request_params": in.RequestParams,
		"ip_address":     in.IpAddress,
		"user_agent":     in.UserAgent,
		"status":         in.Status,
		"error_message":  in.ErrorMessage,
	}

	result, err := dao.OperationLog.Ctx(ctx).Data(data).Insert()
	if err != nil {
		return 0, err
	}
	id64, _ := result.LastInsertId()
	return int(id64), nil
}

// GetList 获取操作日志列表
func (s *sOperationLog) GetList(ctx context.Context, userId int, operationType, module string, page, size int) (list []*entity.OperationLog, total int, err error) {
	m := dao.OperationLog.Ctx(ctx)

	if userId > 0 {
		m = m.Where(dao.OperationLog.Columns().UserId, userId)
	}
	if operationType != "" {
		m = m.Where(dao.OperationLog.Columns().OperationType, operationType)
	}
	if module != "" {
		m = m.Where(dao.OperationLog.Columns().Module, module)
	}

	total, err = m.Count()
	if err != nil {
		return
	}

	err = m.OrderDesc(dao.OperationLog.Columns().CreatedAt).Page(page, size).Scan(&list)
	return
}

// LogOperation 记录操作日志（便捷方法）
func (s *sOperationLog) LogOperation(ctx context.Context, userId int, username, operationType, module, operationDesc, requestMethod, requestPath, ipAddress, userAgent string, status int, requestParams interface{}, errMsg string) error {
	var paramsStr string
	if requestParams != nil {
		paramsBytes, err := json.Marshal(requestParams)
		if err == nil {
			paramsStr = string(paramsBytes)
		}
	}

	log := &entity.OperationLog{
		UserId:        userId,
		Username:      username,
		OperationType: operationType,
		Module:        module,
		OperationDesc: operationDesc,
		RequestMethod: requestMethod,
		RequestPath:   requestPath,
		RequestParams: paramsStr,
		IpAddress:     ipAddress,
		UserAgent:     userAgent,
		Status:        status,
		ErrorMessage:  errMsg,
	}

	_, err := s.Create(ctx, log)
	return err
}

