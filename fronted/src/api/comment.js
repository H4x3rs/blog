import request from '../utils/request'

// 创建评论
export function createComment(data) {
  return request({
    url: '/comment/create',
    method: 'post',
    data
  })
}

// 获取评论列表
export function getCommentList(data) {
  return request({
    url: '/comment/getList',
    method: 'post',
    data
  })
}

// 更新评论
export function updateComment(data) {
  return request({
    url: '/comment/update',
    method: 'post',
    data
  })
}

// 删除评论
export function deleteComment(data) {
  return request({
    url: '/comment/delete',
    method: 'post',
    data
  })
}

// 获取管理评论列表
export function getCommentManageList(data) {
  return request({
    url: '/comment/getManageList',
    method: 'post',
    data
  })
}
