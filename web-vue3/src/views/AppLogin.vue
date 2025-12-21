<template>
  <div class="app-login-container">
    <div class="background-animation"></div>
    <div class="floating-shapes">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
    </div>
    
    <div class="login-content">
      <div class="login-header">
        <n-button text @click="backToAppSelection" class="back-button">
          <template #icon>
            <n-icon size="20"><arrow-back /></n-icon>
          </template>
          返回应用选择
        </n-button>
      </div>
      
      <div v-if="selectedApp" class="selected-app-info">
        <div class="app-card">
          <div class="app-header">
            <div class="app-icon">
              <n-icon size="48" :color="getAppColor()">
                <apps />
              </n-icon>
            </div>
            <div class="app-details">
              <h2>{{ selectedApp.name }}</h2>
              <p>{{ selectedApp.code }}</p>
            </div>
          </div>
          <div class="app-status">
            <n-tag type="info" size="large" :bordered="false">
              <template #icon>
                <n-icon><checkmark-circle /></n-icon>
              </template>
              已选择应用
            </n-tag>
          </div>
        </div>
      </div>

      <n-card class="login-card">
        <template #header>
          <div class="card-header">
            <h2>应用登录</h2>
            <n-tag type="success" size="small" round>安全验证</n-tag>
          </div>
        </template>

        <n-form ref="formRef" :model="loginForm" :rules="loginRules" size="large">
          <n-form-item path="appUuid">
            <n-input v-model:value="loginForm.appUuid" placeholder="请输入应用标识符" clearable>
              <template #prefix>
                <n-icon color="#666">
                  <key />
                </n-icon>
              </template>
            </n-input>
          </n-form-item>

          <n-form-item path="appSecret">
            <n-input v-model:value="loginForm.appSecret" type="password" placeholder="请输入应用密钥"
              show-password-on="click" clearable>
              <template #prefix>
                <n-icon color="#666">
                  <lock-closed />
                </n-icon>
              </template>
            </n-input>
          </n-form-item>

          <div class="form-actions">
            <n-button type="primary" size="large" :loading="loading" @click="handleLogin" class="login-btn">
              <template #icon>
                <n-icon><log-in /></n-icon>
              </template>
              登录应用
            </n-button>
            <n-button type="error" size="large" @click="handleDeleteApp" :disabled="!selectedApp" class="delete-btn">
              <template #icon>
                <n-icon>
                  <trash-outline />
                </n-icon>
              </template>
              删除应用
            </n-button>
          </div>
        </n-form>
      </n-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  NCard, NButton, NIcon, NForm, NFormItem, NInput, NDivider, NTag
} from 'naive-ui'
import {
  ArrowBack,
  Key,
  LockClosed,
  TrashOutline,
  CheckmarkCircle,
  LogIn,
  Apps
} from '@vicons/ionicons5'
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
  appUuid: '',
  appSecret: ''
})

// 获取应用颜色
const getAppColor = () => {
  const colors = ['#667eea', '#f56565', '#48bb78', '#ed8936', '#9f7aea', '#38b2ac', '#ed64a6']
  const index = selectedApp.value?.id ? selectedApp.value.id.toString().charCodeAt(0) % colors.length : 0
  return colors[index]
}

