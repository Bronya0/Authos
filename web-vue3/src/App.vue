<template>
  <n-config-provider :theme="pageTheme">
    <n-message-provider>
      <n-modal-provider>
        <AppContent />
      </n-modal-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { darkTheme } from 'naive-ui'
import { useAuthStore } from './stores/auth'
import AppContent from './AppContent.vue'

const route = useRoute()
const authStore = useAuthStore()

// 判断当前是否在登录后的界面（非登录相关页面）
const isLoginPage = computed(() => {
  const loginPages = ['Login', 'SystemLogin', 'AppLogin', 'AppSelection']
  return loginPages.includes(route.name)
})

// 只在非登录页面应用主题
const pageTheme = computed(() => {
  return isLoginPage.value ? null : (authStore.isDark ? darkTheme : null)
})
</script>