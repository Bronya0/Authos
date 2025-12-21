<template>
  <div>
    <n-card>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <n-space align="center">
            <span>菜单列表</span>
            <n-input v-model:value="searchName" placeholder="搜索菜单名称" clearable @update:value="handleSearch" style="width: 200px" />
            <n-button @click="handleReset">重置</n-button>
          </n-space>
          <n-button type="primary" @click="showAddModal = true">
            <template #icon>
              <n-icon>
                <add />
              </n-icon>
            </template>
            添加菜单
          </n-button>
        </div>
      </template>

      <!-- 批量操作栏 -->
      <div v-if="checkedKeys.length > 0"
        style="margin-bottom: 16px; padding: 12px; background-color: #f5f5f5; border-radius: 6px;">
        <n-space align="center">
          <span>已选择 {{ checkedKeys.length }} 项</span>
          <n-button size="small" type="error" @click="handleBatchDelete">
            <template #icon>
              <n-icon>
                <Trash />
              </n-icon>
            </template>
            批量删除
          </n-button>
          <n-button size="small" type="warning" @click="handleBatchEdit">
            <template #icon>
              <n-icon>
                <Create />
              </n-icon>
            </template>
            批量编辑
          </n-button>
          <n-button size="small" @click="clearSelection">
            清除选择
          </n-button>
        </n-space>
      </div>

      <n-tree :data="menuTree" :render-label="renderLabel" :render-suffix="renderSuffix" :expand-on-click="true"
        :checkable="true" :cascade="true" :check-strategy="'all'" :checked-keys="checkedKeys"
        :indeterminate-keys="indeterminateKeys" :expanded-keys="expandedKeys" :selected-keys="selectedKeys"
        :default-expanded-keys="defaultExpandedKeys"
        :show-line="true" :indent="20" block-line key-field="ID" label-field="name" children-field="children"
        @update:checked-keys="handleCheckedKeysChange" @update:indeterminate-keys="handleIndeterminateKeysChange"
        @update:expanded-keys="handleExpandedKeysChange" @update:selected-keys="handleSelectedKeysChange" />
    </n-card>

    <!-- 添加/编辑菜单模态框 -->
    <n-modal v-model:show="showModal" :title="isEdit ? '编辑菜单' : '添加菜单'" preset="dialog" :show-icon="false"
      @after-leave="resetForm">
      <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="80px">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="form.name" placeholder="请输入菜单名称" />
        </n-form-item>

        <n-form-item label="路径" path="path">
          <n-input v-model:value="form.path" placeholder="请输入菜单路径" />
        </n-form-item>

        <n-form-item label="组件" path="component">
          <n-input v-model:value="form.component" placeholder="请输入组件路径" />
        </n-form-item>

        <n-form-item label="类型" path="type">
          <n-select v-model:value="form.type" placeholder="请选择菜单类型" :options="typeOptions" />
        </n-form-item>

        <n-form-item label="父菜单" path="parentId">
          <n-tree-select v-model:value="form.parentId" placeholder="请选择父菜单" :options="menuTreeOptions" clearable />
        </n-form-item>

        <n-form-item label="排序" path="sort">
          <n-input-number v-model:value="form.sort" :min="0" placeholder="请输入排序值" />
        </n-form-item>

        <n-form-item label="是否隐藏" path="hidden">
          <n-switch v-model:value="form.hidden" />
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
import { menuAPI } from '../api'
import { useAppStore } from '../stores/app'
import { List, Add, Create, Trash } from '@vicons/ionicons5'
import {
  NIcon, NButton, NSpace, NTag, NTreeSelect,
  NCard, NTree, NModal, NForm, NFormItem, NInput, NSelect,
  NInputNumber, NSwitch
} from 'naive-ui'

const appStore = useAppStore()

// 数据转换函数，确保符合n-tree组件要求
const transformMenuData = (data) => {
  if (!Array.isArray(data)) return []

  return data.map(item => {
    // 收集所有菜单ID作为默认展开的键
    if (item.ID || item.id) {
      defaultExpandedKeys.value.push(item.ID || item.id)
    }
    
    return {
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
    }
  })
}

const menuTree = ref([])
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const isEdit = ref(false)
const currentMenuId = ref(null)
const searchName = ref('')

// 处理搜索
const handleSearch = () => {
  loadMenus()
}

const handleReset = () => {
  searchName.value = ''
  loadMenus()
}

// 级联选择相关状态
const checkedKeys = ref([])
const indeterminateKeys = ref([])
const expandedKeys = ref([])
const selectedKeys = ref([])

