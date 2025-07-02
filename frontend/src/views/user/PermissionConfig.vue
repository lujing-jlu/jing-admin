<template>
  <div class="page-container">
    <!-- 页面头部 -->
    <div class="page-header fade-in">
      <div class="header-info">
        <h1 class="page-title">权限配置</h1>
        <p class="page-subtitle">管理角色权限分配和权限关系</p>
      </div>
      <div class="header-actions">
        <el-button 
          type="primary" 
          :icon="Refresh" 
          @click="refreshData"
          :loading="loading"
          size="default"
        >
          刷新数据
        </el-button>
        <el-button :icon="Download" @click="exportPermissions" style="margin-left: 8px;">
          导出权限
        </el-button>
        <el-button :icon="Upload" @click="importDialogVisible = true" style="margin-left: 8px;">
          导入权限
        </el-button>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-grid fade-in" style="animation-delay: 0.1s">
      <el-row :gutter="20">
        <!-- 左侧：角色列表 -->
        <el-col :xl="8" :lg="10" :md="24" :sm="24" :xs="24">
          <div class="card-container roles-card">
            <div class="card-header">
              <div class="card-title">
                <el-icon><UserFilled /></el-icon>
                <span>角色列表</span>
              </div>
              <el-tag type="info" size="small">{{ roles.length }} 个角色</el-tag>
            </div>
            <div class="card-content">
              <div class="role-list">
                <div 
                  v-for="role in roles" 
                  :key="role.id"
                  class="role-item"
                  :class="{ 'active': selectedRole?.id === role.id }"
                  @click="selectRole(role)"
                >
                  <div class="role-info">
                    <div class="role-name">{{ role.name }}</div>
                    <div class="role-desc">{{ role.description }}</div>
                    <div class="role-meta">
                      <el-tag 
                        :type="role.status ? 'success' : 'danger'" 
                        size="small"
                      >
                        {{ role.status ? '启用' : '禁用' }}
                      </el-tag>
                      <span class="permission-count">
                        {{ role.permissions?.length || 0 }} 个权限
                      </span>
                    </div>
                  </div>
                  <el-icon class="role-arrow">
                    <ArrowRight />
                  </el-icon>
                </div>
              </div>
            </div>
          </div>
        </el-col>

        <!-- 右侧：权限配置 -->
        <el-col :xl="16" :lg="14" :md="24" :sm="24" :xs="24">
          <div class="card-container permissions-card">
            <div class="card-header">
              <div class="card-title">
                <el-icon><Key /></el-icon>
                <span>
                  {{ selectedRole ? `配置"${selectedRole.name}"权限` : '权限管理' }}
                </span>
              </div>
              <div v-if="selectedRole" class="permission-actions">
                <el-button 
                  type="success" 
                  :icon="Check" 
                  @click="savePermissions"
                  :loading="saveLoading"
                  size="small"
                >
                  保存配置
                </el-button>
                <el-button 
                  type="warning" 
                  :icon="RefreshLeft" 
                  @click="resetPermissions"
                  size="small"
                >
                  重置
                </el-button>
              </div>
            </div>
            <div class="card-content">
              <!-- 未选择角色时的提示 -->
              <div v-if="!selectedRole" class="empty-state">
                <el-icon :size="64" color="#cbd5e1"><Key /></el-icon>
                <h3>请选择角色</h3>
                <p>从左侧选择一个角色来配置其权限</p>
              </div>

              <!-- 权限配置区域 -->
              <div v-else class="permissions-config">
                <!-- 权限统计 -->
                <div class="permission-stats">
                  <div class="stat-item">
                    <span class="stat-label">已选择权限:</span>
                    <span class="stat-value">{{ selectedPermissions.length }}</span>
                  </div>
                  <div class="stat-item">
                    <span class="stat-label">总权限数:</span>
                    <span class="stat-value">{{ allPermissions.length }}</span>
                  </div>
                  <div class="actions">
                    <el-button 
                      type="primary" 
                      :icon="Check" 
                      @click="selectAllPermissions"
                      size="small"
                      link
                    >
                      全选
                    </el-button>
                    <el-button 
                      type="warning" 
                      :icon="Close" 
                      @click="clearAllPermissions"
                      size="small"
                      link
                    >
                      清空
                    </el-button>
                  </div>
                </div>

                <!-- 按资源分组的权限列表 -->
                <div class="permission-groups">
                  <div 
                    v-for="(permissions, resource) in groupedPermissions" 
                    :key="resource"
                    class="permission-group"
                  >
                    <div class="group-header">
                      <el-checkbox
                        :modelValue="isGroupSelected(resource)"
                        :indeterminate="isGroupIndeterminate(resource)"
                        @change="toggleGroup(resource, $event)"
                      >
                        <div class="group-title">
                          <el-icon class="group-icon">
                            <component :is="getResourceIcon(resource)" />
                          </el-icon>
                          <span>{{ getResourceText(resource) }}</span>
                          <el-tag type="info" size="small">
                            {{ permissions.length }} 项
                          </el-tag>
                        </div>
                      </el-checkbox>
                    </div>
                    <div class="group-content">
                      <el-checkbox-group v-model="selectedPermissions">
                        <div 
                          v-for="permission in permissions" 
                          :key="permission.id"
                          class="permission-item"
                        >
                          <el-checkbox :label="permission.id">
                            <div class="permission-info">
                              <div class="permission-main">
                                <span class="permission-name">{{ permission.display_name }}</span>
                                <el-tag 
                                  :type="getActionTagType(permission.action)"
                                  size="small"
                                >
                                  {{ getActionText(permission.action) }}
                                </el-tag>
                              </div>
                              <div class="permission-desc">{{ permission.description }}</div>
                              <div class="permission-code">{{ permission.name }}</div>
                            </div>
                          </el-checkbox>
                        </div>
                      </el-checkbox-group>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 权限详情对话框 -->
    <el-dialog
      v-model="permissionDetailVisible"
      title="权限详情"
      width="600px"
      destroy-on-close
    >
      <div v-if="selectedPermissionDetail" class="permission-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="权限ID">{{ selectedPermissionDetail.id }}</el-descriptions-item>
          <el-descriptions-item label="权限代码">{{ selectedPermissionDetail.name }}</el-descriptions-item>
          <el-descriptions-item label="显示名称">{{ selectedPermissionDetail.display_name }}</el-descriptions-item>
          <el-descriptions-item label="操作类型">
            <el-tag :type="getActionTagType(selectedPermissionDetail.action)">
              {{ getActionText(selectedPermissionDetail.action) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="资源类型">{{ getResourceText(selectedPermissionDetail.resource) }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDateTime(selectedPermissionDetail.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">{{ selectedPermissionDetail.description }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <!-- 权限导入弹窗 -->
    <el-dialog v-model="importDialogVisible" title="批量导入权限" width="500px" center>
      <FileUpload
        ref="fileUploadRef"
        category="document"
        :allowed-types="['csv']"
        :auto-upload="false"
        :show-button="false"
        :hide-drop-zone="false"
        @success="handleImportSuccess"
        @error="handleImportError"
      />
      <template #footer>
        <el-button @click="importDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="importing" @click="submitImport">开始导入</el-button>
      </template>
      <div v-if="importResult" class="import-result">
        <el-alert
          :title="`导入完成：成功 ${importResult.success} 条，失败 ${importResult.failed} 条`"
          :type="importResult.failed === 0 ? 'success' : 'warning'"
          show-icon
          style="margin-bottom: 12px;"
        />
        <el-scrollbar v-if="importResult.errors && importResult.errors.length" style="max-height: 120px;">
          <ul class="import-errors">
            <li v-for="(err, idx) in importResult.errors" :key="idx" style="color: var(--danger-color); font-size: 13px;">{{ err }}</li>
          </ul>
        </el-scrollbar>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import {
  UserFilled,
  Key,
  Check,
  Close,
  RefreshLeft,
  Refresh,
  ArrowRight,
  User,
  Setting,
  Monitor,
  Lock,
  Download,
  Upload
} from '@element-plus/icons-vue'
import FileUpload from '../../components/FileUpload.vue'

// 类型定义
interface Permission {
  id: number
  name: string
  display_name: string
  resource: string
  action: string
  description: string
  created_at: string
}

interface Role {
  id: number
  name: string
  description: string
  status: boolean
  permissions?: Permission[]
  created_at: string
}

// 响应式数据
const loading = ref(false)
const saveLoading = ref(false)
const roles = ref<Role[]>([])
const allPermissions = ref<Permission[]>([])
const selectedRole = ref<Role | null>(null)
const selectedPermissions = ref<number[]>([])
const originalPermissions = ref<number[]>([])
const permissionDetailVisible = ref(false)
const selectedPermissionDetail = ref<Permission | null>(null)
const importDialogVisible = ref(false)
const importing = ref(false)
const importResult = ref<any>(null)
const fileUploadRef = ref()

// 计算属性
const groupedPermissions = computed(() => {
  const groups: Record<string, Permission[]> = {}
  allPermissions.value.forEach(permission => {
    if (!groups[permission.resource]) {
      groups[permission.resource] = []
    }
    groups[permission.resource].push(permission)
  })
  return groups
})

// 获取角色列表
const getRolesList = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('http://localhost:8081/api/roles', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    const result = await response.json()
    if (result.code === 200) {
      roles.value = result.data.roles || []
    } else {
      ElMessage.error(result.message || '获取角色列表失败')
    }
  } catch (error) {
    console.error('获取角色列表错误:', error)
    ElMessage.error('网络错误，请稍后重试')
  }
}

// 获取权限列表
const getPermissionsList = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('http://localhost:8081/api/permissions', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    const result = await response.json()
    if (result.code === 200) {
      allPermissions.value = result.data.permissions || []
    } else {
      ElMessage.error(result.message || '获取权限列表失败')
    }
  } catch (error) {
    console.error('获取权限列表错误:', error)
    ElMessage.error('网络错误，请稍后重试')
  }
}

// 选择角色
const selectRole = (role: Role) => {
  selectedRole.value = role
  const permissionIds = role.permissions?.map(p => p.id) || []
  selectedPermissions.value = [...permissionIds]
  originalPermissions.value = [...permissionIds]
}

// 保存权限配置
const savePermissions = async () => {
  if (!selectedRole.value) return

  saveLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('http://localhost:8081/api/permissions/assign', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        role_id: selectedRole.value.id,
        permission_ids: selectedPermissions.value
      })
    })

    const result = await response.json()
    if (result.code === 200) {
      ElMessage.success('权限配置保存成功')
      originalPermissions.value = [...selectedPermissions.value]
      // 更新角色权限
      if (selectedRole.value) {
        selectedRole.value.permissions = allPermissions.value.filter(p => 
          selectedPermissions.value.includes(p.id)
        )
      }
      getRolesList() // 刷新角色列表
    } else {
      ElMessage.error(result.message || '保存权限配置失败')
    }
  } catch (error) {
    console.error('保存权限配置错误:', error)
    ElMessage.error('网络错误，请稍后重试')
  } finally {
    saveLoading.value = false
  }
}

