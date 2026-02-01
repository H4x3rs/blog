<template>
  <div class="home-container">
    <!-- Hero Banner -->
    <div class="hero-section">
      <div class="hero-bg-shape shape-1"></div>
      <div class="hero-bg-shape shape-2"></div>
      <div class="hero-content">
        <h1 class="hero-title">
          <span class="highlight" v-html="formatTitle(siteStore.bannerTitle)"></span>
        </h1>
        <p class="hero-subtitle">{{ siteStore.bannerSubtitle }}</p>
        <div class="hero-actions">
          <el-button type="primary" size="large" round class="start-btn" @click="scrollToContent">
            开始阅读 <el-icon class="el-icon--right"><ArrowRight /></el-icon>
          </el-button>
          <el-button size="large" round class="github-btn" @click="openGithub">
            <el-icon style="margin-right: 6px; font-size: 18px;">
              <svg viewBox="0 0 1024 1024" xmlns="http://www.w3.org/2000/svg" width="1em" height="1em">
                <path fill="currentColor" d="M511.6 76.3C264.3 76.2 64 276.4 64 523.5 64 718.9 189.3 885 363.8 946c23.5 5.9 19.9-10.8 19.9-22.2v-77.5c-135.7 15.9-141.2-73.9-150.3-88.9C215 726 171.5 718 184.5 703c30.9-15.9 62.4 4 98.9 57.9 26.4 39.1 77.9 32.5 104 26 5.7-23.5 17.9-44.5 34.7-60.8-140.6-25.2-199.2-111-199.2-213 0-49.5 16.3-95 48.3-131.7-20.4-60.5 1.9-112.3 4.9-120 58.1-5.2 118.5 41.6 123.2 45.3 33-8.9 70.7-13.6 112.9-13.6 42.4 0 80.2 4.9 113.5 13.9 11.3-8.6 67.3-48.8 121.3-43.9 2.9 7.7 24.7 58.3 5.5 118 32.4 36.8 48.9 82.7 48.9 132.3 0 102.2-59 188.1-200 212.9a127.5 127.5 0 0 1 38.1 91v112.5c.8 9 0 17.9 15 17.9 177.1-59.7 304.6-227 304.6-424.1 0-247.2-200.4-447.3-447.5-447.3z"></path>
              </svg>
            </el-icon>
            GitHub
          </el-button>
        </div>
      </div>
      <div class="hero-illustration">
        <div class="code-window">
          <div class="window-header">
            <div class="dot red"></div>
            <div class="dot yellow"></div>
            <div class="dot green"></div>
          </div>
          <div class="window-body">
            <div class="code-line"><span class="keyword">package</span> main</div>
            <div class="code-line">&nbsp;</div>
            <div class="code-line"><span class="keyword">import</span> (</div>
            <div class="code-line indent"><span class="string">"fmt"</span></div>
            <div class="code-line indent"><span class="string">"github.com/gogf/gf/v2"</span></div>
            <div class="code-line">)</div>
            <div class="code-line">&nbsp;</div>
            <div class="code-line"><span class="keyword">func</span> main() {</div>
            <div class="code-line indent"><span class="comment">// Start the journey</span></div>
            <div class="code-line indent">fmt.Println(<span class="string">"Hello, World!"</span>)</div>
            <div class="code-line">}</div>
          </div>
        </div>
      </div>
    </div>

    <el-row :gutter="40" class="content-wrapper" id="content-start">
      <!-- 左侧文章列表 -->
      <el-col :span="16" :xs="24" class="main-column">
        <div class="article-list">
          <div 
            v-for="(item, index) in articles" 
            :key="item.id" 
            class="article-item"
            :style="{ animationDelay: `${index * 0.1}s` }"
          >
            <NuxtLink :to="`/article/${item.id}`" class="article-thumb">
              <img :src="item.cover" loading="lazy" :alt="item.title" />
              <div v-if="item.category" class="article-cat-tag">{{ item.category }}</div>
            </NuxtLink>
            <div class="article-detail">
              <div class="article-meta-top">
                <span class="date"><el-icon><Calendar /></el-icon> {{ item.createdAt }}</span>
                <span class="views"><el-icon><View /></el-icon> {{ item.views }}</span>
              </div>
              <NuxtLink :to="`/article/${item.id}`">
                <h2 class="article-title">{{ item.title }}</h2>
              </NuxtLink>
              <p class="article-summary">{{ item.desc }}</p>
              <div class="article-footer">
                <NuxtLink 
                  v-if="item.publishedByUser"
                  :to="`/author/${item.publishedByUser.id}`"
                  class="author"
                >
                  <el-avatar :size="24" :src="item.publishedByUser?.avatar || 'https://picsum.photos/id/64/100/100'" />
                  <span>{{ item.publishedByUser.nickname || item.publishedByUser.username }}</span>
                </NuxtLink>
                <span v-else class="author">
                  <el-avatar :size="24" src="https://picsum.photos/id/64/100/100" />
                  <span>Admin</span>
                </span>
                <NuxtLink :to="`/article/${item.id}`">
                  <el-button link type="primary" class="read-more">阅读全文</el-button>
                </NuxtLink>
              </div>
            </div>
          </div>
        </div>
        
        <div class="pagination-wrapper">
          <el-pagination 
            background 
            layout="prev, pager, next" 
            :total="total" 
            :current-page="currentPage"
            :page-size="pageSize"
            @current-change="handlePageChange"
            class="pagination-responsive" 
          />
        </div>
      </el-col>

      <!-- 右侧侧边栏 -->
      <el-col :span="8" :xs="24" class="sidebar-column">
        <div class="sidebar-sticky">
          <!-- 热门文章 -->
          <div class="sidebar-widget hot-widget">
            <h4 class="widget-title">热门文章</h4>
            <ul class="hot-list">
              <li v-for="(item, i) in hotArticles" :key="i" class="hot-item">
                <NuxtLink :to="`/article/${item.id}`">
                  <span class="hot-number">{{ i + 1 }}</span>
                  <div class="hot-info">
                     <div class="hot-text">{{ item.title }}</div>
                     <span class="hot-views">{{ formatViews(item.views) }} reads</span>
                  </div>
                </NuxtLink>
              </li>
            </ul>
          </div>

          <!-- 标签云 -->
          <div class="sidebar-widget tag-widget">
             <h4 class="widget-title">标签</h4>
             <div class="tags-cloud">
               <NuxtLink 
                 v-for="tag in tags" 
                 :key="tag.id" 
                 :to="`/tag/${tag.id}`"
                 class="tag-chip"
               >
                 {{ tag.name }}
               </NuxtLink>
             </div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { Calendar, View, ArrowRight } from '@element-plus/icons-vue'
