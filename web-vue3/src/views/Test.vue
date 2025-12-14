<template>
  <div style="padding: 20px;">
    <h1>测试页面</h1>
    <p>如果看到这个页面，说明路由正常工作</p>
    <n-button @click="goToLogin">跳转到登录页</n-button>
    <n-button @click="goToDashboard">跳转到仪表盘</n-button>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goToLogin = () => {
  router.push('/login')
}

const goToDashboard = () => {
  router.push('/')
}
</script>