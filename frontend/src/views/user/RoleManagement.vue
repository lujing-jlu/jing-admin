<template>
  <div class="page-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <h2 class="page-title">角色管理</h2>
      <div class="page-actions">
        <el-button type="primary" :icon="Plus" @click="handleAdd">
          新增角色
        </el-button>
        <el-button :icon="Download" @click="exportRoles" style="margin-left: 8px;">
          导出角色
        </el-button>
        <el-button :icon="Upload" @click="importDialogVisible = true" style="margin-left: 8px;">
          导入角色
        </el-button>
      </div>
    </div>

    <!-- 角色列表 -->
    <div class="table-section">
      <el-table 
        :data="roleStore.roles" 
        :loading="roleStore.loading"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="display_name" label="角色名称" />
        <el-table-column prop="name" label="角色标识" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status ? 'success' : 'danger'">
              {{ row.status ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="权限数量" width="100">
          <template #default="{ row }">
            <el-tag type="info">{{ row.permissions?.length || 0 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="CreatedAt" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button size="small" type="warning" @click="handlePermissions(row)">
              权限
            </el-button>
            <el-button 
              size="small" 
              type="danger" 
              @click="handleDelete(row)"
              :disabled="isSystemRole(row.name)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 错误提示 -->
    <el-alert
      v-if="roleStore.error"
      :title="roleStore.error"
      type="error"
      :closable="false"
      style="margin-top: 20px"
    />

    <!-- 角色表单弹窗 -->
    <RoleForm
      v-model:visible="formVisible"
      :mode="formMode"
      :role-data="currentRole"
      :loading="roleStore.loading"
      @submit="handleFormSubmit"
    />

    <!-- 权限分配弹窗 -->
    <PermissionAssign
      v-model:visible="permissionVisible"
      :role="currentRole"
      :loading="roleStore.loading"
      @submit="handlePermissionSubmit"
    />

    <!-- 角色导入弹窗 -->
    <el-dialog v-model="importDialogVisible" title="批量导入角色" width="500px" center>
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
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Download, Upload } from '@element-plus/icons-vue'
import { useRoleStore, type Role } from '../../stores/role'
import RoleForm from '../../components/RoleForm.vue'
import PermissionAssign from '../../components/PermissionAssign.vue'
import FileUpload from '../../components/FileUpload.vue'

const roleStore = useRoleStore()

// 表单相关
const formVisible = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const currentRole = ref<Role | null>(null)

// 权限分配相关
const permissionVisible = ref(false)

// 导入相关
const importDialogVisible = ref(false)
const importing = ref(false)
const importResult = ref<any>(null)
const fileUploadRef = ref()

// 生命周期
onMounted(() => {
  loadRoles()
})

// 方法
async function loadRoles() {
  try {
    await roleStore.fetchRoles()
  } catch (error) {
    ElMessage.error('加载角色列表失败')
  }
}

function handleAdd() {
  formMode.value = 'create'
  currentRole.value = null
  formVisible.value = true
}

function handleEdit(row: Role) {
  formMode.value = 'edit'
  currentRole.value = row
  formVisible.value = true
}

function handlePermissions(row: Role) {
  currentRole.value = row
  permissionVisible.value = true
}

async function handleDelete(row: Role) {
  if (isSystemRole(row.name)) {
    ElMessage.warning('不能删除系统默认角色')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除角色 "${row.display_name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await roleStore.deleteRole(row.id)
    ElMessage.success('删除成功')
    await loadRoles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 表单提交处理
async function handleFormSubmit(formData: any) {
  try {
    if (formMode.value === 'create') {
      await roleStore.createRole(formData)
      ElMessage.success('角色创建成功')
    } else {
      await roleStore.updateRole(currentRole.value!.id, formData)
      ElMessage.success('角色更新成功')
    }
    formVisible.value = false
    await loadRoles()
  } catch (error) {
    ElMessage.error(formMode.value === 'create' ? '创建角色失败' : '更新角色失败')
  }
}

// 权限分配提交处理
async function handlePermissionSubmit(permissionData: any) {
  try {
    await roleStore.assignPermissions(permissionData.roleId, permissionData.permissionIds)
    ElMessage.success('权限分配成功')
    permissionVisible.value = false
    await loadRoles()
  } catch (error) {
    ElMessage.error('权限分配失败')
  }
}

function isSystemRole(roleName: string): boolean {
  return ['super_admin', 'admin', 'user'].includes(roleName)
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleString('zh-CN')
}

function exportRoles() {
  ElMessage.info('正在导出角色数据...')
  const token = localStorage.getItem('token')
  fetch('/api/export/roles', {
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
      a.download = `roles_${new Date().toISOString().slice(0,10)}.csv`
      document.body.appendChild(a)
      a.click()
      a.remove()
      window.URL.revokeObjectURL(url)
      ElMessage.success('角色数据已导出')
    })
    .catch(() => {
      ElMessage.error('导出失败，您可能没有权限')
    })
}

function submitImport() {
  const files = fileUploadRef.value?.getFiles()
  if (!files || files.length === 0) {
    ElMessage.warning('请先选择CSV文件')
    return
  }
  importing.value = true
  const formData = new FormData()
  formData.append('file', files[0].file)
  const token = localStorage.getItem('token')
  fetch('/api/import/roles', {
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
        ElMessage.success(`成功导入 ${res.data.success} 条角色`)
        loadRoles()
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
</script>

<style scoped>
/* 页面特定样式，现在只需要很少的自定义样式 */
.el-alert {
  margin-top: var(--spacing-xl);
}
</style> 