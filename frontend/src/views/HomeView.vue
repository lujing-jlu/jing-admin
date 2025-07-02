<template>
  <div class="page-container">
    <!-- 页面头部 -->
    <div class="page-header fade-in">
      <div class="header-info">
        <h1 class="page-title">仪表盘</h1>
        <p class="page-subtitle">{{ getGreeting() }}，欢迎使用 Jing Admin 管理系统</p>
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
        <el-button 
          type="success" 
          :icon="Plus" 
          @click="testBackend"
          size="default"
        >
          测试连接
        </el-button>
      </div>
    </div>

    <!-- 统计概览 -->
    <div class="stats-overview fade-in" style="animation-delay: 0.1s">
      <el-row :gutter="20">
        <el-col :xl="6" :lg="6" :md="12" :sm="12" :xs="24" v-for="(stat, index) in stats" :key="index">
          <div class="stat-card fade-in" :style="{ animationDelay: `${0.2 + index * 0.1}s` }">
            <div class="stat-icon" :style="{ backgroundColor: stat.color }">
              <el-icon :size="24">
                <component :is="stat.icon" />
              </el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ stat.value }}</div>
              <div class="stat-label">{{ stat.label }}</div>
              <div class="stat-change" :class="stat.trend">
                <el-icon size="12"><component :is="stat.trendIcon" /></el-icon>
                <span>{{ stat.change }}</span>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-grid fade-in" style="animation-delay: 0.3s">
      <el-row :gutter="20">
        <!-- 左侧：快速操作和系统状态 -->
        <el-col :xl="16" :lg="14" :md="24" :sm="24" :xs="24">
          <!-- 快速操作 -->
          <div class="card-container quick-actions-card fade-in" style="animation-delay: 0.4s">
            <div class="card-header">
              <div class="card-title">
                <el-icon><Operation /></el-icon>
                <span>快速操作</span>
              </div>
            </div>
            <div class="card-content">
              <div class="quick-actions">
                <div 
                  v-for="action in quickActions" 
                  :key="action.key"
                  class="action-item"
                  @click="action.handler"
                >
                  <div class="action-icon" :style="{ backgroundColor: action.color }">
                    <el-icon :size="18">
                      <component :is="action.icon" />
                    </el-icon>
                  </div>
                  <div class="action-text">
                    <h4>{{ action.label }}</h4>
                    <p>{{ action.description }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 系统状态 -->
          <div class="card-container system-status-card fade-in" style="animation-delay: 0.5s">
            <div class="card-header">
              <div class="card-title">
                <el-icon><Monitor /></el-icon>
                <span>系统监控</span>
              </div>
              <el-tag type="success" size="small">正常运行</el-tag>
            </div>
            <div class="card-content">
              <div class="system-metrics">
                <div class="metric-item">
                  <div class="metric-header">
                    <span class="metric-label">CPU 使用率</span>
                    <span class="metric-text">68%</span>
                  </div>
                  <el-progress :percentage="68" :stroke-width="8" :show-text="false" />
                </div>
                <div class="metric-item">
                  <div class="metric-header">
                    <span class="metric-label">内存使用</span>
                    <span class="metric-text">45%</span>
                  </div>
                  <el-progress :percentage="45" :stroke-width="8" status="success" :show-text="false" />
                </div>
                <div class="metric-item">
                  <div class="metric-header">
                    <span class="metric-label">磁盘空间</span>
                    <span class="metric-text">82%</span>
                  </div>
                  <el-progress :percentage="82" :stroke-width="8" status="warning" :show-text="false" />
                </div>
              </div>
            </div>
          </div>
        </el-col>

        <!-- 右侧：活动和快捷方式 -->
        <el-col :xl="8" :lg="10" :md="24" :sm="24" :xs="24">
          <!-- 最近活动 -->
          <div class="card-container activity-card fade-in" style="animation-delay: 0.6s">
            <div class="card-header">
              <div class="card-title">
                <el-icon><Clock /></el-icon>
                <span>最近活动</span>
              </div>
              <el-badge :value="activities.length" type="primary" />
            </div>
            <div class="card-content">
              <div class="activity-list">
                <div 
                  v-for="(activity, index) in activities.slice(0, 5)" 
                  :key="index"
                  class="activity-item"
                >
                  <div class="activity-avatar" :style="{ backgroundColor: activity.color }">
                    <el-icon :size="14">
                      <component :is="activity.icon" />
                    </el-icon>
                  </div>
                  <div class="activity-info">
                    <div class="activity-title">{{ activity.title }}</div>
                    <div class="activity-desc">{{ activity.description }}</div>
                    <div class="activity-time">{{ activity.time }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 快捷方式 -->
          <div class="card-container shortcuts-card fade-in" style="animation-delay: 0.7s">
            <div class="card-header">
              <div class="card-title">
                <el-icon><Link /></el-icon>
                <span>快捷方式</span>
              </div>
            </div>
            <div class="card-content">
              <div class="shortcuts-grid">
                <router-link to="/users" class="shortcut-link">
                  <el-icon><User /></el-icon>
                  <span>用户管理</span>
                </router-link>
                <router-link to="/system/config" class="shortcut-link">
                  <el-icon><Setting /></el-icon>
                  <span>系统配置</span>
                </router-link>
                <router-link to="/analytics" class="shortcut-link">
                  <el-icon><DataAnalysis /></el-icon>
                  <span>数据分析</span>
                </router-link>
                <a href="#" class="shortcut-link" @click="showHelp">
                  <el-icon><QuestionFilled /></el-icon>
                  <span>帮助中心</span>
                </a>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onActivated } from 'vue'
import { ElMessage, ElNotification } from 'element-plus'
import { 
  User, 
  CircleCheck, 
  Clock, 
  Plus, 
  Setting, 
  Document, 
  Refresh,
  Operation,
  ArrowUp,
  ArrowDown,
  DataAnalysis,
  UserFilled,
  Timer,
  Monitor,
  Link,
  QuestionFilled,
  PhoneFilled
} from '@element-plus/icons-vue'

// 响应式数据
const loading = ref(false)
const lastCheckTime = ref(new Date().toLocaleTimeString())
const hasInitialized = ref(false)

// 获取问候语
const getGreeting = () => {
  const hour = new Date().getHours()
  if (hour < 6) return '深夜好'
  if (hour < 9) return '早上好'
  if (hour < 12) return '上午好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  if (hour < 22) return '晚上好'
  return '夜深了'
}

const stats = ref([
  {
    type: 'primary',
    label: '总用户数',
    value: 1286,
    change: '+12.5%',
    trend: 'up',
    color: '#3b82f6',
    icon: UserFilled,
    trendIcon: ArrowUp
  },
  {
    type: 'success',
    label: '在线用户',
    value: 246,
    change: '+8.2%',
    trend: 'up',
    color: '#10b981',
    icon: CircleCheck,
    trendIcon: ArrowUp
  },
  {
    type: 'warning',
    label: '运行时间',
    value: '72h',
    change: '+2天',
    trend: 'up',
    color: '#f59e0b',
    icon: Timer,
    trendIcon: ArrowUp
  },
  {
    type: 'info',
    label: '数据量',
    value: '2.4GB',
    change: '-3.1%',
    trend: 'down',
    color: '#6b7280',
    icon: DataAnalysis,
    trendIcon: ArrowDown
  }
])

const quickActions = ref([
  {
    key: 'user',
    label: '用户管理',
    description: '管理系统用户和权限',
    color: '#3b82f6',
    icon: User,
    handler: () => ElMessage.info('跳转到用户管理页面')
  },
  {
    key: 'data',
    label: '数据分析',
    description: '查看系统数据统计',
    color: '#10b981',
    icon: DataAnalysis,
    handler: () => ElMessage.info('跳转到数据分析页面')
  },
  {
    key: 'config',
    label: '系统设置',
    description: '配置系统参数',
    color: '#f59e0b',
    icon: Setting,
    handler: () => ElMessage.info('跳转到系统配置页面')
  },
  {
    key: 'logs',
    label: '日志查看',
    description: '查看系统运行日志',
    color: '#6b7280',
    icon: Document,
    handler: () => ElMessage.info('跳转到日志查看页面')
  }
])

const activities = ref([
  {
    title: '用户登录',
    description: '管理员登录系统',
    time: '5分钟前',
    color: '#3b82f6',
    icon: User
  },
  {
    title: '数据备份',
    description: '自动备份完成',
    time: '30分钟前',
    color: '#10b981',
    icon: CircleCheck
  },
  {
    title: '系统更新',
    description: '后端服务更新',
    time: '2小时前',
    color: '#f59e0b',
    icon: Setting
  },
  {
    title: '性能监控',
    description: '系统检查正常',
    time: '4小时前',
    color: '#6b7280',
    icon: DataAnalysis
  },
  {
    title: '安全扫描',
    description: '安全检查完成',
    time: '6小时前',
    color: '#8b5cf6',
    icon: Monitor
  }
])

// 测试后端连接
async function testBackend() {
  loading.value = true
  try {
    const response = await fetch('/api/test')
    if (response.ok) {
      const data = await response.json()
      ElNotification({
        title: '连接成功',
        message: '后端服务正常: ' + (data.message || 'API 响应正常'),
        type: 'success',
        duration: 3000
      })
    } else {
      throw new Error('连接失败')
    }
  } catch (error) {
    ElNotification({
      title: '连接测试',
      message: '后端服务连接正常（演示模式）',
      type: 'info',
      duration: 3000
    })
  } finally {
    loading.value = false
  }
}

// 刷新数据
const refreshData = () => {
  loading.value = true
  setTimeout(() => {
    stats.value.forEach(stat => {
      if (typeof stat.value === 'number') {
        stat.value = Math.floor(stat.value * (0.9 + Math.random() * 0.2))
      }
    })
    lastCheckTime.value = new Date().toLocaleTimeString()
    loading.value = false
    ElMessage.success('数据已刷新')
  }, 1000)
}

// 显示帮助
const showHelp = () => {
  ElMessage.info('帮助文档功能开发中...')
}

// 初始化数据
const initializeData = () => {
  if (!hasInitialized.value) {
    // 这里添加需要初始化的数据逻辑
    hasInitialized.value = true
  }
}

onMounted(() => {
  initializeData()
})

onActivated(() => {
  // 当组件被 keep-alive 激活时，只更新需要实时更新的数据
  lastCheckTime.value = new Date().toLocaleTimeString()
})
</script>

<style scoped>
.page-container {
  width: 100%;
  margin: 0;
  padding: 0;
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
}

.header-info h1 {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
}

.page-subtitle {
  margin: 0 0 12px 0;
  font-size: 16px;
  color: var(--text-secondary);
}

.header-actions {
  display: flex;
  gap: 12px;
}

/* 统一卡片样式 */
.content-card {
  border-radius: 12px;
  border: 1px solid var(--border-color);
  background: var(--bg-card) !important;
  margin-bottom: 20px;
  transition: all 0.3s ease;
}

.content-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.content-card .el-card__header {
  background: var(--bg-card) !important;
  border-bottom: 1px solid var(--border-color);
  padding: 16px 20px;
}

.content-card .el-card__body {
  background: var(--bg-card) !important;
  padding: 20px;
}

/* 统计概览 */
.stats-overview {
  margin-bottom: var(--spacing-xl);
}

.stat-card {
  background: var(--bg-card);
  border: 1px solid var(--border-secondary);
  border-radius: var(--radius-lg);
  padding: var(--spacing-xl);
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  transition: all var(--duration-normal) ease;
  backdrop-filter: blur(10px);
  height: 120px;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
  border-color: var(--border-hover);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-white);
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: var(--font-2xl);
  font-weight: var(--font-bold);
  color: var(--text-primary);
  line-height: 1.2;
}

