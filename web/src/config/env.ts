// 环境变量配置
export const ENV_CONFIG = {
  // API基础URL - 根据环境自动选择
  API_BASE_URL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  
  // 是否为生产环境
  IS_PRODUCTION: import.meta.env.MODE === 'production',
  
  // 是否为开发环境
  IS_DEVELOPMENT: import.meta.env.MODE === 'development',
  
  // Cloudflare Pages相关配置
  CF_PAGES_URL: import.meta.env.VITE_CF_PAGES_URL || '',
  
  // 后端API域名（VPS）
  BACKEND_DOMAIN: import.meta.env.VITE_BACKEND_DOMAIN || '',
  
  // 调试模式
  DEBUG: import.meta.env.VITE_DEBUG === 'true',
} as const

// API请求配置
export const API_CONFIG = {
  baseURL: ENV_CONFIG.API_BASE_URL,
  timeout: 30000,
  withCredentials: true, // 支持跨域cookie
} as const

export default ENV_CONFIG