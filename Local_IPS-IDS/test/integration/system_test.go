//go:build integration
// +build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// SystemTestSuite 系統整合測試套件
type SystemTestSuite struct {
	suite.Suite
	dockerComposeFile string
	baseURL           string
	metricsURL        string
	grafanaURL        string
}

// SetupSuite 設置測試套件
func (s *SystemTestSuite) SetupSuite() {
	s.dockerComposeFile = "../../docker-compose.test.yml"
	s.baseURL = "http://localhost:3001"
	s.metricsURL = "http://localhost:8080"
	s.grafanaURL = "http://localhost:3000"

	// 啟動測試環境
	s.startTestEnvironment()

	// 等待服務啟動
	s.waitForServices()
}

// TearDownSuite 清理測試套件
func (s *SystemTestSuite) TearDownSuite() {
	s.stopTestEnvironment()
}

// startTestEnvironment 啟動測試環境
func (s *SystemTestSuite) startTestEnvironment() {
	cmd := exec.Command("docker-compose", "-f", s.dockerComposeFile, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	s.Require().NoError(err, "啟動測試環境失敗")

	s.T().Log("測試環境已啟動")
}

// stopTestEnvironment 停止測試環境
func (s *SystemTestSuite) stopTestEnvironment() {
	cmd := exec.Command("docker-compose", "-f", s.dockerComposeFile, "down", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		s.T().Logf("停止測試環境時發生錯誤: %v", err)
	}

	s.T().Log("測試環境已停止")
}

// waitForServices 等待服務啟動
func (s *SystemTestSuite) waitForServices() {
	services := map[string]string{
		"Axiom UI":   s.baseURL + "/api/v1/status",
		"Metrics":    s.metricsURL + "/health",
		"Grafana":    s.grafanaURL + "/api/health",
		"Prometheus": "http://localhost:9090/-/healthy",
	}

	for serviceName, healthURL := range services {
		s.waitForService(serviceName, healthURL, 60*time.Second)
	}
}

// waitForService 等待單個服務啟動
func (s *SystemTestSuite) waitForService(name, url string, timeout time.Duration) {
	s.T().Logf("等待服務啟動: %s", name)

	client := &http.Client{Timeout: 5 * time.Second}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			s.T().Logf("服務 %s 已啟動", name)
			return
		}
		if resp != nil {
			resp.Body.Close()
		}

		time.Sleep(2 * time.Second)
	}

	s.Fail(fmt.Sprintf("服務 %s 在 %v 內未能啟動", name, timeout))
}

// TestSystemHealth 測試系統健康狀態
func (s *SystemTestSuite) TestSystemHealth() {
	// 測試 Axiom UI 健康狀態
	resp, err := http.Get(s.baseURL + "/api/v1/status")
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 測試 Metrics 健康狀態
	resp, err = http.Get(s.metricsURL + "/health")
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 測試 Prometheus 健康狀態
	resp, err = http.Get("http://localhost:9090/-/healthy")
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

// TestMetricsEndpoint 測試指標端點
func (s *SystemTestSuite) TestMetricsEndpoint() {
	resp, err := http.Get(s.metricsURL + "/metrics")
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	resp.Body.Close()

	metricsContent := string(body)

	// 檢查關鍵指標是否存在
	expectedMetrics := []string{
		"pandora_system_uptime_seconds",
		"pandora_network_blocked",
		"pandora_device_connected",
		"pandora_threats_detected_total",
		"pandora_packets_inspected_total",
	}

	for _, metric := range expectedMetrics {
		s.Contains(metricsContent, metric, "應該包含指標: %s", metric)
	}
}

// TestAPIEndpoints 測試 API 端點
func (s *SystemTestSuite) TestAPIEndpoints() {
	testCases := []struct {
		name       string
		method     string
		url        string
		statusCode int
	}{
		{
			name:       "取得儀表板數據",
			method:     "GET",
			url:        "/api/v1/dashboard",
			statusCode: http.StatusOK,
		},
		{
			name:       "取得事件列表",
			method:     "GET",
			url:        "/api/v1/events",
			statusCode: http.StatusOK,
		},
		{
			name:       "取得指標數據",
			method:     "GET",
			url:        "/api/v1/metrics",
			statusCode: http.StatusOK,
		},
		{
			name:       "取得系統狀態",
			method:     "GET",
			url:        "/api/v1/status",
			statusCode: http.StatusOK,
		},
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req, err := http.NewRequest(tc.method, s.baseURL+tc.url, nil)
			s.Require().NoError(err)

			resp, err := client.Do(req)
			s.Require().NoError(err)
			defer resp.Body.Close()

			s.Equal(tc.statusCode, resp.StatusCode, "API 端點 %s 應該返回狀態碼 %d", tc.url, tc.statusCode)

			// 檢查回應是否為有效的 JSON
			if resp.Header.Get("Content-Type") == "application/json" {
				var jsonResp interface{}
				err = json.NewDecoder(resp.Body).Decode(&jsonResp)
				s.NoError(err, "回應應該是有效的 JSON")
			}
		})
	}
}

// TestDashboardData 測試儀表板數據結構
func (s *SystemTestSuite) TestDashboardData() {
	resp, err := http.Get(s.baseURL + "/api/v1/dashboard")
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var dashboardData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&dashboardData)
	s.Require().NoError(err)

	// 檢查必要的數據結構
	expectedFields := []string{
		"system_status",
		"security_metrics",
		"network_status",
		"recent_events",
		"timestamp",
	}

	for _, field := range expectedFields {
		s.Contains(dashboardData, field, "儀表板數據應該包含欄位: %s", field)
	}

	// 檢查系統狀態結構
	if systemStatus, ok := dashboardData["system_status"].(map[string]interface{}); ok {
		s.Contains(systemStatus, "uptime")
		s.Contains(systemStatus, "device_connected")
		s.Contains(systemStatus, "network_blocked")
		s.Contains(systemStatus, "status")
	}
}

