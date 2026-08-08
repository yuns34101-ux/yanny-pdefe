import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { post, put } from '@/utils/request'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref(null)
  const token = ref(uni.getStorageSync('mp_token') || '')
  const platform = ref('wechat') // wechat | douyin | h5 | mock
  const ready = ref(!!token.value) // 登录就绪状态：token 已存在或静默登录已完成
  const loginError = ref('')
  const showProfileSetup = ref(false) // 首次静默登录成功后是否弹出头像/昵称采集

  const isLoggedIn = computed(() => !!token.value)
  const userId = computed(() => userInfo.value?.user_id || 0)

  // 检测当前平台
  function detectPlatform() {
    // #ifdef MP-WEIXIN
    return 'wechat'
    // #endif
    // #ifdef MP-DOUYIN
    return 'douyin'
    // #endif
    // #ifdef H5
    return 'h5'
    // #endif
    return 'mock'
  }

  // 静默登录（进场时调用，无 UI）
  async function silentLogin() {
    ready.value = false
    loginError.value = ''
    try {
      const p = detectPlatform()
      platform.value = p

      if (p === 'wechat') {
        // 微信小程序：wx.login 获取 code
        await new Promise((resolve, reject) => {
          uni.login({
            provider: 'weixin',
            success: async (res) => {
              try {
                await doLogin(res.code, 'wechat')
                resolve()
              } catch (e) { reject(e) }
            },
            fail: (e) => reject(e),
          })
        })
      } else {
        // H5 / mock / 其他：直接使用设备标识
        const deviceId = uni.getStorageSync('device_id') || `dev_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
        uni.setStorageSync('device_id', deviceId)
        await doLogin(deviceId, p)
      }
      ready.value = true
    } catch (e) {
      loginError.value = e.message || '登录失败'
      throw e
    }
  }

  // 实际登录请求
  async function doLogin(code, platformType) {
    const appId = getAppId()
    const inviter = Number(uni.getStorageSync('pending_inviter')) || 0
    const res = await post('/login', {
      code,
      platform: platformType,
      app_id: appId,
      inviter_user_id: inviter,
    })
    token.value = res.data.token
    userInfo.value = {
      user_id: res.data.user_id,
      nickname: res.data.nickname,
      avatar_url: res.data.avatar_url,
    }
    uni.setStorageSync('mp_token', res.data.token)
    if (inviter) uni.removeStorageSync('pending_inviter')
    // 微信新用户首次登录：昵称/头像还是占位值，弹窗采集真实资料
    if (platformType === 'wechat' && res.data.is_new_user) {
      showProfileSetup.value = true
    }
    return res.data
  }

  // 绑定手机号（已登录状态下调用，用登录时留存的 session_key 解密）
  async function updatePhone(e) {
    if (!e?.detail?.encryptedData) return
    try {
      const res = await post('/user/phone', {
        encrypted_data: e.detail.encryptedData,
        iv: e.detail.iv,
      })
      if (userInfo.value) userInfo.value.phone = res.data.phone
      uni.showToast({ title: '手机号绑定成功', icon: 'success' })
    } catch (e) {
      console.error('手机号绑定失败', e)
      uni.showToast({ title: '手机号绑定失败', icon: 'none' })
      throw e
    }
  }

  // 提交头像/昵称（chooseAvatar + nickname 采集弹窗确认后调用）
  async function updateProfile(nickname, avatarURL) {
    const res = await put('/user/info', { nickname, avatar_url: avatarURL })
    if (userInfo.value) {
      userInfo.value.nickname = res.data.nickname
      userInfo.value.avatar_url = res.data.avatar_url
    }
    showProfileSetup.value = false
    return res.data
  }

  // 获取 AppID
  function getAppId() {
    // #ifdef MP-WEIXIN
    return 'wxff0ecb7fddca4ecc'
    // #endif
    // #ifdef H5
    return 'h5-yanny'
    // #endif
    return ''
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    uni.removeStorageSync('mp_token')
  }

  // 等待登录就绪（页面发起业务请求前调用，避免 token 未就绪时被 401）
  function waitForReady() {
    if (ready.value) return Promise.resolve()
    return new Promise((resolve) => {
      const timer = setInterval(() => {
        if (ready.value) {
          clearInterval(timer)
          resolve()
        }
      }, 100)
    })
  }

  return {
    userInfo, token, platform, ready, loginError, showProfileSetup, isLoggedIn, userId,
    silentLogin, doLogin, updatePhone, updateProfile, logout, detectPlatform, waitForReady,
  }
})
