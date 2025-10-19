package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"pandora_box_console_ids_ips/internal/logging"
	"pandora_box_console_ids_ips/internal/metrics"
	"pandora_box_console_ids_ips/internal/tpm"
)

// AuthHandler 認證處理器
type AuthHandler struct {
	logger        *logrus.Logger
	centralLogger *logging.CentralLogger
	metrics       *metrics.PrometheusMetrics
	challengeMgr  *tpm.ChallengeManager

	// 預共享金鑰 (生產環境中應該從安全的配置中讀取)
	preSharedKeys map[string]string
}

// PCVerifyRequest PC身份驗證請求
type PCVerifyRequest struct {
	PCIdentifier string      `json:"pc_identifier" binding:"required"`
	PreSharedKey string      `json:"pre_shared_key" binding:"required"`
	Timestamp    int64       `json:"timestamp" binding:"required"`
	EventData    interface{} `json:"event_data,omitempty"`
}

// PCVerifyResponse PC身份驗證響應
type PCVerifyResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// NewAuthHandler 創建新的認證處理器
func NewAuthHandler(logger *logrus.Logger, centralLogger *logging.CentralLogger, metrics *metrics.PrometheusMetrics) *AuthHandler {
	handler := &AuthHandler{
		logger:        logger,
		centralLogger: centralLogger,
		metrics:       metrics,
		challengeMgr:  tpm.NewChallengeManager(logger),
		preSharedKeys: make(map[string]string),
	}

	// 載入預共享金鑰 (在實際部署中，這些應該從安全配置中讀取)
	handler.loadPreSharedKeys()

	// 啟動TPM挑戰清理協程
	handler.challengeMgr.StartCleanupRoutine()

	return handler
}

// loadPreSharedKeys 載入預共享金鑰配置
func (h *AuthHandler) loadPreSharedKeys() {
	// TODO: 從安全配置檔案或環境變數載入
	// 這裡使用硬編碼的測試金鑰，生產環境中應該從安全儲存中載入

	h.preSharedKeys["test-pc-001"] = "test-key-001"
	h.preSharedKeys["test-pc-002"] = "test-key-002"
	h.preSharedKeys["default"] = "pandora-default-psk-2024"

	h.logger.Infof("載入了 %d 個預共享金鑰", len(h.preSharedKeys))
}

// VerifyPC PC身份驗證端點
func (h *AuthHandler) VerifyPC(c *gin.Context) {
	var req PCVerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("解析PC驗證請求失敗: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 驗證時間戳 (防止重放攻擊)
	now := time.Now().Unix()
	timeDiff := abs(now - req.Timestamp)
	if timeDiff > 300 { // 5分鐘容差
		h.logger.Warnf("PC驗證請求時間戳無效 - PC: %s, 時差: %d秒", req.PCIdentifier, timeDiff)
		h.centralLogger.LogSecurityEvent(req.PCIdentifier, "PC_VERIFY", "fail",
			"時間戳驗證失敗", map[string]interface{}{
				"timestamp_diff": timeDiff,
				"max_allowed":    300,
			})

		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"error":  "Invalid timestamp",
		})
		return
	}

	// 驗證預共享金鑰
	expectedKey := h.getPreSharedKey(req.PCIdentifier)
	if expectedKey == "" || expectedKey != req.PreSharedKey {
		h.logger.Warnf("PC驗證失敗 - 預共享金鑰不匹配: %s", req.PCIdentifier)
		h.centralLogger.LogSecurityEvent(req.PCIdentifier, "PC_VERIFY", "fail",
			"預共享金鑰驗證失敗", map[string]interface{}{
				"reason": "key_mismatch",
			})

		// 更新指標
		if h.metrics != nil {
			h.metrics.AuthAttempts.WithLabelValues("fail").Inc()
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"error":  "Authentication failed",
		})
		return
	}

	// 驗證成功
	now = time.Now().Unix()
	h.logger.Infof("PC身份驗證成功 - PC: %s", req.PCIdentifier)

	// 記錄成功的認證事件
	h.centralLogger.LogAuthEvent(req.PCIdentifier, "PC_VERIFY", "success",
		"PC身份驗證成功", map[string]interface{}{
			"timestamp":  req.Timestamp,
			"event_data": req.EventData,
		})

	// 更新指標
	if h.metrics != nil {
		h.metrics.AuthAttempts.WithLabelValues("success").Inc()
	}

	c.JSON(http.StatusOK, PCVerifyResponse{
		Status:    "success",
		Message:   "PC verified successfully",
		Timestamp: now,
	})
}

