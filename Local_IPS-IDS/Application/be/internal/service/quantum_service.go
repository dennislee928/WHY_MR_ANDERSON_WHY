package service

import (
	"context"
	"fmt"
	"time"

	"axiom-backend/internal/client"
	"axiom-backend/internal/database"
	"axiom-backend/internal/dto"
	"axiom-backend/internal/model"
	"axiom-backend/internal/vo"
	apperrors "axiom-backend/internal/errors"
	
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// QuantumService 量子服務
type QuantumService struct {
	BaseService
	httpClient *client.HTTPClient
	db         *database.Database
}

// NewQuantumService 創建量子服務
func NewQuantumService(baseURL string, db *database.Database) *QuantumService {
	return &QuantumService{
		BaseService: BaseService{
			Name:    "cyber-ai-quantum",
			BaseURL: baseURL,
		},
		httpClient: client.NewHTTPClient(&client.Config{
			BaseURL: baseURL,
			Timeout: 120 * time.Second, // 量子計算可能需要更長時間
		}),
		db: db,
	}
}

// HealthCheck 健康檢查
func (s *QuantumService) HealthCheck(ctx context.Context) error {
	_, err := s.httpClient.Get(ctx, "/health")
	if err != nil {
		return apperrors.Wrap(err, "quantum service health check failed")
	}
	return nil
}

// GetStatus 獲取服務狀態
func (s *QuantumService) GetStatus(ctx context.Context) (map[string]interface{}, error) {
	var status map[string]interface{}
	err := s.httpClient.GetJSON(ctx, "/health", &status)
	if err != nil {
		return nil, apperrors.Wrap(err, "get quantum service status failed")
	}
	
	status["name"] = s.Name
	status["base_url"] = s.BaseURL
	status["timestamp"] = time.Now()
	
	return status, nil
}

// GenerateQKD 生成量子密鑰
func (s *QuantumService) GenerateQKD(ctx context.Context, req *dto.QuantumQKDRequest) (*vo.QuantumQKDVO, error) {
	// 創建作業記錄
	jobID := fmt.Sprintf("qkd-%s", uuid.New().String()[:8])
	job := &model.QuantumJob{
		JobID:       jobID,
		Type:        "qkd",
		Status:      "pending",
		Backend:     req.Backend,
		Shots:       req.Shots,
		SubmittedAt: time.Now(),
		InputData:   datatypes.JSON(fmt.Sprintf(`{"key_length": %d}`, req.KeyLength)),
	}
	
	if err := s.db.PG.Create(job).Error; err != nil {
		return nil, apperrors.Wrap(err, "create quantum job failed")
	}

	// 調用 cyber-ai-quantum 服務
	requestBody := map[string]interface{}{
		"key_length": req.KeyLength,
	}
	if req.Backend != "" {
		requestBody["backend"] = req.Backend
	}

	var result map[string]interface{}
	err := s.httpClient.PostJSON(ctx, "/api/v1/quantum/qkd/generate", requestBody, &result)
	if err != nil {
		// 更新作業狀態為失敗
		s.db.PG.Model(job).Updates(map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
		})
		return nil, apperrors.Wrap(err, "quantum qkd generation failed")
	}

	// 更新作業狀態為完成
	now := time.Now()
	s.db.PG.Model(job).Updates(map[string]interface{}{
		"status":       "completed",
		"completed_at": &now,
		"result":       datatypes.JSON(fmt.Sprintf("%v", result)),
	})

	return &vo.QuantumQKDVO{
		JobID:       jobID,
		Key:         fmt.Sprintf("%v", result["key"]),
		KeyLength:   req.KeyLength,
		Status:      "completed",
		SubmittedAt: job.SubmittedAt,
	}, nil
}

// ClassifyQSVM 執行 QSVM 分類
func (s *QuantumService) ClassifyQSVM(ctx context.Context, req *dto.QuantumQSVMRequest) (*vo.QuantumClassifyVO, error) {
	jobID := fmt.Sprintf("qsvm-%s", uuid.New().String()[:8])
	job := &model.QuantumJob{
		JobID:       jobID,
		Type:        "qsvm",
		Status:      "pending",
		Backend:     req.Backend,
		Shots:       req.Shots,
		SubmittedAt: time.Now(),
	}
	
	if err := s.db.PG.Create(job).Error; err != nil {
		return nil, apperrors.Wrap(err, "create quantum job failed")
	}

	requestBody := map[string]interface{}{
		"features": req.Features,
	}

	var result map[string]interface{}
	err := s.httpClient.PostJSON(ctx, "/api/v1/quantum/qsvm/classify", requestBody, &result)
	if err != nil {
		s.db.PG.Model(job).Updates(map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
		})
		return nil, apperrors.Wrap(err, "qsvm classification failed")
	}

	now := time.Now()
	s.db.PG.Model(job).Updates(map[string]interface{}{
		"status":       "completed",
		"completed_at": &now,
		"result":       datatypes.JSON(fmt.Sprintf("%v", result)),
	})

	return &vo.QuantumClassifyVO{
		JobID:       jobID,
		Prediction:  int(result["prediction"].(float64)),
		Probability: result["probability"].(float64),
		Status:      "completed",
		SubmittedAt: job.SubmittedAt,
	}, nil
}