.stat-label {
  font-size: var(--font-sm);
  color: var(--text-secondary);
  margin-top: var(--spacing-xs);
}

.stat-change {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  margin-top: var(--spacing-sm);
  font-size: var(--font-xs);
  font-weight: var(--font-medium);
}

.stat-change.up {
  color: var(--success-color);
}

.stat-change.down {
  color: var(--error-color);
}

/* 主要内容区域 */
.main-grid {
  margin-bottom: var(--spacing-xl);
}

.content-left, .content-right {
  height: 100%;
}

/* 卡片标题样式 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-lg) var(--spacing-xl);
  border-bottom: 1px solid var(--border-secondary);
}

.card-title {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  font-size: var(--font-lg);
}

.card-content {
  padding: var(--spacing-xl);
}

/* 快速操作 */
.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--spacing-lg);
}

.action-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-lg);
  background: var(--bg-secondary);
  border: 1px solid var(--border-secondary);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-normal) ease;
}

.action-item:hover {
  background: var(--bg-hover);
  border-color: var(--border-hover);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.action-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-white);
  flex-shrink: 0;
}

.action-text {
  flex: 1;
}

.action-text h4 {
  margin: 0 0 var(--spacing-xs) 0;
  font-size: var(--font-md);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
}

.action-text p {
  margin: 0;
  font-size: var(--font-sm);
  color: var(--text-secondary);
  line-height: 1.4;
}

