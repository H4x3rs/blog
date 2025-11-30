<template>
  <el-container class="layout-container">
    <el-header class="header">
      <div class="header-inner">
        <div class="logo" @click="$router.push('/')">
          <span class="logo-icon">B</span>
          <span class="logo-text">{{ siteName || 'Blog System' }}</span>
        </div>
        
        <!-- Desktop Nav -->
        <div class="nav hidden-xs-only">
          <el-menu 
            mode="horizontal" 
            :router="true" 
            :default-active="$route.path"
            :ellipsis="false"
            class="nav-menu"
          >
            <el-menu-item index="/">首页</el-menu-item>
            <el-menu-item index="/category">分类</el-menu-item>
            <el-menu-item index="/tag">标签</el-menu-item>
            <el-menu-item index="/topic">专题</el-menu-item>
            <el-menu-item index="/about">关于</el-menu-item>
          </el-menu>
        </div>

        <div class="actions hidden-xs-only">
          <!-- 已登录：显示用户信息 -->
          <el-dropdown v-if="isLoggedIn" trigger="click" class="user-dropdown">
            <span class="user-info-link">
              <el-avatar :size="32" :src="userAvatar">
                <el-icon v-if="!userAvatar"><User /></el-icon>
              </el-avatar>
              <span class="username-text">{{ displayName }}</span>
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="$router.push('/admin/dashboard')">管理后台</el-dropdown-item>
                <el-dropdown-item @click="$router.push('/admin/profile')">个人中心</el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <!-- 未登录：显示登录按钮 -->
          <el-button v-else class="login-btn" round @click="$router.push('/login')">登录 / 注册</el-button>
        </div>

        <!-- Mobile Menu Toggle -->
        <div class="mobile-menu-toggle hidden-sm-and-up" @click="mobileMenuVisible = !mobileMenuVisible">
           <el-icon :size="24"><Menu /></el-icon>
        </div>
      </div>
    </el-header>

    <!-- Mobile Menu Drawer -->
    <el-drawer
      v-model="mobileMenuVisible"
      direction="rtl"
      size="70%"
      :show-close="false"
      class="mobile-drawer"
    >
       <div class="mobile-nav">
          <el-menu 
            :router="true" 
            :default-active="$route.path"
            class="mobile-nav-menu"
          >
            <el-menu-item index="/" @click="mobileMenuVisible = false">首页</el-menu-item>
            <el-menu-item index="/category" @click="mobileMenuVisible = false">分类</el-menu-item>
            <el-menu-item index="/tag" @click="mobileMenuVisible = false">标签</el-menu-item>
            <el-menu-item index="/topic" @click="mobileMenuVisible = false">专题</el-menu-item>
            <el-menu-item index="/about" @click="mobileMenuVisible = false">关于</el-menu-item>
          </el-menu>
          <div class="mobile-actions">
            <!-- 已登录：显示用户信息 -->
            <div v-if="isLoggedIn" class="mobile-user-info">
              <el-dropdown trigger="click" @visible-change="mobileMenuVisible = false">
                <div class="mobile-user-card">
                  <el-avatar :size="40" :src="userAvatar">
                    <el-icon v-if="!userAvatar"><User /></el-icon>
                  </el-avatar>
                  <div class="mobile-user-details">
                    <div class="mobile-username">{{ displayName }}</div>
                    <div class="mobile-user-email">{{ userInfo.username }}</div>
                  </div>
                  <el-icon><ArrowDown /></el-icon>
                </div>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="$router.push('/admin/dashboard')">管理后台</el-dropdown-item>
                    <el-dropdown-item @click="$router.push('/admin/profile')">个人中心</el-dropdown-item>
                    <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
            <!-- 未登录：显示登录按钮 -->
            <el-button v-else class="login-btn" block @click="$router.push('/login-new'); mobileMenuVisible = false">登录 / 注册</el-button>
          </div>
       </div>
    </el-drawer>
    
    <el-main class="main-content">
      <router-view />
    </el-main>
    
    <el-footer class="footer">
      <div class="footer-inner">
        <el-row :gutter="40">
          <!-- Logo & Intro -->
          <el-col :span="8" :xs="24" class="footer-col">
            <div class="footer-logo">
              <span class="logo-icon small">B</span>
              <span class="logo-text-light">{{ siteName || 'Blog System' }}</span>
            </div>
            <p class="footer-desc">
              探索技术之美，记录成长足迹。
              分享 Go, Vue, 云原生等前沿技术心得。
            </p>
          </el-col>

          <!-- Quick Links -->
          <el-col :span="8" :xs="24" class="footer-col">
            <h4 class="footer-title">快速链接</h4>
            <ul class="footer-links-list">
              <li><router-link to="/">首页</router-link></li>
              <li><router-link to="/category">文章分类</router-link></li>
              <li><router-link to="/about">关于作者</router-link></li>
              <li><a href="#">RSS 订阅</a></li>
            </ul>
          </el-col>

          <!-- Contact -->
          <el-col :span="8" :xs="24" class="footer-col">
             <h4 class="footer-title">联系我们</h4>
             <ul class="contact-list">
               <li><el-icon><Message /></el-icon> contact@example.com</li>
               <li><el-icon><Location /></el-icon> ShangHai, China</li>
             </ul>
             <div class="social-links">
                <a href="#" class="social-item"><i class="fab fa-github"></i> G</a>
                <a href="#" class="social-item"><i class="fab fa-twitter"></i> T</a>
                <a href="#" class="social-item"><i class="fab fa-weixin"></i> W</a>
             </div>
          </el-col>
        </el-row>
      </div>
      
      <div class="footer-bottom">
        <div class="copyright">
           &copy; 2024 {{ siteName || 'Blog System' }}. All rights reserved. 
           <span class="divider">|</span> 
           <a href="#">京ICP备88888888号</a>
        </div>
        <div class="powered-by">
           Powered by GoFrame & Vue 3
        </div>
      </div>
    </el-footer>
  </el-container>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Menu, Message, Location, User, ArrowDown } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useSiteConfig } from '../store/site'
