package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"pandora_box_console_ids_ips/internal/pubsub"

	"github.com/sirupsen/logrus"
)

// EventPublisher handles publishing events to RabbitMQ
// 事件發布器，負責將 Agent 事件發布到 RabbitMQ
type EventPublisher struct {
	mq     pubsub.MessageQueue
	logger *logrus.Logger
	mu     sync.RWMutex
	closed bool
}

// NewEventPublisher creates a new event publisher
// 創建新的事件發布器
func NewEventPublisher(config *pubsub.Config, logger *logrus.Logger) (*EventPublisher, error) {
	if logger == nil {
		logger = logrus.New()
	}

	mq, err := pubsub.NewRabbitMQ(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ client: %w", err)
	}

	return &EventPublisher{
		mq:     mq,
		logger: logger,
	}, nil
}

// PublishThreatEvent publishes a threat event
// 發布威脅事件
func (p *EventPublisher) PublishThreatEvent(ctx context.Context, event *pubsub.ThreatEvent) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return fmt.Errorf("publisher is closed")
	}

	message, err := pubsub.ToJSON(event)
	if err != nil {
		p.logger.Errorf("Failed to marshal threat event: %v", err)
		return err
	}

	routingKey := pubsub.GetRoutingKey(event.Type)
	err = p.mq.Publish(ctx, "pandora.events", routingKey, message)
	if err != nil {
		p.logger.Errorf("Failed to publish threat event: %v", err)
		return err
	}

	p.logger.Debugf("Published threat event: %s (type: %s, level: %d)",
		event.ID, event.ThreatType, event.ThreatLevel)
	return nil
}

// PublishNetworkEvent publishes a network event
// 發布網路事件
func (p *EventPublisher) PublishNetworkEvent(ctx context.Context, event *pubsub.NetworkEvent) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return fmt.Errorf("publisher is closed")
	}

	message, err := pubsub.ToJSON(event)
	if err != nil {
		p.logger.Errorf("Failed to marshal network event: %v", err)
		return err
	}

	routingKey := pubsub.GetRoutingKey(event.Type)
	err = p.mq.Publish(ctx, "pandora.events", routingKey, message)
	if err != nil {
		p.logger.Errorf("Failed to publish network event: %v", err)
		return err
	}

	p.logger.Debugf("Published network event: %s (subtype: %s)",
		event.ID, event.EventSubType)
	return nil
}

// PublishSystemEvent publishes a system event
// 發布系統事件
func (p *EventPublisher) PublishSystemEvent(ctx context.Context, event *pubsub.SystemEvent) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return fmt.Errorf("publisher is closed")
	}

	message, err := pubsub.ToJSON(event)
	if err != nil {
		p.logger.Errorf("Failed to marshal system event: %v", err)
		return err
	}

	routingKey := pubsub.GetRoutingKey(event.Type)
	err = p.mq.Publish(ctx, "pandora.events", routingKey, message)
	if err != nil {
		p.logger.Errorf("Failed to publish system event: %v", err)
		return err
	}

	p.logger.Debugf("Published system event: %s (component: %s, status: %s)",
		event.ID, event.Component, event.Status)
	return nil
}

// PublishDeviceEvent publishes a device event
// 發布設備事件
func (p *EventPublisher) PublishDeviceEvent(ctx context.Context, event *pubsub.DeviceEvent) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return fmt.Errorf("publisher is closed")
	}

	message, err := pubsub.ToJSON(event)
	if err != nil {
		p.logger.Errorf("Failed to marshal device event: %v", err)
		return err
	}

	routingKey := pubsub.GetRoutingKey(event.Type)
	err = p.mq.Publish(ctx, "pandora.events", routingKey, message)
	if err != nil {
		p.logger.Errorf("Failed to publish device event: %v", err)
		return err
	}

	p.logger.Debugf("Published device event: %s (device: %s, status: %s)",
		event.ID, event.DeviceID, event.Status)
	return nil
}

// Close closes the event publisher
// 關閉事件發布器
func (p *EventPublisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true

	if p.mq != nil {
		return p.mq.Close()
	}

	return nil
}

// Health checks if the publisher is healthy
// 健康檢查
func (p *EventPublisher) Health(ctx context.Context) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return fmt.Errorf("publisher is closed")
	}

	return p.mq.Health(ctx)
}

