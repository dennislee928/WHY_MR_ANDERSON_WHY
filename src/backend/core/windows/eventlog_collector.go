// +build windows

package windows

import (
	"encoding/json"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/sirupsen/logrus"
)

var (
	// Windows DLL
	advapi32 = syscall.NewLazyDLL("advapi32.dll")
	
	// Windows API 函數
	openEventLogW       = advapi32.NewProc("OpenEventLogW")
	closeEventLog       = advapi32.NewProc("CloseEventLog")
	readEventLogW       = advapi32.NewProc("ReadEventLogW")
	getNumberOfEventLogRecords = advapi32.NewProc("GetNumberOfEventLogRecords")
)

const (
	// 事件類型
	EVENTLOG_ERROR_TYPE       = 0x0001
	EVENTLOG_WARNING_TYPE     = 0x0002
	EVENTLOG_INFORMATION_TYPE = 0x0004
	EVENTLOG_AUDIT_SUCCESS    = 0x0008
	EVENTLOG_AUDIT_FAILURE    = 0x0010

	// 讀取標誌
	EVENTLOG_SEQUENTIAL_READ = 0x0001
	EVENTLOG_FORWARDS_READ   = 0x0004
)

// EventLogCollector Windows 事件日誌收集器
type EventLogCollector struct {
	logger      *logrus.Logger
	logTypes    []string // System, Security, Application, Setup
	batchSize   int
	pollInterval time.Duration
}

