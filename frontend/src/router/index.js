import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getInfo } from '@/api/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    hidden: true
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export const allRoutes = [
  {
    path: '/',
    component: () => import('@/layout/index.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '仪表盘', icon: 'SpeedometerOutline' }
      }
    ]
  },
  {
    path: '/profile',
    component: () => import('@/layout/index.vue'),
    children: [
      {
        path: '',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人中心', icon: 'PersonOutline' }
      }
    ]
  },
  {
    path: '/system',
    component: () => import('@/layout/index.vue'),
    redirect: '/system/user',
    meta: { title: '系统管理', icon: 'SettingsOutline' },
    children: [
      {
        path: 'user',
        name: 'User',
        component: () => import('@/views/system/user/index.vue'),
        meta: { title: '用户管理', icon: 'PeopleOutline' }
      },
      {
        path: 'dept',
        name: 'Dept',
        component: () => import('@/views/system/dept/index.vue'),
        meta: { title: '部门管理', icon: 'BusinessOutline' }
      },
      {
        path: 'menu',
        name: 'Menu',
        component: () => import('@/views/system/menu/index.vue'),
        meta: { title: '菜单管理', icon: 'MenuOutline' }
      },
      {
        path: 'role',
        name: 'Role',
        component: () => import('@/views/system/role/index.vue'),
        meta: { title: '角色管理', icon: 'ShieldOutline' }
      },
      {
        path: 'operLog',
        name: 'OperLog',
        component: () => import('@/views/system/operLog/index.vue'),
        meta: { title: '操作日志', icon: 'DocumentTextOutline' }
      },
      {
        path: 'file',
        name: 'File',
        component: () => import('@/views/system/file/index.vue'),
        meta: { title: '文件管理', icon: 'FolderOutline' }
      }
    ]
  }
]

allRoutes.forEach(route => {
  router.addRoute(route)
})

router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()

  if (to.path === '/login') {
    next()
    return
  }

  if (!userStore.token) {
    next('/login')
    return
  }

  // 刷新页面后重新获取用户信息
  if (userStore.menus.length === 0) {
    try {
      const res = await getInfo()
      userStore.setUserInfo(res.userInfo)
      userStore.setPermissions(res.permissions)
      userStore.setMenus(res.menus)
    } catch (error) {
      userStore.logout()
      next('/login')
      return
    }
  }

  next()
})

export default router
