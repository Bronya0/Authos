import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    // 系统管理员认证
    const systemToken = localStorage.getItem('systemToken')
    if (systemToken) {
      config.headers['X-System-Token'] = `Bearer ${systemToken}`
    }
    
    // 应用认证
    const appToken = localStorage.getItem('appToken')
    if (appToken) {
      config.headers['X-App-Token'] = `Bearer ${appToken}`
    }
    
    // 用户认证（多租户）
    const userToken = localStorage.getItem('userToken')
    if (userToken) {
      config.headers['Authorization'] = `Bearer ${userToken}`
    }
    
    // 从localStorage获取应用ID并添加到请求头
    const currentApp = localStorage.getItem('currentApp')
    if (currentApp) {
      try {
        const appData = JSON.parse(currentApp)
        if (appData.id) {
          config.headers['X-App-ID'] = appData.id
        }
      } catch (e) {
        console.error('Failed to parse currentApp from localStorage', e)
      }
    }
    
    return config
  },
  error => Promise.reject(error)
)

// 响应拦截器
api.interceptors.response.use(
  response => response.data,
  error => {
    // 检查是否是登录相关的请求，如果是则不进行全局401跳转，交由页面自行处理
    const isLoginRequest = error.config?.url && (
      error.config.url.includes('/public/app-login') || 
      error.config.url.includes('/public/system-login') || 
      error.config.url.includes('/public/login')
    )

    if (error.response?.status === 401 && !isLoginRequest) {
      // 清除所有认证信息
      localStorage.removeItem('systemToken')
      localStorage.removeItem('systemUser')
      localStorage.removeItem('appToken')
      localStorage.removeItem('currentApp')
      localStorage.removeItem('userToken')
      
      // 重定向到系统登录页面
      window.location.href = '/system-login'
    }
    return Promise.reject(error)
  }
)

export const authAPI = {
  // 系统管理员登录
  systemLogin: (data) => api.post('/public/system-login', data),
  // 应用登录
  appLogin: (data) => api.post('/public/app-login', data),
  // 用户登录
  login: (data) => api.post('/public/login', data),
  // 登出
  logout: () => api.post('/public/logout'),
  // 获取仪表盘统计
  getDashboardStats: () => api.get('/v1/dashboard/stats')
}

export const userAPI = {
  getUsers: (params) => api.get('/v1/users', { params }),
  getUser: (id) => api.get(`/v1/users/${id}`),
  createUser: (data) => api.post('/v1/users', data),
  updateUser: (id, data) => api.put(`/v1/users/${id}`, data),
  deleteUser: (id) => api.delete(`/v1/users/${id}`)
}

export const roleAPI = {
  getRoles: (params) => api.get('/v1/roles', { params }),
  getRole: (id) => api.get(`/v1/roles/${id}`),
  createRole: (data) => api.post('/v1/roles', data),
  updateRole: (id, data) => api.put(`/v1/roles/${id}`, data),
  deleteRole: (id) => api.delete(`/v1/roles/${id}`),
  getRoleMenus: (id) => api.get(`/v1/roles/${id}/menus`),
  updateRoleMenus: (id, data) => api.put(`/v1/roles/${id}/menus`, data),
  getRolePermissions: (id) => api.get(`/v1/roles/${id}/permissions`),
  updateRolePermissions: (id, data) => api.put(`/v1/roles/${id}/permissions`, data),
  // 使用现有的API获取菜单和权限数量
  getRoleMenusCount: async (id) => {
    try {
      const menus = await api.get(`/v1/roles/${id}/menus`)
      return Array.isArray(menus) ? menus.length : 0
    } catch (error) {
      console.error('获取角色菜单数量失败:', error)
      return 0
    }
  },
  getRoleApiPermissionsCount: async (id) => {
    try {
      // 需要先获取角色的UUID
      const role = await api.get(`/v1/roles/${id}`)
      const roleUUID = role.uuid || role.UUID
      if (!roleUUID) {
        console.error('角色UUID不存在')
        return 0
      }
      const permissions = await api.get(`/v1/api-permissions/roles/${roleUUID}`)
      return Array.isArray(permissions) ? permissions.length : 0
    } catch (error) {
      console.error('获取角色接口权限数量失败:', error)
      return 0
    }
  }
}

export const menuAPI = {
  getMenus: (params) => api.get('/v1/menus', { params }),
  getMenuTree: () => api.get('/v1/menus/tree'),
  getNonSystemMenuTree: () => api.get('/v1/menus/non-system-tree'),
  getMenu: (id) => api.get(`/v1/menus/${id}`),
  createMenu: (data) => api.post('/v1/menus', data),
  updateMenu: (id, data) => api.put(`/v1/menus/${id}`, data),
  deleteMenu: (id) => api.delete(`/v1/menus/${id}`)
}

export const permissionAPI = {
  getPermissions: (params) => api.get('/v1/permissions', { params }),
  getPermission: (id) => api.get(`/v1/permissions/${id}`),
  createPermission: (data) => api.post('/v1/permissions', data),
  updatePermission: (id, data) => api.put(`/v1/permissions/${id}`, data),
  deletePermission: (id) => api.delete(`/v1/permissions/${id}`)
}

export const apiPermissionAPI = {
  getApiPermissions: (params) => api.get('/v1/api-permissions', { params }),
  getApiPermission: (id) => api.get(`/v1/api-permissions/${id}`),
  createApiPermission: (data) => api.post('/v1/api-permissions', data),
  updateApiPermission: (id, data) => api.put(`/v1/api-permissions/${id}`, data),
  deleteApiPermission: (id) => api.delete(`/v1/api-permissions/${id}`),
  getApiPermissionsForRole: (roleUUID) => api.get(`/v1/api-permissions/roles/${roleUUID}`),
  addApiPermissionToRole: (roleUUID, data) => api.post(`/v1/api-permissions/roles/${roleUUID}`, data),
  removeApiPermissionFromRole: (roleUUID, data) => api.delete(`/v1/api-permissions/roles/${roleUUID}`, { data })
}

export const configDictionaryAPI = {
  getConfigDictionaries: (params) => api.get('/v1/config-dictionaries', { params }),
  getConfigDictionary: (id) => api.get(`/v1/config-dictionaries/${id}`),
  createConfigDictionary: (data) => api.post('/v1/config-dictionaries', data),
  updateConfigDictionary: (id, data) => api.put(`/v1/config-dictionaries/${id}`, data),
  deleteConfigDictionary: (id) => api.delete(`/v1/config-dictionaries/${id}`)
}

export const applicationAPI = {
  getApplications: (params) => api.get('/v1/applications', { params }),
  getApplication: (id) => api.get(`/v1/applications/${id}`),
  createApplication: (data) => api.post('/v1/applications', data),
  updateApplication: (id, data) => api.put(`/v1/applications/${id}`, data),
  deleteApplication: (id) => api.delete(`/v1/applications/${id}`),
  getApplicationByCode: (code) => api.get(`/v1/applications/by-code/${code}`)
}

export const auditLogAPI = {
  getLogs: (params) => api.get('/v1/audit-logs', { params }),
  getSystemLogs: (params) => api.get('/v1/system/audit-logs', { params })
}

export default api
