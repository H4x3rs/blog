<template>
  <div class="home-container">
    <!-- Hero Banner -->
    <div class="hero-section">
      <div class="hero-bg-shape shape-1"></div>
      <div class="hero-bg-shape shape-2"></div>
      <div class="hero-content">
        <h1 class="hero-title">
          <span class="highlight" v-html="formatTitle(bannerTitle)"></span>
        </h1>
        <p class="hero-subtitle">{{ bannerSubtitle }}</p>
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
            <div class="article-thumb" @click="$router.push(`/article/${item.id}`)">
              <img :src="item.cover" loading="lazy" />
              <div v-if="item.category" class="article-cat-tag">{{ item.category }}</div>
            </div>
            <div class="article-detail">
              <div class="article-meta-top">
                <span class="date"><el-icon><Calendar /></el-icon> {{ item.createdAt }}</span>
                <span class="views"><el-icon><View /></el-icon> {{ item.views }}</span>
              </div>
              <h2 class="article-title" @click="$router.push(`/article/${item.id}`)">{{ item.title }}</h2>
              <p class="article-summary">{{ item.desc }}</p>
              <div class="article-footer">
                <div 
                  class="author" 
                  @click.stop="goToAuthorArticles(item.publishedByUser)"
                >
                  <el-avatar :size="24" :src="item.publishedByUser?.avatar || 'https://picsum.photos/id/64/100/100'" />
                  <span>{{ item.publishedByUser ? (item.publishedByUser.nickname || item.publishedByUser.username) : 'Admin' }}</span>
                </div>
                <el-button link type="primary" class="read-more" @click="$router.push(`/article/${item.id}`)">阅读全文</el-button>
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
            @current-change="(page) => { currentPage = page; loadArticles() }"
            class="pagination-responsive" 
          />
        </div>
      </el-col>

      <!-- 右侧侧边栏 -->
      <el-col :span="8" :xs="24" class="sidebar-column">
        <div class="sidebar-sticky">
          <!-- 个人简介 - 仅登录后显示 -->
          <div v-if="isLoggedIn" class="sidebar-widget profile-widget">
            <div class="profile-bg"></div>
            <div class="profile-content">
              <div class="profile-avatar" @click="$router.push('/admin/profile')">
                <img :src="userAvatar || 'https://picsum.photos/id/64/150/150'" :alt="displayName" />
              </div>
              <h3>{{ displayName }}</h3>
              <p class="job-title">{{ userInfo.username || 'Blog User' }}</p>
              <div class="social-icons">
                 <a href="https://github.com/gogf/gf" target="_blank" class="social-btn"><i class="fab fa-github"></i></a>
                 <a href="https://twitter.com/golang" target="_blank" class="social-btn"><i class="fab fa-twitter"></i></a>
              </div>
              <div class="stats-row">
                <div class="stat-item clickable" @click="goToMyArticles">
                  <strong>{{ myArticleCount }}</strong>
                  <span>Articles</span>
                </div>
                <div class="stat-item">
                  <strong>0</strong>
                  <span>Likes</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 热门文章 -->
          <div class="sidebar-widget hot-widget">
            <h4 class="widget-title">热门文章</h4>
            <ul class="hot-list">
              <li v-for="(item, i) in hotArticles" :key="i" class="hot-item" @click="$router.push(`/article/1`)">
                <span class="hot-number">{{ i + 1 }}</span>
                <div class="hot-info">
                   <div class="hot-text">{{ item.title }}</div>
                   <span class="hot-views">{{ item.views }} reads</span>
                </div>
              </li>
            </ul>
          </div>

          <!-- 标签云 -->
          <div class="sidebar-widget tag-widget">
             <h4 class="widget-title">标签</h4>
             <div class="tags-cloud">
               <span v-for="tag in tags" :key="tag" class="tag-chip" @click="$router.push(`/tag`)">{{ tag }}</span>
             </div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { Calendar, View, ArrowRight } from '@element-plus/icons-vue'
import { getArticleList } from '@/api/article'
import { getTagList } from '@/api/tag'
import { getCurrentUser } from '@/api/user'
import { useSiteConfig } from '@/store/site'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const tags = ref([])
const articles = ref([])
const hotArticles = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const myArticleCount = ref(0)

// 用户信息
const userInfo = ref({
  username: '',
  nickname: '',
  avatar: ''
})

// 是否已登录
const isLoggedIn = computed(() => {
  const token = localStorage.getItem('token')
  return !!token && !!userInfo.value.username
})

// 显示名称（优先显示nickname，如果没有则显示username）
const displayName = computed(() => {
  return userInfo.value.nickname || userInfo.value.username || '用户'
})

