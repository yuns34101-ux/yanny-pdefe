import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录' },
  },
  {
    path: '/',
    component: () => import('@/views/layout/index.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/analytics/dashboard.vue'),
        meta: { title: '数据看板', icon: 'DataAnalysis', perm: 'analytics:view' },
      },
      {
        path: 'entities',
        name: 'Entities',
        component: () => import('@/views/entity/index.vue'),
        meta: { title: '主体管理', icon: 'OfficeBuilding', perm: 'entity:view' },
      },
      {
        path: 'mp-accounts',
        name: 'MpAccounts',
        component: () => import('@/views/mp/index.vue'),
        meta: { title: '小程序账号', icon: 'Cellphone', perm: 'mp_account:view' },
      },
      {
        path: 'videos',
        name: 'Videos',
        component: () => import('@/views/video/index.vue'),
        meta: { title: '视频管理', icon: 'VideoCamera', perm: 'video:view' },
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/user/index.vue'),
        meta: { title: '用户管理', icon: 'User', perm: 'user:view' },
      },
      {
        path: 'admins',
        name: 'Admins',
        component: () => import('@/views/admin/index.vue'),
        meta: { title: '管理员', icon: 'UserFilled', perm: 'admin:view' },
      },
      {
        path: 'roles',
        name: 'Roles',
        component: () => import('@/views/admin/roles.vue'),
        meta: { title: '角色权限', icon: 'Lock', perm: 'role:view' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  if (to.path === '/login') {
    if (authStore.isLoggedIn) {
      return next('/')
    }
    return next()
  }

  if (!authStore.isLoggedIn) {
    return next('/login')
  }

  // 首次进入时加载用户信息
  if (!authStore.userInfo) {
    try {
      await authStore.loadUserInfo()
    } catch {
      authStore.logout()
      return next('/login')
    }
  }

  // 权限校验
  const requiredPerm = to.meta.perm
  if (requiredPerm && !authStore.hasPermission(requiredPerm)) {
    return next('/')
  }

  next()
})

export default router
