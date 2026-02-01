<template>
  <div class="topic-page">
    <!-- 简洁的页面标题 -->
    <div class="page-banner">
      <div class="banner-content">
        <h1 class="page-title">精选专题</h1>
        <p class="page-subtitle">{{ topics.length }} 个专题 · {{ totalArticles }} 篇文章</p>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="content-container">
      <el-row :gutter="30">
        <!-- 专题列表 -->
        <el-col :span="18" :xs="24">
          <!-- 筛选排序 -->
          <div class="filter-header">
            <h3 class="filter-title">📚 全部专题</h3>
            <el-select v-model="sortBy" style="width: 140px">
              <el-option label="最新更新" value="update" />
              <el-option label="最多阅读" value="views" />
              <el-option label="最多文章" value="articles" />
            </el-select>
          </div>

          <!-- 置顶专题 -->
          <div class="featured-section" v-if="featuredTopic">
            <NuxtLink :to="`/topic/${featuredTopic.id}`" class="featured-card">
              <div class="featured-cover">
                <img :src="featuredTopic.coverImage || `https://picsum.photos/id/${featuredTopic.id}/600/400`" alt="" />
              </div>
              <div class="featured-content">
                <div class="featured-tags">
                  <span class="featured-tag">精选</span>
                </div>
                <h2 class="featured-title">{{ featuredTopic.name }}</h2>
                <p class="featured-desc">{{ featuredTopic.description }}</p>
                <div class="featured-meta">
                  <span class="meta-item">
                    <el-icon><Document /></el-icon>
                    {{ featuredTopic.articleCount || 0 }} 篇文章
                  </span>
                  <span class="meta-item">
                    <el-icon><View /></el-icon>
                    {{ formatNumber(featuredTopic.views || 0) }} 阅读
                  </span>
                  <span class="meta-item">
                    <el-icon><Clock /></el-icon>
                    每周更新
                  </span>
                </div>
              </div>
            </NuxtLink>
          </div>

          <div class="topics-grid">
            <NuxtLink
              v-for="(topic, index) in sortedTopics"
              :key="topic.id"
              :to="`/topic/${topic.id}`"
              class="topic-card"
              :style="{ animationDelay: `${index * 0.1}s` }"
            >
              <div class="topic-cover">
                <img :src="topic.coverImage || `https://picsum.photos/id/${topic.id}/600/400`" alt="" />
                <div class="topic-status">
                  连载中
                </div>
              </div>
              <div class="topic-content">
                <h3 class="topic-title">{{ topic.name }}</h3>
                <p class="topic-desc">{{ topic.description }}</p>
                <div class="topic-footer">
                  <div class="topic-meta">
                    <span class="meta-item">
                      <el-icon><Document /></el-icon>
                      {{ topic.articleCount || 0 }}
                    </span>
                    <span class="meta-item">
                      <el-icon><View /></el-icon>
                      {{ formatNumber(topic.views || 0) }}
                    </span>
                    <span class="meta-item">
                      <el-icon><Clock /></el-icon>
                      {{ formatDate(topic.updatedAt) }}
                    </span>
                  </div>
                </div>
              </div>
            </NuxtLink>
          </div>

          <el-empty v-if="topics.length === 0" description="暂无专题" />
        </el-col>

        <!-- 侧边栏 -->
        <el-col :span="6" :xs="24">
          <div class="sidebar">
            <!-- 热门专题 -->
            <div class="sidebar-card">
              <h4 class="card-title">最受欢迎</h4>
              <div class="popular-list">
                <NuxtLink
                  v-for="(item, index) in popularTopics"
                  :key="item.id"
                  :to="`/topic/${item.id}`"
                  class="popular-item"
                >
                  <div class="popular-rank">{{ index + 1 }}</div>
                  <div class="popular-info">
                    <div class="popular-name">{{ item.name }}</div>
                    <div class="popular-meta">
                      {{ item.articleCount || 0 }} 篇 · {{ formatNumber(item.views || 0) }} 阅读
                    </div>
                  </div>
                </NuxtLink>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Document, View, Clock } from '@element-plus/icons-vue'

useSeoMeta({
  title: '专题系列',
  description: '系统性学习，深入理解技术'
})

const topicApi = useTopicApi()
const sortBy = ref('update')

const { data: topicsData } = await useAsyncData('topics', async () => {
  try {
    const res = await topicApi.getList({ page: 1, size: 100 })
    return res.list || []
  } catch (error) {
    console.error('加载专题失败:', error)
    return []
  }
})

const topics = computed(() => topicsData.value || [])

// 格式化日期
const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  try {
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return dateStr
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))
    
    if (days === 0) return '今天'
    if (days === 1) return '昨天'
    if (days < 7) return `${days}天前`
    if (days < 30) return `${Math.floor(days / 7)}周前`
    if (days < 365) return `${Math.floor(days / 30)}个月前`
    
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  } catch {
    return dateStr
  }
}

const featuredTopic = computed(() => {
  const featured = topics.value.filter((t: any) => t.isFeatured)
  if (featured.length === 0) return null
  return featured.sort((a: any, b: any) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())[0]
})

const sortedTopics = computed(() => {
  const featuredId = featuredTopic.value?.id
  const filtered = topics.value.filter((t: any) => t.id !== featuredId)
  const sorted = [...filtered]

  if (sortBy.value === 'views') {
    sorted.sort((a: any, b: any) => (b.views || 0) - (a.views || 0))
  } else if (sortBy.value === 'articles') {
    sorted.sort((a: any, b: any) => (b.articleCount || 0) - (a.articleCount || 0))
  }

  return sorted
})

