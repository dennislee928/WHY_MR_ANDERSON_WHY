package ratelimit

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// TokenBucket implements the token bucket rate limiting algorithm
// Token Bucket 率限制算法實現
type TokenBucket struct {
	capacity     int64         // 桶容量
	tokens       int64         // 當前令牌數
	refillRate   int64         // 每秒補充令牌數
	lastRefill   time.Time     // 上次補充時間
	mu           sync.RWMutex
	logger       *logrus.Logger
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity, refillRate int64, logger *logrus.Logger) *TokenBucket {
	if logger == nil {
		logger = logrus.New()
	}

	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
		logger:     logger,
	}
}

// Allow checks if a request is allowed
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// 補充令牌
	tb.refill()

	// 檢查是否有可用令牌
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// AllowN checks if N requests are allowed
func (tb *TokenBucket) AllowN(n int64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}

	return false
}

// refill refills tokens based on elapsed time
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	// 計算應該補充的令牌數
	tokensToAdd := int64(elapsed.Seconds()) * tb.refillRate

	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}
}

// GetTokens returns the current number of tokens
func (tb *TokenBucket) GetTokens() int64 {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	return tb.tokens
}

// Reset resets the token bucket
func (tb *TokenBucket) Reset() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.tokens = tb.capacity
	tb.lastRefill = time.Now()
}

// RedisTokenBucket implements distributed token bucket using Redis
// 使用 Redis 實現的分散式 Token Bucket
type RedisTokenBucket struct {
	redis      *redis.Client
	key        string
	capacity   int64
	refillRate int64
	logger     *logrus.Logger
}

// NewRedisTokenBucket creates a new Redis-based token bucket
func NewRedisTokenBucket(redis *redis.Client, key string, capacity, refillRate int64, logger *logrus.Logger) *RedisTokenBucket {
	if logger == nil {
		logger = logrus.New()
	}

	return &RedisTokenBucket{
		redis:      redis,
		key:        key,
		capacity:   capacity,
		refillRate: refillRate,
		logger:     logger,
	}
}

// Allow checks if a request is allowed (distributed)
func (rtb *RedisTokenBucket) Allow(ctx context.Context) (bool, error) {
	// 使用 Lua 腳本實現原子操作
	script := `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

-- 獲取當前狀態
local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
local tokens = tonumber(bucket[1]) or capacity
local last_refill = tonumber(bucket[2]) or now

-- 計算補充的令牌
local elapsed = now - last_refill
local tokens_to_add = math.floor(elapsed * refill_rate)

if tokens_to_add > 0 then
    tokens = math.min(tokens + tokens_to_add, capacity)
    last_refill = now
end

-- 檢查是否有可用令牌
if tokens > 0 then
    tokens = tokens - 1
    redis.call('HMSET', key, 'tokens', tokens, 'last_refill', last_refill)
    redis.call('EXPIRE', key, 3600)  -- 1 小時過期
    return 1
else
    return 0
end
`

	result, err := rtb.redis.Eval(ctx, script, []string{rtb.key},
		rtb.capacity, rtb.refillRate, time.Now().Unix()).Int()

	if err != nil {
		rtb.logger.Errorf("Redis token bucket error: %v", err)
		return false, err
	}

	return result == 1, nil
}

// GetTokens returns the current number of tokens (distributed)
func (rtb *RedisTokenBucket) GetTokens(ctx context.Context) (int64, error) {
	result, err := rtb.redis.HGet(ctx, rtb.key, "tokens").Int64()
	if err == redis.Nil {
		return rtb.capacity, nil
	}
	if err != nil {
		return 0, err
	}
	return result, nil
}

// Reset resets the token bucket (distributed)
func (rtb *RedisTokenBucket) Reset(ctx context.Context) error {
	return rtb.redis.HMSet(ctx, rtb.key,
		"tokens", rtb.capacity,
		"last_refill", time.Now().Unix(),
	).Err()
}

// MultiLevelRateLimiter implements multi-level rate limiting
// 多層級率限制器
type MultiLevelRateLimiter struct {
	ipLimiter       map[string]*TokenBucket
	endpointLimiter map[string]*TokenBucket
	userLimiter     map[string]*TokenBucket
	redis           *redis.Client
	logger          *logrus.Logger
	mu              sync.RWMutex
}

