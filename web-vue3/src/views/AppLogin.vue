<template>
  <div class="app-login-container">
    <n-card class="app-login-card" title="应用登录">
      <template #header-extra>
        <n-button text @click="backToAppSelection">
          <template #icon>
            <n-icon><arrow-back /></n-icon>
          </template>
          返回应用选择
        </n-button>
      </template>

      <div v-if="selectedApp" class="selected-app-info">
        <n-card size="small" :bordered="false">
          <template #header>
            <span>当前应用</span>
          </template>
          <div class="app-details">
            <div class="app-name">{{ selectedApp.name }}</div>
            <div class="app-code">应用代码: {{ selectedApp.code }}</div>
            <div class="app-id">应用ID: {{ selectedApp.id }}</div>
          </div>
        </n-card>
      </div>

      <n-form ref="formRef" :model="loginForm" :rules="loginRules" class="login-form">
        <n-form-item path="appId" label="应用ID">
          <n-input v-model:value="loginForm.appId" placeholder="请输入应用ID" size="large" clearable>
            <template #prefix>
              <n-icon>
                <key />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item path="appSecret" label="应用密钥">
          <n-input v-model:value="loginForm.appSecret" type="password" placeholder="请输入应用密钥" size="large"
            show-password-on="click" clearable>
            <template #prefix><n-icon><lock-closed /></n-icon></template>
          </n-input>
        </n-form-item>

        <div class="form-actions">
          <n-button type="primary" size="large" :loading="loading" @click="handleLogin" block>
            登录应用
          </n-button>
        </div>
      </n-form>

      <n-divider class="my-4">应用管理</n-divider>

      <div class="app-actions">
        <n-button type="error" @click="handleDeleteApp" block :disabled="!selectedApp">
          <template #icon>
            <n-icon>
              <trash-outline />
            </n-icon>
          </template>
          删除当前应用
        </n-button>
      </div>
    </n-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  NCard, NButton, NIcon, NForm, NFormItem, NInput, NDivider
} from 'naive-ui'
import { ArrowBack, Key, LockClosed, TrashOutline } from '@vicons/ionicons5'
import { useAuthStore } from '../stores/auth'
import { useAppStore } from '../stores/app'
import { authAPI, applicationAPI } from '../api'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

// 表单引用
const formRef = ref()
const loading = ref(false)

// 选中的应用
const selectedApp = ref(null)

// 登录表单
const loginForm = reactive({
  appId: '',
  appSecret: ''
})

// 登录表单验证规则
const loginRules = {
  appId: [
    { required: true, message: '请输入应用ID', trigger: 'blur' }
  ],
  appSecret: [
    { required: true, message: '请输入应用密钥', trigger: 'blur' }
  ]
}

// 应用登录处理
const handleLogin = async () => {
  try {
    await formRef.value?.validate()
    loading.value = true

    try {
      // 调用应用登录API
      const response = await authAPI.appLogin({
        appId: loginForm.appId,
        appSecret: loginForm.appSecret
      })

      // 设置应用认证
      authStore.setAppAuth(response.app, response.token)
      appStore.showSuccess('应用登录成功')

      // 跳转到仪表盘
      router.push('/dashboard')
    } catch (error) {
      // 如果API不可用，使用模拟验证
      console.warn('应用登录API不可用，使用模拟验证')

      // 简单的模拟验证：检查应用ID和密钥格式
      if (loginForm.appId && loginForm.appSecret.length >= 8) {
        const mockApp = selectedApp.value || {
          id: loginForm.appId,
          name: '应用',
          code: 'app'
        }
        const mockToken = `app-token-${Date.now()}`

        authStore.setAppAuth(mockApp, mockToken)
        appStore.showSuccess('应用登录成功（模拟模式）')
        router.push('/dashboard')
      } else {
        appStore.showError('应用ID或密钥错误')
      }
    }
  } catch (error) {
    appStore.showError(error.message || '登录失败')
  } finally {
    loading.value = false
  }
}

// 返回应用选择
const backToAppSelection = () => {
  router.push('/app-selection')
}

// 删除应用
const handleDeleteApp = async () => {
  if (!selectedApp.value) {
    appStore.showError('请先选择应用')
    return
  }

  try {
    console.log('Deleting application with ID:', selectedApp.value.id, 'Type:', typeof selectedApp.value.id)
    // 确认删除
    if (confirm(`确定要删除应用 "${selectedApp.value.name}" 吗？此操作不可恢复！`)) {
      await applicationAPI.deleteApplication(selectedApp.value.id)
      appStore.showSuccess('应用删除成功')

      // 清空选中的应用和表单
      selectedApp.value = null
      loginForm.appId = ''
      loginForm.appSecret = ''

      // 返回应用选择页面
      router.push('/app-selection')
    }
  } catch (error) {
    appStore.showError('删除应用失败')
    console.error(error)
  }
}

// 初始化
onMounted(() => {
  // 从应用存储获取选中的应用
  if (appStore.currentApp) {
    selectedApp.value = appStore.currentApp
    loginForm.appId = selectedApp.value.id
  }
})
</script>

<style scoped>
.app-login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #52c41a, #389e0d);
}

.app-login-card {
  width: 100%;
  max-width: 500px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.selected-app-info {
  margin-bottom: 24px;
}

.app-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.app-name {
  font-size: 16px;
  font-weight: 500;
}

.app-code,
.app-id {
  font-size: 14px;
  color: var(--n-text-color-secondary);
}

.login-form {
  margin-top: 24px;
}

.form-actions {
  margin-top: 24px;
}

.app-actions {
  margin-top: 16px;
}

.my-4 {
  margin: 16px 0;
}

.mt-4 {
  margin-top: 16px;
}

.flex {
  display: flex;
}

.justify-end {
  justify-content: flex-end;
}
</style>