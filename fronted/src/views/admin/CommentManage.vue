<template>
  <div class="page-container">
    <!-- 筛选区域 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="请选择状态" clearable style="width: 150px" size="default">
            <el-option label="全部" value="" />
            <el-option label="待审核" value="pending" />
            <el-option label="已审核" value="approved" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item label="文章ID">
          <el-input 
            v-model.number="filters.articleId" 
            placeholder="请输入文章ID" 
            clearable 
            style="width: 200px"
            size="default"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch" size="default">查询</el-button>
          <el-button :icon="Refresh" @click="resetFilters" size="default">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格区域 -->
    <el-card shadow="never" class="table-card">
      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户" width="120">
          <template #default="scope">
            <div class="user-info">
              <el-avatar :size="32" :src="scope.row.userAvatar || 'https://picsum.photos/id/64/100/100'" />
              <span class="user-name">{{ scope.row.userName || '匿名用户' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="articleTitle" label="文章" min-width="200">
          <template #default="scope">
            <el-link type="primary" :underline="false">{{ scope.row.articleTitle || `文章ID: ${scope.row.articleId}` }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="评论内容" min-width="300" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button 
              v-if="scope.row.status !== 'approved'" 
              link 
              type="success" 
              size="small" 
              :icon="Check" 
              @click="handleApprove(scope.row)"
            >
              审核通过
            </el-button>
            <el-button 
              v-if="scope.row.status !== 'rejected'" 
              link 
              type="warning" 
              size="small" 
              :icon="Close" 
              @click="handleReject(scope.row)"
            >
              拒绝
            </el-button>
            <el-button 
              link 
              type="danger" 
              size="small" 
              :icon="Delete" 
              @click="handleDelete(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination 
          background 
          layout="total, sizes, prev, pager, next, jumper" 
          :total="totalCount" 
          :page-sizes="[10, 20, 30, 50]"
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Search, Refresh, Delete, Check, Close } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCommentManageList, updateComment, deleteComment } from '@/api/comment'

const filters = ref({
  status: '',
  articleId: null
})

const tableData = ref([])
const loading = ref(false)
const totalCount = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const res = await getCommentManageList({
      status: filters.value.status || '',
      articleId: filters.value.articleId || 0,
      page: currentPage.value,
      size: pageSize.value
    })
    if (res) {
      tableData.value = res.list || []
      totalCount.value = res.total || 0
    }
  } catch (error) {
    console.error('加载评论列表失败:', error)
    ElMessage.error('加载评论列表失败')
  } finally {
    loading.value = false
  }
}

// 查询
const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

// 重置筛选
const resetFilters = () => {
  filters.value = {
    status: '',
    articleId: null
  }
  handleSearch()
}

// 获取状态类型
const getStatusType = (status) => {
  const map = {
    'pending': 'warning',
    'approved': 'success',
    'rejected': 'danger'
  }
  return map[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    'pending': '待审核',
    'approved': '已审核',
    'rejected': '已拒绝'
  }
  return map[status] || status
}

// 审核通过
const handleApprove = async (row) => {
  try {
    await ElMessageBox.confirm('确定要审核通过这条评论吗？', '提示', {
      type: 'warning'
    })
    await updateComment({
      id: row.id,
      status: 'approved'
    })
    ElMessage.success('审核通过成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('审核通过失败:', error)
      ElMessage.error('审核通过失败')
    }
  }
}

// 拒绝
const handleReject = async (row) => {
  try {
    await ElMessageBox.confirm('确定要拒绝这条评论吗？', '提示', {
      type: 'warning'
    })
    await updateComment({
      id: row.id,
      status: 'rejected'
    })
    ElMessage.success('拒绝成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('拒绝失败:', error)
      ElMessage.error('拒绝失败')
    }
  }
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？删除后无法恢复！', '警告', {
      type: 'warning'
    })
    await deleteComment({ id: row.id })
    ElMessage.success('删除成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
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

/* 搜索栏输入框统一样式 */
.filter-form :deep(.el-input__wrapper) {
  height: 32px;
  border-radius: 4px;
}

.filter-form :deep(.el-input__inner) {
  font-size: 14px;
  line-height: 1.5;
}

.table-card { 
  border-radius: 16px; 
  border: none; 
  box-shadow: 0 4px 12px rgba(0,0,0,0.03); 
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-name {
  font-size: 14px;
  color: #303133;
}

/* 分页样式 */
.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
  padding: 20px 0;
}

.pagination-container :deep(.el-pagination) {
  padding: 0 20px;
}

.pagination-container :deep(.el-pagination.is-background .el-pager li) {
  border-radius: 6px;
  font-weight: 500;
}
</style>

