package axiom

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"pandora_box_console_ids_ips/internal/metrics"
)

// AnalysisEngine Axiom分析引擎
type AnalysisEngine struct {
	logger      *logrus.Logger
	metrics     *metrics.PrometheusMetrics
	rules       []*SecurityRule
	blacklist   map[string]time.Time
	whitelist   map[string]bool
	threatCache map[string]*ThreatInfo
	mutex       sync.RWMutex
	running     bool
	stopChan    chan struct{}
}

// SecurityRule 安全規則
type SecurityRule struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // "ip", "port", "pattern", "behavior"
	Pattern     string    `json:"pattern"`
	Action      string    `json:"action"`   // "block", "alert", "log"
	Severity    string    `json:"severity"` // "low", "medium", "high", "critical"
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	regex       *regexp.Regexp
}

// ThreatInfo 威脅資訊
type ThreatInfo struct {
	SourceIP   string                 `json:"source_ip"`
	ThreatType string                 `json:"threat_type"`
	Severity   string                 `json:"severity"`
	FirstSeen  time.Time              `json:"first_seen"`
	LastSeen   time.Time              `json:"last_seen"`
	Count      int                    `json:"count"`
	Details    map[string]interface{} `json:"details"`
	Blocked    bool                   `json:"blocked"`
}

// NetworkPacket 網路封包結構
type NetworkPacket struct {
	SourceIP      string            `json:"source_ip"`
	DestIP        string            `json:"dest_ip"`
	SourcePort    int               `json:"source_port"`
	DestPort      int               `json:"dest_port"`
	Protocol      string            `json:"protocol"`
	Payload       []byte            `json:"payload"`
	PayloadString string            `json:"payload_string"`
	Timestamp     time.Time         `json:"timestamp"`
	Headers       map[string]string `json:"headers"`
}

// AnalysisResult 分析結果
type AnalysisResult struct {
	PacketID    string    `json:"packet_id"`
	Timestamp   time.Time `json:"timestamp"`
	SourceIP    string    `json:"source_ip"`
	ThreatLevel string    `json:"threat_level"`
	ThreatType  string    `json:"threat_type"`
	Action      string    `json:"action"`
	RuleID      string    `json:"rule_id"`
	RuleName    string    `json:"rule_name"`
	Details     string    `json:"details"`
	Blocked     bool      `json:"blocked"`
}

// NewAnalysisEngine 建立新的分析引擎
func NewAnalysisEngine(logger *logrus.Logger, metrics *metrics.PrometheusMetrics) *AnalysisEngine {
	engine := &AnalysisEngine{
		logger:      logger,
		metrics:     metrics,
		rules:       make([]*SecurityRule, 0),
		blacklist:   make(map[string]time.Time),
		whitelist:   make(map[string]bool),
		threatCache: make(map[string]*ThreatInfo),
		stopChan:    make(chan struct{}),
	}

	// 載入預設規則
	engine.loadDefaultRules()

	return engine
}

// Start 啟動分析引擎
func (ae *AnalysisEngine) Start(ctx context.Context) error {
	ae.mutex.Lock()
	if ae.running {
		ae.mutex.Unlock()
		return fmt.Errorf("分析引擎已在運行中")
	}
	ae.running = true
	ae.mutex.Unlock()

	ae.logger.Info("Axiom分析引擎已啟動")

	// 啟動威脅快取清理
	go ae.startThreatCacheCleanup()

	// 啟動黑名單清理
	go ae.startBlacklistCleanup()

	// 等待停止信號
	select {
	case <-ctx.Done():
		ae.Stop()
		return ctx.Err()
	case <-ae.stopChan:
		return nil
	}
}

// Stop 停止分析引擎
func (ae *AnalysisEngine) Stop() {
	ae.mutex.Lock()
	defer ae.mutex.Unlock()

	if !ae.running {
		return
	}

	ae.running = false
	close(ae.stopChan)
	ae.logger.Info("Axiom分析引擎已停止")
}

