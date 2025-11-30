<template>
  <div class="page-container">
    <!-- 筛选区域 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filters" class="filter-form">
        <el-form-item label="用户名">
          <el-input 
            v-model="filters.username" 
            placeholder="请输入用户名" 
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
            <span style="margin-left: 4px;">新增用户</span>
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格区域 -->
    <el-card shadow="never" class="table-card">
      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="头像" width="80">
          <template #default="scope">
            <el-avatar :size="32" :src="scope.row.avatar || 'https://picsum.photos/id/64/100/100'" />
          </template>
        </el-table-column>
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="scope">
            <el-button link type="primary" size="small" :icon="Edit" @click="openDialog('edit', scope.row)">编辑</el-button>
            <el-button link type="success" size="small" :icon="UserFilled" @click="openRoleDialog(scope.row)">角色</el-button>
            <el-button link type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)">删除</el-button>
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
        />
      </div>
    </el-card>

    <el-dialog 
      v-model="dialogVisible" 
      :title="dialogType === 'create' ? '新增用户' : '编辑用户'" 
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
        <!-- 头像预览区域 -->
        <el-form-item label="头像" prop="avatar">
          <div class="avatar-upload-container">
            <el-avatar 
              :size="80" 
              :src="form.avatar || 'https://picsum.photos/id/64/100/100'"
              class="avatar-preview"
            >
              <el-icon><User /></el-icon>
            </el-avatar>
            <div class="avatar-input-wrapper">
              <el-input 
                v-model="form.avatar" 
                placeholder="请输入头像URL"
                clearable
              >
                <template #prefix>
                  <el-icon><Link /></el-icon>
                </template>
              </el-input>
              <div class="avatar-tip">支持图片URL链接</div>
            </div>
          </div>
        </el-form-item>

        <!-- 用户名 -->
        <el-form-item label="用户名" prop="username">
          <el-input 
            v-model="form.username" 
            placeholder="请输入用户名"
            clearable
            :disabled="dialogType === 'edit'"
          >
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <!-- 密码 -->
        <el-form-item 
          label="密码"
          :prop="dialogType === 'create' ? 'password' : ''"
        >
          <el-input 
            v-model="form.password" 
            type="password" 
            :placeholder="dialogType === 'create' ? '请输入密码' : '留空则不修改密码'"
            show-password
            clearable
          >
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <!-- 昵称 -->
        <el-form-item label="昵称" prop="nickname">
          <el-input 
            v-model="form.nickname" 
            placeholder="请输入昵称"
            clearable
          >
            <template #prefix>
              <el-icon><EditPen /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <!-- 邮箱 -->
        <el-form-item label="邮箱" prop="email">
          <el-input 
            v-model="form.email" 
            placeholder="请输入邮箱地址"
            clearable
          >
            <template #prefix>
              <el-icon><Message /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <!-- 状态（仅编辑时显示） -->
        <el-form-item label="状态" prop="status" v-if="dialogType === 'edit'">
          <el-switch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
            active-text="启用"
            inactive-text="禁用"
            inline-prompt
            size="large"
            class="status-switch"
            style="--el-switch-on-color: #67c23a;"
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

    <!-- 角色授权对话框 -->
    <el-dialog 
      v-model="roleDialogVisible" 
      title="分配角色" 
      width="500px"
      :close-on-click-modal="false"
    >
      <el-checkbox-group v-model="selectedRoleIds" style="width: 100%">
        <el-checkbox 
          v-for="role in allRoles" 
          :key="role.id" 
          :label="role.id"
          style="display: block; margin-bottom: 12px; padding: 8px; border-radius: 4px; border: 1px solid #e4e7ed;"
        >
          <div>
            <div style="font-weight: 500; color: #303133;">{{ role.name }}</div>
            <div style="font-size: 12px; color: #909399; margin-top: 4px;">{{ role.desc || '无描述' }}</div>
            <el-tag size="small" style="margin-top: 4px;">{{ role.key }}</el-tag>
          </div>
        </el-checkbox>
      </el-checkbox-group>
      
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveRoles" :loading="savingRoles">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { 
  Plus, Edit, Delete, Search, Refresh, 
  User, Lock, EditPen, Message, Link, UserFilled
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserList, createUser, updateUser, deleteUser, getUser, setUserRoles, getUserRoles } from '@/api/user'
import { getRoleList } from '@/api/role'

const loading = ref(false)
const tableData = ref([])
const totalCount = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const filters = reactive({
  username: ''
})

const dialogVisible = ref(false)
const dialogType = ref('create')
const formRef = ref(null)
const submitting = ref(false)
const form = reactive({ 
  id: null, 
  username: '', 
  password: '',
  nickname: '',
  email: '', 
  avatar: '',
  status: 1
})

// 角色授权相关
const roleDialogVisible = ref(false)
const allRoles = ref([])
const selectedRoleIds = ref([])
const currentUserId = ref(null)
const savingRoles = ref(false)

// 表单验证规则
const formRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  nickname: [
    { max: 50, message: '昵称长度不能超过 50 个字符', trigger: 'blur' }
  ],
  avatar: [
    { type: 'url', message: '请输入正确的URL地址', trigger: 'blur' }
  ]
}

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getUserList({
      page: currentPage.value,
      size: pageSize.value,
      username: filters.username
    })
    tableData.value = res.list || []
    totalCount.value = res.total || 0
  } catch (error) {
    ElMessage.error('获取用户列表失败')
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
  filters.username = ''
  currentPage.value = 1
  fetchList()
}

