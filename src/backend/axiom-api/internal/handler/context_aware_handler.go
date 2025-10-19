package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	apperrors "axiom-backend/internal/errors"
	"axiom-backend/internal/service"
)

// ContextAwareHandler 情境感知處理器
type ContextAwareHandler struct {
	contextAwareService *service.ContextAwareService
}

// NewContextAwareHandler 創建情境感知處理器
func NewContextAwareHandler(contextAwareService *service.ContextAwareService) *ContextAwareHandler {
	return &ContextAwareHandler{
		contextAwareService: contextAwareService,
	}
}

// RouteAlert 智能告警路由
// @Summary 智能告警路由
// @Tags Context Aware
// @Accept json
// @Produce json
// @Param request body map[string]string true "路由請求"
// @Success 200 {object} service.AlertRouting
// @Router /api/v2/context-aware/alert-routing [post]
func (h *ContextAwareHandler) RouteAlert(c *gin.Context) {
	var req struct {
		AlertID   string `json:"alert_id" binding:"required"`
		AlertName string `json:"alert_name" binding:"required"`
		Severity  string `json:"severity" binding:"required"`
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

	result, err := h.contextAwareService.RouteAlert(c.Request.Context(), req.AlertID, req.AlertName, req.Severity)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}


