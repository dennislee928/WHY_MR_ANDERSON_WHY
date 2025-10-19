package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"axiom-backend/internal/service"
)

// TechDebtHandler 技術債務處理器
type TechDebtHandler struct {
	techDebtService *service.TechDebtService
}

// NewTechDebtHandler 創建技術債務處理器
func NewTechDebtHandler(techDebtService *service.TechDebtService) *TechDebtHandler {
	return &TechDebtHandler{
		techDebtService: techDebtService,
	}
}

// ScanTechDebt 掃描技術債務
// @Summary 掃描技術債務
// @Tags Tech Debt
// @Produce json
// @Success 200 {object} service.TechDebtScan
// @Router /api/v2/tech-debt/scan [post]
func (h *TechDebtHandler) ScanTechDebt(c *gin.Context) {
	result, err := h.techDebtService.ScanTechDebt(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GenerateRoadmap 生成修復路線圖
// @Summary 生成技術債務修復路線圖
// @Tags Tech Debt
// @Produce json
// @Success 200 {object} service.RemediationRoadmap
// @Router /api/v2/tech-debt/remediation-roadmap [post]
func (h *TechDebtHandler) GenerateRoadmap(c *gin.Context) {
	result, err := h.techDebtService.GenerateRoadmap(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}


