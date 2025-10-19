package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	apperrors "axiom-backend/internal/errors"
	"axiom-backend/internal/service"
)

// TimeTravelHandler 時間旅行處理器
type TimeTravelHandler struct {
	timeTravelService *service.TimeTravelService
}

// NewTimeTravelHandler 創建時間旅行處理器
func NewTimeTravelHandler(timeTravelService *service.TimeTravelService) *TimeTravelHandler {
	return &TimeTravelHandler{
		timeTravelService: timeTravelService,
	}
}

// CreateSnapshot 創建系統快照
// @Summary 創建系統狀態快照
// @Tags TimeTravel
// @Accept json
// @Produce json
// @Param request body map[string]string true "快照請求"
// @Success 200 {object} service.SystemSnapshot
// @Router /api/v2/time-travel/snapshot/create [post]
func (h *TimeTravelHandler) CreateSnapshot(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
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

	result, err := h.timeTravelService.CreateSnapshot(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetSnapshot 獲取快照
// @Summary 獲取系統快照
// @Tags TimeTravel
// @Produce json
// @Param snapshotId path string true "快照 ID"
// @Success 200 {object} service.SystemSnapshot
// @Router /api/v2/time-travel/snapshot/{snapshotId} [get]
func (h *TimeTravelHandler) GetSnapshot(c *gin.Context) {
	snapshotID := c.Param("snapshotId")

	result, err := h.timeTravelService.GetSnapshot(c.Request.Context(), snapshotID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// CompareSnapshots 比較快照
// @Summary 比較兩個系統快照
// @Tags TimeTravel
// @Produce json
// @Param snapshot1 query string true "快照 1 ID"
// @Param snapshot2 query string true "快照 2 ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/time-travel/snapshot/compare [get]
func (h *TimeTravelHandler) CompareSnapshots(c *gin.Context) {
	snapshot1 := c.Query("snapshot1")
	snapshot2 := c.Query("snapshot2")

	if snapshot1 == "" || snapshot2 == "" {
		handleError(c, apperrors.New(
			apperrors.ErrCodeValidation,
			"Both snapshot IDs are required",
			http.StatusBadRequest,
		))
		return
	}

	result, err := h.timeTravelService.CompareSnapshots(c.Request.Context(), snapshot1, snapshot2)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// WhatIfAnalysis What-If 分析
// @Summary 執行 What-If 分析
// @Tags TimeTravel
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "場景參數"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/time-travel/what-if-analysis [post]
func (h *TimeTravelHandler) WhatIfAnalysis(c *gin.Context) {
	var req struct {
		Scenario   string                 `json:"scenario" binding:"required"`
		Parameters map[string]interface{} `json:"parameters"`
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

	result, err := h.timeTravelService.WhatIfAnalysis(c.Request.Context(), req.Scenario, req.Parameters)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}


