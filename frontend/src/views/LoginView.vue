<template>
  <div class="login-card">
      <div class="card-header">
        <div class="logo-section">
          <el-icon class="logo-icon" size="32"><House /></el-icon>
          <h1 class="system-title">Jing Admin</h1>
        </div>
        <p class="subtitle">现代化管理后台系统</p>
      </div>

      <div class="card-body">
        <!-- 切换标签 -->
        <div class="tab-switcher">
          <button 
            :class="['tab-btn', { active: activeTab === 'login' }]"
            @click="activeTab = 'login'"
          >
            登录
          </button>
          <button 
            :class="['tab-btn', { active: activeTab === 'register' }]"
            @click="activeTab = 'register'"
          >
            注册
          </button>
        </div>

        <!-- 登录表单 -->
        <div v-if="activeTab === 'login'" class="form-container">
          <el-form
            ref="loginFormRef"
            :model="loginForm"
            :rules="loginRules"
            @submit.prevent="handleLogin"
          >
            <el-form-item prop="username">
              <el-input
                v-model="loginForm.username"
                placeholder="用户名"
                :prefix-icon="User"
                size="large"
                clearable
                @keyup.enter="focusPassword"
              />
            </el-form-item>
            
            <el-form-item prop="password">
              <el-input
                ref="passwordInput"
                v-model="loginForm.password"
                placeholder="密码"
                type="password"
                :prefix-icon="Lock"
                size="large"
                show-password
                @keyup.enter="handleLogin"
              />
            </el-form-item>

            <div class="login-options">
              <el-checkbox v-model="rememberMe">记住密码</el-checkbox>
              <el-link type="primary" @click="activeTab = 'forgot'">忘记密码？</el-link>
            </div>

            <el-form-item>
              <el-button
                type="primary"
                size="large"
                style="width: 100%"
                :loading="authStore.loading"
                @click="handleLogin"
              >
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 注册表单 -->
        <div v-if="activeTab === 'register'" class="form-container">
          <el-form
            ref="registerFormRef"
            :model="registerForm"
            :rules="registerRules"
            @submit.prevent="handleRegister"
          >
            <el-form-item prop="username">
              <el-input
                v-model="registerForm.username"
                placeholder="用户名 (3-20字符)"
                :prefix-icon="User"
                size="large"
                clearable
              />
            </el-form-item>
            
            <el-form-item prop="email">
              <el-input
                v-model="registerForm.email"
                placeholder="邮箱地址"
                :prefix-icon="Message"
                size="large"
                clearable
              />
            </el-form-item>
            
            <el-form-item prop="password">
              <el-input
                v-model="registerForm.password"
                placeholder="密码 (至少6位)"
                type="password"
                :prefix-icon="Lock"
                size="large"
                show-password
              />
            </el-form-item>
            
            <el-form-item prop="confirmPassword">
              <el-input
                v-model="registerForm.confirmPassword"
                placeholder="确认密码"
                type="password"
                :prefix-icon="Lock"
                size="large"
                show-password
                @keyup.enter="handleRegister"
              />
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                size="large"
                style="width: 100%"
                :loading="authStore.loading"
                @click="handleRegister"
              >
                注册
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 默认账户提示 -->
        <div class="demo-hint">
          <el-alert
            title="演示账户"
            type="info"
            :closable="false"
            show-icon
          >
            <template #default>
              <p>默认管理员账户：</p>
              <p><strong>用户名：</strong> admin</p>
              <p><strong>密码：</strong> admin123</p>
            </template>
          </el-alert>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import type { RuleItem } from 'element-plus/es/components/form/src/types'
import { User, Lock, House, Message } from '@element-plus/icons-vue'
import { useAuthStore, type LoginForm, type RegisterForm } from '../stores/auth'
import { security } from '../utils/error-handler'

const router = useRouter()
const authStore = useAuthStore()

// 表单引用
const loginFormRef = ref<FormInstance>()
const registerFormRef = ref<FormInstance>()
const passwordInput = ref<any>()

// 记住密码
const rememberMe = ref(!!localStorage.getItem('rememberedUsername'))

// 标签状态
const activeTab = ref<'login' | 'register' | 'forgot'>('login')

// 登录表单
const loginForm = reactive<LoginForm>({
  username: localStorage.getItem('rememberedUsername') || '',
  password: '',
})

// 注册表单
const registerForm = reactive<RegisterForm>({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
})