// AnalyzePacket 分析網路封包
func (ae *AnalysisEngine) AnalyzePacket(packet *NetworkPacket) (*AnalysisResult, error) {
	ae.mutex.RLock()
	defer ae.mutex.RUnlock()

	if !ae.running {
		return nil, fmt.Errorf("分析引擎未運行")
	}

	// 記錄安全事件指標 (簡化版)
	ae.metrics.RecordSecurityEvent("packet_analysis", "success")

	result := &AnalysisResult{
		PacketID:    fmt.Sprintf("%d", time.Now().UnixNano()),
		Timestamp:   packet.Timestamp,
		SourceIP:    packet.SourceIP,
		ThreatLevel: "low",
		Action:      "allow",
		Blocked:     false,
	}

	// 檢查白名單
	if ae.whitelist[packet.SourceIP] {
		result.Details = "IP在白名單中"
		return result, nil
	}

	// 檢查黑名單
	if _, blocked := ae.blacklist[packet.SourceIP]; blocked {
		result.ThreatLevel = "high"
		result.ThreatType = "blacklisted_ip"
		result.Action = "block"
		result.Blocked = true
		result.Details = "IP在黑名單中"

		ae.metrics.RecordSecurityEvent("blacklisted_ip", "blocked")
		return result, nil
	}

	// 應用安全規則
	for _, rule := range ae.rules {
		if !rule.Enabled {
			continue
		}

		if ae.matchRule(rule, packet) {
			result.ThreatLevel = rule.Severity
			result.ThreatType = rule.Type
			result.Action = rule.Action
			result.RuleID = rule.ID
			result.RuleName = rule.Name
			result.Details = fmt.Sprintf("觸發規則: %s", rule.Description)

			if rule.Action == "block" {
				result.Blocked = true
				ae.addToBlacklist(packet.SourceIP, time.Hour) // 暫時加入黑名單1小時
			}

			// 記錄威脅
			ae.recordThreat(packet.SourceIP, rule.Type, rule.Severity, result.Details)
			ae.metrics.RecordSecurityEvent(rule.Type, "detected")

			break
		}
	}

	return result, nil
}

// matchRule 檢查規則是否匹配
func (ae *AnalysisEngine) matchRule(rule *SecurityRule, packet *NetworkPacket) bool {
	switch rule.Type {
	case "ip":
		return ae.matchIPRule(rule, packet.SourceIP)
	case "port":
		return ae.matchPortRule(rule, packet.DestPort)
	case "pattern":
		return ae.matchPatternRule(rule, packet.PayloadString)
	case "behavior":
		return ae.matchBehaviorRule(rule, packet)
	default:
		return false
	}
}

// matchIPRule 匹配IP規則
func (ae *AnalysisEngine) matchIPRule(rule *SecurityRule, ip string) bool {
	if rule.regex != nil {
		return rule.regex.MatchString(ip)
	}
	return strings.Contains(ip, rule.Pattern)
}

// matchPortRule 匹配連接埠規則
func (ae *AnalysisEngine) matchPortRule(rule *SecurityRule, port int) bool {
	portStr := fmt.Sprintf("%d", port)
	if rule.regex != nil {
		return rule.regex.MatchString(portStr)
	}
	return strings.Contains(portStr, rule.Pattern)
}

// matchPatternRule 匹配模式規則
func (ae *AnalysisEngine) matchPatternRule(rule *SecurityRule, payload string) bool {
	if rule.regex != nil {
		return rule.regex.MatchString(payload)
	}
	return strings.Contains(strings.ToLower(payload), strings.ToLower(rule.Pattern))
}

// matchBehaviorRule 匹配行為規則
func (ae *AnalysisEngine) matchBehaviorRule(rule *SecurityRule, packet *NetworkPacket) bool {
	// 實現行為分析邏輯
	switch rule.Pattern {
	case "port_scan":
		return ae.detectPortScan(packet.SourceIP)
	case "brute_force":
		return ae.detectBruteForce(packet)
	case "ddos":
		return ae.detectDDoS(packet.SourceIP)
	default:
		return false
	}
}

