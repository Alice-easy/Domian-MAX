import React, { useState } from 'react'
import { 
  Typography, 
  Card, 
  Form, 
  Input, 
  Button, 
  Result, 
  Space,
  Divider 
} from 'antd'
import { MailOutlined, ArrowLeftOutlined, CheckCircleOutlined } from '@ant-design/icons'
import { Link } from 'react-router-dom'

const { Title, Text } = Typography

export default function ForgotPassword() {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [emailSent, setEmailSent] = useState(false)
  const [sentEmail, setSentEmail] = useState('')

  const handleSubmit = async (values: { email: string }) => {
    setLoading(true)
    try {
      // TODO: 实现发送重置邮件的API调用
      await new Promise(resolve => setTimeout(resolve, 2000)) // 模拟API调用
      setSentEmail(values.email)
      setEmailSent(true)
    } catch (error) {
      console.error('Send reset email failed:', error)
    } finally {
      setLoading(false)
    }
  }

  if (emailSent) {
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
            icon={<CheckCircleOutlined style={{ color: '#52c41a' }} />}
            title="邮件已发送"
            subTitle={
              <div>
                <Text>
                  密码重置邮件已发送到 <Text strong>{sentEmail}</Text>
                </Text>
                <br />
                <Text type="secondary" style={{ fontSize: 14 }}>
                  请检查您的邮箱（包括垃圾邮件文件夹），并点击邮件中的链接重置密码。
                </Text>
              </div>
            }
            extra={
              <Space direction="vertical" style={{ width: '100%' }}>
                <Text type="secondary" style={{ fontSize: 12, textAlign: 'center', display: 'block' }}>
                  没有收到邮件？
                </Text>
                <Space style={{ width: '100%', justifyContent: 'center' }}>
                  <Button 
                    type="link" 
                    onClick={() => {
                      setEmailSent(false)
                      form.setFieldsValue({ email: sentEmail })
                    }}
                    style={{ padding: 0 }}
                  >
                    重新发送
                  </Button>
                  <Divider type="vertical" />
                  <Link to="/login">
                    <Button type="link" style={{ padding: 0 }}>
                      返回登录
                    </Button>
                  </Link>
                </Space>
              </Space>
            }
          />
        </Card>
      </div>
    )
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
            忘记密码
          </Title>
          <Text type="secondary" style={{ fontSize: 16 }}>
            输入您的邮箱地址，我们将发送重置链接
          </Text>
        </div>

        <Form
          form={form}
          name="forgot-password"
          onFinish={handleSubmit}
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
              placeholder="请输入您的邮箱地址"
              autoComplete="email"
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
              发送重置邮件
            </Button>
          </Form.Item>

          <div style={{ textAlign: 'center', marginTop: 16 }}>
            <Link to="/login">
              <Button 
                type="text"
                icon={<ArrowLeftOutlined />}
                style={{ padding: 0 }}
              >
                返回登录
              </Button>
            </Link>
          </div>
        </Form>

        <Divider plain>
          <Text type="secondary" style={{ fontSize: 12 }}>
            记住密码了？
          </Text>
        </Divider>

        <div style={{ textAlign: 'center' }}>
          <Link to="/login">
            <Button 
              type="default"
              style={{ 
                width: '100%',
                height: 44,
                borderRadius: 6,
                fontSize: 16
              }}
            >
              立即登录
            </Button>
          </Link>
        </div>
      </Card>
    </div>
  )
}