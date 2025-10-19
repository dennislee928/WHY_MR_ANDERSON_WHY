package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	apperrors "axiom-backend/internal/errors"
	"axiom-backend/internal/service"
)

// AgentPracticalHandler Agent 實用功能處理器
type AgentPracticalHandler struct {
	agentPracticalService *service.AgentPracticalService
}

// NewAgentPracticalHandler 創建 Agent 實用處理器
func NewAgentPracticalHandler(agentPracticalService *service.AgentPracticalService) *AgentPracticalHandler {
	return &AgentPracticalHandler{
		agentPracticalService: agentPracticalService,
	}
}

// DiscoverAssets 資產發現
// @Summary Agent 資產發現
// @Tags Agent Practical
// @Accept json
// @Produce json
// @Param request body map[string]string true "發現請求"
// @Success 200 {object} service.AssetDiscoveryResult
// @Router /api/v2/agent/practical/discover-assets [post]
func (h *AgentPracticalHandler) DiscoverAssets(c *gin.Context) {
	var req struct {
		AgentID  string `json:"agent_id" binding:"required"`
		ScanType string `json:"scan_type"` // full, quick, network
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

	if req.ScanType == "" {
		req.ScanType = "quick"
	}

	result, err := h.agentPracticalService.DiscoverAssets(c.Request.Context(), req.AgentID, req.ScanType)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// CheckCompliance 合規性檢查
// @Summary Agent 合規性檢查
// @Tags Agent Practical
// @Accept json
// @Produce json
// @Param request body map[string]string true "檢查請求"
// @Success 200 {object} service.ComplianceCheckResult
// @Router /api/v2/agent/practical/check-compliance [post]
func (h *AgentPracticalHandler) CheckCompliance(c *gin.Context) {
	var req struct {
		AgentID   string `json:"agent_id" binding:"required"`
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

	result, err := h.agentPracticalService.CheckCompliance(c.Request.Context(), req.AgentID, req.Framework)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ExecuteRemoteCommand 遠端執行指令
// @Summary Agent 遠端執行指令
// @Tags Agent Practical
// @Accept json
// @Produce json
// @Param request body map[string]string true "執行請求"
// @Success 200 {object} service.RemoteCommandResult
// @Router /api/v2/agent/practical/execute-command [post]
func (h *AgentPracticalHandler) ExecuteRemoteCommand(c *gin.Context) {
	var req struct {
		AgentID string `json:"agent_id" binding:"required"`
		Command string `json:"command" binding:"required"`
		User    string `json:"user"`
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

	if req.User == "" {
		req.User = "system"
	}

	result, err := h.agentPracticalService.ExecuteRemoteCommand(c.Request.Context(), req.AgentID, req.Command, req.User)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetExecutionStatus 獲取執行狀態
// @Summary 獲取命令執行狀態
// @Tags Agent Practical
// @Produce json
// @Param executionId path string true "執行 ID"
// @Success 200 {object} service.RemoteCommandResult
// @Router /api/v2/agent/practical/execution/{executionId} [get]
func (h *AgentPracticalHandler) GetExecutionStatus(c *gin.Context) {
	executionID := c.Param("executionId")

	result, err := h.agentPracticalService.GetExecutionStatus(c.Request.Context(), executionID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

