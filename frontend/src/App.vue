<template>
  <div id="app">
    <!-- 登录页面独立布局 -->
    <template v-if="isLoginPage">
      <router-view />
    </template>

    <!-- 主应用布局 -->
    <template v-else>
      <el-container class="layout-container">
        <!-- 简洁顶部导航栏 -->
        <el-header class="header">
          <div class="header-content">
            <div class="logo-section slide-in-left">
              <el-icon class="logo-icon" size="24"><House /></el-icon>
              <h1 class="title">Jing Admin</h1>
              <el-tag size="small" class="version-tag">v1.0</el-tag>
            </div>
            <div class="header-actions slide-in-right">
              <!-- 主题切换按钮 -->
              <el-tooltip 
                :content="themeStore.themeText" 
                placement="bottom" 
                effect="dark"
                popper-class="custom-tooltip"
              >
                <el-button 
                  :icon="themeStore.isDark ? Sunny : Moon"
                  circle 
                  class="theme-toggle"
                  @click="themeStore.toggleTheme"
                  size="small"
                />
              </el-tooltip>
              
              <el-tooltip 
                content="通知" 
                placement="bottom" 
                effect="dark"
                popper-class="custom-tooltip"
              >
                <el-badge :value="3" class="notification-badge">
                  <el-button :icon="Bell" circle size="small" />
                </el-badge>
              </el-tooltip>
              
              <template v-if="authStore.isAuthenticated">
                <el-dropdown class="user-dropdown">
                  <el-avatar :size="28" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item>
                        <div class="user-info">
                          <div class="username">{{ authStore.currentUser?.username }}</div>
                          <div class="role">{{ authStore.isAdmin ? '管理员' : '普通用户' }}</div>
                        </div>
                      </el-dropdown-item>
                      <el-dropdown-item divided @click="goToProfile">个人中心</el-dropdown-item>
                      <el-dropdown-item @click="goToProfileAndChangePassword">修改密码</el-dropdown-item>
                      <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </template>
            </div>
          </div>
        </el-header>

        <el-container class="main-container">
        <!-- 简洁侧边栏 -->
        <el-aside class="sidebar" :width="'220px'">
          <el-menu
            :default-active="activeMenuIndex"
            class="sidebar-menu"
            :background-color="'transparent'"
            :text-color="'var(--menu-text)'"
            :active-text-color="'var(--menu-active-text)'"
            unique-opened
            router
          >
            <el-menu-item index="/">
              <el-icon><HomeFilled /></el-icon>
              <span>首页</span>
            </el-menu-item>
            
            <el-sub-menu index="/users">
              <template #title>
                <el-icon><User /></el-icon>
                <span>用户管理</span>
              </template>
              <el-menu-item index="/users">用户列表</el-menu-item>
              <el-menu-item index="/users/roles">角色管理</el-menu-item>
              <el-menu-item index="/users/permissions">权限配置</el-menu-item>
            </el-sub-menu>
            
            <el-sub-menu index="/system">
              <template #title>
                <el-icon><Setting /></el-icon>
                <span>系统管理</span>
              </template>
              <el-menu-item index="/system/config">系统配置</el-menu-item>
              <el-menu-item index="/system/logs">操作日志</el-menu-item>
              <el-menu-item index="/system/backup">数据备份</el-menu-item>
            </el-sub-menu>
            
            <el-menu-item index="/analytics">
              <el-icon><DataAnalysis /></el-icon>
              <span>数据分析</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        
        <!-- 主内容区域 -->
        <el-main class="main-content">
          <router-view v-slot="{ Component }">
            <transition 
              name="fade-transform" 
              mode="out-in"
            >
              <keep-alive>
                <component :is="Component" />
              </keep-alive>
            </transition>
          </router-view>
        </el-main>
        </el-container>
      </el-container>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useThemeStore } from './stores/theme'
import { useAuthStore } from './stores/auth'
import { 
  User, 
  Setting, 
  House, 
  Bell, 
  DataAnalysis,
  Moon,
  Sunny,
  HomeFilled
} from '@element-plus/icons-vue'

const themeStore = useThemeStore()
const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()

// 判断是否为登录页面
const isLoginPage = computed(() => route.path === '/login')

// 计算当前活动的菜单项
const activeMenuIndex = computed(() => {
  return route.path
})

// 处理登出
async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}

// 前往个人中心
function goToProfile() {
  router.push('/profile')
}

// 前往个人中心并打开修改密码对话框
function goToProfileAndChangePassword() {
  router.push('/profile?action=change-password')
}

onMounted(() => {
  themeStore.initTheme()
  authStore.initAuth()
})
</script>

<style scoped>
.layout-container {
  height: 100vh;
  width: 100%;
  background: transparent;
  display: flex;
  flex-direction: column;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
}

