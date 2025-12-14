<template>
  <div class="login-container">
    <n-card class="login-card" title="Authos 权限管理系统">
      <template #header-extra>
        <n-tag type="info" size="small">系统管理员登录</n-tag>
      </template>

      <n-form ref="formRef" :model="loginForm" :rules="loginRules">
        <n-form-item path="username" label="管理员账号">
          <n-input v-model:value="loginForm.username" placeholder="请输入系统管理员账号" size="large" clearable>
            <template #prefix>
              <n-icon>
                <person />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item path="password" label="管理员密码">
          <n-input v-model:value="loginForm.password" type="password" placeholder="请输入系统管理员密码" size="large"
            show-password-on="click" clearable>
            <template #prefix><n-icon><lock-closed /></n-icon></template>
          </n-input>
        </n-form-item>

        <div class="form-actions">
          <n-button type="primary" size="large" :loading="loading" @click="handleLogin" block>
            登录系统
          </n-button>
        </div>
      </n-form>

      <n-divider class="my-4">系统说明</n-divider>

      <div class="system-info">
        <n-alert type="info" title="内置管理员账号" :show-icon="false">
          <div class="admin-info">
            <p>默认账号: admin</p>
            <p>默认密码: admin123</p>
          </div>
        </n-alert>

        <n-alert type="warning" title="安全提示" class="mt-4">
          系统管理员登录后，您可以选择或创建应用，然后使用应用 ID 和密钥进行应用级登录。
        </n-alert>
      </div>
    </n-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { Person, LockClosed } from '@vicons/ionicons5'
import { useAuthStore } from '../stores/auth'
import { authAPI } from '../api'

const router = useRouter()
const authStore = useAuthStore()

// 表单引用
const formRef = ref()
const loading = ref(false)

// 登录表单
const loginForm = reactive({
  username: '',
  password: ''
})

// 登录表单验证规则
const loginRules = {
  username: [
    { required: true, message: '请输入管理员账号', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入管理员密码', trigger: 'blur' }
  ]
}

// 系统管理员登录处理
const handleLogin = async () => {
  try {
    await formRef.value?.validate()
    loading.value = true

    try {
      // 调用系统管理员登录API
      const response = await authAPI.systemLogin({
        username: loginForm.username,
        password: loginForm.password
      })

      // 设置系统管理员认证
      authStore.setSystemAuth(response.user, response.token)
      authStore.showSuccess('系统管理员登录成功')

      // 跳转到应用选择页面
      router.push('/app-selection')
    } catch (error) {
      console.error('登录错误:', error)

      if (error.response) {
        // 服务器返回了错误状态码
        const status = error.response.status
        const message = error.response.data?.message || '登录失败'

        if (status === 401) {
          authStore.showError('账号或密码错误')
        } else {
          authStore.showError(`登录失败: ${message}`)
        }
      } else if (error.request) {
        // 请求已发送但没有收到响应
        authStore.showError('网络错误，请检查后端服务是否启动')
      } else {
        // 请求配置错误
        authStore.showError('请求配置错误')
      }
    }
  } catch (error) {
    authStore.showError(error.message || '登录失败')
  } finally {
    loading.value = false
  }
}
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

.form-actions {
  margin-top: 24px;
}

.system-info {
  margin-top: 16px;
}

.admin-info {
  display: flex;
  justify-content: space-between;
}

.admin-info p {
  margin: 0;
  font-family: monospace;
  background: rgba(0, 0, 0, 0.05);
  padding: 4px 8px;
  border-radius: 4px;
}

.my-4 {
  margin: 16px 0;
}

.mt-4 {
  margin-top: 16px;
}
</style>