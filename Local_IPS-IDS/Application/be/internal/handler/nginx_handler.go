package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"axiom-backend/internal/dto"
	apperrors "axiom-backend/internal/errors"
	"axiom-backend/internal/service"
)

// NginxHandler Nginx 處理器
type NginxHandler struct {
	nginxService *service.NginxService
}

// NewNginxHandler 創建 Nginx 處理器
func NewNginxHandler(nginxService *service.NginxService) *NginxHandler {
	return &NginxHandler{
		nginxService: nginxService,
	}
}

// GetConfig 獲取配置
// @Summary 獲取 Nginx 配置
// @Tags Nginx
// @Produce json
// @Success 200 {object} vo.NginxConfigVO
// @Router /api/v2/nginx/config [get]
func (h *NginxHandler) GetConfig(c *gin.Context) {
	result, err := h.nginxService.GetConfig(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// UpdateConfig 更新配置
// @Summary 更新 Nginx 配置
// @Tags Nginx
// @Accept json
// @Produce json
// @Param request body dto.NginxConfigUpdateRequest true "配置更新請求"
// @Success 200 {object} vo.NginxConfigVO
// @Router /api/v2/nginx/config [put]
func (h *NginxHandler) UpdateConfig(c *gin.Context) {
	var req dto.NginxConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.nginxService.UpdateConfig(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// Reload 重載配置
// @Summary 重載 Nginx 配置
// @Tags Nginx
// @Accept json
// @Produce json
// @Success 200 {object} vo.NginxReloadVO
// @Router /api/v2/nginx/reload [post]
func (h *NginxHandler) Reload(c *gin.Context) {
	result, err := h.nginxService.Reload(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetStatus 獲取狀態
// @Summary 獲取 Nginx 狀態
// @Tags Nginx
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/nginx/status [get]
func (h *NginxHandler) GetStatus(c *gin.Context) {
	status, err := h.nginxService.GetStatus(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    status,
	})
}


