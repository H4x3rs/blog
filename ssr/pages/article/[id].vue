<template>
  <div class="article-detail-container" v-loading="pending">
    <!-- Header Area -->
    <div class="article-header">
      <div class="header-content">
        <div class="meta-info">
          <el-tag effect="dark" round class="category-tag">{{ article?.category }}</el-tag>
          <el-tag 
            v-if="article?.topic" 
            effect="dark" 
            round 
            class="topic-tag"
            @click="goToTopic"
          >
            📚 {{ article.topic.name }}
          </el-tag>
          <span class="date">{{ article?.createdAt }}</span>
        </div>
        <h1 class="title">{{ article?.title }}</h1>
        <div class="author-info">
          <el-avatar 
            :size="32" 
            :src="article?.publishedByUser?.avatar || 'https://picsum.photos/id/64/100/100'"
          />
          <NuxtLink 
            v-if="article?.publishedByUser"
            :to="`/author/${article.publishedByUser.id}`"
            class="author-name"
          >
            {{ article.publishedByUser.nickname || article.publishedByUser.username }}
          </NuxtLink>
          <span v-else class="author-name">Admin</span>
          <span class="divider">·</span>
          <span class="read-time">5 min read</span>
          <span class="divider">·</span>
          <span class="views">{{ article?.views }} views</span>
        </div>
      </div>
      <div class="header-image" :style="{ backgroundImage: `url(${article?.cover})` }"></div>
    </div>

    <!-- Content Area -->
    <el-row justify="center" class="content-wrapper">
      <el-col :span="16" :xs="22">
        <!-- 专题信息 -->
        <el-card v-if="article?.topic" class="topic-info-card" shadow="never">
          <NuxtLink :to="`/topic/${article.topic.id}`" class="topic-info">
            <div class="topic-icon">📚</div>
            <div class="topic-content">
              <div class="topic-label">本文属于专题</div>
              <div class="topic-name">{{ article.topic.name }}</div>
              <div class="topic-desc">{{ article.topic.description }}</div>
              <div class="topic-stats">
                <span>{{ article.topic.articleCount }} 篇文章</span>
              </div>
            </div>
            <div class="topic-action">
              <el-icon><ArrowRight /></el-icon>
            </div>
          </NuxtLink>
        </el-card>

        <el-card class="content-card" shadow="never">
          <div class="markdown-body" v-html="renderedContent"></div>
          
          <div class="article-footer">
            <div class="tags-list">
              <NuxtLink 
                v-for="tag in article?.tags" 
                :key="typeof tag === 'object' ? tag.id : tag"
                :to="`/tag/${typeof tag === 'object' ? tag.id : tag}`"
              >
                <el-tag class="tag-item" effect="plain">
                  # {{ typeof tag === 'object' ? tag.name : tag }}
                </el-tag>
              </NuxtLink>
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

        <!-- Comments Area -->
        <div class="comments-section">
          <h3>评论 (0)</h3>
          <el-empty description="暂无评论，快来抢沙发吧！" />
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { Star, Share, ArrowRight } from '@element-plus/icons-vue'
import { marked } from 'marked'
import { useSiteStore } from '~/stores/site'

const route = useRoute()
const router = useRouter()
const siteStore = useSiteStore()
const articleApi = useArticleApi()

const articleId = computed(() => parseInt(route.params.id as string))

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

// 获取文章数据
const { data: articleData, pending } = await useAsyncData(
  `article-${articleId.value}`,
  async () => {
    try {
      const res = await articleApi.getOne({ id: articleId.value })
      return {
        id: res.id,
        title: res.title || '',
        category: res.categoryName || '',
        createdAt: formatDate(res.createdAt),
        views: res.views || 0,
        cover: res.coverImage || `https://picsum.photos/id/${res.id}/1200/600`,
        tags: res.tags && Array.isArray(res.tags) ? res.tags : [],
        content: res.content || '',
        desc: res.desc || '',
        publishedByUser: res.publishedByUser || null,
        topic: res.topic || null
      }
    } catch (error) {
      console.error('加载文章失败:', error)
      return null
    }
  }
)

const article = computed(() => articleData.value)

// 渲染 Markdown 内容
const renderedContent = computed(() => {
  if (!article.value?.content) return ''
  return marked.parse(article.value.content) as string
})

// SEO
useSeoMeta({
  title: () => article.value?.title || '文章详情',
  ogTitle: () => article.value?.title || '文章详情',
  description: () => article.value?.desc || article.value?.title || '',
  ogDescription: () => article.value?.desc || article.value?.title || '',
  ogImage: () => article.value?.cover || '',
  ogType: 'article'
})

// 设置页面标题
useHead({
  title: () => `${article.value?.title || '文章详情'} - ${siteStore.siteName}`
})

// 跳转到专题
const goToTopic = () => {
  if (article.value?.topic?.id) {
    router.push(`/topic/${article.value.topic.id}`)
  }
}
</script>

<style scoped lang="scss">
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
  flex-wrap: wrap;
}

.category-tag {
  background-color: #409eff;
  border-color: #409eff;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.topic-tag {
  cursor: pointer;
  background-color: #67c23a;
  border-color: #67c23a;
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
  flex-wrap: wrap;
}

.author-name {
  margin-left: 8px;
  font-weight: 600;
  color: white;
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
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

.topic-info-card {
  margin-bottom: 20px;
  border-radius: 16px;
  border: none;
}

.topic-info {
  display: flex;
  align-items: center;
  gap: 16px;
  cursor: pointer;
  text-decoration: none;
  color: inherit;

  &:hover .topic-name {
    color: #409eff;
  }
}

.topic-icon {
  font-size: 40px;
}

.topic-content {
  flex: 1;
}

.topic-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
}

.topic-name {
  font-size: 18px;
  font-weight: 700;
  color: #303133;
  transition: color 0.3s;
}

.topic-desc {
  font-size: 14px;
  color: #606266;
  margin-top: 4px;
}

.topic-stats {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}

.topic-action {
  color: #909399;
}

.content-card {
  padding: 40px;
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.05);
  background: white;
  border: none;
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

  a {
    text-decoration: none;
  }
}

.tag-item {
  cursor: pointer;
}

.comments-section {
  margin-top: 30px;

  h3 {
    margin-bottom: 20px;
    font-size: 18px;
    color: #303133;
  }
}

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
}
</style>
