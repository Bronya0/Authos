import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('currentUser') || '{}'))
  const app = ref(JSON.parse(localStorage.getItem('appInfo') || '{}'))
  const appId = ref(localStorage.getItem('appId') || '')
  const isDark = ref(localStorage.getItem('isDark') === 'true')

  const isLoggedIn = computed(() => !!token.value)
  const currentUser = computed(() => {
    if (!user.value || typeof user.value !== 'object') return null
    return user.value
  })
  
  const currentApp = computed(() => {
    if (!app.value || typeof app.value !== 'object') return null
    return app.value
  })

  // 获取用户名字符串的辅助函数
  const getUsername = computed(() => {
    if (!user.value) return 'User'
    
    // 如果是字符串，可能是JSON字符串或直接用户名
    if (typeof user.value === 'string') {
      return user.value || 'User'
    }
    
    // 如果是对象，提取用户名字段
    if (typeof user.value === 'object' && user.value !== null) {
      return user.value.username || user.value.Username || user.value.name || user.value.Name || 'User'
    }
    
    return 'User'
  })

  function setAuth(userData, userToken, appData = null) {
    token.value = userToken
    
    // 确保userData是正确的格式
    if (typeof userData === 'string') {
      // 如果userData是字符串，可能是JSON字符串
      try {
        const parsedUser = JSON.parse(userData)
        user.value = parsedUser
      } catch (e) {
        // 如果解析失败，直接使用字符串作为用户名
        user.value = { username: userData }
      }
    } else {
      user.value = userData || {}
    }
    
    // 设置应用信息
    if (appData) {
      app.value = appData
      appId.value = appData.id || appData.ID || ''
      localStorage.setItem('appInfo', JSON.stringify(app.value))
      localStorage.setItem('appId', appId.value.toString())
    }
    
    localStorage.setItem('token', userToken)
    localStorage.setItem('currentUser', JSON.stringify(user.value))
  }

  function logout() {
    token.value = ''
    user.value = {}
    app.value = {}
    appId.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('currentUser')
    localStorage.removeItem('appInfo')
    localStorage.removeItem('appId')
  }

  function toggleTheme() {
    isDark.value = !isDark.value
    localStorage.setItem('isDark', isDark.value.toString())
  }

  return {
    token,
    user,
    app,
    appId,
    isDark,
    isLoggedIn,
    currentUser,
    currentApp,
    getUsername,
    setAuth,
    logout,
    toggleTheme
  }
})