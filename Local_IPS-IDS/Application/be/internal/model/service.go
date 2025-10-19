package model

import (
	"time"

	"gorm.io/datatypes"
)

// Service 服務狀態表
type Service struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"uniqueIndex;not null;size:100"`
	Type      string         `gorm:"size:50"` // prometheus, grafana, loki, etc.
	Status    string         `gorm:"size:20"` // healthy, unhealthy, unknown
	URL       string         `gorm:"size:255"`
	Version   string         `gorm:"size:50"`
	LastCheck time.Time      `gorm:"index"`
	Config    datatypes.JSON `gorm:"type:jsonb"`
	Metrics   datatypes.JSON `gorm:"type:jsonb"` // 即時指標快照
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 指定表名
func (Service) TableName() string {
	return "services"
}

