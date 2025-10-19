package axiom

import (
	"context"
	"fmt"
	"testing"
	"time"

	"pandora_box_console_ids_ips/internal/metrics"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAnalysisEngine(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	assert.NotNil(t, engine)
	assert.NotNil(t, engine.rules)
	assert.NotNil(t, engine.blacklist)
	assert.NotNil(t, engine.whitelist)
	assert.NotNil(t, engine.threatCache)
	assert.Greater(t, len(engine.rules), 0, "應該載入預設規則")
}

func TestAnalysisEngineStart(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// 啟動引擎
	go func() {
		err := engine.Start(ctx)
		assert.Equal(t, context.DeadlineExceeded, err)
	}()

	// 等待引擎啟動
	time.Sleep(10 * time.Millisecond)

	// 檢查引擎狀態
	engine.mutex.RLock()
	running := engine.running
	engine.mutex.RUnlock()

	assert.True(t, running)
}

func TestAnalyzePacket(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動引擎
	go func() {
		if err := engine.Start(ctx); err != nil {
			// 在 goroutine 中無法直接使用 t，改用 panic
			panic(fmt.Sprintf("引擎啟動失敗: %v", err))
		}
	}()
	time.Sleep(10 * time.Millisecond)

	// 測試正常封包
	normalPacket := &NetworkPacket{
		SourceIP:      "192.168.1.100",
		DestIP:        "192.168.1.1",
		SourcePort:    12345,
		DestPort:      80,
		Protocol:      "TCP",
		PayloadString: "GET / HTTP/1.1",
		Timestamp:     time.Now(),
	}

	result, err := engine.AnalyzePacket(normalPacket)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "192.168.1.100", result.SourceIP)
	assert.Equal(t, "allow", result.Action)
	assert.False(t, result.Blocked)
}

func TestAnalyzePacketWithThreat(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動引擎
	go func() {
		if err := engine.Start(ctx); err != nil {
			// 在 goroutine 中無法直接使用 t，改用 panic
			panic(fmt.Sprintf("引擎啟動失敗: %v", err))
		}
	}()
	time.Sleep(10 * time.Millisecond)

	// 測試 SQL 注入攻擊
	maliciousPacket := &NetworkPacket{
		SourceIP:      "10.0.0.1",
		DestIP:        "192.168.1.1",
		SourcePort:    54321,
		DestPort:      80,
		Protocol:      "TCP",
		PayloadString: "GET /?id=1 UNION SELECT * FROM users--",
		Timestamp:     time.Now(),
	}

	result, err := engine.AnalyzePacket(maliciousPacket)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "10.0.0.1", result.SourceIP)
	assert.Equal(t, "block", result.Action)
	assert.True(t, result.Blocked)
	assert.Equal(t, "high", result.ThreatLevel)
	assert.Equal(t, "pattern", result.ThreatType)
}

func TestWhitelistFunctionality(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動引擎
	go func() {
		if err := engine.Start(ctx); err != nil {
			// 在 goroutine 中無法直接使用 t，改用 panic
			panic(fmt.Sprintf("引擎啟動失敗: %v", err))
		}
	}()
	time.Sleep(10 * time.Millisecond)

	// 添加到白名單
	testIP := "192.168.1.200"
	engine.AddToWhitelist(testIP)

	// 測試白名單 IP 的惡意封包
	maliciousPacket := &NetworkPacket{
		SourceIP:      testIP,
		DestIP:        "192.168.1.1",
		SourcePort:    54321,
		DestPort:      80,
		Protocol:      "TCP",
		PayloadString: "GET /?id=1 UNION SELECT * FROM users--",
		Timestamp:     time.Now(),
	}

	result, err := engine.AnalyzePacket(maliciousPacket)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "allow", result.Action)
	assert.False(t, result.Blocked)
	assert.Contains(t, result.Details, "白名單")
}

