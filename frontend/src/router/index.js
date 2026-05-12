import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'
import { getInfo } from '@/api/auth'

// 路由在编译时静态定义，侧边栏菜单由后端返回的 menus 动态渲染
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

// allRoutes 定义了所有业务路由，在 router 创建后通过 addRoute 动态注册
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
      },
      {
        path: 'onlineUser',
        name: 'OnlineUser',
        component: () => import('@/views/system/onlineUser/index.vue'),
        meta: { title: '在线用户', icon: 'WifiOutline' }
      }
    ]
  }
]

allRoutes.forEach(route => {
  router.addRoute(route)
})

// 未匹配路由重定向到登录页
router.addRoute({
  path: '/:pathMatch(.*)*',
  name: 'NotFound',
  redirect: '/login',
  hidden: true
})

// 全局路由守卫：未登录跳转登录页，已登录但无 menus 时拉取用户信息
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

  // 刷新页面后重新获取用户信息（菜单、权限等）
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
