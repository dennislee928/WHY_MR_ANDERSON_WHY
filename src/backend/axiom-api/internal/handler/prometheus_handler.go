package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"axiom-backend/internal/dto"
	apperrors "axiom-backend/internal/errors"
	"axiom-backend/internal/service"
)

// PrometheusHandler Prometheus 處理器
type PrometheusHandler struct {
	prometheusService *service.PrometheusService
}

// NewPrometheusHandler 創建 Prometheus 處理器
func NewPrometheusHandler(prometheusService *service.PrometheusService) *PrometheusHandler {
	return &PrometheusHandler{
		prometheusService: prometheusService,
	}
}

// Query 執行 PromQL 查詢
// @Summary 執行 PromQL 查詢
// @Tags Prometheus
// @Accept json
// @Produce json
// @Param request body dto.PrometheusQueryRequest true "查詢請求"
// @Success 200 {object} vo.PrometheusQueryVO
// @Failure 400 {object} errors.AppError
// @Router /api/v2/prometheus/query [post]
func (h *PrometheusHandler) Query(c *gin.Context) {
	var req dto.PrometheusQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.prometheusService.Query(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// QueryRange 執行範圍查詢
// @Summary 執行 Prometheus 範圍查詢
// @Tags Prometheus
// @Accept json
// @Produce json
// @Param request body dto.PrometheusQueryRangeRequest true "範圍查詢請求"
// @Success 200 {object} vo.PrometheusQueryVO
// @Failure 400 {object} errors.AppError
// @Router /api/v2/prometheus/query-range [post]
func (h *PrometheusHandler) QueryRange(c *gin.Context) {
	var req dto.PrometheusQueryRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.prometheusService.QueryRange(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetAlertRules 獲取告警規則
// @Summary 獲取所有告警規則
// @Tags Prometheus
// @Produce json
// @Success 200 {object} vo.PrometheusAlertRulesVO
// @Router /api/v2/prometheus/rules [get]
func (h *PrometheusHandler) GetAlertRules(c *gin.Context) {
	result, err := h.prometheusService.GetAlertRules(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetTargets 獲取抓取目標
// @Summary 獲取所有抓取目標
// @Tags Prometheus
// @Produce json
// @Success 200 {object} vo.PrometheusTargetsVO
// @Router /api/v2/prometheus/targets [get]
func (h *PrometheusHandler) GetTargets(c *gin.Context) {
	result, err := h.prometheusService.GetTargets(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// HealthCheck Prometheus 健康檢查
// @Summary Prometheus 健康檢查
// @Tags Prometheus
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/prometheus/health [get]
func (h *PrometheusHandler) HealthCheck(c *gin.Context) {
	err := h.prometheusService.HealthCheck(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Prometheus is healthy",
	})
}

// GetStatus 獲取 Prometheus 狀態
// @Summary 獲取 Prometheus 狀態
// @Tags Prometheus
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/prometheus/status [get]
func (h *PrometheusHandler) GetStatus(c *gin.Context) {
	status, err := h.prometheusService.GetStatus(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    status,
	})
}


