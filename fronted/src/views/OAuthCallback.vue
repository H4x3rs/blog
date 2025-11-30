<template>
  <div class="oauth-callback-container">
    <div class="loading-wrapper">
      <div class="loading-spinner">
        <el-icon class="is-loading" :size="48">
          <Loading />
        </el-icon>
      </div>
      <p class="loading-text">{{ loadingText }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { oauthCallback } from '@/api/user'

const router = useRouter()
const route = useRoute()
const loadingText = ref('正在处理登录...')

onMounted(() => {
  handleOAuthCallback()
})

// 处理OAuth回调 - 从URL获取code，调用后端API
const handleOAuthCallback = async () => {
  const code = route.query.code
  const state = route.query.state
  const error = route.query.error
  const errorDescription = route.query.error_description

  // 如果有错误信息（GitHub可能返回error参数）
  if (error || errorDescription) {
    loadingText.value = '登录失败'
    ElMessage.error(errorDescription || error || 'OAuth登录失败')
    setTimeout(() => {
      router.replace('/login-new')
    }, 2000)
    return
  }

  // 如果没有code，可能是直接访问了这个页面
  if (!code) {
    loadingText.value = '缺少授权码'
    ElMessage.warning('请先进行登录')
    setTimeout(() => {
      router.replace('/login-new')
    }, 2000)
    return
  }

  // 从路由参数获取provider
  const provider = route.params.provider
  
  loadingText.value = '正在验证登录信息...'
  
  try {
    // 调用后端API处理OAuth回调
    const res = await oauthCallback(provider, code, state)
    
    if (res && res.token) {
      loadingText.value = '登录成功，正在跳转...'
      
      // 保存token和用户信息
      localStorage.setItem('token', res.token)
      localStorage.setItem('username', res.username)
      if (res.nickname) {
        localStorage.setItem('nickname', res.nickname)
      }
      if (res.avatar) {
        localStorage.setItem('avatar', res.avatar)
      }
      localStorage.setItem('loginProvider', provider)
      
      ElMessage.success('登录成功！')
      
      // 跳转到管理页面
      setTimeout(() => {
        router.replace('/admin/dashboard')
      }, 500)
    } else {
      throw new Error('登录失败：未获取到token')
    }
  } catch (error) {
    loadingText.value = '登录失败'
    ElMessage.error(error.message || '登录失败，请重试')
    setTimeout(() => {
      router.replace('/login-new')
    }, 2000)
  }
}
</script>

<style scoped>
.oauth-callback-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8f9fa;
}

.loading-wrapper {
  text-align: center;
}

.loading-spinner {
  margin-bottom: 24px;
}

.loading-spinner .el-icon {
  color: #409eff;
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.loading-text {
  font-size: 16px;
  color: #606266;
  margin: 0;
}
</style>

