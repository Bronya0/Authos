<template>
  <div>
    <n-grid :cols="3" :x-gap="24" :y-gap="24">
      <n-grid-item>
        <n-card>
          <n-statistic label="用户总数" :value="stats.users">
            <template #prefix>
              <n-icon size="24" color="#1890ff">
                <people />
              </n-icon>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>

      <n-grid-item>
        <n-card>
          <n-statistic label="角色总数" :value="stats.roles">
            <template #prefix>
              <n-icon size="24" color="#52c41a">
                <person />
              </n-icon>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>

      <n-grid-item>
        <n-card>
          <n-statistic label="菜单总数" :value="stats.menus">
            <template #prefix>
              <n-icon size="24" color="#faad14">
                <list />
              </n-icon>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>
    </n-grid>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { userAPI, roleAPI, menuAPI } from '../api'
import { People, Person, List } from '@vicons/ionicons5'

const stats = ref({
  users: 0,
  roles: 0,
  menus: 0
})

const loadStats = async () => {
  try {
    const [users, roles, menus] = await Promise.all([
      userAPI.getUsers().catch(() => []),
      roleAPI.getRoles().catch(() => []),
      menuAPI.getMenus().catch(() => [])
    ])

    stats.value = {
      users: users.length,
      roles: roles.length,
      menus: menus.length
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

onMounted(() => {
  loadStats()
})
</script>