import { useSiteStore } from '~/stores/site'

// SEO
useSeoMeta({
  title: '首页',
  ogTitle: '首页',
  description: '分享编程心得，记录技术成长。分享 Go, Vue, 云原生等前沿技术心得。',
  ogDescription: '分享编程心得，记录技术成长。分享 Go, Vue, 云原生等前沿技术心得。'
})

const siteStore = useSiteStore()
const articleApi = useArticleApi()
const tagApi = useTagApi()

const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const articles = ref<any[]>([])
const hotArticles = ref<any[]>([])
const tags = ref<any[]>([])

// 格式化日期
const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  try {
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return dateStr
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  } catch {
    return dateStr
  }
}

// 格式化标题（处理换行）
const formatTitle = (title: string) => {
  if (!title || title.trim() === '') {
    return '分享编程心得<br>记录技术成长'
  }
  if (title.includes('\n')) {
    return title.replace(/\n/g, '<br>')
  }
  if (title.length > 10) {
    const mid = Math.floor(title.length / 2)
    let breakPoint = mid
    for (let i = 0; i < 5; i++) {
      if (title[mid + i] === ' ' || title[mid + i] === '，' || title[mid + i] === '。') {
        breakPoint = mid + i + 1
        break
      }
      if (title[mid - i] === ' ' || title[mid - i] === '，' || title[mid - i] === '。') {
        breakPoint = mid - i + 1
        break
      }
    }
    return title.substring(0, breakPoint) + '<br>' + title.substring(breakPoint)
  }
  return title
}

// 格式化浏览量
const formatViews = (views: number) => {
  if (views >= 1000) {
    return (views / 1000).toFixed(1) + 'k'
  }
  return views
}

