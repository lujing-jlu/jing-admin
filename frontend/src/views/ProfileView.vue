<template>
  <div class="page-container">
    <!-- 页面头部 -->
    <div class="page-header fade-in">
      <div class="header-info">
        <h1 class="page-title">个人中心</h1>
        <p class="page-subtitle">管理您的个人信息和账户设置</p>
      </div>
      <div class="header-actions">
        <el-button 
          type="primary" 
          :icon="Edit" 
          @click="editMode = !editMode"
          size="default"
        >
          {{ editMode ? '取消编辑' : '编辑资料' }}
        </el-button>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-grid fade-in" style="animation-delay: 0.1s">
      <el-row :gutter="20">
        <!-- 左侧：用户信息卡片 -->
        <el-col :xl="8" :lg="10" :md="24" :sm="24" :xs="24">
          <div class="card-container profile-card">
            <div class="card-header">
              <div class="card-title">
                <el-icon><User /></el-icon>
                <span>基本信息</span>
              </div>
            </div>
            <div class="card-content">
              <div class="profile-avatar">
                <div class="avatar-upload-wrapper" @click="triggerFileInput" title="点击更换头像">
                  <el-avatar 
                    :size="80" 
                    :src="userInfo.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'"
                    class="user-avatar clickable"
                  />
                  <div class="avatar-actions">
                    <el-button 
                      type="primary" 
                      :icon="Camera" 
                      circle 
                      size="small"
                      @click.stop="triggerFileInput"
                      title="更换头像"
                    />
                  </div>
                  <div class="avatar-hover-mask">点击更换头像</div>
                </div>
                <input ref="fileInput" type="file" accept="image/*" style="display:none" @change="onFileInputChange" />
              </div>
              
              <div class="profile-info">
                <div class="info-item">
                  <span class="info-label">用户名</span>
                  <span class="info-value">{{ userInfo.username }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">角色</span>
                  <el-tag 
                    :type="userInfo.role === 'admin' ? 'danger' : 'primary'"
                    size="small"
                  >
                    {{ getRoleText(userInfo.role) }}
                  </el-tag>
                </div>
                <div class="info-item">
                  <span class="info-label">状态</span>
                  <el-tag 
                    :type="userInfo.status ? 'success' : 'danger'"
                    size="small"
                  >
                    {{ userInfo.status ? '正常' : '禁用' }}
                  </el-tag>
                </div>
                <div class="info-item">
                  <span class="info-label">注册时间</span>
                  <span class="info-value">{{ formatDateTime(userInfo.created_at) }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 快速操作 -->
          <div class="card-container quick-actions-card">
            <div class="card-header">
              <div class="card-title">
                <el-icon><Operation /></el-icon>
                <span>快速操作</span>
              </div>
            </div>
            <div class="card-content">
              <div class="quick-actions">
                <div class="action-item" @click="showChangePasswordDialog">
                  <div class="action-icon" style="background-color: #f59e0b;">
                    <el-icon><Lock /></el-icon>
                  </div>
                  <div class="action-text">
                    <h4>修改密码</h4>
                    <p>更改您的登录密码</p>
                  </div>
                </div>
                <div class="action-item" @click="viewPermissions">
                  <div class="action-icon" style="background-color: #6366f1;">
                    <el-icon><Key /></el-icon>
                  </div>
                  <div class="action-text">
                    <h4>我的权限</h4>
                    <p>查看您的权限列表</p>
                  </div>
                </div>
                <div class="action-item" @click="viewLogs">
                  <div class="action-icon" style="background-color: #10b981;">
                    <el-icon><Document /></el-icon>
                  </div>
                  <div class="action-text">
                    <h4>操作记录</h4>
                    <p>查看您的操作日志</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-col>

        <!-- 右侧：详细设置 -->
        <el-col :xl="16" :lg="14" :md="24" :sm="24" :xs="24">
          <div class="card-container settings-card">
            <div class="card-header">
              <div class="card-title">
                <el-icon><Setting /></el-icon>
                <span>详细信息</span>
              </div>
              <div v-if="editMode" class="settings-actions">
                <el-button 
                  type="success" 
                  :icon="Check" 
                  @click="saveProfile"
                  :loading="saveLoading"
                  size="small"
                >
                  保存更改
                </el-button>
                <el-button 
                  type="warning" 
                  :icon="RefreshLeft" 
                  @click="resetForm"
                  size="small"
                >
                  重置
                </el-button>
              </div>
            </div>
            <div class="card-content">
              <el-form 
                :model="profileForm" 
                :rules="profileRules"
                ref="profileFormRef"
                label-width="100px"
                :disabled="!editMode"
              >
                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="用户名" prop="username">
                      <el-input 
                        v-model="profileForm.username" 
                        disabled
                        placeholder="用户名不可修改"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="邮箱" prop="email">
                      <el-input 
                        v-model="profileForm.email" 
                        placeholder="请输入邮箱"
                        type="email"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="真实姓名" prop="realName">
                      <el-input 
                        v-model="profileForm.realName" 
                        placeholder="请输入真实姓名"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="手机号码" prop="phone">
                      <el-input 
                        v-model="profileForm.phone" 
                        placeholder="请输入手机号码"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="部门" prop="department">
                      <el-input 
                        v-model="profileForm.department" 
                        placeholder="请输入部门"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="职位" prop="position">
                      <el-input 
                        v-model="profileForm.position" 
                        placeholder="请输入职位"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-form-item label="个人简介" prop="bio">
                  <el-input 
                    v-model="profileForm.bio" 
                    type="textarea" 
                    :rows="4"
                    placeholder="请输入个人简介"
                  />
                </el-form-item>
              </el-form>
            </div>
          </div>

          <!-- 安全设置 -->
          <div class="card-container security-card">
            <div class="card-header">
              <div class="card-title">
                <el-icon><Lock /></el-icon>
                <span>安全设置</span>
              </div>
            </div>
            <div class="card-content">
              <div class="security-items">
                <div class="security-item">
                  <div class="security-info">
                    <h4>登录密码</h4>
                    <p>定期更换密码可以提高账户安全性</p>
                  </div>
                  <el-button 
                    type="primary" 
                    @click="showChangePasswordDialog"
                    size="small"
                  >
                    修改密码
                  </el-button>
                </div>
                <div class="security-item">
                  <div class="security-info">
                    <h4>最后登录</h4>
                    <p>{{ userInfo.last_login ? formatDateTime(userInfo.last_login) : '首次登录' }}</p>
                  </div>
                  <el-tag type="info" size="small">安全</el-tag>
                </div>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="changePasswordVisible"
      title="修改密码"
      width="500px"
      destroy-on-close
    >
      <el-form 
        :model="passwordForm" 
        :rules="passwordRules"
        ref="passwordFormRef"
        label-width="100px"
      >
        <el-form-item label="当前密码" prop="currentPassword">
          <el-input 
            v-model="passwordForm.currentPassword" 
            type="password"
            placeholder="请输入当前密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input 
            v-model="passwordForm.newPassword" 
            type="password"
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input 
            v-model="passwordForm.confirmPassword" 
            type="password"
            placeholder="请再次输入新密码"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="changePasswordVisible = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="changePassword"
            :loading="changePasswordLoading"
          >
            确定修改
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 我的权限对话框 -->
    <el-dialog
      v-model="permissionsVisible"
      title="我的权限"
      width="600px"
      destroy-on-close
    >
      <div class="permissions-list">
        <div v-if="myPermissions.length === 0" class="empty-permissions">
          <el-icon :size="48" color="#cbd5e1"><Key /></el-icon>
          <p>暂无权限</p>
        </div>
        <div v-else>
          <el-tag 
            v-for="permission in myPermissions" 
            :key="permission"
            type="primary"
            size="small"
            class="permission-tag"
          >
            {{ permission }}
          </el-tag>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { useAuthStore } from '../stores/auth'
import {
  User,
  Edit,
  Camera,
  Operation,
  Lock,
  Key,
  Document,
  Setting,
  Check,
  RefreshLeft
} from '@element-plus/icons-vue'

// 使用stores
const authStore = useAuthStore()
const router = useRouter()

// 响应式数据
const editMode = ref(false)
const saveLoading = ref(false)
const changePasswordVisible = ref(false)
const changePasswordLoading = ref(false)
const permissionsVisible = ref(false)
const myPermissions = ref<string[]>([])

// 表单引用
const profileFormRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()
const fileInput = ref<HTMLInputElement | null>(null)

// 用户信息
const userInfo = ref({
  id: 0,
  username: '',
  email: '',
  role: '',
  status: true,
  created_at: '',
  last_login: '',
  avatar: ''
})

// 个人资料表单
const profileForm = reactive({
  username: '',
  email: '',
  realName: '',
  phone: '',
  avatar: '',
  department: '',
  position: '',
  bio: ''
})

// 密码修改表单
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 表单验证规则
const profileRules: FormRules = {
  email: [
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ]
}

const passwordRules: FormRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 获取当前用户信息
const getCurrentUser = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('http://localhost:8081/api/me', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    const result = await response.json()
    if (result.code === 200) {
      userInfo.value = result.data
      // 填充表单
      Object.assign(profileForm, {
        username: result.data.username,
        email: result.data.email || '',
        realName: result.data.real_name || '',
        phone: result.data.phone || '',
        avatar: result.data.avatar || '',
        department: result.data.department || '',
        position: result.data.position || '',
        bio: result.data.bio || ''
      })
    } else {
      ElMessage.error(result.message || '获取用户信息失败')
    }
  } catch (error) {
    console.error('获取用户信息错误:', error)
    ElMessage.error('网络错误，请稍后重试')
  }
}

// 保存个人资料
const saveProfile = async () => {
  if (!profileFormRef.value) return

  try {
    await profileFormRef.value.validate()
    saveLoading.value = true

    const token = localStorage.getItem('token')
    const response = await fetch('http://localhost:8081/api/me', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        email: profileForm.email,
        real_name: profileForm.realName,
        phone: profileForm.phone,
        avatar: profileForm.avatar,
        department: profileForm.department,
        position: profileForm.position,
        bio: profileForm.bio
      })
    })

    const result = await response.json()
    if (result.code === 200) {
      ElMessage.success('个人资料保存成功')
      editMode.value = false
      getCurrentUser() // 刷新用户信息
    } else {
      ElMessage.error(result.message || '保存失败')
    }
  } catch (error) {
    if (error !== false) {
      console.error('保存个人资料错误:', error)
      ElMessage.error('保存失败')
    }
  } finally {
    saveLoading.value = false
  }
}

