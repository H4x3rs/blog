<template>
  <div class="login-container">
    <!-- 动态背景 -->
    <div class="bg-animation">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
      <div class="shape shape-4"></div>
    </div>

    <!-- 装饰性圆点 -->
    <div class="decorative-dots">
      <span v-for="i in 50" :key="i" class="dot" :style="getDotStyle(i)"></span>
    </div>

    <!-- 登录卡片 -->
    <div class="login-card">
      <!-- Logo 和标题区域 -->
      <div class="login-header">
        <div class="logo-wrapper">
          <div class="logo-icon">
            <span class="logo-letter">B</span>
            <div class="logo-ring"></div>
          </div>
        </div>
        <h1 class="main-title">欢迎回来</h1>
      </div>

      <!-- 登录表单 -->
      <el-form :model="form" class="login-form" size="large">
        <el-form-item>
          <el-input 
            v-model="form.username" 
            placeholder="用户名 / admin"
            :prefix-icon="User"
            class="custom-input"
          />
        </el-form-item>

        <el-form-item>
          <el-input 
            v-model="form.password" 
            type="password"
            placeholder="密码 / 123456"
            :prefix-icon="Lock"
            show-password
            class="custom-input"
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <div class="form-options">
          <el-checkbox v-model="rememberMe" class="remember-checkbox">
            <span class="checkbox-label">记住登录状态</span>
          </el-checkbox>
          <el-link type="primary" :underline="false" class="forgot-link">
            忘记密码？
          </el-link>
        </div>

        <el-button 
          type="primary" 
          class="login-button" 
          :loading="loading"
          @click="handleLogin"
        >
          <span v-if="!loading">立即登录</span>
          <span v-else>登录中...</span>
        </el-button>
      </el-form>

      <!-- 分隔线 -->
      <div class="divider">
        <span class="divider-text">或使用以下方式登录</span>
      </div>

      <!-- 第三方登录 -->
      <div class="social-login">
        <button class="social-btn github" @click="handleOAuthLogin('github')" :disabled="oauthLoading">
          <svg class="icon" viewBox="0 0 1024 1024" width="20" height="20">
            <path d="M512 42.666667A464.64 464.64 0 0 0 42.666667 502.186667 460.373333 460.373333 0 0 0 363.52 938.666667c23.466667 4.266667 32-9.813333 32-22.186667v-78.08c-130.56 27.733333-158.293333-61.44-158.293333-61.44a122.026667 122.026667 0 0 0-52.053334-67.413333c-42.666667-28.16 3.413333-27.733333 3.413334-27.733334a98.56 98.56 0 0 1 71.68 47.36 101.12 101.12 0 0 0 136.533333 37.973334 99.413333 99.413333 0 0 1 29.866667-61.44c-104.106667-11.52-213.333333-50.773333-213.333334-226.986667a177.066667 177.066667 0 0 1 47.36-124.16 161.28 161.28 0 0 1 4.693334-121.173333s39.68-12.373333 128 46.933333a455.68 455.68 0 0 1 234.666666 0c89.6-59.306667 128-46.933333 128-46.933333a161.28 161.28 0 0 1 4.693334 121.173333A177.066667 177.066667 0 0 1 810.666667 477.866667c0 176.64-110.08 215.466667-213.333334 226.986666a106.666667 106.666667 0 0 1 32 85.333334v125.866666c0 14.933333 8.533333 26.88 32 22.186667A460.8 460.8 0 0 0 981.333333 502.186667 464.64 464.64 0 0 0 512 42.666667" fill="currentColor"/>
          </svg>
          <span>GitHub</span>
        </button>
        <button class="social-btn google" @click="handleOAuthLogin('google')" :disabled="oauthLoading">
          <svg class="icon" viewBox="0 0 1024 1024" width="20" height="20">
            <path d="M214.101333 512c0-32.512 5.546667-63.701333 15.36-92.928L57.173333 290.218667A491.861333 491.861333 0 0 0 4.693333 512c0 79.701333 18.858667 154.88 52.394667 221.610667l172.202667-129.066667A290.56 290.56 0 0 1 214.101333 512" fill="#FBBC05"/>
            <path d="M516.693333 216.192c72.106667 0 136.533333 26.624 187.221334 70.144l140.629333-140.629333C747.605333 59.690667 638.762667 20.608 516.693333 20.608c-202.325333 0-376.234667 113.28-459.52 278.613333l172.373334 128.853334c39.68-118.016 152.832-211.882667 287.146666-211.882667" fill="#EA4335"/>
            <path d="M516.693333 807.808c-134.357333 0-247.509333-93.866667-287.232-211.882667L57.088 724.992C140.416 890.325333 314.325333 1003.605333 516.693333 1003.605333c117.76 0 225.28-38.4 300.928-104.106666l-146.688-113.152c-40.96 27.392-93.653333 43.52-154.24 43.52" fill="#34A853"/>
            <path d="M1005.397333 512c0-29.568-4.693333-61.44-11.648-91.008H516.650667V614.4h274.602666c-13.696 65.962667-51.072 116.650667-104.533333 149.632l146.645333 113.152c85.845333-79.36 140.032-195.84 140.032-365.184" fill="#4285F4"/>
          </svg>
          <span>Google</span>
        </button>
        <button class="social-btn wechat" @click="handleOAuthLogin('wechat')" :disabled="oauthLoading">
          <svg class="icon" viewBox="0 0 1024 1024" width="20" height="20">
            <path d="M664.250054 368.541681c10.015098 0 19.892049 0.732687 29.67281 1.795902-26.647917-122.810047-159.358451-214.077703-310.826188-214.077703-169.353083 0-308.085774 114.232694-308.085774 259.274068 0 83.708416 46.165436 152.460695 123.281791 205.78483l-30.80868 91.730191 107.688651-53.455469c38.558178 7.53665 69.459978 15.308661 107.924012 15.308661 9.66308 0 19.230958-0.470721 28.752858-1.225921-6.025227-20.36584-9.521864-41.723264-9.521864-63.862493C402.328693 476.632491 517.908058 368.541681 664.250054 368.541681zM498.62897 285.87389c23.200398 0 38.557154 15.120372 38.557154 38.061874 0 22.846334-15.356756 38.156018-38.557154 38.156018-23.107277 0-46.260603-15.309684-46.260603-38.156018C452.368366 300.994262 475.522716 285.87389 498.62897 285.87389zM283.016307 362.090758c-23.107277 0-46.402843-15.309684-46.402843-38.156018 0-22.941502 23.295566-38.061874 46.402843-38.061874 23.081695 0 38.46301 15.120372 38.46301 38.061874C321.479317 346.782098 306.098002 362.090758 283.016307 362.090758zM945.448458 606.151333c0-121.888048-123.258255-221.236753-261.683954-221.236753-146.57838 0-262.015505 99.348706-262.015505 221.236753 0 122.06508 115.437126 221.200938 262.015505 221.200938 30.66644 0 61.617359-7.609305 92.423993-15.262685l84.513836 45.786813-23.178909-76.17082C899.379213 735.776599 945.448458 674.90216 945.448458 606.151333zM598.803483 567.994292c-15.332197 0-30.807656-15.096836-30.807656-30.501688 0-15.190981 15.47546-30.477129 30.807656-30.477129 23.295566 0 38.558178 15.286148 38.558178 30.477129C637.361661 552.897456 622.099049 567.994292 598.803483 567.994292zM768.25071 567.994292c-15.213493 0-30.594809-15.096836-30.594809-30.501688 0-15.190981 15.381315-30.477129 30.594809-30.477129 23.107277 0 38.558178 15.286148 38.558178 30.477129C806.808888 552.897456 791.357987 567.994292 768.25071 567.994292z" fill="currentColor"/>
          </svg>
          <span>微信</span>
        </button>
      </div>

      <!-- 底部提示 -->
      <div class="login-footer">
        <p class="footer-text">
          还没有账号？
          <el-link type="primary" class="register-link" :underline="false" @click="goToRegister">
            立即注册
          </el-link>
        </p>
      </div>
    </div>

  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, CircleClose } from '@element-plus/icons-vue'
