package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"axiom-backend/internal/compliance"
	apperrors "axiom-backend/internal/errors"
)

// GDPRHandler GDPR 合規處理器
type GDPRHandler struct {
	gdprService *compliance.GDPRService
}

// NewGDPRHandler 創建 GDPR 處理器
func NewGDPRHandler(gdprService *compliance.GDPRService) *GDPRHandler {
	return &GDPRHandler{
		gdprService: gdprService,
	}
}

// CreateDeletionRequest 創建刪除請求
// @Summary 創建 GDPR 刪除請求
// @Tags GDPR
// @Accept json
// @Produce json
// @Param request body map[string]string true "刪除請求"
// @Success 200 {object} model.GDPRDeletionRequest
// @Router /api/v2/compliance/gdpr/deletion-request [post]
func (h *GDPRHandler) CreateDeletionRequest(c *gin.Context) {
	var req struct {
		SubjectID   string `json:"subject_id" binding:"required"`
		RequestedBy string `json:"requested_by" binding:"required"`
		Notes       string `json:"notes"`
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

	result, err := h.gdprService.CreateDeletionRequest(c.Request.Context(), req.SubjectID, req.RequestedBy, req.Notes)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ListDeletionRequests 列出刪除請求
// @Summary 列出 GDPR 刪除請求
// @Tags GDPR
// @Produce json
// @Param status query string false "狀態過濾"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/compliance/gdpr/deletion-request/list [get]
func (h *GDPRHandler) ListDeletionRequests(c *gin.Context) {
	status := c.Query("status")

	result, err := h.gdprService.ListDeletionRequests(c.Request.Context(), status)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"count":   len(result),
	})
}

// ApproveDeletionRequest 審批刪除請求
// @Summary 審批 GDPR 刪除請求
// @Tags GDPR
// @Accept json
// @Produce json
// @Param requestId path string true "請求 ID"
// @Param request body map[string]string true "審批請求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/compliance/gdpr/deletion-request/{requestId}/approve [post]
func (h *GDPRHandler) ApproveDeletionRequest(c *gin.Context) {
	requestID := c.Param("requestId")
	
	var req struct {
		ApprovedBy string `json:"approved_by" binding:"required"`
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

	err := h.gdprService.ApproveDeletionRequest(c.Request.Context(), requestID, req.ApprovedBy)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Deletion request approved",
	})
}

// ExecuteDeletion 執行刪除
// @Summary 執行 GDPR 刪除
// @Tags GDPR
// @Produce json
// @Param requestId path string true "請求 ID"
// @Success 200 {object} compliance.DeletionResult
// @Router /api/v2/compliance/gdpr/deletion-request/{requestId}/execute [post]
func (h *GDPRHandler) ExecuteDeletion(c *gin.Context) {
	requestID := c.Param("requestId")

	result, err := h.gdprService.ExecuteDeletion(c.Request.Context(), requestID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// VerifyDeletion 驗證刪除
// @Summary 驗證 GDPR 刪除
// @Tags GDPR
// @Produce json
// @Param requestId path string true "請求 ID"
// @Success 200 {object} compliance.DeletionVerification
// @Router /api/v2/compliance/gdpr/deletion-request/{requestId}/verify [get]
func (h *GDPRHandler) VerifyDeletion(c *gin.Context) {
	requestID := c.Param("requestId")

	result, err := h.gdprService.VerifyDeletion(c.Request.Context(), requestID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ExportData 匯出個人數據
// @Summary 匯出個人數據 (GDPR 資料可攜性)
// @Tags GDPR
// @Accept json
// @Produce json
// @Param request body map[string]string true "匯出請求"
// @Success 200 {object} compliance.DataExport
// @Router /api/v2/compliance/gdpr/data-export [post]
func (h *GDPRHandler) ExportData(c *gin.Context) {
	var req struct {
		SubjectID string `json:"subject_id" binding:"required"`
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

	result, err := h.gdprService.ExportData(c.Request.Context(), req.SubjectID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

