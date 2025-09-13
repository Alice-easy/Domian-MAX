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

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 连接数据库（支持远程数据库）
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

	// 添加增强的CORS中间件（支持Cloudflare Pages）
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowedOrigins: cfg.AllowedOrigins, // 从配置文件读取
		IsDevelopment:  cfg.Environment == "development",
	}))
	router.Use(middleware.RateLimitMiddleware())

	// 仅提供API服务，不服务前端文件
	setupAPIRoutes(router, authAPI, dnsAPI, jwtService)

	log.Printf("API服务器启动在端口 %s", cfg.Port)
	log.Printf("允许的来源: %v", cfg.AllowedOrigins)
	log.Fatal(router.Run(":" + cfg.Port))
}

func setupAPIRoutes(router *gin.Engine, authAPI *api.AuthAPI, dnsAPI *api.SimpleDNSAPI, jwtService *utils.JWTService) {
	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "API服务运行正常",
			"version": "1.0.0",
			"mode":    "api-only",
		})
	})

	// 根路径提示
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Domain MAX API Server",
			"version": "1.0.0",
			"docs":    "/api/v1/docs",
			"health":  "/health",
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