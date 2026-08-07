<template>
  <view class="player-page">
    <!-- 全屏视频滑动 -->
    <swiper
      class="video-swiper"
      :vertical="true"
      :current="currentIndex"
      :duration="300"
      @change="onSwipeChange"
      @transition="onSwipeTransition"
    >
      <swiper-item v-for="(v, i) in videoList" :key="v.id">
        <view class="video-wrapper">
          <!-- 视频播放器 -->
          <video
            v-if="Math.abs(i - currentIndex) <= 1"
            :id="'video-' + v.id"
            :src="v.video_url"
            :poster="v.cover_url"
            :autoplay="i === currentIndex"
            :loop="false"
            :controls="false"
            :show-center-play-btn="false"
            :show-play-btn="false"
            :enable-progress-gesture="false"
            object-fit="contain"
            class="video-player"
            @play="onPlay(v, i)"
            @ended="onEnded(v)"
            @error="onError(v)"
            @timeupdate="onTimeUpdate(v, $event)"
            @loadedmetadata="onLoadedMetadata(v, $event)"
          />
          <!-- 占位封面（未播放时） -->
          <image
            v-else
            :src="v.cover_url"
            mode="aspectFill"
            class="video-cover"
          />

          <!-- 底部信息区 -->
          <view class="bottom-bar">
            <!-- 进度条 -->
            <view class="progress-row">
              <text class="progress-time">{{ formatDuration(v._currentTime) }}</text>
              <text class="progress-speed" @click.stop="toggleSpeed(v)">倍数</text>
            </view>
            <view class="progress-track">
              <view class="progress-fill" :style="{ width: (v._progress || 0) + '%' }" />
              <view class="progress-dot" :style="{ left: (v._progress || 0) + '%' }" />
            </view>

            <view class="video-desc">
              <text class="video-title-text">{{ v.title }}</text>
            </view>
            <view class="video-desc" v-if="v.description">
              <text class="video-desc-text">{{ v.description }}</text>
            </view>

            <view class="bottom-row">
              <view class="bottom-entity" v-if="entityInfo">
                <image :src="entityInfo.logo_url || '/static/logo.png'" class="entity-logo" />
                <text class="entity-name">{{ entityInfo.name }}</text>
                <view class="follow-pill" :class="{ followed: entityStore.followed }" @click.stop="handleFollow">
                  <text>{{ entityStore.followed ? '已关注' : '+ 关注' }}</text>
                </view>
              </view>
              <view class="bottom-actions">
                <view class="action-btn" @click.stop="handleLike(v)">
                  <text class="action-icon" :style="{ color: v._liked ? '#409EFF' : '#fff' }">👍</text>
                  <text class="action-count">{{ formatCount(v.like_count) }}</text>
                </view>
                <view class="action-btn" @click.stop="handleShare(v)">
                  <text class="action-icon">↗️</text>
                  <text class="action-count">{{ formatCount(v.share_count) }}</text>
                </view>
                <view class="action-btn" @click.stop="handleFavorite(v)">
                  <text class="action-icon" :style="{ color: v._favored ? '#ff4d4f' : '#fff' }">
                    {{ v._favored ? '❤️' : '🤍' }}
                  </text>
                  <text class="action-count">{{ formatCount(v.collect_count) }}</text>
                </view>
                <view class="action-btn" @click.stop="handleComment(v)">
                  <text class="action-icon">💬</text>
                  <text class="action-count">{{ formatCount(v.comment_count) }}</text>
                </view>
              </view>
            </view>
          </view>
        </view>
      </swiper-item>
    </swiper>

    <!-- 返回按钮 -->
    <view class="back-btn" @click="goBack">
      <text class="back-icon">←</text>
    </view>

    <!-- 评论面板 -->
    <view class="comment-panel" v-if="showComment" @click.stop>
      <view class="comment-header">
        <text class="comment-title">评论 ({{ commentTotal }})</text>
        <text class="comment-close" @click="showComment = false">✕</text>
      </view>
      <scroll-view scroll-y class="comment-list">
        <view v-for="c in comments" :key="c.id" class="comment-item">
          <image :src="c.user?.avatar_url || '/static/avatar.png'" class="comment-avatar" />
          <view class="comment-body">
            <view class="comment-user">
              <text class="comment-nickname">{{ c.user?.nickname || '匿名用户' }}</text>
              <text class="comment-time">{{ formatTime(c.created_at) }}</text>
            </view>
            <text class="comment-content">{{ c.content }}</text>
            <!-- 回复 -->
            <view v-if="c.replies?.length" class="comment-replies">
              <view v-for="r in c.replies" :key="r.id" class="reply-item">
                <text class="reply-user">{{ r.user?.nickname || '匿名' }}：</text>
                <text>{{ r.content }}</text>
              </view>
            </view>
            <text class="reply-btn" @click="onReply(c)">回复</text>
          </view>
        </view>
      </scroll-view>
      <!-- 评论输入框 -->
      <view class="comment-input-bar">
        <input
          v-model="commentText"
          class="comment-input"
          :placeholder="replyTarget ? '回复 @' + replyTarget.user?.nickname : '写评论...'"
          @confirm="submitComment"
        />
        <text class="send-btn" @click="submitComment">发送</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useVideoStore } from '@/store/video'
