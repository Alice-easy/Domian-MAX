// 应用常量配置
export const APP_CONFIG = {
  // 应用信息
  APP_NAME: 'Domain MAX',
  APP_VERSION: '1.0.0',
  APP_DESCRIPTION: '二级域名分发管理系统',
  
  // API配置
  API_TIMEOUT: 30000,
  
  // 分页配置
  DEFAULT_PAGE_SIZE: 20,
  PAGE_SIZE_OPTIONS: ['10', '20', '50', '100'],
  
  // 本地存储键名
  STORAGE_KEYS: {
    AUTH_TOKEN: 'auth_token',
    REFRESH_TOKEN: 'refresh_token',
    USER_INFO: 'user_info',
    REMEMBER_ME: 'remember_me',
  },
  
  // DNS记录类型
  DNS_RECORD_TYPES: [
    { value: 'A', label: 'A记录', color: 'blue' },
    { value: 'AAAA', label: 'AAAA记录', color: 'purple' },
    { value: 'CNAME', label: 'CNAME记录', color: 'orange' },
    { value: 'MX', label: 'MX记录', color: 'green' },
    { value: 'TXT', label: 'TXT记录', color: 'red' },
    { value: 'NS', label: 'NS记录', color: 'cyan' },
    { value: 'SRV', label: 'SRV记录', color: 'magenta' },
    { value: 'PTR', label: 'PTR记录', color: 'gold' },
  ],
  
  // 用户角色
  USER_ROLES: {
    USER: 'user',
    ADMIN: 'admin',
  },
  
  // 响应式断点
  BREAKPOINTS: {
    XS: 480,
    SM: 576,
    MD: 768,
    LG: 992,
    XL: 1200,
    XXL: 1600,
  },
} as const

// DNS记录类型颜色映射
export const getDNSRecordTypeColor = (type: string): string => {
  const recordType = APP_CONFIG.DNS_RECORD_TYPES.find(t => t.value === type)
  return recordType?.color || 'default'
}

// 获取DNS记录类型标签
export const getDNSRecordTypeLabel = (type: string): string => {
  const recordType = APP_CONFIG.DNS_RECORD_TYPES.find(t => t.value === type)
  return recordType?.label || type
}

export default APP_CONFIG