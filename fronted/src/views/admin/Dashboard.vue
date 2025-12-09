<template>
  <div class="dashboard-container">
    <!-- 统计卡片 -->
    <el-row :gutter="20">
      <el-col :span="6" v-for="item in stats" :key="item.title">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-title">{{ item.title }}</div>
              <div class="stat-value">
                 <!-- 数字滚动效果可以用 countup.js，这里简化直接显示 -->
                 {{ item.value }}
              </div>
            </div>
            <div class="stat-icon" :style="{ background: item.color }">
              <el-icon><component :is="item.icon" /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="16">
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="card-header">
              <span>访问趋势</span>
              <el-radio-group v-model="period" size="small" @change="handlePeriodChange">
                <el-radio-button label="week">本周</el-radio-button>
                <el-radio-button label="month">本月</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="lineChartRef" style="height: 350px"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="chart-card">
          <template #header>
            <div class="card-header">
              <span>内容分类占比</span>
            </div>
          </template>
          <div ref="pieChartRef" style="height: 350px"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最新动态 -->
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card shadow="hover" class="chart-card">
          <template #header>最新评论</template>
          <div v-if="latestComments.length === 0" class="empty-comments">
            <el-empty description="暂无评论" :image-size="80" />
          </div>
          <div v-for="comment in latestComments" :key="comment.id" class="comment-item">
            <div class="comment-avatar">
              <el-avatar :size="40" :src="comment.userAvatar || 'https://picsum.photos/id/64/100/100'" />
            </div>
            <div class="comment-content">
              <div class="comment-user">
                {{ comment.userName || '匿名用户' }}
                <span class="comment-time">{{ formatTime(comment.time) }}</span>
              </div>
              <div class="comment-text">{{ comment.content }}</div>
              <div class="comment-article" v-if="comment.articleTitle">
                文章：{{ comment.articleTitle }}
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover" class="chart-card">
           <template #header>系统公告</template>
           <el-timeline>
            <el-timeline-item timestamp="2024-05-20" type="primary" placement="top">
              <el-card class="timeline-card">
                <h4>系统版本更新 v1.0.0</h4>
                <p>正式版上线，欢迎使用</p>
              </el-card>
            </el-timeline-item>
            <el-timeline-item timestamp="2024-05-18" placement="top">
              <el-card class="timeline-card">
                <h4>修复已知 Bug</h4>
                <p>修复了移动端显示的若干问题</p>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, markRaw } from 'vue'
import * as echarts from 'echarts'
import { View, Document, ChatLineRound, User } from '@element-plus/icons-vue'
import { getDashboardStats, getDashboardTrend, getDashboardCategoryRatio, getDashboardLatestComments } from '@/api/dashboard'
import { ElMessage } from 'element-plus'

// 统计数据
const stats = ref([
  { title: '总访问量', value: '0', icon: markRaw(View), color: '#409eff' },
  { title: '文章总数', value: '0', icon: markRaw(Document), color: '#67c23a' },
  { title: '评论总数', value: '0', icon: markRaw(ChatLineRound), color: '#e6a23c' },
  { title: '注册用户', value: '0', icon: markRaw(User), color: '#f56c6c' },
])

const period = ref('week')
const lineChartRef = ref(null)
const pieChartRef = ref(null)
let lineChart = null
let pieChart = null

// 最新评论
const latestComments = ref([])

// 加载统计数据
const loadStats = async () => {
  try {
    const res = await getDashboardStats()
    if (res) {
      stats.value[0].value = formatNumber(res.totalViews)
      stats.value[1].value = formatNumber(res.totalArticles)
      stats.value[2].value = formatNumber(res.totalComments)
      stats.value[3].value = formatNumber(res.totalUsers)
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
    ElMessage.error('加载统计数据失败')
  }
}

// 加载访问趋势
const loadTrend = async () => {
  try {
    const res = await getDashboardTrend({ period: period.value })
    if (res && res.data && lineChartRef.value) {
      const dates = res.data.map(item => {
        const date = new Date(item.date)
        return period.value === 'week' 
          ? ['周一', '周二', '周三', '周四', '周五', '周六', '周日'][date.getDay() === 0 ? 6 : date.getDay() - 1]
          : `${date.getMonth() + 1}/${date.getDate()}`
      })
      const viewsData = res.data.map(item => item.views)
      const readsData = res.data.map(item => item.reads)

      if (!lineChart) {
        lineChart = echarts.init(lineChartRef.value)
      }
      lineChart.setOption({
        tooltip: { trigger: 'axis' },
        grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
        xAxis: { type: 'category', boundaryGap: false, data: dates },
        yAxis: { type: 'value' },
        series: [
          {
            name: '访问量',
            type: 'line',
            smooth: true,
            areaStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: 'rgba(64,158,255,0.5)' },
                { offset: 1, color: 'rgba(64,158,255,0.01)' }
              ])
            },
            data: viewsData,
            itemStyle: { color: '#409eff' }
          },
          {
            name: '阅读量',
            type: 'line',
            smooth: true,
            data: readsData,
            itemStyle: { color: '#67c23a' }
          }
        ]
      })
    }
  } catch (error) {
    console.error('加载访问趋势失败:', error)
  }
}

