import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { get, post } from '@/utils/request'

export const useVideoStore = defineStore('video', () => {
  const videos = ref([])        // 当前视频列表
  const categories = ref([])    // 分类列表
  const currentPage = ref(1)
  const totalCount = ref(0)
  const pageSize = 10
  const loading = ref(false)
  const hasMore = computed(() => videos.value.length < totalCount.value)

  // 加载视频列表（首页双列流）
  async function fetchVideos(categoryId = 0, entityId = 0, page = 1) {
    if (loading.value) return
    loading.value = true
    try {
      const res = await get('/videos', {
        category_id: categoryId,
        entity_id: entityId,
        page,
        page_size: pageSize,
      })
      const list = res.data || []
      if (page === 1) {
        videos.value = list
      } else {
        videos.value = [...videos.value, ...list]
      }
      currentPage.value = page
      totalCount.value = res.meta?.total || 0
    } catch (err) {
      console.error('加载视频失败', err)
    } finally {
      loading.value = false
    }
  }

  // 加载下一页（播放器预加载用）
  async function fetchNextPage() {
    if (!hasMore.value || loading.value) return
    await fetchVideos(0, 0, currentPage.value + 1)
  }

  // 加载分类
  async function fetchCategories(entityId, mpAccountId) {
    try {
      const res = await get('/categories', { entity_id: entityId, mp_account_id: mpAccountId })
      categories.value = res.data || []
    } catch { /* 分类非必须 */ }
  }

  // ========== 互动操作 ==========

  // 切换点赞
  async function toggleLike(targetType, targetId) {
    const res = await post('/like', { target_type: targetType, target_id: targetId })
    return res.data.liked
  }

  // 切换收藏
  async function toggleFavorite(videoId) {
    const res = await post('/favorite', { video_id: videoId })
    return res.data.favored
  }

  // 记录分享
  async function recordShare(videoId, shareType = 'wechat_friend') {
    await post('/share', { video_id: videoId, share_type: shareType })
  }

  // 获取互动状态
  async function getInteractionStatus(videoId) {
    const res = await get('/interaction-status', { video_id: videoId })
    return res.data
  }

  // ========== 评论 ==========

  async function fetchComments(videoId, page = 1) {
    const res = await get('/comments', { video_id: videoId, page, page_size: 20 })
    return res
  }

  async function postComment(videoId, content, parentId = null, replyToUserId = null) {
    const data = { video_id: videoId, content }
    if (parentId) data.parent_id = parentId
    if (replyToUserId) data.reply_to_user_id = replyToUserId
    return await post('/comments', data)
  }

  return {
    videos, categories, currentPage, totalCount, pageSize, loading, hasMore,
    fetchVideos, fetchNextPage, fetchCategories,
    toggleLike, toggleFavorite, recordShare, getInteractionStatus,
    fetchComments, postComment,
  }
})
