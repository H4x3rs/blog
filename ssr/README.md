# Blog SSR - Nuxt 3

基于 Nuxt 3 的博客系统 SSR 版本，提供更好的 SEO 支持和首屏加载性能。

## 特性

- 🚀 **服务端渲染 (SSR)** - 更好的 SEO 和首屏性能
- 📱 **响应式设计** - 完美适配移动端和桌面端
- 🎨 **Element Plus** - 美观的 UI 组件库
- 🔍 **SEO 优化** - 自动生成 sitemap、meta 标签
- ⚡ **Vite 构建** - 极速的开发体验

## 技术栈

- [Nuxt 3](https://nuxt.com/) - Vue 3 全栈框架
- [Vue 3](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Element Plus](https://element-plus.org/) - Vue 3 组件库
- [Pinia](https://pinia.vuejs.org/) - Vue 状态管理
- [TypeScript](https://www.typescriptlang.org/) - 类型安全

## 快速开始

### 安装依赖

```bash
npm install
```

### 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

### 开发模式

```bash
npm run dev
```

访问 http://localhost:3000

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

### 生成静态站点

```bash
npm run generate
```

## 目录结构

```
ssr/
├── assets/              # 静态资源（会被构建处理）
│   └── scss/           # SCSS 样式文件
├── composables/        # 组合式函数
│   └── useApi.ts       # API 请求封装
├── layouts/            # 布局组件
│   └── default.vue     # 默认布局
├── pages/              # 页面组件（自动路由）
│   ├── index.vue       # 首页
│   ├── article/        # 文章相关页面
│   ├── category/       # 分类相关页面
│   ├── tag/            # 标签相关页面
│   ├── topic/          # 专题相关页面
│   └── about.vue       # 关于页面
├── public/             # 静态资源（不会被构建处理）
├── server/             # 服务端代码
│   └── api/            # API 路由
├── stores/             # Pinia 状态管理
├── app.vue             # 应用入口
├── nuxt.config.ts      # Nuxt 配置
└── package.json        # 项目配置
```

## API 代理配置

开发环境下，Nuxt 会自动将 `/api` 请求代理到后端服务器。

生产环境需要配置 Nginx 反向代理：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    location /api {
        proxy_pass http://localhost:8000/api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 部署

### PM2 部署

```bash
# 构建
npm run build

# 使用 PM2 启动
pm2 start .output/server/index.mjs --name blog-ssr
```

### Docker 部署

```dockerfile
FROM node:20-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]
```

## 与原 SPA 版本的区别

| 特性 | SPA (fronted) | SSR (ssr) |
|------|---------------|-----------|
| 首屏渲染 | 客户端渲染 | 服务端渲染 |
| SEO | 需要预渲染 | 原生支持 |
| 首屏性能 | 较慢 | 较快 |
| 服务器要求 | 静态服务器 | Node.js 服务器 |
| 路由模式 | Hash 路由 | History 路由 |

## License

MIT
