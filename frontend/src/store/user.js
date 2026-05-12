import { defineStore } from 'pinia'

// 用户状态管理：Token（localStorage 持久化）、用户信息、权限列表、菜单树
export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: null,
    permissions: [],
    menus: []
  }),
  actions: {
    setToken(token) {
      this.token = token
      localStorage.setItem('token', token)
    },
    setUserInfo(userInfo) {
      this.userInfo = userInfo
    },
    setPermissions(permissions) {
      this.permissions = permissions
    },
    setMenus(menus) {
      this.menus = menus
    },
    logout() {
      // 清除所有用户状态并跳转登录页
      this.token = ''
      this.userInfo = null
      this.permissions = []
      this.menus = []
      localStorage.removeItem('token')
    }
  }
})
