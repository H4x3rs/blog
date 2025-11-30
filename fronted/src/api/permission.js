import request from '../utils/request'

// 获取权限列表
export function getPermissionList(data = {}) {
  return request({
    url: '/permission/getList',
    method: 'post',
    data
  })
}

// 获取单个权限
export function getPermission(data) {
  return request({
    url: '/permission/getOne',
    method: 'post',
    data
  })
}

// 创建权限
export function createPermission(data) {
  return request({
    url: '/permission/create',
    method: 'post',
    data
  })
}

// 更新权限
export function updatePermission(data) {
  return request({
    url: '/permission/update',
    method: 'post',
    data
  })
}

// 删除权限
export function deletePermission(data) {
  return request({
    url: '/permission/delete',
    method: 'post',
    data
  })
}