const popularTopics = computed(() => {
  return [...topics.value].sort((a: any, b: any) => (b.views || 0) - (a.views || 0)).slice(0, 5)
})

const totalArticles = computed(() => {
  return topics.value.reduce((sum: number, t: any) => sum + (t.articleCount || 0), 0)
})

const formatNumber = (num: number) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'W'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num
}
</script>

<style scoped>
.topic-page {
  min-height: 100vh;
  padding-bottom: 60px;
}

/* 页面标题 - 背景图样式 */
.page-banner {
  position: relative;
  height: 280px;
  background: url('https://picsum.photos/id/201/1920/600') center/cover;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 60px;
  overflow: hidden;
}

.page-banner::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 0;
}

.page-banner::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 20% 50%, rgba(255, 255, 255, 0.1) 0%, transparent 50%),
    radial-gradient(circle at 80% 50%, rgba(255, 255, 255, 0.1) 0%, transparent 50%);
}

.banner-content {
  position: relative;
  z-index: 1;
  text-align: center;
  color: white;
  padding: 0 20px;
}

.page-title {
  font-size: 56px;
  font-weight: 800;
  margin-bottom: 16px;
  letter-spacing: -1px;
  text-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  animation: fadeInDown 0.8s ease-out;
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.page-subtitle {
  font-size: 18px;
  font-weight: 500;
  opacity: 0.95;
  animation: fadeInUp 0.8s ease-out 0.2s backwards;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 内容容器 */
.content-container {
  margin: 0 auto;
}

/* 筛选头部 */
.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.filter-title {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

/* 置顶专题 */
.featured-section {
  margin-bottom: 40px;
}

/* 置顶专题 */
.featured-card {
  display: flex;
  background: white;
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  text-decoration: none;
  color: inherit;
}

.featured-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.12);
}

.featured-cover {
  width: 45%;
  height: 280px;
  flex-shrink: 0;
  overflow: hidden;
}

.featured-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.featured-card:hover .featured-cover img {
  transform: scale(1.08);
}

.featured-content {
  flex: 1;
  padding: 32px;
  display: flex;
  flex-direction: column;
}

.featured-tags {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.featured-tag {
  padding: 4px 12px;
  background: #f0f2f5;
  color: #606266;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 600;
}

.featured-title {
  font-size: 28px;
  font-weight: 700;
  color: #303133;
  margin: 0 0 12px 0;
  line-height: 1.3;
}

.featured-desc {
  font-size: 15px;
  color: #606266;
  line-height: 1.6;
  margin: 0 0 auto 0;
}

.featured-meta {
  display: flex;
  gap: 24px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #f0f2f5;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: #909399;
}

/* 专题网格 */

.topics-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.topic-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
  animation: fade-up 0.6s ease-out forwards;
  opacity: 0;
  transform: translateY(20px);
  text-decoration: none;
  color: inherit;
}

@keyframes fade-up {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.topic-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.08);
}

.topic-cover {
  position: relative;
  height: 180px;
  overflow: hidden;
}

.topic-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s;
}

.topic-card:hover .topic-cover img {
  transform: scale(1.08);
}

.topic-status {
  position: absolute;
  top: 12px;
  right: 12px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(4px);
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  color: #67c23a;
}

.topic-content {
  padding: 20px;
}

.topic-tags {
  display: flex;
  gap: 6px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.topic-tag {
  padding: 3px 10px;
  background: #f0f2f5;
  color: #606266;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 600;
}

.topic-title {
  font-size: 18px;
  font-weight: 700;
  color: #303133;
  margin: 0 0 8px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.topic-desc {
  font-size: 14px;
  color: #606266;
  line-height: 1.5;
  margin: 0 0 16px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.topic-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid #f0f2f5;
}

.topic-meta {
  display: flex;
  gap: 16px;
}

.topic-author {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #606266;
}

/* 侧边栏 */
.sidebar {
  position: sticky;
  top: 80px;
}

.sidebar-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.card-title {
  font-size: 18px;
  font-weight: 700;
  margin: 0 0 20px 0;
  padding-left: 10px;
  border-left: 4px solid #409eff;
}

/* 热门专题 */
.popular-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.popular-item {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  transition: all 0.2s;
  text-decoration: none;
  color: inherit;
}

.popular-item:hover {
  opacity: 0.8;
}

.popular-rank {
  width: 28px;
  height: 28px;
  background: #f2f6fc;
  color: #909399;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  flex-shrink: 0;
}

.popular-item:nth-child(1) .popular-rank {
  background: #ffece8;
  color: #f56c6c;
}

.popular-item:nth-child(2) .popular-rank {
  background: #e8f3ff;
  color: #409eff;
}

.popular-item:nth-child(3) .popular-rank {
  background: #fff7e8;
  color: #e6a23c;
}

.popular-info {
  flex: 1;
  min-width: 0;
}

.popular-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 4px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.popular-meta {
  font-size: 12px;
  color: #909399;
}

/* 响应式 */
@media (max-width: 1024px) {
  .topics-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .page-banner {
    height: 220px;
  }

  .page-title {
    font-size: 36px;
  }

  .page-subtitle {
    font-size: 16px;
  }

  .featured-card {
    flex-direction: column;
  }

  .featured-cover {
    width: 100%;
    height: 200px;
  }

  .featured-content {
    padding: 24px;
  }

  .sidebar {
    margin-top: 40px;
    position: static;
  }
}
</style>
