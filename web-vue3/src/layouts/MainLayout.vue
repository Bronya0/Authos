<template>
  <n-layout class="main-layout" has-sider>
    <n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="240" :collapsed="collapsed"
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
        <div class="header-left">
          <n-breadcrumb>
            <n-breadcrumb-item>{{ currentPageTitle }}</n-breadcrumb-item>
          </n-breadcrumb>
        </div>

        <div class="header-right">
          <n-space>
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
                <n-space>
                  <n-icon><person-circle /></n-icon>
                  <span>{{ displayUsername }}</span>
                </n-space>
              </n-button>
            </n-dropdown>
          </n-space>
        </div>
      </n-layout-header>

      <n-layout-content class="layout-content" content-style="padding: 24px;">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup>
import { ref, computed, h } from 'vue'
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

const menuOptions = [
  {
    label: '仪表盘',
    key: 'Dashboard',
    icon: () => h(NIcon, null, { default: () => h(SpeedometerIcon) })
  },
  {
    label: '应用管理',
    key: 'Applications',
    icon: () => h(NIcon, null, { default: () => h(AppsIcon) })
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
    label: '接口文档',
    key: 'ApiDocs',
    icon: () => h(NIcon, null, { default: () => h(BookIcon) })
  }
]

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
  return authStore.getUsername || 'User'
})

const handleMenuSelect = (key) => {
  router.push({ name: key })
}

const handleRefresh = () => {
  appStore.showLoading()
  router.go(0)
  setTimeout(() => appStore.hideLoading(), 500)
}

const handleUserDropdown = (key) => {
  if (key === 'logout') {
    if (confirm('确定要退出登录吗？')) {
      authStore.logout()
      router.push('/login')
    }
  }
}
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

.header-left {
  flex: 1;
}

.header-right {
  display: flex;
  align-items: center;
}

.layout-content {
  height: calc(100vh - 64px);
  overflow-y: auto;
}
</style>