const formRef = ref()
const form = reactive({
  name: '',
  path: '',
  component: '',
  type: 1, // 默认为菜单
  parentId: null,
  sort: 0,
  hidden: false
})

const rules = {
  name: [
    { required: true, message: '请输入菜单名称', trigger: 'blur' }
  ],
  path: [
    { required: true, message: '请输入菜单路径', trigger: 'blur' }
  ]
}

// 菜单类型选项
const typeOptions = [
  { label: '目录', value: 0 },
  { label: '菜单', value: 1 },
  { label: '按钮', value: 2 }
]

// 获取菜单类型文本
const getTypeText = (type) => {
  const typeMap = {
    0: '目录',
    1: '菜单',
    2: '按钮'
  }
  return typeMap[type] || '未知'
}

// 获取菜单类型颜色
const getTypeColor = (type) => {
  const colorMap = {
    0: 'success',
    1: 'info',
    2: 'warning'
  }
  return colorMap[type] || 'default'
}

// 将菜单数据转换为树形选择器选项
const menuTreeOptions = computed(() => {
  const convertToTreeOptions = (menus, level = 0) => {
    return menus.map(menu => ({
      label: menu.name,
      value: menu.ID,
      children: menu.children && menu.children.length > 0
        ? convertToTreeOptions(menu.children, level + 1)
        : undefined
    }))
  }

  return [
    { label: '根菜单', value: 0 },
    ...convertToTreeOptions(menuTree.value)
  ]
})

// 渲染树节点标签
const renderLabel = ({ option }) => {
  return h('div', { style: 'display: flex; align-items: center; gap: 8px;' }, [
    h('span', option.name),
    h(NTag, {
      type: getTypeColor(option.type),
      size: 'small'
    }, () => getTypeText(option.type))
  ])
}

