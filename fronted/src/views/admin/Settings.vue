<template>
  <div class="page-container">
    <el-card shadow="never" class="settings-card">
      <template #header>
        <div class="card-header">
          <el-icon><Setting /></el-icon>
          <span>系统设置</span>
        </div>
      </template>

      <el-tabs v-model="activeTab" class="settings-tabs">
        <!-- 基本设置 -->
        <el-tab-pane label="基本设置" name="basic">
          <el-form 
            ref="formRef" 
            :model="form" 
            :rules="rules" 
            label-width="120px" 
            class="settings-form"
          >
            <!-- 网站信息 -->
            <el-card shadow="never" class="form-section">
              <template #header>
                <div class="section-header">
                  <el-icon><InfoFilled /></el-icon>
                  <span>网站信息</span>
                </div>
              </template>
              
              <el-form-item label="网站名称" prop="siteName">
                <el-input 
                  v-model="form.siteName" 
                  placeholder="请输入网站名称"
                  maxlength="50"
                  show-word-limit
                  clearable
                />
              </el-form-item>

              <el-form-item label="网站标题" prop="title">
                <el-input 
                  v-model="form.title" 
                  placeholder="请输入网站标题（Banner 标题）"
                  maxlength="100"
                  show-word-limit
                  clearable
                />
              </el-form-item>

              <el-form-item label="网站副标题" prop="subtitle">
                <el-input 
                  v-model="form.subtitle" 
                  type="textarea"
                  :rows="2"
                  placeholder="请输入网站副标题"
                  maxlength="200"
                  show-word-limit
                  clearable
                />
              </el-form-item>

              <el-form-item label="网站描述" prop="siteDesc">
                <el-input 
                  v-model="form.siteDesc" 
                  type="textarea"
                  :rows="4"
                  placeholder="请输入网站描述"
                  maxlength="500"
                  show-word-limit
                  clearable
                />
              </el-form-item>
            </el-card>

            <!-- SEO 设置 -->
            <el-card shadow="never" class="form-section">
              <template #header>
                <div class="section-header">
                  <el-icon><Search /></el-icon>
                  <span>SEO 设置</span>
                </div>
              </template>

              <el-form-item label="关键词" prop="keywords">
                <el-input 
                  v-model="form.keywords" 
                  placeholder="请输入关键词，多个关键词用逗号分隔"
                  maxlength="200"
                  show-word-limit
                  clearable
                />
              </el-form-item>
            </el-card>

            <!-- 其他设置 -->
            <el-card shadow="never" class="form-section">
              <template #header>
                <div class="section-header">
                  <el-icon><Document /></el-icon>
                  <span>其他设置</span>
                </div>
              </template>

              <el-form-item label="底部版权" prop="copyright">
                <el-input 
                  v-model="form.copyright" 
                  placeholder="请输入版权信息"
                  maxlength="200"
                  show-word-limit
                  clearable
                />
              </el-form-item>

              <el-form-item label="ICP备案号" prop="icpNumber">
                <el-input 
                  v-model="form.icpNumber" 
                  placeholder="请输入ICP备案号，例如：京ICP备2021032569号"
                  maxlength="50"
                  show-word-limit
                  clearable
                />
                <div class="form-tip">ICP备案号将显示在网站底部</div>
              </el-form-item>
            </el-card>

            <!-- 操作按钮 -->
            <div class="form-actions">
              <el-button type="primary" :loading="saving" @click="handleSave">
                <el-icon><Check /></el-icon>
                保存设置
              </el-button>
              <el-button @click="handleReset">
                <el-icon><Refresh /></el-icon>
                重置
              </el-button>
            </div>
          </el-form>
        </el-tab-pane>

        <!-- 关于我 -->
        <el-tab-pane label="关于我" name="about">
          <el-form 
            ref="aboutFormRef" 
            :model="form" 
            :rules="aboutRules" 
            label-width="120px" 
            class="settings-form"
          >
            <el-card shadow="never" class="form-section">
              <template #header>
                <div class="section-header">
                  <el-icon><User /></el-icon>
                  <span>关于我页面设置</span>
                </div>
              </template>

              <el-form-item label="头像URL" prop="aboutAvatar">
                <el-row :gutter="20" class="avatar-row">
                  <el-col :span="16" :xs="24">
                    <el-input 
                      v-model="form.aboutAvatar" 
                      placeholder="请输入头像图片URL"
                      maxlength="500"
                      clearable
                    />
                  </el-col>
                  <el-col :span="8" :xs="24">
                    <div class="avatar-preview">
                      <div class="preview-label">头像预览</div>
                      <el-avatar 
                        :size="80" 
                        :src="form.aboutAvatar || 'https://picsum.photos/id/64/200/200'"
                        shape="circle"
                      >
                        <el-icon><UserFilled /></el-icon>
                      </el-avatar>
                    </div>
                  </el-col>
                </el-row>
              </el-form-item>

              <el-form-item label="副标题" prop="aboutSubtitle">
                <el-input 
                  v-model="form.aboutSubtitle" 
                  placeholder="Full Stack Developer / Open Source Enthusiast"
                  maxlength="100"
                  clearable
                />
              </el-form-item>

              <el-form-item label="个人简介" prop="aboutContent">
                <div class="markdown-editor-wrapper">
                  <div ref="vditorRef" id="vditor-about-content"></div>
                </div>
              </el-form-item>

              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="邮箱" prop="aboutEmail">
                    <el-input 
                      v-model="form.aboutEmail" 
                      placeholder="admin@example.com"
                      maxlength="100"
                      clearable
                    >
                      <template #prepend>
                        <el-icon><Message /></el-icon>
                      </template>
                    </el-input>
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="GitHub" prop="aboutGithub">
                    <el-input 
                      v-model="form.aboutGithub" 
                      placeholder="@username 或 GitHub URL"
                      maxlength="200"
                      clearable
                    >
                      <template #prepend>
                        <el-icon><Link /></el-icon>
                      </template>
                    </el-input>
                  </el-form-item>
                </el-col>
              </el-row>

              <!-- 操作按钮 -->
              <div class="form-actions">
                <el-button type="primary" :loading="saving" @click="handleSave">
                  <el-icon><Check /></el-icon>
                  保存设置
                </el-button>
                <el-button @click="handleReset">
                  <el-icon><Refresh /></el-icon>
                  重置
                </el-button>
                <el-button text type="primary" @click="handlePreviewAbout">
                  <el-icon><View /></el-icon>
                  预览页面
                </el-button>
              </div>
            </el-card>
          </el-form>
        </el-tab-pane>
        
        <!-- 邮件服务 -->
        <el-tab-pane label="邮件服务" name="mail">
          <el-card shadow="never" class="form-section">
            <template #header>
              <div class="section-header">
                <el-icon><Message /></el-icon>
                <span>邮件服务配置</span>
              </div>
            </template>
            <el-form label-width="120px" class="settings-form">
              <el-form-item label="SMTP 服务器">
                <el-input placeholder="smtp.example.com" />
              </el-form-item>
              <el-form-item label="端口">
                <el-input placeholder="465" />
              </el-form-item>
              <el-form-item label="邮箱账号">
                <el-input placeholder="admin@example.com" />
              </el-form-item>
              <el-form-item label="邮箱密码">
                <el-input type="password" show-password />
              </el-form-item>
              <el-form-item>
                <el-button type="primary">保存设置</el-button>
                <el-button>测试邮件</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Setting, 
  InfoFilled, 
  Search, 
  Document, 
  Message, 
  Check, 
  Refresh,
  User,
  UserFilled,
  Link,
  View
} from '@element-plus/icons-vue'
import { getSettings, updateSettings } from '@/api/settings'
import { useSiteConfig, initSiteConfig } from '@/store/site'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

