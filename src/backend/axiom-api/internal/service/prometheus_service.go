package service

import (
	"context"
	"fmt"
	"time"

	"axiom-backend/internal/client"
	"axiom-backend/internal/dto"
	"axiom-backend/internal/vo"
	apperrors "axiom-backend/internal/errors"
)

// PrometheusService Prometheus 服務
type PrometheusService struct {
	BaseService
	httpClient *client.HTTPClient
}

// NewPrometheusService 創建 Prometheus 服務
func NewPrometheusService(baseURL string) *PrometheusService {
	return &PrometheusService{
		BaseService: BaseService{
			Name:    "prometheus",
			BaseURL: baseURL,
		},
		httpClient: client.NewHTTPClient(&client.Config{
			BaseURL: baseURL,
			Timeout: 30 * time.Second,
		}),
	}
}

// HealthCheck 健康檢查
func (s *PrometheusService) HealthCheck(ctx context.Context) error {
	_, err := s.httpClient.Get(ctx, "/-/healthy")
	if err != nil {
		return apperrors.Wrap(err, "prometheus health check failed")
	}
	return nil
}

// GetStatus 獲取服務狀態
func (s *PrometheusService) GetStatus(ctx context.Context) (map[string]interface{}, error) {
	// 獲取基本狀態
	var statusData map[string]interface{}
	err := s.httpClient.GetJSON(ctx, "/api/v1/status/config", &statusData)
	if err != nil {
		return nil, apperrors.Wrap(err, "get prometheus status failed")
	}

	return map[string]interface{}{
		"name":       s.Name,
		"status":     "healthy",
		"base_url":   s.BaseURL,
		"timestamp":  time.Now(),
		"config":     statusData,
	}, nil
}

// Query 執行 PromQL 查詢
func (s *PrometheusService) Query(ctx context.Context, req *dto.PrometheusQueryRequest) (*vo.PrometheusQueryVO, error) {
	params := map[string]interface{}{
		"query": req.Query,
	}
	if req.Time != "" {
		params["time"] = req.Time
	}

	var result vo.PrometheusQueryVO
	err := s.httpClient.GetJSON(ctx, fmt.Sprintf("/api/v1/query?query=%s", req.Query), &result)
	if err != nil {
		return nil, apperrors.Wrap(err, "prometheus query failed")
	}

	result.Timestamp = time.Now()
	return &result, nil
}

// QueryRange 執行範圍查詢
func (s *PrometheusService) QueryRange(ctx context.Context, req *dto.PrometheusQueryRangeRequest) (*vo.PrometheusQueryVO, error) {
	path := fmt.Sprintf("/api/v1/query_range?query=%s&start=%s&end=%s", req.Query, req.Start, req.End)
	if req.Step != "" {
		path += fmt.Sprintf("&step=%s", req.Step)
	}

	var result vo.PrometheusQueryVO
	err := s.httpClient.GetJSON(ctx, path, &result)
	if err != nil {
		return nil, apperrors.Wrap(err, "prometheus query range failed")
	}

	result.Timestamp = time.Now()
	return &result, nil
}

// GetAlertRules 獲取告警規則
func (s *PrometheusService) GetAlertRules(ctx context.Context) (*vo.PrometheusAlertRulesVO, error) {
	var result vo.PrometheusAlertRulesVO
	err := s.httpClient.GetJSON(ctx, "/api/v1/rules", &result)
	if err != nil {
		return nil, apperrors.Wrap(err, "get alert rules failed")
	}
	return &result, nil
}

// CreateAlertRule 創建告警規則
func (s *PrometheusService) CreateAlertRule(ctx context.Context, req *dto.PrometheusAlertRuleRequest) error {
	// Note: Prometheus 不支援通過 API 直接創建規則
	// 需要通過配置文件或使用 Prometheus Operator
	return apperrors.New(apperrors.ErrCodeBadRequest, "creating alert rules requires configuration file updates", 400)
}

// GetTargets 獲取抓取目標
func (s *PrometheusService) GetTargets(ctx context.Context) (*vo.PrometheusTargetsVO, error) {
	var result vo.PrometheusTargetsVO
	err := s.httpClient.GetJSON(ctx, "/api/v1/targets", &result)
	if err != nil {
		return nil, apperrors.Wrap(err, "get targets failed")
	}
	return &result, nil
}


