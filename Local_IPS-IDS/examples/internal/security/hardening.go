package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// SecurityConfig 安全配置
type SecurityConfig struct {
	RateLimitMaxRequests int
	RateLimitWindow      time.Duration
	SessionTimeout       time.Duration
	EncryptionKey        []byte
}

// SecurityHardening 安全強化模組
type SecurityHardening struct {
	logger         *logrus.Logger
	rateLimiter    *RateLimiter
	ipFilter       *IPFilter
	inputValidator *InputValidator
	sessionManager *SessionManager
	encryptionKey  []byte
}

// RateLimiter 速率限制器
type RateLimiter struct {
	requests map[string]*RequestCounter
	maxReqs  int
	window   time.Duration
	mutex    sync.RWMutex
}

// RequestCounter 請求計數器
type RequestCounter struct {
	count     int
	lastReset time.Time
}

// IPFilter IP過濾器
type IPFilter struct {
	whitelist map[string]bool
	blacklist map[string]time.Time
	mutex     sync.RWMutex
}

// InputValidator 輸入驗證器
type InputValidator struct {
	sqlInjectionPatterns []string
	xssPatterns          []string
	commandPatterns      []string
	compiledRegexes      []*regexp.Regexp
}

// SessionManager 會話管理器
type SessionManager struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
	timeout  time.Duration
}

// Session 會話資訊
type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	LastSeen  time.Time
	IP        string
	UserAgent string
}

// NewSecurityHardening 建立新的安全強化模組
func NewSecurityHardening(config *SecurityConfig, logger *logrus.Logger) (*SecurityHardening, error) {
	if config == nil {
		return nil, fmt.Errorf("配置不能為空")
	}

	if logger == nil {
		logger = logrus.New()
	}

	// 處理加密金鑰
	encKey := config.EncryptionKey
	if len(encKey) == 0 {
		encKey = make([]byte, 32)
		if _, err := rand.Read(encKey); err != nil {
			return nil, fmt.Errorf("生成加密金鑰失敗: %v", err)
		}
	}

	if len(encKey) != 32 {
		hash := sha256.Sum256(encKey)
		encKey = hash[:]
	}

	sh := &SecurityHardening{
		logger: logger,
		rateLimiter: &RateLimiter{
			requests: make(map[string]*RequestCounter),
			maxReqs:  config.RateLimitMaxRequests,
			window:   config.RateLimitWindow,
		},
		ipFilter: &IPFilter{
			whitelist: make(map[string]bool),
			blacklist: make(map[string]time.Time),
		},
		inputValidator: &InputValidator{
			sqlInjectionPatterns: []string{
				`(?i)(union|select|insert|delete|update|drop|exec|script)`,
				`(?i)(or|and)\s+\d+\s*=\s*\d+`,
				`(?i)'.*--`,
				`(?i)\/\*.*\*\/`,
			},
			xssPatterns: []string{
				`(?i)<script.*?>.*?</script>`,
				`(?i)javascript:`,
				`(?i)on\w+\s*=`,
				`(?i)<iframe.*?>`,
			},
			commandPatterns: []string{
				`(?i)(;|\||&|&&|\$\(|\)`,
				`(?i)(rm|del|format|shutdown|reboot)`,
				`(?i)\.\.\/`,
			},
		},
		sessionManager: &SessionManager{
			sessions: make(map[string]*Session),
			timeout:  config.SessionTimeout,
		},
		encryptionKey: encKey,
	}

	// 編譯正則表達式
	if err := sh.inputValidator.compileRegexes(); err != nil {
		return nil, fmt.Errorf("編譯安全規則失敗: %v", err)
	}

	// 啟動清理任務
	go sh.startCleanupTasks()

	return sh, nil
}

// compileRegexes 編譯正則表達式
func (iv *InputValidator) compileRegexes() error {
	allPatterns := append(iv.sqlInjectionPatterns, iv.xssPatterns...)
	allPatterns = append(allPatterns, iv.commandPatterns...)

	iv.compiledRegexes = make([]*regexp.Regexp, 0, len(allPatterns))
	for _, pattern := range allPatterns {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("編譯正則表達式失敗 %s: %v", pattern, err)
		}
		iv.compiledRegexes = append(iv.compiledRegexes, regex)
	}

	return nil
}

// startCleanupTasks 啟動清理任務
func (sh *SecurityHardening) startCleanupTasks() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sh.cleanupExpiredSessions()
		sh.cleanupExpiredIPs()
		sh.cleanupExpiredRequests()
	}
}

