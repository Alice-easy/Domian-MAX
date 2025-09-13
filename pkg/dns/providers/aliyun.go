package providers

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// AliyunProvider 阿里云DNS服务商
type AliyunProvider struct {
	config     ProviderConfig
	httpClient *http.Client
	endpoint   string
}

// NewAliyunProvider 创建阿里云DNS服务商实例
func NewAliyunProvider(config ProviderConfig) (*AliyunProvider, error) {
	if config.APIKey == "" || config.APISecret == "" {
		return nil, fmt.Errorf("阿里云DNS需要API Key和API Secret")
	}
	
	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = "https://alidns.aliyuncs.com"
	}
	
	return &AliyunProvider{
		config:   config,
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// GetName 获取服务商名称
func (p *AliyunProvider) GetName() string {
	return "aliyun"
}

// ValidateConfig 验证API配置
func (p *AliyunProvider) ValidateConfig() error {
	if p.config.APIKey == "" {
		return fmt.Errorf("阿里云DNS API Key不能为空")
	}
	if p.config.APISecret == "" {
		return fmt.Errorf("阿里云DNS API Secret不能为空")
	}
	return nil
}

// TestConnection 测试连接
func (p *AliyunProvider) TestConnection(ctx context.Context) error {
	// 通过获取域名列表来测试连接
	params := map[string]string{
		"Action":  "DescribeDomains",
		"Version": "2015-01-09",
	}
	
	_, err := p.makeRequest(ctx, params)
	return err
}

// ListRecords 获取域名记录列表
func (p *AliyunProvider) ListRecords(ctx context.Context, domain string) ([]DNSRecord, error) {
	params := map[string]string{
		"Action":     "DescribeDomainRecords",
		"Version":    "2015-01-09",
		"DomainName": domain,
		"PageSize":   "500",
	}
	
	response, err := p.makeRequest(ctx, params)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		DomainRecords struct {
			Record []struct {
				RecordId string `json:"RecordId"`
				RR       string `json:"RR"`
				Type     string `json:"Type"`
				Value    string `json:"Value"`
				TTL      int    `json:"TTL"`
				Priority int    `json:"Priority"`
				Line     string `json:"Line"`
				Status   string `json:"Status"`
			} `json:"Record"`
		} `json:"DomainRecords"`
	}
	
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	
	var records []DNSRecord
	for _, record := range result.DomainRecords.Record {
		records = append(records, DNSRecord{
			ID:       record.RecordId,
			Name:     record.RR,
			Type:     record.Type,
			Value:    record.Value,
			TTL:      record.TTL,
			Priority: record.Priority,
			Line:     record.Line,
			Status:   record.Status,
		})
	}
	
	return records, nil
}

// AddRecord 添加DNS记录
func (p *AliyunProvider) AddRecord(ctx context.Context, domain string, record DNSRecord) (*DNSRecord, error) {
	params := map[string]string{
		"Action":     "AddDomainRecord",
		"Version":    "2015-01-09",
		"DomainName": domain,
		"RR":         record.Name,
		"Type":       record.Type,
		"Value":      record.Value,
		"TTL":        fmt.Sprintf("%d", record.TTL),
	}
	
	if record.Priority > 0 && (record.Type == "MX" || record.Type == "SRV") {
		params["Priority"] = fmt.Sprintf("%d", record.Priority)
	}
	
	if record.Line != "" {
		params["Line"] = record.Line
	}
	
	response, err := p.makeRequest(ctx, params)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		RecordId string `json:"RecordId"`
	}
	
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	
	// 返回创建的记录，包含分配的ID
	createdRecord := record
	createdRecord.ID = result.RecordId
	
	return &createdRecord, nil
}

