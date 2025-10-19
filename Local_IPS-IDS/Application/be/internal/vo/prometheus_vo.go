package vo

import "time"

// PrometheusQueryVO Prometheus 查詢響應
type PrometheusQueryVO struct {
	Status    string        `json:"status"` // success, error
	Data      QueryDataVO   `json:"data"`
	ErrorType string        `json:"error_type,omitempty"`
	Error     string        `json:"error,omitempty"`
	Warnings  []string      `json:"warnings,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// QueryDataVO 查詢數據
type QueryDataVO struct {
	ResultType string        `json:"result_type"` // vector, matrix, scalar, string
	Result     []MetricVO    `json:"result"`
}

// MetricVO 指標數據
type MetricVO struct {
	Metric map[string]string `json:"metric"` // 標籤
	Value  []interface{}     `json:"value,omitempty"`  // [timestamp, value]
	Values [][]interface{}   `json:"values,omitempty"` // [[timestamp, value], ...]
}

// PrometheusAlertRulesVO Prometheus 告警規則列表響應
type PrometheusAlertRulesVO struct {
	Groups []RuleGroupVO `json:"groups"`
	Total  int           `json:"total"`
}

// RuleGroupVO 規則組
type RuleGroupVO struct {
	Name     string    `json:"name"`
	File     string    `json:"file"`
	Rules    []RuleVO  `json:"rules"`
	Interval string    `json:"interval"`
	Limit    int       `json:"limit,omitempty"`
}

// RuleVO 告警規則
type RuleVO struct {
	Name        string            `json:"name"`
	Query       string            `json:"query"`
	Duration    string            `json:"duration,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	State       string            `json:"state,omitempty"` // inactive, pending, firing
	Type        string            `json:"type"`            // alerting, recording
	Health      string            `json:"health,omitempty"`
	LastError   string            `json:"last_error,omitempty"`
}

// PrometheusTargetsVO Prometheus 目標列表響應
type PrometheusTargetsVO struct {
	ActiveTargets  []TargetVO `json:"active_targets"`
	DroppedTargets []TargetVO `json:"dropped_targets"`
	TotalActive    int        `json:"total_active"`
	TotalDropped   int        `json:"total_dropped"`
}

// TargetVO 抓取目標
type TargetVO struct {
	DiscoveredLabels map[string]string `json:"discovered_labels"`
	Labels           map[string]string `json:"labels"`
	ScrapePool       string            `json:"scrape_pool"`
	ScrapeURL        string            `json:"scrape_url"`
	GlobalURL        string            `json:"global_url"`
	LastError        string            `json:"last_error,omitempty"`
	LastScrape       time.Time         `json:"last_scrape"`
	LastScrapeDuration float64         `json:"last_scrape_duration"` // 秒
	Health           string            `json:"health"` // up, down, unknown
}

