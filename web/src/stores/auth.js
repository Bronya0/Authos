import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAppStore } from './app'

export const useAuthStore = defineStore('auth', () => {
  const appStore = useAppStore()
  
  // 系统管理员认证
  const systemToken = ref(localStorage.getItem('systemToken') || '')
  const systemUser = ref(JSON.parse(localStorage.getItem('systemUser') || '{}'))
  
  // 应用认证
  const appToken = ref(localStorage.getItem('appToken') || '')
  const currentApp = ref(JSON.parse(localStorage.getItem('currentApp') || '{}'))
  
  // 主题设置
  const isDark = ref(localStorage.getItem('isDark') === 'true')

  // 计算属性
  const isSystemLoggedIn = computed(() => !!systemToken.value)
  const isAppLoggedIn = computed(() => !!appToken.value && !!currentApp.value)
  const isLoggedIn = computed(() => isSystemLoggedIn.value && isAppLoggedIn.value)
  
  const currentUser = computed(() => {
    if (!systemUser.value || typeof systemUser.value !== 'object') return null
    return systemUser.value
  })
  
  const getUsername = computed(() => {
    if (!systemUser.value) return 'Admin'
    
    if (typeof systemUser.value === 'string') {
      return systemUser.value || 'Admin'
    }
    
    if (typeof systemUser.value === 'object' && systemUser.value !== null) {
      return systemUser.value.username || systemUser.value.Username || systemUser.value.name || systemUser.value.Name || 'Admin'
    }
    
    return 'Admin'
  })

  // 系统管理员认证
  function setSystemAuth(userData, userToken) {
    systemToken.value = userToken
    
    if (typeof userData === 'string') {
      try {
        const parsedUser = JSON.parse(userData)
        systemUser.value = parsedUser
      } catch (e) {
        systemUser.value = { username: userData }
      }
    } else {
      systemUser.value = userData || {}
    }
    
    localStorage.setItem('systemToken', userToken)
    localStorage.setItem('systemUser', JSON.stringify(systemUser.value))
  }

  // 应用认证
  function setAppAuth(appData, appTokenValue) {
    appToken.value = appTokenValue
    currentApp.value = appData || {}
    
    localStorage.setItem('appToken', appTokenValue)
    localStorage.setItem('currentApp', JSON.stringify(currentApp.value))
  }

  // 用户认证（多租户）
  function setAuth(userData, userToken, appData) {
    // 设置用户token到Authorization头
    localStorage.setItem('userToken', userToken)
    
    // 设置应用信息
    if (appData) {
      currentApp.value = appData
      localStorage.setItem('currentApp', JSON.stringify(currentApp.value))
    }
  }

  // 清除应用认证
  function clearAppAuth() {
    appToken.value = ''
    currentApp.value = {}
    localStorage.removeItem('appToken')
    localStorage.removeItem('currentApp')
  }

  // 完全登出
  function logout() {
    systemToken.value = ''
    systemUser.value = {}
    appToken.value = ''
    currentApp.value = {}
    
    localStorage.removeItem('systemToken')
    localStorage.removeItem('systemUser')
    localStorage.removeItem('appToken')
    localStorage.removeItem('currentApp')
    localStorage.removeItem('userToken')
  }

  // 主题切换
  function toggleTheme() {
    isDark.value = !isDark.value
    localStorage.setItem('isDark', isDark.value.toString())
  }

  return {
    // 系统认证
    systemToken,
    systemUser,
    isSystemLoggedIn,
    
    // 应用认证
    appToken,
    currentApp,
    isAppLoggedIn,
    
    // 通用
    isDark,
    isLoggedIn,
    currentUser,
    getUsername,
    
    // 方法
    setSystemAuth,
    setAppAuth,
    setAuth,
    clearAppAuth,
    logout,
    toggleTheme,
    showError: appStore.showError,
    showSuccess: appStore.showSuccess
  }
})