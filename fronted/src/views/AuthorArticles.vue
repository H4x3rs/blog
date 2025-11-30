<template>
  <div class="author-articles-container" v-loading="loading">
    <!-- Header Area -->
    <div class="author-header-section">
      <div class="header-bg"></div>
      <div class="header-content">
        <div class="author-avatar-wrapper">
          <el-avatar :size="100" :src="authorInfo.avatar || 'https://picsum.photos/id/64/150/150'" />
        </div>
        <h1 class="page-title">{{ authorInfo.name || '作者' }}</h1>
        <p class="author-bio">{{ authorInfo.bio || '暂无简介' }}</p>
        <div class="header-stats">
          <div class="stat-item">
            <span class="stat-value">{{ total }}</span>
            <span class="stat-label">文章数</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Content Area -->
    <el-row justify="center" class="content-wrapper">
      <el-col :span="18" :xs="22" :sm="20" :md="18" :lg="16">
        <!-- 文章列表 -->
        <el-card class="articles-card" shadow="never">
          <div class="articles-header">
            <h3 class="section-title">
              <el-icon><List /></el-icon>
              全部文章
              <span class="count">({{ articles.length }})</span>
            </h3>
            <el-select v-model="sortBy" style="width: 140px" size="small">
              <el-option label="最新发布" value="date" />
              <el-option label="最多阅读" value="views" />
            </el-select>
          </div>

          <el-empty v-if="articles.length === 0 && !loading" description="该作者暂无文章" />

          <div class="articles-list" v-else>
            <div 
              v-for="article in sortedArticles" 
              :key="article.id"
              class="article-item"
              @click="viewArticle(article)"
            >
              <div class="article-cover" v-if="article.cover">
                <img :src="article.cover" :alt="article.title" />
              </div>
              <div class="article-content">
                <div class="article-top">
                  <h4 class="article-title">{{ article.title }}</h4>
                </div>
                <p class="article-summary">{{ article.summary }}</p>
                <div class="article-meta">
                  <span class="meta-item">
                    <el-icon><Calendar /></el-icon>
                    {{ article.createdAt }}
                  </span>
                  <span class="meta-item">
                    <el-icon><View /></el-icon>
                    {{ formatNumber(article.views) }}
                  </span>
                  <span class="meta-item" v-if="article.category">
                    <el-icon><Folder /></el-icon>
                    {{ article.category }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- 分页 -->
          <div class="pagination-wrapper" v-if="total > pageSize">
            <el-pagination 
              background 
              layout="prev, pager, next" 
              :total="total" 
              :current-page="currentPage"
              :page-size="pageSize"
              @current-change="handlePageChange"
            />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  View, 
  List, 
  Calendar,
  Folder
} from '@element-plus/icons-vue'
import { getArticleList } from '@/api/article'
import { getUser } from '@/api/user'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const authorId = computed(() => parseInt(route.params.id) || 0)

const loading = ref(true)
const articles = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(12)
const sortBy = ref('date')

const authorInfo = ref({
  id: 0,
  name: '',
  avatar: '',
  bio: ''
})

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

// 格式化数字
const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'W'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num
}

// 加载作者信息
const loadAuthorInfo = async () => {
  if (!authorId.value) {
    ElMessage.error('作者ID无效')
    return
  }
  
  try {
    const res = await getUser({ id: authorId.value })
    if (res) {
      authorInfo.value = {
        id: res.id,
        name: res.nickname || res.username || '作者',
        avatar: res.avatar || '',
        bio: res.bio || ''
      }
    }
  } catch (error) {
    console.error('加载作者信息失败:', error)
  }
}

// 加载文章列表
const loadArticles = async () => {
  if (!authorId.value) {
    return
  }
  
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value,
      status: 'published', // 只显示已发布的文章
      authorId: authorId.value // 根据作者ID查询
    }
    
    const res = await getArticleList(params)
    articles.value = (res.list || []).map(item => ({
      id: item.id,
      title: item.title,
      summary: item.desc || '',
      cover: item.coverImage || 'https://picsum.photos/id/' + item.id + '/400/300',
      createdAt: formatDate(item.createdAt),
      views: item.views || 0,
      category: item.categoryName || ''
    }))
    total.value = res.total || 0
  } catch (error) {
    console.error('加载文章失败:', error)
    ElMessage.error('加载文章失败')
  } finally {
    loading.value = false
  }
}

