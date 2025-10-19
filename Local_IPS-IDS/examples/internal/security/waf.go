package security

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// WAF implements Web Application Firewall functionality
type WAF struct {
	rules      []WAFRule
	rulesMu    sync.RWMutex
	logger     *logrus.Logger
	enabled    bool
	blockMode  bool // true: block, false: monitor only
}

// WAFRule represents a WAF rule
type WAFRule struct {
	ID          string
	Name        string
	Description string
	Category    string
	Severity    string
	Pattern     *regexp.Regexp
	Action      string // "block", "log", "challenge"
	Enabled     bool
}

// WAFResult represents the result of WAF inspection
type WAFResult struct {
	Blocked      bool
	MatchedRules []string
	Reasons      []string
	Score        int
	Action       string
}

// NewWAF creates a new WAF instance
func NewWAF(logger *logrus.Logger) *WAF {
	if logger == nil {
		logger = logrus.New()
	}

	waf := &WAF{
		rules:     make([]WAFRule, 0),
		logger:    logger,
		enabled:   true,
		blockMode: true,
	}

	// 載入預設規則
	waf.loadDefaultRules()

	return waf
}

// InspectRequest inspects an HTTP request
func (w *WAF) InspectRequest(req *http.Request) *WAFResult {
	if !w.enabled {
		return &WAFResult{Blocked: false}
	}

	result := &WAFResult{
		Blocked:      false,
		MatchedRules: make([]string, 0),
		Reasons:      make([]string, 0),
		Score:        0,
	}

	w.rulesMu.RLock()
	defer w.rulesMu.RUnlock()

	// 檢查 URL
	w.inspectURL(req.URL.String(), result)

	// 檢查 Headers
	w.inspectHeaders(req.Header, result)

	// 檢查 Query Parameters
	w.inspectQueryParams(req.URL.Query(), result)

	// 檢查 User-Agent
	w.inspectUserAgent(req.UserAgent(), result)

	// 根據分數決定動作
	if result.Score >= 10 {
		result.Blocked = w.blockMode
		result.Action = "block"
		w.logger.Warnf("WAF blocked request: %s %s (score: %d, rules: %v)", 
			req.Method, req.URL.Path, result.Score, result.MatchedRules)
	} else if result.Score >= 5 {
		result.Action = "challenge"
		w.logger.Infof("WAF challenge request: %s %s (score: %d)", 
			req.Method, req.URL.Path, result.Score)
	} else if result.Score > 0 {
		result.Action = "log"
		w.logger.Debugf("WAF suspicious request: %s %s (score: %d)", 
			req.Method, req.URL.Path, result.Score)
	}

	return result
}

// inspectURL checks URL for malicious patterns
func (w *WAF) inspectURL(url string, result *WAFResult) {
	for _, rule := range w.rules {
		if !rule.Enabled || rule.Category != "url" {
			continue
		}

		if rule.Pattern.MatchString(url) {
			result.MatchedRules = append(result.MatchedRules, rule.ID)
			result.Reasons = append(result.Reasons, rule.Name)
			result.Score += w.getSeverityScore(rule.Severity)
		}
	}
}

// inspectHeaders checks HTTP headers
func (w *WAF) inspectHeaders(headers http.Header, result *WAFResult) {
	for key, values := range headers {
		headerStr := fmt.Sprintf("%s: %s", key, strings.Join(values, ", "))
		
		for _, rule := range w.rules {
			if !rule.Enabled || rule.Category != "header" {
				continue
			}

			if rule.Pattern.MatchString(headerStr) {
				result.MatchedRules = append(result.MatchedRules, rule.ID)
				result.Reasons = append(result.Reasons, rule.Name)
				result.Score += w.getSeverityScore(rule.Severity)
			}
		}
	}
}

// inspectQueryParams checks query parameters
func (w *WAF) inspectQueryParams(params map[string][]string, result *WAFResult) {
	for key, values := range params {
		paramStr := fmt.Sprintf("%s=%s", key, strings.Join(values, ","))
		
		for _, rule := range w.rules {
			if !rule.Enabled || rule.Category != "parameter" {
				continue
			}

			if rule.Pattern.MatchString(paramStr) {
				result.MatchedRules = append(result.MatchedRules, rule.ID)
				result.Reasons = append(result.Reasons, rule.Name)
				result.Score += w.getSeverityScore(rule.Severity)
			}
		}
	}
}

// inspectUserAgent checks User-Agent header
func (w *WAF) inspectUserAgent(userAgent string, result *WAFResult) {
	for _, rule := range w.rules {
		if !rule.Enabled || rule.Category != "user-agent" {
			continue
		}

		if rule.Pattern.MatchString(userAgent) {
			result.MatchedRules = append(result.MatchedRules, rule.ID)
			result.Reasons = append(result.Reasons, rule.Name)
			result.Score += w.getSeverityScore(rule.Severity)
		}
	}
}

// getSeverityScore returns score based on severity
func (w *WAF) getSeverityScore(severity string) int {
	switch severity {
	case "critical":
		return 10
	case "high":
		return 7
	case "medium":
		return 4
	case "low":
		return 2
	default:
		return 1
	}
}