const activeTab = ref('basic')
const formRef = ref(null)
const aboutFormRef = ref(null)
const saving = ref(false)
const { siteName } = useSiteConfig()
const vditorRef = ref(null)
let vditorInstance = null

// 默认的关于我内容（Markdown 格式）
const defaultAboutContent = `我是 Admin，一名热爱技术的全栈开发者。

## 技术栈

- **后端**: Go, GoFrame, MySQL, Redis
- **前端**: Vue 3, Element Plus, Vite
- **工具**: Docker, Git, Linux

## 关于这个博客

这个博客系统是我用来记录技术成长、分享编程心得的小天地。在这里我会分享：

- Go 语言的高性能并发模型
- Vue 3 的前端开发实践
- 云原生技术的应用经验
- 系统架构设计思路

欢迎交流学习！`

const form = reactive({
  title: '',
  subtitle: '',
  siteName: '',
  siteDesc: '',
  keywords: '',
  copyright: '',
  icpNumber: '',
  aboutAvatar: '',
  aboutSubtitle: '',
  aboutContent: defaultAboutContent,
  aboutEmail: '',
  aboutGithub: ''
})

// 保存原始数据用于重置
const originalForm = reactive({
  title: '',
  subtitle: '',
  siteName: '',
  siteDesc: '',
  keywords: '',
  copyright: '',
  icpNumber: '',
  aboutAvatar: '',
  aboutSubtitle: '',
  aboutContent: defaultAboutContent,
  aboutEmail: '',
  aboutGithub: ''
})