// WindowsEventLog Windows 事件日誌結構
type WindowsEventLog struct {
	LogType     string                 `json:"log_type"`
	Source      string                 `json:"source"`
	EventID     int                    `json:"event_id"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	TimeCreated time.Time              `json:"time_created"`
	UserID      string                 `json:"user_id,omitempty"`
	ProcessID   int                    `json:"process_id,omitempty"`
	ThreadID    int                    `json:"thread_id,omitempty"`
	Keywords    string                 `json:"keywords,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewEventLogCollector 創建新的事件日誌收集器
func NewEventLogCollector(logger *logrus.Logger) *EventLogCollector {
	return &EventLogCollector{
		logger:       logger,
		logTypes:     []string{"System", "Security", "Application", "Setup"},
		batchSize:    100,
		pollInterval: 30 * time.Second,
	}
}

// CollectLogs 收集日誌
func (c *EventLogCollector) CollectLogs(logType string, maxRecords int) ([]WindowsEventLog, error) {
	// 打開事件日誌
	logNamePtr, err := syscall.UTF16PtrFromString(logType)
	if err != nil {
		return nil, fmt.Errorf("convert log name: %w", err)
	}

	handle, _, err := openEventLogW.Call(
		0, // 本地機器
		uintptr(unsafe.Pointer(logNamePtr)),
	)
	if handle == 0 {
		return nil, fmt.Errorf("open event log failed: %w", err)
	}
	defer closeEventLog.Call(handle)

	var logs []WindowsEventLog
	buffer := make([]byte, 0x10000) // 64KB buffer

	for len(logs) < maxRecords {
		var bytesRead uint32
		var minBytesNeeded uint32

		ret, _, _ := readEventLogW.Call(
			handle,
			EVENTLOG_SEQUENTIAL_READ|EVENTLOG_FORWARDS_READ,
			0, // 從當前位置開始讀取
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(len(buffer)),
			uintptr(unsafe.Pointer(&bytesRead)),
			uintptr(unsafe.Pointer(&minBytesNeeded)),
		)

		if ret == 0 {
			// 沒有更多記錄
			break
		}

		// 解析事件記錄
		parsedLogs, err := c.parseEventRecords(buffer[:bytesRead], logType)
		if err != nil {
			c.logger.Warnf("Failed to parse event records: %v", err)
			continue
		}

		logs = append(logs, parsedLogs...)

		if len(logs) >= maxRecords {
			break
		}
	}

	return logs, nil
}

// parseEventRecords 解析事件記錄
func (c *EventLogCollector) parseEventRecords(buffer []byte, logType string) ([]WindowsEventLog, error) {
	var logs []WindowsEventLog

	// Note: 實際實現需要正確解析 EVENTLOGRECORD 結構
	// 這裡提供簡化版本
	
	// 實際應該使用 Windows Event Log API v2 (EvtQuery)
	// 或使用 golang.org/x/sys/windows 套件

	c.logger.Debugf("Parsed %d event log records from %s", len(logs), logType)
	return logs, nil
}

// getLevelString 獲取級別字符串
func (c *EventLogCollector) getLevelString(eventType uint16) string {
	switch eventType {
	case EVENTLOG_ERROR_TYPE:
		return "Error"
	case EVENTLOG_WARNING_TYPE:
		return "Warning"
	case EVENTLOG_INFORMATION_TYPE:
		return "Information"
	case EVENTLOG_AUDIT_SUCCESS:
		return "Audit Success"
	case EVENTLOG_AUDIT_FAILURE:
		return "Audit Failure"
	default:
		return "Unknown"
	}
}

// StartCollection 開始收集（定期輪詢）
func (c *EventLogCollector) StartCollection(callback func([]WindowsEventLog) error) {
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()

	for range ticker.C {
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
}



package windows

import (
	"encoding/json"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/sirupsen/logrus"
)

var (
	// Windows DLL
	advapi32 = syscall.NewLazyDLL("advapi32.dll")
	
	// Windows API 函數
	openEventLogW       = advapi32.NewProc("OpenEventLogW")
	closeEventLog       = advapi32.NewProc("CloseEventLog")
	readEventLogW       = advapi32.NewProc("ReadEventLogW")
	getNumberOfEventLogRecords = advapi32.NewProc("GetNumberOfEventLogRecords")
)

const (
	// 事件類型
	EVENTLOG_ERROR_TYPE       = 0x0001
	EVENTLOG_WARNING_TYPE     = 0x0002
	EVENTLOG_INFORMATION_TYPE = 0x0004
	EVENTLOG_AUDIT_SUCCESS    = 0x0008
	EVENTLOG_AUDIT_FAILURE    = 0x0010

	// 讀取標誌
	EVENTLOG_SEQUENTIAL_READ = 0x0001
	EVENTLOG_FORWARDS_READ   = 0x0004
)

// EventLogCollector Windows 事件日誌收集器
type EventLogCollector struct {
	logger      *logrus.Logger
	logTypes    []string // System, Security, Application, Setup
	batchSize   int
	pollInterval time.Duration
}

// WindowsEventLog Windows 事件日誌結構
type WindowsEventLog struct {
	LogType     string                 `json:"log_type"`
	Source      string                 `json:"source"`
	EventID     int                    `json:"event_id"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	TimeCreated time.Time              `json:"time_created"`
	UserID      string                 `json:"user_id,omitempty"`
	ProcessID   int                    `json:"process_id,omitempty"`
	ThreadID    int                    `json:"thread_id,omitempty"`
	Keywords    string                 `json:"keywords,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewEventLogCollector 創建新的事件日誌收集器
func NewEventLogCollector(logger *logrus.Logger) *EventLogCollector {
	return &EventLogCollector{
		logger:       logger,
		logTypes:     []string{"System", "Security", "Application", "Setup"},
		batchSize:    100,
		pollInterval: 30 * time.Second,
	}
}

// CollectLogs 收集日誌
func (c *EventLogCollector) CollectLogs(logType string, maxRecords int) ([]WindowsEventLog, error) {
	// 打開事件日誌
	logNamePtr, err := syscall.UTF16PtrFromString(logType)
	if err != nil {
		return nil, fmt.Errorf("convert log name: %w", err)
	}

	handle, _, err := openEventLogW.Call(
		0, // 本地機器
		uintptr(unsafe.Pointer(logNamePtr)),
	)
	if handle == 0 {
		return nil, fmt.Errorf("open event log failed: %w", err)
	}
	defer closeEventLog.Call(handle)

	var logs []WindowsEventLog
	buffer := make([]byte, 0x10000) // 64KB buffer

	for len(logs) < maxRecords {
		var bytesRead uint32
		var minBytesNeeded uint32

		ret, _, _ := readEventLogW.Call(
			handle,
			EVENTLOG_SEQUENTIAL_READ|EVENTLOG_FORWARDS_READ,
			0, // 從當前位置開始讀取
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(len(buffer)),
			uintptr(unsafe.Pointer(&bytesRead)),
			uintptr(unsafe.Pointer(&minBytesNeeded)),
		)

		if ret == 0 {
			// 沒有更多記錄
			break
		}

		// 解析事件記錄
		parsedLogs, err := c.parseEventRecords(buffer[:bytesRead], logType)
		if err != nil {
			c.logger.Warnf("Failed to parse event records: %v", err)
			continue
		}

		logs = append(logs, parsedLogs...)

		if len(logs) >= maxRecords {
			break
		}
	}

	return logs, nil
}

// parseEventRecords 解析事件記錄
func (c *EventLogCollector) parseEventRecords(buffer []byte, logType string) ([]WindowsEventLog, error) {
	var logs []WindowsEventLog

	// Note: 實際實現需要正確解析 EVENTLOGRECORD 結構
	// 這裡提供簡化版本
	
	// 實際應該使用 Windows Event Log API v2 (EvtQuery)
	// 或使用 golang.org/x/sys/windows 套件

	c.logger.Debugf("Parsed %d event log records from %s", len(logs), logType)
	return logs, nil
}

// getLevelString 獲取級別字符串
func (c *EventLogCollector) getLevelString(eventType uint16) string {
	switch eventType {
	case EVENTLOG_ERROR_TYPE:
		return "Error"
	case EVENTLOG_WARNING_TYPE:
		return "Warning"
	case EVENTLOG_INFORMATION_TYPE:
		return "Information"
	case EVENTLOG_AUDIT_SUCCESS:
		return "Audit Success"
	case EVENTLOG_AUDIT_FAILURE:
		return "Audit Failure"
	default:
		return "Unknown"
	}
}

// StartCollection 開始收集（定期輪詢）
func (c *EventLogCollector) StartCollection(callback func([]WindowsEventLog) error) {
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()

	for range ticker.C {
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
}


package windows

import (
	"encoding/json"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/sirupsen/logrus"
)

var (
	// Windows DLL
	advapi32 = syscall.NewLazyDLL("advapi32.dll")
	
	// Windows API 函數
	openEventLogW       = advapi32.NewProc("OpenEventLogW")
	closeEventLog       = advapi32.NewProc("CloseEventLog")
	readEventLogW       = advapi32.NewProc("ReadEventLogW")
	getNumberOfEventLogRecords = advapi32.NewProc("GetNumberOfEventLogRecords")
)

const (
	// 事件類型
	EVENTLOG_ERROR_TYPE       = 0x0001
	EVENTLOG_WARNING_TYPE     = 0x0002
	EVENTLOG_INFORMATION_TYPE = 0x0004
	EVENTLOG_AUDIT_SUCCESS    = 0x0008
	EVENTLOG_AUDIT_FAILURE    = 0x0010

	// 讀取標誌
	EVENTLOG_SEQUENTIAL_READ = 0x0001
	EVENTLOG_FORWARDS_READ   = 0x0004
)

// EventLogCollector Windows 事件日誌收集器
type EventLogCollector struct {
	logger      *logrus.Logger
	logTypes    []string // System, Security, Application, Setup
	batchSize   int
	pollInterval time.Duration
}

// WindowsEventLog Windows 事件日誌結構
type WindowsEventLog struct {
	LogType     string                 `json:"log_type"`
	Source      string                 `json:"source"`
	EventID     int                    `json:"event_id"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	TimeCreated time.Time              `json:"time_created"`
	UserID      string                 `json:"user_id,omitempty"`
	ProcessID   int                    `json:"process_id,omitempty"`
	ThreadID    int                    `json:"thread_id,omitempty"`
	Keywords    string                 `json:"keywords,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewEventLogCollector 創建新的事件日誌收集器
func NewEventLogCollector(logger *logrus.Logger) *EventLogCollector {
	return &EventLogCollector{
		logger:       logger,
		logTypes:     []string{"System", "Security", "Application", "Setup"},
		batchSize:    100,
		pollInterval: 30 * time.Second,
	}
}

// CollectLogs 收集日誌
func (c *EventLogCollector) CollectLogs(logType string, maxRecords int) ([]WindowsEventLog, error) {
	// 打開事件日誌
	logNamePtr, err := syscall.UTF16PtrFromString(logType)
	if err != nil {
		return nil, fmt.Errorf("convert log name: %w", err)
	}

	handle, _, err := openEventLogW.Call(
		0, // 本地機器
		uintptr(unsafe.Pointer(logNamePtr)),
	)
	if handle == 0 {
		return nil, fmt.Errorf("open event log failed: %w", err)
	}
	defer closeEventLog.Call(handle)

	var logs []WindowsEventLog
	buffer := make([]byte, 0x10000) // 64KB buffer

	for len(logs) < maxRecords {
		var bytesRead uint32
		var minBytesNeeded uint32

		ret, _, _ := readEventLogW.Call(
			handle,
			EVENTLOG_SEQUENTIAL_READ|EVENTLOG_FORWARDS_READ,
			0, // 從當前位置開始讀取
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(len(buffer)),
			uintptr(unsafe.Pointer(&bytesRead)),
			uintptr(unsafe.Pointer(&minBytesNeeded)),
		)

		if ret == 0 {
			// 沒有更多記錄
			break
		}

		// 解析事件記錄
		parsedLogs, err := c.parseEventRecords(buffer[:bytesRead], logType)
		if err != nil {
			c.logger.Warnf("Failed to parse event records: %v", err)
			continue
		}

		logs = append(logs, parsedLogs...)

		if len(logs) >= maxRecords {
			break
		}
	}

	return logs, nil
}

// parseEventRecords 解析事件記錄
func (c *EventLogCollector) parseEventRecords(buffer []byte, logType string) ([]WindowsEventLog, error) {
	var logs []WindowsEventLog

	// Note: 實際實現需要正確解析 EVENTLOGRECORD 結構
	// 這裡提供簡化版本
	
	// 實際應該使用 Windows Event Log API v2 (EvtQuery)
	// 或使用 golang.org/x/sys/windows 套件

	c.logger.Debugf("Parsed %d event log records from %s", len(logs), logType)
	return logs, nil
}

// getLevelString 獲取級別字符串
func (c *EventLogCollector) getLevelString(eventType uint16) string {
	switch eventType {
	case EVENTLOG_ERROR_TYPE:
		return "Error"
	case EVENTLOG_WARNING_TYPE:
		return "Warning"
	case EVENTLOG_INFORMATION_TYPE:
		return "Information"
	case EVENTLOG_AUDIT_SUCCESS:
		return "Audit Success"
	case EVENTLOG_AUDIT_FAILURE:
		return "Audit Failure"
	default:
		return "Unknown"
	}
}

// StartCollection 開始收集（定期輪詢）
func (c *EventLogCollector) StartCollection(callback func([]WindowsEventLog) error) {
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()

	for range ticker.C {
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
}



package windows

import (
	"encoding/json"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/sirupsen/logrus"
)

var (
	// Windows DLL
	advapi32 = syscall.NewLazyDLL("advapi32.dll")
	
	// Windows API 函數
	openEventLogW       = advapi32.NewProc("OpenEventLogW")
	closeEventLog       = advapi32.NewProc("CloseEventLog")
	readEventLogW       = advapi32.NewProc("ReadEventLogW")
	getNumberOfEventLogRecords = advapi32.NewProc("GetNumberOfEventLogRecords")
)

const (
	// 事件類型
	EVENTLOG_ERROR_TYPE       = 0x0001
	EVENTLOG_WARNING_TYPE     = 0x0002
	EVENTLOG_INFORMATION_TYPE = 0x0004
	EVENTLOG_AUDIT_SUCCESS    = 0x0008
	EVENTLOG_AUDIT_FAILURE    = 0x0010

	// 讀取標誌
	EVENTLOG_SEQUENTIAL_READ = 0x0001
	EVENTLOG_FORWARDS_READ   = 0x0004
)

// EventLogCollector Windows 事件日誌收集器
type EventLogCollector struct {
	logger      *logrus.Logger
	logTypes    []string // System, Security, Application, Setup
	batchSize   int
	pollInterval time.Duration
}

// WindowsEventLog Windows 事件日誌結構
type WindowsEventLog struct {
	LogType     string                 `json:"log_type"`
	Source      string                 `json:"source"`
	EventID     int                    `json:"event_id"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	TimeCreated time.Time              `json:"time_created"`
	UserID      string                 `json:"user_id,omitempty"`
	ProcessID   int                    `json:"process_id,omitempty"`
	ThreadID    int                    `json:"thread_id,omitempty"`
	Keywords    string                 `json:"keywords,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewEventLogCollector 創建新的事件日誌收集器
func NewEventLogCollector(logger *logrus.Logger) *EventLogCollector {
	return &EventLogCollector{
		logger:       logger,
		logTypes:     []string{"System", "Security", "Application", "Setup"},
		batchSize:    100,
		pollInterval: 30 * time.Second,
	}
}

// CollectLogs 收集日誌
func (c *EventLogCollector) CollectLogs(logType string, maxRecords int) ([]WindowsEventLog, error) {
	// 打開事件日誌
	logNamePtr, err := syscall.UTF16PtrFromString(logType)
	if err != nil {
		return nil, fmt.Errorf("convert log name: %w", err)
	}

	handle, _, err := openEventLogW.Call(
		0, // 本地機器
		uintptr(unsafe.Pointer(logNamePtr)),
	)
	if handle == 0 {
		return nil, fmt.Errorf("open event log failed: %w", err)
	}
	defer closeEventLog.Call(handle)

	var logs []WindowsEventLog
	buffer := make([]byte, 0x10000) // 64KB buffer

	for len(logs) < maxRecords {
		var bytesRead uint32
		var minBytesNeeded uint32

		ret, _, _ := readEventLogW.Call(
			handle,
			EVENTLOG_SEQUENTIAL_READ|EVENTLOG_FORWARDS_READ,
			0, // 從當前位置開始讀取
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(len(buffer)),
			uintptr(unsafe.Pointer(&bytesRead)),
			uintptr(unsafe.Pointer(&minBytesNeeded)),
		)

		if ret == 0 {
			// 沒有更多記錄
			break
		}

		// 解析事件記錄
		parsedLogs, err := c.parseEventRecords(buffer[:bytesRead], logType)
		if err != nil {
			c.logger.Warnf("Failed to parse event records: %v", err)
			continue
		}

		logs = append(logs, parsedLogs...)

		if len(logs) >= maxRecords {
			break
		}
	}

	return logs, nil
}

// parseEventRecords 解析事件記錄
func (c *EventLogCollector) parseEventRecords(buffer []byte, logType string) ([]WindowsEventLog, error) {
	var logs []WindowsEventLog

	// Note: 實際實現需要正確解析 EVENTLOGRECORD 結構
	// 這裡提供簡化版本
	
	// 實際應該使用 Windows Event Log API v2 (EvtQuery)
	// 或使用 golang.org/x/sys/windows 套件

	c.logger.Debugf("Parsed %d event log records from %s", len(logs), logType)
	return logs, nil
}

// getLevelString 獲取級別字符串
func (c *EventLogCollector) getLevelString(eventType uint16) string {
	switch eventType {
	case EVENTLOG_ERROR_TYPE:
		return "Error"
	case EVENTLOG_WARNING_TYPE:
		return "Warning"
	case EVENTLOG_INFORMATION_TYPE:
		return "Information"
	case EVENTLOG_AUDIT_SUCCESS:
		return "Audit Success"
	case EVENTLOG_AUDIT_FAILURE:
		return "Audit Failure"
	default:
		return "Unknown"
	}
}

// StartCollection 開始收集（定期輪詢）
func (c *EventLogCollector) StartCollection(callback func([]WindowsEventLog) error) {
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()

	for range ticker.C {
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
}

