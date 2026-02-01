<template>
  <div class="topic-detail-container">
    <!-- Header Area -->
    <div class="topic-header">
      <div class="header-content">
        <div class="meta-info">
          <span class="status-tag status-updating">
            连载中
          </span>
        </div>
        <h1 class="title">{{ topic?.name }}</h1>
        <p class="description">{{ topic?.description }}</p>
        <div class="topic-meta">
          <span class="meta-item">
            <el-icon><Document /></el-icon>
            {{ articles.length }} 篇文章
          </span>
          <span class="divider">·</span>
          <span class="meta-item">
            <el-icon><View /></el-icon>
            {{ formatNumber(topic?.views || 0) }} 阅读
          </span>
          <span class="divider">·</span>
          <span class="meta-item">
            <el-icon><Clock /></el-icon>
            每周更新
          </span>
        </div>
        <div class="author-info">
          <el-avatar :size="32" src="https://i.pravatar.cc/150?img=1" />
          <span class="author-name">Admin</span>
          <span class="author-role">专题作者</span>
        </div>
      </div>
      <div class="header-image" :style="{ backgroundImage: `url(${topic?.coverImage || 'https://picsum.photos/id/' + topicId + '/1200/600'})` }"></div>
    </div>

    <!-- Content Area -->
    <el-row justify="center" class="content-wrapper">
      <el-col :span="18" :xs="22" :sm="20" :md="18" :lg="16">
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
            <NuxtLink 
              v-for="(article, index) in sortedArticles" 
              :key="article.id"
              :to="`/article/${article.id}`"
              class="article-item"
            >
              <div class="article-number">{{ index + 1 }}</div>
              <div class="article-content">
                <div class="article-top">
                  <h4 class="article-title">{{ article.title }}</h4>
                </div>
                <p class="article-summary">{{ article.desc }}</p>
                <div class="article-meta">
                  <span class="meta-item">
                    <el-icon><Calendar /></el-icon>
                    {{ article.createdAt }}
                  </span>
                  <span class="meta-item">
                    <el-icon><View /></el-icon>
                    {{ formatNumber(article.views) }}
                  </span>
                </div>
              </div>
              <div class="article-cover" v-if="article.cover">
                <img :src="article.cover" :alt="article.title" />
              </div>
            </NuxtLink>
          </div>

          <el-empty v-if="articles.length === 0" description="该专题暂无文章" />
        </el-card>

        <!-- 相关专题 -->
        <el-card class="related-card" shadow="never" v-if="relatedTopics.length > 0">
          <h3 class="section-title">
            <el-icon><CollectionTag /></el-icon>
            相关专题
          </h3>
          <div class="related-list">
            <NuxtLink 
              v-for="item in relatedTopics" 
              :key="item.id"
              :to="`/topic/${item.id}`"
              class="related-item"
            >
              <div class="related-cover">
                <img :src="item.coverImage || `https://picsum.photos/id/${item.id}/300/200`" alt="" />
              </div>
              <div class="related-info">
                <h4 class="related-title">{{ item.name }}</h4>
                <p class="related-desc">{{ item.description }}</p>
                <div class="related-meta">
                  {{ item.articleCount || 0 }} 篇 · {{ formatNumber(item.views || 0) }} 阅读
                </div>
              </div>
            </NuxtLink>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { Document, View, Clock, List, Calendar, CollectionTag } from '@element-plus/icons-vue'

const route = useRoute()
const topicApi = useTopicApi()

const topicId = computed(() => parseInt(route.params.id as string))
const sortBy = ref('default')

// 格式化日期
const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  try {
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return dateStr
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
  } catch {
    return dateStr
  }
}

// 格式化数字
const formatNumber = (num: number) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'W'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num
}

// 获取专题信息
const { data: topic } = await useAsyncData(`topic-${topicId.value}`, async () => {
  try {
    return await topicApi.getOne({ id: topicId.value })
  } catch (error) {
    console.error('加载专题失败:', error)
    return null
  }
})

// 获取专题文章
const { data: articlesData } = await useAsyncData(`topic-articles-${topicId.value}`, async () => {
  try {
    const res = await topicApi.getTopicArticles({
      topicId: topicId.value,
      page: 1,
      size: 100
    })
    return res.list || []
  } catch (error) {
    console.error('加载专题文章失败:', error)
    return []
  }
})

// 获取相关专题
const { data: relatedTopicsData } = await useAsyncData(`related-topics-${topicId.value}`, async () => {
  try {
    const res = await topicApi.getList({ page: 1, size: 10 })
    return (res.list || []).filter((t: any) => t.id !== topicId.value).slice(0, 2)
  } catch (error) {
    console.error('加载相关专题失败:', error)
    return []
  }
})

const relatedTopics = computed(() => relatedTopicsData.value || [])

const articles = computed(() => {
  return (articlesData.value || []).map((item: any) => ({
    id: item.id,
    title: item.title,
    desc: item.desc || '',
    createdAt: formatDate(item.createdAt),
    views: item.views || 0,
    cover: item.coverImage || `https://picsum.photos/id/${item.id}/400/200`
  }))
})

const sortedArticles = computed(() => {
  const sorted = [...articles.value]
  
  if (sortBy.value === 'date') {
    sorted.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
  } else if (sortBy.value === 'views') {
    sorted.sort((a, b) => b.views - a.views)
  }
  
  return sorted
})

useSeoMeta({
  title: () => topic.value?.name || '专题详情',
  description: () => topic.value?.description || '专题文章列表'
})
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
  text-decoration: none;
  color: inherit;
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
  text-decoration: none;
  color: inherit;
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
