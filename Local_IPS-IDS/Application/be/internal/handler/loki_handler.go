package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	
	apperrors "axiom-backend/internal/errors"
	"axiom-backend/internal/service"
)

// LokiHandler Loki 處理器
type LokiHandler struct {
	lokiService *service.LokiService
}

// NewLokiHandler 創建 Loki 處理器
func NewLokiHandler(lokiService *service.LokiService) *LokiHandler {
	return &LokiHandler{
		lokiService: lokiService,
	}
}

// QueryLogs 查詢日誌
// @Summary 查詢 Loki 日誌
// @Tags Loki
// @Accept json
// @Produce json
// @Param query query string true "LogQL 查詢"
// @Param limit query int false "限制數量"
// @Param start query string false "開始時間 (RFC3339)"
// @Param end query string false "結束時間 (RFC3339)"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/loki/query [get]
func (h *LokiHandler) QueryLogs(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		handleError(c, apperrors.New(
			apperrors.ErrCodeValidation,
			"query parameter is required",
			http.StatusBadRequest,
		))
		return
	}

	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}

	start := c.Query("start")
	end := c.Query("end")

	result, err := h.lokiService.QueryLogs(c.Request.Context(), query, limit, start, end)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetLabels 獲取標籤
// @Summary 獲取所有日誌標籤
// @Tags Loki
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/loki/labels [get]
func (h *LokiHandler) GetLabels(c *gin.Context) {
	labels, err := h.lokiService.GetLabels(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    gin.H{"labels": labels},
	})
}

// GetLabelValues 獲取標籤值
// @Summary 獲取指定標籤的所有值
// @Tags Loki
// @Produce json
// @Param label path string true "標籤名稱"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/loki/labels/{label}/values [get]
func (h *LokiHandler) GetLabelValues(c *gin.Context) {
	label := c.Param("label")

	values, err := h.lokiService.GetLabelValues(c.Request.Context(), label)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    gin.H{"values": values},
	})
}

// HealthCheck 健康檢查
// @Summary Loki 健康檢查
// @Tags Loki
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/loki/health [get]
func (h *LokiHandler) HealthCheck(c *gin.Context) {
	err := h.lokiService.HealthCheck(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Loki is healthy",
	})
}


