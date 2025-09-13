import { create } from 'zustand'
import { dnsProviderAPI, dnsRecordAPI, type DNSProvider, type DNSRecord } from '@/services/api'
import { message } from 'antd'

interface DNSState {
  // DNS 提供商状态
  providers: DNSProvider[]
  providersLoading: boolean
  
  // DNS 记录状态
  records: DNSRecord[]
  recordsLoading: boolean
  recordsTotal: number
  recordsPagination: {
    page: number
    limit: number
  }
  
  // 筛选条件
  filters: {
    domain?: string
    provider_id?: string
    record_type?: string
  }
  
  // Actions - DNS 提供商
  fetchProviders: () => Promise<void>
  createProvider: (data: Partial<DNSProvider>) => Promise<DNSProvider>
  updateProvider: (id: string, data: Partial<DNSProvider>) => Promise<DNSProvider>
  deleteProvider: (id: string) => Promise<void>
  testProvider: (id: string) => Promise<{ success: boolean; message: string }>
  
  // Actions - DNS 记录
  fetchRecords: (params?: { domain?: string; provider_id?: string; page?: number; limit?: number }) => Promise<void>
  createRecord: (data: Partial<DNSRecord>) => Promise<DNSRecord>
  updateRecord: (id: string, data: Partial<DNSRecord>) => Promise<DNSRecord>
  deleteRecord: (id: string) => Promise<void>
  batchCreateRecords: (records: Partial<DNSRecord>[]) => Promise<DNSRecord[]>
  batchDeleteRecords: (ids: string[]) => Promise<void>
  
  // Utilities
  setFilters: (filters: Partial<{ domain?: string; provider_id?: string; record_type?: string }>) => void
  setPagination: (page: number, limit?: number) => void
  clearFilters: () => void
}

export const useDNSStore = create<DNSState>((set, get) => ({
  // 初始状态
  providers: [],
  providersLoading: false,
  records: [],
  recordsLoading: false,
  recordsTotal: 0,
  recordsPagination: {
    page: 1,
    limit: 20
  },
  filters: {},
  
  // DNS 提供商管理
  fetchProviders: async () => {
    set({ providersLoading: true })
    try {
      const response = await dnsProviderAPI.list()
      set({ 
        providers: response.data,
        providersLoading: false 
      })
    } catch (error: any) {
      set({ providersLoading: false })
      message.error('获取DNS提供商列表失败')
      throw error
    }
  },
  
  createProvider: async (data: Partial<DNSProvider>) => {
    try {
      const response = await dnsProviderAPI.create(data)
      const newProvider = response.data
      
      set(state => ({
        providers: [...state.providers, newProvider]
      }))
      
      message.success('DNS提供商创建成功')
      return newProvider
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'DNS提供商创建失败'
      message.error(errorMessage)
      throw error
    }
  },
  
  updateProvider: async (id: string, data: Partial<DNSProvider>) => {
    try {
      const response = await dnsProviderAPI.update(id, data)
      const updatedProvider = response.data
      
      set(state => ({
        providers: state.providers.map(p => 
          p.id === id ? updatedProvider : p
        )
      }))
      
      message.success('DNS提供商更新成功')
      return updatedProvider
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'DNS提供商更新失败'
      message.error(errorMessage)
      throw error
    }
  },
  
  deleteProvider: async (id: string) => {
    try {
      await dnsProviderAPI.delete(id)
      
      set(state => ({
        providers: state.providers.filter(p => p.id !== id)
      }))
      
      message.success('DNS提供商删除成功')
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'DNS提供商删除失败'
      message.error(errorMessage)
      throw error
    }
  },
  
  testProvider: async (id: string) => {
    try {
      const response = await dnsProviderAPI.test(id)
      const result = response.data
      
      if (result.success) {
        message.success(result.message || 'DNS提供商连接测试成功')
      } else {
        message.warning(result.message || 'DNS提供商连接测试失败')
      }
      
      return result
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'DNS提供商连接测试失败'
      message.error(errorMessage)
      throw error
    }
  },
  
  // DNS 记录管理
  fetchRecords: async (params) => {
    set({ recordsLoading: true })
    try {
      const { filters, recordsPagination } = get()
      const requestParams = {
        ...filters,
        ...recordsPagination,
        ...params
      }
      
      const response = await dnsRecordAPI.list(requestParams)
      
      set({
        records: response.data.records,
        recordsTotal: response.data.total,
        recordsLoading: false
      })
    } catch (error: any) {
      set({ recordsLoading: false })
      message.error('获取DNS记录列表失败')
      throw error
    }
  },
  
  createRecord: async (data: Partial<DNSRecord>) => {
    try {
      const response = await dnsRecordAPI.create(data)
      const newRecord = response.data
      
      set(state => ({
        records: [newRecord, ...state.records],
        recordsTotal: state.recordsTotal + 1
      }))
      
      message.success('DNS记录创建成功')
      return newRecord
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'DNS记录创建失败'
      message.error(errorMessage)
      throw error
    }
  },
  
  updateRecord: async (id: string, data: Partial<DNSRecord>) => {
    try {
      const response = await dnsRecordAPI.update(id, data)
      const updatedRecord = response.data
      
      set(state => ({
        records: state.records.map(r => 
          r.id === id ? updatedRecord : r
        )
      }))
      
      message.success('DNS记录更新成功')
      return updatedRecord
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'DNS记录更新失败'
      message.error(errorMessage)
      throw error
    }
  },
  
  deleteRecord: async (id: string) => {
    try {
      await dnsRecordAPI.delete(id)
      
      set(state => ({
        records: state.records.filter(r => r.id !== id),
        recordsTotal: state.recordsTotal - 1
      }))
      
      message.success('DNS记录删除成功')
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || 'DNS记录删除失败'
      message.error(errorMessage)
      throw error
    }
  },
  
  batchCreateRecords: async (records: Partial<DNSRecord>[]) => {
    try {
      const response = await dnsRecordAPI.batchCreate(records)
      const newRecords = response.data
      
      set(state => ({
        records: [...newRecords, ...state.records],
        recordsTotal: state.recordsTotal + newRecords.length
      }))
      
      message.success(`批量创建 ${newRecords.length} 条DNS记录成功`)
      return newRecords
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || '批量创建DNS记录失败'
      message.error(errorMessage)
      throw error
    }
  },
  
  batchDeleteRecords: async (ids: string[]) => {
    try {
      await dnsRecordAPI.batchDelete(ids)
      
      set(state => ({
        records: state.records.filter(r => !ids.includes(r.id)),
        recordsTotal: state.recordsTotal - ids.length
      }))
      
      message.success(`批量删除 ${ids.length} 条DNS记录成功`)
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || '批量删除DNS记录失败'
      message.error(errorMessage)
      throw error
    }
  },
  
  // 工具方法
  setFilters: (newFilters) => {
    set(state => ({
      filters: { ...state.filters, ...newFilters }
    }))
  },
  
  setPagination: (page: number, limit?: number) => {
    set(state => ({
      recordsPagination: {
        page,
        limit: limit || state.recordsPagination.limit
      }
    }))
  },
  
  clearFilters: () => {
    set({ 
      filters: {},
      recordsPagination: { page: 1, limit: 20 }
    })
  }
}))