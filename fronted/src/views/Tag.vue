<template>
  <div class="tag-page">
    <!-- 顶部搜索框 -->
    <div class="search-header">
      <div class="search-container">
        <div class="search-wrapper">
          <el-icon class="search-icon"><Search /></el-icon>
          <input
            v-model="searchKeyword"
            type="text"
            placeholder="搜索标签，例如：Vue、React、Docker..."
            class="search-input"
          />
          <el-button v-if="searchKeyword" text class="clear-btn" @click="searchKeyword = ''">
            <el-icon><CircleClose /></el-icon>
          </el-button>
        </div>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="content-container" v-loading="loading">
      <el-row :gutter="30">
        <!-- 左侧：全部标签 -->
        <el-col :span="18" :xs="24">
          <div class="tags-section">
            <div class="section-header">
              <h3 class="section-title">
                <el-icon><PriceTag /></el-icon>
                全部标签
              </h3>
              <span class="tags-count">{{ filteredTags.length }} 个标签</span>
            </div>
            <div class="tags-cloud">
              <div
                v-for="tag in filteredTags"
                :key="tag.id"
                class="tag-item"
                :class="getSizeClass(tag.count)"
                :style="{ '--tag-color': tag.color }"
                @click="selectTag(tag)"
              >
                <span class="tag-hash">#</span>{{ tag.name }}
                <span class="tag-count">{{ tag.count }}</span>
              </div>
            </div>
            <el-empty v-if="filteredTags.length === 0" description="未找到相关标签" />
          </div>
        </el-col>

        <!-- 右侧：热门标签 -->
        <el-col :span="6" :xs="24">
          <div class="sidebar">
            <div class="hot-tags-card">
              <h4 class="card-title">
                <el-icon><TrendCharts /></el-icon>
                热点标签
              </h4>
              <el-divider class="card-divider" />
              <ul class="hot-tags-list">
                <li
                  v-for="(tag, index) in hotTags"
                  :key="tag.id"
                  class="hot-tag-item"
                  @click="selectTag(tag)"
                >
                  <span class="hot-number">{{ index + 1 }}</span>
                  <div class="hot-tag-text">#{{ tag.name }}</div>
                  <el-icon class="arrow-icon"><ArrowRight /></el-icon>
                </li>
              </ul>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, PriceTag, CircleClose, TrendCharts, ArrowRight } from '@element-plus/icons-vue'
import { getTagList } from '@/api/tag'
import { getArticleList } from '@/api/article'
import { ElMessage } from 'element-plus'

const router = useRouter()
const searchKeyword = ref('')
const allTags = ref([])
const loading = ref(true)

