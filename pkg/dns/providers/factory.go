package providers

import (
	"context"
	"fmt"
	"time"
)

// ProviderFactory DNS服务商工厂
type ProviderFactory struct {
	retryConfig RetryConfig
}

// NewProviderFactory 创建DNS服务商工厂
func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		retryConfig: DefaultRetryConfig,
	}
}

// SetRetryConfig 设置重试配置
func (f *ProviderFactory) SetRetryConfig(config RetryConfig) {
	f.retryConfig = config
}

// IsSupported 检查是否支持指定的DNS提供商类型
func (f *ProviderFactory) IsSupported(providerType string) bool {
	supportedTypes := []string{
		"aliyun", "dnspod", "huawei", "baidu", "west",
		"volcengine", "dnsla", "cloudflare", "namesilo", "powerdns",
	}
	
	for _, supported := range supportedTypes {
		if supported == providerType {
			return true
		}
	}
	return false
}

// GetSupportedTypes 获取支持的DNS提供商类型列表
func (f *ProviderFactory) GetSupportedTypes() []string {
	return []string{
		"aliyun", "dnspod", "huawei", "baidu", "west",
		"volcengine", "dnsla", "cloudflare", "namesilo", "powerdns",
	}
}

// CreateProvider 创建DNS服务商实例
func (f *ProviderFactory) CreateProvider(providerType string, config map[string]string) (DNSProvider, error) {
	// 将map[string]string转换为ProviderConfig
	providerConfig := ProviderConfig{
		APIKey:      config["api_key"],
		APISecret:   config["api_secret"],
		Token:       config["token"],
		Region:      config["region"],
		Endpoint:    config["endpoint"],
		ExtraParams: config, // 保存所有原始参数
	}
	
	switch providerType {
	case "aliyun":
		return NewAliyunProvider(providerConfig)
	case "dnspod":
		return NewDNSPodProvider(providerConfig)
	case "huawei":
		return NewHuaweiProvider(providerConfig)
	case "baidu":
		return NewBaiduProvider(providerConfig)
	case "west":
		return NewWestProvider(providerConfig)
	case "volcengine":
		return NewVolcengineProvider(providerConfig)
	case "dnsla":
		return NewDNSLAProvider(providerConfig)
	case "cloudflare":
		return NewCloudflareProvider(providerConfig)
	case "namesilo":
		return NewNamesiloProvider(providerConfig)
	case "powerdns":
		return NewPowerDNSProvider(providerConfig)
	default:
		return nil, fmt.Errorf("不支持的DNS服务商: %s", providerType)
	}
}

// ProviderManager DNS服务商管理器
type ProviderManager struct {
	factory   *ProviderFactory
	providers map[string]DNSProvider
	configs   map[string]ProviderConfig
}

// NewProviderManager 创建DNS服务商管理器
func NewProviderManager() *ProviderManager {
	return &ProviderManager{
		factory:   NewProviderFactory(),
		providers: make(map[string]DNSProvider),
		configs:   make(map[string]ProviderConfig),
	}
}

// RegisterProvider 注册DNS服务商
func (m *ProviderManager) RegisterProvider(name, providerType string, config map[string]string) error {
	provider, err := m.factory.CreateProvider(providerType, config)
	if err != nil {
		return fmt.Errorf("创建DNS服务商失败: %v", err)
	}
	
	// 验证配置
	if err := provider.ValidateConfig(); err != nil {
		return fmt.Errorf("DNS服务商配置验证失败: %v", err)
	}
	
	m.providers[name] = provider
	// 将map配置转换为ProviderConfig保存
	providerConfig := ProviderConfig{
		APIKey:      config["api_key"],
		APISecret:   config["api_secret"],
		Token:       config["token"],
		Region:      config["region"],
		Endpoint:    config["endpoint"],
		ExtraParams: config,
	}
	m.configs[name] = providerConfig
	return nil
}

// GetProvider 获取DNS服务商
func (m *ProviderManager) GetProvider(name string) (DNSProvider, error) {
	provider, exists := m.providers[name]
	if !exists {
		return nil, fmt.Errorf("DNS服务商 '%s' 未注册", name)
	}
	return provider, nil
}

// ListProviders 获取所有已注册的DNS服务商
func (m *ProviderManager) ListProviders() []string {
	var names []string
	for name := range m.providers {
		names = append(names, name)
	}
	return names
}

// TestProvider 测试DNS服务商连接
func (m *ProviderManager) TestProvider(name string) error {
	provider, err := m.GetProvider(name)
	if err != nil {
		return err
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	return provider.TestConnection(ctx)
}

// TestAllProviders 测试所有DNS服务商连接
func (m *ProviderManager) TestAllProviders() map[string]error {
	results := make(map[string]error)
	for name := range m.providers {
		results[name] = m.TestProvider(name)
	}
	return results
}

// RemoveProvider 移除DNS服务商
func (m *ProviderManager) RemoveProvider(name string) {
	delete(m.providers, name)
	delete(m.configs, name)
}

// UpdateProvider 更新DNS服务商配置
func (m *ProviderManager) UpdateProvider(name, providerType string, config map[string]string) error {
	// 先移除旧的
	m.RemoveProvider(name)
	
	// 注册新的
	return m.RegisterProvider(name, providerType, config)
}

// RetryOperation 重试操作的通用方法
func (m *ProviderManager) RetryOperation(operation func() error) error {
	retryConfig := m.factory.retryConfig
	var lastErr error
	
	for attempt := 0; attempt <= retryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(float64(retryConfig.InitialDelay) * 
				pow(retryConfig.BackoffFactor, float64(attempt-1)))
			if delay > retryConfig.MaxDelay {
				delay = retryConfig.MaxDelay
			}
			time.Sleep(delay)
		}
		
		if err := operation(); err != nil {
			lastErr = err
			// 检查是否是可重试的错误
			if !isRetryableError(err) {
				return err
			}
			continue
		}
		
		return nil
	}
	
	return fmt.Errorf("操作失败，已重试%d次: %v", retryConfig.MaxRetries, lastErr)
}

// isRetryableError 判断错误是否可重试
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	// 网络相关错误通常可以重试
	errorStr := err.Error()
	retryableErrors := []string{
		"timeout",
		"connection refused",
		"connection reset",
		"network is unreachable",
		"temporary failure",
		"server error",
		"rate limit",
		"too many requests",
	}
	
	for _, retryable := range retryableErrors {
		if contains(errorStr, retryable) {
			return true
		}
	}
	
	return false
}

// pow 简单的幂运算实现
func pow(base, exp float64) float64 {
	if exp == 0 {
		return 1
	}
	result := base
	for i := 1; i < int(exp); i++ {
		result *= base
	}
	return result
}

// contains 检查字符串是否包含子字符串
func contains(str, substr string) bool {
	return len(str) >= len(substr) && 
		(len(substr) == 0 || str[len(str)-len(substr):] == substr || 
		 str[:len(substr)] == substr ||
		 (len(str) > len(substr) && findSubstring(str, substr)))
}

// findSubstring 查找子字符串
func findSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}