package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"axiom-backend/internal/database"
)

func main() {
	// 初始化 Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	logger.Info("Starting Axiom Backend V3...")

	// 讀取配置
	cfg := loadConfig()

	// 初始化資料庫
	db, err := database.NewDatabase(&database.Config{
		PGHost:     cfg.PostgresHost,
		PGPort:     cfg.PostgresPort,
		PGUser:     cfg.PostgresUser,
		PGPassword: cfg.PostgresPassword,
		PGDatabase: cfg.PostgresDB,
		PGSSL:      "disable",
		RedisHost:  cfg.RedisHost,
		RedisPort:  cfg.RedisPort,
		RedisPassword: cfg.RedisPassword,
		RedisDB:    cfg.RedisDB,
	})
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 自動遷移
	if err := db.AutoMigrate(); err != nil {
		logger.Fatalf("Failed to auto migrate: %v", err)
	}
	logger.Info("Database migration completed")

	// 設置 Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(loggingMiddleware(logger))
	router.Use(corsMiddleware())

	// 健康檢查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "axiom-backend-v3",
			"version": "3.0.0",
			"time":    time.Now(),
		})
	})

	// 設置所有路由
	setupRoutes(router, db, cfg)

	// 啟動服務器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	// 優雅關閉
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	logger.Infof("Axiom Backend V3 started on port %d", cfg.Port)

	// 等待中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

// Config 配置結構
type Config struct {
	Port           int
	PostgresHost   string
	PostgresPort   int
	PostgresUser   string
	PostgresPassword string
	PostgresDB     string
	RedisHost      string
	RedisPort      int
	RedisPassword  string
	RedisDB        int
	PrometheusURL  string
	GrafanaURL     string
	LokiURL        string
	QuantumURL     string
	NginxURL       string
	NginxConfigPath string
}

// loadConfig 載入配置
func loadConfig() *Config {
	return &Config{
		Port:           getEnvInt("PORT", 3001),
		PostgresHost:   getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:   getEnvInt("POSTGRES_PORT", 5432),
		PostgresUser:   getEnv("POSTGRES_USER", "pandora"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "pandora123"),
		PostgresDB:     getEnv("POSTGRES_DB", "pandora_db"),
		RedisHost:      getEnv("REDIS_HOST", "localhost"),
		RedisPort:      getEnvInt("REDIS_PORT", 6379),
		RedisPassword:  getEnv("REDIS_PASSWORD", "pandora123"),
		RedisDB:        getEnvInt("REDIS_DB", 0),
		PrometheusURL:  getEnv("PROMETHEUS_URL", "http://localhost:9090"),
		GrafanaURL:     getEnv("GRAFANA_URL", "http://localhost:3000"),
		LokiURL:        getEnv("LOKI_URL", "http://localhost:3100"),
		QuantumURL:     getEnv("QUANTUM_URL", "http://localhost:8000"),
		NginxURL:       getEnv("NGINX_URL", "http://localhost:80"),
		NginxConfigPath: getEnv("NGINX_CONFIG_PATH", "/etc/nginx/nginx.conf"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		fmt.Sscanf(value, "%d", &intValue)
		return intValue
	}
	return defaultValue
}

// loggingMiddleware 日誌中間件
func loggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.WithFields(logrus.Fields{
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
		}).Info("HTTP Request")
	}
}

// corsMiddleware CORS 中間件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

