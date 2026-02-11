import React, { useState } from 'react'
import { Card, Form, Input, Button, Tabs, message } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { userApi } from '../api'
import { useUserStore } from '../store'

const { TabPane } = Tabs

const Login: React.FC = () => {
  const navigate = useNavigate()
  const { setUser } = useUserStore()
  const [loading, setLoading] = useState(false)
  const [activeTab, setActiveTab] = useState('login')

  const handleLogin = async (values: { username: string; password: string }) => {
    setLoading(true)
    try {
      const res: any = await userApi.login(values)
      if (res.code === 0) {
        setUser(res.data.user, res.data.token)
        message.success('登录成功')
        navigate('/')
      } else {
        message.error(res.message || '登录失败')
      }
    } catch (error) {
      message.error('网络错误')
    } finally {
      setLoading(false)
    }
  }

  const handleRegister = async (values: {
    username: string
    password: string
    nickname: string
  }) => {
    setLoading(true)
    try {
      const res: any = await userApi.register(values)
      if (res.code === 0) {
        message.success('注册成功，请登录')
        setActiveTab('login')
      } else {
        message.error(res.message || '注册失败')
      }
    } catch (error) {
      message.error('网络错误')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        background: '#f0f2f5',
      }}
    >
      <Card style={{ width: 400, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <h2 style={{ textAlign: 'center', marginBottom: 24 }}>Bot Chat</h2>

        <Tabs activeKey={activeTab} onChange={setActiveTab} centered>
          <TabPane tab="登录" key="login">
            <Form onFinish={handleLogin}>
              <Form.Item
                name="username"
                rules={[{ required: true, message: '请输入用户名' }]}
              >
                <Input
                  prefix={<UserOutlined />}
                  placeholder="用户名"
                  size="large"
                />
              </Form.Item>

              <Form.Item
                name="password"
                rules={[{ required: true, message: '请输入密码' }]}
              >
                <Input.Password
                  prefix={<LockOutlined />}
                  placeholder="密码"
                  size="large"
                />
              </Form.Item>

              <Form.Item>
                <Button
                  type="primary"
                  htmlType="submit"
                  size="large"
                  block
                  loading={loading}
                >
                  登录
                </Button>
              </Form.Item>
            </Form>
          </TabPane>

          <TabPane tab="注册" key="register">
            <Form onFinish={handleRegister}>
              <Form.Item
                name="username"
                rules={[
                  { required: true, message: '请输入用户名' },
                  { min: 3, message: '用户名至少3个字符' },
                  { max: 20, message: '用户名最多20个字符' },
                ]}
              >
                <Input
                  prefix={<UserOutlined />}
                  placeholder="用户名"
                  size="large"
                />
              </Form.Item>

              <Form.Item
                name="nickname"
                rules={[{ required: true, message: '请输入昵称' }]}
              >
                <Input placeholder="昵称" size="large" />
              </Form.Item>

              <Form.Item
                name="password"
                rules={[
                  { required: true, message: '请输入密码' },
                  { min: 6, message: '密码至少6个字符' },
                ]}
              >
                <Input.Password
                  prefix={<LockOutlined />}
                  placeholder="密码"
                  size="large"
                />
              </Form.Item>

              <Form.Item
                name="confirmPassword"
                dependencies={['password']}
                rules={[
                  { required: true, message: '请确认密码' },
                  ({ getFieldValue }) => ({
                    validator(_, value) {
                      if (!value || getFieldValue('password') === value) {
                        return Promise.resolve()
                      }
                      return Promise.reject(new Error('两次输入的密码不一致'))
                    },
                  }),
                ]}
              >
                <Input.Password
                  prefix={<LockOutlined />}
                  placeholder="确认密码"
                  size="large"
                />
              </Form.Item>

              <Form.Item>
                <Button
                  type="primary"
                  htmlType="submit"
                  size="large"
                  block
                  loading={loading}
                >
                  注册
                </Button>
              </Form.Item>
            </Form>
          </TabPane>
        </Tabs>
      </Card>
    </div>
  )
}

export default Login
