<template>
  <div class="topic-manage">
    <!-- 筛选区域 -->
    <el-card shadow="never" class="filter-card">
          <el-form :inline="true" :model="topicFilters" class="filter-form">
            <el-form-item label="专题名称">
              <el-input 
                v-model="topicFilters.name" 
                placeholder="请输入专题名称" 
                clearable 
                :prefix-icon="Search"
                style="width: 260px"
                size="default"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :icon="Search" @click="handleTopicSearch" size="default">查询</el-button>
              <el-button :icon="Refresh" @click="resetTopicFilters" size="default">重置</el-button>
            </el-form-item>
            <el-form-item style="margin-left: auto;">
              <el-button type="primary" :icon="Plus" @click="openTopicDialog('create')" size="default">
                <span style="margin-left: 4px;">创建专题</span>
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 专题列表 (Grid Card Layout) -->
        <div class="topic-list" v-loading="topicLoading">
          <el-empty v-if="topicData.length === 0" description="暂无专题数据" />
          <div class="topic-grid" v-else>
            <div v-for="item in topicData" :key="item.id" class="topic-card-wrapper">
              <el-card shadow="hover" class="topic-card" :body-style="{ padding: '0px' }">
                <div class="topic-cover">
                  <el-image :src="item.coverImage" fit="cover" lazy>
                    <template #error>
                      <div class="image-slot">
                        <el-icon><Picture /></el-icon>
                      </div>
                    </template>
                  </el-image>
                  <div class="topic-actions-overlay">
                    <el-button circle type="primary" :icon="Edit" @click.stop="openTopicDialog('edit', item)"></el-button>
                    <el-button circle type="success" :icon="List" @click.stop="manageTopicArticles(item)"></el-button>
                    <el-button circle type="danger" :icon="Delete" @click.stop="handleTopicDelete(item)"></el-button>
                  </div>
                </div>
                <div class="topic-content">
                  <div class="topic-header">
                    <h3 class="topic-name" :title="item.name">{{ item.name }}</h3>
                  </div>
                  <p class="topic-desc">{{ item.description || '暂无描述' }}</p>
                  <div class="topic-footer">
                    <span class="date">
                      <el-icon class="date-icon"><Calendar /></el-icon>
                      <span class="date-text">{{ formatDate(item.createdAt) }}</span>
                    </span>
                    <div class="article-count-badge">
                      <el-icon class="count-icon"><Document /></el-icon>
                      <span class="count-text">{{ item.articleCount || 0 }}篇文章</span>
                    </div>
                  </div>
                </div>
              </el-card>
            </div>
          </div>

          <div class="pagination-container" v-if="topicTotal > topicPageSize">
            <el-pagination
              background
              layout="total, prev, pager, next, jumper"
              :total="topicTotal"
              v-model:current-page="topicCurrentPage"
              v-model:page-size="topicPageSize"
            />
          </div>
        </div>

    <!-- 创建/编辑专题对话框 -->
    <el-dialog
      v-model="topicDialogVisible"
      :title="topicDialogType === 'create' ? '创建专题' : '编辑专题'"
      width="800px"
      align-center
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      destroy-on-close
    >
      <el-form 
        :model="topicForm" 
        :rules="topicRules" 
        ref="topicFormRef"
        label-width="100px"
        label-position="right"
      >
        <el-form-item label="专题名称" prop="name">
          <el-input 
            v-model="topicForm.name" 
            placeholder="请输入专题名称"
            clearable
            size="large"
          >
            <template #prefix>
              <el-icon><EditPen /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item label="专题别名" prop="slug">
          <el-input 
            v-model="topicForm.slug" 
            placeholder="专题别名（URL友好格式，英文）"
            clearable
            size="large"
          >
            <template #prefix>
              <el-icon><Link /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item label="专题描述" prop="description">
          <el-input 
            v-model="topicForm.description" 
            type="textarea" 
            :rows="3" 
            placeholder="请输入专题描述"
            maxlength="500"
            show-word-limit
            clearable
          />
        </el-form-item>
        
        <el-form-item label="专题介绍" prop="intro">
          <el-input 
            v-model="topicForm.intro" 
            type="textarea" 
            :rows="4" 
            placeholder="请输入专题详细介绍"
            clearable
          />
        </el-form-item>
        
        <el-form-item label="封面图片" prop="coverImage">
          <div class="cover-image-container">
            <div class="cover-input-group">
              <el-input 
                v-model="topicForm.coverImage" 
                placeholder="请输入图片URL，支持 http/https 链接"
                clearable
                size="large"
                class="cover-input"
              >
                <template #prefix>
                  <el-icon><Picture /></el-icon>
                </template>
              </el-input>
              <el-button 
                :icon="View" 
                @click="previewCoverImage"
                :disabled="!topicForm.coverImage"
                size="large"
                type="primary"
                class="preview-btn"
              >
                预览
              </el-button>
            </div>
          </div>
        </el-form-item>
        
        <el-form-item label="是否置顶" prop="isFeatured">
          <el-switch
            v-model="topicForm.isFeatured"
            :active-value="1"
            :inactive-value="0"
            active-text="是"
            inactive-text="否"
            inline-prompt
            class="featured-switch"
            style="--el-switch-on-color: #67c23a;"
          />
        </el-form-item>
        
        <el-form-item label="排序" prop="sortOrder">
          <el-input-number 
            v-model="topicForm.sortOrder" 
            :min="0" 
            :max="9999"
            size="large"
            style="width: 200px;"
          />
          <span class="form-tip">数字越小越靠前</span>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <div class="footer-left">
            <el-button 
              type="info" 
              :icon="View" 
              @click="previewTopic" 
              size="large"
              :disabled="!topicForm.name"
            >
              预览专题
            </el-button>
          </div>
          <div class="footer-right">
            <el-button @click="topicDialogVisible = false" size="large">取消</el-button>
            <el-button 
              type="primary" 
              @click="handleTopicSubmit" 
              size="large"
              :loading="submitting"
            >
              {{ topicDialogType === 'create' ? '创建' : '保存' }}
            </el-button>
          </div>
        </div>
      </template>
    </el-dialog>
    
    <!-- 封面图预览对话框 -->
    <el-dialog v-model="coverPreviewVisible" title="封面图预览" width="600px">
      <div class="cover-preview-dialog">
        <el-image 
          v-if="topicForm.coverImage" 
          :src="topicForm.coverImage" 
          alt="封面图预览"
          fit="contain"
          style="max-width: 100%; max-height: 500px;"
        />
        <el-empty v-else description="暂无封面图" />
      </div>
    </el-dialog>

    <!-- 专题完整预览对话框 -->
    <el-dialog
      v-model="topicPreviewVisible"
      title="专题预览"
      width="90%"
      top="5vh"
      destroy-on-close
      class="topic-preview-dialog"
      :close-on-click-modal="false"
    >
      <div class="topic-preview-container" v-loading="previewLoading">
        <!-- 预览头部 -->
        <div class="preview-header" :style="{ backgroundImage: `url(${topicForm.coverImage || 'https://picsum.photos/1200/400'})` }">
          <div class="header-overlay"></div>
          <div class="header-content">
            <div class="preview-meta-info">
              <el-tag v-if="topicForm.isFeatured === 1" effect="dark" round class="featured-tag">
                置顶专题
              </el-tag>
            </div>
            <h1 class="preview-title">{{ topicForm.name || '专题名称' }}</h1>
            <p class="preview-description">{{ topicForm.description || '专题描述' }}</p>
            <div class="preview-meta">
              <span class="meta-item">
                <el-icon><Document /></el-icon>
                {{ previewArticleCount }} 篇文章
              </span>
              <span class="divider">·</span>
              <span class="meta-item">
                <el-icon><View /></el-icon>
                {{ formatNumber(previewTotalViews) }} 阅读
              </span>
            </div>
          </div>
        </div>

        <!-- 预览内容区域 -->
        <div class="preview-content-wrapper">
          <!-- 专题介绍 -->
          <el-card class="preview-intro-card" shadow="never" v-if="topicForm.intro">
            <h3 class="section-title">
              <el-icon><Reading /></el-icon>
              专题介绍
            </h3>
            <div class="markdown-body" v-html="renderedPreviewIntro"></div>
          </el-card>

          <!-- 文章列表 -->
          <el-card class="preview-articles-card" shadow="never">
            <div class="articles-header">
              <h3 class="section-title">
                <el-icon><List /></el-icon>
                文章列表
                <span class="count">({{ previewArticles.length }})</span>
              </h3>
            </div>

            <div class="articles-list" v-if="previewArticles.length > 0">
              <div 
                v-for="(article, index) in previewArticles" 
                :key="article.id"
                class="article-item"
              >
                <div class="article-number">{{ index + 1 }}</div>
                <div class="article-content">
                  <div class="article-top">
                    <h4 class="article-title">{{ article.title }}</h4>
                  </div>
                  <p class="article-summary">{{ article.desc || '暂无摘要' }}</p>
                  <div class="article-meta">
                    <span class="meta-item">
                      <el-icon><Calendar /></el-icon>
                      {{ formatPreviewDate(article.createdAt) }}
                    </span>
                    <span class="meta-item">
                      <el-icon><View /></el-icon>
                      {{ formatNumber(article.views || 0) }}
                    </span>
                  </div>
                </div>
                <div class="article-cover" v-if="article.coverImage">
                  <img :src="article.coverImage" :alt="article.title" />
                </div>
              </div>
            </div>
            <el-empty v-else description="暂无文章，预览时将显示空状态" />
          </el-card>
        </div>
      </div>
    </el-dialog>

    <!-- 专题详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="currentTopic?.name"
      width="85%"
      destroy-on-close
      class="topic-detail-dialog"
      top="5vh"
    >
      <template #header>
        <div class="detail-dialog-header">
          <div class="header-left">
            <el-icon :size="24" color="#409eff"><Reading /></el-icon>
            <div class="header-info">
              <h3 class="topic-title">{{ currentTopic?.name }}</h3>
              <p class="topic-description">{{ currentTopic?.description || '暂无描述' }}</p>
            </div>
          </div>
          <div class="header-right">
            <el-button type="primary" :icon="Plus" @click="openAddArticleDialog">添加文章</el-button>
          </div>
        </div>
      </template>

      <div class="detail-content">
        <div class="detail-stats">
          <div class="stat-item">
            <el-icon><Document /></el-icon>
            <div class="stat-info">
              <span class="stat-value">{{ topicArticles.length }}</span>
              <span class="stat-label">文章数量</span>
            </div>
          </div>
          <div class="stat-item">
            <el-icon><View /></el-icon>
            <div class="stat-info">
              <span class="stat-value">{{ calculateTotalViews() }}</span>
              <span class="stat-label">总阅读量</span>
            </div>
          </div>
          <div class="stat-item">
            <el-icon><Calendar /></el-icon>
            <div class="stat-info">
              <span class="stat-value">{{ formatDate(currentTopic?.createdAt) }}</span>
              <span class="stat-label">创建时间</span>
            </div>
          </div>
        </div>

        <el-divider />

        <div class="articles-section">
          <div class="section-header">
            <h4 class="section-title">
              <el-icon><List /></el-icon>
              <span>专题文章列表</span>
            </h4>
          </div>

          <el-table 
            :data="topicArticles" 
            style="width: 100%" 
            v-loading="articleLoading"
            row-key="id"
            :header-cell-style="{ background: '#fafafa', color: '#303133', fontWeight: '600', fontSize: '13px' }"
            :cell-style="{ fontSize: '13px' }"
            stripe
            empty-text="暂无文章，点击上方按钮添加"
            class="topic-articles-table"
          >
            <el-table-column label="排序" width="80" align="center">
              <template #default="scope">
                <span class="sort-index">{{ scope.$index + 1 + (articlePage - 1) * articlePageSize }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="articleId" label="文章ID" width="85" align="center">
              <template #default="scope">
                <el-tag size="small" effect="plain" type="info" style="font-weight: 600;">{{ scope.row.articleId || '-' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="title" label="文章标题" min-width="250" show-overflow-tooltip>
              <template #default="scope">
                <div class="article-title-cell">
                  <el-icon color="#409eff"><Document /></el-icon>
                  <span class="title-text">{{ scope.row.title }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="category" label="分类" width="120" align="center">
              <template #default="scope">
                <el-tag size="small" type="success" effect="plain">{{ scope.row.category }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="views" label="阅读量" width="100" align="center">
              <template #default="scope">
                <div class="views-cell">
                  <el-icon><View /></el-icon>
                  <span>{{ formatNumber(scope.row.views || 0) }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="发布时间" width="160" align="center">
              <template #default="scope">
                <div class="date-cell">
                  <el-icon class="date-icon-small"><Calendar /></el-icon>
                  <span class="date-text">{{ formatPreviewDate(scope.row.createdAt) }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="scope">
                <div style="display: flex; align-items: center; justify-content: center; gap: 8px; white-space: nowrap;">
                  <el-button 
                    link 
                    size="small" 
                    :icon="ArrowUp"
                    :disabled="scope.$index === 0"
                    @click="handleMoveUp(scope.$index)"
                  >
                    上移
                  </el-button>
                  <el-button 
                    link 
                    size="small" 
                    :icon="ArrowDown"
                    :disabled="scope.$index === topicArticles.length - 1"
                    @click="handleMoveDown(scope.$index)"
                  >
                    下移
                  </el-button>
                  <el-button 
                    link 
                    type="danger" 
                    size="small" 
                    :icon="Delete"
                    @click="handleRemoveArticle(scope.row)"
                  >
                    移除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
          
          <!-- 分页组件 -->
          <div class="pagination-container" v-if="articleTotal > articlePageSize" style="margin-top: 20px; display: flex; justify-content: flex-end;">
            <el-pagination
              v-model:current-page="articlePage"
              v-model:page-size="articlePageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="articleTotal"
              layout="total, sizes, prev, pager, next, jumper"
              @current-change="handleArticlePageChange"
              @size-change="handleArticlePageSizeChange"
            />
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 添加文章到专题对话框 -->
    <el-dialog
      v-model="addArticleDialogVisible"
      title="添加文章到专题"
      width="70%"
      destroy-on-close
      class="article-dialog"
    >
      <div class="dialog-search-area">
        <el-form :inline="true" :model="articleFilters" class="filter-form">
          <el-form-item label="文章标题">
            <el-input 
              v-model="articleFilters.title" 
              placeholder="搜索文章标题" 
              clearable 
              :prefix-icon="Search"
              style="width: 280px"
            />
          </el-form-item>
          <el-form-item label="分类">
            <el-select v-model="articleFilters.categoryId" placeholder="请选择分类" clearable style="width: 200px">
              <el-option 
                v-for="cat in categoryList" 
                :key="cat.id" 
                :label="cat.name" 
                :value="cat.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :icon="Search" @click="searchAvailableArticles">搜索</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table
        :data="availableArticles"
        style="width: 100%"
        v-loading="availableArticlesLoading"
        @selection-change="handleArticleSelectionChange"
        :header-cell-style="{ background: '#fafafa', color: '#303133', fontWeight: '600' }"
        stripe
        max-height="450"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column prop="id" label="ID" width="70" align="center">
          <template #default="scope">
            <el-tag size="small" effect="plain">#{{ scope.row.id }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="250" show-overflow-tooltip>
          <template #default="scope">
            <div class="article-title-cell">
              <el-icon color="#409eff"><Document /></el-icon>
              <span class="title-text">{{ scope.row.title }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="categoryId" label="分类" width="120" align="center">
          <template #default="scope">
            <el-tag size="small" type="success" effect="plain">
              {{ categoryMap[scope.row.categoryId] || '未分类' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="views" label="阅读量" width="110" align="center">
          <template #default="scope">
            <div class="views-cell">
              <el-icon><View /></el-icon>
              <span>{{ scope.row.views }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="发布时间" width="180" align="center">
          <template #default="scope">
            <span class="date-text">{{ scope.row.createdAt }}</span>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <div class="dialog-footer">
          <div class="footer-info">
            <span v-if="selectedArticles.length > 0" class="selection-info">
              <el-icon><SuccessFilled /></el-icon>
              已选择 {{ selectedArticles.length }} 篇文章
            </span>
          </div>
          <div class="footer-actions">
            <el-button size="large" @click="addArticleDialogVisible = false">取消</el-button>
            <el-button 
              type="primary" 
              size="large"
              @click="handleAddArticles"
              :disabled="selectedArticles.length === 0"
              :icon="Plus"
            >
              添加选中文章
            </el-button>
          </div>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { 
  Search, Refresh, Plus, Edit, Delete, Picture, View, List, Upload, 
  Reading, Document, SuccessFilled, Calendar, EditPen, Link, 
  Document as DocumentIcon, CircleCheck, CircleClose, ArrowUp, ArrowDown
} from '@element-plus/icons-vue'
import { marked } from 'marked'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import customParseFormat from 'dayjs/plugin/customParseFormat'
import 'dayjs/locale/zh-cn'
import { getTopicList, createTopic, updateTopic, deleteTopic, getTopicArticles, addArticleToTopic, removeArticleFromTopic, updateTopicArticleSort } from '@/api/topic'
import { getArticleList } from '@/api/article'
import { getCategoryList } from '@/api/category'

// 配置 dayjs
dayjs.extend(relativeTime)
dayjs.extend(customParseFormat)
dayjs.locale('zh-cn')

// ===== 专题管理相关 =====
const topicLoading = ref(false)
const topicFilters = reactive({
  name: ''
})
const topicData = ref([])
const topicTotal = ref(0)
const topicCurrentPage = ref(1)
const topicPageSize = ref(12)

const topicDialogVisible = ref(false)
const topicDialogType = ref('create')
const topicFormRef = ref(null)
const submitting = ref(false)
const coverPreviewVisible = ref(false)
const topicPreviewVisible = ref(false)
const previewLoading = ref(false)
const previewArticles = ref([])
const topicForm = reactive({
  id: null,
  name: '',
  slug: '',
  description: '',
  coverImage: '',
  intro: '',
  isFeatured: 0,
  sortOrder: 0
})

const topicRules = {
  name: [
    { required: true, message: '请输入专题名称', trigger: 'blur' },
    { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  slug: [
    { required: true, message: '请输入专题别名', trigger: 'blur' },
    { pattern: /^[a-z0-9-]+$/, message: '别名只能包含小写字母、数字和连字符', trigger: 'blur' },
    { min: 1, max: 100, message: '别名长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 500, message: '描述不能超过 500 个字符', trigger: 'blur' }
  ],
  coverImage: [
    { type: 'url', message: '请输入正确的URL地址', trigger: 'blur' }
  ],
  sortOrder: [
    { type: 'number', min: 0, max: 9999, message: '排序值必须在 0 到 9999 之间', trigger: 'blur' }
  ]
}

const handleTopicSearch = async () => {
  topicLoading.value = true
  try {
    const res = await getTopicList({
      page: topicCurrentPage.value,
      size: topicPageSize.value,
      name: topicFilters.name
    })
    // 获取每个专题的文章数量
    const topicsWithCount = await Promise.all((res.list || []).map(async (item) => {
      try {
        const articlesRes = await getTopicArticles({ id: item.id, page: 1, size: 100 })
        // 确保响应数据存在，尝试多种可能的字段名
        if (articlesRes) {
          const count = articlesRes.total ?? articlesRes.Total ?? (articlesRes.list?.length ?? 0)
          return {
            ...item,
            articleCount: count
          }
        } else {
          console.warn(`获取专题 ${item.id} 的文章数量失败: 响应数据为空`, articlesRes)
          return {
            ...item,
            articleCount: 0
          }
        }
      } catch (error) {
        console.error(`获取专题 ${item.id} 的文章数量失败:`, error)
        return {
          ...item,
          articleCount: 0
        }
      }
    }))
    topicData.value = topicsWithCount
    topicTotal.value = res.total || 0
  } catch (error) {
    ElMessage.error('获取专题列表失败')
    console.error(error)
  } finally {
    topicLoading.value = false
  }
}

const resetTopicFilters = () => {
  topicFilters.name = ''
  handleTopicSearch()
}

const openTopicDialog = (type, row) => {
  topicDialogType.value = type
  topicDialogVisible.value = true
  
  // 重置表单验证状态
  if (topicFormRef.value) {
    topicFormRef.value.clearValidate()
  }
  
  if (type === 'edit' && row) {
    Object.assign(topicForm, {
      id: row.id,
      name: row.name || '',
      slug: row.slug || '',
      description: row.description || '',
      coverImage: row.coverImage || '',
      intro: row.intro || '',
      isFeatured: row.isFeatured || 0,
      sortOrder: row.sortOrder || 0
    })
  } else {
    Object.assign(topicForm, {
      id: null,
      name: '',
      slug: '',
      description: '',
      coverImage: '',
      intro: '',
      isFeatured: 0,
      sortOrder: 0
    })
  }
}

// 预览封面图
const previewCoverImage = () => {
  if (!topicForm.coverImage) {
    ElMessage.warning('请先输入封面图 URL')
    return
  }
  coverPreviewVisible.value = true
}

// 预览专题
const previewTopic = async () => {
  if (!topicForm.name) {
    ElMessage.warning('请先输入专题名称')
    return
  }
  
  previewLoading.value = true
  topicPreviewVisible.value = true
  
  // 如果是编辑模式且有ID，加载专题文章
  if (topicForm.id) {
    try {
      const res = await getTopicArticles({ id: topicForm.id, page: 1, size: 100 })
      previewArticles.value = res.list || []
    } catch (error) {
      console.error('加载专题文章失败:', error)
      previewArticles.value = []
    }
  } else {
    previewArticles.value = []
  }
  
  previewLoading.value = false
}

// 计算预览文章数量
const previewArticleCount = computed(() => previewArticles.value.length)

// 计算预览总阅读量
const previewTotalViews = computed(() => {
  return previewArticles.value.reduce((total, article) => total + (article.views || 0), 0)
})

// 渲染预览介绍（Markdown）
const renderedPreviewIntro = computed(() => {
  if (!topicForm.intro) return ''
  try {
    return marked.parse(topicForm.intro)
  } catch (error) {
    return topicForm.intro
  }
})

const handleTopicSubmit = async () => {
  if (!topicFormRef.value) return
  
  // 验证表单
  try {
    await topicFormRef.value.validate()
  } catch (error) {
    return
  }
  
  submitting.value = true
  try {
    const data = {
      name: topicForm.name,
      slug: topicForm.slug || topicForm.name.toLowerCase().replace(/\s+/g, '-'),
      description: topicForm.description || '',
      coverImage: topicForm.coverImage || '',
      intro: topicForm.intro || topicForm.description || '',
      isFeatured: topicForm.isFeatured || 0,
      sortOrder: topicForm.sortOrder || 0
    }
    
    if (topicDialogType.value === 'create') {
      await createTopic(data)
      ElMessage.success('专题创建成功')
    } else {
      await updateTopic({
        id: topicForm.id,
        ...data
      })
      ElMessage.success('专题更新成功')
    }
    topicDialogVisible.value = false
    handleTopicSearch()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const handleTopicDelete = (row) => {
  ElMessageBox.confirm(
    '确定要删除这个专题吗？删除后无法恢复。',
    '警告',
    {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await deleteTopic({ id: row.id })
      ElMessage.success('专题删除成功')
      handleTopicSearch()
    } catch (error) {
      ElMessage.error('删除失败')
      console.error(error)
    }
  }).catch(() => {})
}

const manageTopicArticles = (topic) => {
  currentTopic.value = topic
  detailDialogVisible.value = true
  // 重置分页
  articlePage.value = 1
  articlePageSize.value = 10
  loadTopicArticles({ id: topic.id }, articlePage.value, articlePageSize.value)
}

// ===== 专题详情对话框 =====
const detailDialogVisible = ref(false)
const currentTopic = ref(null)

// ===== 专题文章管理相关 =====
const articleLoading = ref(false)
const topicArticles = ref([])
// 分页相关
const articlePage = ref(1)
const articlePageSize = ref(10)
const articleTotal = ref(0)

const addArticleDialogVisible = ref(false)
const availableArticlesLoading = ref(false)
const availableArticles = ref([
  {
    id: 101,
    title: 'GoFrame 最佳实践指南',
    category: 'Golang',
    views: 1230,
    createdAt: '2024-05-03 10:30:00'
  },
  {
    id: 102,
    title: 'Vue 3 + Element Plus 全面解析',
    category: 'Frontend',
    views: 856,
    createdAt: '2024-05-02 14:20:00'
  },
  {
    id: 103,
    title: 'Docker 容器化部署实战',
    category: 'DevOps',
    views: 645,
    createdAt: '2024-05-01 09:15:00'
  },
  {
    id: 104,
    title: 'MySQL 索引优化指南',
    category: 'Database',
    views: 2341,
    createdAt: '2024-04-28 16:40:00'
  }
])
const selectedArticles = ref([])

const articleFilters = reactive({
  title: '',
  categoryId: ''
})

// 分类列表（用于下拉选择）
const categoryList = ref([])

// 分类映射
const categoryMap = ref({})

// 加载分类映射和分类列表
const loadCategoryMap = async () => {
  try {
    const map = {}
    const list = []
    let page = 1
    const size = 100 // 后端限制最大100
    let hasMore = true
    
    // 分页加载所有分类
    while (hasMore) {
      const res = await getCategoryList({ page, size })
      const categoryData = res.list || []
      
      categoryData.forEach(cat => {
        map[cat.id] = cat.name
        list.push(cat)
      })
      
      // 如果返回的数据少于size，说明已经是最后一页
      if (categoryData.length < size) {
        hasMore = false
      } else {
        page++
      }
    }
    
    categoryMap.value = map
    categoryList.value = list
  } catch (error) {
    console.error('加载分类列表失败:', error)
  }
}

const loadTopicArticles = async (data, page = 1, size = 10) => {
  articleLoading.value = true
  try {
    // 确保分类映射已加载
    if (Object.keys(categoryMap.value).length === 0) {
      await loadCategoryMap()
    }
    
    const res = await getTopicArticles({
      ...data,
      page,
      size
    })
    // 为每篇文章添加分类名称
    topicArticles.value = (res.list || []).map(item => ({
      ...item,
      category: categoryMap.value[item.categoryId] || '未分类'
    }))
    // 更新总数
    articleTotal.value = res.total || 0
  } catch (error) {
    ElMessage.error('获取专题文章失败')
    console.error(error)
    topicArticles.value = []
    articleTotal.value = 0
  } finally {
    articleLoading.value = false
  }
}

// 计算总阅读量
const calculateTotalViews = () => {
  return topicArticles.value.reduce((total, article) => total + (article.views || 0), 0)
}

const openAddArticleDialog = async () => {
  addArticleDialogVisible.value = true
  // 确保分类列表已加载
  if (categoryList.value.length === 0) {
    await loadCategoryMap()
  }
  searchAvailableArticles()
}

const searchAvailableArticles = async () => {
  availableArticlesLoading.value = true
  try {
    // 确保分类映射已加载
    if (Object.keys(categoryMap.value).length === 0) {
      await loadCategoryMap()
    }
    
    const res = await getArticleList({
      page: 1,
      size: 50,
      title: articleFilters.title,
      categoryId: articleFilters.categoryId ? parseInt(articleFilters.categoryId) : 0
    })
    availableArticles.value = res.list || []
  } catch (error) {
    ElMessage.error('获取文章列表失败')
    console.error(error)
  } finally {
    availableArticlesLoading.value = false
  }
}

const handleArticleSelectionChange = (selection) => {
  selectedArticles.value = selection
}

const handleAddArticles = async () => {
  if (selectedArticles.value.length === 0) {
    ElMessage.warning('请选择要添加的文章')
    return
  }
  
  // 检查重复文章 - 获取所有文章（不限于当前页）
  // 使用文章ID来判断是否重复
  try {
    const allArticlesRes = await getTopicArticles({ id: currentTopic.value.id, page: 1, size: 100 })
    // 提取专题中已有文章的文章ID集合（item.articleId 对应文章ID）
    const existingArticleIds = new Set((allArticlesRes.list || []).map(item => item.articleId))
    // 检查选中的文章ID是否已在专题中（article.id 对应文章ID）
    const duplicateArticles = selectedArticles.value.filter(article => 
      existingArticleIds.has(article.id)
    )
    
    if (duplicateArticles.length > 0) {
      const duplicateTitles = duplicateArticles.map(a => a.title).join('、')
      ElMessage.warning(`以下文章已在专题中，无法重复添加：${duplicateTitles}`)
      return
    }
    
    // 获取当前最大排序值
    const allArticles = allArticlesRes.list || []
    const maxSortOrder = allArticles.length > 0 
      ? Math.max(...allArticles.map(item => item.sortOrder || 0))
      : -1
    
    let currentSortOrder = maxSortOrder + 1
    const addedCount = []
    const failedArticles = []
    
    for (const article of selectedArticles.value) {
      try {
        await addArticleToTopic({
          id: currentTopic.value.id,
          articleId: article.id,
          sortOrder: currentSortOrder
        })
        addedCount.push(article.title)
        currentSortOrder++
      } catch (error) {
        // 捕获后端返回的重复错误或其他错误
        const errorMessage = error.message || error.response?.data?.message || '添加失败'
        failedArticles.push({ title: article.title, error: errorMessage })
        console.error(`添加文章 ${article.title} 失败:`, error)
      }
    }
    
    // 显示结果
    if (addedCount.length > 0) {
      ElMessage.success(`成功添加 ${addedCount.length} 篇文章到专题`)
    }
    if (failedArticles.length > 0) {
      const failedTitles = failedArticles.map(f => `${f.title}(${f.error})`).join('、')
      ElMessage.warning(`以下文章添加失败：${failedTitles}`)
    }
    
    if (addedCount.length > 0) {
      addArticleDialogVisible.value = false
      selectedArticles.value = []
      // 重新加载当前页
      loadTopicArticles({ id: currentTopic.value.id }, articlePage.value, articlePageSize.value)
      // 更新专题列表中的文章数量
      updateTopicArticleCount(currentTopic.value.id)
    }
  } catch (error) {
    const errorMessage = error.message || error.response?.data?.message || '添加文章失败'
    ElMessage.error(errorMessage)
    console.error(error)
  }
}

const handleRemoveArticle = (row) => {
  ElMessageBox.confirm(
    '确定要将此文章从专题中移除吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await removeArticleFromTopic({ id: currentTopic.value.id, articleId: row.articleId })
      ElMessage.success('文章已从专题中移除')
      // 重新加载当前页
      loadTopicArticles({ id: currentTopic.value.id }, articlePage.value, articlePageSize.value)
      // 更新专题列表中的文章数量
      updateTopicArticleCount(currentTopic.value.id)
    } catch (error) {
      ElMessage.error('移除文章失败')
      console.error(error)
    }
  }).catch(() => {})
}

// 上移文章
const handleMoveUp = async (index) => {
  if (index === 0) return
  
  try {
    const current = topicArticles.value[index]
    const previous = topicArticles.value[index - 1]
    
    // 交换排序值
    const tempSortOrder = current.sortOrder
    current.sortOrder = previous.sortOrder
    previous.sortOrder = tempSortOrder
    
    // 更新数据库
    await Promise.all([
      updateTopicArticleSort({ id: current.id, sortOrder: current.sortOrder }),
      updateTopicArticleSort({ id: previous.id, sortOrder: previous.sortOrder })
    ])
    
    // 重新加载列表以保持顺序
    await loadTopicArticles({ id: currentTopic.value.id }, articlePage.value, articlePageSize.value)
    ElMessage.success('排序已更新')
  } catch (error) {
    ElMessage.error('更新排序失败')
    console.error(error)
    // 重新加载以恢复原状态
    await loadTopicArticles({ id: currentTopic.value.id }, articlePage.value, articlePageSize.value)
  }
}

// 下移文章
const handleMoveDown = async (index) => {
  if (index === topicArticles.value.length - 1) return
  
  try {
    const current = topicArticles.value[index]
    const next = topicArticles.value[index + 1]
    
    // 交换排序值
    const tempSortOrder = current.sortOrder
    current.sortOrder = next.sortOrder
    next.sortOrder = tempSortOrder
    
    // 更新数据库
    await Promise.all([
      updateTopicArticleSort({ id: current.id, sortOrder: current.sortOrder }),
      updateTopicArticleSort({ id: next.id, sortOrder: next.sortOrder })
    ])
    
    // 重新加载列表以保持顺序
    await loadTopicArticles({ id: currentTopic.value.id }, articlePage.value, articlePageSize.value)
    ElMessage.success('排序已更新')
  } catch (error) {
    ElMessage.error('更新排序失败')
    console.error(error)
    // 重新加载以恢复原状态
    await loadTopicArticles({ id: currentTopic.value.id }, articlePage.value, articlePageSize.value)
  }
}

// 分页变化处理
const handleArticlePageChange = (page) => {
  articlePage.value = page
  loadTopicArticles({ id: currentTopic.value.id }, articlePage.value, articlePageSize.value)
}

const handleArticlePageSizeChange = (size) => {
  articlePageSize.value = size
  articlePage.value = 1 // 重置到第一页
  loadTopicArticles({ id: currentTopic.value.id }, articlePage.value, articlePageSize.value)
}

// 更新专题列表中的文章数量（优化：只更新单个专题）
const updateTopicArticleCount = async (topicId) => {
  try {
    const articlesRes = await getTopicArticles({ id: topicId, page: 1, size: 100 })
    // 调试日志
    console.log(`专题 ${topicId} 的文章响应:`, articlesRes)
    // 确保响应数据存在
    if (articlesRes) {
      // 尝试多种可能的字段名
      const count = articlesRes.total ?? articlesRes.Total ?? (articlesRes.list?.length ?? 0)
      // 更新topicData中对应专题的文章数量
      const topic = topicData.value.find(t => t.id === topicId)
      if (topic) {
        topic.articleCount = count
        console.log(`专题 ${topicId} 的文章数量已更新为:`, count)
      }
    } else {
      console.warn(`获取专题 ${topicId} 的文章数量失败: 响应数据为空`, articlesRes)
    }
  } catch (error) {
    console.error(`更新专题 ${topicId} 的文章数量失败:`, error)
  }
}

// 格式化日期（使用 dayjs）
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  try {
    // 处理各种日期格式
    let date = null
    const dateString = String(dateStr).trim()
    
    // 尝试多种日期格式解析
    const formats = [
      'YYYY-MM-DD HH:mm:ss',
      'YYYY-MM-DDTHH:mm:ssZ',
      'YYYY-MM-DDTHH:mm:ss.SSSZ',
      'YYYY-MM-DD HH:mm:ss.SSS',
      'YYYY-MM-DD',
      'YYYY/MM/DD HH:mm:ss',
      'YYYY/MM/DD'
    ]
    
    // 先尝试自动解析
    date = dayjs(dateString)
    
    // 如果自动解析失败，尝试使用自定义格式
    if (!date.isValid()) {
      for (const format of formats) {
        date = dayjs(dateString, format, true) // strict mode
        if (date.isValid()) break
      }
    }
    
    // 如果还是失败，尝试清理字符串后解析
    if (!date.isValid()) {
      const cleanedStr = dateString.replace(/T/, ' ').replace(/Z/, '').replace(/\+.*$/, '').trim()
      date = dayjs(cleanedStr)
    }
    
    // 如果还是失败，尝试提取日期部分
    if (!date.isValid()) {
      const dateMatch = dateString.match(/(\d{4}-\d{2}-\d{2})/)
      if (dateMatch) {
        date = dayjs(dateMatch[1], 'YYYY-MM-DD', true)
      }
    }
    
    // 如果所有解析都失败
    if (!date.isValid()) {
      console.warn('无法解析日期:', dateStr)
      // 返回原始字符串的日期部分
      const match = dateString.match(/(\d{4}-\d{2}-\d{2})/)
      return match ? match[0] : dateStr
    }
    
    // 使用 startOf('day') 确保只比较日期部分，忽略时间
    const now = dayjs().startOf('day')
    const targetDate = date.startOf('day')
    const diffDays = now.diff(targetDate, 'day')
    
    
 
    if (diffDays > 0 && diffDays < 7) return `${diffDays}天前`
    if (diffDays > 0 && diffDays < 30) return `${Math.floor(diffDays / 7)}周前`
    if (diffDays > 0 && diffDays < 365) return `${Math.floor(diffDays / 30)}个月前`
    
    // 如果日期是未来的（不应该发生，但处理一下）
    if (diffDays < 0) {
      return date.format('YYYY-MM-DD')
    }
    
    // 超过一年显示具体日期
    return date.format('YYYY-MM-DD')
  } catch (error) {
    console.error('日期格式化错误:', error, dateStr)
    // 如果格式化失败，尝试返回原始字符串的日期部分
    if (typeof dateStr === 'string') {
      const match = dateStr.match(/(\d{4}-\d{2}-\d{2})/)
      if (match) return match[0]
    }
    return String(dateStr)
  }
}

// 格式化数字（用于阅读量等）
const formatNumber = (num) => {
  if (!num && num !== 0) return '0'
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'W'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toString()
}

// 格式化预览日期（显示完整日期，使用 dayjs）
const formatPreviewDate = (dateStr) => {
  if (!dateStr) return ''
  try {
    let date = dayjs(dateStr)
    
    // 如果解析失败，尝试其他格式
    if (!date.isValid()) {
      const cleanedStr = String(dateStr).replace(/T/, ' ').replace(/Z/, '').replace(/\+.*$/, '')
      date = dayjs(cleanedStr)
      
      if (!date.isValid()) {
        const dateMatch = String(dateStr).match(/(\d{4}-\d{2}-\d{2})/)
        if (dateMatch) {
          date = dayjs(dateMatch[1])
        }
      }
      
      if (!date.isValid()) {
        // 如果还是失败，返回原始字符串的日期部分
        const match = String(dateStr).match(/(\d{4}-\d{2}-\d{2})/)
        return match ? match[0] : dateStr
      }
    }
    
    return date.format('YYYY-MM-DD')
  } catch (error) {
    console.error('日期格式化错误:', error, dateStr)
    // 如果格式化失败，尝试返回原始字符串的日期部分
    if (typeof dateStr === 'string') {
      const match = dateStr.match(/(\d{4}-\d{2}-\d{2})/)
      if (match) return match[0]
    }
    return String(dateStr)
  }
}

// 监听分页变化
watch([topicCurrentPage, topicPageSize], () => {
  handleTopicSearch()
})

onMounted(() => {
  handleTopicSearch()
})
</script>

<style scoped>
.topic-manage {
  /* padding: 20px; handled by layout */
}

/* Topic Detail Dialog */
.topic-detail-dialog {
  --el-dialog-border-radius: 16px;
}

.detail-dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  padding: 0;
}

.detail-dialog-header .header-left {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  flex: 1;
}

.detail-dialog-header .header-info {
  flex: 1;
}

.detail-dialog-header .topic-title {
  margin: 0 0 8px;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  line-height: 1.4;
}

.detail-dialog-header .topic-description {
  margin: 0;
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}

.detail-dialog-header .header-right {
  flex-shrink: 0;
}

.detail-content {
  padding: 0;
}

.detail-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
  border-radius: 12px;
  border: 1px solid #f0f0f0;
  transition: all 0.3s;
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  border-color: #409eff;
}

.stat-item > .el-icon {
  font-size: 32px;
  color: #409eff;
  background: linear-gradient(135deg, #e8f4ff 0%, #d4e7ff 100%);
  padding: 12px;
  border-radius: 10px;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
}

.articles-section {
  margin-top: 20px;
}

.section-header {
  margin-bottom: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.section-title .el-icon {
  color: #409eff;
  font-size: 18px;
}

/* 专题文章表格样式优化 */
.topic-articles-table {
  border-radius: 8px;
  overflow: hidden;
}

.topic-articles-table :deep(.el-table__header) {
  background: #fafafa;
}

.topic-articles-table :deep(.el-table__header th) {
  padding: 12px 8px;
  font-weight: 600;
  font-size: 13px;
}

.topic-articles-table :deep(.el-table__header th.is-center) {
  text-align: center;
}

.topic-articles-table :deep(.el-table__body td) {
  padding: 10px 8px;
  vertical-align: middle;
}

.topic-articles-table :deep(.el-table__body td.is-center) {
  text-align: center;
}

.topic-articles-table :deep(.el-table__row:hover > td) {
  background-color: #f5f7fa !important;
}

.topic-articles-table :deep(.el-table__row) {
  transition: background-color 0.2s ease;
}

.sort-index {
  font-weight: 600;
  color: #606266;
  text-align: center;
}




.topic-articles-table :deep(.el-tag) {
  margin: 0;
  font-weight: 600;
  font-size: 12px;
  padding: 2px 8px;
}

.filter-card {
  margin-bottom: 20px;
  border-radius: 12px;
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  background: linear-gradient(to bottom, #ffffff 0%, #fafbfc 100%);
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}

.filter-form .el-form-item {
  margin-bottom: 0;
  margin-right: 20px;
}

/* 搜索栏输入框和按钮统一样式 */
.filter-form :deep(.el-input__wrapper) {
  height: 32px;
}

.filter-form :deep(.el-button) {
  height: 32px;
  padding: 8px 15px;
}

/* Topic Grid Layout */
.topic-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
  padding: 4px;
}

.topic-card-wrapper {
  animation: fadeInUp 0.4s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.topic-card {
  border: none;
  border-radius: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
  border: 1px solid #f0f0f0;
}

.topic-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
}

.topic-cover {
  height: 150px;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.topic-cover .el-image {
  width: 100%;
  height: 100%;
  transition: transform 0.4s ease;
}

.topic-card:hover .topic-cover .el-image {
  transform: scale(1.08);
}

.image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-size: 48px;
  opacity: 0.8;
}

.topic-actions-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(180deg, rgba(0, 0, 0, 0.3) 0%, rgba(0, 0, 0, 0.7) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  opacity: 0;
  transition: all 0.3s ease;
}

.topic-cover:hover .topic-actions-overlay {
  opacity: 1;
}

.topic-actions-overlay .el-button {
  transform: translateY(10px);
  transition: all 0.3s ease;
}

.topic-cover:hover .topic-actions-overlay .el-button {
  transform: translateY(0);
}

.topic-actions-overlay .el-button:nth-child(1) {
  transition-delay: 0.05s;
}

.topic-actions-overlay .el-button:nth-child(2) {
  transition-delay: 0.1s;
}

.topic-actions-overlay .el-button:nth-child(3) {
  transition-delay: 0.15s;
}

.topic-content {
  padding: 12px;
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #fff;
}

.topic-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
  gap: 12px;
  min-height: 24px;
}

.topic-name {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  flex: 1;
  min-width: 0;
  transition: color 0.3s;
  padding-right: 8px;
}

.topic-card:hover .topic-name {
  color: #409eff;
}

.topic-desc {
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
  margin: 0 0 auto;
  min-height: 60px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  flex: 1;
}

.topic-footer {
  padding-top: 10px;
  margin-top: auto;
  border-top: 1px solid #f0f2f5;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  min-height: 36px;
}

.date {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #606266;
  padding: 5px 10px;
  border-radius: 8px;
  background: linear-gradient(135deg, #f8f9fa 0%, #f0f2f5 100%);
  border: 1px solid rgba(0, 0, 0, 0.04);
}

.date-icon {
  font-size: 13px;
  color: #909399;
  opacity: 0.85;
  line-height: 1;
}

.date-text {
  line-height: 1.3;
  font-weight: 500;
  letter-spacing: 0.2px;
}

.sort-badge {
  font-size: 12px;
  color: #909399;
  background: #f5f7fa;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

.article-count-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  background: linear-gradient(135deg, rgba(64, 158, 255, 0.12) 0%, rgba(64, 158, 255, 0.08) 100%);
  border: 1px solid rgba(64, 158, 255, 0.2);
  border-radius: 8px;
  white-space: nowrap;
  flex-shrink: 0;
  cursor: default;
  line-height: 1.2;
}

.article-count-badge .count-icon {
  font-size: 14px;
  color: #409eff;
  line-height: 1;
  opacity: 0.95;
}

.article-count-badge .count-text {
  font-size: 12px;
  font-weight: 600;
  color: #409eff;
  line-height: 1.3;
  white-space: nowrap;
  letter-spacing: 0.3px;
  text-shadow: 0 1px 2px rgba(64, 158, 255, 0.1);
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}

/* Table Card */
.table-card {
  margin-top: 20px;
  border-radius: 16px;
  border: none;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  overflow: hidden;
  border: 1px solid #f0f0f0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header .header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.card-header .header-icon {
  color: #409eff;
}

.card-header .header-title {
  font-weight: 600;
  font-size: 16px;
  color: #303133;
}

.count-tag {
  background: linear-gradient(135deg, #e8f4ff 0%, #d4e7ff 100%);
  color: #409eff;
  border: none;
  font-weight: 600;
  padding: 8px 16px;
}

.article-title-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
}

.article-title-cell .el-icon {
  flex-shrink: 0;
}

.title-text {
  font-weight: 500;
  color: #303133;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.views-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: #606266;
  font-weight: 500;
}

.views-cell .el-icon {
  color: #909399;
  font-size: 14px;
}

.date-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.date-icon-small {
  font-size: 12px;
  color: #909399;
  opacity: 0.7;
}

.date-text {
  color: #909399;
  font-size: 12px;
  font-weight: 500;
}

/* Form Tips */
.form-tip {
  margin-left: 10px;
  font-size: 12px;
  color: #909399;
}

/* 封面图片容器 */
.cover-image-container {
  width: 100%;
}

/* 输入框组 */
.cover-input-group {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  margin-bottom: 0;
}

.cover-input {
  flex: 1;
  min-width: 0;
}

.preview-btn {
  flex-shrink: 0;
  min-width: 100px;
  white-space: nowrap;
}

.preview-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 图片预览容器 */
.image-preview-container {
  margin-top: 16px;
  padding: 16px;
  background: #fafbfc;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  transition: all 0.3s ease;
}

.image-preview-container:hover {
  border-color: #c0c4cc;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e4e7ed;
}

.preview-label {
  font-size: 14px;
  font-weight: 600;
  color: #606266;
  display: flex;
  align-items: center;
  gap: 6px;
}

.preview-label::before {
  content: '';
  width: 4px;
  height: 16px;
  background: #409eff;
  border-radius: 2px;
}

.preview-actions {
  display: flex;
  gap: 8px;
}

.preview-content {
  display: flex;
  justify-content: center;
  align-items: center;
}

.preview-image {
  width: 100%;
  max-width: 500px;
  height: 280px;
  border-radius: 8px;
  border: 2px solid #e4e7ed;
  transition: all 0.3s ease;
  overflow: hidden;
  cursor: pointer;
  background: #fff;
}

.preview-image:hover {
  border-color: #409eff;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.25);
  transform: translateY(-2px);
}

.preview-image :deep(.el-image__inner) {
  transition: transform 0.3s ease;
}

.preview-image:hover :deep(.el-image__inner) {
  transform: scale(1.03);
}

.image-slot {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8ecf0 100%);
  color: #909399;
  font-size: 14px;
  gap: 8px;
}

.image-slot .el-icon {
  font-size: 48px;
  opacity: 0.6;
}

.image-slot span {
  font-size: 13px;
}

/* Table Styling */
:deep(.el-table) {
  border-radius: 8px;
}

:deep(.el-table th.el-table__cell) {
  background: #fafafa;
  color: #303133;
  font-weight: 600;
}

:deep(.el-table tbody tr:hover > td) {
  background: #f5f7fa !important;
}

/* Button Styling */
:deep(.el-button) {
  border-radius: 8px;
  font-weight: 500;
}

:deep(.el-button--primary) {
  background: linear-gradient(135deg, #409eff 0%, #5c8efd 100%);
  border: none;
}

:deep(.el-button--primary:hover) {
  background: linear-gradient(135deg, #66b1ff 0%, #79a3ff 100%);
}

/* Dialog Styling */
:deep(.el-dialog) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.el-dialog__header) {
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
}

:deep(.el-dialog__title) {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-dialog__footer) {
  padding: 16px 24px;
  border-top: 1px solid #f0f0f0;
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.footer-left {
  flex: 1;
}

.footer-right {
  display: flex;
  gap: 12px;
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

/* 数字输入框样式 */
:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-input-number .el-input__wrapper) {
  width: 100%;
}

/* 开关样式 */
.featured-switch {
  height: 32px;
}

.featured-switch :deep(.el-switch__core) {
  width: 60px !important;
  height: 28px !important;
  border-radius: 14px !important;
}

.featured-switch :deep(.el-switch__core::after) {
  width: 24px !important;
  height: 24px !important;
}

.featured-switch.is-checked :deep(.el-switch__core::after) {
  left: calc(100% - 26px) !important;
}

.featured-switch :deep(.el-switch__label) {
  font-size: 13px;
  font-weight: 500;
}

/* 封面预览对话框 */
.cover-preview-dialog {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
  background: #f5f7fa;
  border-radius: 8px;
  padding: 20px;
}

/* Input & Select Styling */
:deep(.el-select .el-input__wrapper) {
  border-radius: 6px;
}

/* Pagination */
:deep(.el-pagination) {
  justify-content: center;
  margin-top: 32px;
}

:deep(.el-pagination .el-pager li) {
  border-radius: 6px;
}

/* Empty State */
:deep(.el-empty) {
  padding: 60px 0;
}

:deep(.el-empty__description) {
  margin-top: 16px;
  color: #909399;
}

/* Tag Styling */
:deep(.el-tag) {
  border-radius: 6px;
  font-weight: 500;
  padding: 4px 12px;
}

/* Dialog Search Area */
.dialog-search-area {
  background: #f8f9fa;
  padding: 16px;
  border-radius: 12px;
  margin-bottom: 20px;
}

.dialog-search-area .filter-form {
  margin: 0;
}

.dialog-search-area .el-form-item {
  margin-bottom: 0;
}

/* Article Dialog */
.article-dialog :deep(.el-dialog__body) {
  padding: 20px 24px;
}

/* Dialog Footer */
.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 12px 0;
}

.footer-info {
  flex: 1;
}

.selection-info {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #67c23a;
  font-size: 14px;
  font-weight: 500;
  background: #f0f9ff;
  padding: 8px 16px;
  border-radius: 8px;
  border: 1px solid #b3e19d;
}

.footer-actions {
  display: flex;
  gap: 12px;
}

/* Loading State */
.topic-list {
  min-height: 400px;
}

/* Empty State in Dialog */
:deep(.el-table__empty-block) {
  padding: 60px 0;
}

:deep(.el-empty__image) {
  width: 180px;
}

/* Smooth Transitions */
* {
  transition-property: background-color, border-color, color, fill, stroke, opacity, box-shadow, transform;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
}

/* Hover Effects */
.el-button.is-circle:hover {
  transform: scale(1.1) rotate(5deg);
}

/* Input Focus Glow */
:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #409eff, 0 0 0 4px rgba(64, 158, 255, 0.1);
}

/* 专题预览对话框样式 */
.topic-preview-dialog {
  --el-dialog-border-radius: 16px;
}

.preview-dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0;
}

.preview-dialog-header .header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.preview-dialog-header .header-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.preview-dialog-header .header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.topic-preview-container {
  min-height: 600px;
  background: #f8f9fa;
  border-radius: 12px;
  overflow: hidden;
}

/* 预览头部 */
.preview-header {
  position: relative;
  height: 400px;
  background-color: #2c3e50;
  color: white;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  background-size: cover;
  background-position: center;
}

.header-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(180deg, rgba(0, 0, 0, 0.4) 0%, rgba(0, 0, 0, 0.7) 100%);
}

.header-content {
  position: relative;
  z-index: 1;
  text-align: center;
  max-width: 900px;
  padding: 0 20px;
}

.preview-meta-info {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  flex-wrap: wrap;
}

.featured-tag {
  background-color: rgba(103, 194, 58, 0.9);
  border-color: rgba(103, 194, 58, 1);
  backdrop-filter: blur(10px);
  font-weight: 600;
}

.preview-title {
  font-size: 32px;
  font-weight: 800;
  line-height: 1.3;
  margin-bottom: 12px;
  text-shadow: 0 2px 20px rgba(0,0,0,0.3);
}

.preview-description {
  font-size: 15px;
  line-height: 1.5;
  margin-bottom: 16px;
  opacity: 0.95;
  text-shadow: 0 1px 4px rgba(0,0,0,0.2);
}

.preview-meta {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 13px;
  flex-wrap: wrap;
}

.preview-meta .meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.preview-meta .divider {
  opacity: 0.6;
}

/* 预览内容区域 */
.preview-content-wrapper {
  margin-top: -60px;
  position: relative;
  z-index: 2;
  padding: 0 20px 40px;
}

.preview-intro-card,
.preview-articles-card {
  margin-bottom: 24px;
  border-radius: 16px;
  border: none;
  box-shadow: 0 2px 12px rgba(0,0,0,0.08);
}

.preview-intro-card :deep(.el-card__body),
.preview-articles-card :deep(.el-card__body) {
  padding: 32px;
}

.preview-intro-card .section-title,
.preview-articles-card .section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin-bottom: 24px;
}

.preview-intro-card .section-title .el-icon,
.preview-articles-card .section-title .el-icon {
  font-size: 22px;
  color: #409eff;
}

.preview-intro-card .count {
  font-size: 16px;
  font-weight: 500;
  color: #909399;
}

/* Markdown 样式 */
.markdown-body {
  font-size: 16px;
  line-height: 1.8;
  color: #333;
}

.markdown-body :deep(h2) {
  margin-top: 24px;
  margin-bottom: 12px;
  font-size: 22px;
  font-weight: 700;
  color: #1a1a1a;
  border-bottom: 1px solid #eee;
  padding-bottom: 8px;
}

.markdown-body :deep(h3) {
  margin-top: 20px;
  margin-bottom: 10px;
  font-size: 18px;
  font-weight: 700;
  color: #1a1a1a;
}

.markdown-body :deep(p) {
  margin-bottom: 16px;
}

.markdown-body :deep(ul) {
  padding-left: 24px;
  margin-bottom: 16px;
}

.markdown-body :deep(li) {
  margin-bottom: 8px;
}

/* 预览文章列表 */
.preview-articles-card .articles-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.preview-articles-card .articles-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.preview-articles-card .article-item {
  display: flex;
  gap: 20px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 12px;
  transition: all 0.3s ease;
}

.preview-articles-card .article-item:hover {
  background: #e9ecef;
  transform: translateX(4px);
}

.preview-articles-card .article-number {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 8px;
  font-weight: 700;
  font-size: 16px;
}

.preview-articles-card .article-content {
  flex: 1;
  min-width: 0;
}

.preview-articles-card .article-top {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 8px;
}

.preview-articles-card .article-title {
  flex: 1;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin: 0;
  line-height: 1.4;
}

.preview-articles-card .article-summary {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.preview-articles-card .article-meta {
  display: flex;
  align-items: center;
  gap: 20px;
  font-size: 13px;
  color: #909399;
}

.preview-articles-card .article-meta .meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.preview-articles-card .article-cover {
  flex-shrink: 0;
  width: 180px;
  height: 120px;
  border-radius: 8px;
  overflow: hidden;
}

.preview-articles-card .article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.preview-articles-card .article-item:hover .article-cover img {
  transform: scale(1.05);
}

/* 响应式 */
@media (max-width: 1200px) {
  .topic-grid {
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  }
}

@media (max-width: 768px) {
  .topic-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .filter-card {
    margin-bottom: 16px;
  }
  
  .topic-cover {
    height: 160px;
  }
  
  .filter-form {
    flex-direction: column;
  }
  
  .filter-form .el-form-item {
    width: 100%;
    margin-right: 0;
    margin-bottom: 12px;
  }
  
  .filter-form .el-form-item:last-child {
    margin-bottom: 0;
  }
  
  .article-dialog :deep(.el-dialog) {
    width: 95% !important;
  }
  
  .topic-detail-dialog :deep(.el-dialog) {
    width: 95% !important;
  }
  
  .detail-dialog-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .detail-stats {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .stat-item {
    padding: 16px;
  }
  
  .stat-item > .el-icon {
    font-size: 24px;
    padding: 8px;
  }
  
  .stat-value {
    font-size: 20px;
  }
  
  .topic-preview-dialog :deep(.el-dialog) {
    width: 95% !important;
  }
  
  .preview-header {
    height: 300px;
  }
  
  .preview-title {
    font-size: 24px;
  }
  
  .preview-description {
    font-size: 14px;
  }
  
  .preview-content-wrapper {
    padding: 0 16px 30px;
  }
  
  .preview-articles-card .article-item {
    flex-direction: column-reverse;
  }
  
  .preview-articles-card .article-cover {
    width: 100%;
    height: 200px;
  }
}

@media (max-width: 480px) {
  .topic-content {
    padding: 16px;
  }
  
  .topic-header {
    flex-wrap: wrap;
    gap: 8px;
  }
  
  .topic-name {
    font-size: 16px;
    padding-right: 0;
    width: 100%;
  }
  
  .date {
    padding: 3px 6px;
    font-size: 11px;
    gap: 4px;
  }
  
  .date-icon {
    font-size: 11px;
  }
  
  .date-text {
    font-size: 11px;
  }
  
  .article-count-badge {
    padding: 3px 8px;
    gap: 4px;
  }
  
  .article-count-badge .count-icon {
    font-size: 11px;
  }
  
  .article-count-badge .count-text {
    font-size: 11px;
  }
  
  .topic-desc {
    font-size: 13px;
    min-height: 60px;
  }
  
  .topic-cover {
    height: 140px;
  }
  
  .topic-footer {
    flex-wrap: wrap;
    gap: 8px;
  }
}
</style>


