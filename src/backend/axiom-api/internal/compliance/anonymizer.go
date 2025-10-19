package compliance

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// AnonymizationMethod 匿名化方法
type AnonymizationMethod string

const (
	MethodMask          AnonymizationMethod = "mask"          // 遮罩
	MethodHash          AnonymizationMethod = "hash"          // 雜湊
	MethodGeneralize    AnonymizationMethod = "generalize"    // 泛化
	MethodPseudonymize  AnonymizationMethod = "pseudonymize"  // 假名化（可逆）
)

// Anonymizer 匿名化器
type Anonymizer struct {
	salt          string
	encryptionKey []byte
	piiDetector   *PIIDetector
}

// AnonymizationResult 匿名化結果
type AnonymizationResult struct {
	OriginalText  string                 `json:"original_text,omitempty"`
	AnonymizedText string                `json:"anonymized_text"`
	Method        AnonymizationMethod    `json:"method"`
	PIIDetected   []PIIMatch             `json:"pii_detected"`
	Reversible    bool                   `json:"reversible"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// NewAnonymizer 創建匿名化器
func NewAnonymizer(salt string, encryptionKey []byte) *Anonymizer {
	return &Anonymizer{
		salt:          salt,
		encryptionKey: encryptionKey,
		piiDetector:   NewPIIDetector(),
	}
}

// Anonymize 匿名化文本
func (a *Anonymizer) Anonymize(ctx context.Context, text string, method AnonymizationMethod) (*AnonymizationResult, error) {
	// 先檢測 PII
	detection := a.piiDetector.DetectPII(ctx, text)
	
	var anonymizedText string
	var err error
	var reversible bool
	
	switch method {
	case MethodMask:
		anonymizedText = a.maskText(text, detection.Matches)
		reversible = false
		
	case MethodHash:
		anonymizedText = a.hashText(text, detection.Matches)
		reversible = false
		
	case MethodGeneralize:
		anonymizedText = a.generalizeText(text, detection.Matches)
		reversible = false
		
	case MethodPseudonymize:
		anonymizedText, err = a.pseudonymizeText(text, detection.Matches)
		reversible = true
		if err != nil {
			return nil, fmt.Errorf("pseudonymization failed: %w", err)
		}
		
	default:
		return nil, fmt.Errorf("unknown anonymization method: %s", method)
	}
	
	return &AnonymizationResult{
		AnonymizedText: anonymizedText,
		Method:         method,
		PIIDetected:    detection.Matches,
		Reversible:     reversible,
		Metadata: map[string]interface{}{
			"pii_count":  len(detection.Matches),
			"risk_level": detection.RiskLevel,
		},
	}, nil
}

// maskText 遮罩文本
func (a *Anonymizer) maskText(text string, matches []PIIMatch) string {
	result := text
	offset := 0
	
	for _, match := range matches {
		adjustedStart := match.Start + offset
		adjustedEnd := match.End + offset
		
		before := result[:adjustedStart]
		after := result[adjustedEnd:]
		
		result = before + match.Masked + after
		offset += len(match.Masked) - (match.End - match.Start)
	}
	
	return result
}

// hashText 雜湊文本
func (a *Anonymizer) hashText(text string, matches []PIIMatch) string {
	result := text
	offset := 0
	
	for _, match := range matches {
		adjustedStart := match.Start + offset
		adjustedEnd := match.End + offset
		
		// 使用 SHA-256 雜湊
		hash := sha256.New()
		hash.Write([]byte(match.Value + a.salt))
		hashedValue := "REDACTED_" + hex.EncodeToString(hash.Sum(nil))[:16]
		
		before := result[:adjustedStart]
		after := result[adjustedEnd:]
		
		result = before + hashedValue + after
		offset += len(hashedValue) - (match.End - match.Start)
	}
	
	return result
}

// generalizeText 泛化文本
func (a *Anonymizer) generalizeText(text string, matches []PIIMatch) string {
	result := text
	offset := 0
	
	for _, match := range matches {
		adjustedStart := match.Start + offset
		adjustedEnd := match.End + offset
		
		var generalizedValue string
		
		switch match.Type {
		case PIITypeIPAddress:
			// 192.168.1.100 → 192.168.0.0/16
			generalizedValue = "*.*.0.0/16"
		case PIITypeEmail:
			// john@example.com → *@example.com
			generalizedValue = "*@domain.com"
		default:
			generalizedValue = "[" + string(match.Type) + "]"
		}
		
		before := result[:adjustedStart]
		after := result[adjustedEnd:]
		
		result = before + generalizedValue + after
		offset += len(generalizedValue) - (match.End - match.Start)
	}
	
	return result
}

// pseudonymizeText 假名化文本（可逆）
func (a *Anonymizer) pseudonymizeText(text string, matches []PIIMatch) (string, error) {
	result := text
	offset := 0
	
	for _, match := range matches {
		adjustedStart := match.Start + offset
		adjustedEnd := match.End + offset
		
		// 使用 AES 加密
		encrypted, err := a.encrypt(match.Value)
		if err != nil {
			return "", err
		}
		
		token := "TOKEN_" + encrypted[:16]
		
		before := result[:adjustedStart]
		after := result[adjustedEnd:]
		
		result = before + token + after
		offset += len(token) - (match.End - match.Start)
	}
	
	return result, nil
}

// Depseudonymize 反假名化（還原）
func (a *Anonymizer) Depseudonymize(ctx context.Context, token string) (string, error) {
	// 移除 TOKEN_ 前綴
	if len(token) < 23 {
		return "", fmt.Errorf("invalid token format")
	}
	
	encryptedData := token[6:] // 去除 "TOKEN_"
	
	return a.decrypt(encryptedData)
}

// encrypt 加密
func (a *Anonymizer) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(a.encryptionKey)
	if err != nil {
		return "", err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// decrypt 解密
func (a *Anonymizer) decrypt(ciphertext string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	
	block, err := aes.NewCipher(a.encryptionKey)
	if err != nil {
		return "", err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	
	nonce := data[:nonceSize]
	ciphertextBytes := data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}

// AnonymizeBatch 批量匿名化
func (a *Anonymizer) AnonymizeBatch(ctx context.Context, texts []string, method AnonymizationMethod) ([]*AnonymizationResult, error) {
	results := make([]*AnonymizationResult, len(texts))
	
	for i, text := range texts {
		result, err := a.Anonymize(ctx, text, method)
		if err != nil {
			return nil, fmt.Errorf("failed to anonymize text %d: %w", i, err)
		}
		results[i] = result
	}
	
	return results, nil
}

