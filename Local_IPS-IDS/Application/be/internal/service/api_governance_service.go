package service

import (
	"context"
	"time"

	"axiom-backend/internal/database"
)

// APIGovernanceService API 治理服務
type APIGovernanceService struct {
	db *database.Database
}

// APIHealthScore API 健康評分
type APIHealthScore struct {
	APIPath       string             `json:"api_path"`
	HealthScore   int                `json:"health_score"` // 0-100
	Metrics       APIMetrics         `json:"metrics"`
	Issues        []string           `json:"issues"`
	Recommendations []string         `json:"recommendations"`
	Trend         string             `json:"trend"` // improving, stable, degrading
	LastUpdated   time.Time          `json:"last_updated"`
}

// APIMetrics API 指標
type APIMetrics struct {
	TotalRequests     int64   `json:"total_requests"`
	ErrorRate         float64 `json:"error_rate"`
	AvgResponseTime   float64 `json:"avg_response_time_ms"`
	P95ResponseTime   float64 `json:"p95_response_time_ms"`
	P99ResponseTime   float64 `json:"p99_response_time_ms"`
	SuccessRate       float64 `json:"success_rate"`
	RequestsPerSecond float64 `json:"requests_per_second"`
}

// APIUsageAnalytics API 使用分析
type APIUsageAnalytics struct {
	TopEndpoints     []EndpointUsage    `json:"top_endpoints"`
	TopClients       []ClientUsage      `json:"top_clients"`
	UsageByHour      map[string]int     `json:"usage_by_hour"`
	ErrorsByEndpoint map[string]int     `json:"errors_by_endpoint"`
	TotalRequests    int64              `json:"total_requests"`
	TimeRange        string             `json:"time_range"`
}

// EndpointUsage 端點使用情況
type EndpointUsage struct {
	Path          string  `json:"path"`
	Method        string  `json:"method"`
	Count         int64   `json:"count"`
	AvgDuration   float64 `json:"avg_duration_ms"`
	ErrorRate     float64 `json:"error_rate"`
}

// ClientUsage 客戶端使用情況
type ClientUsage struct {
	ClientIP      string `json:"client_ip"`
	RequestCount  int64  `json:"request_count"`
	ErrorCount    int64  `json:"error_count"`
	UserAgent     string `json:"user_agent"`
}

// NewAPIGovernanceService 創建 API 治理服務
func NewAPIGovernanceService(db *database.Database) *APIGovernanceService {
	return &APIGovernanceService{
		db: db,
	}
}

// GetAPIHealth 獲取 API 健康評分
func (s *APIGovernanceService) GetAPIHealth(ctx context.Context, apiPath string) (*APIHealthScore, error) {
	// 從 api_logs 表統計數據
	var totalRequests int64
	var errorCount int64
	
	s.db.PG.Model(&struct{}{}).
		Table("api_logs").
		Where("path = ?", apiPath).
		Count(&totalRequests)
	
	s.db.PG.Model(&struct{}{}).
		Table("api_logs").
		Where("path = ? AND status >= 400", apiPath).
		Count(&errorCount)
	
	errorRate := 0.0
	if totalRequests > 0 {
		errorRate = float64(errorCount) / float64(totalRequests)
	}
	
	// 計算平均響應時間
	var avgDuration float64
	s.db.PG.Model(&struct{}{}).
		Table("api_logs").
		Where("path = ?", apiPath).
		Select("AVG(duration)").
		Scan(&avgDuration)
	
	avgDuration = avgDuration / 1000.0 // 轉換為毫秒
	
	// 計算健康評分
	healthScore := 100
	if errorRate > 0.05 {
		healthScore -= 30
	} else if errorRate > 0.01 {
		healthScore -= 10
	}
	
	if avgDuration > 1000 {
		healthScore -= 20
	} else if avgDuration > 500 {
		healthScore -= 10
	}
	
	issues := []string{}
	recommendations := []string{}
	
	if errorRate > 0.05 {
		issues = append(issues, "High error rate detected")
		recommendations = append(recommendations, "Review error logs and fix issues")
	}
	
	if avgDuration > 500 {
		issues = append(issues, "Slow response time")
		recommendations = append(recommendations, "Optimize database queries or add caching")
	}
	
	return &APIHealthScore{
		APIPath:     apiPath,
		HealthScore: healthScore,
		Metrics: APIMetrics{
			TotalRequests:   totalRequests,
			ErrorRate:       errorRate,
			AvgResponseTime: avgDuration,
		},
		Issues:          issues,
		Recommendations: recommendations,
		Trend:           "stable",
		LastUpdated:     time.Now(),
	}, nil
}

// GetUsageAnalytics 獲取使用分析
func (s *APIGovernanceService) GetUsageAnalytics(ctx context.Context, timeRange string) (*APIUsageAnalytics, error) {
	// 統計 Top Endpoints
	var topEndpoints []EndpointUsage
	s.db.PG.Model(&struct{}{}).
		Table("api_logs").
		Select("path, method, COUNT(*) as count, AVG(duration) as avg_duration").
		Group("path, method").
		Order("count DESC").
		Limit(10).
		Scan(&topEndpoints)
	
	// 統計 Top Clients
	var topClients []ClientUsage
	s.db.PG.Model(&struct{}{}).
		Table("api_logs").
		Select("client_ip, COUNT(*) as request_count, user_agent").
		Group("client_ip, user_agent").
		Order("request_count DESC").
		Limit(10).
		Scan(&topClients)
	
	// 總請求數
	var totalRequests int64
	s.db.PG.Model(&struct{}{}).
		Table("api_logs").
		Count(&totalRequests)
	
	return &APIUsageAnalytics{
		TopEndpoints:  topEndpoints,
		TopClients:    topClients,
		UsageByHour:   map[string]int{},
		TotalRequests: totalRequests,
		TimeRange:     timeRange,
	}, nil
}


