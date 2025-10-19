package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pandora_box_console_ids_ips/internal/services/network"
	"pandora_box_console_ids_ips/internal/pubsub"
	pb "pandora_box_console_ids_ips/api/proto/network"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	ServiceName    = "Network Service"
	ServiceVersion = "1.0.0"
	GRPCPort       = "50052"
	HTTPPort       = "8082"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	logger.Infof("Starting %s v%s", ServiceName, ServiceVersion)

	// 創建 RabbitMQ 連接
	mqConfig := &pubsub.Config{
		URL:                  getEnv("RABBITMQ_URL", "amqp://pandora:pandora123@localhost:5672/"),
		Exchange:             getEnv("RABBITMQ_EXCHANGE", "pandora.events"),
		ConnectionTimeout:    30 * time.Second,
		HeartbeatInterval:    60 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
	}

	mq, err := pubsub.NewRabbitMQ(mqConfig)
	if err != nil {
		logger.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer mq.Close()

	// 創建 Network Service
	networkConfig := &network.Config{
		DefaultInterface: getEnv("NETWORK_INTERFACE", "eth0"),
		SnapshotLength:   1600,
		Promiscuous:      true,
		Timeout:          30 * time.Second,
	}

	networkService := network.NewService(networkConfig, mq, logger)

	// 創建 gRPC 服務器
	grpcServer := grpc.NewServer()
	pb.RegisterNetworkServiceServer(grpcServer, networkService)

	// 註冊健康檢查
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	// 啟動 gRPC 服務器
	lis, err := net.Listen("tcp", ":"+GRPCPort)
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		logger.Infof("gRPC server listening on port %s", GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve: %v", err)
		}
	}()

	// 啟動 HTTP 健康檢查服務器
	go startHTTPHealthServer(logger, networkService)

	// 發布服務啟動事件
	ctx := context.Background()
	startEvent := pubsub.NewSystemEvent(ServiceName, "running", "Service started successfully")
	if message, err := pubsub.ToJSON(startEvent); err == nil {
		mq.Publish(ctx, "pandora.events", "system.started", message)
	}

	// 等待中斷信號
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Service is ready")
	<-sigChan

	logger.Info("Shutting down...")
	grpcServer.GracefulStop()
	logger.Info("Service stopped")
}

func startHTTPHealthServer(logger *logrus.Logger, service *network.Service) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		healthy := service.Health(ctx) == nil
		status := "healthy"
		statusCode := http.StatusOK
		if !healthy {
			status = "unhealthy"
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, `{"status":"%s","service":"%s","version":"%s"}`,
			status, ServiceName, ServiceVersion)
	})

	server := &http.Server{
		Addr:    ":" + HTTPPort,
		Handler: mux,
	}

	logger.Infof("HTTP health server listening on port %s", HTTPPort)
	server.ListenAndServe()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

