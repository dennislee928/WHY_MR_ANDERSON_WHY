package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"axiom-backend/internal/dto"
	apperrors "axiom-backend/internal/errors"
	"axiom-backend/internal/service"
)

// WindowsLogHandler Windows 日誌處理器
type WindowsLogHandler struct {
	windowsLogService *service.WindowsLogService
}

// NewWindowsLogHandler 創建 Windows 日誌處理器
func NewWindowsLogHandler(windowsLogService *service.WindowsLogService) *WindowsLogHandler {
	return &WindowsLogHandler{
		windowsLogService: windowsLogService,
	}
}

// BatchReceive 批量接收日誌
// @Summary 批量接收 Windows 日誌
// @Tags Windows Logs
// @Accept json
// @Produce json
// @Param request body dto.WindowsLogBatchRequest true "批量日誌請求"
// @Success 200 {object} vo.WindowsLogBatchVO
// @Router /api/v2/logs/windows/batch [post]
func (h *WindowsLogHandler) BatchReceive(c *gin.Context) {
	var req dto.WindowsLogBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.windowsLogService.BatchReceive(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// Query 查詢日誌
// @Summary 查詢 Windows 日誌
// @Tags Windows Logs
// @Produce json
// @Param agent_id query string false "Agent ID"
// @Param log_type query string false "日誌類型"
// @Param level query string false "日誌級別"
// @Param page query int false "頁碼"
// @Param page_size query int false "每頁數量"
// @Success 200 {object} vo.WindowsLogsListVO
// @Router /api/v2/logs/windows [get]
func (h *WindowsLogHandler) Query(c *gin.Context) {
	var req dto.WindowsLogQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid query parameters",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.windowsLogService.Query(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetStats 獲取統計
// @Summary 獲取 Windows 日誌統計
// @Tags Windows Logs
// @Produce json
// @Param time_range query string false "時間範圍"
// @Success 200 {object} vo.WindowsLogStatsVO
// @Router /api/v2/logs/windows/stats [get]
func (h *WindowsLogHandler) GetStats(c *gin.Context) {
	timeRange := c.DefaultQuery("time_range", "24h")

	result, err := h.windowsLogService.GetStats(c.Request.Context(), timeRange)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}


