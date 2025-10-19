package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// SmartCache implements intelligent caching with adaptive TTL
type SmartCache struct {
	redis         *redis.Client
	localCache    map[string]*CacheEntry
	mu            sync.RWMutex
	logger        *logrus.Logger
	stats         *CacheStats
	config        *CacheConfig
}

// CacheConfig contains cache configuration
type CacheConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	DefaultTTL    time.Duration
	MaxLocalSize  int
	EnableAdaptiveTTL bool
	EnablePrefetch    bool
}

// CacheEntry represents a cache entry
type CacheEntry struct {
	Key          string
	Value        interface{}
	CreatedAt    time.Time
	ExpiresAt    time.Time
	AccessCount  int64
	LastAccess   time.Time
	HitRate      float64
	Priority     int
}

// CacheStats contains cache statistics
type CacheStats struct {
	TotalHits      int64
	TotalMisses    int64
	LocalHits      int64
	RedisHits      int64
	Evictions      int64
	PrefetchHits   int64
	mu             sync.RWMutex
}

// NewSmartCache creates a new smart cache
func NewSmartCache(config *CacheConfig, logger *logrus.Logger) (*SmartCache, error) {
	if logger == nil {
		logger = logrus.New()
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// 測試連接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	sc := &SmartCache{
		redis:      redisClient,
		localCache: make(map[string]*CacheEntry),
		logger:     logger,
		stats:      &CacheStats{},
		config:     config,
	}

	// 啟動後台任務
	go sc.startEvictionWorker()
	if config.EnablePrefetch {
		go sc.startPrefetchWorker()
	}

	logger.Info("Smart cache initialized")
	return sc, nil
}

// Get retrieves a value from cache
func (sc *SmartCache) Get(ctx context.Context, key string) (interface{}, error) {
	// 1. 檢查本地緩存
	if value, found := sc.getFromLocal(key); found {
		sc.recordHit(true)
		return value, nil
	}

	// 2. 檢查 Redis
	value, err := sc.getFromRedis(ctx, key)
	if err == nil {
		sc.recordHit(false)
		// 提升到本地緩存
		sc.promoteToLocal(key, value)
		return value, nil
	}

	// 3. Cache miss
	sc.recordMiss()
	return nil, fmt.Errorf("cache miss: %s", key)
}

// Set stores a value in cache
func (sc *SmartCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// 計算自適應 TTL
	if sc.config.EnableAdaptiveTTL {
		ttl = sc.calculateAdaptiveTTL(key, ttl)
	}

	// 存儲到 Redis
	if err := sc.setToRedis(ctx, key, value, ttl); err != nil {
		return err
	}

	// 存儲到本地緩存
	sc.setToLocal(key, value, ttl)

	return nil
}

// getFromLocal retrieves from local cache
func (sc *SmartCache) getFromLocal(key string) (interface{}, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	entry, exists := sc.localCache[key]
	if !exists {
		return nil, false
	}

	// 檢查是否過期
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	// 更新訪問統計
	entry.AccessCount++
	entry.LastAccess = time.Now()

	return entry.Value, true
}

// getFromRedis retrieves from Redis
func (sc *SmartCache) getFromRedis(ctx context.Context, key string) (interface{}, error) {
	value, err := sc.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return value, nil
}

// setToLocal stores to local cache
func (sc *SmartCache) setToLocal(key string, value interface{}, ttl time.Duration) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// 檢查本地緩存大小
	if len(sc.localCache) >= sc.config.MaxLocalSize {
		sc.evictLRU()
	}

	entry := &CacheEntry{
		Key:         key,
		Value:       value,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(ttl),
		AccessCount: 0,
		LastAccess:  time.Now(),
		Priority:    1,
	}

	sc.localCache[key] = entry
}

// setToRedis stores to Redis
func (sc *SmartCache) setToRedis(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return sc.redis.Set(ctx, key, value, ttl).Err()
}

// promoteToLocal promotes a Redis entry to local cache
func (sc *SmartCache) promoteToLocal(key string, value interface{}) {
	sc.setToLocal(key, value, sc.config.DefaultTTL)
}

// evictLRU evicts least recently used entry
func (sc *SmartCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time = time.Now()

	for key, entry := range sc.localCache {
		if entry.LastAccess.Before(oldestTime) {
			oldestTime = entry.LastAccess
			oldestKey = key
		}
	}

	if oldestKey != "" {
		delete(sc.localCache, oldestKey)
		sc.stats.mu.Lock()
		sc.stats.Evictions++
		sc.stats.mu.Unlock()
	}
}

