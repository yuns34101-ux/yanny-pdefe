<template>
  <div class="dashboard">
    <h2>数据看板</h2>
    <p class="subtitle">运营数据概览（后续接入实时统计数据）</p>

    <!-- 核心指标卡片 -->
    <el-row :gutter="20" class="stat-row">
      <el-col :span="6" v-for="card in statCards" :key="card.label">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-card-content">
            <div class="stat-card-icon" :style="{ background: card.color }">
              <el-icon :size="24"><component :is="card.icon" /></el-icon>
            </div>
            <div class="stat-card-info">
              <div class="stat-card-label">{{ card.label }}</div>
              <div class="stat-card-value">{{ card.value }}</div>
              <div class="stat-card-trend" v-if="card.trend">
                <span :class="card.trend > 0 ? 'up' : 'down'">
                  {{ card.trend > 0 ? '↑' : '↓' }} {{ Math.abs(card.trend) }}%
                </span>
                <span class="trend-label">较昨日</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card>
          <template #header><span>播放量趋势（近 30 天）</span></template>
          <div class="chart-placeholder">
            <el-icon :size="48" color="#c0c4cc"><DataAnalysis /></el-icon>
            <p>图表区域 — 接入 ECharts 后展示播放量/用户增长趋势</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header><span>热门视频 Top10</span></template>
          <div class="chart-placeholder" style="height: 300px">
            <el-icon :size="48" color="#c0c4cc"><TrophyBase /></el-icon>
            <p>排行榜区域</p>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header><span>用户地域分布</span></template>
          <div class="chart-placeholder" style="height: 300px">
            <el-icon :size="48" color="#c0c4cc"><Location /></el-icon>
            <p>中国地图热力图 — 各省份用户分布</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header><span>用户活跃度</span></template>
          <div class="chart-placeholder" style="height: 300px">
            <el-icon :size="48" color="#c0c4cc"><Timer /></el-icon>
            <p>时段分布 / 人均在线时长 / 留存率</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const statCards = ref([
  { label: '今日播放量', value: '--', color: '#409EFF', icon: 'VideoPlay', trend: null },
  { label: '日活跃用户', value: '--', color: '#67C23A', icon: 'User', trend: null },
  { label: '今日新增用户', value: '--', color: '#E6A23C', icon: 'UserFilled', trend: null },
  { label: '总用户数', value: '--', color: '#F56C6C', icon: 'Avatar', trend: null },
])
</script>

<style scoped>
.dashboard h2 { margin: 0 0 4px; font-size: 20px; }
.subtitle { color: #909399; font-size: 13px; margin: 0 0 20px; }
.stat-card-content { display: flex; align-items: center; gap: 16px; }
.stat-card-icon {
  width: 56px; height: 56px; border-radius: 12px;
  display: flex; align-items: center; justify-content: center; color: #fff;
}
.stat-card-label { font-size: 13px; color: #909399; }
.stat-card-value { font-size: 28px; font-weight: 700; color: #303133; margin: 4px 0; }
.stat-card-trend { font-size: 12px; }
.stat-card-trend .up { color: #67C23A; }
.stat-card-trend .down { color: #F56C6C; }
.trend-label { color: #909399; margin-left: 4px; }
.chart-placeholder {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  height: 350px; color: #909399; font-size: 14px;
}
</style>
