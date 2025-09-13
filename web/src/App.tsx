import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'
import UserLayout from '@/components/userLayout'
import ProtectedRoute from '@/components/protectedRoute'
import Login from '@/pages/login'
import Register from '@/pages/register'
import ResetPassword from '@/pages/reset-password'
import ForgotPassword from '@/pages/forgot-password'
import Dashboard from '@/pages/dashboard'
import DNSRecords from '@/pages/dns-records'
import Profile from '@/pages/profile'
import AdminLayout from '@/components/adminLayout'
import AdminDashboard from '@/pages/admin/dashboard'
import AdminUsers from '@/pages/admin/users'
import AdminDomains from '@/pages/admin/domains'
import AdminProviders from '@/pages/admin/providers'
import AdminSMTPConfigs from '@/pages/admin/smtp-configs'
import { useEffect } from 'react'
import { Spin } from 'antd'

function App() {
  const { user, isInitialized, initAuth } = useAuthStore()

  useEffect(() => {
    if (!isInitialized) {
      initAuth()
    }
  }, [initAuth, isInitialized])

  // 显示加载状态直到认证初始化完成
  if (!isInitialized) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
      }}>
        <Spin size="large" />
      </div>
    )
  }

  return (
    <Routes>
      {/* 公共路由 */}
      <Route path="/login" element={user ? <Navigate to="/" replace /> : <Login />} />
      <Route path="/register" element={user ? <Navigate to="/" replace /> : <Register />} />
      <Route path="/forgot-password" element={<ForgotPassword />} />
      <Route path="/reset-password" element={<ResetPassword />} />
      
      {/* 用户路由 */}
      <Route path="/" element={
        <ProtectedRoute>
          <UserLayout />
        </ProtectedRoute>
      }>
        <Route index element={<Dashboard />} />
        <Route path="dns-records" element={<DNSRecords />} />
        <Route path="profile" element={<Profile />} />
      </Route>

      {/* 管理员路由 */}
      <Route path="/admin" element={
        <ProtectedRoute adminOnly>
          <AdminLayout />
        </ProtectedRoute>
      }>
        <Route index element={<AdminDashboard />} />
        <Route path="users" element={<AdminUsers />} />
        <Route path="domains" element={<AdminDomains />} />
        <Route path="providers" element={<AdminProviders />} />
        <Route path="smtp-configs" element={<AdminSMTPConfigs />} />
      </Route>

      {/* 404 */}
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App