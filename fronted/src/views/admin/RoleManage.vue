<template>
  <div class="page-container">
    <!-- 筛选区域 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="角色名称">
          <el-input 
            v-model="filters.name" 
            placeholder="请输入角色名称" 
            clearable 
            :prefix-icon="Search"
            style="width: 260px"
            size="default"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch" size="default">查询</el-button>
          <el-button :icon="Refresh" @click="resetFilters" size="default">重置</el-button>
        </el-form-item>
        <el-form-item style="margin-left: auto;">
          <el-button type="primary" :icon="Plus" @click="openDialog('create')" size="default">
            <span style="margin-left: 4px;">新增角色</span>
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格区域 -->
    <el-card shadow="never" class="table-card">
      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" />
        <el-table-column prop="key" label="权限字符" />
        <el-table-column prop="desc" label="描述" />
        <el-table-column prop="createdAt" label="创建时间" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="scope">
            <el-button link type="primary" size="small" :icon="Edit" @click="openDialog('edit', scope.row)">编辑</el-button>
            <el-button link type="success" size="small" :icon="Key" @click="openPermissionDialog(scope.row)">权限</el-button>
            <el-button link type="warning" size="small" :icon="Menu" @click="openMenuDialog(scope.row)">菜单</el-button>
            <el-button link type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogType === 'create' ? '新增角色' : '编辑角色'" width="40%">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="权限字符">
          <el-input v-model="form.key" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.desc" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>

    <!-- 权限授权对话框 -->
    <el-dialog 
      v-model="permissionDialogVisible" 
      title="分配权限" 
      width="600px"
      :close-on-click-modal="false"
    >
      <el-tree
        ref="permissionTreeRef"
        :data="permissionTree"
        :props="{ label: 'name', children: 'children' }"
        show-checkbox
        node-key="id"
        :default-checked-keys="selectedPermissionIds"
        style="max-height: 400px; overflow-y: auto;"
      >
        <template #default="{ node, data }">
          <div style="display: flex; align-items: center; width: 100%;">
            <el-icon v-if="data.icon" style="margin-right: 8px; color: #409eff;">
              <component :is="data.icon" />
            </el-icon>
            <span style="flex: 1;">{{ data.name }}</span>
            <el-tag size="small" type="info" style="margin-left: 8px;">{{ data.code }}</el-tag>
            <el-tag 
              size="small" 
              :type="data.type === 'menu' ? 'primary' : 'warning'"
              style="margin-left: 8px;"
            >
              {{ data.type === 'menu' ? '菜单' : '按钮' }}
            </el-tag>
          </div>
        </template>
      </el-tree>
      
      <template #footer>
        <el-button @click="permissionDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSavePermissions" :loading="savingPermissions">保存</el-button>
      </template>
    </el-dialog>

    <!-- 菜单授权对话框 -->
    <el-dialog 
      v-model="menuDialogVisible" 
      title="分配菜单" 
      width="600px"
      :close-on-click-modal="false"
    >
      <el-tree
        ref="menuTreeRef"
        :data="menuTree"
        :props="{ label: 'title', children: 'children' }"
        show-checkbox
        node-key="id"
        :default-checked-keys="selectedMenuIds"
        style="max-height: 400px; overflow-y: auto;"
      >
        <template #default="{ node, data }">
          <div style="display: flex; align-items: center; width: 100%;">
            <el-icon v-if="data.icon" style="margin-right: 8px; color: #409eff;">
              <component :is="data.icon" />
            </el-icon>
            <span style="flex: 1;">{{ data.title }}</span>
            <el-tag 
              size="small" 
              :type="data.type === 'M' ? 'info' : (data.type === 'C' ? 'success' : 'warning')"
              style="margin-left: 8px;"
            >
              {{ data.type === 'M' ? '目录' : (data.type === 'C' ? '菜单' : '按钮') }}
            </el-tag>
            <el-tag size="small" type="info" style="margin-left: 8px;" v-if="data.path">{{ data.path }}</el-tag>
          </div>
        </template>
      </el-tree>
      
      <template #footer>
        <el-button @click="menuDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveMenus" :loading="savingMenus">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { Plus, Edit, Delete, Search, Refresh, Key, Menu } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getRoleList, createRole, updateRole, deleteRole, getRole, setRolePermissions, getRolePermissions } from '@/api/role'
import { getPermissionList } from '@/api/permission'
import { getMenuList } from '@/api/menu'

const loading = ref(false)
const tableData = ref([])
const filters = reactive({
  name: ''
})

const dialogVisible = ref(false)
const dialogType = ref('create')
const form = reactive({ id: null, name: '', key: '', desc: '' })

// 权限授权相关
const permissionDialogVisible = ref(false)
const permissionTree = ref([])
const permissionTreeRef = ref(null)
const selectedPermissionIds = ref([])
const currentRoleId = ref(null)
const savingPermissions = ref(false)

