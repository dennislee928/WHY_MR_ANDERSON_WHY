package ratelimit

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// TokenBucketLimiter Token Bucket 速率限制器
type TokenBucketLimiter struct {
	config  *Config
	logger  *logrus.Logger
	buckets map[string]*bucket
	mu      sync.RWMutex
	stopCh  chan struct{}
}

// bucket Token Bucket 實例
type bucket struct {
	tokens         float64
	lastRefillTime time.Time
	failedAttempts int
	lockedUntil    time.Time
	blockedUntil   time.Time
	mu             sync.Mutex
}

// Config 速率限制配置
type Config struct {
	// Token Bucket 配置
	Enabled    bool          `yaml:"enabled" json:"enabled"`
	Rate       int           `yaml:"rate" json:"rate"`               // 每秒允許的請求數
	Burst      int           `yaml:"burst" json:"burst"`             // 桶容量（突發流量）
	WindowSize time.Duration `yaml:"window_size" json:"window_size"` // 時間窗口大小

	// 暴力攻擊防護
	MaxAttempts int           `yaml:"max_attempts" json:"max_attempts"` // 最大失敗嘗試次數
	LockoutTime time.Duration `yaml:"lockout_time" json:"lockout_time"` // 鎖定時間

	// IP 封鎖
	BlockEnabled bool          `yaml:"block_enabled" json:"block_enabled"` // 是否啟用封鎖
	BlockTime    time.Duration `yaml:"block_time" json:"block_time"`       // 封鎖時間

	// 清理配置
	CleanupInterval time.Duration `yaml:"cleanup_interval" json:"cleanup_interval"` // 清理過期 bucket 的間隔
}

// Status 限制器狀態
type Status struct {
	Allowed        bool      `json:"allowed"`
	Remaining      int       `json:"remaining"`
	FailedAttempts int       `json:"failed_attempts"`
	IsLocked       bool      `json:"is_locked"`
	IsBlocked      bool      `json:"is_blocked"`
	LockedUntil    time.Time `json:"locked_until,omitempty"`
	BlockedUntil   time.Time `json:"blocked_until,omitempty"`
}

// NewTokenBucketLimiter 創建新的 Token Bucket 限制器
func NewTokenBucketLimiter(config *Config, logger *logrus.Logger) *TokenBucketLimiter {
	if logger == nil {
		logger = logrus.New()
	}

	// 設定預設值
	if config.Rate == 0 {
		config.Rate = 100 // 每秒 100 個請求
	}
	if config.Burst == 0 {
		config.Burst = 200 // 允許突發 200 個請求
	}
	if config.WindowSize == 0 {
		config.WindowSize = time.Second
	}
	if config.MaxAttempts == 0 {
		config.MaxAttempts = 5
	}
	if config.LockoutTime == 0 {
		config.LockoutTime = 15 * time.Minute
	}
	if config.BlockTime == 0 {
		config.BlockTime = 24 * time.Hour
	}
	if config.CleanupInterval == 0 {
		config.CleanupInterval = 5 * time.Minute
	}

	limiter := &TokenBucketLimiter{
		config:  config,
		logger:  logger,
		buckets: make(map[string]*bucket),
		stopCh:  make(chan struct{}),
	}

	// 啟動清理協程
	go limiter.cleanupRoutine()

	return limiter
}

// Allow 檢查是否允許請求
func (l *TokenBucketLimiter) Allow(key string) (bool, error) {
	if !l.config.Enabled {
		return true, nil
	}

	l.mu.Lock()
	b, exists := l.buckets[key]
	if !exists {
		b = &bucket{
			tokens:         float64(l.config.Burst),
			lastRefillTime: time.Now(),
		}
		l.buckets[key] = b
	}
	l.mu.Unlock()

	b.mu.Lock()
	defer b.mu.Unlock()

	// 檢查是否被封鎖
	if l.config.BlockEnabled && !b.blockedUntil.IsZero() && time.Now().Before(b.blockedUntil) {
		return false, fmt.Errorf("IP 已被封鎖至 %s", b.blockedUntil.Format(time.RFC3339))
	}

	// 檢查是否被鎖定
	if !b.lockedUntil.IsZero() && time.Now().Before(b.lockedUntil) {
		return false, fmt.Errorf("帳號已被鎖定至 %s", b.lockedUntil.Format(time.RFC3339))
	}

	// 補充 Token
	now := time.Now()
	elapsed := now.Sub(b.lastRefillTime)
	tokensToAdd := elapsed.Seconds() * float64(l.config.Rate)
	b.tokens += tokensToAdd
	if b.tokens > float64(l.config.Burst) {
		b.tokens = float64(l.config.Burst)
	}
	b.lastRefillTime = now

	// 檢查是否有足夠的 Token
	if b.tokens < 1 {
		l.logger.Debugf("速率限制 [%s]: tokens=%.2f", key, b.tokens)
		return false, nil
	}

	// 消耗一個 Token
	b.tokens--
	return true, nil
}

