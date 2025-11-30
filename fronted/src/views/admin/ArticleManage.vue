<template>
  <div class="article-manage">
    <!-- 筛选区域 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="文章标题">
          <el-input v-model="filters.title" placeholder="请输入标题关键字" clearable />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="filters.category" placeholder="请选择分类" clearable style="width: 180px">
            <el-option label="Golang" value="golang" />
            <el-option label="Frontend" value="frontend" />
            <el-option label="Database" value="database" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
           <el-select v-model="filters.status" placeholder="状态" clearable style="width: 120px">
            <el-option label="已发布" value="published" />
            <el-option label="草稿" value="draft" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
          <el-button :icon="Refresh" @click="resetFilters">重置</el-button>
        </el-form-item>
        <el-form-item style="margin-left: auto;">
           <el-button type="primary" :icon="Plus" @click="handleCreate">发布文章</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 列表区域 (Grid Card Layout) -->
    <div class="article-list" v-loading="loading">
      <el-empty v-if="tableData.length === 0" description="暂无文章数据" />
      <div class="article-grid" v-else>
        <div v-for="item in tableData" :key="item.id" class="article-card-wrapper">
          <el-card shadow="hover" class="article-card" :body-style="{ padding: '0px' }">
            <div class="article-cover">
              <el-image :src="item.cover" fit="cover" lazy>
                <template #error>
                  <div class="image-slot">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <div class="article-status">
                 <el-tag :type="item.status === 'published' ? 'success' : 'info'" effect="dark" size="small">
                   {{ item.status === 'published' ? '已发布' : '草稿' }}
                 </el-tag>
              </div>
              <div class="article-actions-overlay">
                 <el-button circle type="primary" :icon="Edit" @click.stop="openDialog('edit', item)"></el-button>
                 <el-button circle type="danger" :icon="Delete" @click.stop="handleDelete(item)"></el-button>
              </div>
            </div>
            <div class="article-content">
              <div class="article-meta-top">
                <el-tag size="small" effect="plain" class="cat-tag">{{ item.category }}</el-tag>
                <span class="date">{{ item.date }}</span>
              </div>
              <h3 class="article-title" :title="item.title">{{ item.title }}</h3>
              <p class="article-desc">{{ item.desc }}</p>
              <div class="article-footer">
                <div class="stats">
                  <span><el-icon><View /></el-icon> {{ item.views }}</span>
                  <span v-if="item.publishedBy" class="publisher">
                    <el-icon><User /></el-icon>
                    {{ item.publishedBy }}
                  </span>
                </div>
                <el-button link type="primary" size="small">查看详情</el-button>
              </div>
            </div>
          </el-card>
        </div>
      </div>
      
      <div class="pagination-container" v-if="totalCount > pageSize">
        <el-pagination 
          background 
          layout="total, sizes, prev, pager, next, jumper" 
          :total="totalCount" 
          :page-sizes="[12, 24, 36, 48]"
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
        />
      </div>
    </div>

    <!-- 编辑/新增 弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'create' ? '发布文章' : '编辑文章'"
      width="60%"
      destroy-on-close
      align-center
    >
      <el-form :model="form" label-width="80px">
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="请输入文章标题" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="form.categoryId" placeholder="请选择分类" style="width: 100%">
            <el-option label="Golang" :value="1" />
            <el-option label="Frontend" :value="2" />
            <el-option label="Database" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" placeholder="请选择状态" style="width: 100%">
            <el-option label="草稿" value="draft" />
            <el-option label="已发布" value="published" />
          </el-select>
        </el-form-item>
        <el-form-item label="摘要">
          <el-input v-model="form.desc" type="textarea" :rows="3" placeholder="请输入文章摘要" />
        </el-form-item>
        <el-form-item label="内容">
          <!-- 这里应该是一个富文本编辑器，暂用 Textarea 替代 -->
          <el-input 
            v-model="form.content" 
            type="textarea" 
            :rows="10" 
            placeholder="Markdown / Rich Text Content" 
            style="font-family: monospace"
          />
        </el-form-item>
        <el-form-item label="封面图">
             <el-input v-model="form.coverImage" placeholder="图片 URL" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Refresh, Plus, Edit, Delete, Picture, View, User } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getArticleList, createArticle, updateArticle, deleteArticle } from '@/api/article'

const router = useRouter()

// Filters
const filters = reactive({
  title: '',
  category: '',
  status: ''
})

// Pagination
const currentPage = ref(1)
const pageSize = ref(12)
const totalCount = ref(0)
const loading = ref(false)

// Data
const tableData = ref([])

// Dialog
const dialogVisible = ref(false)
const dialogType = ref('create')
const form = reactive({
  id: null,
  title: '',
  categoryId: null,
  desc: '',
  content: '',
  coverImage: '',
  status: 'draft'
})