// 菜单授权相关
const menuDialogVisible = ref(false)
const menuTree = ref([])
const menuTreeRef = ref(null)
const selectedMenuIds = ref([])
const savingMenus = ref(false)

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

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getRoleList({ page: 1, size: 100, name: filters.name })
    tableData.value = (res.list || []).map(item => ({
      ...item,
      createdAt: formatDate(item.createdAt)
    }))
  } catch (error) {
    ElMessage.error('获取角色列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  fetchList()
}

const resetFilters = () => {
  filters.name = ''
  fetchList()
}

const openDialog = async (type, row) => {
  dialogType.value = type
  dialogVisible.value = true
  if (type === 'edit' && row) {
    try {
      const res = await getRole({ id: row.id })
      Object.assign(form, res)
    } catch (error) {
      ElMessage.error('获取角色信息失败')
      console.error(error)
    }
  } else {
    Object.assign(form, { id: null, name: '', key: '', desc: '' })
  }
}

const handleSubmit = async () => {
  if (!form.name || !form.key) {
    ElMessage.warning('请填写完整信息')
    return
  }
  try {
    if (dialogType.value === 'create') {
      await createRole({ name: form.name, key: form.key, desc: form.desc })
      ElMessage.success('创建成功')
    } else {
      await updateRole({ id: form.id, name: form.name, key: form.key, desc: form.desc })
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
        await deleteRole({ id: row.id })
        ElMessage.success('删除成功')
        fetchList()
      } catch (error) {
        ElMessage.error('删除失败')
        console.error(error)
      }
    })
    .catch(() => {})
}

// 构建权限树（已废弃，后端直接返回树形结构）
// const buildPermissionTree = (permissions) => {
//   const map = {}
//   const roots = []
//   
//   // 创建映射
//   permissions.forEach(perm => {
//     map[perm.id] = { ...perm, children: [] }
//   })
//   
//   // 构建树
//   permissions.forEach(perm => {
//     if (perm.parentId === 0) {
//       roots.push(map[perm.id])
//     } else {
//       if (map[perm.parentId]) {
//         map[perm.parentId].children.push(map[perm.id])
//       }
//     }
//   })
//   
//   return roots
// }

const openPermissionDialog = async (row) => {
  currentRoleId.value = row.id
  permissionDialogVisible.value = true
  selectedPermissionIds.value = []
  
  // 加载所有权限（返回的是树形结构）
  try {
    const res = await getPermissionList({})
    // 权限列表已经是树形结构，直接使用
    permissionTree.value = res.list || []
  } catch (error) {
    ElMessage.error('获取权限列表失败')
    console.error(error)
    return
  }
  
  // 加载角色已有权限
  try {
    const res = await getRolePermissions({ id: row.id })
    selectedPermissionIds.value = res.permissionIds || []
    
    // 等待树渲染完成后设置选中状态
    await nextTick()
    if (permissionTreeRef.value) {
      permissionTreeRef.value.setCheckedKeys(selectedPermissionIds.value)
    }
  } catch (error) {
    ElMessage.error('获取角色权限失败')
    console.error(error)
  }
}

const handleSavePermissions = async () => {
  if (!currentRoleId.value) return
  
  const checkedKeys = permissionTreeRef.value.getCheckedKeys()
  const halfCheckedKeys = permissionTreeRef.value.getHalfCheckedKeys()
  const allCheckedKeys = [...checkedKeys, ...halfCheckedKeys]
  
  savingPermissions.value = true
  try {
    await setRolePermissions({
      id: currentRoleId.value,
      permissionIds: allCheckedKeys
    })
    ElMessage.success('权限分配成功')
    permissionDialogVisible.value = false
  } catch (error) {
    ElMessage.error('权限分配失败')
    console.error(error)
  } finally {
    savingPermissions.value = false
  }
}

// 收集菜单对应的权限ID（递归遍历树）
const collectMenuPermissionIds = (menus) => {
  const permissionIds = []
  const traverse = (items) => {
    items.forEach(item => {
      if (item.permission) {
        // 通过权限标识查找权限ID
        const permission = allPermissions.value.find(p => p.code === item.permission)
        if (permission && !permissionIds.includes(permission.id)) {
          permissionIds.push(permission.id)
        }
      }
      if (item.children && item.children.length > 0) {
        traverse(item.children)
      }
    })
  }
  traverse(menus)
  return permissionIds
}

// 从权限ID反推菜单ID（递归遍历树）
const getMenuIdsFromPermissionIds = (permissionIds) => {
  const menuIds = []
  const traverse = (items) => {
    items.forEach(item => {
      let matched = false
      
      // 如果菜单有关联的权限标识
      if (item.permission) {
        const permission = allPermissions.value.find(p => p.code === item.permission)
        if (permission && permissionIds.includes(permission.id)) {
          menuIds.push(item.id)
          matched = true
        }
      }
      
      // 如果通过permission没找到，尝试通过path查找
      if (!matched && item.path) {
        const pathPermission = allPermissions.value.find(p => p.path === item.path && p.type === 'menu')
        if (pathPermission && permissionIds.includes(pathPermission.id)) {
          menuIds.push(item.id)
        }
      }
      
      if (item.children && item.children.length > 0) {
        traverse(item.children)
      }
    })
  }
  traverse(menuTree.value)
  return menuIds
}

const allPermissions = ref([])