// loadDefaultRules loads default WAF rules
func (w *WAF) loadDefaultRules() {
	rules := []WAFRule{
		// SQL Injection
		{
			ID:          "WAF-001",
			Name:        "SQL Injection",
			Description: "Detects SQL injection attempts",
			Category:    "parameter",
			Severity:    "critical",
			Pattern:     regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute|script|javascript|<script)`),
			Action:      "block",
			Enabled:     true,
		},
		// XSS
		{
			ID:          "WAF-002",
			Name:        "Cross-Site Scripting (XSS)",
			Description: "Detects XSS attempts",
			Category:    "parameter",
			Severity:    "high",
			Pattern:     regexp.MustCompile(`(?i)(<script|javascript:|onerror=|onload=|<iframe|<object|<embed)`),
			Action:      "block",
			Enabled:     true,
		},
		// Path Traversal
		{
			ID:          "WAF-003",
			Name:        "Path Traversal",
			Description: "Detects directory traversal attempts",
			Category:    "url",
			Severity:    "high",
			Pattern:     regexp.MustCompile(`(\.\.\/|\.\.\\|%2e%2e%2f|%2e%2e\/|\.\.%2f)`),
			Action:      "block",
			Enabled:     true,
		},
		// Command Injection
		{
			ID:          "WAF-004",
			Name:        "Command Injection",
			Description: "Detects OS command injection",
			Category:    "parameter",
			Severity:    "critical",
			Pattern:     regexp.MustCompile(`(?i)(;|\||&|`+"`"+`|\$\(|>\||<\||&&|\|\|)`),
			Action:      "block",
			Enabled:     true,
		},
		// LFI/RFI
		{
			ID:          "WAF-005",
			Name:        "Local/Remote File Inclusion",
			Description: "Detects file inclusion attempts",
			Category:    "url",
			Severity:    "high",
			Pattern:     regexp.MustCompile(`(?i)(file://|php://|data://|expect://|zip://)`),
			Action:      "block",
			Enabled:     true,
		},
		// Malicious User-Agent
		{
			ID:          "WAF-006",
			Name:        "Malicious User-Agent",
			Description: "Detects known malicious user agents",
			Category:    "user-agent",
			Severity:    "medium",
			Pattern:     regexp.MustCompile(`(?i)(nikto|sqlmap|nmap|masscan|acunetix|nessus|openvas|metasploit)`),
			Action:      "block",
			Enabled:     true,
		},
		// Scanner Detection
		{
			ID:          "WAF-007",
			Name:        "Security Scanner",
			Description: "Detects security scanning tools",
			Category:    "user-agent",
			Severity:    "medium",
			Pattern:     regexp.MustCompile(`(?i)(burp|zap|w3af|dirbuster|gobuster|wfuzz|ffuf)`),
			Action:      "block",
			Enabled:     true,
		},
		// Suspicious Headers
		{
			ID:          "WAF-008",
			Name:        "Suspicious Headers",
			Description: "Detects suspicious HTTP headers",
			Category:    "header",
			Severity:    "low",
			Pattern:     regexp.MustCompile(`(?i)(X-Forwarded-For: 127\.0\.0\.1|X-Real-IP: localhost)`),
			Action:      "log",
			Enabled:     true,
		},
	}

	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()

	w.rules = append(w.rules, rules...)
	w.logger.Infof("Loaded %d WAF rules", len(w.rules))
}

// AddRule adds a custom WAF rule
func (w *WAF) AddRule(rule WAFRule) error {
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()

	// 檢查 ID 是否已存在
	for _, r := range w.rules {
		if r.ID == rule.ID {
			return fmt.Errorf("rule ID %s already exists", rule.ID)
		}
	}

	w.rules = append(w.rules, rule)
	w.logger.Infof("Added WAF rule: %s (%s)", rule.ID, rule.Name)
	return nil
}

// RemoveRule removes a WAF rule
func (w *WAF) RemoveRule(ruleID string) error {
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()

	for i, rule := range w.rules {
		if rule.ID == ruleID {
			w.rules = append(w.rules[:i], w.rules[i+1:]...)
			w.logger.Infof("Removed WAF rule: %s", ruleID)
			return nil
		}
	}

	return fmt.Errorf("rule ID %s not found", ruleID)
}

// EnableRule enables a WAF rule
func (w *WAF) EnableRule(ruleID string) error {
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()

	for i, rule := range w.rules {
		if rule.ID == ruleID {
			w.rules[i].Enabled = true
			w.logger.Infof("Enabled WAF rule: %s", ruleID)
			return nil
		}
	}

	return fmt.Errorf("rule ID %s not found", ruleID)
}

// DisableRule disables a WAF rule
func (w *WAF) DisableRule(ruleID string) error {
	w.rulesMu.Lock()
	defer w.rulesMu.Unlock()

	for i, rule := range w.rules {
		if rule.ID == ruleID {
			w.rules[i].Enabled = false
			w.logger.Infof("Disabled WAF rule: %s", ruleID)
			return nil
		}
	}

	return fmt.Errorf("rule ID %s not found", ruleID)
}

// SetBlockMode sets the WAF block mode
func (w *WAF) SetBlockMode(block bool) {
	w.blockMode = block
	mode := "monitor"
	if block {
		mode = "block"
	}
	w.logger.Infof("WAF mode set to: %s", mode)
}

// Enable enables the WAF
func (w *WAF) Enable() {
	w.enabled = true
	w.logger.Info("WAF enabled")
}

// Disable disables the WAF
func (w *WAF) Disable() {
	w.enabled = false
	w.logger.Info("WAF disabled")
}

// GetRules returns all WAF rules
func (w *WAF) GetRules() []WAFRule {
	w.rulesMu.RLock()
	defer w.rulesMu.RUnlock()

	rules := make([]WAFRule, len(w.rules))
	copy(rules, w.rules)
	return rules
}

