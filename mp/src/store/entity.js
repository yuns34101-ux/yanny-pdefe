import { defineStore } from 'pinia'
import { ref } from 'vue'
import { get, post } from '@/utils/request'

export const useEntityStore = defineStore('entity', () => {
  const entity = ref(null)
  const videoCount = ref(0)
  const followerCount = ref(0)
  const followed = ref(false)
  const loading = ref(false)

  // 获取主体详情 + 实时统计（video_count / follower_count / followed）
  async function fetchEntityInfo(entityId = 0) {
    loading.value = true
    try {
      const res = await get('/entity', { entity_id: entityId })
      entity.value = res.data
      videoCount.value = res.data?.video_count || 0
      followerCount.value = res.data?.follower_count || 0
      followed.value = !!res.data?.followed
    } catch (err) {
      console.error('加载主体信息失败', err)
    } finally {
      loading.value = false
    }
  }

  // 切换关注（乐观更新计数）
  async function toggleFollow() {
    if (!entity.value?.id) return
    const res = await post('/follow', { entity_id: entity.value.id })
    followed.value = res.data.followed
    followerCount.value += followed.value ? 1 : -1
    return followed.value
  }

  return {
    entity, videoCount, followerCount, followed, loading,
    fetchEntityInfo, toggleFollow,
  }
})