// PredictZeroTrust 執行 Zero Trust 預測
func (s *QuantumService) PredictZeroTrust(ctx context.Context, req *dto.ZeroTrustPredictRequest) (*vo.ZeroTrustPredictVO, error) {
	requestBody := map[string]interface{}{
		"user_id":    req.UserID,
		"ip_address": req.IPAddress,
	}
	if req.DeviceID != "" {
		requestBody["device_id"] = req.DeviceID
	}
	if req.Features != nil {
		requestBody["features"] = req.Features
	}

	var result map[string]interface{}
	err := s.httpClient.PostJSON(ctx, "/api/v1/zerotrust/predict", requestBody, &result)
	if err != nil {
		return nil, apperrors.Wrap(err, "zero trust prediction failed")
	}

	return &vo.ZeroTrustPredictVO{
		TrustScore:  result["trust_score"].(float64),
		RiskLevel:   result["risk_level"].(string),
		Decision:    result["decision"].(string),
		Confidence:  result["confidence"].(float64),
		UsedQuantum: req.UseQuantum,
		Timestamp:   time.Now(),
	}, nil
}

// GetQuantumJob 獲取量子作業
func (s *QuantumService) GetQuantumJob(ctx context.Context, jobID string) (*vo.QuantumJobVO, error) {
	var job model.QuantumJob
	if err := s.db.PG.Where("job_id = ?", jobID).First(&job).Error; err != nil {
		return nil, apperrors.New(apperrors.ErrCodeNotFound, "Quantum job not found", 404)
	}

	return &vo.QuantumJobVO{
		JobID:       job.JobID,
		Type:        job.Type,
		Status:      job.Status,
		Backend:     job.Backend,
		BackendName: job.BackendName,
		Qubits:      job.Qubits,
		Depth:       job.Depth,
		Shots:       job.Shots,
		SubmittedAt: job.SubmittedAt,
		CompletedAt: job.CompletedAt,
		Error:       job.Error,
	}, nil
}

// ListQuantumJobs 列出量子作業
func (s *QuantumService) ListQuantumJobs(ctx context.Context, req *dto.QuantumJobQueryRequest) (*vo.QuantumJobsListVO, error) {
	query := s.db.PG.Model(&model.QuantumJob{})

	// 應用過濾器
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Backend != "" {
		query = query.Where("backend = ?", req.Backend)
	}

	// 計算總數
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, apperrors.Wrap(err, "count quantum jobs failed")
	}

	// 分頁
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize
	var jobs []model.QuantumJob
	if err := query.Offset(offset).Limit(req.PageSize).Order("submitted_at DESC").Find(&jobs).Error; err != nil {
		return nil, apperrors.Wrap(err, "list quantum jobs failed")
	}

	// 轉換為 VO
	jobVOs := make([]vo.QuantumJobVO, len(jobs))
	for i, job := range jobs {
		jobVOs[i] = vo.QuantumJobVO{
			JobID:       job.JobID,
			Type:        job.Type,
			Status:      job.Status,
			Backend:     job.Backend,
			Qubits:      job.Qubits,
			Shots:       job.Shots,
			SubmittedAt: job.SubmittedAt,
			CompletedAt: job.CompletedAt,
		}
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		totalPages++
	}

	return &vo.QuantumJobsListVO{
		Jobs:       jobVOs,
		Total:      int(total),
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetQuantumStats 獲取量子統計
func (s *QuantumService) GetQuantumStats(ctx context.Context) (*vo.QuantumStatsVO, error) {
	var stats struct {
		Total     int64
		Completed int64
		Failed    int64
		Running   int64
	}

	// 總作業數
	s.db.PG.Model(&model.QuantumJob{}).Count(&stats.Total)
	
	// 完成的作業
	s.db.PG.Model(&model.QuantumJob{}).Where("status = ?", "completed").Count(&stats.Completed)
	
	// 失敗的作業
	s.db.PG.Model(&model.QuantumJob{}).Where("status = ?", "failed").Count(&stats.Failed)
	
	// 運行中的作業
	s.db.PG.Model(&model.QuantumJob{}).Where("status IN ?", []string{"pending", "running"}).Count(&stats.Running)

	// 按類型統計
	var jobsByType []struct {
		Type  string
		Count int64
	}
	s.db.PG.Model(&model.QuantumJob{}).
		Select("type, count(*) as count").
		Group("type").
		Scan(&jobsByType)

	jobsByTypeMap := make(map[string]int)
	for _, item := range jobsByType {
		jobsByTypeMap[item.Type] = int(item.Count)
	}

	// 計算成功率
	successRate := 0.0
	if stats.Total > 0 {
		successRate = float64(stats.Completed) / float64(stats.Total)
	}

	return &vo.QuantumStatsVO{
		TotalJobs:     int(stats.Total),
		CompletedJobs: int(stats.Completed),
		FailedJobs:    int(stats.Failed),
		RunningJobs:   int(stats.Running),
		JobsByType:    jobsByTypeMap,
		SuccessRate:   successRate,
	}, nil
}