// 加载文章列表
const loadArticles = async () => {
  try {
    const res = await articleApi.getList({
      page: currentPage.value,
      size: pageSize.value,
      status: 'published'
    })
    articles.value = (res.list || []).map((item: any) => ({
      id: item.id,
      title: item.title,
      desc: item.desc || '',
      category: item.categoryName || '',
      createdAt: formatDate(item.createdAt),
      views: item.views || 0,
      cover: item.coverImage || `https://picsum.photos/id/${item.id}/600/400`,
      publishedByUser: item.publishedByUser || null
    }))
    total.value = res.total || 0
    
    // 设置热门文章（取前4篇）
    hotArticles.value = articles.value.slice(0, 4)
  } catch (error) {
    console.error('加载文章失败:', error)
  }
}

// 加载标签列表
const loadTags = async () => {
  try {
    const res = await tagApi.getList({ page: 1, size: 20 })
    tags.value = (res.list || []).slice(0, 9)
  } catch (error) {
    console.error('加载标签失败:', error)
  }
}

// 分页变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  loadArticles()
}

const scrollToContent = () => {
  if (import.meta.client) {
    document.getElementById('content-start')?.scrollIntoView({ behavior: 'smooth' })
  }
}

const openGithub = () => {
  if (import.meta.client) {
    window.open('https://github.com/gogf/gf', '_blank')
  }
}

// SSR 数据获取
await Promise.all([loadArticles(), loadTags()])

// 客户端加载网站配置
onMounted(() => {
  siteStore.loadConfig()
})
</script>

<style scoped lang="scss">
.home-container {
  padding-bottom: 60px;
}

/* Hero Section */
.hero-section {
  position: relative;
  min-height: 500px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 40px;
  margin-bottom: 60px;
  background-color: #f8f9fa;
  overflow: hidden;
  border-radius: 20px;
}

.hero-bg-shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  z-index: 0;
  opacity: 0.6;
}

.shape-1 {
  width: 400px;
  height: 400px;
  background: #409eff;
  top: -100px;
  left: -100px;
  animation: float 10s infinite ease-in-out;
}

.shape-2 {
  width: 300px;
  height: 300px;
  background: #36d1dc;
  bottom: -50px;
  right: 10%;
  animation: float 12s infinite ease-in-out reverse;
}

.hero-content {
  position: relative;
  z-index: 1;
  max-width: 600px;
  padding-left: 40px;
}

.hero-title {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  font-size: 56px;
  font-weight: 800;
  line-height: 1.2;
  color: #1a1a1a;
  margin-bottom: 24px;
  letter-spacing: -1px;
}

.hero-title .highlight {
  color: #1b1e1b;
}

.hero-subtitle {
  font-size: 20px;
  color: #606266;
  margin-bottom: 40px;
  line-height: 1.6;
}

.hero-actions {
  display: flex;
  gap: 20px;
}

.start-btn {
  padding: 12px 32px;
  font-weight: 600;
  font-size: 16px;
  box-shadow: 0 10px 20px rgba(64,158,255,0.3);
  transition: transform 0.2s;

  &:hover {
    transform: translateY(-2px);
  }
}

.github-btn {
  font-weight: 600;
}

.hero-illustration {
  position: relative;
  z-index: 1;
  width: 45%;
  text-align: center;
  display: flex;
  justify-content: center;
}

/* Code Window */
.code-window {
  width: 100%;
  max-width: 400px;
  background: #1e1e1e;
  border-radius: 12px;
  box-shadow: 0 20px 50px rgba(0,0,0,0.2);
  overflow: hidden;
  font-family: 'Fira Code', 'Consolas', monospace;
  text-align: left;
  transform: perspective(1000px) rotateY(-5deg) rotateX(2deg);
  transition: transform 0.3s ease;

  &:hover {
    transform: perspective(1000px) rotateY(0) rotateX(0);
  }
}

.window-header {
  background: #2d2d2d;
  padding: 12px 16px;
  display: flex;
  gap: 8px;
  border-bottom: 1px solid #333;
}

