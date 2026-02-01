<template>
  <div class="article-detail-container" v-loading="loading">
    <!-- Header Area -->
    <div class="article-header">
      <div class="header-content">
        <div class="meta-info">
          <el-tag effect="dark" round class="category-tag">{{ article.category }}</el-tag>
          <el-tag 
            v-if="article.topic" 
            effect="dark" 
            round 
            class="topic-tag"
            @click="goToTopic"
          >
            📚 {{ article.topic.name }}
          </el-tag>
          <span class="date">{{ article.createdAt }}</span>
        </div>
        <h1 class="title">{{ article.title }}</h1>
        <div class="author-info">
          <el-avatar 
            :size="32" 
            :src="article.publishedByUser?.avatar || 'https://picsum.photos/id/64/100/100'"
            @click="goToAuthorArticles(article.publishedByUser)"
            style="cursor: pointer;"
          />
          <span 
            class="author-name"
            @click="goToAuthorArticles(article.publishedByUser)"
            style="cursor: pointer;"
          >
            {{ article.publishedByUser ? (article.publishedByUser.nickname || article.publishedByUser.username) : 'Admin' }}
          </span>
          <span class="divider">·</span>
          <span class="read-time">5 min read</span>
          <span class="divider">·</span>
          <span class="views">{{ article.views }} views</span>
        </div>
      </div>
      <div class="header-image" :style="{ backgroundImage: `url(${article.cover})` }"></div>
    </div>

    <!-- Content Area -->
    <el-row justify="center" class="content-wrapper">
      <el-col :span="16" :xs="22">
        <!-- 专题信息 -->
        <el-card v-if="article.topic" class="topic-info-card" shadow="never">
          <div class="topic-info" @click="goToTopic">
            <div class="topic-icon">📚</div>
            <div class="topic-content">
              <div class="topic-label">本文属于专题</div>
              <div class="topic-name">{{ article.topic.name }}</div>
              <div class="topic-desc">{{ article.topic.description }}</div>
              <div class="topic-stats">
                <span>{{ article.topic.articleCount }} 篇文章</span>
                <span class="divider">·</span>
                <span>{{ article.topic.progress }}</span>
              </div>
            </div>
            <div class="topic-action">
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>
        </el-card>

        <el-card class="content-card" shadow="never">
          <div class="markdown-body" v-html="renderedContent"></div>
          
          <div class="article-footer">
            <div class="tags-list">
              <el-tag 
                v-for="tag in article.tags" 
                :key="typeof tag === 'object' ? tag.id : tag" 
                class="tag-item" 
                effect="plain"
                @click="goToTag(tag)"
                style="cursor: pointer;"
              >
                # {{ typeof tag === 'object' ? tag.name : tag }}
              </el-tag>
            </div>
            <div class="actions">
               <el-button type="primary" plain round size="small">
                  <el-icon><Star /></el-icon> 点赞
               </el-button>
               <el-button round size="small">
                  <el-icon><Share /></el-icon> 分享
               </el-button>
            </div>
          </div>
        </el-card>

        <!-- Comments Area (Placeholder) -->
        <div class="comments-section">
          <h3>评论 (0)</h3>
          <el-empty description="暂无评论，快来抢沙发吧！" />
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Star, Share } from '@element-plus/icons-vue'
import { marked } from 'marked'
import { getArticle } from '@/api/article'
import { ElMessage } from 'element-plus'
import { updateSEO, generateArticleSEO } from '@/utils/seo'
import { useSiteConfig } from '@/store/site'

const route = useRoute()
const router = useRouter()
const articleId = route.params.id

const article = ref({
  id: articleId,
  title: '',
  category: '',
  createdAt: '',
  views: 0,
  cover: '',
  tags: [],
  content: ''
})

const loading = ref(true)
const { siteName, bannerSubtitle } = useSiteConfig()

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  try {
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return dateStr
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  } catch (error) {
    if (typeof dateStr === 'string' && /^\d{4}-\d{2}-\d{2}/.test(dateStr)) {
      return dateStr.split(' ')[0] || dateStr.split('T')[0]
    }
    return dateStr
  }
}

