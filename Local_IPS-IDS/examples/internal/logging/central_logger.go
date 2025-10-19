package logging

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// CentralLogger 中央日誌記錄器
type CentralLogger struct {
	logger *logrus.Logger

	// 日誌事件通道
	eventChan chan LogEvent
	stopChan  chan struct{}
	wg        sync.WaitGroup

	// 日誌統計
	stats     LogStats
	statMutex sync.RWMutex
}

// LogEvent 日誌事件結構
type LogEvent struct {
	Timestamp    time.Time              `json:"timestamp"`
	Level        string                 `json:"level"`
	Source       string                 `json:"source"`     // ESP32, Agent, Console
	EventType    string                 `json:"event_type"` // AUTH, SECURITY, NETWORK, SYSTEM
	Action       string                 `json:"action"`     // LOGIN, LOGOUT, OTP_VERIFY, etc.
	Status       string                 `json:"status"`     // SUCCESS, FAIL, ERROR
	PCIdentifier string                 `json:"pc_identifier,omitempty"`
	Message      string                 `json:"message"`
	Details      map[string]interface{} `json:"details,omitempty"`
	IPAddress    string                 `json:"ip_address,omitempty"`
	UserAgent    string                 `json:"user_agent,omitempty"`
}

// LogStats 日誌統計資訊
type LogStats struct {
	TotalEvents    int64     `json:"total_events"`
	AuthEvents     int64     `json:"auth_events"`
	SecurityEvents int64     `json:"security_events"`
	NetworkEvents  int64     `json:"network_events"`
	SystemEvents   int64     `json:"system_events"`
	ErrorEvents    int64     `json:"error_events"`
	LastEventTime  time.Time `json:"last_event_time"`
	EventsPerHour  int64     `json:"events_per_hour"`
	HourStartTime  time.Time `json:"hour_start_time"`
}

// NewCentralLogger 創建新的中央日誌記錄器
func NewCentralLogger() *CentralLogger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	cl := &CentralLogger{
		logger:    logger,
		eventChan: make(chan LogEvent, 1000), // 緩衝1000個事件
		stopChan:  make(chan struct{}),
		stats: LogStats{
			HourStartTime: time.Now(),
		},
	}

	// 啟動日誌處理協程
	cl.wg.Add(1)
	go cl.processEvents()

	return cl
}

// LogAuthEvent 記錄認證事件
func (cl *CentralLogger) LogAuthEvent(pcID, action, status, message string, details map[string]interface{}) {
	event := LogEvent{
		Timestamp:    time.Now(),
		Level:        cl.getLogLevel(status),
		Source:       "Console",
		EventType:    "AUTH",
		Action:       action,
		Status:       status,
		PCIdentifier: pcID,
		Message:      message,
		Details:      details,
	}

	cl.logEvent(event)
}

// LogSecurityEvent 記錄安全事件
func (cl *CentralLogger) LogSecurityEvent(pcID, action, status, message string, details map[string]interface{}) {
	event := LogEvent{
		Timestamp:    time.Now(),
		Level:        "WARN", // 安全事件通常為警告級別
		Source:       "Console",
		EventType:    "SECURITY",
		Action:       action,
		Status:       status,
		PCIdentifier: pcID,
		Message:      message,
		Details:      details,
	}

	cl.logEvent(event)
}

// LogNetworkEvent 記錄網路事件
func (cl *CentralLogger) LogNetworkEvent(pcID, action, status, message string, details map[string]interface{}) {
	event := LogEvent{
		Timestamp:    time.Now(),
		Level:        cl.getLogLevel(status),
		Source:       "Console",
		EventType:    "NETWORK",
		Action:       action,
		Status:       status,
		PCIdentifier: pcID,
		Message:      message,
		Details:      details,
	}

	cl.logEvent(event)
}

// LogSystemEvent 記錄系統事件
func (cl *CentralLogger) LogSystemEvent(source, action, status, message string, details map[string]interface{}) {
	event := LogEvent{
		Timestamp: time.Now(),
		Level:     cl.getLogLevel(status),
		Source:    source,
		EventType: "SYSTEM",
		Action:    action,
		Status:    status,
		Message:   message,
		Details:   details,
	}

	cl.logEvent(event)
}

// LogAgentEvent 記錄Agent事件 (從Agent接收)
func (cl *CentralLogger) LogAgentEvent(event LogEvent) {
	event.Timestamp = time.Now()
	event.Source = "Agent"
	cl.logEvent(event)
}

// LogESP32Event 記錄ESP32事件 (從Agent轉發)
func (cl *CentralLogger) LogESP32Event(event LogEvent) {
	event.Timestamp = time.Now()
	event.Source = "ESP32"
	cl.logEvent(event)
}

// logEvent 內部日誌事件處理
func (cl *CentralLogger) logEvent(event LogEvent) {
	select {
	case cl.eventChan <- event:
		// 事件已加入佇列
	default:
		// 佇列已滿，丟棄事件並記錄警告
		cl.logger.Warn("日誌事件佇列已滿，丟棄事件")
	}
}