// 重置权限配置
const resetPermissions = () => {
  selectedPermissions.value = [...originalPermissions.value]
  ElMessage.info('已重置为保存前的状态')
}

// 全选权限
const selectAllPermissions = () => {
  selectedPermissions.value = allPermissions.value.map(p => p.id)
}

// 清空权限
const clearAllPermissions = () => {
  selectedPermissions.value = []
}

// 判断资源组是否全选
const isGroupSelected = (resource: string) => {
  const groupPermissions = groupedPermissions.value[resource] || []
  return groupPermissions.length > 0 && 
         groupPermissions.every(p => selectedPermissions.value.includes(p.id))
}

// 判断资源组是否部分选择
const isGroupIndeterminate = (resource: string) => {
  const groupPermissions = groupedPermissions.value[resource] || []
  const selectedCount = groupPermissions.filter(p => selectedPermissions.value.includes(p.id)).length
  return selectedCount > 0 && selectedCount < groupPermissions.length
}

// 切换资源组选择状态
const toggleGroup = (resource: string, checked: boolean) => {
  const groupPermissions = groupedPermissions.value[resource] || []
  const groupIds = groupPermissions.map(p => p.id)
  
  if (checked) {
    // 添加该组的所有权限
    groupIds.forEach(id => {
      if (!selectedPermissions.value.includes(id)) {
        selectedPermissions.value.push(id)
      }
    })
  } else {
    // 移除该组的所有权限
    selectedPermissions.value = selectedPermissions.value.filter(id => !groupIds.includes(id))
  }
}

