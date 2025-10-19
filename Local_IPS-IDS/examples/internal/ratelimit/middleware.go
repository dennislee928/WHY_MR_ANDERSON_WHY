package ratelimit

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Middleware Rate Limiting 中間件
type Middleware struct {
	limiter *TokenBucketLimiter
	logger  *logrus.Logger
	config  *MiddlewareConfig
}

// MiddlewareConfig 中間件配置
type MiddlewareConfig struct {
	// Key 生成策略
	KeyStrategy string `yaml:"key_strategy" json:"key_strategy"` // "ip", "user", "ip+user"

	// 響應配置
	StatusCode   int    `yaml:"status_code" json:"status_code"`     // HTTP 狀態碼
	ErrorMessage string `yaml:"error_message" json:"error_message"` // 錯誤訊息
	RetryAfter   bool   `yaml:"retry_after" json:"retry_after"`     // 是否返回 Retry-After header

	// 豁免配置
	WhitelistIPs  []string `yaml:"whitelist_ips" json:"whitelist_ips"`   // IP 白名單
	WhitelistKeys []string `yaml:"whitelist_keys" json:"whitelist_keys"` // Key 白名單
}

// NewMiddleware 建立新的中間件
func NewMiddleware(limiter *TokenBucketLimiter, config *MiddlewareConfig, logger *logrus.Logger) *Middleware {
	if logger == nil {
		logger = logrus.New()
	}

	if config == nil {
		config = &MiddlewareConfig{
			KeyStrategy:  "ip",
			StatusCode:   http.StatusTooManyRequests,
			ErrorMessage: "Too many requests",
			RetryAfter:   true,
		}
	}

	return &Middleware{
		limiter: limiter,
		logger:  logger,
		config:  config,
	}
}

// Handler Gin 中間件處理器
func (m *Middleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成限制 key
		key := m.generateKey(c)

		// 檢查白名單
		if m.isWhitelisted(key, c.ClientIP()) {
			c.Next()
			return
		}

		// 檢查速率限制
		allowed, err := m.limiter.Allow(key)
		if err != nil {
			m.logger.Errorf("檢查速率限制失敗 [%s]: %v", key, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			c.Abort()
			return
		}

		if !allowed {
			// 獲取狀態信息
			status := m.limiter.GetStatus(key)

			// 設定 Retry-After header
			if m.config.RetryAfter {
				c.Header("Retry-After", "60") // 60 秒後重試
			}

			// 設定 X-RateLimit headers
			remaining, _ := m.limiter.GetRemaining(key)
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", m.limiter.config.Rate))
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

			m.logger.Warnf("請求被限制 [%s]: %s - %s", key, c.ClientIP(), c.Request.URL.Path)

			c.JSON(m.config.StatusCode, gin.H{
				"error":  m.config.ErrorMessage,
				"status": status,
				"key":    key,
			})
			c.Abort()
			return
		}

		// 設定速率限制資訊 headers
		remaining, _ := m.limiter.GetRemaining(key)
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", m.limiter.config.Rate))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		c.Next()
	}
}

// BruteForceProtection 暴力攻擊防護中間件
func (m *Middleware) BruteForceProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := m.generateKey(c)

		// 檢查是否被鎖定
		if m.limiter.IsLocked(key) {
			status := m.limiter.GetStatus(key)
			m.logger.Warnf("請求被鎖定 [%s]: %s", key, c.ClientIP())

			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Account locked due to too many failed attempts",
				"status": status,
			})
			c.Abort()
			return
		}

		// 檢查是否被封鎖
		if m.limiter.IsBlocked(key) {
			status := m.limiter.GetStatus(key)
			m.logger.Errorf("請求被封鎖 [%s]: %s", key, c.ClientIP())

			c.JSON(http.StatusForbidden, gin.H{
				"error":  "IP blocked due to suspicious activity",
				"status": status,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// generateKey 生成限制 key
func (m *Middleware) generateKey(c *gin.Context) string {
	switch m.config.KeyStrategy {
	case "ip":
		return c.ClientIP()
	case "user":
		// 從認證 token 或 session 獲取用戶 ID
		if userID, exists := c.Get("user_id"); exists {
			return fmt.Sprintf("user:%v", userID)
		}
		return c.ClientIP()
	case "ip+user":
		userID := c.ClientIP()
		if uid, exists := c.Get("user_id"); exists {
			userID = fmt.Sprintf("%s:%v", userID, uid)
		}
		return userID
	default:
		return c.ClientIP()
	}
}

// isWhitelisted 檢查是否在白名單中
func (m *Middleware) isWhitelisted(key, ip string) bool {
	// 檢查 IP 白名單
	for _, whiteIP := range m.config.WhitelistIPs {
		if ip == whiteIP {
			return true
		}
	}

	// 檢查 Key 白名單
	for _, whiteKey := range m.config.WhitelistKeys {
		if key == whiteKey {
			return true
		}
	}

	return false
}

// RecordFailedAuth 記錄失敗的認證嘗試（用於認證處理器）
func (m *Middleware) RecordFailedAuth(key string) error {
	return m.limiter.RecordFailedAttempt(key)
}

// ResetLimit 重置限制（用於成功認證後）
func (m *Middleware) ResetLimit(key string) error {
	return m.limiter.Reset(key)
}
