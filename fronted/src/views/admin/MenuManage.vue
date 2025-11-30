<template>
  <div class="page-container">
    <!-- 筛选区域 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="菜单名称">
          <el-input 
            v-model="filters.title" 
            placeholder="请输入菜单名称" 
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
            <el-option label="目录" value="M" />
            <el-option label="菜单" value="C" />
            <el-option label="按钮" value="F" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select 
            v-model="filters.status" 
            placeholder="请选择状态" 
            clearable
            style="width: 140px"
            size="default"
          >
            <el-option label="启用" value="active" />
            <el-option label="禁用" value="inactive" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch" size="default">查询</el-button>
          <el-button :icon="Refresh" @click="resetFilters" size="default">重置</el-button>
        </el-form-item>
        <el-form-item style="margin-left: auto;">
          <el-button type="primary" :icon="Plus" @click="openDialog('create')" size="default">
            <span style="margin-left: 4px;">新增菜单</span>
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
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
      >
        <el-table-column prop="title" label="菜单名称" min-width="200" />
        <el-table-column prop="type" label="类型" width="100" align="center">
          <template #default="scope">
            <el-tag 
              :type="scope.row.type === 'M' ? 'info' : scope.row.type === 'C' ? 'primary' : 'warning'"
              size="small"
            >
              {{ scope.row.type === 'M' ? '目录' : scope.row.type === 'C' ? '菜单' : '按钮' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="icon" label="图标" width="80" align="center">
          <template #default="scope">
            <el-icon v-if="scope.row.icon"><component :is="scope.row.icon" /></el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="order" label="排序" width="80" />
        <el-table-column prop="path" label="路由地址" show-overflow-tooltip />
        <el-table-column prop="component" label="组件路径" show-overflow-tooltip />
        <el-table-column prop="permission" label="权限字符" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
           <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'info'">{{ scope.row.status === 'active' ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button link type="primary" size="small" :icon="Plus" @click="openDialog('create', scope.row)">新增</el-button>
            <el-button link type="primary" size="small" :icon="Edit" @click="openDialog('edit', scope.row)">编辑</el-button>
            <el-button link type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogType === 'create' ? '新增菜单' : '编辑菜单'" width="40%">
      <el-form :model="form" label-width="100px">
        <el-form-item label="上级菜单">
          <el-tree-select
            v-model="form.parentId"
            :data="menuOptions"
            check-strictly
            :render-after-expand="false"
            placeholder="选择上级菜单"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="菜单类型">
           <el-radio-group v-model="form.type">
             <el-radio label="M">目录</el-radio>
             <el-radio label="C">菜单</el-radio>
             <el-radio label="F">按钮</el-radio>
           </el-radio-group>
        </el-form-item>
        <el-form-item label="菜单名称">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="图标" v-if="form.type !== 'F'">
          <el-input v-model="form.icon" placeholder="Element Plus Icon Name" />
        </el-form-item>
        <el-form-item label="路由地址" v-if="form.type !== 'F'">
          <el-input v-model="form.path" />
        </el-form-item>
        <el-form-item label="组件路径" v-if="form.type === 'C'">
          <el-input v-model="form.component" />
        </el-form-item>
        <el-form-item label="权限字符" v-if="form.type !== 'M'">
          <el-input v-model="form.permission" />
        </el-form-item>
        <el-form-item label="显示排序">
          <el-input-number v-model="form.order" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
           <el-radio-group v-model="form.status">
             <el-radio label="active">启用</el-radio>
             <el-radio label="inactive">禁用</el-radio>
           </el-radio-group>
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
import { Plus, Edit, Delete, Search, Refresh, Setting, User, Document, Odometer } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getMenuList, createMenu, updateMenu, deleteMenu, getMenu } from '@/api/menu'

const loading = ref(false)
const tableData = ref([])
const menuOptions = ref([{ value: 0, label: '主类目' }])
const filters = reactive({
  title: '',
  type: '',
  status: ''
})

const dialogVisible = ref(false)
const dialogType = ref('create')
const form = reactive({
  id: null,
  parentId: 0,
  type: 'C',
  title: '',
  icon: '',
  path: '',
  component: '',
  permission: '',
  order: 0,
  status: 'active'
})

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getMenuList()
    tableData.value = res.list || []
    // 构建菜单选项树（扁平化处理）
    const flattenMenus = (menus, result = []) => {
      menus.forEach(menu => {
        result.push(menu)
        if (menu.children && menu.children.length > 0) {
          flattenMenus(menu.children, result)
        }
      })
      return result
    }
    const flatList = flattenMenus(res.list || [])
    buildMenuOptions(flatList)
  } catch (error) {
    ElMessage.error('获取菜单列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const buildMenuOptions = (list, parentId = 0, prefix = '') => {
  if (parentId === 0) {
    menuOptions.value = [{ value: 0, label: '主类目' }]
  }
  list.forEach(item => {
    if (item.parentId === parentId) {
      const label = prefix + item.title
      menuOptions.value.push({ value: item.id, label })
      if (item.children && item.children.length > 0) {
        buildMenuOptions(item.children, item.id, prefix + '  ')
      }
    }
  })
}

// 递归过滤树形数据
const filterTreeData = (data, titleFilter, typeFilter, statusFilter) => {
  return data.filter(item => {
    // 检查当前项是否匹配
    const titleMatch = !titleFilter || item.title?.toLowerCase().includes(titleFilter.toLowerCase())
    const typeMatch = !typeFilter || item.type === typeFilter
    const statusMatch = !statusFilter || item.status === statusFilter
    
    // 递归过滤子项
    let filteredChildren = []
    if (item.children && item.children.length > 0) {
      filteredChildren = filterTreeData(item.children, titleFilter, typeFilter, statusFilter)
    }
    
    // 如果当前项匹配，或者有匹配的子项，则保留
    return (titleMatch && typeMatch && statusMatch) || filteredChildren.length > 0
  }).map(item => {
    // 递归处理子项
    if (item.children && item.children.length > 0) {
      const filteredChildren = filterTreeData(item.children, titleFilter, typeFilter, statusFilter)
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
  if (!filters.title && !filters.type && !filters.status) {
    return tableData.value
  }
  return filterTreeData(tableData.value, filters.title, filters.type, filters.status)
})

const handleSearch = () => {
  // 过滤逻辑已在 computed 中处理
}

const resetFilters = () => {
  filters.title = ''
  filters.type = ''
  filters.status = ''
}

const openDialog = async (type, row) => {
  dialogType.value = type
  dialogVisible.value = true
  if (type === 'create' && row) {
     Object.assign(form, { 
        id: null, parentId: row.id, type: 'C', title: '', icon: '', 
        path: '', component: '', permission: '', order: 0, status: 'active' 
     })
  } else if (type === 'edit' && row) {
    try {
      const res = await getMenu({ id: row.id })
      Object.assign(form, res)
    } catch (error) {
      ElMessage.error('获取菜单信息失败')
      console.error(error)
    }
  } else {
    Object.assign(form, { 
        id: null, parentId: 0, type: 'M', title: '', icon: '', 
        path: '', component: '', permission: '', order: 0, status: 'active' 
     })
  }
}

const handleSubmit = async () => {
  if (!form.title) {
    ElMessage.warning('请填写菜单名称')
    return
  }
  try {
    if (dialogType.value === 'create') {
      await createMenu({
        parentId: form.parentId,
        type: form.type,
        title: form.title,
        icon: form.icon,
        path: form.path,
        component: form.component,
        permission: form.permission,
        order: form.order,
        status: form.status
      })
      ElMessage.success('创建成功')
    } else {
      await updateMenu({
        id: form.id,
        parentId: form.parentId,
        type: form.type,
        title: form.title,
        icon: form.icon,
        path: form.path,
        component: form.component,
        permission: form.permission,
        order: form.order,
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

const handleDelete = (row) => {
  ElMessageBox.confirm('确认删除？', '提示', { type: 'warning' })
    .then(async () => {
      try {
        await deleteMenu({ id: row.id })
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

