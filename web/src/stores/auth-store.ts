import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { authAPI, type User, type LoginRequest, type RegisterRequest } from '@/services/api'
import { message } from 'antd'

interface AuthState {
  user: User | null
  token: string | null
  refreshToken: string | null
  isLoading: boolean
  isInitialized: boolean
  
  // Actions
  login: (data: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
  logout: () => void
  initAuth: () => Promise<void>
  refreshAccessToken: () => Promise<void>
  updateProfile: (data: Partial<User>) => void
  changePassword: (oldPassword: string, newPassword: string) => Promise<void>
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      refreshToken: null,
      isLoading: false,
      isInitialized: false,
      
      login: async (data: LoginRequest) => {
        set({ isLoading: true })
        try {
          const response = await authAPI.login(data)
          const { token, refresh_token, user } = response.data
          
          // 保存到 localStorage
          localStorage.setItem('auth_token', token)
          localStorage.setItem('refresh_token', refresh_token)
          
          set({ 
            user, 
            token, 
            refreshToken: refresh_token, 
            isLoading: false 
          })
          
          message.success('登录成功')
        } catch (error: any) {
          set({ isLoading: false })
          const errorMessage = error.response?.data?.message || '登录失败'
          message.error(errorMessage)
          throw error
        }
      },
      
      register: async (data: RegisterRequest) => {
        set({ isLoading: true })
        try {
          const response = await authAPI.register(data)
          const { token, refresh_token, user } = response.data
          
          // 保存到 localStorage
          localStorage.setItem('auth_token', token)
          localStorage.setItem('refresh_token', refresh_token)
          
          set({ 
            user, 
            token, 
            refreshToken: refresh_token, 
            isLoading: false 
          })
          
          message.success('注册成功')
        } catch (error: any) {
          set({ isLoading: false })
          const errorMessage = error.response?.data?.message || '注册失败'
          message.error(errorMessage)
          throw error
        }
      },
      
      logout: () => {
        // 清除本地存储
        localStorage.removeItem('auth_token')
        localStorage.removeItem('refresh_token')
        localStorage.removeItem('user_info')
        
        set({ 
          user: null, 
          token: null, 
          refreshToken: null 
        })
        
        message.success('已退出登录')
      },
      
      initAuth: async () => {
        try {
          const token = localStorage.getItem('auth_token')
          const refreshToken = localStorage.getItem('refresh_token')
          
          if (token && refreshToken) {
            set({ token, refreshToken })
            
            // 验证 token 并获取用户信息
            const response = await authAPI.profile()
            set({ 
              user: response.data, 
              isInitialized: true 
            })
          } else {
            set({ isInitialized: true })
          }
        } catch (error: any) {
          // Token 可能已过期，尝试刷新
          const refreshToken = localStorage.getItem('refresh_token')
          if (refreshToken) {
            try {
              await get().refreshAccessToken()
              set({ isInitialized: true })
            } catch (refreshError) {
              // 刷新失败，清除状态
              get().logout()
              set({ isInitialized: true })
            }
          } else {
            set({ isInitialized: true })
          }
        }
      },
      
      refreshAccessToken: async () => {
        const { refreshToken } = get()
        if (!refreshToken) {
          throw new Error('No refresh token available')
        }
        
        try {
          const response = await authAPI.refreshToken(refreshToken)
          const { token, refresh_token, user } = response.data
          
          // 更新本地存储
          localStorage.setItem('auth_token', token)
          localStorage.setItem('refresh_token', refresh_token)
          
          set({ 
            token, 
            refreshToken: refresh_token, 
            user 
          })
        } catch (error) {
          // 刷新失败，清除状态
          get().logout()
          throw error
        }
      },
      
      updateProfile: (data: Partial<User>) => {
        const { user } = get()
        if (user) {
          set({ user: { ...user, ...data } })
        }
      },
      
      changePassword: async (oldPassword: string, newPassword: string) => {
        set({ isLoading: true })
        try {
          await authAPI.changePassword({
            old_password: oldPassword,
            new_password: newPassword
          })
          set({ isLoading: false })
          message.success('密码修改成功')
        } catch (error: any) {
          set({ isLoading: false })
          const errorMessage = error.response?.data?.message || '密码修改失败'
          message.error(errorMessage)
          throw error
        }
      }
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ 
        user: state.user,
        token: state.token,
        refreshToken: state.refreshToken
      }),
    }
  )
)