package database

import (
	authmodels "domain-max/pkg/auth/models"
	dnsmodels "domain-max/pkg/dns/models"
	emailmodels "domain-max/pkg/email/models"
	"log"

	"gorm.io/gorm"
)

// Migrate 执行数据库迁移
func Migrate(db *gorm.DB) error {
	log.Println("开始数据库迁移...")
	
	// 用户相关表
	if err := db.AutoMigrate(
		&authmodels.User{},
		&authmodels.EmailVerification{},
		&authmodels.PasswordReset{},
	); err != nil {
		log.Printf("用户表迁移失败: %v", err)
		return err
	}
	
	// DNS相关表
	if err := db.AutoMigrate(
		&dnsmodels.Domain{},
		&dnsmodels.SubDomain{},
		&dnsmodels.DNSRecord{},
		&dnsmodels.DNSProvider{},
	); err != nil {
		log.Printf("DNS表迁移失败: %v", err)
		return err
	}
	
	// 邮件相关表
	if err := db.AutoMigrate(
		&emailmodels.SMTPConfig{},
	); err != nil {
		log.Printf("邮件表迁移失败: %v", err)
		return err
	}
	
	// 创建索引
	if err := createIndexes(db); err != nil {
		log.Printf("创建索引失败: %v", err)
		return err
	}
	
	// 插入默认数据
	if err := insertDefaultData(db); err != nil {
		log.Printf("插入默认数据失败: %v", err)
		return err
	}
	
	log.Println("数据库迁移完成")
	return nil
}

// createIndexes 创建额外的索引
func createIndexes(db *gorm.DB) error {
	indexes := []struct {
		table string
		sql   string
	}{
		// 用户表索引
		{"users", "CREATE INDEX IF NOT EXISTS idx_users_email_active ON users(email, is_active) WHERE deleted_at IS NULL"},
		{"users", "CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC)"},
		
		// 域名表索引
		{"domains", "CREATE INDEX IF NOT EXISTS idx_domains_user_platform ON domains(user_id, platform)"},
		{"domains", "CREATE INDEX IF NOT EXISTS idx_domains_platform_active ON domains(platform, is_active)"},
		
		// 子域名表索引
		{"sub_domains", "CREATE INDEX IF NOT EXISTS idx_subdomains_domain_type_status ON sub_domains(domain_id, record_type, status)"},
		{"sub_domains", "CREATE INDEX IF NOT EXISTS idx_subdomains_name_type ON sub_domains(sub_domain_name, record_type)"},
		
		// DNS服务商表索引
		{"dns_providers", "CREATE INDEX IF NOT EXISTS idx_providers_type_active ON dns_providers(type, is_active)"},
	}
	
	for _, idx := range indexes {
		if err := db.Exec(idx.sql).Error; err != nil {
			log.Printf("创建索引失败 %s: %v", idx.sql, err)
			// 索引创建失败不应该中断迁移，只记录日志
		}
	}
	
	return nil
}

// insertDefaultData 插入默认数据
func insertDefaultData(db *gorm.DB) error {
	// 插入默认DNS服务商配置
	providers := []dnsmodels.DNSProvider{
		{
			Name:        "阿里云DNS",
			Type:        "aliyun",
			Description: "阿里云云解析DNS服务",
			IsActive:    true,
			SortOrder:   1,
		},
		{
			Name:        "腾讯云DNSPod",
			Type:        "dnspod",
			Description: "腾讯云DNSPod服务",
			IsActive:    true,
			SortOrder:   2,
		},
		{
			Name:        "CloudFlare",
			Type:        "cloudflare",
			Description: "CloudFlare DNS服务",
			IsActive:    true,
			SortOrder:   3,
		},
		{
			Name:        "华为云DNS",
			Type:        "huawei",
			Description: "华为云云解析服务",
			IsActive:    false,
			SortOrder:   4,
		},
		{
			Name:        "百度云DNS",
			Type:        "baidu",
			Description: "百度智能云DNS服务",
			IsActive:    false,
			SortOrder:   5,
		},
		{
			Name:        "西部数码DNS",
			Type:        "west",
			Description: "西部数码域名解析服务",
			IsActive:    false,
			SortOrder:   6,
		},
		{
			Name:        "火山引擎DNS",
			Type:        "volcengine",
			Description: "火山引擎TrafficRoute DNS",
			IsActive:    false,
			SortOrder:   7,
		},
		{
			Name:        "DNSLA",
			Type:        "dnsla",
			Description: "DNSLA专业DNS解析服务",
			IsActive:    false,
			SortOrder:   8,
		},
		{
			Name:        "Namesilo",
			Type:        "namesilo",
			Description: "Namesilo域名DNS服务",
			IsActive:    false,
			SortOrder:   9,
		},
		{
			Name:        "PowerDNS",
			Type:        "powerdns",
			Description: "PowerDNS开源DNS服务器",
			IsActive:    false,
			SortOrder:   10,
		},
	}
	
	for _, provider := range providers {
		var existing dnsmodels.DNSProvider
		result := db.Where("type = ?", provider.Type).First(&existing)
		if result.Error != nil && result.Error.Error() == "record not found" {
			if err := db.Create(&provider).Error; err != nil {
				log.Printf("插入默认DNS服务商失败 %s: %v", provider.Name, err)
			}
		}
	}
	
	return nil
}