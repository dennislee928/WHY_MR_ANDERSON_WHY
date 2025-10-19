package agent

import (
	"time"
)

// AgentMode Agent 連接模式
type AgentMode string

const (
	AgentModeExternal AgentMode = "external" // 外部連接，通過 Nginx
	AgentModeInternal AgentMode = "internal" // 內部直連
)

// AgentConfig Agent 配置
type AgentConfig struct {
	// 基本信息
	Mode              AgentMode `json:"mode"`
	ID                string    `json:"id"`
	Hostname          string    `json:"hostname"`
	IPAddress         string    `json:"ip_address"`
	
	// 連接配置
	Endpoint          string    `json:"endpoint"`
	ThroughNginx      bool      `json:"through_nginx"`
	DirectConnect     bool      `json:"direct_connect"`
	
	// 認證配置
	AuthMethod        string    `json:"auth_method"` // "mtls" or "api_key"
	APIKey            string    `json:"api_key"`
	ClientCert        string    `json:"client_cert,omitempty"`
	ClientKey         string    `json:"client_key,omitempty"`
	CACert            string    `json:"ca_cert,omitempty"`
	
	// 上傳配置
	UploadMethod      string        `json:"upload_method"` // "streaming" or "batch"
	BatchSize         int           `json:"batch_size"`
	FlushInterval     time.Duration `json:"flush_interval"`
	MaxRetry          int           `json:"max_retry"`
	RetryBackoff      string        `json:"retry_backoff"` // "exponential" or "linear"
	Compression       bool          `json:"compression"`
	
	// 緩衝配置
	BufferType        string `json:"buffer_type"` // "persistent" or "memory"
	BufferPath        string `json:"buffer_path,omitempty"`
	BufferMaxSize     int64  `json:"buffer_max_size"`
	OverflowStrategy  string `json:"overflow_strategy"` // "drop_oldest", "drop_newest", "block"
	
	// 安全配置
	EncryptInTransit  bool          `json:"encrypt_in_transit"`
	EncryptAtRest     bool          `json:"encrypt_at_rest"`
	KeyRotation       time.Duration `json:"key_rotation,omitempty"`
	
	// 能力
	Capabilities      []string `json:"capabilities"`
}

// DefaultExternalConfig 默認外部配置
func DefaultExternalConfig() *AgentConfig {
	return &AgentConfig{
		Mode:              AgentModeExternal,
		ThroughNginx:      true,
		DirectConnect:     false,
		AuthMethod:        "mtls",
		UploadMethod:      "streaming",
		BatchSize:         100,
		FlushInterval:     10 * time.Second,
		MaxRetry:          5,
		RetryBackoff:      "exponential",
		Compression:       true,
		BufferType:        "persistent",
		BufferMaxSize:     1024 * 1024 * 1024, // 1GB
		OverflowStrategy:  "drop_oldest",
		EncryptInTransit:  true,
		EncryptAtRest:     true,
		KeyRotation:       24 * time.Hour,
		Capabilities:      []string{"windows_logs", "metrics", "compliance_scan"},
	}
}

// DefaultInternalConfig 默認內部配置
func DefaultInternalConfig() *AgentConfig {
	return &AgentConfig{
		Mode:              AgentModeInternal,
		ThroughNginx:      false,
		DirectConnect:     true,
		AuthMethod:        "api_key",
		UploadMethod:      "streaming",
		BatchSize:         500,
		FlushInterval:     5 * time.Second,
		MaxRetry:          3,
		RetryBackoff:      "linear",
		Compression:       false,
		BufferType:        "memory",
		BufferMaxSize:     256 * 1024 * 1024, // 256MB
		OverflowStrategy:  "block",
		EncryptInTransit:  false,
		EncryptAtRest:     false,
		Capabilities:      []string{"windows_logs", "metrics"},
	}
}

