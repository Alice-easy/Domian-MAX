package providers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DNSPodProvider 腾讯云DNSPod服务商
type DNSPodProvider struct {
	config     ProviderConfig
	httpClient *http.Client
	endpoint   string
}

// NewDNSPodProvider 创建腾讯云DNSPod服务商实例
func NewDNSPodProvider(config ProviderConfig) (*DNSPodProvider, error) {
	if config.APIKey == "" || config.APISecret == "" {
		return nil, fmt.Errorf("腾讯云DNSPod需要API Key和API Secret")
	}
	
	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = "https://dnspod.tencentcloudapi.com"
	}
	
	return &DNSPodProvider{
		config:   config,
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// GetName 获取服务商名称
func (p *DNSPodProvider) GetName() string {
	return "dnspod"
}

// ValidateConfig 验证API配置
func (p *DNSPodProvider) ValidateConfig() error {
	if p.config.APIKey == "" {
		return fmt.Errorf("腾讯云DNSPod API Key不能为空")
	}
	if p.config.APISecret == "" {
		return fmt.Errorf("腾讯云DNSPod API Secret不能为空")
	}
	return nil
}

// TestConnection 测试连接
func (p *DNSPodProvider) TestConnection(ctx context.Context) error {
	// 通过获取域名列表来测试连接
	_, err := p.makeRequest(ctx, "DescribeDomainList", map[string]interface{}{
		"Limit": 1,
	})
	return err
}

// ListRecords 获取域名记录列表
func (p *DNSPodProvider) ListRecords(ctx context.Context, domain string) ([]DNSRecord, error) {
	// 首先获取域名ID
	domainID, err := p.getDomainID(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("获取域名ID失败: %v", err)
	}
	
	params := map[string]interface{}{
		"Domain":   domain,
		"DomainId": domainID,
		"Limit":    3000,
	}
	
	response, err := p.makeRequest(ctx, "DescribeRecordList", params)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		Response struct {
			RecordList []struct {
				RecordId int    `json:"RecordId"`
				Name     string `json:"Name"`
				Type     string `json:"Type"`
				Value    string `json:"Value"`
				TTL      int    `json:"TTL"`
				MX       int    `json:"MX"`
				Line     string `json:"Line"`
				Status   string `json:"Status"`
			} `json:"RecordList"`
		} `json:"Response"`
	}
	
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	
	var records []DNSRecord
	for _, record := range result.Response.RecordList {
		records = append(records, DNSRecord{
			ID:       strconv.Itoa(record.RecordId),
			Name:     record.Name,
			Type:     record.Type,
			Value:    record.Value,
			TTL:      record.TTL,
			Priority: record.MX,
			Line:     record.Line,
			Status:   record.Status,
		})
	}
	
	return records, nil
}

// AddRecord 添加DNS记录
func (p *DNSPodProvider) AddRecord(ctx context.Context, domain string, record DNSRecord) (*DNSRecord, error) {
	// 获取域名ID
	domainID, err := p.getDomainID(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("获取域名ID失败: %v", err)
	}
	
	params := map[string]interface{}{
		"Domain":     domain,
		"DomainId":   domainID,
		"SubDomain":  record.Name,
		"RecordType": record.Type,
		"Value":      record.Value,
		"TTL":        record.TTL,
	}
	
	if record.Priority > 0 && (record.Type == "MX" || record.Type == "SRV") {
		params["MX"] = record.Priority
	}
	
	if record.Line != "" {
		params["RecordLine"] = record.Line
	}
	
	response, err := p.makeRequest(ctx, "CreateRecord", params)
	if err != nil {
		return nil, err
	}
	
	var result struct {
		Response struct {
			RecordId int `json:"RecordId"`
		} `json:"Response"`
	}
	
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}
	
	// 返回创建的记录
	createdRecord := record
	createdRecord.ID = strconv.Itoa(result.Response.RecordId)
	
	return &createdRecord, nil
}

// UpdateRecord 更新DNS记录
func (p *DNSPodProvider) UpdateRecord(ctx context.Context, domain string, recordID string, record DNSRecord) error {
	// 获取域名ID
	domainID, err := p.getDomainID(ctx, domain)
	if err != nil {
		return fmt.Errorf("获取域名ID失败: %v", err)
	}
	
	recordIDInt, err := strconv.Atoi(recordID)
	if err != nil {
		return fmt.Errorf("无效的记录ID: %s", recordID)
	}
	
	params := map[string]interface{}{
		"Domain":     domain,
		"DomainId":   domainID,
		"RecordId":   recordIDInt,
		"SubDomain":  record.Name,
		"RecordType": record.Type,
		"Value":      record.Value,
		"TTL":        record.TTL,
	}
	
	if record.Priority > 0 && (record.Type == "MX" || record.Type == "SRV") {
		params["MX"] = record.Priority
	}
	
	if record.Line != "" {
		params["RecordLine"] = record.Line
	}
	
	_, err = p.makeRequest(ctx, "ModifyRecord", params)
	return err
}

