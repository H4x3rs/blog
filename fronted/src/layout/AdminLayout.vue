<template>
  <el-container class="admin-container">
    <!-- Desktop Sidebar -->
    <el-aside :width="isCollapse ? '64px' : '220px'" class="admin-aside hidden-xs-only">
      <div class="aside-header">
        <span class="logo-icon">B</span>
        <span class="logo-text" v-show="!isCollapse">{{ siteName || 'Blog Admin' }}</span>
      </div>
      <el-menu
        active-text-color="#409EFF"
        background-color="#ffffff"
        class="el-menu-vertical"
        :default-active="$route.path"
        text-color="#606266"
        :collapse="isCollapse"
        :router="true"
        border="false"
        v-loading="menuLoading"
      >
        <!-- 动态渲染菜单 -->
        <template v-for="menu in menuList" :key="menu.id">
          <!-- 有子菜单的情况 -->
          <el-sub-menu v-if="menu.children && menu.children.length > 0" :index="menu.path || menu.id.toString()">
            <template #title>
              <el-icon v-if="menu.icon"><component :is="getIconComponent(menu.icon)" /></el-icon>
              <span>{{ menu.title }}</span>
            </template>
            <template v-for="child in menu.children" :key="child.id">
              <!-- 子菜单还有子菜单 -->
              <el-sub-menu v-if="child.children && child.children.length > 0" :index="child.path || child.id.toString()">
                <template #title>
                  <el-icon v-if="child.icon"><component :is="getIconComponent(child.icon)" /></el-icon>
                  <span>{{ child.title }}</span>
                </template>
                <el-menu-item v-for="grandchild in child.children" :key="grandchild.id" :index="grandchild.path">
                  <el-icon v-if="grandchild.icon"><component :is="getIconComponent(grandchild.icon)" /></el-icon>
                  <template #title>{{ grandchild.title }}</template>
                </el-menu-item>
              </el-sub-menu>
              <!-- 子菜单是普通菜单项 -->
              <el-menu-item v-else :index="child.path">
                <el-icon v-if="child.icon"><component :is="getIconComponent(child.icon)" /></el-icon>
                <template #title>{{ child.title }}</template>
              </el-menu-item>
            </template>
          </el-sub-menu>
          <!-- 没有子菜单的情况 -->
          <el-menu-item v-else :index="menu.path">
            <el-icon v-if="menu.icon"><component :is="getIconComponent(menu.icon)" /></el-icon>
            <template #title>{{ menu.title }}</template>
          </el-menu-item>
        </template>

        <!-- 返回首页（始终显示） -->
        <el-menu-item index="/">
          <el-icon><HomeFilled /></el-icon>
          <template #title>返回首页</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- Mobile Sidebar Drawer -->
    <el-drawer
      v-model="mobileSidebarVisible"
      direction="ltr"
      size="220px"
      :with-header="false"
      class="mobile-sidebar"
    >
       <div class="aside-header mobile">
        <span class="logo-icon">B</span>
        <span class="logo-text">{{ siteName || 'Blog Admin' }}</span>
      </div>
       <el-menu
        active-text-color="#409EFF"
        background-color="#ffffff"
        class="el-menu-vertical mobile"
        :default-active="$route.path"
        text-color="#606266"
        :router="true"
        border="false"
        v-loading="menuLoading"
      >
        <!-- 动态渲染菜单 -->
        <template v-for="menu in menuList" :key="menu.id">
          <!-- 有子菜单的情况 -->
          <el-sub-menu v-if="menu.children && menu.children.length > 0" :index="menu.path || menu.id.toString()">
            <template #title>
              <el-icon v-if="menu.icon"><component :is="getIconComponent(menu.icon)" /></el-icon>
              <span>{{ menu.title }}</span>
            </template>
            <template v-for="child in menu.children" :key="child.id">
              <!-- 子菜单还有子菜单 -->
              <el-sub-menu v-if="child.children && child.children.length > 0" :index="child.path || child.id.toString()">
                <template #title>
                  <el-icon v-if="child.icon"><component :is="getIconComponent(child.icon)" /></el-icon>
                  <span>{{ child.title }}</span>
                </template>
                <el-menu-item v-for="grandchild in child.children" :key="grandchild.id" :index="grandchild.path" @click="mobileSidebarVisible = false">
                  <el-icon v-if="grandchild.icon"><component :is="getIconComponent(grandchild.icon)" /></el-icon>
                  <span>{{ grandchild.title }}</span>
                </el-menu-item>
              </el-sub-menu>
              <!-- 子菜单是普通菜单项 -->
              <el-menu-item v-else :index="child.path" @click="mobileSidebarVisible = false">
                <el-icon v-if="child.icon"><component :is="getIconComponent(child.icon)" /></el-icon>
                <span>{{ child.title }}</span>
              </el-menu-item>
            </template>
          </el-sub-menu>
          <!-- 没有子菜单的情况 -->
          <el-menu-item v-else :index="menu.path" @click="mobileSidebarVisible = false">
            <el-icon v-if="menu.icon"><component :is="getIconComponent(menu.icon)" /></el-icon>
            <span>{{ menu.title }}</span>
          </el-menu-item>
        </template>

        <!-- 返回首页（始终显示） -->
        <el-menu-item index="/" @click="mobileSidebarVisible = false">
          <el-icon><HomeFilled /></el-icon>
          <span>返回首页</span>
        </el-menu-item>
      </el-menu>
    </el-drawer>
    
    <el-container>
      <el-header class="admin-header">
        <div class="header-left">
           <!-- Collapse Trigger -->
           <div class="collapse-btn hidden-xs-only" @click="isCollapse = !isCollapse">
             <el-icon :size="20"><component :is="isCollapse ? Expand : Fold" /></el-icon>
           </div>
           
           <!-- Mobile Trigger -->
           <div class="collapse-btn hidden-sm-and-up" @click="mobileSidebarVisible = true">
             <el-icon :size="20"><Expand /></el-icon>
           </div>

           <el-breadcrumb separator="/" class="hidden-xs-only">
            <el-breadcrumb-item :to="{ path: '/admin/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>{{ $route.meta.title }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown trigger="click">
            <span class="el-dropdown-link user-info">
              <el-avatar size="small" :src="userAvatar">
                <el-icon v-if="!userAvatar"><User /></el-icon>
              </el-avatar>
              <span class="username hidden-xs-only">{{ displayName || 'Admin' }}</span>
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="$router.push('/admin/profile')">个人中心</el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-main class="admin-main">
        <router-view v-slot="{ Component }">
          <transition name="fade-transform" mode="out-in">
             <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Odometer, Document, HomeFilled, ArrowDown, Fold, Expand,
  Files, CollectionTag, Reading, User, Avatar, Key, Setting, DocumentCopy, Menu
} from '@element-plus/icons-vue'
import { useSiteConfig } from '../store/site'
import { getCurrentUser } from '../api/user'
import { getCurrentUserMenus } from '../api/menu'
import 'element-plus/theme-chalk/display.css'

