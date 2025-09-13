package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTService JWT服务
type JWTService struct {
	secretKey   []byte
	expireHours int
}

// NewJWTService 创建JWT服务
func NewJWTService(secretKey string, expireHours int) *JWTService {
	return &JWTService{
		secretKey:   []byte(secretKey),
		expireHours: expireHours,
	}
}

// GenerateToken 生成JWT令牌
func (j *JWTService) GenerateToken(userID uint, username, email, role string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(j.expireHours) * time.Hour)
	
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "domain-max",
			Subject:   fmt.Sprintf("user:%d", userID),
			ID:        generateJTI(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	return tokenString, expiresAt, err
}

// ValidateToken 验证JWT令牌
func (j *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法是HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// RefreshToken 刷新JWT令牌
func (j *JWTService) RefreshToken(tokenString string) (string, time.Time, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", time.Time{}, err
	}

	// 检查令牌是否即将过期（剩余时间小于1小时）
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return "", time.Time{}, fmt.Errorf("token还未到刷新时间")
	}

	// 生成新令牌
	return j.GenerateToken(claims.UserID, claims.Username, claims.Email, claims.Role)
}

// generateJTI 生成JWT ID
func generateJTI() string {
	// 使用crypto/rand生成更安全的随机数
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// 如果随机数生成失败，使用时间戳作为备用
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%d_%s", time.Now().UnixNano(), n.String())
}

// PasswordService 密码服务
type PasswordService struct{}

// NewPasswordService 创建密码服务
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// HashPassword 哈希密码
func (p *PasswordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func (p *PasswordService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// VerifyPassword 验证密码（CheckPassword的别名，为了API兼容性）
func (p *PasswordService) VerifyPassword(password, hash string) bool {
	return p.CheckPassword(password, hash)
}

// ValidatePasswordStrength 验证密码强度
func (p *PasswordService) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度至少8位")
	}

	if len(password) > 100 {
		return errors.New("密码长度不能超过100位")
	}

	// 检查是否包含字母
	hasLetter, _ := regexp.MatchString(`[a-zA-Z]`, password)
	if !hasLetter {
		return errors.New("密码必须包含字母")
	}

	// 检查是否包含数字
	hasNumber, _ := regexp.MatchString(`[0-9]`, password)
	if !hasNumber {
		return errors.New("密码必须包含数字")
	}

	// 检查是否包含特殊字符
	hasSpecial, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]`, password)

	// 如果密码长度小于12，则必须包含特殊字符
	if len(password) < 12 && !hasSpecial {
		return errors.New("密码长度小于12位时必须包含特殊字符")
	}

	return nil
}

// EncryptionService 加密服务
type EncryptionService struct {
	key []byte
}

// NewEncryptionService 创建加密服务
func NewEncryptionService(key string) (*EncryptionService, error) {
	// 确保密钥长度为32字节（AES-256）
	keyBytes := sha256.Sum256([]byte(key))
	return &EncryptionService{key: keyBytes[:]}, nil
}

// Encrypt 加密数据
func (e *EncryptionService) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密数据
func (e *EncryptionService) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptJSON 加密JSON格式的map数据
func (e *EncryptionService) EncryptJSON(data map[string]string) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return e.Encrypt(string(jsonBytes))
}

// DecryptJSON 解密JSON格式的数据到map
func (e *EncryptionService) DecryptJSON(ciphertext string) (map[string]string, error) {
	plaintext, err := e.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}
	
	var data map[string]string
	err = json.Unmarshal([]byte(plaintext), &data)
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

// ValidationService 验证服务
type ValidationService struct{}

// NewValidationService 创建验证服务
func NewValidationService() *ValidationService {
	return &ValidationService{}
}

// ValidateStruct 验证结构体
func (v *ValidationService) ValidateStruct(s interface{}) error {
	// 这里可以使用反射或第三方库进行复杂验证
	// 为了简化，目前返回nil，实际使用时应该实现具体的验证逻辑
	return nil
}

// ValidateEmail 验证邮箱格式
func (v *ValidationService) ValidateEmail(email string) error {
	if len(email) == 0 {
		return errors.New("邮箱不能为空")
	}

	if len(email) > 255 {
		return errors.New("邮箱长度不能超过255位")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("邮箱格式不正确")
	}

	// 检查危险字符
	dangerousChars := []string{"<", ">", "\"", "'", "&", ";", "|", "`", "$"}
	for _, char := range dangerousChars {
		if strings.Contains(email, char) {
			return errors.New("邮箱包含不安全字符")
		}
	}

	return nil
}

