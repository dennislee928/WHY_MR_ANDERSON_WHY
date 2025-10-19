package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pandora_box_console_ids_ips/internal/agent"
	"pandora_box_console_ids_ips/internal/pubsub"

	"github.com/sirupsen/logrus"
)

// 這是一個示例，展示如何在 Pandora Agent 中整合 RabbitMQ
func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// 創建 RabbitMQ 配置
	config := &pubsub.Config{
		URL:                  getEnv("RABBITMQ_URL", "amqp://pandora:pandora123@localhost:5672/"),
		Exchange:             getEnv("RABBITMQ_EXCHANGE", "pandora.events"),
		ConnectionTimeout:    30 * time.Second,
		HeartbeatInterval:    60 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
	}

	// 創建事件發布器
	publisher, err := agent.NewEventPublisher(config, logger)
	if err != nil {
		log.Fatalf("Failed to create event publisher: %v", err)
	}
	defer publisher.Close()

	logger.Info("Event publisher initialized")

	// 創建 context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 發布 Agent 啟動事件
	if err := publisher.PublishAgentStarted(ctx); err != nil {
		logger.Errorf("Failed to publish agent started event: %v", err)
	}

	// 啟動定期健康檢查（每 30 秒）
	go publisher.StartPeriodicHealthCheck(ctx, 30*time.Second)

	// 模擬威脅檢測
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 模擬檢測到 DDoS 攻擊
				err := publisher.PublishThreatDetected(
					ctx,
					"ddos",
					"192.168.1.100",
					"Detected DDoS attack with 10000 requests/sec",
					8, // 威脅等級
				)
				if err != nil {
					logger.Errorf("Failed to publish threat: %v", err)
				} else {
					logger.Info("Published threat detection event")
				}
			}
		}
	}()

	// 模擬網路攻擊檢測
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 模擬檢測到端口掃描
				err := publisher.PublishNetworkAttack(
					ctx,
					"port_scan",
					"192.168.1.101",
					"10.0.0.1",
					"tcp",
				)
				if err != nil {
					logger.Errorf("Failed to publish network attack: %v", err)
				} else {
					logger.Info("Published network attack event")
				}
			}
		}
	}()

	// 模擬設備連接
	err = publisher.PublishDeviceConnected(ctx, "usb-001", "usb-serial", "/dev/ttyUSB0")
	if err != nil {
		logger.Errorf("Failed to publish device connected: %v", err)
	}

	// 等待中斷信號
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Agent running, press Ctrl+C to stop...")
	<-sigChan

	logger.Info("Shutting down...")

	// 發布 Agent 停止事件
	if err := publisher.PublishAgentStopped(ctx); err != nil {
		logger.Errorf("Failed to publish agent stopped event: %v", err)
	}

	logger.Info("Agent stopped")
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

