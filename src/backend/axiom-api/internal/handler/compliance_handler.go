package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"axiom-backend/internal/compliance"
	apperrors "axiom-backend/internal/errors"
)

// ComplianceHandler 合規性處理器
type ComplianceHandler struct {
	piiDetector *compliance.PIIDetector
	anonymizer  *compliance.Anonymizer
}

// NewComplianceHandler 創建合規性處理器
func NewComplianceHandler(piiDetector *compliance.PIIDetector, anonymizer *compliance.Anonymizer) *ComplianceHandler {
	return &ComplianceHandler{
		piiDetector: piiDetector,
		anonymizer:  anonymizer,
	}
}

// DetectPII 檢測 PII
// @Summary 檢測個人資料 (PII)
// @Tags Compliance
// @Accept json
// @Produce json
// @Param request body map[string]string true "檢測請求"
// @Success 200 {object} compliance.PIIDetectionResult
// @Router /api/v2/compliance/pii/detect [post]
func (h *ComplianceHandler) DetectPII(c *gin.Context) {
	var req struct {
		Text string `json:"text" binding:"required"`
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

	result := h.piiDetector.DetectPII(c.Request.Context(), req.Text)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// AnonymizeData 匿名化資料
// @Summary 匿名化資料
// @Tags Compliance
// @Accept json
// @Produce json
// @Param request body map[string]string true "匿名化請求"
// @Success 200 {object} compliance.AnonymizationResult
// @Router /api/v2/compliance/pii/anonymize [post]
func (h *ComplianceHandler) AnonymizeData(c *gin.Context) {
	var req struct {
		Text   string                         `json:"text" binding:"required"`
		Method compliance.AnonymizationMethod `json:"method" binding:"required"`
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

	result, err := h.anonymizer.Anonymize(c.Request.Context(), req.Text, req.Method)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// Depseudonymize 反假名化
// @Summary 反假名化（還原）
// @Tags Compliance
// @Accept json
// @Produce json
// @Param request body map[string]string true "反假名化請求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/compliance/pii/depseudonymize [post]
func (h *ComplianceHandler) Depseudonymize(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
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

	originalValue, err := h.anonymizer.Depseudonymize(c.Request.Context(), req.Token)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]string{
			"original_value": originalValue,
		},
	})
}

// GetSupportedPIITypes 獲取支援的 PII 類型
// @Summary 獲取支援的 PII 類型
// @Tags Compliance
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/compliance/pii/types [get]
func (h *ComplianceHandler) GetSupportedPIITypes(c *gin.Context) {
	types := h.piiDetector.GetSupportedTypes()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"supported_types": types,
			"count":           len(types),
		},
	})
}

