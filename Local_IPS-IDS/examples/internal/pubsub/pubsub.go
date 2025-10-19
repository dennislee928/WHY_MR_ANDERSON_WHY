package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// PubSub 發布訂閱接口
type PubSub interface {
	Publish(ctx context.Context, topic string, message interface{}) error
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	Unsubscribe(ctx context.Context, topic string) error
	Close() error
}

// MessageHandler 訊息處理函數
type MessageHandler func(topic string, message []byte) error

// Config Pub/Sub 配置
type Config struct {
	// 類型: "redis", "memory"
	Type string `yaml:"type" json:"type"`

	// Redis 配置
	RedisAddr     string `yaml:"redis_addr" json:"redis_addr"`
	RedisPassword string `yaml:"redis_password" json:"redis_password"`
	RedisDB       int    `yaml:"redis_db" json:"redis_db"`

	// 通用配置
	BufferSize int `yaml:"buffer_size" json:"buffer_size"`
}

// Event 事件結構
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// ========================
// Redis Pub/Sub 實作
// ========================

// RedisPubSub Redis 實作
type RedisPubSub struct {
	client   *redis.Client
	logger   *logrus.Logger
	handlers map[string][]MessageHandler
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewRedisPubSub 建立 Redis Pub/Sub
func NewRedisPubSub(config *Config, logger *logrus.Logger) (*RedisPubSub, error) {
	if logger == nil {
		logger = logrus.New()
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// 測試連接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis 連接失敗: %w", err)
	}

	ctx, cancel = context.WithCancel(context.Background())

	return &RedisPubSub{
		client:   client,
		logger:   logger,
		handlers: make(map[string][]MessageHandler),
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// Publish 發布訊息
func (r *RedisPubSub) Publish(ctx context.Context, topic string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化訊息失敗: %w", err)
	}

	if err := r.client.Publish(ctx, topic, data).Err(); err != nil {
		return fmt.Errorf("發布訊息失敗: %w", err)
	}

	r.logger.Debugf("發布訊息到主題 [%s]: %d bytes", topic, len(data))
	return nil
}

// Subscribe 訂閱主題
func (r *RedisPubSub) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
	r.mu.Lock()
	r.handlers[topic] = append(r.handlers[topic], handler)
	isFirstSubscriber := len(r.handlers[topic]) == 1
	r.mu.Unlock()

	if isFirstSubscriber {
		// 只有第一個訂閱者時才創建 Redis 訂閱
		pubsub := r.client.Subscribe(ctx, topic)

		r.wg.Add(1)
		go r.listenToChannel(pubsub, topic)

		r.logger.Infof("訂閱主題 [%s]", topic)
	}

	return nil
}

// Unsubscribe 取消訂閱
func (r *RedisPubSub) Unsubscribe(ctx context.Context, topic string) error {
	r.mu.Lock()
	delete(r.handlers, topic)
	r.mu.Unlock()

	r.logger.Infof("取消訂閱主題 [%s]", topic)
	return nil
}

// Close 關閉連接
func (r *RedisPubSub) Close() error {
	r.logger.Info("關閉 Redis Pub/Sub...")
	r.cancel()
	r.wg.Wait()
	return r.client.Close()
}

// listenToChannel 監聽 Redis 頻道
func (r *RedisPubSub) listenToChannel(pubsub *redis.PubSub, topic string) {
	defer r.wg.Done()
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case <-r.ctx.Done():
			return
		case msg := <-ch:
			if msg == nil {
				return
			}

			r.handleMessage(topic, []byte(msg.Payload))
		}
	}
}

// handleMessage 處理訊息
func (r *RedisPubSub) handleMessage(topic string, data []byte) {
	r.mu.RLock()
	handlers := r.handlers[topic]
	r.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(topic, data); err != nil {
			r.logger.Errorf("處理訊息失敗 [%s]: %v", topic, err)
		}
	}
}

