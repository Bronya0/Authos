import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

// Global message instance that will be set by the App component
let globalMessage = null

export function setGlobalMessage(message) {
  globalMessage = message
}

export const useAppStore = defineStore('app', () => {
  const loading = ref(false)
  const currentApp = ref(null) // 当前选中的应用
  const hasSelectedApp = computed(() => !!currentApp.value) // 是否已选择应用

  // 设置当前应用
  function setCurrentApp(app) {
    console.log('设置当前应用:', app)
    currentApp.value = app
    if (app) {
      localStorage.setItem('currentApp', JSON.stringify(app))
      localStorage.setItem('appId', app.id)
      localStorage.setItem('appCode', app.code)
    } else {
      localStorage.removeItem('currentApp')
      localStorage.removeItem('appId')
      localStorage.removeItem('appCode')
    }
  }

  // 从本地存储加载应用
  function loadCurrentApp() {
    const appStr = localStorage.getItem('currentApp')
    console.log('从本地存储加载应用:', appStr)
    if (appStr) {
      try {
        currentApp.value = JSON.parse(appStr)
        console.log('成功加载应用:', currentApp.value)
      } catch (e) {
        console.error('加载当前应用失败', e)
        currentApp.value = null
      }
    } else {
      console.log('本地存储中没有应用信息')
      currentApp.value = null
    }
  }

  // 清除当前应用
  function clearCurrentApp() {
    setCurrentApp(null)
  }

  function showLoading() {
    loading.value = true
  }

  function hideLoading() {
    loading.value = false
  }

  function showSuccess(msg) {
    if (globalMessage) globalMessage.success(msg)
  }

  function showError(msg) {
    if (globalMessage) globalMessage.error(msg)
  }

  function showWarning(msg) {
    if (globalMessage) globalMessage.warning(msg)
  }

  function showInfo(msg) {
    if (globalMessage) globalMessage.info(msg)
  }

  return {
    loading,
    currentApp,
    hasSelectedApp,
    setCurrentApp,
    loadCurrentApp,
    clearCurrentApp,
    showLoading,
    hideLoading,
    showSuccess,
    showError,
    showWarning,
    showInfo
  }
})