<template>
  <div class="app-selection-container">
    <div class="background-pattern"></div>
    
    <div class="app-selection-header">
      <div class="header-content">
        <div class="logo-section">
          <div class="logo">
            <n-icon size="40" color="#667eea">
              <shield-checkmark />
            </n-icon>
          </div>
          <div class="title-section">
            <h1>应用管理中心</h1>
            <p>管理和配置您的所有应用</p>
          </div>
        </div>
        <div class="user-section">
          <n-dropdown trigger="hover" :options="userMenuOptions" @select="handleUserMenuSelect">
            <n-button text>
              <template #icon>
                <n-icon size="20"><person-circle /></n-icon>
              </template>
              系统管理员
              <template #suffix>
                <n-icon size="14"><chevron-down /></n-icon>
              </template>
            </n-button>
          </n-dropdown>
        </div>
      </div>
    </div>

    <div class="app-selection-content">
      <div class="current-app-section" v-if="authStore.currentApp && authStore.isAppLoggedIn">
        <n-card class="current-app-card" size="large">
          <div class="current-app-header">
            <div class="app-icon-large">
              <n-icon size="48" color="#52c41a">
                <checkmark-circle />
              </n-icon>
            </div>
            <div class="current-app-info">
              <h2>当前登录应用</h2>
              <div class="app-name">{{ authStore.currentApp.name }}</div>
              <div class="app-code">{{ authStore.currentApp.code }}</div>
              <div class="app-description" v-if="authStore.currentApp.description">
                {{ authStore.currentApp.description }}
              </div>
            </div>
          </div>
          <div class="current-app-actions">
            <n-button type="primary" size="large" @click="handleContinueToDashboard">
              <template #icon>
                <n-icon><arrow-forward /></n-icon>
              </template>
              进入管理界面
            </n-button>
          </div>
        </n-card>
      </div>

      <div class="app-list-section">
        <div class="section-header">
          <div class="section-title">
            <h2>我的应用</h2>
            <p>选择或创建新的应用</p>
          </div>
          <n-button type="primary" size="large" @click="showCreateModal = true" :disabled="!canCreateApp" class="create-btn">
            <template #icon>
              <n-icon>
                <add-icon />
              </n-icon>
            </template>
            创建新应用
          </n-button>
        </div>

        <div class="app-grid" v-if="applications.length > 0">
          <div v-for="app in applications" :key="app.id" class="app-card"
            :class="{ 'selected': authStore.currentApp && app.id === authStore.currentApp.id }"
            @click="handleSelectApp(app)">
            <div class="app-card-background"></div>
            <div class="app-card-content">
              <div class="app-card-header">
                <div class="app-icon">
                  <n-icon size="32" :color="getAppColor(app)">
                    <apps-icon />
                  </n-icon>
                </div>
                <div class="app-info">
                  <div class="app-name">{{ app.name }}</div>
                  <div class="app-code">{{ app.code }}</div>
                </div>
                <div v-if="authStore.currentApp && app.id === authStore.currentApp.id" class="selected-badge">
                  <n-tag type="success" size="small">当前应用</n-tag>
                </div>
              </div>
              <div class="app-description" v-if="app.description">
                {{ app.description }}
              </div>
              <div class="app-meta">
                <n-tag size="small" type="info" :bordered="false">
                  <template #icon>
                    <n-icon size="12"><calendar /></n-icon>
                  </template>
                  {{ formatDate(app.createdAt) }}
                </n-tag>
              </div>
              <div class="app-actions">
                <n-button size="small" type="primary" @click.stop="handleLoginToApp(app)">
                  <template #icon>
                    <n-icon><log-in /></n-icon>
                  </template>
                  登录应用
                </n-button>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="empty-state">
          <div class="empty-content">
            <div class="empty-icon">
              <n-icon size="80" color="#d9d9d9">
                <apps />
              </n-icon>
            </div>
            <h3>暂无应用</h3>
            <p>创建您的第一个应用开始管理权限</p>
            <n-button type="primary" size="large" @click="showCreateModal = true" :disabled="!canCreateApp">
              <template #icon>
                <n-icon>
                  <add-icon />
                </n-icon>
              </template>
              创建第一个应用
            </n-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建应用模态框 -->
    <n-modal v-model:show="showCreateModal">
      <n-card style="width: 500px" title="创建新应用" :bordered="false" size="huge" role="dialog" aria-modal="true">
        <n-form ref="formRef" :model="createForm" :rules="rules">
          <n-form-item path="code" label="应用代码">
            <n-input v-model:value="createForm.code" placeholder="请输入应用代码，如: myapp" @blur="generateName" />
          </n-form-item>
          <n-form-item path="name" label="应用名称">
            <n-input v-model:value="createForm.name" placeholder="请输入应用名称" />
          </n-form-item>
          <n-form-item path="description" label="应用描述">
            <n-input v-model:value="createForm.description" type="textarea" placeholder="请输入应用描述" :rows="3" />
          </n-form-item>
        </n-form>

        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showCreateModal = false">取消</n-button>
            <n-button type="primary" :loading="creating" @click="handleCreate">
              创建
            </n-button>
          </div>
        </template>
      </n-card>
    </n-modal>

    <!-- 应用创建成功模态框 -->
    <n-modal v-model:show="showSuccessModal" :mask-closable="false" preset="dialog" title="应用创建成功">
      <div class="success-content">
        <div class="success-icon">
          <n-icon size="64" color="#52c41a">
            <checkmark-circle />
          </n-icon>
        </div>
        <h3>应用创建成功！</h3>
        <p class="success-message">请保存以下信息，应用密钥仅在创建时显示一次！</p>

        <div class="credentials-container">
          <n-card class="credential-card" title="应用标识" size="small">
            <div class="credential-item">
              <span class="label">应用UUID:</span>
              <div class="value-container">
                <n-input :value="createdApp.uuid" readonly class="credential-input" />
                <n-button circle type="primary" size="small" @click="copyToClipboard(createdApp.uuid)" class="copy-btn">
                  <template #icon>
                    <n-icon><copy-outline /></n-icon>
                  </template>
                </n-button>
              </div>
            </div>
          </n-card>

          <n-card class="credential-card" title="安全凭证" size="small">
            <div class="credential-item">
              <span class="label">应用密钥:</span>
              <div class="value-container">
                <n-input :value="createdApp.secretKey" readonly type="password" show-password-on="click" class="credential-input" />
                <n-button circle type="primary" size="small" @click="copyToClipboard(createdApp.secretKey)" class="copy-btn">
                  <template #icon>
                    <n-icon><copy-outline /></n-icon>
                  </template>
                </n-button>
              </div>
            </div>
          </n-card>
        </div>
      </div>

      <template #footer>
        <div class="success-footer">
          <n-button type="primary" size="large" @click="showSuccessModal = false" block>
            我已保存，继续管理
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, h } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, NButton, NIcon, NTag, NEmpty, NCard, NModal, NForm, NFormItem, NInput, NAlert, NDropdown } from 'naive-ui'
import { useAppStore } from '../stores/app'
import { useAuthStore } from '../stores/auth'
import { applicationAPI } from '../api'
import {
  Add as AddIcon,
  Apps as AppsIcon,
  LogOut,
  CopyOutline,
  ShieldCheckmark,
  PersonCircle,
  ChevronDown,
  CheckmarkCircle,
  ArrowForward,
  LogIn,
  Calendar,
  Close,
  Apps
} from '@vicons/ionicons5'