import { login } from '@/api/user'
import { getOAuthLoginUrl, oauthCallback } from '@/api/user'

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const oauthLoading = ref(false)
const rememberMe = ref(false)

const form = reactive({
  username: '',
  password: ''
})

// 处理普通登录
const handleLogin = async () => {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  
  loading.value = true
  try {
    const res = await login({
      username: form.username,
      password: form.password
    })
    
    if (res && res.token) {
      // 保存token和用户信息到localStorage
      localStorage.setItem('token', res.token)
      localStorage.setItem('username', res.username)
      if (res.nickname) {
        localStorage.setItem('nickname', res.nickname)
      }
      if (res.avatar) {
        localStorage.setItem('avatar', res.avatar)
      }
      ElMessage.success('登录成功！')
      // 检查是否有重定向路径
      const redirect = route.query.redirect || '/admin/dashboard'
      router.push(redirect)
    } else {
      ElMessage.error('登录失败，请检查用户名和密码')
    }
  } catch (error) {
    ElMessage.error(error.message || '登录失败，请重试')
  } finally {
    loading.value = false
  }
}

// 处理OAuth登录
const handleOAuthLogin = async (provider) => {
  oauthLoading.value = true
  try {
    const res = await getOAuthLoginUrl(provider)
    if (res && res.redirectUrl) {
      // 跳转到第三方登录页面
      window.location.href = res.redirectUrl
    } else {
      ElMessage.error('获取登录链接失败')
      oauthLoading.value = false
    }
  } catch (error) {
    ElMessage.error(error.message || '登录失败，请重试')
    oauthLoading.value = false
  }
}

