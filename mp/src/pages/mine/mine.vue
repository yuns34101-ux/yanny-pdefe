<template>
  <view class="mine-page">
    <!-- 用户信息 -->
    <view class="user-card" v-if="userStore.isLoggedIn">
      <image :src="userStore.userInfo?.avatar_url || '/static/avatar.png'" class="avatar" />
      <view class="user-info">
        <text class="nickname">{{ userStore.userInfo?.nickname || '未设置' }}</text>
        <text class="phone">{{ userStore.userInfo?.phone || '未绑定手机号' }}</text>
      </view>
    </view>
    <view class="user-card" v-else @click="goLogin">
      <image src="/static/avatar.png" class="avatar" />
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
import { useUserStore } from '@/store/user'

const userStore = useUserStore()

function goLogin() { uni.navigateTo({ url: '/pages/login/login' }) }
function goFavorites() {
  if (!userStore.isLoggedIn) return goLogin()
  uni.showToast({ title: '功能开发中', icon: 'none' })
}
function goHistory() { uni.showToast({ title: '功能开发中', icon: 'none' }) }
function goSettings() { uni.showToast({ title: '功能开发中', icon: 'none' }) }
function handleLogout() {
  uni.showModal({
    title: '确认退出',
    success: (res) => {
      if (res.confirm) {
        userStore.logout()
        uni.showToast({ title: '已退出', icon: 'success' })
      }
    },
  })
}
</script>

<style scoped>
.mine-page { min-height: 100vh; background: #f5f5f5; }
.user-card { display: flex; align-items: center; padding: 40rpx 30rpx; background: #fff; margin-bottom: 20rpx; }
.avatar { width: 100rpx; height: 100rpx; border-radius: 50rpx; margin-right: 20rpx; }
.nickname { font-size: 32rpx; color: #333; font-weight: 500; }
.phone { font-size: 24rpx; color: #999; margin-top: 6rpx; display: block; }
.menu-list { background: #fff; }
.menu-item { display: flex; justify-content: space-between; padding: 28rpx 30rpx; border-bottom: 1rpx solid #f5f5f5; font-size: 28rpx; }
.arrow { color: #ccc; font-size: 36rpx; }
.logout-btn { margin: 40rpx 30rpx; padding: 24rpx; text-align: center; background: #fff; border-radius: 12rpx; color: #ff4d4f; font-size: 28rpx; }
</style>
