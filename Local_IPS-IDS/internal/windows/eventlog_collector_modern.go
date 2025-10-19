// +build windows

package windows

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ModernEventLogCollector 使用 PowerShell 的現代事件日誌收集器
type ModernEventLogCollector struct {
	logger       *logrus.Logger
	logTypes     []string
	batchSize    int
	pollInterval time.Duration
	lastCollectTime map[string]time.Time
}

// EventRecord PowerShell 事件記錄結構
type EventRecord struct {
	XMLName      xml.Name `xml:"Event"`
	System       SystemData
	EventData    string `xml:"EventData"`
	UserData     string `xml:"UserData"`
}

// SystemData 系統數據
type SystemData struct {
	Provider      ProviderData
	EventID       int       `xml:"EventID"`
	Version       int       `xml:"Version"`
	Level         int       `xml:"Level"`
	Task          int       `xml:"Task"`
	Opcode        int       `xml:"Opcode"`
	Keywords      string    `xml:"Keywords"`
	TimeCreated   TimeData  `xml:"TimeCreated"`
	EventRecordID int64     `xml:"EventRecordID"`
	Correlation   string    `xml:"Correlation"`
	Execution     ExecutionData
	Channel       string    `xml:"Channel"`
	Computer      string    `xml:"Computer"`
	Security      SecurityData
}

// ProviderData 提供者數據
type ProviderData struct {
	Name string `xml:"Name,attr"`
	Guid string `xml:"Guid,attr"`
}

// TimeData 時間數據
type TimeData struct {
	SystemTime string `xml:"SystemTime,attr"`
}

// ExecutionData 執行數據
type ExecutionData struct {
	ProcessID int `xml:"ProcessID,attr"`
	ThreadID  int `xml:"ThreadID,attr"`
}

// SecurityData 安全數據
type SecurityData struct {
	UserID string `xml:"UserID,attr"`
}

// NewModernEventLogCollector 創建現代事件日誌收集器
func NewModernEventLogCollector(logger *logrus.Logger) *ModernEventLogCollector {
	return &ModernEventLogCollector{
		logger:          logger,
		logTypes:        []string{"System", "Security", "Application", "Setup"},
		batchSize:       100,
		pollInterval:    30 * time.Second,
		lastCollectTime: make(map[string]time.Time),
	}
}

// CollectLogs 使用 PowerShell 收集日誌
func (c *ModernEventLogCollector) CollectLogs(logType string, maxRecords int) ([]WindowsEventLog, error) {
	// 獲取上次收集時間
	lastTime := c.lastCollectTime[logType]
	if lastTime.IsZero() {
		// 首次收集，只收集最近 1 小時的日誌
		lastTime = time.Now().Add(-1 * time.Hour)
	}

	// 構建 PowerShell 命令
	// 使用 Get-WinEvent 替代舊的 Get-EventLog
	timeFilter := lastTime.Format("2006-01-02T15:04:05")
	
	psScript := fmt.Sprintf(`
		$logs = Get-WinEvent -LogName %s -MaxEvents %d -ErrorAction SilentlyContinue | Where-Object { $_.TimeCreated -gt [datetime]'%s' }
		$logs | ForEach-Object {
			[PSCustomObject]@{
				LogName = $_.LogName
				Source = $_.ProviderName
				EventID = $_.Id
				Level = $_.Level
				LevelDisplayName = $_.LevelDisplayName
				Message = $_.Message
				TimeCreated = $_.TimeCreated.ToString("o")
				UserID = $_.UserId
				ProcessID = $_.ProcessId
				ThreadID = $_.ThreadId
				MachineName = $_.MachineName
			}
		} | ConvertTo-Json -Compress
	`, logType, maxRecords, timeFilter)

	// 執行 PowerShell
	cmd := exec.Command("powershell", "-Command", psScript)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("powershell execution failed: %w, stderr: %s", err, stderr.String())
	}

	output := stdout.String()
	if output == "" || output == "null" {
		return nil, nil // 沒有新日誌
	}

	// 解析 JSON 結果
	var psLogs []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &psLogs); err != nil {
		// 可能是單個對象，嘗試解析為單個對象
		var singleLog map[string]interface{}
		if err := json.Unmarshal([]byte(output), &singleLog); err != nil {
			return nil, fmt.Errorf("parse json failed: %w", err)
		}
		psLogs = []map[string]interface{}{singleLog}
	}

	// 轉換為 WindowsEventLog
	var logs []WindowsEventLog
	for _, psLog := range psLogs {
		timeCreated, _ := time.Parse(time.RFC3339, psLog["TimeCreated"].(string))
		
		log := WindowsEventLog{
			LogType:     logType,
			Source:      fmt.Sprintf("%v", psLog["Source"]),
			EventID:     int(psLog["EventID"].(float64)),
			Level:       c.getLevelName(int(psLog["Level"].(float64))),
			Message:     fmt.Sprintf("%v", psLog["Message"]),
			TimeCreated: timeCreated,
			ProcessID:   int(psLog["ProcessID"].(float64)),
			ThreadID:    int(psLog["ThreadID"].(float64)),
			Metadata:    psLog,
		}

		if userID, ok := psLog["UserID"]; ok && userID != nil {
			log.UserID = fmt.Sprintf("%v", userID)
		}

		logs = append(logs, log)
	}

	// 更新最後收集時間
	if len(logs) > 0 {
		c.lastCollectTime[logType] = time.Now()
	}

	return logs, nil
}

