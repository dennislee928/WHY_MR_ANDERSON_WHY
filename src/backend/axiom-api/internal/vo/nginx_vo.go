package vo

import "time"

// NginxConfigVO Nginx 配置響應
type NginxConfigVO struct {
	Config       string    `json:"config"`
	ConfigPath   string    `json:"config_path"`
	LastModified time.Time `json:"last_modified"`
	Size         int64     `json:"size"`
	Valid        bool      `json:"valid"`
	Version      string    `json:"version,omitempty"`
}

// NginxStatusVO Nginx 狀態響應
type NginxStatusVO struct {
	Version            string    `json:"version"`
	Address            string    `json:"address"`
	ActiveConnections  int       `json:"active_connections"`
	ServerAccepts      int64     `json:"server_accepts"`
	ServerHandled      int64     `json:"server_handled"`
	ServerRequests     int64     `json:"server_requests"`
	Reading            int       `json:"reading"`
	Writing            int       `json:"writing"`
	Waiting            int       `json:"waiting"`
	Uptime             string    `json:"uptime"`
	RequestsPerSecond  float64   `json:"requests_per_second"`
	ConnectionsPerSecond float64 `json:"connections_per_second"`
	Timestamp          time.Time `json:"timestamp"`
}

// NginxReloadVO Nginx 重載響應
type NginxReloadVO struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Duration  int       `json:"duration"` // 毫秒
	Timestamp time.Time `json:"timestamp"`
}

// NginxUpstreamVO Nginx 上游服務響應
type NginxUpstreamVO struct {
	Name    string               `json:"name"`
	Servers []NginxServerStatusVO `json:"servers"`
	Total   int                  `json:"total"`
	Active  int                  `json:"active"`
	Backup  int                  `json:"backup"`
	Down    int                  `json:"down"`
}

// NginxServerStatusVO Nginx 服務器狀態
type NginxServerStatusVO struct {
	Address     string    `json:"address"`
	Status      string    `json:"status"` // up, down, backup
	Weight      int       `json:"weight"`
	MaxFails    int       `json:"max_fails"`
	FailTimeout int       `json:"fail_timeout"`
	Fails       int       `json:"fails"`
	Unavail     int       `json:"unavail"`
	Requests    int64     `json:"requests"`
	Received    int64     `json:"received"`
	Sent        int64     `json:"sent"`
	LastCheck   time.Time `json:"last_check"`
}

// NginxAccessLogStatsVO Nginx 訪問日誌統計
type NginxAccessLogStatsVO struct {
	TotalRequests   int64              `json:"total_requests"`
	RequestsByStatus map[string]int64  `json:"requests_by_status"`
	RequestsByMethod map[string]int64  `json:"requests_by_method"`
	TopPaths        []PathCount        `json:"top_paths"`
	TopIPs          []IPCount          `json:"top_ips"`
	AverageResponseTime float64        `json:"average_response_time"` // 毫秒
	TimeRange       string             `json:"time_range"`
	Timestamp       time.Time          `json:"timestamp"`
}

// PathCount 路徑計數
type PathCount struct {
	Path  string `json:"path"`
	Count int64  `json:"count"`
}

// IPCount IP 計數
type IPCount struct {
	IP    string `json:"ip"`
	Count int64  `json:"count"`
}