// 获取列表
const fetchList = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value,
      title: filters.title,
      categoryId: filters.category ? parseInt(filters.category) : 0,
      status: filters.status
    }
    const res = await getArticleList(params)
    // 转换数据格式以适配卡片显示
    tableData.value = (res.list || []).map(item => ({
      id: item.id,
      title: item.title,
      desc: item.desc,
      cover: item.coverImage,
      status: item.status,
      views: item.views || 0,
      date: formatDate(item.createdAt),
      category: item.categoryId || '',
      content: item.content,
      publishedBy: item.publishedByUser ? (item.publishedByUser.nickname || item.publishedByUser.username) : ''
    }))
    totalCount.value = res.total || 0
  } catch (error) {
    ElMessage.error('获取文章列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

const resetFilters = () => {
  filters.title = ''
  filters.category = ''
  filters.status = ''
  currentPage.value = 1
  fetchList()
}

// 新建文章 - 跳转到编辑页面
const handleCreate = () => {
  router.push('/admin/articles/edit')
}

const openDialog = (type, row) => {
  if (type === 'edit' && row) {
    // 编辑文章 - 跳转到编辑页面
    router.push(`/admin/articles/edit/${row.id}`)
  } else {
    // 新建文章 - 跳转到编辑页面
    handleCreate()
  }
}

const handleDelete = (row) => {
  ElMessageBox.confirm(
    '确定要删除这篇文章吗？删除后无法恢复。',
    '警告',
    {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await deleteArticle({ id: row.id })
      ElMessage.success('删除成功')
      fetchList()
    } catch (error) {
      ElMessage.error('删除失败')
      console.error(error)
    }
  }).catch(() => {})
}

const handleSubmit = async () => {
  if (!form.title || !form.content) {
    ElMessage.warning('请填写标题和内容')
    return
  }
  
  try {
    if (dialogType.value === 'create') {
      await createArticle({
        title: form.title,
        content: form.content,
        desc: form.desc,
        coverImage: form.coverImage,
        categoryId: form.categoryId || 0,
        status: form.status
      })
      ElMessage.success('发布成功')
    } else {
      await updateArticle({
        id: form.id,
        title: form.title,
        content: form.content,
        desc: form.desc,
        coverImage: form.coverImage,
        categoryId: form.categoryId || 0,
        status: form.status
      })
      ElMessage.success('更新成功')
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    ElMessage.error('操作失败')
    console.error(error)
  }
}

// 日期格式化函数
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
    // 如果已经是格式化的字符串，直接返回
    if (typeof dateStr === 'string' && /^\d{4}-\d{2}-\d{2}/.test(dateStr)) {
      return dateStr.split(' ')[0] || dateStr.split('T')[0]
    }
    return dateStr
  }
}

// 监听分页变化
watch([currentPage, pageSize], () => {
  fetchList()
})

// 初始化
onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.article-manage {
  /* padding: 20px; handled by layout */
  min-height: 100vh;
}

.filter-card {
  margin-bottom: 20px;
  border-radius: 12px;
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.filter-form {
  display: flex;
  align-items: center;
}

.filter-form .el-form-item {
  margin-bottom: 0;
  margin-right: 20px;
}

/* Article Grid Layout */
.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.article-card {
  border: none;
  border-radius: 12px;
  transition: all 0.3s;
  overflow: hidden;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.article-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.08);
}

.article-cover {
  height: 160px;
  position: relative;
  overflow: hidden;
}

.article-cover .el-image {
  width: 100%;
  height: 100%;
}

.image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  color: #909399;
  font-size: 24px;
}

.article-status {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 2;
}

.article-actions-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  opacity: 0;
  transition: opacity 0.3s;
}

.article-cover:hover .article-actions-overlay {
  opacity: 1;
}

.article-content {
  padding: 16px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.article-meta-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.date {
  font-size: 12px;
  color: #909399;
}

.article-title {
  margin: 0 0 8px;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  line-height: 1.4;
  height: 44px; /* 2 lines */
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-desc {
  font-size: 13px;
  color: #606266;
  line-height: 1.5;
  margin: 0 0 16px;
  height: 40px; /* 2 lines approx */
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex: 1;
}

.article-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f2f5;
  font-size: 13px;
  color: #909399;
}

.stats {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stats .publisher {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #909399;
  font-size: 13px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 30px;
  padding: 20px 0;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

/* 分页样式优化 */
.pagination-container :deep(.el-pagination) {
  padding: 0 20px;
}

.pagination-container :deep(.el-pagination.is-background .el-pager li) {
  border-radius: 6px;
  font-weight: 500;
}

.pagination-container :deep(.el-pagination.is-background .el-pager li:not(.is-disabled).is-active) {
  background: #409eff;
  color: #fff;
}

.pagination-container :deep(.el-pagination.is-background .btn-next),
.pagination-container :deep(.el-pagination.is-background .btn-prev) {
  border-radius: 6px;
  font-weight: 500;
}

.pagination-container :deep(.el-pagination__total) {
  font-weight: 500;
  color: #606266;
}

.pagination-container :deep(.el-pagination__sizes) {
  margin-right: 16px;
}

.pagination-container :deep(.el-select .el-input__wrapper) {
  border-radius: 6px;
}

.pagination-container :deep(.el-pagination__jump) {
  margin-left: 16px;
}

.pagination-container :deep(.el-input__wrapper) {
  border-radius: 6px;
}
</style>
