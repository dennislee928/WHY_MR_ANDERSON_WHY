package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pandora_box_console_ids_ips/internal/engine"
	"pandora_box_console_ids_ips/internal/pubsub"

	"github.com/sirupsen/logrus"
)

// 這是一個示例，展示如何在 Axiom Engine 中整合 RabbitMQ
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

	// 創建 Engine 實例（佔位符）
	eng := &engine.Engine{}

	// 創建事件訂閱器
	subscriber, err := engine.NewEventSubscriber(config, eng, logger)
	if err != nil {
		log.Fatalf("Failed to create event subscriber: %v", err)
	}
	defer subscriber.Close()

	logger.Info("Event subscriber initialized")

	// 創建 context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動訂閱器
	if err := subscriber.Start(ctx); err != nil {
		log.Fatalf("Failed to start subscriber: %v", err)
	}

	logger.Info("Engine is now listening for events...")
	logger.Info("Subscribed to:")
	logger.Info("  - threat_events (威脅事件)")
	logger.Info("  - network_events (網路事件)")
	logger.Info("  - system_events (系統事件)")
	logger.Info("  - device_events (設備事件)")

	// 等待中斷信號
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Press Ctrl+C to stop...")
	<-sigChan

	logger.Info("Shutting down...")
	cancel()

	// 等待清理完成
	time.Sleep(2 * time.Second)
	logger.Info("Engine stopped")
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

