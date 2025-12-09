<template>
  <div class="article-detail-container" v-loading="loading">
    <!-- Header Area -->
    <div class="article-header">
      <div class="header-content">
        <div class="meta-info">
          <el-tag effect="dark" round class="category-tag">{{ article.category }}</el-tag>
          <el-tag 
            v-if="article.topic" 
            effect="dark" 
            round 
            class="topic-tag"
            @click="goToTopic"
          >
            📚 {{ article.topic.name }}
          </el-tag>
          <span class="date">{{ article.createdAt }}</span>
        </div>
        <h1 class="title">{{ article.title }}</h1>
        <div class="author-info">
          <el-avatar 
            :size="32" 
            :src="article.publishedByUser?.avatar || 'https://picsum.photos/id/64/100/100'"
            @click="goToAuthorArticles(article.publishedByUser)"
            style="cursor: pointer;"
          />
          <span 
            class="author-name"
            @click="goToAuthorArticles(article.publishedByUser)"
            style="cursor: pointer;"
          >
            {{ article.publishedByUser ? (article.publishedByUser.nickname || article.publishedByUser.username) : 'Admin' }}
          </span>
          <span class="divider">·</span>
          <span class="read-time">5 min read</span>
          <span class="divider">·</span>
          <span class="views">{{ article.views }} views</span>
        </div>
      </div>
      <div class="header-image" :style="{ backgroundImage: `url(${article.cover})` }"></div>
    </div>

    <!-- Content Area -->
    <el-row justify="center" class="content-wrapper">
      <el-col :span="16" :xs="22">
        <!-- 专题信息 -->
        <el-card v-if="article.topic" class="topic-info-card" shadow="never">
          <div class="topic-info" @click="goToTopic">
            <div class="topic-icon">📚</div>
            <div class="topic-content">
              <div class="topic-label">本文属于专题</div>
              <div class="topic-name">{{ article.topic.name }}</div>
              <div class="topic-desc">{{ article.topic.description }}</div>
              <div class="topic-stats">
                <span>{{ article.topic.articleCount }} 篇文章</span>
                <span class="divider">·</span>
                <span>{{ article.topic.progress }}</span>
              </div>
            </div>
            <div class="topic-action">
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>
        </el-card>

        <el-card class="content-card" shadow="never">
          <div class="markdown-body" v-html="renderedContent"></div>
          
          <div class="article-footer">
            <div class="tags-list">
              <el-tag
                v-for="tag in article.tags"
                :key="typeof tag === 'object' ? tag.id : tag"
                class="tag-item"
                effect="plain"
                :style="typeof tag === 'object' && tag.color ? { borderColor: tag.color, color: tag.color } : {}"
                @click="goToTag(tag)"
                style="cursor: pointer;"
              >
                # {{ typeof tag === 'object' ? tag.name : tag }}
              </el-tag>
            </div>
            <div class="actions">
               <el-button type="primary" plain round size="small">
                  <el-icon><Star /></el-icon> 点赞
               </el-button>
               <el-button round size="small">
                  <el-icon><Share /></el-icon> 分享
               </el-button>
            </div>
          </div>
        </el-card>

        <!-- Comments Area -->
        <div class="comments-section">
          <h3>评论 ({{ commentTotal }})</h3>
          
          <!-- 评论输入框 -->
          <div v-if="isLoggedIn" class="comment-input-area">
            <div class="comment-input-header">
              <el-icon class="comment-icon"><ChatLineRound /></el-icon>
              <span class="comment-input-title">发表评论</span>
            </div>
            <el-input
              v-model="commentContent"
              type="textarea"
              :rows="4"
              placeholder="写下你的想法，与大家分享..."
              maxlength="1000"
              show-word-limit
              class="comment-textarea"
            />
            <div class="comment-actions">
              <el-button type="primary" @click="submitComment" :loading="submitting" :disabled="!commentContent.trim()">
                <el-icon v-if="!submitting"><Promotion /></el-icon>
                <span>{{ submitting ? '发表中...' : '发表评论' }}</span>
              </el-button>
            </div>
          </div>
          
          <!-- 未登录提示 -->
          <div v-else class="comment-login-prompt">
            <div class="login-prompt-card">
              <el-icon class="login-prompt-icon"><Lock /></el-icon>
              <div class="login-prompt-content">
                <p class="login-prompt-text">请先登录后即可发表评论</p>
                <el-button 
                  type="primary" 
                  @click="goToLogin"
                  class="login-prompt-button"
                >
                  立即登录
                </el-button>
              </div>
            </div>
          </div>

          <!-- 评论列表 -->
          <div class="comments-list" v-loading="commentsLoading">
            <div v-if="comments.length === 0 && !commentsLoading" class="empty-comments">
              <el-empty description="暂无评论，快来抢沙发吧！" />
            </div>
            
            <div v-for="comment in comments" :key="comment.id" class="comment-item">
              <div class="comment-header">
                <el-avatar 
                  :size="40" 
                  :src="comment.user?.avatar || 'https://picsum.photos/id/64/100/100'"
                />
                <div class="comment-user-info">
                  <div class="comment-user-name">
                    {{ comment.user ? (comment.user.nickname || comment.user.username) : '匿名用户' }}
                  </div>
                  <div class="comment-time">{{ formatTime(comment.createdAt) }}</div>
                </div>
              </div>
              <div class="comment-content">{{ comment.content }}</div>
              
              <!-- 回复按钮 -->
              <div v-if="isLoggedIn" class="comment-actions-bar">
                <button 
                  @click="showReplyInput(comment.id)"
                  class="reply-button"
                >
                  回复
                </button>
              </div>

              <!-- 回复输入框 -->
              <div v-if="replyToCommentId === comment.id" class="reply-input-area">
                <div class="reply-input-header">
                  <el-icon class="reply-icon"><ChatLineRound /></el-icon>
                  <span>回复评论</span>
                </div>
                <el-input
                  v-model="replyContent"
                  type="textarea"
                  :rows="2"
                  placeholder="写下你的回复..."
                  maxlength="1000"
                />
                <div class="reply-actions">
                  <el-button size="small" @click="cancelReply">取消</el-button>
                  <el-button 
                    type="primary" 
                    size="small" 
                    @click="submitReply(comment.id)" 
                    :loading="submitting"
                    :disabled="!replyContent.trim()"
                  >
                    <el-icon v-if="!submitting"><Promotion /></el-icon>
                    <span>{{ submitting ? '发表中...' : '发表回复' }}</span>
                  </el-button>
                </div>
              </div>

              <!-- 回复列表 -->
              <div v-if="comment.replies && comment.replies.length > 0" class="replies-list">
                <div v-for="reply in comment.replies" :key="reply.id" class="reply-item">
                  <div class="reply-header">
                    <el-avatar 
                      :size="36" 
                      :src="reply.user?.avatar || 'https://picsum.photos/id/64/100/100'"
                    />
                    <div class="reply-user-info">
                      <div class="reply-user-name">
                        {{ reply.user ? (reply.user.nickname || reply.user.username) : '匿名用户' }}
                      </div>
                      <div class="reply-time">{{ formatTime(reply.createdAt) }}</div>
                    </div>
                  </div>
                  <div class="reply-text">{{ reply.content }}</div>
                </div>
              </div>
            </div>

            <!-- 分页 -->
            <el-pagination
              v-if="commentTotal > 0"
              v-model:current-page="commentPage"
              v-model:page-size="commentPageSize"
              :total="commentTotal"
              :page-sizes="[10, 20, 50]"
              layout="total, prev, pager, next"
              @current-change="loadComments"
              @size-change="handlePageSizeChange"
              class="comment-pagination"
            />
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Star, Share, ArrowRight, ChatLineRound, Promotion, Lock, User } from '@element-plus/icons-vue'
import { marked } from 'marked'
import { getArticle } from '@/api/article'
import { createComment, getCommentList } from '@/api/comment'
import { ElMessage } from 'element-plus'
import { updateSEO, generateArticleSEO } from '@/utils/seo'
import { useSiteConfig } from '@/store/site'

