package api

import (
	"domain-max/pkg/database"
	"domain-max/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthAPI 认证API控制器
type AuthAPI struct {
	db                *gorm.DB
	jwtService        *utils.JWTService
	passwordService   *utils.PasswordService
	validationService *utils.ValidationService
}

// NewAuthAPI 创建认证API控制器
func NewAuthAPI(db *gorm.DB, jwtService *utils.JWTService, passwordService *utils.PasswordService, validationService *utils.ValidationService) *AuthAPI {
	return &AuthAPI{
		db:                db,
		jwtService:        jwtService,
		passwordService:   passwordService,
		validationService: validationService,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	Token string               `json:"token"`
	User  *database.User       `json:"user"`
}

// Login 用户登录
func (a *AuthAPI) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
		})
		return
	}

	// 清理输入
	req.Username = a.validationService.SanitizeInput(req.Username)
	req.Password = a.validationService.SanitizeInput(req.Password)

	// 查找用户
	var user database.User
	if err := a.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户名或密码错误",
		})
		return
	}

	// 检查用户是否激活
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "账户已被禁用",
		})
		return
	}

	// 验证密码
	valid, err := a.passwordService.VerifyPassword(req.Password, user.Password)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "用户名或密码错误",
		})
		return
	}

	// 生成JWT令牌
	token, err := a.jwtService.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "令牌生成失败",
		})
		return
	}

	// 清除密码字段
	user.Password = ""

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  &user,
	})
}

// Register 用户注册
func (a *AuthAPI) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
		})
		return
	}

	// 清理输入
	req.Username = a.validationService.SanitizeInput(req.Username)
	req.Email = a.validationService.SanitizeInput(req.Email)
	req.Password = a.validationService.SanitizeInput(req.Password)

	// 验证用户名
	if valid, errors := a.validationService.ValidateUsername(req.Username); !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "用户名格式错误",
			"details": errors,
		})
		return
	}

	// 验证邮箱
	if !a.validationService.ValidateEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "邮箱格式错误",
		})
		return
	}

	// 验证密码
	if valid, errors := a.validationService.ValidatePassword(req.Password); !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "密码格式错误",
			"details": errors,
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser database.User
	if err := a.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "用户名已存在",
		})
		return
	}

	// 检查邮箱是否已存在
	if err := a.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "邮箱已存在",
		})
		return
	}

	// 哈希密码
	hashedPassword, err := a.passwordService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "密码加密失败",
		})
		return
	}

	// 创建用户
	user := database.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
		IsActive: true,
	}

	if err := a.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用户创建失败",
		})
		return
	}

	// 生成JWT令牌
	token, err := a.jwtService.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "令牌生成失败",
		})
		return
	}

	// 清除密码字段
	user.Password = ""

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  &user,
	})
}

// GetProfile 获取用户资料
func (a *AuthAPI) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未找到用户信息",
		})
		return
	}

	var user database.User
	if err := a.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 清除密码字段
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// ChangePassword 修改密码
func (a *AuthAPI) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未找到用户信息",
		})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
		})
		return
	}

	// 清理输入
	req.OldPassword = a.validationService.SanitizeInput(req.OldPassword)
	req.NewPassword = a.validationService.SanitizeInput(req.NewPassword)

	// 验证新密码
	if valid, errors := a.validationService.ValidatePassword(req.NewPassword); !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "新密码格式错误",
			"details": errors,
		})
		return
	}

	// 查找用户
	var user database.User
	if err := a.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	// 验证旧密码
	valid, err := a.passwordService.VerifyPassword(req.OldPassword, user.Password)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "旧密码错误",
		})
		return
	}

	// 哈希新密码
	hashedPassword, err := a.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "密码加密失败",
		})
		return
	}

	// 更新密码
	if err := a.db.Model(&user).Update("password", hashedPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "密码更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "密码修改成功",
	})
}

// RefreshToken 刷新令牌
func (a *AuthAPI) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "缺少认证令牌",
		})
		return
	}

	// 提取令牌
	tokenParts := strings.SplitN(authHeader, " ", 2)
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "无效的认证令牌格式",
		})
		return
	}

	token := tokenParts[1]
	newToken, err := a.jwtService.RefreshToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "令牌刷新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}

// Logout 用户登出
func (a *AuthAPI) Logout(c *gin.Context) {
	// 在实际应用中，可以将令牌加入黑名单
	// 这里简单返回成功消息
	c.JSON(http.StatusOK, gin.H{
		"message": "登出成功",
	})
}