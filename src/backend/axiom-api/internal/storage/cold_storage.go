package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ColdStorage Cold Storage (PostgreSQL) - 90天歷史
type ColdStorage struct {
	db *gorm.DB
}

// ColdLogEntry 冷儲存日誌條目
type ColdLogEntry struct {
	ID            int64                  `gorm:"primaryKey;autoIncrement"`
	Timestamp     time.Time              `gorm:"not null;index:idx_timestamp"`
	AgentID       string                 `gorm:"type:varchar(64);not null;index:idx_agent_timestamp"`
	AgentMode     string                 `gorm:"type:varchar(16);not null"`
	EventType     string                 `gorm:"type:varchar(64);not null;index:idx_event_type"`
	Source        string                 `gorm:"type:varchar(128)"`
	EventID       int                    `gorm:"index:idx_event_id"`
	Level         string                 `gorm:"type:varchar(32);index:idx_level"`
	Computer      string                 `gorm:"type:varchar(256);index:idx_computer"`
	Message       string                 `gorm:"type:text"`
	RawData       string                 `gorm:"type:jsonb"` // JSONB for PostgreSQL
	
	// 合規性欄位
	RetentionUntil time.Time `gorm:"index:idx_retention"`
	Archived       bool      `gorm:"default:false;index:idx_archived"`
	IntegrityHash  string    `gorm:"type:varchar(64);index:idx_integrity"`
	
	// 元數據
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (ColdLogEntry) TableName() string {
	return "event_logs"
}

// NewColdStorage 創建 Cold Storage
func NewColdStorage(db *gorm.DB) *ColdStorage {
	return &ColdStorage{
		db: db,
	}
}

// Write 寫入日誌（批量）
func (c *ColdStorage) Write(ctx context.Context, logs []LogEntry) error {
	coldLogs := make([]ColdLogEntry, len(logs))
	
	for i, log := range logs {
		// 計算完整性 Hash
		hash := c.calculateIntegrityHash(log)
		
		// 計算保留期限（默認 90 天）
		retentionUntil := log.Timestamp.Add(90 * 24 * time.Hour)
		
		coldLogs[i] = ColdLogEntry{
			Timestamp:      log.Timestamp,
			AgentID:        log.AgentID,
			AgentMode:      log.AgentMode,
			EventType:      log.EventType,
			Source:         log.Source,
			Level:          log.Level,
			Message:        log.Message,
			RawData:        c.marshalJSON(log.Data),
			RetentionUntil: retentionUntil,
			Archived:       false,
			IntegrityHash:  hash,
		}
	}
	
	// 批量插入
	return c.db.WithContext(ctx).CreateInBatches(coldLogs, 1000).Error
}

// Query 查詢日誌
func (c *ColdStorage) Query(ctx context.Context, filter QueryFilter) ([]ColdLogEntry, error) {
	query := c.db.WithContext(ctx)
	
	// 應用過濾條件
	if filter.AgentID != "" {
		query = query.Where("agent_id = ?", filter.AgentID)
	}
	if filter.EventType != "" {
		query = query.Where("event_type = ?", filter.EventType)
	}
	if filter.Level != "" {
		query = query.Where("level = ?", filter.Level)
	}
	if !filter.StartTime.IsZero() {
		query = query.Where("timestamp >= ?", filter.StartTime)
	}
	if !filter.EndTime.IsZero() {
		query = query.Where("timestamp <= ?", filter.EndTime)
	}
	if filter.MessageContains != "" {
		query = query.Where("message LIKE ?", "%"+filter.MessageContains+"%")
	}
	
	// 排序和分頁
	query = query.Order("timestamp DESC")
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	
	var logs []ColdLogEntry
	err := query.Find(&logs).Error
	return logs, err
}

// GetLogsOlderThan 獲取舊日誌（用於封存）
func (c *ColdStorage) GetLogsOlderThan(ctx context.Context, duration time.Duration) ([]ColdLogEntry, error) {
	cutoffTime := time.Now().Add(-duration)
	
	var logs []ColdLogEntry
	err := c.db.WithContext(ctx).
		Where("timestamp < ?", cutoffTime).
		Where("archived = ?", false).
		Find(&logs).Error
	
	return logs, err
}

// MarkAsArchived 標記為已封存
func (c *ColdStorage) MarkAsArchived(ctx context.Context, logIDs []int64) error {
	return c.db.WithContext(ctx).
		Model(&ColdLogEntry{}).
		Where("id IN ?", logIDs).
		Update("archived", true).Error
}

// VerifyIntegrity 驗證完整性
func (c *ColdStorage) VerifyIntegrity(ctx context.Context) ([]int64, error) {
	var tamperedLogs []int64
	
	// 分批讀取並驗證
	batchSize := 1000
	offset := 0
	
	for {
		var logs []ColdLogEntry
		err := c.db.WithContext(ctx).
			Limit(batchSize).
			Offset(offset).
			Find(&logs).Error
		
		if err != nil {
			return nil, err
		}
		
		if len(logs) == 0 {
			break
		}
		
		for _, log := range logs {
			expectedHash := c.calculateIntegrityHashFromDB(log)
			if log.IntegrityHash != expectedHash {
				tamperedLogs = append(tamperedLogs, log.ID)
			}
		}
		
		offset += batchSize
	}
	
	return tamperedLogs, nil
}

// CreatePartition 創建月度分區
func (c *ColdStorage) CreatePartition(ctx context.Context, date time.Time) error {
	partitionName := fmt.Sprintf("event_logs_%s", date.Format("2006_01"))
	startDate := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)
	
	sql := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s PARTITION OF event_logs
		FOR VALUES FROM ('%s') TO ('%s')
	`, partitionName, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	
	return c.db.WithContext(ctx).Exec(sql).Error
}

// GetStats 獲取統計信息
func (c *ColdStorage) GetStats(ctx context.Context) (map[string]interface{}, error) {
	var totalCount int64
	var archivedCount int64
	
	c.db.WithContext(ctx).Model(&ColdLogEntry{}).Count(&totalCount)
	c.db.WithContext(ctx).Model(&ColdLogEntry{}).Where("archived = ?", true).Count(&archivedCount)
	
	return map[string]interface{}{
		"total_logs":     totalCount,
		"archived_logs":  archivedCount,
		"active_logs":    totalCount - archivedCount,
		"retention":      "90 days",
		"partitioned":    true,
	}, nil
}

// calculateIntegrityHash 計算完整性 Hash
func (c *ColdStorage) calculateIntegrityHash(log LogEntry) string {
	data := fmt.Sprintf("%s%s%s%s",
		log.Timestamp.Format(time.RFC3339Nano),
		log.AgentID,
		log.EventType,
		log.Message,
	)
	
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// calculateIntegrityHashFromDB 從資料庫記錄計算 Hash
func (c *ColdStorage) calculateIntegrityHashFromDB(log ColdLogEntry) string {
	data := fmt.Sprintf("%s%s%s%s",
		log.Timestamp.Format(time.RFC3339Nano),
		log.AgentID,
		log.EventType,
		log.Message,
	)
	
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// marshalJSON 序列化 JSON
func (c *ColdStorage) marshalJSON(data map[string]interface{}) string {
	// 簡化實現，實際應使用 json.Marshal
	return "{}"
}

// QueryFilter 查詢過濾器
type QueryFilter struct {
	AgentID         string
	EventType       string
	Level           string
	StartTime       time.Time
	EndTime         time.Time
	MessageContains string
	Limit           int
	Offset          int
}

