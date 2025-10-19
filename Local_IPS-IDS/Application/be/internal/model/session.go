package model

import (
	"time"

	"gorm.io/datatypes"
)

// Session 會話表（也會存在 Redis 中）
type Session struct {
	ID        uint           `gorm:"primaryKey"`
	SessionID string         `gorm:"uniqueIndex;not null;size:128"` // 會話 ID
	UserID    uint           `gorm:"index;not null"`
	User      User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Token     string         `gorm:"uniqueIndex;not null;size:512"` // JWT Token
	ClientIP  string         `gorm:"size:45"`
	UserAgent string         `gorm:"size:500"`
	ExpiresAt time.Time      `gorm:"not null;index"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"` // 額外元數據
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName 指定表名
func (Session) TableName() string {
	return "sessions"
}

// IsExpired 檢查會話是否過期
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// TimeToExpire 返回距離過期的時間
func (s *Session) TimeToExpire() time.Duration {
	return time.Until(s.ExpiresAt)
}

