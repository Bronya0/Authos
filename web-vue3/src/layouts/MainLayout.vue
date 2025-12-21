<template>
  <n-layout class="main-layout" has-sider>
    <n-layout-sider v-if="authStore.isAppLoggedIn" bordered collapse-mode="width" :collapsed-width="64" :width="240" :collapsed="collapsed"
      show-trigger @collapse="collapsed = true" @expand="collapsed = false">
      <div class="sidebar-header">
        <n-icon size="24" color="#1890ff">
          <shield-lock-icon />
        </n-icon>
        <span v-if="!collapsed" class="sidebar-title">Authos</span>
      </div>

      <n-menu :collapsed="collapsed" :collapsed-width="64" :collapsed-icon-size="22" :options="menuOptions"
        :value="activeKey" @update:value="handleMenuSelect" />
    </n-layout-sider>

    <n-layout>
      <n-layout-header bordered class="layout-header">
        <div class="header-content">
          <div class="header-left">
            <n-breadcrumb>
              <n-breadcrumb-item>{{ currentPageTitle }}</n-breadcrumb-item>
            </n-breadcrumb>
          </div>

          <div class="header-center">
            <!-- 当前应用显示与切换 -->
            <n-button 
              v-if="authStore.isAppLoggedIn && authStore.currentApp"
              quaternary 
              @click="handleSwitchApp"
            >
              <template #icon>
                <n-icon>
                  <AppsIcon />
                </n-icon>
              </template>
              应用管理中心：{{ authStore.currentApp.name }} ({{ authStore.currentApp.code }})
            </n-button>
          </div>

          <div class="header-right">
            <n-space :size="4">
              <n-button quaternary circle @click="handleRefresh">
                <template #icon>
                  <n-icon>
                    <RefreshIcon />
                  </n-icon>
                </template>
              </n-button>

              <n-button quaternary circle @click="authStore.toggleTheme">
                <template #icon>
                  <n-icon>
                    <component :is="authStore.isDark ? SunnyIcon : MoonIcon" />
                  </n-icon>
                </template>
              </n-button>

              <n-dropdown trigger="click" :options="userDropdownOptions" @select="handleUserDropdown">
                <n-button quaternary>
                  <n-space :size="8">
                    <n-icon><person-circle /></n-icon>
                    <span>{{ displayUsername }}</span>
                  </n-space>
                </n-button>
              </n-dropdown>
            </n-space>
          </div>
        </div>
      </n-layout-header>

      <n-layout-content class="layout-content" content-style="padding: 24px;">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup>
import { ref, computed, h, onMounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useAppStore } from '../stores/app'
import { NIcon } from 'naive-ui'
import {
  ShieldCheckmark as ShieldLockIcon,
  Speedometer as SpeedometerIcon,
  People as PeopleIcon,
  Person as PersonIcon,
  List as ListIcon,
  Key as KeyIcon,
  Apps as AppsIcon,
  AppsSharp as AppsShuffleIcon,
  PeopleCircle as PeopleCircleIcon,
  Book as BookIcon,
  Refresh as RefreshIcon,
  PersonCircle as PersonCircleIcon,
  Sunny as SunnyIcon,
  Moon as MoonIcon,
  LogOut as LogOutIcon
} from '@vicons/ionicons5'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const appStore = useAppStore()

const collapsed = ref(false)
const shieldLock = ShieldLockIcon
const refresh = RefreshIcon
const personCircle = PersonCircleIcon

const menuOptions = computed(() => {
  // 只有登录了应用后才显示用户、角色等管理菜单
  if (authStore.isAppLoggedIn) {
    return [
      {
        label: '仪表盘',
        key: 'Dashboard',
        icon: () => h(NIcon, null, { default: () => h(SpeedometerIcon) })
      },
      {
        label: '用户管理',
        key: 'Users',
        icon: () => h(NIcon, null, { default: () => h(PeopleIcon) })
      },
      {
        label: '角色管理',
        key: 'Roles',
        icon: () => h(NIcon, null, { default: () => h(PersonIcon) })
      },
      {
        label: '菜单管理',
        key: 'Menus',
        icon: () => h(NIcon, null, { default: () => h(ListIcon) })
      },
      {
        label: '权限配置',
        key: 'Permissions',
        icon: () => h(NIcon, null, { default: () => h(KeyIcon) })
      },
      {
        label: '审计日志',
        key: 'AuditLogs',
        icon: () => h(NIcon, null, { default: () => h(ListIcon) })
      },
      {
        label: '接口文档',
        key: 'ApiDocs',
        icon: () => h(NIcon, null, { default: () => h(BookIcon) })
      }
    ]
  }

  // 未登录应用时只显示仪表盘
  return [
    {
      label: '仪表盘',
      key: 'Dashboard',
      icon: () => h(NIcon, null, { default: () => h(SpeedometerIcon) })
    }
  ]
})

const userDropdownOptions = [
  {
    label: '退出登录',
    key: 'logout',
    icon: () => h(NIcon, null, { default: () => h(LogOutIcon) })
  }
]

const activeKey = computed(() => route.name)
const currentPageTitle = computed(() => route.meta?.title || 'Authos')

const displayUsername = computed(() => {
  // 优先显示系统管理员用户名
  if (authStore.systemUser && authStore.systemUser.username) {
    return authStore.systemUser.username
  }
  return 'Admin'
})

const handleMenuSelect = (key) => {
  router.push({ name: key })
}

const handleRefresh = () => {
  appStore.showLoading()
  router.go(0)
  setTimeout(() => appStore.hideLoading(), 500)
}

const handleSwitchApp = () => {
  if (confirm('确定要切换应用吗？当前的操作可能会被中断。')) {
    // 清除应用认证信息，保留系统管理员认证
    authStore.clearAppAuth()
    router.push('/app-selection')
  }
}

const handleUserDropdown = (key) => {
  if (key === 'logout') {
    if (confirm('确定要退出登录吗？')) {
      // 清除所有认证信息
      authStore.logout()
      appStore.clearCurrentApp()
      router.push('/system-login')
    }
  }
}

// 组件挂载时检查认证状态
onMounted(() => {
  // 等待DOM更新后检查
  nextTick(() => {
    console.log('检查认证状态:', {
      isSystemLoggedIn: authStore.isSystemLoggedIn,
      isAppLoggedIn: authStore.isAppLoggedIn,
      currentApp: authStore.currentApp,
      currentPath: route.path
    })

    // 如果系统管理员未登录，跳转到系统登录页面
    if (!authStore.isSystemLoggedIn && route.path !== '/system-login') {
      console.log('跳转到系统登录页面')
      router.push('/system-login')
      return
    }

    // 如果系统管理员已登录但未登录应用，跳转到应用选择页面
    if (authStore.isSystemLoggedIn && !authStore.isAppLoggedIn &&
      route.path !== '/app-selection' && route.path !== '/app-login') {
      console.log('跳转到应用选择页面')
      router.push('/app-selection')
    }
  })
})
</script>

<style scoped>
.main-layout {
  height: 100vh;
}

.sidebar-header {
  display: flex;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid var(--n-border-color);
}

.sidebar-title {
  margin-left: 12px;
  font-size: 18px;
  font-weight: bold;
  color: #1890ff;
}

.layout-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  height: 64px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  gap: 24px;
}

.header-left {
  flex: 0 0 auto;
  display: flex;
  align-items: center;
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.header-right {
  flex: 0 0 auto;
  display: flex;
  align-items: center;
}

.layout-content {
  height: calc(100vh - 64px);
  overflow-y: auto;
}
</style>