<template>
  <el-container class="layout-container">
    <el-header class="header">
      <div class="header-inner">
        <NuxtLink to="/" class="logo">
          <span class="logo-icon">B</span>
          <span class="logo-text">{{ siteStore.siteName }}</span>
        </NuxtLink>
        
        <!-- Desktop Nav -->
        <div class="nav hidden-xs-only">
          <el-menu 
            mode="horizontal" 
            :default-active="route.path"
            :ellipsis="false"
            class="nav-menu"
            router
          >
            <el-menu-item index="/">首页</el-menu-item>
            <el-menu-item index="/category">分类</el-menu-item>
            <el-menu-item index="/tag">标签</el-menu-item>
            <el-menu-item index="/topic">专题</el-menu-item>
            <el-menu-item index="/about">关于</el-menu-item>
          </el-menu>
        </div>

        <div class="actions hidden-xs-only">
          <template v-if="isLoggedIn">
            <el-dropdown @command="handleCommand">
              <div class="user-info">
                <el-avatar :size="32" :src="userAvatar">
                  <el-icon v-if="!userAvatar"><User /></el-icon>
                </el-avatar>
                <span class="username">{{ displayName }}</span>
                <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item command="myArticles">我的文章</el-dropdown-item>
                  <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <NuxtLink to="/login">
              <el-button class="login-btn" round>登录 / 注册</el-button>
            </NuxtLink>
          </template>
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
            router
            :default-active="route.path"
            class="mobile-nav-menu"
          >
            <el-menu-item index="/" @click="mobileMenuVisible = false">首页</el-menu-item>
            <el-menu-item index="/category" @click="mobileMenuVisible = false">分类</el-menu-item>
            <el-menu-item index="/tag" @click="mobileMenuVisible = false">标签</el-menu-item>
            <el-menu-item index="/topic" @click="mobileMenuVisible = false">专题</el-menu-item>
            <el-menu-item index="/about" @click="mobileMenuVisible = false">关于</el-menu-item>
          </el-menu>
          <div class="mobile-actions">
            <template v-if="isLoggedIn">
              <div class="mobile-user-info">
                <el-avatar :size="40" :src="userAvatar">
                  <el-icon v-if="!userAvatar"><User /></el-icon>
                </el-avatar>
                <div class="mobile-user-details">
                  <div class="mobile-username">{{ displayName }}</div>
                </div>
              </div>
              <el-button class="mobile-btn" @click="handleCommand('logout'); mobileMenuVisible = false">退出登录</el-button>
            </template>
            <template v-else>
              <NuxtLink to="/login" @click="mobileMenuVisible = false">
                <el-button class="login-btn">登录 / 注册</el-button>
              </NuxtLink>
            </template>
          </div>
       </div>
    </el-drawer>
    
    <el-main class="main-content">
      <slot />
    </el-main>
    
    <el-footer class="footer">
      <div class="footer-inner">
        <el-row :gutter="40">
          <!-- Logo & Intro -->
          <el-col :span="8" :xs="24" class="footer-col">
            <div class="footer-logo">
              <span class="logo-icon small">B</span>
              <span class="logo-text-light">{{ siteStore.siteName }}</span>
            </div>
            <p class="footer-desc">
              {{ siteStore.bannerSubtitle }}
            </p>
          </el-col>

          <!-- Quick Links -->
          <el-col :span="8" :xs="24" class="footer-col">
            <h4 class="footer-title">快速链接</h4>
            <ul class="footer-links-list">
              <li><NuxtLink to="/">首页</NuxtLink></li>
              <li><NuxtLink to="/category">文章分类</NuxtLink></li>
              <li><NuxtLink to="/about">关于作者</NuxtLink></li>
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
                <a href="https://github.com" target="_blank" class="social-item">G</a>
                <a href="#" class="social-item">T</a>
                <a href="#" class="social-item">W</a>
             </div>
          </el-col>
        </el-row>
      </div>
      
      <div class="footer-bottom">
        <div class="copyright">
           &copy; {{ new Date().getFullYear() }} {{ siteStore.siteName }}. All rights reserved. 
           <template v-if="siteStore.icpNumber">
             <span class="divider">|</span> 
             <a href="https://beian.miit.gov.cn/" target="_blank">{{ siteStore.icpNumber }}</a>
           </template>
        </div>
        <div class="powered-by">
           Powered by Nuxt 3 & GoFrame
        </div>
      </div>
    </el-footer>
  </el-container>
