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

	"pandora_box_console_ids_ips/internal/services/device"
	"pandora_box_console_ids_ips/internal/pubsub"
	pb "pandora_box_console_ids_ips/api/proto/device"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	ServiceName    = "Device Service"
	ServiceVersion = "1.0.0"
	GRPCPort       = "50051"
	HTTPPort       = "8081"
)

func main() {
	// 初始化日誌
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	logger.Infof("Starting %s v%s", ServiceName, ServiceVersion)

	// 載入配置
	loadConfig(logger)

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

	logger.Info("Connected to RabbitMQ")

	// 創建 Device Service
	deviceConfig := &device.Config{
		DefaultPort:     getEnv("DEVICE_PORT", "/dev/ttyUSB0"),
		DefaultBaudRate: 115200,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    10 * time.Second,
	}

	deviceService := device.NewService(deviceConfig, mq, logger)

	// 創建 gRPC 服務器
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(logger)),
		grpc.StreamInterceptor(streamLoggingInterceptor(logger)),
	)

	// 註冊服務
	pb.RegisterDeviceServiceServer(grpcServer, deviceService)

	// 註冊健康檢查
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	// 啟用反射（用於 grpcurl 等工具）
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
	go startHTTPHealthServer(logger, deviceService)

	// 發布服務啟動事件
	ctx := context.Background()
	startEvent := pubsub.NewSystemEvent(ServiceName, "running", "Service started successfully")
	startEvent.Type = pubsub.EventTypeSystemStarted
	if message, err := pubsub.ToJSON(startEvent); err == nil {
		mq.Publish(ctx, "pandora.events", "system.started", message)
	}

	// 等待中斷信號
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Service is ready")
	sig := <-sigChan
	logger.Infof("Received signal %v, shutting down...", sig)

	// 優雅關閉
	healthServer.SetServingStatus(ServiceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	// 發布服務停止事件
	stopEvent := pubsub.NewSystemEvent(ServiceName, "stopped", "Service stopped gracefully")
	stopEvent.Type = pubsub.EventTypeSystemStopped
	if message, err := pubsub.ToJSON(stopEvent); err == nil {
		mq.Publish(ctx, "pandora.events", "system.stopped", message)
	}

	// 停止 gRPC 服務器
	grpcServer.GracefulStop()

	logger.Info("Service stopped")
}

// loadConfig loads configuration from file and environment
func loadConfig(logger *logrus.Logger) {
	viper.SetConfigName("device-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app/configs")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Warnf("Failed to read config file: %v, using defaults", err)
	} else {
		logger.Infof("Loaded config file: %s", viper.ConfigFileUsed())
	}
}

// startHTTPHealthServer starts the HTTP health check server
func startHTTPHealthServer(logger *logrus.Logger, service *device.Service) {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Check service health
		healthy := service.Health(ctx) == nil

		status := "healthy"
		statusCode := http.StatusOK
		if !healthy {
			status = "unhealthy"
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, `{"status":"%s","service":"%s","version":"%s","timestamp":"%s"}`,
			status, ServiceName, ServiceVersion, time.Now().Format(time.RFC3339))
	})

	// Readiness check
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"ready":true}`)
	})

	// Liveness check
	mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"alive":true}`)
	})

	server := &http.Server{
		Addr:         ":" + HTTPPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Infof("HTTP health server listening on port %s", HTTPPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("HTTP server error: %v", err)
	}
}

// loggingInterceptor logs gRPC requests
func loggingInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		if err != nil {
			logger.Errorf("RPC %s failed: %v (duration: %v)", info.FullMethod, err, duration)
		} else {
			logger.Debugf("RPC %s completed (duration: %v)", info.FullMethod, duration)
		}

		return resp, err
	}
}

// streamLoggingInterceptor logs gRPC streaming requests
func streamLoggingInterceptor(logger *logrus.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		err := handler(srv, ss)

		duration := time.Since(start)
		if err != nil {
			logger.Errorf("Stream %s failed: %v (duration: %v)", info.FullMethod, err, duration)
		} else {
			logger.Debugf("Stream %s completed (duration: %v)", info.FullMethod, duration)
		}

		return err
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

