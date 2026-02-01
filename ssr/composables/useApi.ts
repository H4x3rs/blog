// API 请求封装
export const useApi = () => {
  const config = useRuntimeConfig()
  
  // 服务端使用完整URL，客户端使用相对路径
  const getBaseUrl = () => {
    if (import.meta.server) {
      return config.apiBase
    }
    return config.public.apiBase
  }

  const request = async <T>(url: string, options: {
    method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
    data?: any
    headers?: Record<string, string>
  } = {}): Promise<T> => {
    const { method = 'POST', data, headers = {} } = options
    
    // 获取 token（仅客户端）
    let token = ''
    if (import.meta.client) {
      token = localStorage.getItem('token') || ''
    }

    const fetchOptions: any = {
      method,
      headers: {
        'Content-Type': 'application/json',
        ...headers,
        ...(token ? { Authorization: `Bearer ${token}` } : {})
      }
    }

    if (data && method !== 'GET') {
      fetchOptions.body = JSON.stringify(data)
    }

    const response = await $fetch<{
      code: number
      message: string
      data: T
    }>(`${getBaseUrl()}${url}`, fetchOptions)

    // GoFrame 响应格式: {code, message, data}
    if (response.code === 0 || response.code === undefined) {
      return (response.data !== undefined ? response.data : response) as T
    } else {
      throw new Error(response.message || '请求失败')
    }
  }

  return { request }
}

// 文章相关 API
export const useArticleApi = () => {
  const { request } = useApi()

  return {
    // 获取文章列表
    getList: (data: {
      page?: number
      size?: number
      status?: string
      categoryId?: number
      tagId?: number
      authorId?: number
      onlyMine?: boolean
    } = {}) => request<{
      list: any[]
      total: number
    }>('/article/getList', { data }),

    // 获取单篇文章
    getOne: (data: { id: number }) => request<any>('/article/getOne', { data })
  }
}

// 分类相关 API
export const useCategoryApi = () => {
  const { request } = useApi()

  return {
    getList: (data: { page?: number; size?: number } = {}) => 
      request<{ list: any[]; total: number }>('/category/getList', { data }),
    
    getOne: (data: { id?: number; slug?: string }) => 
      request<any>('/category/getOne', { data }),
    
    getBySlug: (data: { slug: string }) => 
      request<any>('/category/getBySlug', { data })
  }
}

// 标签相关 API
export const useTagApi = () => {
  const { request } = useApi()

  return {
    getList: (data: { page?: number; size?: number } = {}) => 
      request<{ list: any[]; total: number }>('/tag/getList', { data }),
    
    getOne: (data: { id?: number; slug?: string }) => 
      request<any>('/tag/getOne', { data }),
    
    getBySlug: (data: { slug: string }) => 
      request<any>('/tag/getBySlug', { data })
  }
}

// 专题相关 API
export const useTopicApi = () => {
  const { request } = useApi()

  return {
    getList: (data: { page?: number; size?: number } = {}) => 
      request<{ list: any[]; total: number }>('/topic/getList', { data }),
    
    getOne: (data: { id: number }) => 
      request<any>('/topic/getOne', { data }),
    
    getTopicArticles: (data: { topicId: number; page?: number; size?: number }) => 
      request<{ list: any[]; total: number }>('/topic/getTopicArticles', { data })
  }
}

// 设置相关 API
export const useSettingsApi = () => {
  const { request } = useApi()

  return {
    getSettings: () => request<any>('/settings/getSettings', { data: {} }),
    getBanner: () => request<any>('/settings/getBanner', { data: {} }),
    getAbout: () => request<any>('/settings/getAbout', { data: {} })
  }
}