const openMenuDialog = async (row) => {
  currentRoleId.value = row.id
  menuDialogVisible.value = true
  selectedMenuIds.value = []
  
  // 先加载所有权限（用于菜单和权限的映射）
  try {
    const res = await getPermissionList({})
    // 权限列表返回的是树形结构，需要扁平化
    const flattenPermissions = (items) => {
      const result = []
      items.forEach(item => {
        result.push(item)
        if (item.children && item.children.length > 0) {
          result.push(...flattenPermissions(item.children))
        }
      })
      return result
    }
    allPermissions.value = flattenPermissions(res.list || [])
    console.log('加载的权限列表:', allPermissions.value)
  } catch (error) {
    ElMessage.error('获取权限列表失败')
    console.error(error)
    return
  }
  
  // 加载所有菜单
  try {
    const res = await getMenuList({})
    menuTree.value = res.list || []
    console.log('加载的菜单列表:', menuTree.value)
  } catch (error) {
    ElMessage.error('获取菜单列表失败')
    console.error(error)
    return
  }
  
  // 加载角色已有权限，并转换为菜单ID
  try {
    const res = await getRolePermissions({ id: row.id })
    const permissionIds = res.permissionIds || []
    console.log('角色已有权限ID:', permissionIds)
    selectedMenuIds.value = getMenuIdsFromPermissionIds(permissionIds)
    console.log('转换后的菜单ID:', selectedMenuIds.value)
    
    // 等待树渲染完成后设置选中状态
    await nextTick()
    if (menuTreeRef.value) {
      menuTreeRef.value.setCheckedKeys(selectedMenuIds.value)
    }
  } catch (error) {
    ElMessage.error('获取角色权限失败')
    console.error(error)
  }
}

const handleSaveMenus = async () => {
  if (!currentRoleId.value || !menuTreeRef.value) {
    ElMessage.warning('请先选择菜单')
    return
  }
  
  try {
    const checkedKeys = menuTreeRef.value.getCheckedKeys()
    const halfCheckedKeys = menuTreeRef.value.getHalfCheckedKeys()
    
    console.log('选中的菜单ID:', checkedKeys)
    console.log('半选中的菜单ID:', halfCheckedKeys)
    console.log('所有权限列表:', allPermissions.value)
    console.log('菜单树:', menuTree.value)
    
    // 如果没有选中任何菜单
    if (checkedKeys.length === 0 && halfCheckedKeys.length === 0) {
      ElMessage.warning('请至少选择一个菜单')
      return
    }
    
    // 收集所有选中和半选中菜单的权限ID
    const permissionIds = []
    
    // 递归遍历菜单树，收集权限ID
    const collectPermissions = (items) => {
      items.forEach(item => {
        // 如果菜单被选中或半选中
        if (checkedKeys.includes(item.id) || halfCheckedKeys.includes(item.id)) {
          console.log(`处理菜单: ${item.title}, ID: ${item.id}, permission: ${item.permission}, path: ${item.path}`)
          
          // 如果菜单有关联的权限，则通过权限标识查找权限ID
          if (item.permission) {
            const permission = allPermissions.value.find(p => p.code === item.permission)
            if (permission) {
              console.log(`找到权限: ${permission.name}, ID: ${permission.id}`)
              if (!permissionIds.includes(permission.id)) {
                permissionIds.push(permission.id)
              }
            } else {
              console.warn(`菜单 ${item.title} 的权限标识 ${item.permission} 未找到对应的权限`)
              console.warn('可用权限代码:', allPermissions.value.map(p => p.code))
            }
          } else {
            // 如果菜单没有permission字段，尝试通过path查找对应的权限
            const pathPermission = allPermissions.value.find(p => p.path === item.path && p.type === 'menu')
            if (pathPermission) {
              console.log(`通过path找到权限: ${pathPermission.name}, ID: ${pathPermission.id}`)
              if (!permissionIds.includes(pathPermission.id)) {
                permissionIds.push(pathPermission.id)
              }
            } else {
              console.warn(`菜单 ${item.title} (${item.path}) 没有关联的权限，无法授权`)
              console.warn('可用权限path:', allPermissions.value.filter(p => p.type === 'menu').map(p => p.path))
            }
          }
        }
        // 继续遍历子菜单
        if (item.children && item.children.length > 0) {
          collectPermissions(item.children)
        }
      })
    }
    collectPermissions(menuTree.value)
    
    console.log('收集到的权限ID:', permissionIds)
    
    if (permissionIds.length === 0) {
      ElMessage.warning('选中的菜单没有关联的权限，请检查菜单配置或联系管理员')
      return
    }
    
    savingMenus.value = true
    await setRolePermissions({
      id: currentRoleId.value,
      permissionIds: permissionIds
    })
    ElMessage.success('菜单分配成功')
    menuDialogVisible.value = false
  } catch (error) {
    console.error('菜单分配失败:', error)
    ElMessage.error(error.response?.data?.message || error.message || '菜单分配失败')
  } finally {
    savingMenus.value = false
  }
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

.table-card { 
  border-radius: 16px; 
  border: none; 
  box-shadow: 0 4px 12px rgba(0,0,0,0.03); 
}
</style>

