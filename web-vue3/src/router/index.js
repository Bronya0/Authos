import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/system-login',
    name: 'SystemLogin',
    component: () => import('../views/SystemLogin.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/app-login',
    name: 'AppLogin',
    component: () => import('../views/AppLogin.vue'),
    meta: { requiresSystemAuth: true }
  },
  {
    path: '/app-selection',
    name: 'AppSelection',
    component: () => import('../views/AppSelection.vue'),
    meta: { requiresSystemAuth: true }
  },
  {
    path: '/login',
    redirect: '/system-login'
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('../layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { title: '仪表盘', icon: 'speedometer' }
      },
      
      {
        path: 'users',
        name: 'Users',
        component: () => import('../views/Users.vue'),
        meta: { title: '用户管理', icon: 'people' }
      },
      {
        path: 'roles',
        name: 'Roles',
        component: () => import('../views/Roles.vue'),
        meta: { title: '角色管理', icon: 'person' }
      },
      {
        path: 'menus',
        name: 'Menus',
        component: () => import('../views/Menus.vue'),
        meta: { title: '菜单管理', icon: 'list' }
      },
      {
        path: 'permissions',
        name: 'Permissions',
        component: () => import('../views/Permissions.vue'),
        meta: { title: '权限配置', icon: 'key' }
      },
      {
        path: 'api-docs',
        name: 'ApiDocs',
        component: () => import('../views/ApiDocs.vue'),
        meta: { title: '接口文档', icon: 'book' }
      }
    ]
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫 - 支持多级认证
router.beforeEach((to, from, next) => {
  // 检查是否需要系统管理员认证
  if (to.meta.requiresSystemAuth) {
    const systemToken = localStorage.getItem('systemToken')
    if (!systemToken) {
      next('/system-login')
      return
    }
  }
  
  // 检查是否需要完整认证（系统+应用）
  if (to.meta.requiresAuth) {
    const systemToken = localStorage.getItem('systemToken')
    const appToken = localStorage.getItem('appToken')
    
    if (!systemToken) {
      next('/system-login')
      return
    }
    
    if (!appToken) {
      next('/app-selection')
      return
    }
  }
  
  // 其他情况继续导航
  next()
})

export default router