<template>
  <div class="page-container">
    <el-row :gutter="20">
      <!-- 左侧个人信息卡片 -->
      <el-col :span="8" :xs="24">
        <el-card shadow="never" class="profile-card">
          <div class="user-profile">
            <div class="avatar-wrapper">
              <el-upload
                class="avatar-uploader"
                :show-file-list="false"
                :before-upload="beforeAvatarUpload"
                :http-request="handleAvatarUpload"
                accept="image/*"
              >
                <el-avatar :size="100" :src="user.avatar">
                  <el-icon><User /></el-icon>
                </el-avatar>
                <div class="avatar-edit-overlay">
                  <el-icon><Camera /></el-icon>
                  <span>更换头像</span>
                </div>
              </el-upload>
              <div class="role-tag">
                <el-tag type="success" size="small" effect="dark">{{ user.role }}</el-tag>
              </div>
            </div>
            <h3 class="username">{{ user.nickname }}</h3>
            <p class="bio">{{ user.bio }}</p>
            
            <div class="user-meta">
              <div class="meta-item">
                <el-icon><User /></el-icon>
                <span>{{ user.username }}</span>
              </div>
              <div class="meta-item">
                <el-icon><Message /></el-icon>
                <span>{{ user.email }}</span>
              </div>
              <div class="meta-item">
                <el-icon><Location /></el-icon>
                <span>{{ user.location }}</span>
              </div>
            </div>

            <div class="divider"></div>
            
            <div class="skills">
              <div class="skills-title">标签</div>
              <div class="tags">
                 <el-tag v-for="tag in user.tags" :key="tag" class="tag-item" size="small">{{ tag }}</el-tag>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧设置/活动区域 -->
      <el-col :span="16" :xs="24">
        <el-card shadow="never" class="settings-card">
          <el-tabs v-model="activeTab">
            <el-tab-pane label="基本资料" name="info">
              <el-form :model="form" label-width="80px" class="profile-form">
                <el-form-item label="头像">
                  <div class="avatar-form-item">
                    <el-avatar :size="80" :src="form.avatar || user.avatar">
                      <el-icon><User /></el-icon>
                    </el-avatar>
                    <div class="avatar-input-group">
                      <el-input 
                        v-model="form.avatar" 
                        placeholder="请输入头像URL或上传图片"
                        clearable
                        size="large"
                      >
                        <template #prefix>
                          <el-icon><Link /></el-icon>
                        </template>
                      </el-input>
                      <el-upload
                        class="avatar-url-uploader"
                        :show-file-list="false"
                        :before-upload="beforeAvatarUpload"
                        :http-request="handleAvatarUpload"
                        accept="image/*"
                      >
                        <el-button type="primary" :icon="Upload" size="large">上传图片</el-button>
                      </el-upload>
                    </div>
                  </div>
                </el-form-item>
                <el-form-item label="昵称">
                  <el-input v-model="form.nickname" size="large" />
                </el-form-item>
                <el-form-item label="邮箱">
                  <el-input v-model="form.email" size="large" />
                </el-form-item>
                <el-form-item label="个人简介">
                  <el-input v-model="form.bio" type="textarea" :rows="3" size="large" />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="handleUpdateInfo" size="large">保存修改</el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>
            
            <el-tab-pane label="修改密码" name="password">
              <el-form :model="pwdForm" label-width="80px" class="profile-form">
                <el-form-item label="旧密码">
                  <el-input v-model="pwdForm.old" type="password" show-password />
                </el-form-item>
                <el-form-item label="新密码">
                  <el-input v-model="pwdForm.new" type="password" show-password />
                </el-form-item>
                <el-form-item label="确认密码">
                  <el-input v-model="pwdForm.confirm" type="password" show-password />
                </el-form-item>
                 <el-form-item>
                  <el-button type="primary" @click="handleChangePassword">修改密码</el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>

             <el-tab-pane label="登录日志" name="log">
                <el-timeline>
                  <el-timeline-item timestamp="2024-05-28 10:23:00" placement="top">
                    <el-card shadow="hover">
                      <h4>登录成功</h4>
                      <p>IP: 192.168.1.1 (Chrome / Windows 10)</p>
                    </el-card>
                  </el-timeline-item>
                   <el-timeline-item timestamp="2024-05-27 18:45:00" placement="top">
                    <el-card shadow="hover">
                      <h4>登录成功</h4>
                      <p>IP: 192.168.1.1 (Chrome / Windows 10)</p>
                    </el-card>
                  </el-timeline-item>
                </el-timeline>
             </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { User, Message, Location, Camera, Link, Upload } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { getCurrentUser, updateProfile, changePassword, getCurrentUserRoles } from '@/api/user'
import { getRoleList } from '@/api/role'

const activeTab = ref('info')

const user = ref({
  avatar: '',
  nickname: '',
  username: '',
  role: '用户',
  email: '',
  location: 'Shanghai, China',
  bio: '',
  tags: ['Golang', 'Vue 3', 'Kubernetes', 'DevOps']
})

const form = reactive({
  avatar: '',
  nickname: '',
  email: '',
  bio: ''
})

const pwdForm = reactive({
  old: '',
  new: '',
  confirm: ''
})

