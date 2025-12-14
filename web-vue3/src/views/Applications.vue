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
    <n-modal v-model:show="showCreateModal">
      <n-card style="width: 500px" title="创建应用" :bordered="false" size="huge" role="dialog" aria-modal="true">
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

    <!-- 编辑应用模态框 -->
    <n-modal v-model:show="showEditModal">
      <n-card style="width: 500px" title="编辑应用" :bordered="false" size="huge" role="dialog" aria-modal="true">
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

        <template #footer>
          <div class="flex justify-end gap-2">
            <n-button @click="showEditModal = false">取消</n-button>
            <n-button type="primary" :loading="updating" @click="handleUpdate">
              更新
            </n-button>
          </div>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, h } from 'vue'
import { useMessage, NButton, NIcon, NTag } from 'naive-ui'
import { Add, PencilOutline, TrashOutline, Play } from '@vicons/ionicons5'
import { applicationAPI } from '../api'

const message = useMessage()

// 响应式数据
const applications = ref([])
const loading = ref(false)
const creating = ref(false)
const updating = ref(false)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const currentApp = ref(null)

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
    }
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
    applications.value = response.items || response
    pagination.itemCount = response.total || applications.value.length
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
    await formRef.value?.validate()
    creating.value = true

    await applicationAPI.createApplication(createForm)
    message.success('应用创建成功')
    showCreateModal.value = false
    Object.assign(createForm, { code: '', name: '', description: '' })
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
  localStorage.setItem('appId', row.id)
  localStorage.setItem('appCode', row.code)
  localStorage.setItem('appInfo', JSON.stringify(row))
  message.success(`已选择应用: ${row.name}`)
}

// 分页更新
const handlePageUpdate = (page) => {
  pagination.page = page
  fetchApplications()
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