const route = useRoute()
const router = useRouter()
const articleId = route.params.id
const { siteName, bannerSubtitle } = useSiteConfig()
const siteUrl = window.location.origin

const article = ref({
  id: articleId,
  title: '',
  category: '',
  createdAt: '',
  views: 0,
  cover: '',
  tags: [],
  content: ''
})

const loading = ref(true)

// 是否已登录
const isLoggedIn = computed(() => {
  const token = localStorage.getItem('token')
  return !!token
})

// 评论相关
const comments = ref([])
const commentTotal = ref(0)
const commentPage = ref(1)
const commentPageSize = ref(10)
const commentsLoading = ref(false)
const commentContent = ref('')
const replyContent = ref('')
const replyToCommentId = ref(0)
const submitting = ref(false)

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  try {
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return dateStr
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  } catch (error) {
    if (typeof dateStr === 'string' && /^\d{4}-\d{2}-\d{2}/.test(dateStr)) {
      return dateStr.split(' ')[0] || dateStr.split('T')[0]
    }
    return dateStr
  }
}

// 加载文章详情
const loadArticle = async () => {
  loading.value = true
  try {
    const res = await getArticle({ id: parseInt(articleId) })
    if (res) {
      // 处理tags：如果是对象数组，保留对象；如果是字符串数组，保留字符串；如果是字符串，转换为数组
      let tags = []
      if (res.tags) {
        if (Array.isArray(res.tags)) {
          // 已经是数组，直接使用
          tags = res.tags
        } else if (typeof res.tags === 'string') {
          // 字符串，按逗号分割
          tags = res.tags.split(',').map(t => t.trim()).filter(t => t)
        }
      }
      
      article.value = {
        id: res.id,
        title: res.title || '',
        category: res.categoryName || '',
        createdAt: formatDate(res.createdAt),
        views: res.views || 0,
        cover: res.coverImage || 'https://picsum.photos/id/' + res.id + '/1200/600',
        tags: tags,
        content: res.content || '',
        publishedByUser: res.publishedByUser || null,
        topic: res.topic || null,
        // 保存原始数据用于SEO
        desc: res.desc || '',
        coverImage: res.coverImage || '',
        updatedAt: res.updatedAt || res.createdAt
      }
      
      // 设置文章页面的SEO
      const seoConfig = generateArticleSEO({
        ...res,
        tags: tags,
        category: res.categoryName || ''
      }, siteName.value || 'Blog System', bannerSubtitle.value || '', siteUrl)
      
      updateSEO(seoConfig)
    }
  } catch (error) {
    console.error('加载文章失败:', error)
    ElMessage.error('加载文章失败')
  } finally {
    loading.value = false
  }
}

