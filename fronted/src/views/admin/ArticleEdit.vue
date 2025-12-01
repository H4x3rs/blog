<template>
  <div class="article-edit-container">
    <el-card shadow="never" class="edit-card">
      <template #header>
        <div class="card-header">
          <el-icon><Document /></el-icon>
          <span>{{ isEdit ? '编辑文章' : '新建文章' }}</span>
        </div>
      </template>

      <el-form 
        ref="formRef" 
        :model="form" 
        :rules="rules" 
        label-width="100px"
        class="article-form"
      >
        <el-row :gutter="20">
          <el-col :span="16">
            <el-form-item label="文章标题" prop="title">
              <el-input 
                v-model="form.title" 
                placeholder="请输入文章标题"
                maxlength="200"
                show-word-limit
                clearable
                size="large"
              >
                <template #prefix>
                  <el-icon><EditPen /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="状态" prop="status">
              <el-select 
                v-model="form.status" 
                placeholder="请选择状态" 
                style="width: 100%"
                size="large"
              >
                <el-option label="草稿" value="draft">
                  <span style="display: flex; align-items: center; gap: 8px">
                    <el-icon><DocumentCopy /></el-icon>
                    <span>草稿</span>
                  </span>
                </el-option>
                <el-option label="已发布" value="published">
                  <span style="display: flex; align-items: center; gap: 8px">
                    <el-icon><CircleCheck /></el-icon>
                    <span>已发布</span>
                  </span>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="分类" prop="categoryId">
              <el-select 
                v-model="form.categoryId" 
                placeholder="请选择分类" 
                style="width: 100%"
                size="large"
              >
                <el-option 
                  v-for="cat in categories" 
                  :key="cat.id" 
                  :label="cat.name" 
                  :value="cat.id"
                >
                  <span style="display: flex; align-items: center; gap: 8px">
                    <el-icon><Folder /></el-icon>
                    <span>{{ cat.name }}</span>
                  </span>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="封面图">
              <el-input 
                v-model="form.coverImage" 
                placeholder="请输入封面图 URL"
                clearable
                size="large"
              >
                <template #prefix>
                  <el-icon><Picture /></el-icon>
                </template>
                <template #append>
                  <el-button @click="previewCover" :icon="View">预览</el-button>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="标签">
          <el-select
            v-model="form.tagIds"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="请选择或输入标签"
            style="width: 100%"
            size="large"
          >
            <el-option
              v-for="tag in tags"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            >
              <span style="display: flex; align-items: center; gap: 8px">
                <el-tag :color="tag.color" effect="plain" size="small">{{ tag.name }}</el-tag>
              </span>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="文章摘要" prop="desc">
          <el-input 
            v-model="form.desc" 
            type="textarea"
            :rows="3"
            placeholder="请输入文章摘要"
            maxlength="500"
            show-word-limit
            clearable
          />
        </el-form-item>

        <el-form-item label="文章内容" prop="content">
          <div ref="editorRef" id="cherry-editor-container" class="cherry-editor-wrapper"></div>
        </el-form-item>

        <div class="form-actions">
          <el-button @click="handleCancel" size="large">取消</el-button>
          <el-button type="info" @click="handleSaveDraft" size="large" :icon="DocumentCopy">
            保存草稿
          </el-button>
          <el-button type="primary" :loading="saving" @click="handleSubmit" size="large" :icon="CircleCheck">
            {{ isEdit ? '更新文章' : '发布文章' }}
          </el-button>
        </div>
      </el-form>
    </el-card>

    <!-- 封面图预览对话框 -->
    <el-dialog v-model="coverPreviewVisible" title="封面图预览" width="600px">
      <div class="cover-preview">
        <img v-if="form.coverImage" :src="form.coverImage" alt="封面图预览" />
        <el-empty v-else description="暂无封面图" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Document, EditPen, Folder, Picture, 
  View, DocumentCopy, CircleCheck 
} from '@element-plus/icons-vue'
import Cherry from 'cherry-markdown'
import 'cherry-markdown/dist/cherry-markdown.min.css'
import * as echarts from 'echarts'
import { getArticle, createArticle, updateArticle } from '@/api/article'
import { getCategoryList } from '@/api/category'
import { getTagList } from '@/api/tag'
import { uploadImage } from '@/api/upload'