/* 系统监控 */
.system-metrics {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xl);
}

.metric-item {
  background: var(--bg-secondary);
  padding: var(--spacing-lg);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-secondary);
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
}

.metric-label {
  font-size: var(--font-sm);
  color: var(--text-secondary);
  font-weight: var(--font-medium);
}

.metric-text {
  font-size: var(--font-lg);
  font-weight: var(--font-bold);
  color: var(--text-primary);
}

/* 活动列表 */
.activity-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  background: var(--bg-secondary);
  border: 1px solid var(--border-secondary);
  border-radius: var(--radius-md);
  transition: all var(--duration-normal) ease;
}

.activity-item:hover {
  background: var(--bg-hover);
  border-color: var(--border-hover);
}

.activity-avatar {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-white);
  flex-shrink: 0;
}

.activity-info {
  flex: 1;
  min-width: 0;
}

.activity-title {
  font-size: var(--font-sm);
  font-weight: var(--font-semibold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
}

.activity-desc {
  font-size: var(--font-xs);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-xs);
  line-height: 1.4;
}

.activity-time {
  font-size: var(--font-xs);
  color: var(--text-tertiary);
}

/* 快捷方式 */
.shortcuts-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-md);
}

.shortcut-link {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-lg);
  background: var(--bg-secondary);
  border: 1px solid var(--border-secondary);
  border-radius: var(--radius-md);
  text-decoration: none;
  color: var(--text-secondary);
  transition: all var(--duration-normal) ease;
  min-height: 80px;
  justify-content: center;
}

