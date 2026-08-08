<script setup>
import { onLaunch } from '@dcloudio/uni-app'
import { useUserStore } from '@/store/user'
import { get } from '@/utils/request'

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

// 从服务端拉取当前用户完整信息（昵称/头像/手机号），解决 token 存在但 userInfo 丢失的问题
async function restoreUserInfo(userStore) {
  try {
    const res = await get('/user/me')
    if (res.data) {
      userStore.userInfo = res.data
    }
  } catch {
    // 静默失败，login 时 doLogin 会兜底设置
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
    // 通过 API 拉取用户信息，填补重启后 userInfo 为空的缺口
    restoreUserInfo(userStore)
  }
})
</script>

<style>
page { background-color: #f5f5f5; }
</style>
