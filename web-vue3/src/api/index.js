import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
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
    if (error.response?.status === 401) {
      // 清除所有认证信息
      localStorage.removeItem('systemToken')
      localStorage.removeItem('systemUser')
      localStorage.removeItem('appToken')
      localStorage.removeItem('currentApp')
      
      // 重定向到系统登录页面
      window.location.href = '/system-login'
    }
    return Promise.reject(error)
  }
)

export const authAPI = {
  // 系统管理员登录
  systemLogin: (data) => api.post('/system-login', data),
  // 应用登录
  appLogin: (data) => api.post('/app-login', data),
  // 登出
  logout: () => api.post('/logout')
}

export const userAPI = {
  getUsers: () => api.get('/users'),
  getUser: (id) => api.get(`/users/${id}`),
  createUser: (data) => api.post('/users', data),
  updateUser: (id, data) => api.put(`/users/${id}`, data),
  deleteUser: (id) => api.delete(`/users/${id}`)
}

export const roleAPI = {
  getRoles: () => api.get('/roles'),
  getRole: (id) => api.get(`/roles/${id}`),
  createRole: (data) => api.post('/roles', data),
  updateRole: (id, data) => api.put(`/roles/${id}`, data),
  deleteRole: (id) => api.delete(`/roles/${id}`),
  getRoleMenus: (id) => api.get(`/roles/${id}/menus`),
  updateRoleMenus: (id, data) => api.put(`/roles/${id}/menus`, data),
  getRolePermissions: (id) => api.get(`/roles/${id}/permissions`),
  updateRolePermissions: (id, data) => api.put(`/roles/${id}/permissions`, data)
}

export const menuAPI = {
  getMenus: () => api.get('/menus'),
  getMenuTree: () => api.get('/menus/tree'),
  getNonSystemMenuTree: () => api.get('/menus/non-system-tree'),
  getMenu: (id) => api.get(`/menus/${id}`),
  createMenu: (data) => api.post('/menus', data),
  updateMenu: (id, data) => api.put(`/menus/${id}`, data),
  deleteMenu: (id) => api.delete(`/menus/${id}`)
}

export const permissionAPI = {
  getPermissions: () => api.get('/permissions'),
  getPermission: (id) => api.get(`/permissions/${id}`),
  createPermission: (data) => api.post('/permissions', data),
  updatePermission: (id, data) => api.put(`/permissions/${id}`, data),
  deletePermission: (id) => api.delete(`/permissions/${id}`)
}

export const apiPermissionAPI = {
  getApiPermissions: () => api.get('/api-permissions'),
  getApiPermission: (id) => api.get(`/api-permissions/${id}`),
  createApiPermission: (data) => api.post('/api-permissions', data),
  updateApiPermission: (id, data) => api.put(`/api-permissions/${id}`, data),
  deleteApiPermission: (id) => api.delete(`/api-permissions/${id}`),
  getApiPermissionsForRole: (roleUUID) => api.get(`/api-permissions/roles/${roleUUID}`),
  addApiPermissionToRole: (roleUUID, data) => api.post(`/api-permissions/roles/${roleUUID}`, data),
  removeApiPermissionFromRole: (roleUUID, data) => api.delete(`/api-permissions/roles/${roleUUID}`, { data })
}

export const applicationAPI = {
  getApplications: () => api.get('/applications'),
  getApplication: (id) => api.get(`/applications/${id}`),
  createApplication: (data) => api.post('/applications', data),
  updateApplication: (id, data) => api.put(`/applications/${id}`, data),
  deleteApplication: (id) => api.delete(`/applications/${id}`),
  getApplicationByCode: (code) => api.get(`/applications/${code}`)
}

export default api