package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pandora_box_console_ids_ips/internal/pubsub"
	"github.com/sirupsen/logrus"
)

// EventProcessor 事件處理器，模擬 Axiom Engine
type EventProcessor struct {
	logger *logrus.Logger
	mq     *pubsub.RabbitMQ
}

// NewEventProcessor 創建新的事件處理器
func NewEventProcessor(mq *pubsub.RabbitMQ, logger *logrus.Logger) *EventProcessor {
	return &EventProcessor{
		logger: logger,
		mq:     mq,
	}
}

// ProcessThreatEvent 處理威脅事件
func (ep *EventProcessor) ProcessThreatEvent(routingKey string, message []byte) error {
	var event pubsub.ThreatEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal threat event: %w", err)
	}

	ep.logger.WithFields(logrus.Fields{
		"event_type":   event.Type,
		"source_ip":    event.SourceIP,
		"threat_level": event.ThreatLevel,
		"action":       event.Action,
	}).Info("處理威脅事件")

	// 模擬威脅分析邏輯
	switch event.Type {
	case "ddos":
		ep.logger.Warn("檢測到 DDoS 攻擊，啟動防護機制")
		// 這裡可以調用實際的防護邏輯
	case "port_scan":
		ep.logger.Warn("檢測到端口掃描，記錄並分析")
		// 這裡可以調用實際的分析邏輯
	case "brute_force":
		ep.logger.Warn("檢測到暴力破解攻擊，暫時封鎖 IP")
		// 這裡可以調用實際的封鎖邏輯
	default:
		ep.logger.Info("處理其他類型威脅")
	}

	return nil
}

// ProcessNetworkEvent 處理網路事件
func (ep *EventProcessor) ProcessNetworkEvent(routingKey string, message []byte) error {
	var event pubsub.NetworkEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal network event: %w", err)
	}

	ep.logger.WithFields(logrus.Fields{
		"event_type":    event.Type,
		"source_ip":     event.SourceIP,
		"dest_ip":       event.DestIP,
		"protocol":      event.Protocol,
		"bytes":         event.Bytes,
	}).Info("處理網路事件")

	// 模擬網路分析邏輯
	if event.Bytes > 1000000 { // 大於 1MB
		ep.logger.Warn("檢測到大流量傳輸，需要進一步分析")
	}

	return nil
}

// ProcessSystemEvent 處理系統事件
func (ep *EventProcessor) ProcessSystemEvent(routingKey string, message []byte) error {
	var event pubsub.SystemEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal system event: %w", err)
	}

	ep.logger.WithFields(logrus.Fields{
		"event_type": event.Type,
		"component":  event.Component,
		"status":     event.Status,
		"message":    event.Message,
	}).Info("處理系統事件")

	return nil
}

// ProcessDeviceEvent 處理設備事件
func (ep *EventProcessor) ProcessDeviceEvent(routingKey string, message []byte) error {
	var event pubsub.DeviceEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return fmt.Errorf("failed to unmarshal device event: %w", err)
	}

	ep.logger.WithFields(logrus.Fields{
		"event_type": event.Type,
		"device_id":  event.DeviceID,
		"device_type": event.DeviceType,
		"status":     event.Status,
	}).Info("處理設備事件")

	return nil
}

// Start 啟動事件處理器
func (ep *EventProcessor) Start(ctx context.Context) error {
	ep.logger.Info("啟動 Axiom Engine 事件處理器...")

	// 訂閱威脅事件
	if err := ep.mq.Subscribe(ctx, "threat_events", ep.ProcessThreatEvent); err != nil {
		return fmt.Errorf("failed to subscribe to threat events: %w", err)
	}

	// 訂閱網路事件
	if err := ep.mq.Subscribe(ctx, "network_events", ep.ProcessNetworkEvent); err != nil {
		return fmt.Errorf("failed to subscribe to network events: %w", err)
	}

	// 訂閱系統事件
	if err := ep.mq.Subscribe(ctx, "system_events", ep.ProcessSystemEvent); err != nil {
		return fmt.Errorf("failed to subscribe to system events: %w", err)
	}

	// 訂閱設備事件
	if err := ep.mq.Subscribe(ctx, "device_events", ep.ProcessDeviceEvent); err != nil {
		return fmt.Errorf("failed to subscribe to device events: %w", err)
	}

	ep.logger.Info("Axiom Engine 已訂閱所有事件隊列")
	return nil
}

