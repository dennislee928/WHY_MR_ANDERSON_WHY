package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Auth USB Token認證系統
type Auth struct {
	logger      *logrus.Logger
	isEnabled   bool
	tokenPath   string
	validTokens map[string]time.Time
	mu          sync.RWMutex
}

// NewAuth 建立新的USB Token認證系統
func NewAuth(logger *logrus.Logger) *Auth {
	return &Auth{
		logger:      logger,
		validTokens: make(map[string]time.Time),
	}
}

// Initialize 初始化USB Token認證系統
func (a *Auth) Initialize() error {
	// 偵測USB Token裝置
	if err := a.detectUSBToken(); err != nil {
		a.logger.Warnf("USB Token偵測失敗: %v", err)
		a.isEnabled = false
		return nil // 不強制要求USB Token
	}

	a.isEnabled = true
	a.logger.Info("USB Token認證系統初始化完成")
	return nil
}

// detectUSBToken 偵測USB Token裝置
func (a *Auth) detectUSBToken() error {
	// 在Linux系統中偵測USB裝置
	if err := a.detectLinuxUSB(); err != nil {
		return err
	}

	// 在Windows系統中偵測USB裝置
	if err := a.detectWindowsUSB(); err != nil {
		return err
	}

	return fmt.Errorf("未找到USB Token裝置")
}

// detectLinuxUSB 偵測Linux USB裝置
func (a *Auth) detectLinuxUSB() error {
	// 檢查常見的USB Token掛載點
	possiblePaths := []string{
		"/media/usb",
		"/mnt/usb",
		"/dev/sd*",
		"/dev/usb*",
	}

	for _, path := range possiblePaths {
		if matches, err := filepath.Glob(path); err == nil {
			for _, match := range matches {
				if a.isTokenDevice(match) {
					a.tokenPath = match
					a.logger.Infof("找到USB Token裝置: %s", match)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("未找到Linux USB Token裝置")
}

// detectWindowsUSB 偵測Windows USB裝置
func (a *Auth) detectWindowsUSB() error {
	// 檢查Windows USB裝置
	drives := []string{"D:", "E:", "F:", "G:", "H:"}

	for _, drive := range drives {
		if a.isTokenDevice(drive) {
			a.tokenPath = drive
			a.logger.Infof("找到USB Token裝置: %s", drive)
			return nil
		}
	}

	return fmt.Errorf("未找到Windows USB Token裝置")
}

// isTokenDevice 檢查是否為Token裝置
func (a *Auth) isTokenDevice(path string) bool {
	// 檢查是否存在token.key檔案
	tokenFile := filepath.Join(path, "token.key")
	if _, err := os.Stat(tokenFile); err == nil {
		return true
	}

	// 檢查是否存在.auth目錄
	authDir := filepath.Join(path, ".auth")
	if stat, err := os.Stat(authDir); err == nil && stat.IsDir() {
		return true
	}

	return false
}

// IsEnabled 檢查USB Token認證是否啟用
func (a *Auth) IsEnabled() bool {
	return a.isEnabled
}

// ValidateToken 驗證USB Token
func (a *Auth) ValidateToken() bool {
	if !a.isEnabled {
		return true // 如果未啟用，直接通過
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	// 檢查Token檔案
	tokenData, err := a.readTokenData()
	if err != nil {
		a.logger.Errorf("讀取Token資料失敗: %v", err)
		return false
	}

	// 驗證Token
	if !a.verifyTokenData(tokenData) {
		a.logger.Warn("Token驗證失敗")
		return false
	}

	// 檢查Token是否在有效清單中
	tokenHash := a.hashToken(tokenData)
	if expiry, exists := a.validTokens[tokenHash]; exists {
		if time.Now().After(expiry) {
			a.logger.Warn("Token已過期")
			delete(a.validTokens, tokenHash)
			return false
		}
	} else {
		// 新Token，需要註冊
		if !a.registerNewToken(tokenHash) {
			a.logger.Warn("新Token註冊失敗")
			return false
		}
	}

	a.logger.Info("USB Token驗證成功")
	return true
}

// readTokenData 讀取Token資料
func (a *Auth) readTokenData() ([]byte, error) {
	tokenFile := filepath.Join(a.tokenPath, "token.key")

	file, err := os.Open(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("無法開啟Token檔案: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("讀取Token檔案失敗: %v", err)
	}

	return data, nil
}

// verifyTokenData 驗證Token資料
func (a *Auth) verifyTokenData(data []byte) bool {
	// 檢查Token格式
	if len(data) < 32 {
		return false
	}

	// 檢查Token簽章（這裡使用簡單的檢查）
	// 實際應用中應該使用更複雜的加密驗證
	return len(data) >= 32 && len(data) <= 1024
}

// hashToken 對Token進行雜湊
func (a *Auth) hashToken(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// registerNewToken 註冊新Token
func (a *Auth) registerNewToken(tokenHash string) bool {
	// 檢查是否為授權的Token
	if !a.isAuthorizedToken(tokenHash) {
		return false
	}

	// 設定Token過期時間（24小時）
	a.validTokens[tokenHash] = time.Now().Add(24 * time.Hour)
	a.logger.Infof("新Token已註冊: %s", tokenHash[:8]+"...")
	return true
}

// isAuthorizedToken 檢查是否為授權的Token
func (a *Auth) isAuthorizedToken(tokenHash string) bool {
	// 這裡應該檢查預先註冊的Token清單
	// 暫時使用簡單的檢查
	authorizedTokens := []string{
		"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6",
		"b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a1",
		// 可以從設定檔或資料庫載入更多授權Token
	}

	for _, authorized := range authorizedTokens {
		if tokenHash == authorized {
			return true
		}
	}

	return false
}

// GenerateToken 生成新的Token（用於管理員）
func (a *Auth) GenerateToken() (string, error) {
	// 生成32位元組的隨機Token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("生成Token失敗: %v", err)
	}

	tokenStr := hex.EncodeToString(tokenBytes)
	a.logger.Infof("生成新Token: %s", tokenStr)
	return tokenStr, nil
}

// WriteTokenToUSB 將Token寫入USB裝置
func (a *Auth) WriteTokenToUSB(token string) error {
	if !a.isEnabled {
		return fmt.Errorf("USB Token系統未啟用")
	}

	tokenFile := filepath.Join(a.tokenPath, "token.key")

	file, err := os.Create(tokenFile)
	if err != nil {
		return fmt.Errorf("建立Token檔案失敗: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(token)
	if err != nil {
		return fmt.Errorf("寫入Token失敗: %v", err)
	}

	a.logger.Infof("Token已寫入USB裝置: %s", tokenFile)
	return nil
}

// RevokeToken 撤銷Token
func (a *Auth) RevokeToken(tokenHash string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	delete(a.validTokens, tokenHash)
	a.logger.Infof("Token已撤銷: %s", tokenHash[:8]+"...")
}

// GetValidTokens 取得有效Token清單
func (a *Auth) GetValidTokens() map[string]time.Time {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// 回傳副本
	tokens := make(map[string]time.Time)
	for k, v := range a.validTokens {
		tokens[k] = v
	}

	return tokens
}

// CleanExpiredTokens 清理過期的Token
func (a *Auth) CleanExpiredTokens() {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	for tokenHash, expiry := range a.validTokens {
		if now.After(expiry) {
			delete(a.validTokens, tokenHash)
			a.logger.Debugf("已清理過期Token: %s", tokenHash[:8]+"...")
		}
	}
}
