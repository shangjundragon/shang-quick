import axios from 'axios'

/**
 * 1. 创建axios实例
 */
const service = axios.create({
  baseURL: '/api',
  timeout: 30 * 1000, // 请求超时时间
  headers: {
    'Content-Type': 'application/json;charset=utf-8'
  }
})

/**
 * 2. 请求拦截器：统一注入Token
 */
service.interceptors.request.use(
  async (config) => {
    // config.headers.Authorization = `Bearer ${token}`
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

/**
 * 3. 响应拦截器：统一处理响应、状态码、文件流
 */
service.interceptors.response.use(
  async (response) => {
    const { data, config } = response
    const processCode = config.processCode || {}
    // ====================== 场景1：处理文件流下载 ======================
    if (response.request.responseType === 'blob' || data instanceof Blob) {
      // 直接返回完整响应，用于文件下载
      return response
    }

    // ====================== 场景2：处理JSON标准响应 ======================
    const { code, data: resData, message } = data

    if (code && code !== 200 && processCode[code]) {
      return processCode[code](data)
    }

    if (code && code !== 200 && processCode['unSuccessAll']) {
      return processCode['unSuccessAll'](data)
    }

    switch (code) {
      // 200：业务成功 → 直接返回data
      case 200:
        return resData

      // 401：Token过期/未登录 → 清除Token + 提示 + 跳转登录
      case 401: {
        return Promise.reject(new Error(message || '登录已过期'))
      }

      case 403: {
        return Promise.reject(new Error(message || '无权限'))
      }

      // 500：服务端业务异常 → 提示错误
      case 500:
        window?.$message.error(message || '服务器异常')
        return Promise.reject(new Error(message || '服务器异常'))

      // 其他状态码 → 统一提示
      default:
        window?.$message.error(message || '请求失败')
        return Promise.reject(new Error(message || '请求失败'))
    }
  },
  // ====================== HTTP网络错误处理 ======================
  (error) => {
    console.log('onRejected error', error)
    console.log('onRejected error.message', error.message)
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
    console.error('axios onRejected error.message', message)
      window?.$message.error(message || '请求失败')
    return Promise.reject(error)
  }
)



/**
 * 4. 导出常用请求方法
 */
export default service