// 加载分类占比
const loadCategoryRatio = async () => {
  try {
    const res = await getDashboardCategoryRatio()
    if (res && res.data && pieChartRef.value) {
      if (!pieChart) {
        pieChart = echarts.init(pieChartRef.value)
      }
      
      const data = res.data.map(item => ({ value: item.value, name: item.name }))
      const dataLength = data.length
      
      // 根据分类数量调整配置
      // 分类较多时，将图例放在右侧，并调整饼图位置
      const legendConfig = dataLength > 5 
        ? {
            orient: 'vertical',
            right: '5%',
            top: 'middle',
            itemWidth: 12,
            itemHeight: 8,
            textStyle: {
              fontSize: 12
            },
            formatter: (name) => {
              const item = data.find(d => d.name === name)
              return item ? `${name} (${item.value})` : name
            }
          }
        : {
            bottom: '5%',
            left: 'center',
            itemWidth: 12,
            itemHeight: 8,
            textStyle: {
              fontSize: 12
            }
          }
      
      // 根据分类数量调整饼图位置和大小
      const pieRadius = dataLength > 5 
        ? ['35%', '60%']  // 分类多时，饼图小一些，给右侧图例留空间
        : ['40%', '70%']  // 分类少时，使用原来的大小
      
      const pieCenter = dataLength > 5
        ? ['35%', '50%']  // 分类多时，饼图左移
        : ['50%', '50%']  // 分类少时，居中
      
      pieChart.setOption({
        tooltip: { 
          trigger: 'item',
          formatter: '{b}: {c} ({d}%)'
        },
        legend: legendConfig,
        series: [
          {
            name: '文章数',
            type: 'pie',
            radius: pieRadius,
            center: pieCenter,
            avoidLabelOverlap: true,
            itemStyle: {
              borderRadius: 10,
              borderColor: '#fff',
              borderWidth: 2
            },
            label: { 
              show: dataLength <= 5,  // 分类少时显示标签
              position: 'outside',
              formatter: '{b}: {c}\n({d}%)',
              fontSize: 11
            },
            labelLine: { 
              show: dataLength <= 5,
              length: 15,
              length2: 10
            },
            emphasis: {
              label: { 
                show: true, 
                fontSize: 14, 
                fontWeight: 'bold' 
              },
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
            },
            data: data
          }
        ]
      })
    }
  } catch (error) {
    console.error('加载分类占比失败:', error)
  }
}

// 加载最新评论
const loadLatestComments = async () => {
  try {
    const res = await getDashboardLatestComments({ limit: 3 })
    if (res && res.list) {
      latestComments.value = res.list
    }
  } catch (error) {
    console.error('加载最新评论失败:', error)
  }
}

// 格式化数字
const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  }
  return num.toLocaleString()
}

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  
  // 处理 GoFrame 格式的时间字符串 (Y-m-d H:i:s)
  let date
  if (typeof timeStr === 'string' && /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/.test(timeStr)) {
    // GoFrame 格式: 2024-01-01 12:00:00，需要转换为 ISO 格式
    date = new Date(timeStr.replace(' ', 'T') + '+08:00')
  } else {
    date = new Date(timeStr)
  }
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return timeStr
  }
  
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

const initCharts = () => {
  loadTrend()
  loadCategoryRatio()
}

const handleResize = () => {
  lineChart?.resize()
  pieChart?.resize()
}

// 监听时间段变化
const handlePeriodChange = () => {
  loadTrend()
}

onMounted(() => {
  loadStats()
  initCharts()
  loadLatestComments()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  lineChart?.dispose()
  pieChart?.dispose()
})
</script>

<style scoped>
.dashboard-container {
  /* padding: 20px; 已在 Layout 中处理 */
}

.stat-card {
  margin-bottom: 20px;
  border: none;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.03);
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 20px rgba(0,0,0,0.06);
}

.stat-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 24px;
  box-shadow: 0 4px 10px rgba(0,0,0,0.1);
}

.chart-row {
  margin-bottom: 20px;
}

.chart-card {
  border: none;
  border-radius: 16px;
  height: 100%;
  box-shadow: 0 4px 12px rgba(0,0,0,0.03);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.comment-item {
  display: flex;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f2f5;
}

.comment-item:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.comment-avatar {
  margin-right: 15px;
}

.comment-content {
  flex: 1;
}

.comment-user {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 5px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  color: #303133;
}

.comment-time {
  font-weight: normal;
  color: #909399;
  font-size: 12px;
  margin-left: 10px;
}

.comment-text {
  font-size: 14px;
  color: #606266;
  line-height: 1.5;
  margin-top: 4px;
}

.comment-article {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.empty-comments {
  padding: 20px 0;
}

.timeline-card {
  border-radius: 12px;
  border: none;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}
</style>
