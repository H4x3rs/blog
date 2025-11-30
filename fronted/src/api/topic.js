import request from '../utils/request'

// 获取专题列表
export function getTopicList(data) {
  return request({
    url: '/topic/getList',
    method: 'post',
    data
  })
}

// 获取单个专题
export function getTopic(data) {
  return request({
    url: '/topic/getOne',
    method: 'post',
    data
  })
}

// 创建专题
export function createTopic(data) {
  return request({
    url: '/topic/create',
    method: 'post',
    data
  })
}

// 更新专题
export function updateTopic(data) {
  return request({
    url: '/topic/update',
    method: 'post',
    data
  })
}

// 删除专题
export function deleteTopic(data) {
  return request({
    url: '/topic/delete',
    method: 'post',
    data
  })
}

// 获取专题的文章列表
export function getTopicArticles(data) {
  return request({
    url: '/topic/getTopicArticles',
    method: 'post',
    data
  })
}

// 添加文章到专题
export function addArticleToTopic(data) {
  return request({
    url: '/topic/addArticleToTopic',
    method: 'post',
    data
  })
}

// 从专题移除文章
export function removeArticleFromTopic(data) {
  return request({
    url: '/topic/removeArticleFromTopic',
    method: 'post',
    data
  })
}

// 更新专题文章排序
export function updateTopicArticleSort(data) {
  return request({
    url: '/topic/updateTopicArticleSort',
    method: 'post',
    data
  })
}