// 登录表单验证规则
const loginRules = {
  appUuid: [
    {
      required: true,
      validator: (rule, value) => {
        if (value === null || value === undefined || value === '') {
          return new Error('请输入应用标识符')
        }
        if (typeof value === 'string') {
          return value.trim() !== '' ? true : new Error('请输入应用标识符')
        }
        return true
      },
      trigger: 'blur'
    }
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
      // 准备登录参数，使用UUID
      const loginParams = {
        appUuid: loginForm.appUuid.trim(),
        appSecret: loginForm.appSecret.trim()
      }

      console.log('发送应用登录请求:', loginParams)

      // 调用应用登录API
      const response = await authAPI.appLogin(loginParams)

      console.log('应用登录响应:', response)

      // 设置应用认证
      authStore.setAppAuth(response.app, response.token)
      appStore.showSuccess('应用登录成功')

      // 跳转到仪表盘
      router.push('/dashboard')
    } catch (error) {
      console.error('应用登录API错误:', error)

      // 显示详细的错误信息
      if (error.response) {
        const status = error.response.status
        const message = error.response.data?.message || error.message
        console.error(`API错误状态: ${status}, 消息: ${message}`)
        appStore.showError(`登录失败: ${message}`)
      } else if (error.request) {
        appStore.showError('网络错误，请检查连接')
      } else {
        appStore.showError('登录失败，请重试')
      }

      // 如果API不可用，使用模拟验证作为备选
      if (error.code === 'NETWORK_ERROR' || error.message?.includes('Network Error')) {
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
          appStore.showError('应用ID或密钥格式不正确')
        }
      }
    }
  } catch (error) {
    console.error('表单验证错误:', error)
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
  console.log('AppLogin 初始化 - appStore.currentApp:', appStore.currentApp)
  console.log('AppLogin 初始化 - localStorage.currentApp:', localStorage.getItem('currentApp'))

  // 优先从localStorage获取应用信息（作为备用方案）
  const storedApp = localStorage.getItem('currentApp')
  if (storedApp) {
    try {
      const appData = JSON.parse(storedApp)
      selectedApp.value = appData
      console.log('从 localStorage 加载应用信息:', appData)
    } catch (e) {
      console.error('解析 localStorage 应用信息失败:', e)
    }
  }

  // 如果应用存储中有当前应用，也使用它
  if (appStore.currentApp) {
    selectedApp.value = appStore.currentApp
    console.log('从 appStore 加载应用信息:', appStore.currentApp)
  }

  // 自动填入应用UUID，但不直接显示在UI上
  if (selectedApp.value) {
    // 确保应用UUID是字符串类型，但不直接显示在UI上
    const appUuid = selectedApp.value.uuid || selectedApp.value.UUID || selectedApp.value.id || ''
    loginForm.appUuid = typeof appUuid === 'number' ? appUuid.toString() : appUuid
    console.log('设置应用UUID:', loginForm.appUuid, '类型:', typeof loginForm.appUuid)
  }
})
</script>

<style scoped>
.app-login-container {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  overflow: hidden;
  user-select: none;
}

.background-animation {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #52c41a 0%, #389e0d 100%);
  z-index: -2;
}

.floating-shapes {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: -1;
  overflow: hidden;
}

.shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  background: #fff;
}

.shape-1 {
  width: 300px;
  height: 300px;
  top: -150px;
  right: -100px;
  animation: float 15s infinite ease-in-out;
}

.shape-2 {
  width: 200px;
  height: 200px;
  bottom: -100px;
  left: -50px;
  animation: float 20s infinite ease-in-out reverse;
}

.shape-3 {
  width: 150px;
  height: 150px;
  top: 50%;
  left: 10%;
  animation: float 25s infinite ease-in-out;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(10deg);
  }
}

.login-content {
  width: 100%;
  max-width: 450px;
  padding: 20px;
  z-index: 1;
  user-select: text;
}

.login-header {
  display: flex;
  justify-content: flex-start;
  margin-bottom: 24px;
}

.back-button {
  color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 8px 16px;
  transition: all 0.3s ease;
}

.back-button:hover {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}

.selected-app-info {
  margin-bottom: 32px;
}

.app-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.app-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.app-icon {
  margin-right: 16px;
  padding: 12px;
  background: rgba(82, 196, 26, 0.1);
  border-radius: 12px;
}

.app-details h2 {
  margin: 0 0 4px 0;
  font-size: 20px;
  font-weight: 600;
  color: #1a202c;
}

.app-details p {
  margin: 0;
  font-size: 14px;
  color: #718096;
}

.app-status {
  display: flex;
  justify-content: center;
}

.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  margin-bottom: 24px;
}

.login-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

.form-actions {
  margin-top: 32px;
  display: flex;
  gap: 12px;
}

.login-btn {
  flex: 2;
}

.delete-btn {
  flex: 1;
  background: rgba(245, 101, 101, 0.2);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(245, 101, 101, 0.5);
  color: #f56565;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.delete-btn:hover {
  background: #f56565;
  color: #fff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(245, 101, 101, 0.4);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-content {
    padding: 16px;
  }
  
  .app-header {
    flex-direction: column;
    text-align: center;
  }
  
  .app-icon {
    margin-right: 0;
    margin-bottom: 16px;
  }
}
</style>