/* 顶部导航栏样式 */
.header {
  background: var(--bg-header);
  color: var(--text-primary);
  display: flex;
  align-items: center;
  padding: 0 var(--spacing-xl);
  border-bottom: 1px solid var(--border-primary);
  height: 64px;
  backdrop-filter: blur(20px);
  box-shadow: var(--shadow-sm);
  flex-shrink: 0;
  width: 100% !important;
  margin: 0 !important;
  position: relative;
  z-index: 1000;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  color: var(--accent-color);
}

.title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.version-tag {
  background: var(--hover-bg);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.theme-toggle {
  transition: all 0.3s ease;
}

.notification-badge {
  margin-right: 4px;
}

.user-dropdown {
  cursor: pointer;
  margin-left: 4px;
  transition: transform 0.3s ease;
}

.user-dropdown:hover {
  transform: scale(1.05);
}

.user-info {
  padding: 4px 0;
}

.username {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
}

.role {
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: 2px;
}

/* 主容器样式 */
.main-container {
  flex: 1;
  display: flex;
  background: transparent;
  overflow: hidden;
  width: 100%;
  margin: 0;
  padding: 0;
}

/* 侧边栏样式 */
.sidebar {
  background: var(--bg-sidebar);
  backdrop-filter: blur(20px);
  border-right: 1px solid var(--border-primary);
  padding: var(--spacing-md);
  transition: all var(--duration-normal) ease;
  position: relative;
  overflow-y: auto;
  overflow-x: hidden;
  height: calc(100vh - 64px);
  flex-shrink: 0;
}

.sidebar-menu {
  border: none;
  background: transparent;
}

.sidebar-menu :deep(.el-menu-item),
.sidebar-menu :deep(.el-sub-menu__title) {
  height: 44px;
  line-height: 44px;
  border-radius: var(--radius-lg);
  margin: 4px 0;
  color: var(--text-secondary);
  transition: all var(--duration-normal) ease;
}

.sidebar-menu :deep(.el-menu-item .el-icon),
.sidebar-menu :deep(.el-sub-menu__title .el-icon) {
  margin-right: var(--spacing-md);
  font-size: var(--font-md);
  transition: all var(--duration-normal) ease;
}

.sidebar-menu :deep(.el-menu-item:hover),
.sidebar-menu :deep(.el-sub-menu__title:hover) {
  background: var(--bg-hover);
  color: var(--primary-color);
}

.sidebar-menu :deep(.el-menu-item:hover .el-icon),
.sidebar-menu :deep(.el-sub-menu__title:hover .el-icon) {
  color: var(--primary-color);
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background: var(--bg-selected);
  color: var(--primary-color);
  font-weight: var(--font-medium);
}

.sidebar-menu :deep(.el-menu-item.is-active .el-icon) {
  color: var(--primary-color);
}

.sidebar-menu :deep(.el-sub-menu.is-active .el-sub-menu__title) {
  color: var(--primary-color);
}

.sidebar-menu :deep(.el-sub-menu.is-active .el-sub-menu__title .el-icon) {
  color: var(--primary-color);
}

.sidebar-menu :deep(.el-sub-menu__title) {
  padding-left: var(--spacing-md) !important;
}

.sidebar-menu :deep(.el-menu-item) {
  padding-left: var(--spacing-md) !important;
}

.sidebar-menu :deep(.el-sub-menu .el-menu-item) {
  padding-left: calc(var(--spacing-md) + var(--spacing-xl)) !important;
  min-width: auto;
}

/* 主内容区域样式 */
.main-content {
  background: var(--bg-page);
  padding: var(--spacing-xl);
  overflow-y: auto;
  overflow-x: hidden;
  flex: 1;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: -220px;
    top: 64px;
    bottom: 0;
    z-index: 1000;
    transition: left var(--duration-normal) ease;
  }

  .sidebar.is-open {
    left: 0;
  }
}

@media (max-width: 480px) {
  .logo-section {
    gap: var(--spacing-sm);
  }
  
  .title {
    font-size: var(--font-sm);
  }
  
  .header-actions {
    gap: var(--spacing-xs);
  }
  
  .main-content {
    padding: var(--spacing-sm);
  }
}

/* 路由过渡动画 */
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all var(--duration-normal) ease;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

/* 移除页面内组件的重复动画 */
:deep(.fade-in),
:deep(.slide-in-left),
:deep(.slide-in-right),
:deep(.scale-in) {
  animation: none;
}

/* 自定义tooltip样式 */
:deep(.custom-tooltip) {
  padding: 8px 12px;
  font-size: 14px;
  line-height: 1.4;
  border-radius: 4px;
  z-index: 9999;
}
</style>
