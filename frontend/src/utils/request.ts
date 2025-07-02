import { ElMessage } from 'element-plus'
import { useAuthStore } from '../stores/auth'

interface RequestConfig extends RequestInit {
  retry?: number
  retryDelay?: number
}

interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081/api'
const DEFAULT_RETRY_TIMES = 2
const DEFAULT_RETRY_DELAY = 1000

export class ApiError extends Error {
  code: number
  constructor(message: string, code: number) {
    super(message)
    this.code = code
    this.name = 'ApiError'
  }
}

export async function request<T>(endpoint: string, config: RequestConfig = {}): Promise<T> {
  const authStore = useAuthStore()
  const retryTimes = config.retry ?? DEFAULT_RETRY_TIMES
  const retryDelay = config.retryDelay ?? DEFAULT_RETRY_DELAY
  let attempts = 0

  const executeRequest = async (): Promise<T> => {
    try {
      // 使用 authStore.getAuthHeaders() 获取认证头
      const headers = new Headers({
        ...authStore.getAuthHeaders(),
        ...config.headers
      })

      // 合并配置
      const finalConfig = {
        ...config,
        headers
      }

      // 发送请求
      const response = await fetch(`${API_BASE_URL}${endpoint}`, finalConfig)
      const result: ApiResponse<T> = await response.json()

      // 处理业务错误
      if (result.code !== 200) {
        // token过期
        if (result.code === 401) {
          authStore.clearAuth()
          window.location.href = '/login'
          throw new ApiError('登录已过期，请重新登录', 401)
        }
        throw new ApiError(result.message || '请求失败', result.code)
      }

      return result.data
    } catch (error) {
      // 网络错误或解析错误时重试
      if (!(error instanceof ApiError) && attempts < retryTimes) {
        attempts++
        await new Promise(resolve => setTimeout(resolve, retryDelay))
        return executeRequest()
      }
      throw error
    }
  }

  return executeRequest()
}

// 请求方法封装
export const http = {
  async get<T>(endpoint: string, config: RequestConfig = {}) {
    return request<T>(endpoint, { ...config, method: 'GET' })
  },

  async post<T>(endpoint: string, data?: any, config: RequestConfig = {}) {
    return request<T>(endpoint, {
      ...config,
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...config.headers,
      },
      body: JSON.stringify(data),
    })
  },

  async put<T>(endpoint: string, data?: any, config: RequestConfig = {}) {
    return request<T>(endpoint, {
      ...config,
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        ...config.headers,
      },
      body: JSON.stringify(data),
    })
  },

  async delete<T>(endpoint: string, config: RequestConfig = {}) {
    return request<T>(endpoint, { ...config, method: 'DELETE' })
  },
}

// 错误处理中间件
export function handleApiError(error: unknown): never {
  if (error instanceof ApiError) {
    ElMessage.error(error.message)
  } else if (error instanceof Error) {
    ElMessage.error('网络错误，请稍后重试')
    console.error('Network Error:', error)
  } else {
    ElMessage.error('未知错误')
    console.error('Unknown Error:', error)
  }
  throw error
} 