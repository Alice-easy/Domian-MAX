// 环境变量配置
export const config = {
  // API基础URL - 根据环境自动选择
  API_BASE_URL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  
  // 是否为生产环境
  IS_PRODUCTION: import.meta.env.MODE === 'production',
  
  // Cloudflare Pages相关配置
  CF_PAGES_URL: import.meta.env.VITE_CF_PAGES_URL || '',
  
  // 后端API域名（VPS）
  BACKEND_DOMAIN: import.meta.env.VITE_BACKEND_DOMAIN || '',
}

// API请求配置
export const apiConfig = {
  baseURL: config.API_BASE_URL,
  timeout: 10000,
  withCredentials: true, // 支持跨域cookie
}

export default config