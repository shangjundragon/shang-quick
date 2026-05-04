import { defineStore } from 'pinia'

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
      this.token = ''
      this.userInfo = null
      this.permissions = []
      this.menus = []
      localStorage.removeItem('token')
    }
  }
})