// Helper functions for common event scenarios
// 常見事件場景的輔助函數

// PublishThreatDetected publishes a threat detected event
// 發布威脅檢測事件
func (p *EventPublisher) PublishThreatDetected(ctx context.Context, threatType, sourceIP, description string, level int) error {
	event := pubsub.NewThreatEvent(threatType, sourceIP, description, "detected", level)
	event.Type = pubsub.EventTypeThreatDetected
	return p.PublishThreatEvent(ctx, event)
}

// PublishThreatBlocked publishes a threat blocked event
// 發布威脅阻斷事件
func (p *EventPublisher) PublishThreatBlocked(ctx context.Context, threatType, sourceIP, description string, level int) error {
	event := pubsub.NewThreatEvent(threatType, sourceIP, description, "blocked", level)
	event.Type = pubsub.EventTypeThreatBlocked
	return p.PublishThreatEvent(ctx, event)
}

// PublishNetworkAttack publishes a network attack event
// 發布網路攻擊事件
func (p *EventPublisher) PublishNetworkAttack(ctx context.Context, attackType, sourceIP, destIP, protocol string) error {
	event := pubsub.NewNetworkEvent(attackType, sourceIP, destIP, protocol)
	event.Type = pubsub.EventTypeNetworkAttack
	return p.PublishNetworkEvent(ctx, event)
}

// PublishAgentStarted publishes an agent started event
// 發布 Agent 啟動事件
func (p *EventPublisher) PublishAgentStarted(ctx context.Context) error {
	event := pubsub.NewSystemEvent("pandora-agent", "running", "Agent started successfully")
	event.Type = pubsub.EventTypeSystemStarted
	return p.PublishSystemEvent(ctx, event)
}

// PublishAgentStopped publishes an agent stopped event
// 發布 Agent 停止事件
func (p *EventPublisher) PublishAgentStopped(ctx context.Context) error {
	event := pubsub.NewSystemEvent("pandora-agent", "stopped", "Agent stopped gracefully")
	event.Type = pubsub.EventTypeSystemStopped
	return p.PublishSystemEvent(ctx, event)
}

// PublishAgentError publishes an agent error event
// 發布 Agent 錯誤事件
func (p *EventPublisher) PublishAgentError(ctx context.Context, errorMsg string) error {
	event := pubsub.NewSystemEvent("pandora-agent", "error", errorMsg)
	event.Type = pubsub.EventTypeSystemError
	event.Severity = "high"
	event.ErrorDetails = errorMsg
	return p.PublishSystemEvent(ctx, event)
}

// PublishDeviceConnected publishes a device connected event
// 發布設備連接事件
func (p *EventPublisher) PublishDeviceConnected(ctx context.Context, deviceID, deviceType, port string) error {
	event := pubsub.NewDeviceEvent(deviceID, deviceType, "connected")
	event.Type = pubsub.EventTypeDeviceConnected
	event.Port = port
	return p.PublishDeviceEvent(ctx, event)
}

// PublishDeviceDisconnected publishes a device disconnected event
// 發布設備斷開事件
func (p *EventPublisher) PublishDeviceDisconnected(ctx context.Context, deviceID, deviceType, port string) error {
	event := pubsub.NewDeviceEvent(deviceID, deviceType, "disconnected")
	event.Type = pubsub.EventTypeDeviceDisconnect
	event.Port = port
	return p.PublishDeviceEvent(ctx, event)
}

// PublishHealthCheck publishes a periodic health check event
// 發布定期健康檢查事件
func (p *EventPublisher) PublishHealthCheck(ctx context.Context, metrics map[string]interface{}) error {
	event := pubsub.NewSystemEvent("pandora-agent", "healthy", "Health check passed")
	event.Type = pubsub.EventTypeSystemHealthy
	event.Metrics = metrics
	return p.PublishSystemEvent(ctx, event)
}

// StartPeriodicHealthCheck starts publishing periodic health check events
// 啟動定期健康檢查事件發布
func (p *EventPublisher) StartPeriodicHealthCheck(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metrics := map[string]interface{}{
				"timestamp": time.Now().Unix(),
				"uptime":    time.Since(time.Now()).Seconds(), // 需要實際的啟動時間
			}

			if err := p.PublishHealthCheck(ctx, metrics); err != nil {
				p.logger.Warnf("Failed to publish health check: %v", err)
			}
		}
	}
}

