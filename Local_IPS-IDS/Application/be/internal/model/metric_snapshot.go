package model

import (
	"time"

	"gorm.io/datatypes"
)

// MetricSnapshot 指標快照表（用於歷史數據分析）
type MetricSnapshot struct {
	ID          uint           `gorm:"primaryKey"`
	ServiceID   uint           `gorm:"index;not null"`
	Service     Service        `gorm:"foreignKey:ServiceID;constraint:OnDelete:CASCADE"`
	MetricName  string         `gorm:"size:255;not null;index"` // CPU, Memory, Disk, RequestRate, etc.
	MetricType  string         `gorm:"size:50"`                 // gauge, counter, histogram, summary
	Value       float64        `gorm:"not null"`
	Unit        string         `gorm:"size:20"`                 // %, MB, requests/s, etc.
	Labels      datatypes.JSON `gorm:"type:jsonb"`              // Prometheus 標籤
	Timestamp   time.Time      `gorm:"not null;index:idx_service_metric_time"`
	CreatedAt   time.Time
}

// TableName 指定表名
func (MetricSnapshot) TableName() string {
	return "metric_snapshots"
}

