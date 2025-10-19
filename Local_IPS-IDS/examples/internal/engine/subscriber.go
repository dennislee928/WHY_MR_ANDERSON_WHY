package engine

import (
	"context"
	"fmt"
	"sync"

	"pandora_box_console_ids_ips/internal/pubsub"

	"github.com/sirupsen/logrus"
)

// EventSubscriber handles subscribing to events from RabbitMQ
// 事件訂閱器，負責從 RabbitMQ 訂閱和處理事件
type EventSubscriber struct {
	mq     pubsub.MessageQueue
	logger *logrus.Logger
	engine *Engine
	mu     sync.RWMutex
	closed bool
}

// NewEventSubscriber creates a new event subscriber
// 創建新的事件訂閱器
func NewEventSubscriber(config *pubsub.Config, engine *Engine, logger *logrus.Logger) (*EventSubscriber, error) {
	if logger == nil {
		logger = logrus.New()
	}

	mq, err := pubsub.NewRabbitMQ(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ client: %w", err)
	}

	return &EventSubscriber{
		mq:     mq,
		logger: logger,
		engine: engine,
	}, nil
}

// Start starts subscribing to all event queues
// 開始訂閱所有事件隊列
func (s *EventSubscriber) Start(ctx context.Context) error {
	s.logger.Info("Starting event subscriber...")

	// Subscribe to threat events
	if err := s.subscribeThreatEvents(ctx); err != nil {
		return fmt.Errorf("failed to subscribe to threat events: %w", err)
	}

	// Subscribe to network events
	if err := s.subscribeNetworkEvents(ctx); err != nil {
		return fmt.Errorf("failed to subscribe to network events: %w", err)
	}

	// Subscribe to system events
	if err := s.subscribeSystemEvents(ctx); err != nil {
		return fmt.Errorf("failed to subscribe to system events: %w", err)
	}

	// Subscribe to device events
	if err := s.subscribeDeviceEvents(ctx); err != nil {
		return fmt.Errorf("failed to subscribe to device events: %w", err)
	}

	s.logger.Info("Event subscriber started successfully")
	return nil
}

// subscribeThreatEvents subscribes to threat events
// 訂閱威脅事件
func (s *EventSubscriber) subscribeThreatEvents(ctx context.Context) error {
	handler := func(ctx context.Context, msg *pubsub.Message) error {
		var event pubsub.ThreatEvent
		if err := pubsub.FromJSON(msg.Body, &event); err != nil {
			s.logger.Errorf("Failed to parse threat event: %v", err)
			return err
		}

		s.logger.Infof("Received threat event: %s from %s (level: %d)",
			event.ThreatType, event.SourceIP, event.ThreatLevel)

		// Process threat event
		return s.handleThreatEvent(ctx, &event)
	}

	return s.mq.Subscribe(ctx, "threat_events", handler)
}

// subscribeNetworkEvents subscribes to network events
// 訂閱網路事件
func (s *EventSubscriber) subscribeNetworkEvents(ctx context.Context) error {
	handler := func(ctx context.Context, msg *pubsub.Message) error {
		var event pubsub.NetworkEvent
		if err := pubsub.FromJSON(msg.Body, &event); err != nil {
			s.logger.Errorf("Failed to parse network event: %v", err)
			return err
		}

		s.logger.Infof("Received network event: %s from %s to %s",
			event.EventSubType, event.SourceIP, event.DestIP)

		// Process network event
		return s.handleNetworkEvent(ctx, &event)
	}

	return s.mq.Subscribe(ctx, "network_events", handler)
}

// subscribeSystemEvents subscribes to system events
// 訂閱系統事件
func (s *EventSubscriber) subscribeSystemEvents(ctx context.Context) error {
	handler := func(ctx context.Context, msg *pubsub.Message) error {
		var event pubsub.SystemEvent
		if err := pubsub.FromJSON(msg.Body, &event); err != nil {
			s.logger.Errorf("Failed to parse system event: %v", err)
			return err
		}

		s.logger.Infof("Received system event: %s - %s (%s)",
			event.Component, event.Status, event.Message)

		// Process system event
		return s.handleSystemEvent(ctx, &event)
	}

	return s.mq.Subscribe(ctx, "system_events", handler)
}

// subscribeDeviceEvents subscribes to device events
// 訂閱設備事件
func (s *EventSubscriber) subscribeDeviceEvents(ctx context.Context) error {
	handler := func(ctx context.Context, msg *pubsub.Message) error {
		var event pubsub.DeviceEvent
		if err := pubsub.FromJSON(msg.Body, &event); err != nil {
			s.logger.Errorf("Failed to parse device event: %v", err)
			return err
		}

		s.logger.Infof("Received device event: %s (%s) - %s",
			event.DeviceID, event.DeviceType, event.Status)

		// Process device event
		return s.handleDeviceEvent(ctx, &event)
	}

	return s.mq.Subscribe(ctx, "device_events", handler)
}

// Event handlers
// 事件處理函數

// handleThreatEvent processes a threat event
// 處理威脅事件
func (s *EventSubscriber) handleThreatEvent(ctx context.Context, event *pubsub.ThreatEvent) error {
	s.logger.Debugf("Processing threat event: %s", event.ID)

	// TODO: 實現威脅分析邏輯
	// 1. 分析威脅類型和等級
	// 2. 更新威脅資料庫
	// 3. 觸發告警（如果需要）
	// 4. 執行自動響應（如果配置）

	// 示例：記錄到資料庫
	if s.engine != nil {
		// s.engine.AnalyzeThreat(event)
		s.logger.Debugf("Threat analyzed: %s", event.ThreatType)
	}

	return nil
}

// handleNetworkEvent processes a network event
// 處理網路事件
func (s *EventSubscriber) handleNetworkEvent(ctx context.Context, event *pubsub.NetworkEvent) error {
	s.logger.Debugf("Processing network event: %s", event.ID)

	// TODO: 實現網路事件分析邏輯
	// 1. 分析網路流量模式
	// 2. 檢測異常行為
	// 3. 更新網路統計
	// 4. 觸發相關告警

	return nil
}

// handleSystemEvent processes a system event
// 處理系統事件
func (s *EventSubscriber) handleSystemEvent(ctx context.Context, event *pubsub.SystemEvent) error {
	s.logger.Debugf("Processing system event: %s", event.ID)

	// TODO: 實現系統事件處理邏輯
	// 1. 更新系統狀態
	// 2. 記錄系統指標
	// 3. 觸發健康檢查告警（如果需要）

	return nil
}

// handleDeviceEvent processes a device event
// 處理設備事件
func (s *EventSubscriber) handleDeviceEvent(ctx context.Context, event *pubsub.DeviceEvent) error {
	s.logger.Debugf("Processing device event: %s", event.ID)

	// TODO: 實現設備事件處理邏輯
	// 1. 更新設備狀態
	// 2. 記錄設備數據
	// 3. 觸發設備告警（如果需要）

	return nil
}

// Close closes the event subscriber
// 關閉事件訂閱器
func (s *EventSubscriber) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true

	if s.mq != nil {
		return s.mq.Close()
	}

	return nil
}

// Health checks if the subscriber is healthy
// 健康檢查
func (s *EventSubscriber) Health(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return fmt.Errorf("subscriber is closed")
	}

	return s.mq.Health(ctx)
}

// Engine is a placeholder for the actual engine implementation
// Engine 佔位符，實際實現需要根據現有代碼調整
type Engine struct {
	// Add engine fields here
}

