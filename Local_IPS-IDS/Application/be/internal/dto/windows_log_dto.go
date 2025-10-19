package dto

import "time"

// WindowsLogBatchRequest Windows 日誌批量上報請求
type WindowsLogBatchRequest struct {
	AgentID  string           `json:"agent_id" binding:"required"`
	Computer string           `json:"computer"`
	Logs     []WindowsLogItem `json:"logs" binding:"required,min=1,max=1000"` // 最多1000條
}

// WindowsLogItem 單條 Windows 日誌
type WindowsLogItem struct {
	LogType     string                 `json:"log_type" binding:"required"` // System, Security, Application, Setup
	Source      string                 `json:"source"`
	EventID     int                    `json:"event_id"`
	Level       string                 `json:"level"` // Information, Warning, Error, Critical
	Message     string                 `json:"message"`
	TimeCreated time.Time              `json:"time_created"`
	UserID      string                 `json:"user_id,omitempty"`
	ProcessID   int                    `json:"process_id,omitempty"`
	ThreadID    int                    `json:"thread_id,omitempty"`
	Keywords    string                 `json:"keywords,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"` // 完整元數據
}

// WindowsLogQueryRequest Windows 日誌查詢請求
type WindowsLogQueryRequest struct {
	AgentID     string   `form:"agent_id"`
	LogType     string   `form:"log_type"` // System, Security, Application
	Source      string   `form:"source"`
	EventID     int      `form:"event_id"`
	Level       string   `form:"level"` // Information, Warning, Error, Critical
	StartTime   string   `form:"start_time"` // RFC3339
	EndTime     string   `form:"end_time"`   // RFC3339
	Keyword     string   `form:"keyword"`    // 在 message 中搜索
	Page        int      `form:"page" binding:"min=1"`
	PageSize    int      `form:"page_size" binding:"min=1,max=1000"`
	SortBy      string   `form:"sort_by"`      // time_created, event_id
	SortOrder   string   `form:"sort_order"`   // asc, desc
}

// WindowsLogSearchRequest Windows 日誌全文搜索請求
type WindowsLogSearchRequest struct {
	Query     string            `json:"query" binding:"required"`      // 搜索關鍵字
	Filters   WindowsLogFilters `json:"filters"`
	TimeRange string            `json:"time_range"` // 1h, 24h, 7d, 30d, custom
	StartTime string            `json:"start_time,omitempty"` // RFC3339, for custom range
	EndTime   string            `json:"end_time,omitempty"`   // RFC3339, for custom range
	Page      int               `json:"page" binding:"min=1"`
	PageSize  int               `json:"page_size" binding:"min=1,max=1000"`
}

// WindowsLogFilters 日誌過濾器
type WindowsLogFilters struct {
	AgentIDs  []string `json:"agent_ids,omitempty"`
	LogTypes  []string `json:"log_types,omitempty"`
	Sources   []string `json:"sources,omitempty"`
	EventIDs  []int    `json:"event_ids,omitempty"`
	Levels    []string `json:"levels,omitempty"`
}

