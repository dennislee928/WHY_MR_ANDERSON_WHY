package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"axiom-backend/internal/dto"
	apperrors "axiom-backend/internal/errors"
	"axiom-backend/internal/service"
)

// QuantumHandler 量子功能處理器
type QuantumHandler struct {
	quantumService *service.QuantumService
}

// NewQuantumHandler 創建量子處理器
func NewQuantumHandler(quantumService *service.QuantumService) *QuantumHandler {
	return &QuantumHandler{
		quantumService: quantumService,
	}
}

// GenerateQKD 生成量子密鑰
// @Summary 生成量子密鑰分發
// @Tags Quantum
// @Accept json
// @Produce json
// @Param request body dto.QuantumQKDRequest true "QKD 請求"
// @Success 200 {object} vo.QuantumQKDVO
// @Router /api/v2/quantum/qkd/generate [post]
func (h *QuantumHandler) GenerateQKD(c *gin.Context) {
	var req dto.QuantumQKDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.quantumService.GenerateQKD(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ClassifyQSVM 執行 QSVM 分類
// @Summary 執行量子支持向量機分類
// @Tags Quantum
// @Accept json
// @Produce json
// @Param request body dto.QuantumQSVMRequest true "QSVM 請求"
// @Success 200 {object} vo.QuantumClassifyVO
// @Router /api/v2/quantum/qsvm/classify [post]
func (h *QuantumHandler) ClassifyQSVM(c *gin.Context) {
	var req dto.QuantumQSVMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.quantumService.ClassifyQSVM(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// PredictZeroTrust 執行 Zero Trust 預測
// @Summary 執行 Zero Trust 安全預測
// @Tags Quantum
// @Accept json
// @Produce json
// @Param request body dto.ZeroTrustPredictRequest true "Zero Trust 請求"
// @Success 200 {object} vo.ZeroTrustPredictVO
// @Router /api/v2/quantum/zerotrust/predict [post]
func (h *QuantumHandler) PredictZeroTrust(c *gin.Context) {
	var req dto.ZeroTrustPredictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.quantumService.PredictZeroTrust(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetJob 獲取量子作業
// @Summary 獲取量子作業詳情
// @Tags Quantum
// @Produce json
// @Param jobId path string true "作業 ID"
// @Success 200 {object} vo.QuantumJobVO
// @Router /api/v2/quantum/jobs/{jobId} [get]
func (h *QuantumHandler) GetJob(c *gin.Context) {
	jobID := c.Param("jobId")

	result, err := h.quantumService.GetQuantumJob(c.Request.Context(), jobID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ListJobs 列出量子作業
// @Summary 列出所有量子作業
// @Tags Quantum
// @Produce json
// @Param type query string false "作業類型"
// @Param status query string false "作業狀態"
// @Param page query int false "頁碼"
// @Param page_size query int false "每頁數量"
// @Success 200 {object} vo.QuantumJobsListVO
// @Router /api/v2/quantum/jobs [get]
func (h *QuantumHandler) ListJobs(c *gin.Context) {
	var req dto.QuantumJobQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid query parameters",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.quantumService.ListQuantumJobs(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetStats 獲取量子統計
// @Summary 獲取量子作業統計
// @Tags Quantum
// @Produce json
// @Success 200 {object} vo.QuantumStatsVO
// @Router /api/v2/quantum/stats [get]
func (h *QuantumHandler) GetStats(c *gin.Context) {
	result, err := h.quantumService.GetQuantumStats(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// HealthCheck 健康檢查
// @Summary 量子服務健康檢查
// @Tags Quantum
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/quantum/health [get]
func (h *QuantumHandler) HealthCheck(c *gin.Context) {
	err := h.quantumService.HealthCheck(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Quantum service is healthy",
	})
}


