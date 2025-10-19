package metrics

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// PrometheusMetrics Prometheus指標收集器
type PrometheusMetrics struct {
	logger *logrus.Logger

	// 系統指標
	NetworkBlockedGauge    prometheus.Gauge
	DeviceConnectionGauge  prometheus.Gauge
	PinValidationCounter   prometheus.Counter
	TokenValidationCounter prometheus.Counter
	NetworkEventsCounter   *prometheus.CounterVec
	SystemUptimeGauge      prometheus.Gauge
	AuthAttempts           *prometheus.CounterVec

	// 基礎安全指標 (簡化版)
	SecurityEventsCounter *prometheus.CounterVec

	// 效能指標
	ResponseTimeHistogram *prometheus.HistogramVec
	ActiveSessionsGauge   prometheus.Gauge
	DataThroughputGauge   *prometheus.GaugeVec

	registry *prometheus.Registry
}

// NewPrometheusMetrics 建立新的Prometheus指標收集器
func NewPrometheusMetrics(logger *logrus.Logger) *PrometheusMetrics {
	registry := prometheus.NewRegistry()

	// 系統指標
	networkBlockedGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pandora_network_blocked",
		Help: "網路是否被阻斷 (1=blocked, 0=allowed)",
	})

	deviceConnectionGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pandora_device_connected",
		Help: "IoT裝置連接狀態 (1=connected, 0=disconnected)",
	})

	pinValidationCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pandora_pin_validations_total",
		Help: "PIN碼驗證總次數",
	})

	tokenValidationCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pandora_token_validations_total",
		Help: "Token驗證總次數",
	})

	networkEventsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pandora_network_events_total",
			Help: "網路事件總數",
		},
		[]string{"event_type", "status"},
	)

	systemUptimeGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pandora_system_uptime_seconds",
		Help: "系統運行時間（秒）",
	})

	authAttempts := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pandora_auth_attempts_total",
			Help: "認證嘗試總數",
		},
		[]string{"status"},
	)

	// 基礎安全指標 (簡化版)
	securityEventsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pandora_security_events_total",
			Help: "安全事件總數",
		},
		[]string{"event_type", "status"},
	)

	// 效能指標
	responseTimeHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "pandora_response_time_seconds",
			Help:    "回應時間分布",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "status"},
	)

	activeSessionsGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pandora_active_sessions",
		Help: "當前活躍會話數量",
	})

	dataThroughputGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pandora_data_throughput_bytes_per_second",
			Help: "數據吞吐量（bytes/秒）",
		},
		[]string{"direction"}, // "inbound", "outbound"
	)

	// 註冊所有指標 (簡化版)
	registry.MustRegister(
		networkBlockedGauge,
		deviceConnectionGauge,
		pinValidationCounter,
		tokenValidationCounter,
		networkEventsCounter,
		systemUptimeGauge,
		authAttempts,
		securityEventsCounter,
		responseTimeHistogram,
		activeSessionsGauge,
		dataThroughputGauge,
	)

	return &PrometheusMetrics{
		logger:                 logger,
		NetworkBlockedGauge:    networkBlockedGauge,
		DeviceConnectionGauge:  deviceConnectionGauge,
		PinValidationCounter:   pinValidationCounter,
		TokenValidationCounter: tokenValidationCounter,
		NetworkEventsCounter:   networkEventsCounter,
		SystemUptimeGauge:      systemUptimeGauge,
		AuthAttempts:           authAttempts,
		SecurityEventsCounter:  securityEventsCounter,
		ResponseTimeHistogram:  responseTimeHistogram,
		ActiveSessionsGauge:    activeSessionsGauge,
		DataThroughputGauge:    dataThroughputGauge,
		registry:               registry,
	}
}

// StartMetricsServer 啟動Prometheus指標伺服器
func (pm *PrometheusMetrics) StartMetricsServer(port string) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(pm.registry, promhttp.HandlerOpts{})))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})

	// System info endpoint
	router.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":          "Pandora Box Console IDS-IPS",
			"version":          "1.0.0",
			"metrics_endpoint": "/metrics",
		})
	})

	pm.logger.Infof("Prometheus指標伺服器啟動於端口: %s", port)
	return router.Run(":" + port)
}

// RecordNetworkBlock 記錄網路阻斷事件
func (pm *PrometheusMetrics) RecordNetworkBlock(blocked bool) {
	if blocked {
		pm.NetworkBlockedGauge.Set(1)
		pm.NetworkEventsCounter.WithLabelValues("block", "success").Inc()
	} else {
		pm.NetworkBlockedGauge.Set(0)
		pm.NetworkEventsCounter.WithLabelValues("unblock", "success").Inc()
	}
}

// RecordDeviceConnection 記錄裝置連接狀態
func (pm *PrometheusMetrics) RecordDeviceConnection(connected bool) {
	if connected {
		pm.DeviceConnectionGauge.Set(1)
	} else {
		pm.DeviceConnectionGauge.Set(0)
	}
}

// RecordPinValidation 記錄PIN碼驗證
func (pm *PrometheusMetrics) RecordPinValidation() {
	pm.PinValidationCounter.Inc()
}

// RecordTokenValidation 記錄Token驗證
func (pm *PrometheusMetrics) RecordTokenValidation() {
	pm.TokenValidationCounter.Inc()
}

// RecordSecurityEvent 記錄安全事件 (簡化版)
func (pm *PrometheusMetrics) RecordSecurityEvent(eventType, status string) {
	pm.SecurityEventsCounter.WithLabelValues(eventType, status).Inc()
}

// RecordResponseTime 記錄回應時間
func (pm *PrometheusMetrics) RecordResponseTime(operation, status string, duration time.Duration) {
	pm.ResponseTimeHistogram.WithLabelValues(operation, status).Observe(duration.Seconds())
}

// UpdateSystemUptime 更新系統運行時間
func (pm *PrometheusMetrics) UpdateSystemUptime(startTime time.Time) {
	uptime := time.Since(startTime).Seconds()
	pm.SystemUptimeGauge.Set(uptime)
}

// UpdateActiveSessions 更新活躍會話數
func (pm *PrometheusMetrics) UpdateActiveSessions(count int) {
	pm.ActiveSessionsGauge.Set(float64(count))
}

// UpdateDataThroughput 更新數據吞吐量
func (pm *PrometheusMetrics) UpdateDataThroughput(direction string, bytesPerSecond float64) {
	pm.DataThroughputGauge.WithLabelValues(direction).Set(bytesPerSecond)
}
