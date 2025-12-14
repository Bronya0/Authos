import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { create, NButton, NInput, NForm, NFormItem, NCard, NLayout, NLayoutSider, NLayoutHeader, NLayoutContent, NMenu, NBreadcrumb, NBreadcrumbItem, NIcon, NDropdown, NAvatar, NTag, NSelect, NModal, NDataTable, NSpace, NGrid, NGridItem, NStatistic, NSpin, NAlert, NBadge, NTree, NCheckbox, NRadio, NRadioGroup, NPagination, NDivider, NTooltip, NPopconfirm, NInputNumber, NConfigProvider, NMessageProvider, NModalProvider } from 'naive-ui'
import App from './App.vue'
import router from './router'

const naive = create({
  components: [
    NButton, NInput, NForm, NFormItem, NCard, NLayout, NLayoutSider, NLayoutHeader, 
    NLayoutContent, NMenu, NBreadcrumb, NBreadcrumbItem, NIcon, NDropdown, NAvatar, 
    NTag, NSelect, NModal, NDataTable, NSpace, NGrid, NGridItem, NStatistic, 
    NSpin, NAlert, NBadge, NTree, NCheckbox, NRadio, NRadioGroup, NPagination, 
    NDivider, NTooltip, NPopconfirm, NInputNumber, NConfigProvider, NMessageProvider,
    NModalProvider
  ]
})

const pinia = createPinia()
const app = createApp(App)

app.use(naive)
app.use(router)
app.use(pinia)

app.mount('#app')