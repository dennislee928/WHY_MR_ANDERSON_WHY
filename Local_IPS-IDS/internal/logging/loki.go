package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// LokiClient Loki日誌客戶端
type LokiClient struct {
	logger     *logrus.Logger
	baseURL    string
	username   string
	password   string
	labels     map[string]string
	httpClient *http.Client
}

// LogEntry Loki日誌條目
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Labels    map[string]string      `json:"labels"`
	Fields    map[string]interface{} `json:"fields"`
}

// LokiPushRequest Loki推送請求格式
type LokiPushRequest struct {
	Streams []LokiStream `json:"streams"`
}

// LokiStream Loki串流格式
type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// NewLokiClient 建立新的Loki客戶端
func NewLokiClient(logger *logrus.Logger, baseURL, username, password string) *LokiClient {
	return &LokiClient{
		logger:   logger,
		baseURL:  baseURL,
		username: username,
		password: password,
		labels: map[string]string{
			"service":     "pandora-box-console",
			"environment": "production",
		},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SetLabels 設定全域標籤
func (lc *LokiClient) SetLabels(labels map[string]string) {
	for k, v := range labels {
		lc.labels[k] = v
	}
}

// PushLog 推送日誌到Loki
func (lc *LokiClient) PushLog(entry LogEntry) error {
	// 合併標籤
	streamLabels := make(map[string]string)
	for k, v := range lc.labels {
		streamLabels[k] = v
	}
	for k, v := range entry.Labels {
		streamLabels[k] = v
	}
	streamLabels["level"] = entry.Level

	// 建立日誌訊息
	logMessage := map[string]interface{}{
		"message": entry.Message,
		"fields":  entry.Fields,
	}

	messageBytes, err := json.Marshal(logMessage)
	if err != nil {
		return fmt.Errorf("序列化日誌訊息失敗: %v", err)
	}

	// 建立Loki推送請求
	pushRequest := LokiPushRequest{
		Streams: []LokiStream{
			{
				Stream: streamLabels,
				Values: [][]string{
					{
						strconv.FormatInt(entry.Timestamp.UnixNano(), 10),
						string(messageBytes),
					},
				},
			},
		},
	}

	return lc.sendToLoki(pushRequest)
}

// PushBatchLogs 批次推送日誌到Loki
func (lc *LokiClient) PushBatchLogs(entries []LogEntry) error {
	if len(entries) == 0 {
		return nil
	}

	streams := make(map[string]*LokiStream)

	for _, entry := range entries {
		// 合併標籤
		streamLabels := make(map[string]string)
		for k, v := range lc.labels {
			streamLabels[k] = v
		}
		for k, v := range entry.Labels {
			streamLabels[k] = v
		}
		streamLabels["level"] = entry.Level

		// 建立串流鍵
		streamKey := lc.createStreamKey(streamLabels)

		// 建立日誌訊息
		logMessage := map[string]interface{}{
			"message": entry.Message,
			"fields":  entry.Fields,
		}

		messageBytes, err := json.Marshal(logMessage)
		if err != nil {
			lc.logger.Errorf("序列化日誌訊息失敗: %v", err)
			continue
		}

		// 添加到對應的串流
		if streams[streamKey] == nil {
			streams[streamKey] = &LokiStream{
				Stream: streamLabels,
				Values: [][]string{},
			}
		}

		streams[streamKey].Values = append(streams[streamKey].Values, []string{
			strconv.FormatInt(entry.Timestamp.UnixNano(), 10),
			string(messageBytes),
		})
	}

	// 轉換為切片
	streamSlice := make([]LokiStream, 0, len(streams))
	for _, stream := range streams {
		streamSlice = append(streamSlice, *stream)
	}

	pushRequest := LokiPushRequest{
		Streams: streamSlice,
	}

	return lc.sendToLoki(pushRequest)
}

// sendToLoki 發送請求到Loki
func (lc *LokiClient) sendToLoki(pushRequest LokiPushRequest) error {
	jsonData, err := json.Marshal(pushRequest)
	if err != nil {
		return fmt.Errorf("序列化Loki請求失敗: %v", err)
	}

	url := fmt.Sprintf("%s/loki/api/v1/push", lc.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("建立HTTP請求失敗: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if lc.username != "" && lc.password != "" {
		req.SetBasicAuth(lc.username, lc.password)
	}

	resp, err := lc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("發送HTTP請求失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Loki API回應錯誤: %d", resp.StatusCode)
	}

	return nil
}

// createStreamKey 建立串流鍵
func (lc *LokiClient) createStreamKey(labels map[string]string) string {
	var key bytes.Buffer
	for k, v := range labels {
		key.WriteString(k)
		key.WriteString("=")
		key.WriteString(v)
		key.WriteString(",")
	}
	return key.String()
}

// LogrusHook Logrus Hook for Loki
type LogrusHook struct {
	client     *LokiClient
	levels     []logrus.Level
	batchSize  int
	buffer     []LogEntry
	flushTimer *time.Timer
}

// NewLogrusHook 建立新的Logrus Hook
func NewLogrusHook(client *LokiClient, levels []logrus.Level) *LogrusHook {
	hook := &LogrusHook{
		client:     client,
		levels:     levels,
		batchSize:  100,
		buffer:     make([]LogEntry, 0),
		flushTimer: time.NewTimer(5 * time.Second),
	}

	// 啟動定期刷新
	go hook.periodicFlush()

	return hook
}

// Levels 返回支援的日誌等級
func (hook *LogrusHook) Levels() []logrus.Level {
	return hook.levels
}

// Fire 處理日誌事件
func (hook *LogrusHook) Fire(entry *logrus.Entry) error {
	logEntry := LogEntry{
		Timestamp: entry.Time,
		Level:     entry.Level.String(),
		Message:   entry.Message,
		Labels:    make(map[string]string),
		Fields:    make(map[string]interface{}),
	}

	// 轉換欄位
	for k, v := range entry.Data {
		if strVal, ok := v.(string); ok && (k == "component" || k == "module" || k == "operation") {
			logEntry.Labels[k] = strVal
		} else {
			logEntry.Fields[k] = v
		}
	}

	hook.buffer = append(hook.buffer, logEntry)

	// 如果緩衝區滿了，立即刷新
	if len(hook.buffer) >= hook.batchSize {
		return hook.flush()
	}

	return nil
}

// flush 刷新緩衝區
func (hook *LogrusHook) flush() error {
	if len(hook.buffer) == 0 {
		return nil
	}

	err := hook.client.PushBatchLogs(hook.buffer)
	hook.buffer = hook.buffer[:0] // 清空緩衝區
	return err
}

// periodicFlush 定期刷新緩衝區
func (hook *LogrusHook) periodicFlush() {
	for {
		<-hook.flushTimer.C
		hook.flush()
		hook.flushTimer.Reset(5 * time.Second)
	}
}

// SecurityLogger 安全事件專用日誌記錄器
type SecurityLogger struct {
	client *LokiClient
	logger *logrus.Logger
}

// NewSecurityLogger 建立新的安全事件日誌記錄器
func NewSecurityLogger(client *LokiClient, logger *logrus.Logger) *SecurityLogger {
	return &SecurityLogger{
		client: client,
		logger: logger,
	}
}

// LogThreatDetection 記錄威脅偵測事件
func (sl *SecurityLogger) LogThreatDetection(threatType, severity, sourceIP, details string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "warning",
		Message:   fmt.Sprintf("威脅偵測: %s", threatType),
		Labels: map[string]string{
			"event_type": "threat_detection",
			"severity":   severity,
			"source_ip":  sourceIP,
		},
		Fields: map[string]interface{}{
			"threat_type": threatType,
			"details":     details,
			"timestamp":   time.Now().Unix(),
		},
	}

	if err := sl.client.PushLog(entry); err != nil {
		sl.logger.Errorf("記錄威脅偵測事件失敗: %v", err)
	}
}

// LogSecurityEvent 記錄一般安全事件
func (sl *SecurityLogger) LogSecurityEvent(eventType, action, details string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   fmt.Sprintf("安全事件: %s - %s", eventType, action),
		Labels: map[string]string{
			"event_type": "security_event",
			"action":     action,
		},
		Fields: map[string]interface{}{
			"event_type": eventType,
			"details":    details,
			"timestamp":  time.Now().Unix(),
		},
	}

	if err := sl.client.PushLog(entry); err != nil {
		sl.logger.Errorf("記錄安全事件失敗: %v", err)
	}
}

// LogNetworkEvent 記錄網路事件
func (sl *SecurityLogger) LogNetworkEvent(eventType, status, reason string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   fmt.Sprintf("網路事件: %s - %s", eventType, status),
		Labels: map[string]string{
			"event_type": "network_event",
			"status":     status,
		},
		Fields: map[string]interface{}{
			"event_type": eventType,
			"reason":     reason,
			"timestamp":  time.Now().Unix(),
		},
	}

	if err := sl.client.PushLog(entry); err != nil {
		sl.logger.Errorf("記錄網路事件失敗: %v", err)
	}
}
