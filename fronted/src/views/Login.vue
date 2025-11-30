<template>
  <div class="login-container">
    <div class="login-wrapper">
      <!-- Left Side: Brand Area -->
      <div class="login-left">
        <div class="login-left-mask"></div>
        <div class="login-left-content">
          <div class="brand-logo">
            <span class="logo-icon">B</span>
          </div>
          <h2 class="brand-title">Blog System</h2>
          <p class="brand-desc">分享知识，激发创意。<br>探索代码的无限可能。</p>
        </div>
      </div>

      <!-- Right Side: Login Form -->
      <div class="login-right">
        <div class="login-box">
          <div class="login-header">
            <h2>Welcome Back</h2>
            <p>请输入您的账号密码进行登录</p>
          </div>
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
                @keyup.enter="handleLogin" 
                show-password
                class="custom-input"
              />
            </el-form-item>
            <div class="form-actions">
               <el-checkbox v-model="rememberMe" class="custom-checkbox">记住我</el-checkbox>
               <el-link type="primary" :underline="false" class="forgot-pwd">忘记密码？</el-link>
            </div>
            <el-form-item>
              <el-button type="primary" @click="handleLogin" class="submit-btn" :loading="loading">
                登 录
              </el-button>
            </el-form-item>
          </el-form>
          <div class="login-footer">
            <p>还没有账号？ <el-link type="primary" class="register-link" @click="goToRegister">立即注册</el-link></p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'

const router = useRouter()
const loading = ref(false)
const rememberMe = ref(false)
const form = reactive({
  username: '',
  password: ''
})

const handleLogin = () => {
  loading.value = true
  setTimeout(() => {
    // Mock login logic
    if (form.username === 'admin' && form.password === '123456') {
      ElMessage.success('登录成功')
      router.push('/admin/dashboard')
    } else {
      ElMessage.error('用户名或密码错误')
    }
    loading.value = false
  }, 800)
}

const goToRegister = () => {
  router.push('/register')
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f0f2f5;
  padding: 20px;
}

.login-wrapper {
  display: flex;
  width: 1000px;
  height: 600px;
  background: white;
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 20px 50px rgba(0,0,0,0.1);
}

/* Left Side */
.login-left {
  flex: 1;
  position: relative;
  background-image: url('https://images.unsplash.com/photo-1497294815431-9365093b7331?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1950&q=80');
  background-size: cover;
  background-position: center;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.login-left-mask {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(40, 50, 70, 0.6); /* Dark overlay */
  backdrop-filter: blur(2px);
}

.login-left-content {
  position: relative;
  z-index: 1;
  text-align: center;
  padding: 40px;
}

.brand-logo {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.logo-icon {
  width: 60px;
  height: 60px;
  background: white;
  color: #409eff;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: bold;
  box-shadow: 0 4px 15px rgba(0,0,0,0.2);
}

.brand-title {
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 16px;
  letter-spacing: 1px;
}

.brand-desc {
  font-size: 18px;
  line-height: 1.6;
  opacity: 0.9;
  font-weight: 300;
}

/* Right Side */
.login-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.login-box {
  width: 100%;
  max-width: 360px;
}

.login-header {
  margin-bottom: 40px;
}

.login-header h2 {
  margin: 0 0 12px;
  color: #303133;
  font-size: 32px;
  font-weight: 700;
}

.login-header p {
  margin: 0;
  color: #909399;
  font-size: 15px;
}

/* Input Styles */
.custom-input :deep(.el-input__wrapper) {
  background-color: #f5f7fa;
  box-shadow: none !important;
  border-radius: 8px;
  padding: 8px 15px;
  transition: all 0.3s;
}

.custom-input :deep(.el-input__wrapper.is-focus),
.custom-input :deep(.el-input__wrapper:hover) {
  background-color: #fff;
  box-shadow: 0 0 0 1px #409eff !important;
}

.custom-input :deep(.el-input__inner) {
  height: 36px;
}

/* Checkbox Styles */
.custom-checkbox {
  color: #606266;
  font-weight: 400;
}
.custom-checkbox :deep(.el-checkbox__inner) {
  border-radius: 4px;
}

/* Link Styles */
.forgot-pwd {
  font-size: 14px;
  color: #909399;
  transition: color 0.3s;
}
.forgot-pwd:hover {
  color: #409eff;
}

.register-link {
  font-weight: 600;
  margin-left: 4px;
  position: relative;
}
.register-link::after {
  content: '';
  position: absolute;
  width: 0;
  height: 2px;
  bottom: 0;
  left: 0;
  background-color: #409eff;
  transition: width 0.3s;
}
.register-link:hover::after {
  width: 100%;
}

.submit-btn {
  width: 100%;
  font-size: 16px;
  height: 48px;
  border-radius: 8px;
  margin-top: 10px;
  background: linear-gradient(135deg, #409eff 0%, #3a8ee6 100%);
  border: none;
  transition: all 0.3s;
  font-weight: 600;
  letter-spacing: 1px;
}

.submit-btn:hover, .submit-btn:focus {
  background: linear-gradient(135deg, #66b1ff 0%, #409eff 100%);
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(64, 158, 255, 0.3);
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.login-footer {
  margin-top: 24px;
  font-size: 14px;
  color: #606266;
  text-align: center;
}

/* Responsive */
@media (max-width: 900px) {
  .login-wrapper {
    flex-direction: column;
    width: 100%;
    max-width: 450px;
    height: auto;
  }
  
  .login-left {
    display: none;
  }
  
  .login-right {
    padding: 40px 30px;
  }
}
</style>
