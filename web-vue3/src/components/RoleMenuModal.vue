<template>
  <n-drawer v-model:show="showModal" :width="500" placement="right" @after-leave="resetForm">
    <n-drawer-content title="菜单权限分配" closable>
      <n-spin :show="loading">
        <n-tree 
          :data="menuTree" 
          :checked-keys="selectedMenuIds" 
          :on-update:checked-keys="handleCheck"
          :expand-on-click="true" 
          :selectable="false" 
          :checkable="true" 
          :show-line="true" 
          :indent="20" 
          :cascade="false"
          :expanded-keys="expandedKeys"
          @update:expanded-keys="(keys) => expandedKeys = keys"
          block-line 
          key-field="ID" 
          label-field="name" 
          children-field="children" 
        />
      </n-spin>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showModal = false">取消</n-button>
          <n-button type="primary" :loading="saving" @click="handleSave">
            保存
          </n-button>
        </n-space>
      </template>
    </n-drawer-content>
  </n-drawer>
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

const showModal = ref(false)
const loading = ref(false)
const saving = ref(false)
const menuTree = ref([])
const selectedMenuIds = ref([])
const expandedKeys = ref([])

// 获取所有菜单ID用于默认展开
const getAllMenuIds = (menus) => {
  let ids = []
  menus.forEach(menu => {
    ids.push(menu.ID)
    if (menu.children && menu.children.length > 0) {
      ids = ids.concat(getAllMenuIds(menu.children))
    }
  })
  return ids
}

// 获取子节点的所有ID
const getChildrenIds = (node) => {
  let ids = []
  if (node.children && node.children.length > 0) {
    node.children.forEach(child => {
      ids.push(child.ID)
      ids = ids.concat(getChildrenIds(child))
    })
  }
  return ids
}

// 查找节点
const findNode = (menus, id) => {
  for (const menu of menus) {
    if (menu.ID === id) return menu
    if (menu.children && menu.children.length > 0) {
      const found = findNode(menu.children, id)
      if (found) return found
    }
  }
  return null
}

// 监听visible变化
watch(() => props.visible, (val) => {
  showModal.value = val
  if (val && props.roleId) {
    loadMenus()
    loadRoleMenus()
  }
})

// 监听showModal变化，同步到父组件
watch(showModal, (val) => {
  emit('update:visible', val)
})

// 加载菜单树
const loadMenus = async () => {
  loading.value = true
  try {
    const data = await menuAPI.getMenuTree()
    const transformedData = transformMenuData(data || [])
    menuTree.value = transformedData
    expandedKeys.value = getAllMenuIds(transformedData)
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
  if (meta.action === 'check') {
    // 选中时，同步选中所有子节点
    const node = findNode(menuTree.value, meta.node.ID)
    if (node) {
      const childIds = getChildrenIds(node)
      selectedMenuIds.value = [...new Set([...keys, ...childIds])]
    } else {
      selectedMenuIds.value = keys
    }
  } else if (meta.action === 'uncheck') {
    // 取消选中时，同步取消选中所有子节点
    const childIds = getChildrenIds(meta.node)
    selectedMenuIds.value = keys.filter(key => !childIds.includes(key))
  } else {
    selectedMenuIds.value = keys
  }
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

const resetForm = () => {
  selectedMenuIds.value = []
}
</script>