// 加载标签列表
const loadTags = async () => {
  loading.value = true
  try {
    const res = await getTagList({
      page: 1,
      size: 100
    })
    
    // 获取每个标签的文章数量
    const tagsWithCount = await Promise.all((res.list || []).map(async (item) => {
      try {
        // 这里需要通过标签查询文章，但后端可能不支持，暂时设为0
        return {
          id: item.id,
          name: item.name,
          slug: item.slug,
          color: item.color || '#409EFF',
          count: item.count || 0 // 使用后端返回的count字段
        }
      } catch (error) {
        return {
          id: item.id,
          name: item.name,
          slug: item.slug,
          color: item.color || '#409EFF',
          count: 0
        }
      }
    }))
    
    allTags.value = tagsWithCount
  } catch (error) {
    console.error('加载标签失败:', error)
    ElMessage.error('加载标签失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadTags()
})

const hotTags = computed(() => {
  return [...allTags.value].sort((a, b) => b.count - a.count).slice(0, 5)
})

const filteredTags = computed(() => {
  if (!searchKeyword.value) return allTags.value
  return allTags.value.filter(tag =>
    tag.name.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
})

const getSizeClass = (count) => {
  if (count >= 40) return 'size-xl'
  if (count >= 30) return 'size-lg'
  if (count >= 20) return 'size-md'
  return 'size-sm'
}

const selectTag = (tag) => {
  // 跳转到标签详情页面，优先使用slug，如果没有则使用name
  const identifier = tag.slug || tag.name.toLowerCase().replace(/\s+/g, '-')
  router.push(`/tag/${identifier}`)
}
</script>

<style scoped>
.tag-page {
  min-height: 100vh;
  padding-bottom: 60px;
}

/* 搜索头部 */
.search-header {
  padding: 30px 0;
  margin-bottom: 40px;
}

.search-container {
  width: 100%;
}

.search-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  background: white;
  border-radius: 50px;
  padding: 0 20px;
  box-shadow: 
    0 6px 20px rgba(0, 0, 0, 0.06),
    0 0 0 1px rgba(0, 0, 0, 0.02);
  transition: all 0.3s;
}

.search-wrapper:hover {
  box-shadow: 
    0 15px 50px rgba(0, 0, 0, 0.12),
    0 0 0 1px rgba(64, 158, 255, 0.2);
}

.search-wrapper:focus-within {
  box-shadow: 
    0 15px 50px rgba(64, 158, 255, 0.15),
    0 0 0 2px #409eff;
}

.search-icon {
  font-size: 20px;
  color: #909399;
  margin-right: 10px;
  flex-shrink: 0;
}

.search-wrapper:focus-within .search-icon {
  color: #409eff;
}

.search-input {
  flex: 1;
  height: 48px;
  border: none;
  outline: none;
  font-size: 15px;
  color: #303133;
  background: transparent;
}

.search-input::placeholder {
  color: #c0c4cc;
}

.clear-btn {
  padding: 8px;
  font-size: 20px;
  color: #909399;
  flex-shrink: 0;
}

.clear-btn:hover {
  color: #606266;
}

/* 内容容器 */
.content-container {
  margin: 0 auto;
}

/* 标签区域 */
.tags-section {
  background: white;
  border-radius: 16px;
  padding: 30px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 2px solid #f0f2f5;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 22px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.section-title .el-icon {
  font-size: 24px;
  color: #409eff;
}

.tags-count {
  font-size: 14px;
  color: #909399;
  background: #f4f4f5;
  padding: 6px 16px;
  border-radius: 20px;
  font-weight: 500;
}

/* 标签云 */
.tags-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.tag-item {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 18px;
  background: white;
  border: 2px solid #e4e7ed;
  border-radius: 24px;
  cursor: pointer;
  transition: all 0.3s;
  font-weight: 600;
  color: #606266;
}

.tag-item:hover {
  border-color: var(--tag-color);
  color: var(--tag-color);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.tag-hash {
  font-weight: 700;
  opacity: 0.6;
}

.tag-count {
  background: #f0f2f5;
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 12px;
  color: #909399;
}

/* 标签大小变化 */
.tag-item.size-xl {
  font-size: 20px;
  padding: 12px 24px;
}

.tag-item.size-lg {
  font-size: 18px;
  padding: 11px 20px;
}

.tag-item.size-md {
  font-size: 16px;
  padding: 10px 18px;
}

.tag-item.size-sm {
  font-size: 15px;
  padding: 9px 16px;
}

/* 侧边栏 */
.sidebar {
  position: sticky;
  top: 80px;
}

.hot-tags-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 700;
  margin: 0;
  color: #303133;
}

.card-title .el-icon {
  font-size: 20px;
  color: #409eff;
}

.card-divider {
  margin: 16px 0 20px 0;
}

/* 热门标签列表 - 简化版 */
.hot-tags-list {
  padding: 0;
  margin: 0;
  list-style: none;
}

.hot-tag-item {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  cursor: pointer;
}

.hot-tag-item:last-child {
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
  flex-shrink: 0;
}

.hot-tag-item:nth-child(1) .hot-number { 
  background: #ffece8;
  color: #f56c6c;
}

.hot-tag-item:nth-child(2) .hot-number { 
  background: #e8f3ff;
  color: #409eff;
}

.hot-tag-item:nth-child(3) .hot-number { 
  background: #fff7e8;
  color: #e6a23c;
}

.hot-tag-text {
  flex: 1;
  font-size: 14px;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.2s;
}

.hot-tag-item:hover .hot-tag-text {
  color: #409eff;
}

.arrow-icon {
  font-size: 14px;
  color: #dcdfe6;
  transition: all 0.3s;
}

.hot-tag-item:hover .arrow-icon {
  color: #409eff;
  transform: translateX(2px);
}

/* 响应式 */
@media (max-width: 768px) {
  .search-header {
    padding: 30px 0;
  }


  .tags-section {
    padding: 20px;
  }

  .section-title {
    font-size: 18px;
  }

  .tags-cloud {
    gap: 12px;
  }

  .tag-item.size-xl {
    font-size: 16px;
    padding: 10px 18px;
  }

  .tag-item.size-lg {
    font-size: 15px;
    padding: 9px 16px;
  }

  .tag-item.size-md {
    font-size: 14px;
    padding: 8px 14px;
  }

  .tag-item.size-sm {
    font-size: 13px;
    padding: 7px 12px;
  }

  .sidebar {
    margin-top: 30px;
    position: static;
  }
}
</style>
