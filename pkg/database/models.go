package database

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"default:user" json:"role"` // user, admin
	IsActive bool   `gorm:"default:true" json:"is_active"`
	
	// 关联
	DNSProviders []DNSProvider `json:"dns_providers,omitempty"`
	Domains      []Domain      `json:"domains,omitempty"`
}

// DNSProvider DNS提供商模型
type DNSProvider struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	UserID   uint   `gorm:"not null" json:"user_id"`
	Name     string `gorm:"not null" json:"name"`
	Type     string `gorm:"not null" json:"type"` // cloudflare, aliyun, etc.
	Config   string `gorm:"type:text" json:"-"`   // 加密的配置信息
	IsActive bool   `gorm:"default:true" json:"is_active"`
	
	// 关联
	User       User        `json:"user,omitempty"`
	DNSRecords []DNSRecord `json:"dns_records,omitempty"`
}

// Domain 域名模型
type Domain struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	UserID       uint   `gorm:"not null" json:"user_id"`
	Name         string `gorm:"uniqueIndex;not null" json:"name"`
	Description  string `json:"description"`
	IsActive     bool   `gorm:"default:true" json:"is_active"`
	DNSProviderID *uint `json:"dns_provider_id"`
	
	// 关联
	User        User         `json:"user,omitempty"`
	DNSProvider *DNSProvider `json:"dns_provider,omitempty"`
	DNSRecords  []DNSRecord  `json:"dns_records,omitempty"`
}

// DNSRecord DNS记录模型
type DNSRecord struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	DomainID      uint   `gorm:"not null" json:"domain_id"`
	DNSProviderID uint   `gorm:"not null" json:"dns_provider_id"`
	Name          string `gorm:"not null" json:"name"`
	Type          string `gorm:"not null" json:"type"` // A, AAAA, CNAME, MX, TXT, etc.
	Value         string `gorm:"not null" json:"value"`
	TTL           int    `gorm:"default:300" json:"ttl"`
	Priority      *int   `json:"priority,omitempty"` // for MX records
	IsActive      bool   `gorm:"default:true" json:"is_active"`
	
	// 关联
	Domain      Domain      `json:"domain,omitempty"`
	DNSProvider DNSProvider `json:"dns_provider,omitempty"`
}

// SMTPConfig SMTP配置模型
type SMTPConfig struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	Name     string `gorm:"not null" json:"name"`
	Host     string `gorm:"not null" json:"host"`
	Port     int    `gorm:"not null" json:"port"`
	Username string `gorm:"not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // 加密存储
	FromName string `json:"from_name"`
	FromEmail string `gorm:"not null" json:"from_email"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
	IsDefault bool  `gorm:"default:false" json:"is_default"`
}