// 用户头像
const userAvatar = computed(() => {
  return userInfo.value.avatar || ''
})

// 加载用户信息
const loadUserInfo = async () => {
  // 先从localStorage读取
  const username = localStorage.getItem('username')
  const nickname = localStorage.getItem('nickname')
  const avatar = localStorage.getItem('avatar')
  
  // 立即显示localStorage中的信息（快速显示）
  if (username) {
    userInfo.value.username = username
    userInfo.value.nickname = nickname || ''
    userInfo.value.avatar = avatar || ''
  }
  
  // 如果有token，尝试从后端获取完整用户信息
  const token = localStorage.getItem('token')
  if (token) {
    try {
      const res = await getCurrentUser()
      if (res) {
        // 处理响应数据（可能是嵌套的User对象）
        const userData = res.user || res
        userInfo.value.username = userData.username || username || ''
        userInfo.value.nickname = userData.nickname || nickname || ''
        userInfo.value.avatar = userData.avatar || avatar || ''
        
        // 更新localStorage
        if (userData.username) {
          localStorage.setItem('username', userData.username)
        }
        if (userData.nickname) {
          localStorage.setItem('nickname', userData.nickname)
        }
        if (userData.avatar) {
          localStorage.setItem('avatar', userData.avatar)
        }
      }
    } catch (error) {
      // 如果获取失败，使用localStorage中的信息
      console.warn('获取用户信息失败:', error)
      // 确保至少显示localStorage中的信息
      if (!userInfo.value.username && username) {
        userInfo.value.username = username
        userInfo.value.nickname = nickname || ''
        userInfo.value.avatar = avatar || ''
      }
    }
  } else {
    // 没有token，清空用户信息
    userInfo.value.username = ''
    userInfo.value.nickname = ''
    userInfo.value.avatar = ''
  }
}

// 使用全局网站配置
const { bannerTitle, bannerSubtitle } = useSiteConfig()

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

// 格式化标题（处理换行）
const formatTitle = (title) => {
  if (!title || title.trim() === '') {
    return '分享编程心得<br>记录技术成长'
  }
  // 如果包含换行符，直接返回；否则尝试按长度分割
  if (title.includes('\n')) {
    return title.replace(/\n/g, '<br>')
  }
  // 如果标题较长，尝试在中间位置换行
  if (title.length > 10) {
    const mid = Math.floor(title.length / 2)
    // 查找最近的空格或标点符号
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

// Banner 数据已在全局 store 中初始化，这里不需要单独加载

// 加载文章列表
const loadArticles = async () => {
  try {
    const res = await getArticleList({
      page: currentPage.value,
      size: pageSize.value,
      status: 'published' // 只获取已发布的文章
    })
    articles.value = (res.list || []).map(item => ({
      id: item.id,
      title: item.title,
      desc: item.desc || '',
      category: item.categoryName || '',
      createdAt: formatDate(item.createdAt),
      views: item.views || 0,
      cover: item.coverImage || 'https://picsum.photos/id/' + item.id + '/600/400',
      publishedByUser: item.publishedByUser || null
    }))
    total.value = res.total || 0
    
    // 设置热门文章（取前4篇）
    hotArticles.value = articles.value.slice(0, 4).map(item => ({
      title: item.title,
      views: item.views >= 1000 ? (item.views / 1000).toFixed(1) + 'k' : item.views
    }))
  } catch (error) {
    console.error('加载文章失败:', error)
    ElMessage.error('加载文章失败')
  }
}

// 加载标签列表
const loadTags = async () => {
  try {
    const res = await getTagList({
      page: 1,
      size: 20
    })
    tags.value = (res.list || []).slice(0, 9).map(item => item.name)
  } catch (error) {
    console.error('加载标签失败:', error)
  }
}

// 加载我的文章数量
const loadMyArticleCount = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    myArticleCount.value = 0
    return
  }
  
  try {
    const res = await getArticleList({
      page: 1,
      size: 1,
      status: '', // 不传status，获取所有状态的文章
      onlyMine: true // 只查询当前用户的文章
    })
    myArticleCount.value = res.total || 0
  } catch (error) {
    console.error('加载我的文章数量失败:', error)
    myArticleCount.value = 0
  }
}

// 跳转到我的文章页面
const goToMyArticles = () => {
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.warning('请先登录')
    router.push('/login-new')
    return
  }
  router.push('/my-articles')
}

onMounted(() => {
  loadArticles()
  loadTags()
  loadUserInfo()
  loadMyArticleCount()
})

// 监听路由变化，重新加载用户信息
watch(() => route.path, () => {
  if (route.path === '/') {
    loadUserInfo()
    loadMyArticleCount()
  }
})