const fetchUserInfo = async () => {
  try {
    const res = await getCurrentUser()
    
    // 获取当前用户角色
    let roleName = '用户' // 默认角色
    try {
      const roleRes = await getCurrentUserRoles()
      if (roleRes && roleRes.roleIds && roleRes.roleIds.length > 0) {
        // 获取角色列表
        const roleListRes = await getRoleList({ page: 1, size: 100 })
        if (roleListRes && roleListRes.list) {
          // 找到第一个角色ID对应的角色名称
          const role = roleListRes.list.find(r => r.id === roleRes.roleIds[0])
          if (role) {
            roleName = role.name
          }
        }
      }
    } catch (error) {
      console.warn('获取用户角色失败:', error)
    }
    
    user.value = {
      avatar: res.avatar || 'https://picsum.photos/id/64/200/200',
      nickname: res.nickname || res.username,
      username: res.username,
      role: roleName,
      email: res.email || '',
      location: 'Shanghai, China',
      bio: '',
      tags: ['Golang', 'Vue 3', 'Kubernetes', 'DevOps']
    }
    form.avatar = res.avatar || ''
    form.nickname = res.nickname || res.username
    form.email = res.email || ''
    form.bio = ''
  } catch (error) {
    ElMessage.error('获取用户信息失败')
    console.error(error)
  }
}

const handleUpdateInfo = async () => {
  try {
    await updateProfile({
      nickname: form.nickname,
      email: form.email,
      avatar: form.avatar,
      bio: form.bio
    })
    ElMessage.success('个人资料更新成功')
    fetchUserInfo()
  } catch (error) {
    ElMessage.error('更新失败')
    console.error(error)
  }
}

const handleChangePassword = async () => {
  if (!pwdForm.old || !pwdForm.new) {
    ElMessage.warning('请填写完整信息')
    return
  }
  if (pwdForm.new !== pwdForm.confirm) {
    ElMessage.error('两次输入的密码不一致')
    return
  }
  try {
    await changePassword({
      oldPassword: pwdForm.old,
      newPassword: pwdForm.new
    })
    ElMessage.success('密码修改成功，下次登录生效')
    pwdForm.old = ''
    pwdForm.new = ''
    pwdForm.confirm = ''
  } catch (error) {
    ElMessage.error('密码修改失败')
    console.error(error)
  }
}

// 头像上传前的验证
const beforeAvatarUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB!')
    return false
  }
  return true
}

// 处理头像上传
const handleAvatarUpload = async (options) => {
  const file = options.file
  
  try {
    // 先压缩图片，减少base64数据大小
    ElMessage.info('正在处理图片...')
    const compressedBase64 = await compressImage(file, 0.7, 400) // 压缩到70%质量，最大宽度400px（头像不需要太大）
    
    form.avatar = compressedBase64
    user.value.avatar = compressedBase64
    
    // 自动保存
    try {
      await updateProfile({
        nickname: form.nickname,
        email: form.email,
        avatar: compressedBase64,
        bio: form.bio
      })
      ElMessage.success('头像更新成功')
      // 更新localStorage（只保存前500字符，避免localStorage过大）
      const avatarPreview = compressedBase64.substring(0, 500)
      localStorage.setItem('avatar', avatarPreview)
    } catch (error) {
      ElMessage.error('头像保存失败：' + (error.message || '未知错误'))
      console.error(error)
    }
  } catch (error) {
    ElMessage.error('头像上传失败')
    console.error(error)
  }
}

// 压缩图片
const compressImage = (file, quality = 0.8, maxWidth = 800) => {
  return new Promise((resolve) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      const img = new Image()
      img.onload = () => {
        const canvas = document.createElement('canvas')
        let width = img.width
        let height = img.height
        
        // 如果宽度超过最大宽度，按比例缩放
        if (width > maxWidth) {
          height = (height * maxWidth) / width
          width = maxWidth
        }
        
        canvas.width = width
        canvas.height = height
        
        const ctx = canvas.getContext('2d')
        ctx.drawImage(img, 0, 0, width, height)
        
        // 转换为base64
        const compressedBase64 = canvas.toDataURL('image/jpeg', quality)
        resolve(compressedBase64)
      }
      img.src = e.target.result
    }
    reader.readAsDataURL(file)
  })
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped>
.page-container {
}

.profile-card, .settings-card {
  border-radius: 16px;
  border: none;
  box-shadow: 0 4px 12px rgba(0,0,0,0.03);
  margin-bottom: 20px;
}

.user-profile {
  text-align: center;
  padding: 20px 0;
}

.avatar-wrapper {
  position: relative;
  display: inline-block;
  margin-bottom: 16px;
}

.avatar-uploader {
  position: relative;
  cursor: pointer;
}

.avatar-uploader:hover .avatar-edit-overlay {
  opacity: 1;
}

.avatar-edit-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
  opacity: 0;
  transition: opacity 0.3s;
  font-size: 12px;
}

.avatar-edit-overlay .el-icon {
  font-size: 20px;
  margin-bottom: 4px;
}

.role-tag {
  position: absolute;
  bottom: 0;
  right: -10px;
}

.username {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px;
}

.bio {
  font-size: 14px;
  color: #909399;
  margin: 0 0 20px;
  padding: 0 20px;
}

.user-meta {
  text-align: left;
  padding: 0 20px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
  color: #606266;
  font-size: 14px;
}

.divider {
  height: 1px;
  background-color: #f0f2f5;
  margin: 20px;
}

.skills-title {
  text-align: left;
  padding: 0 20px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #303133;
}

.tags {
  padding: 0 20px;
  text-align: left;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.profile-form {
  max-width: 500px;
  margin-top: 20px;
}

.avatar-form-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.avatar-input-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.avatar-url-uploader {
  display: inline-block;
}
</style>