// ========================
// In-Memory Pub/Sub 實作
// ========================

// MemoryPubSub 記憶體實作（用於測試或單機部署）
type MemoryPubSub struct {
	logger     *logrus.Logger
	channels   map[string]chan []byte
	handlers   map[string][]MessageHandler
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	bufferSize int
}

// NewMemoryPubSub 建立記憶體 Pub/Sub
func NewMemoryPubSub(config *Config, logger *logrus.Logger) *MemoryPubSub {
	if logger == nil {
		logger = logrus.New()
	}

	bufferSize := config.BufferSize
	if bufferSize == 0 {
		bufferSize = 100
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &MemoryPubSub{
		logger:     logger,
		channels:   make(map[string]chan []byte),
		handlers:   make(map[string][]MessageHandler),
		ctx:        ctx,
		cancel:     cancel,
		bufferSize: bufferSize,
	}
}

// Publish 發布訊息
func (m *MemoryPubSub) Publish(ctx context.Context, topic string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化訊息失敗: %w", err)
	}

	m.mu.RLock()
	ch, exists := m.channels[topic]
	m.mu.RUnlock()

	if !exists {
		m.logger.Warnf("主題 [%s] 沒有訂閱者", topic)
		return nil
	}

	select {
	case ch <- data:
		m.logger.Debugf("發布訊息到主題 [%s]: %d bytes", topic, len(data))
	case <-ctx.Done():
		return ctx.Err()
	default:
		m.logger.Warnf("主題 [%s] 緩衝區已滿，丟棄訊息", topic)
	}

	return nil
}

// Subscribe 訂閱主題
func (m *MemoryPubSub) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.handlers[topic] = append(m.handlers[topic], handler)

	// 如果是第一個訂閱者，創建頻道
	if _, exists := m.channels[topic]; !exists {
		ch := make(chan []byte, m.bufferSize)
		m.channels[topic] = ch

		m.wg.Add(1)
		go m.listenToChannel(ch, topic)

		m.logger.Infof("訂閱主題 [%s]", topic)
	}

	return nil
}

// Unsubscribe 取消訂閱
func (m *MemoryPubSub) Unsubscribe(ctx context.Context, topic string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.handlers, topic)

	if ch, exists := m.channels[topic]; exists {
		close(ch)
		delete(m.channels, topic)
	}

	m.logger.Infof("取消訂閱主題 [%s]", topic)
	return nil
}

// Close 關閉
func (m *MemoryPubSub) Close() error {
	m.logger.Info("關閉 Memory Pub/Sub...")
	m.cancel()

	m.mu.Lock()
	for _, ch := range m.channels {
		close(ch)
	}
	m.mu.Unlock()

	m.wg.Wait()
	return nil
}

// listenToChannel 監聽頻道
func (m *MemoryPubSub) listenToChannel(ch chan []byte, topic string) {
	defer m.wg.Done()

	for {
		select {
		case <-m.ctx.Done():
			return
		case data, ok := <-ch:
			if !ok {
				return
			}
			m.handleMessage(topic, data)
		}
	}
}

// handleMessage 處理訊息
func (m *MemoryPubSub) handleMessage(topic string, data []byte) {
	m.mu.RLock()
	handlers := m.handlers[topic]
	m.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(topic, data); err != nil {
			m.logger.Errorf("處理訊息失敗 [%s]: %v", topic, err)
		}
	}
}

// ========================
// Factory 函數
// ========================

// NewPubSub 建立 Pub/Sub（根據配置選擇實作）
func NewPubSub(config *Config, logger *logrus.Logger) (PubSub, error) {
	switch config.Type {
	case "redis":
		return NewRedisPubSub(config, logger)
	case "memory":
		return NewMemoryPubSub(config, logger), nil
	default:
		return nil, fmt.Errorf("不支援的 Pub/Sub 類型: %s", config.Type)
	}
}
