package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"pandora_box_console_ids_ips/internal/token"
)

// Client Grafana客戶端
type Client struct {
	logger     *logrus.Logger
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient 建立新的Grafana客戶端
func NewClient(logger *logrus.Logger) *Client {
	return &Client{
		logger: logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Initialize 初始化Grafana客戶端
func (c *Client) Initialize(baseURL string) error {
	c.baseURL = baseURL
	c.apiKey = "your-grafana-api-key" // 應該從設定檔或環境變數讀取

	c.logger.Infof("Grafana客戶端初始化完成: %s", baseURL)
	return nil
}

// LogEvent 記錄事件到Grafana
func (c *Client) LogEvent(eventType, message string, data map[string]interface{}) error {
	event := map[string]interface{}{
		"timestamp":  time.Now().Format(time.RFC3339),
		"event_type": eventType,
		"message":    message,
		"data":       data,
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("序列化事件資料失敗: %v", err)
	}

	// 發送到Grafana API
	url := fmt.Sprintf("%s/api/annotations", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("建立HTTP請求失敗: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("發送HTTP請求失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Grafana API回應錯誤: %d", resp.StatusCode)
	}

	c.logger.Infof("已記錄事件到Grafana: %s", eventType)
	return nil
}

// VerifyAgent 驗證Agent
func (c *Client) VerifyAgent(tokenAuth *token.Auth) error {
	c.logger.Info("開始向Grafana伺服器驗證Agent...")

	// 準備驗證資料
	verificationData := map[string]interface{}{
		"timestamp":     time.Now().Format(time.RFC3339),
		"agent_id":      "internet-blocker-agent",
		"token_valid":   tokenAuth.IsEnabled() && tokenAuth.ValidateToken(),
		"device_status": "connected",
	}

	// 發送驗證請求
	jsonData, err := json.Marshal(verificationData)
	if err != nil {
		return fmt.Errorf("序列化驗證資料失敗: %v", err)
	}

	url := fmt.Sprintf("%s/api/agent/verify", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("建立驗證請求失敗: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("發送驗證請求失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Grafana驗證失敗: %d", resp.StatusCode)
	}

	c.logger.Info("Agent驗證成功")
	return nil
}

// SendNetworkStatus 發送網路狀態到Grafana
func (c *Client) SendNetworkStatus(isBlocked bool, reason string) error {
	statusData := map[string]interface{}{
		"timestamp":  time.Now().Format(time.RFC3339),
		"is_blocked": isBlocked,
		"reason":     reason,
		"agent_id":   "internet-blocker-agent",
	}

	return c.LogEvent("network_status", "網路狀態更新", statusData)
}

// SendPinCodeEvent 發送PIN碼事件到Grafana
func (c *Client) SendPinCodeEvent(eventType, pinCode string, success bool) error {
	pinData := map[string]interface{}{
		"timestamp":  time.Now().Format(time.RFC3339),
		"event_type": eventType,
		"pin_code":   pinCode,
		"success":    success,
		"agent_id":   "internet-blocker-agent",
	}

	return c.LogEvent("pin_code_event", fmt.Sprintf("PIN碼事件: %s", eventType), pinData)
}

// SendTokenEvent 發送Token事件到Grafana
func (c *Client) SendTokenEvent(eventType string, success bool) error {
	tokenData := map[string]interface{}{
		"timestamp":  time.Now().Format(time.RFC3339),
		"event_type": eventType,
		"success":    success,
		"agent_id":   "internet-blocker-agent",
	}

	return c.LogEvent("token_event", fmt.Sprintf("Token事件: %s", eventType), tokenData)
}

// GetAgentStatus 取得Agent狀態
func (c *Client) GetAgentStatus() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/agent/status", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("建立狀態請求失敗: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("發送狀態請求失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("取得Agent狀態失敗: %d", resp.StatusCode)
	}

	var status map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("解析狀態回應失敗: %v", err)
	}

	return status, nil
}

// SendHeartbeat 發送心跳到Grafana
func (c *Client) SendHeartbeat() error {
	heartbeatData := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"agent_id":  "internet-blocker-agent",
		"status":    "alive",
	}

	return c.LogEvent("heartbeat", "Agent心跳", heartbeatData)
}

// CreateDashboard 建立儀表板
func (c *Client) CreateDashboard() error {
	dashboard := map[string]interface{}{
		"dashboard": map[string]interface{}{
			"title": "Internet Blocker Agent",
			"panels": []map[string]interface{}{
				{
					"title": "網路狀態",
					"type":  "stat",
					"targets": []map[string]interface{}{
						{
							"expr": "internet_blocker_network_status",
						},
					},
				},
				{
					"title": "PIN碼事件",
					"type":  "graph",
					"targets": []map[string]interface{}{
						{
							"expr": "rate(internet_blocker_pin_events_total[5m])",
						},
					},
				},
				{
					"title": "Token驗證",
					"type":  "stat",
					"targets": []map[string]interface{}{
						{
							"expr": "internet_blocker_token_valid",
						},
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(dashboard)
	if err != nil {
		return fmt.Errorf("序列化儀表板資料失敗: %v", err)
	}

	url := fmt.Sprintf("%s/api/dashboards/db", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("建立儀表板請求失敗: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("發送儀表板請求失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("建立儀表板失敗: %d", resp.StatusCode)
	}

	c.logger.Info("Grafana儀表板建立成功")
	return nil
}
