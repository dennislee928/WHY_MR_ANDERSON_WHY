package dto

// QuantumQKDRequest 量子密鑰分發請求
type QuantumQKDRequest struct {
	KeyLength int    `json:"key_length" binding:"required,min=64,max=2048"` // 64-2048 bits
	Backend   string `json:"backend"`                                       // ibm_quantum, simulator
	Shots     int    `json:"shots" binding:"min=1,max=8192"`                // 1-8192
}

// QuantumEncryptRequest 量子加密請求
type QuantumEncryptRequest struct {
	Plaintext string `json:"plaintext" binding:"required"`
	Backend   string `json:"backend"` // ibm_quantum, simulator
}

// QuantumDecryptRequest 量子解密請求
type QuantumDecryptRequest struct {
	Ciphertext string `json:"ciphertext" binding:"required"`
	Key        string `json:"key" binding:"required"`
	Backend    string `json:"backend"`
}

// QuantumQSVMRequest 量子支持向量機請求
type QuantumQSVMRequest struct {
	Features    []float64 `json:"features" binding:"required,min=2,max=100"`
	Backend     string    `json:"backend"`
	FeatureDim  int       `json:"feature_dim" binding:"min=2,max=20"`
	Shots       int       `json:"shots" binding:"min=1,max=8192"`
}

// QuantumQAOARequest 量子近似優化算法請求
type QuantumQAOARequest struct {
	Graph    map[string]interface{} `json:"graph" binding:"required"`    // 圖結構
	PLayers  int                    `json:"p_layers" binding:"min=1,max=10"`
	Backend  string                 `json:"backend"`
	Shots    int                    `json:"shots" binding:"min=1,max=8192"`
}

// QuantumWalkRequest 量子漫步算法請求
type QuantumWalkRequest struct {
	Nodes     int    `json:"nodes" binding:"required,min=4,max=100"`
	Steps     int    `json:"steps" binding:"min=1,max=50"`
	Backend   string `json:"backend"`
	Shots     int    `json:"shots" binding:"min=1,max=8192"`
}

// ZeroTrustPredictRequest Zero Trust 預測請求
type ZeroTrustPredictRequest struct {
	UserID      string                 `json:"user_id" binding:"required"`
	IPAddress   string                 `json:"ip_address" binding:"required"`
	DeviceID    string                 `json:"device_id"`
	Timestamp   string                 `json:"timestamp"`
	Features    map[string]interface{} `json:"features"` // 額外特徵
	UseQuantum  bool                   `json:"use_quantum"` // 是否使用量子算法
}

// QuantumJobQueryRequest 量子作業查詢請求
type QuantumJobQueryRequest struct {
	JobID     string   `form:"job_id"`
	Type      string   `form:"type"`      // qkd, qsvm, qaoa, etc.
	Status    string   `form:"status"`    // pending, running, completed, failed
	Backend   string   `form:"backend"`
	StartDate string   `form:"start_date"` // YYYY-MM-DD
	EndDate   string   `form:"end_date"`   // YYYY-MM-DD
	Page      int      `form:"page" binding:"min=1"`
	PageSize  int      `form:"page_size" binding:"min=1,max=100"`
}