// RecordFailedAttempt 記錄失敗的認證嘗試
func (l *TokenBucketLimiter) RecordFailedAttempt(key string) error {
	l.mu.Lock()
	b, exists := l.buckets[key]
	if !exists {
		b = &bucket{
			tokens:         float64(l.config.Burst),
			lastRefillTime: time.Now(),
		}
		l.buckets[key] = b
	}
	l.mu.Unlock()

	b.mu.Lock()
	defer b.mu.Unlock()

	b.failedAttempts++
	l.logger.Warnf("認證失敗 [%s]: %d/%d 次", key, b.failedAttempts, l.config.MaxAttempts)

	// 檢查是否需要鎖定
	if b.failedAttempts >= l.config.MaxAttempts {
		b.lockedUntil = time.Now().Add(l.config.LockoutTime)
		l.logger.Errorf("帳號已鎖定 [%s]: 失敗次數超過 %d 次，鎖定至 %s",
			key, l.config.MaxAttempts, b.lockedUntil.Format(time.RFC3339))

		// 如果啟用封鎖且失敗次數過多，進行封鎖
		if l.config.BlockEnabled && b.failedAttempts >= l.config.MaxAttempts*2 {
			b.blockedUntil = time.Now().Add(l.config.BlockTime)
			l.logger.Errorf("IP 已封鎖 [%s]: 失敗次數超過 %d 次，封鎖至 %s",
				key, l.config.MaxAttempts*2, b.blockedUntil.Format(time.RFC3339))
		}
	}

	return nil
}

// Reset 重置限制器（成功認證後調用）
func (l *TokenBucketLimiter) Reset(key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if b, exists := l.buckets[key]; exists {
		b.mu.Lock()
		b.failedAttempts = 0
		b.lockedUntil = time.Time{}
		b.mu.Unlock()
		l.logger.Infof("重置限制器 [%s]", key)
	}

	return nil
}

// IsLocked 檢查是否被鎖定
func (l *TokenBucketLimiter) IsLocked(key string) bool {
	l.mu.RLock()
	b, exists := l.buckets[key]
	l.mu.RUnlock()

	if !exists {
		return false
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	return !b.lockedUntil.IsZero() && time.Now().Before(b.lockedUntil)
}

// IsBlocked 檢查是否被封鎖
func (l *TokenBucketLimiter) IsBlocked(key string) bool {
	if !l.config.BlockEnabled {
		return false
	}

	l.mu.RLock()
	b, exists := l.buckets[key]
	l.mu.RUnlock()

	if !exists {
		return false
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	return !b.blockedUntil.IsZero() && time.Now().Before(b.blockedUntil)
}

// GetStatus 獲取限制器狀態
func (l *TokenBucketLimiter) GetStatus(key string) Status {
	l.mu.RLock()
	b, exists := l.buckets[key]
	l.mu.RUnlock()

	if !exists {
		return Status{
			Allowed:   true,
			Remaining: l.config.Burst,
		}
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	return Status{
		Allowed:        b.tokens >= 1,
		Remaining:      int(b.tokens),
		FailedAttempts: b.failedAttempts,
		IsLocked:       !b.lockedUntil.IsZero() && time.Now().Before(b.lockedUntil),
		IsBlocked:      !b.blockedUntil.IsZero() && time.Now().Before(b.blockedUntil),
		LockedUntil:    b.lockedUntil,
		BlockedUntil:   b.blockedUntil,
	}
}

// GetRemaining 獲取剩餘 Token 數量
func (l *TokenBucketLimiter) GetRemaining(key string) (int, error) {
	l.mu.RLock()
	b, exists := l.buckets[key]
	l.mu.RUnlock()

	if !exists {
		return l.config.Burst, nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	return int(b.tokens), nil
}

// cleanupRoutine 清理過期的 bucket
func (l *TokenBucketLimiter) cleanupRoutine() {
	ticker := time.NewTicker(l.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.cleanup()
		case <-l.stopCh:
			return
		}
	}
}

// cleanup 清理過期的 bucket
func (l *TokenBucketLimiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for key, b := range l.buckets {
		b.mu.Lock()
		// 如果 bucket 已經很久沒使用且沒有被鎖定/封鎖，則刪除
		if now.Sub(b.lastRefillTime) > l.config.CleanupInterval &&
			(b.lockedUntil.IsZero() || now.After(b.lockedUntil)) &&
			(b.blockedUntil.IsZero() || now.After(b.blockedUntil)) {
			delete(l.buckets, key)
			l.logger.Debugf("清理過期 bucket: %s", key)
		}
		b.mu.Unlock()
	}
}

// Stop 停止限制器
func (l *TokenBucketLimiter) Stop() {
	close(l.stopCh)
	l.logger.Info("速率限制器已停止")
}