// 重置表单
const resetForm = () => {
  getCurrentUser()
  editMode.value = false
}

// 显示修改密码对话框
const showChangePasswordDialog = () => {
  changePasswordVisible.value = true
  // 重置表单
  Object.assign(passwordForm, {
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  })
}

// 修改密码
const changePassword = async () => {
  if (!passwordFormRef.value) return

  try {
    await passwordFormRef.value.validate()
    changePasswordLoading.value = true

    const token = localStorage.getItem('token')
    const response = await fetch('http://localhost:8081/api/change-password', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({
        current_password: passwordForm.currentPassword,
        new_password: passwordForm.newPassword
      })
    })

    const result = await response.json()
    if (result.code === 200) {
      ElMessage.success('密码修改成功')
      changePasswordVisible.value = false
      // 重置表单
      Object.assign(passwordForm, {
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      })
    } else {
      ElMessage.error(result.message || '密码修改失败')
    }
  } catch (error) {
    if (error !== false) {
      console.error('修改密码错误:', error)
      ElMessage.error('修改密码失败')
    }
  } finally {
    changePasswordLoading.value = false
  }
}

// 更换头像
function triggerFileInput() {
  fileInput.value?.click()
}

function onFileInputChange(e: Event) {
  const files = (e.target as HTMLInputElement).files
  if (!files || files.length === 0) return
  const file = files[0]
  if (!file.type.startsWith('image/')) {
    ElMessage.error('只能上传图片文件')
    return
  }
  if (file.size / 1024 / 1024 > 2) {
    ElMessage.error('图片大小不能超过2MB')
    return
  }
  const formData = new FormData()
  formData.append('file', file)
  formData.append('category', 'avatar')
  const token = localStorage.getItem('token')
  fetch('http://localhost:8081/api/files/upload', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    },
    body: formData
  })
    .then(res => res.json())
    .then(result => {
      if (result.code === 200 && result.data && result.data.url) {
        ElMessage.success('头像上传成功')
        getCurrentUser()
      } else {
        ElMessage.error(result.message || '头像上传失败')
      }
    })
    .catch(() => {
      ElMessage.error('网络错误，头像上传失败')
    })
}