// 表单验证规则
const loginRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const registerRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度为3-20个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule: RuleItem, value: string, callback: (error?: Error) => void) => {
        if (value !== registerForm.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 聚焦密码输入框
function focusPassword() {
  passwordInput.value?.focus()
}

// 登录尝试限制
const LOGIN_MAX_ATTEMPTS = 5
const LOGIN_LOCKOUT_TIME = 15 * 60 * 1000 // 15分钟
const loginAttempts = ref(Number(localStorage.getItem('loginAttempts') || '0'))
const loginLockoutUntil = ref<number>(Number(localStorage.getItem('loginLockoutUntil') || '0'))

// 检查是否被锁定
function isLoginLocked(): boolean {
  const now = Date.now()
  if (loginLockoutUntil.value > now) {
    const remainingMinutes = Math.ceil((loginLockoutUntil.value - now) / 60000)
    ElMessage.warning(`登录已被锁定，请${remainingMinutes}分钟后重试`)
    return true
  }
  if (loginLockoutUntil.value && loginLockoutUntil.value <= now) {
    // 锁定时间已过，重置尝试次数
    loginAttempts.value = 0
    loginLockoutUntil.value = 0
    localStorage.removeItem('loginAttempts')
    localStorage.removeItem('loginLockoutUntil')
  }
  return false
}

// 记录登录尝试
function recordLoginAttempt(success: boolean) {
  if (success) {
    // 登录成功，重置尝试次数
    loginAttempts.value = 0
    loginLockoutUntil.value = 0
    localStorage.removeItem('loginAttempts')
    localStorage.removeItem('loginLockoutUntil')
  } else {
    // 登录失败，增加尝试次数
    loginAttempts.value++
    localStorage.setItem('loginAttempts', loginAttempts.value.toString())
    
    // 检查是否需要锁定
    if (loginAttempts.value >= LOGIN_MAX_ATTEMPTS) {
      loginLockoutUntil.value = Date.now() + LOGIN_LOCKOUT_TIME
      localStorage.setItem('loginLockoutUntil', loginLockoutUntil.value.toString())
      ElMessage.error(`登录尝试次数过多，请${LOGIN_LOCKOUT_TIME / 60000}分钟后重试`)
    } else {
      const remainingAttempts = LOGIN_MAX_ATTEMPTS - loginAttempts.value
      ElMessage.warning(`登录失败，还剩${remainingAttempts}次尝试机会`)
    }
  }
}

// 密码强度检查
const passwordStrength = ref({ score: 0, feedback: [''] })

function checkPasswordStrength() {
  if (registerForm.password) {
    passwordStrength.value = security.checkPasswordStrength(registerForm.password)
  }
}

// 处理登录
async function handleLogin() {
  if (!loginFormRef.value) return
  
  // 检查是否被锁定
  if (isLoginLocked()) return
  
  await loginFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      const success = await authStore.login(loginForm)
      recordLoginAttempt(success)
      if (success) {
        // 处理记住密码
        if (rememberMe.value) {
          localStorage.setItem('rememberedUsername', loginForm.username)
        } else {
          localStorage.removeItem('rememberedUsername')
        }
        router.push('/')
      }
    }
  })
}

// 处理注册
async function handleRegister() {
  if (!registerFormRef.value) return
  
  await registerFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      // 检查密码强度
      const strength = security.checkPasswordStrength(registerForm.password)
      if (strength.score < 3) {
        ElMessage.warning('密码强度不足，请按照建议修改密码')
        return
      }
      
      const success = await authStore.register(registerForm)
      if (success) {
        router.push('/')
      }
    }
  })
}

// 生命周期钩子
onMounted(() => {
  // 如果有记住的用户名，自动聚焦密码输入框
  if (loginForm.username) {
    focusPassword()
  }
})

onBeforeUnmount(() => {
  // 清理认证store的定时器
  authStore.cleanup()
})
</script>

<style scoped>
/* 登录卡片容器样式已移至 AuthLayout */

.login-card {
  position: relative;
  z-index: 1;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 40px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
}

/* 深色模式下的登录卡片 */
html.dark .login-card {
  background: rgba(26, 26, 26, 0.95);
  border: 1px solid rgba(64, 64, 64, 0.3);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
}

.card-header {
  text-align: center;
  margin-bottom: 30px;
}

.logo-section {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 8px;
}

.logo-icon {
  color: #667eea;
}

.system-title {
  margin: 0;
  font-size: 28px;
  font-weight: 700;
  color: #2d3748;
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 深色模式下的标题 */
html.dark .system-title {
  color: #ffffff;
  background: linear-gradient(135deg, #60a5fa, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  margin: 0;
  color: #718096;
  font-size: 14px;
}

/* 深色模式下的副标题 */
html.dark .subtitle {
  color: #a3a3a3;
}

.tab-switcher {
  display: flex;
  background: #f7fafc;
  border-radius: 8px;
  padding: 4px;
  margin-bottom: 24px;
}

.tab-btn {
  flex: 1;
  padding: 8px 16px;
  border: none;
  background: transparent;
  color: #718096;
  font-size: 14px;
  font-weight: 500;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.tab-btn.active {
  background: #667eea;
  color: white;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.tab-btn:hover:not(.active) {
  color: #4a5568;
  background: rgba(102, 126, 234, 0.1);
}

.form-container {
  margin-bottom: 20px;
}

.demo-hint {
  margin-top: 20px;
}

/* Element Plus 组件样式覆盖 */
:deep(.el-form-item) {
  margin-bottom: 16px;
}

:deep(.el-input__wrapper) {
  background: rgba(247, 250, 252, 0.8);
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  box-shadow: none;
  transition: all 0.3s ease;
}

:deep(.el-input__wrapper:hover) {
  border-color: #667eea;
}

:deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

:deep(.el-button--primary) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

:deep(.el-button--primary:hover) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
}

:deep(.el-alert) {
  border-radius: 8px;
  border: none;
  background: rgba(102, 126, 234, 0.1);
}

:deep(.el-alert__title) {
  font-size: 14px;
  font-weight: 600;
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style> 