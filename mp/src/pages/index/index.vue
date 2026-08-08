<template>
  <view class="home">
    <!-- 主体介绍卡片 -->
    <view class="entity-card">
      <view class="entity-header">
        <image class="entity-logo" :src="entity?.logo_url || '/static/logo.png'" mode="aspectFill" />
        <view class="entity-basic">
          <text class="entity-name">{{ entity?.name || '格物山夏山野民艺' }}</text>
        </view>
      </view>
      <text class="entity-desc" v-if="entity?.description">{{ entity.description }}</text>
      <text class="entity-line" v-if="entity?.address" @click="navigateToAddress">地址: {{ entity.address }}</text>
      <text class="entity-line entity-contact" v-if="entity?.contact_phone" @click="callPhone">联系方式：{{ entity.contact_phone }}</text>
      <text class="entity-line">{{ entityStore.videoCount }}条原创内容</text>
      <text class="entity-line">{{ entityStore.followerCount }}个朋友关注</text>

      <view class="entity-actions">
        <view class="action-btn" :class="{ followed: entityStore.followed }" @click="handleFollow">
          <text>{{ entityStore.followed ? '✓ 已关注' : '+ 关注' }}</text>
        </view>
        <button class="action-btn outline contact-btn" open-type="contact">
          <text>私信</text>
        </button>
      </view>
    </view>

    <!-- 视频区块标题 -->
    <view class="section-title">视频</view>

    <!-- 分类标签 -->
    <view class="category-row">
      <view
        v-for="c in allCategory"
        :key="c.id"
        class="category-tag"
        :class="{ active: activeCategory === c.id }"
        @click="switchCategory(c.id)"
      >
        <image v-if="c.icon_url" class="category-icon" :src="c.icon_url" mode="aspectFit" />
        <text class="category-icon" v-else>▤</text>
        <text>{{ c.name }}</text>
      </view>
    </view>

    <!-- 视频三列流 -->
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

    <ProfileSetupModal />
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { onReachBottom, onShow, onShareAppMessage } from '@dcloudio/uni-app'
import { useVideoStore } from '@/store/video'
import { useUserStore } from '@/store/user'
import { useEntityStore } from '@/store/entity'
import ProfileSetupModal from '@/components/ProfileSetupModal.vue'

const videoStore = useVideoStore()
const userStore = useUserStore()
const entityStore = useEntityStore()

const entity = computed(() => entityStore.entity)
const activeCategory = ref(0)

const videos = computed(() => videoStore.videos)
const loading = computed(() => videoStore.loading)
const hasMore = computed(() => videoStore.hasMore)

const allCategory = computed(() => [
  { id: 0, name: '全部' },
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

function requireLogin() {
  if (!userStore.isLoggedIn) {
    uni.navigateTo({ url: '/pages/login/login' })
    return false
  }
  return true
}

async function handleFollow() {
  if (!requireLogin()) return
  if (!entityStore.entity?.id) return
  await entityStore.toggleFollow()
}

function callPhone() {
  if (!entity.value?.contact_phone) return
  uni.makePhoneCall({ phoneNumber: entity.value.contact_phone })
}

function navigateToAddress() {
  if (!entity.value?.address) return
  if (entity.value.latitude && entity.value.longitude) {
    uni.openLocation({
      latitude: entity.value.latitude,
      longitude: entity.value.longitude,
      name: entity.value.name,
      address: entity.value.address,
    })
  } else {
    uni.showToast({ title: '该主体暂未配置精确位置', icon: 'none' })
  }
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

let mounted = false

onMounted(async () => {
  await userStore.waitForReady()
  await entityStore.fetchEntityInfo()
  // uni.setNavigationBarTitle({ title: entityStore.entity?.name || '扬昵 Yanny' })
  videoStore.fetchVideos(activeCategory.value)
  videoStore.fetchCategories(1, 1)
  mounted = true
})

// 从播放页返回时静默刷新视频列表，同步点赞/收藏/分享/评论等计数
onShow(() => {
  if (!mounted) return
  videoStore.fetchVideos(activeCategory.value)
})

// 首页分享：卡片用主体信息（logo + 名称），带入邀请人 ID 实现分享裂变
onShareAppMessage(() => ({
  title: entity.value?.name || '格物山夏山野民艺',
  imageUrl: entity.value?.logo_url || '',
  path: `/pages/index/index?inviter=${userStore.userId}`,
}))
</script>

<style scoped>
.home { padding-bottom: 20rpx; background: #fff; min-height: 100vh; }

/* 主体介绍卡片 */
.entity-card { padding: 24rpx; }
.entity-header { display: flex; align-items: center; }
.entity-logo { width: 96rpx; height: 96rpx; border-radius: 50%; flex-shrink: 0; }
.entity-basic { flex: 1; margin-left: 20rpx; min-width: 0; }
.entity-name { display: block; font-size: 34rpx; font-weight: 700; color: #333; }
.entity-desc { display: block; font-size: 26rpx; color: #333; margin-top: 20rpx; line-height: 1.6; white-space: pre-wrap; }
.entity-line { display: block; font-size: 26rpx; color: #333; margin-top: 12rpx; line-height: 1.6; }
.entity-contact { color: #576b95; }
.entity-actions { display: flex; gap: 20rpx; margin-top: 24rpx; }
.action-btn { flex: 1; text-align: center; padding: 16rpx 0; border-radius: 8rpx; background: #f2f2f2; }
.action-btn text { font-size: 26rpx; color: #333; }
.action-btn.followed { background: #f2f2f2; }
.action-btn.outline { background: #f2f2f2; }
.contact-btn { margin: 0; line-height: normal; border: none; }
.contact-btn::after { border: none; }
.section-title { font-size: 30rpx; font-weight: 700; color: #333; padding: 24rpx 24rpx 16rpx; border-bottom: 1rpx solid #eee; margin: 0 24rpx; }
.category-row { display: flex; flex-wrap: wrap; padding: 20rpx 24rpx; gap: 32rpx 40rpx; }
.category-tag { display: flex; align-items: center; gap: 8rpx; font-size: 26rpx; color: #666; padding-bottom: 8rpx; border-bottom: 4rpx solid transparent; }
.category-tag.active { color: #333; font-weight: 700; border-bottom-color: #333; }
.category-icon { width: 28rpx; height: 28rpx; font-size: 24rpx; }
.video-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 4rpx; padding: 0 4rpx; }
.video-card { position: relative; aspect-ratio: 3/4; background: #f5f5f5; overflow: hidden; }
.video-cover { width: 100%; height: 100%; display: block; }
.video-like { position: absolute; left: 8rpx; bottom: 8rpx; display: flex; align-items: center; gap: 4rpx; }
.like-icon { font-size: 22rpx; color: #fff; }
.like-count { font-size: 20rpx; color: #fff; }
.load-more { text-align: center; padding: 20rpx; color: #999; font-size: 24rpx; }
.empty { text-align: center; padding: 80rpx; color: #999; }
</style>
