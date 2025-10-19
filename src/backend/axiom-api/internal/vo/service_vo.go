package vo

import "time"

// ServiceStatusVO 服務狀態響應
type ServiceStatusVO struct {
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Status    string                 `json:"status"`
	URL       string                 `json:"url"`
	Version   string                 `json:"version"`
	LastCheck time.Time              `json:"last_check"`
	Metrics   map[string]interface{} `json:"metrics,omitempty"`
	Uptime    string                 `json:"uptime,omitempty"`
}

// ServicesListVO 服務列表響應
type ServicesListVO struct {
	Services  []ServiceStatusVO `json:"services"`
	Total     int               `json:"total"`
	Healthy   int               `json:"healthy"`
	Unhealthy int               `json:"unhealthy"`
	Unknown   int               `json:"unknown"`
	Timestamp time.Time         `json:"timestamp"`
}

// ServiceHealthCheckVO 服務健康檢查響應
type ServiceHealthCheckVO struct {
	ServiceName string    `json:"service_name"`
	Status      string    `json:"status"`
	Healthy     bool      `json:"healthy"`
	Message     string    `json:"message,omitempty"`
	CheckedAt   time.Time `json:"checked_at"`
}

// ServiceRestartVO 服務重啟響應
type ServiceRestartVO struct {
	ServiceName string    `json:"service_name"`
	Status      string    `json:"status"` // success, failed, pending
	Message     string    `json:"message"`
	StartedAt   time.Time `json:"started_at"`
}

