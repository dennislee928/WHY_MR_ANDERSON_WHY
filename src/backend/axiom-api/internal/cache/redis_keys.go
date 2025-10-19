package cache

import (
	"fmt"
	"time"
)

// Redis Key 常量定義

const (
	// 服務健康狀態快取 (TTL: 30s)
	ServiceHealthKeyPrefix = "service:health:"
	ServiceHealthTTL       = 30 * time.Second

	// 即時指標快取 (TTL: 10s)
	MetricsRealtimeKeyPrefix = "metrics:realtime:"
	MetricsRealtimeTTL       = 10 * time.Second

	// 量子作業狀態快取 (TTL: 5min)
	QuantumJobKeyPrefix = "quantum:job:"
	QuantumJobTTL       = 5 * time.Minute

	// API 速率限制 (TTL: 1min)
	RateLimitAPIKeyPrefix = "ratelimit:api:"
	RateLimitAPITTL       = 1 * time.Minute

	// 會話管理 (TTL: 24h)
	SessionKeyPrefix = "session:"
	SessionTTL       = 24 * time.Hour

	// 即時統計計數器 (no TTL)
	CounterAPIRequestsKey    = "counter:api:requests:total"
	CounterThreatsKey        = "counter:threats:detected:total"
	CounterQuantumJobsKey    = "counter:quantum:jobs:completed"
	CounterWindowsLogsKey    = "counter:windows:logs:received"

	// 服務配置快取 (TTL: 5min)
	ServiceConfigKeyPrefix = "service:config:"
	ServiceConfigTTL       = 5 * time.Minute

	// Nginx 狀態快取 (TTL: 5s)
	NginxStatusKey = "nginx:status"
	NginxStatusTTL = 5 * time.Second

	// 告警快取 (TTL: 1min)
	AlertsActiveKey = "alerts:active"
	AlertsTTL       = 1 * time.Minute

	// Windows 日誌統計快取 (TTL: 1min)
	WindowsLogStatsKey = "windows_logs:stats"
	WindowsLogStatsTTL = 1 * time.Minute

	// 量子統計快取 (TTL: 30s)
	QuantumStatsKey = "quantum:stats"
	QuantumStatsTTL = 30 * time.Second

	// 用戶登入嘗試計數 (TTL: 15min)
	UserLoginAttemptsKeyPrefix = "user:login:attempts:"
	UserLoginAttemptsTTL       = 15 * time.Minute

	// API Token 驗證快取 (TTL: 5min)
	APITokenKeyPrefix = "api:token:"
	APITokenTTL       = 5 * time.Minute
)

// Key 生成函數

// ServiceHealthKey 生成服務健康狀態 Key
func ServiceHealthKey(serviceName string) string {
	return fmt.Sprintf("%s%s", ServiceHealthKeyPrefix, serviceName)
}

// MetricsRealtimeKey 生成即時指標 Key
func MetricsRealtimeKey(serviceName string) string {
	return fmt.Sprintf("%s%s", MetricsRealtimeKeyPrefix, serviceName)
}

// QuantumJobKey 生成量子作業 Key
func QuantumJobKey(jobID string) string {
	return fmt.Sprintf("%s%s", QuantumJobKeyPrefix, jobID)
}

// RateLimitAPIKey 生成 API 速率限制 Key
func RateLimitAPIKey(clientIP, endpoint string) string {
	return fmt.Sprintf("%s%s:%s", RateLimitAPIKeyPrefix, clientIP, endpoint)
}

// SessionKey 生成會話 Key
func SessionKey(sessionID string) string {
	return fmt.Sprintf("%s%s", SessionKeyPrefix, sessionID)
}

// ServiceConfigKey 生成服務配置 Key
func ServiceConfigKey(serviceName string) string {
	return fmt.Sprintf("%s%s", ServiceConfigKeyPrefix, serviceName)
}

// UserLoginAttemptsKey 生成用戶登入嘗試 Key
func UserLoginAttemptsKey(username string) string {
	return fmt.Sprintf("%s%s", UserLoginAttemptsKeyPrefix, username)
}

// APITokenKey 生成 API Token Key
func APITokenKey(token string) string {
	return fmt.Sprintf("%s%s", APITokenKeyPrefix, token)
}

