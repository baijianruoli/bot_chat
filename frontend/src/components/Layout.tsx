import React from 'react'
import { Layout, Menu, Button, Avatar, Dropdown, Space, Typography } from 'antd'
import {
  UserOutlined,
  TeamOutlined,
  SettingOutlined,
  LogoutOutlined,
  MessageOutlined,
} from '@ant-design/icons'
import { useUserStore } from '../store'
import { useNavigate, useLocation } from 'react-router-dom'

const { Header, Sider, Content } = Layout
const { Title } = Typography

interface AppLayoutProps {
  children: React.ReactNode
}

const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useUserStore()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人资料',
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '设置',
    },
    {
      type: 'divider' as const,
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ]

  const sidebarItems = [
    {
      key: '/',
      icon: <MessageOutlined />,
      label: '聊天室',
    },
    {
      key: '/rooms',
      icon: <TeamOutlined />,
      label: '房间列表',
    },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          padding: '0 24px',
          background: '#1677ff',
        }}
      >
        <Space>
          <MessageOutlined style={{ fontSize: 24, color: 'white' }} />
          <Title level={4} style={{ color: 'white', margin: 0 }}>
            Bot Chat
          </Title>
        </Space>

        <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
          <Space style={{ cursor: 'pointer', color: 'white' }}>
            <Avatar icon={<UserOutlined />} src={user?.avatar} />
            <span>{user?.nickname || user?.username}</span>
          </Space>
        </Dropdown>
      </Header>

      <Layout>
        <Sider width={200} theme="light">
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            style={{ height: '100%', borderRight: 0 }}
            items={sidebarItems}
            onClick={({ key }) => navigate(key)}
          />
        </Sider>

        <Content style={{ padding: 24, background: '#f5f5f5' }}>
          {children}
        </Content>
      </Layout>
    </Layout>
  )
}

export default AppLayout