// 加载评论列表
const loadComments = async () => {
  commentsLoading.value = true
  try {
    const res = await getCommentList({
      articleId: parseInt(articleId),
      page: commentPage.value,
      size: commentPageSize.value
    })
    if (res) {
      comments.value = res.list || []
      commentTotal.value = res.total || 0
    }
  } catch (error) {
    console.error('加载评论失败:', error)
    ElMessage.error('加载评论失败')
  } finally {
    commentsLoading.value = false
  }
}

// 提交评论
const submitComment = async () => {
  if (!isLoggedIn.value) {
    ElMessage.warning('请先登录后再发表评论')
    goToLogin()
    return
  }

  if (!commentContent.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }

  submitting.value = true
  try {
    await createComment({
      articleId: parseInt(articleId),
      parentId: 0,
      content: commentContent.value.trim()
    })
    ElMessage.success('评论发表成功')
    commentContent.value = ''
    commentPage.value = 1
    loadComments()
  } catch (error) {
    console.error('发表评论失败:', error)
    const errorMsg = error.response?.data?.message || error.message || '发表评论失败'
    ElMessage.error(errorMsg)
    // 如果是未登录错误，跳转到登录页
    if (errorMsg.includes('登录')) {
      goToLogin()
    }
  } finally {
    submitting.value = false
  }
}

// 显示回复输入框
const showReplyInput = (commentId) => {
  if (!isLoggedIn.value) {
    ElMessage.warning('请先登录后再回复')
    goToLogin()
    return
  }
  replyToCommentId.value = commentId
  replyContent.value = ''
}

// 取消回复
const cancelReply = () => {
  replyToCommentId.value = 0
  replyContent.value = ''
}

