package compliance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"axiom-backend/internal/model"
	
	"gorm.io/gorm"
)

// GDPRService GDPR 合規服務
type GDPRService struct {
	db         *gorm.DB
	anonymizer *Anonymizer
}

// NewGDPRService 創建 GDPR 服務
func NewGDPRService(db *gorm.DB, anonymizer *Anonymizer) *GDPRService {
	return &GDPRService{
		db:         db,
		anonymizer: anonymizer,
	}
}

// CreateDeletionRequest 創建刪除請求
func (s *GDPRService) CreateDeletionRequest(ctx context.Context, subjectID, requestedBy, notes string) (*model.GDPRDeletionRequest, error) {
	request := &model.GDPRDeletionRequest{
		SubjectIdentifier: subjectID,
		RequestedBy:       requestedBy,
		RequestedAt:       time.Now(),
		Status:            "pending",
		Notes:             notes,
	}
	
	err := s.db.WithContext(ctx).Create(request).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create deletion request: %w", err)
	}
	
	return request, nil
}

// ApproveDeletionRequest 審批刪除請求
func (s *GDPRService) ApproveDeletionRequest(ctx context.Context, requestID, approvedBy string) error {
	now := time.Now()
	
	result := s.db.WithContext(ctx).
		Model(&model.GDPRDeletionRequest{}).
		Where("request_id = ? AND status = ?", requestID, "pending").
		Updates(map[string]interface{}{
			"status":      "approved",
			"approved_by": approvedBy,
			"approved_at": now,
		})
	
	if result.Error != nil {
		return fmt.Errorf("failed to approve request: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("request not found or already processed")
	}
	
	return nil
}

// ExecuteDeletion 執行刪除
func (s *GDPRService) ExecuteDeletion(ctx context.Context, requestID string) (*DeletionResult, error) {
	// 獲取請求
	var request model.GDPRDeletionRequest
	err := s.db.WithContext(ctx).
		Where("request_id = ? AND status = ?", requestID, "approved").
		First(&request).Error
	
	if err != nil {
		return nil, fmt.Errorf("request not found or not approved: %w", err)
	}
	
	// 執行刪除
	deletedCount := 0
	affectedTables := []string{}
	
	// 1. 刪除 event_logs 中的相關數據
	result := s.db.WithContext(ctx).
		Exec("DELETE FROM event_logs WHERE message LIKE ?", "%"+request.SubjectIdentifier+"%")
	
	if result.Error == nil {
		deletedCount += int(result.RowsAffected)
		if result.RowsAffected > 0 {
			affectedTables = append(affectedTables, "event_logs")
		}
	}
	
	// 2. 刪除其他表的數據（如需要）
	// 可以根據實際需求刪除更多表
	
	// 3. 生成驗證 Hash
	verificationHash := s.generateVerificationHash(requestID, deletedCount)
	
	// 4. 更新請求狀態
	now := time.Now()
	s.db.WithContext(ctx).
		Model(&request).
		Updates(map[string]interface{}{
			"status":            "completed",
			"completion_date":   now,
			"deleted_count":     deletedCount,
			"verification_hash": verificationHash,
		})
	
	return &DeletionResult{
		RequestID:        requestID,
		DeletedCount:     deletedCount,
		AffectedTables:   affectedTables,
		VerificationHash: verificationHash,
		CompletedAt:      now,
	}, nil
}

// VerifyDeletion 驗證刪除
func (s *GDPRService) VerifyDeletion(ctx context.Context, requestID string) (*DeletionVerification, error) {
	var request model.GDPRDeletionRequest
	err := s.db.WithContext(ctx).
		Where("request_id = ?", requestID).
		First(&request).Error
	
	if err != nil {
		return nil, fmt.Errorf("request not found: %w", err)
	}
	
	// 驗證數據是否真的被刪除
	var remainingCount int64
	s.db.WithContext(ctx).
		Model(&struct{}{}).
		Table("event_logs").
		Where("message LIKE ?", "%"+request.SubjectIdentifier+"%").
		Count(&remainingCount)
	
	verified := remainingCount == 0
	
	return &DeletionVerification{
		RequestID:       requestID,
		Status:          request.Status,
		DeletedCount:    request.DeletedCount,
		RemainingCount:  int(remainingCount),
		Verified:        verified,
		VerifiedAt:      time.Now(),
	}, nil
}

// ListDeletionRequests 列出刪除請求
func (s *GDPRService) ListDeletionRequests(ctx context.Context, status string) ([]model.GDPRDeletionRequest, error) {
	var requests []model.GDPRDeletionRequest
	
	query := s.db.WithContext(ctx)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Order("requested_at DESC").Find(&requests).Error
	return requests, err
}

// ExportData 匯出個人數據（GDPR 資料可攜性）
func (s *GDPRService) ExportData(ctx context.Context, subjectID string) (*DataExport, error) {
	// 查詢所有相關數據
	var logs []map[string]interface{}
	
	err := s.db.WithContext(ctx).
		Table("event_logs").
		Where("message LIKE ?", "%"+subjectID+"%").
		Limit(10000).
		Find(&logs).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to export data: %w", err)
	}
	
	return &DataExport{
		SubjectID:    subjectID,
		RecordCount:  len(logs),
		Data:         logs,
		ExportedAt:   time.Now(),
		Format:       "JSON",
		Regulation:   "GDPR",
	}, nil
}

// generateVerificationHash 生成驗證 Hash
func (s *GDPRService) generateVerificationHash(requestID string, deletedCount int) string {
	data := fmt.Sprintf("%s:%d:%s", requestID, deletedCount, time.Now().Format(time.RFC3339))
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// DeletionResult 刪除結果
type DeletionResult struct {
	RequestID        string    `json:"request_id"`
	DeletedCount     int       `json:"deleted_count"`
	AffectedTables   []string  `json:"affected_tables"`
	VerificationHash string    `json:"verification_hash"`
	CompletedAt      time.Time `json:"completed_at"`
}

// DeletionVerification 刪除驗證
type DeletionVerification struct {
	RequestID      string    `json:"request_id"`
	Status         string    `json:"status"`
	DeletedCount   int       `json:"deleted_count"`
	RemainingCount int       `json:"remaining_count"`
	Verified       bool      `json:"verified"`
	VerifiedAt     time.Time `json:"verified_at"`
}

// DataExport 數據匯出
type DataExport struct {
	SubjectID   string                   `json:"subject_id"`
	RecordCount int                      `json:"record_count"`
	Data        []map[string]interface{} `json:"data"`
	ExportedAt  time.Time                `json:"exported_at"`
	Format      string                   `json:"format"`
	Regulation  string                   `json:"regulation"`
}