// calculateAdaptiveTTL calculates adaptive TTL based on access patterns
func (sc *SmartCache) calculateAdaptiveTTL(key string, defaultTTL time.Duration) time.Duration {
	sc.mu.RLock()
	entry, exists := sc.localCache[key]
	sc.mu.RUnlock()

	if !exists {
		return defaultTTL
	}

	// 根據訪問頻率調整 TTL
	accessRate := float64(entry.AccessCount) / time.Since(entry.CreatedAt).Minutes()

	if accessRate > 10 {
		// 高頻訪問：延長 TTL
		return defaultTTL * 2
	} else if accessRate < 1 {
		// 低頻訪問：縮短 TTL
		return defaultTTL / 2
	}

	return defaultTTL
}

// startEvictionWorker starts background eviction worker
func (sc *SmartCache) startEvictionWorker() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sc.evictExpired()
	}
}

// evictExpired removes expired entries
func (sc *SmartCache) evictExpired() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	now := time.Now()
	for key, entry := range sc.localCache {
		if now.After(entry.ExpiresAt) {
			delete(sc.localCache, key)
			sc.stats.mu.Lock()
			sc.stats.Evictions++
			sc.stats.mu.Unlock()
		}
	}
}

// startPrefetchWorker starts background prefetch worker
func (sc *SmartCache) startPrefetchWorker() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sc.prefetchHotKeys()
	}
}

// prefetchHotKeys prefetches frequently accessed keys
func (sc *SmartCache) prefetchHotKeys() {
	sc.mu.RLock()
	hotKeys := make([]string, 0)

	for key, entry := range sc.localCache {
		// 識別熱點數據
		if entry.AccessCount > 100 {
			hotKeys = append(hotKeys, key)
		}
	}
	sc.mu.RUnlock()

	// 預取熱點數據
	ctx := context.Background()
	for _, key := range hotKeys {
		if _, err := sc.getFromRedis(ctx, key); err == nil {
			sc.stats.mu.Lock()
			sc.stats.PrefetchHits++
			sc.stats.mu.Unlock()
		}
	}

	if len(hotKeys) > 0 {
		sc.logger.Debugf("Prefetched %d hot keys", len(hotKeys))
	}
}

// recordHit records a cache hit
func (sc *SmartCache) recordHit(local bool) {
	sc.stats.mu.Lock()
	defer sc.stats.mu.Unlock()

	sc.stats.TotalHits++
	if local {
		sc.stats.LocalHits++
	} else {
		sc.stats.RedisHits++
	}
}

// recordMiss records a cache miss
func (sc *SmartCache) recordMiss() {
	sc.stats.mu.Lock()
	defer sc.stats.mu.Unlock()

	sc.stats.TotalMisses++
}

// GetStats returns cache statistics
func (sc *SmartCache) GetStats() *CacheStats {
	sc.stats.mu.RLock()
	defer sc.stats.mu.RUnlock()

	return &CacheStats{
		TotalHits:    sc.stats.TotalHits,
		TotalMisses:  sc.stats.TotalMisses,
		LocalHits:    sc.stats.LocalHits,
		RedisHits:    sc.stats.RedisHits,
		Evictions:    sc.stats.Evictions,
		PrefetchHits: sc.stats.PrefetchHits,
	}
}

// GetHitRate returns cache hit rate
func (sc *SmartCache) GetHitRate() float64 {
	sc.stats.mu.RLock()
	defer sc.stats.mu.RUnlock()

	total := sc.stats.TotalHits + sc.stats.TotalMisses
	if total == 0 {
		return 0.0
	}

	return float64(sc.stats.TotalHits) / float64(total)
}

// Clear clears all cache
func (sc *SmartCache) Clear(ctx context.Context) error {
	sc.mu.Lock()
	sc.localCache = make(map[string]*CacheEntry)
	sc.mu.Unlock()

	// 清除 Redis（謹慎使用）
	// return sc.redis.FlushDB(ctx).Err()

	return nil
}

// Close closes the cache
func (sc *SmartCache) Close() error {
	if sc.redis != nil {
		return sc.redis.Close()
	}
	return nil
}

// DefaultConfig returns default cache configuration
func DefaultConfig() *CacheConfig {
	return &CacheConfig{
		RedisAddr:         "localhost:6379",
		RedisPassword:     "",
		RedisDB:           0,
		DefaultTTL:        5 * time.Minute,
		MaxLocalSize:      1000,
		EnableAdaptiveTTL: true,
		EnablePrefetch:    true,
	}
}

