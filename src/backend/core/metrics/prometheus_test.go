package metrics

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPrometheusMetrics(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // 減少測試輸出

	metrics := NewPrometheusMetrics(logger)
	assert.NotNil(t, metrics)
	assert.NotNil(t, metrics.NetworkBlockedGauge)
	assert.NotNil(t, metrics.DeviceConnectionGauge)
	assert.NotNil(t, metrics.SecurityEventsCounter)
	assert.NotNil(t, metrics.registry)
}

func TestRecordNetworkBlock(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	metrics := NewPrometheusMetrics(logger)

	// 測試網路阻斷記錄
	metrics.RecordNetworkBlock(true)

	// 建立測試用的 HTTP 處理器
	handler := promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{})
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	response := w.Body.String()
	assert.Contains(t, response, "pandora_network_blocked 1")
	assert.Contains(t, response, "pandora_network_events_total")

	// 測試網路解除阻斷
	metrics.RecordNetworkBlock(false)
	req = httptest.NewRequest("GET", "/metrics", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	response = w.Body.String()
	assert.Contains(t, response, "pandora_network_blocked 0")
}

func TestRecordDeviceConnection(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	metrics := NewPrometheusMetrics(logger)

	// 測試裝置連接記錄
	metrics.RecordDeviceConnection(true)

	handler := promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{})
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	response := w.Body.String()
	assert.Contains(t, response, "pandora_device_connected 1")

	// 測試裝置斷線
	metrics.RecordDeviceConnection(false)
	req = httptest.NewRequest("GET", "/metrics", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	response = w.Body.String()
	assert.Contains(t, response, "pandora_device_connected 0")
}

func TestRecordSecurityEvent(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	metrics := NewPrometheusMetrics(logger)

	// 記錄安全事件 (簡化版)
	metrics.RecordSecurityEvent("brute_force", "detected")
	metrics.RecordSecurityEvent("ddos", "blocked")

	handler := promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{})
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	response := w.Body.String()
	assert.Contains(t, response, "pandora_security_events_total")
	assert.Contains(t, response, `event_type="brute_force"`)
	assert.Contains(t, response, `status="detected"`)
}

func TestRecordResponseTime(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	metrics := NewPrometheusMetrics(logger)

	// 記錄回應時間
	duration := 250 * time.Millisecond
	metrics.RecordResponseTime("api_call", "success", duration)

	handler := promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{})
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	response := w.Body.String()
	assert.Contains(t, response, "pandora_response_time_seconds")
	assert.Contains(t, response, `operation="api_call"`)
	assert.Contains(t, response, `status="success"`)
}

func TestUpdateSystemUptime(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	metrics := NewPrometheusMetrics(logger)

	// 模擬系統啟動時間
	startTime := time.Now().Add(-1 * time.Hour)
	metrics.UpdateSystemUptime(startTime)

	handler := promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{})
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	response := w.Body.String()
	assert.Contains(t, response, "pandora_system_uptime_seconds")
	// 檢查運行時間大約是 1 小時 (3600 秒左右)
	assert.Contains(t, response, "3")
}

func TestUpdateActiveSessions(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	metrics := NewPrometheusMetrics(logger)

	// 更新活躍會話數
	metrics.UpdateActiveSessions(5)

	handler := promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{})
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	response := w.Body.String()
	assert.Contains(t, response, "pandora_active_sessions 5")
}

func TestUpdateDataThroughput(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	metrics := NewPrometheusMetrics(logger)

	// 更新數據吞吐量
	metrics.UpdateDataThroughput("inbound", 1024.5)
	metrics.UpdateDataThroughput("outbound", 512.3)

	handler := promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{})
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	response := w.Body.String()
	assert.Contains(t, response, "pandora_data_throughput_bytes_per_second")
	assert.Contains(t, response, `direction="inbound"`)
	assert.Contains(t, response, `direction="outbound"`)
}

func TestMultipleMetricsRecording(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	metrics := NewPrometheusMetrics(logger)

	// 記錄多種指標
	metrics.RecordNetworkBlock(true)
	metrics.RecordDeviceConnection(true)
	metrics.RecordSecurityEvent("port_scan", "detected")
	metrics.RecordPinValidation()
	metrics.RecordTokenValidation()
	metrics.UpdateActiveSessions(10)

	handler := promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{})
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	response := w.Body.String()

	// 檢查所有指標都有出現
	assert.Contains(t, response, "pandora_network_blocked")
	assert.Contains(t, response, "pandora_device_connected")
	assert.Contains(t, response, "pandora_threats_detected_total")
	assert.Contains(t, response, "pandora_packets_inspected_total")
	assert.Contains(t, response, "pandora_pin_validations_total")
	assert.Contains(t, response, "pandora_token_validations_total")
	assert.Contains(t, response, "pandora_active_sessions")
	assert.Contains(t, response, "pandora_blocked_connections")
}

func TestMetricsServerStartup(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	metrics := NewPrometheusMetrics(logger)

	// 測試指標伺服器啟動 (使用隨機端口)
	go func() {
		err := metrics.StartMetricsServer("0") // 使用端口 0 讓系統自動分配
		if err != nil && !strings.Contains(err.Error(), "Server closed") {
			t.Errorf("指標伺服器啟動失敗: %v", err)
		}
	}()

	// 等待伺服器啟動
	time.Sleep(100 * time.Millisecond)
}

func BenchmarkRecordSecurityEvent(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	metrics := NewPrometheusMetrics(logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordSecurityEvent("benchmark_event", "detected")
	}
}

func BenchmarkRecordPinValidation(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	metrics := NewPrometheusMetrics(logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordPinValidation()
	}
}

func BenchmarkRecordResponseTime(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	metrics := NewPrometheusMetrics(logger)

	duration := 100 * time.Millisecond

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordResponseTime("benchmark_op", "success", duration)
	}
}

func TestPrometheusRegistryIntegrity(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	metrics := NewPrometheusMetrics(logger)

	// 檢查註冊表不為空
	gatherer := metrics.registry
	metricFamilies, err := gatherer.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metricFamilies), 0, "應該有註冊的指標")

	// 檢查特定指標是否存在
	var foundMetrics []string
	for _, mf := range metricFamilies {
		foundMetrics = append(foundMetrics, mf.GetName())
	}

	expectedMetrics := []string{
		"pandora_network_blocked",
		"pandora_device_connected",
		"pandora_threats_detected_total",
		"pandora_packets_inspected_total",
		"pandora_system_uptime_seconds",
	}

	for _, expected := range expectedMetrics {
		assert.Contains(t, foundMetrics, expected, "應該包含指標: %s", expected)
	}
}
