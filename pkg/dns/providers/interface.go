package providers

import (
	"context"
	"time"
)

// DNSProvider DNS服务商接口
type DNSProvider interface {
	// GetName 获取服务商名称
	GetName() string
	
	// ValidateConfig 验证API配置
	ValidateConfig() error
	
	// ListRecords 获取域名记录列表
	ListRecords(ctx context.Context, domain string) ([]DNSRecord, error)
	
	// AddRecord 添加DNS记录
	AddRecord(ctx context.Context, domain string, record DNSRecord) (*DNSRecord, error)
	
	// UpdateRecord 更新DNS记录
	UpdateRecord(ctx context.Context, domain string, recordID string, record DNSRecord) error
	
	// DeleteRecord 删除DNS记录
	DeleteRecord(ctx context.Context, domain string, recordID string) error
	
	// GetRecord 获取单个记录详情
	GetRecord(ctx context.Context, domain string, recordID string) (*DNSRecord, error)
	
	// BatchAddRecords 批量添加DNS记录
	BatchAddRecords(ctx context.Context, domain string, records []DNSRecord) ([]DNSRecord, error)
	
	// TestConnection 测试连接
	TestConnection(ctx context.Context) error
}

// DNSRecord DNS记录结构
type DNSRecord struct {
	ID       string `json:"id"`
	Name     string `json:"name"`     // 子域名
	Type     string `json:"type"`     // 记录类型
	Value    string `json:"value"`    // 记录值
	TTL      int    `json:"ttl"`      // TTL值
	Priority int    `json:"priority"` // MX记录优先级
	Weight   int    `json:"weight"`   // SRV记录权重
	Port     int    `json:"port"`     // SRV记录端口
	Line     string `json:"line"`     // 解析线路
	Status   string `json:"status"`   // 记录状态
}

// ProviderConfig DNS服务商配置
type ProviderConfig struct {
	APIKey      string            `json:"api_key"`
	APISecret   string            `json:"api_secret"`
	Token       string            `json:"token"`
	Region      string            `json:"region"`
	Endpoint    string            `json:"endpoint"`
	ExtraParams map[string]string `json:"extra_params"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries    int           `json:"max_retries"`
	InitialDelay  time.Duration `json:"initial_delay"`
	MaxDelay      time.Duration `json:"max_delay"`
	BackoffFactor float64       `json:"backoff_factor"`
}

// DefaultRetryConfig 默认重试配置
var DefaultRetryConfig = RetryConfig{
	MaxRetries:    3,
	InitialDelay:  time.Second,
	MaxDelay:      time.Minute,
	BackoffFactor: 2.0,
}

// SupportedProviders 支持的DNS服务商列表
var SupportedProviders = []string{
	"aliyun",
	"dnspod",
	"huawei",
	"baidu",
	"west",
	"volcengine",
	"dnsla",
	"cloudflare",
	"namesilo",
	"powerdns",
}

// ProviderFeatures DNS服务商功能特性
type ProviderFeatures struct {
	SupportedRecordTypes []string `json:"supported_record_types"`
	SupportsBatch        bool     `json:"supports_batch"`
	SupportsLineTypes    bool     `json:"supports_line_types"`
	MaxRecordsPerDomain  int      `json:"max_records_per_domain"`
	MinTTL               int      `json:"min_ttl"`
	MaxTTL               int      `json:"max_ttl"`
}

// GetProviderFeatures 获取服务商功能特性
func GetProviderFeatures(providerType string) ProviderFeatures {
	features := map[string]ProviderFeatures{
		"aliyun": {
			SupportedRecordTypes: []string{"A", "AAAA", "CNAME", "MX", "TXT", "NS", "SRV", "CAA"},
			SupportsBatch:        true,
			SupportsLineTypes:    true,
			MaxRecordsPerDomain:  10000,
			MinTTL:               1,
			MaxTTL:               604800,
		},
		"dnspod": {
			SupportedRecordTypes: []string{"A", "AAAA", "CNAME", "MX", "TXT", "NS", "SRV"},
			SupportsBatch:        true,
			SupportsLineTypes:    true,
			MaxRecordsPerDomain:  10000,
			MinTTL:               1,
			MaxTTL:               604800,
		},
		"cloudflare": {
			SupportedRecordTypes: []string{"A", "AAAA", "CNAME", "MX", "TXT", "NS", "SRV", "CAA"},
			SupportsBatch:        false,
			SupportsLineTypes:    false,
			MaxRecordsPerDomain:  20000,
			MinTTL:               60,
			MaxTTL:               604800,
		},
	}
	
	if feature, exists := features[providerType]; exists {
		return feature
	}
	
	// 默认功能特性
	return ProviderFeatures{
		SupportedRecordTypes: []string{"A", "AAAA", "CNAME", "MX", "TXT", "NS"},
		SupportsBatch:        false,
		SupportsLineTypes:    false,
		MaxRecordsPerDomain:  1000,
		MinTTL:               300,
		MaxTTL:               86400,
	}
}