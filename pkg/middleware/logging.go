package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware 日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	})
}

// RequestResponseLogger 请求响应日志中间件
func RequestResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 创建响应写入器
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)

		// 构建日志信息
		logData := map[string]interface{}{
			"timestamp":    start.Format(time.RFC3339),
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"query":        c.Request.URL.RawQuery,
			"status_code":  c.Writer.Status(),
			"latency":      latency.String(),
			"client_ip":    c.ClientIP(),
			"user_agent":   c.Request.UserAgent(),
			"request_size": c.Request.ContentLength,
		}

		// 添加用户信息（如果已认证）
		if userID, exists := c.Get("user_id"); exists {
			logData["user_id"] = userID
		}
		if username, exists := c.Get("username"); exists {
			logData["username"] = username
		}

		// 记录请求体（仅非敏感接口）
		if shouldLogRequestBody(c.Request.URL.Path) && len(requestBody) > 0 {
			logData["request_body"] = string(requestBody)
		}

		// 记录响应体（仅开发环境）
		if gin.Mode() == gin.DebugMode && shouldLogResponseBody(c.Request.URL.Path) {
			logData["response_body"] = blw.body.String()
		}

		// 输出日志
		logJSON, _ := json.Marshal(logData)
		log.Printf("API_LOG: %s", string(logJSON))
	}
}

// bodyLogWriter 响应体日志写入器
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// shouldLogRequestBody 判断是否应该记录请求体
func shouldLogRequestBody(path string) bool {
	// 不记录敏感接口的请求体
	sensitiveEndpoints := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/auth/change-password",
	}

	for _, endpoint := range sensitiveEndpoints {
		if path == endpoint {
			return false
		}
	}
	return true
}

// shouldLogResponseBody 判断是否应该记录响应体
func shouldLogResponseBody(path string) bool {
	// 不记录敏感接口的响应体
	sensitiveEndpoints := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/refresh",
	}

	for _, endpoint := range sensitiveEndpoints {
		if path == endpoint {
			return false
		}
	}
	return true
}