import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { ElMessage } from 'element-plus'

// 用户信息类型
export interface UserInfo {
  id: number
  username: string
  email: string
  role: string
  status: boolean
  createdAt: string
  updatedAt: string
}

// 登录请求类型
export interface LoginForm {
  username: string
  password: string
}

// 注册请求类型
export interface RegisterForm {
  username: string
  email: string
  password: string
  confirmPassword?: string
}

// 登录响应类型
interface LoginResponse {
  token: string
  user_info: UserInfo
}

// API响应类型
interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref<string | null>(localStorage.getItem('token'))
  const userInfo = ref<UserInfo | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const tokenCheckInterval = ref<number | null>(null)

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!userInfo.value)
  const isAdmin = computed(() => userInfo.value?.role === 'admin')
  const currentUser = computed(() => userInfo.value)

  // API基础配置
  const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081/api'
  const TOKEN_CHECK_INTERVAL = 5 * 60 * 1000 // 5分钟检查一次token

  // 设置认证头
  function getAuthHeaders() {
    return {
      'Content-Type': 'application/json',
      'Authorization': token.value ? `Bearer ${token.value}` : '',
    }
  }

  // 启动token定时检查
  function startTokenCheck() {
    stopTokenCheck() // 先清除可能存在的定时器
    tokenCheckInterval.value = window.setInterval(async () => {
      if (token.value) {
        const isValid = await getCurrentUser()
        if (!isValid) {
          stopTokenCheck()
          ElMessage.warning('登录已过期，请重新登录')
        }
      }
    }, TOKEN_CHECK_INTERVAL)
  }

  // 停止token检查
  function stopTokenCheck() {
    if (tokenCheckInterval.value) {
      window.clearInterval(tokenCheckInterval.value)
      tokenCheckInterval.value = null
    }
  }

  // 清除认证状态
  function clearAuth() {
    token.value = null
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    stopTokenCheck()
  }

  // Actions
  async function login(loginForm: LoginForm) {
    loading.value = true
    error.value = null
    try {
      const response = await fetch(`${API_BASE_URL}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(loginForm),
      })

      const result: ApiResponse<LoginResponse> = await response.json()
      
      if (result.code === 200) {
        token.value = result.data.token
        userInfo.value = result.data.user_info
        
        // 保存到本地存储
        localStorage.setItem('token', result.data.token)
        localStorage.setItem('userInfo', JSON.stringify(result.data.user_info))
        
        // 启动token检查
        startTokenCheck()
        
        ElMessage.success('登录成功')
        return true
      } else {
        throw new Error(result.message || '登录失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      ElMessage.error(error.value)
      return false
    } finally {
      loading.value = false
    }
  }

  async function register(registerForm: RegisterForm) {
    loading.value = true
    error.value = null
    try {
      const { confirmPassword, ...submitData } = registerForm
      
      const response = await fetch(`${API_BASE_URL}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(submitData),
      })

      const result: ApiResponse<LoginResponse> = await response.json()
      
      if (result.code === 200) {
        token.value = result.data.token
        userInfo.value = result.data.user_info
        
        // 保存到本地存储
        localStorage.setItem('token', result.data.token)
        localStorage.setItem('userInfo', JSON.stringify(result.data.user_info))
        
        ElMessage.success('注册成功')
        return true
      } else {
        throw new Error(result.message || '注册失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      ElMessage.error(error.value)
      return false
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    loading.value = true
    try {
      // 调用后端登出接口
      if (token.value) {
        await fetch(`${API_BASE_URL}/auth/logout`, {
          method: 'POST',
          headers: getAuthHeaders(),
        })
      }
    } catch (err) {
      console.warn('Logout API call failed:', err)
    } finally {
      clearAuth()
      loading.value = false
      ElMessage.success('已退出登录')
    }
  }

  async function getCurrentUser() {
    if (!token.value) return false
    
    loading.value = true
    error.value = null
    try {
      const response = await fetch(`${API_BASE_URL}/me`, {
        headers: getAuthHeaders(),
      })

      const result: ApiResponse<UserInfo> = await response.json()
      
      if (result.code === 200) {
        userInfo.value = result.data
        localStorage.setItem('userInfo', JSON.stringify(result.data))
        return true
      } else {
        // Token可能过期，清除认证状态但不调用logout API
        console.warn('Get current user failed:', result.message)
        token.value = null
        userInfo.value = null
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        return false
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      console.error('Get current user network error:', err)
      // 网络错误时不清除token，可能是临时问题
      return false
    } finally {
      loading.value = false
    }
  }

  async function changePassword(oldPassword: string, newPassword: string) {
    loading.value = true
    error.value = null
    try {
      const response = await fetch(`${API_BASE_URL}/change-password`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({
          old_password: oldPassword,
          new_password: newPassword,
        }),
      })

      const result: ApiResponse<any> = await response.json()
      
      if (result.code === 200) {
        ElMessage.success('密码修改成功')
        return true
      } else {
        throw new Error(result.message || '密码修改失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      ElMessage.error(error.value)
      return false
    } finally {
      loading.value = false
    }
  }

  // 初始化认证状态
  function initAuth() {
    const savedToken = localStorage.getItem('token')
    const savedUserInfo = localStorage.getItem('userInfo')
    
    if (savedToken && savedUserInfo) {
      try {
        token.value = savedToken
        userInfo.value = JSON.parse(savedUserInfo)
        // 启动token检查
        startTokenCheck()
      } catch (err) {
        console.error('Failed to parse saved user info:', err)
        clearAuth()
      }
    }
  }

  // 检查权限
  function hasPermission(permission: string): boolean {
    if (!userInfo.value) return false
    
    // 管理员拥有所有权限
    if (userInfo.value.role === 'admin') return true
    
    // 根据具体需求实现权限检查逻辑
    // 这里简化处理
    return userInfo.value.role === permission
  }

  // 清除错误
  function clearError() {
    error.value = null
  }

  // 在组件卸载时清理
  function cleanup() {
    stopTokenCheck()
  }

  return {
    // 状态
    token,
    userInfo,
    loading,
    error,
    
    // Getters
    isAuthenticated,
    isAdmin,
    currentUser,
    
    // Actions
    login,
    register,
    logout,
    getCurrentUser,
    changePassword,
    getAuthHeaders,
    initAuth,
    cleanup,
    clearError,
    hasPermission,
    clearAuth,
  }
}) 