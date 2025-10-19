package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"pandora_box_console_ids_ips/internal/logging"
	"pandora_box_console_ids_ips/internal/metrics"
	"pandora_box_console_ids_ips/internal/mqtt"
	"pandora_box_console_ids_ips/internal/pubsub"
	"pandora_box_console_ids_ips/internal/ratelimit"
)

// AuthHandlerExtended 擴展的認證處理器（包含新模組）
type AuthHandlerExtended struct {
	*AuthHandler // 嵌入原始處理器

	// 新增模組
	rateLimiter *ratelimit.Middleware
	pubsub      pubsub.PubSub
	mqttBroker  *mqtt.Broker
}

// NewAuthHandlerWithModules 建立包含新模組的認證處理器
// 注意：Rate Limiting 應該在路由層作為 middleware 使用，而不是在處理器內部
func NewAuthHandlerWithModules(
	logger *logrus.Logger,
	centralLogger *logging.CentralLogger,
	metrics *metrics.PrometheusMetrics,
	rateLimiter *ratelimit.Middleware,
	pubsubInstance pubsub.PubSub,
	mqttBroker *mqtt.Broker,
) (*AuthHandler, *AuthHandlerExtended) {
	// 建立基礎處理器
	baseHandler := NewAuthHandler(logger, centralLogger, metrics)

	// 建立擴展處理器
	extendedHandler := &AuthHandlerExtended{
		AuthHandler: baseHandler,
		rateLimiter: rateLimiter,
		pubsub:      pubsubInstance,
		mqttBroker:  mqttBroker,
	}

	// 如果沒有新模組，直接返回基礎處理器
	if rateLimiter == nil && pubsubInstance == nil && mqttBroker == nil {
		return baseHandler, extendedHandler
	}

	// 設定 Pub/Sub 事件發布
	if pubsubInstance != nil {
		// 訂閱認證事件
		ctx := context.Background()
		pubsubInstance.Subscribe(ctx, "auth.events", func(topic string, message []byte) error {
			logger.Infof("收到認證事件: %s", string(message))
			return nil
		})

		// 訂閱安全事件
		pubsubInstance.Subscribe(ctx, "security.events", func(topic string, message []byte) error {
			logger.Warnf("收到安全事件: %s", string(message))
			return nil
		})
	}

	// 設定 MQTT 訂閱
	if mqttBroker != nil && mqttBroker.IsConnected() {
		// 訂閱設備認證請求
		mqttBroker.Subscribe("device/+/auth/request", func(topic string, payload []byte) error {
			logger.Infof("收到 MQTT 認證請求 [%s]: %s", topic, string(payload))

			// TODO: 處理 MQTT 認證請求
			// 1. 解析請求
			// 2. 驗證設備
			// 3. 回應結果

			return nil
		})

		// 訂閱設備狀態更新
		mqttBroker.Subscribe("device/+/status", func(topic string, payload []byte) error {
			logger.Debugf("收到設備狀態 [%s]: %s", topic, string(payload))
			return nil
		})
	}

	return baseHandler, extendedHandler
}

// GetRateLimitMiddleware 取得 Rate Limit 中間件（用於路由註冊）
func (h *AuthHandlerExtended) GetRateLimitMiddleware() gin.HandlerFunc {
	if h.rateLimiter != nil {
		return h.rateLimiter.Handler()
	}
	// 返回空中間件
	return func(c *gin.Context) {
		c.Next()
	}
}

// PublishAuthEvent 發布認證事件到 Pub/Sub 和 MQTT
// 注意：這是一個示例方法，實際使用需要透過 AuthHandlerExtended
func (h *AuthHandler) PublishAuthEvent(eventType, pcIdentifier, status string, data map[string]interface{}) {
	// 將 status 放入 data 中
	if data == nil {
		data = make(map[string]interface{})
	}
	data["status"] = status
	data["pc_identifier"] = pcIdentifier

	event := pubsub.Event{
		Type:      eventType,
		Source:    pcIdentifier,
		Timestamp: time.Now(),
		Data:      data,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		h.logger.Errorf("序列化事件失敗: %v", err)
		return
	}

	h.logger.Debugf("發布認證事件: %s - %s - %s", eventType, pcIdentifier, status)

	// 實際發布邏輯需要透過 AuthHandlerExtended 訪問 pubsub 和 mqttBroker 實例
	_ = eventJSON
}

// PublishAuthEventExtended 透過擴展處理器發布認證事件
func (h *AuthHandlerExtended) PublishAuthEventExtended(eventType, pcIdentifier, status string, data map[string]interface{}) error {
	// 將 status 放入 data 中
	if data == nil {
		data = make(map[string]interface{})
	}
	data["status"] = status
	data["pc_identifier"] = pcIdentifier

	event := pubsub.Event{
		Type:      eventType,
		Source:    pcIdentifier,
		Timestamp: time.Now(),
		Data:      data,
	}

	ctx := context.Background()

	// 發布到 Pub/Sub（如果可用）
	if h.pubsub != nil {
		// 將 Event 轉換為 Message
		message := &pubsub.Message{
			ID:         event.ID,
			RoutingKey: "auth.events",
			Body:       []byte(fmt.Sprintf(`{"type":"%s","source":"%s","data":%v}`, event.Type, event.Source, event.Data)),
			Timestamp:  event.Timestamp,
		}
		if err := h.pubsub.Publish(ctx, "auth.events", message); err != nil {
			h.logger.Errorf("發布到 Pub/Sub 失敗: %v", err)
			return fmt.Errorf("發布到 Pub/Sub 失敗: %w", err)
		}
		h.logger.Debugf("已發布認證事件到 Pub/Sub: %s - %s - %s", eventType, pcIdentifier, status)
	}

	// 發布到 MQTT（如果可用）
	if h.mqttBroker != nil && h.mqttBroker.IsConnected() {
		topic := fmt.Sprintf("auth/%s/%s", pcIdentifier, eventType)
		eventJSON, err := json.Marshal(event)
		if err != nil {
			h.logger.Errorf("序列化事件失敗: %v", err)
			return fmt.Errorf("序列化事件失敗: %w", err)
		}

		// QoS 1 (至少一次傳遞), retained = false (不保留訊息)
		if err := h.mqttBroker.Publish(topic, eventJSON, 1, false); err != nil {
			h.logger.Errorf("發布到 MQTT 失敗: %v", err)
			return fmt.Errorf("發布到 MQTT 失敗: %w", err)
		}
		h.logger.Debugf("已發布到 MQTT 主題: %s", topic)
	}

	return nil
}