</template>

<script setup lang="ts">
import { Menu, Message, Location, User, ArrowDown } from '@element-plus/icons-vue'
import { useSiteStore } from '~/stores/site'

const route = useRoute()
const router = useRouter()
const siteStore = useSiteStore()
const mobileMenuVisible = ref(false)

// 用户信息
const userInfo = ref({
  username: '',
  nickname: '',
  avatar: ''
})

// 检查是否已登录
const isLoggedIn = computed(() => {
  if (import.meta.server) return false
  const token = localStorage.getItem('token')
  return !!token && !!userInfo.value.username
})

// 显示名称
const displayName = computed(() => {
  return userInfo.value.nickname || userInfo.value.username || '用户'
})

// 用户头像
const userAvatar = computed(() => {
  return userInfo.value.avatar || ''
})

// 加载用户信息
const loadUserInfo = () => {
  if (import.meta.server) return
  
  const token = localStorage.getItem('token')
  if (!token) {
    userInfo.value = { username: '', nickname: '', avatar: '' }
    return
  }

  userInfo.value = {
    username: localStorage.getItem('username') || '',
    nickname: localStorage.getItem('nickname') || '',
    avatar: localStorage.getItem('avatar') || ''
  }
}

// 处理下拉菜单命令
const handleCommand = async (command: string) => {
  if (command === 'logout') {
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    localStorage.removeItem('nickname')
    localStorage.removeItem('avatar')
    userInfo.value = { username: '', nickname: '', avatar: '' }
    ElMessage.success('已退出登录')
    router.push('/')
  } else if (command === 'profile') {
    router.push('/admin/profile')
  } else if (command === 'myArticles') {
    router.push('/my-articles')
  }
}

// 初始化
onMounted(() => {
  loadUserInfo()
  siteStore.loadConfig()
})
</script>

<style scoped lang="scss">
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
  text-decoration: none;
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

  &.small {
    width: 24px;
    height: 24px;
    font-size: 14px;
    border-radius: 6px;
  }
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
  gap: 12px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 12px;
  border-radius: 20px;
  transition: all 0.3s;

  &:hover {
    background-color: #f5f7fa;
  }
}

.username {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}

.dropdown-icon {
  font-size: 12px;
  color: #909399;
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

.mobile-user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  margin-bottom: 16px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.mobile-username {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.mobile-btn {
  width: 100%;
  margin-bottom: 12px;
}

.main-content {
  flex: 1;
  width: 100%;
  max-width: 1200px;
  margin: 20px auto;
  padding: 0 20px;
  overflow: visible;
  box-sizing: border-box;
}

.login-btn {
  background-color: transparent;
  border: 1px solid #409eff;
  color: #409eff;
  transition: all 0.3s;
  height: 34px;
  line-height: 1;
  padding: 0 20px;

  &:hover, &:focus {
    background-color: #409eff;
    color: white;
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  }
}

// Footer
.footer {
  background-color: #232323;
  color: #a0a0a0;
  padding: 60px 0 20px;
  margin-top: auto;
  font-size: 14px;
  height: auto;
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

  &:hover {
    color: #409eff;
  }
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

  &:hover {
    background: #409eff;
    transform: translateY(-3px);
  }
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

  &:hover {
    color: white;
  }
}

.divider {
  margin: 0 10px;
}

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
