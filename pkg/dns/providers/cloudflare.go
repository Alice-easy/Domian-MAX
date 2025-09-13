package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// CloudflareProvider CloudFlare DNS服务商
type CloudflareProvider struct {
	config     ProviderConfig
	httpClient *http.Client
	endpoint   string
}

// NewCloudflareProvider 创建CloudFlare DNS服务商实例
func NewCloudflareProvider(config ProviderConfig) (*CloudflareProvider, error) {
	if config.Token == "" && config.APIKey == "" {
		return nil, fmt.Errorf("CloudFlare DNS需要API Token或Global API Key")
	}
	
	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = "https://api.cloudflare.com/client/v4"
	}
	
	return &CloudflareProvider{
		config:   config,
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// GetName 获取服务商名称
func (p *CloudflareProvider) GetName() string {
	return "cloudflare"
}

// ValidateConfig 验证API配置
func (p *CloudflareProvider) ValidateConfig() error {
	if p.config.Token == "" && p.config.APIKey == "" {
		return fmt.Errorf("CloudFlare DNS需要API Token或Global API Key")
	}
	return nil
}

// TestConnection 测试连接
func (p *CloudflareProvider) TestConnection(ctx context.Context) error {
	// 通过获取用户信息来测试连接
	_, err := p.makeRequest(ctx, "GET", "/user", nil)
	return err
}

// ListRecords 获取域名记录列表
func (p *CloudflareProvider) ListRecords(ctx context.Context, domain string) ([]DNSRecord, error) {
	// 首先获取域名的Zone ID
	zoneID, err := p.getZoneID(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("获取Zone ID失败: %v", err)
	}
	
	// 获取DNS记录
	path := fmt.Sprintf("/zones/%s/dns_records", zoneID)
	response, err := p.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		Result []struct {
			ID       string      `json:"id"`
			Name     string      `json:"name"`
			Type     string      `json:"type"`
			Content  string      `json:"content"`
			TTL      int         `json:"ttl"`
			Priority interface{} `json:"priority"`
			Data     struct {
				Priority int `json:"priority"`
				Weight   int `json:"weight"`
				Port     int `json:"port"`
			} `json:"data"`
			Proxied bool `json:"proxied"`
		} `json:"result"`
		Success bool `json:"success"`
	}
	
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	
	if !result.Success {
		return nil, fmt.Errorf("CloudFlare API返回失败")
	}
	
	var records []DNSRecord
	for _, record := range result.Result {
		// 提取子域名（去掉主域名部分）
		name := record.Name
		if strings.HasSuffix(name, "."+domain) {
			name = strings.TrimSuffix(name, "."+domain)
		} else if name == domain {
			name = "@"
		}
		
		priority := 0
		if record.Priority != nil {
			if p, ok := record.Priority.(float64); ok {
				priority = int(p)
			}
		}
		if priority == 0 && record.Data.Priority > 0 {
			priority = record.Data.Priority
		}
		
		status := "active"
		if record.Proxied {
			status = "proxied"
		}
		
		records = append(records, DNSRecord{
			ID:       record.ID,
			Name:     name,
			Type:     record.Type,
			Value:    record.Content,
			TTL:      record.TTL,
			Priority: priority,
			Weight:   record.Data.Weight,
			Port:     record.Data.Port,
			Status:   status,
		})
	}
	
	return records, nil
}