// 刷新数据
const refreshData = async () => {
  loading.value = true
  try {
    await Promise.all([getRolesList(), getPermissionsList()])
  } finally {
    loading.value = false
  }
}

// 格式化时间
const formatDateTime = (dateTime: string) => {
  if (!dateTime) return ''
  return new Date(dateTime).toLocaleString('zh-CN')
}

// 获取资源图标
const getResourceIcon = (resource: string) => {
  const iconMap: Record<string, any> = {
    user: User,
    role: UserFilled,
    permission: Key,
    system: Setting,
    auth: Lock
  }
  return iconMap[resource] || Key
}

// 获取资源文本
const getResourceText = (resource: string) => {
  const textMap: Record<string, string> = {
    user: '用户管理',
    role: '角色管理',
    permission: '权限管理',
    system: '系统管理',
    auth: '认证管理'
  }
  return textMap[resource] || resource
}

// 获取操作类型标签样式
const getActionTagType = (action: string) => {
  const typeMap: Record<string, string> = {
    read: 'info',
    write: 'warning',
    delete: 'danger'
  }
  return typeMap[action] || 'info'
}

// 获取操作类型文本
const getActionText = (action: string) => {
  const textMap: Record<string, string> = {
    read: '读取',
    write: '写入',
    delete: '删除'
  }
  return textMap[action] || action
}

