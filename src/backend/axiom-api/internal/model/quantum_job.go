package model

import (
	"time"

	"gorm.io/datatypes"
)

// QuantumJob 量子作業表
type QuantumJob struct {
	ID            uint           `gorm:"primaryKey"`
	JobID         string         `gorm:"uniqueIndex;not null;size:100"` // 唯一作業 ID
	Type          string         `gorm:"size:50;not null;index"`        // qkd, qsvm, qaoa, qwalk, zerotrust, etc.
	Status        string         `gorm:"size:20;not null;index"`        // pending, running, completed, failed, cancelled
	Backend       string         `gorm:"size:50"`                       // ibm_quantum, simulator
	BackendName   string         `gorm:"size:100"`                      // 實際後端名稱，如 ibm_kyoto
	Circuit       string         `gorm:"type:text"`                     // QASM 電路代碼
	InputData     datatypes.JSON `gorm:"type:jsonb"`                    // 輸入參數
	Result        datatypes.JSON `gorm:"type:jsonb"`                    // 執行結果
	Shots         int            `gorm:"default:1024"`                  // 測量次數
	Qubits        int            // 使用的量子位元數
	Depth         int            // 電路深度
	EstimatedTime int            // 預估執行時間（秒）
	ActualTime    int            // 實際執行時間（秒）
	SubmittedBy   string         `gorm:"size:100"` // 提交者
	SubmittedAt   time.Time      `gorm:"not null;index"`
	StartedAt     *time.Time
	CompletedAt   *time.Time `gorm:"index"`
	Error         string     `gorm:"type:text"`
	Metadata      datatypes.JSON `gorm:"type:jsonb"` // 額外元數據
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TableName 指定表名
func (QuantumJob) TableName() string {
	return "quantum_jobs"
}

// IsCompleted 檢查作業是否完成
func (q *QuantumJob) IsCompleted() bool {
	return q.Status == "completed" || q.Status == "failed" || q.Status == "cancelled"
}

// IsRunning 檢查作業是否正在執行
func (q *QuantumJob) IsRunning() bool {
	return q.Status == "running" || q.Status == "pending"
}

