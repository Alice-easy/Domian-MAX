import React, { useState, ReactNode } from 'react'
import { Layout as AntLayout, Menu, Avatar, Dropdown, Button, Drawer } from 'antd'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { 
  LogoutOutlined,
  MenuOutlined,
  SettingOutlined,
  UserOutlined
} from '@ant-design/icons'
import { useAuthStore } from '@/stores/authStore'

const { Header, Sider, Content } = AntLayout

interface MenuItem {
  key: string
  icon: ReactNode
  label: string
}

interface BaseLayoutProps {
  title: string
  menuItems: MenuItem[]
  showMobileDrawer?: boolean
  headerExtra?: ReactNode
}

export const BaseLayout: React.FC<BaseLayoutProps> = ({
  title,
  menuItems,
  showMobileDrawer = true,
  headerExtra
}) => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()
  const [mobileDrawerOpen, setMobileDrawerOpen] = useState(false)

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
      <div className="layout-logo">
        {title}
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
        style={{ display: 'none' }}
        onBreakpoint={(broken: boolean) => {
          const siderElement = document.querySelector('.desktop-sider') as HTMLElement
          if (siderElement) {
            siderElement.style.display = broken ? 'none' : 'block'
          }
        }}
      >
        {sidebarContent}
      </Sider>

      {/* 移动端抽屉 */}
      {showMobileDrawer && (
        <Drawer
          title={<div className="drawer-title">{title}</div>}
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
      )}

      <AntLayout>
        <Header className="layout-header">
          {/* 移动端菜单按钮 */}
          {showMobileDrawer && (
            <Button
              type="text"
              icon={<MenuOutlined />}
              onClick={() => setMobileDrawerOpen(true)}
              className="mobile-menu-button"
            />
          )}

          {/* 头部额外内容 */}
          {headerExtra && <div className="header-extra">{headerExtra}</div>}

          {/* 用户菜单 */}
          <div className="user-menu">
            <Dropdown 
              menu={{ items: userMenuItems }} 
              placement="bottomRight"
              trigger={['click']}
            >
              <Button type="text" className="user-menu-button">
                <Avatar 
                  icon={<UserOutlined />} 
                  size="small"
                  style={{ backgroundColor: '#1890ff' }}
                />
                <span className="username">
                  {user?.username || user?.email}
                </span>
              </Button>
            </Dropdown>
          </div>
        </Header>

        <Content style={{ margin: '24px 16px 0' }}>
          <div className="content-wrapper">
            <Outlet />
          </div>
        </Content>
      </AntLayout>

      <style>{`
        .layout-logo {
          height: 64px;
          margin: 16px;
          background: rgba(255, 255, 255, 0.2);
          border-radius: 6px;
          display: flex;
          align-items: center;
          justify-content: center;
          color: white;
          font-size: 18px;
          font-weight: bold;
        }

        .drawer-title {
          color: white;
          font-size: 18px;
          font-weight: bold;
        }

        .layout-header {
          padding: 0 24px;
          background: #fff;
          display: flex;
          justify-content: space-between;
          align-items: center;
          box-shadow: 0 2px 8px rgba(0,0,0,0.06);
        }

        .mobile-menu-button {
          display: none;
          font-size: 18px;
        }

        .header-extra {
          flex: 1;
          display: flex;
          justify-content: center;
        }

        .user-menu {
          margin-left: auto;
        }

        .user-menu-button {
          display: flex;
          align-items: center;
          gap: 8px;
          padding: 4px 8px;
          height: auto;
        }

        .username {
          max-width: 120px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .content-wrapper {
          padding: 24px;
          min-height: 360px;
          background: #fff;
          border-radius: 8px;
          box-shadow: 0 2px 8px rgba(0,0,0,0.06);
        }

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

export default BaseLayout