// processEvents 處理日誌事件協程
func (cl *CentralLogger) processEvents() {
	defer cl.wg.Done()

	for {
		select {
		case event := <-cl.eventChan:
			cl.handleEvent(event)
		case <-cl.stopChan:
			// 處理剩餘事件
			for len(cl.eventChan) > 0 {
				event := <-cl.eventChan
				cl.handleEvent(event)
			}
			return
		}
	}
}

// handleEvent 處理單個日誌事件
func (cl *CentralLogger) handleEvent(event LogEvent) {
	// 更新統計
	cl.updateStats(event)

	// 輸出到logrus
	logEntry := cl.logger.WithFields(logrus.Fields{
		"source":        event.Source,
		"event_type":    event.EventType,
		"action":        event.Action,
		"status":        event.Status,
		"pc_identifier": event.PCIdentifier,
		"ip_address":    event.IPAddress,
		"user_agent":    event.UserAgent,
	})

	if event.Details != nil {
		logEntry = logEntry.WithField("details", event.Details)
	}

	switch event.Level {
	case "DEBUG":
		logEntry.Debug(event.Message)
	case "INFO":
		logEntry.Info(event.Message)
	case "WARN":
		logEntry.Warn(event.Message)
	case "ERROR":
		logEntry.Error(event.Message)
	default:
		logEntry.Info(event.Message)
	}

	// TODO: 發送到外部日誌系統 (Loki, ElasticSearch等)
	cl.sendToExternalSystems(event)
}

// updateStats 更新日誌統計
func (cl *CentralLogger) updateStats(event LogEvent) {
	cl.statMutex.Lock()
	defer cl.statMutex.Unlock()

	cl.stats.TotalEvents++
	cl.stats.LastEventTime = event.Timestamp

	// 更新每小時計數
	now := time.Now()
	if now.Sub(cl.stats.HourStartTime) > time.Hour {
		cl.stats.HourStartTime = now
		cl.stats.EventsPerHour = 1
	} else {
		cl.stats.EventsPerHour++
	}

	// 按事件類型計數
	switch event.EventType {
	case "AUTH":
		cl.stats.AuthEvents++
	case "SECURITY":
		cl.stats.SecurityEvents++
	case "NETWORK":
		cl.stats.NetworkEvents++
	case "SYSTEM":
		cl.stats.SystemEvents++
	}

	// 按狀態計數
	if event.Status == "error" || event.Status == "fail" {
		cl.stats.ErrorEvents++
	}
}

// getLogLevel 根據狀態決定日誌級別
func (cl *CentralLogger) getLogLevel(status string) string {
	switch status {
	case "success":
		return "INFO"
	case "fail", "invalid":
		return "WARN"
	case "error", "timeout", "locked":
		return "ERROR"
	default:
		return "INFO"
	}
}

// sendToExternalSystems 發送到外部日誌系統
func (cl *CentralLogger) sendToExternalSystems(event LogEvent) {
	// TODO: 實作發送到Loki, ElasticSearch等
	// 這裡可以根據配置選擇發送到不同的外部系統

	// 範例: 發送到Loki
	// cl.sendToLoki(event)

	// 範例: 發送到ElasticSearch
	// cl.sendToElasticsearch(event)
}

// GetStats 獲取日誌統計資訊
func (cl *CentralLogger) GetStats() LogStats {
	cl.statMutex.RLock()
	defer cl.statMutex.RUnlock()
	return cl.stats
}

// GetStatsJSON 獲取JSON格式的統計資訊
func (cl *CentralLogger) GetStatsJSON() (string, error) {
	stats := cl.GetStats()
	jsonData, err := json.Marshal(stats)
	if err != nil {
		return "", fmt.Errorf("序列化統計資訊失敗: %v", err)
	}
	return string(jsonData), nil
}

// QueryEvents 查詢日誌事件 (簡單實作)
func (cl *CentralLogger) QueryEvents(filters map[string]interface{}, limit int) []LogEvent {
	// TODO: 實作更複雜的查詢邏輯
	// 這裡可以整合資料庫或搜索引擎
	return []LogEvent{}
}

// Stop 停止中央日誌記錄器
func (cl *CentralLogger) Stop() {
	close(cl.stopChan)
	cl.wg.Wait()
	cl.logger.Info("中央日誌記錄器已停止")
}

// SetLevel 設定日誌級別
func (cl *CentralLogger) SetLevel(level string) {
	switch level {
	case "debug":
		cl.logger.SetLevel(logrus.DebugLevel)
	case "info":
		cl.logger.SetLevel(logrus.InfoLevel)
	case "warn":
		cl.logger.SetLevel(logrus.WarnLevel)
	case "error":
		cl.logger.SetLevel(logrus.ErrorLevel)
	default:
		cl.logger.SetLevel(logrus.InfoLevel)
	}
}
