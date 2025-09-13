import React, { useEffect, useState } from 'react'
import {
  Typography,
  Table,
  Button,
  Space,
  Input,
  Select,
  Card,
  Tag,
  Modal,
  Form,
  InputNumber,
  message,
  Popconfirm,
  Tooltip,
  Alert,
  Row,
  Col,
  Divider
} from 'antd'
import {
  PlusOutlined,
  SearchOutlined,
  ReloadOutlined,
  EditOutlined,
  DeleteOutlined,
  CopyOutlined,
  ExportOutlined,
  ImportOutlined,
  FilterOutlined
} from '@ant-design/icons'
import { useDNSStore } from '@/stores/dns-store'
import type { DNSRecord, DNSProvider } from '@/services/api'
import dayjs from 'dayjs'

const { Title, Text } = Typography
const { Search } = Input
const { Option } = Select

interface DNSRecordFormData {
  domain: string
  record_type: string
  name: string
  value: string
  ttl: number
  priority?: number
  weight?: number
  port?: number
  provider_id: string
}

export default function DNSRecords() {
  const {
    records,
    providers,
    recordsLoading,
    recordsTotal,
    recordsPagination,
    filters,
    fetchRecords,
    fetchProviders,
    createRecord,
    updateRecord,
    deleteRecord,
    batchDeleteRecords,
    setFilters,
    setPagination
  } = useDNSStore()

  const [modalVisible, setModalVisible] = useState(false)
  const [editingRecord, setEditingRecord] = useState<DNSRecord | null>(null)
  const [selectedRowKeys, setSelectedRowKeys] = useState<string[]>([])
  const [searchText, setSearchText] = useState('')
  const [form] = Form.useForm<DNSRecordFormData>()

  useEffect(() => {
    fetchProviders()
    fetchRecords()
  }, [fetchProviders, fetchRecords])

  const recordTypes = [
    'A', 'AAAA', 'CNAME', 'MX', 'TXT', 'NS', 'SRV', 'PTR', 'SOA'
  ]

  const getRecordTypeColor = (type: string) => {
    const colors: { [key: string]: string } = {
      'A': 'blue',
      'AAAA': 'purple',
      'CNAME': 'orange',
      'MX': 'green',
      'TXT': 'red',
      'NS': 'cyan',
      'SRV': 'magenta',
      'PTR': 'gold',
      'SOA': 'lime'
    }
    return colors[type] || 'default'
  }

  const columns = [
    {
      title: '域名',
      dataIndex: 'domain',
      key: 'domain',
      sorter: true,
      render: (domain: string) => (
        <Text strong style={{ color: '#1890ff' }}>{domain}</Text>
      )
    },
    {
      title: '记录类型',
      dataIndex: 'record_type',
      key: 'record_type',
      width: 100,
      filters: recordTypes.map(type => ({ text: type, value: type })),
      render: (type: string) => (
        <Tag color={getRecordTypeColor(type)}>{type}</Tag>
      )
    },
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      ellipsis: true,
      render: (name: string) => (
        <Tooltip title={name}>
          <Text>{name || '@'}</Text>
        </Tooltip>
      )
    },
    {
      title: '值',
      dataIndex: 'value',
      key: 'value',
      ellipsis: true,
      render: (value: string) => (
        <Tooltip title={value}>
          <Text copyable={{ text: value }}>{value}</Text>
        </Tooltip>
      )
    },
    {
      title: 'TTL',
      dataIndex: 'ttl',
      key: 'ttl',
      width: 80,
      sorter: true,
      render: (ttl: number) => (
        <Text type="secondary">{ttl}s</Text>
      )
    },
    {
      title: 'DNS提供商',
      dataIndex: 'provider_id',
      key: 'provider_id',
      width: 120,
      render: (providerId: string) => {
        const provider = providers.find(p => p.id === providerId)
        return (
          <Tag color="processing">
            {provider?.name || '未知'}
          </Tag>
        )
      }
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 120,
      sorter: true,
      render: (date: string) => (
        <Text type="secondary">
          {dayjs(date).format('MM-DD HH:mm')}
        </Text>
      )
    },
    {
      title: '操作',
      key: 'actions',
      width: 120,
      render: (_, record: DNSRecord) => (
        <Space size="small">
          <Tooltip title="编辑">
            <Button
              type="text"
              size="small"
              icon={<EditOutlined />}
              onClick={() => handleEdit(record)}
            />
          </Tooltip>
          <Tooltip title="复制">
            <Button
              type="text"
              size="small"
              icon={<CopyOutlined />}
              onClick={() => handleCopy(record)}
            />
          </Tooltip>
          <Popconfirm
            title="确定要删除这条记录吗？"
            description="删除后无法恢复"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Tooltip title="删除">
              <Button
                type="text"
                size="small"
                icon={<DeleteOutlined />}
                danger
              />
            </Tooltip>
          </Popconfirm>
        </Space>
      )
    }
  ]

  const handleSearch = (value: string) => {
    setSearchText(value)
    setFilters({ domain: value })
    fetchRecords({ domain: value, page: 1 })
  }

  const handleFilterChange = (field: string, value: any) => {
    const newFilters = { ...filters, [field]: value }
    setFilters(newFilters)
    fetchRecords({ ...newFilters, page: 1 })
  }

  const handleTableChange = (pagination: any, tableFilters: any, sorter: any) => {
    const { current, pageSize } = pagination
    setPagination(current, pageSize)
    fetchRecords({
      ...filters,
      page: current,
      limit: pageSize
    })
  }

  const handleCreate = () => {
    setEditingRecord(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (record: DNSRecord) => {
    setEditingRecord(record)
    form.setFieldsValue(record)
    setModalVisible(true)
  }

  const handleCopy = (record: DNSRecord) => {
    const { id, created_at, updated_at, ...copyData } = record
    form.setFieldsValue({
      ...copyData,
      name: `${copyData.name}-copy`
    })
    setEditingRecord(null)
    setModalVisible(true)
  }

  const handleDelete = async (id: string) => {
    try {
      await deleteRecord(id)
    } catch (error) {
      // Error handling is done in the store
    }
  }

  const handleBatchDelete = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要删除的记录')
      return
    }

    try {
      await batchDeleteRecords(selectedRowKeys)
      setSelectedRowKeys([])
    } catch (error) {
      // Error handling is done in the store
    }
  }

  const handleSubmit = async (values: DNSRecordFormData) => {
    try {
      if (editingRecord) {
        await updateRecord(editingRecord.id, values)
      } else {
        await createRecord(values)
      }
      setModalVisible(false)
      form.resetFields()
    } catch (error) {
      // Error handling is done in the store
    }
  }

  const rowSelection = {
    selectedRowKeys,
    onChange: (keys: React.Key[]) => {
      setSelectedRowKeys(keys as string[])
    }
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
            DNS 记录管理
          </Title>
          <Text type="secondary">
            管理您的域名解析记录
          </Text>
        </div>
        <Space>
          <Button 
            icon={<ExportOutlined />}
            onClick={() => message.info('导出功能开发中')}
          >
            导出
          </Button>
          <Button 
            icon={<ImportOutlined />}
            onClick={() => message.info('导入功能开发中')}
          >
            导入
          </Button>
          <Button 
            icon={<PlusOutlined />}
            type="primary"
            onClick={handleCreate}
          >
            添加记录
          </Button>
        </Space>
      </div>

      {/* 筛选区域 */}
      <Card size="small" style={{ marginBottom: 16 }}>
        <Row gutter={16} align="middle">
          <Col xs={24} sm={8} md={6}>
            <Search
              placeholder="搜索域名"
              value={searchText}
              onChange={(e) => setSearchText(e.target.value)}
              onSearch={handleSearch}
              enterButton={<SearchOutlined />}
              allowClear
            />
          </Col>
          <Col xs={24} sm={8} md={4}>
            <Select
              placeholder="记录类型"
              style={{ width: '100%' }}
              value={filters.record_type}
              onChange={(value) => handleFilterChange('record_type', value)}
              allowClear
            >
              {recordTypes.map(type => (
                <Option key={type} value={type}>{type}</Option>
              ))}
            </Select>
          </Col>
          <Col xs={24} sm={8} md={4}>
            <Select
              placeholder="DNS提供商"
              style={{ width: '100%' }}
              value={filters.provider_id}
              onChange={(value) => handleFilterChange('provider_id', value)}
              allowClear
            >
              {providers.map(provider => (
                <Option key={provider.id} value={provider.id}>
                  {provider.name}
                </Option>
              ))}
            </Select>
          </Col>
          <Col xs={24} sm={24} md={10}>
            <Space>
              <Button 
                icon={<FilterOutlined />}
                onClick={() => {
                  setFilters({})
                  setSearchText('')
                  fetchRecords({ page: 1 })
                }}
              >
                清除筛选
              </Button>
              <Button 
                icon={<ReloadOutlined />}
                onClick={() => fetchRecords()}
                loading={recordsLoading}
              >
                刷新
              </Button>
              {selectedRowKeys.length > 0 && (
                <Popconfirm
                  title={`确定要删除选中的 ${selectedRowKeys.length} 条记录吗？`}
                  description="删除后无法恢复"
                  onConfirm={handleBatchDelete}
                  okText="确定"
                  cancelText="取消"
                >
                  <Button danger>
                    批量删除 ({selectedRowKeys.length})
                  </Button>
                </Popconfirm>
              )}
            </Space>
          </Col>
        </Row>
      </Card>

      {/* 提示信息 */}
      {providers.length === 0 && (
        <Alert
          message="暂无可用的DNS提供商"
          description="请先配置DNS提供商才能管理DNS记录。"
          type="warning"
          showIcon
          style={{ marginBottom: 16 }}
          action={
            <Button size="small" type="primary">
              配置提供商
            </Button>
          }
        />
      )}

      {/* 数据表格 */}
      <Card>
        <Table
          rowSelection={rowSelection}
          columns={columns}
          dataSource={records}
          loading={recordsLoading}
          rowKey="id"
          size="small"
          pagination={{
            current: recordsPagination.page,
            pageSize: recordsPagination.limit,
            total: recordsTotal,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) => 
              `第 ${range[0]}-${range[1]} 条，共 ${total} 条记录`,
            pageSizeOptions: ['10', '20', '50', '100']
          }}
          onChange={handleTableChange}
          scroll={{ x: 1000 }}
        />
      </Card>

      {/* 添加/编辑模态框 */}
      <Modal
        title={editingRecord ? '编辑DNS记录' : '添加DNS记录'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        onOk={() => form.submit()}
        width={600}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={{
            ttl: 600,
            record_type: 'A'
          }}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="domain"
                label="域名"
                rules={[
                  { required: true, message: '请输入域名' },
                  { 
                    pattern: /^([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$/,
                    message: '请输入有效的域名格式'
                  }
                ]}
              >
                <Input placeholder="example.com" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="record_type"
                label="记录类型"
                rules={[{ required: true, message: '请选择记录类型' }]}
              >
                <Select>
                  {recordTypes.map(type => (
                    <Option key={type} value={type}>{type}</Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="name"
                label="记录名称"
                tooltip="留空表示根域名 (@)"
              >
                <Input placeholder="www 或留空" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="ttl"
                label="TTL (秒)"
                rules={[
                  { required: true, message: '请输入TTL值' },
                  { type: 'number', min: 60, max: 86400, message: 'TTL范围: 60-86400秒' }
                ]}
              >
                <InputNumber
                  style={{ width: '100%' }}
                  placeholder="600"
                  min={60}
                  max={86400}
                />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="value"
            label="记录值"
            rules={[{ required: true, message: '请输入记录值' }]}
          >
            <Input.TextArea
              rows={3}
              placeholder="根据记录类型输入相应的值"
            />
          </Form.Item>

          <Form.Item
            name="provider_id"
            label="DNS提供商"
            rules={[{ required: true, message: '请选择DNS提供商' }]}
          >
            <Select placeholder="选择DNS提供商">
              {providers.filter(p => p.is_active).map(provider => (
                <Option key={provider.id} value={provider.id}>
                  {provider.name} ({provider.provider_type})
                </Option>
              ))}
            </Select>
          </Form.Item>

          {/* MX记录的优先级字段 */}
          <Form.Item
            noStyle
            shouldUpdate={(prevValues, currentValues) =>
              prevValues.record_type !== currentValues.record_type
            }
          >
            {({ getFieldValue }) => {
              const recordType = getFieldValue('record_type')
              if (recordType === 'MX') {
                return (
                  <Form.Item
                    name="priority"
                    label="优先级"
                    rules={[
                      { required: true, message: '请输入MX记录优先级' },
                      { type: 'number', min: 0, max: 65535 }
                    ]}
                  >
                    <InputNumber
                      style={{ width: '100%' }}
                      placeholder="10"
                      min={0}
                      max={65535}
                    />
                  </Form.Item>
                )
              }
              return null
            }}
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}