// ValidateDomainName 验证域名格式
func (v *ValidationService) ValidateDomainName(domain string) error {
	if len(domain) == 0 {
		return errors.New("域名不能为空")
	}

	if len(domain) > 253 {
		return errors.New("域名长度不能超过253个字符")
	}

	domainPattern := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-\.]*[a-zA-Z0-9])?$`)
	if !domainPattern.MatchString(domain) {
		return errors.New("域名格式不正确")
	}

	// 验证每个标签
	labels := strings.Split(domain, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return errors.New("域名标签长度必须在1-63个字符之间")
		}
		if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return errors.New("域名标签不能以连字符开头或结尾")
		}
	}

	return nil
}

// ValidatePort 验证端口号
func ValidatePort(port string) error {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return errors.New("端口必须是数字")
	}

	if portNum < 1 || portNum > 65535 {
		return errors.New("端口范围必须在1-65535之间")
	}

	return nil
}

// ValidateConfigValue 验证配置值
func ValidateConfigValue(key, value string, isProduction bool) error {
	switch key {
	case "JWT_SECRET":
		if len(value) < 32 {
			return errors.New("JWT密钥长度至少32字符")
		}
		if isProduction && containsCommonWords(value) {
			return errors.New("生产环境不能使用常见词汇作为JWT密钥")
		}

	case "ENCRYPTION_KEY":
		if len(value) < 32 {
			return errors.New("加密密钥长度至少32字符")
		}
		if isProduction && containsCommonWords(value) {
			return errors.New("生产环境不能使用常见词汇作为加密密钥")
		}

	case "DB_PASSWORD":
		if len(value) < 8 {
			return errors.New("数据库密码长度至少8字符")
		}
		if isProduction && containsCommonWords(value) {
			return errors.New("生产环境不能使用常见词汇作为数据库密码")
		}

	case "SMTP_PASSWORD":
		if len(value) < 6 {
			return errors.New("SMTP密码长度至少6字符")
		}
	}

	return nil
}

// containsCommonWords 检查是否包含常见词汇
func containsCommonWords(value string) bool {
	lowerValue := strings.ToLower(value)
	commonWords := []string{
		"password", "secret", "key", "token", "admin", "user",
		"test", "demo", "example", "sample", "default",
		"123456", "qwerty", "abc123", "password123",
	}

	for _, word := range commonWords {
		if strings.Contains(lowerValue, word) {
			return true
		}
	}

	return false
}

// SanitizeInput 清理用户输入
func SanitizeInput(input string) string {
	// 移除可能的脚本标签
	scriptPattern := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	input = scriptPattern.ReplaceAllString(input, "")

	// 移除其他潜在危险的HTML标签
	dangerousTags := []string{
		`(?i)<iframe[^>]*>.*?</iframe>`,
		`(?i)<object[^>]*>.*?</object>`,
		`(?i)<embed[^>]*>.*?</embed>`,
		`(?i)<link[^>]*>`,
		`(?i)<meta[^>]*>`,
	}

	for _, pattern := range dangerousTags {
		re := regexp.MustCompile(pattern)
		input = re.ReplaceAllString(input, "")
	}

	// 移除潜在的事件处理器
	eventPattern := regexp.MustCompile(`(?i)on\w+\s*=\s*["\'][^"\']*["\']`)
	input = eventPattern.ReplaceAllString(input, "")

	// 移除javascript:伪协议
	jsPattern := regexp.MustCompile(`(?i)javascript\s*:`)
	input = jsPattern.ReplaceAllString(input, "")

	return strings.TrimSpace(input)
}