const route = useRoute()
const router = useRouter()
const isCollapse = ref(false)
const mobileSidebarVisible = ref(false)
const { siteName } = useSiteConfig()

// 图标映射表
const iconMap = {
  'Odometer': Odometer,
  'Document': Document,
  'DocumentCopy': DocumentCopy,
  'Files': Files,
  'CollectionTag': CollectionTag,
  'Reading': Reading,
  'User': User,
  'Avatar': Avatar,
  'Key': Key,
  'Setting': Setting,
  'Menu': Menu
}

// 获取图标组件
const getIconComponent = (iconName) => {
  return iconMap[iconName] || null
}

// 用户信息
const userInfo = ref({
  username: '',
  nickname: '',
  avatar: ''
})

// 菜单列表
const menuList = ref([])
const menuLoading = ref(false)

// 显示名称（优先显示nickname，如果没有则显示username）
const displayName = computed(() => {
  return userInfo.value.nickname || userInfo.value.username || 'Admin'
})

// 用户头像
const userAvatar = computed(() => {
  return userInfo.value.avatar || ''
})

// 加载用户信息
const loadUserInfo = async () => {
  // 先从localStorage读取
  const username = localStorage.getItem('username')
  const nickname = localStorage.getItem('nickname')
  const avatar = localStorage.getItem('avatar')
  
  // 立即显示localStorage中的信息（快速显示）
  if (username) {
    userInfo.value.username = username
    userInfo.value.nickname = nickname || ''
    userInfo.value.avatar = avatar || ''
  }
  
  // 如果有token，尝试从后端获取完整用户信息
  const token = localStorage.getItem('token')
  if (token) {
    try {
      const res = await getCurrentUser()
      if (res) {
        // 处理响应数据（可能是嵌套的User对象）
        const userData = res.user || res
        userInfo.value.username = userData.username || username || ''
        userInfo.value.nickname = userData.nickname || nickname || ''
        userInfo.value.avatar = userData.avatar || avatar || ''
        
        // 更新localStorage
        if (userData.username) {
          localStorage.setItem('username', userData.username)
        }
        if (userData.nickname) {
          localStorage.setItem('nickname', userData.nickname)
        }
        if (userData.avatar) {
          localStorage.setItem('avatar', userData.avatar)
        }
      }
    } catch (error) {
      // 如果获取失败，使用localStorage中的信息
      console.warn('获取用户信息失败:', error)
      // 确保至少显示localStorage中的信息
      if (!userInfo.value.username && username) {
        userInfo.value.username = username
        userInfo.value.nickname = nickname || ''
        userInfo.value.avatar = avatar || ''
      }
    }
  } else {
    // 没有token，清空用户信息
    userInfo.value.username = ''
    userInfo.value.nickname = ''
    userInfo.value.avatar = ''
  }
}