const route = useRoute()
const router = useRouter()

const formRef = ref(null)
const editorRef = ref(null)
const saving = ref(false)
const coverPreviewVisible = ref(false)
let cherryInstance = null

const articleId = ref(null)
const isEdit = ref(false)

// 分类列表
const categories = ref([])
// 标签列表
const tags = ref([])

// 表单数据
const form = reactive({
  title: '',
  categoryId: null,
  desc: '',
  content: '',
  coverImage: '',
  status: 'draft',
  tagIds: []
})

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入文章标题', trigger: 'blur' },
    { max: 200, message: '标题长度不能超过 200 个字符', trigger: 'blur' }
  ],
  categoryId: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入文章内容', trigger: 'blur' }
  ],
  desc: [
    { max: 500, message: '摘要长度不能超过 500 个字符', trigger: 'blur' }
  ]
}

// 初始化 Cherry Markdown 编辑器
const initEditor = () => {
  if (!editorRef.value) {
    console.error('编辑器容器未找到')
    return
  }

  // 如果编辑器已存在，先销毁
  if (cherryInstance) {
    try {
      cherryInstance.destroy()
    } catch (e) {
      console.warn('销毁旧编辑器实例失败:', e)
    }
    cherryInstance = null
  }

  try {
    // 确保编辑器容器有 ID
    const editorId = editorRef.value.id || 'cherry-editor-container'
    if (!editorRef.value.id) {
      editorRef.value.id = editorId
    }

    // Cherry Markdown 配置
    const cherryConfig = {
      id: editorId,
      value: form.content || '',
      editor: {
        height: '600px',
        defaultModel: 'edit&preview', // 编辑和预览模式：'editOnly' | 'previewOnly' | 'edit&preview'
      },
      toolbars: {
        toolbar: [
          'bold',
          'italic',
          'strikethrough',
          '|',
          'header',
          '|',
          'list',
          '|',
          'quote',
          'hr',
          'code',
          'inlineCode',
          '|',
          'link',
          'image',
          'table'
        ]
      },
      // 配置外部依赖，提供 echarts 给表格转图表插件
      externals: {
        echarts: echarts
      },
      // 配置文件上传（Cherry Markdown 的图片上传配置）
      // fileUpload 应该是一个函数，接收 file 和 callback 参数
      fileUpload: async (file, callback) => {
        try {
          // 只处理图片文件
          if (!file.type.startsWith('image/')) {
            ElMessage.warning('请上传图片文件')
            return
          }
          
          ElMessage.info('正在上传图片...')
          console.log('上传文件:', file, '文件类型:', file.type, '文件大小:', file.size)
          
          // 确保 file 是 File 对象
          if (!(file instanceof File)) {
            ElMessage.error('文件格式错误')
            return
          }
          
          // 上传图片到服务器
          const res = await uploadImage(file)
          console.log('上传响应:', res)
          
          if (res && res.url) {
                // 构建完整的图片URL
                // 后端返回的 URL 如果是相对路径（如 /uploads/images/xxx.jpg），直接使用
                // 如果是绝对路径，直接使用
                let imageUrl = res.url
                // 如果后端返回的是相对路径，直接使用（浏览器会自动使用当前域名）
                // 不需要手动添加域名，因为前后端在同一域名下
                if (!imageUrl.startsWith('http://') && !imageUrl.startsWith('https://')) {
                  // 相对路径，直接使用（浏览器会自动使用当前域名）
                  // 例如：/uploads/images/xxx.jpg 会被解析为 https://blog.renhj.cc/uploads/images/xxx.jpg
                }
            
            // 调用 callback 函数，将图片URL传递给编辑器
            // callback 函数会将图片插入到编辑器中
            callback(imageUrl)
            ElMessage.success('图片上传成功')
          } else {
            console.error('上传响应格式错误:', res)
            ElMessage.error('图片上传失败：响应格式错误')
          }
        } catch (error) {
          console.error('图片上传失败:', error)
          console.error('错误详情:', error.response?.data || error.message)
          ElMessage.error('图片上传失败：' + (error.message || '未知错误'))
        }
      },
      callback: {
        afterChange: (content) => {
          form.content = content
        }
      }
    }

    cherryInstance = new Cherry(cherryConfig)
  } catch (error) {
    console.error('Cherry Markdown 初始化失败:', error)
    console.error('错误详情:', error.message, error.stack)
    ElMessage.error('编辑器初始化失败：' + (error.message || '未知错误'))
    
    // 如果初始化失败，尝试使用最简单的配置
    try {
      console.log('尝试使用默认配置重新初始化...')
      cherryInstance = new Cherry({
        id: editorRef.value.id || 'cherry-editor-container',
        value: form.content || '',
        editor: {
          height: '600px',
          defaultModel: 'edit&preview'
        }
      })
      ElMessage.success('编辑器已使用默认配置初始化')
    } catch (retryError) {
      console.error('使用默认配置也失败:', retryError)
      ElMessage.error('编辑器初始化失败，请刷新页面重试')
    }
  }
}

