<template>
  <div class="page-container">
    <!-- 筛选区域 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="权限名称">
          <el-input 
            v-model="filters.name" 
            placeholder="请输入权限名称" 
            clearable 
            :prefix-icon="Search"
            style="width: 260px"
            size="default"
          />
        </el-form-item>
        <el-form-item label="类型">
          <el-select 
            v-model="filters.type" 
            placeholder="请选择类型" 
            clearable
            style="width: 140px"
            size="default"
          >
            <el-option label="菜单" value="menu" />
            <el-option label="按钮" value="btn" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch" size="default">查询</el-button>
          <el-button :icon="Refresh" @click="resetFilters" size="default">重置</el-button>
        </el-form-item>
        <el-form-item style="margin-left: auto;">
          <el-button type="primary" :icon="Plus" @click="openDialog('create')" size="default">
            <span style="margin-left: 4px;">新增权限</span>
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格区域 -->
    <el-card shadow="never" class="table-card">
      <el-table 
        :data="filteredTableData" 
        style="width: 100%" 
        row-key="id" 
        default-expand-all 
        v-loading="loading"
      >
        <el-table-column prop="name" label="权限名称" min-width="200" />
        <el-table-column prop="code" label="权限标识" min-width="220" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100" align="center">
          <template #default="scope">
            <el-tag 
              size="small" 
              :type="scope.row.type === 'menu' ? 'info' : 'warning'"
            >
              {{ scope.row.type === 'menu' ? '菜单' : '按钮' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路由地址" show-overflow-tooltip />
        <el-table-column label="操作" width="220">
          <template #default="scope">
            <el-button link type="primary" size="small" :icon="Plus" @click="openDialog('create', scope.row)" v-if="scope.row.type === 'menu'">新增</el-button>
            <el-button link type="primary" size="small" :icon="Edit" @click="openDialog('edit', scope.row)">编辑</el-button>
            <el-button link type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogType === 'create' ? '新增权限' : '编辑权限'" width="40%">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="标识">
          <el-input v-model="form.code" />
        </el-form-item>
         <el-form-item label="类型">
          <el-radio-group v-model="form.type">
            <el-radio label="menu">菜单</el-radio>
            <el-radio label="btn">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="路径" v-if="form.type === 'menu'">
          <el-input v-model="form.path" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus, Edit, Delete, Search, Refresh } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getPermissionList, createPermission, updatePermission, deletePermission, getPermission } from '@/api/permission'

const loading = ref(false)
const tableData = ref([])
const filters = reactive({
  name: '',
  type: ''
})

const dialogVisible = ref(false)
const dialogType = ref('create')
const form = reactive({ id: null, name: '', code: '', type: 'menu', path: '', parentId: 0 })

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getPermissionList()
    tableData.value = res.list || []
  } catch (error) {
    ElMessage.error('获取权限列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 递归过滤树形数据
const filterTreeData = (data, nameFilter, typeFilter) => {
  return data.filter(item => {
    // 检查当前项是否匹配
    const nameMatch = !nameFilter || item.name?.toLowerCase().includes(nameFilter.toLowerCase())
    const typeMatch = !typeFilter || item.type === typeFilter
    
    // 递归过滤子项
    let filteredChildren = []
    if (item.children && item.children.length > 0) {
      filteredChildren = filterTreeData(item.children, nameFilter, typeFilter)
    }
    
    // 如果当前项匹配，或者有匹配的子项，则保留
    return (nameMatch && typeMatch) || filteredChildren.length > 0
  }).map(item => {
    // 递归处理子项
    if (item.children && item.children.length > 0) {
      const filteredChildren = filterTreeData(item.children, nameFilter, typeFilter)
      return {
        ...item,
        children: filteredChildren
      }
    }
    return item
  })
}

// 过滤后的表格数据
const filteredTableData = computed(() => {
  if (!filters.name && !filters.type) {
    return tableData.value
  }
  return filterTreeData(tableData.value, filters.name, filters.type)
})

const handleSearch = () => {
  // 过滤逻辑已在 computed 中处理
}

const resetFilters = () => {
  filters.name = ''
  filters.type = ''
}

const openDialog = async (type, row) => {
  dialogType.value = type
  dialogVisible.value = true
  if (type === 'edit' && row) {
    try {
      const res = await getPermission({ id: row.id })
      Object.assign(form, {
        id: res.id,
        name: res.name,
        code: res.code,
        type: res.type,
        path: res.path || '',
        parentId: res.parentId || 0
      })
    } catch (error) {
      ElMessage.error('获取权限信息失败')
      console.error(error)
    }
  } else if (type === 'create' && row) {
    // 创建子权限
    Object.assign(form, { id: null, name: '', code: '', type: 'btn', path: '', parentId: row.id })
  } else {
    Object.assign(form, { id: null, name: '', code: '', type: 'menu', path: '', parentId: 0 })
  }
}

const handleSubmit = async () => {
  if (!form.name || !form.code) {
    ElMessage.warning('请填写完整信息')
    return
  }
  try {
    if (dialogType.value === 'create') {
      await createPermission({ 
        name: form.name, 
        code: form.code, 
        type: form.type, 
        path: form.path,
        parentId: form.parentId || 0
      })
      ElMessage.success('创建成功')
    } else {
      await updatePermission({ 
        id: form.id,
        name: form.name, 
        code: form.code, 
        type: form.type, 
        path: form.path,
        parentId: form.parentId || 0
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

const handleDelete = (row) => {
  ElMessageBox.confirm('确认删除？', '提示', { type: 'warning' })
    .then(async () => {
      try {
        await deletePermission({ id: row.id })
        ElMessage.success('删除成功')
        fetchList()
      } catch (error) {
        ElMessage.error('删除失败')
        console.error(error)
      }
    })
    .catch(() => {})
}

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
  background: linear-gradient(to bottom, #ffffff 0%, #fafbfc 100%);
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}

.filter-form .el-form-item {
  margin-bottom: 0;
  margin-right: 20px;
}

/* 搜索栏输入框和按钮统一样式 */
.filter-form :deep(.el-input__wrapper) {
  height: 32px;
}

.filter-form :deep(.el-button) {
  height: 32px;
  padding: 8px 15px;
}

.filter-form :deep(.el-select .el-input__wrapper) {
  height: 32px;
}

.table-card { 
  border-radius: 16px; 
  border: none; 
  box-shadow: 0 4px 12px rgba(0,0,0,0.03); 
}
</style>

