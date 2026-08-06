<template>
  <div class="dashboard">
    <h2>数据看板</h2>

    <!-- 核心指标 -->
    <el-row :gutter="20" class="stat-row">
      <el-col :xs="12" :sm="6" v-for="card in statCards" :key="card.key">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" :style="{ background: card.color }">
            <el-icon :size="24"><component :is="card.icon" /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">{{ card.label }}</div>
            <div class="stat-value">{{ card.value }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 趋势图 + 排行 -->
    <el-row :gutter="20" style="margin-top:20px">
      <el-col :xs="24" :lg="16">
        <el-card>
          <template #header><span>播放量趋势（近 7 天）</span></template>
          <div class="chart-area">
            <div v-if="loading" class="chart-empty"><el-icon :size="40"><Loading /></el-icon></div>
            <div v-else-if="!trendData.length" class="chart-empty">暂无数据</div>
            <div v-else class="bar-chart">
              <div v-for="d in trendData" :key="d.date" class="bar-col">
                <div class="bar-val">{{ formatNum(d.view_count) }}</div>
                <div class="bar-fill" :style="{ height: barHeight(d.view_count) + '%' }"></div>
                <div class="bar-label">{{ d.date.slice(5) }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="8">
        <el-card>
          <template #header><span>热门视频 Top10</span></template>
          <div class="rank-list" v-if="topVideos.length">
            <div v-for="(v, i) in topVideos" :key="v.video_id" class="rank-item">
              <span class="rank-num" :class="{ top3: i < 3 }">{{ i + 1 }}</span>
              <span class="rank-title">{{ v.video_id }}</span>
              <span class="rank-count">{{ formatNum(v.view_count) }}</span>
            </div>
          </div>
          <div class="chart-empty" v-else>暂无数据</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 地域分布 -->
    <el-row :gutter="20" style="margin-top:20px">
      <el-col :span="24">
        <el-card>
          <template #header><span>用户地域分布（近 7 天）</span></template>
          <div class="region-table" v-if="regionData.length">
            <div v-for="r in regionData" :key="r.province" class="region-row">
              <span class="region-name">{{ r.province }}</span>
              <div class="region-bar-wrap">
                <div class="region-bar" :style="{ width: regionWidth(r.view_count) + '%' }"></div>
              </div>
              <span class="region-val">{{ formatNum(r.view_count) }}</span>
            </div>
          </div>
          <div class="chart-empty" v-else>暂无数据</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '@/utils/request'

const loading = ref(false)
const statCards = ref([
  { key: 'views', label: '今日播放量', value: '--', color: '#409EFF', icon: 'VideoPlay' },
  { key: 'users', label: '日活跃用户', value: '--', color: '#67C23A', icon: 'User' },
  { key: 'newUsers', label: '今日新增', value: '--', color: '#E6A23C', icon: 'UserFilled' },
  { key: 'total', label: '总用户数', value: '--', color: '#F56C6C', icon: 'Avatar' },
])
const trendData = ref([])
const topVideos = ref([])
const regionData = ref([])

const barHeight = (v) => {
  const max = Math.max(...trendData.value.map(d => d.total_views || 0), 1)
  return (v / max) * 100
}
const regionWidth = (v) => {
  const max = Math.max(...regionData.value.map(d => d.view_count || 0), 1)
  return (v / max) * 100
}
const formatNum = (n) => {
  if (!n) return '0'
  if (n >= 10000) return (n / 10000).toFixed(1) + 'w'
  return String(n)
}

const fetchDashboard = async () => {
  loading.value = true
  try {
    const res = await request.get('/stats/dashboard')
    const { platform, top_videos, regions } = res.data || {}

    // 指标卡：取最后一天数据
    if (platform?.length) {
      const last = platform[platform.length - 1]
      statCards.value[0].value = formatNum(last.total_views)
      statCards.value[1].value = formatNum(last.active_users)
      statCards.value[2].value = formatNum(last.new_users)
    }
    // 总用户
    const totalRes = await request.get('/stats/dashboard')
    statCards.value[3].value = '--'

    trendData.value = platform || []
    topVideos.value = top_videos || []
    regionData.value = regions || []
  } catch { /* 暂无数据 */ }
  loading.value = false
}

onMounted(fetchDashboard)
</script>

<style scoped>
.dashboard h2 { margin: 0 0 16px; font-size: 20px; }
.stat-row { display: flex; flex-wrap: wrap; }
.stat-card { margin-bottom: 16px; }
.stat-card :deep(.el-card__body) { display: flex; align-items: center; gap: 16px; padding: 20px; }
.stat-icon { width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; color: #fff; flex-shrink: 0; }
.stat-label { font-size: 13px; color: #909399; }
.stat-value { font-size: 28px; font-weight: 700; color: #303133; margin-top: 4px; }
.chart-area { min-height: 280px; }
.chart-empty { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 240px; color: #c0c4cc; }
.bar-chart { display: flex; align-items: flex-end; justify-content: space-around; height: 240px; padding: 0 10px; }
.bar-col { display: flex; flex-direction: column; align-items: center; flex: 1; max-width: 60px; }
.bar-val { font-size: 11px; color: #909399; margin-bottom: 4px; }
.bar-fill { width: 32px; background: linear-gradient(180deg, #409EFF, #79bbff); border-radius: 4px 4px 0 0; min-height: 4px; transition: height 0.5s; }
.bar-label { font-size: 11px; color: #909399; margin-top: 6px; }
.rank-list { padding: 0; }
.rank-item { display: flex; align-items: center; padding: 8px 0; border-bottom: 1px solid #f5f5f5; font-size: 13px; }
.rank-num { width: 24px; text-align: center; color: #909399; }
.rank-num.top3 { color: #e6a23c; font-weight: 700; }
.rank-title { flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.rank-count { color: #909399; margin-left: 8px; }
.region-table { max-height: 400px; overflow-y: auto; }
.region-row { display: flex; align-items: center; padding: 6px 0; font-size: 13px; gap: 12px; }
.region-name { width: 80px; text-align: right; flex-shrink: 0; }
.region-bar-wrap { flex: 1; height: 16px; background: #f0f2f5; border-radius: 8px; overflow: hidden; }
.region-bar { height: 100%; background: linear-gradient(90deg, #409EFF, #79bbff); border-radius: 8px; transition: width 0.5s; min-width: 4px; }
.region-val { width: 50px; color: #909399; flex-shrink: 0; }
</style>
