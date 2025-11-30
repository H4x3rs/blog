<template>
  <div class="page-container">
    <!-- 筛选区域 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="标签名称">
          <el-input v-model="filters.name" placeholder="请输入标签名称" clearable size="default" style="width: 200px" />
        </el-form-item>
        <el-form-item label="别名">
          <el-input v-model="filters.slug" placeholder="请输入别名" clearable size="default" style="width: 200px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
          <el-button :icon="Refresh" @click="resetFilters">重置</el-button>
        </el-form-item>
        <el-form-item style="margin-left: auto;">
          <el-button type="primary" :icon="Plus" @click="openDialog('create')">新增标签</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格区域 -->
    <el-card shadow="never" class="table-card">
      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="slug" label="别名" min-width="150" show-overflow-tooltip />
        <el-table-column prop="color" label="颜色">
          <template #default="scope">
            <el-tag :color="scope.row.color" effect="dark" style="border: none">{{ scope.row.color }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <el-button link type="primary" size="small" :icon="Edit" @click="openDialog('edit', scope.row)">编辑</el-button>
            <el-button link type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)">删除</el-button>
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

    <el-dialog 
      v-model="dialogVisible" 
      :title="dialogType === 'create' ? '新增标签' : '编辑标签'" 
      width="600px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <el-form 
        ref="formRef"
        :model="form" 
        :rules="formRules"
        label-width="100px"
        label-position="right"
      >
        <el-form-item label="名称" prop="name">
          <el-input 
            v-model="form.name" 
            placeholder="请输入标签名称"
            clearable
            size="large"
          >
            <template #prefix>
              <el-icon><EditPen /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item label="别名" prop="slug">
          <el-input 
            v-model="form.slug" 
            placeholder="请输入别名（英文，用于URL）"
            clearable
            size="large"
          >
            <template #prefix>
              <el-icon><Link /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item label="颜色" prop="color">
          <div class="color-picker-wrapper">
            <el-color-picker 
              v-model="form.color" 
              :predefine="predefineColors"
              show-alpha
            />
            <div class="color-preview">
              <el-tag 
                :color="form.color" 
                effect="dark" 
                style="border: none; margin-left: 12px;"
              >
                {{ form.color }}
              </el-tag>
            </div>
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false" size="large">取消</el-button>
          <el-button 
            type="primary" 
            @click="handleSubmit" 
            size="large"
            :loading="submitting"
          >
            {{ dialogType === 'create' ? '创建' : '保存' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { 
  Plus, Edit, Delete, Search, Refresh,
  EditPen, Link, Brush
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getTagList, createTag, updateTag, deleteTag } from '@/api/tag'

// 筛选条件
const filters = reactive({
  name: '',
  slug: ''
})

// 分页
const currentPage = ref(1)
const pageSize = ref(10)
const totalCount = ref(0)

const loading = ref(false)
const tableData = ref([])

const dialogVisible = ref(false)
const dialogType = ref('create')
const formRef = ref(null)
const submitting = ref(false)
const form = reactive({ id: null, name: '', slug: '', color: '#409EFF' })

// 预定义颜色
const predefineColors = [
  '#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399',
  '#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399',
  '#FF4500', '#FF8C00', '#FFD700', '#ADFF2F', '#00CED1',
  '#1E90FF', '#9370DB', '#FF1493', '#00FA9A', '#FF69B4'
]

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入标签名称', trigger: 'blur' },
    { min: 1, max: 50, message: '标签名称长度在 1 到 50 个字符', trigger: 'blur' }
  ],
  slug: [
    { required: true, message: '请输入别名', trigger: 'blur' },
    { pattern: /^[a-z0-9-]+$/, message: '别名只能包含小写字母、数字和连字符', trigger: 'blur' },
    { min: 1, max: 100, message: '别名长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  color: [
    { required: true, message: '请选择颜色', trigger: 'change' }
  ]
}

// 获取列表
const fetchList = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value,
      name: filters.name,
      slug: filters.slug
    }
    const res = await getTagList(params)
    tableData.value = res.list || []
    totalCount.value = res.total || 0
  } catch (error) {
    ElMessage.error('获取标签列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 查询
const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

// 重置
const resetFilters = () => {
  filters.name = ''
  filters.slug = ''
  currentPage.value = 1
  fetchList()
}

// 打开对话框
const openDialog = (type, row) => {
  dialogType.value = type
  dialogVisible.value = true
  
  // 重置表单验证状态
  if (formRef.value) {
    formRef.value.clearValidate()
  }
  
  if (type === 'edit' && row) {
    Object.assign(form, {
      id: row.id,
      name: row.name || '',
      slug: row.slug || '',
      color: row.color || '#409EFF'
    })
  } else {
    Object.assign(form, { id: null, name: '', slug: '', color: '#409EFF' })
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  // 验证表单
  try {
    await formRef.value.validate()
  } catch (error) {
    return
  }
  
  submitting.value = true
  try {
    if (dialogType.value === 'create') {
      await createTag({ name: form.name, slug: form.slug, color: form.color })
      ElMessage.success('创建成功')
    } else {
      await updateTag({ id: form.id, name: form.name, slug: form.slug, color: form.color })
      ElMessage.success('更新成功')
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
    console.error(error)
  } finally {
    submitting.value = false
  }
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm('确认删除该标签？', '提示', { 
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  })
    .then(async () => {
      try {
        await deleteTag({ id: row.id })
        ElMessage.success('删除成功')
        fetchList()
      } catch (error) {
        ElMessage.error('删除失败')
        console.error(error)
      }
    })
    .catch(() => {})
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

.table-card { 
  border-radius: 16px; 
  border: none; 
  box-shadow: 0 4px 12px rgba(0,0,0,0.03); 
}

.card-header { 
  display: flex; 
  justify-content: space-between; 
  align-items: center; 
  font-weight: 600; 
}

/* 分页样式 */
.pagination-container {
  display: flex;
  justify-content: center;
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

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 表单样式统一 */
:deep(.el-form-item__label) {
  font-weight: 500;
  color: #606266;
  font-size: 14px;
}

:deep(.el-form-item) {
  margin-bottom: 22px;
}

/* 输入框统一样式 */
:deep(.el-input__wrapper) {
  border-radius: 6px;
  padding: 0 12px;
  height: 40px;
  font-size: 14px;
  transition: all 0.3s ease;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

:deep(.el-input--large .el-input__wrapper) {
  height: 40px;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

:deep(.el-input.is-focus .el-input__wrapper) {
  box-shadow: 0 0 0 1px #409eff inset, 0 0 0 3px rgba(64, 158, 255, 0.1);
}

:deep(.el-input__inner) {
  font-size: 14px;
  color: #606266;
  line-height: 1.5;
}

:deep(.el-input__prefix) {
  color: #909399;
  margin-right: 8px;
}

:deep(.el-input.is-focus .el-input__prefix) {
  color: #409eff;
}

/* 颜色选择器样式 */
.color-picker-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.color-preview {
  display: flex;
  align-items: center;
}

:deep(.el-color-picker) {
  height: 40px;
}

:deep(.el-color-picker__trigger) {
  width: 40px;
  height: 40px;
  border-radius: 6px;
  border: 1px solid #dcdfe6;
  transition: all 0.3s ease;
}

:deep(.el-color-picker__trigger:hover) {
  border-color: #409eff;
  box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
}
</style>