import { getCurrentUser } from '../api/user'
import 'element-plus/theme-chalk/display.css' // Helper classes for responsive

const router = useRouter()
const route = useRoute()
const mobileMenuVisible = ref(false)
const { siteName } = useSiteConfig()

// 用户信息
const userInfo = ref({
  username: '',
  nickname: '',
  avatar: ''
})

// 是否已登录
const isLoggedIn = computed(() => {
  const token = localStorage.getItem('token')
  return !!token && !!userInfo.value.username
})

// 显示名称（优先显示nickname，如果没有则显示username）
const displayName = computed(() => {
  return userInfo.value.nickname || userInfo.value.username || '用户'
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
    
    // 清空用户信息
    userInfo.value.username = ''
    userInfo.value.nickname = ''
    userInfo.value.avatar = ''
    
    ElMessage.success('已退出登录')
    
    // 如果在admin页面，跳转到首页
    if (route.path.startsWith('/admin')) {
      router.push('/')
    }
  } catch {
    // 用户取消
  }
}

// 页面加载时加载用户信息
onMounted(() => {
  loadUserInfo()
})

// 监听路由变化，重新加载用户信息
watch(() => route.path, () => {
  loadUserInfo()
})
</script>

<style scoped>
.layout-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  position: sticky;
  top: 0;
  z-index: 100;
  padding: 0;
  height: 60px;
}

