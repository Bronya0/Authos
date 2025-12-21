<template>
  <div>
    <n-grid :cols="5" :x-gap="24" :y-gap="24">
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

      <n-grid-item>
        <n-card>
          <n-statistic label="接口权限" :value="stats.apiPerms">
            <template #prefix>
              <n-icon size="24" color="#eb2f96">
                <shield-checkmark />
              </n-icon>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>

      <n-grid-item>
        <n-card>
          <n-statistic label="审计日志" :value="stats.auditLogs">
            <template #prefix>
              <n-icon size="24" color="#722ed1">
                <document-text />
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
import { authAPI } from '../api'
import { People, Person, List, ShieldCheckmark, DocumentText } from '@vicons/ionicons5'

const stats = ref({
  users: 0,
  roles: 0,
  menus: 0,
  apiPerms: 0,
  auditLogs: 0
})

const loadStats = async () => {
  try {
    const data = await authAPI.getDashboardStats()
    stats.value = data
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

onMounted(() => {
  loadStats()
})
</script>