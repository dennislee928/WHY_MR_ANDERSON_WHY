package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// EnhancedClient 增強版Grafana客戶端
type EnhancedClient struct {
	*Client
	dashboardManager  *DashboardManager
	alertManager      *AlertManager
	dataSourceManager *DataSourceManager
}

// DashboardManager 儀表板管理器
type DashboardManager struct {
	client *EnhancedClient
}

// AlertManager 告警管理器
type AlertManager struct {
	client *EnhancedClient
}

// DataSourceManager 資料源管理器
type DataSourceManager struct {
	client *EnhancedClient
}

// Dashboard 儀表板結構
type Dashboard struct {
	ID            int                    `json:"id"`
	UID           string                 `json:"uid"`
	Title         string                 `json:"title"`
	Tags          []string               `json:"tags"`
	Timezone      string                 `json:"timezone"`
	Panels        []Panel                `json:"panels"`
	Time          map[string]interface{} `json:"time"`
	Templating    map[string]interface{} `json:"templating"`
	Annotations   map[string]interface{} `json:"annotations"`
	Refresh       string                 `json:"refresh"`
	SchemaVersion int                    `json:"schemaVersion"`
	Version       int                    `json:"version"`
}

// Panel 儀表板面板
type Panel struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Type        string                 `json:"type"`
	Targets     []Target               `json:"targets"`
	GridPos     map[string]int         `json:"gridPos"`
	Options     map[string]interface{} `json:"options"`
	FieldConfig map[string]interface{} `json:"fieldConfig"`
}

// Target 查詢目標
type Target struct {
	Expr         string `json:"expr"`
	LegendFormat string `json:"legendFormat"`
	RefID        string `json:"refId"`
}

// Alert 告警結構
type Alert struct {
	ID                  int              `json:"id"`
	Name                string           `json:"name"`
	Message             string           `json:"message"`
	Frequency           string           `json:"frequency"`
	Conditions          []AlertCondition `json:"conditions"`
	ExecutionErrorState string           `json:"executionErrorState"`
	NoDataState         string           `json:"noDataState"`
	For                 string           `json:"for"`
}

// AlertCondition 告警條件
type AlertCondition struct {
	Query     AlertQuery             `json:"query"`
	Reducer   map[string]interface{} `json:"reducer"`
	Evaluator map[string]interface{} `json:"evaluator"`
}

// AlertQuery 告警查詢
type AlertQuery struct {
	QueryType string                 `json:"queryType"`
	RefID     string                 `json:"refId"`
	Model     map[string]interface{} `json:"model"`
}

// DataSource 資料源結構
type DataSource struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	URL      string                 `json:"url"`
	Access   string                 `json:"access"`
	Database string                 `json:"database"`
	JSONData map[string]interface{} `json:"jsonData"`
}

// NewEnhancedClient 建立增強版Grafana客戶端
func NewEnhancedClient(logger *logrus.Logger) *EnhancedClient {
	baseClient := NewClient(logger)
	enhanced := &EnhancedClient{
		Client: baseClient,
	}

	enhanced.dashboardManager = &DashboardManager{client: enhanced}
	enhanced.alertManager = &AlertManager{client: enhanced}
	enhanced.dataSourceManager = &DataSourceManager{client: enhanced}

	return enhanced
}

// Initialize 初始化增強版客戶端
func (ec *EnhancedClient) Initialize(baseURL, apiKey string) error {
	ec.baseURL = baseURL
	ec.apiKey = apiKey

	// 測試連接
	if err := ec.testConnection(); err != nil {
		return fmt.Errorf("Grafana連接測試失敗: %v", err)
	}

	ec.logger.Infof("增強版Grafana客戶端初始化完成: %s", baseURL)
	return nil
}

// testConnection 測試連接
func (ec *EnhancedClient) testConnection() error {
	url := fmt.Sprintf("%s/api/health", ec.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+ec.apiKey)

	resp, err := ec.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Grafana健康檢查失敗: %d", resp.StatusCode)
	}

	return nil
}

// SetupPandoraBoxDashboards 設置Pandora Box專用儀表板
func (ec *EnhancedClient) SetupPandoraBoxDashboards() error {
	dashboards := []Dashboard{
		ec.createSystemOverviewDashboard(),
		ec.createSecurityDashboard(),
		ec.createNetworkDashboard(),
		ec.createPerformanceDashboard(),
	}

	for _, dashboard := range dashboards {
		if err := ec.dashboardManager.CreateOrUpdateDashboard(dashboard); err != nil {
			ec.logger.Errorf("建立儀表板失敗 (%s): %v", dashboard.Title, err)
		} else {
			ec.logger.Infof("儀表板建立成功: %s", dashboard.Title)
		}
	}

	return nil
}