// 退出登录
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要退出登录吗？',
      '确认退出',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 清除localStorage中的用户信息
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    localStorage.removeItem('nickname')
    localStorage.removeItem('avatar')
    localStorage.removeItem('loginProvider')
    
    ElMessage.success('已退出登录')
    
    // 跳转到登录页面
    router.push('/login-new')
  } catch {
    // 用户取消
  }
}

// 加载用户菜单
const loadUserMenus = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    menuList.value = []
    return
  }
  
  menuLoading.value = true
  try {
    const res = await getCurrentUserMenus()
    // 确保菜单列表保持顺序（按order字段排序）
    const menus = res.list || []
    // 后端已经排序，这里直接使用
    menuList.value = menus
    console.log('加载的菜单列表:', menuList.value)
  } catch (error) {
    console.error('获取用户菜单失败:', error)
    menuList.value = []
  } finally {
    menuLoading.value = false
  }
}

// 页面加载时加载用户信息和菜单
onMounted(() => {
  loadUserInfo()
  loadUserMenus()
})

// 监听路由变化，重新加载用户信息和菜单（防止从其他页面跳转过来时信息丢失）
watch(() => route.path, () => {
  if (route.path.startsWith('/admin')) {
    loadUserInfo()
    loadUserMenus()
  }
})
</script>

<style scoped>
.admin-container {
  height: 100vh;
}

.admin-aside {
  background-color: #ffffff;
  color: #303133;
  display: flex;
  flex-direction: column;
  transition: width 0.3s;
  overflow-x: hidden;
  box-shadow: 2px 0 8px rgba(0,0,0,0.05);
  z-index: 10;
}

.aside-header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #ffffff;
  white-space: nowrap;
  overflow: hidden;
  border-bottom: 1px solid #f0f2f5;
}

.logo-icon {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #409eff 0%, #36d1dc 100%);
  color: white;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  margin-right: 10px;
  flex-shrink: 0;
  font-size: 18px;
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  color: #303133;
}

.el-menu-vertical {
  border-right: none;
  border-top: 1px solid transparent; /* Fix for potential double border */
}

.el-menu-vertical:not(.el-menu--collapse) {
  width: 220px;
}

.admin-header {
  background-color: white;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
}

.header-left {
  display: flex;
  align-items: center;
}

.collapse-btn {
  margin-right: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.username {
  margin: 0 8px;
  font-weight: 500;
  color: #606266;
}

.admin-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}

/* Mobile Sidebar Styles */
.mobile-sidebar :deep(.el-drawer__body) {
    padding: 0;
    background-color: #ffffff;
}
.el-menu-vertical.mobile {
    width: 100%;
}
.aside-header.mobile {
    justify-content: center;
    color: #303133;
}

/* Transition */
.fade-transform-leave-active,
.fade-transform-enter-active {
  transition: all 0.5s;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

@media (max-width: 768px) {
    .admin-main {
        padding: 10px;
    }
}
</style>
