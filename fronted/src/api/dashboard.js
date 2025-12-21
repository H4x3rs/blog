import request from '../utils/request'

// 获取统计数据
export function getDashboardStats() {
  return request({
    url: '/dashboard/stats',
    method: 'post',
    data: {}
  })
}

// 获取访问趋势
export function getDashboardTrend(data) {
  return request({
    url: '/dashboard/trend',
    method: 'post',
    data
  })
}

// 获取分类占比
export function getDashboardCategoryRatio() {
  return request({
    url: '/dashboard/categoryRatio',
    method: 'post',
    data: {}
  })
}

// 获取最新评论
export function getDashboardLatestComments(data) {
  return request({
    url: '/dashboard/latestComments',
    method: 'post',
    data
  })
}


