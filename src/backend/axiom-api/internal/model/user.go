package model

import (
	"time"

	"gorm.io/datatypes"
)

// User 使用者表
type User struct {
	ID           uint           `gorm:"primaryKey"`
	Username     string         `gorm:"uniqueIndex;not null;size:100"`
	Email        string         `gorm:"uniqueIndex;size:255"`
	PasswordHash string         `gorm:"not null;size:255"` // bcrypt hash
	Role         string         `gorm:"size:50;not null;index;default:'viewer'"` // admin, operator, viewer
	Status       string         `gorm:"size:20;not null;default:'active'"` // active, disabled, locked
	APIKey       string         `gorm:"uniqueIndex;size:64"` // API Key for programmatic access
	Permissions  datatypes.JSON `gorm:"type:jsonb"` // 細粒度權限
	LastLoginAt  *time.Time
	LastLoginIP  string `gorm:"size:45"`
	LoginAttempts int `gorm:"default:0"` // 失敗登入次數
	LockedUntil  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"` // 軟刪除
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// IsAdmin 檢查是否為管理員
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsActive 檢查使用者是否啟用
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// IsLocked 檢查使用者是否被鎖定
func (u *User) IsLocked() bool {
	if u.Status == "locked" {
		return true
	}
	if u.LockedUntil != nil && u.LockedUntil.After(time.Now()) {
		return true
	}
	return false
}