// 排序后的文章列表
const sortedArticles = computed(() => {
  const sorted = [...articles.value]
  
  if (sortBy.value === 'date') {
    sorted.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
  } else if (sortBy.value === 'views') {
    sorted.sort((a, b) => b.views - a.views)
  }
  
  return sorted
})

// 分页变化
const handlePageChange = (page) => {
  currentPage.value = page
  loadArticles()
}

// 查看文章
const viewArticle = (article) => {
  router.push(`/article/${article.id}`)
}

onMounted(async () => {
  if (!authorId.value) {
    ElMessage.error('作者ID无效')
    router.push('/')
    return
  }
  
  await loadAuthorInfo()
  await loadArticles()
})

// 监听路由变化
watch(() => route.params.id, async (newId) => {
  if (newId) {
    await loadAuthorInfo()
    await loadArticles()
  }
})
</script>

<style scoped>
.author-articles-container {
  background-color: #f8f9fa;
  min-height: 100vh;
  padding-bottom: 60px;
}

/* Header Section */
.author-header-section {
  position: relative;
  min-height: 450px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  padding: 20px 0;
}

.header-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: 
    radial-gradient(circle at 20% 50%, rgba(255, 255, 255, 0.1) 0%, transparent 50%),
    radial-gradient(circle at 80% 80%, rgba(255, 255, 255, 0.1) 0%, transparent 50%);
  z-index: 0;
  pointer-events: none;
}

.header-content {
  position: relative;
  z-index: 2;
  text-align: center;
  max-width: 900px;
  padding: 40px 20px;
}

.author-avatar-wrapper {
  margin-bottom: 24px;
  display: flex;
  justify-content: center;
}

.author-avatar-wrapper :deep(.el-avatar) {
  border: 4px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.page-title {
  font-size: 42px;
  font-weight: 800;
  line-height: 1.3;
  margin-bottom: 16px;
  text-shadow: 0 2px 20px rgba(0,0,0,0.3);
}

.author-bio {
  font-size: 16px;
  line-height: 1.5;
  margin-bottom: 30px;
  opacity: 0.95;
  text-shadow: 0 1px 4px rgba(0,0,0,0.2);
}

.header-stats {
  display: flex;
  justify-content: center;
  gap: 30px;
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  line-height: 1;
  margin-bottom: 8px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
}

/* Content Area */
.content-wrapper {
  margin-top: -40px;
  position: relative;
  z-index: 2;
}

.articles-card {
  border-radius: 16px;
  border: none;
  box-shadow: 0 2px 12px rgba(0,0,0,0.08);
}

:deep(.el-card__body) {
  padding: 32px;
}

.articles-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.section-title .el-icon {
  font-size: 22px;
  color: #409eff;
}

.count {
  font-size: 16px;
  font-weight: 500;
  color: #909399;
}

/* Articles List */
.articles-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.article-item {
  display: flex;
  gap: 20px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.article-item:hover {
  background: #e9ecef;
  transform: translateX(4px);
}

.article-cover {
  flex-shrink: 0;
  width: 200px;
  height: 140px;
  border-radius: 8px;
  overflow: hidden;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.article-item:hover .article-cover img {
  transform: scale(1.05);
}

.article-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.article-top {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 8px;
}

.article-title {
  flex: 1;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin: 0;
  line-height: 1.4;
}

.article-summary {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex: 1;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 20px;
  font-size: 13px;
  color: #909399;
  flex-wrap: wrap;
  margin-top: auto;
}

.article-meta .meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.pagination-wrapper {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}

/* 响应式 */
@media (max-width: 768px) {
  .author-header-section {
    padding: 40px 20px;
  }
  
  .author-info-header {
    flex-direction: column;
    text-align: center;
  }
  
  .page-title {
    font-size: 28px;
  }
  
  :deep(.el-card__body) {
    padding: 20px;
  }
  
  .article-item {
    flex-direction: column;
  }
  
  .article-cover {
    width: 100%;
    height: 200px;
  }
  
  .articles-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
}
</style>

