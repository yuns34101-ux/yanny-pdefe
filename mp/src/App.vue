<script setup>
import { onLaunch } from '@dcloudio/uni-app'
import { useUserStore } from '@/store/user'

onLaunch(async () => {
  console.log('Yanny 启动')
  const userStore = useUserStore()

  // 进场即静默登录 —— 获取 Token 用于后续埋点链路
  const token = uni.getStorageSync('mp_token')
  if (!token) {
    try {
      await userStore.silentLogin()
      console.log('静默登录成功')
    } catch (e) {
      console.log('静默登录失败，稍后重试', e.message)
    }
  } else {
    userStore.token = token
  }
})
</script>

<style>
page { background-color: #f5f5f5; }
</style>
