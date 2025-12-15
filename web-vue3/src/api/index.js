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
    if (error.response?.status === 401) {
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
  logout: () => api.post('/public/logout')  
}

export const userAPI = {
  getUsers: () => api.get('/v1/users'),
  getUser: (id) => api.get(`/v1/users/${id}`),
  createUser: (data) => api.post('/v1/users', data),
  updateUser: (id, data) => api.put(`/v1/users/${id}`, data),
  deleteUser: (id) => api.delete(`/v1/users/${id}`)
}

export const roleAPI = {
  getRoles: () => api.get('/v1/roles'),
  getRole: (id) => api.get(`/v1/roles/${id}`),
  createRole: (data) => api.post('/v1/roles', data),
  updateRole: (id, data) => api.put(`/v1/roles/${id}`, data),
  deleteRole: (id) => api.delete(`/v1/roles/${id}`),
  getRoleMenus: (id) => api.get(`/v1/roles/${id}/menus`),
  updateRoleMenus: (id, data) => api.put(`/v1/roles/${id}/menus`, data),
  getRolePermissions: (id) => api.get(`/v1/roles/${id}/permissions`),
  updateRolePermissions: (id, data) => api.put(`/v1/roles/${id}/permissions`, data)
}

export const menuAPI = {
  getMenus: () => api.get('/v1/menus'),
  getMenuTree: () => api.get('/v1/menus/tree'),
  getNonSystemMenuTree: () => api.get('/v1/menus/non-system-tree'),
  getMenu: (id) => api.get(`/v1/menus/${id}`),
  createMenu: (data) => api.post('/v1/menus', data),
  updateMenu: (id, data) => api.put(`/v1/menus/${id}`, data),
  deleteMenu: (id) => api.delete(`/v1/menus/${id}`)
}

export const permissionAPI = {
  getPermissions: () => api.get('/v1/permissions'),
  getPermission: (id) => api.get(`/v1/permissions/${id}`),
  createPermission: (data) => api.post('/v1/permissions', data),
  updatePermission: (id, data) => api.put(`/v1/permissions/${id}`, data),
  deletePermission: (id) => api.delete(`/v1/permissions/${id}`)
}

export const apiPermissionAPI = {
  getApiPermissions: () => api.get('/v1/api-permissions'),
  getApiPermission: (id) => api.get(`/v1/api-permissions/${id}`),
  createApiPermission: (data) => api.post('/v1/api-permissions', data),
  updateApiPermission: (id, data) => api.put(`/v1/api-permissions/${id}`, data),
  deleteApiPermission: (id) => api.delete(`/v1/api-permissions/${id}`),
  getApiPermissionsForRole: (roleUUID) => api.get(`/v1/api-permissions/roles/${roleUUID}`),
  addApiPermissionToRole: (roleUUID, data) => api.post(`/api-permissions/roles/${roleUUID}`, data),
  removeApiPermissionFromRole: (roleUUID, data) => api.delete(`/api-permissions/roles/${roleUUID}`, { data })
}

export const applicationAPI = {
  getApplications: () => api.get('/v1/applications'),
  getApplication: (id) => api.get(`/v1/applications/${id}`),
  createApplication: (data) => api.post('/v1/applications', data),
  updateApplication: (id, data) => api.put(`/v1/applications/${id}`, data),
  deleteApplication: (id) => api.delete(`/v1/applications/${id}`),
  getApplicationByCode: (code) => api.get(`/v1/applications/${code}`)
}

export default api