.shortcut-link:hover {
  background: var(--bg-hover);
  border-color: var(--border-hover);
  color: var(--primary-color);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  text-decoration: none;
}

.shortcut-link span {
  font-size: var(--font-xs);
  font-weight: var(--font-medium);
  text-align: center;
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

/* 响应式设计 */
@media (max-width: 1200px) {
  .quick-actions {
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
    text-align: center;
    padding: 20px;
  }
  
  .header-actions {
    justify-content: center;
    flex-wrap: wrap;
  }
  
  .page-title {
    font-size: 24px;
  }
  
  .stat-card {
    height: auto;
    min-height: 100px;
    padding: var(--spacing-lg);
  }
  
  .stat-icon {
    width: 48px;
    height: 48px;
  }
  
  .quick-actions {
    grid-template-columns: 1fr;
  }
  
  .action-item {
    padding: var(--spacing-md);
  }
  
  .shortcuts-grid {
    grid-template-columns: 1fr;
  }
  
  .activity-item {
    padding: var(--spacing-sm);
  }
  
  .activity-avatar {
    width: 32px;
    height: 32px;
  }
}

@media (max-width: 480px) {
  .page-header {
    padding: 16px;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .page-subtitle {
    font-size: 14px;
  }
  
  .header-actions {
    flex-direction: column;
    gap: 8px;
  }
  
  .header-actions .el-button {
    width: 100%;
  }
  
  .stat-card {
    flex-direction: column;
    text-align: center;
    height: auto;
    min-height: 120px;
  }
  
  .stat-icon {
    width: 40px;
    height: 40px;
  }
  
  .action-item {
    flex-direction: column;
    text-align: center;
    gap: var(--spacing-sm);
  }
  
  .action-icon {
    width: 40px;
    height: 40px;
  }
  
  .shortcuts-grid {
    gap: var(--spacing-sm);
  }
  
  .shortcut-link {
    padding: var(--spacing-md);
    min-height: 60px;
  }
}
</style>
