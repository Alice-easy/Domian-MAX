package config

import (
	"domain-max/pkg/utils"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// 应用配置
	App struct {
		Port    string
		Mode    string // development, production
		BaseURL string
	}
	
	// 数据库配置
	Database struct {
		Type     string // postgres, mysql
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	
	// JWT配置
	JWT struct {
		Secret          string
		ExpirationHours int
	}
	
	// 安全配置
	Security struct {
		EncryptionKey string
	}
	
	// 邮件配置
	SMTP struct {
		Host     string
		Port     int
		User     string
		Password string
		From     string
	}
	
	// 兼容性字段（保持向后兼容）
	Port          string
	Environment   string
	BaseURL       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBType        string
	JWTSecret     string
	EncryptionKey string
	SMTPHost      string
	SMTPPort      int
	SMTPUser      string
	SMTPPassword  string
	SMTPFrom      string
	DNSPodToken   string
}

func Load() *Config {
	// 尝试加载.env文件
	godotenv.Load()

	cfg := &Config{}
	
	// 新的结构化配置
	cfg.App.Port = getEnv("PORT", "8080")
	cfg.App.Mode = getEnv("APP_MODE", "development")
	cfg.App.BaseURL = getEnv("BASE_URL", "")
	
	cfg.Database.Type = getEnv("DB_TYPE", "postgres")
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.Name = getEnv("DB_NAME", "domain_manager")
	
	cfg.JWT.Secret = getEnv("JWT_SECRET", "")
	cfg.JWT.ExpirationHours = getEnvInt("JWT_EXPIRATION_HOURS", 24)
	
	cfg.Security.EncryptionKey = getEnv("ENCRYPTION_KEY", "")
	
	cfg.SMTP.Host = getEnv("SMTP_HOST", "smtp.gmail.com")
	cfg.SMTP.Port = getEnvInt("SMTP_PORT", 587)
	cfg.SMTP.User = getEnv("SMTP_USER", "")
	cfg.SMTP.Password = getEnv("SMTP_PASSWORD", "")
	cfg.SMTP.From = getEnv("SMTP_FROM", "noreply@example.com")

	// 兼容性字段映射
	cfg.Port = cfg.App.Port
	cfg.Environment = cfg.App.Mode
	cfg.BaseURL = cfg.App.BaseURL
	cfg.DBHost = cfg.Database.Host
	cfg.DBPort = cfg.Database.Port
	cfg.DBUser = cfg.Database.User
	cfg.DBPassword = cfg.Database.Password
	cfg.DBName = cfg.Database.Name
	cfg.DBType = cfg.Database.Type
	cfg.JWTSecret = cfg.JWT.Secret
	cfg.EncryptionKey = cfg.Security.EncryptionKey
	cfg.SMTPHost = cfg.SMTP.Host
	cfg.SMTPPort = cfg.SMTP.Port
	cfg.SMTPUser = cfg.SMTP.User
	cfg.SMTPPassword = cfg.SMTP.Password
	cfg.SMTPFrom = cfg.SMTP.From
	cfg.DNSPodToken = getEnv("DNSPOD_TOKEN", "")

	// 如果没有设置BASE_URL，根据环境和端口自动生成
	if cfg.App.BaseURL == "" {
		if cfg.App.Mode == "development" {
			cfg.App.BaseURL = fmt.Sprintf("http://localhost:%s", cfg.App.Port)
		} else {
			cfg.App.BaseURL = fmt.Sprintf("https://localhost:%s", cfg.App.Port)
		}
		cfg.BaseURL = cfg.App.BaseURL
	}

	// 验证必要的配置项
	if err := cfg.validate(); err != nil {
		panic("配置验证失败: " + err.Error())
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// validate 验证配置项的有效性
func (c *Config) validate() error {
	isProduction := c.Environment == "production"
	
	// 验证端口
	if err := utils.ValidatePort(c.Port); err != nil {
		return fmt.Errorf("端口配置错误: %v", err)
	}
	
	// 验证必要的配置项
	requiredConfigs := map[string]string{
		"DB_PASSWORD":    c.DBPassword,
		"JWT_SECRET":     c.JWTSecret,
		"ENCRYPTION_KEY": c.EncryptionKey,
	}
	
	for key, value := range requiredConfigs {
		if value == "" {
			if isProduction {
				return fmt.Errorf("%s 不能为空，请设置相应的环境变量", key)
			} else {
				// 开发环境使用默认值
				switch key {
				case "DB_PASSWORD":
					c.DBPassword = "dev_password_123"
				case "JWT_SECRET":
					c.JWTSecret = "dev_jwt_secret_key_for_development_only_not_for_production_use_this_is_64_chars"
				case "ENCRYPTION_KEY":
					c.EncryptionKey = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
				}
				continue
			}
		}
		
		// 使用统一的配置验证
		if err := utils.ValidateConfigValue(key, value, isProduction); err != nil {
			return fmt.Errorf("%s 配置错误: %v", key, err)
		}
	}
	
	// 验证可选的SMTP配置
	if c.SMTPPassword != "" {
		if err := utils.ValidateConfigValue("SMTP_PASSWORD", c.SMTPPassword, isProduction); err != nil {
			return fmt.Errorf("SMTP_PASSWORD 配置错误: %v", err)
		}
	}
	
	// 生产环境额外安全检查
	if isProduction {
		if err := c.validateProductionSecurity(); err != nil {
			return err
		}
	}
	
	// 验证数据库类型
	validDBTypes := []string{"postgres", "mysql"}
	if !contains(validDBTypes, c.DBType) {
		return fmt.Errorf("不支持的数据库类型: %s，支持的类型: %s", c.DBType, strings.Join(validDBTypes, ", "))
	}
	
	return nil
}

// validateProductionSecurity 生产环境安全验证
func (c *Config) validateProductionSecurity() error {
	// 检查是否使用了明显的测试/默认值
	dangerousValues := []string{
		"test", "demo", "example", "sample", "default",
		"localhost", "127.0.0.1", "your_", "change_this",
		"password", "secret", "key", "token",
	}
	
	securityConfigs := map[string]string{
		"数据库密码":   c.DBPassword,
		"JWT密钥":   c.JWTSecret,
		"加密密钥":    c.EncryptionKey,
	}
	
	for configName, configValue := range securityConfigs {
		lowerValue := strings.ToLower(configValue)
		for _, dangerous := range dangerousValues {
			if strings.Contains(lowerValue, dangerous) {
				return fmt.Errorf("生产环境的%s包含不安全的词汇 '%s'，请使用随机生成的密钥", configName, dangerous)
			}
		}
	}
	
	// 检查BaseURL是否配置为生产域名
	if c.BaseURL != "" && strings.Contains(c.BaseURL, "localhost") {
		return errors.New("生产环境不能使用localhost作为BaseURL，请配置正确的域名")
	}
	
	return nil
}

// contains 检查切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}