.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.red { background: #ff5f56; }
.yellow { background: #ffbd2e; }
.green { background: #27c93f; }

.window-body {
  padding: 20px;
  color: #d4d4d4;
  font-size: 14px;
  line-height: 1.5;
}

.code-line {
  white-space: pre;
}

.indent {
  padding-left: 20px;
}

.keyword { color: #569cd6; }
.string { color: #ce9178; }
.comment { color: #6a9955; font-style: italic; }

/* Article List */
.article-item {
  display: flex;
  gap: 24px;
  margin-bottom: 30px;
  padding: 24px;
  background: white;
  border-radius: 16px;
  transition: all 0.3s ease;
  border: 1px solid transparent;
  animation: fade-up 0.6s ease-out forwards;
  opacity: 0;
  transform: translateY(20px);

  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 15px 30px rgba(0,0,0,0.06);
    border-color: rgba(0,0,0,0.03);
  }
}

.article-thumb {
  position: relative;
  width: 260px;
  height: 180px;
  flex-shrink: 0;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  display: block;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.5s;
  }
}

.article-item:hover .article-thumb img {
  transform: scale(1.08);
}

.article-cat-tag {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 10;
  background: rgba(255,255,255,0.95);
  backdrop-filter: blur(4px);
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 700;
  color: #303133;
  box-shadow: 0 2px 8px rgba(0,0,0,0.15);
}

.article-detail {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.article-meta-top {
  display: flex;
  gap: 16px;
  color: #999;
  font-size: 13px;
  margin-bottom: 10px;

  span {
    display: flex;
    align-items: center;
    gap: 4px;
  }
}

.article-title {
  font-size: 22px;
  font-weight: 700;
  color: #303133;
  margin: 0 0 12px 0;
  line-height: 1.4;
  cursor: pointer;
  transition: color 0.2s;

  &:hover {
    color: #409eff;
  }
}

.article-summary {
  color: #606266;
  font-size: 15px;
  line-height: 1.6;
  margin: 0 0 auto 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 16px;
}

.author {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #606266;
  font-weight: 500;
  cursor: pointer;
  transition: opacity 0.2s;
  text-decoration: none;

  &:hover {
    opacity: 0.8;
    color: #409eff;
  }
}

.read-more {
  font-weight: 600;
}

/* Sidebar */
.sidebar-sticky {
  position: sticky;
  top: 80px;
}

.sidebar-widget {
  background: white;
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.03);
}

.widget-title {
  font-size: 18px;
  font-weight: 700;
  margin: 0 0 20px 0;
  padding-left: 10px;
  border-left: 4px solid #409eff;
}

.hot-list {
  padding: 0;
  margin: 0;
  list-style: none;
}

.hot-item {
  margin-bottom: 16px;

  &:last-child {
    margin-bottom: 0;
  }

  a {
    display: flex;
    align-items: center;
    gap: 12px;
    text-decoration: none;
  }

  &:nth-child(1) .hot-number { background: #ffece8; color: #f56c6c; }
  &:nth-child(2) .hot-number { background: #e8f3ff; color: #409eff; }
  &:nth-child(3) .hot-number { background: #fff7e8; color: #e6a23c; }
}

.hot-number {
  width: 24px;
  height: 24px;
  background: #f2f6fc;
  color: #909399;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
}

.hot-info {
  flex: 1;
  overflow: hidden;
}

.hot-text {
  font-size: 14px;
  color: #303133;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.2s;
}

.hot-item:hover .hot-text {
  color: #409eff;
}

.hot-views {
  font-size: 12px;
  color: #c0c4cc;
}

.tags-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.tag-chip {
  padding: 6px 12px;
  background: #f4f4f5;
  color: #606266;
  border-radius: 20px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  text-decoration: none;

  &:hover {
    background: #409eff;
    color: white;
    transform: translateY(-2px);
  }
}

.pagination-wrapper {
  margin-top: 40px;
  display: flex;
  justify-content: center;
}

/* 响应式 */
@media (max-width: 768px) {
  .hero-section {
    flex-direction: column;
    padding: 40px 20px;
    text-align: center;
    min-height: auto;
  }

  .hero-content {
    padding-left: 0;
    margin-bottom: 40px;
  }

  .hero-title {
    font-size: 32px;
  }

  .hero-subtitle {
    font-size: 16px;
  }

  .hero-actions {
    justify-content: center;
  }

  .hero-illustration {
    width: 80%;
  }

  .code-window {
    transform: none;
  }

  .article-item {
    flex-direction: column;
    padding: 16px;
  }

  .article-thumb {
    width: 100%;
    height: 180px;
  }

  .article-detail {
    padding-top: 16px;
  }

  .sidebar-column {
    margin-top: 40px;
  }

  .content-wrapper {
    margin-left: 0 !important;
    margin-right: 0 !important;
  }

  .content-wrapper > .el-col {
    padding-left: 0 !important;
    padding-right: 0 !important;
  }
}
</style>
