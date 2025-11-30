<template>
  <div class="register-container">
    <div class="register-wrapper">
      <!-- Left Side: Brand Area -->
      <div class="register-left">
        <div class="register-left-mask"></div>
        <div class="register-left-content">
          <div class="brand-logo">
            <span class="logo-icon">B</span>
          </div>
          <h2 class="brand-title">Blog System</h2>
          <p class="brand-desc">加入我们，开启你的创作之旅。<br>分享知识，激发创意。</p>
        </div>
      </div>

      <!-- Right Side: Register Form -->
      <div class="register-right">
        <div class="register-box">
          <div class="register-header">
            <h2>创建账号</h2>
            <p>填写以下信息完成注册</p>
          </div>
          <el-form 
            ref="formRef"
            :model="form" 
            :rules="rules"
            class="register-form" 
            size="large"
          >
            <el-form-item prop="username">
              <el-input 
                v-model="form.username" 
                placeholder="用户名（3-20个字符）" 
                :prefix-icon="User"
                class="custom-input"
              />
            </el-form-item>
            <el-form-item prop="password">
              <el-input 
                v-model="form.password" 
                type="password" 
                placeholder="密码（6-20个字符）" 
                :prefix-icon="Lock"
                @keyup.enter="handleRegister" 
                show-password
                class="custom-input"
              />
            </el-form-item>
            <el-form-item prop="confirmPassword">
              <el-input 
                v-model="form.confirmPassword" 
                type="password" 
                placeholder="确认密码" 
                :prefix-icon="Lock"
                @keyup.enter="handleRegister" 
                show-password
                class="custom-input"
              />
            </el-form-item>
            <el-form-item prop="nickname">
              <el-input 
                v-model="form.nickname" 
                placeholder="昵称（可选）" 
                :prefix-icon="UserFilled"
                class="custom-input"
              />
            </el-form-item>
            <el-form-item prop="email">
              <el-input 
                v-model="form.email" 
                placeholder="邮箱" 
                :prefix-icon="Message"
                class="custom-input"
              />
            </el-form-item>
            <el-form-item>
              <el-button 
                type="primary" 
                @click="handleRegister" 
                class="submit-btn" 
                :loading="loading"
              >
                注 册
              </el-button>
            </el-form-item>
          </el-form>
          <div class="register-footer">
            <p>已有账号？ <el-link type="primary" @click="goToLogin" class="login-link">立即登录</el-link></p>
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
import { User, Lock, UserFilled, Message } from '@element-plus/icons-vue'
import { register } from '../api/user'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)
const form = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  nickname: '',
  email: ''
})

// 验证确认密码
const validateConfirmPassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度必须在3-20个字符之间', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度必须在6-20个字符之间', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    try {
      const res = await register({
        username: form.username,
        password: form.password,
        nickname: form.nickname || undefined,
        email: form.email
      })
      
      // 响应拦截器已经处理了成功响应，返回的是 data 字段
      // 如果请求成功，res 就是 {id: 16} 这样的数据
      // 如果请求失败，会抛出异常进入 catch 块
      if (res && res.id) {
        ElMessage.success('注册成功，请登录')
        router.push('/login-new')
      } else {
        ElMessage.error('注册失败')
      }
    } catch (error) {
      ElMessage.error(error.message || '注册失败，请稍后重试')
    } finally {
      loading.value = false
    }
  })
}

const goToLogin = () => {
  router.push('/login-new')
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f0f2f5;
  padding: 20px;
}

.register-wrapper {
  display: flex;
  width: 1000px;
  min-height: 600px;
  background: white;
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 20px 50px rgba(0,0,0,0.1);
}

/* Left Side */
.register-left {
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

.register-left-mask {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(40, 50, 70, 0.6);
  backdrop-filter: blur(2px);
}

.register-left-content {
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
.register-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.register-box {
  width: 100%;
  max-width: 400px;
}

.register-header {
  margin-bottom: 40px;
}

.register-header h2 {
  margin: 0 0 12px;
  color: #303133;
  font-size: 32px;
  font-weight: 700;
}

.register-header p {
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

/* Link Styles */
.login-link {
  font-weight: 600;
  margin-left: 4px;
  position: relative;
}
.login-link::after {
  content: '';
  position: absolute;
  width: 0;
  height: 2px;
  bottom: 0;
  left: 0;
  background-color: #409eff;
  transition: width 0.3s;
}
.login-link:hover::after {
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

.register-footer {
  margin-top: 24px;
  font-size: 14px;
  color: #606266;
  text-align: center;
}

/* Responsive */
@media (max-width: 900px) {
  .register-wrapper {
    flex-direction: column;
    width: 100%;
    max-width: 450px;
    min-height: auto;
  }
  
  .register-left {
    display: none;
  }
  
  .register-right {
    padding: 40px 30px;
  }
}
</style>