// 跳转到注册页面
const goToRegister = () => {
  router.push('/register')
}

// 页面加载时的初始化逻辑（如果需要）
onMounted(() => {
  // 可以在这里添加其他初始化逻辑
})

// 生成装饰性圆点样式
const getDotStyle = (index) => {
  const size = Math.random() * 4 + 2
  return {
    width: `${size}px`,
    height: `${size}px`,
    left: `${Math.random() * 100}%`,
    top: `${Math.random() * 100}%`,
    animationDelay: `${Math.random() * 3}s`,
    animationDuration: `${Math.random() * 3 + 2}s`
  }
}
</script>

<style scoped>
/* 全局优化 */
* {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* 容器 */
.login-container {
  position: relative;
  min-height: 100vh;
  min-height: 100dvh; /* 动态视口高度，避免移动端地址栏影响 */
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: #f8f9fa; /* 与首页一致的浅色背景 */
  padding: 20px;
  -webkit-tap-highlight-color: transparent; /* 移除 iOS 点击高亮 */
}

/* 动态背景 */
.bg-animation {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.5;
  animation: float 20s infinite ease-in-out;
}

.shape-1 {
  width: 500px;
  height: 500px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  top: -250px;
  left: -250px;
  animation-delay: 0s;
}

.shape-2 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  bottom: -200px;
  right: -200px;
  animation-delay: 5s;
}

.shape-3 {
  width: 300px;
  height: 300px;
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  top: 50%;
  right: -150px;
  animation-delay: 10s;
}

.shape-4 {
  width: 350px;
  height: 350px;
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  bottom: 30%;
  left: -175px;
  animation-delay: 15s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
  }
  33% {
    transform: translate(30px, -50px) rotate(120deg);
  }
  66% {
    transform: translate(-20px, 20px) rotate(240deg);
  }
}

/* 装饰性圆点 */
.decorative-dots {
  position: absolute;
  width: 100%;
  height: 100%;
}

.dot {
  position: absolute;
  background: rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  animation: twinkle 3s infinite ease-in-out;
}

@keyframes twinkle {
  0%, 100% {
    opacity: 0.3;
    transform: scale(1);
  }
  50% {
    opacity: 0.8;
    transform: scale(1.2);
  }
}

/* 登录卡片 */
.login-card {
  position: relative;
  width: 100%;
  max-width: 550px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px); /* Safari 支持 */
  border-radius: 24px;
  padding: 48px 50px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: slideUp 0.6s ease-out;
  z-index: 10;
  will-change: transform, opacity; /* 性能优化 */
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Logo 和标题 */
.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo-wrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 24px;
}

.logo-icon {
  position: relative;
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
  animation: logoFloat 3s infinite ease-in-out;
}

@keyframes logoFloat {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.logo-letter {
  font-size: 40px;
  font-weight: 800;
  color: white;
  z-index: 2;
}

.logo-ring {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-radius: 20px;
  animation: pulse 2s infinite ease-in-out;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.1);
    opacity: 0.5;
  }
}

.main-title {
  font-size: 32px;
  font-weight: 700;
  color: #303133;
  margin: 0 0 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  font-size: 15px;
  color: #909399;
  margin: 0;
}

/* 表单 */
.login-form {
  margin-bottom: 24px;
}

.login-form :deep(.el-form-item) {
  margin-bottom: 20px;
}

.login-form :deep(.el-form-item__content) {
  line-height: normal;
}

/* 输入框样式（参照旧版） */
.custom-input :deep(.el-input__wrapper) {
  background-color: #fff;
  box-shadow: 0 0 0 1px #e4e7ed !important;
  border-radius: 12px;
  padding: 8px 16px;
  transition: all 0.3s;
}

