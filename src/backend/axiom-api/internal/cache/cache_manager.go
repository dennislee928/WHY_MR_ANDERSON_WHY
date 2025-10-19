package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Manager Redis 快取管理器
type Manager struct {
	client *redis.Client
}

// NewManager 創建新的快取管理器
func NewManager(client *redis.Client) *Manager {
	return &Manager{
		client: client,
	}
}

// Set 設置快取（帶 TTL）
func (m *Manager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return m.client.Set(ctx, key, data, ttl).Err()
}

// Get 獲取快取
func (m *Manager) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := m.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete 刪除快取
func (m *Manager) Delete(ctx context.Context, keys ...string) error {
	return m.client.Del(ctx, keys...).Err()
}

// Exists 檢查 Key 是否存在
func (m *Manager) Exists(ctx context.Context, key string) (bool, error) {
	result, err := m.client.Exists(ctx, key).Result()
	return result > 0, err
}

// Increment 增加計數器
func (m *Manager) Increment(ctx context.Context, key string) (int64, error) {
	return m.client.Incr(ctx, key).Result()
}

// IncrementBy 增加計數器指定值
func (m *Manager) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return m.client.IncrBy(ctx, key, value).Result()
}

// Decrement 減少計數器
func (m *Manager) Decrement(ctx context.Context, key string) (int64, error) {
	return m.client.Decr(ctx, key).Result()
}

// GetCounter 獲取計數器值
func (m *Manager) GetCounter(ctx context.Context, key string) (int64, error) {
	return m.client.Get(ctx, key).Int64()
}

// SetExpire 設置 Key 過期時間
func (m *Manager) SetExpire(ctx context.Context, key string, ttl time.Duration) error {
	return m.client.Expire(ctx, key, ttl).Err()
}

// GetTTL 獲取 Key 剩餘生存時間
func (m *Manager) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return m.client.TTL(ctx, key).Result()
}

// Keys 獲取符合模式的所有 Key
func (m *Manager) Keys(ctx context.Context, pattern string) ([]string, error) {
	return m.client.Keys(ctx, pattern).Result()
}

// FlushPattern 刪除符合模式的所有 Key
func (m *Manager) FlushPattern(ctx context.Context, pattern string) error {
	keys, err := m.Keys(ctx, pattern)
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return m.Delete(ctx, keys...)
	}
	return nil
}

// SetNX 只在 Key 不存在時設置（分布式鎖）
func (m *Manager) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return m.client.SetNX(ctx, key, data, ttl).Result()
}

// GetMultiple 批量獲取快取
func (m *Manager) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	pipe := m.client.Pipeline()
	cmds := make(map[string]*redis.StringCmd)

	for _, key := range keys {
		cmds[key] = pipe.Get(ctx, key)
	}

	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for key, cmd := range cmds {
		data, err := cmd.Bytes()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			continue
		}
		
		var value interface{}
		if err := json.Unmarshal(data, &value); err == nil {
			result[key] = value
		}
	}

	return result, nil
}

// SetMultiple 批量設置快取
func (m *Manager) SetMultiple(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	pipe := m.client.Pipeline()

	for key, value := range items {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		pipe.Set(ctx, key, data, ttl)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// Client 獲取底層 Redis 客戶端
func (m *Manager) Client() *redis.Client {
	return m.client
}