// 加载分类列表
const loadCategories = async () => {
  try {
    const res = await getCategoryList({ page: 1, size: 100 })
    if (res && res.list) {
      categories.value = res.list
    }
  } catch (error) {
    console.error('加载分类列表失败:', error)
  }
}

// 加载标签列表
const loadTags = async () => {
  try {
    const res = await getTagList({ page: 1, size: 100 })
    if (res && res.list) {
      tags.value = res.list
    }
  } catch (error) {
    console.error('加载标签列表失败:', error)
  }
}

// 加载文章数据
const loadArticle = async () => {
  if (!articleId.value) return
  
  try {
    const res = await getArticle({ id: articleId.value })
    if (res) {
      form.title = res.title || ''
      form.categoryId = res.categoryId || null
      form.desc = res.desc || ''
      form.content = res.content || ''
      form.coverImage = res.coverImage || ''
      form.status = res.status || 'draft'
      
      // 加载文章标签
      if (res.tags && Array.isArray(res.tags)) {
        form.tagIds = res.tags.map(tag => tag.id)
      } else {
        form.tagIds = []
      }
      
      // 更新编辑器内容，需要等待编辑器初始化完成
      if (cherryInstance) {
        try {
          cherryInstance.setMarkdown(form.content)
        } catch (e) {
          console.warn('设置编辑器内容失败:', e)
          // 如果设置失败，尝试重新初始化编辑器
          setTimeout(() => {
            if (cherryInstance) {
              cherryInstance.setMarkdown(form.content)
            }
          }, 200)
        }
      }
    }
  } catch (error) {
    ElMessage.error('加载文章失败')
    console.error(error)
  }
}

// 保存草稿
const handleSaveDraft = async () => {
  form.status = 'draft'
  await handleSubmit()
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    
    // 获取编辑器内容
    if (cherryInstance) {
      form.content = cherryInstance.getMarkdown()
    }
    
    if (!form.content || form.content.trim() === '') {
      ElMessage.warning('请输入文章内容')
      return
    }

    saving.value = true

    if (isEdit.value) {
      await updateArticle({
        id: articleId.value,
        ...form
      })
      ElMessage.success('文章更新成功')
    } else {
      await createArticle(form)
      ElMessage.success('文章发布成功')
    }

    // 返回文章列表
    router.push('/admin/articles')
  } catch (error) {
    if (error !== false) {
      ElMessage.error('操作失败：' + (error.message || '未知错误'))
      console.error(error)
    }
  } finally {
    saving.value = false
  }
}

// 取消编辑
const handleCancel = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要离开吗？未保存的更改将丢失。',
      '确认离开',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    router.back()
  } catch {
    // 用户取消
  }
}

// 预览封面图
const previewCover = () => {
  if (!form.coverImage) {
    ElMessage.warning('请先输入封面图 URL')
    return
  }
  coverPreviewVisible.value = true
}

