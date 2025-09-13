import React, { useState } from 'react'
import { Form, Input, Button, Card, message, Space, Divider, Typography, Checkbox } from 'antd'
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons'
import { useAuthStore } from '@/stores/auth-store'
import { useNavigate, Link } from 'react-router-dom'

const { Title, Text } = Typography

export default function Login() {
  const { login, isLoading } = useAuthStore()
  const navigate = useNavigate()
  const [form] = Form.useForm()
  const [rememberMe, setRememberMe] = useState(false)

  const onFinish = async (values: { email: string; password: string }) => {
    try {
      await login(values)
      
      // 如果选择记住我，设置更长的本地存储时间
      if (rememberMe) {
        localStorage.setItem('remember_me', 'true')
      }
      
      navigate('/')
    } catch (error) {
      // 错误已在 store 中处理
    }
  }

  const handleForgotPassword = () => {
    navigate('/forgot-password')
  }

  return (
    <div style={{ 
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center', 
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      padding: '20px'
    }}>
      <Card 
        style={{ 
          width: '100%',
          maxWidth: 400,
          boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1), 0 1px 3px rgba(0, 0, 0, 0.08)',
          borderRadius: 12
        }}
        bodyStyle={{ padding: '32px' }}
      >
        <div style={{ textAlign: 'center', marginBottom: 32 }}>
          <Title level={2} style={{ margin: 0, color: '#1f2937' }}>
            Domain MAX
          </Title>
          <Text type="secondary" style={{ fontSize: 16 }}>
            二级域名分发管理系统
          </Text>
        </div>

        <Form
          form={form}
          name="login"
          onFinish={onFinish}
          autoComplete="off"
          layout="vertical"
          size="large"
        >
          <Form.Item
            name="email"
            label="邮箱地址"
            rules={[
              { required: true, message: '请输入邮箱地址!' },
              { type: 'email', message: '请输入有效的邮箱地址!' }
            ]}
          >
            <Input 
              prefix={<MailOutlined style={{ color: '#9ca3af' }} />} 
              placeholder="请输入邮箱地址"
              autoComplete="email"
            />
          </Form.Item>

          <Form.Item
            name="password"
            label="密码"
            rules={[
              { required: true, message: '请输入密码!' },
              { min: 6, message: '密码至少6位字符!' }
            ]}
          >
            <Input.Password
              prefix={<LockOutlined style={{ color: '#9ca3af' }} />}
              placeholder="请输入密码"
              autoComplete="current-password"
            />
          </Form.Item>

          <Form.Item>
            <div style={{ 
              display: 'flex', 
              justifyContent: 'space-between', 
              alignItems: 'center',
              marginBottom: 16
            }}>
              <Checkbox 
                checked={rememberMe}
                onChange={(e) => setRememberMe(e.target.checked)}
              >
                记住我
              </Checkbox>
              
              <Button 
                type="link" 
                onClick={handleForgotPassword}
                style={{ padding: 0 }}
              >
                忘记密码？
              </Button>
            </div>

            <Button 
              type="primary" 
              htmlType="submit" 
              loading={isLoading}
              style={{ 
                width: '100%',
                height: 44,
                borderRadius: 6,
                fontSize: 16,
                fontWeight: 500
              }}
            >
              登录
            </Button>
          </Form.Item>

          <Divider plain>
            <Text type="secondary">还没有账号？</Text>
          </Divider>

          <div style={{ textAlign: 'center' }}>
            <Link to="/register">
              <Button 
                type="default"
                style={{ 
                  width: '100%',
                  height: 44,
                  borderRadius: 6,
                  fontSize: 16
                }}
              >
                立即注册
              </Button>
            </Link>
          </div>
        </Form>
      </Card>
    </div>
  )
}