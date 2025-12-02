<template>
  <div class="chat-container">
    <el-card shadow="never" class="chat-card">
      <template #header>
        <div class="chat-header">
          <div class="header-left">
            <el-icon class="chat-icon"><ChatLineRound /></el-icon>
            <span class="chat-title">AI 智能对话</span>
          </div>
          <div class="header-right">
            <el-button 
              type="danger" 
              size="small" 
              :disabled="!currentSessionId || loading"
              @click="handleClearSession"
            >
              <el-icon><Delete /></el-icon>
              清除会话
            </el-button>
          </div>
        </div>
      </template>

      <div class="chat-content">
        <!-- 消息列表 -->
        <div class="messages-container" ref="messagesContainer">
          <div v-if="messages.length === 0" class="empty-message">
            <el-icon class="empty-icon"><ChatDotRound /></el-icon>
            <p>开始与 AI 对话吧！</p>
          </div>
          
          <div
            v-for="(message, index) in messages"
            :key="index"
            :class="['message-item', message.role === 'user' ? 'user-message' : 'ai-message']"
          >
            <div class="message-avatar">
              <el-avatar v-if="message.role === 'user'" :size="32">
                <el-icon><User /></el-icon>
              </el-avatar>
              <el-avatar v-else :size="32" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
                <el-icon><Service /></el-icon>
              </el-avatar>
            </div>
            <div class="message-content">
              <div class="message-text" v-html="formatMessage(message.content)"></div>
              <div class="message-time">{{ formatTime(message.timestamp) }}</div>
            </div>
          </div>

          <!-- 加载中提示 -->
          <div v-if="loading" class="message-item ai-message">
            <div class="message-avatar">
              <el-avatar :size="32" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
                <el-icon><Service /></el-icon>
              </el-avatar>
            </div>
            <div class="message-content">
              <div class="message-text loading-text">
                <el-icon class="is-loading"><Loading /></el-icon>
                <span>AI 正在思考中...</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="input-container">
          <el-input
            v-model="inputMessage"
            type="textarea"
            :rows="3"
            placeholder="输入您的问题..."
            :disabled="loading"
            @keydown.ctrl.enter="handleSend"
            @keydown.meta.enter="handleSend"
            class="message-input"
          />
          <div class="input-actions">
            <div class="input-tips">
              <span>按 Ctrl+Enter 或 Cmd+Enter 发送</span>
            </div>
            <el-button
              type="primary"
              :loading="loading"
              :disabled="!inputMessage.trim()"
              @click="handleSend"
              class="send-button"
            >
              <el-icon v-if="!loading"><Promotion /></el-icon>
              {{ loading ? '发送中...' : '发送' }}
            </el-button>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  ChatLineRound, 
  ChatDotRound, 
  User, 
  Service, 
  Loading, 
  Promotion,
  Delete
} from '@element-plus/icons-vue'
import { sendChatMessage, sendChatMessageStream, clearChatSession } from '../../api/chat'
import { marked } from 'marked'

// 消息列表
const messages = ref([])
const inputMessage = ref('')
const loading = ref(false)
const currentSessionId = ref('')
const messagesContainer = ref(null)

// 发送消息
const handleSend = async () => {
  const message = inputMessage.value.trim()
  if (!message || loading.value) {
    return
  }

  // 添加用户消息
  const userMessage = {
    role: 'user',
    content: message,
    timestamp: Date.now()
  }
  messages.value.push(userMessage)
  inputMessage.value = ''
  
  // 滚动到底部
  scrollToBottom()

  // 发送到后端（流式）
  loading.value = true
  
  // 创建AI消息占位符
  const aiMessageIndex = messages.value.length
  const aiMessage = {
    role: 'assistant',
    content: '',
    timestamp: Date.now()
  }
  messages.value.push(aiMessage)
  scrollToBottom()

  try {
    await sendChatMessageStream({
      message: message,
      sessionId: currentSessionId.value || undefined
    }, (chunk) => {
      // 实时更新AI回复内容
      console.log('收到数据块:', chunk)
      if (messages.value[aiMessageIndex]) {
        messages.value[aiMessageIndex].content += chunk
        scrollToBottom()
      }
    }, (sessionId) => {
      // 更新会话ID
      console.log('会话完成，sessionId:', sessionId)
      if (sessionId) {
        currentSessionId.value = sessionId
      }
    })
  } catch (error) {
    console.error('发送消息失败:', error)
    // 移除AI消息（如果发送失败）
    if (messages.value[aiMessageIndex] && messages.value[aiMessageIndex].content === '') {
      messages.value.splice(aiMessageIndex, 1)
    }
    // 移除用户消息（如果发送失败）
    if (messages.value.length > 0 && messages.value[messages.value.length - 1].role === 'user') {
      messages.value.pop()
    }
    // 检查是否是超时错误
    if (error.message && error.message.includes('timeout')) {
      ElMessage.error('请求超时，AI响应时间较长，请稍后重试')
    } else {
      ElMessage.error(error.message || '发送消息失败，请稍后重试')
    }
  } finally {
    loading.value = false
    // 确保AI消息有时间戳
    if (messages.value[aiMessageIndex]) {
      messages.value[aiMessageIndex].timestamp = Date.now()
    }
  }
}

