import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

// 角色类型定义
export interface Role {
  id: number
  name: string
  display_name: string
  description: string
  status: boolean
  permissions?: Permission[]
  CreatedAt: string
  UpdatedAt: string
}

// 权限类型定义
export interface Permission {
  id: number
  name: string
  display_name: string
  resource: string
  action: string
  description: string
  CreatedAt: string
  UpdatedAt: string
}

// API响应类型
interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export const useRoleStore = defineStore('role', () => {
  // 状态
  const roles = ref<Role[]>([])
  const permissions = ref<Permission[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const total = ref(0)

  // Getters
  const roleCount = computed(() => roles.value.length)
  const activeRoles = computed(() => roles.value.filter(role => role.status))

  // API基础配置
  const API_BASE_URL = 'http://localhost:8081/api'

  // Actions
  async function fetchRoles() {
    loading.value = true
    error.value = null
    try {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      
      const response = await fetch(`${API_BASE_URL}/roles`, {
        headers: authStore.getAuthHeaders(),
      })
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const result: ApiResponse<{ roles: Role[]; total: number }> = await response.json()
      
      if (result.code === 200) {
        roles.value = result.data.roles
        total.value = result.data.total
      } else {
        throw new Error(result.message || '获取角色列表失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      console.error('Error fetching roles:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchPermissions() {
    loading.value = true
    error.value = null
    try {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      
      const response = await fetch(`${API_BASE_URL}/permissions`, {
        headers: authStore.getAuthHeaders(),
      })
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const result: ApiResponse<{ permissions: Permission[]; total: number }> = await response.json()
      
      if (result.code === 200) {
        permissions.value = result.data.permissions
      } else {
        throw new Error(result.message || '获取权限列表失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      console.error('Error fetching permissions:', err)
    } finally {
      loading.value = false
    }
  }

  async function createRole(roleData: Omit<Role, 'id' | 'CreatedAt' | 'UpdatedAt'>) {
    loading.value = true
    error.value = null
    try {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      
      const response = await fetch(`${API_BASE_URL}/roles`, {
        method: 'POST',
        headers: authStore.getAuthHeaders(),
        body: JSON.stringify(roleData),
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const result: ApiResponse<Role> = await response.json()
      if (result.code === 200) {
        roles.value.push(result.data)
        total.value++
        return result.data
      } else {
        throw new Error(result.message || '创建角色失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      console.error('Error creating role:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateRole(id: number, roleData: Partial<Role>) {
    loading.value = true
    error.value = null
    try {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      
      const response = await fetch(`${API_BASE_URL}/roles/${id}`, {
        method: 'PUT',
        headers: authStore.getAuthHeaders(),
        body: JSON.stringify(roleData),
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const result: ApiResponse<Role> = await response.json()
      if (result.code === 200) {
        const index = roles.value.findIndex(role => role.id === id)
        if (index !== -1) {
          roles.value[index] = result.data
        }
        return result.data
      } else {
        throw new Error(result.message || '更新角色失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      console.error('Error updating role:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteRole(id: number) {
    loading.value = true
    error.value = null
    try {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      
      const response = await fetch(`${API_BASE_URL}/roles/${id}`, {
        method: 'DELETE',
        headers: authStore.getAuthHeaders(),
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const result: ApiResponse<any> = await response.json()
      if (result.code === 200) {
        roles.value = roles.value.filter(role => role.id !== id)
        total.value--
      } else {
        throw new Error(result.message || '删除角色失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      console.error('Error deleting role:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function assignPermissions(roleId: number, permissionIds: number[]) {
    loading.value = true
    error.value = null
    try {
      const { useAuthStore } = await import('./auth')
      const authStore = useAuthStore()
      
      const response = await fetch(`${API_BASE_URL}/permissions/assign`, {
        method: 'POST',
        headers: authStore.getAuthHeaders(),
        body: JSON.stringify({
          role_id: roleId,
          permission_ids: permissionIds,
        }),
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const result: ApiResponse<any> = await response.json()
      if (result.code === 200) {
        // 更新本地角色数据
        const roleIndex = roles.value.findIndex(role => role.id === roleId)
        if (roleIndex !== -1 && result.data.role) {
          roles.value[roleIndex] = result.data.role
        }
      } else {
        throw new Error(result.message || '权限分配失败')
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : '网络错误'
      console.error('Error assigning permissions:', err)
      throw err
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
    roles.value = []
    permissions.value = []
    loading.value = false
    error.value = null
    total.value = 0
  }

  return {
    // 状态
    roles,
    permissions,
    loading,
    error,
    total,
    // Getters
    roleCount,
    activeRoles,
    // Actions
    fetchRoles,
    fetchPermissions,
    createRole,
    updateRole,
    deleteRole,
    assignPermissions,
    clearError,
    resetState,
  }
}) 