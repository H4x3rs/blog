<template>
  <div class="topic-detail-container" v-loading="loading">
    <!-- Header Area -->
    <div class="topic-header">
      <div class="header-content">
        <div class="meta-info">
          <el-tag 
            v-for="tag in topic.tags" 
            :key="tag" 
            effect="dark" 
            round 
            class="topic-tag"
          >
            {{ tag }}
          </el-tag>
          <span class="status-tag" :class="topic.status === '连载中' ? 'status-updating' : 'status-finished'">
            {{ topic.status }}
          </span>
        </div>
        <h1 class="title">{{ topic.name }}</h1>
        <p class="description">{{ topic.description }}</p>
        <div class="topic-meta">
          <span class="meta-item">
            <el-icon><Document /></el-icon>
            {{ topic.articleCount }} 篇文章
          </span>
          <span class="divider">·</span>
          <span class="meta-item">
            <el-icon><View /></el-icon>
            {{ formatNumber(topic.views) }} 阅读
          </span>
          <span class="divider">·</span>
          <span class="meta-item">
            <el-icon><Clock /></el-icon>
            {{ topic.updateFrequency }}
          </span>
        </div>
        <div class="author-info">
          <el-avatar :size="32" :src="topic.author.avatar" />
          <span class="author-name">{{ topic.author.name }}</span>
          <span class="author-role">专题作者</span>
        </div>
      </div>
      <div class="header-image" :style="{ backgroundImage: `url(${topic.cover})` }"></div>
    </div>

    <!-- Content Area -->
    <el-row justify="center" class="content-wrapper">
      <el-col :span="18" :xs="22" :sm="20" :md="18" :lg="16">
        <!-- 专题简介 -->
        <el-card class="intro-card" shadow="never" v-if="topic.intro">
          <h3 class="section-title">
            <el-icon><Reading /></el-icon>
            专题介绍
          </h3>
          <div class="markdown-body" v-html="renderedIntro"></div>
        </el-card>

        <!-- 文章列表 -->
        <el-card class="articles-card" shadow="never">
          <div class="articles-header">
            <h3 class="section-title">
              <el-icon><List /></el-icon>
              文章列表
              <span class="count">({{ articles.length }})</span>
            </h3>
            <el-select v-model="sortBy" style="width: 140px" size="small">
              <el-option label="默认排序" value="default" />
              <el-option label="最新发布" value="date" />
              <el-option label="最多阅读" value="views" />
            </el-select>
          </div>

          <div class="articles-list">
            <div 
              v-for="(article, index) in sortedArticles" 
              :key="article.id"
              class="article-item"
              @click="viewArticle(article)"
            >
              <div class="article-number">{{ index + 1 }}</div>
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
              <div class="article-cover" v-if="article.cover">
                <img :src="article.cover" :alt="article.title" />
              </div>
            </div>
          </div>
        </el-card>

        <!-- 相关专题 -->
        <el-card class="related-card" shadow="never" v-if="relatedTopics.length > 0">
          <h3 class="section-title">
            <el-icon><CollectionTag /></el-icon>
            相关专题
          </h3>
          <div class="related-list">
            <div 
              v-for="item in relatedTopics" 
              :key="item.id"
              class="related-item"
              @click="viewTopic(item)"
            >
              <div class="related-cover">
                <img :src="item.cover" alt="" />
              </div>
              <div class="related-info">
                <h4 class="related-title">{{ item.name }}</h4>
                <p class="related-desc">{{ item.description }}</p>
                <div class="related-meta">
                  {{ item.articleCount }} 篇 · {{ formatNumber(item.views) }} 阅读
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
  Clock, 
  Reading, 
  List, 
  Calendar, 
  ChatDotRound,
  CollectionTag 
} from '@element-plus/icons-vue'
import { marked } from 'marked'
import { getTopic, getTopicList } from '@/api/topic'
import { getTopicArticles } from '@/api/topic'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const topicId = computed(() => route.params.id)

const sortBy = ref('default')
const loading = ref(true)

const topic = ref({
  id: topicId.value,
  name: '',
  description: '',
  cover: '',
  tags: [],
  articleCount: 0,
  views: 0,
  status: '连载中',
  updateFrequency: '每周更新',
  author: {
    name: 'Admin',
    avatar: 'https://i.pravatar.cc/150?img=1'
  },
  intro: ''
})

const articles = ref([])
const relatedTopics = ref([])

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

// 加载专题信息
const loadTopic = async () => {
  try {
    const res = await getTopic({ id: parseInt(topicId.value) })
    if (res) {
      topic.value = {
        id: res.id,
        name: res.name || '',
        description: res.description || '',
        cover: res.coverImage || 'https://picsum.photos/id/' + res.id + '/1200/600',
        tags: [],
        articleCount: 0,
        views: 0,
        status: '连载中',
        updateFrequency: '每周更新',
        author: {
          name: 'Admin',
          avatar: 'https://i.pravatar.cc/150?img=1'
        },
        intro: res.intro || ''
      }
      await loadArticles()
    }
  } catch (error) {
    console.error('加载专题失败:', error)
    ElMessage.error('加载专题失败')
  }
}

