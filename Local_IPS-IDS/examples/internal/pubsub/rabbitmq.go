package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ implements the MessageQueue interface using RabbitMQ
// RabbitMQ 實現，提供可靠的消息隊列功能
type RabbitMQ struct {
	config *Config
	conn   *amqp.Connection
	ch     *amqp.Channel
	mu     sync.RWMutex

	// reconnect control
	reconnecting bool
	closed       bool
	closeChan    chan struct{}
}

// NewRabbitMQ creates a new RabbitMQ message queue instance
// 創建新的 RabbitMQ 實例
func NewRabbitMQ(config *Config) (*RabbitMQ, error) {
	if config == nil {
		config = DefaultConfig()
	}

	mq := &RabbitMQ{
		config:    config,
		closeChan: make(chan struct{}),
	}

	if err := mq.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Start connection monitor
	go mq.monitorConnection()

	return mq, nil
}

// connect establishes a connection to RabbitMQ
// 建立與 RabbitMQ 的連接
func (mq *RabbitMQ) connect() error {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	// Create connection
	conn, err := amqp.DialConfig(mq.config.URL, amqp.Config{
		Heartbeat: mq.config.HeartbeatInterval,
		Locale:    "en_US",
	})
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	// Create channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		mq.config.Exchange, // name
		"topic",            // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	mq.conn = conn
	mq.ch = ch

	log.Printf("[RabbitMQ] Connected to %s", mq.config.URL)
	return nil
}

// monitorConnection monitors the connection and reconnects if necessary
// 監控連接狀態，斷線時自動重連
func (mq *RabbitMQ) monitorConnection() {
	for {
		select {
		case <-mq.closeChan:
			return
		case <-mq.conn.NotifyClose(make(chan *amqp.Error)):
			if mq.closed {
				return
			}
			log.Println("[RabbitMQ] Connection lost, attempting to reconnect...")
			mq.reconnect()
		}
	}
}

// reconnect attempts to reconnect to RabbitMQ
// 重新連接到 RabbitMQ
func (mq *RabbitMQ) reconnect() {
	mq.mu.Lock()
	if mq.reconnecting {
		mq.mu.Unlock()
		return
	}
	mq.reconnecting = true
	mq.mu.Unlock()

	defer func() {
		mq.mu.Lock()
		mq.reconnecting = false
		mq.mu.Unlock()
	}()

	for attempt := 1; attempt <= mq.config.MaxReconnectAttempts; attempt++ {
		if mq.closed {
			return
		}

		log.Printf("[RabbitMQ] Reconnect attempt %d/%d", attempt, mq.config.MaxReconnectAttempts)

		if err := mq.connect(); err != nil {
			log.Printf("[RabbitMQ] Reconnect failed: %v", err)
			time.Sleep(mq.config.ReconnectDelay * time.Duration(attempt))
			continue
		}

		log.Println("[RabbitMQ] Reconnected successfully")
		return
	}

	log.Println("[RabbitMQ] Failed to reconnect after maximum attempts")
}

// Publish sends a message to the specified exchange with a routing key
// 發布消息到指定的交換機
func (mq *RabbitMQ) Publish(ctx context.Context, exchange, routingKey string, message []byte) error {
	mq.mu.RLock()
	defer mq.mu.RUnlock()

	if mq.ch == nil {
		return fmt.Errorf("channel is not initialized")
	}

	opts := DefaultPublishOptions()

	return mq.ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  opts.ContentType,
			Body:         message,
			DeliveryMode: amqp.Persistent,
			Priority:     opts.Priority,
			Timestamp:    time.Now(),
			Headers:      amqp.Table(opts.Headers),
		},
	)
}

// Subscribe listens to messages from the specified queue
// 訂閱指定隊列的消息
func (mq *RabbitMQ) Subscribe(ctx context.Context, queue string, handler MessageHandler) error {
	mq.mu.RLock()
	ch := mq.ch
	mq.mu.RUnlock()

	if ch == nil {
		return fmt.Errorf("channel is not initialized")
	}

	opts := DefaultSubscribeOptions()

	// Set QoS
	err := ch.Qos(
		opts.PrefetchCount, // prefetch count
		0,                  // prefetch size
		false,              // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Start consuming
	msgs, err := ch.Consume(
		queue,
		"",                // consumer tag
		opts.AutoAck,      // auto-ack
		opts.Exclusive,    // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	// Process messages
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-mq.closeChan:
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				message := &Message{
					ID:          msg.MessageId,
					RoutingKey:  msg.RoutingKey,
					Body:        msg.Body,
					Timestamp:   msg.Timestamp,
					Headers:     make(map[string]interface{}),
					DeliveryTag: msg.DeliveryTag,
					Redelivered: msg.Redelivered,
				}

				// Convert AMQP headers to map
				for k, v := range msg.Headers {
					message.Headers[k] = v
				}

				// Handle message
				if err := handler(ctx, message); err != nil {
					log.Printf("[RabbitMQ] Error handling message: %v", err)
					// Nack the message for retry
					msg.Nack(false, true)
				} else {
					// Ack the message
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

// Close gracefully shuts down the message queue connection
// 優雅關閉連接
func (mq *RabbitMQ) Close() error {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if mq.closed {
		return nil
	}

	mq.closed = true
	close(mq.closeChan)

	if mq.ch != nil {
		if err := mq.ch.Close(); err != nil {
			log.Printf("[RabbitMQ] Error closing channel: %v", err)
		}
	}

	if mq.conn != nil {
		if err := mq.conn.Close(); err != nil {
			log.Printf("[RabbitMQ] Error closing connection: %v", err)
		}
	}

	log.Println("[RabbitMQ] Connection closed")
	return nil
}

// Health checks if the message queue connection is healthy
// 健康檢查
func (mq *RabbitMQ) Health(ctx context.Context) error {
	mq.mu.RLock()
	defer mq.mu.RUnlock()

	if mq.conn == nil || mq.conn.IsClosed() {
		return fmt.Errorf("connection is closed")
	}

	if mq.ch == nil {
		return fmt.Errorf("channel is not initialized")
	}

	return nil
}

// PublishJSON is a helper function to publish JSON messages
// 發布 JSON 消息的輔助函數
func (mq *RabbitMQ) PublishJSON(ctx context.Context, exchange, routingKey string, data interface{}) error {
	message, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return mq.Publish(ctx, exchange, routingKey, message)
}