// cleanupExpiredSessions 清理過期會話
func (sh *SecurityHardening) cleanupExpiredSessions() {
	sh.sessionManager.mutex.Lock()
	defer sh.sessionManager.mutex.Unlock()

	now := time.Now()
	for id, session := range sh.sessionManager.sessions {
		if now.Sub(session.LastSeen) > sh.sessionManager.timeout {
			delete(sh.sessionManager.sessions, id)
		}
	}
}

// cleanupExpiredIPs 清理過期IP
func (sh *SecurityHardening) cleanupExpiredIPs() {
	sh.ipFilter.mutex.Lock()
	defer sh.ipFilter.mutex.Unlock()

	now := time.Now()
	for ip, expireTime := range sh.ipFilter.blacklist {
		if now.After(expireTime) {
			delete(sh.ipFilter.blacklist, ip)
		}
	}
}

// cleanupExpiredRequests 清理過期請求
func (sh *SecurityHardening) cleanupExpiredRequests() {
	sh.rateLimiter.mutex.Lock()
	defer sh.rateLimiter.mutex.Unlock()

	now := time.Now()
	for ip, counter := range sh.rateLimiter.requests {
		if now.Sub(counter.lastReset) > sh.rateLimiter.window {
			delete(sh.rateLimiter.requests, ip)
		}
	}
}

// ValidateInput 驗證輸入
func (sh *SecurityHardening) ValidateInput(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("輸入不能為空")
	}

	if len(input) > 10000 {
		return fmt.Errorf("輸入長度超過限制")
	}

	// 檢查危險模式
	for _, regex := range sh.inputValidator.compiledRegexes {
		if regex.MatchString(input) {
			return fmt.Errorf("輸入包含危險內容")
		}
	}

	return nil
}

// RateLimit 速率限制
func (sh *SecurityHardening) RateLimit(ip string) error {
	sh.rateLimiter.mutex.Lock()
	defer sh.rateLimiter.mutex.Unlock()

	now := time.Now()
	counter, exists := sh.rateLimiter.requests[ip]

	if !exists {
		sh.rateLimiter.requests[ip] = &RequestCounter{
			count:     1,
			lastReset: now,
		}
		return nil
	}

	// 檢查是否需要重置計數器
	if now.Sub(counter.lastReset) > sh.rateLimiter.window {
		counter.count = 1
		counter.lastReset = now
		return nil
	}

	// 檢查是否超過限制
	if counter.count >= sh.rateLimiter.maxReqs {
		return fmt.Errorf("請求頻率過高")
	}

	counter.count++
	return nil
}

// IsIPBlocked 檢查IP是否被阻擋
func (sh *SecurityHardening) IsIPBlocked(ip string) bool {
	sh.ipFilter.mutex.RLock()
	defer sh.ipFilter.mutex.RUnlock()

	// 檢查白名單
	if sh.ipFilter.whitelist[ip] {
		return false
	}

	// 檢查黑名單
	expireTime, exists := sh.ipFilter.blacklist[ip]
	if exists && time.Now().Before(expireTime) {
		return true
	}

	return false
}

// BlockIP 阻擋IP
func (sh *SecurityHardening) BlockIP(ip string, duration time.Duration) {
	sh.ipFilter.mutex.Lock()
	defer sh.ipFilter.mutex.Unlock()

	sh.ipFilter.blacklist[ip] = time.Now().Add(duration)
}

// AllowIP 允許IP
func (sh *SecurityHardening) AllowIP(ip string) {
	sh.ipFilter.mutex.Lock()
	defer sh.ipFilter.mutex.Unlock()

	sh.ipFilter.whitelist[ip] = true
	delete(sh.ipFilter.blacklist, ip)
}

// CreateSession 建立會話
func (sh *SecurityHardening) CreateSession(userID, ip, userAgent string) (string, error) {
	sessionID := sh.generateSessionID()
	now := time.Now()

	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: now,
		LastSeen:  now,
		IP:        ip,
		UserAgent: userAgent,
	}

	sh.sessionManager.mutex.Lock()
	defer sh.sessionManager.mutex.Unlock()

	sh.sessionManager.sessions[sessionID] = session
	return sessionID, nil
}

// ValidateSession 驗證會話
func (sh *SecurityHardening) ValidateSession(sessionID string) (*Session, error) {
	sh.sessionManager.mutex.RLock()
	defer sh.sessionManager.mutex.RUnlock()

	session, exists := sh.sessionManager.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("會話不存在")
	}

	if time.Since(session.LastSeen) > sh.sessionManager.timeout {
		return nil, fmt.Errorf("會話已過期")
	}

	return session, nil
}

