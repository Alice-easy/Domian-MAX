package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	visitors map[string]*Visitor
	mutex    sync.RWMutex
	rate     time.Duration
	capacity int
}

// Visitor 访问者信息
type Visitor struct {
	limiter  chan struct{}
	lastSeen time.Time
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(rate time.Duration, capacity int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		capacity: capacity,
	}
	
	// 启动清理协程
	go rl.cleanupVisitors()
	
	return rl
}

// getVisitor 获取访问者
func (rl *RateLimiter) getVisitor(ip string) *Visitor {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	visitor, exists := rl.visitors[ip]
	if !exists {
		visitor = &Visitor{
			limiter:  make(chan struct{}, rl.capacity),
			lastSeen: time.Now(),
		}
		rl.visitors[ip] = visitor
		
		// 填充令牌桶
		for i := 0; i < rl.capacity; i++ {
			visitor.limiter <- struct{}{}
		}
		
		// 启动令牌补充协程
		go rl.refillTokens(visitor)
	}
	
	visitor.lastSeen = time.Now()
	return visitor
}

// refillTokens 补充令牌
func (rl *RateLimiter) refillTokens(visitor *Visitor) {
	ticker := time.NewTicker(rl.rate)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			select {
			case visitor.limiter <- struct{}{}:
			default:
				// 令牌桶已满
			}
		}
	}
}

// cleanupVisitors 清理过期访问者
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			rl.mutex.Lock()
			for ip, visitor := range rl.visitors {
				if time.Since(visitor.lastSeen) > 3*time.Minute {
					delete(rl.visitors, ip)
				}
			}
			rl.mutex.Unlock()
		}
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(ip string) bool {
	visitor := rl.getVisitor(ip)
	
	select {
	case <-visitor.limiter:
		return true
	default:
		return false
	}
}

// 全局速率限制器
var globalRateLimiter = NewRateLimiter(time.Second, 10) // 每秒10个请求

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		if !globalRateLimiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// CustomRateLimitMiddleware 自定义速率限制中间件
func CustomRateLimitMiddleware(rate time.Duration, capacity int) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, capacity)
	
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}