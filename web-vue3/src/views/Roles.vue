<template>
    <div>
        <n-card>
            <template #header>
                <div style="display: flex; justify-content: space-between; align-items: center;">
                    <span>角色列表</span>
                    <n-button type="primary" @click="showAddModal = true">
                        <template #icon>
                            <n-icon>
                                <add />
                            </n-icon>
                        </template>
                        添加角色
                    </n-button>
                </div>
            </template>

            <n-data-table :columns="columns" :data="roles" :loading="loading" :pagination="pagination" />
        </n-card>

        <!-- 添加/编辑角色模态框 -->
        <n-modal v-model:show="showModal" :title="isEdit ? '编辑角色' : '添加角色'" preset="dialog" :show-icon="false"
            @after-leave="resetForm">
            <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="80px">
                <n-form-item label="名称" path="name">
                    <n-input v-model:value="form.name" placeholder="请输入角色名称" />
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

        <!-- 菜单权限分配模态框 -->
        <RoleMenuModal v-model:visible="showMenuModal" :role-id="currentRoleId" :role-name="currentRoleName"
            @saved="loadRoles" />

        <!-- 接口权限分配模态框 -->
        <RoleApiPermissionModal v-model:visible="showApiPermissionModal" :role-u-u-i-d="currentRoleUUID"
            :role-name="currentRoleName" @saved="loadRoles" />
    </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, h } from 'vue'
import { roleAPI } from '../api'
import { useAppStore } from '../stores/app'
import { Person, Add, Create, Trash, List, Key } from '@vicons/ionicons5'
import { NButton, NIcon, NSpace } from 'naive-ui'
import RoleMenuModal from '../components/RoleMenuModal.vue'
import RoleApiPermissionModal from '../components/RoleApiPermissionModal.vue'

const appStore = useAppStore()

const roles = ref([])
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const isEdit = ref(false)
const currentRoleId = ref(null)
const currentRoleUUID = ref('')
const currentRoleName = ref('')
const showMenuModal = ref(false)
const showApiPermissionModal = ref(false)

const formRef = ref()
const form = reactive({
    name: ''
})

const rules = {
    name: [
        { required: true, message: '请输入角色名称', trigger: 'blur' }
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
        title: '名称',
        key: 'name',
        width: 500
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
                                if (confirm('确定要删除该角色吗？')) {
                                    handleDelete(row)
                                }
                            }
                        },
                        {
                            default: () => '删除',
                            icon: () => h(NIcon, null, { default: () => h(Trash) })
                        }
                    ),
                    h(
                        NButton,
                        {
                            size: 'small',
                            type: 'info',
                            quaternary: true,
                            onClick: () => handleRoleMenus(row)
                        },
                        {
                            default: () => '菜单',
                            icon: () => h(NIcon, null, { default: () => h(List) })
                        }
                    ),
                    h(
                        NButton,
                        {
                            size: 'small',
                            type: 'info',
                            quaternary: true,
                            onClick: () => handleRolePermissions(row)
                        },
                        {
                            default: () => '权限',
                            icon: () => h(NIcon, null, { default: () => h(Key) })
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

const loadRoles = async () => {
    loading.value = true
    try {
        const data = await roleAPI.getRoles()
        roles.value = data
    } catch (error) {
        appStore.showError('加载角色列表失败')
    } finally {
        loading.value = false
    }
}

const handleEdit = (row) => {
    isEdit.value = true
    currentRoleId.value = row.ID || row.id
    currentRoleUUID.value = row.uuid || row.UUID || ''

    // 填充表单数据
    form.name = row.name

    showModal.value = true
}

const handleSave = async () => {
    try {
        await formRef.value?.validate()
        saving.value = true

        const data = {
            name: form.name
        }

        if (isEdit.value) {
            await roleAPI.updateRole(currentRoleId.value, data)
            appStore.showSuccess('角色更新成功')
        } else {
            await roleAPI.createRole(data)
            appStore.showSuccess('角色创建成功')
        }

        showModal.value = false
        loadRoles()
    } catch (error) {
        appStore.showError(error.message || '保存失败')
    } finally {
        saving.value = false
    }
}

const handleDelete = async (row) => {
    try {
        const id = row.ID || row.id
        await roleAPI.deleteRole(id)
        appStore.showSuccess('角色删除成功')
        loadRoles()
    } catch (error) {
        appStore.showError('删除失败')
    }
}

const handleRoleMenus = (row) => {
    currentRoleId.value = row.ID || row.id
    currentRoleName.value = row.name || row.Name
    showMenuModal.value = true
}

const handleRolePermissions = (row) => {
    currentRoleUUID.value = row.uuid || row.UUID
    currentRoleName.value = row.name || row.Name
    showApiPermissionModal.value = true
}

const resetForm = () => {
    form.name = ''
    isEdit.value = false
    currentRoleId.value = null
    currentRoleUUID.value = ''
    currentRoleName.value = ''
}

onMounted(() => {
    loadRoles()
})
</script>