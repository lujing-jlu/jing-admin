import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { title: '仪表盘', requiresAuth: true }
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
      meta: { title: '关于' }
    },
    // 个人中心路由
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { title: '个人中心', requiresAuth: true }
    },
    // 用户管理路由
    {
      path: '/users',
      name: 'users',
      children: [
        {
          path: '',
          name: 'user-list',
          component: () => import('../views/user/UserList.vue'),
          meta: { title: '用户列表', requiresAuth: true }
        },
        {
          path: 'roles',
          name: 'roles',
          component: () => import('../views/user/RoleManagement.vue'),
          meta: { title: '角色管理', requiresAuth: true }
        },
        {
          path: 'permissions',
          name: 'permissions',
          component: () => import('../views/user/PermissionConfig.vue'),
          meta: { title: '权限配置', requiresAuth: true }
        }
      ]
    },
    // 系统管理路由
    {
      path: '/system',
      name: 'system',
      children: [
        {
          path: 'config',
          name: 'system-config',
          component: () => import('../views/system/SystemConfig.vue'),
          meta: { title: '系统配置', requiresAuth: true }
        },
        {
          path: 'logs',
          name: 'operation-logs',
          component: () => import('../views/system/OperationLogs.vue'),
          meta: { title: '操作日志', requiresAuth: true }
        },
        {
          path: 'backup',
          name: 'data-backup',
          component: () => import('../views/system/DataBackup.vue'),
          meta: { title: '数据备份', requiresAuth: true }
        }
      ]
    },
    // 数据分析路由
    {
      path: '/analytics',
      name: 'analytics',
      component: () => import('../views/analytics/DataAnalytics.vue'),
      meta: { title: '数据分析', requiresAuth: true }
    },
    // 登录页面
    {
      path: '/login',
      name: 'login',
      component: () => import('../layouts/AuthLayout.vue'),
      meta: { title: '登录', requiresGuest: true },
      children: [
        {
          path: '',
          component: () => import('../views/LoginView.vue')
        }
      ]
    },
    // 404页面
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('../views/NotFound.vue'),
      meta: { title: '页面未找到' }
    }
  ],
})

// 路由守卫 - 认证检查和页面标题
router.beforeEach(async (to, from, next) => {
  // 设置页面标题
  if (to.meta?.title) {
    document.title = `${to.meta.title} - Jing Admin`
  } else {
    document.title = 'Jing Admin'
  }

  // 直接用localStorage判断登录态
  const localToken = localStorage.getItem('token')
  const isAuthenticated = !!localToken
  const requiresAuth = to.meta?.requiresAuth === true  // 明确指定需要认证的路由
  const requiresGuest = to.meta?.requiresGuest === true // 仅游客访问（如登录页）

  // 调试日志
  console.log('[路由守卫]', { path: to.path, isAuthenticated, requiresAuth, requiresGuest })

  // 已登录用户访问游客页面（如登录页）
  if (requiresGuest && isAuthenticated) {
    next('/')
    return
  }

  // 未登录用户访问需要认证的页面
  if (requiresAuth && !isAuthenticated) {
    next('/login')
    return
  }

  // 其他情况正常放行
  next()
})

export default router
