import type { App } from 'vue'
import { ElMessage } from 'element-plus'
import type { Router } from 'vue-router'
import { ApiError } from './request'

interface ErrorHandlerOptions {
  router: Router
}

export function setupErrorHandler(app: App, options: ErrorHandlerOptions) {
  const { router } = options

  // Vue错误处理
  app.config.errorHandler = (error, instance, info) => {
    console.error('Vue Error:', error)
    console.error('Error Info:', info)
    
    if (error instanceof ApiError) {
      // API错误已经在request中处理
      return
    }
    
    ElMessage.error('系统错误，请刷新页面重试')
  }

  // 全局Promise错误处理
  window.addEventListener('unhandledrejection', event => {
    console.error('Unhandled Promise Rejection:', event.reason)
    
    if (event.reason instanceof ApiError) {
      // API错误已经在request中处理
      return
    }
    
    ElMessage.error('操作失败，请稍后重试')
  })

  // 全局JS错误处理
  window.addEventListener('error', event => {
    console.error('Global Error:', event.error)
    
    if (event.error instanceof ApiError) {
      // API错误已经在request中处理
      return
    }
    
    ElMessage.error('系统错误，请刷新页面重试')
  })

  // 路由错误处理
  router.onError((error) => {
    console.error('Router Error:', error)
    ElMessage.error('页面加载失败，请刷新重试')
  })
}

// 安全相关工具函数
export const security = {
  // XSS防护：HTML转义
  escapeHtml(html: string): string {
    const div = document.createElement('div')
    div.textContent = html
    return div.innerHTML
  },

  // XSS防护：URL编码
  escapeUrl(url: string): string {
    return encodeURIComponent(url)
  },

  // 密码强度检查
  checkPasswordStrength(password: string): {
    score: number;
    feedback: string[];
  } {
    const feedback: string[] = []
    let score = 0

    // 长度检查
    if (password.length < 6) {
      feedback.push('密码长度至少6位')
    } else if (password.length >= 8) {
      score += 1
    }

    // 包含数字
    if (/\d/.test(password)) {
      score += 1
    } else {
      feedback.push('建议包含数字')
    }

    // 包含小写字母
    if (/[a-z]/.test(password)) {
      score += 1
    } else {
      feedback.push('建议包含小写字母')
    }

    // 包含大写字母
    if (/[A-Z]/.test(password)) {
      score += 1
    } else {
      feedback.push('建议包含大写字母')
    }

    // 包含特殊字符
    if (/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
      score += 1
    } else {
      feedback.push('建议包含特殊字符')
    }

    return {
      score,
      feedback: feedback.length ? feedback : ['密码强度良好']
    }
  },

  // CSRF Token管理
  csrf: {
    getToken(): string {
      return localStorage.getItem('csrf_token') || ''
    },

    setToken(token: string): void {
      localStorage.setItem('csrf_token', token)
    },

    removeToken(): void {
      localStorage.removeItem('csrf_token')
    }
  }
} 