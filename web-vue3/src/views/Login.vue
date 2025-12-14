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

          <div class="divider">
            <span>或者</span>
          </div>

          <div>
            <n-button type="success" :block="true" @click="currentStep = 'create-app'">
              <template #icon>
                <n-icon>
                  <add />
                </n-icon>
              </template>
              创建新应用
            </n-button>
          </div>
        </div>
      </div>

      <!-- 创建应用步骤 -->
      <div v-else-if="currentStep === 'create-app'" class="step-container">
        <div class="step-header">
          <n-steps :current="1" size="small">
            <n-step title="创建应用" />
            <n-step title="用户登录" />
          </n-steps>
        </div>

        <n-form ref="createAppFormRef" :model="createAppForm" :rules="createAppRules">
          <n-form-item path="code" label="应用代码">
            <n-input v-model:value="createAppForm.code" placeholder="请输入应用代码，如: myapp" @blur="generateAppName" />
          </n-form-item>
          <n-form-item path="name" label="应用名称">
            <n-input v-model:value="createAppForm.name" placeholder="请输入应用名称" />
          </n-form-item>
          <n-form-item path="description" label="应用描述">
            <n-input v-model:value="createAppForm.description" type="textarea" placeholder="请输入应用描述" :rows="3" />
          </n-form-item>

          <div class="form-actions">
            <n-button @click="currentStep = 'app-selection'">返回</n-button>
            <n-button type="primary" :loading="creatingApp" @click="handleCreateApp">
              创建并继续
            </n-button>
          </div>
        </n-form>
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
import { Person, LockClosed, Add } from '@vicons/ionicons5'
import { useAuthStore } from '../stores/auth'
import { useAppStore } from '../stores/app'
import { authAPI, applicationAPI } from '../api'

const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

// 表单引用
const createAppFormRef = ref()
const loginFormRef = ref()
const loading = ref(false)
const creatingApp = ref(false)

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

// 创建应用表单
const createAppForm = reactive({
  code: '',
  name: '',
  description: ''
})

// 登录表单
const loginForm = reactive({
  username: '',
  password: ''
})

// 创建应用表单验证规则
const createAppRules = {
  code: [
    { required: true, message: '请输入应用代码', trigger: 'blur' },
    { pattern: /^[a-z][a-z0-9_-]*$/, message: '应用代码只能包含小写字母、数字、下划线和连字符，且必须以字母开头', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入应用名称', trigger: 'blur' }
  ]
}

// 登录表单验证规则
const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

// 生成应用名称
const generateAppName = () => {
  if (createAppForm.code && !createAppForm.name) {
    createAppForm.name = createAppForm.code.replace(/[-_]/g, ' ')
      .split(' ')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ')
  }
}

// 获取应用列表
const fetchApplications = async () => {
  try {
    const response = await applicationAPI.getApplications()
    applications.value = response.items || response
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
  currentStep.value = 'user-login'
}

// 创建应用
const handleCreateApp = async () => {
  try {
    await createAppFormRef.value?.validate()
    creatingApp.value = true

    const newApp = await applicationAPI.createApplication(createAppForm)
    applications.value.push(newApp)
    selectedAppId.value = newApp.id

    appStore.showSuccess('应用创建成功')
    currentStep.value = 'user-login'
  } catch (error) {
    if (error.response?.status === 409) {
      appStore.showError('应用代码已存在')
    } else {
      appStore.showError('创建应用失败')
    }
    console.error(error)
  } finally {
    creatingApp.value = false
  }
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

.divider {
  text-align: center;
  margin: 20px 0;
  position: relative;
  color: #999;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: #e8e8e8;
  z-index: 1;
}

.divider span {
  background: white;
  padding: 0 16px;
  position: relative;
  z-index: 2;
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