.custom-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #d0d3d9 !important;
}

.custom-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #667eea !important;
}

.custom-input :deep(.el-input__inner) {
  height: 40px;
  font-size: 15px;
  color: #303133;
}

.custom-input :deep(.el-input__inner::placeholder) {
  color: #c0c4cc;
  font-size: 14px;
}

/* 表单选项 */
.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.remember-checkbox {
  color: #606266;
  font-weight: 400;
}

.checkbox-label {
  font-size: 14px;
  user-select: none;
}

.remember-checkbox :deep(.el-checkbox__inner) {
  border-radius: 4px;
}

.remember-checkbox :deep(.el-checkbox__label) {
  color: #606266;
  font-size: 14px;
}

.forgot-link {
  font-size: 14px;
  color: #667eea;
  font-weight: 500;
  transition: all 0.3s;
}

.forgot-link:hover {
  color: #764ba2;
}

/* 登录按钮 */
.login-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  letter-spacing: 1px;
  transition: all 0.3s;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.login-button:hover,
.login-button:focus {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-button:active {
  transform: translateY(0);
}

.login-button.is-loading {
  transform: none;
}

/* 分隔线 */
.divider {
  position: relative;
  text-align: center;
  margin: 32px 0;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: #e4e7ed;
}

.divider-text {
  position: relative;
  display: inline-block;
  padding: 0 16px;
  background: rgba(255, 255, 255, 0.95);
  color: #909399;
  font-size: 13px;
}

/* 社交登录 */
.social-login {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 28px;
}

.social-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 16px 12px;
  border: 2px solid #e4e7ed;
  border-radius: 12px;
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 13px;
  font-weight: 500;
  color: inherit;
}

.social-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.social-btn.github {
  color: #333;
}

.social-btn.github:hover {
  border-color: #333;
  background: #333;
  color: white;
}

.social-btn.google {
  color: #4285F4;
}

.social-btn.google:hover {
  border-color: #4285F4;
  background: #4285F4;
  color: white;
}

.social-btn.wechat {
  color: #07C160;
}

.social-btn.wechat:hover {
  border-color: #07C160;
  background: #07C160;
  color: white;
}

.social-btn .icon {
  transition: transform 0.3s ease;
  flex-shrink: 0;
}

.social-btn:hover .icon {
  transform: scale(1.1);
}

.social-btn span {
  transition: color 0.3s ease;
}

/* 底部 */
.login-footer {
  text-align: center;
}

.footer-text {
  font-size: 14px;
  color: #606266;
  margin: 0;
}

.register-link {
  font-weight: 600;
  color: #667eea;
}

.register-link:hover {
  color: #764ba2;
}


/* 响应式设计 */

/* 大屏幕 (1920px+) */
@media (min-width: 1920px) {
  .login-card {
    max-width: 650px;
    padding: 56px 60px;
  }

  .main-title {
    font-size: 36px;
  }

  .subtitle {
    font-size: 16px;
  }

  .logo-icon {
    width: 90px;
    height: 90px;
  }

  .logo-letter {
    font-size: 45px;
  }
}

/* 平板横屏 (768px - 1024px) */
@media (max-width: 1024px) and (min-width: 769px) {
  .login-card {
    max-width: 520px;
  }

  .shape {
    filter: blur(60px);
  }
}

/* 平板竖屏 (481px - 768px) */
@media (max-width: 768px) and (min-width: 481px) {
  .login-container {
    padding: 16px;
  }

  .login-card {
    padding: 40px 36px;
    max-width: 500px;
  }

  .main-title {
    font-size: 28px;
  }

  .subtitle {
    font-size: 15px;
  }

  .logo-icon {
    width: 70px;
    height: 70px;
  }

  .logo-letter {
    font-size: 36px;
  }

  .social-login {
    grid-template-columns: 1fr;
    gap: 10px;
  }

  .custom-input :deep(.el-input__wrapper) {
    height: 50px;
  }

  .login-button {
    height: 50px;
  }

  .shape {
    filter: blur(60px);
  }
}

