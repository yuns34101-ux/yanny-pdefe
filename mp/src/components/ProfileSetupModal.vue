<template>
  <view class="profile-mask" v-if="userStore.showProfileSetup">
    <view class="profile-panel" @click.stop>
      <text class="profile-title">完善你的资料</text>
      <text class="profile-sub">头像昵称将展示在评论区，方便大家认识你</text>

      <button class="avatar-picker" open-type="chooseAvatar" @chooseavatar="onChooseAvatar">
        <image class="avatar-preview" :src="avatarUrl || '/static/avatar.png'" mode="aspectFill" />
        <text class="avatar-hint">点击选择头像</text>
      </button>

      <input
        class="nickname-input"
        type="nickname"
        v-model="nickname"
        placeholder="请输入昵称"
      />

      <button
        v-if="userStore.platform === 'wechat'"
        class="phone-btn"
        :class="{ bound: phoneBound }"
        open-type="getPhoneNumber"
        @getphonenumber="onGetPhoneNumber"
      >
        <text>{{ phoneBound ? '✓ 手机号已绑定' : '一键绑定手机号（可选）' }}</text>
      </button>

      <view class="profile-actions">
        <text class="skip-btn" @click="handleSkip">暂不设置</text>
        <text class="confirm-btn" :class="{ disabled: submitting }" @click="handleConfirm">
          {{ submitting ? '提交中...' : '确定' }}
        </text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { useUserStore } from '@/store/user'
import { post, uploadFile } from '@/utils/request'

const userStore = useUserStore()
const nickname = ref('')
const avatarUrl = ref('')
const avatarTempFile = ref('')
const submitting = ref(false)
const phoneBound = ref(false)

function onChooseAvatar(e) {
  avatarTempFile.value = e.detail.avatarUrl
  avatarUrl.value = e.detail.avatarUrl
}

async function onGetPhoneNumber(e) {
  if (!e.detail?.encryptedData) return
  try {
    await userStore.updatePhone(e)
    phoneBound.value = true
  } catch { }
}

// 通过后端代理上传（小程序不直传七牛）
async function uploadAvatar(filePath) {
  const result = await uploadFile('/upload/avatar', filePath)
  return result.url
}

async function handleConfirm() {
  if (submitting.value) return
  if (!nickname.value.trim()) {
    uni.showToast({ title: '请输入昵称', icon: 'none' })
    return
  }
  submitting.value = true
  try {
    let finalAvatarUrl = ''
    if (avatarTempFile.value) {
      finalAvatarUrl = await uploadAvatar(avatarTempFile.value)
    }
    await userStore.updateProfile(nickname.value.trim(), finalAvatarUrl)
    uni.showToast({ title: '设置成功', icon: 'success' })
  } catch (err) {
    uni.showToast({ title: err.message || '设置失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}

function handleSkip() {
  userStore.showProfileSetup = false
}
</script>

<style scoped>
.profile-mask {
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.5); z-index: 100;
  display: flex; align-items: center; justify-content: center;
}
.profile-panel {
  width: 600rpx; background: #fff; border-radius: 24rpx;
  padding: 48rpx 40rpx; display: flex; flex-direction: column; align-items: center;
}
.profile-title { font-size: 32rpx; font-weight: 700; color: #333; }
.profile-sub { font-size: 22rpx; color: #999; margin-top: 12rpx; text-align: center; }
.avatar-picker {
  margin-top: 40rpx; background: none; border: none; padding: 0; line-height: normal;
  display: flex; flex-direction: column; align-items: center;
}
.avatar-picker::after { border: none; }
.avatar-preview { width: 140rpx; height: 140rpx; border-radius: 70rpx; background: #f5f5f5; }
.avatar-hint { font-size: 22rpx; color: #999; margin-top: 12rpx; }
.nickname-input {
  width: 100%; height: 80rpx; background: #f5f5f5; border-radius: 12rpx;
  margin-top: 32rpx; padding: 0 24rpx; font-size: 28rpx; text-align: center;
}
.phone-btn {
  width: 100%; height: 76rpx; margin-top: 24rpx; padding: 0; line-height: 76rpx;
  background: #f0f9eb; color: #07c160; border-radius: 12rpx; border: none; font-size: 26rpx;
}
.phone-btn::after { border: none; }
.phone-btn.bound { background: #f5f5f5; color: #999; }
.profile-actions { display: flex; width: 100%; margin-top: 40rpx; gap: 24rpx; }
.skip-btn, .confirm-btn {
  flex: 1; text-align: center; padding: 20rpx 0; border-radius: 40rpx; font-size: 28rpx;
}
.skip-btn { background: #f2f2f2; color: #666; }
.confirm-btn { background: #409EFF; color: #fff; }
.confirm-btn.disabled { opacity: 0.6; }
</style>
