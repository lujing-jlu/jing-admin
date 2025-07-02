<template>
  <div class="page-container">
    <!-- 页面头部 -->
    <div class="page-header fade-in">
      <div class="header-info">
        <h1 class="page-title">操作日志</h1>
        <p class="page-subtitle">查看和管理系统操作记录</p>
      </div>
      <div class="header-actions">
        <el-button 
          type="warning" 
          :icon="Delete" 
          @click="showClearDialog"
          size="default"
        >
          清理旧日志
        </el-button>
        <el-button 
          type="primary" 
          :icon="Refresh" 
          @click="refreshData"
          :loading="loading"
          size="default"
        >
          刷新数据
        </el-button>
      </div>
    </div>

    <!-- 统计概览 -->
    <div class="stats-overview fade-in" style="animation-delay: 0.1s">
      <el-row :gutter="20">
        <el-col :xl="6" :lg="6" :md="12" :sm="12" :xs="24">
          <div class="stat-card">
            <div class="stat-icon" style="background-color: #3b82f6;">
              <el-icon :size="24"><Document /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.totalLogs || 0 }}</div>
              <div class="stat-label">总日志数</div>
            </div>
          </div>
        </el-col>
        <el-col :xl="6" :lg="6" :md="12" :sm="12" :xs="24">
          <div class="stat-card">
            <div class="stat-icon" style="background-color: #10b981;">
              <el-icon :size="24"><Clock /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.todayLogs || 0 }}</div>
              <div class="stat-label">今日日志</div>
            </div>
          </div>
        </el-col>
        <el-col :xl="6" :lg="6" :md="12" :sm="12" :xs="24">
          <div class="stat-card">
            <div class="stat-icon" style="background-color: #f59e0b;">
              <el-icon :size="24"><Calendar /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stats.weekLogs || 0 }}</div>
              <div class="stat-label">本周日志</div>
            </div>
          </div>
        </el-col>
        <el-col :xl="6" :lg="6" :md="12" :sm="12" :xs="24">
          <div class="stat-card">
            <div class="stat-icon" style="background-color: #6366f1;">
              <el-icon :size="24"><User /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ activeUsers || 0 }}</div>
              <div class="stat-label">活跃用户</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索筛选 -->
    <div class="card-container search-section fade-in" style="animation-delay: 0.2s">
      <div class="search-form">
        <el-form :model="searchForm" :inline="true">
          <el-form-item label="用户名">
            <el-input 
              v-model="searchForm.username" 
              placeholder="请输入用户名"
              clearable
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item label="操作类型">
            <el-select 
              v-model="searchForm.action" 
              placeholder="请选择操作类型"
              clearable
              style="width: 150px"
            >
              <el-option label="创建" value="create" />
              <el-option label="读取" value="read" />
              <el-option label="更新" value="update" />
              <el-option label="删除" value="delete" />
            </el-select>
          </el-form-item>
          <el-form-item label="资源类型">
            <el-select 
              v-model="searchForm.resource" 
              placeholder="请选择资源类型"
              clearable
              style="width: 150px"
            >
              <el-option label="用户" value="user" />
              <el-option label="角色" value="role" />
              <el-option label="权限" value="permission" />
              <el-option label="系统" value="system" />
              <el-option label="认证" value="auth" />
            </el-select>
          </el-form-item>
          <el-form-item label="时间范围">
            <el-date-picker
              v-model="searchForm.dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              style="width: 250px"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
            <el-button :icon="RefreshLeft" @click="resetSearch">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <!-- 日志表格 -->
    <div class="card-container table-section fade-in" style="animation-delay: 0.3s">
      <div class="card-header">
        <div class="card-title">
          <el-icon><List /></el-icon>
          <span>操作日志列表</span>
        </div>
        <div class="table-actions">
          <el-button
            v-if="selectedLogs.length > 0"
            type="danger"
            :icon="Delete"
            @click="batchDelete"
            size="small"
          >
            批量删除 ({{ selectedLogs.length }})
          </el-button>
        </div>
      </div>

      <el-table
        :data="logs"
        v-loading="loading"
        @selection-change="handleSelectionChange"
        row-key="id"
        stripe
        border
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户" width="120" />
        <el-table-column prop="action" label="操作" width="80">
          <template #default="{ row }">
            <el-tag 
              :type="getActionTagType(row.action)"
              size="small"
            >
              {{ getActionText(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource" label="资源" width="100">
          <template #default="{ row }">
            <el-tag 
              :type="getResourceTagType(row.resource)"
              size="small"
              effect="plain"
            >
              {{ getResourceText(row.resource) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="method" label="方法" width="80" />
        <el-table-column prop="path" label="路径" min-width="200" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP地址" width="120" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag 
              :type="getStatusTagType(row.status)"
              size="small"
            >
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="操作时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              :icon="View"
              @click="viewDetail(row)"
              size="small"
              link
            >
              详情
            </el-button>
            <el-button
              type="danger"
              :icon="Delete"
              @click="deleteLog(row)"
              size="small"
              link
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-section">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 日志详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="操作日志详情"
      width="800px"
      destroy-on-close
    >
      <div v-if="selectedLog" class="log-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="日志ID">{{ selectedLog.id }}</el-descriptions-item>
          <el-descriptions-item label="用户">{{ selectedLog.username }}</el-descriptions-item>
          <el-descriptions-item label="操作类型">
            <el-tag :type="getActionTagType(selectedLog.action)">
              {{ getActionText(selectedLog.action) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="资源类型">
            <el-tag :type="getResourceTagType(selectedLog.resource)" effect="plain">
              {{ getResourceText(selectedLog.resource) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="HTTP方法">{{ selectedLog.method }}</el-descriptions-item>
          <el-descriptions-item label="状态码">
            <el-tag :type="getStatusTagType(selectedLog.status)">
              {{ selectedLog.status }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="请求路径" :span="2">{{ selectedLog.path }}</el-descriptions-item>
          <el-descriptions-item label="IP地址">{{ selectedLog.ip }}</el-descriptions-item>
          <el-descriptions-item label="操作时间">{{ formatDateTime(selectedLog.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="用户代理" :span="2">
            <div class="user-agent">{{ selectedLog.user_agent }}</div>
          </el-descriptions-item>
          <el-descriptions-item v-if="selectedLog.details" label="详细信息" :span="2">
            <div class="details">{{ selectedLog.details }}</div>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <!-- 清理确认对话框 -->
    <el-dialog
      v-model="clearDialogVisible"
      title="清理旧日志"
      width="400px"
    >
      <p>确定要清理30天前的旧日志吗？此操作不可撤销。</p>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="clearDialogVisible = false">取消</el-button>
          <el-button type="warning" @click="clearOldLogs" :loading="clearLoading">
            确定清理
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Document,
  Clock,
  Calendar,
  User,
  Search,
  RefreshLeft,
  List,
  View,
  Delete,
  Refresh
} from '@element-plus/icons-vue'

// 类型定义
interface OperationLog {
  id: number
  user_id: number
  username: string
  action: string
  resource: string
  resource_id?: string
  method: string
  path: string
  ip: string
  user_agent: string
  status: number
  details?: string
  created_at: string
}

interface LogStats {
  totalLogs: number
  todayLogs: number
  weekLogs: number
  actionStats: Array<{ action: string; count: number }>
  resourceStats: Array<{ resource: string; count: number }>
  userStats: Array<{ username: string; count: number }>
}

// 响应式数据
const loading = ref(false)
const clearLoading = ref(false)
const logs = ref<OperationLog[]>([])
const stats = ref<LogStats>({} as LogStats)
const selectedLogs = ref<OperationLog[]>([])
const detailDialogVisible = ref(false)
const clearDialogVisible = ref(false)
const selectedLog = ref<OperationLog | null>(null)

// 搜索表单
const searchForm = ref({
  username: '',
  action: '',
  resource: '',
  dateRange: [] as string[]
})

// 分页
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

// 计算属性
const activeUsers = computed(() => {
  if (!stats.value.userStats || !Array.isArray(stats.value.userStats)) {
    return 0
  }
  return stats.value.userStats.length
})

// 获取操作日志列表
const getLogsList = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      page: pagination.value.page.toString(),
      page_size: pagination.value.pageSize.toString()
    })

    // 添加搜索条件
    if (searchForm.value.username) {
      params.append('username', searchForm.value.username)
    }
    if (searchForm.value.action) {
      params.append('action', searchForm.value.action)
    }
    if (searchForm.value.resource) {
      params.append('resource', searchForm.value.resource)
    }
    if (searchForm.value.dateRange && searchForm.value.dateRange.length === 2) {
      params.append('start_date', searchForm.value.dateRange[0])
      params.append('end_date', searchForm.value.dateRange[1])
    }

    const token = localStorage.getItem('token')
    const response = await fetch(`http://localhost:8081/api/logs?${params}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    const result = await response.json()
    if (result.code === 200) {
      logs.value = result.data.logs || []
      pagination.value.total = result.data.total || 0
    } else {
      ElMessage.error(result.message || '获取日志列表失败')
    }
  } catch (error) {
    console.error('获取日志列表错误:', error)
    ElMessage.error('网络错误，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 获取统计信息
const getStats = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('http://localhost:8081/api/logs/stats', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    const result = await response.json()
    if (result.code === 200) {
      stats.value = result.data
    }
  } catch (error) {
    console.error('获取统计信息错误:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.value.page = 1
  getLogsList()
}

// 重置搜索
const resetSearch = () => {
  searchForm.value = {
    username: '',
    action: '',
    resource: '',
    dateRange: []
  }
  pagination.value.page = 1
  getLogsList()
}

// 分页处理
const handleCurrentChange = (page: number) => {
  pagination.value.page = page
  getLogsList()
}

const handleSizeChange = (size: number) => {
  pagination.value.pageSize = size
  pagination.value.page = 1
  getLogsList()
}

// 选择处理
const handleSelectionChange = (selection: OperationLog[]) => {
  selectedLogs.value = selection
}

// 查看详情
const viewDetail = (log: OperationLog) => {
  selectedLog.value = log
  detailDialogVisible.value = true
}

// 删除单条日志
const deleteLog = async (log: OperationLog) => {
  try {
    await ElMessageBox.confirm('确定要删除这条日志吗？', '删除确认', {
      type: 'warning'
    })

    const token = localStorage.getItem('token')
    const response = await fetch(`http://localhost:8081/api/logs/${log.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    const result = await response.json()
    if (result.code === 200) {
      ElMessage.success('删除成功')
      getLogsList()
      getStats()
    } else {
      ElMessage.error(result.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除日志错误:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 批量删除
const batchDelete = async () => {
  if (selectedLogs.value.length === 0) {
    ElMessage.warning('请选择要删除的日志')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedLogs.value.length} 条日志吗？`,
      '批量删除确认',
      { type: 'warning' }
    )

    const token = localStorage.getItem('token')
    const response = await fetch('http://localhost:8081/api/logs/batch-delete', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        ids: selectedLogs.value.map(log => log.id)
      })
    })

    const result = await response.json()
    if (result.code === 200) {
      ElMessage.success(`成功删除 ${result.data.deleted} 条日志`)
      selectedLogs.value = []
      getLogsList()
      getStats()
    } else {
      ElMessage.error(result.message || '批量删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除错误:', error)
      ElMessage.error('批量删除失败')
    }
  }
}

// 显示清理对话框
const showClearDialog = () => {
  clearDialogVisible.value = true
}

// 清理旧日志
const clearOldLogs = async () => {
  clearLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('http://localhost:8081/api/logs/clear-old', {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    const result = await response.json()
    if (result.code === 200) {
      ElMessage.success(`成功清理 ${result.data.deleted} 条旧日志`)
      clearDialogVisible.value = false
      getLogsList()
      getStats()
    } else {
      ElMessage.error(result.message || '清理失败')
    }
  } catch (error) {
    console.error('清理旧日志错误:', error)
    ElMessage.error('清理失败')
  } finally {
    clearLoading.value = false
  }
}

// 刷新数据
const refreshData = () => {
  getLogsList()
  getStats()
}

// 格式化时间
const formatDateTime = (dateTime: string) => {
  if (!dateTime) return ''
  return new Date(dateTime).toLocaleString('zh-CN')
}

// 获取操作类型标签样式
const getActionTagType = (action: string) => {
  const typeMap: Record<string, string> = {
    create: 'success',
    read: 'info',
    update: 'warning',
    delete: 'danger'
  }
  return typeMap[action] || 'info'
}

// 获取操作类型文本
const getActionText = (action: string) => {
  const textMap: Record<string, string> = {
    create: '创建',
    read: '读取',
    update: '更新',
    delete: '删除'
  }
  return textMap[action] || action
}

// 获取资源类型标签样式
const getResourceTagType = (resource: string) => {
  const typeMap: Record<string, string> = {
    user: 'primary',
    role: 'success',
    permission: 'warning',
    system: 'info',
    auth: 'danger'
  }
  return typeMap[resource] || 'info'
}

// 获取资源类型文本
const getResourceText = (resource: string) => {
  const textMap: Record<string, string> = {
    user: '用户',
    role: '角色',
    permission: '权限',
    system: '系统',
    auth: '认证'
  }
  return textMap[resource] || resource
}

// 获取状态码标签样式
const getStatusTagType = (status: number) => {
  if (status >= 200 && status < 300) return 'success'
  if (status >= 400 && status < 500) return 'warning'
  if (status >= 500) return 'danger'
  return 'info'
}

// 初始化
onMounted(() => {
  getLogsList()
  getStats()
})
</script>

<style scoped>
.stats-overview {
  margin-bottom: var(--spacing-xl);
}

.stat-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-xl);
  border: 1px solid var(--border-secondary);
  backdrop-filter: blur(10px);
  transition: all var(--duration-normal) ease;
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
}

.stat-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: var(--font-2xl);
  font-weight: var(--font-bold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
}

.stat-label {
  font-size: var(--font-sm);
  color: var(--text-secondary);
}

.table-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.log-detail {
  padding: var(--spacing-md);
}

.user-agent {
  word-break: break-all;
  font-family: monospace;
  font-size: var(--font-sm);
  color: var(--text-secondary);
}

.details {
  white-space: pre-wrap;
  font-family: monospace;
  font-size: var(--font-sm);
  color: var(--text-secondary);
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


</style> 