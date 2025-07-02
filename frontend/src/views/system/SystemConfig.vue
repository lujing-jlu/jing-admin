<template>
  <div class="page-container fade-in">
    <div class="page-header">
      <div>
        <h1 class="page-title">系统配置</h1>
        <p class="page-subtitle">管理系统的基础设置、邮件配置、安全策略等</p>
      </div>
      <div class="page-actions">
        <el-button :icon="Refresh" @click="refreshConfigs">刷新</el-button>
        <el-button type="primary" :icon="Check" @click="saveAllConfigs" :loading="saving">
          保存所有更改
        </el-button>
      </div>
    </div>

    <!-- 配置统计 -->
    <div class="stats-overview">
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon" style="background: linear-gradient(45deg, #667eea 0%, #764ba2 100%)">
              <el-icon><Setting /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ totalConfigs }}</div>
              <div class="stat-label">配置项总数</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon" style="background: linear-gradient(45deg, #f093fb 0%, #f5576c 100%)">
              <el-icon><Edit /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ editableConfigs }}</div>
              <div class="stat-label">可编辑配置</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon" style="background: linear-gradient(45deg, #4facfe 0%, #00f2fe 100%)">
              <el-icon><View /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ publicConfigs }}</div>
              <div class="stat-label">公开配置</div>
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="stat-card">
            <div class="stat-icon" style="background: linear-gradient(45deg, #43e97b 0%, #38f9d7 100%)">
              <el-icon><DataLine /></el-icon>
            </div>
            <div class="stat-content">
              <div class="stat-value">{{ Object.keys(changedConfigs).length }}</div>
              <div class="stat-label">待保存更改</div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 配置分类标签 -->
    <div class="category-tabs">
      <el-tabs v-model="activeCategory" @tab-click="handleCategoryChange">
        <el-tab-pane label="基础配置" name="basic">
          <template #label>
            <span><el-icon><House /></el-icon> 基础配置</span>
          </template>
        </el-tab-pane>
        <el-tab-pane label="邮件配置" name="mail">
          <template #label>
            <span><el-icon><Message /></el-icon> 邮件配置</span>
          </template>
        </el-tab-pane>
        <el-tab-pane label="安全配置" name="security">
          <template #label>
            <span><el-icon><Lock /></el-icon> 安全配置</span>
          </template>
        </el-tab-pane>
        <el-tab-pane label="系统配置" name="system">
          <template #label>
            <span><el-icon><Monitor /></el-icon> 系统配置</span>
          </template>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 配置表单 -->
    <div v-loading="loading" class="config-forms">
      <div v-for="(configs, category) in groupedConfigs" :key="category">
        <div v-if="activeCategory === 'all' || activeCategory === category" class="config-category">
          <div class="category-content">
            <el-form :model="formData" label-width="140px" size="large">
              <div v-for="config in configs" :key="config.key" class="config-item">
                <el-form-item :label="config.display_name">
                  <template #label>
                    <div class="config-label">
                      <span>{{ config.display_name }}</span>
                      <el-tooltip v-if="config.description" :content="config.description" placement="top">
                        <el-icon class="help-icon"><QuestionFilled /></el-icon>
                      </el-tooltip>
                      <el-tag v-if="!config.is_editable" size="small" type="info">只读</el-tag>
                      <el-tag v-if="config.is_public" size="small" type="success">公开</el-tag>
                    </div>
                  </template>

                  <!-- 根据配置类型渲染不同的输入控件 -->
                  <div class="config-input">
                    <!-- 布尔类型 -->
                    <el-switch
                      v-if="config.type === 'boolean'"
                      v-model="formData[config.key]"
                      :disabled="!config.is_editable"
                      active-text="启用"
                      inactive-text="禁用"
                      :active-value="true"
                      :inactive-value="false"
                      @change="onConfigChange(config.key, $event)"
                    />
                    
                    <!-- 数字类型 -->
                    <el-input-number
                      v-else-if="config.type === 'number'"
                      v-model="formData[config.key]"
                      :disabled="!config.is_editable"
                      :controls="true"
                      :min="0"
                      style="width: 200px"
                      @change="onConfigChange(config.key, $event)"
                    />
                    
                    <!-- 密码类型 -->
                    <el-input
                      v-else-if="config.key.includes('password')"
                      v-model="formData[config.key]"
                      type="password"
                      :disabled="!config.is_editable"
                      placeholder="留空表示不更改"
                      show-password
                      @input="onConfigChange(config.key, $event)"
                    />
                    
                    <!-- 文本域类型 -->
                    <el-input
                      v-else-if="config.key.includes('description') || config.key.includes('bio')"
                      v-model="formData[config.key]"
                      type="textarea"
                      :disabled="!config.is_editable"
                      :rows="3"
                      @input="onConfigChange(config.key, $event)"
                    />
                    
                    <!-- 普通文本类型 -->
                    <el-input
                      v-else
                      v-model="formData[config.key]"
                      :disabled="!config.is_editable"
                      @input="onConfigChange(config.key, $event)"
                    />
                    
                    <!-- 重置按钮 -->
                    <el-button
                      v-if="config.is_editable && changedConfigs[config.key]"
                      link
                      type="primary"
                      :icon="RefreshLeft"
                      @click="resetConfig(config.key)"
                    >
                      重置
                    </el-button>
                  </div>

                  <div v-if="config.description" class="config-description">
                    {{ config.description }}
                  </div>
                </el-form-item>
              </div>
            </el-form>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部操作栏 -->
    <div v-if="Object.keys(changedConfigs).length > 0" class="bottom-actions">
      <div class="changed-configs">
        <span>已修改 {{ Object.keys(changedConfigs).length }} 项配置</span>
        <el-button @click="discardChanges">放弃更改</el-button>
        <el-button type="primary" @click="saveAllConfigs" :loading="saving">
          保存更改
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Setting,
  Edit,
  View,
  DataLine,
  House,
  Message,
  Lock,
  Monitor,
  QuestionFilled,
  Refresh,
  Check,
  RefreshLeft
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

