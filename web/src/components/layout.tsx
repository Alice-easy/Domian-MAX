import React, { useState } from 'react'
import { Layout as AntLayout, Menu, Avatar, Dropdown, Button, Drawer } from 'antd'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { 
  DashboardOutlined, 
  DatabaseOutlined, 
  UserOutlined, 
  LogoutOutlined,
  MenuOutlined,
  SettingOutlined
} from '@ant-design/icons'
import { useAuthStore } from '@/stores/auth-store'

const { Header, Sider, Content } = AntLayout

export default function Layout() {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()
  const [mobileDrawerOpen, setMobileDrawerOpen] = useState(false)

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

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人资料',
      onClick: () => navigate('/profile'),
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '设置',
      onClick: () => navigate('/settings'),
    },
    {
      key: 'divider',
      type: 'divider' as const,
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      danger: true,
      onClick: () => {
        logout()
        navigate('/login')
      },
    },
  ]

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key)
    setMobileDrawerOpen(false)
  }

  const sidebarContent = (
    <>
      <div style={{ 
        height: 64, 
        margin: 16, 
        background: 'rgba(255, 255, 255, 0.2)',
        borderRadius: 6,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        color: 'white',
        fontSize: '18px',
        fontWeight: 'bold'
      }}>
        Domain MAX
      </div>
      <Menu
        theme="dark"
        mode="inline"
        selectedKeys={[location.pathname]}
        items={menuItems}
        onClick={handleMenuClick}
      />
    </>
  )

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      {/* 桌面侧边栏 */}
      <Sider
        breakpoint="lg"
        collapsedWidth="0"
        className="desktop-sider"
        style={{ 
          display: 'none',
        }}
        onBreakpoint={(broken: boolean) => {
          // 在大屏幕上显示侧边栏
          const siderElement = document.querySelector('.desktop-sider') as HTMLElement
          if (siderElement) {
            siderElement.style.display = broken ? 'none' : 'block'
          }
        }}
      >
        {sidebarContent}
      </Sider>

      {/* 移动端抽屉 */}
      <Drawer
        title={
          <div style={{ 
            color: 'white',
            fontSize: '18px',
            fontWeight: 'bold'
          }}>
            Domain MAX
          </div>
        }
        placement="left"
        onClose={() => setMobileDrawerOpen(false)}
        open={mobileDrawerOpen}
        width={250}
        bodyStyle={{ padding: 0 }}
        headerStyle={{ 
          background: '#001529',
          color: 'white',
          borderBottom: '1px solid #303030'
        }}
      >
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={handleMenuClick}
          style={{ border: 'none' }}
        />
      </Drawer>

      <AntLayout>
        <Header style={{ 
          padding: '0 24px', 
          background: '#fff',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          boxShadow: '0 2px 8px rgba(0,0,0,0.06)'
        }}>
          {/* 移动端菜单按钮 */}
          <Button
            type="text"
            icon={<MenuOutlined />}
            onClick={() => setMobileDrawerOpen(true)}
            className="mobile-menu-button"
            style={{ 
              display: 'none',
              fontSize: '18px'
            }}
          />

          {/* 用户菜单 */}
          <div style={{ marginLeft: 'auto' }}>
            <Dropdown 
              menu={{ items: userMenuItems }} 
              placement="bottomRight"
              trigger={['click']}
            >
              <Button 
                type="text" 
                style={{ 
                  display: 'flex', 
                  alignItems: 'center', 
                  gap: 8,
                  padding: '4px 8px',
                  height: 'auto'
                }}
              >
                <Avatar 
                  icon={<UserOutlined />} 
                  size="small"
                  style={{ backgroundColor: '#1890ff' }}
                />
                <span style={{ 
                  maxWidth: 120, 
                  overflow: 'hidden', 
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap'
                }}>
                  {user?.username || user?.email}
                </span>
              </Button>
            </Dropdown>
          </div>
        </Header>

        <Content style={{ margin: '24px 16px 0' }}>
          <div style={{ 
            padding: 24, 
            minHeight: 360, 
            background: '#fff',
            borderRadius: 8,
            boxShadow: '0 2px 8px rgba(0,0,0,0.06)'
          }}>
            <Outlet />
          </div>
        </Content>
      </AntLayout>

      <style>{`
        @media (min-width: 992px) {
          .desktop-sider {
            display: block !important;
          }
          .mobile-menu-button {
            display: none !important;
          }
        }
        
        @media (max-width: 991px) {
          .mobile-menu-button {
            display: inline-flex !important;
          }
        }
      `}</style>
    </AntLayout>
  )
}