// detectPortScan 偵測連接埠掃描
func (ae *AnalysisEngine) detectPortScan(sourceIP string) bool {
	// 簡化的連接埠掃描偵測邏輯
	threat, exists := ae.threatCache[sourceIP]
	if !exists {
		return false
	}

	// 如果在短時間內連接多個不同連接埠，可能是掃描
	if time.Since(threat.LastSeen) < 1*time.Minute && threat.Count > 10 {
		return true
	}

	return false
}

// detectBruteForce 偵測暴力破解
func (ae *AnalysisEngine) detectBruteForce(packet *NetworkPacket) bool {
	// 檢查是否為SSH、FTP、HTTP認證等常見服務
	commonPorts := []int{22, 21, 80, 443, 3389}
	for _, port := range commonPorts {
		if packet.DestPort == port {
			threat, exists := ae.threatCache[packet.SourceIP]
			if exists && threat.Count > 5 && time.Since(threat.FirstSeen) < 5*time.Minute {
				return true
			}
		}
	}
	return false
}

// detectDDoS 偵測DDoS攻擊
func (ae *AnalysisEngine) detectDDoS(sourceIP string) bool {
	threat, exists := ae.threatCache[sourceIP]
	if !exists {
		return false
	}

	// 如果在短時間內有大量連接，可能是DDoS
	if time.Since(threat.FirstSeen) < 1*time.Minute && threat.Count > 100 {
		return true
	}

	return false
}

// recordThreat 記錄威脅資訊
func (ae *AnalysisEngine) recordThreat(sourceIP, threatType, severity, details string) {
	threat, exists := ae.threatCache[sourceIP]
	if !exists {
		threat = &ThreatInfo{
			SourceIP:   sourceIP,
			ThreatType: threatType,
			Severity:   severity,
			FirstSeen:  time.Now(),
			Count:      1,
			Details:    make(map[string]interface{}),
		}
		ae.threatCache[sourceIP] = threat
	} else {
		threat.Count++
	}

	threat.LastSeen = time.Now()
	threat.Details["last_details"] = details
}

// AddRule 添加安全規則
func (ae *AnalysisEngine) AddRule(rule *SecurityRule) error {
	ae.mutex.Lock()
	defer ae.mutex.Unlock()

	// 編譯正則表達式
	if strings.HasPrefix(rule.Pattern, "regex:") {
		regex, err := regexp.Compile(strings.TrimPrefix(rule.Pattern, "regex:"))
		if err != nil {
			return fmt.Errorf("編譯正則表達式失敗: %v", err)
		}
		rule.regex = regex
	}

	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()
	ae.rules = append(ae.rules, rule)

	ae.logger.Infof("已添加安全規則: %s", rule.Name)
	return nil
}

// RemoveRule 移除安全規則
func (ae *AnalysisEngine) RemoveRule(ruleID string) error {
	ae.mutex.Lock()
	defer ae.mutex.Unlock()

	for i, rule := range ae.rules {
		if rule.ID == ruleID {
			ae.rules = append(ae.rules[:i], ae.rules[i+1:]...)
			ae.logger.Infof("已移除安全規則: %s", rule.Name)
			return nil
		}
	}

	return fmt.Errorf("找不到規則ID: %s", ruleID)
}

// addToBlacklist 添加到黑名單
func (ae *AnalysisEngine) addToBlacklist(ip string, duration time.Duration) {
	ae.blacklist[ip] = time.Now().Add(duration)
	ae.logger.Warnf("已將IP %s 添加到黑名單，持續時間: %v", ip, duration)
}

// AddToWhitelist 添加到白名單
func (ae *AnalysisEngine) AddToWhitelist(ip string) {
	ae.mutex.Lock()
	defer ae.mutex.Unlock()

	ae.whitelist[ip] = true
	ae.logger.Infof("已將IP %s 添加到白名單", ip)
}

