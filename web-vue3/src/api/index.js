import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    // 从localStorage获取应用ID并添加到请求头
    const appId = localStorage.getItem('appId')
    if (appId) {
      config.headers['X-App-ID'] = appId
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
      localStorage.removeItem('token')
      localStorage.removeItem('currentUser')
      localStorage.removeItem('appId')
      localStorage.removeItem('appInfo')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const authAPI = {
  login: (data) => api.post('/login', data),
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
  getApplications: () => api.get('/api/applications'),
  getApplication: (id) => api.get(`/api/applications/${id}`),
  createApplication: (data) => api.post('/api/applications', data),
  updateApplication: (id, data) => api.put(`/api/applications/${id}`, data),
  deleteApplication: (id) => api.delete(`/api/applications/${id}`),
  getApplicationByCode: (code) => api.get(`/api/applications/${code}`)
}

export default api