package model

import (
	"time"
)

// ConfigHistory 配置歷史表
type ConfigHistory struct {
	ID         uint   `gorm:"primaryKey"`
	ServiceID  uint   `gorm:"index;not null"`
	Service    Service `gorm:"foreignKey:ServiceID;constraint:OnDelete:CASCADE"`
	ConfigType string  `gorm:"size:50;not null"` // nginx, agent, prometheus, etc.
	Content    string  `gorm:"type:text;not null"`
	AppliedBy  string  `gorm:"size:100"` // 執行者
	AppliedAt  time.Time `gorm:"not null;index"`
	Status     string  `gorm:"size:20;not null"` // pending, applied, failed, rollback
	Error      string  `gorm:"type:text"` // 錯誤訊息（如果失敗）
	CreatedAt  time.Time
}

// TableName 指定表名
func (ConfigHistory) TableName() string {
	return "config_histories"
}