// 基本设置表单验证规则
const rules = {
  siteName: [
    { required: true, message: '请输入网站名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  title: [
    { required: true, message: '请输入网站标题', trigger: 'blur' },
    { max: 100, message: '长度不能超过 100 个字符', trigger: 'blur' }
  ],
  subtitle: [
    { max: 200, message: '长度不能超过 200 个字符', trigger: 'blur' }
  ],
  siteDesc: [
    { max: 500, message: '长度不能超过 500 个字符', trigger: 'blur' }
  ],
  keywords: [
    { max: 200, message: '长度不能超过 200 个字符', trigger: 'blur' }
  ],
  copyright: [
    { max: 200, message: '长度不能超过 200 个字符', trigger: 'blur' }
  ],
  icpNumber: [
    { max: 50, message: '长度不能超过 50 个字符', trigger: 'blur' }
  ]
}

// 关于我表单验证规则
const aboutRules = {
  aboutAvatar: [
    { max: 500, message: '长度不能超过 500 个字符', trigger: 'blur' }
  ],
  aboutSubtitle: [
    { max: 100, message: '长度不能超过 100 个字符', trigger: 'blur' }
  ],
  aboutContent: [
    { max: 1000, message: '长度不能超过 1000 个字符', trigger: 'change' }
  ],
  aboutEmail: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' },
    { max: 100, message: '长度不能超过 100 个字符', trigger: 'blur' }
  ],
  aboutGithub: [
    { max: 200, message: '长度不能超过 200 个字符', trigger: 'blur' }
  ]
}

const fetchSettings = async () => {
  try {
    const res = await getSettings()
    console.log('获取系统设置:', res) // 调试日志
    Object.assign(form, res)
    console.log('表单数据:', form) // 调试日志
    // 如果 aboutContent 为空，使用默认值
    if (!form.aboutContent || form.aboutContent.trim() === '') {
      form.aboutContent = defaultAboutContent
    }
    // 保存原始数据
    Object.assign(originalForm, form)
    // 如果编辑器已初始化，更新编辑器内容
    if (vditorInstance) {
      vditorInstance.setValue(form.aboutContent)
    }
  } catch (error) {
    ElMessage.error('获取系统设置失败')
    console.error(error)
  }
}

const handleSave = async () => {
  // 根据当前tab选择对应的表单ref
  const currentFormRef = activeTab.value === 'about' ? aboutFormRef.value : formRef.value
  if (!currentFormRef) return
  
  try {
    await currentFormRef.validate()
    saving.value = true
    
    await updateSettings(form)
    
    // 更新原始数据
    Object.assign(originalForm, form)
    
    // 更新全局网站配置
    await initSiteConfig()
    
    ElMessage.success('系统设置已保存')
  } catch (error) {
    if (error !== false) { // 验证失败时 error 为 false
      ElMessage.error('保存失败：' + (error.message || '未知错误'))
      console.error(error)
    }
  } finally {
    saving.value = false
  }
}

const handleReset = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要重置所有修改吗？未保存的更改将丢失。',
      '确认重置',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 重置表单
    Object.assign(form, originalForm)
    formRef.value?.clearValidate()
    aboutFormRef.value?.clearValidate()
    ElMessage.success('已重置为原始值')
  } catch {
    // 用户取消
  }
}

const handlePreviewAbout = () => {
  const baseUrl = window.location.origin
  window.open(`${baseUrl}/about`, '_blank')
}

