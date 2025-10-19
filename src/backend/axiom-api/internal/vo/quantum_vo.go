package vo

import "time"

// QuantumJobVO 量子作業響應
type QuantumJobVO struct {
	JobID         string                 `json:"job_id"`
	Type          string                 `json:"type"`
	Status        string                 `json:"status"`
	Backend       string                 `json:"backend"`
	BackendName   string                 `json:"backend_name,omitempty"`
	Result        map[string]interface{} `json:"result,omitempty"`
	Qubits        int                    `json:"qubits,omitempty"`
	Depth         int                    `json:"depth,omitempty"`
	Shots         int                    `json:"shots,omitempty"`
	EstimatedTime int                    `json:"estimated_time,omitempty"` // 秒
	ActualTime    int                    `json:"actual_time,omitempty"`    // 秒
	Progress      int                    `json:"progress,omitempty"`       // 0-100%
	SubmittedAt   time.Time              `json:"submitted_at"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty"`
	Error         string                 `json:"error,omitempty"`
}

// QuantumQKDVO 量子密鑰分發響應
type QuantumQKDVO struct {
	JobID       string    `json:"job_id"`
	Key         string    `json:"key,omitempty"`         // Base64 encoded
	KeyLength   int       `json:"key_length"`
	ErrorRate   float64   `json:"error_rate,omitempty"`
	Status      string    `json:"status"`
	SubmittedAt time.Time `json:"submitted_at"`
	Message     string    `json:"message,omitempty"`
}

// QuantumEncryptVO 量子加密響應
type QuantumEncryptVO struct {
	Ciphertext string `json:"ciphertext"` // Base64 encoded
	Key        string `json:"key"`        // Base64 encoded
	Algorithm  string `json:"algorithm"`
	Message    string `json:"message,omitempty"`
}

// QuantumClassifyVO 量子分類響應（QSVM）
type QuantumClassifyVO struct {
	JobID       string    `json:"job_id"`
	Prediction  int       `json:"prediction"`           // 0 or 1
	Probability float64   `json:"probability"`          // 0.0-1.0
	Confidence  float64   `json:"confidence,omitempty"` // 0.0-1.0
	Status      string    `json:"status"`
	SubmittedAt time.Time `json:"submitted_at"`
	Message     string    `json:"message,omitempty"`
}

// ZeroTrustPredictVO Zero Trust 預測響應
type ZeroTrustPredictVO struct {
	TrustScore  float64                `json:"trust_score"`  // 0.0-1.0
	RiskLevel   string                 `json:"risk_level"`   // low, medium, high, critical
	Decision    string                 `json:"decision"`     // allow, deny, mfa_required
	Confidence  float64                `json:"confidence"`   // 0.0-1.0
	Factors     map[string]interface{} `json:"factors"`      // 影響因素
	UsedQuantum bool                   `json:"used_quantum"` // 是否使用了量子算法
	Timestamp   time.Time              `json:"timestamp"`
	Message     string                 `json:"message,omitempty"`
}

// QuantumJobsListVO 量子作業列表響應
type QuantumJobsListVO struct {
	Jobs       []QuantumJobVO `json:"jobs"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// QuantumStatsVO 量子統計響應
type QuantumStatsVO struct {
	TotalJobs       int            `json:"total_jobs"`
	CompletedJobs   int            `json:"completed_jobs"`
	FailedJobs      int            `json:"failed_jobs"`
	RunningJobs     int            `json:"running_jobs"`
	AverageTime     float64        `json:"average_time"` // 秒
	JobsByType      map[string]int `json:"jobs_by_type"`
	JobsByBackend   map[string]int `json:"jobs_by_backend"`
	SuccessRate     float64        `json:"success_rate"` // 0.0-1.0
	Last24Hours     int            `json:"last_24_hours"`
}

