package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
	"axiom-backend/internal/agent"
	apperrors "axiom-backend/internal/errors"
)

// AgentHandler Agent 管理處理器
type AgentHandler struct {
	agentManager *agent.AgentManager
}

// NewAgentHandler 創建 Agent 處理器
func NewAgentHandler(agentManager *agent.AgentManager) *AgentHandler {
	return &AgentHandler{
		agentManager: agentManager,
	}
}

// RegisterAgent Agent 註冊
// @Summary Agent 註冊
// @Tags Agent
// @Accept json
// @Produce json
// @Param request body agent.AgentRegistrationRequest true "註冊請求"
// @Success 200 {object} agent.AgentRegistrationResponse
// @Router /api/v2/agent/register [post]
func (h *AgentHandler) RegisterAgent(c *gin.Context) {
	var req agent.AgentRegistrationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	result, err := h.agentManager.RegisterAgent(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// Heartbeat Agent 心跳
// @Summary Agent 心跳
// @Tags Agent
// @Accept json
// @Produce json
// @Param request body agent.HeartbeatData true "心跳數據"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/agent/heartbeat [post]
func (h *AgentHandler) Heartbeat(c *gin.Context) {
	var req agent.HeartbeatData

	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	err := h.agentManager.ProcessHeartbeat(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Heartbeat received",
	})
}

// GetAgent 獲取 Agent 信息
// @Summary 獲取 Agent 信息
// @Tags Agent
// @Produce json
// @Param agentId path string true "Agent ID"
// @Success 200 {object} agent.AgentInfo
// @Router /api/v2/agent/{agentId}/status [get]
func (h *AgentHandler) GetAgent(c *gin.Context) {
	agentID := c.Param("agentId")

	result, err := h.agentManager.GetAgent(c.Request.Context(), agentID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ListAgents 列出 Agents
// @Summary 列出所有 Agents
// @Tags Agent
// @Produce json
// @Param mode query string false "Agent 模式 (external/internal)"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/agent/list [get]
func (h *AgentHandler) ListAgents(c *gin.Context) {
	mode := agent.AgentMode(c.Query("mode"))

	result, err := h.agentManager.ListAgents(c.Request.Context(), mode)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"count":   len(result),
	})
}

// DeregisterAgent 註銷 Agent
// @Summary 註銷 Agent
// @Tags Agent
// @Produce json
// @Param agentId path string true "Agent ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/agent/{agentId} [delete]
func (h *AgentHandler) DeregisterAgent(c *gin.Context) {
	agentID := c.Param("agentId")

	err := h.agentManager.DeregisterAgent(c.Request.Context(), agentID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Agent deregistered successfully",
	})
}

// UpdateConfig 更新 Agent 配置
// @Summary 更新 Agent 配置
// @Tags Agent
// @Accept json
// @Produce json
// @Param agentId path string true "Agent ID"
// @Param request body agent.AgentConfig true "配置"
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/agent/{agentId}/config [put]
func (h *AgentHandler) UpdateConfig(c *gin.Context) {
	agentID := c.Param("agentId")
	var config agent.AgentConfig

	if err := c.ShouldBindJSON(&config); err != nil {
		handleError(c, apperrors.NewWithDetails(
			apperrors.ErrCodeValidation,
			"Invalid request",
			http.StatusBadRequest,
			err.Error(),
		))
		return
	}

	err := h.agentManager.UpdateAgentConfig(c.Request.Context(), agentID, &config)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Config updated successfully",
	})
}

// CheckHealth 檢查所有 Agent 健康狀態
// @Summary 檢查所有 Agent 健康狀態
// @Tags Agent
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v2/agent/health [get]
func (h *AgentHandler) CheckHealth(c *gin.Context) {
	healthStatus := h.agentManager.CheckAgentHealth(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    healthStatus,
	})
}