// 查看权限
const viewPermissions = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('http://localhost:8081/api/my-permissions', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    const result = await response.json()
    if (result.code === 200) {
      myPermissions.value = result.data.permissions || []
      permissionsVisible.value = true
    } else {
      ElMessage.error(result.message || '获取权限失败')
    }
  } catch (error) {
    console.error('获取权限错误:', error)
    ElMessage.error('网络错误，请稍后重试')
  }
}

// 查看操作记录
const viewLogs = () => {
  router.push('/system/logs')
}

// 获取角色文本
const getRoleText = (role: string) => {
  const roleMap: Record<string, string> = {
    admin: '超级管理员',
    manager: '管理员',
    user: '普通用户'
  }
  return roleMap[role] || role
}

// 格式化时间
const formatDateTime = (dateTime: string) => {
  if (!dateTime) return ''
  return new Date(dateTime).toLocaleString('zh-CN')
}

// 初始化
onMounted(() => {
  getCurrentUser()
  
  // 检查URL参数，如果有action=change-password则自动打开修改密码对话框
  const urlParams = new URLSearchParams(window.location.search)
  if (urlParams.get('action') === 'change-password') {
    showChangePasswordDialog()
  }
})
</script>

<style scoped>
.main-grid {
  margin-top: var(--spacing-xl);
}