// UpdateSession 更新會話
func (sh *SecurityHardening) UpdateSession(sessionID string) error {
	sh.sessionManager.mutex.Lock()
	defer sh.sessionManager.mutex.Unlock()

	session, exists := sh.sessionManager.sessions[sessionID]
	if !exists {
		return fmt.Errorf("會話不存在")
	}

	session.LastSeen = time.Now()
	return nil
}

// DeleteSession 刪除會話
func (sh *SecurityHardening) DeleteSession(sessionID string) {
	sh.sessionManager.mutex.Lock()
	defer sh.sessionManager.mutex.Unlock()

	delete(sh.sessionManager.sessions, sessionID)
}

// generateSessionID 生成會話ID
func (sh *SecurityHardening) generateSessionID() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// 如果產生隨機數失敗，使用時間戳作為備用方案
		return fmt.Sprintf("%x", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// EncryptData 加密資料
func (sh *SecurityHardening) EncryptData(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("資料不能為空")
	}

	// 使用簡單的XOR加密（實際應用中應使用更安全的加密方法）
	encrypted := make([]byte, len(data))
	for i, b := range data {
		encrypted[i] = b ^ sh.encryptionKey[i%len(sh.encryptionKey)]
	}

	return encrypted, nil
}

// DecryptData 解密資料
func (sh *SecurityHardening) DecryptData(encryptedData []byte) ([]byte, error) {
	if len(encryptedData) == 0 {
		return nil, fmt.Errorf("加密資料不能為空")
	}

	// 使用簡單的XOR解密
	decrypted := make([]byte, len(encryptedData))
	for i, b := range encryptedData {
		decrypted[i] = b ^ sh.encryptionKey[i%len(sh.encryptionKey)]
	}

	return decrypted, nil
}

// HashPassword 雜湊密碼
func (sh *SecurityHardening) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密碼雜湊失敗: %v", err)
	}
	return string(hash), nil
}

// VerifyPassword 驗證密碼
func (sh *SecurityHardening) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSecureToken 生成安全令牌
func (sh *SecurityHardening) GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("生成令牌失敗: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// ValidateCSRFToken 驗證CSRF令牌
func (sh *SecurityHardening) ValidateCSRFToken(token, sessionID string) bool {
	// 簡單的CSRF驗證（實際應用中應使用更安全的方法）
	expectedToken := sh.generateCSRFToken(sessionID)
	return subtle.ConstantTimeCompare([]byte(token), []byte(expectedToken)) == 1
}

// generateCSRFToken 生成CSRF令牌
func (sh *SecurityHardening) generateCSRFToken(sessionID string) string {
	data := sessionID + string(sh.encryptionKey)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// SanitizeInput 清理輸入
func (sh *SecurityHardening) SanitizeInput(input string) string {
	// 移除危險字符
	dangerousChars := []string{"<", ">", "\"", "'", "&", ";", "|", "`", "$", "(", ")"}
	sanitized := input

	for _, char := range dangerousChars {
		sanitized = strings.ReplaceAll(sanitized, char, "")
	}

	return strings.TrimSpace(sanitized)
}

// ValidateEmail 驗證電子郵件
func (sh *SecurityHardening) ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return regex.MatchString(email)
}

// ValidateURL 驗證URL
func (sh *SecurityHardening) ValidateURL(url string) bool {
	pattern := `^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return regex.MatchString(url)
}

// GetSecurityStats 取得安全統計
func (sh *SecurityHardening) GetSecurityStats() map[string]interface{} {
	sh.rateLimiter.mutex.RLock()
	sh.ipFilter.mutex.RLock()
	sh.sessionManager.mutex.RLock()
	defer sh.rateLimiter.mutex.RUnlock()
	defer sh.ipFilter.mutex.RUnlock()
	defer sh.sessionManager.mutex.RUnlock()

	return map[string]interface{}{
		"active_sessions":   len(sh.sessionManager.sessions),
		"blocked_ips":       len(sh.ipFilter.blacklist),
		"whitelisted_ips":   len(sh.ipFilter.whitelist),
		"active_requests":   len(sh.rateLimiter.requests),
		"rate_limit_window": sh.rateLimiter.window.String(),
		"session_timeout":   sh.sessionManager.timeout.String(),
		"max_requests":      sh.rateLimiter.maxReqs,
	}
}
