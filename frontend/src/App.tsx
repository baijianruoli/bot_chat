import React from 'react'
import { Layout, Typography } from 'antd'
import './App.css'

const { Header, Content } = Layout
const { Title } = Typography

function App() {
  return (
    <Layout className="app-layout">
      <Header className="app-header">
        <Title level={3} style={{ color: 'white', margin: 0 }}>
          🤖 Bot Chat
        </Title>
      </Header>
      <Content className="app-content">
        <div className="welcome-card">
          <Title level={2}>欢迎使用 Bot Chat</Title>
          <p>一个基于 Golang + Kitex + React 的实时聊天室</p>
          <p>项目正在开发中，敬请期待...</p>
        </div>
      </Content>
    </Layout>
  )
}

export default App
