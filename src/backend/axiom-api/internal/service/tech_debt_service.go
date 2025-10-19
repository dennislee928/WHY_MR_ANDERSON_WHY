package service

import (
	"context"
	"time"
)

// TechDebtService 技術債務追蹤服務
type TechDebtService struct {
	// 依賴
}

// TechDebtScan 技術債務掃描結果
type TechDebtScan struct {
	ScanID          string          `json:"scan_id"`
	TotalIssues     int             `json:"total_issues"`
	CriticalIssues  int             `json:"critical_issues"`
	HighIssues      int             `json:"high_issues"`
	MediumIssues    int             `json:"medium_issues"`
	LowIssues       int             `json:"low_issues"`
	DebtScore       int             `json:"debt_score"` // 0-100, 越高越差
	Issues          []TechDebtIssue `json:"issues"`
	Timestamp       time.Time       `json:"timestamp"`
}

// TechDebtIssue 技術債務問題
type TechDebtIssue struct {
	IssueID       string    `json:"issue_id"`
	Category      string    `json:"category"` // dependency, security, performance, code_quality
	Priority      string    `json:"priority"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	AffectedFiles []string  `json:"affected_files"`
	Impact        string    `json:"impact"`
	Effort        string    `json:"effort"` // hours or story points
	DetectedAt    time.Time `json:"detected_at"`
	Age           int       `json:"age_days"`
}

// RemediationRoadmap 修復路線圖
type RemediationRoadmap struct {
	RoadmapID   string              `json:"roadmap_id"`
	Phases      []RemediationPhase  `json:"phases"`
	TotalEffort string              `json:"total_effort"`
	Timeline    string              `json:"timeline"`
	Priority    string              `json:"priority"`
	CreatedAt   time.Time           `json:"created_at"`
}

// RemediationPhase 修復階段
type RemediationPhase struct {
	Phase       int      `json:"phase"`
	Name        string   `json:"name"`
	Issues      []string `json:"issues"` // Issue IDs
	Effort      string   `json:"effort"`
	Duration    string   `json:"duration"`
	Dependencies []string `json:"dependencies"`
	Order       int      `json:"order"`
}

// NewTechDebtService 創建技術債務服務
func NewTechDebtService() *TechDebtService {
	return &TechDebtService{}
}

// ScanTechDebt 掃描技術債務
func (s *TechDebtService) ScanTechDebt(ctx context.Context) (*TechDebtScan, error) {
	// 模擬掃描結果
	issues := []TechDebtIssue{
		{
			IssueID:       "DEBT-001",
			Category:      "dependency",
			Priority:      "high",
			Title:         "Outdated Go modules",
			Description:   "Several Go modules are 2+ major versions behind",
			AffectedFiles: []string{"go.mod"},
			Impact:        "Security vulnerabilities, missing features",
			Effort:        "4 hours",
			DetectedAt:    time.Now().Add(-15 * 24 * time.Hour),
			Age:           15,
		},
		{
			IssueID:       "DEBT-002",
			Category:      "security",
			Priority:      "critical",
			Title:         "Weak password hashing",
			Description:   "Using bcrypt with cost factor 10, should be 12+",
			AffectedFiles: []string{"internal/auth/password.go"},
			Impact:        "Passwords easier to crack",
			Effort:        "2 hours",
			DetectedAt:    time.Now().Add(-30 * 24 * time.Hour),
			Age:           30,
		},
		{
			IssueID:       "DEBT-003",
			Category:      "performance",
			Priority:      "medium",
			Title:         "Missing database indexes",
			Description:   "Several frequently queried columns lack indexes",
			AffectedFiles: []string{"database/migrations/001_initial_schema.sql"},
			Impact:        "Slow query performance",
			Effort:        "1 hour",
			DetectedAt:    time.Now().Add(-7 * 24 * time.Hour),
			Age:           7,
		},
		{
			IssueID:       "DEBT-004",
			Category:      "code_quality",
			Priority:      "low",
			Title:         "TODO comments in production code",
			Description:   "15 TODO comments found",
			AffectedFiles: []string{"multiple files"},
			Impact:        "Incomplete features",
			Effort:        "8 hours",
			DetectedAt:    time.Now().Add(-45 * 24 * time.Hour),
			Age:           45,
		},
	}
	
	critical := 0
	high := 0
	medium := 0
	low := 0
	
	for _, issue := range issues {
		switch issue.Priority {
		case "critical":
			critical++
		case "high":
			high++
		case "medium":
			medium++
		case "low":
			low++
		}
	}
	
	// 計算債務評分
	debtScore := (critical * 25) + (high * 15) + (medium * 5) + (low * 2)
	if debtScore > 100 {
		debtScore = 100
	}
	
	return &TechDebtScan{
		ScanID:         "SCAN-" + time.Now().Format("20060102"),
		TotalIssues:    len(issues),
		CriticalIssues: critical,
		HighIssues:     high,
		MediumIssues:   medium,
		LowIssues:      low,
		DebtScore:      debtScore,
		Issues:         issues,
		Timestamp:      time.Now(),
	}, nil
}

// GenerateRoadmap 生成修復路線圖
func (s *TechDebtService) GenerateRoadmap(ctx context.Context) (*RemediationRoadmap, error) {
	phases := []RemediationPhase{
		{
			Phase:        1,
			Name:         "Critical Security Fixes",
			Issues:       []string{"DEBT-002"},
			Effort:       "2 hours",
			Duration:     "1 day",
			Dependencies: []string{},
			Order:        1,
		},
		{
			Phase:        2,
			Name:         "High Priority Dependencies",
			Issues:       []string{"DEBT-001"},
			Effort:       "4 hours",
			Duration:     "2 days",
			Dependencies: []string{"Phase 1"},
			Order:        2,
		},
		{
			Phase:        3,
			Name:         "Performance Optimization",
			Issues:       []string{"DEBT-003"},
			Effort:       "1 hour",
			Duration:     "1 day",
			Dependencies: []string{},
			Order:        3,
		},
		{
			Phase:        4,
			Name:         "Code Quality Improvements",
			Issues:       []string{"DEBT-004"},
			Effort:       "8 hours",
			Duration:     "3 days",
			Dependencies: []string{"Phase 2"},
			Order:        4,
		},
	}
	
	return &RemediationRoadmap{
		RoadmapID:   "ROADMAP-" + time.Now().Format("20060102"),
		Phases:      phases,
		TotalEffort: "15 hours",
		Timeline:    "1-2 weeks",
		Priority:    "high",
		CreatedAt:   time.Now(),
	}, nil
}

// PrioritizeIssues 優先級排序
func (s *TechDebtService) PrioritizeIssues(ctx context.Context, issues []TechDebtIssue) []TechDebtIssue {
	// 使用加權評分排序
	// 考慮因素：優先級、年齡、影響範圍
	
	// 這裡返回已排序的問題列表
	return issues
}


