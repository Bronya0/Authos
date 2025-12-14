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
  const types = ['default', 'info', 'success', 'warning', 'error']
  return types[index % types.length]
}