const openDialog = async (type, row) => {
  dialogType.value = type
  dialogVisible.value = true
  
  // 重置表单验证状态
  if (formRef.value) {
    formRef.value.clearValidate()
  }
  
  if (type === 'edit' && row) {
    try {
      const res = await getUser({ id: row.id })
      Object.assign(form, {
        id: res.id,
        username: res.username,
        nickname: res.nickname || '',
        email: res.email || '',
        avatar: res.avatar || '',
        status: res.status || 1,
        password: '' // 编辑时不显示密码
      })
    } catch (error) {
      ElMessage.error('获取用户信息失败')
      console.error(error)
      dialogVisible.value = false
    }
  } else {
    Object.assign(form, { 
      id: null, 
      username: '', 
      password: '',
      nickname: '',
      email: '', 
      avatar: '',
      status: 1
    })
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  // 验证表单
  try {
    await formRef.value.validate()
  } catch (error) {
    return
  }
  
  // 创建时密码必填（表单验证已处理，这里作为双重检查）
  if (dialogType.value === 'create' && !form.password) {
    ElMessage.warning('请输入密码')
    return
  }
  
  // 编辑时如果填写了密码，需要验证密码格式
  if (dialogType.value === 'edit' && form.password) {
    if (form.password.length < 6 || form.password.length > 20) {
      ElMessage.warning('密码长度在 6 到 20 个字符')
      return
    }
  }
  
  submitting.value = true
  try {
    if (dialogType.value === 'create') {
      await createUser({
        username: form.username,
        password: form.password,
        nickname: form.nickname,
        email: form.email,
        avatar: form.avatar
      })
      ElMessage.success('创建成功')
    } else {
      const updateData = {
        username: form.username,
        nickname: form.nickname,
        email: form.email,
        avatar: form.avatar,
        status: form.status
      }
      if (form.password) {
        updateData.password = form.password
      }
      await updateUser({ id: form.id, ...updateData })
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

const handleDelete = (row) => {
  ElMessageBox.confirm('确认删除该用户？', '提示', { 
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  })
    .then(async () => {
      try {
        await deleteUser({ id: row.id })
        ElMessage.success('删除成功')
        fetchList()
      } catch (error) {
        ElMessage.error('删除失败')
        console.error(error)
      }
    })
    .catch(() => {})
}

const openRoleDialog = async (row) => {
  currentUserId.value = row.id
  roleDialogVisible.value = true
  selectedRoleIds.value = []
  
  // 加载所有角色
  try {
    const res = await getRoleList({ page: 1, size: 100 })
    allRoles.value = res.list || []
  } catch (error) {
    ElMessage.error('获取角色列表失败')
    console.error(error)
    return
  }
  
  // 加载用户已有角色
  try {
    const res = await getUserRoles({ id: row.id })
    selectedRoleIds.value = res.roleIds || []
  } catch (error) {
    ElMessage.error('获取用户角色失败')
    console.error(error)
  }
}

const handleSaveRoles = async () => {
  if (!currentUserId.value) return
  
  savingRoles.value = true
  try {
    await setUserRoles({
      id: currentUserId.value,
      roleIds: selectedRoleIds.value
    })
    ElMessage.success('角色分配成功')
    roleDialogVisible.value = false
  } catch (error) {
    ElMessage.error('角色分配失败')
    console.error(error)
  } finally {
    savingRoles.value = false
  }
}

watch([currentPage, pageSize], () => {
  fetchList()
})

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
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

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  padding: 20px 0;
}

/* 头像上传区域 */
.avatar-upload-container {
  display: flex;
  align-items: flex-start;
  gap: 20px;
}

.avatar-preview {
  flex-shrink: 0;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
}

.avatar-input-wrapper {
  flex: 1;
}

.avatar-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
  line-height: 1.5;
}

/* 表单样式优化 */
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

/* 密码输入框特殊样式 */
:deep(.el-input--password .el-input__wrapper) {
  padding-right: 40px;
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

:deep(.el-radio-group) {
  display: flex;
  gap: 20px;
}

:deep(.el-radio) {
  display: flex;
  align-items: center;
  margin-right: 0;
}

/* 对话框底部 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 10px;
}

/* 对话框样式 */
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

/* 状态开关优化样式 */
.status-switch {
  height: 32px;
}

.status-switch :deep(.el-switch__core) {
  width: 60px !important;
  height: 28px !important;
  background-color: #dcdfe6 !important;
  border-color: #dcdfe6 !important;
  border-radius: 14px !important;
  transition: all 0.3s ease;
}

.status-switch.is-checked :deep(.el-switch__core) {
  background-color: #67c23a !important;
  border-color: #67c23a !important;
}

.status-switch :deep(.el-switch__core::after) {
  width: 24px !important;
  height: 24px !important;
  background-color: #ffffff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
  transition: all 0.3s ease;
}

.status-switch.is-checked :deep(.el-switch__core::after) {
  left: calc(100% - 26px) !important;
}

/* 文字样式优化 */
.status-switch :deep(.el-switch__label) {
  font-size: 13px;
  font-weight: 500;
  color: #606266;
  transition: color 0.3s ease;
}

.status-switch :deep(.el-switch__label.is-active) {
  color: #67c23a;
  font-weight: 600;
}

.status-switch :deep(.el-switch__label--left) {
  margin-right: 8px;
}

.status-switch :deep(.el-switch__label--right) {
  margin-left: 8px;
}

/* 确保激活状态时显示绿色 */
.status-switch :deep(.el-switch.is-checked .el-switch__core) {
  background-color: #67c23a !important;
  border-color: #67c23a !important;
}

/* 悬停效果 */
.status-switch:hover :deep(.el-switch__core) {
  box-shadow: 0 0 0 4px rgba(103, 194, 58, 0.1);
}

.status-switch.is-checked:hover :deep(.el-switch__core) {
  box-shadow: 0 0 0 4px rgba(103, 194, 58, 0.2);
}
</style>

