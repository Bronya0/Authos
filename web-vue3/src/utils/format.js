export const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

export const formatStatus = (status) => {
  return status === 1 ? '启用' : '禁用'
}

export const getStatusType = (status) => {
  return status === 1 ? 'success' : 'error'
}

export const getRoleTagType = (index) => {
  const types = ['info', 'success', 'warning', 'error', 'default', 'primary']
  return types[index % types.length]
}

export const getRoleColor = (index) => {
  const colors = [
    '#2080f0', // 蓝色
    '#18a058', // 绿色
    '#f0a020', // 橙色
    '#d03050', // 红色
    '#909399', // 灰色
    '#722ed1', // 紫色
    '#13c2c2', // 青色
    '#eb2f96'  // 粉色
  ]
  return colors[index % colors.length]
}