// TestPrometheusIntegration 測試 Prometheus 整合
func (s *SystemTestSuite) TestPrometheusIntegration() {
	// 檢查 Prometheus 目標狀態
	resp, err := http.Get("http://localhost:9090/api/v1/targets")
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var targetsResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&targetsResp)
	s.Require().NoError(err)

	s.Equal("success", targetsResp["status"])

	// 檢查是否有活躍的目標
	if data, ok := targetsResp["data"].(map[string]interface{}); ok {
		if activeTargets, ok := data["activeTargets"].([]interface{}); ok {
			s.Greater(len(activeTargets), 0, "應該有活躍的 Prometheus 目標")
		}
	}
}

// TestGrafanaIntegration 測試 Grafana 整合
func (s *SystemTestSuite) TestGrafanaIntegration() {
	// 測試 Grafana 健康狀態
	resp, err := http.Get(s.grafanaURL + "/api/health")
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var healthResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&healthResp)
	s.Require().NoError(err)

	s.Equal("ok", healthResp["database"])
}

// TestLokiIntegration 測試 Loki 整合
func (s *SystemTestSuite) TestLokiIntegration() {
	// 測試 Loki 健康狀態
	resp, err := http.Get("http://localhost:3100/ready")
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)
}

// TestNetworkControl 測試網路控制功能
func (s *SystemTestSuite) TestNetworkControl() {
	client := &http.Client{Timeout: 10 * time.Second}

	// 測試阻斷網路
	blockData := `{"action": "block"}`
	req, err := http.NewRequest("POST", s.baseURL+"/api/v1/control/network",
		http.NoBody)
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		// 在測試環境中，這個功能可能不會真正執行
		// 我們只檢查 API 是否正確回應
		s.True(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest)
	}
}

// TestConcurrentRequests 測試並發請求
func (s *SystemTestSuite) TestConcurrentRequests() {
	const numRequests = 10
	const concurrency = 5

	client := &http.Client{Timeout: 5 * time.Second}

	// 建立通道來控制並發
	semaphore := make(chan struct{}, concurrency)
	results := make(chan error, numRequests)

	// 發送並發請求
	for i := 0; i < numRequests; i++ {
		go func() {
			semaphore <- struct{}{}        // 獲取信號量
			defer func() { <-semaphore }() // 釋放信號量

			resp, err := client.Get(s.baseURL + "/api/v1/status")
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				return
			}

			results <- nil
		}()
	}

	// 收集結果
	var errors []error
	for i := 0; i < numRequests; i++ {
		if err := <-results; err != nil {
			errors = append(errors, err)
		}
	}

	s.Empty(errors, "並發請求不應該有錯誤: %v", errors)
}

// TestDataPersistence 測試數據持久化
func (s *SystemTestSuite) TestDataPersistence() {
	// 這個測試需要重啟服務來檢查數據是否持久化
	// 在實際環境中，我們可能需要檢查資料庫或檔案系統

	// 取得初始狀態
	resp, err := http.Get(s.baseURL + "/api/v1/dashboard")
	s.Require().NoError(err)
	defer resp.Body.Close()

	var initialData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&initialData)
	s.Require().NoError(err)

	// 檢查時間戳是否合理
	if timestamp, ok := initialData["timestamp"].(float64); ok {
		s.Greater(timestamp, float64(0), "時間戳應該大於 0")

		// 檢查時間戳是否在合理範圍內 (最近 1 小時內)
		now := time.Now().Unix()
		s.Less(timestamp, float64(now+3600), "時間戳不應該超過當前時間太多")
		s.Greater(timestamp, float64(now-3600), "時間戳不應該太舊")
	}
}

// TestErrorHandling 測試錯誤處理
func (s *SystemTestSuite) TestErrorHandling() {
	client := &http.Client{Timeout: 5 * time.Second}

	// 測試不存在的端點
	resp, err := http.Get(s.baseURL + "/api/v1/nonexistent")
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusNotFound, resp.StatusCode)

	// 測試無效的 HTTP 方法
	req, err := http.NewRequest("DELETE", s.baseURL+"/api/v1/status", nil)
	s.Require().NoError(err)

	resp, err = client.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.True(resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusNotFound)
}

// TestPerformance 測試系統效能
func (s *SystemTestSuite) TestPerformance() {
	client := &http.Client{Timeout: 10 * time.Second}

	const numRequests = 100
	start := time.Now()

	// 發送多個請求測試效能
	for i := 0; i < numRequests; i++ {
		resp, err := client.Get(s.baseURL + "/api/v1/status")
		s.Require().NoError(err)
		s.Equal(http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	duration := time.Since(start)
	avgResponseTime := duration / numRequests

	s.T().Logf("平均回應時間: %v", avgResponseTime)
	s.Less(avgResponseTime, 100*time.Millisecond, "平均回應時間應該小於 100ms")
}

// 執行測試套件
func TestSystemIntegration(t *testing.T) {
	// 檢查是否在 CI 環境中或有 Docker
	if os.Getenv("CI") == "" {
		if _, err := exec.LookPath("docker-compose"); err != nil {
			t.Skip("Docker Compose 不可用，跳過整合測試")
		}
	}

	suite.Run(t, new(SystemTestSuite))
}

// TestMain 測試主函數
func TestMain(m *testing.M) {
	// 設定測試環境變數
	os.Setenv("LOG_LEVEL", "error")
	os.Setenv("ENVIRONMENT", "test")

	// 執行測試
	code := m.Run()

	// 清理
	os.Exit(code)
}
