package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	apperrors "axiom-backend/internal/errors"
	"axiom-backend/internal/service"
)

// SelfHealingHandler 自癒系統處理器
type SelfHealingHandler struct {
	selfHealingService *service.SelfHealingService
}

// NewSelfHealingHandler 創建自癒處理器
func NewSelfHealingHandler(selfHealingService *service.SelfHealingService) *SelfHealingHandler {
	return &SelfHealingHandler{
		selfHealingService: selfHealingService,
	}
}

// Remediate 執行自動修復
// @Summary 執行自動修復
// @Description 智能診斷並自動修復系統問題
// @Tags SelfHealing
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "修復請求"
// @Success 200 {object} service.HealingResult
// @Router /api/v2/combined/self-healing/remediate [post]
func (h *SelfHealingHandler) Remediate(c *gin.Context) {
	var req struct {
		IncidentType string                 `json:"incident_type" binding:"required"`
		Parameters   map[string]interface{} `json:"parameters"`
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

	result, err := h.selfHealingService.Remediate(c.Request.Context(), req.IncidentType, req.Parameters)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetSuccessRate 獲取成功率
// @Summary 獲取自癒成功率統計
// @Tags SelfHealing
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/combined/self-healing/success-rate [get]
func (h *SelfHealingHandler) GetSuccessRate(c *gin.Context) {
	result, err := h.selfHealingService.GetSuccessRate(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}