// 提交回复
const submitReply = async (parentId) => {
  if (!isLoggedIn.value) {
    ElMessage.warning('请先登录后再回复')
    goToLogin()
    return
  }

  if (!replyContent.value.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }

  submitting.value = true
  try {
    await createComment({
      articleId: parseInt(articleId),
      parentId: parentId,
      content: replyContent.value.trim()
    })
    ElMessage.success('回复发表成功')
    cancelReply()
    loadComments()
  } catch (error) {
    console.error('发表回复失败:', error)
    const errorMsg = error.response?.data?.message || error.message || '发表回复失败'
    ElMessage.error(errorMsg)
    // 如果是未登录错误，跳转到登录页
    if (errorMsg.includes('登录')) {
      goToLogin()
    }
  } finally {
    submitting.value = false
  }
}

// 跳转到登录页
const goToLogin = () => {
  router.push({
    path: '/login',
    query: { redirect: route.fullPath }
  })
}

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now - date
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 分页大小改变
const handlePageSizeChange = () => {
  commentPage.value = 1
  loadComments()
}

// 处理专题链接点击
onMounted(() => {
  loadArticle()
  loadComments()
  nextTick(() => {
    document.addEventListener('click', handleLinkClick)
  })
})

const handleLinkClick = (e) => {
  const target = e.target.closest('a.topic-link')
  if (target) {
    e.preventDefault()
    const topicId = target.getAttribute('data-topic-id')
    if (topicId) {
      router.push(`/topic/${topicId}`)
    }
  }
}

// 跳转到作者文章列表
const goToAuthorArticles = (user) => {
  if (user && user.id) {
    router.push(`/author/${user.id}`)
  }
}

// 跳转到标签详情页
const goToTag = (tag) => {
  if (tag) {
    const slug = typeof tag === 'object' ? tag.slug : tag
    if (slug) {
      router.push(`/tag/${slug}`)
    }
  }
}

// 跳转到专题详情页
const goToTopic = () => {
  if (article.value.topic && article.value.topic.id) {
    router.push(`/topic/${article.value.topic.id}`)
  }
}

const renderedContent = computed(() => {
  if (!article.value.content) return ''
  let html = marked.parse(article.value.content)
  
  // 处理专题链接，确保使用路由导航而不是页面刷新
  html = html.replace(
    /<a href="\/topic\/(\d+)"/g, 
    '<a href="/topic/$1" class="topic-link" data-topic-id="$1"'
  )
  
  return html
})
</script>

<style scoped>
.article-detail-container {
  background-color: #f9f9f9;
  min-height: 100vh;
  padding-bottom: 60px;
}

.article-header {
  position: relative;
  height: 400px;
  background-color: #2c3e50;
  color: white;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}

.header-image {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-size: cover;
  background-position: center;
  opacity: 0.4;
  filter: blur(4px);
  transform: scale(1.1);
}

.header-content {
  position: relative;
  z-index: 1;
  text-align: center;
  max-width: 800px;
  padding: 0 20px;
}

.meta-info {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
}

.category-tag {
  background-color: #409eff;
  border-color: #409eff;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.title {
  font-size: 36px;
  font-weight: 800;
  line-height: 1.3;
  margin-bottom: 24px;
  text-shadow: 0 2px 10px rgba(0,0,0,0.3);
}

.author-info {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  opacity: 0.9;
  flex-wrap: wrap; /* 允许换行 */
}

.author-name {
  margin-left: 8px;
  font-weight: 600;
}

.divider {
  margin: 0 10px;
  opacity: 0.6;
}

.content-wrapper {
  margin-top: -60px;
  position: relative;
  z-index: 2;
}

.content-card {
  padding: 40px;
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.05);
  background: white;
  border: none;
}

