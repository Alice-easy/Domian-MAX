import React, { useState } from 'react'
import {
  Typography,
  Card,
  Form,
  Input,
  Button,
  Avatar,
  Space,
  Divider,
  Row,
  Col,
  Tag,
  Alert,
  Descriptions,
  Modal,
  Progress
} from 'antd'
import {
  UserOutlined,
  MailOutlined,
  LockOutlined,
  EditOutlined,
  SaveOutlined,
  KeyOutlined,
  InfoCircleOutlined,
  CheckCircleOutlined
} from '@ant-design/icons'
import { useAuthStore } from '@/stores/auth-store'
import dayjs from 'dayjs'

const { Title, Text } = Typography

interface ChangePasswordForm {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

interface ProfileForm {
  username: string
  email: string
}

export default function Profile() {
  const { user, updateProfile, changePassword, isLoading } = useAuthStore()
  const [editMode, setEditMode] = useState(false)
  const [passwordModalVisible, setPasswordModalVisible] = useState(false)
  const [profileForm] = Form.useForm<ProfileForm>()
  const [passwordForm] = Form.useForm<ChangePasswordForm>()

  if (!user) {
    return null
  }

  const handleProfileUpdate = async (values: ProfileForm) => {
    try {
      updateProfile(values)
      setEditMode(false)
      // 在实际应用中，这里应该调用API更新用户信息
    } catch (error) {
      console.error('Profile update failed:', error)
    }
  }

  const handlePasswordChange = async (values: ChangePasswordForm) => {
    try {
      await changePassword(values.oldPassword, values.newPassword)
      setPasswordModalVisible(false)
      passwordForm.resetFields()
    } catch (error) {
      // Error handling is done in the store
    }
  }

  const getPasswordStrength = (password: string) => {
    let score = 0
    if (password.length >= 8) score += 25
    if (/[a-z]/.test(password)) score += 25
    if (/[A-Z]/.test(password)) score += 25
    if (/[0-9]/.test(password) && /[^A-Za-z0-9]/.test(password)) score += 25
    return score
  }

  const watchedPassword = Form.useWatch('newPassword', passwordForm) || ''
  const passwordStrength = getPasswordStrength(watchedPassword)

  const getPasswordStrengthColor = (strength: number) => {
    if (strength < 25) return '#ff4d4f'
    if (strength < 50) return '#ff7a45'
    if (strength < 75) return '#ffa940'
    return '#52c41a'
  }

  const getPasswordStrengthText = (strength: number) => {
    if (strength < 25) return '弱'
    if (strength < 50) return '一般'
    if (strength < 75) return '良好'
    return '强'
  }

  return (
    <div>
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        marginBottom: 24 
      }}>
        <div>
          <Title level={2} style={{ margin: 0 }}>
            个人资料
          </Title>
          <Text type="secondary">
            管理您的账户信息和安全设置
          </Text>
        </div>
      </div>

      <Row gutter={[24, 24]}>
        {/* 基本信息 */}
        <Col xs={24} lg={16}>
          <Card
            title={
              <Space>
                <UserOutlined />
                <span>基本信息</span>
              </Space>
            }
            extra={
              !editMode ? (
                <Button
                  type="text"
                  icon={<EditOutlined />}
                  onClick={() => {
                    setEditMode(true)
                    profileForm.setFieldsValue({
                      username: user.username,
                      email: user.email
                    })
                  }}
                >
                  编辑
                </Button>
              ) : (
                <Space>
                  <Button 
                    size="small"
                    onClick={() => setEditMode(false)}
                  >
                    取消
                  </Button>
                  <Button
                    type="primary"
                    size="small"
                    icon={<SaveOutlined />}
                    onClick={() => profileForm.submit()}
                    loading={isLoading}
                  >
                    保存
                  </Button>
                </Space>
              )
            }
          >
            {!editMode ? (
              <Descriptions column={1} size="small">
                <Descriptions.Item label="用户名">
                  <Text strong>{user.username}</Text>
                </Descriptions.Item>
                <Descriptions.Item label="邮箱地址">
                  <Space>
                    <Text>{user.email}</Text>
                    <Tag color="green" icon={<CheckCircleOutlined />}>
                      已验证
                    </Tag>
                  </Space>
                </Descriptions.Item>
                <Descriptions.Item label="用户角色">
                  <Tag color={user.role === 'admin' ? 'red' : 'blue'}>
                    {user.role === 'admin' ? '管理员' : '普通用户'}
                  </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="注册时间">
                  <Text type="secondary">
                    {dayjs(user.created_at).format('YYYY年MM月DD日 HH:mm')}
                  </Text>
                </Descriptions.Item>
                <Descriptions.Item label="最后更新">
                  <Text type="secondary">
                    {dayjs(user.updated_at).format('YYYY年MM月DD日 HH:mm')}
                  </Text>
                </Descriptions.Item>
              </Descriptions>
            ) : (
              <Form
                form={profileForm}
                layout="vertical"
                onFinish={handleProfileUpdate}
              >
                <Form.Item
                  name="username"
                  label="用户名"
                  rules={[
                    { required: true, message: '请输入用户名' },
                    { min: 3, message: '用户名至少3个字符' },
                    { max: 50, message: '用户名最多50个字符' },
                    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线' }
                  ]}
                >
                  <Input 
                    prefix={<UserOutlined />} 
                    placeholder="请输入用户名"
                  />
                </Form.Item>

                <Form.Item
                  name="email"
                  label="邮箱地址"
                  rules={[
                    { required: true, message: '请输入邮箱地址' },
                    { type: 'email', message: '请输入有效的邮箱地址' }
                  ]}
                >
                  <Input 
                    prefix={<MailOutlined />} 
                    placeholder="请输入邮箱地址"
                  />
                </Form.Item>
              </Form>
            )}
          </Card>
        </Col>

