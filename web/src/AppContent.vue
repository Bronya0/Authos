<template>
  <router-view />
</template>

<script setup>
import { useMessage } from 'naive-ui'
import { setGlobalMessage } from './stores/app'

// Initialize global message - safe to use here as we're inside message provider
const message = useMessage()
setGlobalMessage(message)
</script>