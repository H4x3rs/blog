<template>
  <div class="author-page">
    <div class="page-header">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item>
          <NuxtLink to="/">首页</NuxtLink>
        </el-breadcrumb-item>
        <el-breadcrumb-item>作者文章</el-breadcrumb-item>
      </el-breadcrumb>
      <h1 class="page-title">
        <el-icon><User /></el-icon>
        作者文章
      </h1>
      <p class="page-desc">共 {{ total }} 篇文章</p>
    </div>

    <div class="article-list">
      <div 
        v-for="article in articles" 
        :key="article.id" 
        class="article-item"
      >
        <NuxtLink :to="`/article/${article.id}`" class="article-thumb">
          <img :src="article.cover" loading="lazy" :alt="article.title" />
        </NuxtLink>
        <div class="article-detail">
          <div class="article-meta">
            <span class="date">{{ article.createdAt }}</span>
            <span class="views">{{ article.views }} 阅读</span>
          </div>
          <NuxtLink :to="`/article/${article.id}`">
            <h2 class="article-title">{{ article.title }}</h2>
          </NuxtLink>
          <p class="article-summary">{{ article.desc }}</p>
        </div>
      </div>
    </div>

    <el-empty v-if="articles.length === 0" description="该作者暂无文章" />

    <div v-if="total > pageSize" class="pagination-wrapper">
      <el-pagination 
        background 
        layout="prev, pager, next" 
        :total="total" 
        :current-page="currentPage"
        :page-size="pageSize"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { User } from '@element-plus/icons-vue'

const route = useRoute()
const articleApi = useArticleApi()

const authorId = computed(() => parseInt(route.params.id as string))
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

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

// 获取作者文章
const { data: articlesData } = await useAsyncData(`author-articles-${authorId.value}`, async () => {
  try {
    const res = await articleApi.getList({
      page: currentPage.value,
      size: pageSize.value,
      status: 'published',
      authorId: authorId.value
    })
    return res
  } catch (error) {
    console.error('加载文章失败:', error)
    return { list: [], total: 0 }
  }
})

const articles = computed(() => {
  return (articlesData.value?.list || []).map((item: any) => ({
    id: item.id,
    title: item.title,
    desc: item.desc || '',
    createdAt: formatDate(item.createdAt),
    views: item.views || 0,
    cover: item.coverImage || `https://picsum.photos/id/${item.id}/600/400`
  }))
})

watch(articlesData, (data) => {
  total.value = data?.total || 0
}, { immediate: true })

const handlePageChange = async (page: number) => {
  currentPage.value = page
  await refreshNuxtData(`author-articles-${authorId.value}`)
}

useSeoMeta({
  title: '作者文章',
  description: '作者发布的所有文章'
})
</script>

<style scoped lang="scss">
.author-page {
  padding: 20px 0;
}

.page-header {
  margin-bottom: 40px;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  color: #303133;
  margin: 20px 0 12px;
  display: flex;
  align-items: center;
  gap: 12px;

  .el-icon {
    color: #409eff;
  }
}

.page-desc {
  font-size: 16px;
  color: #909399;
}

.article-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.article-item {
  display: flex;
  gap: 24px;
  padding: 24px;
  background: white;
  border-radius: 16px;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 12px 24px rgba(0,0,0,0.06);
  }
}

.article-thumb {
  width: 200px;
  height: 140px;
  flex-shrink: 0;
  border-radius: 12px;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.5s;
  }

  &:hover img {
    transform: scale(1.05);
  }
}

.article-detail {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.article-meta {
  display: flex;
  gap: 16px;
  color: #909399;
  font-size: 13px;
  margin-bottom: 12px;
}

.article-title {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 12px;
  transition: color 0.3s;

  &:hover {
    color: #409eff;
  }
}

.article-summary {
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.pagination-wrapper {
  margin-top: 40px;
  display: flex;
  justify-content: center;
}

@media (max-width: 768px) {
  .article-item {
    flex-direction: column;
  }

  .article-thumb {
    width: 100%;
    height: 180px;
  }
}
</style>