        {/* 侧边栏信息 */}
        <Col xs={24} lg={8}>
          <Space direction="vertical" style={{ width: '100%' }} size="large">
            {/* 头像卡片 */}
            <Card size="small">
              <div style={{ textAlign: 'center' }}>
                <Avatar 
                  size={80} 
                  icon={<UserOutlined />}
                  style={{ backgroundColor: '#1890ff', marginBottom: 16 }}
                />
                <div>
                  <Text strong style={{ fontSize: 16 }}>
                    {user.username}
                  </Text>
                  <br />
                  <Text type="secondary">
                    {user.email}
                  </Text>
                </div>
              </div>
            </Card>

            {/* 安全设置 */}
            <Card 
              title={
                <Space>
                  <LockOutlined />
                  <span>安全设置</span>
                </Space>
              }
              size="small"
            >
              <Space direction="vertical" style={{ width: '100%' }}>
                <div style={{ 
                  display: 'flex', 
                  justifyContent: 'space-between', 
                  alignItems: 'center' 
                }}>
                  <div>
                    <Text>登录密码</Text>
                    <br />
                    <Text type="secondary" style={{ fontSize: 12 }}>
                      定期更换密码可提高安全性
                    </Text>
                  </div>
                  <Button
                    type="link"
                    icon={<KeyOutlined />}
                    onClick={() => setPasswordModalVisible(true)}
                  >
                    修改
                  </Button>
                </div>
              </Space>
            </Card>

            {/* 账户状态 */}
            <Card 
              title={
                <Space>
                  <InfoCircleOutlined />
                  <span>账户状态</span>
                </Space>
              }
              size="small"
            >
              <Space direction="vertical" style={{ width: '100%' }} size="small">
                <div style={{ 
                  display: 'flex', 
                  justifyContent: 'space-between', 
                  alignItems: 'center' 
                }}>
                  <Text>邮箱验证</Text>
                  <Tag color="green">已验证</Tag>
                </div>
                <div style={{ 
                  display: 'flex', 
                  justifyContent: 'space-between', 
                  alignItems: 'center' 
                }}>
                  <Text>账户状态</Text>
                  <Tag color="green">正常</Tag>
                </div>
                <div style={{ 
                  display: 'flex', 
                  justifyContent: 'space-between', 
                  alignItems: 'center' 
                }}>
                  <Text>最后登录</Text>
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    {dayjs().format('MM-DD HH:mm')}
                  </Text>
                </div>
              </Space>
            </Card>
          </Space>
        </Col>
      </Row>

      {/* 修改密码模态框 */}
      <Modal
        title={
          <Space>
            <KeyOutlined />
            <span>修改密码</span>
          </Space>
        }
        open={passwordModalVisible}
        onCancel={() => {
          setPasswordModalVisible(false)
          passwordForm.resetFields()
        }}
        onOk={() => passwordForm.submit()}
        confirmLoading={isLoading}
        width={500}
      >
        <Alert
          message="密码安全提示"
          description="建议使用包含大小写字母、数字和特殊字符的复杂密码，长度至少8位。"
          type="info"
          showIcon
          style={{ marginBottom: 24 }}
        />

        <Form
          form={passwordForm}
          layout="vertical"
          onFinish={handlePasswordChange}
        >
          <Form.Item
            name="oldPassword"
            label="当前密码"
            rules={[{ required: true, message: '请输入当前密码' }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="请输入当前密码"
              autoComplete="current-password"
            />
          </Form.Item>

          <Form.Item
            name="newPassword"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 8, message: '密码至少8位字符' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('oldPassword') !== value) {
                    return Promise.resolve()
                  }
                  return Promise.reject(new Error('新密码不能与当前密码相同'))
                }
              })
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="请输入新密码"
              autoComplete="new-password"
            />
          </Form.Item>

          {watchedPassword && (
            <div style={{ marginBottom: 16 }}>
              <div style={{ 
                display: 'flex', 
                justifyContent: 'space-between', 
                marginBottom: 8 
              }}>
                <Text style={{ fontSize: 12 }}>密码强度:</Text>
                <Text 
                  style={{ 
                    fontSize: 12, 
                    color: getPasswordStrengthColor(passwordStrength) 
                  }}
                >
                  {getPasswordStrengthText(passwordStrength)}
                </Text>
              </div>
              <Progress
                percent={passwordStrength}
                strokeColor={getPasswordStrengthColor(passwordStrength)}
                showInfo={false}
                size="small"
              />
            </div>
          )}

          <Form.Item
            name="confirmPassword"
            label="确认新密码"
            dependencies={['newPassword']}
            rules={[
              { required: true, message: '请确认新密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('newPassword') === value) {
                    return Promise.resolve()
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'))
                }
              })
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="请再次输入新密码"
              autoComplete="new-password"
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}