package service

import (
	"context"
	"time"
)

// ContextAwareService 情境感知服務
type ContextAwareService struct {
	// 依賴
}

// AlertRouting 告警路由決策
type AlertRouting struct {
	AlertID      string           `json:"alert_id"`
	AlertName    string           `json:"alert_name"`
	Severity     string           `json:"severity"`
	RouteTo      []string         `json:"route_to"` // 接收人列表
	EscalationPath []EscalationStep `json:"escalation_path"`
	NotificationChannels []string `json:"notification_channels"`
	Context      AlertContext     `json:"context"`
	Timestamp    time.Time        `json:"timestamp"`
}

// AlertContext 告警上下文
type AlertContext struct {
	TimeZone          string    `json:"timezone"`
	CurrentTime       time.Time `json:"current_time"`
	IsBusinessHours   bool      `json:"is_business_hours"`
	OnCallEngineers   []OnCallEngineer `json:"oncall_engineers"`
	WorkloadLevel     string    `json:"workload_level"` // low, medium, high
	SimilarIncidents  int       `json:"similar_incidents"`
}

// OnCallEngineer 值班工程師
type OnCallEngineer struct {
	Name              string    `json:"name"`
	Skills            []string  `json:"skills"`
	CurrentWorkload   int       `json:"current_workload"`
	SuccessRate       float64   `json:"success_rate"`
	Available         bool      `json:"available"`
	TimeZone          string    `json:"timezone"`
	PreferredChannel  string    `json:"preferred_channel"`
}

// EscalationStep 升級步驟
type EscalationStep struct {
	Level        int      `json:"level"`
	WaitMinutes  int      `json:"wait_minutes"`
	Recipients   []string `json:"recipients"`
	Channels     []string `json:"channels"`
}

// NewContextAwareService 創建情境感知服務
func NewContextAwareService() *ContextAwareService {
	return &ContextAwareService{}
}

// RouteAlert 智能告警路由
func (s *ContextAwareService) RouteAlert(ctx context.Context, alertID, alertName, severity string) (*AlertRouting, error) {
	// 1. 收集上下文
	context := s.gatherContext()
	
	// 2. 根據上下文選擇路由
	routeTo := s.selectRecipients(severity, context)
	
	// 3. 生成升級路徑
	escalationPath := s.generateEscalationPath(severity)
	
	// 4. 選擇通知渠道
	channels := s.selectChannels(severity, context)
	
	return &AlertRouting{
		AlertID:              alertID,
		AlertName:            alertName,
		Severity:             severity,
		RouteTo:              routeTo,
		EscalationPath:       escalationPath,
		NotificationChannels: channels,
		Context:              context,
		Timestamp:            time.Now(),
	}, nil
}

// gatherContext 收集告警上下文
func (s *ContextAwareService) gatherContext() AlertContext {
	now := time.Now()
	hour := now.Hour()
	isBusinessHours := hour >= 9 && hour < 18
	
	onCallEngineers := []OnCallEngineer{
		{
			Name:             "張工程師",
			Skills:           []string{"backend", "database", "quantum"},
			CurrentWorkload:  3,
			SuccessRate:      0.95,
			Available:        true,
			TimeZone:         "Asia/Taipei",
			PreferredChannel: "slack",
		},
		{
			Name:             "李工程師",
			Skills:           []string{"frontend", "monitoring"},
			CurrentWorkload:  1,
			SuccessRate:      0.88,
			Available:        true,
			TimeZone:         "Asia/Taipei",
			PreferredChannel: "email",
		},
	}
	
	return AlertContext{
		TimeZone:         "Asia/Taipei",
		CurrentTime:      now,
		IsBusinessHours:  isBusinessHours,
		OnCallEngineers:  onCallEngineers,
		WorkloadLevel:    "medium",
		SimilarIncidents: 3,
	}
}

// selectRecipients 選擇接收人
func (s *ContextAwareService) selectRecipients(severity string, ctx AlertContext) []string {
	recipients := []string{}
	
	// 根據嚴重程度和上下文選擇
	if severity == "critical" {
		// 嚴重告警，通知所有值班人員
		for _, eng := range ctx.OnCallEngineers {
			if eng.Available {
				recipients = append(recipients, eng.Name)
			}
		}
	} else if severity == "high" {
		// 高優先級，選擇最合適的工程師
		if len(ctx.OnCallEngineers) > 0 {
			// 選擇工作負載最低的
			bestEng := ctx.OnCallEngineers[0]
			for _, eng := range ctx.OnCallEngineers {
				if eng.CurrentWorkload < bestEng.CurrentWorkload && eng.Available {
					bestEng = eng
				}
			}
			recipients = append(recipients, bestEng.Name)
		}
	}
	
	return recipients
}

// generateEscalationPath 生成升級路徑
func (s *ContextAwareService) generateEscalationPath(severity string) []EscalationStep {
	if severity == "critical" {
		return []EscalationStep{
			{Level: 1, WaitMinutes: 5, Recipients: []string{"oncall-primary"}, Channels: []string{"phone", "sms"}},
			{Level: 2, WaitMinutes: 10, Recipients: []string{"oncall-secondary", "team-lead"}, Channels: []string{"phone", "slack"}},
			{Level: 3, WaitMinutes: 15, Recipients: []string{"manager", "director"}, Channels: []string{"phone", "email"}},
		}
	}
	
	return []EscalationStep{
		{Level: 1, WaitMinutes: 15, Recipients: []string{"oncall-primary"}, Channels: []string{"slack"}},
		{Level: 2, WaitMinutes: 30, Recipients: []string{"team-lead"}, Channels: []string{"slack", "email"}},
	}
}

// selectChannels 選擇通知渠道
func (s *ContextAwareService) selectChannels(severity string, ctx AlertContext) []string {
	channels := []string{"slack"} // 默認使用 Slack
	
	if severity == "critical" {
		channels = append(channels, "phone", "sms")
	}
	
	if !ctx.IsBusinessHours {
		channels = append(channels, "phone") // 非工作時間使用電話
	}
	
	return channels
}


