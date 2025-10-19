package automation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// N8NClient implements n8n workflow automation client
type N8NClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     *logrus.Logger
}

// WorkflowExecution represents a workflow execution request
type WorkflowExecution struct {
	WorkflowID string                 `json:"workflowId"`
	Data       map[string]interface{} `json:"data"`
}

// WorkflowResponse represents the response from n8n
type WorkflowResponse struct {
	ExecutionID string                 `json:"executionId"`
	Status      string                 `json:"status"`
	Data        map[string]interface{} `json:"data"`
	Error       string                 `json:"error,omitempty"`
}

// ThreatEvent represents a threat event to be processed
type ThreatEvent struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// NewN8NClient creates a new n8n client
func NewN8NClient(baseURL, apiKey string, logger *logrus.Logger) *N8NClient {
	if logger == nil {
		logger = logrus.New()
	}

	return &N8NClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// TriggerWorkflow triggers an n8n workflow
func (c *N8NClient) TriggerWorkflow(ctx context.Context, workflowID string, data map[string]interface{}) (*WorkflowResponse, error) {
	execution := WorkflowExecution{
		WorkflowID: workflowID,
		Data:       data,
	}

	jsonData, err := json.Marshal(execution)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal execution data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", 
		fmt.Sprintf("%s/webhook/%s", c.baseURL, workflowID), 
		bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-N8N-API-KEY", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("workflow execution failed: %s", string(body))
	}

	var workflowResp WorkflowResponse
	if err := json.Unmarshal(body, &workflowResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.logger.Infof("Workflow triggered: %s (execution: %s)", workflowID, workflowResp.ExecutionID)
	return &workflowResp, nil
}

// ProcessThreatEvent processes a threat event through n8n workflow
func (c *N8NClient) ProcessThreatEvent(ctx context.Context, event *ThreatEvent) error {
	workflowID := c.getWorkflowIDForThreat(event.Type, event.Severity)
	
	data := map[string]interface{}{
		"event": event,
		"actions": c.getRecommendedActions(event),
	}

	resp, err := c.TriggerWorkflow(ctx, workflowID, data)
	if err != nil {
		return fmt.Errorf("failed to process threat event: %w", err)
	}

	c.logger.Infof("Threat event processed: %s (status: %s)", event.ID, resp.Status)
	return nil
}

// getWorkflowIDForThreat returns the appropriate workflow ID based on threat type and severity
func (c *N8NClient) getWorkflowIDForThreat(threatType, severity string) string {
	// 根據威脅類型和嚴重程度選擇工作流程
	workflows := map[string]map[string]string{
		"malware": {
			"critical": "malware-critical-response",
			"high":     "malware-high-response",
			"medium":   "malware-medium-response",
			"low":      "malware-low-response",
		},
		"intrusion": {
			"critical": "intrusion-critical-response",
			"high":     "intrusion-high-response",
			"medium":   "intrusion-medium-response",
			"low":      "intrusion-low-response",
		},
		"ddos": {
			"critical": "ddos-mitigation",
			"high":     "ddos-mitigation",
			"medium":   "ddos-monitoring",
			"low":      "ddos-monitoring",
		},
		"anomaly": {
			"critical": "anomaly-investigation",
			"high":     "anomaly-investigation",
			"medium":   "anomaly-monitoring",
			"low":      "anomaly-logging",
		},
	}

	if typeWorkflows, exists := workflows[threatType]; exists {
		if workflowID, exists := typeWorkflows[severity]; exists {
			return workflowID
		}
	}

	// 預設工作流程
	return "default-threat-response"
}

// getRecommendedActions returns recommended actions for a threat
func (c *N8NClient) getRecommendedActions(event *ThreatEvent) []string {
	actions := []string{}

	switch event.Severity {
	case "critical":
		actions = append(actions, 
			"block_ip",
			"isolate_system",
			"notify_soc",
			"create_incident",
			"start_forensics",
		)
	case "high":
		actions = append(actions,
			"block_ip",
			"notify_admin",
			"create_ticket",
			"increase_monitoring",
		)
	case "medium":
		actions = append(actions,
			"log_event",
			"notify_admin",
			"monitor_closely",
		)
	case "low":
		actions = append(actions,
			"log_event",
			"update_statistics",
		)
	}

	return actions
}

// SendNotification sends a notification through n8n
func (c *N8NClient) SendNotification(ctx context.Context, channel, message string, metadata map[string]interface{}) error {
	data := map[string]interface{}{
		"channel":  channel,
		"message":  message,
		"metadata": metadata,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	_, err := c.TriggerWorkflow(ctx, "notification-dispatcher", data)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	c.logger.Infof("Notification sent to %s", channel)
	return nil
}

// CreateIncident creates an incident ticket through n8n
func (c *N8NClient) CreateIncident(ctx context.Context, incident map[string]interface{}) (string, error) {
	resp, err := c.TriggerWorkflow(ctx, "incident-creation", incident)
	if err != nil {
		return "", fmt.Errorf("failed to create incident: %w", err)
	}

	incidentID, ok := resp.Data["incident_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid incident ID in response")
	}

	c.logger.Infof("Incident created: %s", incidentID)
	return incidentID, nil
}

// ExecuteRemediationAction executes a remediation action
func (c *N8NClient) ExecuteRemediationAction(ctx context.Context, action string, params map[string]interface{}) error {
	data := map[string]interface{}{
		"action": action,
		"params": params,
	}

	resp, err := c.TriggerWorkflow(ctx, "remediation-executor", data)
	if err != nil {
		return fmt.Errorf("failed to execute remediation: %w", err)
	}

	if resp.Status != "success" {
		return fmt.Errorf("remediation failed: %s", resp.Error)
	}

	c.logger.Infof("Remediation action executed: %s", action)
	return nil
}

