<template>
  <div class="category-page">
    <!-- Banner with background image -->
    <div class="page-banner">
      <div class="banner-content">
        <h1 class="page-title">文章分类</h1>
        <p class="page-subtitle">Browse articles by category</p>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="content-container">
      <div class="category-grid">
        <NuxtLink 
          v-for="(category, index) in categories" 
          :key="category.id"
          :to="`/category/${category.id}`"
          class="category-card"
          :style="{ animationDelay: `${index * 0.08}s` }"
        >
          <div class="card-header" :style="{ 
            backgroundImage: `url(https://picsum.photos/id/${(category.id % 100)}/400/300)`
          }">
            <div class="card-overlay"></div>
            <div class="category-icon">
              <el-icon :size="28"><Folder /></el-icon>
            </div>
          </div>
          <div class="card-body">
            <h3 class="category-name">{{ category.name }}</h3>
            <p class="category-desc">{{ category.description || category.name + '相关技术文章' }}</p>
            <div class="category-footer">
              <span class="article-count">
                <el-icon><Document /></el-icon>
                {{ category.articleCount || 0 }} 篇文章
              </span>
              <el-icon class="arrow-icon"><ArrowRight /></el-icon>
            </div>
          </div>
        </NuxtLink>
      </div>

      <el-empty v-if="categories.length === 0" description="暂无分类" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Folder, Document, ArrowRight } from '@element-plus/icons-vue'

useSeoMeta({
  title: '文章分类',
  description: '按分类浏览所有文章'
})

const categoryApi = useCategoryApi()

const { data: categoriesData } = await useAsyncData('categories', async () => {
  try {
    const res = await categoryApi.getList({ page: 1, size: 100 })
    return res.list || []
  } catch (error) {
    console.error('加载分类失败:', error)
    return []
  }
})

const categories = computed(() => categoriesData.value || [])
</script>

<style scoped>
.category-page {
  min-height: 100vh;
  padding-bottom: 60px;
}

/* Banner with background image */
.page-banner {
  position: relative;
  height: 280px;
  background: url('https://picsum.photos/id/1/1920/600') center/cover;
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
  font-size: 20px;
  font-weight: 500;
  opacity: 0.95;
  animation: fadeInUp 0.8s ease-out 0.2s backwards;
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

/* 内容容器 - 与首页保持一致 */
.content-container {
  margin: 0 auto;
}

/* 分类网格 - 3列布局，充满容器 */
.category-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 30px;
}

.category-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  animation: fadeInScale 0.6s ease-out forwards;
  opacity: 0;
  transform: translateY(20px) scale(0.95);
  text-decoration: none;
}

@keyframes fadeInScale {
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.category-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.12);
}

/* 卡片头部 - 背景图片 */
.card-header {
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  transition: all 0.3s;
  overflow: hidden;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

/* 遮罩层 - 确保图标清晰可见 */
.card-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 0;
  transition: background 0.3s;
}

.category-card:hover .card-overlay {
  background: rgba(0, 0, 0, 0.5);
}

.card-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 30% 40%, rgba(255, 255, 255, 0.15) 0%, transparent 60%);
  opacity: 0;
  transition: opacity 0.3s;
}

.category-card:hover .card-header::before {
  opacity: 1;
}

.category-icon {
  color: white;
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 60px;
  height: 60px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  transition: all 0.3s;
}

.category-card:hover .category-icon {
  transform: scale(1.1) rotate(5deg);
  background: rgba(255, 255, 255, 0.3);
}

/* 卡片主体 */
.card-body {
  padding: 20px;
}

.category-name {
  font-size: 18px;
  font-weight: 700;
  color: #303133;
  margin: 0 0 10px 0;
  line-height: 1.3;
}

.category-desc {
  font-size: 14px;
  color: #606266;
  line-height: 1.5;
  margin: 0 0 16px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  min-height: 42px;
}

.category-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid #f0f2f5;
}

.article-count {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #909399;
  font-weight: 500;
}

.arrow-icon {
  color: #dcdfe6;
  transition: all 0.3s;
}

.category-card:hover .arrow-icon {
  color: #409eff;
  transform: translateX(4px);
}

/* 响应式 */
@media (max-width: 1024px) {
  .category-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 24px;
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

  .category-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .card-header {
    height: 90px;
  }

  .category-icon {
    width: 50px;
    height: 50px;
  }

  .card-body {
    padding: 16px;
  }

  .category-name {
    font-size: 16px;
  }
}
</style>
