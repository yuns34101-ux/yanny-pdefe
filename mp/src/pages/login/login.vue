<template>
  <view class="login-page">
    <!-- 广告位 / 品牌图预留 -->
    <view class="ad-slot">
      <image
        v-if="brandImage"
        :src="brandImage"
        mode="aspectFill"
        class="brand-img"
      />
      <view v-else class="brand-placeholder">
        <text class="brand-emoji">🎬</text>
        <text class="brand-text">Yanny</text>
        <text class="brand-slogan">发现精彩短视频</text>
      </view>
    </view>

    <!-- 登录卡片 -->
    <view class="login-card">
      <text class="login-title">绑定手机号，体验更多功能</text>

      <!-- 微信手机号绑定 -->
      <button
        v-if="platform === 'wechat'"
        class="platform-btn wx-btn"
        open-type="getPhoneNumber"
        @getphonenumber="handleWxPhoneLogin"
        :loading="loading"
      >
        <text class="btn-icon">💬</text>
        <text>微信一键绑定手机号</text>
      </button>

      <!-- 抖音登录 -->
      <button
        v-if="platform === 'douyin'"
        class="platform-btn dy-btn"
        @click="handleDouyinLogin"
        :loading="loading"
      >
        <text class="btn-icon">🎵</text>
        <text>抖音授权登录</text>
      </button>

      <!-- H5 / Web 登录 -->
      <view v-if="platform === 'h5' || platform === 'mock'" class="h5-login">
        <view class="phone-input-row">
          <text class="phone-prefix">+86</text>
          <input
            v-model="phone"
            class="phone-input"
            type="number"
            maxlength="11"
            placeholder="请输入手机号"
          />
        </view>
        <button class="platform-btn h5-btn" @click="handleH5Login" :loading="loading">
          <text>手机号登录 / 注册</text>
        </button>
      </view>

      <!-- 暂不登录 -->
      <view class="skip-login" @click="goBack">
        <text>暂不登录，先看看</text>
      </view>

      <!-- 协议 -->
      <view class="agreement">
        <text>登录即同意《用户协议》和《隐私政策》</text>
      </view>
    </view>

    <!-- 底部广告位 -->
    <view class="bottom-ad">
      <text class="ad-label">— 广告 —</text>
      <view class="ad-box">
        <text class="ad-placeholder">广告位预留</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()
const loading = ref(false)
const phone = ref('')
const brandImage = ref('') // 可配置的品牌宣传图

const platform = computed(() => userStore.platform)

// 微信手机号绑定（登录早已由 App.vue 静默完成，这里只做手机号授权+绑定）
async function handleWxPhoneLogin(e) {
  if (!e.detail?.encryptedData) {
    uni.showToast({ title: '未授权手机号', icon: 'none' })
    return
  }
  loading.value = true
  try {
    await userStore.updatePhone(e)
    setTimeout(() => uni.navigateBack(), 1000)
  } catch (err) {
    // updatePhone 内部已 toast 失败提示
  } finally {
    loading.value = false
  }
}

// 抖音登录
async function handleDouyinLogin() {
  loading.value = true
  try {
    // #ifdef MP-DOUYIN
    const res = await tt.login()
    await userStore.doLogin(res.code, 'douyin')
    // #endif
    uni.showToast({ title: '登录成功', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 1000)
  } catch (err) {
    uni.showToast({ title: '登录失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

// H5 手机号登录
async function handleH5Login() {
  if (!phone.value || phone.value.length < 11) {
    uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
    return
  }
  loading.value = true
  try {
    await userStore.doLogin(phone.value, 'h5')
    uni.showToast({ title: '登录成功', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 1000)
  } catch (err) {
    uni.showToast({ title: '登录失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

function goBack() {
  uni.navigateBack()
}
</script>

<style scoped>
.login-page {
  min-height: 100vh; background: #fff;
  display: flex; flex-direction: column;
}
/* 广告位/品牌区 */
.ad-slot {
  height: 480rpx; position: relative; overflow: hidden;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.brand-img { width: 100%; height: 100%; }
.brand-placeholder {
  display: flex; flex-direction: column; align-items: center;
  justify-content: center; height: 100%;
}
.brand-emoji { font-size: 80rpx; margin-bottom: 16rpx; }
.brand-text { font-size: 40rpx; color: #fff; font-weight: 700; }
.brand-slogan { font-size: 24rpx; color: rgba(255,255,255,0.8); margin-top: 8rpx; }
/* 登录卡片 */
.login-card {
  padding: 40rpx 48rpx; flex: 1;
  display: flex; flex-direction: column; align-items: center;
}
.login-title { font-size: 32rpx; color: #333; font-weight: 600; margin-bottom: 40rpx; }
.platform-btn {
  width: 100%; height: 96rpx; border-radius: 48rpx;
  display: flex; align-items: center; justify-content: center;
  font-size: 30rpx; border: none; margin-bottom: 24rpx; color: #fff;
}
.wx-btn { background: #07c160; }
.dy-btn { background: #111; }
.h5-btn { background: linear-gradient(135deg, #409EFF, #337ecc); }
.btn-icon { margin-right: 12rpx; font-size: 36rpx; }
/* H5 登录 */
.h5-login { width: 100%; }
.phone-input-row {
  display: flex; align-items: center; border: 2rpx solid #e0e0e0;
  border-radius: 48rpx; padding: 0 24rpx; height: 96rpx; margin-bottom: 24rpx;
}
.phone-prefix { font-size: 30rpx; color: #333; padding-right: 16rpx; border-right: 2rpx solid #eee; }
.phone-input { flex: 1; margin-left: 16rpx; font-size: 30rpx; }
.skip-login { padding: 20rpx; }
.skip-login text { color: #999; font-size: 26rpx; }
.agreement { margin-top: 40rpx; text-align: center; }
.agreement text { font-size: 22rpx; color: #ccc; }
/* 底部广告 */
.bottom-ad { padding: 24rpx 48rpx 40rpx; text-align: center; }
.ad-label { font-size: 22rpx; color: #ccc; }
.ad-box {
  margin-top: 12rpx; height: 120rpx; background: #f5f5f5;
  border-radius: 12rpx; display: flex; align-items: center; justify-content: center;
}
.ad-placeholder { color: #ccc; font-size: 24rpx; }
</style>