// 初始化 Vditor 编辑器
const initVditor = () => {
  if (!vditorRef.value) return
  
  vditorInstance = new Vditor('vditor-about-content', {
    height: 500,
    placeholder: '请输入个人简介内容，支持 Markdown 语法',
    mode: 'wysiwyg', // 所见即所得模式，可选：'wysiwyg' | 'sv' | 'ir' (即时渲染，类似 Typora)
    toolbar: [
      'headings',
      'bold',
      'italic',
      'strike',
      '|',
      'line',
      'quote',
      'list',
      'ordered-list',
      'check',
      '|',
      'code',
      'inline-code',
      'link',
      'table',
      '|',
      'undo',
      'redo',
      '|',
      'outline',
      'preview',
      'fullscreen'
    ],
    value: form.aboutContent || defaultAboutContent,
    upload: {
      accept: 'image/*',
      handler: (files) => {
        // 图片上传处理（如果需要）
        console.log('上传图片:', files)
        return ''
      }
    },
    after: () => {
      // 编辑器初始化完成后的回调
      const content = form.aboutContent || defaultAboutContent
      vditorInstance.setValue(content)
      form.aboutContent = content
    },
    input: (value) => {
      // 内容变化时更新表单数据
      form.aboutContent = value
    },
    focus: (value) => {
      // 聚焦时的回调
    },
    blur: (value) => {
      // 失焦时的回调
      form.aboutContent = value
    }
  })
}

onMounted(() => {
  fetchSettings().then(() => {
    // 等待设置加载完成后再初始化编辑器
    setTimeout(() => {
      initVditor()
    }, 100)
  })
})

// 组件卸载时销毁编辑器
onBeforeUnmount(() => {
  if (vditorInstance) {
    vditorInstance.destroy()
    vditorInstance = null
  }
})
</script>

<style scoped>
.page-container {
  padding: 0;
  min-height: calc(100vh - 60px);
}

.settings-card {
  border-radius: 0;
  border: none;
  box-shadow: none;
  height: 100%;
  min-height: calc(100vh - 60px);
}

.settings-card :deep(.el-card__body) {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 18px;
  color: #303133;
}

.card-header .el-icon {
  font-size: 20px;
  color: #409eff;
}

.settings-tabs {
  margin-top: 10px;
}

.settings-tabs :deep(.el-tabs__header) {
  margin-bottom: 24px;
}

.settings-tabs :deep(.el-tabs__item) {
  font-size: 15px;
  font-weight: 500;
}

.settings-form {
  max-width: 800px;
}

.form-section {
  margin-bottom: 24px;
  border-radius: 8px;
  border: 1px solid #ebeef5;
}

.form-section :deep(.el-card__header) {
  padding: 16px 20px;
  background-color: #fafafa;
  border-bottom: 1px solid #ebeef5;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 16px;
  color: #303133;
}

.section-header .el-icon {
  font-size: 18px;
  color: #409eff;
}

.form-section :deep(.el-card__body) {
  padding: 24px;
}

.form-section :deep(.el-form-item) {
  margin-bottom: 24px;
}

.form-section :deep(.el-form-item:last-child) {
  margin-bottom: 0;
}

.form-section :deep(.el-form-item__content) {
  line-height: normal;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.5;
}

.form-actions {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #ebeef5;
  display: flex;
  gap: 12px;
}

.form-actions .el-button {
  padding: 10px 24px;
}

.form-actions .el-button .el-icon {
  margin-right: 4px;
}

.avatar-row {
  align-items: flex-start;
  width: 100%;
}

.avatar-row :deep(.el-col) {
  display: flex;
  flex-direction: column;
}

.avatar-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
  min-height: 120px;
  width: 100%;
  box-sizing: border-box;
}

.preview-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 12px;
  font-weight: 500;
  text-align: center;
}

/* 响应式 */
@media (max-width: 768px) {
  .page-container {
    padding: 10px;
  }
  
  .settings-form {
    max-width: 100%;
  }
  
  .form-section :deep(.el-card__body) {
    padding: 16px;
  }

  .avatar-row {
    margin: 0;
  }

  .avatar-row .el-col {
    margin-bottom: 12px;
  }

  .avatar-preview {
    margin-top: 0;
    min-height: 100px;
    width: 100%;
  }

  .markdown-editor-wrapper {
    width: 100%;
  }
}

.markdown-editor-wrapper {
  width: 100%;
}

.markdown-editor-wrapper :deep(#vditor-about-content) {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.markdown-editor-wrapper :deep(#vditor-about-content:hover) {
  border-color: #c0c4cc;
}

.markdown-editor-wrapper :deep(#vditor-about-content:focus-within) {
  border-color: #409eff;
}
</style>
