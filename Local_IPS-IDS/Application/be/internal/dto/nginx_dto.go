package dto

// NginxConfigUpdateRequest Nginx 配置更新請求
type NginxConfigUpdateRequest struct {
	Config   string `json:"config" binding:"required"`
	Validate bool   `json:"validate"` // 是否先驗證配置
	Backup   bool   `json:"backup"`   // 是否備份舊配置
}

// NginxReloadRequest Nginx 重載請求
type NginxReloadRequest struct {
	Force   bool `json:"force"`   // 強制重載
	Timeout int  `json:"timeout"` // 超時時間（秒）
}

// NginxUpstreamUpdateRequest Nginx 上游服務更新請求
type NginxUpstreamUpdateRequest struct {
	UpstreamName string              `json:"upstream_name" binding:"required"`
	Servers      []NginxServerConfig `json:"servers" binding:"required,min=1"`
}

// NginxServerConfig Nginx 服務器配置
type NginxServerConfig struct {
	Address string `json:"address" binding:"required"` // host:port
	Weight  int    `json:"weight"`                     // 權重
	MaxFails int   `json:"max_fails"`                  // 最大失敗次數
	FailTimeout int `json:"fail_timeout"`              // 失敗超時（秒）
	Backup  bool   `json:"backup"`                     // 是否為備用服務器
	Down    bool   `json:"down"`                       // 是否停用
}

