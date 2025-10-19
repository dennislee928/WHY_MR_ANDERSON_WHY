package model

import (
	"time"

	"gorm.io/gorm"
)

// RetentionPolicy 保留策略
type RetentionPolicy struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	EventType       string    `gorm:"type:varchar(64);not null;index" json:"event_type"`
	AgentMode       string    `gorm:"type:varchar(16);index" json:"agent_mode"` // external, internal
	RetentionDays   int       `gorm:"not null" json:"retention_days"`
	LegalHold       bool      `gorm:"default:false" json:"legal_hold"`
	Regulation      string    `gorm:"type:varchar(32);not null" json:"regulation"` // GDPR, PCI-DSS, HIPAA, SOX, ISO27001
	AutoDelete      bool      `gorm:"default:true" json:"auto_delete"`
	ArchiveRequired bool      `gorm:"default:false" json:"archive_required"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// GDPRDeletionRequest GDPR 刪除請求
type GDPRDeletionRequest struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	RequestID          string    `gorm:"type:uuid;uniqueIndex;default:gen_random_uuid()" json:"request_id"`
	SubjectIdentifier  string    `gorm:"type:varchar(256);not null" json:"subject_identifier"` // Email, user ID
	RequestedBy        string    `gorm:"type:varchar(128)" json:"requested_by"`
	RequestedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"requested_at"`
	ApprovedBy         string    `gorm:"type:varchar(128)" json:"approved_by,omitempty"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty"`
	Status             string    `gorm:"type:varchar(32);default:'pending';index" json:"status"` // pending, approved, completed, rejected
	CompletionDate     *time.Time `json:"completion_date,omitempty"`
	VerificationHash   string    `gorm:"type:varchar(64)" json:"verification_hash,omitempty"`
	DeletedCount       int       `json:"deleted_count,omitempty"`
	Notes              string    `gorm:"type:text" json:"notes,omitempty"`
}

// AuditAccessLog 審計訪問日誌
type AuditAccessLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Timestamp     time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"timestamp"`
	UserID        string    `gorm:"type:varchar(64);not null;index" json:"user_id"`
	Action        string    `gorm:"type:varchar(64);not null;index" json:"action"` // query, export, delete, update
	ResourceType  string    `gorm:"type:varchar(64);index" json:"resource_type"`
	ResourceID    string    `gorm:"type:varchar(256)" json:"resource_id"`
	QueryText     string    `gorm:"type:text" json:"query_text,omitempty"`
	RecordCount   int       `json:"record_count,omitempty"`
	IPAddress     string    `gorm:"type:inet" json:"ip_address"`
	UserAgent     string    `gorm:"type:text" json:"user_agent,omitempty"`
	Justification string    `gorm:"type:text" json:"justification,omitempty"` // Required for GDPR
	ApprovedBy    string    `gorm:"type:varchar(64)" json:"approved_by,omitempty"`
	SessionID     string    `gorm:"type:varchar(128);index" json:"session_id"`
}

// PIIPattern PII 模式
type PIIPattern struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PatternType string    `gorm:"type:varchar(64);not null;uniqueIndex" json:"pattern_type"`
	Regex       string    `gorm:"type:varchar(512);not null" json:"regex"`
	Description string    `gorm:"type:text" json:"description"`
	Severity    string    `gorm:"type:varchar(20)" json:"severity"` // low, medium, high, critical
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PIIOccurrence PII 出現記錄
type PIIOccurrence struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	LogID      int64     `gorm:"index" json:"log_id"`
	PIIType    string    `gorm:"type:varchar(64);index" json:"pii_type"`
	FieldName  string    `gorm:"type:varchar(128)" json:"field_name"`
	DetectedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"detected_at"`
	Masked     bool      `gorm:"default:false" json:"masked"`
	Hash       string    `gorm:"type:varchar(64)" json:"hash"`
}

// AutoMigrate 自動遷移合規性表
func AutoMigrateComplianceTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&RetentionPolicy{},
		&GDPRDeletionRequest{},
		&AuditAccessLog{},
		&PIIPattern{},
		&PIIOccurrence{},
	)
}

// SeedDefaultPolicies 初始化默認策略
func SeedDefaultPolicies(db *gorm.DB) error {
	policies := []RetentionPolicy{
		{
			EventType:       "windows_event_log",
			AgentMode:       "external",
			RetentionDays:   90,
			LegalHold:       false,
			Regulation:      "PCI-DSS",
			AutoDelete:      true,
			ArchiveRequired: true,
		},
		{
			EventType:       "security_alert",
			AgentMode:       "external",
			RetentionDays:   365,
			LegalHold:       false,
			Regulation:      "GDPR",
			AutoDelete:      true,
			ArchiveRequired: true,
		},
		{
			EventType:       "access_log",
			AgentMode:       "internal",
			RetentionDays:   180,
			LegalHold:       false,
			Regulation:      "SOX",
			AutoDelete:      true,
			ArchiveRequired: false,
		},
		{
			EventType:       "compliance_scan",
			AgentMode:       "external",
			RetentionDays:   2555, // 7 years
			LegalHold:       false,
			Regulation:      "HIPAA",
			AutoDelete:      false,
			ArchiveRequired: true,
		},
	}
	
	for _, policy := range policies {
		// 檢查是否已存在
		var existing RetentionPolicy
		result := db.Where("event_type = ? AND agent_mode = ?", policy.EventType, policy.AgentMode).First(&existing)
		
		if result.Error == gorm.ErrRecordNotFound {
			// 創建新策略
			if err := db.Create(&policy).Error; err != nil {
				return err
			}
		}
	}
	
	return nil
}

// SeedDefaultPIIPatterns 初始化默認 PII 模式
func SeedDefaultPIIPatterns(db *gorm.DB) error {
	patterns := []PIIPattern{
		{
			PatternType: "email",
			Regex:       `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`,
			Description: "Email addresses",
			Severity:    "medium",
			Enabled:     true,
		},
		{
			PatternType: "credit_card",
			Regex:       `\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`,
			Description: "Credit card numbers",
			Severity:    "critical",
			Enabled:     true,
		},
		{
			PatternType: "ssn",
			Regex:       `\b\d{3}-\d{2}-\d{4}\b`,
			Description: "Social Security Numbers",
			Severity:    "critical",
			Enabled:     true,
		},
		{
			PatternType: "ip_address",
			Regex:       `\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`,
			Description: "IP addresses",
			Severity:    "low",
			Enabled:     true,
		},
		{
			PatternType: "phone",
			Regex:       `\b\d{3}[-.]?\d{3}[-.]?\d{4}\b`,
			Description: "Phone numbers",
			Severity:    "medium",
			Enabled:     true,
		},
	}
	
	for _, pattern := range patterns {
		var existing PIIPattern
		result := db.Where("pattern_type = ?", pattern.PatternType).First(&existing)
		
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&pattern).Error; err != nil {
				return err
			}
		}
	}
	
	return nil
}

