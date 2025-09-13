import React, { useState, useEffect } from 'react'
import { 
  Typography, 
  Card, 
  Form, 
  Input, 
  Button, 
  Result, 
  Progress,
  Space
} from 'antd'
import { LockOutlined, CheckOutlined, CloseOutlined } from '@ant-design/icons'
import { Link, useSearchParams, useNavigate } from 'react-router-dom'

const { Title, Text } = Typography

interface ResetPasswordForm {
  password: string
  confirmPassword: string
}

interface PasswordStrength {
  score: number
  label: string
  color: string
}

export default function ResetPassword() {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [resetSuccess, setResetSuccess] = useState(false)
  const [tokenValid, setTokenValid] = useState(true)
  const [passwordStrength, setPasswordStrength] = useState<PasswordStrength>({ score: 0, label: '', color: '' })
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()

  const token = searchParams.get('token')
  const email = searchParams.get('email')

  useEffect(() => {
    // 验证重置令牌
    if (!token || !email) {
      setTokenValid(false)
      return
    }

    // TODO: 调用API验证令牌有效性
    // validateResetToken(token, email)
  }, [token, email])

  const calculatePasswordStrength = (password: string): PasswordStrength => {
    let score = 0
    let label = '非常弱'
    let color = '#ff4d4f'

    if (password.length >= 8) score += 1
    if (/[a-z]/.test(password)) score += 1
    if (/[A-Z]/.test(password)) score += 1
    if (/[0-9]/.test(password)) score += 1
    if (/[^A-Za-z0-9]/.test(password)) score += 1

    switch (score) {
      case 0:
      case 1:
        label = '非常弱'
        color = '#ff4d4f'
        break
      case 2:
        label = '弱'
        color = '#ff7a45'
        break
      case 3:
        label = '中等'
        color = '#ffa940'
        break
      case 4:
        label = '强'
        color = '#52c41a'
        break
      case 5:
        label = '非常强'
        color = '#389e0d'
        break
    }

    return { score: (score / 5) * 100, label, color }
  }

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const password = e.target.value
    setPasswordStrength(calculatePasswordStrength(password))
  }

  const getPasswordRequirements = (password: string) => {
    const requirements = [
      { label: '至少8个字符', met: password.length >= 8 },
      { label: '包含小写字母', met: /[a-z]/.test(password) },
      { label: '包含大写字母', met: /[A-Z]/.test(password) },
      { label: '包含数字', met: /[0-9]/.test(password) },
      { label: '包含特殊字符', met: /[^A-Za-z0-9]/.test(password) },
    ]

    return requirements
  }

  const handleSubmit = async (values: ResetPasswordForm) => {
    setLoading(true)
    try {
      // TODO: 调用API重置密码
      await new Promise(resolve => setTimeout(resolve, 2000)) // 模拟API调用
      setResetSuccess(true)
    } catch (error) {
      console.error('Reset password failed:', error)
    } finally {
      setLoading(false)
    }
  }

  if (!tokenValid) {
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
            maxWidth: 450,
            boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1), 0 1px 3px rgba(0, 0, 0, 0.08)',
            borderRadius: 12
          }}
          bodyStyle={{ padding: '32px' }}
        >
          <Result
            status="error"
            title="链接无效或已过期"
            subTitle="重置密码链接无效或已过期，请重新申请密码重置。"
            extra={
              <Space>
                <Link to="/forgot-password">
                  <Button type="primary">重新申请</Button>
                </Link>
                <Link to="/login">
                  <Button>返回登录</Button>
                </Link>
              </Space>
            }
          />
        </Card>
      </div>
    )
  }

  if (resetSuccess) {
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
            maxWidth: 450,
            boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1), 0 1px 3px rgba(0, 0, 0, 0.08)',
            borderRadius: 12
          }}
          bodyStyle={{ padding: '32px' }}
        >
          <Result
            status="success"
            title="密码重置成功"
            subTitle="您的密码已成功重置，现在可以使用新密码登录。"
            extra={
              <Button 
                type="primary" 
                onClick={() => navigate('/login')}
                style={{ 
                  height: 44,
                  borderRadius: 6,
                  fontSize: 16,
                  fontWeight: 500
                }}
              >
                立即登录
              </Button>
            }
          />
        </Card>
      </div>
    )
  }

  const watchedPassword = Form.useWatch('password', form) || ''
  const requirements = getPasswordRequirements(watchedPassword)

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
          maxWidth: 450,
          boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1), 0 1px 3px rgba(0, 0, 0, 0.08)',
          borderRadius: 12
        }}
        bodyStyle={{ padding: '32px' }}
      >
        <div style={{ textAlign: 'center', marginBottom: 32 }}>
          <Title level={2} style={{ margin: 0, color: '#1f2937' }}>
            重置密码
          </Title>
          <Text type="secondary" style={{ fontSize: 16 }}>
            为 {email} 设置新密码
          </Text>
        </div>

        <Form
          form={form}
          name="reset-password"
          onFinish={handleSubmit}
          autoComplete="off"
          layout="vertical"
          size="large"
        >
          <Form.Item
            name="password"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码!' },
              { min: 8, message: '密码至少8位字符!' },
              { 
                validator: (_, value) => {
                  if (!value) return Promise.resolve()
                  const strength = calculatePasswordStrength(value)
                  if (strength.score < 40) {
                    return Promise.reject(new Error('密码强度太弱，请设置更复杂的密码'))
                  }
                  return Promise.resolve()
                }
              }
            ]}
          >
            <Input.Password
              prefix={<LockOutlined style={{ color: '#9ca3af' }} />}
              placeholder="请输入新密码"
              autoComplete="new-password"
              onChange={handlePasswordChange}
            />
          </Form.Item>

          {watchedPassword && (
            <div style={{ marginBottom: 24 }}>
              <div style={{ marginBottom: 8 }}>
                <Text style={{ fontSize: 12, color: passwordStrength.color }}>
                  密码强度: {passwordStrength.label}
                </Text>
              </div>
              <Progress 
                percent={passwordStrength.score} 
                strokeColor={passwordStrength.color}
                showInfo={false}
                size="small"
              />
              
              <div style={{ marginTop: 12 }}>
                <Space direction="vertical" size={4} style={{ width: '100%' }}>
                  {requirements.map((req, index) => (
                    <div key={index} style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      {req.met ? (
                        <CheckOutlined style={{ color: '#52c41a', fontSize: 12 }} />
                      ) : (
                        <CloseOutlined style={{ color: '#ff4d4f', fontSize: 12 }} />
                      )}
                      <Text 
                        style={{ 
                          fontSize: 12, 
                          color: req.met ? '#52c41a' : '#9ca3af' 
                        }}
                      >
                        {req.label}
                      </Text>
                    </div>
                  ))}
                </Space>
              </div>
            </div>
          )}

          <Form.Item
            name="confirmPassword"
            label="确认新密码"
            dependencies={['password']}
            rules={[
              { required: true, message: '请确认新密码!' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('password') === value) {
                    return Promise.resolve()
                  }
                  return Promise.reject(new Error('两次输入的密码不一致!'))
                },
              }),
            ]}
          >
            <Input.Password
              prefix={<LockOutlined style={{ color: '#9ca3af' }} />}
              placeholder="请再次输入新密码"
              autoComplete="new-password"
            />
          </Form.Item>

          <Form.Item>
            <Button 
              type="primary" 
              htmlType="submit" 
              loading={loading}
              style={{ 
                width: '100%',
                height: 44,
                borderRadius: 6,
                fontSize: 16,
                fontWeight: 500
              }}
            >
              重置密码
            </Button>
          </Form.Item>

          <div style={{ textAlign: 'center', marginTop: 16 }}>
            <Link to="/login">
              <Button 
                type="text"
                style={{ padding: 0 }}
              >
                返回登录
              </Button>
            </Link>
          </div>
        </Form>
      </Card>
    </div>
  )
}