// EventPublisher 事件發布器，模擬 Pandora Agent
type EventPublisher struct {
	logger *logrus.Logger
	mq     *pubsub.RabbitMQ
}

// NewEventPublisher 創建新的事件發布器
func NewEventPublisher(mq *pubsub.RabbitMQ, logger *logrus.Logger) *EventPublisher {
	return &EventPublisher{
		logger: logger,
		mq:     mq,
	}
}

// PublishThreatEvent 發布威脅事件
func (ep *EventPublisher) PublishThreatEvent(threatType, sourceIP, action string, threatLevel int) error {
	event := pubsub.NewThreatEvent(threatType, sourceIP, "威脅檢測", action, threatLevel)
	event.TargetIP = "10.0.0.1"
	event.TargetPort = 80

	message, err := pubsub.ToJSON(event)
	if err != nil {
		return fmt.Errorf("failed to marshal threat event: %w", err)
	}

	routingKey := fmt.Sprintf("threat.%s", threatType)
	if err := ep.mq.Publish(context.Background(), "pandora.events", routingKey, message); err != nil {
		return fmt.Errorf("failed to publish threat event: %w", err)
	}

	ep.logger.WithFields(logrus.Fields{
		"threat_type": threatType,
		"source_ip":   sourceIP,
		"routing_key": routingKey,
	}).Info("發布威脅事件")

	return nil
}

// PublishNetworkEvent 發布網路事件
func (ep *EventPublisher) PublishNetworkEvent(eventType, sourceIP, destIP, protocol string, bytes int64) error {
	event := pubsub.NewNetworkEvent(eventType, sourceIP, destIP, protocol, bytes)

	message, err := pubsub.ToJSON(event)
	if err != nil {
		return fmt.Errorf("failed to marshal network event: %w", err)
	}

	routingKey := fmt.Sprintf("network.%s", eventType)
	if err := ep.mq.Publish(context.Background(), "pandora.events", routingKey, message); err != nil {
		return fmt.Errorf("failed to publish network event: %w", err)
	}

	ep.logger.WithFields(logrus.Fields{
		"event_type": eventType,
		"source_ip":  sourceIP,
		"dest_ip":    destIP,
		"protocol":   protocol,
		"bytes":      bytes,
	}).Info("發布網路事件")

	return nil
}

// PublishSystemEvent 發布系統事件
func (ep *EventPublisher) PublishSystemEvent(eventType, component, status, message string) error {
	event := pubsub.NewSystemEvent(eventType, component, status, message)

	msgBytes, err := pubsub.ToJSON(event)
	if err != nil {
		return fmt.Errorf("failed to marshal system event: %w", err)
	}

	routingKey := fmt.Sprintf("system.%s", eventType)
	if err := ep.mq.Publish(context.Background(), "pandora.events", routingKey, msgBytes); err != nil {
		return fmt.Errorf("failed to publish system event: %w", err)
	}

	ep.logger.WithFields(logrus.Fields{
		"event_type": eventType,
		"component":  component,
		"status":     status,
	}).Info("發布系統事件")

	return nil
}

