/**
 * SEO工具函数
 * 用于动态更新页面的SEO meta标签
 */

/**
 * 更新页面SEO信息
 * @param {Object} options SEO配置选项
 * @param {string} options.title 页面标题
 * @param {string} options.description 页面描述
 * @param {string} options.keywords 关键词，多个用逗号分隔
 * @param {string} options.image 页面图片（用于Open Graph）
 * @param {string} options.url 页面URL
 * @param {string} options.type 页面类型（article, website等）
 * @param {Object} options.article 文章信息（如果是文章页）
 */
export function updateSEO(options = {}) {
  const {
    title = '',
    description = '',
    keywords = '',
    image = '',
    url = window.location.href,
    type = 'website',
    article = null
  } = options

  // 更新页面标题（title标签不是meta标签，需要特殊处理）
  if (title) {
    document.title = title
  }
  
  // 更新基础meta标签
  updateMetaTag('description', description, 'name')
  updateMetaTag('keywords', keywords, 'name')

  // 更新Open Graph标签
  updateMetaTag('og:title', title, 'property')
  updateMetaTag('og:description', description, 'property')
  updateMetaTag('og:image', image, 'property') // 使用文章的 coverImage
  updateMetaTag('og:url', url, 'property')
  updateMetaTag('og:type', type, 'property')
  
  // 添加 og:image:width 和 og:image:height（可选，有助于社交媒体平台显示）
  if (image) {
    updateMetaTag('og:image:width', '1200', 'property')
    updateMetaTag('og:image:height', '630', 'property')
  }

  // 更新Twitter Card标签
  updateMetaTag('twitter:card', 'summary_large_image', 'name')
  updateMetaTag('twitter:title', title, 'name')
  updateMetaTag('twitter:description', description, 'name')
  updateMetaTag('twitter:image', image, 'name') // 使用文章的 coverImage

  // 如果是文章，添加文章相关的meta标签
  if (article && type === 'article') {
    updateMetaTag('article:published_time', article.publishedTime, 'property')
    updateMetaTag('article:modified_time', article.modifiedTime, 'property')
    updateMetaTag('article:author', article.author, 'property')
    if (article.tags && article.tags.length > 0) {
      article.tags.forEach((tag, index) => {
        const tagName = typeof tag === 'object' ? tag.name : tag
        updateMetaTag(`article:tag`, tagName, 'property', index)
      })
    }
    if (article.category) {
      updateMetaTag('article:section', article.category, 'property')
    }
  }

  // 添加结构化数据（JSON-LD）
  addStructuredData(options)
}

/**
 * 更新或创建meta标签
 * @param {string} name meta标签名称
 * @param {string} content 内容
 * @param {string} attr 属性名（name或property）
 * @param {number} index 索引（用于多个相同名称的标签）
 */
function updateMetaTag(name, content, attr = 'name', index = null) {
  // 对于图片等标签，即使为空也要更新（设置为空字符串）
  // 但对于其他标签，如果内容为空则跳过
  const allowEmpty = ['og:image', 'twitter:image', 'image'].includes(name)
  if (!content && !allowEmpty) return

  const selector = index !== null 
    ? `meta[${attr}="${name}"]:nth-of-type(${index + 1})`
    : `meta[${attr}="${name}"]`
  
  let meta = document.querySelector(selector)
  
  if (!meta) {
    meta = document.createElement('meta')
    meta.setAttribute(attr, name)
    document.head.appendChild(meta)
  }
  
  meta.setAttribute('content', content || '')
}

/**
 * 添加结构化数据（JSON-LD）
 * @param {Object} options SEO配置选项
 */
function addStructuredData(options) {
  // 移除旧的JSON-LD
  const oldScript = document.querySelector('script[type="application/ld+json"]')
  if (oldScript) {
    oldScript.remove()
  }

  const { title, description, image, url, type, article } = options

  let structuredData = {}

  if (type === 'article' && article) {
    // 文章结构化数据
    structuredData = {
      '@context': 'https://schema.org',
      '@type': 'BlogPosting',
      headline: title,
      description: description,
      image: image ? [image] : [],
      datePublished: article.publishedTime,
      dateModified: article.modifiedTime || article.publishedTime,
      author: {
        '@type': 'Person',
        name: article.author || 'Admin'
      },
      publisher: {
        '@type': 'Organization',
        name: article.siteName || 'Blog System',
        logo: {
          '@type': 'ImageObject',
          url: image || ''
        }
      },
      mainEntityOfPage: {
        '@type': 'WebPage',
        '@id': url
      }
    }

    if (article.category) {
      structuredData.articleSection = article.category
    }

    if (article.tags && article.tags.length > 0) {
      structuredData.keywords = article.tags.map(tag => 
        typeof tag === 'object' ? tag.name : tag
      ).join(', ')
    }
  } else {
    // 网站结构化数据
    structuredData = {
      '@context': 'https://schema.org',
      '@type': 'WebSite',
      name: article?.siteName || 'Blog System',
      description: description,
      url: url
    }
  }

  // 添加JSON-LD脚本
  const script = document.createElement('script')
  script.type = 'application/ld+json'
  script.textContent = JSON.stringify(structuredData)
  document.head.appendChild(script)
}

/**
 * 设置页面标题
 * @param {string} title 页面标题
 * @param {string} siteName 网站名称
 */
export function setPageTitle(title, siteName = 'Blog System') {
  if (title) {
    document.title = `${title} - ${siteName}`
  } else {
    document.title = siteName
  }
}

/**
 * 生成文章页面的SEO配置
 * @param {Object} article 文章对象
 * @param {string} siteName 网站名称
 * @param {string} siteDesc 网站描述
 * @param {string} siteUrl 网站URL
 * @returns {Object} SEO配置对象
 */
export function generateArticleSEO(article, siteName, siteDesc, siteUrl) {
  // 使用 hash 路由格式（因为使用的是 createWebHashHistory）
  const articleUrl = `${siteUrl}/#/article/${article.id}`
  
  // 处理文章封面图片（coverImage），确保使用正确的图片路径
  // 优先使用文章的 coverImage，如果没有则使用默认logo
  let articleImage = article.coverImage || null
  
  // 检查 coverImage 是否有效（不为 null、undefined、空字符串）
  const hasValidCoverImage = articleImage && 
                             typeof articleImage === 'string' && 
                             articleImage.trim() !== ''
  
  if (hasValidCoverImage) {
    // 如果已经是完整的URL（http:// 或 https://），直接使用
    if (articleImage.startsWith('http://') || articleImage.startsWith('https://')) {
      // 已经是完整URL，直接使用
    } else {
      // 相对路径，需要转换为绝对路径
      // 如果路径不是以 / 开头，添加 /
      if (!articleImage.startsWith('/')) {
        articleImage = `/${articleImage}`
      }
      articleImage = `${siteUrl}${articleImage}`
    }
  } else {
    // 文章没有 coverImage 或 coverImage 为空，使用默认logo
    articleImage = `${siteUrl}/logo-512x512.png`
  }

  return {
    title: article.title,
    description: article.desc || article.title,
    keywords: article.tags 
      ? article.tags.map(tag => typeof tag === 'object' ? tag.name : tag).join(', ')
      : '',
    image: articleImage,
    url: articleUrl,
    type: 'article',
    article: {
      publishedTime: article.createdAt,
      modifiedTime: article.updatedAt || article.createdAt,
      author: article.publishedByUser 
        ? (article.publishedByUser.nickname || article.publishedByUser.username)
        : 'Admin',
      category: article.category || '',
      tags: article.tags || [],
      siteName: siteName
    }
  }
}