// 导出权限
const exportPermissions = () => {
  ElMessage.info('正在导出权限数据...')
  const token = localStorage.getItem('token')
  fetch('/api/export/permissions', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  })
    .then(res => {
      if (!res.ok) throw new Error('导出失败')
      return res.blob()
    })
    .then(blob => {
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `permissions_${new Date().toISOString().slice(0,10)}.csv`
      document.body.appendChild(a)
      a.click()
      a.remove()
      window.URL.revokeObjectURL(url)
      ElMessage.success('权限数据已导出')
    })
    .catch(() => {
      ElMessage.error('导出失败，您可能没有权限')
    })
}

// 提交导入
const submitImport = () => {
  const files = fileUploadRef.value?.getFiles()
  if (!files || files.length === 0) {
    ElMessage.warning('请先选择CSV文件')
    return
  }
  importing.value = true
  const formData = new FormData()
  formData.append('file', files[0].file)
  const token = localStorage.getItem('token')
  fetch('/api/import/permissions', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    },
    body: formData
  })
    .then(res => res.json())
    .then(res => {
      importResult.value = res.data
      if (res.data.success > 0) {
        ElMessage.success(`成功导入 ${res.data.success} 条权限`)
        refreshData()
      }
      if (res.data.failed > 0) {
        ElMessage.warning(`有 ${res.data.failed} 条导入失败`)
      }
    })
    .catch(() => {
      ElMessage.error('导入失败')
    })
    .finally(() => {
      importing.value = false
    })
}

function handleImportSuccess() {}
function handleImportError() {}

// 初始化
onMounted(() => {
  refreshData()
})
</script>

<style scoped>
.main-grid {
  margin-top: var(--spacing-xl);
}

.roles-card,
.permissions-card {
  height: calc(100vh - 280px);
  display: flex;
  flex-direction: column;
}

.card-content {
  flex: 1;
  overflow-y: auto;
}

.role-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.role-item {
  padding: var(--spacing-lg);
  border: 1px solid var(--border-secondary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-normal) ease;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--bg-card);
}

.role-item:hover {
  border-color: var(--primary-color);
  box-shadow: var(--shadow-sm);
}

.role-item.active {
  border-color: var(--primary-color);
  background: var(--bg-selected);
}

.role-info {
  flex: 1;
}

.role-name {
  font-weight: var(--font-medium);
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
}

.role-desc {
  font-size: var(--font-sm);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-sm);
}

.role-meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.permission-count {
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

.role-arrow {
  color: var(--text-tertiary);
  font-size: var(--font-lg);
}

.empty-state {
  text-align: center;
  padding: var(--spacing-2xl);
  color: var(--text-secondary);
}

.empty-state h3 {
  margin: var(--spacing-lg) 0 var(--spacing-sm);
  color: var(--text-primary);
}

.permissions-config {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.permission-stats {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-lg);
  background: var(--bg-secondary);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-lg);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.stat-label {
  font-size: var(--font-sm);
  color: var(--text-secondary);
}

.stat-value {
  font-weight: var(--font-medium);
  color: var(--primary-color);
}

.actions {
  display: flex;
  gap: var(--spacing-xs);
}

.permission-groups {
  flex: 1;
  overflow-y: auto;
}

.permission-group {
  margin-bottom: var(--spacing-xl);
}

.group-header {
  margin-bottom: var(--spacing-lg);
}

.group-title {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-weight: var(--font-medium);
  color: var(--text-primary);
}

.group-icon {
  color: var(--primary-color);
}

.group-content {
  padding-left: var(--spacing-xl);
}

.permission-item {
  margin-bottom: var(--spacing-lg);
  padding: var(--spacing-lg);
  border: 1px solid var(--border-secondary);
  border-radius: var(--radius-md);
  transition: all var(--duration-normal) ease;
}

.permission-item:hover {
  border-color: var(--primary-color);
  background: var(--bg-hover);
}

.permission-info {
  margin-left: var(--spacing-lg);
}

.permission-main {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-xs);
}

.permission-name {
  font-weight: var(--font-medium);
  color: var(--text-primary);
}

.permission-desc {
  font-size: var(--font-sm);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-xs);
}

.permission-code {
  font-size: var(--font-xs);
  color: var(--text-tertiary);
  font-family: monospace;
}

.permission-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.permission-detail {
  padding: var(--spacing-md);
}

/* 动画效果 */
.fade-in {
  animation: fadeIn 0.6s ease-out;
  animation-fill-mode: both;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式调整 */
@media (max-width: 768px) {
  .roles-card,
  .permissions-card {
    height: auto;
    min-height: 400px;
    margin-bottom: var(--spacing-lg);
  }

  .permission-stats {
    flex-direction: column;
    gap: var(--spacing-sm);
    align-items: stretch;
  }

  .actions {
    justify-content: center;
  }
}

.import-result {
  padding: var(--spacing-md);
}

.import-errors {
  list-style: none;
  padding-left: 0;
}
</style> 