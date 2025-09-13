import React from 'react'
import { 
  DashboardOutlined, 
  DatabaseOutlined, 
  UserOutlined
} from '@ant-design/icons'
import { BaseLayout } from './common/BaseLayout'

const UserLayout: React.FC = () => {
  const menuItems = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: '仪表板',
    },
    {
      key: '/dns-records',
      icon: <DatabaseOutlined />,
      label: 'DNS 记录',
    },
    {
      key: '/profile',
      icon: <UserOutlined />,
      label: '个人资料',
    },
  ]

  return (
    <BaseLayout
      title="Domain MAX"
      menuItems={menuItems}
      showMobileDrawer={true}
    />
  )
}

export default UserLayout