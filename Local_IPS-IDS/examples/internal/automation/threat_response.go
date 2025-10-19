package automation

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ThreatResponseSystem implements automated threat response
type ThreatResponseSystem struct {
	n8nClient      *N8NClient
	responseRules  []ResponseRule
	rulesMu        sync.RWMutex
	logger         *logrus.Logger
	enabled        bool
	dryRun         bool
}

// ResponseRule defines an automated response rule
type ResponseRule struct {
	ID          string
	Name        string
	ThreatType  string
	Severity    string
	Actions     []ResponseAction
	Conditions  []Condition
	Enabled     bool
	Priority    int
}

// ResponseAction defines an action to take
type ResponseAction struct {
	Type        string // "block_ip", "isolate_host", "kill_process", "notify", etc.
	Params      map[string]interface{}
	Timeout     time.Duration
	RetryCount  int
}

// Condition defines a condition for rule matching
type Condition struct {
	Field    string
	Operator string // "equals", "contains", "greater_than", etc.
	Value    interface{}
}

// ResponseResult represents the result of automated response
type ResponseResult struct {
	RuleID         string
	ActionsExecuted []string
	Success        bool
	Errors         []string
	Duration       time.Duration
	Timestamp      time.Time
}

// NewThreatResponseSystem creates a new threat response system
func NewThreatResponseSystem(n8nClient *N8NClient, logger *logrus.Logger) *ThreatResponseSystem {
	if logger == nil {
		logger = logrus.New()
	}

	trs := &ThreatResponseSystem{
		n8nClient:     n8nClient,
		responseRules: make([]ResponseRule, 0),
		logger:        logger,
		enabled:       true,
		dryRun:        false,
	}

	// 載入預設響應規則
	trs.loadDefaultRules()

	return trs
}

// ProcessThreat processes a threat and executes automated response
func (trs *ThreatResponseSystem) ProcessThreat(ctx context.Context, event *ThreatEvent) (*ResponseResult, error) {
	if !trs.enabled {
		return nil, fmt.Errorf("threat response system is disabled")
	}

	startTime := time.Now()
	result := &ResponseResult{
		ActionsExecuted: make([]string, 0),
		Errors:          make([]string, 0),
		Timestamp:       startTime,
	}

	// 找到匹配的規則
	rule := trs.findMatchingRule(event)
	if rule == nil {
		trs.logger.Debugf("No matching rule found for threat: %s", event.ID)
		return result, nil
	}

	result.RuleID = rule.ID
	trs.logger.Infof("Processing threat %s with rule %s", event.ID, rule.ID)

	// 執行響應動作
	for _, action := range rule.Actions {
		if err := trs.executeAction(ctx, event, action); err != nil {
			errMsg := fmt.Sprintf("Action %s failed: %v", action.Type, err)
			result.Errors = append(result.Errors, errMsg)
			trs.logger.Errorf(errMsg)
			continue
		}

		result.ActionsExecuted = append(result.ActionsExecuted, action.Type)
		trs.logger.Infof("Action executed: %s", action.Type)
	}

	result.Success = len(result.Errors) == 0
	result.Duration = time.Since(startTime)

	// 記錄響應結果
	trs.logResponse(event, result)

	return result, nil
}

// findMatchingRule finds the first matching rule for a threat
func (trs *ThreatResponseSystem) findMatchingRule(event *ThreatEvent) *ResponseRule {
	trs.rulesMu.RLock()
	defer trs.rulesMu.RUnlock()

	var matchedRule *ResponseRule
	highestPriority := -1

	for i := range trs.responseRules {
		rule := &trs.responseRules[i]
		
		if !rule.Enabled {
			continue
		}

		// 檢查威脅類型和嚴重程度
		if rule.ThreatType != "" && rule.ThreatType != event.Type {
			continue
		}

		if rule.Severity != "" && rule.Severity != event.Severity {
			continue
		}

		// 檢查額外條件
		if !trs.checkConditions(event, rule.Conditions) {
			continue
		}

		// 選擇優先級最高的規則
		if rule.Priority > highestPriority {
			matchedRule = rule
			highestPriority = rule.Priority
		}
	}

	return matchedRule
}

