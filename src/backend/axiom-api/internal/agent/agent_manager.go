package agent

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"axiom-backend/internal/database"
)

// AgentInfo Agent 信息
type AgentInfo struct {
	AgentID       string                 `json:"agent_id"`
	Mode          AgentMode              `json:"mode"`
	Hostname      string                 `json:"hostname"`
	IPAddress     string                 `json:"ip_address"`
	Status        string                 `json:"status"` // active, inactive, offline
	APIKey        string                 `json:"api_key,omitempty"`
	ClientCert    string                 `json:"client_cert,omitempty"`
	Capabilities  []string               `json:"capabilities"`
	LastHeartbeat time.Time              `json:"last_heartbeat"`
	RegisteredAt  time.Time              `json:"registered_at"`
	Config        *AgentConfig           `json:"config"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// HeartbeatData 心跳數據
type HeartbeatData struct {
	AgentID   string                 `json:"agent_id"`
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Metrics   map[string]interface{} `json:"metrics"`
}

// AgentManager Agent 管理器
type AgentManager struct {
	db     *database.Database
	agents sync.Map // map[string]*AgentInfo
	mu     sync.RWMutex
}

// NewAgentManager 創建 Agent 管理器
func NewAgentManager(db *database.Database) *AgentManager {
	return &AgentManager{
		db: db,
	}
}

// RegisterAgent 註冊 Agent
func (m *AgentManager) RegisterAgent(ctx context.Context, req *AgentRegistrationRequest) (*AgentRegistrationResponse, error) {
	agentID := req.AgentID
	if agentID == "" {
		agentID = generateAgentID(req.Mode)
	}
	
	// 生成 API Key
	apiKey := generateAPIKey()
	
	// 根據模式生成配置
	var config *AgentConfig
	if req.Mode == AgentModeExternal {
		config = DefaultExternalConfig()
		// 為外部 Agent 生成證書（實際應該使用真實的證書簽發）
		config.ClientCert = "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"
		config.ClientKey = "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----"
	} else {
		config = DefaultInternalConfig()
	}
	
	config.ID = agentID
	config.Hostname = req.Hostname
	config.IPAddress = req.IPAddress
	config.Endpoint = m.getEndpoint(req.Mode)
	
	// 創建 Agent 信息
	agentInfo := &AgentInfo{
		AgentID:       agentID,
		Mode:          req.Mode,
		Hostname:      req.Hostname,
		IPAddress:     req.IPAddress,
		Status:        "active",
		APIKey:        apiKey,
		Capabilities:  req.Capabilities,
		LastHeartbeat: time.Now(),
		RegisteredAt:  time.Now(),
		Config:        config,
		Metadata:      req.Metadata,
	}
	
	// 存儲到內存和數據庫
	m.agents.Store(agentID, agentInfo)
	
	// TODO: 保存到數據庫
	// m.db.PG.Create(&AgentModel{...})
	
	return &AgentRegistrationResponse{
		AgentID:          agentID,
		APIKey:           apiKey,
		ClientCert:       config.ClientCert,
		ClientKey:        config.ClientKey,
		CACert:           config.CACert,
		HeartbeatInterval: 30, // 秒
		Config:           config,
	}, nil
}

// ProcessHeartbeat 處理心跳
func (m *AgentManager) ProcessHeartbeat(ctx context.Context, heartbeat *HeartbeatData) error {
	value, ok := m.agents.Load(heartbeat.AgentID)
	if !ok {
		return fmt.Errorf("agent not found: %s", heartbeat.AgentID)
	}
	
	agentInfo := value.(*AgentInfo)
	agentInfo.LastHeartbeat = heartbeat.Timestamp
	agentInfo.Status = heartbeat.Status
	
	m.agents.Store(heartbeat.AgentID, agentInfo)
	
	// TODO: 更新數據庫
	
	return nil
}

// GetAgent 獲取 Agent 信息
func (m *AgentManager) GetAgent(ctx context.Context, agentID string) (*AgentInfo, error) {
	value, ok := m.agents.Load(agentID)
	if !ok {
		return nil, fmt.Errorf("agent not found: %s", agentID)
	}
	
	return value.(*AgentInfo), nil
}

// ListAgents 列出所有 Agent
func (m *AgentManager) ListAgents(ctx context.Context, mode AgentMode) ([]*AgentInfo, error) {
	var agents []*AgentInfo
	
	m.agents.Range(func(key, value interface{}) bool {
		agentInfo := value.(*AgentInfo)
		if mode == "" || agentInfo.Mode == mode {
			agents = append(agents, agentInfo)
		}
		return true
	})
	
	return agents, nil
}

// DeregisterAgent 註銷 Agent
func (m *AgentManager) DeregisterAgent(ctx context.Context, agentID string) error {
	m.agents.Delete(agentID)
	
	// TODO: 從數據庫刪除
	
	return nil
}

// UpdateAgentConfig 更新 Agent 配置
func (m *AgentManager) UpdateAgentConfig(ctx context.Context, agentID string, config *AgentConfig) error {
	value, ok := m.agents.Load(agentID)
	if !ok {
		return fmt.Errorf("agent not found: %s", agentID)
	}
	
	agentInfo := value.(*AgentInfo)
	agentInfo.Config = config
	
	m.agents.Store(agentID, agentInfo)
	
	// TODO: 更新數據庫
	
	return nil
}

// CheckAgentHealth 檢查 Agent 健康狀態
func (m *AgentManager) CheckAgentHealth(ctx context.Context) map[string]string {
	healthStatus := make(map[string]string)
	timeout := 60 * time.Second
	
	m.agents.Range(func(key, value interface{}) bool {
		agentInfo := value.(*AgentInfo)
		agentID := agentInfo.AgentID
		
		if time.Since(agentInfo.LastHeartbeat) > timeout {
			healthStatus[agentID] = "offline"
			agentInfo.Status = "offline"
			m.agents.Store(agentID, agentInfo)
		} else {
			healthStatus[agentID] = agentInfo.Status
		}
		
		return true
	})
	
	return healthStatus
}

// getEndpoint 根據模式獲取端點
func (m *AgentManager) getEndpoint(mode AgentMode) string {
	if mode == AgentModeExternal {
		return "https://axiom.yourdomain.com/api/v2/agent"
	}
	return "http://axiom-backend.internal:8080/api/v2/agent"
}

// generateAgentID 生成 Agent ID
func generateAgentID(mode AgentMode) string {
	prefix := "ext"
	if mode == AgentModeInternal {
		prefix = "int"
	}
	
	bytes := make([]byte, 8)
	rand.Read(bytes)
	
	return fmt.Sprintf("agent-%s-%s", prefix, hex.EncodeToString(bytes)[:12])
}

// generateAPIKey 生成 API Key
func generateAPIKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// AgentRegistrationRequest Agent 註冊請求
type AgentRegistrationRequest struct {
	AgentID      string                 `json:"agent_id,omitempty"`
	Mode         AgentMode              `json:"mode"`
	Hostname     string                 `json:"hostname"`
	IPAddress    string                 `json:"ip_address"`
	Capabilities []string               `json:"capabilities"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// AgentRegistrationResponse Agent 註冊響應
type AgentRegistrationResponse struct {
	AgentID           string       `json:"agent_id"`
	APIKey            string       `json:"api_key"`
	ClientCert        string       `json:"client_cert,omitempty"`
	ClientKey         string       `json:"client_key,omitempty"`
	CACert            string       `json:"ca_cert,omitempty"`
	HeartbeatInterval int          `json:"heartbeat_interval"` // seconds
	Config            *AgentConfig `json:"config"`
}

