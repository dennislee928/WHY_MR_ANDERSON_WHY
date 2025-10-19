package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"axiom-backend/internal/storage"
)

// StorageHandler 儲存管理處理器
type StorageHandler struct {
	tieringPipeline *storage.TieringPipeline
}

// NewStorageHandler 創建儲存處理器
func NewStorageHandler(tieringPipeline *storage.TieringPipeline) *StorageHandler {
	return &StorageHandler{
		tieringPipeline: tieringPipeline,
	}
}

// GetTierStats 獲取各層統計
// @Summary 獲取儲存層統計
// @Tags Storage
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/storage/tiers/stats [get]
func (h *StorageHandler) GetTierStats(c *gin.Context) {
	stats, err := h.tieringPipeline.GetStats(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// TriggerTransfer 手動觸發轉移
// @Summary 手動觸發資料轉移
// @Tags Storage
// @Accept json
// @Produce json
// @Param request body map[string]string true "轉移請求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/storage/tier/transfer [post]
func (h *StorageHandler) TriggerTransfer(c *gin.Context) {
	var req struct {
		From string `json:"from" binding:"required"` // hot, warm, cold
		To   string `json:"to" binding:"required"`   // warm, cold, archive
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// TODO: 實現手動轉移邏輯
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transfer triggered",
		"from":    req.From,
		"to":      req.To,
	})
}

