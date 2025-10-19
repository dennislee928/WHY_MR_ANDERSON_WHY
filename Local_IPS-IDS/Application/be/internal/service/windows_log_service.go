package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"axiom-backend/internal/database"
	"axiom-backend/internal/dto"
	"axiom-backend/internal/model"
	"axiom-backend/internal/vo"
	
	"gorm.io/datatypes"
)

// WindowsLogService Windows 日誌服務
type WindowsLogService struct {
	db *database.Database
}

// NewWindowsLogService 創建 Windows 日誌服務
func NewWindowsLogService(db *database.Database) *WindowsLogService {
	return &WindowsLogService{
		db: db,
	}
}

// BatchReceive 批量接收日誌
func (s *WindowsLogService) BatchReceive(ctx context.Context, req *dto.WindowsLogBatchRequest) (*vo.WindowsLogBatchVO, error) {
	receivedCount := len(req.Logs)
	savedCount := 0
	failedCount := 0
	var errors []string

	// 批量插入日誌
	for _, logItem := range req.Logs {
		log := &model.WindowsLog{
			AgentID:     req.AgentID,
			LogType:     logItem.LogType,
			Source:      logItem.Source,
			EventID:     logItem.EventID,
			Level:       logItem.Level,
			Message:     logItem.Message,
			TimeCreated: logItem.TimeCreated,
			ReceivedAt:  time.Now(),
			Computer:    req.Computer,
			UserID:      logItem.UserID,
			ProcessID:   logItem.ProcessID,
			ThreadID:    logItem.ThreadID,
			Keywords:    logItem.Keywords,
		}

		// 轉換 metadata 為 JSONB
		if logItem.Metadata != nil {
			metadataBytes, _ := json.Marshal(logItem.Metadata)
			log.Metadata = datatypes.JSON(metadataBytes)
		}

		if err := s.db.PG.Create(log).Error; err != nil {
			failedCount++
			errors = append(errors, err.Error())
		} else {
			savedCount++
		}
	}

	return &vo.WindowsLogBatchVO{
		ReceivedCount: receivedCount,
		SavedCount:    savedCount,
		FailedCount:   failedCount,
		Errors:        errors,
		Timestamp:     time.Now(),
		Message:       "Logs processed successfully",
	}, nil
}

// Query 查詢日誌
func (s *WindowsLogService) Query(ctx context.Context, req *dto.WindowsLogQueryRequest) (*vo.WindowsLogsListVO, error) {
	query := s.db.PG.Model(&model.WindowsLog{})

	// 應用過濾器
	if req.AgentID != "" {
		query = query.Where("agent_id = ?", req.AgentID)
	}
	if req.LogType != "" {
		query = query.Where("log_type = ?", req.LogType)
	}
	if req.Source != "" {
		query = query.Where("source = ?", req.Source)
	}
	if req.EventID > 0 {
		query = query.Where("event_id = ?", req.EventID)
	}
	if req.Level != "" {
		query = query.Where("level = ?", req.Level)
	}
	if req.Keyword != "" {
		query = query.Where("message ILIKE ?", "%"+req.Keyword+"%")
	}

	// 時間範圍
	if req.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, req.StartTime)
		if err == nil {
			query = query.Where("time_created >= ?", startTime)
		}
	}
	if req.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, req.EndTime)
		if err == nil {
			query = query.Where("time_created <= ?", endTime)
		}
	}

	// 計算總數
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("count windows logs failed: %w", err)
	}

	// 分頁
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 50
	}

	offset := (req.Page - 1) * req.PageSize
	
	// 排序
	orderBy := "time_created DESC"
	if req.SortBy != "" {
		order := "DESC"
		if req.SortOrder == "asc" {
			order = "ASC"
		}
		orderBy = fmt.Sprintf("%s %s", req.SortBy, order)
	}

	var logs []model.WindowsLog
	if err := query.Offset(offset).Limit(req.PageSize).Order(orderBy).Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("query windows logs failed: %w", err)
	}

	// 轉換為 VO
	logVOs := make([]vo.WindowsLogVO, len(logs))
	for i, log := range logs {
		logVOs[i] = vo.WindowsLogVO{
			ID:          log.ID,
			AgentID:     log.AgentID,
			LogType:     log.LogType,
			Source:      log.Source,
			EventID:     log.EventID,
			Level:       log.Level,
			Message:     log.Message,
			TimeCreated: log.TimeCreated,
			ReceivedAt:  log.ReceivedAt,
			Computer:    log.Computer,
			UserID:      log.UserID,
		}
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		totalPages++
	}

	return &vo.WindowsLogsListVO{
		Logs:       logVOs,
		Total:      int(total),
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
		Timestamp:  time.Now(),
	}, nil
}

// GetStats 獲取統計信息
func (s *WindowsLogService) GetStats(ctx context.Context, timeRange string) (*vo.WindowsLogStatsVO, error) {
	query := s.db.PG.Model(&model.WindowsLog{})

	// 應用時間範圍
	if timeRange == "24h" {
		query = query.Where("time_created >= ?", time.Now().Add(-24*time.Hour))
	}

	// 總日誌數
	var totalLogs int64
	query.Count(&totalLogs)

	// 按類型統計
	var logsByType []struct {
		LogType string
		Count   int64
	}
	s.db.PG.Model(&model.WindowsLog{}).
		Select("log_type, count(*) as count").
		Group("log_type").
		Scan(&logsByType)

	logsByTypeMap := make(map[string]int)
	for _, item := range logsByType {
		logsByTypeMap[item.LogType] = int(item.Count)
	}

	// 按級別統計
	var logsByLevel []struct {
		Level string
		Count int64
	}
	s.db.PG.Model(&model.WindowsLog{}).
		Select("level, count(*) as count").
		Group("level").
		Scan(&logsByLevel)

	logsByLevelMap := make(map[string]int)
	for _, item := range logsByLevel {
		logsByLevelMap[item.Level] = int(item.Count)
	}

	return &vo.WindowsLogStatsVO{
		TotalLogs:     int(totalLogs),
		LogsByType:    logsByTypeMap,
		LogsByLevel:   logsByLevelMap,
		CriticalCount: logsByLevelMap["Critical"],
		ErrorCount:    logsByLevelMap["Error"],
		WarningCount:  logsByLevelMap["Warning"],
		TimeRange:     timeRange,
	}, nil
}