// 监听登录状态变化，更新文章数量
watch(isLoggedIn, (newVal) => {
  if (newVal) {
    loadMyArticleCount()
  } else {
    myArticleCount.value = 0
  }
})

const scrollToContent = () => {
  document.getElementById('content-start').scrollIntoView({ behavior: 'smooth' })
}

const openGithub = () => {
  window.open('https://github.com/gogf/gf', '_blank')
}

// 跳转到作者文章列表
const goToAuthorArticles = (user) => {
  if (user && user.id) {
    router.push(`/author/${user.id}`)
  } else {
    // 如果没有作者信息，跳转到关于页面
    router.push('/about')
  }
}
</script>

<style scoped>
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
  background-color: #f8f9fa; /* Fallback */
  overflow: hidden;
  border-radius: 20px;
}

/* 动态背景图形 */
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

@keyframes float {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(30px, 40px); }
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
}
.start-btn:hover {
  transform: translateY(-2px);
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

/* CSS Code Window */
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
}

.code-window:hover {
  transform: perspective(1000px) rotateY(0) rotateX(0);
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

/* Syntax Highlighting */
.keyword { color: #569cd6; }
.string { color: #ce9178; }
.comment { color: #6a9955; font-style: italic; }

@media (max-width: 768px) {
  .hero-illustration {
    width: 100%;
    margin-top: 30px;
  }
  .code-window {
    transform: none;
  }
}

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
}

@keyframes fade-up {
  to { opacity: 1; transform: translateY(0); }
}

.article-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 30px rgba(0,0,0,0.06);
  border-color: rgba(0,0,0,0.03);
}

.article-thumb {
  position: relative;
  width: 260px;
  height: 180px;
  flex-shrink: 0;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
}

.article-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
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
  white-space: nowrap;
  max-width: calc(100% - 24px);
  overflow: hidden;
  text-overflow: ellipsis;
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
}
.article-meta-top span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.article-title {
  font-size: 22px;
  font-weight: 700;
  color: #303133;
  margin: 0 0 12px 0;
  line-height: 1.4;
  cursor: pointer;
  transition: color 0.2s;
}
.article-title:hover {
  color: #409eff;
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
}

.author:hover {
  opacity: 0.8;
  color: #409eff;
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

.profile-widget {
  padding: 0;
  overflow: hidden;
  text-align: center;
}
.profile-bg {
  height: 100px;
  background: linear-gradient(135deg, #a1c4fd 0%, #c2e9fb 100%);
}
.profile-content {
  padding: 0 24px 24px;
  position: relative;
}
.profile-avatar {
  width: 80px;
  height: 80px;
  margin: -40px auto 16px;
  border: 4px solid white;
  border-radius: 50%;
  overflow: hidden;
  box-shadow: 0 4px 10px rgba(0,0,0,0.1);
  cursor: pointer;
  transition: transform 0.3s;
}

.profile-avatar:hover {
  transform: rotate(360deg);
}
.profile-avatar img {
  width: 100%;
  height: 100%;
}
.job-title {
  color: #909399;
  font-size: 14px;
  margin-bottom: 16px;
}
.stats-row {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #f0f2f5;
}
.stat-item {
  display: flex;
  flex-direction: column;
}
.stat-item strong {
  font-size: 18px;
  color: #303133;
}
.stat-item span {
  font-size: 12px;
  color: #909399;
}

.stat-item.clickable {
  cursor: pointer;
  transition: all 0.2s;
}

.stat-item.clickable:hover {
  transform: translateY(-2px);
}

.stat-item.clickable:hover strong {
  color: #409eff;
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
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  cursor: pointer;
}
.hot-item:last-child {
  margin-bottom: 0;
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
.hot-item:nth-child(1) .hot-number { background: #ffece8; color: #f56c6c; }
.hot-item:nth-child(2) .hot-number { background: #e8f3ff; color: #409eff; }
.hot-item:nth-child(3) .hot-number { background: #fff7e8; color: #e6a23c; }

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
}
.tag-chip:hover {
  background: #409eff;
  color: white;
  transform: translateY(-2px);
}

.pagination-wrapper {
  margin-top: 40px;
  display: flex;
  justify-content: center;
}

/* 响应式调整 */
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
  /* 修复 el-row gutter 在移动端的溢出问题 */
  .content-wrapper {
    margin-left: 0 !important;
    margin-right: 0 !important;
  }
  .content-wrapper > .el-col {
    padding-left: 0 !important;
    padding-right: 0 !important;
  }
  
  .pagination-responsive :deep(.el-pagination__jump) {
    display: none;
  }
}
</style>
