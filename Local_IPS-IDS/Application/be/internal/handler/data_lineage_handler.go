package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	apperrors "axiom-backend/internal/errors"
	"axiom-backend/internal/service"
)

// DataLineageHandler 資料血緣處理器
type DataLineageHandler struct {
	dataLineageService *service.DataLineageService
}

// NewDataLineageHandler 創建資料血緣處理器
func NewDataLineageHandler(dataLineageService *service.DataLineageService) *DataLineageHandler {
	return &DataLineageHandler{
		dataLineageService: dataLineageService,
	}
}

// TraceDataLineage 追蹤資料血緣
// @Summary 追蹤資料血緣
// @Tags Data Lineage
// @Accept json
// @Produce json
// @Param request body map[string]string true "追蹤請求"
// @Success 200 {object} service.DataLineageTrace
// @Router /api/v2/data-lineage/trace [post]
func (h *DataLineageHandler) TraceDataLineage(c *gin.Context) {
	var req struct {
		DataAsset string `json:"data_asset" binding:"required"`
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

	result, err := h.dataLineageService.TraceDataLineage(c.Request.Context(), req.DataAsset)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// AnalyzeImpact 分析影響
// @Summary 分析變更影響
// @Tags Data Lineage
// @Accept json
// @Produce json
// @Param request body map[string]string true "影響分析請求"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/data-lineage/impact-analysis [post]
func (h *DataLineageHandler) AnalyzeImpact(c *gin.Context) {
	var req struct {
		DataAsset  string `json:"data_asset" binding:"required"`
		ChangeType string `json:"change_type" binding:"required"`
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

	result, err := h.dataLineageService.AnalyzeImpact(c.Request.Context(), req.DataAsset, req.ChangeType)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}


