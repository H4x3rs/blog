<template>
  <div class="category-detail-container" v-loading="loading">
    <!-- Header Area -->
    <div class="category-header">
      <div class="header-content">
        <div class="category-icon-large">
          <el-icon :size="60">
            <component :is="category.icon || Folder" />
          </el-icon>
        </div>
        <h1 class="title">{{ category.name }}</h1>
        <p class="description">{{ category.description }}</p>
        <div class="category-meta">
          <span class="meta-item">
            <el-icon><Document /></el-icon>
            {{ articles.length }} 篇文章
          </span>
          <span class="divider">·</span>
          <span class="meta-item">
            <el-icon><View /></el-icon>
            {{ formatNumber(totalViews) }} 阅读
          </span>
        </div>
      </div>
      <div class="header-image" :style="{ backgroundImage: `url(${category.bgImage})` }"></div>
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
                  <span class="meta-item tags">
                    <el-tag 
                      v-for="tag in article.tags.slice(0, 3)" 
                      :key="tag" 
                      size="small" 
                      effect="plain"
                    >
                      {{ tag }}
                    </el-tag>
                  </span>
                </div>
              </div>
            </div>
          </div>
        </el-card>

        <!-- 相关分类 -->
        <el-card class="related-card" shadow="never" v-if="relatedCategories.length > 0">
          <h3 class="section-title">
            <el-icon><CollectionTag /></el-icon>
            相关分类
          </h3>
          <div class="related-list">
            <div 
              v-for="item in relatedCategories" 
              :key="item.id"
              class="related-item"
              @click="viewCategory(item)"
            >
              <div class="related-icon">
                <el-icon :size="32">
                  <component :is="item.icon || Folder" />
                </el-icon>
              </div>
              <div class="related-info">
                <h4 class="related-title">{{ item.name }}</h4>
                <p class="related-desc">{{ item.description }}</p>
                <div class="related-meta">
                  {{ item.count }} 篇文章
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
import { ref, computed, onMounted, watch, markRaw } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  Document, 
  View, 
  List, 
  Calendar, 
  ChatDotRound,
  CollectionTag,
  Folder,
  Platform,
  Monitor,
  Cpu,
  Coin,
  Setting,
  DataAnalysis,
  Collection,
  Grid,
  Menu,
  Operation,
  Tools,
  Connection,
  Link,
  Share,
  Star,
  Trophy
} from '@element-plus/icons-vue'
import { getCategoryBySlug, getCategoryList } from '@/api/category'
import { getArticleList } from '@/api/article'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const categorySlug = computed(() => route.params.slug)

const sortBy = ref('date')
const loading = ref(true)

// 图标组件映射（支持所有Element Plus图标）- 使用 markRaw 避免响应式化
const iconComponentMap = {
  Folder: markRaw(Folder),
  Platform: markRaw(Platform),
  Monitor: markRaw(Monitor),
  Cpu: markRaw(Cpu),
  Coin: markRaw(Coin),
  Setting: markRaw(Setting),
  Document: markRaw(Document),
  DataAnalysis: markRaw(DataAnalysis),
  Collection: markRaw(Collection),
  Grid: markRaw(Grid),
  List: markRaw(List),
  Menu: markRaw(Menu),
  Operation: markRaw(Operation),
  Tools: markRaw(Tools),
  Connection: markRaw(Connection),
  Link: markRaw(Link),
  Share: markRaw(Share),
  Star: markRaw(Star),
  Trophy: markRaw(Trophy)
}

// 获取图标组件（如果图标不存在则使用默认的Folder图标）
const getIconComponent = (iconName) => {
  if (!iconName) return iconComponentMap.Folder
  return iconComponentMap[iconName] || iconComponentMap.Folder
}

const category = ref({
  id: 0,
  slug: categorySlug.value,
  name: '',
  description: '',
  icon: markRaw(Folder),
  bgImage: 'https://picsum.photos/id/1084/1920/600',
  count: 0
})

const articles = ref([])
const relatedCategories = ref([])

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

