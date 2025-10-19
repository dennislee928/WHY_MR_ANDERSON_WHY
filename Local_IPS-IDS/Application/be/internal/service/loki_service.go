package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"axiom-backend/internal/client"
	apperrors "axiom-backend/internal/errors"
)

// LokiService Loki 服務
type LokiService struct {
	BaseService
	httpClient *client.HTTPClient
}

// NewLokiService 創建 Loki 服務
func NewLokiService(baseURL string) *LokiService {
	return &LokiService{
		BaseService: BaseService{
			Name:    "loki",
			BaseURL: baseURL,
		},
		httpClient: client.NewHTTPClient(&client.Config{
			BaseURL: baseURL,
			Timeout: 60 * time.Second, // Loki 查詢可能較慢
		}),
	}
}

// HealthCheck 健康檢查
func (s *LokiService) HealthCheck(ctx context.Context) error {
	_, err := s.httpClient.Get(ctx, "/ready")
	if err != nil {
		return apperrors.Wrap(err, "loki health check failed")
	}
	return nil
}

// GetStatus 獲取服務狀態
func (s *LokiService) GetStatus(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{
		"name":      s.Name,
		"status":    "healthy",
		"base_url":  s.BaseURL,
		"timestamp": time.Now(),
	}, nil
}

// QueryLogs 查詢日誌
func (s *LokiService) QueryLogs(ctx context.Context, query string, limit int, start, end string) (interface{}, error) {
	params := url.Values{}
	params.Add("query", query)
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	if start != "" {
		params.Add("start", start)
	}
	if end != "" {
		params.Add("end", end)
	}

	path := fmt.Sprintf("/loki/api/v1/query_range?%s", params.Encode())

	var result interface{}
	err := s.httpClient.GetJSON(ctx, path, &result)
	if err != nil {
		return nil, apperrors.Wrap(err, "loki query failed")
	}

	return result, nil
}

// GetLabels 獲取標籤列表
func (s *LokiService) GetLabels(ctx context.Context) ([]string, error) {
	var result struct {
		Data []string `json:"data"`
	}

	err := s.httpClient.GetJSON(ctx, "/loki/api/v1/labels", &result)
	if err != nil {
		return nil, apperrors.Wrap(err, "get labels failed")
	}

	return result.Data, nil
}

// GetLabelValues 獲取標籤值
func (s *LokiService) GetLabelValues(ctx context.Context, label string) ([]string, error) {
	path := fmt.Sprintf("/loki/api/v1/label/%s/values", label)

	var result struct {
		Data []string `json:"data"`
	}

	err := s.httpClient.GetJSON(ctx, path, &result)
	if err != nil {
		return nil, apperrors.Wrap(err, "get label values failed")
	}

	return result.Data, nil
}