// 加载文章详情
const loadArticle = async () => {
  loading.value = true
  try {
    const res = await getArticle({ id: parseInt(articleId) })
    if (res) {
      article.value = {
        id: res.id,
        title: res.title || '',
        category: res.categoryName || '',
        createdAt: formatDate(res.createdAt),
        views: res.views || 0,
        cover: res.coverImage || 'https://picsum.photos/id/' + res.id + '/1200/600',
        tags: res.tags && Array.isArray(res.tags) ? res.tags : [],
        content: res.content || '',
        publishedByUser: res.publishedByUser || null,
        topic: res.topic || null,
        updatedAt: res.updatedAt || res.createdAt
      }
      
      // 更新SEO信息
      updateArticleSEO(res)
    }
  } catch (error) {
    console.error('加载文章失败:', error)
    ElMessage.error('加载文章失败')
  } finally {
    loading.value = false
  }
}

// 更新文章SEO信息
const updateArticleSEO = (articleData) => {
  const siteUrl = window.location.origin
  
  // 使用文章的 coverImage 字段
  const coverImage = articleData.coverImage || null
  
  // 使用当前页面的完整URL（包含hash路由）
  const articleUrl = window.location.href
  
  const seoConfig = generateArticleSEO(
    {
      id: articleData.id,
      title: articleData.title || '',
      desc: articleData.desc || articleData.title || '',
      coverImage: coverImage,
      tags: articleData.tags ? (typeof articleData.tags === 'string' ? articleData.tags.split(',') : articleData.tags) : [],
      category: articleData.categoryName || '',
      createdAt: articleData.createdAt || '',
      updatedAt: articleData.updatedAt || articleData.createdAt || '',
      publishedByUser: articleData.publishedByUser || null
    },
    siteName.value || 'Blog System',
    bannerSubtitle.value || '分享编程心得，记录技术成长',
    siteUrl
  )
  
  // 覆盖URL为当前页面的完整URL（包含hash）
  seoConfig.url = articleUrl
  
  // 更新SEO meta标签
  // - title: 使用文章标题
  // - image: 使用文章的 coverImage（如果存在），否则使用默认logo
  // - Open Graph 和 Twitter Card 都会使用文章的封面图片
  updateSEO(seoConfig)
  
  // 浏览器标签页的 title 使用带网站名称的完整标题
  document.title = `${articleData.title || '文章详情'} - ${siteName.value || 'Blog System'}`
}

// 监听路由参数变化，确保切换文章时也能更新SEO
watch(() => route.params.id, (newId) => {
  if (newId) {
    loadArticle()
  }
}, { immediate: false })

// 处理专题链接点击
onMounted(() => {
  loadArticle()
  nextTick(() => {
    document.addEventListener('click', handleLinkClick)
  })
})

const handleLinkClick = (e) => {
  const target = e.target.closest('a.topic-link')
  if (target) {
    e.preventDefault()
    const topicId = target.getAttribute('data-topic-id')
    if (topicId) {
      router.push(`/topic/${topicId}`)
    }
  }
}

// 跳转到作者文章列表
const goToAuthorArticles = (user) => {
  if (user && user.id) {
    router.push(`/author/${user.id}`)
  }
}

// 跳转到标签页面
const goToTag = (tag) => {
  if (!tag) return
  const tagId = typeof tag === 'object' ? tag.id : tag
  const tagName = typeof tag === 'object' ? tag.name : tag
  if (tagId) {
    router.push(`/tag/${tagId}`)
  } else if (tagName) {
    router.push(`/tag?name=${encodeURIComponent(tagName)}`)
  }
}

const renderedContent = computed(() => {
  if (!article.value.content) return ''
  let html = marked.parse(article.value.content)
  
  // 处理专题链接，确保使用路由导航而不是页面刷新
  html = html.replace(
    /<a href="\/topic\/(\d+)"/g, 
    '<a href="/topic/$1" class="topic-link" data-topic-id="$1"'
  )
  
  return html
})
</script>

<style scoped>
.article-detail-container {
  background-color: #f9f9f9;
  min-height: 100vh;
  padding-bottom: 60px;
}

.article-header {
  position: relative;
  height: 400px;
  background-color: #2c3e50;
  color: white;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}

.header-image {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-size: cover;
  background-position: center;
  opacity: 0.4;
  filter: blur(4px);
  transform: scale(1.1);
}

.header-content {
  position: relative;
  z-index: 1;
  text-align: center;
  max-width: 800px;
  padding: 0 20px;
}

.meta-info {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
}

.category-tag {
  background-color: #409eff;
  border-color: #409eff;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.title {
  font-size: 36px;
  font-weight: 800;
  line-height: 1.3;
  margin-bottom: 24px;
  text-shadow: 0 2px 10px rgba(0,0,0,0.3);
}

