import React from 'react'
import { Navigate } from 'react-router-dom'
import { useAuthStore } from '@/stores/auth-store'
import { Spin, Result } from 'antd'

interface ProtectedRouteProps {
  children: React.ReactNode
  adminOnly?: boolean
}

export default function ProtectedRoute({ children, adminOnly = false }: ProtectedRouteProps) {
  const { user, isInitialized, isLoading } = useAuthStore()

  // 如果还在初始化认证状态，显示加载
  if (!isInitialized || isLoading) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        background: '#f0f2f5'
      }}>
        <Spin size="large" />
      </div>
    )
  }

  // 如果用户未登录，重定向到登录页
  if (!user) {
    return <Navigate to="/login" replace />
  }

  // 如果需要管理员权限但用户不是管理员
  if (adminOnly && user.role !== 'admin') {
    return (
      <Result
        status="403"
        title="403"
        subTitle="抱歉，您没有权限访问此页面。"
        extra={
          <Navigate to="/" replace />
        }
      />
    )
  }

  return <>{children}</>
}