// 渲染树节点后缀（操作按钮）
const renderSuffix = ({ option }) => {
  return h(NSpace, { size: 'small' }, {
    default: () => [
      // 只有目录和菜单类型才能添加子菜单
      (option.type === 0 || option.type === 1) ? h(
        NButton,
        {
          size: 'small',
          type: 'primary',
          quaternary: true,
          onClick: (e) => {
            e.stopPropagation()
            handleAddChild(option)
          }
        },
        {
          default: () => '添加子菜单',
          icon: () => h(NIcon, null, { default: () => h(Add) })
        }
      ) : null,
      h(
        NButton,
        {
          size: 'small',
          type: 'warning',
          quaternary: true,
          onClick: (e) => {
            e.stopPropagation()
            handleEdit(option)
          }
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
          onClick: (e) => {
            e.stopPropagation()
            if (confirm('确定要删除该菜单吗？')) {
              handleDelete(option)
            }
          }
        },
        {
          default: () => '删除',
          icon: () => h(NIcon, null, { default: () => h(Trash) })
        }
      )
    ]
  })
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

const getRowKey = (row) => row.id || row.ID || row.key || row.uuid || row.Id

const columns = [
  {
    title: 'ID',
    key: 'ID',
    width: 80,
    render: (row) => row.ID || row.id
  },
  {
    title: '名称',
    key: 'Name',
    render: (row) => row.Name || row.name || '-'
  },
  {
    title: '路径',
    key: 'Path',
    render: (row) => row.Path || row.path || '-'
  },
  {
    title: '图标',
    key: 'Icon',
    render: (row) => row.Icon || row.icon || '-'
  },
  {
    title: '父菜单',
    key: 'ParentID',
    render: (row) => {
      const parentId = row.ParentID || row.parentId || row.parent_id || 0
      if (parentId === 0) return '根菜单'
      const parent = menus.value.find(m => (m.ID || m.id) === parentId)
      return parent ? (parent.Name || parent.name || '未知菜单') : '-'
    }
  },
  {
    title: '排序',
    key: 'Sort',
    render: (row) => row.Sort || row.sort || 0,
    width: 100
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
                if (confirm('确定要删除该菜单吗？')) {
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

const loadMenus = async () => {
  loading.value = true
  // 清空之前的展开键
  defaultExpandedKeys.value = []
  try {
    // 使用标准菜单树API（现在所有数据都是用户自定义的）
    // 如果有搜索，使用列表API
    let data;
    if (searchName.value) {
      data = await menuAPI.getMenus({ name: searchName.value })
      // 搜索模式下展示扁平结构或构建临时树
      menuTree.value = transformMenuData(data || [])
    } else {
      data = await menuAPI.getMenuTree()
      const transformedData = transformMenuData(data || [])
      menuTree.value = transformedData
      expandedKeys.value = getAllMenuIds(transformedData)
    }
  } catch (error) {
    appStore.showError('加载菜单列表失败')
  } finally {
    loading.value = false
  }
}

const handleEdit = (row) => {
  isEdit.value = true
  currentMenuId.value = row.ID

  // 填充表单数据
  form.name = row.name || ''
  form.path = row.path || ''
  form.component = row.component || ''
  form.type = row.type || 1
  form.parentId = row.parentId || 0
  form.sort = row.sort || 0
  form.hidden = row.hidden || false

  showModal.value = true
}

// 处理添加子菜单
const handleAddChild = (parentOption) => {
  isEdit.value = false
  currentMenuId.value = null

  // 重置表单
  resetForm()

  // 设置父菜单ID
  form.parentId = parentOption.ID

  showModal.value = true
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
    saving.value = true

    const data = {
      name: form.name,
      path: form.path,
      component: form.component,
      type: form.type,
      parentId: form.parentId || 0,
      sort: form.sort,
      hidden: form.hidden
    }

    if (isEdit.value) {
      await menuAPI.updateMenu(currentMenuId.value, data)
      appStore.showSuccess('菜单更新成功')
    } else {
      await menuAPI.createMenu(data)
      appStore.showSuccess('菜单创建成功')
    }

    showModal.value = false
    loadMenus()
  } catch (error) {
    appStore.showError(error.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const handleDelete = async (row) => {
  try {
    const id = row.ID
    await menuAPI.deleteMenu(id)
    appStore.showSuccess('菜单删除成功')
    loadMenus()
  } catch (error) {
    appStore.showError('删除失败')
  }
}

const resetForm = () => {
  form.name = ''
  form.path = ''
  form.component = ''
  form.type = 1 // 默认为菜单
  form.parentId = null
  form.sort = 0
  form.hidden = false
  isEdit.value = false
  currentMenuId.value = null
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

// 级联选择事件处理
const handleCheckedKeysChange = (keys, options, meta) => {
  if (meta.action === 'check') {
    // 选中时，同步选中所有子节点
    const node = options.find(opt => opt.ID === meta.node.ID)
    if (node) {
      const childIds = getChildrenIds(node)
      const newKeys = [...new Set([...keys, ...childIds])]
      checkedKeys.value = newKeys
    } else {
      checkedKeys.value = keys
    }
  } else if (meta.action === 'uncheck') {
    // 取消选中时，同步取消选中所有子节点
    const childIds = getChildrenIds(meta.node)
    const newKeys = keys.filter(key => !childIds.includes(key))
    checkedKeys.value = newKeys
  } else {
    checkedKeys.value = keys
  }
}

const handleIndeterminateKeysChange = (keys) => {
  indeterminateKeys.value = keys
}

const handleExpandedKeysChange = (keys, options, meta) => {
  expandedKeys.value = keys
}

const handleSelectedKeysChange = (keys, options, meta) => {
  selectedKeys.value = keys
}

// 清除选择
const clearSelection = () => {
  checkedKeys.value = []
  indeterminateKeys.value = []
  selectedKeys.value = []
  // 保持菜单展开状态，不清空展开键
}

// 批量删除
const handleBatchDelete = async () => {
  if (checkedKeys.value.length === 0) {
    appStore.showWarning('请先选择要删除的菜单')
    return
  }

  if (confirm(`确定要删除选中的 ${checkedKeys.value.length} 个菜单吗？此操作不可撤销！`)) {
    try {
      const deletePromises = checkedKeys.value.map(id => menuAPI.deleteMenu(id))
      await Promise.all(deletePromises)
      appStore.showSuccess(`成功删除 ${checkedKeys.value.length} 个菜单`)
      clearSelection()
      loadMenus()
    } catch (error) {
      appStore.showError('批量删除失败')
    }
  }
}

// 批量编辑
const handleBatchEdit = () => {
  if (checkedKeys.value.length === 0) {
    appStore.showWarning('请先选择要编辑的菜单')
    return
  }

  if (checkedKeys.value.length > 1) {
    appStore.showWarning('批量编辑功能暂时只支持单个菜单的编辑，请选择一个菜单')
    return
  }

  // 获取选中的菜单数据
  const selectedMenu = findMenuById(menuTree.value, checkedKeys.value[0])
  if (selectedMenu) {
    handleEdit(selectedMenu)
  }
}

// 根据ID查找菜单
const findMenuById = (menus, id) => {
  for (const menu of menus) {
    if (menu.ID === id) {
      return menu
    }
    if (menu.children && menu.children.length > 0) {
      const found = findMenuById(menu.children, id)
      if (found) return found
    }
  }
  return null
}

onMounted(() => {
  loadMenus()
})
</script>