// getLevelName 獲取級別名稱
func (c *ModernEventLogCollector) getLevelName(level int) string {
	switch level {
	case 1:
		return "Critical"
	case 2:
		return "Error"
	case 3:
		return "Warning"
	case 4:
		return "Information"
	case 0:
		return "Verbose"
	default:
		return "Unknown"
	}
}

// StartCollection 開始收集
func (c *ModernEventLogCollector) StartCollection(callback func([]WindowsEventLog) error) {
	c.logger.Info("Starting Windows Event Log collection...")
	
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()

	// 立即執行一次
	c.collectAndSend(callback)

	for range ticker.C {
		c.collectAndSend(callback)
	}
}

// collectAndSend 收集並發送日誌
func (c *ModernEventLogCollector) collectAndSend(callback func([]WindowsEventLog) error) {
	for _, logType := range c.logTypes {
		logs, err := c.CollectLogs(logType, c.batchSize)
		if err != nil {
			c.logger.Warnf("Failed to collect %s logs: %v", logType, err)
			continue
		}

		if len(logs) > 0 {
			c.logger.Infof("Collected %d logs from %s", len(logs), logType)
			if err := callback(logs); err != nil {
				c.logger.Errorf("Failed to process logs: %v", err)
			}
		}
	}
}

// SetPollInterval 設置輪詢間隔
func (c *ModernEventLogCollector) SetPollInterval(interval time.Duration) {
	c.pollInterval = interval
}

// SetBatchSize 設置批量大小
func (c *ModernEventLogCollector) SetBatchSize(size int) {
	c.batchSize = size
}



package windows

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ModernEventLogCollector 使用 PowerShell 的現代事件日誌收集器
type ModernEventLogCollector struct {
	logger       *logrus.Logger
	logTypes     []string
	batchSize    int
	pollInterval time.Duration
	lastCollectTime map[string]time.Time
}

// EventRecord PowerShell 事件記錄結構
type EventRecord struct {
	XMLName      xml.Name `xml:"Event"`
	System       SystemData
	EventData    string `xml:"EventData"`
	UserData     string `xml:"UserData"`
}

// SystemData 系統數據
type SystemData struct {
	Provider      ProviderData
	EventID       int       `xml:"EventID"`
	Version       int       `xml:"Version"`
	Level         int       `xml:"Level"`
	Task          int       `xml:"Task"`
	Opcode        int       `xml:"Opcode"`
	Keywords      string    `xml:"Keywords"`
	TimeCreated   TimeData  `xml:"TimeCreated"`
	EventRecordID int64     `xml:"EventRecordID"`
	Correlation   string    `xml:"Correlation"`
	Execution     ExecutionData
	Channel       string    `xml:"Channel"`
	Computer      string    `xml:"Computer"`
	Security      SecurityData
}