const router = useRouter()
const message = useMessage()
const appStore = useAppStore()
const authStore = useAuthStore()

// 响应式数据
const applications = ref([])
const loading = ref(false)
const creating = ref(false)
const showCreateModal = ref(false)
const showSuccessModal = ref(false)
const createdApp = ref({})

// 权限检查
const canCreateApp = computed(() => {
  // 系统管理员已登录可以创建应用
  return authStore.isSystemLoggedIn
})

// 用户菜单选项
const userMenuOptions = [
  {
    label: '退出系统',
    key: 'logout',
    icon: () => h(NIcon, null, { default: () => h(LogOut) })
  }
]

// 处理用户菜单选择
const handleUserMenuSelect = (key) => {
  if (key === 'logout') {
    handleLogout()
  }
}

// 获取应用颜色
const getAppColor = (app) => {
  const colors = ['#667eea', '#f56565', '#48bb78', '#ed8936', '#9f7aea', '#38b2ac', '#ed64a6']
  const index = app.id ? app.id.toString().charCodeAt(0) % colors.length : 0
  return colors[index]
}

// 创建表单
const createForm = reactive({
  code: '',
  name: '',
  description: ''
})

// 表单引用
const formRef = ref()

// 表单验证规则
const rules = {
  code: [
    { required: true, message: '请输入应用代码', trigger: 'blur' },
    { pattern: /^[a-z][a-z0-9_-]*$/, message: '应用代码只能包含小写字母、数字、下划线和连字符，且必须以字母开头', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入应用名称', trigger: 'blur' }
  ]
}

// 生成应用名称
const generateName = () => {
  if (createForm.code && !createForm.name) {
    createForm.name = createForm.code.replace(/[-_]/g, ' ')
      .split(' ')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ')
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString('zh-CN')
}

// 获取应用列表
const fetchApplications = async () => {
  try {
    loading.value = true
    const response = await applicationAPI.getApplications()
    applications.value = response.apps || response.items || response
  } catch (error) {
    message.error('获取应用列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 选择应用
const handleSelectApp = (app) => {
  appStore.setCurrentApp(app)
  message.success(`已选择应用: ${app.name}`)
}

// 登录到应用
const handleLoginToApp = (app) => {
  console.log('选择登录应用:', app)
  appStore.setCurrentApp(app)

  // 存储应用信息到localStorage作为备用
  localStorage.setItem('currentApp', JSON.stringify(app))

  router.push('/app-login')
}

// 继续到仪表盘
const handleContinueToDashboard = () => {
  router.push('/dashboard')
}

// 系统登出
const handleLogout = () => {
  authStore.logout()
  appStore.clearCurrentApp()
  router.push('/system-login')
}

// 创建应用
const handleCreate = async () => {
  try {
    await formRef.value?.validate()
    creating.value = true

    const response = await applicationAPI.createApplication(createForm)

    // 保存创建的应用信息，包括ID和密钥
    const appInfo = {
      id: response.appId || response.app?.id || '',
      uuid: response.appUuid || response.app?.uuid || '',
      secretKey: response.secretKey || '',
      code: createForm.code,
      name: createForm.name,
      description: createForm.description || ''
    }

    createdApp.value = appInfo

    // 关闭创建模态框，显示成功模态框
    showCreateModal.value = false
    showSuccessModal.value = true

    // 清空表单
    Object.assign(createForm, { code: '', name: '', description: '' })

    // 重新获取应用列表
    await fetchApplications()

  } catch (error) {
    if (error.response?.status === 409) {
      message.error('应用代码已存在')
    } else {
      message.error('创建应用失败')
    }
    console.error(error)
  } finally {
    creating.value = false
  }
}

// 复制到剪贴板
const copyToClipboard = (text) => {
  if (!text) {
    message.warning('复制内容为空')
    return
  }

  // 1. 尝试使用 Clipboard API (仅在安全上下文有效: HTTPS 或 localhost)
  if (navigator.clipboard && navigator.clipboard.writeText) {
    navigator.clipboard.writeText(text)
      .then(() => {
        message.success('已复制到剪贴板')
      })
      .catch((err) => {
        console.error('Clipboard API error:', err)
        // Clipboard API 失败时尝试降级方案
        // 注意：如果是用户拒绝权限导致的异步错误，execCommand 可能因失去用户手势而失效
        fallbackCopy(text)
      })
  } else {
    // 2. HTTP 环境或不支持 Clipboard API，直接使用降级方案 (同步执行，确保 execCommand 有效)
    fallbackCopy(text)
  }
}

// 降级复制方案
const fallbackCopy = (text) => {
  try {
    const textArea = document.createElement('textarea')
    textArea.value = text
    
    // 样式调整，避免页面抖动，同时保证可见性（有些浏览器不复制完全隐藏的元素）
    textArea.style.position = 'fixed'
    textArea.style.top = '0'
    textArea.style.left = '0'
    textArea.style.width = '2em'
    textArea.style.height = '2em'
    textArea.style.padding = '0'
    textArea.style.border = 'none'
    textArea.style.outline = 'none'
    textArea.style.boxShadow = 'none'
    textArea.style.background = 'transparent'
    
    document.body.appendChild(textArea)
    
    textArea.select()
    textArea.setSelectionRange(0, 99999) // 兼容移动端
    
    const successful = document.execCommand('copy')
    document.body.removeChild(textArea)
    
    if (successful) {
      message.success('已复制到剪贴板')
    } else {
      message.error('复制失败，请手动复制')
    }
  } catch (err) {
    console.error('Fallback copy error:', err)
    message.error('复制失败，请手动复制')
  }
}

// 初始化
onMounted(() => {
  // 检查系统管理员是否已登录
  if (!authStore.isSystemLoggedIn) {
    router.push('/system-login')
    return
  }

  fetchApplications()
})
</script>

<style scoped>
.app-selection-container {
  min-height: 100vh;
  background: #f5f7fa;
  position: relative;
}

.background-pattern {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image:
    radial-gradient(circle at 20% 30%, rgba(102, 126, 234, 0.05) 0%, transparent 50%),
    radial-gradient(circle at 80% 70%, rgba(118, 75, 162, 0.05) 0%, transparent 50%);
  z-index: 0;
}

.app-selection-header {
  position: relative;
  z-index: 1;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  padding: 0 24px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
  height: 80px;
}

.logo-section {
  display: flex;
  align-items: center;
}

.logo {
  margin-right: 16px;
  padding: 12px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 12px;
}

.title-section h1 {
  margin: 0 0 4px 0;
  font-size: 24px;
  font-weight: 600;
  color: #1a202c;
}

.title-section p {
  margin: 0;
  font-size: 14px;
  color: #718096;
}

.user-section {
  display: flex;
  align-items: center;
}

.app-selection-content {
  position: relative;
  z-index: 1;
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}

.current-app-section {
  margin-bottom: 32px;
}

.current-app-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 10px 25px rgba(102, 126, 234, 0.2);
}

.current-app-header {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
}

.app-icon-large {
  margin-right: 20px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  backdrop-filter: blur(10px);
}

.current-app-info h2 {
  margin: 0 0 12px 0;
  font-size: 20px;
  font-weight: 600;
  color: #fff;
}

.current-app-info .app-name {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
}

.current-app-info .app-code {
  font-size: 16px;
  opacity: 0.8;
  margin-bottom: 8px;
}

.current-app-info .app-description {
  font-size: 14px;
  opacity: 0.7;
}

.current-app-actions {
  display: flex;
  justify-content: flex-end;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 32px;
}

.section-title h2 {
  margin: 0 0 4px 0;
  font-size: 24px;
  font-weight: 600;
  color: #1a202c;
}

.section-title p {
  margin: 0;
  font-size: 14px;
  color: #718096;
}

.create-btn {
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.app-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

.app-card {
  position: relative;
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #fff;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

.app-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.1);
}

.app-card.selected {
  border-color: #667eea;
  box-shadow: 0 8px 16px rgba(102, 126, 234, 0.2);
}

.app-card-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 4px;
  background: linear-gradient(90deg, #667eea, #764ba2);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.app-card:hover .app-card-background,
.app-card.selected .app-card-background {
  opacity: 1;
}

.app-card-content {
  padding: 24px;
}

.app-card-header {
  display: flex;
  align-items: flex-start;
  margin-bottom: 16px;
  position: relative;
}

.app-icon {
  margin-right: 16px;
  padding: 12px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 12px;
}

.app-info {
  flex: 1;
}

.app-info .app-name {
  font-size: 18px;
  font-weight: 600;
  color: #1a202c;
  margin-bottom: 4px;
}

.app-info .app-code {
  font-size: 14px;
  color: #718096;
}

.selected-badge {
  position: absolute;
  top: 0;
  right: 0;
}

.app-description {
  font-size: 14px;
  color: #4a5568;
  margin-bottom: 16px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.app-meta {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 16px;
}

.app-actions {
  display: flex;
  justify-content: center;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.empty-content {
  text-align: center;
  max-width: 400px;
}

.empty-icon {
  margin-bottom: 24px;
}

.empty-content h3 {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
  color: #1a202c;
}

.empty-content p {
  margin: 0 0 24px 0;
  font-size: 14px;
  color: #718096;
}

.success-content {
  text-align: center;
  padding: 16px 0;
}

.success-icon {
  margin-bottom: 16px;
}

.success-content h3 {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
  color: #1a202c;
}

.success-message {
  margin: 0 0 24px 0;
  font-size: 14px;
  color: #718096;
}

.credentials-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.credential-card {
  border-radius: 8px;
}

.credential-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.label {
  font-size: 12px;
  font-weight: 500;
  color: #718096;
}

.value-container {
  display: flex;
  gap: 8px;
}

.credential-input {
  flex: 1;
}

.copy-btn {
  flex-shrink: 0;
}

.success-footer {
  margin-top: 16px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-selection-header {
    padding: 0 16px;
  }
  
  .header-content {
    flex-direction: column;
    height: auto;
    padding: 16px 0;
  }
  
  .logo-section {
    margin-bottom: 16px;
  }
  
  .app-selection-content {
    padding: 24px 16px;
  }
  
  .app-grid {
    grid-template-columns: 1fr;
  }
  
  .current-app-header {
    flex-direction: column;
    text-align: center;
  }
  
  .app-icon-large {
    margin-right: 0;
    margin-bottom: 16px;
  }
  
  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
}
</style>