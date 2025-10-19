package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"axiom-backend/internal/service"
)

// APIGovernanceHandler API 治理處理器
type APIGovernanceHandler struct {
	apiGovernanceService *service.APIGovernanceService
}

// NewAPIGovernanceHandler 創建 API 治理處理器
func NewAPIGovernanceHandler(apiGovernanceService *service.APIGovernanceService) *APIGovernanceHandler {
	return &APIGovernanceHandler{
		apiGovernanceService: apiGovernanceService,
	}
}

// GetAPIHealth 獲取 API 健康評分
// @Summary 獲取 API 健康評分
// @Tags API Governance
// @Produce json
// @Param apiPath path string true "API 路徑"
// @Success 200 {object} service.APIHealthScore
// @Router /api/v2/governance/api-health/{apiPath} [get]
func (h *APIGovernanceHandler) GetAPIHealth(c *gin.Context) {
	apiPath := c.Param("apiPath")

	result, err := h.apiGovernanceService.GetAPIHealth(c.Request.Context(), apiPath)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetUsageAnalytics 獲取使用分析
// @Summary 獲取 API 使用分析
// @Tags API Governance
// @Produce json
// @Param time_range query string false "時間範圍"
// @Success 200 {object} service.APIUsageAnalytics
// @Router /api/v2/governance/api-usage-analytics [get]
func (h *APIGovernanceHandler) GetUsageAnalytics(c *gin.Context) {
	timeRange := c.DefaultQuery("time_range", "24h")

	result, err := h.apiGovernanceService.GetUsageAnalytics(c.Request.Context(), timeRange)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}