// GetChallenge 獲取TPM挑戰
func (h *AuthHandler) GetChallenge(c *gin.Context) {
	var req struct {
		PCIdentifier string `json:"pc_identifier" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("解析TPM挑戰請求失敗: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 生成挑戰
	challenge, err := h.challengeMgr.GenerateChallenge(req.PCIdentifier)
	if err != nil {
		h.logger.Errorf("生成TPM挑戰失敗: %v", err)
		h.centralLogger.LogSecurityEvent(req.PCIdentifier, "CHALLENGE_GENERATION", "fail",
			"TPM挑戰生成失敗", map[string]interface{}{
				"error": err.Error(),
			})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate challenge"})
		return
	}

	// 記錄日誌
	h.centralLogger.LogAuthEvent(req.PCIdentifier, "TPM_CHALLENGE_REQUEST", "success",
		"TPM挑戰生成成功", map[string]interface{}{
			"challenge_id": challenge.ID,
			"expires_at":   challenge.ExpiresAt,
		})

	// 更新指標
	if h.metrics != nil {
		h.metrics.RecordSecurityEvent("tpm_challenge", "generated")
	}

	c.JSON(http.StatusOK, gin.H{
		"challenge_id":    challenge.ID,
		"nonce":           challenge.Nonce,
		"expires_at":      challenge.ExpiresAt.Unix(),
		"timeout_seconds": 300, // 5分鐘
	})
}

// VerifyTPMSignature 驗證TPM簽章
func (h *AuthHandler) VerifyTPMSignature(c *gin.Context) {
	var req tpm.ChallengeResponse

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("解析TPM簽章驗證請求失敗: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 驗證挑戰回應
	valid, err := h.challengeMgr.VerifyChallenge(&req)
	if err != nil {
		h.logger.Errorf("TPM挑戰驗證失敗: %v", err)
		h.centralLogger.LogSecurityEvent(req.PCIdentifier, "TPM_VERIFICATION", "fail",
			"TPM挑戰驗證失敗", map[string]interface{}{
				"challenge_id": req.ChallengeID,
				"error":        err.Error(),
			})

		// 更新指標
		if h.metrics != nil {
			h.metrics.AuthAttempts.WithLabelValues("fail").Inc()
			h.metrics.RecordSecurityEvent("tpm_verification", "failed")
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"error":  "TPM verification failed",
		})
		return
	}

	if !valid {
		h.logger.Warnf("TPM簽章驗證失敗 - PC: %s", req.PCIdentifier)
		h.centralLogger.LogSecurityEvent(req.PCIdentifier, "TPM_VERIFICATION", "fail",
			"TPM簽章無效", map[string]interface{}{
				"challenge_id": req.ChallengeID,
			})

		// 更新指標
		if h.metrics != nil {
			h.metrics.AuthAttempts.WithLabelValues("fail").Inc()
			h.metrics.RecordSecurityEvent("tpm_verification", "invalid")
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "fail",
			"error":  "Invalid TPM signature",
		})
		return
	}

	// 驗證成功
	h.logger.Infof("TPM認證成功 - PC: %s, 挑戰: %s", req.PCIdentifier, req.ChallengeID)
	h.centralLogger.LogAuthEvent(req.PCIdentifier, "TPM_VERIFICATION", "success",
		"TPM認證成功", map[string]interface{}{
			"challenge_id": req.ChallengeID,
			"timestamp":    req.Timestamp,
		})

	// 更新指標
	if h.metrics != nil {
		h.metrics.AuthAttempts.WithLabelValues("success").Inc()
		h.metrics.RecordSecurityEvent("tpm_verification", "success")
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"pc_identifier": req.PCIdentifier,
		"verified_at":   time.Now().Unix(),
		"message":       "TPM authentication successful",
	})
}

// HealthCheck 健康檢查端點
func (h *AuthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "pandora-console-auth",
		"version":   "1.0.0",
	})
}

// getPreSharedKey 獲取PC的預共享金鑰
func (h *AuthHandler) getPreSharedKey(pcIdentifier string) string {
	if key, exists := h.preSharedKeys[pcIdentifier]; exists {
		return key
	}

	// 如果找不到特定的金鑰，使用預設金鑰
	if defaultKey, exists := h.preSharedKeys["default"]; exists {
		return defaultKey
	}

	return ""
}

// abs 計算絕對值
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