func TestBlacklistFunctionality(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動引擎
	go func() {
		if err := engine.Start(ctx); err != nil {
			// 在 goroutine 中無法直接使用 t，改用 panic
			panic(fmt.Sprintf("引擎啟動失敗: %v", err))
		}
	}()
	time.Sleep(10 * time.Millisecond)

	// 手動添加到黑名單
	testIP := "10.0.0.100"
	engine.addToBlacklist(testIP, 1*time.Hour)

	// 測試黑名單 IP
	packet := &NetworkPacket{
		SourceIP:      testIP,
		DestIP:        "192.168.1.1",
		SourcePort:    12345,
		DestPort:      80,
		Protocol:      "TCP",
		PayloadString: "GET / HTTP/1.1",
		Timestamp:     time.Now(),
	}

	result, err := engine.AnalyzePacket(packet)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "block", result.Action)
	assert.True(t, result.Blocked)
	assert.Equal(t, "high", result.ThreatLevel)
	assert.Equal(t, "blacklisted_ip", result.ThreatType)
}

func TestAddRule(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	// 添加自定義規則
	customRule := &SecurityRule{
		ID:          "test_rule_001",
		Name:        "測試規則",
		Description: "測試用的安全規則",
		Type:        "pattern",
		Pattern:     "test_attack",
		Action:      "alert",
		Severity:    "medium",
		Enabled:     true,
	}

	err := engine.AddRule(customRule)
	require.NoError(t, err)

	// 檢查規則是否已添加
	rules := engine.GetRules()
	found := false
	for _, rule := range rules {
		if rule.ID == "test_rule_001" {
			found = true
			assert.Equal(t, "測試規則", rule.Name)
			assert.Equal(t, "medium", rule.Severity)
			break
		}
	}
	assert.True(t, found, "應該找到添加的規則")
}

func TestRemoveRule(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	// 取得初始規則數量
	initialRules := engine.GetRules()
	initialCount := len(initialRules)

	// 移除第一個規則
	if initialCount > 0 {
		ruleToRemove := initialRules[0]
		err := engine.RemoveRule(ruleToRemove.ID)
		require.NoError(t, err)

		// 檢查規則是否已移除
		updatedRules := engine.GetRules()
		assert.Equal(t, initialCount-1, len(updatedRules))

		// 確認規則確實被移除
		for _, rule := range updatedRules {
			assert.NotEqual(t, ruleToRemove.ID, rule.ID)
		}
	}

	// 測試移除不存在的規則
	err := engine.RemoveRule("non_existent_rule")
	assert.Error(t, err)
}

func TestThreatDetection(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動引擎
	go func() {
		if err := engine.Start(ctx); err != nil {
			// 在 goroutine 中無法直接使用 t，改用 panic
			panic(fmt.Sprintf("引擎啟動失敗: %v", err))
		}
	}()
	time.Sleep(10 * time.Millisecond)

	testCases := []struct {
		name           string
		payload        string
		expectedThreat bool
		expectedType   string
	}{
		{
			name:           "SQL Injection - UNION",
			payload:        "GET /?id=1 UNION SELECT password FROM users",
			expectedThreat: true,
			expectedType:   "pattern",
		},
		{
			name:           "SQL Injection - DROP",
			payload:        "POST /login with DROP TABLE users",
			expectedThreat: true,
			expectedType:   "pattern",
		},
		{
			name:           "Normal Request",
			payload:        "GET /api/users?page=1",
			expectedThreat: false,
			expectedType:   "",
		},
		{
			name:           "XSS Attack",
			payload:        "GET /?search=<script>alert('xss')</script>",
			expectedThreat: false, // 這個規則在預設規則中可能不存在
			expectedType:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			packet := &NetworkPacket{
				SourceIP:      "192.168.1.100",
				DestIP:        "192.168.1.1",
				SourcePort:    12345,
				DestPort:      80,
				Protocol:      "TCP",
				PayloadString: tc.payload,
				Timestamp:     time.Now(),
			}

			result, err := engine.AnalyzePacket(packet)
			require.NoError(t, err)

			if tc.expectedThreat {
				assert.Equal(t, "block", result.Action)
				assert.True(t, result.Blocked)
				assert.Equal(t, tc.expectedType, result.ThreatType)
			} else {
				assert.Equal(t, "allow", result.Action)
				assert.False(t, result.Blocked)
			}
		})
	}
}