// createSystemOverviewDashboard 建立系統總覽儀表板
func (ec *EnhancedClient) createSystemOverviewDashboard() Dashboard {
	return Dashboard{
		UID:     "pandora-system-overview",
		Title:   "Pandora Box - 系統總覽",
		Tags:    []string{"pandora", "system", "overview"},
		Refresh: "5s",
		Panels: []Panel{
			{
				ID:    1,
				Title: "系統狀態",
				Type:  "stat",
				Targets: []Target{
					{
						Expr:         "pandora_system_uptime_seconds",
						LegendFormat: "運行時間",
						RefID:        "A",
					},
				},
				GridPos: map[string]int{"h": 8, "w": 12, "x": 0, "y": 0},
			},
			{
				ID:    2,
				Title: "網路狀態",
				Type:  "stat",
				Targets: []Target{
					{
						Expr:         "pandora_network_blocked",
						LegendFormat: "網路阻斷狀態",
						RefID:        "A",
					},
				},
				GridPos: map[string]int{"h": 8, "w": 12, "x": 12, "y": 0},
			},
			{
				ID:    3,
				Title: "裝置連接狀態",
				Type:  "stat",
				Targets: []Target{
					{
						Expr:         "pandora_device_connected",
						LegendFormat: "裝置連接",
						RefID:        "A",
					},
				},
				GridPos: map[string]int{"h": 8, "w": 24, "x": 0, "y": 8},
			},
		},
		Time: map[string]interface{}{
			"from": "now-1h",
			"to":   "now",
		},
		SchemaVersion: 27,
	}
}

// createSecurityDashboard 建立安全儀表板
func (ec *EnhancedClient) createSecurityDashboard() Dashboard {
	return Dashboard{
		UID:     "pandora-security",
		Title:   "Pandora Box - 安全監控",
		Tags:    []string{"pandora", "security", "threats"},
		Refresh: "10s",
		Panels: []Panel{
			{
				ID:    1,
				Title: "威脅偵測統計",
				Type:  "graph",
				Targets: []Target{
					{
						Expr:         "rate(pandora_threats_detected_total[5m])",
						LegendFormat: "威脅偵測率",
						RefID:        "A",
					},
				},
				GridPos: map[string]int{"h": 8, "w": 12, "x": 0, "y": 0},
			},
			{
				ID:    2,
				Title: "安全事件",
				Type:  "table",
				Targets: []Target{
					{
						Expr:         "pandora_security_events_total",
						LegendFormat: "{{event_type}} - {{severity}}",
						RefID:        "A",
					},
				},
				GridPos: map[string]int{"h": 8, "w": 12, "x": 12, "y": 0},
			},
			{
				ID:    3,
				Title: "被阻斷的連線",
				Type:  "stat",
				Targets: []Target{
					{
						Expr:         "pandora_blocked_connections",
						LegendFormat: "阻斷連線數",
						RefID:        "A",
					},
				},
				GridPos: map[string]int{"h": 6, "w": 24, "x": 0, "y": 8},
			},
		},
		Time: map[string]interface{}{
			"from": "now-6h",
			"to":   "now",
		},
		SchemaVersion: 27,
	}
}

// createNetworkDashboard 建立網路儀表板
func (ec *EnhancedClient) createNetworkDashboard() Dashboard {
	return Dashboard{
		UID:     "pandora-network",
		Title:   "Pandora Box - 網路監控",
		Tags:    []string{"pandora", "network", "traffic"},
		Refresh: "5s",
		Panels: []Panel{
			{
				ID:    1,
				Title: "網路流量",
				Type:  "graph",
				Targets: []Target{
					{
						Expr:         "pandora_data_throughput_bytes_per_second",
						LegendFormat: "{{direction}}",
						RefID:        "A",
					},
				},
				GridPos: map[string]int{"h": 8, "w": 24, "x": 0, "y": 0},
			},
			{
				ID:    2,
				Title: "網路事件",
				Type:  "table",
				Targets: []Target{
					{
						Expr:         "pandora_network_events_total",
						LegendFormat: "{{event_type}} - {{status}}",
						RefID:        "A",
					},
				},
				GridPos: map[string]int{"h": 8, "w": 24, "x": 0, "y": 8},
			},
		},
		Time: map[string]interface{}{
			"from": "now-1h",
			"to":   "now",
		},
		SchemaVersion: 27,
	}
}

