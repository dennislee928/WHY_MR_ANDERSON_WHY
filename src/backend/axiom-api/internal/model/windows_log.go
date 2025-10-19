package model

import (
	"time"

	"gorm.io/datatypes"
)

// WindowsLog Windows 日誌表
type WindowsLog struct {
	ID          uint           `gorm:"primaryKey"`
	AgentID     string         `gorm:"index;not null;size:100"` // Agent 識別碼
	LogType     string         `gorm:"size:50;not null;index"`  // System, Security, Application, Setup
	Source      string         `gorm:"size:255;index"`          // 日誌來源
	EventID     int            `gorm:"index"`                   // Windows Event ID
	Level       string         `gorm:"size:20;index"`           // Information, Warning, Error, Critical
	Message     string         `gorm:"type:text"`               // 日誌訊息
	TimeCreated time.Time      `gorm:"not null;index"`          // Windows 日誌產生時間
	ReceivedAt  time.Time      `gorm:"not null;index"`          // Axiom 接收時間
	Computer    string         `gorm:"size:255"`                // 電腦名稱
	UserID      string         `gorm:"size:100"`                // 使用者 SID
	ProcessID   int                                             // Process ID
	ThreadID    int                                             // Thread ID
	Keywords    string         `gorm:"size:255"`                // 關鍵字
	Metadata    datatypes.JSON `gorm:"type:jsonb"`              // 額外元數據（完整 XML）
	CreatedAt   time.Time
}

// TableName 指定表名
func (WindowsLog) TableName() string {
	return "windows_logs"
}

// IsCritical 檢查是否為關鍵日誌
func (w *WindowsLog) IsCritical() bool {
	return w.Level == "Error" || w.Level == "Critical"
}

// IsSecurityLog 檢查是否為安全日誌
func (w *WindowsLog) IsSecurityLog() bool {
	return w.LogType == "Security"
}

