import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

// 用户类型定义
export interface User {
  id: number
  username: string
  email: string
  role: string
  status: boolean
  createdAt: string
  updatedAt: string
}

// API响应类型
interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export const useUserStore = defineStore('user', () => {
  // 状态
  const users = ref<User[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const total = ref(0)

  // Getters
  const userCount = computed(() => users.value.length)
  const adminUsers = computed(() => users.value.filter(user => user.role === 'admin'))
  const activeUsers = computed(() => users.value.filter(user => user.status))

  // API基础配置
  const API_BASE_URL = 'http://localhost:8081/api'

  // Actions
  async function fetchUsers() {
    loading.value = true
    error.value = null
    try {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      
      const response = await fetch(`${API_BASE_URL}/users`, {
        headers: authStore.getAuthHeaders(),
      })
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const result: ApiResponse<{ users: User[]; total: number }> = await response.json()
      
      if (result.code === 200) {
        users.value = result.data.users
        total.value = result.data.total
      } else {
        throw new Error(result.message || '获取用户列表失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      console.error('Error fetching users:', err)
    } finally {
      loading.value = false
    }
  }

  async function createUser(userData: Omit<User, 'id' | 'createdAt' | 'updatedAt'>) {
    loading.value = true
    error.value = null
    try {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      
      const response = await fetch(`${API_BASE_URL}/users`, {
        method: 'POST',
        headers: authStore.getAuthHeaders(),
        body: JSON.stringify(userData),
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const result: ApiResponse<User> = await response.json()
      if (result.code === 200) {
        users.value.push(result.data)
        total.value++
        return result.data
      } else {
        throw new Error(result.message || '创建用户失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      console.error('Error creating user:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateUser(id: number, userData: Partial<User>) {
    loading.value = true
    error.value = null
    try {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      
      const response = await fetch(`${API_BASE_URL}/users/${id}`, {
        method: 'PUT',
        headers: authStore.getAuthHeaders(),
        body: JSON.stringify(userData),
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const result: ApiResponse<User> = await response.json()
      if (result.code === 200) {
        const index = users.value.findIndex(user => user.id === id)
        if (index !== -1) {
          users.value[index] = result.data
        }
        return result.data
      } else {
        throw new Error(result.message || '更新用户失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      console.error('Error updating user:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteUser(id: number) {
    loading.value = true
    error.value = null
    try {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      
      const response = await fetch(`${API_BASE_URL}/users/${id}`, {
        method: 'DELETE',
        headers: authStore.getAuthHeaders(),
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const result: ApiResponse<any> = await response.json()
      if (result.code === 200) {
        users.value = users.value.filter(user => user.id !== id)
        total.value--
      } else {
        throw new Error(result.message || '删除用户失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      console.error('Error deleting user:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function getUserById(id: number): Promise<User | null> {
    loading.value = true
    error.value = null
    try {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      
      const response = await fetch(`${API_BASE_URL}/users/${id}`, {
        headers: authStore.getAuthHeaders(),
      })
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const result: ApiResponse<User> = await response.json()
      if (result.code === 200) {
        return result.data
      } else {
        throw new Error(result.message || '获取用户信息失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      console.error('Error fetching user:', err)
      return null
    } finally {
      loading.value = false
    }
  }

  // 清除错误
  function clearError() {
    error.value = null
  }

  // 重置状态
  function resetState() {
    users.value = []
    loading.value = false
    error.value = null
    total.value = 0
  }

  return {
    // 状态
    users,
    loading,
    error,
    total,
    // Getters
    userCount,
    adminUsers,
    activeUsers,
    // Actions
    fetchUsers,
    createUser,
    updateUser,
    deleteUser,
    getUserById,
    clearError,
    resetState,
  }
}) 