// ProviderData 提供者數據
type ProviderData struct {
	Name string `xml:"Name,attr"`
	Guid string `xml:"Guid,attr"`
}

// TimeData 時間數據
type TimeData struct {
	SystemTime string `xml:"SystemTime,attr"`
}

// ExecutionData 執行數據
type ExecutionData struct {
	ProcessID int `xml:"ProcessID,attr"`
	ThreadID  int `xml:"ThreadID,attr"`
}

// SecurityData 安全數據
type SecurityData struct {
	UserID string `xml:"UserID,attr"`
}

// NewModernEventLogCollector 創建現代事件日誌收集器
func NewModernEventLogCollector(logger *logrus.Logger) *ModernEventLogCollector {
	return &ModernEventLogCollector{
		logger:          logger,
		logTypes:        []string{"System", "Security", "Application", "Setup"},
		batchSize:       100,
		pollInterval:    30 * time.Second,
		lastCollectTime: make(map[string]time.Time),
	}
}

// CollectLogs 使用 PowerShell 收集日誌
func (c *ModernEventLogCollector) CollectLogs(logType string, maxRecords int) ([]WindowsEventLog, error) {
	// 獲取上次收集時間
	lastTime := c.lastCollectTime[logType]
	if lastTime.IsZero() {
		// 首次收集，只收集最近 1 小時的日誌
		lastTime = time.Now().Add(-1 * time.Hour)
	}

	// 構建 PowerShell 命令
	// 使用 Get-WinEvent 替代舊的 Get-EventLog
	timeFilter := lastTime.Format("2006-01-02T15:04:05")
	
	psScript := fmt.Sprintf(`
		$logs = Get-WinEvent -LogName %s -MaxEvents %d -ErrorAction SilentlyContinue | Where-Object { $_.TimeCreated -gt [datetime]'%s' }
		$logs | ForEach-Object {
			[PSCustomObject]@{
				LogName = $_.LogName
				Source = $_.ProviderName
				EventID = $_.Id
				Level = $_.Level
				LevelDisplayName = $_.LevelDisplayName
				Message = $_.Message
				TimeCreated = $_.TimeCreated.ToString("o")
				UserID = $_.UserId
				ProcessID = $_.ProcessId
				ThreadID = $_.ThreadId
				MachineName = $_.MachineName
			}
		} | ConvertTo-Json -Compress
	`, logType, maxRecords, timeFilter)

	// 執行 PowerShell
	cmd := exec.Command("powershell", "-Command", psScript)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("powershell execution failed: %w, stderr: %s", err, stderr.String())
	}

	output := stdout.String()
	if output == "" || output == "null" {
		return nil, nil // 沒有新日誌
	}

	// 解析 JSON 結果
	var psLogs []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &psLogs); err != nil {
		// 可能是單個對象，嘗試解析為單個對象
		var singleLog map[string]interface{}
		if err := json.Unmarshal([]byte(output), &singleLog); err != nil {
			return nil, fmt.Errorf("parse json failed: %w", err)
		}
		psLogs = []map[string]interface{}{singleLog}
	}

	// 轉換為 WindowsEventLog
	var logs []WindowsEventLog
	for _, psLog := range psLogs {
		timeCreated, _ := time.Parse(time.RFC3339, psLog["TimeCreated"].(string))
		
		log := WindowsEventLog{
			LogType:     logType,
			Source:      fmt.Sprintf("%v", psLog["Source"]),
			EventID:     int(psLog["EventID"].(float64)),
			Level:       c.getLevelName(int(psLog["Level"].(float64))),
			Message:     fmt.Sprintf("%v", psLog["Message"]),
			TimeCreated: timeCreated,
			ProcessID:   int(psLog["ProcessID"].(float64)),
			ThreadID:    int(psLog["ThreadID"].(float64)),
			Metadata:    psLog,
		}

		if userID, ok := psLog["UserID"]; ok && userID != nil {
			log.UserID = fmt.Sprintf("%v", userID)
		}

		logs = append(logs, log)
	}

	// 更新最後收集時間
	if len(logs) > 0 {
		c.lastCollectTime[logType] = time.Now()
	}

	return logs, nil
}

