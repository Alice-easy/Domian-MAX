import axios, { AxiosInstance, AxiosResponse } from 'axios'
import { message } from 'antd'

// API 基础配置
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

// 创建 axios 实例
const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response
  },
  (error) => {
    const { response } = error
    
    if (response) {
      switch (response.status) {
        case 401:
          localStorage.removeItem('auth_token')
          localStorage.removeItem('user_info')
          window.location.href = '/login'
          message.error('认证已过期，请重新登录')
          break
        case 403:
          message.error('权限不足')
          break
        case 404:
          message.error('请求的资源不存在')
          break
        case 500:
          message.error('服务器内部错误')
          break
        default:
          message.error(response.data?.message || '请求失败')
      }
    } else {
      message.error('网络连接失败')
    }
    
    return Promise.reject(error)
  }
)

// 类型定义
export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface User {
  id: string
  username: string
  email: string
  role: 'user' | 'admin'
  created_at: string
  updated_at: string
}

export interface LoginResponse {
  token: string
  refresh_token: string
  user: User
}

export interface DNSProvider {
  id: string
  name: string
  provider_type: string
  config: Record<string, any>
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface DNSRecord {
  id: string
  domain: string
  record_type: string
  name: string
  value: string
  ttl: number
  priority?: number
  weight?: number
  port?: number
  provider_id: string
  created_at: string
  updated_at: string
}

// API 方法
export const authAPI = {
  // 登录
  login: (data: LoginRequest): Promise<AxiosResponse<LoginResponse>> =>
    api.post('/auth/login', data),
  
  // 注册
  register: (data: RegisterRequest): Promise<AxiosResponse<LoginResponse>> =>
    api.post('/auth/register', data),
  
  // 获取当前用户信息
  profile: (): Promise<AxiosResponse<User>> =>
    api.get('/auth/profile'),
  
  // 刷新 token
  refreshToken: (refreshToken: string): Promise<AxiosResponse<LoginResponse>> =>
    api.post('/auth/refresh', { refresh_token: refreshToken }),
  
  // 修改密码
  changePassword: (data: { old_password: string; new_password: string }): Promise<AxiosResponse> =>
    api.post('/auth/change-password', data),
}

export const dnsProviderAPI = {
  // 获取DNS提供商列表
  list: (): Promise<AxiosResponse<DNSProvider[]>> =>
    api.get('/dns-providers'),
  
  // 创建DNS提供商
  create: (data: Partial<DNSProvider>): Promise<AxiosResponse<DNSProvider>> =>
    api.post('/dns-providers', data),
  
  // 更新DNS提供商
  update: (id: string, data: Partial<DNSProvider>): Promise<AxiosResponse<DNSProvider>> =>
    api.put(`/dns-providers/${id}`, data),
  
  // 删除DNS提供商
  delete: (id: string): Promise<AxiosResponse> =>
    api.delete(`/dns-providers/${id}`),
  
  // 测试DNS提供商连接
  test: (id: string): Promise<AxiosResponse<{ success: boolean; message: string }>> =>
    api.post(`/dns-providers/${id}/test`),
}

export const dnsRecordAPI = {
  // 获取DNS记录列表
  list: (params?: { domain?: string; provider_id?: string; page?: number; limit?: number }): Promise<AxiosResponse<{ records: DNSRecord[]; total: number }>> =>
    api.get('/dns-records', { params }),
  
  // 创建DNS记录
  create: (data: Partial<DNSRecord>): Promise<AxiosResponse<DNSRecord>> =>
    api.post('/dns-records', data),
  
  // 更新DNS记录
  update: (id: string, data: Partial<DNSRecord>): Promise<AxiosResponse<DNSRecord>> =>
    api.put(`/dns-records/${id}`, data),
  
  // 删除DNS记录
  delete: (id: string): Promise<AxiosResponse> =>
    api.delete(`/dns-records/${id}`),
  
  // 批量操作
  batchCreate: (records: Partial<DNSRecord>[]): Promise<AxiosResponse<DNSRecord[]>> =>
    api.post('/dns-records/batch', { records }),
  
  batchDelete: (ids: string[]): Promise<AxiosResponse> =>
    api.delete('/dns-records/batch', { data: { ids } }),
}

export const adminAPI = {
  // 用户管理
  users: {
    list: (params?: { page?: number; limit?: number; search?: string }): Promise<AxiosResponse<{ users: User[]; total: number }>> =>
      api.get('/admin/users', { params }),
    
    create: (data: RegisterRequest): Promise<AxiosResponse<User>> =>
      api.post('/admin/users', data),
    
    update: (id: string, data: Partial<User>): Promise<AxiosResponse<User>> =>
      api.put(`/admin/users/${id}`, data),
    
    delete: (id: string): Promise<AxiosResponse> =>
      api.delete(`/admin/users/${id}`),
    
    resetPassword: (id: string, newPassword: string): Promise<AxiosResponse> =>
      api.post(`/admin/users/${id}/reset-password`, { new_password: newPassword }),
  },
  
  // 系统统计
  stats: (): Promise<AxiosResponse<{
    total_users: number
    total_providers: number
    total_records: number
    active_providers: number
  }>> =>
    api.get('/admin/stats'),
}

export default api