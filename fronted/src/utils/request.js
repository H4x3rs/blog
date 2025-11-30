import axios from 'axios'

const service = axios.create({
  baseURL: '/api', // Proxy will handle this
  timeout: 5000
})

service.interceptors.request.use(
  config => {
    // 自动添加token到请求头
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    // 如果是 FormData，不设置 Content-Type，让浏览器自动设置（包含 boundary）
    if (config.data instanceof FormData) {
      delete config.headers['Content-Type']
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

service.interceptors.response.use(
  response => {
    const res = response.data
    // GoFrame 响应格式: {code, message, data}
    if (res.code === 0 || res.code === undefined) {
      // 成功响应，返回 data 字段，如果没有 data 字段则返回整个响应
      return res.data !== undefined ? res.data : res
    } else {
      // 业务错误，抛出错误信息
      const error = new Error(res.message || '请求失败')
      error.code = res.code
      return Promise.reject(error)
    }
  },
  error => {
    // HTTP 错误
    const message = error.response?.data?.message || error.message || '网络错误'
    return Promise.reject(new Error(message))
  }
)

export default service

