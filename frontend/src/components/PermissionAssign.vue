<template>
  <el-dialog
    :model-value="visible"
    title="权限分配"
    width="600px"
    @update:model-value="handleClose"
    @close="handleClose"
  >
    <div v-if="role" class="permission-assign">
      <div class="role-info">
        <h4>角色：{{ role.display_name }}</h4>
        <p class="role-desc">{{ role.description }}</p>
      </div>

      <el-divider />

      <div class="permissions-section">
        <h4>权限配置</h4>
        <div v-loading="permissionLoading" class="permissions-grid">
          <div
            v-for="resource in groupedPermissions"
            :key="resource.name"
            class="resource-group"
          >
            <div class="resource-header">
              <h5>{{ getResourceDisplayName(resource.name) }}</h5>
              <el-checkbox
                :model-value="isAllResourceSelected(resource.name)"
                :indeterminate="isResourceIndeterminate(resource.name)"
                @change="handleResourceSelectAll(resource.name, $event)"
              >
                全选
              </el-checkbox>
            </div>
            <div class="permissions-list">
              <el-checkbox
                v-for="permission in resource.permissions"
                :key="permission.id"
                :model-value="selectedPermissions.includes(permission.id)"
                @change="handlePermissionChange(permission.id, $event)"
              >
                {{ permission.display_name }}
                <span class="permission-desc">{{ permission.description }}</span>
              </el-checkbox>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          保存
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, type PropType } from 'vue'
import { ElMessage } from 'element-plus'
import type { Role, Permission } from '../stores/role'
import { useRoleStore } from '../stores/role'

// Props
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  role: {
    type: Object as PropType<Role | null>,
    default: null
  },
  loading: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits<{
  'update:visible': [value: boolean]
  'submit': [data: { roleId: number; permissionIds: number[] }]
}>()

const roleStore = useRoleStore()

// 状态
const selectedPermissions = ref<number[]>([])
const permissionLoading = ref(false)

// 计算属性
const groupedPermissions = computed(() => {
  const grouped = new Map<string, { name: string; permissions: Permission[] }>()
  
  roleStore.permissions.forEach(permission => {
    if (!grouped.has(permission.resource)) {
      grouped.set(permission.resource, {
        name: permission.resource,
        permissions: []
      })
    }
    grouped.get(permission.resource)!.permissions.push(permission)
  })
  
  return Array.from(grouped.values())
})

// 生命周期
onMounted(async () => {
  if (roleStore.permissions.length === 0) {
    permissionLoading.value = true
    try {
      await roleStore.fetchPermissions()
    } catch (error) {
      ElMessage.error('加载权限列表失败')
    } finally {
      permissionLoading.value = false
    }
  }
})

// 监听角色变化，初始化已选权限
watch(() => props.role, (newRole) => {
  if (newRole && newRole.permissions) {
    selectedPermissions.value = newRole.permissions.map(p => p.id)
  } else {
    selectedPermissions.value = []
  }
}, { immediate: true })

// 监听弹窗显示状态
watch(() => props.visible, (visible) => {
  if (visible && props.role) {
    // 重新初始化选中状态
    selectedPermissions.value = props.role.permissions?.map(p => p.id) || []
  }
})

// 方法
function getResourceDisplayName(resource: string): string {
  const resourceMap: Record<string, string> = {
    user: '用户管理',
    role: '角色管理',
    permission: '权限管理',
    system: '系统管理'
  }
  return resourceMap[resource] || resource
}

function isAllResourceSelected(resource: string): boolean {
  const resourcePermissions = groupedPermissions.value.find(g => g.name === resource)?.permissions || []
  return resourcePermissions.length > 0 && resourcePermissions.every(p => selectedPermissions.value.includes(p.id))
}

function isResourceIndeterminate(resource: string): boolean {
  const resourcePermissions = groupedPermissions.value.find(g => g.name === resource)?.permissions || []
  const selectedCount = resourcePermissions.filter(p => selectedPermissions.value.includes(p.id)).length
  return selectedCount > 0 && selectedCount < resourcePermissions.length
}

function handleResourceSelectAll(resource: string, checked: boolean) {
  const resourcePermissions = groupedPermissions.value.find(g => g.name === resource)?.permissions || []
  
  if (checked) {
    // 添加该资源下所有权限
    resourcePermissions.forEach(permission => {
      if (!selectedPermissions.value.includes(permission.id)) {
        selectedPermissions.value.push(permission.id)
      }
    })
  } else {
    // 移除该资源下所有权限
    selectedPermissions.value = selectedPermissions.value.filter(id => 
      !resourcePermissions.some(p => p.id === id)
    )
  }
}

function handlePermissionChange(permissionId: number, checked: boolean) {
  if (checked) {
    if (!selectedPermissions.value.includes(permissionId)) {
      selectedPermissions.value.push(permissionId)
    }
  } else {
    selectedPermissions.value = selectedPermissions.value.filter(id => id !== permissionId)
  }
}

function handleClose() {
  emit('update:visible', false)
}

function handleSubmit() {
  if (!props.role) return
  
  emit('submit', {
    roleId: props.role.id,
    permissionIds: selectedPermissions.value
  })
}
</script>

<style scoped>
.permission-assign {
  max-height: 60vh;
  overflow-y: auto;
}

.role-info h4 {
  margin: 0 0 8px 0;
  color: var(--text-primary);
}

.role-desc {
  margin: 0;
  color: var(--text-secondary);
  font-size: 14px;
}

.permissions-section h4 {
  margin: 0 0 16px 0;
  color: var(--text-primary);
}

.permissions-grid {
  display: grid;
  gap: 20px;
}

.resource-group {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 16px;
  background: var(--card-bg);
}

.resource-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-color);
}

.resource-header h5 {
  margin: 0;
  color: var(--text-primary);
  font-weight: 600;
}

.permissions-list {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 8px;
}

.permissions-list :deep(.el-checkbox) {
  margin-right: 0;
  margin-bottom: 8px;
  flex-direction: column;
  align-items: flex-start;
}

.permissions-list :deep(.el-checkbox__label) {
  padding-left: 0;
  font-size: 14px;
  line-height: 1.4;
}

.permission-desc {
  display: block;
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 2px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 