// RateLimitConfig contains rate limit configuration
type RateLimitConfig struct {
	// IP 層級限制
	IPRequestsPerMinute int64

	// 端點層級限制
	EndpointRequestsPerMinute int64

	// 用戶層級限制
	UserRequestsPerHour int64

	// Redis 配置（用於分散式限制）
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

// NewMultiLevelRateLimiter creates a new multi-level rate limiter
func NewMultiLevelRateLimiter(config *RateLimitConfig, logger *logrus.Logger) *MultiLevelRateLimiter {
	if logger == nil {
		logger = logrus.New()
	}

	var redisClient *redis.Client
	if config.RedisAddr != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.RedisAddr,
			Password: config.RedisPassword,
			DB:       config.RedisDB,
		})
	}

	return &MultiLevelRateLimiter{
		ipLimiter:       make(map[string]*TokenBucket),
		endpointLimiter: make(map[string]*TokenBucket),
		userLimiter:     make(map[string]*TokenBucket),
		redis:           redisClient,
		logger:          logger,
	}
}

// AllowIP checks if an IP is allowed
func (mrl *MultiLevelRateLimiter) AllowIP(ctx context.Context, ip string) (bool, error) {
	mrl.mu.RLock()
	limiter, exists := mrl.ipLimiter[ip]
	mrl.mu.RUnlock()

	if !exists {
		mrl.mu.Lock()
		limiter = NewTokenBucket(60, 1, mrl.logger) // 60 requests per minute
		mrl.ipLimiter[ip] = limiter
		mrl.mu.Unlock()
	}

	// 如果有 Redis，使用分散式限制
	if mrl.redis != nil {
		redisLimiter := NewRedisTokenBucket(mrl.redis, fmt.Sprintf("ratelimit:ip:%s", ip), 60, 1, mrl.logger)
		return redisLimiter.Allow(ctx)
	}

	return limiter.Allow(), nil
}

// AllowEndpoint checks if an endpoint is allowed
func (mrl *MultiLevelRateLimiter) AllowEndpoint(ctx context.Context, endpoint string) (bool, error) {
	mrl.mu.RLock()
	limiter, exists := mrl.endpointLimiter[endpoint]
	mrl.mu.RUnlock()

	if !exists {
		mrl.mu.Lock()
		limiter = NewTokenBucket(10, int64(math.Round(1000.0/6.0)), mrl.logger) // 10 requests per minute for sensitive endpoints
		mrl.endpointLimiter[endpoint] = limiter
		mrl.mu.Unlock()
	}

	return limiter.Allow(), nil
}

// AllowUser checks if a user is allowed
func (mrl *MultiLevelRateLimiter) AllowUser(ctx context.Context, userID string) (bool, error) {
	mrl.mu.RLock()
	limiter, exists := mrl.userLimiter[userID]
	mrl.mu.RUnlock()

	if !exists {
		mrl.mu.Lock()
		limiter = NewTokenBucket(1000, int64(math.Round(1000000.0/3600.0)), mrl.logger) // 1000 requests per hour
		mrl.userLimiter[userID] = limiter
		mrl.mu.Unlock()
	}

	return limiter.Allow(), nil
}

// CheckAll checks all rate limit levels
func (mrl *MultiLevelRateLimiter) CheckAll(ctx context.Context, ip, endpoint, userID string) (bool, string, error) {
	// 檢查 IP 層級
	if allowed, err := mrl.AllowIP(ctx, ip); err != nil {
		return false, "error", err
	} else if !allowed {
		return false, "ip_rate_limit_exceeded", nil
	}

	// 檢查端點層級
	if allowed, err := mrl.AllowEndpoint(ctx, endpoint); err != nil {
		return false, "error", err
	} else if !allowed {
		return false, "endpoint_rate_limit_exceeded", nil
	}

	// 檢查用戶層級
	if userID != "" {
		if allowed, err := mrl.AllowUser(ctx, userID); err != nil {
			return false, "error", err
		} else if !allowed {
			return false, "user_rate_limit_exceeded", nil
		}
	}

	return true, "", nil
}

// Close closes the rate limiter
func (mrl *MultiLevelRateLimiter) Close() error {
	if mrl.redis != nil {
		return mrl.redis.Close()
	}
	return nil
}

