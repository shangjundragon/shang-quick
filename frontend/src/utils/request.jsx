import axios from 'axios'
import { useUserStore } from '@/store/user'

const service = axios.create({
  baseURL: '/api',
  timeout: 30 * 1000,
  headers: {
    'Content-Type': 'application/json;charset=utf-8'
  }
})

service.interceptors.request.use(
  async (config) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

service.interceptors.response.use(
  async (response) => {
    const { data, config } = response
    const processCode = config.processCode || {}

    if (response.request.responseType === 'blob' || data instanceof Blob) {
      return response
    }

    const { code, data: resData, message } = data

    if (code && code !== 200 && processCode[code]) {
      return processCode[code](data)
    }

    if (code && code !== 200 && processCode['unSuccessAll']) {
      return processCode['unSuccessAll'](data)
    }

    switch (code) {
      case 200:
        return resData
      case 401: {
        const userStore = useUserStore()
        userStore.logout()
        window.location.href = '/login'
        return Promise.reject(new Error(message || '登录已过期'))
      }
      case 403: {
        window?.$message?.error(message || '无权限')
        return Promise.reject(new Error(message || '无权限'))
      }
      case 500:
        window?.$message?.error(message || '服务器异常')
        return Promise.reject(new Error(message || '服务器异常'))
      default:
        window?.$message?.error(message || '请求失败')
        return Promise.reject(new Error(message || '请求失败'))
    }
  },
  (error) => {
    let message
    if (error.response) {
      const { status } = error.response
      switch (status) {
        case 404:
          message = '请求地址不存在'
          break
        case 502:
        case 503:
        case 504:
          message = '服务器开小差了~'
          break
        default:
          message = '网络请求异常'
      }
    } else if (error.message.includes('timeout')) {
      message = '请求超时，请重试'
    } else if (error.message.includes('Network Error')) {
      message = '服务端连接错误'
    } else {
      message = '网络连接失败'
    }
    window?.$message?.error(message || '请求失败')
    return Promise.reject(error)
  }
)

export default service