// 清除会话
const handleClearSession = async () => {
  if (!currentSessionId.value) {
    return
  }

  try {
    await ElMessageBox.confirm(
      '确定要清除当前会话吗？清除后将无法恢复对话历史。',
      '确认清除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await clearChatSession({
      sessionId: currentSessionId.value
    })

    messages.value = []
    currentSessionId.value = ''
    ElMessage.success('会话已清除')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '清除会话失败')
    }
  }
}

// 格式化消息（支持Markdown）
const formatMessage = (content) => {
  try {
    return marked.parse(content)
  } catch (error) {
    return content.replace(/\n/g, '<br>')
  }
}

// 格式化时间
const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date

  if (diff < 60000) {
    return '刚刚'
  } else if (diff < 3600000) {
    return `${Math.floor(diff / 60000)}分钟前`
  } else if (diff < 86400000) {
    return `${Math.floor(diff / 3600000)}小时前`
  } else {
    return date.toLocaleString('zh-CN', {
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

onMounted(() => {
  scrollToBottom()
})
</script>

<style scoped>
.chat-container {
  height: calc(100vh - 140px);
  min-height: 600px;
  display: flex;
  flex-direction: column;
}

.chat-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.chat-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 0;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.chat-icon {
  font-size: 20px;
  color: #409eff;
}

.chat-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.chat-content {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 200px);
  min-height: 500px;
  overflow: hidden;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;
  margin-bottom: 0;
  min-height: 0;
}

.empty-message {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 20px;
  opacity: 0.5;
}

.empty-message p {
  font-size: 16px;
  margin: 0;
}

.message-item {
  display: flex;
  margin-bottom: 20px;
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.user-message {
  flex-direction: row-reverse;
}

.user-message .message-content {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  margin-right: 12px;
}

.ai-message .message-content {
  background: white;
  color: #303133;
  margin-left: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.message-avatar {
  flex-shrink: 0;
}

.message-content {
  max-width: 70%;
  padding: 12px 16px;
  border-radius: 12px;
  word-wrap: break-word;
}

.message-text {
  line-height: 1.6;
  font-size: 14px;
}

.message-text :deep(p) {
  margin: 0 0 8px 0;
}

.message-text :deep(p:last-child) {
  margin-bottom: 0;
}

.message-text :deep(code) {
  background: rgba(0, 0, 0, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
}

.message-text :deep(pre) {
  background: rgba(0, 0, 0, 0.1);
  padding: 12px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 8px 0;
}

.message-text :deep(pre code) {
  background: none;
  padding: 0;
}

.message-text :deep(ul),
.message-text :deep(ol) {
  margin: 8px 0;
  padding-left: 24px;
}

.message-text :deep(li) {
  margin: 4px 0;
}

.loading-text {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #909399;
}

.message-time {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  margin-top: 8px;
}

.ai-message .message-time {
  color: #909399;
}

.input-container {
  border-top: 1px solid #e4e7ed;
  padding: 16px;
  background: white;
  flex-shrink: 0;
}

.message-input {
  margin-bottom: 12px;
}

.input-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.input-tips {
  font-size: 12px;
  color: #909399;
}

.send-button {
  min-width: 100px;
}

/* 滚动条样式 */
.messages-container::-webkit-scrollbar {
  width: 6px;
}

.messages-container::-webkit-scrollbar-track {
  background: transparent;
}

.messages-container::-webkit-scrollbar-thumb {
  background: #c0c4cc;
  border-radius: 3px;
}

.messages-container::-webkit-scrollbar-thumb:hover {
  background: #a0a4ac;
}

/* 响应式 */
@media (max-width: 768px) {
  .chat-container {
    height: calc(100vh - 120px);
  }

  .chat-content {
    height: calc(100vh - 180px);
    min-height: 400px;
  }

  .message-content {
    max-width: 85%;
  }

  .input-tips {
    display: none;
  }
}
</style>

