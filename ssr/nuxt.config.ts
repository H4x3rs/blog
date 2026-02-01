// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  
  modules: [
    '@element-plus/nuxt',
    '@pinia/nuxt'
  ],

  // SSR 配置
  ssr: true,

  // 运行时配置
  runtimeConfig: {
    // 服务端私有配置（直接请求后端，无 /api 前缀）
    apiBase: process.env.NUXT_API_BASE || 'http://localhost:8000',
    // 客户端公开配置（通过代理）
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || '/api'
    }
  },

  // Nitro 服务器配置 - 客户端请求 /api/* 代理到后端
  nitro: {
    devProxy: {
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,
        rewrite: (path: string) => path.replace(/^\/api/, '')
      }
    }
  },

  // Element Plus 配置
  elementPlus: {
    icon: 'ElIcon',
    importStyle: 'css'
  },

  // 应用配置
  app: {
    head: {
      title: 'Blog System',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: '分享编程心得，记录技术成长。分享 Go, Vue, 云原生等前沿技术心得。' },
        { name: 'keywords', content: 'Go, Vue, Cloud Native, 编程, 技术博客' }
      ],
      link: [
        { rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' }
      ]
    }
  },

  // CSS
  css: [
    '@/assets/css/main.css'
  ],

  compatibilityDate: '2025-01-01'
})
