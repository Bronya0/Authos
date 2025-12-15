<template>
  <div class="app-selection-container">
    <n-card title="应用管理" class="app-selection-card">
      <template #header-extra>
        <n-button text @click="handleLogout">
          <template #icon>
            <n-icon><log-out /></n-icon>
          </template>
          退出系统
        </n-button>
      </template>

      <div class="app-selection-content">
        <div class="current-app-section" v-if="authStore.currentApp && authStore.isAppLoggedIn">
          <n-alert type="success" class="mb-4">
            <template #header>
              当前登录应用
            </template>
            <div class="current-app-info">
              <div class="app-name">{{ authStore.currentApp.name }}</div>
              <div class="app-code">{{ authStore.currentApp.code }}</div>
              <div class="app-description" v-if="authStore.currentApp.description">
                {{ authStore.currentApp.description }}
              </div>
            </div>
            <template #footer>
              <n-button type="primary" @click="handleContinueToDashboard">
                进入管理界面
              </n-button>
            </template>
          </n-alert>
        </div>

        <div class="app-list-section">
          <div class="section-header">
            <h3>选择应用</h3>
            <n-button type="success" @click="showCreateModal = true" :disabled="!canCreateApp">
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
              <div class="app-card-header">
                <div class="app-icon">
                  <n-icon size="32">
                    <apps-icon />
                  </n-icon>
                </div>
                <div class="app-info">
                  <div class="app-name">{{ app.name }}</div>
                  <div class="app-code">{{ app.code }}</div>
                </div>
              </div>
              <div class="app-description" v-if="app.description">
                {{ app.description }}
              </div>
              <div class="app-meta">
                <n-tag size="small" type="info">
                  创建时间: {{ formatDate(app.createdAt) }}
                </n-tag>
              </div>
              <div class="app-actions">
                <n-button size="small" type="primary" @click.stop="handleLoginToApp(app)">
                  登录应用
                </n-button>
              </div>
            </div>
          </div>

          <div v-else class="empty-state">
            <n-empty description="暂无应用" />
            <n-button type="primary" @click="showCreateModal = true" :disabled="!canCreateApp">
              创建第一个应用
            </n-button>
          </div>
        </div>
      </div>
    </n-card>

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
    <n-modal v-model:show="showSuccessModal" preset="dialog" title="应用创建成功">
      <n-alert type="success" class="mb-4">
        应用已成功创建，请保存以下信息，应用密钥仅在创建时显示一次！
      </n-alert>

      <n-form label-placement="left" label-width="100px">
        <n-form-item label="应用ID">
          <n-input :value="createdApp.id" readonly />
          <template #suffix>
            <n-button text @click="copyToClipboard(createdApp.id)">
              <template #icon>
                <n-icon><copy-outline /></n-icon>
              </template>
            </n-button>
          </template>
        </n-form-item>

        <n-form-item label="应用UUID">
          <n-input :value="createdApp.uuid" readonly />
          <template #suffix>
            <n-button text @click="copyToClipboard(createdApp.uuid)">
              <template #icon>
                <n-icon><copy-outline /></n-icon>
              </template>
            </n-button>
          </template>
        </n-form-item>

        <n-form-item label="应用密钥">
          <n-input :value="createdApp.secretKey" readonly type="password" show-password-on="click" />
          <template #suffix>
            <n-button text @click="copyToClipboard(createdApp.secretKey)">
              <template #icon>
                <n-icon><copy-outline /></n-icon>
              </template>
            </n-button>
          </template>
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="flex justify-end gap-2">
          <n-button type="primary" @click="showSuccessModal = false">
            我已保存
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, NButton, NIcon, NTag, NEmpty, NCard, NModal, NForm, NFormItem, NInput, NAlert } from 'naive-ui'
import { useAppStore } from '../stores/app'
import { useAuthStore } from '../stores/auth'
import { applicationAPI } from '../api'
import { Add as AddIcon, Apps as AppsIcon, LogOut, CopyOutline } from '@vicons/ionicons5'

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
const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    message.success('已复制到剪贴板')
  } catch (error) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    message.success('已复制到剪贴板')
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
  padding: 20px;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.app-selection-card {
  max-width: 1200px;
  width: 100%;
}

.app-selection-content {
  padding: 20px 0;
}

.current-app-section {
  margin-bottom: 32px;
}

.current-app-info {
  padding: 16px 0;
}

.current-app-info .app-name {
  font-size: 18px;
  font-weight: bold;
  color: var(--n-text-color);
  margin-bottom: 4px;
}

.current-app-info .app-code {
  font-size: 14px;
  color: var(--n-text-color-disabled);
  margin-bottom: 8px;
}

.current-app-info .app-description {
  font-size: 14px;
  color: var(--n-text-color-secondary);
  margin-bottom: 16px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.section-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 500;
}

.app-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.app-card {
  border: 1px solid var(--n-border-color);
  border-radius: 8px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  background: var(--n-color);
}

.app-card:hover {
  border-color: var(--n-color-primary);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.app-card.selected {
  border-color: var(--n-color-primary);
  background: rgba(24, 144, 255, 0.05);
}

.app-card-header {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.app-icon {
  margin-right: 12px;
  color: var(--n-color-primary);
}

.app-info .app-name {
  font-size: 16px;
  font-weight: 500;
  color: var(--n-text-color);
  margin-bottom: 2px;
}

.app-info .app-code {
  font-size: 12px;
  color: var(--n-text-color-disabled);
}

.app-description {
  font-size: 14px;
  color: var(--n-text-color-secondary);
  margin-bottom: 12px;
  line-height: 1.5;
}

.app-meta {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 12px;
}

.app-actions {
  display: flex;
  justify-content: center;
}

.empty-state {
  text-align: center;
  padding: 60px 0;
}

.empty-state .n-empty {
  margin-bottom: 24px;
}

.mb-4 {
  margin-bottom: 16px;
}

.flex {
  display: flex;
}

.justify-end {
  justify-content: flex-end;
}

.gap-2 {
  gap: 8px;
}
</style>