/* 手机横屏 (568px - 812px 高度较小) */
@media (max-width: 812px) and (max-height: 480px) and (orientation: landscape) {
  .login-container {
    padding: 20px 40px;
    min-height: auto;
  }

  .login-card {
    padding: 24px 36px;
    max-width: 700px;
  }

  .login-header {
    margin-bottom: 20px;
  }

  .logo-wrapper {
    margin-bottom: 12px;
  }

  .logo-icon {
    width: 50px;
    height: 50px;
  }

  .logo-letter {
    font-size: 26px;
  }

  .main-title {
    font-size: 22px;
    margin-bottom: 6px;
  }

  .subtitle {
    font-size: 13px;
  }

  .login-form :deep(.el-form-item) {
    margin-bottom: 12px;
  }

  .custom-input :deep(.el-input__wrapper) {
    height: 44px;
  }

  .login-button {
    height: 44px;
    font-size: 14px;
  }

  .form-options {
    margin: 6px 0 16px;
  }

  .divider {
    margin: 16px 0;
  }

  .social-login {
    grid-template-columns: repeat(3, 1fr);
    gap: 8px;
    margin-bottom: 16px;
  }

  .social-btn {
    padding: 10px 8px;
    font-size: 11px;
    gap: 4px;
  }

  .social-btn .icon {
    width: 16px;
    height: 16px;
  }

  .bg-animation {
    display: none;
  }

  .decorative-dots {
    display: none;
  }
}

/* 手机竖屏 (320px - 480px) */
@media (max-width: 480px) {
  .login-container {
    padding: 12px;
  }

  .login-card {
    padding: 32px 24px;
    max-width: 100%;
    border-radius: 20px;
  }

  .login-header {
    margin-bottom: 28px;
  }

  .logo-wrapper {
    margin-bottom: 16px;
  }

  .logo-icon {
    width: 64px;
    height: 64px;
  }

  .logo-letter {
    font-size: 32px;
  }

  .main-title {
    font-size: 24px;
  }

  .subtitle {
    font-size: 14px;
  }

  .login-form :deep(.el-form-item) {
    margin-bottom: 16px;
  }

  .custom-input :deep(.el-input__wrapper) {
    height: 48px;
  }

  .custom-input :deep(.el-input__inner) {
    font-size: 14px;
  }

  .input-prefix-icon {
    font-size: 16px;
  }

  .login-button {
    height: 48px;
    font-size: 15px;
  }

  .form-options {
    margin: 6px 0 20px;
  }

  .remember-checkbox :deep(.el-checkbox__label),
  .forgot-link {
    font-size: 13px;
  }

  .divider {
    margin: 24px 0;
  }

  .divider-text {
    font-size: 12px;
    padding: 0 12px;
  }

  .social-login {
    grid-template-columns: 1fr;
    gap: 10px;
    margin-bottom: 20px;
  }

  .social-btn {
    flex-direction: row;
    justify-content: center;
    padding: 14px 16px;
    font-size: 14px;
    gap: 10px;
  }

  .login-footer {
    margin-top: 20px;
  }

  .footer-text {
    font-size: 13px;
  }

  .shape {
    filter: blur(50px);
  }
}

/* 小屏手机 (< 375px) */
@media (max-width: 374px) {
  .login-card {
    padding: 28px 20px;
  }

  .main-title {
    font-size: 22px;
  }

  .subtitle {
    font-size: 13px;
  }

  .logo-icon {
    width: 56px;
    height: 56px;
  }

  .logo-letter {
    font-size: 28px;
  }

  .custom-input :deep(.el-input__wrapper) {
    height: 46px;
  }

  .login-button {
    height: 46px;
    font-size: 14px;
  }

  .social-btn {
    padding: 12px 14px;
    font-size: 13px;
  }
}

/* 超小屏 (< 320px) */
@media (max-width: 320px) {
  .login-card {
    padding: 24px 16px;
    border-radius: 16px;
  }

  .main-title {
    font-size: 20px;
  }
}

/* 横屏优化 (所有尺寸) */
@media (orientation: landscape) and (max-height: 600px) {
  .login-container {
    min-height: 100vh;
    padding: 16px;
  }

  .login-card {
    margin: 0 auto;
  }
}

/* 触摸设备优化 */
@media (hover: none) and (pointer: coarse) {
  .social-btn,
  .login-button,
  .forgot-link {
    min-height: 44px; /* iOS 推荐的最小触摸区域 */
  }

  .remember-checkbox :deep(.el-checkbox__inner) {
    width: 20px;
    height: 20px;
  }
}

/* 打印样式 */
@media print {
  .bg-animation,
  .decorative-dots,
  .social-login {
    display: none;
  }

  .login-container {
    background: white;
  }

  .login-card {
    box-shadow: none;
    border: 1px solid #e4e7ed;
  }
}
</style>

