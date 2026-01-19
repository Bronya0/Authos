<template>
  <div>
    <n-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <n-space align="center">
            <span>配置字典</span>
            <n-input v-model:value="searchParams.key" placeholder="按key搜索" clearable @update:value="handleSearch" style="width: 220px" />
            <n-button @click="handleReset">重置</n-button>
          </n-space>
          <n-button type="primary" @click="showAddModal = true">
            <template #icon>
              <n-icon>
                <add />
              </n-icon>
            </template>
            新增字典
          </n-button>
        </div>
      </template>

      <n-data-table :columns="columns" :data="items || []" :loading="loading" :pagination="pagination"
        :row-key="row => row.id || row.ID" />
    </n-card>

    <n-modal v-model:show="showModal" :title="isEdit ? '编辑字典' : '新增字典'" preset="dialog" :show-icon="false"
      @after-leave="resetForm">
      <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="80px">
        <n-form-item label="key" path="key">
          <n-input v-model:value="form.key" placeholder="请输入字典key" />
        </n-form-item>

        <n-form-item label="value" path="value">
          <n-input v-model:value="form.value" placeholder="请输入字典value" />
        </n-form-item>

        <n-form-item label="描述" path="desc">
          <n-input v-model:value="form.desc" type="textarea" placeholder="请输入描述" />
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
import { configDictionaryAPI } from '../api'
import { useAppStore } from '../stores/app'
import { Add, Create, Trash } from '@vicons/ionicons5'
import { NIcon, NButton, NSpace } from 'naive-ui'
import { formatDate } from '../utils/format'

const appStore = useAppStore()

const items = ref([])
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const searchParams = reactive({
  key: ''
})

const handleSearch = () => {
  loadItems()
}

const handleReset = () => {
  searchParams.key = ''
  loadItems()
}

const formRef = ref()
const form = reactive({
  key: '',
  value: '',
  desc: ''
})

const rules = {
  key: [
    { required: true, message: '请输入字典key', trigger: 'blur' }
  ],
  value: [
    { required: true, message: '请输入字典value', trigger: 'blur' }
  ]
}

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
    title: 'key',
    key: 'key',
    width: 220
  },
  {
    title: 'value',
    key: 'value',
    ellipsis: true
  },
  {
    title: '描述',
    key: 'desc',
    render: (row) => row.desc || '-'
  },
  {
    title: '创建时间',
    key: 'CreatedAt',
    width: 180,
    render: (row) => formatDate(row.CreatedAt)
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
                if (confirm('确定要删除该字典吗？')) {
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

const loadItems = async () => {
  loading.value = true
  try {
    const data = await configDictionaryAPI.getConfigDictionaries(searchParams)
    if (data && Array.isArray(data)) {
      items.value = data
    } else if (data && data.items && Array.isArray(data.items)) {
      items.value = data.items
    } else if (data && Array.isArray(data)) {
      items.value = data
    } else {
      items.value = []
    }
  } catch (error) {
    appStore.showError('加载配置字典失败')
    items.value = []
  } finally {
    loading.value = false
  }
}

const handleEdit = (row) => {
  isEdit.value = true
  currentId.value = row.ID || row.id
  form.key = row.key
  form.value = row.value
  form.desc = row.desc || ''
  showModal.value = true
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
    saving.value = true
    if (isEdit.value && currentId.value) {
      await configDictionaryAPI.updateConfigDictionary(currentId.value, {
        key: form.key,
        value: form.value,
        desc: form.desc
      })
      appStore.showSuccess('配置字典更新成功')
    } else {
      await configDictionaryAPI.createConfigDictionary({
        key: form.key,
        value: form.value,
        desc: form.desc
      })
      appStore.showSuccess('配置字典创建成功')
    }
    showModal.value = false
    loadItems()
  } catch (error) {
    appStore.showError(error?.response?.data?.message || '保存配置字典失败')
  } finally {
    saving.value = false
  }
}

const handleDelete = async (row) => {
  try {
    await configDictionaryAPI.deleteConfigDictionary(row.ID || row.id)
    appStore.showSuccess('配置字典删除成功')
    loadItems()
  } catch (error) {
    appStore.showError(error?.response?.data?.message || '删除配置字典失败')
  }
}

const resetForm = () => {
  form.key = ''
  form.value = ''
  form.desc = ''
  isEdit.value = false
  currentId.value = null
}

onMounted(() => {
  loadItems()
})
</script>
