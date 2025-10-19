package dto

// PrometheusQueryRequest Prometheus 查詢請求
type PrometheusQueryRequest struct {
	Query string `json:"query" binding:"required"` // PromQL 查詢
	Time  string `json:"time"`                     // RFC3339, 可選
}

// PrometheusQueryRangeRequest Prometheus 範圍查詢請求
type PrometheusQueryRangeRequest struct {
	Query string `json:"query" binding:"required"` // PromQL 查詢
	Start string `json:"start" binding:"required"` // RFC3339
	End   string `json:"end" binding:"required"`   // RFC3339
	Step  string `json:"step"`                     // 步長，如 "15s", "1m"
}

// PrometheusAlertRuleRequest Prometheus 告警規則請求
type PrometheusAlertRuleRequest struct {
	Name        string            `json:"name" binding:"required"`
	Expr        string            `json:"expr" binding:"required"` // PromQL expression
	For         string            `json:"for"`                     // Duration, 如 "5m"
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Severity    string            `json:"severity"` // critical, warning, info
}

// PrometheusScrapeTargetRequest Prometheus 抓取目標請求
type PrometheusScrapeTargetRequest struct {
	JobName       string            `json:"job_name" binding:"required"`
	StaticConfigs []StaticConfig    `json:"static_configs" binding:"required,min=1"`
	ScrapeInterval string           `json:"scrape_interval,omitempty"` // 如 "30s"
	ScrapeTimeout  string           `json:"scrape_timeout,omitempty"`  // 如 "10s"
	MetricsPath    string           `json:"metrics_path,omitempty"`    // 默認 "/metrics"
}

// StaticConfig 靜態配置
type StaticConfig struct {
	Targets []string          `json:"targets" binding:"required,min=1"` // ["host:port"]
	Labels  map[string]string `json:"labels,omitempty"`
}