// 类型定义
interface SystemConfig {
  id: number
  key: string
  value: string
  type: string
  category: string
  display_name: string
  description: string
  is_public: boolean
  is_editable: boolean
}

// 数据
const authStore = useAuthStore()
const loading = ref(false)
const saving = ref(false)
const activeCategory = ref('basic')
const configs = ref<SystemConfig[]>([])
const formData = reactive<Record<string, string | number | boolean>>({})
const changedConfigs = reactive<Record<string, any>>({})
const originalValues = reactive<Record<string, any>>({})

// 计算属性
const groupedConfigs = computed(() => {
  const grouped: Record<string, SystemConfig[]> = {}
  configs.value.forEach(config => {
    if (!grouped[config.category]) {
      grouped[config.category] = []
    }
    grouped[config.category].push(config)
  })
  return grouped
})

const totalConfigs = computed(() => configs.value.length)
const editableConfigs = computed(() => configs.value.filter(c => c.is_editable).length)
const publicConfigs = computed(() => configs.value.filter(c => c.is_public).length)

// 方法
const getConfigs = async () => {
  loading.value = true
  try {
    const response = await fetch('/api/config', {
      headers: {
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (!response.ok) {
      throw new Error('获取配置失败')
    }
    
    const data = await response.json()
    configs.value = data.data.configs
    
    // 初始化表单数据
    configs.value.forEach(config => {
      let value = config.value
      
      // 类型转换
      if (config.type === 'boolean') {
        value = config.value === 'true'
      } else if (config.type === 'number') {
        value = parseFloat(config.value) || 0
      }
      
      formData[config.key] = value
      originalValues[config.key] = value
    })
  } catch (error) {
    console.error('获取配置错误:', error)
    ElMessage.error('获取系统配置失败')
  } finally {
    loading.value = false
  }
}

const onConfigChange = (key: string, value: any) => {
  if (value !== originalValues[key]) {
    changedConfigs[key] = value
  } else {
    delete changedConfigs[key]
  }
}

const resetConfig = (key: string) => {
  formData[key] = originalValues[key]
  delete changedConfigs[key]
}

const discardChanges = () => {
  ElMessageBox.confirm(
    '确定要放弃所有未保存的更改吗？',
    '确认放弃',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    // 恢复所有更改
    Object.keys(changedConfigs).forEach(key => {
      formData[key] = originalValues[key]
    })
    
    // 清空更改记录
    Object.keys(changedConfigs).forEach(key => {
      delete changedConfigs[key]
    })
    
    ElMessage.success('已放弃所有更改')
  }).catch(() => {
    // 取消操作
  })
}

const saveAllConfigs = async () => {
  if (Object.keys(changedConfigs).length === 0) {
    ElMessage.info('没有需要保存的更改')
    return
  }

  saving.value = true
  try {
         // 准备保存数据
     const saveData: Record<string, string> = {}
     Object.keys(changedConfigs).forEach(key => {
       let value: string = String(formData[key])
       // 转换为字符串
               if (typeof formData[key] === 'boolean') {
          value = (formData[key] as boolean).toString()
        } else if (typeof formData[key] === 'number') {
          value = (formData[key] as number).toString()
        } else {
          value = formData[key] as string
        }
       saveData[key] = value
     })

    const response = await fetch('/api/config/batch', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ configs: saveData })
    })

    if (!response.ok) {
      throw new Error('保存配置失败')
    }

    // 更新原始值
    Object.keys(changedConfigs).forEach(key => {
      originalValues[key] = formData[key]
    })

    // 清空更改记录
    Object.keys(changedConfigs).forEach(key => {
      delete changedConfigs[key]
    })

    ElMessage.success('系统配置保存成功')
  } catch (error) {
    console.error('保存配置错误:', error)
    ElMessage.error('保存系统配置失败')
  } finally {
    saving.value = false
  }
}

