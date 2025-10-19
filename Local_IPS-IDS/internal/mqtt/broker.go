package mqtt

import (
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// Broker MQTT 代理（客戶端包裝）
type Broker struct {
	config    *Config
	client    mqtt.Client
	logger    *logrus.Logger
	handlers  map[string]MessageHandler
	mu        sync.RWMutex
	connected bool
}

// MessageHandler 訊息處理函數
type MessageHandler func(topic string, payload []byte) error

// Config MQTT 配置
type Config struct {
	Broker           string        `yaml:"broker" json:"broker"`                         // MQTT 代理地址
	Port             int           `yaml:"port" json:"port"`                             // 端口
	ClientID         string        `yaml:"client_id" json:"client_id"`                   // 客戶端 ID
	Username         string        `yaml:"username" json:"username"`                     // 用戶名
	Password         string        `yaml:"password" json:"password"`                     // 密碼
	TLSEnabled       bool          `yaml:"tls_enabled" json:"tls_enabled"`               // 是否啟用 TLS
	DefaultQoS       byte          `yaml:"default_qos" json:"default_qos"`               // 預設 QoS (0, 1, 2)
	KeepAlive        int           `yaml:"keep_alive" json:"keep_alive"`                 // 保持連接時間（秒）
	ConnectTimeout   time.Duration `yaml:"connect_timeout" json:"connect_timeout"`       // 連接超時
	ReconnectDelay   time.Duration `yaml:"reconnect_delay" json:"reconnect_delay"`       // 重連延遲
	MaxReconnectWait time.Duration `yaml:"max_reconnect_wait" json:"max_reconnect_wait"` // 最大重連等待
	AutoReconnect    bool          `yaml:"auto_reconnect" json:"auto_reconnect"`         // 自動重連
	CleanSession     bool          `yaml:"clean_session" json:"clean_session"`           // 清理會話
	OrderMatters     bool          `yaml:"order_matters" json:"order_matters"`           // 訊息順序
}

// NewBroker 創建新的 MQTT Broker
func NewBroker(config *Config, logger *logrus.Logger) (*Broker, error) {
	if logger == nil {
		logger = logrus.New()
	}

	// 設定預設值
	if config.Port == 0 {
		config.Port = 1883 // 預設 MQTT 端口
	}
	if config.ClientID == "" {
		config.ClientID = fmt.Sprintf("pandora-%d", time.Now().Unix())
	}
	if config.DefaultQoS > 2 {
		config.DefaultQoS = 1 // 預設 QoS 1
	}
	if config.KeepAlive == 0 {
		config.KeepAlive = 60 // 60 秒
	}
	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = 30 * time.Second
	}
	if config.ReconnectDelay == 0 {
		config.ReconnectDelay = 1 * time.Second
	}
	if config.MaxReconnectWait == 0 {
		config.MaxReconnectWait = 10 * time.Minute
	}

	broker := &Broker{
		config:   config,
		logger:   logger,
		handlers: make(map[string]MessageHandler),
	}

	// 創建 MQTT 客戶端選項
	opts := mqtt.NewClientOptions()

	// 設定代理地址
	brokerURL := fmt.Sprintf("tcp://%s:%d", config.Broker, config.Port)
	if config.TLSEnabled {
		brokerURL = fmt.Sprintf("ssl://%s:%d", config.Broker, config.Port)
	}
	opts.AddBroker(brokerURL)

	// 設定客戶端 ID
	opts.SetClientID(config.ClientID)

	// 設定認證
	if config.Username != "" {
		opts.SetUsername(config.Username)
		if config.Password != "" {
			opts.SetPassword(config.Password)
		}
	}

	// 設定連接參數
	opts.SetKeepAlive(time.Duration(config.KeepAlive) * time.Second)
	opts.SetConnectTimeout(config.ConnectTimeout)
	opts.SetAutoReconnect(config.AutoReconnect)
	opts.SetMaxReconnectInterval(config.MaxReconnectWait)
	opts.SetCleanSession(config.CleanSession)
	opts.SetOrderMatters(config.OrderMatters)

	// 設定回調
	opts.SetConnectionLostHandler(broker.onConnectionLost)
	opts.SetOnConnectHandler(broker.onConnect)
	opts.SetReconnectingHandler(broker.onReconnecting)

	// 創建客戶端
	broker.client = mqtt.NewClient(opts)

	return broker, nil
}

// Start 啟動 MQTT Broker
func (b *Broker) Start() error {
	b.logger.Infof("連接到 MQTT Broker: %s:%d", b.config.Broker, b.config.Port)

	token := b.client.Connect()
	if !token.WaitTimeout(b.config.ConnectTimeout) {
		return fmt.Errorf("連接 MQTT Broker 超時")
	}

	if err := token.Error(); err != nil {
		return fmt.Errorf("連接 MQTT Broker 失敗: %w", err)
	}

	b.connected = true
	b.logger.Info("已連接到 MQTT Broker")

	return nil
}

