<template>
  <div class="applications-container">
    <n-card title="应用管理" class="mb-4">
      <div class="mb-4">
        <n-button type="primary" @click="showCreateModal = true">
          <template #icon>
            <n-icon>
              <add />
            </n-icon>
          </template>
          创建应用
        </n-button>
      </div>

      <n-data-table :columns="columns" :data="applications" :loading="loading" :pagination="pagination"
        @update:page="handlePageUpdate" />
    </n-card>

    <!-- 创建应用模态框 -->
    <n-modal v-model:show="showCreateModal" preset="dialog" title="创建应用">
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

      <template #action>
        <div class="flex justify-end gap-2">
          <n-button @click="showCreateModal = false">取消</n-button>
          <n-button type="primary" :loading="creating" @click="handleCreate">
            创建
          </n-button>
        </div>
      </template>
    </n-modal>

    <!-- 编辑应用模态框 -->
    <n-modal v-model:show="showEditModal" preset="dialog" title="编辑应用">
      <n-form ref="editFormRef" :model="editForm" :rules="rules">
        <n-form-item path="code" label="应用代码">
          <n-input v-model:value="editForm.code" disabled />
        </n-form-item>
        <n-form-item path="name" label="应用名称">
          <n-input v-model:value="editForm.name" />
        </n-form-item>
        <n-form-item path="description" label="应用描述">
          <n-input v-model:value="editForm.description" type="textarea" :rows="3" />
        </n-form-item>
      </n-form>

      <template #action>
        <div class="flex justify-end gap-2">
          <n-button @click="showEditModal = false">取消</n-button>
          <n-button type="primary" :loading="updating" @click="handleUpdate">
            更新
          </n-button>
        </div>
      </template>
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

      <template #action>
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
import { ref, reactive, onMounted, h } from 'vue'
import { useMessage, NButton, NIcon, NTag, NForm, NFormItem, NInput, NModal, NCard, NAlert, NText } from 'naive-ui'
import { Add, PencilOutline, TrashOutline, Play, CopyOutline } from '@vicons/ionicons5'
import { applicationAPI } from '../api'

const message = useMessage()

// 响应式数据
const applications = ref([])
const loading = ref(false)
const creating = ref(false)
const updating = ref(false)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showSuccessModal = ref(false)
const currentApp = ref(null)
const createdApp = ref({})

// 表单引用
const formRef = ref()
const editFormRef = ref()

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50]
})

// 创建表单
const createForm = reactive({
  code: '',
  name: '',
  description: ''
})

// 编辑表单
const editForm = reactive({
  id: null,
  code: '',
  name: '',
  description: ''
})

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

// 表格列配置
const columns = [
  {
    title: '应用代码',
    key: 'code',
    width: 150
  },
  {
    title: '应用名称',
    key: 'name',
    width: 150
  },
  {
    title: '描述',
    key: 'description',
    ellipsis: {
      tooltip: true
    },
    width: 250
  },
  {
    title: '创建时间',
    key: 'createdAt',
    width: 180,
    render: (row) => {
      return new Date(row.createdAt).toLocaleString('zh-CN')
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render: (row) => {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            ghost: true,
            onClick: () => handleUseApp(row)
          },
          {
            default: () => '使用',
            icon: () => h(NIcon, null, { default: () => h(Play) })
          }
        ),
        h(
          NButton,
          {
            size: 'small',
            style: 'margin-left: 8px',
            onClick: () => handleEdit(row)
          },
          {
            default: () => '编辑',
            icon: () => h(NIcon, null, { default: () => h(PencilOutline) })
          }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            style: 'margin-left: 8px',
            onClick: () => handleDelete(row)
          },
          {
            default: () => '删除',
            icon: () => h(NIcon, null, { default: () => h(TrashOutline) })
          }
        )
      ]
    }
  }
]

// 生成应用名称
const generateName = () => {
  if (createForm.code && !createForm.name) {
    createForm.name = createForm.code.replace(/[-_]/g, ' ')
      .split(' ')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ')
  }
}

// 获取应用列表
const fetchApplications = async () => {
  try {
    loading.value = true
    const response = await applicationAPI.getApplications()
    applications.value = response.apps || response.items || response
    pagination.itemCount = response.total || applications.value.length

    // 调试输出
    console.log('Fetched applications:', applications.value)
    if (applications.value.length > 0) {
      console.log('First application ID:', applications.value[0].id, 'Type:', typeof applications.value[0].id)
    }
  } catch (error) {
    message.error('获取应用列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 创建应用
const handleCreate = async () => {
  try {
    console.log('Starting application creation with form:', createForm)
    await formRef.value?.validate()
    creating.value = true

    console.log('Sending create application request...')
    const response = await applicationAPI.createApplication(createForm)
    console.log('Create application response:', response)
    message.success('应用创建成功')
    console.log('API Response:', response);

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
    console.log('Created app info:', appInfo)

    // 关闭创建模态框，显示成功模态框
    showCreateModal.value = false
    showSuccessModal.value = true

    // 重置表单
    Object.assign(createForm, { code: '', name: '', description: '' })

    // 刷新应用列表
    fetchApplications()
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

// 编辑应用
const handleEdit = (row) => {
  currentApp.value = row
  Object.assign(editForm, {
    id: row.id,
    code: row.code,
    name: row.name,
    description: row.description || ''
  })
  showEditModal.value = true
}

// 更新应用
const handleUpdate = async () => {
  try {
    await editFormRef.value?.validate()
    updating.value = true

    await applicationAPI.updateApplication(editForm.id, {
      name: editForm.name,
      description: editForm.description
    })
    message.success('应用更新成功')
    showEditModal.value = false
    fetchApplications()
  } catch (error) {
    message.error('更新应用失败')
    console.error(error)
  } finally {
    updating.value = false
  }
}

// 删除应用
const handleDelete = async (row) => {
  try {
    console.log('Deleting application with ID:', row.id, 'Type:', typeof row.id)
    await applicationAPI.deleteApplication(row.id)
    message.success('应用删除成功')
    fetchApplications()
  } catch (error) {
    message.error('删除应用失败')
    console.error(error)
  }
}

// 使用应用
const handleUseApp = (row) => {
  message.success(`已选择应用: ${row.name}`)
}

// 分页更新
const handlePageUpdate = (page) => {
  pagination.page = page
  fetchApplications()
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
  fetchApplications()
})
</script>

<style scoped>
.applications-container {
  padding: 20px;
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