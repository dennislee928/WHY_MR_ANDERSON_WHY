package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	apperrors "axiom-backend/internal/errors"
	"axiom-backend/internal/service"
)

// AdaptiveSecurityHandler 自適應安全處理器
type AdaptiveSecurityHandler struct {
	adaptiveSecurityService *service.AdaptiveSecurityService
}

// NewAdaptiveSecurityHandler 創建自適應安全處理器
func NewAdaptiveSecurityHandler(adaptiveSecurityService *service.AdaptiveSecurityService) *AdaptiveSecurityHandler {
	return &AdaptiveSecurityHandler{
		adaptiveSecurityService: adaptiveSecurityService,
	}
}

// CalculateRisk 計算風險評分
// @Summary 計算動態風險評分
// @Tags AdaptiveSecurity
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "風險評估請求"
// @Success 200 {object} service.RiskScore
// @Router /api/v2/adaptive-security/risk/calculate [post]
func (h *AdaptiveSecurityHandler) CalculateRisk(c *gin.Context) {
	var req struct {
		UserID    string                 `json:"user_id" binding:"required"`
		IPAddress string                 `json:"ip_address" binding:"required"`
		Context   map[string]interface{} `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.adaptiveSecurityService.CalculateRisk(c.Request.Context(), req.UserID, req.IPAddress, req.Context)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// EvaluateAccess 評估訪問
// @Summary 評估訪問請求
// @Tags AdaptiveSecurity
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "訪問評估請求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/adaptive-security/access/evaluate [post]
func (h *AdaptiveSecurityHandler) EvaluateAccess(c *gin.Context) {
	var req map[string]interface{}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.adaptiveSecurityService.EvaluateAccess(c.Request.Context(), req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetTrustScore 獲取信任分數
// @Summary 獲取實體信任分數
// @Tags AdaptiveSecurity
// @Produce json
// @Param entityId path string true "實體 ID"
// @Success 200 {object} service.TrustScore
// @Router /api/v2/adaptive-security/access/trust-score/{entityId} [get]
func (h *AdaptiveSecurityHandler) GetTrustScore(c *gin.Context) {
	entityID := c.Param("entityId")

	result, err := h.adaptiveSecurityService.GetTrustScore(c.Request.Context(), entityID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// DeployHoneypot 部署蜜罐
// @Summary 自動部署蜜罐
// @Tags AdaptiveSecurity
// @Accept json
// @Produce json
// @Param request body map[string]string true "蜜罐部署請求"
// @Success 200 {object} service.HoneypotDeployment
// @Router /api/v2/adaptive-security/honeypot/deploy [post]
func (h *AdaptiveSecurityHandler) DeployHoneypot(c *gin.Context) {
	var req struct {
		Type string `json:"type" binding:"required"` // ssh, http, database
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.adaptiveSecurityService.DeployHoneypot(c.Request.Context(), req.Type)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetHoneypotInteractions 獲取蜜罐互動
// @Summary 獲取蜜罐互動記錄
// @Tags AdaptiveSecurity
// @Produce json
// @Param honeypotId path string true "蜜罐 ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/adaptive-security/honeypot/{honeypotId}/interactions [get]
func (h *AdaptiveSecurityHandler) GetHoneypotInteractions(c *gin.Context) {
	honeypotID := c.Param("honeypotId")

	result, err := h.adaptiveSecurityService.GetHoneypotInteractions(c.Request.Context(), honeypotID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// AnalyzeAttacker 分析攻擊者
// @Summary 分析攻擊者行為
// @Tags AdaptiveSecurity
// @Accept json
// @Produce json
// @Param request body map[string]string true "攻擊者分析請求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/adaptive-security/honeypot/analyze-attacker [post]
func (h *AdaptiveSecurityHandler) AnalyzeAttacker(c *gin.Context) {
	var req struct {
		AttackerIP string `json:"attacker_ip" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.adaptiveSecurityService.AnalyzeAttacker(c.Request.Context(), req.AttackerIP)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}