// 初始化
onMounted(async () => {
  // 检查是否是编辑模式
  const id = route.params.id
  if (id) {
    articleId.value = parseInt(id)
    isEdit.value = true
  }

  // 加载分类列表和标签列表
  await Promise.all([loadCategories(), loadTags()])

  // 如果是编辑模式，先加载文章数据
  if (isEdit.value) {
    await loadArticle()
  }

  // 初始化编辑器（延迟确保 DOM 已渲染）
  setTimeout(() => {
    initEditor()
    // 如果是编辑模式，编辑器初始化后再次设置内容
    if (isEdit.value && form.content && cherryInstance) {
      setTimeout(() => {
        if (cherryInstance) {
          cherryInstance.setMarkdown(form.content)
        }
      }, 100)
    }
  }, 100)
})

// 组件卸载时销毁编辑器
onBeforeUnmount(() => {
  if (cherryInstance) {
    cherryInstance.destroy()
    cherryInstance = null
  }
})
</script>

<style scoped>
.article-edit-container {
  padding: 0;
  min-height: calc(100vh - 60px);
}

.edit-card {
  border-radius: 16px;
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
  height: 100%;
  min-height: calc(100vh - 60px);
}

.edit-card :deep(.el-card__header) {
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.edit-card :deep(.el-card__body) {
  padding: 24px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
  font-size: 18px;
  color: #303133;
}

.card-header .el-icon {
  font-size: 22px;
  color: #409eff;
}

.article-form {
  max-width: 100%;
}

/* 表单样式统一 */
:deep(.el-form-item__label) {
  font-weight: 500;
  color: #606266;
  font-size: 14px;
}

:deep(.el-form-item) {
  margin-bottom: 22px;
}

/* 输入框统一样式 */
:deep(.el-input__wrapper) {
  border-radius: 6px;
  padding: 0 12px;
  height: 40px;
  font-size: 14px;
  transition: all 0.3s ease;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

:deep(.el-input--large .el-input__wrapper) {
  height: 40px;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

:deep(.el-input.is-focus .el-input__wrapper) {
  box-shadow: 0 0 0 1px #409eff inset, 0 0 0 3px rgba(64, 158, 255, 0.1);
}

:deep(.el-input__inner) {
  font-size: 14px;
  color: #606266;
  line-height: 1.5;
}

:deep(.el-input__prefix) {
  color: #909399;
  margin-right: 8px;
}

:deep(.el-input.is-focus .el-input__prefix) {
  color: #409eff;
}

/* 选择器样式 */
:deep(.el-select__wrapper) {
  border-radius: 6px;
  padding: 0 12px;
  height: 40px;
  font-size: 14px;
  transition: all 0.3s ease;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

:deep(.el-select__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

:deep(.el-select.is-focused .el-select__wrapper) {
  box-shadow: 0 0 0 1px #409eff inset, 0 0 0 3px rgba(64, 158, 255, 0.1);
}

/* 文本域样式 */
:deep(.el-textarea__inner) {
  border-radius: 6px;
  font-size: 14px;
  padding: 12px;
  line-height: 1.5;
  transition: all 0.3s ease;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

:deep(.el-textarea__inner:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

:deep(.el-textarea.is-focus .el-textarea__inner) {
  box-shadow: 0 0 0 1px #409eff inset, 0 0 0 3px rgba(64, 158, 255, 0.1);
}

/* 编辑器样式 */
.cherry-editor-wrapper {
  width: 100%;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.cherry-editor-wrapper:hover {
  border-color: #c0c4cc;
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

.cherry-editor-wrapper:focus-within {
  border-color: #409eff;
  box-shadow: 0 0 0 1px #409eff inset, 0 0 0 3px rgba(64, 158, 255, 0.1);
}

/* 操作按钮区域 */
.form-actions {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #ebeef5;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.form-actions .el-button {
  min-width: 100px;
  font-weight: 500;
}

/* 封面预览对话框 */
.cover-preview {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
  background: #f5f7fa;
  border-radius: 8px;
}

.cover-preview img {
  max-width: 100%;
  max-height: 500px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* 响应式 */
@media (max-width: 768px) {
  .edit-card :deep(.el-card__body) {
    padding: 16px;
  }

  .form-actions {
    flex-direction: column;
  }

  .form-actions .el-button {
    width: 100%;
  }
}
</style>

