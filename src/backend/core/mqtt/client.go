package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// Client MQTT 客戶端（簡化版）
type Client struct {
	broker *Broker
	logger *logrus.Logger
}

// NewClient 創建新的 MQTT 客戶端
func NewClient(config *Config, logger *logrus.Logger) (*Client, error) {
	broker, err := NewBroker(config, logger)
	if err != nil {
		return nil, err
	}

	return &Client{
		broker: broker,
		logger: logger,
	}, nil
}

// Connect 連接到 MQTT Broker
func (c *Client) Connect() error {
	return c.broker.Start()
}

// Disconnect 斷開連接
func (c *Client) Disconnect() {
	c.broker.Stop()
}

// Publish 發布訊息
func (c *Client) Publish(topic string, payload []byte) error {
	return c.broker.Publish(topic, payload, c.broker.config.DefaultQoS, false)
}

// PublishWithQoS 發布訊息（指定 QoS）
func (c *Client) PublishWithQoS(topic string, payload []byte, qos byte, retained bool) error {
	return c.broker.Publish(topic, payload, qos, retained)
}

// PublishJSON 發布 JSON 訊息
func (c *Client) PublishJSON(topic string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化 JSON 失敗: %w", err)
	}

	return c.Publish(topic, payload)
}

// Subscribe 訂閱主題
func (c *Client) Subscribe(topic string, handler func(topic string, payload []byte) error) error {
	return c.broker.Subscribe(topic, handler)
}

// Unsubscribe 取消訂閱
func (c *Client) Unsubscribe(topic string) error {
	return c.broker.Unsubscribe(topic)
}

// IsConnected 檢查是否已連接
func (c *Client) IsConnected() bool {
	return c.broker.IsConnected()
}

// WaitForConnection 等待連接建立
func (c *Client) WaitForConnection(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if c.IsConnected() {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("等待連接超時")
}
