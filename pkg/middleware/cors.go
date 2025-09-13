package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins []string
	IsDevelopment  bool
}

// CORSWithConfig 带配置的CORS中间件
func CORSWithConfig(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// 检查是否允许该来源
		allowed := false
		if config.IsDevelopment {
			// 开发环境允许所有localhost和127.0.0.1
			if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
				allowed = true
			}
		}
		
		// 检查配置的允许来源
		for _, allowedOrigin := range config.AllowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")
		
		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// DefaultCORS 默认CORS中间件
func DefaultCORS() gin.HandlerFunc {
	return CORSWithConfig(CORSConfig{
		AllowedOrigins: []string{"*"},
		IsDevelopment:  true,
	})
}