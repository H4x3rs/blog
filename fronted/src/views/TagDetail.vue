<template>
  <div class="tag-detail-container" v-loading="loading">
    <!-- Header Area -->
    <div class="tag-header">
      <div class="header-content">
        <div class="tag-icon-large">
          <el-icon :size="60">
            <PriceTag />
          </el-icon>
        </div>
        <h1 class="title"># {{ tag.name }}</h1>
        <p class="description">{{ tag.description }}</p>
        <div class="tag-meta">
          <span class="meta-item">
            <el-icon><Document /></el-icon>
            {{ articles.length }} 篇文章
          </span>
          <span class="divider">·</span>
          <span class="meta-item">
            <el-icon><View /></el-icon>
            {{ formatNumber(totalViews) }} 阅读
          </span>
          <span class="divider">·</span>
          <span class="meta-item">
            <el-icon><User /></el-icon>
            {{ tag.followers }} 关注
          </span>
        </div>
        <el-button type="primary" round size="large" class="follow-btn">
          <el-icon><Plus /></el-icon>
          关注标签
        </el-button>
      </div>
      <div class="header-bg"></div>
    </div>

    <!-- Content Area -->
    <el-row justify="center" class="content-wrapper">
      <el-col :span="16" :xs="22">
        <!-- 文章列表 -->
        <el-card class="articles-card" shadow="never">
          <div class="articles-header">
            <h3 class="section-title">
              <el-icon><List /></el-icon>
              全部文章
              <span class="count">({{ sortedArticles.length }})</span>
            </h3>
            <el-select v-model="sortBy" style="width: 140px" size="small">
              <el-option label="最新发布" value="date" />
              <el-option label="最多阅读" value="views" />
              <el-option label="最多评论" value="comments" />
            </el-select>
          </div>

          <div class="articles-list">
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
                  <div class="article-badges">
                    <el-tag v-if="article.isNew" type="success" size="small" effect="plain">新</el-tag>
                    <el-tag v-if="article.isHot" type="danger" size="small" effect="plain">热</el-tag>
                  </div>
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
                  <span class="meta-item">
                    <el-icon><ChatDotRound /></el-icon>
                    {{ article.comments }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </el-card>

        <!-- 相关标签 -->
        <el-card class="related-card" shadow="never" v-if="relatedTags.length > 0">
          <h3 class="section-title">
            <el-icon><CollectionTag /></el-icon>
            相关标签
          </h3>
          <div class="related-list">
            <div 
              v-for="item in relatedTags" 
              :key="item.id"
              class="related-item"
              @click="viewTag(item)"
            >
              <div class="related-content">
                <div class="related-name"># {{ item.name }}</div>
                <div class="related-desc">{{ item.description }}</div>
              </div>
              <div class="related-stats">
                <div class="stat-item">
                  <span class="stat-value">{{ item.count }}</span>
                  <span class="stat-label">文章</span>
                </div>
                <div class="stat-divider"></div>
                <div class="stat-item">
                  <span class="stat-value">{{ formatNumber(item.views) }}</span>
                  <span class="stat-label">阅读</span>
                </div>
              </div>
            </div>
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
  Document, 
  View, 
  User,
  Plus,
  List, 
  Calendar, 
  ChatDotRound,
  CollectionTag,
  PriceTag
} from '@element-plus/icons-vue'
import { getTagBySlug, getTagList } from '@/api/tag'
import { getArticleList } from '@/api/article'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const tagSlug = computed(() => route.params.slug)

const sortBy = ref('date')
const loading = ref(true)

const tag = ref({
  id: 0,
  slug: tagSlug.value,
  name: '',
  description: '',
  followers: 0
})

const articles = ref([])
const relatedTags = ref([])

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

// 加载标签信息
const loadTag = async () => {
  try {
    const res = await getTagBySlug({ slug: tagSlug.value })
    if (res) {
      const tagData = res
      tag.value = {
        id: tagData.id,
        slug: tagData.slug,
        name: tagData.name,
        description: tagData.description || tagData.name + '相关技术文章',
        followers: 0 // 关注功能暂未实现
      }
      // 注意：后端可能不支持通过标签查询文章，这里暂时设为空
      articles.value = []
    } else {
      ElMessage.error('标签不存在')
    }
  } catch (error) {
    console.error('加载标签失败:', error)
    ElMessage.error('加载标签失败')
  }
}

