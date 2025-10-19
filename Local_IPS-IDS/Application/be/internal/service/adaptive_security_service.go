package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// AdaptiveSecurityService 自適應安全服務
type AdaptiveSecurityService struct {
	quantumService *QuantumService
}

// RiskScore 風險評分
type RiskScore struct {
	Score       int                    `json:"score"` // 0-100, 越高越危險
	Level       string                 `json:"level"` // low, medium, high, critical
	Factors     map[string]float64     `json:"factors"`
	Threshold   int                    `json:"threshold"`
	Action      string                 `json:"action"` // allow, challenge, block
	Timestamp   time.Time              `json:"timestamp"`
	Details     map[string]interface{} `json:"details"`
}

// TrustScore 信任分數
type TrustScore struct {
	Score      float64                `json:"score"` // 0.0-1.0
	Components map[string]float64     `json:"components"`
	History    []HistoricalScore      `json:"history"`
	Timestamp  time.Time              `json:"timestamp"`
}

// HistoricalScore 歷史分數
type HistoricalScore struct {
	Timestamp time.Time `json:"timestamp"`
	Score     float64   `json:"score"`
	Event     string    `json:"event"`
}

// NewAdaptiveSecurityService 創建自適應安全服務
func NewAdaptiveSecurityService(quantum *QuantumService) *AdaptiveSecurityService {
	return &AdaptiveSecurityService{
		quantumService: quantum,
	}
}

// CalculateRisk 計算動態風險評分
func (s *AdaptiveSecurityService) CalculateRisk(ctx context.Context, userID, ipAddress string, context map[string]interface{}) (*RiskScore, error) {
	// 1. 收集風險因素
	factors := map[string]float64{
		"ip_reputation":      0.15, // IP 信譽
		"geolocation_risk":   0.10, // 地理位置風險
		"time_of_day":        0.05, // 時間異常
		"failed_auth":        0.25, // 失敗登入次數
		"behavioral_anomaly": 0.30, // 行為異常
		"device_trust":       0.15, // 設備信任度
	}
	
	// 2. 計算總風險分數 (0-100)
	totalScore := 0.0
	for _, weight := range factors {
		// 模擬風險評分
		totalScore += weight * float64(rand.Intn(100))
	}
	
	score := int(totalScore)
	
	// 3. 確定風險等級
	level := "low"
	action := "allow"
	threshold := 30
	
	if score >= 70 {
		level = "critical"
		action = "block"
		threshold = 70
	} else if score >= 50 {
		level = "high"
		action = "challenge"
		threshold = 50
	} else if score >= 30 {
		level = "medium"
		action = "challenge"
		threshold = 30
	}
	
	return &RiskScore{
		Score:     score,
		Level:     level,
		Factors:   factors,
		Threshold: threshold,
		Action:    action,
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"user_id":    userID,
			"ip_address": ipAddress,
		},
	}, nil
}

// EvaluateAccess 評估訪問請求
func (s *AdaptiveSecurityService) EvaluateAccess(ctx context.Context, request map[string]interface{}) (map[string]interface{}, error) {
	userID := fmt.Sprintf("%v", request["user_id"])
	ipAddress := fmt.Sprintf("%v", request["ip_address"])
	
	// 計算風險
	risk, err := s.CalculateRisk(ctx, userID, ipAddress, request)
	if err != nil {
		return nil, err
	}
	
	result := map[string]interface{}{
		"decision":     risk.Action,
		"risk_score":   risk.Score,
		"risk_level":   risk.Level,
		"require_mfa":  risk.Score >= 30,
		"session_ttl":  s.calculateSessionTTL(risk.Score),
		"restrictions": s.getRestrictions(risk.Score),
		"timestamp":    time.Now(),
	}
	
	return result, nil
}

// calculateSessionTTL 根據風險計算會話時長
func (s *AdaptiveSecurityService) calculateSessionTTL(riskScore int) string {
	if riskScore >= 50 {
		return "15m" // 高風險，短會話
	} else if riskScore >= 30 {
		return "1h" // 中風險，中等會話
	}
	return "24h" // 低風險，長會話
}

// getRestrictions 根據風險獲取訪問限制
func (s *AdaptiveSecurityService) getRestrictions(riskScore int) []string {
	if riskScore >= 70 {
		return []string{"block_all", "notify_security_team"}
	} else if riskScore >= 50 {
		return []string{"require_mfa", "limit_api_access", "enhanced_logging"}
	} else if riskScore >= 30 {
		return []string{"require_mfa", "monitor_closely"}
	}
	return []string{"normal_access"}
}