// getLevelName 獲取級別名稱
func (c *ModernEventLogCollector) getLevelName(level int) string {
	switch level {
	case 1:
		return "Critical"
	case 2:
		return "Error"
	case 3:
		return "Warning"
	case 4:
		return "Information"
	case 0:
		return "Verbose"
	default:
		return "Unknown"
	}
}

// StartCollection 開始收集
func (c *ModernEventLogCollector) StartCollection(callback func([]WindowsEventLog) error) {
	c.logger.Info("Starting Windows Event Log collection...")
	
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()

	// 立即執行一次
	c.collectAndSend(callback)

	for range ticker.C {
		c.collectAndSend(callback)
	}
}

// collectAndSend 收集並發送日誌
func (c *ModernEventLogCollector) collectAndSend(callback func([]WindowsEventLog) error) {
	for _, logType := range c.logTypes {
		logs, err := c.CollectLogs(logType, c.batchSize)
		if err != nil {
			c.logger.Warnf("Failed to collect %s logs: %v", logType, err)
			continue
		}

		if len(logs) > 0 {
			c.logger.Infof("Collected %d logs from %s", len(logs), logType)
			if err := callback(logs); err != nil {
				c.logger.Errorf("Failed to process logs: %v", err)
			}
		}
	}
}

// SetPollInterval 設置輪詢間隔
func (c *ModernEventLogCollector) SetPollInterval(interval time.Duration) {
	c.pollInterval = interval
}

// SetBatchSize 設置批量大小
func (c *ModernEventLogCollector) SetBatchSize(size int) {
	c.batchSize = size
}


package windows

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ModernEventLogCollector 使用 PowerShell 的現代事件日誌收集器
type ModernEventLogCollector struct {
	logger       *logrus.Logger
	logTypes     []string
	batchSize    int
	pollInterval time.Duration
	lastCollectTime map[string]time.Time
}

// EventRecord PowerShell 事件記錄結構
type EventRecord struct {
	XMLName      xml.Name `xml:"Event"`
	System       SystemData
	EventData    string `xml:"EventData"`
	UserData     string `xml:"UserData"`
}

// SystemData 系統數據
type SystemData struct {
	Provider      ProviderData
	EventID       int       `xml:"EventID"`
	Version       int       `xml:"Version"`
	Level         int       `xml:"Level"`
	Task          int       `xml:"Task"`
	Opcode        int       `xml:"Opcode"`
	Keywords      string    `xml:"Keywords"`
	TimeCreated   TimeData  `xml:"TimeCreated"`
	EventRecordID int64     `xml:"EventRecordID"`
	Correlation   string    `xml:"Correlation"`
	Execution     ExecutionData
	Channel       string    `xml:"Channel"`
	Computer      string    `xml:"Computer"`
	Security      SecurityData
}

// ProviderData 提供者數據
type ProviderData struct {
	Name string `xml:"Name,attr"`
	Guid string `xml:"Guid,attr"`
}

// TimeData 時間數據
type TimeData struct {
	SystemTime string `xml:"SystemTime,attr"`
}

// ExecutionData 執行數據
type ExecutionData struct {
	ProcessID int `xml:"ProcessID,attr"`
	ThreadID  int `xml:"ThreadID,attr"`
}

// SecurityData 安全數據
type SecurityData struct {
	UserID string `xml:"UserID,attr"`
}

