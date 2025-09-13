package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// Config 应用配置结构
type Config struct {
	// 服务器配置
	Port        string `json:"port"`
	Environment string `json:"environment"`
	
	// 数据库配置
	DatabaseURL      string `json:"database_url"`
	DatabaseType     string `json:"database_type"`
	DatabaseHost     string `json:"database_host"`
	DatabasePort     string `json:"database_port"`
	DatabaseName     string `json:"database_name"`
	DatabaseUser     string `json:"database_user"`
	DatabasePassword string `json:"database_password"`
	
	// JWT配置
	JWTSecret          string `json:"jwt_secret"`
	JWTExpirationHours int    `json:"jwt_expiration_hours"`
	
	// 加密配置
	EncryptionKey string `json:"encryption_key"`
	
	// CORS配置
	AllowedOrigins []string `json:"allowed_origins"`
	
	// 邮件配置
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	SMTPFrom     string `json:"smtp_from"`
}

// Load 加载配置
func Load() *Config {
	config := &Config{
		Port:               getEnv("PORT", "8080"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		DatabaseType:       getEnv("DATABASE_TYPE", "sqlite"),
		DatabaseHost:       getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:       getEnv("DATABASE_PORT", "3306"),
		DatabaseName:       getEnv("DATABASE_NAME", "domain_max"),
		DatabaseUser:       getEnv("DATABASE_USER", "root"),
		DatabasePassword:   getEnv("DATABASE_PASSWORD", ""),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		EncryptionKey:      getEnv("ENCRYPTION_KEY", "your-32-char-encryption-key-here"),
		AllowedOrigins:     getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		SMTPHost:           getEnv("SMTP_HOST", ""),
		SMTPPort:           getEnvAsInt("SMTP_PORT", 587),
		SMTPUsername:       getEnv("SMTP_USERNAME", ""),
		SMTPPassword:       getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:           getEnv("SMTP_FROM", ""),
	}

	// 验证必需的配置
	if config.JWTSecret == "your-secret-key" {
		log.Println("警告: 使用默认JWT密钥，生产环境请设置JWT_SECRET环境变量")
	}

	return config
}

// IsProduction 判断是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment 判断是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsSlice 获取环境变量并转换为字符串切片
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}