package service

import (
	"context"
)

// ServiceManager 服務管理器接口
type ServiceManager interface {
	// GetName 獲取服務名稱
	GetName() string
	
	// HealthCheck 健康檢查
	HealthCheck(ctx context.Context) error
	
	// GetStatus 獲取服務狀態
	GetStatus(ctx context.Context) (map[string]interface{}, error)
}

// BaseService 基礎服務結構
type BaseService struct {
	Name    string
	BaseURL string
}

// GetName 實現 ServiceManager 接口
func (s *BaseService) GetName() string {
	return s.Name
}

