import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'
import router from './router'
import { setupErrorHandler } from './utils/error-handler'

// 导入统一设计系统
import './styles/design-system.css'

const app = createApp(App)

// 注册Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 配置Element Plus
app.use(ElementPlus, {
  locale: zhCn,
})

// 配置Pinia
const pinia = createPinia()
app.use(pinia)

// 配置路由
app.use(router)

// 配置错误处理
setupErrorHandler(app, { router })

// 初始化认证状态
import { useAuthStore } from './stores/auth'
const authStore = useAuthStore()
authStore.initAuth()

app.mount('#app')