// NewModernEventLogCollector 創建現代事件日誌收集器
func NewModernEventLogCollector(logger *logrus.Logger) *ModernEventLogCollector {
	return &ModernEventLogCollector{
		logger:          logger,
		logTypes:        []string{"System", "Security", "Application", "Setup"},
		batchSize:       100,
		pollInterval:    30 * time.Second,
		lastCollectTime: make(map[string]time.Time),
	}
}

// CollectLogs 使用 PowerShell 收集日誌
func (c *ModernEventLogCollector) CollectLogs(logType string, maxRecords int) ([]WindowsEventLog, error) {
	// 獲取上次收集時間
	lastTime := c.lastCollectTime[logType]
	if lastTime.IsZero() {
		// 首次收集，只收集最近 1 小時的日誌
		lastTime = time.Now().Add(-1 * time.Hour)
	}

	// 構建 PowerShell 命令
	// 使用 Get-WinEvent 替代舊的 Get-EventLog
	timeFilter := lastTime.Format("2006-01-02T15:04:05")
	
	psScript := fmt.Sprintf(`
		$logs = Get-WinEvent -LogName %s -MaxEvents %d -ErrorAction SilentlyContinue | Where-Object { $_.TimeCreated -gt [datetime]'%s' }
		$logs | ForEach-Object {
			[PSCustomObject]@{
				LogName = $_.LogName
				Source = $_.ProviderName
				EventID = $_.Id
				Level = $_.Level
				LevelDisplayName = $_.LevelDisplayName
				Message = $_.Message
				TimeCreated = $_.TimeCreated.ToString("o")
				UserID = $_.UserId
				ProcessID = $_.ProcessId
				ThreadID = $_.ThreadId
				MachineName = $_.MachineName
			}
		} | ConvertTo-Json -Compress
	`, logType, maxRecords, timeFilter)

	// 執行 PowerShell
	cmd := exec.Command("powershell", "-Command", psScript)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("powershell execution failed: %w, stderr: %s", err, stderr.String())
	}

	output := stdout.String()
	if output == "" || output == "null" {
		return nil, nil // 沒有新日誌
	}

	// 解析 JSON 結果
	var psLogs []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &psLogs); err != nil {
		// 可能是單個對象，嘗試解析為單個對象
		var singleLog map[string]interface{}
		if err := json.Unmarshal([]byte(output), &singleLog); err != nil {
			return nil, fmt.Errorf("parse json failed: %w", err)
		}
		psLogs = []map[string]interface{}{singleLog}
	}

	// 轉換為 WindowsEventLog
	var logs []WindowsEventLog
	for _, psLog := range psLogs {
		timeCreated, _ := time.Parse(time.RFC3339, psLog["TimeCreated"].(string))
		
		log := WindowsEventLog{
			LogType:     logType,
			Source:      fmt.Sprintf("%v", psLog["Source"]),
			EventID:     int(psLog["EventID"].(float64)),
			Level:       c.getLevelName(int(psLog["Level"].(float64))),
			Message:     fmt.Sprintf("%v", psLog["Message"]),
			TimeCreated: timeCreated,
			ProcessID:   int(psLog["ProcessID"].(float64)),
			ThreadID:    int(psLog["ThreadID"].(float64)),
			Metadata:    psLog,
		}

		if userID, ok := psLog["UserID"]; ok && userID != nil {
			log.UserID = fmt.Sprintf("%v", userID)
		}

		logs = append(logs, log)
	}

	// 更新最後收集時間
	if len(logs) > 0 {
		c.lastCollectTime[logType] = time.Now()
	}

	return logs, nil
}

// getLevelName 獲取級別名稱
func (c *ModernEventLogCollector) getLevelName(level int) string {
	switch level {
	case 1:
		return "Critical"
	case 2:
		return "Error"
	case 3:
		return "Warning"
	case 4:
		return "Information"
	case 0:
		return "Verbose"
	default:
		return "Unknown"
	}
}

// StartCollection 開始收集
func (c *ModernEventLogCollector) StartCollection(callback func([]WindowsEventLog) error) {
	c.logger.Info("Starting Windows Event Log collection...")
	
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()

	// 立即執行一次
	c.collectAndSend(callback)

	for range ticker.C {
		c.collectAndSend(callback)
	}
}

