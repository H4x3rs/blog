import request from '../utils/request'

// 用户登录
export function login(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

// 获取用户列表
export function getUserList(data) {
  return request({
    url: '/user/getList',
    method: 'post',
    data
  })
}

// 获取单个用户
export function getUser(data) {
  return request({
    url: '/user/getOne',
    method: 'post',
    data
  })
}

// 创建用户
export function createUser(data) {
  return request({
    url: '/user/create',
    method: 'post',
    data
  })
}

// 更新用户
export function updateUser(data) {
  return request({
    url: '/user/update',
    method: 'post',
    data
  })
}

// 删除用户
export function deleteUser(data) {
  return request({
    url: '/user/delete',
    method: 'post',
    data
  })
}

// 更新个人资料
export function updateProfile(data) {
  return request({
    url: '/user/updateProfile',
    method: 'post',
    data
  })
}

// 修改密码
export function changePassword(data) {
  return request({
    url: '/user/changePassword',
    method: 'post',
    data
  })
}

// 获取当前用户信息
export function getCurrentUser(data = {}) {
  return request({
    url: '/user/getCurrentUser',
    method: 'post',
    data
  })
}

// OAuth 登录 - 获取授权URL
export function getOAuthLoginUrl(provider) {
  return request({
    url: `/oauth/${provider}/login`,
    method: 'get',
    params: {}
  })
}

// OAuth 回调
export function oauthCallback(provider, code, state) {
  return request({
    url: `/oauth/${provider}/callback`,
    method: 'get',
    params: { code, state }
  })
}

// 设置用户角色
export function setUserRoles(data) {
  return request({
    url: '/user/setRoles',
    method: 'post',
    data
  })
}

// 获取用户角色
export function getUserRoles(data) {
  return request({
    url: '/user/getRoles',
    method: 'post',
    data
  })
}

// 获取当前用户权限
export function getCurrentUserPermissions(data = {}) {
  return request({
    url: '/user/getCurrentUserPermissions',
    method: 'post',
    data
  })
}

// 获取当前用户角色
export function getCurrentUserRoles(data = {}) {
  return request({
    url: '/user/getCurrentUserRoles',
    method: 'post',
    data
  })
}

// 用户注册
export function register(data) {
  return request({
    url: '/user/register',
    method: 'post',
    data
  })
}