// 加载分类信息
const loadCategory = async () => {
  try {
    const res = await getCategoryBySlug({ slug: categorySlug.value })
    if (res) {
      const cat = res
      category.value = {
        id: cat.id,
        slug: cat.slug,
        name: cat.name,
        description: cat.description || cat.name + '相关技术文章',
        icon: getIconComponent(cat.icon || 'Folder'),
        bgImage: 'https://picsum.photos/id/' + (cat.id % 100) + '/1920/600',
        count: 0
      }
      await loadArticles()
    } else {
      ElMessage.error('分类不存在')
    }
  } catch (error) {
    console.error('加载分类失败:', error)
    ElMessage.error('加载分类失败')
  }
}

// 加载文章列表
const loadArticles = async () => {
  try {
    const res = await getArticleList({
      categoryId: category.value.id,
      page: 1,
      size: 100,
      status: 'published'
    })
    articles.value = (res.list || []).map(item => ({
      id: item.id,
      title: item.title,
      summary: item.desc || '',
      cover: item.coverImage || 'https://picsum.photos/id/' + item.id + '/400/300',
      createdAt: formatDate(item.createdAt),
      views: item.views || 0,
      comments: 0, // 评论功能暂未实现
      tags: item.tags ? (typeof item.tags === 'string' ? item.tags.split(',') : item.tags) : [],
      isNew: false,
      isHot: false
    }))
    category.value.count = articles.value.length
  } catch (error) {
    console.error('加载文章失败:', error)
  }
}

// 加载相关分类
const loadRelatedCategories = async () => {
  try {
    const res = await getCategoryList({
      page: 1,
      size: 10
    })
    const related = (res.list || [])
      .filter(cat => cat.id !== category.value.id)
      .slice(0, 3)
      .map(item => ({
        id: item.id,
        slug: item.slug,
        name: item.name,
        description: item.description || item.name + '相关技术文章',
        icon: getIconComponent(item.icon || 'Folder'),
        count: 0
      }))
    
    // 获取每个分类的文章数量
    for (let item of related) {
      try {
        const articlesRes = await getArticleList({
          categoryId: item.id,
          page: 1,
          size: 1
        })
        item.count = articlesRes.total || 0
      } catch (error) {
        item.count = 0
      }
    }
    
    relatedCategories.value = related
  } catch (error) {
    console.error('加载相关分类失败:', error)
  }
}

onMounted(async () => {
  loading.value = true
  await loadCategory()
  await loadRelatedCategories()
  loading.value = false
})

watch(() => route.params.slug, async (newSlug) => {
  if (newSlug) {
    loading.value = true
    await loadCategory()
    await loadRelatedCategories()
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

const viewCategory = (category) => {
  router.push(`/category/${category.slug}`)
}
</script>

<style scoped>
.category-detail-container {
  background-color: #f8f9fa;
  min-height: 100vh;
  padding-bottom: 60px;
}

/* Header Styles */
.category-header {
  position: relative;
  min-height: 450px;
  background-color: #2c3e50;
  color: white;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  padding: 20px 0;
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
  padding: 40px 20px;
}

.category-icon-large {
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

.category-meta {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 14px;
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

.article-meta .meta-item.tags {
  display: flex;
  gap: 6px;
}

/* Related Categories */
.related-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.related-item {
  display: flex;
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

.related-icon {
  flex-shrink: 0;
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border-radius: 12px;
  color: #409eff;
}

.related-info {
  flex: 1;
  min-width: 0;
}

.related-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 6px 0;
  line-height: 1.4;
}

.related-desc {
  font-size: 13px;
  color: #606266;
  line-height: 1.5;
  margin: 0 0 8px 0;
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
  .category-header {
    min-height: 380px;
    padding: 20px 0;
  }
  
  .header-content {
    padding: 30px 20px;
  }
  
  .category-icon-large {
    width: 80px;
    height: 80px;
    margin-bottom: 16px;
  }
  
  .category-icon-large .el-icon {
    font-size: 48px;
  }
  
  .title {
    font-size: 28px;
  }
  
  .description {
    font-size: 14px;
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
}
</style>