// collectAndSend 收集並發送日誌
func (c *ModernEventLogCollector) collectAndSend(callback func([]WindowsEventLog) error) {
	for _, logType := range c.logTypes {
		logs, err := c.CollectLogs(logType, c.batchSize)
		if err != nil {
			c.logger.Warnf("Failed to collect %s logs: %v", logType, err)
			continue
		}

		if len(logs) > 0 {
			c.logger.Infof("Collected %d logs from %s", len(logs), logType)
			if err := callback(logs); err != nil {
				c.logger.Errorf("Failed to process logs: %v", err)
			}
		}
	}
}

// SetPollInterval 設置輪詢間隔
func (c *ModernEventLogCollector) SetPollInterval(interval time.Duration) {
	c.pollInterval = interval
}

// SetBatchSize 設置批量大小
func (c *ModernEventLogCollector) SetBatchSize(size int) {
	c.batchSize = size
}



package windows

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ModernEventLogCollector 使用 PowerShell 的現代事件日誌收集器
type ModernEventLogCollector struct {
	logger       *logrus.Logger
	logTypes     []string
	batchSize    int
	pollInterval time.Duration
	lastCollectTime map[string]time.Time
}

// EventRecord PowerShell 事件記錄結構
type EventRecord struct {
	XMLName      xml.Name `xml:"Event"`
	System       SystemData
	EventData    string `xml:"EventData"`
	UserData     string `xml:"UserData"`
}

// SystemData 系統數據
type SystemData struct {
	Provider      ProviderData
	EventID       int       `xml:"EventID"`
	Version       int       `xml:"Version"`
	Level         int       `xml:"Level"`
	Task          int       `xml:"Task"`
	Opcode        int       `xml:"Opcode"`
	Keywords      string    `xml:"Keywords"`
	TimeCreated   TimeData  `xml:"TimeCreated"`
	EventRecordID int64     `xml:"EventRecordID"`
	Correlation   string    `xml:"Correlation"`
	Execution     ExecutionData
	Channel       string    `xml:"Channel"`
	Computer      string    `xml:"Computer"`
	Security      SecurityData
}

// ProviderData 提供者數據
type ProviderData struct {
	Name string `xml:"Name,attr"`
	Guid string `xml:"Guid,attr"`
}

// TimeData 時間數據
type TimeData struct {
	SystemTime string `xml:"SystemTime,attr"`
}

// ExecutionData 執行數據
type ExecutionData struct {
	ProcessID int `xml:"ProcessID,attr"`
	ThreadID  int `xml:"ThreadID,attr"`
}

// SecurityData 安全數據
type SecurityData struct {
	UserID string `xml:"UserID,attr"`
}

// NewModernEventLogCollector 創建現代事件日誌收集器
func NewModernEventLogCollector(logger *logrus.Logger) *ModernEventLogCollector {
	return &ModernEventLogCollector{
		logger:          logger,
		logTypes:        []string{"System", "Security", "Application", "Setup"},
		batchSize:       100,
		pollInterval:    30 * time.Second,
		lastCollectTime: make(map[string]time.Time),
	}
}

