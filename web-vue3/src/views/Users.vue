<template>
  <div>
    <n-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span>用户列表</span>
          <n-button type="primary" @click="showAddModal = true">
            <template #icon>
              <n-icon>
                <add />
              </n-icon>
            </template>
            添加用户
          </n-button>
        </div>
      </template>

      <n-data-table :columns="columns" :data="users" :loading="loading" :pagination="pagination" />
    </n-card>

    <!-- 添加/编辑用户模态框 -->
    <n-modal v-model:show="showModal" :title="isEdit ? '编辑用户' : '添加用户'" preset="dialog" :show-icon="false"
      @after-leave="resetForm">
      <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="80px">
        <n-form-item label="用户名" path="username">
          <n-input v-model:value="form.username" :disabled="isEdit" placeholder="请输入用户名" />
        </n-form-item>

        <n-form-item :label="isEdit ? '新密码' : '密码'" path="password">
          <n-input v-model:value="form.password" type="password" :placeholder="isEdit ? '留空不修改' : '请输入密码'" />
        </n-form-item>

        <n-form-item label="状态" path="status">
          <n-radio-group v-model:value="form.status">
            <n-radio :value="1">启用</n-radio>
            <n-radio :value="0">禁用</n-radio>
          </n-radio-group>
        </n-form-item>

        <n-form-item label="角色" path="roleIds">
          <n-select v-model:value="form.roleIds" multiple placeholder="请选择角色" :options="roleOptions" />
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
import { userAPI, roleAPI } from '../api'
import { useAppStore } from '../stores/app'
import { useAuthStore } from '../stores/auth'
import { formatDate, formatStatus, getStatusType } from '../utils/format'
import { People, Add, Create, Trash } from '@vicons/ionicons5'
import { NIcon, NButton, NSpace } from 'naive-ui'

const appStore = useAppStore()
const authStore = useAuthStore()

const users = ref([])
const roles = ref([])
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const isEdit = ref(false)
const currentUserId = ref(null)

const formRef = ref()
const form = reactive({
  username: '',
  password: '',
  status: 1,
  roleIds: []
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: !isEdit.value, message: '请输入密码', trigger: 'blur' }
  ]
}

const roleOptions = computed(() =>
  roles.value.map(role => ({
    label: role.name,
    value: role.id || role.ID
  }))
)

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
    title: '用户名',
    key: 'username'
  },
  {
    title: '状态',
    key: 'status',
    render: (row) => h(
      'n-tag',
      { type: getStatusType(row.status) },
      { default: () => formatStatus(row.status) }
    )
  },
  {
    title: '角色',
    key: 'roles',
    render: (row) => {
      const roles = row.roles || []
      return h(
        'n-space',
        null,
        {
          default: () => roles.map((role, index) =>
            h(
              'n-tag',
              {
                type: 'info',
                size: 'small',
                round: true
              },
              { default: () => role.name }
            )
          )
        }
      )
    }
  },
  {
    title: '创建时间',
    key: 'createdAt',
    render: (row) => formatDate(row.CreatedAt || row.createdAt)
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
                if (confirm('确定要删除该用户吗？')) {
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

const loadUsers = async () => {
  loading.value = true
  try {
    const data = await userAPI.getUsers()
    users.value = data
  } catch (error) {
    appStore.showError('加载用户列表失败')
  } finally {
    loading.value = false
  }
}

const loadRoles = async () => {
  try {
    const data = await roleAPI.getRoles()
    roles.value = data
  } catch (error) {
    appStore.showError('加载角色列表失败')
  }
}

const handleEdit = (row) => {
  isEdit.value = true
  currentUserId.value = row.ID || row.id

  // 填充表单数据
  form.username = row.username
  form.status = row.status
  form.password = ''

  // 设置角色
  const userRoles = (row.roles || []).filter(r => r).map(r => r.ID || r.id)
  form.roleIds = userRoles

  showModal.value = true
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
    saving.value = true

    const data = {
      username: form.username,
      status: form.status,
      roleIds: form.roleIds
    }

    if (form.password) {
      data.password = form.password
    }

    if (isEdit.value) {
      await userAPI.updateUser(currentUserId.value, data)
      appStore.showSuccess('用户更新成功')
    } else {
      await userAPI.createUser(data)
      appStore.showSuccess('用户创建成功')
    }

    showModal.value = false
    loadUsers()
  } catch (error) {
    appStore.showError(error.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const handleDelete = async (row) => {
  try {
    const id = row.ID || row.id
    await userAPI.deleteUser(id)
    appStore.showSuccess('用户删除成功')
    loadUsers()
  } catch (error) {
    appStore.showError('删除失败')
  }
}

const resetForm = () => {
  form.username = ''
  form.password = ''
  form.status = 1
  form.roleIds = []
  isEdit.value = false
  currentUserId.value = null
}

onMounted(() => {
  loadUsers()
  loadRoles()
})
</script>

<style scoped>
/* 确保输入框在亮色和暗色主题下都有良好的对比度 */
:deep(.n-input__input),
:deep(.n-input-group__input) {
  background-color: var(--n-color);
  color: var(--n-text-color);
}

:deep(.n-input__input::placeholder) {
  color: var(--n-placeholder-color);
}

:deep(.n-select) {
  background-color: var(--n-color);
}
</style>