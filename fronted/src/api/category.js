import request from '../utils/request'

// 获取分类列表
export function getCategoryList(data) {
  return request({
    url: '/category/getList',
    method: 'post',
    data
  })
}

// 获取单个分类
export function getCategory(data) {
  return request({
    url: '/category/getOne',
    method: 'post',
    data
  })
}

// 创建分类
export function createCategory(data) {
  return request({
    url: '/category/create',
    method: 'post',
    data
  })
}

// 更新分类
export function updateCategory(data) {
  return request({
    url: '/category/update',
    method: 'post',
    data
  })
}

// 删除分类
export function deleteCategory(data) {
  return request({
    url: '/category/delete',
    method: 'post',
    data
  })
}

// 通过 slug 获取分类
export function getCategoryBySlug(data) {
  return request({
    url: '/category/getBySlug',
    method: 'post',
    data
  })
}
