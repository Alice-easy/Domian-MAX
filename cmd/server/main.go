package main

import (
	"domain-max/pkg/api"
	"domain-max/pkg/config"
	"domain-max/pkg/database"
	"domain-max/pkg/dns/providers"
	"domain-max/pkg/middleware"
	"domain-max/pkg/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 连接数据库
	var db *gorm.DB
	var err error
	
	if cfg.Environment == "development" {
		log.Println("开发环境：连接数据库")
	}
	
	db, err = database.Connect(cfg)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 运行数据库迁移
	if err := database.Migrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化安全服务
	jwtService := utils.NewJWTService(cfg.JWTSecret, 24) // 24小时过期
	passwordService := utils.NewPasswordService()
	encryptionService, err := utils.NewEncryptionService(cfg.EncryptionKey)
	if err != nil {
		log.Fatalf("加密服务初始化失败: %v", err)
	}
	validationService := utils.NewValidationService()

	// 初始化DNS提供商工厂
	providerFactory := providers.NewProviderFactory()

	// 初始化API控制器
	authAPI := api.NewAuthAPI(db, jwtService, passwordService, validationService)
	dnsAPI := api.NewSimpleDNSAPI(db, providerFactory, encryptionService, validationService)

	// 设置Gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// 添加中间件
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RateLimitMiddleware())

	// API路由
	setupAPIRoutes(router, authAPI, dnsAPI, jwtService)

	// 前端静态文件服务
	setupWebRoutes(router)

	log.Printf("服务器启动在端口 %s", cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}

func setupWebRoutes(router *gin.Engine) {
	// 检查web/dist目录是否存在
	webDistPath := "web/dist"
	if _, err := os.Stat(webDistPath); os.IsNotExist(err) {
		log.Printf("警告: web/dist 目录不存在，跳过静态文件服务")
		return
	}

	// 静态文件服务 - 处理构建后的静态资源
	router.Static("/static", filepath.Join(webDistPath, "static"))

	// 处理所有其他路由，返回index.html (用于React Router)
	router.NoRoute(func(c *gin.Context) {
		// 如果是API请求，返回404
		if len(c.Request.URL.Path) > 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}

		// 尝试读取index.html
		indexPath := filepath.Join(webDistPath, "index.html")
		indexHTML, err := os.ReadFile(indexPath)
		if err != nil {
			c.String(500, "无法加载前端页面")
			return
		}

		c.Data(200, "text/html; charset=utf-8", indexHTML)
	})
}

func setupAPIRoutes(router *gin.Engine, authAPI *api.AuthAPI, dnsAPI *api.SimpleDNSAPI, jwtService *utils.JWTService) {
	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
			"version": "1.0.0",
		})
	})

	// API版本1
	v1 := router.Group("/api/v1")

	// 认证路由（无需认证）
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authAPI.Login)
		auth.POST("/register", authAPI.Register)
		auth.POST("/logout", authAPI.Logout)
	}

	// 需要认证的路由
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		// 用户相关
		user := protected.Group("/user")
		{
			user.GET("/profile", authAPI.GetProfile)
			user.POST("/change-password", authAPI.ChangePassword)
			user.POST("/refresh-token", authAPI.RefreshToken)
		}

		// DNS提供商管理
		providers := protected.Group("/dns-providers")
		{
			providers.GET("", dnsAPI.GetDNSProviders)
			providers.POST("", dnsAPI.CreateDNSProvider)
			providers.POST("/:id/test", dnsAPI.TestDNSProvider)
			providers.GET("/supported", dnsAPI.ListSupportedProviders)
		}

		// 域名管理
		domains := protected.Group("/domains")
		{
			domains.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "域名管理功能待实现",
					"todo":    "实现域名列表、添加域名、域名配置等功能",
				})
			})
		}

		// 管理员路由
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminRequiredMiddleware())
		{
			// 用户管理
			adminUsers := admin.Group("/users")
			{
				adminUsers.GET("", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "用户管理功能待实现",
						"todo":    "实现用户列表、用户详情、用户禁用等功能",
					})
				})
			}

			// 系统管理
			system := admin.Group("/system")
			{
				system.GET("/stats", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{
						"message": "系统统计功能待实现",
						"todo":    "实现用户统计、域名统计、DNS记录统计等",
					})
				})
			}
		}
	}
}