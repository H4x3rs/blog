import request from '../utils/request'

// 获取标签列表
export function getTagList(data) {
  return request({
    url: '/tag/getList',
    method: 'post',
    data
  })
}

// 获取单个标签
export function getTag(data) {
  return request({
    url: '/tag/getOne',
    method: 'post',
    data
  })
}

// 创建标签
export function createTag(data) {
  return request({
    url: '/tag/create',
    method: 'post',
    data
  })
}

// 更新标签
export function updateTag(data) {
  return request({
    url: '/tag/update',
    method: 'post',
    data
  })
}

// 删除标签
export function deleteTag(data) {
  return request({
    url: '/tag/delete',
    method: 'post',
    data
  })
}

// 通过 slug 获取标签
export function getTagBySlug(data) {
  return request({
    url: '/tag/getBySlug',
    method: 'post',
    data
  })
}
