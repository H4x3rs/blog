<template>
  <div class="page-container">
    <!-- 筛选区域 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="分类名称">
          <el-input v-model="filters.name" placeholder="请输入分类名称" clearable size="default" style="width: 200px" />
        </el-form-item>
        <el-form-item label="别名">
          <el-input v-model="filters.slug" placeholder="请输入别名" clearable size="default" style="width: 200px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
          <el-button :icon="Refresh" @click="resetFilters">重置</el-button>
        </el-form-item>
        <el-form-item style="margin-left: auto;">
          <el-button type="primary" :icon="Plus" @click="openDialog('create')">新增分类</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格区域 -->
    <el-card shadow="never" class="table-card">
      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="图标" width="80" align="center">
          <template #default="scope">
            <el-icon v-if="scope.row.icon" :size="20">
              <component :is="getIconComponent(scope.row.icon)" />
            </el-icon>
            <span v-else style="color: #ccc;">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="slug" label="别名" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="count" label="文章数" width="100" />
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
      :title="dialogType === 'create' ? '新增分类' : '编辑分类'" 
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
            placeholder="请输入分类名称"
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
        
        <el-form-item label="图标" prop="icon">
          <el-select 
            v-model="form.icon" 
            placeholder="请选择图标" 
            clearable 
            filterable 
            style="width: 100%"
            size="large"
          >
            <el-option
              v-for="icon in iconOptions"
              :key="icon.value"
              :label="icon.label"
              :value="icon.value"
            >
              <div style="display: flex; align-items: center; gap: 8px">
                <el-icon><component :is="icon.component" /></el-icon>
                <span>{{ icon.label }}</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="form.description" 
            type="textarea" 
            :rows="3"
            placeholder="请输入分类描述"
            maxlength="500"
            show-word-limit
            clearable
          />
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
import { ref, reactive, onMounted, watch, markRaw } from 'vue'
import { 
  Plus, Edit, Delete, Search, Refresh,
  Folder, Platform, Monitor, Cpu, Coin, Setting,
  Document, DataAnalysis, Collection, Grid, List,
  Menu, Operation, Tools, Connection, Link, Share,
  Star, Trophy
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCategoryList, createCategory, updateCategory, deleteCategory } from '@/api/category'

// 图标选项 - 使用 markRaw 避免响应式化
const iconOptions = [
  { value: 'Folder', label: '文件夹', component: markRaw(Folder) },
  { value: 'Platform', label: '平台', component: markRaw(Platform) },
  { value: 'Monitor', label: '显示器', component: markRaw(Monitor) },
  { value: 'Cpu', label: 'CPU', component: markRaw(Cpu) },
  { value: 'Coin', label: '硬币', component: markRaw(Coin) },
  { value: 'Setting', label: '设置', component: markRaw(Setting) },
  { value: 'Document', label: '文档', component: markRaw(Document) },
  { value: 'DataAnalysis', label: '数据分析', component: markRaw(DataAnalysis) },
  { value: 'Collection', label: '收藏', component: markRaw(Collection) },
  { value: 'Grid', label: '网格', component: markRaw(Grid) },
  { value: 'List', label: '列表', component: markRaw(List) },
  { value: 'Menu', label: '菜单', component: markRaw(Menu) },
  { value: 'Operation', label: '操作', component: markRaw(Operation) },
  { value: 'Tools', label: '工具', component: markRaw(Tools) },
  { value: 'Connection', label: '连接', component: markRaw(Connection) },
  { value: 'Link', label: '链接', component: markRaw(Link) },
  { value: 'Share', label: '分享', component: markRaw(Share) },
  { value: 'Star', label: '星星', component: markRaw(Star) },
  { value: 'Trophy', label: '奖杯', component: markRaw(Trophy) }
]

// 图标组件映射 - 使用 markRaw 避免响应式化
const iconMap = {
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

// 获取图标组件
const getIconComponent = (iconName) => {
  return iconMap[iconName] || iconMap.Folder
}

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
const form = reactive({ id: null, name: '', slug: '', icon: '', description: '' })

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 1, max: 50, message: '分类名称长度在 1 到 50 个字符', trigger: 'blur' }
  ],
  slug: [
    { required: true, message: '请输入别名', trigger: 'blur' },
    { pattern: /^[a-z0-9-]+$/, message: '别名只能包含小写字母、数字和连字符', trigger: 'blur' },
    { min: 1, max: 100, message: '别名长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 500, message: '描述长度不能超过 500 个字符', trigger: 'blur' }
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
    const res = await getCategoryList(params)
    tableData.value = res.list || []
    totalCount.value = res.total || 0
  } catch (error) {
    ElMessage.error('获取分类列表失败')
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
      icon: row.icon || '',
      description: row.description || ''
    })
  } else {
    Object.assign(form, { id: null, name: '', slug: '', icon: '', description: '' })
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
    const data = {
      name: form.name,
      slug: form.slug,
      icon: form.icon || '',
      description: form.description || ''
    }
    
    if (dialogType.value === 'create') {
      await createCategory(data)
      ElMessage.success('创建成功')
    } else {
      await updateCategory({ id: form.id, ...data })
      ElMessage.success('更新成功')
    }
    dialogVisible.value = false
    // 延迟一下再刷新列表，确保数据库更新完成
    setTimeout(() => {
      fetchList()
    }, 100)
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
    console.error('操作失败:', error)
  } finally {
    submitting.value = false
  }
}

// 删除
const handleDelete = (row) => {
  ElMessageBox.confirm('确认删除该分类？', '提示', { 
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  })
    .then(async () => {
      try {
        await deleteCategory({ id: row.id })
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

/* 选择器样式 */
:deep(.el-select__wrapper) {
  border-radius: 6px;
  padding: 0 12px;
  height: 40px;
  font-size: 14px;
  transition: all 0.3s ease;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

:deep(.el-select__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

:deep(.el-select.is-focused .el-select__wrapper) {
  box-shadow: 0 0 0 1px #409eff inset, 0 0 0 3px rgba(64, 158, 255, 0.1);
}

/* 文本域样式 */
:deep(.el-textarea__inner) {
  border-radius: 6px;
  font-size: 14px;
  padding: 12px;
  line-height: 1.5;
  transition: all 0.3s ease;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

:deep(.el-textarea__inner:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

:deep(.el-textarea.is-focus .el-textarea__inner) {
  box-shadow: 0 0 0 1px #409eff inset, 0 0 0 3px rgba(64, 158, 255, 0.1);
}
</style>
