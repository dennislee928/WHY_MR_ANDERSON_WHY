package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	apperrors "axiom-backend/internal/errors"
	"axiom-backend/internal/service"
)

// CombinedHandler 組合功能處理器
type CombinedHandler struct {
	combinedService *service.CombinedService
}

// NewCombinedHandler 創建組合處理器
func NewCombinedHandler(combinedService *service.CombinedService) *CombinedHandler {
	return &CombinedHandler{
		combinedService: combinedService,
	}
}

// InvestigateIncident 一鍵事件調查
// @Summary 一鍵事件調查
// @Description 整合 Loki、Prometheus、AlertManager、Agent、AI 進行全面事件調查
// @Tags Combined
// @Accept json
// @Produce json
// @Param request body map[string]string true "調查請求"
// @Success 200 {object} service.IncidentInvestigation
// @Router /api/v2/combined/incident/investigate [post]
func (h *CombinedHandler) InvestigateIncident(c *gin.Context) {
	var req struct {
		AlertID   string `json:"alert_id" binding:"required"`
		TimeRange string `json:"time_range"`
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

	if req.TimeRange == "" {
		req.TimeRange = "1h"
	}

	result, err := h.combinedService.InvestigateIncident(c.Request.Context(), req.AlertID, req.TimeRange)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// AnalyzePerformance 全棧性能分析
// @Summary 全棧性能分析
// @Description 整合 Prometheus、Loki、Grafana、PostgreSQL、Redis 進行性能分析
// @Tags Combined
// @Produce json
// @Success 200 {object} service.PerformanceAnalysis
// @Router /api/v2/combined/performance/analyze [post]
func (h *CombinedHandler) AnalyzePerformance(c *gin.Context) {
	result, err := h.combinedService.AnalyzePerformance(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetUnifiedObservability 統一可觀測性儀表板
// @Summary 統一可觀測性儀表板
// @Description 整合指標、日誌、追蹤到統一視圖
// @Tags Combined
// @Produce json
// @Success 200 {object} service.ObservabilityDashboard
// @Router /api/v2/combined/observability/dashboard/unified [get]
func (h *CombinedHandler) GetUnifiedObservability(c *gin.Context) {
	result, err := h.combinedService.GetUnifiedObservability(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// IntelligentAlertGrouping 智能告警聚合
// @Summary 智能告警聚合和降噪
// @Description 使用 AI 分析告警模式，自動聚類和抑制
// @Tags Combined
// @Produce json
// @Success 200 {object} service.AlertNoiseReduction
// @Router /api/v2/combined/alerts/intelligent-grouping [post]
func (h *CombinedHandler) IntelligentAlertGrouping(c *gin.Context) {
	result, err := h.combinedService.IntelligentAlertGrouping(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// FullComplianceAudit 端到端合規檢查
// @Summary 端到端合規檢查
// @Description 執行完整的合規性審計並生成報告
// @Tags Combined
// @Accept json
// @Produce json
// @Param request body map[string]string true "審計請求"
// @Success 200 {object} service.ComplianceAuditResult
// @Router /api/v2/combined/compliance/full-audit [post]
func (h *CombinedHandler) FullComplianceAudit(c *gin.Context) {
	var req struct {
		Framework string `json:"framework" binding:"required"` // CIS, NIST, PCI-DSS
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

	result, err := h.combinedService.FullComplianceAudit(c.Request.Context(), req.Framework)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}


