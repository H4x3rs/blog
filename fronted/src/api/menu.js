import request from '../utils/request'

// 获取菜单列表
export function getMenuList(data = {}) {
  return request({
    url: '/menu/getList',
    method: 'post',
    data
  })
}

// 获取单个菜单
export function getMenu(data) {
  return request({
    url: '/menu/getOne',
    method: 'post',
    data
  })
}

// 创建菜单
export function createMenu(data) {
  return request({
    url: '/menu/create',
    method: 'post',
    data
  })
}

// 更新菜单
export function updateMenu(data) {
  return request({
    url: '/menu/update',
    method: 'post',
    data
  })
}

// 删除菜单
export function deleteMenu(data) {
  return request({
    url: '/menu/delete',
    method: 'post',
    data
  })
}

// 获取当前用户的菜单（根据权限过滤）
export function getCurrentUserMenus(data = {}) {
  return request({
    url: '/menu/getCurrentUserMenus',
    method: 'post',
    data
  })
}