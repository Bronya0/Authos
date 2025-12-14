<template>
  <n-modal v-model:show="showModal" preset="dialog" title="接口权限分配" style="width: 800px;">
    <div style="margin-bottom: 16px;">
      <n-alert type="info" :show-icon="false">
        为角色【{{ roleName }}】分配接口权限
      </n-alert>
    </div>

    <n-spin :show="loading">
      <n-data-table :columns="columns" :data="apiPermissions" :row-key="getRowKey"
        :checked-row-keys="selectedPermissionIds" @update:checked-row-keys="handleCheck" :pagination="pagination" />
    </n-spin>

    <template #action>
      <n-space>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" :loading="saving" @click="handleSave">
          保存
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script setup>
import { ref, watch, h } from 'vue'
import { roleAPI, apiPermissionAPI } from '../api'
import { useAppStore } from '../stores/app'
import { NTag } from 'naive-ui'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  roleUUID: {
    type: String,
    default: ''
  },
  roleName: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:visible', 'saved'])

const appStore = useAppStore()

const showModal = ref(props.visible)
const loading = ref(false)
const saving = ref(false)
const apiPermissions = ref([])
const selectedPermissionIds = ref([])

// 监听visible变化
watch(() => props.visible, (val) => {
  showModal.value = val
  if (val && props.roleUUID) {
    loadApiPermissions()
    loadRoleApiPermissions()
  }
})

// 监 showModal变化，同步到父组件
watch(showModal, (val) => {
  emit('update:visible', val)
})

// 获取HTTP方法颜色
const getMethodColor = (method) => {
  const colorMap = {
    'GET': 'success',
    'POST': 'info',
    'PUT': 'warning',
    'DELETE': 'error',
    'PATCH': 'default',
    'HEAD': 'default',
    'OPTIONS': 'default'
  }
  return colorMap[method] || 'default'
}

// 表格列定义
const columns = [
  {
    type: 'selection'
  },
  {
    title: '权限名称',
    key: 'name',
    width: 200
  },
  {
    title: '接口路径',
    key: 'path',
    width: 250
  },
  {
    title: 'HTTP方法',
    key: 'method',
    width: 120,
    render: (row) => h(NTag, {
      type: getMethodColor(row.method)
    }, () => row.method)
  },
  {
    title: '描述',
    key: 'description'
  }
]

// 分页配置
const pagination = {
  pageSize: 10
}

// 获取行键值
const getRowKey = (row) => row.uuid || row.UUID

// 加载所有接口权限
const loadApiPermissions = async () => {
  loading.value = true
  try {
    const data = await apiPermissionAPI.getApiPermissions()
    apiPermissions.value = data || []
  } catch (error) {
    appStore.showError('加载接口权限列表失败')
  } finally {
    loading.value = false
  }
}

// 加载角色已有的接口权限
const loadRoleApiPermissions = async () => {
  try {
    const data = await apiPermissionAPI.getApiPermissionsForRole(props.roleUUID)
    // 兼容null或undefined返回
    if (data && Array.isArray(data)) {
      selectedPermissionIds.value = data.map(p => p.uuid || p.UUID)
    } else {
      selectedPermissionIds.value = []
    }
  } catch (error) {
    appStore.showError('加载角色接口权限失败')
    selectedPermissionIds.value = []
  }
}

// 处理权限选择
const handleCheck = (keys) => {
  selectedPermissionIds.value = keys
}

// 保存权限分配
const handleSave = async () => {
  saving.value = true
  try {
    // 获取当前角色已有的权限
    const currentPermissions = await apiPermissionAPI.getApiPermissionsForRole(props.roleUUID)
    const currentPermissionIds = currentPermissions && Array.isArray(currentPermissions)
      ? currentPermissions.map(p => p.uuid || p.UUID)
      : []

    // 计算需要添加和删除的权限
    const toAdd = selectedPermissionIds.value.filter(id => !currentPermissionIds.includes(id))
    const toRemove = currentPermissionIds.filter(id => !selectedPermissionIds.value.includes(id))

    // 添加新权限
    for (const permissionId of toAdd) {
      await apiPermissionAPI.addApiPermissionToRole(props.roleUUID, { permissionUUID: permissionId })
    }

    // 删除不再需要的权限
    for (const permissionId of toRemove) {
      await apiPermissionAPI.removeApiPermissionFromRole(props.roleUUID, { permissionUUID: permissionId })
    }

    appStore.showSuccess('接口权限分配成功')
    showModal.value = false
    emit('saved')
  } catch (error) {
    appStore.showError(error.message || '接口权限分配失败')
  } finally {
    saving.value = false
  }
}
</script>