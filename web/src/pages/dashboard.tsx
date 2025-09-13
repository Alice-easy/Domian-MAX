import React, { useEffect, useState } from 'react'
import { 
  Typography, 
  Row, 
  Col, 
  Card, 
  Statistic, 
  Table, 
  Tag, 
  Button, 
  Space,
  Alert,
  Skeleton,
  Progress
} from 'antd'
import { 
  UserOutlined, 
  DatabaseOutlined, 
  GlobalOutlined, 
  CloudServerOutlined,
  PlusOutlined,
  ReloadOutlined,
  SettingOutlined
} from '@ant-design/icons'
import { useAuthStore } from '@/stores/authStore'
import { useDNSStore } from '@/stores/dnsStore'
import { useNavigate } from 'react-router-dom'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'
import { getDNSRecordTypeColor } from '@/config/constants'

// 配置 dayjs
dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const { Title, Text } = Typography

export default function Dashboard() {
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const { 
    providers, 
    records, 
    providersLoading, 
    recordsLoading,
    fetchProviders, 
    fetchRecords 
  } = useDNSStore()
  
  const [stats, setStats] = useState({
    totalDomains: 0,
    totalRecords: 0,
    activeProviders: 0,
    recentActivity: 0
  })

  useEffect(() => {
    // 获取数据
    fetchProviders().catch(console.error)
    fetchRecords({ limit: 10 }).catch(console.error)
  }, [fetchProviders, fetchRecords])

  useEffect(() => {
    // 计算统计数据
    const uniqueDomains = new Set(records.map(r => r.domain)).size
    const activeProviders = providers.filter(p => p.is_active).length
    const recentActivity = records.filter(r => 
      dayjs().diff(dayjs(r.created_at), 'day') <= 7
    ).length

    setStats({
      totalDomains: uniqueDomains,
      totalRecords: records.length,
      activeProviders,
      recentActivity
    })
  }, [providers, records])

  const recentRecordsColumns = [
    {
      title: '域名',
      dataIndex: 'domain',
      key: 'domain',
      render: (domain: string) => (
        <Text strong>{domain}</Text>
      )
    },
    {
      title: '记录类型',
      dataIndex: 'record_type',
      key: 'record_type',
      render: (type: string) => (
        <Tag color={getRecordTypeColor(type)}>{type}</Tag>
      )
    },
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name'
    },
    {
      title: '值',
      dataIndex: 'value',
      key: 'value',
      ellipsis: true
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => dayjs(date).format('MM-DD HH:mm')
    }
  ]

  const providerStatusColumns = [
    {
      title: '提供商',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: any) => (
        <Space>
          <CloudServerOutlined />
          <Text strong>{name}</Text>
        </Space>
      )
    },
    {
      title: '类型',
      dataIndex: 'provider_type',
      key: 'provider_type',
      render: (type: string) => (
        <Tag>{type}</Tag>
      )
    },
    {
      title: '状态',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (isActive: boolean) => (
        <Tag color={isActive ? 'green' : 'red'}>
          {isActive ? '活跃' : '禁用'}
        </Tag>
      )
    },
    {
      title: '最后更新',
      dataIndex: 'updated_at',
      key: 'updated_at',
      render: (date: string) => dayjs(date).fromNow()
    }
  ]

  const getRecordTypeColor = (type: string) => {
    return getDNSRecordTypeColor(type)
  }

  const refreshData = () => {
    fetchProviders()
    fetchRecords({ limit: 10 })
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
            欢迎回来，{user?.username}！
          </Title>
          <Text type="secondary">
            这是您的 DNS 管理概览
          </Text>
        </div>
        <Space>
          <Button 
            icon={<PlusOutlined />}
            type="primary"
            onClick={() => navigate('/dns-records?action=create')}
          >
            添加记录
          </Button>
          <Button 
            icon={<ReloadOutlined />}
            onClick={refreshData}
            loading={providersLoading || recordsLoading}
          >
            刷新
          </Button>
        </Space>
      </div>

      {/* 欢迎提示 */}
      {providers.length === 0 && (
        <Alert
          message="开始使用 Domain MAX"
          description={
            <div>
              您还没有配置任何 DNS 提供商。
              <Button 
                type="link" 
                style={{ padding: 0, marginLeft: 4 }}
                onClick={() => navigate('/admin/providers')}
              >
                立即配置
              </Button>
            </div>
          }
          type="info"
          showIcon
          style={{ marginBottom: 24 }}
        />
      )}

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="管理的域名"
              value={stats.totalDomains}
              prefix={<GlobalOutlined style={{ color: '#1890ff' }} />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="DNS 记录"
              value={stats.totalRecords}
              prefix={<DatabaseOutlined style={{ color: '#52c41a' }} />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="活跃提供商"
              value={stats.activeProviders}
              suffix={`/ ${providers.length}`}
              prefix={<CloudServerOutlined style={{ color: '#722ed1' }} />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="近7天活动"
              value={stats.recentActivity}
              prefix={<UserOutlined style={{ color: '#fa8c16' }} />}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        {/* 最近的 DNS 记录 */}
        <Col xs={24} lg={14}>
          <Card 
            title="最近的 DNS 记录" 
            extra={
              <Button 
                type="link" 
                onClick={() => navigate('/dns-records')}
              >
                查看全部
              </Button>
            }
          >
            {recordsLoading ? (
              <Skeleton active />
            ) : (
              <Table
                dataSource={records.slice(0, 5)}
                columns={recentRecordsColumns}
                pagination={false}
                size="small"
                rowKey="id"
                locale={{ emptyText: '暂无 DNS 记录' }}
              />
            )}
          </Card>
        </Col>

        {/* DNS 提供商状态 */}
        <Col xs={24} lg={10}>
          <Card 
            title="DNS 提供商状态"
            extra={
              <Button 
                type="link" 
                icon={<SettingOutlined />}
                onClick={() => navigate('/admin/providers')}
              >
                管理
              </Button>
            }
          >
            {providersLoading ? (
              <Skeleton active />
            ) : (
              <Table
                dataSource={providers.slice(0, 4)}
                columns={providerStatusColumns}
                pagination={false}
                size="small"
                rowKey="id"
                locale={{ emptyText: '暂无 DNS 提供商' }}
              />
            )}
          </Card>
        </Col>
      </Row>

      {/* 系统健康状态 */}
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col span={24}>
          <Card title="系统状态">
            <Row gutter={16}>
              <Col xs={24} md={8}>
                <div style={{ marginBottom: 16 }}>
                  <Text>DNS 提供商连通性</Text>
                  <Progress 
                    percent={Math.round((stats.activeProviders / Math.max(providers.length, 1)) * 100)}
                    strokeColor="#52c41a"
                    size="small"
                  />
                </div>
              </Col>
              <Col xs={24} md={8}>
                <div style={{ marginBottom: 16 }}>
                  <Text>记录管理效率</Text>
                  <Progress 
                    percent={85}
                    strokeColor="#1890ff"
                    size="small"
                  />
                </div>
              </Col>
              <Col xs={24} md={8}>
                <div style={{ marginBottom: 16 }}>
                  <Text>系统性能</Text>
                  <Progress 
                    percent={92}
                    strokeColor="#722ed1"
                    size="small"
                  />
                </div>
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>
    </div>
  )
}