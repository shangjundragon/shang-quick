import axios from 'axios'
import { useUserStore } from '@/store/user'

// Axios 实例：baseURL /api，通过 Vite 代理转发到后端
const service = axios.create({
  baseURL: '/api',
  timeout: 30 * 1000,
  headers: {
    'Content-Type': 'application/json;charset=utf-8'
  }
})

// 请求拦截器：自动注入 Bearer Token
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

// 响应拦截器：统一处理业务状态码（401/403/500 等）
service.interceptors.response.use(
  async (response) => {
    const { data, config } = response
    const processCode = config.processCode || {}

    // 文件下载类请求透传原始响应
    if (response.request.responseType === 'blob' || data instanceof Blob) {
      return response
    }

    const { code, data: resData, message } = data

    // 自定义状态码处理器（通过 config.processCode 传入）
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
    // HTTP 网络层错误处理
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
