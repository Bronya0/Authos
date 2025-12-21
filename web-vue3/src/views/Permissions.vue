<template>
  <div>
    <n-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <n-space align="center">
            <span>权限列表</span>
            <n-input v-model:value="searchParams.name" placeholder="搜索权限名称" clearable @update:value="handleSearch" style="width: 180px" />
            <n-input v-model:value="searchParams.path" placeholder="搜索路径" clearable @update:value="handleSearch" style="width: 180px" />
            <n-select v-model:value="searchParams.method" placeholder="方法" clearable :options="methodOptions" @update:value="handleSearch" style="width: 120px" />
            <n-button @click="handleReset">重置</n-button>
          </n-space>
          <n-button type="primary" @click="showAddModal = true">
            <template #icon>
              <n-icon>
                <add />
              </n-icon>
            </template>
            添加权限
          </n-button>
        </div>
      </template>

      <n-data-table :columns="columns" :data="permissions || []" :loading="loading" :pagination="pagination"
        :row-key="row => row.id || row.ID" />
    </n-card>

    <!-- 添加/编辑权限模态框 -->
    <n-modal v-model:show="showModal" :title="isEdit ? '编辑权限' : '添加权限'" preset="dialog" :show-icon="false"
      @after-leave="resetForm">
      <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="80px">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="form.name" placeholder="请输入权限名称" />
        </n-form-item>

        <n-form-item label="路径" path="path">
          <n-input v-model:value="form.path" placeholder="请输入API路径，如/api/v1/users" />
        </n-form-item>

        <n-form-item label="方法" path="method">
          <n-select v-model:value="form.method" placeholder="请选择HTTP方法" :options="methodOptions" />
        </n-form-item>

        <n-form-item label="描述" path="description">
          <n-input v-model:value="form.description" type="textarea" placeholder="请输入权限描述" />
        </n-form-item>
      </n-form>

      <template #action>
        <n-space>
          <n-button @click="showModal = false">取消</n-button>
          <n-button type="primary" :loading="saving" @click="handleSave">
            保存
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, h } from 'vue'
import { apiPermissionAPI } from '../api'
import { useAppStore } from '../stores/app'
import { Key, Add, Create, Trash } from '@vicons/ionicons5'
import { NIcon, NButton, NSpace } from 'naive-ui'

const appStore = useAppStore()

const permissions = ref([])
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const isEdit = ref(false)
const currentPermissionId = ref(null)
const searchParams = reactive({
  name: '',
  path: '',
  method: null
})

const handleSearch = () => {
  loadPermissions()
}

const handleReset = () => {
  searchParams.name = ''
  searchParams.path = ''
  searchParams.method = null
  loadPermissions()
}

const formRef = ref()
const form = reactive({
  name: '',
  path: '',
  method: '',
  description: ''
})

const rules = {
  name: [
    { required: true, message: '请输入权限名称', trigger: 'blur' }
  ],
  path: [
    { required: true, message: '请输入API路径', trigger: 'blur' }
  ],
  method: [
    { required: true, message: '请选择HTTP方法', trigger: 'change' }
  ]
}

const actionOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' }
]

const methodOptions = [
  { label: '全部', value: '*' },
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' }
]

const showAddModal = computed({
  get: () => showModal.value && !isEdit.value,
  set: (val) => {
    if (val) {
      isEdit.value = false
      showModal.value = true
    } else {
      showModal.value = false
    }
  }
})

const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
    render: (row) => row.ID || row.id
  },
  {
    title: '名称',
    key: 'name'
  },
  {
    title: '路径',
    key: 'path'
  },
  {
    title: '方法',
    key: 'method',
    width: 80
  },
  {
    title: '描述',
    key: 'description',
    render: (row) => row.description || '-'
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render: (row) => h(
      NSpace,
      null,
      {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              type: 'warning',
              quaternary: true,
              onClick: () => handleEdit(row)
            },
            {
              default: () => '编辑',
              icon: () => h(NIcon, null, { default: () => h(Create) })
            }
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'error',
              quaternary: true,
              onClick: () => {
                if (confirm('确定要删除该权限吗？')) {
                  handleDelete(row)
                }
              }
            },
            {
              default: () => '删除',
              icon: () => h(NIcon, null, { default: () => h(Trash) })
            }
          )
        ]
      }
    )
  }
]

const pagination = {
  pageSize: 10
}

const loadPermissions = async () => {
  loading.value = true
  try {
    const data = await apiPermissionAPI.getApiPermissions(searchParams)
    // 确保数据是数组格式，处理null和undefined的情况
    if (data && Array.isArray(data)) {
      permissions.value = data
    } else if (data && data.items && Array.isArray(data.items)) {
      permissions.value = data.items
    } else if (data && Array.isArray(data)) {
      permissions.value = data
    } else {
      permissions.value = []
    }
  } catch (error) {
    console.error('加载权限列表失败:', error)
    appStore.showError('加载权限列表失败')
    // 设置空数组防止错误
    permissions.value = []
  } finally {
    loading.value = false
  }
}

const handleEdit = (row) => {
  isEdit.value = true
  currentPermissionId.value = row.ID || row.id

  // 填充表单数据
  form.name = row.name
  form.path = row.path
  form.method = row.method
  form.description = row.description || ''

  showModal.value = true
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
    saving.value = true

    const data = {
      name: form.name,
      path: form.path,
      method: form.method,
      description: form.description
    }

    if (isEdit.value) {
      await apiPermissionAPI.updateApiPermission(currentPermissionId.value, data)
      appStore.showSuccess('权限更新成功')
    } else {
      await apiPermissionAPI.createApiPermission(data)
      appStore.showSuccess('权限创建成功')
    }

    showModal.value = false
    loadPermissions()
  } catch (error) {
    appStore.showError(error.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const handleDelete = async (row) => {
  try {
    const id = row.id || row.ID
    await apiPermissionAPI.deleteApiPermission(id)
    appStore.showSuccess('权限删除成功')
    loadPermissions()
  } catch (error) {
    appStore.showError('删除失败')
  }
}

const resetForm = () => {
  form.name = ''
  form.path = ''
  form.method = ''
  form.description = ''
  isEdit.value = false
  currentPermissionId.value = null
}

onMounted(() => {
  loadPermissions()
})
</script>