// AddRecord 添加DNS记录
func (p *CloudflareProvider) AddRecord(ctx context.Context, domain string, record DNSRecord) (*DNSRecord, error) {
	// 获取Zone ID
	zoneID, err := p.getZoneID(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("获取Zone ID失败: %v", err)
	}
	
	// 构建完整的记录名称
	recordName := record.Name
	if recordName == "@" {
		recordName = domain
	} else {
		recordName = record.Name + "." + domain
	}
	
	// 构建请求体
	data := map[string]interface{}{
		"type":    record.Type,
		"name":    recordName,
		"content": record.Value,
		"ttl":     record.TTL,
	}
	
	// 设置优先级（MX和SRV记录）
	if record.Priority > 0 && (record.Type == "MX" || record.Type == "SRV") {
		data["priority"] = record.Priority
	}
	
	// SRV记录需要特殊处理
	if record.Type == "SRV" {
		data["data"] = map[string]interface{}{
			"priority": record.Priority,
			"weight":   record.Weight,
			"port":     record.Port,
			"target":   record.Value,
		}
		// SRV记录的content字段格式为：priority weight port target
		data["content"] = fmt.Sprintf("%d %d %d %s", record.Priority, record.Weight, record.Port, record.Value)
	}
	
	path := fmt.Sprintf("/zones/%s/dns_records", zoneID)
	response, err := p.makeRequest(ctx, "POST", path, data)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		Result struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Type    string `json:"type"`
			Content string `json:"content"`
			TTL     int    `json:"ttl"`
		} `json:"result"`
		Success bool `json:"success"`
	}
	
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	
	if !result.Success {
		return nil, fmt.Errorf("CloudFlare API返回失败")
	}
	
	// 返回创建的记录
	createdRecord := record
	createdRecord.ID = result.Result.ID
	
	return &createdRecord, nil
}

// UpdateRecord 更新DNS记录
func (p *CloudflareProvider) UpdateRecord(ctx context.Context, domain string, recordID string, record DNSRecord) error {
	// 获取Zone ID
	zoneID, err := p.getZoneID(ctx, domain)
	if err != nil {
		return fmt.Errorf("获取Zone ID失败: %v", err)
	}
	
	// 构建完整的记录名称
	recordName := record.Name
	if recordName == "@" {
		recordName = domain
	} else {
		recordName = record.Name + "." + domain
	}
	
	// 构建请求体
	data := map[string]interface{}{
		"type":    record.Type,
		"name":    recordName,
		"content": record.Value,
		"ttl":     record.TTL,
	}
	
	// 设置优先级（MX和SRV记录）
	if record.Priority > 0 && (record.Type == "MX" || record.Type == "SRV") {
		data["priority"] = record.Priority
	}
	
	// SRV记录需要特殊处理
	if record.Type == "SRV" {
		data["data"] = map[string]interface{}{
			"priority": record.Priority,
			"weight":   record.Weight,
			"port":     record.Port,
			"target":   record.Value,
		}
		data["content"] = fmt.Sprintf("%d %d %d %s", record.Priority, record.Weight, record.Port, record.Value)
	}
	
	path := fmt.Sprintf("/zones/%s/dns_records/%s", zoneID, recordID)
	_, err = p.makeRequest(ctx, "PUT", path, data)
	return err
}

// DeleteRecord 删除DNS记录
func (p *CloudflareProvider) DeleteRecord(ctx context.Context, domain string, recordID string) error {
	// 获取Zone ID
	zoneID, err := p.getZoneID(ctx, domain)
	if err != nil {
		return fmt.Errorf("获取Zone ID失败: %v", err)
	}
	
	path := fmt.Sprintf("/zones/%s/dns_records/%s", zoneID, recordID)
	_, err = p.makeRequest(ctx, "DELETE", path, nil)
	return err
}