// 加载专题文章列表
const loadArticles = async () => {
  try {
    const res = await getTopicArticles({ id: parseInt(topicId.value), page: 1, size: 100 })
    if (res && res.list) {
      articles.value = res.list.map((item, index) => {
        const article = item.article || item
        return {
          id: article.id,
          title: article.title || '',
          summary: article.desc || '',
          cover: article.coverImage || 'https://picsum.photos/id/' + article.id + '/400/300',
          createdAt: formatDate(article.createdAt),
          views: article.views || 0,
          comments: 0,
          isNew: index < 3,
          isHot: false
        }
      })
      // 使用API返回的total而不是list.length，因为可能分页了
      topic.value.articleCount = res.total || articles.value.length
    }
  } catch (error) {
    console.error('加载专题文章失败:', error)
  }
}

// 加载相关专题
const loadRelatedTopics = async () => {
  try {
    const res = await getTopicList({
      page: 1,
      size: 10
    })
    relatedTopics.value = (res.list || [])
      .filter(t => t.id !== topic.value.id)
      .slice(0, 2)
      .map(item => ({
        id: item.id,
        name: item.name,
        description: item.description || '',
        cover: item.coverImage || 'https://picsum.photos/id/' + item.id + '/300/200',
        articleCount: 0,
        views: 0
      }))
  } catch (error) {
    console.error('加载相关专题失败:', error)
  }
}

onMounted(async () => {
  loading.value = true
  await loadTopic()
  await loadRelatedTopics()
  loading.value = false
})

watch(() => route.params.id, async (newId) => {
  if (newId) {
    loading.value = true
    await loadTopic()
    await loadRelatedTopics()
    loading.value = false
  }
})

const renderedIntro = computed(() => {
  if (!topic.value.intro) return ''
  return marked.parse(topic.value.intro)
})

const sortedArticles = computed(() => {
  const sorted = [...articles.value]
  
  if (sortBy.value === 'date') {
    sorted.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
  } else if (sortBy.value === 'views') {
    sorted.sort((a, b) => b.views - a.views)
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

const viewTopic = (topic) => {
  router.push(`/topic/${topic.id}`)
}
</script>

<style scoped>
.topic-detail-container {
  background-color: #f8f9fa;
  min-height: 100vh;
  padding-bottom: 60px;
}

/* Header Styles */
.topic-header {
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
  opacity: 0.3;
  filter: blur(8px);
  transform: scale(1.1);
}

.header-content {
  position: relative;
  z-index: 1;
  text-align: center;
  max-width: 900px;
}

.meta-info {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  flex-wrap: wrap;
}

.topic-tag {
  background-color: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
  font-weight: 600;
}

.status-tag {
  padding: 4px 16px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
}

.status-updating {
  background-color: rgba(103, 194, 58, 0.9);
}

.status-finished {
  background-color: rgba(144, 147, 153, 0.9);
}

.title {
  font-size: 32px;
  font-weight: 800;
  line-height: 1.3;
  margin-bottom: 12px;
  text-shadow: 0 2px 20px rgba(0,0,0,0.3);
}

.description {
  font-size: 15px;
  line-height: 1.5;
  margin-bottom: 16px;
  opacity: 0.95;
  text-shadow: 0 1px 4px rgba(0,0,0,0.2);
}

.topic-meta {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 13px;
  margin-bottom: 14px;
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

.author-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding-top: 10px;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
}

.author-name {
  font-weight: 600;
  font-size: 14px;
}

.author-role {
  font-size: 12px;
  opacity: 0.8;
}

/* Content Area */
.content-wrapper {
  margin-top: -60px;
  position: relative;
  z-index: 2;
}

/* Cards */
.intro-card,
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

/* Markdown Styles */
.markdown-body {
  font-size: 16px;
  line-height: 1.8;
  color: #333;
}

:deep(.markdown-body h2) {
  margin-top: 24px;
  margin-bottom: 12px;
  font-size: 22px;
  font-weight: 700;
  color: #1a1a1a;
  border-bottom: 1px solid #eee;
  padding-bottom: 8px;
}

:deep(.markdown-body h3) {
  margin-top: 20px;
  margin-bottom: 10px;
  font-size: 18px;
  font-weight: 700;
  color: #1a1a1a;
}

:deep(.markdown-body p) {
  margin-bottom: 16px;
}

:deep(.markdown-body ul) {
  padding-left: 24px;
  margin-bottom: 16px;
}

:deep(.markdown-body li) {
  margin-bottom: 8px;
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

.article-number {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 8px;
  font-weight: 700;
  font-size: 16px;
}

.article-content {
  flex: 1;
  min-width: 0;
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
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 20px;
  font-size: 13px;
  color: #909399;
}

.article-meta .meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.article-cover {
  flex-shrink: 0;
  width: 180px;
  height: 120px;
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

/* Related Topics */
.related-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.related-item {
  display: flex;
  flex-direction: column;
  background: #f8f9fa;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
}

.related-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.related-cover {
  width: 100%;
  height: 160px;
  overflow: hidden;
}

.related-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.related-item:hover .related-cover img {
  transform: scale(1.05);
}

.related-info {
  padding: 16px;
}

.related-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.related-desc {
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.related-meta {
  font-size: 12px;
  color: #909399;
}

/* 响应式 */
@media (max-width: 768px) {
  .topic-header {
    height: 300px;
  }
  
  .title {
    font-size: 24px;
  }
  
  .description {
    font-size: 14px;
  }
  
  :deep(.el-card__body) {
    padding: 20px;
  }
  
  .article-item {
    flex-direction: column-reverse;
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
}
</style>

