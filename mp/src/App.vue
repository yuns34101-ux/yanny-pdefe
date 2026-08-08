<script setup>
import { onLaunch } from '@dcloudio/uni-app'
import { useUserStore } from '@/store/user'

// 接收分享裂变的邀请人参数（微信小程序 onLaunch 的 query 里带 scene/path 参数）
function capturePendingInviter(options) {
  const inviter = options?.query?.inviter
  if (inviter) uni.setStorageSync('pending_inviter', inviter)
}

async function loginWithRetry(userStore) {
  try {
    await userStore.silentLogin()
  } catch (e) {
    uni.showModal({
      title: '登录失败',
      content: '进入小程序需要先登录，请重试',
      showCancel: false,
      confirmText: '重试',
      success: () => loginWithRetry(userStore),
    })
  }
}

onLaunch(async (options) => {
  const userStore = useUserStore()
  capturePendingInviter(options)

  // 进场即静默登录 —— 登录后所有接口才可访问，登录失败强制重试，不允许降级浏览
  const token = uni.getStorageSync('mp_token')
  if (!token) {
    await loginWithRetry(userStore)
  } else {
    userStore.token = token
    userStore.ready = true
    // 从 storage 恢复用户信息，避免重启后 userInfo 为空导致"我的"页显示未绑定
    const cached = uni.getStorageSync('mp_user_info')
    if (cached) {
      try { userStore.userInfo = JSON.parse(cached) } catch { /* ignore */ }
    }
  }
})
</script>

<style>
page { background-color: #f5f5f5; }
</style>
