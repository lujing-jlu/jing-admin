<template>
  <div class="page-container">
    <!-- 页面头部 -->
    <div class="page-header fade-in">
      <div>
        <h2 class="page-title">用户管理</h2>
        <p class="page-subtitle">管理系统用户账户和权限设置</p>
      </div>
      <div class="page-actions">
        <el-button type="primary" :icon="Plus" @click="handleAdd">
          新增用户
        </el-button>
        <el-button :icon="Download" @click="exportUsers" style="margin-left: 8px;">
          导出用户
        </el-button>
        <el-button :icon="Upload" @click="importDialogVisible = true" style="margin-left: 8px;">
          导入用户
        </el-button>
      </div>
    </div>

    <!-- 搜索区域 -->
    <div class="search-section fade-in" style="animation-delay: 0.1s">
      <el-form :inline="true" class="search-form">
        <el-form-item label="用户名">
          <el-input 
            v-model="searchForm.username" 
            placeholder="请输入用户名" 
            clearable 
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="角色">
          <el-select 
            v-model="searchForm.role" 
            placeholder="请选择角色" 
            clearable
            style="width: 150px"
          >
            <el-option label="管理员" value="admin" />
            <el-option label="普通用户" value="user" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select 
            v-model="searchForm.status" 
            placeholder="请选择状态" 
            clearable
            style="width: 120px"
          >
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">
            搜索
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 数据表格 -->
    <div class="table-section fade-in" style="animation-delay: 0.2s">
      <el-table 
        :data="userStore.users" 
        :loading="userStore.loading"
        stripe
        style="width: 100%"
        empty-text="暂无用户数据"
        v-loading="userStore.loading"
        element-loading-text="正在加载用户数据..."
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="username" label="用户名" min-width="120">
          <template #default="{ row }">
            <div style="display: flex; align-items: center; gap: 8px;">
              <el-avatar :size="32" :src="`https://api.dicebear.com/7.x/initials/svg?seed=${row.username}`" />
              <span style="font-weight: 500;">{{ row.username }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="180">
          <template #default="{ row }">
            <span style="color: var(--text-secondary);">{{ row.email }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="role" label="角色" width="120" align="center">
          <template #default="{ row }">
            <el-tag 
              :type="row.role === 'admin' ? 'danger' : row.role === 'super_admin' ? 'warning' : 'primary'"
              size="small"
            >
              {{ getRoleDisplayName(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status ? 'success' : 'danger'" size="small">
              {{ row.status ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180" align="center">
          <template #default="{ row }">
            <span style="color: var(--text-tertiary); font-size: 13px;">
              {{ formatDate(row.createdAt) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <div style="display: flex; justify-content: center; gap: 8px;">
              <el-button 
                size="small" 
                type="primary" 
                :icon="Edit"
                @click="handleEdit(row)"
                plain
              >
                编辑
              </el-button>
              <el-button 
                size="small" 
                type="danger" 
                :icon="Delete"
                @click="handleDelete(row)"
                plain
                :disabled="row.username === 'admin'"
              >
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-section">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="userStore.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 错误提示 -->
    <el-alert
      v-if="userStore.error"
      :title="userStore.error"
      type="error"
      :closable="false"
      style="margin-top: 20px"
      class="fade-in"
    />

    <!-- 用户表单弹窗 -->
    <UserForm
      v-model:visible="formVisible"
      :mode="formMode"
      :user-data="currentUser"
      :loading="userStore.loading"
      @submit="handleFormSubmit"
    />

    <!-- 用户导入弹窗 -->
    <el-dialog v-model="importDialogVisible" title="批量导入用户" width="500px" center>
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
import { Plus, Search, Edit, Delete, Download, Upload } from '@element-plus/icons-vue'
import { useUserStore, type User } from '../../stores/user'
import UserForm from '../../components/UserForm.vue'
import FileUpload from '../../components/FileUpload.vue'

const userStore = useUserStore()

// 搜索表单
const searchForm = ref({
  username: '',
  role: '',
  status: ''
})

// 分页
const currentPage = ref(1)
const pageSize = ref(10)

// 表单相关
const formVisible = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const currentUser = ref<User | null>(null)

// 导入相关
const importDialogVisible = ref(false)
const importing = ref(false)
const importResult = ref<any>(null)
const fileUploadRef = ref()

// 生命周期
onMounted(() => {
  loadUsers()
})

// 方法
async function loadUsers() {
  try {
    await userStore.fetchUsers()
  } catch (error) {
    ElMessage.error('加载用户列表失败')
  }
}

function handleAdd() {
  formMode.value = 'create'
  currentUser.value = null
  formVisible.value = true
}

function handleEdit(row: User) {
  formMode.value = 'edit'
  currentUser.value = row
  formVisible.value = true
}

async function handleDelete(row: any) {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${row.username}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await userStore.deleteUser(row.id)
    ElMessage.success('删除成功')
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
      await userStore.createUser(formData)
      ElMessage.success('用户创建成功')
    } else {
      await userStore.updateUser(currentUser.value!.id, formData)
      ElMessage.success('用户更新成功')
    }
    formVisible.value = false
    await loadUsers() // 重新加载用户列表
  } catch (error) {
    ElMessage.error(formMode.value === 'create' ? '创建用户失败' : '更新用户失败')
  }
}

function handleSearch() {
  // 实现搜索功能
  loadUsers()
}

function handleReset() {
  searchForm.value = {
    username: '',
    role: '',
    status: ''
  }
  loadUsers()
}

function handleSizeChange(size: number) {
  pageSize.value = size
  loadUsers()
}

function handleCurrentChange(page: number) {
  currentPage.value = page
  loadUsers()
}

function getRoleDisplayName(role: string) {
  const roleMap: Record<string, string> = {
    'super_admin': '超级管理员',
    'admin': '管理员',
    'user': '普通用户'
  }
  return roleMap[role] || role
}

function formatDate(dateString: string) {
  if (!dateString) return '--'
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return '--'
  }
}

function exportUsers() {
  ElMessage.info('正在导出用户数据...')
  const token = localStorage.getItem('token')
  fetch('/api/export/users', {
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
      a.download = `users_${new Date().toISOString().slice(0,10)}.csv`
      document.body.appendChild(a)
      a.click()
      a.remove()
      window.URL.revokeObjectURL(url)
      ElMessage.success('用户数据已导出')
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
  fetch('/api/import/users', {
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
        ElMessage.success(`成功导入 ${res.data.success} 条用户`)
        loadUsers()
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
/* 页面特定样式可以在这里添加，但要使用设计系统的变量 */
.search-form .el-form-item {
  margin-bottom: 0;
}

.table-section .el-table {
  margin-bottom: 0;
}
</style> 