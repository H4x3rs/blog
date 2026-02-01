// 生成 sitemap.xml
export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const baseUrl = config.public.siteUrl || 'https://example.com'
  const apiBase = config.apiBase || 'http://localhost:8000/api'

  // 获取所有文章
  let articles: any[] = []
  try {
    const res = await $fetch<{ code: number; data: { list: any[] } }>(`${apiBase}/article/getList`, {
      method: 'POST',
      body: { page: 1, size: 1000, status: 'published' }
    })
    if (res.code === 0) {
      articles = res.data.list || []
    }
  } catch (error) {
    console.error('获取文章列表失败:', error)
  }

  // 获取所有分类
  let categories: any[] = []
  try {
    const res = await $fetch<{ code: number; data: { list: any[] } }>(`${apiBase}/category/getList`, {
      method: 'POST',
      body: { page: 1, size: 100 }
    })
    if (res.code === 0) {
      categories = res.data.list || []
    }
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }

  // 获取所有标签
  let tags: any[] = []
  try {
    const res = await $fetch<{ code: number; data: { list: any[] } }>(`${apiBase}/tag/getList`, {
      method: 'POST',
      body: { page: 1, size: 100 }
    })
    if (res.code === 0) {
      tags = res.data.list || []
    }
  } catch (error) {
    console.error('获取标签列表失败:', error)
  }

  // 获取所有专题
  let topics: any[] = []
  try {
    const res = await $fetch<{ code: number; data: { list: any[] } }>(`${apiBase}/topic/getList`, {
      method: 'POST',
      body: { page: 1, size: 100 }
    })
    if (res.code === 0) {
      topics = res.data.list || []
    }
  } catch (error) {
    console.error('获取专题列表失败:', error)
  }

  // 生成 sitemap XML
  const urls: string[] = []

  // 静态页面
  const staticPages = [
    { loc: '/', priority: '1.0', changefreq: 'daily' },
    { loc: '/category', priority: '0.8', changefreq: 'weekly' },
    { loc: '/tag', priority: '0.8', changefreq: 'weekly' },
    { loc: '/topic', priority: '0.8', changefreq: 'weekly' },
    { loc: '/about', priority: '0.6', changefreq: 'monthly' }
  ]

  for (const page of staticPages) {
    urls.push(`
    <url>
      <loc>${baseUrl}${page.loc}</loc>
      <changefreq>${page.changefreq}</changefreq>
      <priority>${page.priority}</priority>
    </url>`)
  }

  // 文章页面
  for (const article of articles) {
    const lastmod = article.updatedAt || article.createdAt
    urls.push(`
    <url>
      <loc>${baseUrl}/article/${article.id}</loc>
      <lastmod>${new Date(lastmod).toISOString().split('T')[0]}</lastmod>
      <changefreq>weekly</changefreq>
      <priority>0.9</priority>
    </url>`)
  }

  // 分类页面
  for (const category of categories) {
    urls.push(`
    <url>
      <loc>${baseUrl}/category/${category.id}</loc>
      <changefreq>weekly</changefreq>
      <priority>0.7</priority>
    </url>`)
  }

  // 标签页面
  for (const tag of tags) {
    urls.push(`
    <url>
      <loc>${baseUrl}/tag/${tag.id}</loc>
      <changefreq>weekly</changefreq>
      <priority>0.7</priority>
    </url>`)
  }

  // 专题页面
  for (const topic of topics) {
    urls.push(`
    <url>
      <loc>${baseUrl}/topic/${topic.id}</loc>
      <changefreq>weekly</changefreq>
      <priority>0.8</priority>
    </url>`)
  }

  const sitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${urls.join('')}
</urlset>`

  setResponseHeader(event, 'Content-Type', 'application/xml')
  return sitemap
})
