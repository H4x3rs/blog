<template>
  <div class="page-container">
    <!-- 筛选区域 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="操作类型">
          <el-select 
            v-model="filters.operationType" 
            placeholder="请选择操作类型" 
            clearable 
            style="width: 200px"
            size="default"
          >
            <el-option label="登录" value="login" />
            <el-option label="创建" value="create" />
            <el-option label="更新" value="update" />
            <el-option label="删除" value="delete" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作模块">
          <el-select 
            v-model="filters.module" 
            placeholder="请选择操作模块" 
            clearable 
            style="width: 200px"
            size="default"
          >
            <el-option label="用户" value="user" />
            <el-option label="文章" value="article" />
            <el-option label="分类" value="category" />
            <el-option label="标签" value="tag" />
            <el-option label="角色" value="role" />
            <el-option label="权限" value="permission" />
            <el-option label="OAuth" value="oauth" />
          </el-select>
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
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="operationType" label="操作类型" width="100">
          <template #default="scope">
            <el-tag :type="getOperationTypeTag(scope.row.operationType)">
              {{ getOperationTypeText(scope.row.operationType) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="module" label="操作模块" width="100">
          <template #default="scope">
            <el-tag>{{ getModuleText(scope.row.module) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="operationDesc" label="操作描述" min-width="150" />
        <el-table-column prop="requestMethod" label="请求方法" width="100" />
        <el-table-column prop="requestPath" label="请求路径" min-width="200" />
        <el-table-column prop="ipAddress" label="IP地址" width="140" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="操作时间" width="180" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button link type="primary" size="small" @click="viewDetail(scope.row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container" v-if="totalCount > pageSize">
        <el-pagination 
          background 
          layout="total, sizes, prev, pager, next, jumper" 
          :total="totalCount" 
          :page-sizes="[10, 20, 30, 50]"
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog 
      v-model="detailDialogVisible" 
      title="操作日志详情" 
      width="800px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <el-descriptions :column="2" border v-if="currentDetail">
        <el-descriptions-item label="ID">{{ currentDetail.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ currentDetail.username || '-' }}</el-descriptions-item>
        <el-descriptions-item label="操作类型">
          <el-tag :type="getOperationTypeTag(currentDetail.operationType)">
            {{ getOperationTypeText(currentDetail.operationType) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="操作模块">
          <el-tag>{{ getModuleText(currentDetail.module) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="操作描述" :span="2">{{ currentDetail.operationDesc }}</el-descriptions-item>
        <el-descriptions-item label="请求方法">{{ currentDetail.requestMethod }}</el-descriptions-item>
        <el-descriptions-item label="请求路径">{{ currentDetail.requestPath }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ currentDetail.ipAddress }}</el-descriptions-item>
        <el-descriptions-item label="操作状态">
          <el-tag :type="currentDetail.status === 1 ? 'success' : 'danger'">
            {{ currentDetail.status === 1 ? '成功' : '失败' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="操作时间" :span="2">{{ currentDetail.createdAt }}</el-descriptions-item>
        <el-descriptions-item label="用户代理" :span="2">
          <el-text truncated style="max-width: 100%">{{ currentDetail.userAgent || '-' }}</el-text>
        </el-descriptions-item>
        <el-descriptions-item label="请求参数" :span="2">
          <pre v-if="currentDetail.requestParams" style="max-height: 200px; overflow: auto; background: #f5f5f5; padding: 10px; border-radius: 4px;">{{ formatJson(currentDetail.requestParams) }}</pre>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="错误信息" :span="2" v-if="currentDetail.errorMessage">
          <el-text type="danger">{{ currentDetail.errorMessage }}</el-text>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import { getOperationLogList } from '@/api/operationLog'

const loading = ref(false)
const tableData = ref([])
const totalCount = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const detailDialogVisible = ref(false)
const currentDetail = ref(null)

const filters = reactive({
  operationType: '',
  module: ''
})

// 获取操作类型文本
const getOperationTypeText = (type) => {
  const map = {
    'login': '登录',
    'create': '创建',
    'update': '更新',
    'delete': '删除'
  }
  return map[type] || type
}

// 获取操作类型标签
const getOperationTypeTag = (type) => {
  const map = {
    'login': 'primary',
    'create': 'success',
    'update': 'warning',
    'delete': 'danger'
  }
  return map[type] || ''
}

// 获取模块文本
const getModuleText = (module) => {
  const map = {
    'user': '用户',
    'article': '文章',
    'category': '分类',
    'tag': '标签',
    'role': '角色',
    'permission': '权限',
    'oauth': 'OAuth'
  }
  return map[module] || module
}

// 格式化JSON
const formatJson = (jsonStr) => {
  if (!jsonStr) return ''
  try {
    const obj = JSON.parse(jsonStr)
    return JSON.stringify(obj, null, 2)
  } catch (e) {
    return jsonStr
  }
}

// 获取列表数据
const fetchList = async () => {
  loading.value = true
  try {
    const res = await getOperationLogList({
      userId: 0,
      operationType: filters.operationType || undefined,
      module: filters.module || undefined,
      page: currentPage.value,
      size: pageSize.value
    })
    if (res && res.list) {
      tableData.value = res.list
      totalCount.value = res.total || 0
    }
  } catch (error) {
    ElMessage.error('获取日志列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

// 重置筛选
const resetFilters = () => {
  filters.operationType = ''
  filters.module = ''
  handleSearch()
}


// 查看详情
const viewDetail = (row) => {
  currentDetail.value = row
  detailDialogVisible.value = true
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

.filter-form :deep(.el-select__wrapper) {
  height: 32px;
  border-radius: 4px;
}

.table-card { 
  border-radius: 16px; 
  border: none; 
  box-shadow: 0 4px 12px rgba(0,0,0,0.03); 
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

/* 对话框样式优化 */
:deep(.el-dialog__header) {
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-dialog__footer) {
  padding: 16px 24px;
  border-top: 1px solid #f0f0f0;
}

/* 描述列表样式 */
:deep(.el-descriptions__label) {
  font-weight: 500;
  color: #606266;
}

:deep(.el-descriptions__content) {
  color: #303133;
}

pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.6;
}
</style>

