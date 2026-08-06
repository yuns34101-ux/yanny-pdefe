import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, getAdminInfo } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const userInfo = ref(null)
  const permissions = ref([])

  const isLoggedIn = computed(() => !!token.value)

  const hasPermission = (permCode) => {
    return permissions.value.includes(permCode)
  }

  const doLogin = async (username, password) => {
    const res = await login(username, password)
    token.value = res.data.token
    localStorage.setItem('admin_token', res.data.token)
    await loadUserInfo()
    return res
  }

  const loadUserInfo = async () => {
    try {
      const res = await getAdminInfo()
      userInfo.value = res.data
      permissions.value = res.data.perms || []
    } catch {
      // token 失效
    }
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    permissions.value = []
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_perm_cache')
  }

  return { token, userInfo, permissions, isLoggedIn, hasPermission, doLogin, loadUserInfo, logout }
})
