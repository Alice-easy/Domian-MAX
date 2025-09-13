import React from 'react'
import { Button } from 'antd'
import { useNavigate } from 'react-router-dom'
import { 
  DashboardOutlined, 
  UserOutlined, 
  DatabaseOutlined,
  SettingOutlined,
  MailOutlined
} from '@ant-design/icons'
import { BaseLayout } from './common/BaseLayout'
import { useAuthStore } from '@/stores/authStore'

const AdminLayout: React.FC = () => {
  const navigate = useNavigate()
  const { user } = useAuthStore()

  const menuItems = [
    {
      key: '/admin',
      icon: <DashboardOutlined />,
      label: '管理仪表板',
    },
    {
      key: '/admin/users',
      icon: <UserOutlined />,
      label: '用户管理',
    },
    {
      key: '/admin/domains',
      icon: <DatabaseOutlined />,
      label: '域名管理',
    },
    {
      key: '/admin/providers',
      icon: <SettingOutlined />,
      label: 'DNS 提供商',
    },
    {
      key: '/admin/smtp-configs',
      icon: <MailOutlined />,
      label: 'SMTP 配置',
    },
  ]

  const headerExtra = (
    <Button 
      type="link" 
      onClick={() => navigate('/')}
      style={{ color: '#666' }}
    >
      返回用户界面
    </Button>
  )

  return (
    <BaseLayout
      title="Admin Panel"
      menuItems={menuItems}
      showMobileDrawer={false}
      headerExtra={headerExtra}
    />
  )
}

export default AdminLayout