/* Markdown Styles Override - Simple Version */
.markdown-body {
  font-size: 16px;
  line-height: 1.8;
  color: #333;
  overflow-wrap: break-word; /* 防止长单词溢出 */
}
:deep(.markdown-body h2) {
  margin-top: 30px;
  margin-bottom: 15px;
  font-size: 24px;
  font-weight: 700;
  color: #1a1a1a;
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
}
:deep(.markdown-body p) {
  margin-bottom: 16px;
}
:deep(.markdown-body pre) {
  background: #f6f8fa;
  padding: 16px;
  border-radius: 8px;
  overflow: auto;
  font-family: monospace;
}
:deep(.markdown-body code) {
  background: #f0f2f5;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.9em;
  color: #d63384;
}
:deep(.markdown-body pre code) {
  background: transparent;
  color: inherit;
  padding: 0;
}
:deep(.markdown-body ul) {
  padding-left: 20px;
  margin-bottom: 16px;
}
:deep(.markdown-body img) {
    max-width: 100%;
    height: auto;
    border-radius: 8px;
}
:deep(.markdown-body blockquote) {
  margin: 20px 0;
  padding: 16px 20px;
  background: #f0f9ff;
  border-left: 4px solid #409eff;
  border-radius: 4px;
  color: #606266;
}
:deep(.markdown-body a.topic-link) {
  color: #409eff;
  text-decoration: none;
  font-weight: 600;
  border-bottom: 1px dashed #409eff;
  transition: all 0.3s ease;
}
:deep(.markdown-body a.topic-link:hover) {
  color: #66b1ff;
  border-bottom-style: solid;
}


.article-footer {
  margin-top: 40px;
  padding-top: 20px;
  border-top: 1px solid #f0f2f5;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 20px;
}

.tags-list {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.comments-section {
  margin-top: 40px;
}

.comments-section h3 {
  margin-bottom: 24px;
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  display: flex;
  align-items: center;
  gap: 8px;
}

.comments-section h3::before {
  content: '';
  width: 4px;
  height: 20px;
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  border-radius: 2px;
}

.comment-input-area {
  margin-bottom: 32px;
  padding: 20px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #e4e7ed;
  transition: all 0.3s ease;
}

.comment-input-area:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border-color: #e4e7ed;
  transform: translateY(-2px);
}

.comment-input-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f2f5;
}

.comment-icon {
  font-size: 18px;
  color: #409eff;
}