const refreshConfigs = () => {
  if (Object.keys(changedConfigs).length > 0) {
    ElMessageBox.confirm(
      '刷新配置将丢失未保存的更改，确定要继续吗？',
      '确认刷新',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    ).then(() => {
      getConfigs()
    })
  } else {
    getConfigs()
  }
}

const handleCategoryChange = () => {
  // 标签页切换逻辑
}

// 生命周期
onMounted(() => {
  getConfigs()
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
  font-size: 24px;
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

.category-tabs {
  margin-bottom: var(--spacing-xl);
}

.category-tabs :deep(.el-tabs__item) {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.config-forms {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-secondary);
  overflow: hidden;
}

.config-category {
  border-bottom: 1px solid var(--border-secondary);
}

.config-category:last-child {
  border-bottom: none;
}

.category-content {
  padding: var(--spacing-xl);
}

.config-item {
  margin-bottom: var(--spacing-lg);
  padding: var(--spacing-lg);
  border: 1px solid var(--border-secondary);
  border-radius: var(--radius-md);
  transition: all var(--duration-normal) ease;
}

.config-item:hover {
  border-color: var(--primary-color);
  background: var(--bg-hover);
}

.config-label {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.help-icon {
  color: var(--text-tertiary);
  cursor: help;
}

.config-input {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.config-input .el-input,
.config-input .el-input-number {
  flex: 1;
  max-width: 400px;
}

.config-description {
  font-size: var(--font-sm);
  color: var(--text-secondary);
  margin-top: var(--spacing-xs);
  line-height: 1.5;
}

.bottom-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: var(--bg-card);
  border-top: 1px solid var(--border-secondary);
  padding: var(--spacing-lg) var(--spacing-xl);
  backdrop-filter: blur(10px);
  z-index: 1000;
}

.changed-configs {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 1200px;
  margin: 0 auto;
}

.changed-configs span {
  color: var(--text-secondary);
  font-size: var(--font-sm);
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
  .stat-card {
    flex-direction: column;
    text-align: center;
    gap: var(--spacing-md);
  }

  .stat-icon {
    width: 48px;
    height: 48px;
    font-size: 20px;
  }

  .config-input {
    flex-direction: column;
    align-items: stretch;
    gap: var(--spacing-sm);
  }

  .config-input .el-input,
  .config-input .el-input-number {
    max-width: none;
  }

  .changed-configs {
    flex-direction: column;
    gap: var(--spacing-md);
    text-align: center;
  }
}
</style> 