// 加载相关标签
const loadRelatedTags = async () => {
  try {
    const res = await getTagList({
      page: 1,
      size: 10
    })
    relatedTags.value = (res.list || [])
      .filter(t => t.id !== tag.value.id)
      .slice(0, 4)
      .map(item => ({
        id: item.id,
        slug: item.slug,
        name: item.name,
        description: item.description || item.name + '相关技术文章',
        count: 0,
        views: 0
      }))
  } catch (error) {
    console.error('加载相关标签失败:', error)
  }
}

onMounted(async () => {
  loading.value = true
  await loadTag()
  await loadRelatedTags()
  loading.value = false
})

watch(() => route.params.slug, async (newSlug) => {
  if (newSlug) {
    loading.value = true
    await loadTag()
    await loadRelatedTags()
    loading.value = false
  }
})

const totalViews = computed(() => {
  return articles.value.reduce((sum, article) => sum + article.views, 0)
})

const sortedArticles = computed(() => {
  const sorted = [...articles.value]
  
  if (sortBy.value === 'date') {
    sorted.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
  } else if (sortBy.value === 'views') {
    sorted.sort((a, b) => b.views - a.views)
  } else if (sortBy.value === 'comments') {
    sorted.sort((a, b) => b.comments - a.comments)
  }
  
  return sorted
})

const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'W'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num
}

const viewArticle = (article) => {
  router.push(`/article/${article.id}`)
}

const viewTag = (tag) => {
  router.push(`/tag/${tag.slug}`)
}
</script>

<style scoped>
.tag-detail-container {
  background-color: #f8f9fa;
  min-height: 100vh;
  padding-bottom: 60px;
}

/* Header Styles */
.tag-header {
  position: relative;
  min-height: 480px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: visible;
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

.tag-icon-large {
  width: 100px;
  height: 100px;
  margin: 0 auto 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.title {
  font-size: 42px;
  font-weight: 800;
  line-height: 1.3;
  margin-bottom: 16px;
  text-shadow: 0 2px 20px rgba(0,0,0,0.3);
}

.description {
  font-size: 16px;
  line-height: 1.5;
  margin-bottom: 20px;
  opacity: 0.95;
  text-shadow: 0 1px 4px rgba(0,0,0,0.2);
}

.tag-meta {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 14px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.divider {
  opacity: 0.6;
}

.follow-btn {
  margin-top: 8px;
  padding: 12px 32px;
  font-size: 15px;
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

/* Content Area */
.content-wrapper {
  margin-top: -60px;
  position: relative;
  z-index: 2;
}

/* Cards */
.articles-card,
.related-card {
  margin-bottom: 24px;
  border-radius: 16px;
  border: none;
  box-shadow: 0 2px 12px rgba(0,0,0,0.08);
}

:deep(.el-card__body) {
  padding: 32px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin-bottom: 24px;
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
.articles-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

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

.article-badges {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
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

/* Related Tags */
.related-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.related-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  background: #f8f9fa;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.related-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  background: #e9ecef;
}

.related-content {
  flex: 1;
  min-width: 0;
}

.related-name {
  font-size: 18px;
  font-weight: 600;
  color: #409eff;
  margin-bottom: 6px;
}

.related-desc {
  font-size: 13px;
  color: #606266;
  line-height: 1.5;
}

.related-stats {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 18px;
  font-weight: 700;
  color: #303133;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #909399;
}

.stat-divider {
  width: 1px;
  height: 30px;
  background: #dcdfe6;
}

/* 响应式 */
@media (max-width: 768px) {
  .tag-header {
    min-height: 400px;
    padding: 20px 0;
  }
  
  .header-content {
    padding: 30px 20px;
  }
  
  .tag-icon-large {
    width: 80px;
    height: 80px;
    margin-bottom: 16px;
  }
  
  .tag-icon-large .el-icon {
    font-size: 48px;
  }
  
  .title {
    font-size: 28px;
  }
  
  .description {
    font-size: 14px;
  }
  
  .follow-btn {
    font-size: 14px;
    padding: 10px 24px;
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
  
  .related-list {
    grid-template-columns: 1fr;
  }
  
  .related-item {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .related-stats {
    width: 100%;
    justify-content: center;
  }
}
</style>

