package vo

import "time"

// WindowsLogVO Windows 日誌響應
type WindowsLogVO struct {
	ID          uint                   `json:"id"`
	AgentID     string                 `json:"agent_id"`
	LogType     string                 `json:"log_type"`
	Source      string                 `json:"source"`
	EventID     int                    `json:"event_id"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	TimeCreated time.Time              `json:"time_created"`
	ReceivedAt  time.Time              `json:"received_at"`
	Computer    string                 `json:"computer,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// WindowsLogsListVO Windows 日誌列表響應
type WindowsLogsListVO struct {
	Logs       []WindowsLogVO `json:"logs"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
	Timestamp  time.Time      `json:"timestamp"`
}

// WindowsLogBatchVO Windows 日誌批量上報響應
type WindowsLogBatchVO struct {
	ReceivedCount int       `json:"received_count"`
	SavedCount    int       `json:"saved_count"`
	FailedCount   int       `json:"failed_count"`
	Errors        []string  `json:"errors,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
	Message       string    `json:"message"`
}

// WindowsLogStatsVO Windows 日誌統計響應
type WindowsLogStatsVO struct {
	TotalLogs         int                 `json:"total_logs"`
	LogsByType        map[string]int      `json:"logs_by_type"`
	LogsByLevel       map[string]int      `json:"logs_by_level"`
	LogsByAgent       map[string]int      `json:"logs_by_agent"`
	CriticalCount     int                 `json:"critical_count"`
	ErrorCount        int                 `json:"error_count"`
	WarningCount      int                 `json:"warning_count"`
	Last24Hours       int                 `json:"last_24_hours"`
	TopSources        []LogSourceCount    `json:"top_sources"`
	TopEventIDs       []LogEventIDCount   `json:"top_event_ids"`
	TimeRange         string              `json:"time_range"`
}

// LogSourceCount 日誌來源計數
type LogSourceCount struct {
	Source string `json:"source"`
	Count  int    `json:"count"`
}

// LogEventIDCount 事件 ID 計數
type LogEventIDCount struct {
	EventID int    `json:"event_id"`
	Count   int    `json:"count"`
	Message string `json:"message,omitempty"`
}

