import { defineStore } from 'pinia'
import { ref } from 'vue'

// Global message instance that will be set by the App component
let globalMessage = null

export function setGlobalMessage(message) {
  globalMessage = message
}

export const useAppStore = defineStore('app', () => {
  const loading = ref(false)

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
    showLoading,
    hideLoading,
    showSuccess,
    showError,
    showWarning,
    showInfo
  }
})