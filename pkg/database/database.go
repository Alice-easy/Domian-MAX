package database

import (
	"domain-max/pkg/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect 连接数据库
func Connect(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	
	// 根据配置选择数据库驱动
	switch cfg.DatabaseType {
	case "mysql":
		if cfg.DatabaseURL != "" {
			dialector = mysql.Open(cfg.DatabaseURL)
		} else {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName)
			dialector = mysql.Open(dsn)
		}
	case "postgres":
		if cfg.DatabaseURL != "" {
			dialector = postgres.Open(cfg.DatabaseURL)
		} else {
			dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
				cfg.DatabaseHost, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseName, cfg.DatabasePort)
			dialector = postgres.Open(dsn)
		}
	case "sqlite":
		dialector = sqlite.Open("domain_max.db")
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", cfg.DatabaseType)
	}

	// 配置GORM
	gormConfig := &gorm.Config{}
	if cfg.IsDevelopment() {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	log.Printf("数据库连接成功 (类型: %s)", cfg.DatabaseType)
	return db, nil
}

// Migrate 执行数据库迁移
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&DNSProvider{},
		&DNSRecord{},
		&Domain{},
		&SMTPConfig{},
	)
}