// CollectLogs 使用 PowerShell 收集日誌
func (c *ModernEventLogCollector) CollectLogs(logType string, maxRecords int) ([]WindowsEventLog, error) {
	// 獲取上次收集時間
	lastTime := c.lastCollectTime[logType]
	if lastTime.IsZero() {
		// 首次收集，只收集最近 1 小時的日誌
		lastTime = time.Now().Add(-1 * time.Hour)
	}

	// 構建 PowerShell 命令
	// 使用 Get-WinEvent 替代舊的 Get-EventLog
	timeFilter := lastTime.Format("2006-01-02T15:04:05")
	
	psScript := fmt.Sprintf(`
		$logs = Get-WinEvent -LogName %s -MaxEvents %d -ErrorAction SilentlyContinue | Where-Object { $_.TimeCreated -gt [datetime]'%s' }
		$logs | ForEach-Object {
			[PSCustomObject]@{
				LogName = $_.LogName
				Source = $_.ProviderName
				EventID = $_.Id
				Level = $_.Level
				LevelDisplayName = $_.LevelDisplayName
				Message = $_.Message
				TimeCreated = $_.TimeCreated.ToString("o")
				UserID = $_.UserId
				ProcessID = $_.ProcessId
				ThreadID = $_.ThreadId
				MachineName = $_.MachineName
			}
		} | ConvertTo-Json -Compress
	`, logType, maxRecords, timeFilter)

	// 執行 PowerShell
	cmd := exec.Command("powershell", "-Command", psScript)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("powershell execution failed: %w, stderr: %s", err, stderr.String())
	}

	output := stdout.String()
	if output == "" || output == "null" {
		return nil, nil // 沒有新日誌
	}

	// 解析 JSON 結果
	var psLogs []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &psLogs); err != nil {
		// 可能是單個對象，嘗試解析為單個對象
		var singleLog map[string]interface{}
		if err := json.Unmarshal([]byte(output), &singleLog); err != nil {
			return nil, fmt.Errorf("parse json failed: %w", err)
		}
		psLogs = []map[string]interface{}{singleLog}
	}

	// 轉換為 WindowsEventLog
	var logs []WindowsEventLog
	for _, psLog := range psLogs {
		timeCreated, _ := time.Parse(time.RFC3339, psLog["TimeCreated"].(string))
		
		log := WindowsEventLog{
			LogType:     logType,
			Source:      fmt.Sprintf("%v", psLog["Source"]),
			EventID:     int(psLog["EventID"].(float64)),
			Level:       c.getLevelName(int(psLog["Level"].(float64))),
			Message:     fmt.Sprintf("%v", psLog["Message"]),
			TimeCreated: timeCreated,
			ProcessID:   int(psLog["ProcessID"].(float64)),
			ThreadID:    int(psLog["ThreadID"].(float64)),
			Metadata:    psLog,
		}

		if userID, ok := psLog["UserID"]; ok && userID != nil {
			log.UserID = fmt.Sprintf("%v", userID)
		}

		logs = append(logs, log)
	}

	// 更新最後收集時間
	if len(logs) > 0 {
		c.lastCollectTime[logType] = time.Now()
	}

	return logs, nil
}

// getLevelName 獲取級別名稱
func (c *ModernEventLogCollector) getLevelName(level int) string {
	switch level {
	case 1:
		return "Critical"
	case 2:
		return "Error"
	case 3:
		return "Warning"
	case 4:
		return "Information"
	case 0:
		return "Verbose"
	default:
		return "Unknown"
	}
}

// StartCollection 開始收集
func (c *ModernEventLogCollector) StartCollection(callback func([]WindowsEventLog) error) {
	c.logger.Info("Starting Windows Event Log collection...")
	
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()

	// 立即執行一次
	c.collectAndSend(callback)

	for range ticker.C {
		c.collectAndSend(callback)
	}
}

// collectAndSend 收集並發送日誌
func (c *ModernEventLogCollector) collectAndSend(callback func([]WindowsEventLog) error) {
	for _, logType := range c.logTypes {
		logs, err := c.CollectLogs(logType, c.batchSize)
		if err != nil {
			c.logger.Warnf("Failed to collect %s logs: %v", logType, err)
			continue
		}

		if len(logs) > 0 {
			c.logger.Infof("Collected %d logs from %s", len(logs), logType)
			if err := callback(logs); err != nil {
				c.logger.Errorf("Failed to process logs: %v", err)
			}
		}
	}
}

// SetPollInterval 設置輪詢間隔
func (c *ModernEventLogCollector) SetPollInterval(interval time.Duration) {
	c.pollInterval = interval
}

// SetBatchSize 設置批量大小
func (c *ModernEventLogCollector) SetBatchSize(size int) {
	c.batchSize = size
}

