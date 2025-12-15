<template>
  <div class="login-container">
    <div class="background-animation"></div>
    <div class="floating-shapes">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
    </div>
    
    <div class="login-content">
      <div class="login-header">
        <div class="logo-container">
          <div class="logo">
            <n-icon size="48" color="#fff">
              <shield-checkmark />
            </n-icon>
          </div>
          <h1 class="system-title">Authos</h1>
          <p class="system-subtitle">统一权限管理系统</p>
        </div>
      </div>
      
      <n-card class="login-card">
        <template #header>
          <div class="card-header">
            <h2>系统管理员登录</h2>
            <n-tag type="info" size="small" round>安全登录</n-tag>
          </div>
        </template>

        <n-form ref="formRef" :model="loginForm" :rules="loginRules" size="large">
          <n-form-item path="username">
            <n-input v-model:value="loginForm.username" placeholder="请输入系统管理员账号" clearable>
              <template #prefix>
                <n-icon color="#666">
                  <person />
                </n-icon>
              </template>
            </n-input>
          </n-form-item>

          <n-form-item path="password">
            <n-input v-model:value="loginForm.password" type="password" placeholder="请输入系统管理员密码"
              show-password-on="click" clearable>
              <template #prefix>
                <n-icon color="#666">
                  <lock-closed />
                </n-icon>
              </template>
            </n-input>
          </n-form-item>

          <div class="form-actions">
            <n-button type="primary" size="large" :loading="loading" @click="handleLogin" block>
              <template #icon>
                <n-icon><log-in /></n-icon>
              </template>
              登录系统
            </n-button>
          </div>
        </n-form>
      </n-card>
      
      <div class="login-footer">
        <p>© 2024 Authos. 保留所有权利.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { Person, LockClosed, ShieldCheckmark, LogIn } from '@vicons/ionicons5'
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
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  overflow: hidden;
}

.background-animation {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.logo {
  margin-bottom: 16px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.system-title {
  font-size: 32px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px 0;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}

.system-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
  font-weight: 300;
}

.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
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
}

.login-footer {
  text-align: center;
  margin-top: 32px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-content {
    padding: 16px;
  }
  
  .system-title {
    font-size: 28px;
  }
  
  .system-subtitle {
    font-size: 14px;
  }
}
</style>