package main

import (
	"domain-max/pkg/api"
	"domain-max/pkg/config"
	"domain-max/pkg/database"
	"domain-max/pkg/middleware"
	"domain-max/pkg/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 连接数据库
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 运行数据库迁移
	if err := database.Migrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化服务
	jwtService := utils.NewJWTService(cfg.JWTSecret, cfg.JWTExpirationHours)
	passwordService := utils.NewPasswordService()
	encryptionService, err := utils.NewEncryptionService(cfg.EncryptionKey)
	if err != nil {
		log.Fatalf("加密服务初始化失败: %v", err)
	}
	validationService := utils.NewValidationService()

	// 初始化API控制器
	authAPI := api.NewAuthAPI(db, jwtService, passwordService, validationService)

	// 设置Gin模式
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// 添加中间件
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowedOrigins: cfg.AllowedOrigins,
		IsDevelopment:  cfg.IsDevelopment(),
	}))
	router.Use(middleware.RateLimitMiddleware())
	
	if cfg.IsDevelopment() {
		router.Use(middleware.LoggingMiddleware())
	}

	// 设置路由
	setupAPIRoutes(router, authAPI, jwtService)

	log.Printf("API服务器启动在端口 %s", cfg.Port)
	log.Printf("环境: %s", cfg.Environment)
	log.Printf("允许的来源: %v", cfg.AllowedOrigins)
	log.Fatal(router.Run(":" + cfg.Port))
}

func setupAPIRoutes(router *gin.Engine, authAPI *api.AuthAPI, jwtService *utils.JWTService) {
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
		auth.POST("/refresh", authAPI.RefreshToken)
	}

	// 需要认证的路由
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		// 用户相关
		protected.GET("/auth/profile", authAPI.GetProfile)
		protected.POST("/auth/change-password", authAPI.ChangePassword)

		// DNS提供商管理
		providers := protected.Group("/dns-providers")
		{
			providers.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "DNS提供商功能待实现",
				})
			})
		}

		// DNS记录管理
		records := protected.Group("/dns-records")
		{
			records.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "DNS记录功能待实现",
				})
			})
		}

		// 管理员路由
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminRequiredMiddleware())
		{
			// 用户管理
			admin.GET("/users", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "用户管理功能待实现",
				})
			})

			// 系统统计
			admin.GET("/stats", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "系统统计功能待实现",
				})
			})
		}
	}
}