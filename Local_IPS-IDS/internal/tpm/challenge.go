package tpm

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ChallengeManager TPM挑戰管理器
type ChallengeManager struct {
	logger     *logrus.Logger
	challenges map[string]*Challenge
	mutex      sync.RWMutex

	// 配置參數
	challengeTimeout time.Duration
	maxChallenges    int
}

// Challenge TPM挑戰結構
type Challenge struct {
	ID           string    `json:"id"`
	Nonce        string    `json:"nonce"`
	PCIdentifier string    `json:"pc_identifier"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Used         bool      `json:"used"`

	// TPM相關
	ExpectedPublicKey *rsa.PublicKey `json:"-"`
}

// ChallengeResponse TPM挑戰回應結構
type ChallengeResponse struct {
	ChallengeID  string `json:"challenge_id" binding:"required"`
	PCIdentifier string `json:"pc_identifier" binding:"required"`
	Signature    string `json:"signature" binding:"required"`
	PublicKey    string `json:"public_key" binding:"required"`
	Timestamp    int64  `json:"timestamp" binding:"required"`
}

// NewChallengeManager 創建新的挑戰管理器
func NewChallengeManager(logger *logrus.Logger) *ChallengeManager {
	return &ChallengeManager{
		logger:           logger,
		challenges:       make(map[string]*Challenge),
		challengeTimeout: 5 * time.Minute, // 5分鐘過期
		maxChallenges:    1000,            // 最大挑戰數量
	}
}

// GenerateChallenge 生成新的TPM挑戰
func (cm *ChallengeManager) GenerateChallenge(pcIdentifier string) (*Challenge, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 檢查挑戰數量限制
	if len(cm.challenges) >= cm.maxChallenges {
		// 清理過期挑戰
		cm.cleanupExpiredChallenges()

		if len(cm.challenges) >= cm.maxChallenges {
			return nil, fmt.Errorf("挑戰數量已達上限")
		}
	}

	// 生成唯一ID和隨機nonce
	challengeID := cm.generateChallengeID()
	nonce, err := cm.generateSecureNonce(32) // 32字節隨機數
	if err != nil {
		return nil, fmt.Errorf("生成nonce失敗: %v", err)
	}

	now := time.Now()
	challenge := &Challenge{
		ID:           challengeID,
		Nonce:        nonce,
		PCIdentifier: pcIdentifier,
		CreatedAt:    now,
		ExpiresAt:    now.Add(cm.challengeTimeout),
		Used:         false,
	}

	cm.challenges[challengeID] = challenge

	cm.logger.Infof("為PC %s 生成TPM挑戰: %s", pcIdentifier, challengeID)
	return challenge, nil
}

// VerifyChallenge 驗證TPM挑戰回應
func (cm *ChallengeManager) VerifyChallenge(response *ChallengeResponse) (bool, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 查找挑戰
	challenge, exists := cm.challenges[response.ChallengeID]
	if !exists {
		return false, fmt.Errorf("挑戰不存在: %s", response.ChallengeID)
	}

	// 檢查是否過期
	if time.Now().After(challenge.ExpiresAt) {
		delete(cm.challenges, response.ChallengeID)
		return false, fmt.Errorf("挑戰已過期")
	}

	// 檢查是否已使用
	if challenge.Used {
		return false, fmt.Errorf("挑戰已被使用")
	}

	// 檢查PC標識符
	if challenge.PCIdentifier != response.PCIdentifier {
		return false, fmt.Errorf("PC標識符不匹配")
	}

	// 解析公鑰
	publicKey, err := cm.parsePublicKey(response.PublicKey)
	if err != nil {
		return false, fmt.Errorf("解析公鑰失敗: %v", err)
	}

	// 驗證簽章
	valid, err := cm.verifySignature(challenge.Nonce, response.Signature, publicKey)
	if err != nil {
		return false, fmt.Errorf("驗證簽章失敗: %v", err)
	}

	if !valid {
		cm.logger.Warnf("TPM簽章驗證失敗 - PC: %s, 挑戰: %s",
			response.PCIdentifier, response.ChallengeID)
		return false, fmt.Errorf("簽章驗證失敗")
	}

	// 標記挑戰為已使用
	challenge.Used = true
	challenge.ExpectedPublicKey = publicKey

	cm.logger.Infof("TPM挑戰驗證成功 - PC: %s, 挑戰: %s",
		response.PCIdentifier, response.ChallengeID)

	return true, nil
}

// GetChallenge 獲取挑戰信息
func (cm *ChallengeManager) GetChallenge(challengeID string) (*Challenge, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	challenge, exists := cm.challenges[challengeID]
	if !exists {
		return nil, fmt.Errorf("挑戰不存在")
	}

	return challenge, nil
}

// generateChallengeID 生成挑戰ID
func (cm *ChallengeManager) generateChallengeID() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	return fmt.Sprintf("tpm_%d_%x", timestamp, randomBytes)
}

// generateSecureNonce 生成安全的隨機nonce
func (cm *ChallengeManager) generateSecureNonce(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// parsePublicKey 解析公鑰 (簡化版本)
func (cm *ChallengeManager) parsePublicKey(publicKeyStr string) (*rsa.PublicKey, error) {
	// TODO: 實作完整的公鑰解析邏輯
	// 這裡需要根據實際的TPM公鑰格式進行解析
	// 可能需要處理PEM、DER或其他格式

	cm.logger.Debug("解析TPM公鑰 (簡化版本)")

	// 暫時返回nil，表示需要完整實作
	return nil, fmt.Errorf("公鑰解析功能需要完整實作")
}

// verifySignature 驗證簽章
func (cm *ChallengeManager) verifySignature(nonce, signature string, publicKey *rsa.PublicKey) (bool, error) {
	// TODO: 實作完整的RSA簽章驗證
	// 這裡需要：
	// 1. 將nonce進行SHA256哈希
	// 2. 解碼簽章
	// 3. 使用RSA公鑰驗證簽章

	cm.logger.Debug("驗證TPM簽章 (簡化版本)")

	if publicKey == nil {
		return false, fmt.Errorf("公鑰為空")
	}

	// 計算nonce的哈希
	nonceBytes, err := hex.DecodeString(nonce)
	if err != nil {
		return false, fmt.Errorf("解碼nonce失敗: %v", err)
	}

	hash := sha256.Sum256(nonceBytes)

	// 解碼簽章
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("解碼簽章失敗: %v", err)
	}

	// TODO: 使用rsa.VerifyPKCS1v15進行實際驗證
	// err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], sigBytes)
	// return err == nil, err

	// 暫時返回基本驗證（僅用於開發測試）
	// 使用變數避免編譯器警告
	_ = hash
	_ = sigBytes

	cm.logger.Warn("使用簡化版本的TPM簽章驗證")
	return len(signature) > 0 && len(nonce) > 0, nil
}

// cleanupExpiredChallenges 清理過期挑戰
func (cm *ChallengeManager) cleanupExpiredChallenges() {
	now := time.Now()
	toDelete := []string{}

	for id, challenge := range cm.challenges {
		if now.After(challenge.ExpiresAt) {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range toDelete {
		delete(cm.challenges, id)
		cm.logger.Debugf("清理過期挑戰: %s", id)
	}

	if len(toDelete) > 0 {
		cm.logger.Infof("清理了 %d 個過期挑戰", len(toDelete))
	}
}

// GetStats 獲取挑戰管理器統計資訊
func (cm *ChallengeManager) GetStats() map[string]interface{} {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	active := 0
	expired := 0
	used := 0
	now := time.Now()

	for _, challenge := range cm.challenges {
		if now.After(challenge.ExpiresAt) {
			expired++
		} else if challenge.Used {
			used++
		} else {
			active++
		}
	}

	return map[string]interface{}{
		"total_challenges": len(cm.challenges),
		"active":           active,
		"expired":          expired,
		"used":             used,
		"timeout_minutes":  cm.challengeTimeout.Minutes(),
	}
}

// StartCleanupRoutine 啟動清理協程
func (cm *ChallengeManager) StartCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute) // 每分鐘清理一次
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cm.mutex.Lock()
				cm.cleanupExpiredChallenges()
				cm.mutex.Unlock()
			}
		}
	}()

	cm.logger.Info("TPM挑戰清理協程已啟動")
}
