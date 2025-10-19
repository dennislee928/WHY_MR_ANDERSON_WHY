// Package pubsub provides message queue abstractions for Pandora Box Console IDS-IPS
// 基於專家反饋實現的消息隊列接口，用於解耦 Agent 和 Engine
package pubsub

import (
	"context"
	"time"
)

// MessageQueue defines the interface for message queue operations
// 消息隊列接口，支援發布/訂閱模式
type MessageQueue interface {
	// Publish sends a message to the specified exchange with a routing key
	// exchange: 交換機名稱（如 "pandora.events"）
	// routingKey: 路由鍵（如 "threat.detected", "network.attack"）
	// message: 消息內容（JSON 格式的字節數組）
	Publish(ctx context.Context, exchange, routingKey string, message []byte) error

	// Subscribe listens to messages from the specified queue
	// queue: 隊列名稱（如 "threat_events", "network_events"）
	// handler: 消息處理函數，返回 error 表示處理失敗（會重試）
	Subscribe(ctx context.Context, queue string, handler MessageHandler) error

	// Close gracefully shuts down the message queue connection
	// 優雅關閉連接，確保所有消息都已處理
	Close() error

	// Health checks if the message queue connection is healthy
	// 健康檢查，用於監控和告警
	Health(ctx context.Context) error
}

// PubSub defines the interface for publish-subscribe operations
// 發布訂閱操作接口
type PubSub interface {
	Publish(ctx context.Context, topic string, message *Message) error
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	Close() error
}

// MessageHandler is a function that processes incoming messages
// 消息處理函數類型
type MessageHandler func(topic string, message []byte) error

// Event represents an event in the system
// 系統事件結構
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// Message represents a message in the queue
// 消息結構，包含元數據和內容
type Message struct {
	// ID is the unique identifier for the message
	ID string

	// RoutingKey is the routing key used to publish the message
	RoutingKey string

	// Body is the message content (usually JSON)
	Body []byte

	// Timestamp is when the message was created
	Timestamp time.Time

	// Headers contains additional metadata
	Headers map[string]interface{}

	// DeliveryTag is used for acknowledging the message
	DeliveryTag uint64

	// Redelivered indicates if this message has been redelivered
	Redelivered bool
}

// PublishOptions contains options for publishing messages
// 發布選項，用於配置消息屬性
type PublishOptions struct {
	// ContentType of the message (default: "application/json")
	ContentType string

	// Priority of the message (0-9, higher is more important)
	Priority uint8

	// Expiration time in milliseconds (0 = no expiration)
	Expiration string

	// Persistent indicates if the message should be persisted to disk
	Persistent bool

	// Headers contains additional metadata
	Headers map[string]interface{}
}

// SubscribeOptions contains options for subscribing to messages
// 訂閱選項，用於配置消費者行為
type SubscribeOptions struct {
	// AutoAck automatically acknowledges messages (default: false)
	// 建議設為 false，手動確認以確保消息不丟失
	AutoAck bool

	// Exclusive indicates if this is an exclusive consumer
	Exclusive bool

	// PrefetchCount limits how many unacknowledged messages can be in flight
	// 預取數量，控制並發處理的消息數
	PrefetchCount int

	// RetryPolicy defines how to handle failed messages
	RetryPolicy *RetryPolicy
}

// RetryPolicy defines how to retry failed messages
// 重試策略，用於處理失敗的消息
type RetryPolicy struct {
	// MaxRetries is the maximum number of retry attempts
	MaxRetries int

	// InitialInterval is the initial retry delay
	InitialInterval time.Duration

	// MaxInterval is the maximum retry delay
	MaxInterval time.Duration

	// Multiplier is the backoff multiplier
	Multiplier float64
}

// Config contains configuration for the message queue
// 消息隊列配置
type Config struct {
	// Type is the pub/sub type (rabbitmq, redis, etc.)
	Type string

	// URL is the connection string (e.g., "amqp://user:pass@host:port/vhost")
	URL string

	// Exchange is the default exchange name
	Exchange string

	// RedisAddr is the Redis server address
	RedisAddr string

	// RedisPassword is the Redis password
	RedisPassword string

	// RedisDB is the Redis database number
	RedisDB int

	// BufferSize is the message buffer size
	BufferSize int

	// ConnectionTimeout is the timeout for establishing a connection
	ConnectionTimeout time.Duration

	// HeartbeatInterval is the interval for sending heartbeats
	HeartbeatInterval time.Duration

	// ReconnectDelay is the delay before attempting to reconnect
	ReconnectDelay time.Duration

	// MaxReconnectAttempts is the maximum number of reconnect attempts
	MaxReconnectAttempts int
}

// DefaultConfig returns a default configuration
// 默認配置，適用於開發環境
func DefaultConfig() *Config {
	return &Config{
		URL:                  "amqp://pandora:pandora123@localhost:5672/",
		Exchange:             "pandora.events",
		ConnectionTimeout:    30 * time.Second,
		HeartbeatInterval:    60 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
	}
}

// DefaultPublishOptions returns default publish options
// 默認發布選項
func DefaultPublishOptions() *PublishOptions {
	return &PublishOptions{
		ContentType: "application/json",
		Priority:    5,
		Persistent:  true,
		Headers:     make(map[string]interface{}),
	}
}

// DefaultSubscribeOptions returns default subscribe options
// 默認訂閱選項
func DefaultSubscribeOptions() *SubscribeOptions {
	return &SubscribeOptions{
		AutoAck:       false,
		Exclusive:     false,
		PrefetchCount: 10,
		RetryPolicy: &RetryPolicy{
			MaxRetries:      3,
			InitialInterval: 1 * time.Second,
			MaxInterval:     30 * time.Second,
			Multiplier:      2.0,
		},
	}
}

