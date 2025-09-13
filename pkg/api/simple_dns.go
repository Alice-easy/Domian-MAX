package api

import (
	"context"
	"domain-max/pkg/dns/models"
	"domain-max/pkg/dns/providers"
	"domain-max/pkg/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SimpleDNSAPI 简化的DNS管理API控制器
type SimpleDNSAPI struct {
	DB               *gorm.DB
	ProviderFactory  *providers.ProviderFactory
	EncryptionService *utils.EncryptionService
	Validator        *utils.ValidationService
}

// NewSimpleDNSAPI 创建简化DNS API实例
func NewSimpleDNSAPI(db *gorm.DB, factory *providers.ProviderFactory, encService *utils.EncryptionService, validator *utils.ValidationService) *SimpleDNSAPI {
	return &SimpleDNSAPI{
		DB:               db,
		ProviderFactory:  factory,
		EncryptionService: encService,
		Validator:        validator,
	}
}

// DNSProviderRequest DNS提供商配置请求
type DNSProviderRequest struct {
	Name        string            `json:"name" binding:"required"`
	Type        string            `json:"type" binding:"required"`
	Config      map[string]string `json:"config" binding:"required"`
	Description string            `json:"description,omitempty"`
	IsDefault   bool              `json:"is_default,omitempty"`
}

// GetDNSProviders 获取DNS提供商列表
func (d *SimpleDNSAPI) GetDNSProviders(c *gin.Context) {
	userID, _, _, role, ok := getUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未找到用户信息",
			"code":  "USER_NOT_FOUND",
		})
		return
	}

	var providers []models.DNSProvider
	query := d.DB.Model(&models.DNSProvider{})

	// 非管理员只能看到自己的和默认的提供商
	if role != "admin" {
		query = query.Where("user_id = ? OR is_default = ?", userID, true)
	}

	if err := query.Find(&providers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取DNS提供商失败",
			"code":    "PROVIDER_FETCH_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    providers,
	})
}

// CreateDNSProvider 创建DNS提供商配置
func (d *SimpleDNSAPI) CreateDNSProvider(c *gin.Context) {
	userID, _, _, role, ok := getUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未找到用户信息",
			"code":  "USER_NOT_FOUND",
		})
		return
	}

	var req DNSProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		})
		return
	}

	// 检查提供商类型是否支持
	if !d.ProviderFactory.IsSupported(req.Type) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "不支持的DNS提供商类型",
			"code":    "UNSUPPORTED_PROVIDER",
			"message": "当前不支持该DNS提供商",
		})
		return
	}

	// 验证配置的有效性
	provider, err := d.ProviderFactory.CreateProvider(req.Type, req.Config)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "DNS提供商配置无效",
			"code":    "INVALID_PROVIDER_CONFIG",
			"message": err.Error(),
		})
		return
	}

	// 测试连接
	ctx := context.Background()
	if err := provider.TestConnection(ctx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "DNS提供商连接测试失败",
			"code":    "PROVIDER_CONNECTION_ERROR",
			"message": err.Error(),
		})
		return
	}

	// 加密配置信息
	encryptedConfig, err := d.EncryptionService.EncryptJSON(req.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "配置加密失败",
			"code":    "CONFIG_ENCRYPTION_ERROR",
			"message": "服务器内部错误",
		})
		return
	}

	// 创建DNS提供商配置
	dnsProvider := models.DNSProvider{
		UserID:      userID,
		Name:        req.Name,
		Type:        req.Type,
		Config:      encryptedConfig,
		Description: req.Description,
		IsDefault:   req.IsDefault && role == "admin", // 只有管理员可以设置默认
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := d.DB.Create(&dnsProvider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "DNS提供商创建失败",
			"code":    "PROVIDER_CREATION_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "DNS提供商创建成功",
		"data":    dnsProvider,
	})
}

// TestDNSProvider 测试DNS提供商连接
func (d *SimpleDNSAPI) TestDNSProvider(c *gin.Context) {
	userID, _, _, role, ok := getUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未找到用户信息",
			"code":  "USER_NOT_FOUND",
		})
		return
	}

	providerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "无效的提供商ID",
			"code":    "INVALID_PROVIDER_ID",
			"message": "提供商ID必须是数字",
		})
		return
	}

	// 查找DNS提供商
	var dnsProvider models.DNSProvider
	query := d.DB.Where("id = ?", providerID)
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&dnsProvider).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "DNS提供商不存在",
			"code":    "PROVIDER_NOT_FOUND",
			"message": "未找到指定的DNS提供商",
		})
		return
	}

	// 解密配置
	config, err := d.EncryptionService.DecryptJSON(dnsProvider.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "配置解密失败",
			"code":    "CONFIG_DECRYPT_ERROR",
			"message": "服务器内部错误",
		})
		return
	}

	// 创建提供商实例并测试连接
	provider, err := d.ProviderFactory.CreateProvider(dnsProvider.Type, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "创建提供商实例失败",
			"code":    "PROVIDER_CREATE_ERROR",
			"message": err.Error(),
		})
		return
	}

	ctx := context.Background()
	if err := provider.TestConnection(ctx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "连接测试失败",
			"code":    "CONNECTION_TEST_FAILED",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "连接测试成功",
	})
}

// ListSupportedProviders 获取支持的DNS提供商类型
func (d *SimpleDNSAPI) ListSupportedProviders(c *gin.Context) {
	supportedTypes := d.ProviderFactory.GetSupportedTypes()
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    supportedTypes,
	})
}