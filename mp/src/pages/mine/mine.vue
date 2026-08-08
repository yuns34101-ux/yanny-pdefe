<template>
  <view class="mine-page">
    <!-- 用户信息（可编辑） -->
    <view class="user-card" v-if="userStore.isLoggedIn">
      <button class="avatar-btn" open-type="chooseAvatar" @chooseavatar="onChooseAvatar">
        <image :src="avatarUrl || '/static/avatar.png'" class="avatar" mode="aspectFill" />
      </button>
      <view class="user-info">
        <input
          class="nickname-input"
          type="nickname"
          v-model="nickname"
          placeholder="点击设置昵称"
          @blur="saveProfile"
        />
        <text class="phone" @click="goBindPhone">{{ phone || '未绑定手机号 ›' }}</text>
      </view>
    </view>
    <view class="user-card" v-else @click="goLogin">
      <image src="/static/avatar.png" class="avatar" mode="aspectFill" />
      <view class="user-info">
        <text class="nickname">点击登录</text>
        <text class="phone">登录后享受更多功能</text>
      </view>
    </view>

    <!-- 功能列表 -->
    <view class="menu-list">
      <view class="menu-item" @click="goFavorites">
        <text>⭐ 我的收藏</text>
        <text class="arrow">›</text>
      </view>
      <view class="menu-item" @click="goHistory">
        <text>🕐 观看历史</text>
        <text class="arrow">›</text>
      </view>
      <view class="menu-item" @click="goSettings">
        <text>⚙️ 设置</text>
        <text class="arrow">›</text>
      </view>
    </view>

    <!-- 退出 -->
    <view class="logout-btn" v-if="userStore.isLoggedIn" @click="handleLogout">
      <text>退出登录</text>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useUserStore } from '@/store/user'
import { get, put, uploadFile } from '@/utils/request'

const userStore = useUserStore()

const nickname = ref('')
const avatarUrl = ref('')
const phone = ref('')
const saving = ref(false)

// 拉取用户最新信息
async function loadUserInfo() {
  if (!userStore.isLoggedIn) return
  try {
    const res = await get('/user/me')
    if (res.data) {
      nickname.value = res.data.nickname || ''
      avatarUrl.value = res.data.avatar_url || ''
      phone.value = res.data.phone || ''
    }
  } catch (err) {
    console.error('获取用户信息失败', err)
  }
}

// 选择头像 → 直传七牛 → 保存
async function onChooseAvatar(e) {
  const filePath = e.detail.avatarUrl
  if (!filePath) return
  avatarUrl.value = filePath // 先本地预览
  await uploadAndSave(filePath)
}

// 昵称失焦时保存
async function saveProfile() {
  if (!nickname.value.trim()) return
  await uploadAndSave(null)
}

// 上传头像（有则传）并保存昵称/头像到服务端
async function uploadAndSave(avatarFilePath) {
  if (saving.value) return
  saving.value = true
  try {
    let finalAvatarUrl = avatarUrl.value
    // 如果是本地路径，需要上传到七牛
    if (avatarFilePath && !avatarFilePath.startsWith('http')) {
      finalAvatarUrl = await uploadAvatar(avatarFilePath)
      avatarUrl.value = finalAvatarUrl
    }

    const res = await put('/user/info', {
      nickname: nickname.value.trim(),
      avatar_url: finalAvatarUrl,
    })

    nickname.value = res.data.nickname || nickname.value
    avatarUrl.value = res.data.avatar_url || avatarUrl.value

    // 同步回 store
    if (userStore.userInfo) {
      userStore.userInfo.nickname = nickname.value
      userStore.userInfo.avatar_url = avatarUrl.value
    }

    uni.showToast({ title: '已保存', icon: 'success' })
  } catch (err) {
    console.error('保存用户信息失败', err)
    uni.showToast({ title: '保存失败，请重试', icon: 'none' })
  } finally {
    saving.value = false
  }
}

// 通过后端代理上传头像到七牛（小程序不直传外部域名）
async function uploadAvatar(filePath) {
  const result = await uploadFile('/upload/avatar', filePath)
  return result.url
}

onShow(() => {
  loadUserInfo()
})

function goLogin() { uni.navigateTo({ url: '/pages/login/login' }) }
function goBindPhone() {
  if (!phone.value) uni.navigateTo({ url: '/pages/login/login' })
}
function goFavorites() {
  if (!userStore.isLoggedIn) return goLogin()
  uni.navigateTo({ url: '/pages/favorites/favorites' })
}
function goHistory() {
  if (!userStore.isLoggedIn) return goLogin()
  uni.navigateTo({ url: '/pages/history/history' })
}
function goSettings() { uni.showToast({ title: '功能开发中', icon: 'none' }) }
function handleLogout() {
  uni.showModal({
    title: '确认退出',
    success: (res) => {
      if (res.confirm) {
        userStore.logout()
        nickname.value = ''
        avatarUrl.value = ''
        phone.value = ''
        uni.showToast({ title: '已退出', icon: 'success' })
      }
    },
  })
}
</script>

<style scoped>
.mine-page { min-height: 100vh; background: #f5f5f5; }
.user-card { display: flex; align-items: center; padding: 40rpx 30rpx; background: #fff; margin-bottom: 20rpx; }
.avatar-btn {
  width: 100rpx; height: 100rpx; border-radius: 50rpx; margin-right: 20rpx;
  padding: 0; border: none; background: none; line-height: normal; flex-shrink: 0;
}
.avatar-btn::after { border: none; }
.avatar { width: 100rpx; height: 100rpx; border-radius: 50rpx; background: #f5f5f5; display: block; }
.user-info { flex: 1; min-width: 0; }
.nickname-input {
  font-size: 32rpx; color: #333; font-weight: 500; height: 48rpx;
  padding: 0; background: none; width: 100%;
}
.nickname { font-size: 32rpx; color: #333; font-weight: 500; }
.phone { font-size: 24rpx; color: #999; margin-top: 6rpx; display: block; }
.menu-list { background: #fff; }
.menu-item { display: flex; justify-content: space-between; padding: 28rpx 30rpx; border-bottom: 1rpx solid #f5f5f5; font-size: 28rpx; }
.arrow { color: #ccc; font-size: 36rpx; }
.logout-btn { margin: 40rpx 30rpx; padding: 24rpx; text-align: center; background: #fff; border-radius: 12rpx; color: #ff4d4f; font-size: 28rpx; }
</style>
