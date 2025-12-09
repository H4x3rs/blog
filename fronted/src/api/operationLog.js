import request from '../utils/request'

// 获取操作日志列表
export function getOperationLogList(data) {
  return request({
    url: '/operationLog/getList',
    method: 'post',
    data
  })
}

