package model

import (
	"time"

	"gorm.io/datatypes"
)

// Alert 告警表
type Alert struct {
	ID          uint           `gorm:"primaryKey"`
	AlertName   string         `gorm:"index;not null;size:255"` // 告警名稱
	Fingerprint string         `gorm:"uniqueIndex;size:64"`     // 告警指紋（用於去重）
	Severity    string         `gorm:"size:20;not null;index"`  // critical, high, medium, low, info
	Source      string         `gorm:"size:100;not null;index"` // prometheus, loki, manual, quantum, etc.
	Category    string         `gorm:"size:50;index"`           // network, security, system, quantum, etc.
	Message     string         `gorm:"type:text;not null"`      // 告警訊息
	Description string         `gorm:"type:text"`               // 詳細描述
	Status      string         `gorm:"size:20;not null;index"`  // active, resolved, acknowledged, suppressed
	Priority    int            `gorm:"default:0"`               // 優先級（數字越大越重要）
	Count       int            `gorm:"default:1"`               // 重複次數
	Labels      datatypes.JSON `gorm:"type:jsonb"`              // 標籤（key-value pairs）
	Annotations datatypes.JSON `gorm:"type:jsonb"`              // 註解
	CreatedAt   time.Time      `gorm:"not null;index"`
	UpdatedAt   time.Time
	ResolvedAt  *time.Time `gorm:"index"`
	ResolvedBy  string     `gorm:"size:100"`
	AcknowledgedAt *time.Time
	AcknowledgedBy string `gorm:"size:100"`
	LastOccurredAt time.Time `gorm:"index"` // 最後一次發生時間
}

// TableName 指定表名
func (Alert) TableName() string {
	return "alerts"
}

// IsActive 檢查告警是否活躍
func (a *Alert) IsActive() bool {
	return a.Status == "active"
}

// IsResolved 檢查告警是否已解決
func (a *Alert) IsResolved() bool {
	return a.Status == "resolved"
}

// IsCritical 檢查是否為嚴重告警
func (a *Alert) IsCritical() bool {
	return a.Severity == "critical"
}