.comment-input-title {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.comment-input-area :deep(.el-textarea__inner) {
  padding: 14px 18px;
  font-size: 15px;
  line-height: 1.8;
  resize: none;
  color: #303133;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
  min-height: 120px;
}

.comment-input-area :deep(.el-textarea__inner)::placeholder {
  color: #909399;
  font-size: 14px;
}

.comment-input-area :deep(.el-input__count) {
  color: #909399;
  font-size: 12px;
  background: transparent;
  bottom: 8px;
  right: 12px;
}

.comment-actions {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.comment-actions .el-button {
  padding: 10px 24px;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.comment-actions .el-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.comments-list {
  margin-top: 24px;
}

.empty-comments {
  padding: 60px 0;
  text-align: center;
}

.comment-item {
  padding: 24px;
  margin-bottom: 16px;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  border: 1px solid #f0f2f5;
  transition: all 0.3s ease;
}

.comment-item:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border-color: #e4e7ed;
  transform: translateY(-2px);
}

.comment-header {
  display: flex;
  align-items: flex-start;
  margin-bottom: 16px;
}

.comment-header .el-avatar {
  flex-shrink: 0;
  border: 2px solid #f0f2f5;
  transition: all 0.3s ease;
}

.comment-item:hover .comment-header .el-avatar {
  border-color: #409eff;
  transform: scale(1.05);
}

.comment-user-info {
  margin-left: 14px;
  flex: 1;
  min-width: 0;
}

.comment-user-name {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 6px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.comment-user-name::after {
  content: '';
  width: 4px;
  height: 4px;
  background: #c0c4cc;
  border-radius: 50%;
}

.comment-time {
  font-size: 12px;
  color: #909399;
  display: flex;
  align-items: center;
  gap: 4px;
}

.comment-content {
  font-size: 15px;
  color: #606266;
  line-height: 1.8;
  margin-left: 54px;
  margin-bottom: 16px;
  word-break: break-word;
  white-space: pre-wrap;
}

.comment-actions-bar {
  margin-left: 54px;
  margin-bottom: 12px;
}

.reply-button {
  padding: 0;
  font-size: 13px;
  font-weight: 500;
  color: #909399;
  background: none;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  line-height: 1.5;
}

.reply-button:hover {
  color: #409eff;
}

.reply-input-area {
  margin-left: 54px;
  margin-top: 16px;
  padding: 16px;
  background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
  border-radius: 10px;
  border: 1px solid #e4e7ed;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
  animation: slideDown 0.3s ease;
}

.reply-input-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 12px;
  font-size: 14px;
  font-weight: 600;
  color: #606266;
}

.reply-icon {
  font-size: 16px;
  color: #409eff;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.reply-input-area :deep(.el-textarea__inner) {
  padding: 12px 16px;
  font-size: 14px;
  line-height: 1.7;
  resize: none;
  color: #303133;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
}

.reply-input-area :deep(.el-textarea__inner)::placeholder {
  color: #909399;
  font-size: 13px;
}

.reply-actions {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.reply-actions .el-button {
  padding: 8px 20px;
  border-radius: 6px;
  font-size: 13px;
}

.replies-list {
  margin-left: 54px;
  margin-top: 20px;
  padding-left: 24px;
  border-left: 2px solid #e4e7ed;
  position: relative;
}

.replies-list::before {
  content: '';
  position: absolute;
  left: -2px;
  top: 0;
  width: 2px;
  height: 24px;
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  border-radius: 0 1px 1px 0;
}

.reply-item {
  margin-bottom: 16px;
  padding: 16px;
  background: #ffffff;
  border-radius: 10px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
  border: 1px solid #f0f2f5;
  transition: all 0.3s ease;
}

.reply-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border-color: #e4e7ed;
  transform: translateX(2px);
}

.reply-item:last-child {
  margin-bottom: 0;
}

.reply-header {
  display: flex;
  align-items: flex-start;
  margin-bottom: 12px;
}

.reply-header .el-avatar {
  flex-shrink: 0;
  border: 2px solid #f0f2f5;
  transition: all 0.3s ease;
}

.reply-item:hover .reply-header .el-avatar {
  border-color: #409eff;
  transform: scale(1.05);
}

.reply-user-info {
  margin-left: 12px;
  flex: 1;
  min-width: 0;
}

.reply-user-name {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.reply-user-name::after {
  content: '';
  width: 3px;
  height: 3px;
  background: #c0c4cc;
  border-radius: 50%;
}

.reply-time {
  font-size: 12px;
  color: #909399;
  display: flex;
  align-items: center;
  gap: 4px;
}

.reply-text {
  font-size: 14px;
  color: #606266;
  line-height: 1.7;
  word-break: break-word;
  white-space: pre-wrap;
  margin-left: 48px;
}

.comment-pagination {
  margin-top: 32px;
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

.comment-login-prompt {
  margin-bottom: 32px;
}

.login-prompt-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e4e7ed;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  transition: all 0.3s ease;
}

.login-prompt-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border-color: #c0c4cc;
}

.login-prompt-icon {
  font-size: 24px;
  color: #909399;
  flex-shrink: 0;
}

.login-prompt-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.login-prompt-text {
  font-size: 14px;
  color: #606266;
  margin: 0;
}

.login-prompt-button {
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 500;
  font-size: 14px;
  flex-shrink: 0;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .article-header {
      height: 300px;
  }
  .title {
      font-size: 24px;
  }
  .content-card {
      padding: 20px;
  }
  .article-footer {
      flex-direction: column;
      align-items: flex-start;
  }
  .actions {
      width: 100%;
      display: flex;
      justify-content: space-between;
  }
  .meta-info {
      flex-wrap: wrap;
  }
  
  /* 评论区域响应式 */
  .comments-section h3 {
      font-size: 18px;
  }
  
  .comment-input-area {
      padding: 16px;
  }
  
  .comment-item {
      padding: 16px;
      margin-bottom: 12px;
  }
  
  .comment-content {
      margin-left: 0;
      margin-top: 12px;
  }
  
  .comment-actions-bar {
      margin-left: 0;
  }
  
  .reply-input-area {
      margin-left: 0;
  }
  
  .replies-list {
      margin-left: 0;
      padding-left: 12px;
  }
  
  .login-prompt-card {
      flex-direction: column;
      text-align: center;
      padding: 20px;
      gap: 16px;
  }
  
  .login-prompt-content {
      flex-direction: column;
      width: 100%;
  }
  
  .login-prompt-button {
      width: 100%;
  }
}
</style>
