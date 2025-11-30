<template>
  <div class="about-container">
    <div class="about-hero">
      <div class="avatar-wrapper">
         <img :src="aboutInfo.avatar" alt="Avatar" />
      </div>
      <h1>关于我</h1>
      <p class="subtitle">{{ aboutInfo.subtitle }}</p>
    </div>

    <el-row justify="center" class="about-content">
       <el-col :span="16" :xs="22">
          <el-card shadow="never" class="about-card">
             <div class="section">
                <h2>👋 Hello, World</h2>
                <div class="markdown-content" v-html="renderMarkdown(aboutInfo.content)"></div>
             </div>

             <div class="section">
                <h2>🛠️ 技术栈</h2>
                <div class="tech-stack">
                   <el-tag size="large">Golang</el-tag>
                   <el-tag size="large" type="success">Vue 3</el-tag>
                   <el-tag size="large" type="warning">MySQL</el-tag>
                   <el-tag size="large" type="danger">Redis</el-tag>
                   <el-tag size="large" type="info">Docker</el-tag>
                   <el-tag size="large">Kubernetes</el-tag>
                </div>
             </div>

             <div class="section">
                <h2>📫 联系方式</h2>
                <p>如果你对我的文章感兴趣，或者有任何技术问题想要交流，欢迎通过以下方式联系我：</p>
                <ul class="contact-list">
                   <li v-if="aboutInfo.email">
                     <el-icon><Message /></el-icon> 
                     Email: <a :href="`mailto:${aboutInfo.email}`">{{ aboutInfo.email }}</a>
                   </li>
                   <li v-if="aboutInfo.github">
                     <el-icon><img src="https://unpkg.com/simple-icons@v9/icons/github.svg" width="16" /></el-icon> 
                     GitHub: 
                     <a 
                       :href="formatGithubUrl(aboutInfo.github)" 
                       target="_blank"
                       rel="noopener noreferrer"
                     >
                       {{ aboutInfo.github }}
                     </a>
                   </li>
                </ul>
             </div>
          </el-card>
       </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@element-plus/icons-vue'
import { getAbout } from '@/api/settings'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'

const aboutInfo = ref({
  avatar: 'https://picsum.photos/id/64/200/200',
  title: '关于我',
  subtitle: 'Full Stack Developer / Open Source Enthusiast',
  content: '我是 Admin，一名热爱技术的全栈开发者。热衷于探索 Go 语言的高性能并发模型，也喜欢 Vue 3 带来的丝滑前端体验。\n\n这个博客系统是我用来记录技术成长、分享编程心得的小天地。',
  email: 'admin@example.com',
  github: ''
})

// 配置 marked
marked.setOptions({
  breaks: true, // 支持 GitHub 风格的换行
  gfm: true, // 启用 GitHub 风格的 Markdown
})

// 渲染 Markdown 内容
const renderMarkdown = (content) => {
  if (!content) return ''
  try {
    return marked.parse(content)
  } catch (error) {
    console.error('Markdown 渲染错误:', error)
    return content.replace(/\n/g, '<br>')
  }
}

// 格式化 GitHub URL
const formatGithubUrl = (github) => {
  if (!github) return '#'
  if (github.startsWith('http://') || github.startsWith('https://')) {
    return github
  }
  if (github.startsWith('@')) {
    return `https://github.com/${github.substring(1)}`
  }
  return `https://github.com/${github}`
}

// 加载关于我信息
const loadAbout = async () => {
  try {
    const res = await getAbout()
    if (res) {
      aboutInfo.value = {
        avatar: res.avatar || aboutInfo.value.avatar,
        subtitle: res.subtitle || aboutInfo.value.subtitle,
        content: res.content || aboutInfo.value.content,
        email: res.email || aboutInfo.value.email,
        github: res.github || ''
      }
    }
  } catch (error) {
    console.error('加载关于我信息失败:', error)
    // 使用默认值，不显示错误提示
  }
}

onMounted(() => {
  loadAbout()
})
</script>

<style scoped>
.about-container {
  min-height: 100vh;
  padding-bottom: 60px;
}

.about-hero {
  background: linear-gradient(135deg, #2c3e50 0%, #4ca1af 100%);
  color: white;
  text-align: center;
  padding: 80px 20px 120px;
  clip-path: polygon(0 0, 100% 0, 100% 85%, 0 100%);
}

.avatar-wrapper img {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: 4px solid rgba(255,255,255,0.3);
  margin-bottom: 20px;
}

.about-hero h1 {
  font-size: 36px;
  margin: 0 0 10px;
}

.subtitle {
  font-size: 18px;
  opacity: 0.9;
}

.about-content {
  margin-top: -60px;
}

.about-card {
  padding: 40px;
  border-radius: 16px;
  border: none;
  box-shadow: 0 10px 30px rgba(0,0,0,0.05);
}

.section {
  margin-bottom: 40px;
}

.section h2 {
  font-size: 24px;
  color: #303133;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.section p {
  line-height: 1.8;
  color: #606266;
  font-size: 16px;
}

.markdown-content {
  line-height: 1.8;
  color: #606266;
  font-size: 16px;
}

.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3),
.markdown-content :deep(h4),
.markdown-content :deep(h5),
.markdown-content :deep(h6) {
  color: #303133;
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
}

.markdown-content :deep(h1) {
  font-size: 28px;
}

.markdown-content :deep(h2) {
  font-size: 24px;
}

.markdown-content :deep(h3) {
  font-size: 20px;
}

.markdown-content :deep(p) {
  margin-bottom: 16px;
  line-height: 1.8;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  margin-bottom: 16px;
  padding-left: 24px;
}

.markdown-content :deep(li) {
  margin-bottom: 8px;
}

.markdown-content :deep(code) {
  background-color: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  color: #e83e8c;
}

.markdown-content :deep(pre) {
  background-color: #f5f7fa;
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin-bottom: 16px;
}

.markdown-content :deep(pre code) {
  background-color: transparent;
  padding: 0;
  color: #303133;
}

.markdown-content :deep(blockquote) {
  border-left: 4px solid #409eff;
  padding-left: 16px;
  margin: 16px 0;
  color: #606266;
  font-style: italic;
}

.markdown-content :deep(a) {
  color: #409eff;
  text-decoration: none;
}

.markdown-content :deep(a:hover) {
  text-decoration: underline;
}

.markdown-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 16px 0;
}

.markdown-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 16px 0;
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
  border: 1px solid #ebeef5;
  padding: 8px 12px;
  text-align: left;
}

.markdown-content :deep(th) {
  background-color: #f5f7fa;
  font-weight: 600;
}

.tech-stack {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.contact-list {
  list-style: none;
  padding: 0;
}

.contact-list li {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
  color: #606266;
  font-size: 16px;
}
</style>