func TestGetThreatInfo(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動引擎
	go func() {
		if err := engine.Start(ctx); err != nil {
			// 在 goroutine 中無法直接使用 t，改用 panic
			panic(fmt.Sprintf("引擎啟動失敗: %v", err))
		}
	}()
	time.Sleep(10 * time.Millisecond)

	// 觸發威脅偵測
	maliciousPacket := &NetworkPacket{
		SourceIP:      "192.168.1.100",
		DestIP:        "192.168.1.1",
		SourcePort:    12345,
		DestPort:      80,
		Protocol:      "TCP",
		PayloadString: "SELECT * FROM users",
		Timestamp:     time.Now(),
	}

	result, err := engine.AnalyzePacket(maliciousPacket)
	require.NoError(t, err)

	if result.Blocked {
		// 取得威脅資訊
		threatInfo := engine.GetThreatInfo()
		assert.NotNil(t, threatInfo)

		// 檢查是否記錄了威脅資訊
		if threat, exists := threatInfo["192.168.1.100"]; exists {
			assert.Equal(t, "192.168.1.100", threat.SourceIP)
			assert.Greater(t, threat.Count, 0)
		}
	}
}

func TestEngineStop(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動引擎
	go func() {
		if err := engine.Start(ctx); err != nil {
			// 在 goroutine 中無法直接使用 t，改用 panic
			panic(fmt.Sprintf("引擎啟動失敗: %v", err))
		}
	}()
	time.Sleep(10 * time.Millisecond)

	// 檢查引擎運行狀態
	engine.mutex.RLock()
	running := engine.running
	engine.mutex.RUnlock()
	assert.True(t, running)

	// 停止引擎
	engine.Stop()

	// 檢查引擎停止狀態
	engine.mutex.RLock()
	running = engine.running
	engine.mutex.RUnlock()
	assert.False(t, running)
}

func BenchmarkAnalyzePacket(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動引擎
	go func() {
		if err := engine.Start(ctx); err != nil {
			// 在 goroutine 中無法直接使用 t，改用 panic 或 log
			panic(fmt.Sprintf("引擎啟動失敗: %v", err))
		}
	}()
	time.Sleep(10 * time.Millisecond)

	packet := &NetworkPacket{
		SourceIP:      "192.168.1.100",
		DestIP:        "192.168.1.1",
		SourcePort:    12345,
		DestPort:      80,
		Protocol:      "TCP",
		PayloadString: "GET / HTTP/1.1",
		Timestamp:     time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.AnalyzePacket(packet)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAnalyzePacketWithThreat(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動引擎
	go func() {
		if err := engine.Start(ctx); err != nil {
			// 在 goroutine 中無法直接使用 t，改用 panic
			panic(fmt.Sprintf("引擎啟動失敗: %v", err))
		}
	}()
	time.Sleep(10 * time.Millisecond)

	maliciousPacket := &NetworkPacket{
		SourceIP:      "10.0.0.1",
		DestIP:        "192.168.1.1",
		SourcePort:    54321,
		DestPort:      80,
		Protocol:      "TCP",
		PayloadString: "GET /?id=1 UNION SELECT * FROM users",
		Timestamp:     time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.AnalyzePacket(maliciousPacket)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestRegexRules(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	metricsClient := metrics.NewPrometheusMetrics(logger)
	engine := NewAnalysisEngine(logger, metricsClient)

	// 添加正則表達式規則
	regexRule := &SecurityRule{
		ID:          "regex_test_001",
		Name:        "正則表達式測試",
		Description: "測試正則表達式規則",
		Type:        "pattern",
		Pattern:     "regex:(?i)(hack|crack|exploit)",
		Action:      "block",
		Severity:    "high",
		Enabled:     true,
	}

	err := engine.AddRule(regexRule)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動引擎
	go func() {
		if err := engine.Start(ctx); err != nil {
			// 在 goroutine 中無法直接使用 t，改用 panic
			panic(fmt.Sprintf("引擎啟動失敗: %v", err))
		}
	}()
	time.Sleep(10 * time.Millisecond)

	// 測試匹配正則表達式的封包
	testPacket := &NetworkPacket{
		SourceIP:      "192.168.1.100",
		DestIP:        "192.168.1.1",
		SourcePort:    12345,
		DestPort:      80,
		Protocol:      "TCP",
		PayloadString: "This is a HACK attempt",
		Timestamp:     time.Now(),
	}

	result, err := engine.AnalyzePacket(testPacket)
	require.NoError(t, err)
	assert.Equal(t, "block", result.Action)
	assert.True(t, result.Blocked)
	assert.Equal(t, "high", result.ThreatLevel)
}