.profile-card,
.settings-card,
.security-card,
.quick-actions-card {
  margin-bottom: var(--spacing-xl);
}

.profile-avatar {
  text-align: center;
  margin-bottom: var(--spacing-xl);
  position: relative;
}

.avatar-upload-wrapper {
  position: relative;
  display: inline-block;
  cursor: pointer;
}

.user-avatar.clickable {
  transition: box-shadow 0.2s;
}

.avatar-upload-wrapper:hover .user-avatar.clickable {
  box-shadow: 0 0 0 3px #409eff33;
}

.avatar-actions {
  position: absolute;
  bottom: 0;
  right: 0;
  z-index: 2;
}

.avatar-hover-mask {
  display: none;
  position: absolute;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.3);
  color: #fff;
  font-size: 14px;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  z-index: 1;
  pointer-events: none;
}

.avatar-upload-wrapper:hover .avatar-hover-mask {
  display: flex;
}

.profile-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm) 0;
  border-bottom: 1px solid var(--border-secondary);
}

.info-label {
  font-weight: var(--font-medium);
  color: var(--text-secondary);
}

.info-value {
  color: var(--text-primary);
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.action-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  padding: var(--spacing-lg);
  border: 1px solid var(--border-secondary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-normal) ease;
}

.action-item:hover {
  border-color: var(--primary-color);
  box-shadow: var(--shadow-sm);
}

.action-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.action-text h4 {
  margin: 0 0 var(--spacing-xs);
  color: var(--text-primary);
}

.action-text p {
  margin: 0;
  font-size: var(--font-sm);
  color: var(--text-secondary);
}

.settings-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.security-items {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-lg);
  border: 1px solid var(--border-secondary);
  border-radius: var(--radius-lg);
}

.security-info h4 {
  margin: 0 0 var(--spacing-xs);
  color: var(--text-primary);
}

.security-info p {
  margin: 0;
  font-size: var(--font-sm);
  color: var(--text-secondary);
}

.permissions-list {
  max-height: 400px;
  overflow-y: auto;
}

.empty-permissions {
  text-align: center;
  padding: var(--spacing-2xl);
  color: var(--text-secondary);
}

.permission-tag {
  margin: var(--spacing-xs);
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
  .profile-card,
  .settings-card,
  .security-card,
  .quick-actions-card {
    margin-bottom: var(--spacing-lg);
  }

  .info-item {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-xs);
  }

  .security-item {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-md);
  }

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

/* 表单组件样式 */

:deep(.el-form) {
  background: transparent !important;
}

:deep(.el-form-item__label) {
  color: var(--text-secondary) !important;
}

:deep(.el-input__wrapper) {
  background: var(--bg-input) !important;
  border: 1px solid var(--border-secondary) !important;
  border-radius: var(--radius-md) !important;
}

:deep(.el-input__wrapper:hover) {
  border-color: var(--primary-color) !important;
}

:deep(.el-input__wrapper.is-focus) {
  border-color: var(--primary-color) !important;
  box-shadow: 0 0 0 2px var(--primary-color-light) !important;
}

:deep(.el-textarea__inner) {
  background: var(--bg-input) !important;
  border: 1px solid var(--border-secondary) !important;
  border-radius: var(--radius-md) !important;
  color: var(--text-primary) !important;
}

:deep(.el-textarea__inner:hover) {
  border-color: var(--primary-color) !important;
}

:deep(.el-textarea__inner:focus) {
  border-color: var(--primary-color) !important;
  box-shadow: 0 0 0 2px var(--primary-color-light) !important;
}
</style> 