// GetRecord 获取单个记录详情
func (p *CloudflareProvider) GetRecord(ctx context.Context, domain string, recordID string) (*DNSRecord, error) {
	// 获取Zone ID
	zoneID, err := p.getZoneID(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("获取Zone ID失败: %v", err)
	}
	
	path := fmt.Sprintf("/zones/%s/dns_records/%s", zoneID, recordID)
	response, err := p.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		Result struct {
			ID       string      `json:"id"`
			Name     string      `json:"name"`
			Type     string      `json:"type"`
			Content  string      `json:"content"`
			TTL      int         `json:"ttl"`
			Priority interface{} `json:"priority"`
			Data     struct {
				Priority int `json:"priority"`
				Weight   int `json:"weight"`
				Port     int `json:"port"`
			} `json:"data"`
			Proxied bool `json:"proxied"`
		} `json:"result"`
		Success bool `json:"success"`
	}
	
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	
	if !result.Success {
		return nil, fmt.Errorf("CloudFlare API返回失败")
	}
	
	record := result.Result
	
	// 提取子域名
	name := record.Name
	if strings.HasSuffix(name, "."+domain) {
		name = strings.TrimSuffix(name, "."+domain)
	} else if name == domain {
		name = "@"
	}
	
	priority := 0
	if record.Priority != nil {
		if p, ok := record.Priority.(float64); ok {
			priority = int(p)
		}
	}
	if priority == 0 && record.Data.Priority > 0 {
		priority = record.Data.Priority
	}
	
	status := "active"
	if record.Proxied {
		status = "proxied"
	}
	
	return &DNSRecord{
		ID:       record.ID,
		Name:     name,
		Type:     record.Type,
		Value:    record.Content,
		TTL:      record.TTL,
		Priority: priority,
		Weight:   record.Data.Weight,
		Port:     record.Data.Port,
		Status:   status,
	}, nil
}

// BatchAddRecords 批量添加DNS记录
func (p *CloudflareProvider) BatchAddRecords(ctx context.Context, domain string, records []DNSRecord) ([]DNSRecord, error) {
	// CloudFlare不支持批量操作，逐个添加
	var results []DNSRecord
	var errors []error
	
	for _, record := range records {
		result, err := p.AddRecord(ctx, domain, record)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		results = append(results, *result)
	}
	
	if len(errors) > 0 {
		return results, fmt.Errorf("批量添加记录时发生错误: %v", errors)
	}
	
	return results, nil
}

// getZoneID 获取域名的Zone ID
func (p *CloudflareProvider) getZoneID(ctx context.Context, domain string) (string, error) {
	path := "/zones?name=" + domain
	response, err := p.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return "", err
	}
	
	var result struct {
		Result []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"result"`
		Success bool `json:"success"`
	}
	
	if err := json.Unmarshal(response, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}
	
	if !result.Success {
		return "", fmt.Errorf("CloudFlare API返回失败")
	}
	
	for _, zone := range result.Result {
		if zone.Name == domain {
			return zone.ID, nil
		}
	}
	
	return "", fmt.Errorf("域名不存在: %s", domain)
}

// makeRequest 发起API请求
func (p *CloudflareProvider) makeRequest(ctx context.Context, method, path string, data interface{}) ([]byte, error) {
	var body io.Reader
	
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("序列化请求数据失败: %v", err)
		}
		body = strings.NewReader(string(jsonData))
	}
	
	// 创建请求
	req, err := http.NewRequestWithContext(ctx, method, p.endpoint+path, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	
	// 设置认证头
	if p.config.Token != "" {
		// 使用API Token
		req.Header.Set("Authorization", "Bearer "+p.config.Token)
	} else if p.config.APIKey != "" {
		// 使用Global API Key（需要邮箱）
		if email, exists := p.config.ExtraParams["email"]; exists {
			req.Header.Set("X-Auth-Email", email)
			req.Header.Set("X-Auth-Key", p.config.APIKey)
		} else {
			return nil, fmt.Errorf("使用Global API Key时需要提供邮箱地址")
		}
	}
	
	// 发起请求
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发起请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	
	// 检查HTTP状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(responseBody))
	}
	
	// 检查API错误
	var errorResp struct {
		Success bool `json:"success"`
		Errors  []struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
	}
	
	if err := json.Unmarshal(responseBody, &errorResp); err == nil && !errorResp.Success && len(errorResp.Errors) > 0 {
		var errorMsgs []string
		for _, apiErr := range errorResp.Errors {
			errorMsgs = append(errorMsgs, fmt.Sprintf("[%d] %s", apiErr.Code, apiErr.Message))
		}
		return nil, fmt.Errorf("CloudFlare API错误: %s", strings.Join(errorMsgs, "; "))
	}
	
	return responseBody, nil
}