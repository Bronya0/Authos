<template>
  <n-config-provider :theme="theme">
    <n-message-provider>
      <n-modal-provider>
        <AppContent />
      </n-modal-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup>
import { computed } from 'vue'
import { darkTheme } from 'naive-ui'
import { useAuthStore } from './stores/auth'
import AppContent from './AppContent.vue'

const authStore = useAuthStore()
const theme = computed(() => authStore.isDark ? darkTheme : null)
</script>