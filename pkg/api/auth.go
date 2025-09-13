package api

import (
	"domain-max/pkg/auth/models"
	"domain-max/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthAPI 认证API控制器
type AuthAPI struct {
	DB          *gorm.DB
	JWTService  *utils.JWTService
	PassService *utils.PasswordService
	Validator   *utils.ValidationService
}

// NewAuthAPI 创建认证API实例
func NewAuthAPI(db *gorm.DB, jwtService *utils.JWTService, passService *utils.PasswordService, validator *utils.ValidationService) *AuthAPI {
	return &AuthAPI{
		DB:          db,
		JWTService:  jwtService,
		PassService: passService,
		Validator:   validator,
	}
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=6" validate:"required,min=6"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=50" validate:"required,min=3,max=50"`
	Email           string `json:"email" binding:"required,email" validate:"required,email"`
	Password        string `json:"password" binding:"required,min=6" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"required"`
	InviteCode      string `json:"invite_code,omitempty" validate:"omitempty"`
}

// ChangePasswordRequest 修改密码请求结构
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" validate:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"required"`
}

// AuthResponse 认证响应结构
type AuthResponse struct {
	Token     string      `json:"token"`
	ExpiresAt int64       `json:"expires_at"`
	User      UserProfile `json:"user"`
}

// UserProfile 用户资料结构
type UserProfile struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

// Login 用户登录
func (a *AuthAPI) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		})
		return
	}

	// 验证请求数据
	if err := a.Validator.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "数据验证失败",
			"code":    "VALIDATION_ERROR",
			"message": err.Error(),
		})
		return
	}

	// 查找用户
	var user models.User
	if err := a.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "邮箱或密码错误",
			"code":    "INVALID_CREDENTIALS",
			"message": "请检查您的登录信息",
		})
		return
	}

	// 检查用户状态
	if user.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "账户已被禁用",
			"code":    "ACCOUNT_DISABLED",
			"message": "请联系管理员激活账户",
		})
		return
	}

	// 验证密码
	if !a.PassService.VerifyPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "邮箱或密码错误",
			"code":    "INVALID_CREDENTIALS",
			"message": "请检查您的登录信息",
		})
		return
	}

	// 更新最后登录时间
	now := time.Now()
	a.DB.Model(&user).Update("last_login_at", now)

	// 生成JWT令牌
	token, expiresAt, err := a.JWTService.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "令牌生成失败",
			"code":    "TOKEN_GENERATION_ERROR",
			"message": "服务器内部错误",
		})
		return
	}

	// 返回认证响应
	response := AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt.Unix(),
		User: UserProfile{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Status:   user.Status,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"data":    response,
	})
}

// Register 用户注册
func (a *AuthAPI) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		})
		return
	}

	// 验证请求数据
	if err := a.Validator.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "数据验证失败",
			"code":    "VALIDATION_ERROR",
			"message": err.Error(),
		})
		return
	}

	// 验证密码确认
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "密码不匹配",
			"code":    "PASSWORD_MISMATCH",
			"message": "两次输入的密码不一致",
		})
		return
	}

	// 检查邮箱是否已存在
	var existingUser models.User
	if err := a.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "邮箱已被使用",
			"code":    "EMAIL_EXISTS",
			"message": "该邮箱已被其他用户注册",
		})
		return
	}

	// 检查用户名是否已存在
	if err := a.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "用户名已被使用",
			"code":    "USERNAME_EXISTS",
			"message": "该用户名已被其他用户使用",
		})
		return
	}

	// 验证邀请码（如果需要）
	// TODO: 实现邀请码验证逻辑

	// 加密密码
	hashedPassword, err := a.PassService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "密码加密失败",
			"code":    "PASSWORD_HASH_ERROR",
			"message": "服务器内部错误",
		})
		return
	}

	// 创建新用户
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      "user", // 默认角色
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "用户创建失败",
			"code":    "USER_CREATION_ERROR",
			"message": "服务器内部错误",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "注册成功",
		"data": UserProfile{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Status:   user.Status,
		},
	})
}

// GetProfile 获取用户资料
func (a *AuthAPI) GetProfile(c *gin.Context) {
	userID, _, _, _, ok := getUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "未找到用户信息",
			"code":    "USER_NOT_FOUND",
			"message": "请重新登录",
		})
		return
	}

	// 从数据库获取最新用户信息
	var user models.User
	if err := a.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "用户不存在",
			"code":    "USER_NOT_FOUND",
			"message": "用户信息已被删除",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": UserProfile{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Status:   user.Status,
		},
	})
}

// ChangePassword 修改密码
func (a *AuthAPI) ChangePassword(c *gin.Context) {
	userID, _, _, _, ok := getUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "未找到用户信息",
			"code":    "USER_NOT_FOUND",
			"message": "请重新登录",
		})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		})
		return
	}

	// 验证请求数据
	if err := a.Validator.ValidateStruct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "数据验证失败",
			"code":    "VALIDATION_ERROR",
			"message": err.Error(),
		})
		return
	}

	// 验证新密码确认
	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "新密码不匹配",
			"code":    "PASSWORD_MISMATCH",
			"message": "两次输入的新密码不一致",
		})
		return
	}

	// 获取用户信息
	var user models.User
	if err := a.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "用户不存在",
			"code":    "USER_NOT_FOUND",
			"message": "用户信息已被删除",
		})
		return
	}

	// 验证当前密码
	if !a.PassService.VerifyPassword(req.CurrentPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "当前密码错误",
			"code":    "INVALID_CURRENT_PASSWORD",
			"message": "请输入正确的当前密码",
		})
		return
	}

	// 加密新密码
	newHashedPassword, err := a.PassService.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "密码加密失败",
			"code":    "PASSWORD_HASH_ERROR",
			"message": "服务器内部错误",
		})
		return
	}

	// 更新密码
	if err := a.DB.Model(&user).Update("password", newHashedPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "密码更新失败",
			"code":    "PASSWORD_UPDATE_ERROR",
			"message": "服务器内部错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "密码修改成功",
	})
}

// Logout 用户登出
func (a *AuthAPI) Logout(c *gin.Context) {
	// JWT是无状态的，这里只是简单的成功响应
	// 实际的令牌失效应该在客户端处理
	// 如果需要服务端令牌黑名单，可以在这里实现
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登出成功",
	})
}

// RefreshToken 刷新令牌
func (a *AuthAPI) RefreshToken(c *gin.Context) {
	userID, _, _, _, ok := getUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "未找到用户信息",
			"code":    "USER_NOT_FOUND",
			"message": "请重新登录",
		})
		return
	}

	// 验证用户是否仍然存在且处于活跃状态
	var user models.User
	if err := a.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "用户不存在",
			"code":    "USER_NOT_FOUND",
			"message": "用户信息已被删除",
		})
		return
	}

	if user.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "账户已被禁用",
			"code":    "ACCOUNT_DISABLED",
			"message": "请联系管理员激活账户",
		})
		return
	}

	// 生成新的JWT令牌
	token, expiresAt, err := a.JWTService.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "令牌生成失败",
			"code":    "TOKEN_GENERATION_ERROR",
			"message": "服务器内部错误",
		})
		return
	}

	// 返回新的令牌
	response := AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt.Unix(),
		User: UserProfile{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Status:   user.Status,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "令牌刷新成功",
		"data":    response,
	})
}

// getUserFromContext 从Gin上下文获取用户信息
func getUserFromContext(c *gin.Context) (uint, string, string, string, bool) {
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