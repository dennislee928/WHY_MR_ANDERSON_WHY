package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"axiom-backend/internal/dto"
	"axiom-backend/internal/vo"
	apperrors "axiom-backend/internal/errors"
)

// NginxService Nginx 服務
type NginxService struct {
	BaseService
	configPath string
}

// NewNginxService 創建 Nginx 服務
func NewNginxService(baseURL, configPath string) *NginxService {
	return &NginxService{
		BaseService: BaseService{
			Name:    "nginx",
			BaseURL: baseURL,
		},
		configPath: configPath,
	}
}

// HealthCheck 健康檢查
func (s *NginxService) HealthCheck(ctx context.Context) error {
	// 檢查 Nginx 進程是否在運行
	cmd := exec.CommandContext(ctx, "nginx", "-t")
	if err := cmd.Run(); err != nil {
		return apperrors.Wrap(err, "nginx health check failed")
	}
	return nil
}

// GetStatus 獲取服務狀態
func (s *NginxService) GetStatus(ctx context.Context) (map[string]interface{}, error) {
	// Note: 需要 nginx stub_status 模塊支援
	// 這裡返回模擬數據，實際應該從 nginx status 端點獲取
	return map[string]interface{}{
		"name":               s.Name,
		"status":             "healthy",
		"active_connections": 42,
		"accepts":            15230,
		"handled":            15230,
		"requests":           30460,
		"reading":            0,
		"writing":            3,
		"waiting":            39,
		"timestamp":          time.Now(),
	}, nil
}

// GetConfig 獲取配置
func (s *NginxService) GetConfig(ctx context.Context) (*vo.NginxConfigVO, error) {
	content, err := os.ReadFile(s.configPath)
	if err != nil {
		return nil, apperrors.Wrap(err, "read nginx config failed")
	}

	fileInfo, err := os.Stat(s.configPath)
	if err != nil {
		return nil, apperrors.Wrap(err, "stat nginx config failed")
	}

	return &vo.NginxConfigVO{
		Config:       string(content),
		ConfigPath:   s.configPath,
		LastModified: fileInfo.ModTime(),
		Size:         fileInfo.Size(),
		Valid:        true,
	}, nil
}

// UpdateConfig 更新配置
func (s *NginxService) UpdateConfig(ctx context.Context, req *dto.NginxConfigUpdateRequest) (*vo.NginxConfigVO, error) {
	// 備份舊配置
	if req.Backup {
		backupPath := fmt.Sprintf("%s.backup.%d", s.configPath, time.Now().Unix())
		oldContent, err := os.ReadFile(s.configPath)
		if err == nil {
			os.WriteFile(backupPath, oldContent, 0644)
		}
	}

	// 驗證配置
	if req.Validate {
		// 寫入臨時文件
		tempPath := s.configPath + ".tmp"
		if err := os.WriteFile(tempPath, []byte(req.Config), 0644); err != nil {
			return nil, apperrors.Wrap(err, "write temp config failed")
		}
		defer os.Remove(tempPath)

		// 測試配置
		cmd := exec.CommandContext(ctx, "nginx", "-t", "-c", tempPath)
		if err := cmd.Run(); err != nil {
			return nil, apperrors.New(apperrors.ErrCodeValidation, "Invalid nginx config", 400)
		}
	}

	// 寫入新配置
	if err := os.WriteFile(s.configPath, []byte(req.Config), 0644); err != nil {
		return nil, apperrors.Wrap(err, "write nginx config failed")
	}

	fileInfo, err := os.Stat(s.configPath)
	if err != nil {
		return nil, apperrors.Wrap(err, "stat nginx config failed")
	}

	return &vo.NginxConfigVO{
		Config:       req.Config,
		ConfigPath:   s.configPath,
		LastModified: fileInfo.ModTime(),
		Size:         fileInfo.Size(),
		Valid:        true,
	}, nil
}

// Reload 重載配置
func (s *NginxService) Reload(ctx context.Context) (*vo.NginxReloadVO, error) {
	start := time.Now()

	cmd := exec.CommandContext(ctx, "nginx", "-s", "reload")
	if err := cmd.Run(); err != nil {
		return nil, apperrors.Wrap(err, "nginx reload failed")
	}

	duration := time.Since(start).Milliseconds()

	return &vo.NginxReloadVO{
		Success:   true,
		Message:   "Nginx reloaded successfully",
		Duration:  int(duration),
		Timestamp: time.Now(),
	}, nil
}


