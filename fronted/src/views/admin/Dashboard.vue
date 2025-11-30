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
              <el-radio-group v-model="period" size="small">
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
          <div v-for="i in 3" :key="i" class="comment-item">
            <div class="comment-avatar">
              <el-avatar :size="40" src="https://picsum.photos/id/64/100/100" />
            </div>
            <div class="comment-content">
              <div class="comment-user">用户{{ i }} <span class="comment-time">2小时前</span></div>
              <div class="comment-text">这篇文章写得很好，对我帮助很大！</div>
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

// 统计数据
const stats = [
  { title: '总访问量', value: '128,930', icon: markRaw(View), color: '#409eff' },
  { title: '文章总数', value: '1,203', icon: markRaw(Document), color: '#67c23a' },
  { title: '评论总数', value: '5,602', icon: markRaw(ChatLineRound), color: '#e6a23c' },
  { title: '注册用户', value: '890', icon: markRaw(User), color: '#f56c6c' },
]

const period = ref('week')
const lineChartRef = ref(null)
const pieChartRef = ref(null)
let lineChart = null
let pieChart = null

const initCharts = () => {
  // 折线图
  if (lineChartRef.value) {
    lineChart = echarts.init(lineChartRef.value)
    lineChart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', boundaryGap: false, data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'] },
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
          data: [820, 932, 901, 934, 1290, 1330, 1320],
          itemStyle: { color: '#409eff' }
        },
        {
          name: '阅读量',
          type: 'line',
          smooth: true,
          data: [620, 712, 701, 734, 1090, 1130, 1120],
          itemStyle: { color: '#67c23a' }
        }
      ]
    })
  }

  // 饼图
  if (pieChartRef.value) {
    pieChart = echarts.init(pieChartRef.value)
    pieChart.setOption({
      tooltip: { trigger: 'item' },
      legend: { bottom: '5%', left: 'center' },
      series: [
        {
          name: 'Access From',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: { show: false, position: 'center' },
          emphasis: {
            label: { show: true, fontSize: 20, fontWeight: 'bold' }
          },
          labelLine: { show: false },
          data: [
            { value: 1048, name: 'Golang' },
            { value: 735, name: 'Vue 3' },
            { value: 580, name: 'Database' },
            { value: 484, name: 'DevOps' },
            { value: 300, name: 'Other' }
          ]
        }
      ]
    })
  }
}

const handleResize = () => {
  lineChart?.resize()
  pieChart?.resize()
}

onMounted(() => {
  initCharts()
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

.comment-user {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 5px;
  display: flex;
  justify-content: space-between;
  width: 100%;
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
}

.timeline-card {
  border-radius: 12px;
  border: none;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}
</style>
