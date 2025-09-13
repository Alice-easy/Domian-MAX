package utils

import (
	"net"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

// ValidationService 验证服务
type ValidationService struct{}

// NewValidationService 创建验证服务
func NewValidationService() *ValidationService {
	return &ValidationService{}
}

// ValidateEmail 验证邮箱格式
func (v *ValidationService) ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// ValidatePassword 验证密码强度
func (v *ValidationService) ValidatePassword(password string) (bool, []string) {
	var errors []string
	
	if len(password) < 8 {
		errors = append(errors, "密码长度至少8位")
	}
	
	if len(password) > 128 {
		errors = append(errors, "密码长度不能超过128位")
	}
	
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	if !hasUpper {
		errors = append(errors, "密码必须包含大写字母")
	}
	if !hasLower {
		errors = append(errors, "密码必须包含小写字母")
	}
	if !hasNumber {
		errors = append(errors, "密码必须包含数字")
	}
	if !hasSpecial {
		errors = append(errors, "密码必须包含特殊字符")
	}
	
	return len(errors) == 0, errors
}

// ValidateUsername 验证用户名
func (v *ValidationService) ValidateUsername(username string) (bool, []string) {
	var errors []string
	
	if len(username) < 3 {
		errors = append(errors, "用户名长度至少3位")
	}
	
	if len(username) > 50 {
		errors = append(errors, "用户名长度不能超过50位")
	}
	
	// 只允许字母、数字、下划线和连字符
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
	if !matched {
		errors = append(errors, "用户名只能包含字母、数字、下划线和连字符")
	}
	
	// 不能以数字开头
	if len(username) > 0 && unicode.IsNumber(rune(username[0])) {
		errors = append(errors, "用户名不能以数字开头")
	}
	
	return len(errors) == 0, errors
}

// ValidateDomain 验证域名格式
func (v *ValidationService) ValidateDomain(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}
	
	// 移除末尾的点
	domain = strings.TrimSuffix(domain, ".")
	
	// 域名正则表达式
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)
	return domainRegex.MatchString(domain)
}

// ValidateIPv4 验证IPv4地址
func (v *ValidationService) ValidateIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

// ValidateIPv6 验证IPv6地址
func (v *ValidationService) ValidateIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() == nil
}

// ValidateDNSRecordType 验证DNS记录类型
func (v *ValidationService) ValidateDNSRecordType(recordType string) bool {
	validTypes := []string{"A", "AAAA", "CNAME", "MX", "TXT", "NS", "SRV", "PTR", "CAA"}
	for _, validType := range validTypes {
		if recordType == validType {
			return true
		}
	}
	return false
}

// ValidateDNSRecordValue 验证DNS记录值
func (v *ValidationService) ValidateDNSRecordValue(recordType, value string) (bool, string) {
	switch recordType {
	case "A":
		if !v.ValidateIPv4(value) {
			return false, "A记录值必须是有效的IPv4地址"
		}
	case "AAAA":
		if !v.ValidateIPv6(value) {
			return false, "AAAA记录值必须是有效的IPv6地址"
		}
	case "CNAME", "NS":
		if !v.ValidateDomain(value) {
			return false, recordType + "记录值必须是有效的域名"
		}
	case "MX":
		// MX记录格式: priority domain
		parts := strings.Fields(value)
		if len(parts) != 2 {
			return false, "MX记录值格式错误，应为: priority domain"
		}
		if !v.ValidateDomain(parts[1]) {
			return false, "MX记录的域名部分无效"
		}
	case "TXT":
		if len(value) > 255 {
			return false, "TXT记录值长度不能超过255字符"
		}
	case "SRV":
		// SRV记录格式: priority weight port target
		parts := strings.Fields(value)
		if len(parts) != 4 {
			return false, "SRV记录值格式错误，应为: priority weight port target"
		}
		if !v.ValidateDomain(parts[3]) {
			return false, "SRV记录的目标域名无效"
		}
	}
	
	return true, ""
}

// ValidateTTL 验证TTL值
func (v *ValidationService) ValidateTTL(ttl int) bool {
	return ttl >= 60 && ttl <= 86400 // 1分钟到1天
}

// SanitizeInput 清理输入字符串
func (v *ValidationService) SanitizeInput(input string) string {
	// 移除前后空白字符
	input = strings.TrimSpace(input)
	
	// 移除控制字符
	var result strings.Builder
	for _, r := range input {
		if !unicode.IsControl(r) {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}