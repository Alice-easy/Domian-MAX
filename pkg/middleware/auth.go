package middleware

import (
	"domain-max/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(jwtService *utils.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "未提供认证令牌",
				"code":    "MISSING_TOKEN",
				"message": "请在请求头中包含Authorization字段",
			})
			c.Abort()
			return
		}

		// 检查Bearer格式
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "认证令牌格式错误",
				"code":    "INVALID_TOKEN_FORMAT",
				"message": "Authorization头必须以'Bearer '开头",
			})
			c.Abort()
			return
		}

		// 提取令牌
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "认证令牌为空",
				"code":    "EMPTY_TOKEN",
				"message": "请提供有效的JWT令牌",
			})
			c.Abort()
			return
		}

		// 验证令牌
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "认证令牌无效",
				"code":    "INVALID_TOKEN",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// AdminRequiredMiddleware 管理员权限中间件
func AdminRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "未找到用户角色信息",
				"code":    "MISSING_ROLE",
				"message": "请确保已通过认证中间件",
			})
			c.Abort()
			return
		}

		if roleStr, ok := role.(string); !ok || roleStr != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "权限不足",
				"code":    "INSUFFICIENT_PRIVILEGES",
				"message": "此操作需要管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件（不强制要求认证）
func OptionalAuthMiddleware(jwtService *utils.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		const bearerPrefix = "Bearer "
		if strings.HasPrefix(authHeader, bearerPrefix) {
			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
			if tokenString != "" {
				if claims, err := jwtService.ValidateToken(tokenString); err == nil {
					c.Set("user_id", claims.UserID)
					c.Set("username", claims.Username)
					c.Set("email", claims.Email)
					c.Set("role", claims.Role)
					c.Set("jwt_claims", claims)
				}
			}
		}

		c.Next()
	}
}

// GetCurrentUserID 从上下文获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	if id, ok := userID.(uint); ok {
		return id, true
	}

	return 0, false
}

// GetCurrentUser 从上下文获取当前用户信息
func GetCurrentUser(c *gin.Context) (uint, string, string, string, bool) {
	userID, userIDExists := c.Get("user_id")
	username, usernameExists := c.Get("username")
	email, emailExists := c.Get("email")
	role, roleExists := c.Get("role")

	if !userIDExists || !usernameExists || !emailExists || !roleExists {
		return 0, "", "", "", false
	}

	if id, ok := userID.(uint); ok {
		if uname, ok := username.(string); ok {
			if em, ok := email.(string); ok {
				if r, ok := role.(string); ok {
					return id, uname, em, r, true
				}
			}
		}
	}

	return 0, "", "", "", false
}

// IsAdmin 检查当前用户是否为管理员
func IsAdmin(c *gin.Context) bool {
	role, exists := c.Get("role")
	if !exists {
		return false
	}

	if roleStr, ok := role.(string); ok {
		return roleStr == "admin"
	}

	return false
}