<template>
  <view class="home">
    <!-- 品牌区 -->
    <view class="brand-bar">
      <image class="brand-logo" :src="entity.logo_url || '/static/logo.png'" mode="aspectFit" />
      <text class="brand-name">{{ entity.name || '扬昵 Yanny' }}</text>
    </view>

    <!-- 搜索 -->
    <view class="search-bar" @click="onSearch">
      <text class="search-icon">🔍</text>
      <text class="search-placeholder">搜索视频</text>
    </view>

    <!-- 功能入口 -->
    <view class="func-row">
      <view class="func-item" v-for="f in funcBtns" :key="f.label" @click="f.action">
        <view class="func-icon">{{ f.icon }}</view>
        <text class="func-label">{{ f.label }}</text>
      </view>
    </view>

    <!-- 分类标签 -->
    <scroll-view scroll-x class="category-scroll" :show-scrollbar="false">
      <view class="category-list">
        <view
          v-for="c in allCategory"
          :key="c.id"
          class="category-tag"
          :class="{ active: activeCategory === c.id }"
          @click="switchCategory(c.id)"
        >
          {{ c.name }}
        </view>
      </view>
    </scroll-view>

    <!-- 视频双列流 -->
    <view class="video-grid" v-if="videos.length">
      <view
        v-for="v in videos"
        :key="v.id"
        class="video-card"
        @click="goPlay(v.id)"
      >
        <image class="video-cover" :src="v.cover_url" mode="aspectFill" lazy-load />
        <view class="video-info">
          <text class="video-title">{{ v.title }}</text>
          <view class="video-meta">
            <text class="meta-item">▶ {{ formatCount(v.view_count) }}</text>
            <text class="meta-item">❤ {{ formatCount(v.like_count) }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 加载更多 -->
    <view class="load-more" v-if="loading">
      <text>加载中...</text>
    </view>
    <view class="load-more" v-else-if="hasMore" @click="loadMore">
      <text>加载更多</text>
    </view>

    <!-- 空状态 -->
    <view class="empty" v-if="!loading && !videos.length">
      <text>暂无视频</text>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted, onReachBottom } from 'vue'
import { useVideoStore } from '@/store/video'
import { useUserStore } from '@/store/user'

const videoStore = useVideoStore()
const userStore = useUserStore()

const entity = ref({ name: '扬昵 Yanny', logo_url: '' })
const activeCategory = ref(0)
const funcBtns = [
  { icon: '🔥', label: '热议话题', action: () => {} },
  { icon: '💬', label: '圈子广场', action: () => {} },
  { icon: '🎧', label: '在线客服', action: () => {} },
  { icon: 'ℹ️', label: '关于我们', action: () => {} },
]

const videos = computed(() => videoStore.videos)
const loading = computed(() => videoStore.loading)
const hasMore = computed(() => videoStore.hasMore)

const allCategory = computed(() => [
  { id: 0, name: '精选' },
  ...videoStore.categories,
])

function switchCategory(id) {
  activeCategory.value = id
  videoStore.fetchVideos(id)
}

function loadMore() {
  videoStore.fetchNextPage()
}

function goPlay(videoId) {
  uni.navigateTo({ url: `/pages/player/player?videoId=${videoId}&categoryId=${activeCategory.value}` })
}

function onSearch() {
  // TODO: 搜索页
}

function formatCount(n) {
  if (!n) return '0'
  if (n >= 10000) return (n / 10000).toFixed(1) + 'w'
  return String(n)
}

// 上拉触底
onReachBottom(() => {
  if (hasMore.value) loadMore()
})

onMounted(() => {
  videoStore.fetchVideos(0)
  videoStore.fetchCategories(1, 1)
})
</script>

<style scoped>
.home { padding-bottom: 20rpx; background: #fff; min-height: 100vh; }
.brand-bar { display: flex; align-items: center; padding: 20rpx 24rpx; }
.brand-logo { width: 72rpx; height: 72rpx; border-radius: 16rpx; }
.brand-name { font-size: 36rpx; font-weight: 700; margin-left: 16rpx; color: #333; }
.search-bar { margin: 0 24rpx 16rpx; padding: 16rpx 20rpx; background: #f5f5f5; border-radius: 32rpx; display: flex; align-items: center; }
.search-icon { font-size: 28rpx; margin-right: 8rpx; }
.search-placeholder { color: #999; font-size: 28rpx; }
.func-row { display: flex; justify-content: space-around; padding: 8rpx 24rpx 20rpx; }
.func-item { display: flex; flex-direction: column; align-items: center; }
.func-icon { font-size: 44rpx; margin-bottom: 6rpx; }
.func-label { font-size: 22rpx; color: #666; }
.category-scroll { white-space: nowrap; padding: 0 24rpx 16rpx; }
.category-list { display: inline-flex; gap: 16rpx; }
.category-tag { padding: 10rpx 28rpx; border-radius: 28rpx; font-size: 26rpx; color: #666; background: #f5f5f5; display: inline-block; }
.category-tag.active { background: #333; color: #fff; }
.video-grid { display: flex; flex-wrap: wrap; padding: 0 16rpx; gap: 12rpx; }
.video-card { width: calc(50% - 6rpx); border-radius: 12rpx; overflow: hidden; background: #f5f5f5; }
.video-cover { width: 100%; height: 360rpx; display: block; }
.video-info { padding: 12rpx; }
.video-title { font-size: 26rpx; color: #333; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.video-meta { display: flex; gap: 20rpx; margin-top: 6rpx; }
.meta-item { font-size: 20rpx; color: #999; }
.load-more { text-align: center; padding: 20rpx; color: #999; font-size: 24rpx; }
.empty { text-align: center; padding: 80rpx; color: #999; }
</style>
