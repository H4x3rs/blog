import request from '../utils/request'

// 获取文章列表
export function getArticleList(data) {
  return request({
    url: '/article/getList',
    method: 'post',
    data
  })
}

export function getArticle(data) {
  return request({
    url: '/article/getOne',
    method: 'post',
    data
  })
}

export function createArticle(data) {
  return request({
    url: '/article/create',
    method: 'post',
    data
  })
}

export function updateArticle(data) {
  return request({
    url: '/article/update',
    method: 'post',
    data
  })
}

export function deleteArticle(data) {
  return request({
    url: '/article/delete',
    method: 'post',
    data
  })
}