// PublishDeviceEvent 發布設備事件
func (ep *EventPublisher) PublishDeviceEvent(eventType, deviceID, deviceType, status string) error {
	event := pubsub.NewDeviceEvent(eventType, deviceID, deviceType, status)

	message, err := pubsub.ToJSON(event)
	if err != nil {
		return fmt.Errorf("failed to marshal device event: %w", err)
	}

	routingKey := fmt.Sprintf("device.%s", eventType)
	if err := ep.mq.Publish(context.Background(), "pandora.events", routingKey, message); err != nil {
		return fmt.Errorf("failed to publish device event: %w", err)
	}

	ep.logger.WithFields(logrus.Fields{
		"event_type":  eventType,
		"device_id":   deviceID,
		"device_type": deviceType,
		"status":      status,
	}).Info("發布設備事件")

	return nil
}

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logger.Info("=== Pandora Box Console IDS-IPS RabbitMQ 整合示範 ===")

	// 創建 RabbitMQ 配置
	config := &pubsub.Config{
		URL:                  getEnv("RABBITMQ_URL", "amqp://pandora:pandora123@localhost:5672/"),
		Exchange:             getEnv("RABBITMQ_EXCHANGE", "pandora.events"),
		ConnectionTimeout:    30 * time.Second,
		HeartbeatInterval:    60 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
	}

	// 創建 RabbitMQ 連接
	mq, err := pubsub.NewRabbitMQ(config)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ connection: %v", err)
	}
	defer mq.Close()

	logger.Info("RabbitMQ 連接已建立")

	// 創建事件處理器（Axiom Engine）
	processor := NewEventProcessor(mq, logger)

	// 創建事件發布器（Pandora Agent）
	publisher := NewEventPublisher(mq, logger)

	// 創建 context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動事件處理器
	if err := processor.Start(ctx); err != nil {
		log.Fatalf("Failed to start event processor: %v", err)
	}

	logger.Info("Axiom Engine 事件處理器已啟動")

	// 等待一秒讓訂閱建立
	time.Sleep(1 * time.Second)

	// 發布測試事件
	logger.Info("開始發布測試事件...")

	// 發布威脅事件
	if err := publisher.PublishThreatEvent("ddos", "192.168.1.100", "blocked", 8); err != nil {
		logger.Errorf("Failed to publish threat event: %v", err)
	}

	if err := publisher.PublishThreatEvent("port_scan", "192.168.1.101", "detected", 6); err != nil {
		logger.Errorf("Failed to publish threat event: %v", err)
	}

	if err := publisher.PublishThreatEvent("brute_force", "192.168.1.102", "blocked", 9); err != nil {
		logger.Errorf("Failed to publish threat event: %v", err)
	}

	// 發布網路事件
	if err := publisher.PublishNetworkEvent("connection", "192.168.1.50", "10.0.0.1", "TCP", 1500); err != nil {
		logger.Errorf("Failed to publish network event: %v", err)
	}

	if err := publisher.PublishNetworkEvent("large_transfer", "192.168.1.51", "10.0.0.2", "TCP", 2000000); err != nil {
		logger.Errorf("Failed to publish network event: %v", err)
	}

	// 發布系統事件
	if err := publisher.PublishSystemEvent("service_start", "axiom-ui", "running", "Axiom UI 服務已啟動"); err != nil {
		logger.Errorf("Failed to publish system event: %v", err)
	}

	if err := publisher.PublishSystemEvent("config_update", "prometheus", "updated", "Prometheus 配置已更新"); err != nil {
		logger.Errorf("Failed to publish system event: %v", err)
	}

	// 發布設備事件
	if err := publisher.PublishDeviceEvent("device_online", "device_001", "serial", "online"); err != nil {
		logger.Errorf("Failed to publish device event: %v", err)
	}

	if err := publisher.PublishDeviceEvent("device_error", "device_002", "network", "error"); err != nil {
		logger.Errorf("Failed to publish device event: %v", err)
	}

	logger.Info("所有測試事件已發布")

	// 等待中斷信號
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("按 Ctrl+C 停止...")
	logger.Info("")
	logger.Info("=== 事件流架構 ===")
	logger.Info("Pandora Agent ──發布事件──> RabbitMQ ──路由事件──> Axiom Engine")
	logger.Info("")
	logger.Info("交換機: pandora.events (Topic)")
	logger.Info("隊列:")
	logger.Info("  - threat_events (routing: threat.*)")
	logger.Info("  - network_events (routing: network.*)")
	logger.Info("  - system_events (routing: system.*)")
	logger.Info("  - device_events (routing: device.*)")
	logger.Info("")

	<-sigChan

	logger.Info("正在關閉...")
	cancel()

	// 等待清理完成
	time.Sleep(2 * time.Second)
	logger.Info("RabbitMQ 整合示範已停止")
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
