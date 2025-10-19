package dto

// ServiceHealthCheckRequest 服務健康檢查請求
type ServiceHealthCheckRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
}

// ServiceRestartRequest 服務重啟請求
type ServiceRestartRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Force       bool   `json:"force"`        // 強制重啟
	Timeout     int    `json:"timeout"`      // 超時時間（秒）
}

// ServiceConfigUpdateRequest 服務配置更新請求
type ServiceConfigUpdateRequest struct {
	ServiceName string                 `json:"service_name" binding:"required"`
	ConfigType  string                 `json:"config_type" binding:"required"` // nginx, agent, etc.
	Content     string                 `json:"content" binding:"required"`
	Validate    bool                   `json:"validate"` // 是否先驗證配置
	AppliedBy   string                 `json:"applied_by"`
}

