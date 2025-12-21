<template>
  <div style="display: flex; flex-direction: column; gap: 16px;">
    <n-card>
      <n-space vertical>
        <n-space align="center">
          <n-select v-model:value="searchParams.action" placeholder="操作类型" clearable :options="actionOptions" style="width: 150px" @update:value="handleSearch" />
          <n-select v-model:value="searchParams.resource" placeholder="资源类型" clearable :options="resourceOptions" style="width: 150px" @update:value="handleSearch" />
          <n-input v-model:value="searchParams.username" placeholder="操作人" clearable style="width: 150px" @update:value="handleSearch" />
          <n-button type="primary" @click="handleSearch">查询</n-button>
        </n-space>
        
        <n-data-table :columns="columns" :data="logs" :loading="loading" :pagination="pagination" />
      </n-space>
    </n-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, h, watch } from 'vue'
import { useRoute } from 'vue-router'
import { auditLogAPI } from '../api'
import { formatDate } from '../utils/format'
import { NTag, NCode, NSpace, NCard, NDataTable, NSelect, NInput, NButton } from 'naive-ui'

const route = useRoute()
const isSystem = ref(route.meta.system || false)

const logs = ref([])
const loading = ref(false)
const searchParams = reactive({
  action: null,
  resource: null,
  username: ''
})

watch(() => route.meta.system, (val) => {
  isSystem.value = val || false
  loadLogs()
})

const actionOptions = [
  { label: '登录', value: 'LOGIN' },
  { label: '登出', value: 'LOGOUT' },
  { label: '创建', value: 'CREATE' },
  { label: '更新', value: 'UPDATE' },
  { label: '删除', value: 'DELETE' },
  { label: '分配', value: 'ASSIGN' },
  { label: '取消分配', value: 'UNASSIGN' },
  { label: '系统登录', value: 'SYSTEM_LOGIN' },
  { label: '应用登录', value: 'APP_LOGIN' }
]

const resourceOptions = [
  { label: '用户', value: 'USER' },
  { label: '角色', value: 'ROLE' },
  { label: '菜单', value: 'MENU' },
  { label: '接口权限', value: 'API_PERMISSION' },
  { label: '应用', value: 'APPLICATION' },
  { label: '角色权限关联', value: 'ROLE_PERMISSION' }
]

const getActionColor = (action) => {
  if (action.includes('LOGIN')) return 'success'
  if (action === 'CREATE') return 'info'
  if (action === 'UPDATE') return 'warning'
  if (action === 'DELETE') return 'error'
  return 'default'
}

const columns = [
  { title: '时间', key: 'CreatedAt', width: 180, render: (row) => formatDate(row.CreatedAt) },
  { title: '操作人', key: 'username', width: 120 },
  { 
    title: '动作', 
    key: 'action', 
    width: 120,
    render: (row) => h(NTag, { type: getActionColor(row.action), size: 'small' }, { default: () => row.action })
  },
  { title: '资源', key: 'resource', width: 120 },
  { title: '资源ID', key: 'resourceId', width: 100 },
  { title: '内容', key: 'content' },
  { title: 'IP', key: 'ip', width: 130 },
  { 
    title: '状态', 
    key: 'status', 
    width: 80,
    render: (row) => h(NTag, { type: row.status === 1 ? 'success' : 'error', size: 'small' }, { default: () => row.status === 1 ? '成功' : '失败' })
  }
]

const pagination = { pageSize: 15 }

const loadLogs = async () => {
  loading.value = true
  try {
    const apiFunc = isSystem.value ? auditLogAPI.getSystemLogs : auditLogAPI.getLogs
    const data = await apiFunc(searchParams)
    logs.value = data || []
  } catch (error) {
    console.error('Failed to load audit logs:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  loadLogs()
}

onMounted(loadLogs)
</script>