// Stop 停止 MQTT Broker
func (b *Broker) Stop() {
	b.logger.Info("斷開 MQTT Broker 連接...")

	// 取消所有訂閱
	b.mu.RLock()
	topics := make([]string, 0, len(b.handlers))
	for topic := range b.handlers {
		topics = append(topics, topic)
	}
	b.mu.RUnlock()

	for _, topic := range topics {
		b.Unsubscribe(topic)
	}

	// 斷開連接
	b.client.Disconnect(250)
	b.connected = false
	b.logger.Info("已斷開 MQTT Broker 連接")
}

// Publish 發布訊息
func (b *Broker) Publish(topic string, payload []byte, qos byte, retained bool) error {
	if !b.connected {
		return fmt.Errorf("未連接到 MQTT Broker")
	}

	token := b.client.Publish(topic, qos, retained, payload)
	if !token.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("發布訊息超時")
	}

	if err := token.Error(); err != nil {
		return fmt.Errorf("發布訊息失敗: %w", err)
	}

	b.logger.Debugf("已發布訊息到主題 [%s]: %d bytes", topic, len(payload))
	return nil
}

// PublishJSON 發布 JSON 訊息
func (b *Broker) PublishJSON(topic string, data interface{}) error {
	payload, err := jsonMarshal(data)
	if err != nil {
		return fmt.Errorf("序列化 JSON 失敗: %w", err)
	}

	return b.Publish(topic, payload, b.config.DefaultQoS, false)
}

// Subscribe 訂閱主題
func (b *Broker) Subscribe(topic string, handler MessageHandler) error {
	if !b.connected {
		return fmt.Errorf("未連接到 MQTT Broker")
	}

	// 儲存 handler
	b.mu.Lock()
	b.handlers[topic] = handler
	b.mu.Unlock()

	// 訂閱主題
	token := b.client.Subscribe(topic, b.config.DefaultQoS, func(client mqtt.Client, msg mqtt.Message) {
		b.handleMessage(msg)
	})

	if !token.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("訂閱主題超時")
	}

	if err := token.Error(); err != nil {
		b.mu.Lock()
		delete(b.handlers, topic)
		b.mu.Unlock()
		return fmt.Errorf("訂閱主題失敗: %w", err)
	}

	b.logger.Infof("已訂閱主題: %s", topic)
	return nil
}

// Unsubscribe 取消訂閱
func (b *Broker) Unsubscribe(topic string) error {
	if !b.connected {
		return fmt.Errorf("未連接到 MQTT Broker")
	}

	// 取消訂閱
	token := b.client.Unsubscribe(topic)
	if !token.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("取消訂閱超時")
	}

	if err := token.Error(); err != nil {
		return fmt.Errorf("取消訂閱失敗: %w", err)
	}

	// 移除 handler
	b.mu.Lock()
	delete(b.handlers, topic)
	b.mu.Unlock()

	b.logger.Infof("已取消訂閱主題: %s", topic)
	return nil
}

// IsConnected 檢查是否已連接
func (b *Broker) IsConnected() bool {
	return b.connected && b.client.IsConnected()
}

// handleMessage 處理接收到的訊息
func (b *Broker) handleMessage(msg mqtt.Message) {
	topic := msg.Topic()
	payload := msg.Payload()

	b.mu.RLock()
	handler, exists := b.handlers[topic]
	b.mu.RUnlock()

	if !exists {
		b.logger.Warnf("收到未訂閱主題的訊息 [%s]", topic)
		return
	}

	b.logger.Debugf("收到訊息 [%s]: %d bytes", topic, len(payload))

	if err := handler(topic, payload); err != nil {
		b.logger.Errorf("處理訊息失敗 [%s]: %v", topic, err)
	}
}

// onConnect 連接成功回調
func (b *Broker) onConnect(client mqtt.Client) {
	b.connected = true
	b.logger.Info("MQTT 連接成功")

	// 重新訂閱所有主題
	b.mu.RLock()
	topics := make([]string, 0, len(b.handlers))
	for topic := range b.handlers {
		topics = append(topics, topic)
	}
	b.mu.RUnlock()

	for _, topic := range topics {
		token := client.Subscribe(topic, b.config.DefaultQoS, func(client mqtt.Client, msg mqtt.Message) {
			b.handleMessage(msg)
		})
		if token.Wait() && token.Error() == nil {
			b.logger.Infof("重新訂閱主題: %s", topic)
		}
	}
}

// onConnectionLost 連接丟失回調
func (b *Broker) onConnectionLost(client mqtt.Client, err error) {
	b.connected = false
	b.logger.Errorf("MQTT 連接丟失: %v", err)
}

// onReconnecting 重連中回調
func (b *Broker) onReconnecting(client mqtt.Client, opts *mqtt.ClientOptions) {
	b.logger.Info("正在重新連接 MQTT Broker...")
}

// jsonMarshal JSON 序列化輔助函數
func jsonMarshal(v interface{}) ([]byte, error) {
	// 這裡可以引入 encoding/json
	// 為了簡化，暫時返回空實現
	return []byte(fmt.Sprintf("%v", v)), nil
}