.author-info {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  opacity: 0.9;
  flex-wrap: wrap; /* 允许换行 */
}

.author-name {
  margin-left: 8px;
  font-weight: 600;
}

.divider {
  margin: 0 10px;
  opacity: 0.6;
}

.content-wrapper {
  margin-top: -60px;
  position: relative;
  z-index: 2;
}

.content-card {
  padding: 40px;
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.05);
  background: white;
  border: none;
}

/* Markdown Styles Override - Simple Version */
.markdown-body {
  font-size: 16px;
  line-height: 1.8;
  color: #333;
  overflow-wrap: break-word; /* 防止长单词溢出 */
}
:deep(.markdown-body h2) {
  margin-top: 30px;
  margin-bottom: 15px;
  font-size: 24px;
  font-weight: 700;
  color: #1a1a1a;
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
}
:deep(.markdown-body p) {
  margin-bottom: 16px;
}
:deep(.markdown-body pre) {
  background: #f6f8fa;
  padding: 16px;
  border-radius: 8px;
  overflow: auto;
  font-family: monospace;
}
:deep(.markdown-body code) {
  background: #f0f2f5;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.9em;
  color: #d63384;
}
:deep(.markdown-body pre code) {
  background: transparent;
  color: inherit;
  padding: 0;
}
/* 无序列表样式 */
:deep(.markdown-body ul) {
  padding-left: 20px;
  margin-bottom: 16px;
  list-style-type: disc;
  list-style-position: outside;
}
:deep(.markdown-body ul li) {
  margin-bottom: 8px;
  line-height: 1.8;
  list-style-type: disc;
}
:deep(.markdown-body ul ul) {
  list-style-type: circle;
  margin-top: 8px;
  margin-bottom: 8px;
}
:deep(.markdown-body ul ul ul) {
  list-style-type: square;
}
/* 有序列表样式 */
:deep(.markdown-body ol) {
  padding-left: 20px;
  margin-bottom: 16px;
  list-style-type: decimal;
  list-style-position: outside;
}
:deep(.markdown-body ol li) {
  margin-bottom: 8px;
  line-height: 1.8;
  list-style-type: decimal;
}
:deep(.markdown-body ol ol) {
  list-style-type: lower-alpha;
  margin-top: 8px;
  margin-bottom: 8px;
}
:deep(.markdown-body ol ol ol) {
  list-style-type: lower-roman;
}
:deep(.markdown-body img) {
    max-width: 100%;
    height: auto;
    border-radius: 8px;
}
:deep(.markdown-body blockquote) {
  margin: 20px 0;
  padding: 16px 20px;
  background: #f0f9ff;
  border-left: 4px solid #409eff;
  border-radius: 4px;
  color: #606266;
}
/* 普通超链接样式 */
:deep(.markdown-body a:not(.topic-link)) {
  color: #409eff;
  text-decoration: underline;
  text-decoration-color: rgba(64, 158, 255, 0.3);
  transition: all 0.3s ease;
}
:deep(.markdown-body a:not(.topic-link):hover) {
  color: #66b1ff;
  text-decoration-color: rgba(102, 177, 255, 0.6);
}
:deep(.markdown-body a:not(.topic-link):visited) {
  color: #7c3aed;
}
/* 专题链接样式 */
:deep(.markdown-body a.topic-link) {
  color: #409eff;
  text-decoration: none;
  font-weight: 600;
  border-bottom: 1px dashed #409eff;
  transition: all 0.3s ease;
}
:deep(.markdown-body a.topic-link:hover) {
  color: #66b1ff;
  border-bottom-style: solid;
}


.article-footer {
  margin-top: 40px;
  padding-top: 20px;
  border-top: 1px solid #f0f2f5;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 20px;
}

.tags-list {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.comments-section {
  margin-top: 30px;
}
.comments-section h3 {
  margin-bottom: 20px;
  font-size: 18px;
  color: #303133;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .article-header {
      height: 300px;
  }
  .title {
      font-size: 24px;
  }
  .content-card {
      padding: 20px;
  }
  .article-footer {
      flex-direction: column;
      align-items: flex-start;
  }
  .actions {
      width: 100%;
      display: flex;
      justify-content: space-between;
  }
  .meta-info {
      flex-wrap: wrap;
  }
}
</style>