import { useUserStore } from '@/store/user'
import { useEntityStore } from '@/store/entity'
import { createTrackPayload } from '@/utils/sign'
import { get, post } from '@/utils/request'

const videoStore = useVideoStore()
const userStore = useUserStore()
const entityStore = useEntityStore()

// 视频数据
const videoList = ref([])
const currentIndex = ref(0)
const entityInfo = computed(() => entityStore.entity)
const page = ref(1)
const totalCount = ref(0)
const pageSize = 10
const loadingMore = ref(false)
const PRELOAD_THRESHOLD = 3 // 距末尾 3 个时预加载

// 评论
const showComment = ref(false)
const comments = ref([])
const commentTotal = ref(0)
const commentText = ref('')
const replyTarget = ref(null)
const currentVideo = ref(null)

// ========== 视频加载 + 分页预加载 ==========

const hasMore = computed(() => videoList.value.length < totalCount.value)

async function loadVideos(pageNum = 1) {
  if (loadingMore.value && pageNum > 1) return
  loadingMore.value = true
  try {
    const res = await get('/videos', { page: pageNum, page_size: pageSize })
    const list = (res.data || []).map(v => ({ ...v, _liked: false, _favored: false }))
    if (pageNum === 1) {
      videoList.value = list
    } else {
      videoList.value = [...videoList.value, ...list]
    }
    page.value = pageNum
    totalCount.value = res.meta?.total || 0
  } catch (err) {
    console.error('加载视频失败', err)
  } finally {
    loadingMore.value = false
  }
}

// 预加载下一页（滑动到阈值时触发）
async function preloadNextPage() {
  const needPreload = currentIndex.value >= videoList.value.length - PRELOAD_THRESHOLD
  if (needPreload && hasMore.value && !loadingMore.value) {
    await loadVideos(page.value + 1)
  }
}

// 当前视频互动状态
async function loadInteractionStatus(videoId) {
  if (!userStore.isLoggedIn) return
  try {
    const res = await videoStore.getInteractionStatus(videoId)
    const v = videoList.value.find(v => v.id === videoId)
    if (v) {
      v._liked = res.liked
      v._favored = res.favored
    }
    entityStore.followed = !!res.followed
  } catch { /* ignore */ }
}

// 加载当前视频所属主体信息（logo/名称/关注状态）
function loadEntityInfo(v) {
  if (!v?.entity_id) return
  if (entityStore.entity?.id === v.entity_id) return
  entityStore.fetchEntityInfo(v.entity_id)
}

// ========== 滑动事件 ==========

function onSwipeChange(e) {
  const newIndex = e.detail.current
  const oldIndex = currentIndex.value
  currentIndex.value = newIndex

  // 上报上一个视频的观看数据
  if (oldIndex !== newIndex && videoList.value[oldIndex]) {
    reportView(videoList.value[oldIndex])
  }

  // 加载新视频互动状态
  const newVideo = videoList.value[newIndex]
  if (newVideo) {
    currentVideo.value = newVideo
    loadInteractionStatus(newVideo.id)
    loadEntityInfo(newVideo)
  }

  // 预加载
  preloadNextPage()
}

function onSwipeTransition(e) {
  // 过渡中可预加载
  preloadNextPage()
}

// ========== 互动操作（需登录检查） ==========

function requireLogin(action) {
  if (!userStore.isLoggedIn) {
    uni.navigateTo({ url: '/pages/login/login' })
    return false
  }
  return true
}

async function handleLike(v) {
  if (!requireLogin()) return
  const liked = await videoStore.toggleLike('video', v.id)
  v._liked = liked
  v.like_count += liked ? 1 : -1
}

async function handleFavorite(v) {
  if (!requireLogin()) return
  const favored = await videoStore.toggleFavorite(v.id)
  v._favored = favored
  v.collect_count += favored ? 1 : -1
}

