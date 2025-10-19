package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"pandora_box_console_ids_ips/internal/handlers"
	"pandora_box_console_ids_ips/internal/loadbalancer"
	"pandora_box_console_ids_ips/internal/logging"
	"pandora_box_console_ids_ips/internal/metrics"
	"pandora_box_console_ids_ips/internal/mqtt"
	"pandora_box_console_ids_ips/internal/pubsub"
	"pandora_box_console_ids_ips/internal/ratelimit"
)

const (
	AppName    = "Pandora Box Console"
	AppVersion = "1.0.0"
)

func main() {
	// 命令列參數
	var (
		configFile = flag.String("config", "configs/console-config.yaml", "設定檔路徑")
		logLevel   = flag.String("log-level", "info", "日誌等級")
		port       = flag.Int("port", 3001, "服務器埠號")
		version    = flag.Bool("version", false, "顯示版本資訊")
	)
	flag.Parse()

	if *version {
		fmt.Printf("%s v%s\n", AppName, AppVersion)
		os.Exit(0)
	}

	// 設定日誌
	logger := setupLogging(*logLevel)
	logger.Infof("啟動 %s v%s", AppName, AppVersion)

	// 載入配置
	if err := loadConfig(*configFile); err != nil {
		logger.Fatalf("載入配置失敗: %v", err)
	}

	// 創建中央日誌記錄器
	centralLogger := logging.NewCentralLogger()

	// 創建指標收集器
	metricsCollector := metrics.NewPrometheusMetrics(logger)

	// ========== 新增模組初始化 ==========

	// 1. 初始化 Rate Limiter（暴力攻擊防護）
	rateLimitConfig := &ratelimit.Config{
		Enabled:         viper.GetBool("ratelimit.enabled"),
		Rate:            viper.GetInt("ratelimit.rate"),
		Burst:           viper.GetInt("ratelimit.burst"),
		WindowSize:      viper.GetDuration("ratelimit.window_size"),
		MaxAttempts:     viper.GetInt("ratelimit.max_attempts"),
		LockoutTime:     viper.GetDuration("ratelimit.lockout_time"),
		BlockEnabled:    viper.GetBool("ratelimit.block_enabled"),
		BlockTime:       viper.GetDuration("ratelimit.block_time"),
		CleanupInterval: viper.GetDuration("ratelimit.cleanup_interval"),
	}
	rateLimiter := ratelimit.NewTokenBucketLimiter(rateLimitConfig, logger)
	defer rateLimiter.Stop()

	rateLimitMiddlewareConfig := &ratelimit.MiddlewareConfig{
		KeyStrategy:  viper.GetString("ratelimit.key_strategy"),
		StatusCode:   http.StatusTooManyRequests,
		ErrorMessage: "Too many requests, please try again later",
		RetryAfter:   true,
		WhitelistIPs: viper.GetStringSlice("ratelimit.whitelist_ips"),
	}
	rateLimitMiddleware := ratelimit.NewMiddleware(rateLimiter, rateLimitMiddlewareConfig, logger)

	// 2. 初始化 Pub/Sub 系統
	var pubsubInstance pubsub.PubSub
	if viper.GetBool("pubsub.enabled") {
		pubsubConfig := &pubsub.Config{
			Type:          viper.GetString("pubsub.type"),
			RedisAddr:     viper.GetString("pubsub.redis_addr"),
			RedisPassword: viper.GetString("pubsub.redis_password"),
			RedisDB:       viper.GetInt("pubsub.redis_db"),
			BufferSize:    viper.GetInt("pubsub.buffer_size"),
		}
		var err error
		pubsubInstance, err = pubsub.NewPubSub(pubsubConfig, logger)
		if err != nil {
			logger.Errorf("初始化 Pub/Sub 失敗: %v", err)
		} else {
			defer pubsubInstance.Close()
			logger.Info("Pub/Sub 系統已啟動")
		}
	}

	// 3. 初始化 MQTT Broker
	var mqttBroker *mqtt.Broker
	if viper.GetBool("mqtt.enabled") {
		mqttConfig := &mqtt.Config{
			Broker:           viper.GetString("mqtt.broker"),
			Port:             viper.GetInt("mqtt.port"),
			ClientID:         viper.GetString("mqtt.client_id"),
			Username:         viper.GetString("mqtt.username"),
			Password:         viper.GetString("mqtt.password"),
			TLSEnabled:       viper.GetBool("mqtt.tls_enabled"),
			DefaultQoS:       byte(viper.GetInt("mqtt.default_qos")),
			KeepAlive:        viper.GetInt("mqtt.keep_alive"),
			ConnectTimeout:   viper.GetDuration("mqtt.connect_timeout"),
			ReconnectDelay:   viper.GetDuration("mqtt.reconnect_delay"),
			MaxReconnectWait: viper.GetDuration("mqtt.max_reconnect_wait"),
			AutoReconnect:    viper.GetBool("mqtt.auto_reconnect"),
			CleanSession:     viper.GetBool("mqtt.clean_session"),
			OrderMatters:     viper.GetBool("mqtt.order_matters"),
		}
		var err error
		mqttBroker, err = mqtt.NewBroker(mqttConfig, logger)
		if err != nil {
			logger.Errorf("初始化 MQTT Broker 失敗: %v", err)
		} else {
			if err := mqttBroker.Start(); err != nil {
				logger.Errorf("啟動 MQTT Broker 失敗: %v", err)
			} else {
				defer mqttBroker.Stop()
				logger.Info("MQTT Broker 已啟動")
			}
		}
	}

	// 4. 初始化 Load Balancer（如果啟用）
	var lb *loadbalancer.LoadBalancer
	if viper.GetBool("loadbalancer.enabled") {
		lbConfig := &loadbalancer.Config{
			Backends:            viper.GetStringSlice("loadbalancer.backends"),
			Strategy:            viper.GetString("loadbalancer.strategy"),
			HealthCheckEnabled:  viper.GetBool("loadbalancer.health_check_enabled"),
			HealthCheckInterval: viper.GetDuration("loadbalancer.health_check_interval"),
			HealthCheckTimeout:  viper.GetDuration("loadbalancer.health_check_timeout"),
			HealthCheckPath:     viper.GetString("loadbalancer.health_check_path"),
			MaxRetries:          viper.GetInt("loadbalancer.max_retries"),
			RetryDelay:          viper.GetDuration("loadbalancer.retry_delay"),
		}
		var err error
		lb, err = loadbalancer.NewLoadBalancer(lbConfig, logger)
		if err != nil {
			logger.Errorf("初始化 Load Balancer 失敗: %v", err)
		} else {
			defer lb.Stop()
			logger.Info("Load Balancer 已啟動")
		}
	}

	// 創建認證處理器
	authHandler := handlers.NewAuthHandler(logger, centralLogger, metricsCollector)

	// 記錄已啟動的模組
	logger.Info("======= 模組啟動狀態 =======")
	if rateLimitMiddleware != nil {
		logger.Info("✓ Rate Limiter 已啟動")
	}
	if pubsubInstance != nil {
		logger.Info("✓ Pub/Sub 系統已啟動")
	}
	if mqttBroker != nil {
		logger.Info("✓ MQTT Broker 已啟動")
	}
	if lb != nil {
		logger.Info("✓ Load Balancer 已啟動")
	}
	logger.Info("===========================")

	// 創建HTTP服務器
	server := setupHTTPServer(*port, authHandler, rateLimitMiddleware, lb, logger)

	// 優雅啟動和關閉

	// 啟動服務器
	go func() {
		logger.Infof("Pandora Box Console 服務器啟動於 :%d", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("服務器啟動失敗: %v", err)
		}
	}()

	// 等待中斷信號
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("收到關閉信號，開始優雅關閉...")

	// 優雅關閉
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("服務器關閉失敗: %v", err)
	} else {
		logger.Info("服務器已優雅關閉")
	}
}

