import request from '../utils/request'

// 获取角色列表
export function getRoleList(data) {
  return request({
    url: '/role/getList',
    method: 'post',
    data
  })
}

// 获取单个角色
export function getRole(data) {
  return request({
    url: '/role/getOne',
    method: 'post',
    data
  })
}

// 创建角色
export function createRole(data) {
  return request({
    url: '/role/create',
    method: 'post',
    data
  })
}

// 更新角色
export function updateRole(data) {
  return request({
    url: '/role/update',
    method: 'post',
    data
  })
}

// 删除角色
export function deleteRole(data) {
  return request({
    url: '/role/delete',
    method: 'post',
    data
  })
}

// 设置角色权限
export function setRolePermissions(data) {
  return request({
    url: '/role/setPermissions',
    method: 'post',
    data
  })
}

// 获取角色权限
export function getRolePermissions(data) {
  return request({
    url: '/role/getPermissions',
    method: 'post',
    data
  })
}