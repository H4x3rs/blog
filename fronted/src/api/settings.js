import request from '../utils/request'

// 获取系统设置
export function getSettings(data = {}) {
  return request({
    url: '/settings/getSettings',
    method: 'post',
    data
  })
}

// 更新系统设置
export function updateSettings(data) {
  return request({
    url: '/settings/updateSettings',
    method: 'post',
    data
  })
}

// 获取前台 Banner 标题和副标题
export function getBanner(data = {}) {
  return request({
    url: '/settings/getBanner',
    method: 'post',
    data
  })
}

// 获取关于我页面信息
export function getAbout(data = {}) {
  return request({
    url: '/settings/getAbout',
    method: 'post',
    data
  })
}