// DeleteRecord 删除DNS记录
func (p *DNSPodProvider) DeleteRecord(ctx context.Context, domain string, recordID string) error {
	// 获取域名ID
	domainID, err := p.getDomainID(ctx, domain)
	if err != nil {
		return fmt.Errorf("获取域名ID失败: %v", err)
	}
	
	recordIDInt, err := strconv.Atoi(recordID)
	if err != nil {
		return fmt.Errorf("无效的记录ID: %s", recordID)
	}
	
	params := map[string]interface{}{
		"Domain":   domain,
		"DomainId": domainID,
		"RecordId": recordIDInt,
	}
	
	_, err = p.makeRequest(ctx, "DeleteRecord", params)
	return err
}

// GetRecord 获取单个记录详情
func (p *DNSPodProvider) GetRecord(ctx context.Context, domain string, recordID string) (*DNSRecord, error) {
	// DNSPod API没有单独的获取记录接口，需要通过列表接口获取
	records, err := p.ListRecords(ctx, domain)
	if err != nil {
		return nil, err
	}
	
	for _, record := range records {
		if record.ID == recordID {
			return &record, nil
		}
	}
	
	return nil, fmt.Errorf("记录不存在: %s", recordID)
}

// BatchAddRecords 批量添加DNS记录
func (p *DNSPodProvider) BatchAddRecords(ctx context.Context, domain string, records []DNSRecord) ([]DNSRecord, error) {
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

// getDomainID 获取域名ID
func (p *DNSPodProvider) getDomainID(ctx context.Context, domain string) (int, error) {
	params := map[string]interface{}{
		"Limit": 3000,
	}
	
	response, err := p.makeRequest(ctx, "DescribeDomainList", params)
	if err != nil {
		return 0, err
	}
	
	var result struct {
		Response struct {
			DomainList []struct {
				DomainId int    `json:"DomainId"`
				Name     string `json:"Name"`
			} `json:"DomainList"`
		} `json:"Response"`
	}
	
	if err := json.Unmarshal(response, &result); err != nil {
		return 0, fmt.Errorf("解析响应失败: %v", err)
	}
	
	for _, domainInfo := range result.Response.DomainList {
		if domainInfo.Name == domain {
			return domainInfo.DomainId, nil
		}
	}
	
	return 0, fmt.Errorf("域名不存在: %s", domain)
}

// makeRequest 发起API请求
func (p *DNSPodProvider) makeRequest(ctx context.Context, action string, params map[string]interface{}) ([]byte, error) {
	// 构建请求体
	requestBody := map[string]interface{}{}
	for k, v := range params {
		requestBody[k] = v
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}
	
	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	
	// 设置请求头
	timestamp := time.Now().Unix()
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Host", "dnspod.tencentcloudapi.com")
	req.Header.Set("X-TC-Action", action)
	req.Header.Set("X-TC-Timestamp", strconv.FormatInt(timestamp, 10))
	req.Header.Set("X-TC-Version", "2021-03-23")
	req.Header.Set("X-TC-Region", "ap-beijing")
	
	// 生成签名
	signature := p.generateSignature(req, string(jsonBody), timestamp)
	req.Header.Set("Authorization", signature)
	
	// 发起请求
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
		Response struct {
			Error struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			} `json:"Error"`
			RequestId string `json:"RequestId"`
		} `json:"Response"`
	}
	
	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Response.Error.Code != "" {
		return nil, fmt.Errorf("腾讯云DNSPod API错误: %s - %s (RequestId: %s)", 
			errorResp.Response.Error.Code, errorResp.Response.Error.Message, errorResp.Response.RequestId)
	}
	
	return body, nil
}

// generateSignature 生成腾讯云API签名
func (p *DNSPodProvider) generateSignature(req *http.Request, payload string, timestamp int64) string {
	// 第一步：拼接规范请求串
	httpRequestMethod := req.Method
	canonicalURI := "/"
	canonicalQueryString := ""
	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s\n", 
		req.Header.Get("Content-Type"), req.Header.Get("Host"))
	signedHeaders := "content-type;host"
	hashedRequestPayload := sha256Hex(payload)
	
	canonicalRequest := httpRequestMethod + "\n" +
		canonicalURI + "\n" +
		canonicalQueryString + "\n" +
		canonicalHeaders + "\n" +
		signedHeaders + "\n" +
		hashedRequestPayload
		
	// 第二步：拼接待签名字符串
	algorithm := "TC3-HMAC-SHA256"
	service := "dnspod"
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := date + "/" + service + "/" + "tc3_request"
	hashedCanonicalRequest := sha256Hex(canonicalRequest)
	
	stringToSign := algorithm + "\n" +
		strconv.FormatInt(timestamp, 10) + "\n" +
		credentialScope + "\n" +
		hashedCanonicalRequest
		
	// 第三步：计算签名
	secretDate := hmacSha256([]byte("TC3"+p.config.APISecret), date)
	secretService := hmacSha256(secretDate, service)
	secretSigning := hmacSha256(secretService, "tc3_request")
	signature := hex.EncodeToString(hmacSha256(secretSigning, stringToSign))
	
	// 第四步：拼接Authorization
	authorization := algorithm + 
		" Credential=" + p.config.APIKey + "/" + credentialScope +
		", SignedHeaders=" + signedHeaders +
		", Signature=" + signature
		
	return authorization
}

// sha256Hex 计算SHA256哈希值并返回十六进制字符串
func sha256Hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

// hmacSha256 计算HMAC-SHA256
func hmacSha256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}