// RemoveFromWhitelist 從白名單移除
func (ae *AnalysisEngine) RemoveFromWhitelist(ip string) {
	ae.mutex.Lock()
	defer ae.mutex.Unlock()

	delete(ae.whitelist, ip)
	ae.logger.Infof("已將IP %s 從白名單移除", ip)
}

// GetThreatInfo 取得威脅資訊
func (ae *AnalysisEngine) GetThreatInfo() map[string]*ThreatInfo {
	ae.mutex.RLock()
	defer ae.mutex.RUnlock()

	// 複製威脅快取
	result := make(map[string]*ThreatInfo)
	for k, v := range ae.threatCache {
		result[k] = v
	}

	return result
}

// GetRules 取得所有規則
func (ae *AnalysisEngine) GetRules() []*SecurityRule {
	ae.mutex.RLock()
	defer ae.mutex.RUnlock()

	// 複製規則清單
	result := make([]*SecurityRule, len(ae.rules))
	copy(result, ae.rules)

	return result
}

// startThreatCacheCleanup 啟動威脅快取清理
func (ae *AnalysisEngine) startThreatCacheCleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ae.cleanupThreatCache()
		case <-ae.stopChan:
			return
		}
	}
}

// cleanupThreatCache 清理威脅快取
func (ae *AnalysisEngine) cleanupThreatCache() {
	ae.mutex.Lock()
	defer ae.mutex.Unlock()

	now := time.Now()
	for ip, threat := range ae.threatCache {
		// 移除1小時前的威脅記錄
		if now.Sub(threat.LastSeen) > 1*time.Hour {
			delete(ae.threatCache, ip)
		}
	}

	ae.logger.Debugf("威脅快取清理完成，當前記錄數: %d", len(ae.threatCache))
}

// startBlacklistCleanup 啟動黑名單清理
func (ae *AnalysisEngine) startBlacklistCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ae.cleanupBlacklist()
		case <-ae.stopChan:
			return
		}
	}
}

// cleanupBlacklist 清理黑名單
func (ae *AnalysisEngine) cleanupBlacklist() {
	ae.mutex.Lock()
	defer ae.mutex.Unlock()

	now := time.Now()
	for ip, expiry := range ae.blacklist {
		if now.After(expiry) {
			delete(ae.blacklist, ip)
			ae.logger.Infof("已從黑名單移除過期IP: %s", ip)
		}
	}
}

// loadDefaultRules 載入預設規則
func (ae *AnalysisEngine) loadDefaultRules() {
	defaultRules := []*SecurityRule{
		{
			ID:          "rule_001",
			Name:        "阻擋惡意IP",
			Description: "阻擋已知惡意IP位址",
			Type:        "ip",
			Pattern:     "regex:^(10\\.0\\.0\\.|192\\.168\\.1\\.)",
			Action:      "block",
			Severity:    "high",
			Enabled:     true,
		},
		{
			ID:          "rule_002",
			Name:        "SSH暴力破解偵測",
			Description: "偵測SSH暴力破解嘗試",
			Type:        "behavior",
			Pattern:     "brute_force",
			Action:      "alert",
			Severity:    "medium",
			Enabled:     true,
		},
		{
			ID:          "rule_003",
			Name:        "連接埠掃描偵測",
			Description: "偵測連接埠掃描行為",
			Type:        "behavior",
			Pattern:     "port_scan",
			Action:      "alert",
			Severity:    "medium",
			Enabled:     true,
		},
		{
			ID:          "rule_004",
			Name:        "SQL注入偵測",
			Description: "偵測SQL注入攻擊模式",
			Type:        "pattern",
			Pattern:     "regex:(?i)(union|select|insert|delete|update|drop|exec|script)",
			Action:      "block",
			Severity:    "high",
			Enabled:     true,
		},
	}

	for _, rule := range defaultRules {
		if err := ae.AddRule(rule); err != nil {
			ae.logger.Errorf("載入預設規則失敗: %v", err)
		}
	}

	ae.logger.Infof("已載入 %d 個預設安全規則", len(defaultRules))
}