async function handleFollow() {
  if (!requireLogin()) return
  await entityStore.toggleFollow()
}

// ========== 进度条 ==========

function onLoadedMetadata(v, e) {
  v._duration = e.detail.duration || 0
}

function onTimeUpdate(v, e) {
  v._currentTime = e.detail.currentTime || 0
  const duration = e.detail.duration || v._duration || 0
  v._progress = duration ? (v._currentTime / duration) * 100 : 0
}

function toggleSpeed(v) {
  const rates = [1, 1.25, 1.5, 2]
  const idx = rates.indexOf(v._playbackRate || 1)
  v._playbackRate = rates[(idx + 1) % rates.length]
  const ctx = uni.createVideoContext('video-' + v.id)
  ctx?.playbackRate(v._playbackRate)
}

function formatDuration(sec) {
  if (!sec) return '00:00'
  const m = Math.floor(sec / 60)
  const s = Math.floor(sec % 60)
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

function handleComment(v) {
  if (!requireLogin()) return
  currentVideo.value = v
  showComment.value = true
  loadComments(v.id)
}

function handleShare(v) {
  // 分享不强制登录，但需记录
  const shareType = 'wechat_friend'
  if (userStore.isLoggedIn) {
    videoStore.recordShare(v.id, shareType)
    v.share_count++
  }
  // 触发微信分享
  uni.share({ provider: 'weixin', type: 0, title: v.title, imageUrl: v.cover_url })
}

// ========== 评论 ==========

async function loadComments(videoId) {
  try {
    const res = await videoStore.fetchComments(videoId)
    comments.value = res.data || []
    commentTotal.value = res.meta?.total || 0
  } catch { /* ignore */ }
}

function onReply(comment) {
  replyTarget.value = comment
  commentText.value = ''
}

async function submitComment() {
  if (!commentText.value.trim()) return
  if (!currentVideo.value) return
  try {
    await videoStore.postComment(
      currentVideo.value.id,
      commentText.value,
      replyTarget.value?.id || null,
      replyTarget.value?.user_id || null,
    )
    commentText.value = ''
    replyTarget.value = null
    uni.showToast({ title: '评论成功', icon: 'success' })
    loadComments(currentVideo.value.id)
  } catch (err) {
    uni.showToast({ title: err.message || '评论失败', icon: 'none' })
  }
}

// ========== 埋点上报（签名防篡改） ==========

let viewTimer = null
function onPlay(v, index) {
  // 播放 1 秒后上报
  clearTimeout(viewTimer)
  viewTimer = setTimeout(() => reportView(v), 1000)
}

function reportView(v) {
  const payload = createTrackPayload({
    video_id: v.id,
    watch_duration: 0,
    is_complete: 0,
    source: 'swipe',
  })
  post('/track/view', payload).catch(() => {})
}

function onEnded(v) {
  const payload = createTrackPayload({
    video_id: v.id,
    watch_duration: v.duration || 0,
    is_complete: 1,
    source: 'swipe',
  })
  post('/track/view', payload).catch(() => {})
}

function onError(v) {
  console.error('视频播放失败', v.id)
}

// ========== 工具 ==========

function formatCount(n) {
  if (!n) return '0'
  if (n >= 10000) return (n / 10000).toFixed(1) + 'w'
  return String(n)
}

function formatTime(t) {
  if (!t) return ''
  const d = new Date(t)
  return `${d.getMonth() + 1}-${d.getDate()} ${d.getHours()}:${String(d.getMinutes()).padStart(2, '0')}`
}

function goBack() {
  uni.navigateBack()
}

// ========== 生命周期 ==========

onLoad((query) => {
  // 接收首页传来的初始视频 ID 和分类
  const videoId = query.videoId
  const categoryId = query.categoryId || 0
  loadVideos(1).then(() => {
    if (videoId) {
      const idx = videoList.value.findIndex(v => v.id == videoId)
      if (idx >= 0) {
        currentIndex.value = idx
        currentVideo.value = videoList.value[idx]
        loadInteractionStatus(parseInt(videoId))
        loadEntityInfo(currentVideo.value)
      }
    } else if (videoList.value.length) {
      currentVideo.value = videoList.value[0]
      loadInteractionStatus(currentVideo.value.id)
      loadEntityInfo(currentVideo.value)
    }
  })
})

onUnmounted(() => {
  clearTimeout(viewTimer)
})
</script>

<style scoped>
.player-page { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: #000; }
.video-swiper { width: 100%; height: 100%; }
.video-wrapper { width: 100%; height: 100%; position: relative; }
.video-player { width: 100%; height: 100%; }
.video-cover { width: 100%; height: 100%; }

/* 底部信息 */
.bottom-bar { position: absolute; left: 24rpx; right: 24rpx; bottom: 30rpx; z-index: 10; }

/* 进度条 */
.progress-row { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10rpx; }
.progress-time { font-size: 22rpx; color: #fff; }
.progress-speed { font-size: 22rpx; color: #fff; }
.progress-track { position: relative; height: 4rpx; background: rgba(255,255,255,0.3); border-radius: 2rpx; margin-bottom: 16rpx; }
.progress-fill { position: absolute; left: 0; top: 0; bottom: 0; background: #fff; border-radius: 2rpx; }
.progress-dot { position: absolute; top: 50%; width: 16rpx; height: 16rpx; margin-left: -8rpx; margin-top: -8rpx; border-radius: 50%; background: #fff; }

.video-title-text { font-size: 28rpx; color: #fff; font-weight: 500; }
.video-desc-text { font-size: 24rpx; color: rgba(255,255,255,0.8); margin-top: 6rpx; }

/* 底部：主体信息 + 互动按钮 */
.bottom-row { display: flex; align-items: flex-end; justify-content: space-between; margin-top: 16rpx; }
.bottom-entity { display: flex; align-items: center; }
.entity-logo { width: 56rpx; height: 56rpx; border-radius: 28rpx; margin-right: 12rpx; }
.entity-name { font-size: 26rpx; color: #fff; margin-right: 16rpx; }
.follow-pill { padding: 4rpx 20rpx; border-radius: 24rpx; background: rgba(255,255,255,0.15); }
.follow-pill text { font-size: 22rpx; color: #fff; }
.follow-pill.followed { background: rgba(255,255,255,0.08); }
.follow-pill.followed text { color: rgba(255,255,255,0.6); }
.bottom-actions { display: flex; align-items: center; gap: 32rpx; }
.action-btn { display: flex; flex-direction: column; align-items: center; }
.action-icon { font-size: 44rpx; }
.action-count { font-size: 20rpx; color: #fff; margin-top: 4rpx; }

/* 返回 */
.back-btn { position: absolute; top: 60rpx; left: 20rpx; z-index: 20; width: 60rpx; height: 60rpx; border-radius: 30rpx; background: rgba(0,0,0,0.3); display: flex; align-items: center; justify-content: center; }
.back-icon { color: #fff; font-size: 36rpx; }

/* 评论面板 */
.comment-panel { position: absolute; bottom: 0; left: 0; right: 0; height: 60vh; background: #fff; border-radius: 24rpx 24rpx 0 0; z-index: 30; display: flex; flex-direction: column; }
.comment-header { display: flex; justify-content: space-between; align-items: center; padding: 20rpx 30rpx; border-bottom: 1rpx solid #eee; }
.comment-title { font-size: 30rpx; font-weight: 600; }
.comment-close { font-size: 32rpx; color: #999; padding: 10rpx; }
.comment-list { flex: 1; padding: 0 30rpx; overflow-y: auto; }
.comment-item { display: flex; padding: 20rpx 0; border-bottom: 1rpx solid #f5f5f5; }
.comment-avatar { width: 64rpx; height: 64rpx; border-radius: 32rpx; margin-right: 16rpx; flex-shrink: 0; }
.comment-body { flex: 1; }
.comment-user { display: flex; justify-content: space-between; }
.comment-nickname { font-size: 26rpx; color: #666; }
.comment-time { font-size: 22rpx; color: #ccc; }
.comment-content { font-size: 28rpx; color: #333; margin-top: 8rpx; display: block; }
.comment-replies { background: #f5f5f5; border-radius: 8rpx; padding: 12rpx; margin-top: 10rpx; }
.reply-item { font-size: 26rpx; color: #333; margin-bottom: 4rpx; }
.reply-user { color: #409EFF; }
.reply-btn { font-size: 22rpx; color: #999; margin-top: 8rpx; display: inline-block; }
.comment-input-bar { display: flex; align-items: center; padding: 16rpx 30rpx; border-top: 1rpx solid #eee; background: #fff; }
.comment-input { flex: 1; height: 64rpx; background: #f5f5f5; border-radius: 32rpx; padding: 0 24rpx; font-size: 26rpx; }
.send-btn { color: #409EFF; font-size: 28rpx; margin-left: 20rpx; }
</style>
