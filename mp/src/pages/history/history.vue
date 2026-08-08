<template>
  <view class="history-page">
    <view class="video-grid" v-if="videos.length">
      <view
        v-for="v in videos"
        :key="v.id"
        class="video-card"
        @click="goPlay(v.id)"
      >
        <image class="video-cover" :src="v.cover_url" mode="aspectFill" lazy-load />
        <view class="video-like">
          <text class="like-icon">❤</text>
          <text class="like-count">{{ formatCount(v.like_count) }}</text>
        </view>
      </view>
    </view>

    <view class="load-more" v-if="loading">
      <text>加载中...</text>
    </view>
    <view class="load-more" v-else-if="hasMore" @click="loadMore">
      <text>加载更多</text>
    </view>

    <view class="empty" v-if="!loading && !videos.length">
      <text>暂无观看记录</text>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { onReachBottom } from '@dcloudio/uni-app'
import { get } from '@/utils/request'

const videos = ref([])
const page = ref(1)
const total = ref(0)
const pageSize = 10
const loading = ref(false)
const hasMore = ref(false)

async function loadHistory(pageNum = 1) {
  if (loading.value) return
  loading.value = true
  try {
    const res = await get('/history', { page: pageNum, page_size: pageSize })
    const list = res.data || []
    videos.value = pageNum === 1 ? list : [...videos.value, ...list]
    page.value = pageNum
    total.value = res.meta?.total || 0
    hasMore.value = videos.value.length < total.value
  } finally {
    loading.value = false
  }
}

function loadMore() {
  loadHistory(page.value + 1)
}

function goPlay(videoId) {
  uni.navigateTo({ url: `/pages/player/player?videoId=${videoId}` })
}

function formatCount(n) {
  if (!n) return '0'
  if (n >= 10000) return (n / 10000).toFixed(1) + 'w'
  return String(n)
}

onReachBottom(() => {
  if (hasMore.value) loadMore()
})

onMounted(() => loadHistory(1))
</script>

<style scoped>
.history-page { padding-bottom: 20rpx; background: #fff; min-height: 100vh; }
.video-grid { display: flex; flex-wrap: wrap; padding: 4rpx; gap: 4rpx; }
.video-card { position: relative; width: calc(33.333% - 8rpx/3); margin: 0 4rpx 4rpx 0; aspect-ratio: 3/4; background: #f5f5f5; overflow: hidden; }
.video-cover { width: 100%; height: 100%; display: block; }
.video-like { position: absolute; left: 8rpx; bottom: 8rpx; display: flex; align-items: center; gap: 4rpx; }
.like-icon { font-size: 22rpx; color: #fff; }
.like-count { font-size: 20rpx; color: #fff; }
.load-more { text-align: center; padding: 20rpx; color: #999; font-size: 24rpx; }
.empty { text-align: center; padding: 80rpx; color: #999; }
</style>