.header-inner {
  max-width: 1200px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.logo {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.logo-icon {
  background: linear-gradient(135deg, #409eff 0%, #36d1dc 100%);
  color: white;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 18px;
  margin-right: 10px;
}

.logo-icon.small {
  width: 24px;
  height: 24px;
  font-size: 14px;
  border-radius: 6px;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
}

.nav {
  flex: 1;
  display: flex;
  justify-content: center;
}

.nav-menu {
  border-bottom: none !important;
  background: transparent;
}

.actions {
  display: flex;
  align-items: center;
}

.mobile-menu-toggle {
  cursor: pointer;
  color: #606266;
  display: flex;
  align-items: center;
}

.mobile-nav {
  padding: 20px 0;
}

.mobile-nav-menu {
  border-right: none;
}

.mobile-actions {
  padding: 20px;
}

.main-content {
  flex: 1;
  width: 100%;
  max-width: 1200px;
  margin: 20px auto;
  padding: 0 20px;
  overflow: visible; /* 允许子元素溢出，比如阴影 */
  box-sizing: border-box; /* 修复宽度溢出 */
}

/* Login Button Style - Outline Style */
.login-btn {
  background-color: transparent;
  border: 1px solid #409eff;
  color: #409eff;
  transition: all 0.3s;
  height: 34px;
  line-height: 1;
  padding: 0 20px;
}

.login-btn:hover, .login-btn:focus {
  background-color: #409eff;
  color: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

/* User Dropdown Style */
.user-dropdown {
  cursor: pointer;
}

.user-info-link {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 12px;
  border-radius: 20px;
  transition: all 0.3s;
  color: #303133;
}

.user-info-link:hover {
  background-color: #f5f7fa;
}

.username-text {
  font-size: 14px;
  font-weight: 500;
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Mobile User Info */
.mobile-user-info {
  padding: 16px 20px;
}

.mobile-user-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  background-color: #f5f7fa;
  cursor: pointer;
  transition: all 0.3s;
}

.mobile-user-card:hover {
  background-color: #e9ecef;
}

.mobile-user-details {
  flex: 1;
  min-width: 0;
}

.mobile-username {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.mobile-user-email {
  font-size: 12px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Footer Styling */
.footer {
  background-color: #232323;
  color: #a0a0a0;
  padding: 60px 0 20px;
  margin-top: auto;
  font-size: 14px;
  height: auto; /* 覆盖 element-plus 默认高度 */
}

.footer-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.footer-col {
  margin-bottom: 30px;
}

.footer-logo {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  color: white;
}

.logo-text-light {
  font-size: 18px;
  font-weight: 700;
  margin-left: 10px;
}

.footer-desc {
  line-height: 1.6;
  margin-bottom: 20px;
  max-width: 300px;
}

.footer-title {
  color: white;
  font-size: 16px;
  margin-bottom: 20px;
  font-weight: 600;
}

.footer-links-list, .contact-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.footer-links-list li, .contact-list li {
  margin-bottom: 12px;
}

.footer-links-list a {
  color: #a0a0a0;
  transition: color 0.3s;
  text-decoration: none;
}
.footer-links-list a:hover {
  color: #409eff;
}

.contact-list li {
  display: flex;
  align-items: center;
  gap: 8px;
}

.social-links {
  margin-top: 20px;
  display: flex;
  gap: 15px;
}

.social-item {
  width: 36px;
  height: 36px;
  background: rgba(255,255,255,0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  transition: all 0.3s;
  text-decoration: none;
}
.social-item:hover {
  background: #409eff;
  transform: translateY(-3px);
}

.footer-bottom {
  max-width: 1200px;
  margin: 40px auto 0;
  padding: 20px 20px 0;
  border-top: 1px solid rgba(255,255,255,0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
}

.copyright a {
  color: #a0a0a0;
  text-decoration: none;
}
.copyright a:hover {
  color: white;
}

.divider {
  margin: 0 10px;
}

/* 响应式辅助 */
@media (max-width: 768px) {
  .logo-text {
    font-size: 18px;
  }
  .main-content {
    padding: 0 10px;
  }
  .footer {
    padding: 40px 0 20px;
  }
  .footer-bottom {
    flex-direction: column;
    text-align: center;
    gap: 10px;
  }
}
</style>