// setupLogging 設定日誌系統
func setupLogging(level string) *logrus.Logger {
	logger := logrus.New()

	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	return logger
}

// loadConfig 載入配置
func loadConfig(configFile string) error {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("PANDORA_CONSOLE")

	// 設定預設值
	viper.SetDefault("server.port", 3001)
	viper.SetDefault("server.timeout", "30s")
	viper.SetDefault("auth.default_psk", "pandora-default-key")
	viper.SetDefault("logging.level", "info")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("讀取配置檔案錯誤: %v", err)
		}
		// 配置檔案不存在，使用預設值
		logrus.Warn("配置檔案不存在，使用預設設定")
	}

	return nil
}

// setupHTTPServer 設定HTTP服務器
func setupHTTPServer(port int, authHandler *handlers.AuthHandler,
	rateLimitMW *ratelimit.Middleware, lb *loadbalancer.LoadBalancer, logger *logrus.Logger) *http.Server {
	// 設定Gin模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// 全域中間件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(securityHeadersMiddleware())

	// 全域速率限制（可選）
	if rateLimitMW != nil {
		router.Use(rateLimitMW.Handler())
	}

	// 根路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": AppName,
			"version": AppVersion,
			"status":  "running",
			"time":    time.Now().UTC(),
		})
	})

	// API路由
	v1 := router.Group("/api/v1")
	{
		// 健康檢查（無需速率限制）
		v1.GET("/health", authHandler.HealthCheck)

		// PC身份驗證（帶暴力攻擊防護）
		pcGroup := v1.Group("/verify")
		if rateLimitMW != nil {
			pcGroup.Use(rateLimitMW.BruteForceProtection())
		}
		{
			pcGroup.POST("/pc", authHandler.VerifyPC)
		}

		// TPM認證（帶暴力攻擊防護）
		auth := v1.Group("/auth")
		if rateLimitMW != nil {
			auth.Use(rateLimitMW.BruteForceProtection())
		}
		{
			auth.POST("/challenge", authHandler.GetChallenge)
			auth.POST("/verify-tpm", authHandler.VerifyTPMSignature)
		}

		// 日誌接收端點 (供Agent發送日誌)
		v1.POST("/logs", handleLogs)

		// 指標接收端點 (供Agent發送指標)
		v1.POST("/metrics", handleMetrics)

		// 系統管理端點
		admin := v1.Group("/admin")
		{
			// Rate Limit 狀態
			admin.GET("/ratelimit/status/:key", func(c *gin.Context) {
				key := c.Param("key")
				if rateLimitMW != nil {
					// TODO: 實作獲取狀態的方法
					c.JSON(http.StatusOK, gin.H{"key": key})
				} else {
					c.JSON(http.StatusNotImplemented, gin.H{"error": "Rate limit not enabled"})
				}
			})

			// Load Balancer 狀態
			admin.GET("/loadbalancer/stats", func(c *gin.Context) {
				if lb != nil {
					status := lb.GetStatus()
					c.JSON(http.StatusOK, gin.H{"loadbalancer": status})
				} else {
					c.JSON(http.StatusNotImplemented, gin.H{"error": "Load balancer not enabled"})
				}
			})
		}
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// corsMiddleware CORS中間件
func corsMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

// securityHeadersMiddleware 安全標頭中間件
func securityHeadersMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	})
}

// handleLogs 處理日誌接收
func handleLogs(c *gin.Context) {
	var logEntry interface{}

	if err := c.ShouldBindJSON(&logEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid log format"})
		return
	}

	// TODO: 處理接收到的日誌
	// 可以發送到Loki或其他日誌系統

	c.JSON(http.StatusAccepted, gin.H{"status": "log received"})
}

// handleMetrics 處理指標接收
func handleMetrics(c *gin.Context) {
	var metrics interface{}

	if err := c.ShouldBindJSON(&metrics); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metrics format"})
		return
	}

	// TODO: 處理接收到的指標
	// 可以發送到Prometheus或其他指標系統

	c.JSON(http.StatusAccepted, gin.H{"status": "metrics received"})
}