// UpdateRecord 更新DNS记录
func (p *AliyunProvider) UpdateRecord(ctx context.Context, domain string, recordID string, record DNSRecord) error {
	params := map[string]string{
		"Action":   "UpdateDomainRecord",
		"Version":  "2015-01-09",
		"RecordId": recordID,
		"RR":       record.Name,
		"Type":     record.Type,
		"Value":    record.Value,
		"TTL":      fmt.Sprintf("%d", record.TTL),
	}
	
	if record.Priority > 0 && (record.Type == "MX" || record.Type == "SRV") {
		params["Priority"] = fmt.Sprintf("%d", record.Priority)
	}
	
	if record.Line != "" {
		params["Line"] = record.Line
	}
	
	_, err := p.makeRequest(ctx, params)
	return err
}

// DeleteRecord 删除DNS记录
func (p *AliyunProvider) DeleteRecord(ctx context.Context, domain string, recordID string) error {
	params := map[string]string{
		"Action":   "DeleteDomainRecord",
		"Version":  "2015-01-09",
		"RecordId": recordID,
	}
	
	_, err := p.makeRequest(ctx, params)
	return err
}

// GetRecord 获取单个记录详情
func (p *AliyunProvider) GetRecord(ctx context.Context, domain string, recordID string) (*DNSRecord, error) {
	params := map[string]string{
		"Action":     "DescribeDomainRecordInfo",
		"Version":    "2015-01-09",
		"RecordId":   recordID,
	}
	
	response, err := p.makeRequest(ctx, params)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		RecordId   string `json:"RecordId"`
		RR         string `json:"RR"`
		Type       string `json:"Type"`
		Value      string `json:"Value"`
		TTL        int    `json:"TTL"`
		Priority   int    `json:"Priority"`
		Line       string `json:"Line"`
		Status     string `json:"Status"`
		DomainName string `json:"DomainName"`
	}
	
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	
	return &DNSRecord{
		ID:       result.RecordId,
		Name:     result.RR,
		Type:     result.Type,
		Value:    result.Value,
		TTL:      result.TTL,
		Priority: result.Priority,
		Line:     result.Line,
		Status:   result.Status,
	}, nil
}

// BatchAddRecords 批量添加DNS记录
func (p *AliyunProvider) BatchAddRecords(ctx context.Context, domain string, records []DNSRecord) ([]DNSRecord, error) {
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

// makeRequest 发起API请求
func (p *AliyunProvider) makeRequest(ctx context.Context, params map[string]string) ([]byte, error) {
	// 添加公共参数
	params["Format"] = "JSON"
	params["AccessKeyId"] = p.config.APIKey
	params["SignatureMethod"] = "HMAC-SHA1"
	params["SignatureVersion"] = "1.0"
	params["SignatureNonce"] = generateNonce()
	params["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	
	// 生成签名
	signature := p.generateSignature(params)
	params["Signature"] = signature
	
	// 构建请求URL
	requestURL := p.buildRequestURL(params)
	
	// 发起请求
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发起请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	
	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}
	
	// 检查API错误
	var errorResp struct {
		Code      string `json:"Code"`
		Message   string `json:"Message"`
		RequestId string `json:"RequestId"`
	}
	
	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Code != "" {
		return nil, fmt.Errorf("阿里云DNS API错误: %s - %s (RequestId: %s)", 
			errorResp.Code, errorResp.Message, errorResp.RequestId)
	}
	
	return body, nil
}

// generateSignature 生成API签名
func (p *AliyunProvider) generateSignature(params map[string]string) string {
	// 排序参数
	var keys []string
	for k := range params {
		if k != "Signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	
	// 构建规范化查询字符串
	var parts []string
	for _, key := range keys {
		parts = append(parts, url.QueryEscape(key)+"="+url.QueryEscape(params[key]))
	}
	canonicalQueryString := strings.Join(parts, "&")
	
	// 构建待签名字符串
	stringToSign := "GET&" + url.QueryEscape("/") + "&" + url.QueryEscape(canonicalQueryString)
	
	// 计算签名
	key := p.config.APISecret + "&"
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	
	return signature
}

// buildRequestURL 构建请求URL
func (p *AliyunProvider) buildRequestURL(params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return p.endpoint + "/?" + values.Encode()
}

// generateNonce 生成随机字符串
func generateNonce() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}