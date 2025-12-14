<template>
  <n-modal v-model:show="showModal" preset="dialog" title="菜单权限分配" style="width: 600px;">
    <div style="margin-bottom: 16px;">
      <n-alert type="info" :show-icon="false">
        为角色【{{ roleName }}】分配菜单权限
      </n-alert>
    </div>

    <n-spin :show="loading">
      <n-tree :data="menuTree" :checked-keys="selectedMenuIds" :on-update:checked-keys="handleCheck"
        :expand-on-click="true" :selectable="false" :checkable="true" :show-line="true" :indent="20" :cascade="true"
        block-line key-field="ID" label-field="name" children-field="children" />
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
import { ref, watch } from 'vue'
import { roleAPI, menuAPI } from '../api'
import { useAppStore } from '../stores/app'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  roleId: {
    type: [String, Number],
    default: null
  },
  roleName: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:visible', 'saved'])

const appStore = useAppStore()

// 数据转换函数，确保符合n-tree组件要求
const transformMenuData = (data) => {
  if (!Array.isArray(data)) return []

  return data.map(item => ({
    // 确保关键字段存在
    ID: item.ID || item.id,
    name: item.name || item.Name,
    type: item.type || item.Type,
    parentId: item.parentId || item.ParentID,
    path: item.path || item.Path,
    component: item.component || item.Component,
    sort: item.sort || item.Sort,
    hidden: item.hidden || item.Hidden,
    isSystem: item.isSystem || item.IsSystem,
    // 递归处理children
    children: item.children && item.children.length > 0
      ? transformMenuData(item.children)
      : []
  }))
}

const showModal = ref(props.visible)
const loading = ref(false)
const saving = ref(false)
const menuTree = ref([])
const selectedMenuIds = ref([])
const allMenuIds = ref([])

// 监听visible变化
watch(() => props.visible, (val) => {
  showModal.value = val
  if (val && props.roleId) {
    loadMenus()
    loadRoleMenus()
  }
})

// 监 showModal变化，同步到父组件
watch(showModal, (val) => {
  emit('update:visible', val)
})

// 加载菜单树
const loadMenus = async () => {
  loading.value = true
  try {
    const data = await menuAPI.getMenuTree()
    menuTree.value = transformMenuData(data || [])

    // 收集所有菜单ID
    const collectIds = (menus) => {
      menus.forEach(menu => {
        allMenuIds.value.push(menu.ID)
        if (menu.children && menu.children.length > 0) {
          collectIds(menu.children)
        }
      })
    }
    collectIds(menuTree.value)
  } catch (error) {
    appStore.showError('加载菜单列表失败')
  } finally {
    loading.value = false
  }
}

// 加载角色已有的菜单
const loadRoleMenus = async () => {
  try {
    const data = await roleAPI.getRoleMenus(props.roleId)
    // 转换角色菜单数据并获取ID
    const transformedData = transformMenuData(data || [])
    selectedMenuIds.value = transformedData.map(menu => menu.ID)
  } catch (error) {
    appStore.showError('加载角色菜单失败')
  }
}

// 处理菜单选择
const handleCheck = (keys, options, meta) => {
  selectedMenuIds.value = keys
  console.log("菜单选择变化:", keys, options, meta)
}

// 保存菜单分配
const handleSave = async () => {
  saving.value = true
  try {
    await roleAPI.updateRoleMenus(props.roleId, { menuIds: selectedMenuIds.value })
    appStore.showSuccess('菜单权限分配成功')
    showModal.value = false
    emit('saved')
  } catch (error) {
    appStore.showError(error.message || '菜单权限分配失败')
  } finally {
    saving.value = false
  }
}
</script>