// createPerformanceDashboard 建立效能儀表板
func (ec *EnhancedClient) createPerformanceDashboard() Dashboard {
	return Dashboard{
		UID:     "pandora-performance",
		Title:   "Pandora Box - 效能監控",
		Tags:    []string{"pandora", "performance", "metrics"},
		Refresh: "5s",
		Panels: []Panel{
			{
				ID:    1,
				Title: "回應時間",
				Type:  "graph",
				Targets: []Target{
					{
						Expr:         "pandora_response_time_seconds",
						LegendFormat: "{{operation}} - {{status}}",
						RefID:        "A",
					},
				},
				GridPos: map[string]int{"h": 8, "w": 12, "x": 0, "y": 0},
			},
			{
				ID:    2,
				Title: "活躍會話數",
				Type:  "stat",
				Targets: []Target{
					{
						Expr:         "pandora_active_sessions",
						LegendFormat: "活躍會話",
						RefID:        "A",
					},
				},
				GridPos: map[string]int{"h": 8, "w": 12, "x": 12, "y": 0},
			},
		},
		Time: map[string]interface{}{
			"from": "now-30m",
			"to":   "now",
		},
		SchemaVersion: 27,
	}
}

// SetupDataSources 設置資料源
func (ec *EnhancedClient) SetupDataSources(prometheusURL, lokiURL string) error {
	dataSources := []DataSource{
		{
			Name:   "Prometheus",
			Type:   "prometheus",
			URL:    prometheusURL,
			Access: "proxy",
			JSONData: map[string]interface{}{
				"httpMethod": "POST",
			},
		},
		{
			Name:   "Loki",
			Type:   "loki",
			URL:    lokiURL,
			Access: "proxy",
			JSONData: map[string]interface{}{
				"maxLines": 1000,
			},
		},
	}

	for _, ds := range dataSources {
		if err := ec.dataSourceManager.CreateOrUpdateDataSource(ds); err != nil {
			ec.logger.Errorf("建立資料源失敗 (%s): %v", ds.Name, err)
		} else {
			ec.logger.Infof("資料源建立成功: %s", ds.Name)
		}
	}

	return nil
}

// SetupAlerts 設置告警
func (ec *EnhancedClient) SetupAlerts() error {
	alerts := []Alert{
		{
			Name:      "高威脅偵測告警",
			Message:   "偵測到高威脅等級的安全事件",
			Frequency: "10s",
			Conditions: []AlertCondition{
				{
					Query: AlertQuery{
						QueryType: "",
						RefID:     "A",
						Model: map[string]interface{}{
							"expr": "rate(pandora_threats_detected_total{severity=\"high\"}[1m]) > 0.1",
						},
					},
				},
			},
			ExecutionErrorState: "alerting",
			NoDataState:         "no_data",
			For:                 "1m",
		},
		{
			Name:      "裝置連接中斷告警",
			Message:   "IoT裝置連接已中斷",
			Frequency: "30s",
			Conditions: []AlertCondition{
				{
					Query: AlertQuery{
						QueryType: "",
						RefID:     "A",
						Model: map[string]interface{}{
							"expr": "pandora_device_connected == 0",
						},
					},
				},
			},
			ExecutionErrorState: "alerting",
			NoDataState:         "no_data",
			For:                 "2m",
		},
	}

	for _, alert := range alerts {
		if err := ec.alertManager.CreateOrUpdateAlert(alert); err != nil {
			ec.logger.Errorf("建立告警失敗 (%s): %v", alert.Name, err)
		} else {
			ec.logger.Infof("告警建立成功: %s", alert.Name)
		}
	}

	return nil
}

// CreateOrUpdateDashboard 建立或更新儀表板
func (dm *DashboardManager) CreateOrUpdateDashboard(dashboard Dashboard) error {
	payload := map[string]interface{}{
		"dashboard": dashboard,
		"overwrite": true,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化儀表板失敗: %v", err)
	}

	url := fmt.Sprintf("%s/api/dashboards/db", dm.client.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("建立HTTP請求失敗: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+dm.client.apiKey)

	resp, err := dm.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("發送HTTP請求失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Grafana API回應錯誤: %d", resp.StatusCode)
	}

	return nil
}

// CreateOrUpdateDataSource 建立或更新資料源
func (dsm *DataSourceManager) CreateOrUpdateDataSource(dataSource DataSource) error {
	jsonData, err := json.Marshal(dataSource)
	if err != nil {
		return fmt.Errorf("序列化資料源失敗: %v", err)
	}

	url := fmt.Sprintf("%s/api/datasources", dsm.client.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("建立HTTP請求失敗: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+dsm.client.apiKey)

	resp, err := dsm.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("發送HTTP請求失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Grafana API回應錯誤: %d", resp.StatusCode)
	}

	return nil
}

// CreateOrUpdateAlert 建立或更新告警
func (am *AlertManager) CreateOrUpdateAlert(alert Alert) error {
	jsonData, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("序列化告警失敗: %v", err)
	}

	url := fmt.Sprintf("%s/api/alerts", am.client.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("建立HTTP請求失敗: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+am.client.apiKey)

	resp, err := am.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("發送HTTP請求失敗: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Grafana API回應錯誤: %d", resp.StatusCode)
	}

	return nil
}
