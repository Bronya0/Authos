import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { requiresAuth: false }
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
        meta: { title: '角色管理', icon: 'person-badge' }
      },
      {
        path: 'menus',
        name: 'Menus',
        component: () => import('../views/Menus.vue'),
        meta: { title: '菜单管理', icon: 'list-ul' }
      },
      {
        path: 'permissions',
        name: 'Permissions',
        component: () => import('../views/Permissions.vue'),
        meta: { title: '权限配置', icon: 'key' }
      },
      {
        path: 'applications',
        name: 'Applications',
        component: () => import('../views/Applications.vue'),
        meta: { title: '应用管理', icon: 'apps' }
      },
      {
        path: 'api-docs',
        name: 'ApiDocs',
        component: () => import('../views/ApiDocs.vue'),
        meta: { title: '接口文档', icon: 'book' }
      },
      {
        path: 'test',
        name: 'Test',
        component: () => import('../views/Test.vue'),
        meta: { title: '测试', icon: 'test' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next('/login')
  } else if (to.path === '/login' && authStore.isLoggedIn) {
    next('/')
  } else if (to.path === '/' && authStore.isLoggedIn) {
    // 根路径重定向到仪表盘
    next('/dashboard')
  } else {
    next()
  }
})

export default router