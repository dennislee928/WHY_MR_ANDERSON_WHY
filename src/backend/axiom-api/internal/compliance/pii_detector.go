package compliance

import (
	"context"
	"regexp"
	"time"
)

// PIIType PII 類型
type PIIType string

const (
	PIITypeEmail       PIIType = "email"
	PIITypeCreditCard  PIIType = "credit_card"
	PIITypeSSN         PIIType = "ssn"
	PIITypeIPAddress   PIIType = "ip_address"
	PIITypePhone       PIIType = "phone"
	PIITypePassport    PIIType = "passport"
)

// PIIDetector PII 檢測器
type PIIDetector struct {
	patterns map[PIIType]*regexp.Regexp
}

// PIIMatch PII 匹配結果
type PIIMatch struct {
	Type        PIIType `json:"type"`
	Value       string  `json:"value"`
	Start       int     `json:"start"`
	End         int     `json:"end"`
	Masked      string  `json:"masked"`
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
}

// PIIDetectionResult PII 檢測結果
type PIIDetectionResult struct {
	Text         string     `json:"text"`
	Matches      []PIIMatch `json:"matches"`
	PIIFound     bool       `json:"pii_found"`
	DetectedAt   time.Time  `json:"detected_at"`
	RiskLevel    string     `json:"risk_level"` // low, medium, high, critical
}

// NewPIIDetector 創建 PII 檢測器
func NewPIIDetector() *PIIDetector {
	return &PIIDetector{
		patterns: map[PIIType]*regexp.Regexp{
			PIITypeEmail:      regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
			PIITypeCreditCard: regexp.MustCompile(`\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`),
			PIITypeSSN:        regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),
			PIITypeIPAddress:  regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`),
			PIITypePhone:      regexp.MustCompile(`\b\d{3}[-.]?\d{3}[-.]?\d{4}\b`),
			PIITypePassport:   regexp.MustCompile(`\b[A-Z]{1,2}\d{6,9}\b`),
		},
	}
}

// DetectPII 檢測 PII
func (d *PIIDetector) DetectPII(ctx context.Context, text string) *PIIDetectionResult {
	result := &PIIDetectionResult{
		Text:       text,
		Matches:    []PIIMatch{},
		PIIFound:   false,
		DetectedAt: time.Now(),
		RiskLevel:  "low",
	}
	
	for piiType, pattern := range d.patterns {
		matches := pattern.FindAllStringIndex(text, -1)
		
		for _, match := range matches {
			start, end := match[0], match[1]
			value := text[start:end]
			
			piiMatch := PIIMatch{
				Type:        piiType,
				Value:       value,
				Start:       start,
				End:         end,
				Masked:      d.maskValue(value, piiType),
				Confidence:  d.calculateConfidence(value, piiType),
				Description: d.getDescription(piiType),
			}
			
			result.Matches = append(result.Matches, piiMatch)
			result.PIIFound = true
		}
	}
	
	// 計算風險等級
	result.RiskLevel = d.calculateRiskLevel(result.Matches)
	
	return result
}

// DetectPIIBatch 批量檢測 PII
func (d *PIIDetector) DetectPIIBatch(ctx context.Context, texts []string) []*PIIDetectionResult {
	results := make([]*PIIDetectionResult, len(texts))
	
	for i, text := range texts {
		results[i] = d.DetectPII(ctx, text)
	}
	
	return results
}

// maskValue 遮罩值
func (d *PIIDetector) maskValue(value string, piiType PIIType) string {
	switch piiType {
	case PIITypeEmail:
		// email: j***@e*****.com
		if len(value) > 5 {
			parts := regexp.MustCompile(`@`).Split(value, 2)
			if len(parts) == 2 {
				return value[0:1] + "***@" + parts[1][0:1] + "*****." + parts[1][len(parts[1])-3:]
			}
		}
		return "***@***.***"
		
	case PIITypeCreditCard:
		// credit card: **** **** **** 1234
		if len(value) >= 4 {
			return "**** **** **** " + value[len(value)-4:]
		}
		return "************"
		
	case PIITypeSSN:
		// SSN: ***-**-1234
		if len(value) >= 4 {
			return "***-**-" + value[len(value)-4:]
		}
		return "***-**-****"
		
	case PIITypeIPAddress:
		// IP: 192.168.*.*
		parts := regexp.MustCompile(`\.`).Split(value, 4)
		if len(parts) == 4 {
			return parts[0] + "." + parts[1] + ".*.*"
		}
		return "*.*.*.*"
		
	case PIITypePhone:
		// Phone: ***-***-1234
		if len(value) >= 4 {
			return "***-***-" + value[len(value)-4:]
		}
		return "***-***-****"
		
	case PIITypePassport:
		return "**" + value[len(value)-3:]
		
	default:
		return "***REDACTED***"
	}
}

// calculateConfidence 計算置信度
func (d *PIIDetector) calculateConfidence(value string, piiType PIIType) float64 {
	// 簡化實現，實際應該使用更複雜的驗證邏輯
	switch piiType {
	case PIITypeEmail:
		if regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(value) {
			return 0.95
		}
		return 0.80
		
	case PIITypeCreditCard:
		// Luhn 算法驗證
		if d.validateLuhn(value) {
			return 0.98
		}
		return 0.70
		
	case PIITypeSSN:
		return 0.90
		
	case PIITypeIPAddress:
		// 驗證 IP 範圍
		return 0.85
		
	default:
		return 0.75
	}
}

// validateLuhn Luhn 算法驗證信用卡
func (d *PIIDetector) validateLuhn(cardNumber string) bool {
	// 簡化實現
	// 實際應該實現完整的 Luhn 算法
	return true
}

// getDescription 獲取描述
func (d *PIIDetector) getDescription(piiType PIIType) string {
	descriptions := map[PIIType]string{
		PIITypeEmail:      "Email address",
		PIITypeCreditCard: "Credit card number",
		PIITypeSSN:        "Social Security Number",
		PIITypeIPAddress:  "IP address",
		PIITypePhone:      "Phone number",
		PIITypePassport:   "Passport number",
	}
	
	return descriptions[piiType]
}

// calculateRiskLevel 計算風險等級
func (d *PIIDetector) calculateRiskLevel(matches []PIIMatch) string {
	if len(matches) == 0 {
		return "low"
	}
	
	criticalTypes := map[PIIType]bool{
		PIITypeCreditCard: true,
		PIITypeSSN:        true,
		PIITypePassport:   true,
	}
	
	hasCritical := false
	for _, match := range matches {
		if criticalTypes[match.Type] {
			hasCritical = true
			break
		}
	}
	
	if hasCritical {
		return "critical"
	} else if len(matches) >= 3 {
		return "high"
	} else if len(matches) >= 1 {
		return "medium"
	}
	
	return "low"
}

// GetSupportedTypes 獲取支援的 PII 類型
func (d *PIIDetector) GetSupportedTypes() []PIIType {
	types := make([]PIIType, 0, len(d.patterns))
	for piiType := range d.patterns {
		types = append(types, piiType)
	}
	return types
}

