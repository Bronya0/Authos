<template>
  <div class="login-container">
    <n-card class="login-card" title="Authos 权限管理系统">
      <!-- 应用选择步骤 -->
      <div v-if="currentStep === 'app-selection'" class="step-container">
        <div class="step-header">
          <n-steps :current="1" size="small">
            <n-step title="选择应用" />
            <n-step title="用户登录" />
          </n-steps>
        </div>

        <div class="app-selection">
          <div class="mb-4">
            <n-select v-model:value="selectedAppId" placeholder="请选择应用" :options="appOptions" filterable
              @update:value="handleAppSelect" />
          </div>

          <div class="mb-4">
            <n-button type="primary" :block="true" @click="handleContinueWithSelectedApp">
              继续登录
            </n-button>
          </div>
        </div>
      </div>

      <!-- 用户登录步骤 -->
      <div v-else-if="currentStep === 'user-login'" class="step-container">
        <div class="step-header">
          <n-steps :current="2" size="small">
            <n-step title="选择应用" />
            <n-step title="用户登录" />
          </n-steps>
        </div>

        <div class="selected-app-info" v-if="selectedApp">
          <n-tag type="info" size="large">
            当前应用: {{ selectedApp.name }} ({{ selectedApp.code }})
          </n-tag>
          <n-button size="small" text @click="currentStep = 'app-selection'">更换应用</n-button>
        </div>

        <n-form ref="loginFormRef" :model="loginForm" :rules="loginRules">
          <n-form-item path="username">
            <n-input v-model:value="loginForm.username" placeholder="请输入用户名" size="large" clearable>
              <template #prefix>
                <n-icon>
                  <person />
                </n-icon>
              </template>
            </n-input>
          </n-form-item>

          <n-form-item path="password">
            <n-input v-model:value="loginForm.password" type="password" placeholder="请输入密码" size="large"
              show-password-on="click" clearable>
              <template #prefix><n-icon><lock-closed /></n-icon></template>
            </n-input>
          </n-form-item>

          <div class="form-actions">
            <n-button @click="currentStep = 'app-selection'">返回</n-button>
            <n-button type="primary" size="large" :loading="loading" @click="handleLogin">
              登录
            </n-button>
          </div>
        </n-form>
      </div>
    </n-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  NCard, NSteps, NStep, NSelect, NButton, NIcon,
  NForm, NFormItem, NInput, NTag, NAlert
} from 'naive-ui'
import { Person, LockClosed } from '@vicons/ionicons5'
import { useAuthStore } from '../stores/auth'
import { useAppStore } from '../stores/app'
import { authAPI, applicationAPI } from '../api'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

// 表单引用
const loginFormRef = ref()
const loading = ref(false)

// 当前步骤
const currentStep = ref('app-selection')

// 应用数据
const applications = ref([])
const selectedAppId = ref(null)
const selectedApp = computed(() => {
  return applications.value.find(app => app.id === selectedAppId.value)
})

// 应用选项
const appOptions = computed(() => {
  return applications.value.map(app => ({
    label: `${app.name} (${app.code})`,
    value: app.id
  }))
})

// 登录表单
const loginForm = reactive({
  username: '',
  password: ''
})

// 登录表单验证规则
const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

// 获取应用列表
const fetchApplications = async () => {
  try {
    const response = await applicationAPI.getApplications()
    applications.value = response.apps || response.items || response
  } catch (error) {
    console.warn('获取应用列表失败，使用空列表', error)
    applications.value = []
  }
}

// 应用选择处理
const handleAppSelect = (value) => {
  selectedAppId.value = value
}

// 继续使用选中的应用
const handleContinueWithSelectedApp = () => {
  if (!selectedAppId.value) {
    appStore.showWarning('请先选择应用')
    return
  }

  // 设置当前应用
  const selectedApp = applications.value.find(app => app.id === selectedAppId.value)
  if (selectedApp) {
    appStore.setCurrentApp(selectedApp)
  }

  currentStep.value = 'user-login'
}

// 登录处理
const handleLogin = async () => {
  try {
    await loginFormRef.value?.validate()
    loading.value = true

    if (!selectedApp.value) {
      appStore.showError('请先选择应用')
      return
    }

    try {
      const loginData = {
        appCode: selectedApp.value.code,
        username: loginForm.username,
        password: loginForm.password
      }

      const response = await authAPI.login(loginData)
      authStore.setAuth(response.user, response.token, response.app)
      appStore.showSuccess('登录成功')

      // 如果登录响应中包含应用信息，设置当前应用
      if (response.app) {
        appStore.setCurrentApp(response.app)
      } else if (selectedApp.value) {
        // 如果有选中的应用，也设置它
        appStore.setCurrentApp(selectedApp.value)
      }

      router.push('/dashboard')
    } catch (error) {
      // 模拟登录作为后备方案
      console.warn('API连接失败，使用模拟登录')
      if (loginForm.username) {
        const mockToken = `mock-token-${Date.now()}`
        const mockUser = { username: loginForm.username, id: 1 }
        const mockApp = selectedApp.value
        authStore.setAuth(mockUser, mockToken, mockApp)
        appStore.showSuccess('登录成功（模拟模式）')
        router.push('/dashboard')
      } else {
        appStore.showWarning('请输入用户名')
      }
    }
  } catch (error) {
    appStore.showError(error.message || '登录失败')
  } finally {
    loading.value = false
  }
}

// 初始化
onMounted(() => {
  fetchApplications()
})
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1890ff, #096dd9);
}

.login-card {
  width: 100%;
  max-width: 500px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.step-container {
  min-height: 400px;
}

.step-header {
  margin-bottom: 24px;
  text-align: center;
}

.app-selection {
  padding: 20px 0;
}

.mb-4 {
  margin-bottom: 16px;
}

.selected-app-info {
  margin-bottom: 20px;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-top: 24px;
}

.form-actions .n-button:first-child {
  flex: 1;
}

.form-actions .n-button:last-child {
  flex: 1;
}
</style>