// checkConditions checks if all conditions are met
func (trs *ThreatResponseSystem) checkConditions(event *ThreatEvent, conditions []Condition) bool {
	for _, cond := range conditions {
		if !trs.evaluateCondition(event, cond) {
			return false
		}
	}
	return true
}

// evaluateCondition evaluates a single condition
func (trs *ThreatResponseSystem) evaluateCondition(event *ThreatEvent, cond Condition) bool {
	// 簡化版條件評估
	// 實際實現需要更複雜的邏輯
	return true
}

// executeAction executes a response action
func (trs *ThreatResponseSystem) executeAction(ctx context.Context, event *ThreatEvent, action ResponseAction) error {
	if trs.dryRun {
		trs.logger.Infof("[DRY-RUN] Would execute action: %s", action.Type)
		return nil
	}

	// 設置超時
	actionCtx := ctx
	if action.Timeout > 0 {
		var cancel context.CancelFunc
		actionCtx, cancel = context.WithTimeout(ctx, action.Timeout)
		defer cancel()
	}

	switch action.Type {
	case "block_ip":
		return trs.blockIP(actionCtx, event.Source)
	
	case "isolate_host":
		return trs.isolateHost(actionCtx, event.Source)
	
	case "kill_process":
		return trs.killProcess(actionCtx, action.Params)
	
	case "notify_soc":
		return trs.notifySOC(actionCtx, event)
	
	case "create_incident":
		return trs.createIncident(actionCtx, event)
	
	case "quarantine_file":
		return trs.quarantineFile(actionCtx, action.Params)
	
	case "update_firewall":
		return trs.updateFirewall(actionCtx, action.Params)
	
	case "collect_forensics":
		return trs.collectForensics(actionCtx, event)
	
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// Action implementations

func (trs *ThreatResponseSystem) blockIP(ctx context.Context, ip string) error {
	params := map[string]interface{}{
		"ip":     ip,
		"action": "block",
		"duration": "24h",
	}

	return trs.n8nClient.ExecuteRemediationAction(ctx, "block_ip", params)
}

func (trs *ThreatResponseSystem) isolateHost(ctx context.Context, host string) error {
	params := map[string]interface{}{
		"host":   host,
		"action": "isolate",
	}

	return trs.n8nClient.ExecuteRemediationAction(ctx, "isolate_host", params)
}

func (trs *ThreatResponseSystem) killProcess(ctx context.Context, params map[string]interface{}) error {
	return trs.n8nClient.ExecuteRemediationAction(ctx, "kill_process", params)
}

func (trs *ThreatResponseSystem) notifySOC(ctx context.Context, event *ThreatEvent) error {
	message := fmt.Sprintf("🚨 Critical Threat Detected: %s\nSeverity: %s\nSource: %s\nDescription: %s",
		event.Type, event.Severity, event.Source, event.Description)

	metadata := map[string]interface{}{
		"event_id":  event.ID,
		"timestamp": event.Timestamp,
	}

	return trs.n8nClient.SendNotification(ctx, "soc-alerts", message, metadata)
}

func (trs *ThreatResponseSystem) createIncident(ctx context.Context, event *ThreatEvent) error {
	incident := map[string]interface{}{
		"title":       fmt.Sprintf("Security Incident: %s", event.Type),
		"description": event.Description,
		"severity":    event.Severity,
		"source":      event.Source,
		"timestamp":   event.Timestamp,
		"metadata":    event.Metadata,
	}

	_, err := trs.n8nClient.CreateIncident(ctx, incident)
	return err
}

func (trs *ThreatResponseSystem) quarantineFile(ctx context.Context, params map[string]interface{}) error {
	return trs.n8nClient.ExecuteRemediationAction(ctx, "quarantine_file", params)
}

func (trs *ThreatResponseSystem) updateFirewall(ctx context.Context, params map[string]interface{}) error {
	return trs.n8nClient.ExecuteRemediationAction(ctx, "update_firewall", params)
}

func (trs *ThreatResponseSystem) collectForensics(ctx context.Context, event *ThreatEvent) error {
	params := map[string]interface{}{
		"event_id":  event.ID,
		"source":    event.Source,
		"timestamp": event.Timestamp,
	}

	return trs.n8nClient.ExecuteRemediationAction(ctx, "collect_forensics", params)
}

// loadDefaultRules loads default response rules
func (trs *ThreatResponseSystem) loadDefaultRules() {
	rules := []ResponseRule{
		{
			ID:         "RULE-001",
			Name:       "Critical Malware Response",
			ThreatType: "malware",
			Severity:   "critical",
			Priority:   100,
			Enabled:    true,
			Actions: []ResponseAction{
				{Type: "isolate_host", Timeout: 30 * time.Second},
				{Type: "notify_soc", Timeout: 10 * time.Second},
				{Type: "create_incident", Timeout: 20 * time.Second},
				{Type: "collect_forensics", Timeout: 60 * time.Second},
			},
		},
		{
			ID:         "RULE-002",
			Name:       "Critical Intrusion Response",
			ThreatType: "intrusion",
			Severity:   "critical",
			Priority:   100,
			Enabled:    true,
			Actions: []ResponseAction{
				{Type: "block_ip", Timeout: 15 * time.Second},
				{Type: "notify_soc", Timeout: 10 * time.Second},
				{Type: "create_incident", Timeout: 20 * time.Second},
				{Type: "update_firewall", Timeout: 30 * time.Second},
			},
		},
		{
			ID:         "RULE-003",
			Name:       "DDoS Mitigation",
			ThreatType: "ddos",
			Severity:   "high",
			Priority:   90,
			Enabled:    true,
			Actions: []ResponseAction{
				{Type: "block_ip", Timeout: 15 * time.Second},
				{Type: "update_firewall", Timeout: 30 * time.Second},
				{Type: "notify_soc", Timeout: 10 * time.Second},
			},
		},
	}

	trs.rulesMu.Lock()
	defer trs.rulesMu.Unlock()

	trs.responseRules = append(trs.responseRules, rules...)
	trs.logger.Infof("Loaded %d response rules", len(trs.responseRules))
}

// logResponse logs the response result
func (trs *ThreatResponseSystem) logResponse(event *ThreatEvent, result *ResponseResult) {
	trs.logger.WithFields(logrus.Fields{
		"event_id":         event.ID,
		"rule_id":          result.RuleID,
		"actions_executed": result.ActionsExecuted,
		"success":          result.Success,
		"duration":         result.Duration,
		"errors":           result.Errors,
	}).Info("Threat response completed")
}

// AddRule adds a custom response rule
func (trs *ThreatResponseSystem) AddRule(rule ResponseRule) {
	trs.rulesMu.Lock()
	defer trs.rulesMu.Unlock()

	trs.responseRules = append(trs.responseRules, rule)
	trs.logger.Infof("Added response rule: %s", rule.ID)
}

// Enable enables the threat response system
func (trs *ThreatResponseSystem) Enable() {
	trs.enabled = true
	trs.logger.Info("Threat response system enabled")
}

// Disable disables the threat response system
func (trs *ThreatResponseSystem) Disable() {
	trs.enabled = false
	trs.logger.Info("Threat response system disabled")
}

// SetDryRun sets dry-run mode
func (trs *ThreatResponseSystem) SetDryRun(dryRun bool) {
	trs.dryRun = dryRun
	mode := "normal"
	if dryRun {
		mode = "dry-run"
	}
	trs.logger.Infof("Threat response mode set to: %s", mode)
}