// GetTrustScore 獲取信任分數
func (s *AdaptiveSecurityService) GetTrustScore(ctx context.Context, entityID string) (*TrustScore, error) {
	// 計算多維度信任分數
	components := map[string]float64{
		"authentication_history": 0.85,
		"behavioral_pattern":     0.78,
		"device_security":        0.92,
		"network_reputation":     0.75,
	}
	
	// 加權平均
	totalScore := 0.0
	for _, score := range components {
		totalScore += score
	}
	avgScore := totalScore / float64(len(components))
	
	// 歷史分數
	history := []HistoricalScore{
		{Timestamp: time.Now().Add(-2 * time.Hour), Score: 0.82, Event: "normal"},
		{Timestamp: time.Now().Add(-1 * time.Hour), Score: 0.79, Event: "location_change"},
		{Timestamp: time.Now(), Score: avgScore, Event: "current"},
	}
	
	return &TrustScore{
		Score:      avgScore,
		Components: components,
		History:    history,
		Timestamp:  time.Now(),
	}, nil
}

// DeployHoneypot 部署蜜罐
type HoneypotDeployment struct {
	HoneypotID   string                 `json:"honeypot_id"`
	Type         string                 `json:"type"` // ssh, http, database
	IPAddress    string                 `json:"ip_address"`
	Port         int                    `json:"port"`
	Status       string                 `json:"status"`
	Interactions int                    `json:"interactions"`
	Attackers    []string               `json:"attackers"`
	DeployedAt   time.Time              `json:"deployed_at"`
	Config       map[string]interface{} `json:"config"`
}

// DeployHoneypot 自動部署蜜罐
func (s *AdaptiveSecurityService) DeployHoneypot(ctx context.Context, honeypotType string) (*HoneypotDeployment, error) {
	honeypotID := fmt.Sprintf("HONEY-%d", time.Now().Unix())
	
	// 選擇端口
	port := 22
	if honeypotType == "http" {
		port = 8080
	} else if honeypotType == "database" {
		port = 3306
	}
	
	deployment := &HoneypotDeployment{
		HoneypotID:   honeypotID,
		Type:         honeypotType,
		IPAddress:    "192.168.100.10",
		Port:         port,
		Status:       "active",
		Interactions: 0,
		Attackers:    []string{},
		DeployedAt:   time.Now(),
		Config: map[string]interface{}{
			"logging":      true,
			"alert_on_access": true,
			"fake_data":    true,
		},
	}
	
	return deployment, nil
}

// GetHoneypotInteractions 獲取蜜罐互動記錄
func (s *AdaptiveSecurityService) GetHoneypotInteractions(ctx context.Context, honeypotID string) ([]map[string]interface{}, error) {
	// 模擬蜜罐互動記錄
	interactions := []map[string]interface{}{
		{
			"timestamp":   time.Now().Add(-1 * time.Hour),
			"source_ip":   "203.0.113.45",
			"action":      "ssh_login_attempt",
			"username":    "admin",
			"successful":  false,
			"fingerprint": "SSH-2.0-OpenSSH_7.4",
		},
		{
			"timestamp":   time.Now().Add(-30 * time.Minute),
			"source_ip":   "203.0.113.45",
			"action":      "port_scan",
			"ports":       []int{22, 80, 443, 3306},
			"successful":  true,
		},
	}
	
	return interactions, nil
}

// AnalyzeAttacker 分析攻擊者行為
func (s *AdaptiveSecurityService) AnalyzeAttacker(ctx context.Context, attackerIP string) (map[string]interface{}, error) {
	analysis := map[string]interface{}{
		"attacker_ip": attackerIP,
		"threat_level": "high",
		"attack_patterns": []string{
			"SSH brute force",
			"Port scanning",
			"SQL injection attempts",
		},
		"geolocation": map[string]string{
			"country": "Unknown",
			"region":  "Unknown",
		},
		"reputation": map[string]interface{}{
			"blacklisted":     true,
			"spam_score":      85,
			"bot_probability": 0.92,
		},
		"recommended_action": "block_immediately",
		"confidence":         0.95,
	}
	
	return analysis, nil
}


