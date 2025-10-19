package model

import (
	"time"
)

// APILog API 請求日誌表
type APILog struct {
	ID           uint   `gorm:"primaryKey"`
	Method       string `gorm:"size:10;not null;index"`     // GET, POST, PUT, DELETE, etc.
	Path         string `gorm:"size:500;not null;index"`    // API 路徑
	Status       int    `gorm:"not null;index"`             // HTTP 狀態碼
	Duration     int64  `gorm:"not null"`                   // 請求持續時間（微秒）
	ClientIP     string `gorm:"size:45;index"`              // 客戶端 IP（支援 IPv6）
	UserAgent    string `gorm:"size:500"`                   // User Agent
	RequestBody  string `gorm:"type:text"`                  // 請求 Body（可選，大請求不記錄）
	ResponseBody string `gorm:"type:text"`                  // 響應 Body（可選）
	Error        string `gorm:"type:text"`                  // 錯誤訊息
	UserID       string `gorm:"size:100;index"`             // 使用者 ID
	APIKey       string `gorm:"size:100;index"`             // API Key（前8字元）
	CreatedAt    time.Time `gorm:"not null;index:idx_created_at_status"` // 聯合索引
}

// TableName 指定表名
func (APILog) TableName() string {
	return "api_logs"
}

// IsSuccess 檢查請求是否成功
func (a *APILog) IsSuccess() bool {
	return a.Status >= 200 && a.Status < 300
}

// IsError 檢查是否為錯誤請求
func (a *APILog) IsError() bool {
	return a.Status >= 400
}

// DurationMs 返回持續時間（毫秒）
func (a *APILog) DurationMs() float64 {
	return float64(a.Duration) / 1000.0
}

