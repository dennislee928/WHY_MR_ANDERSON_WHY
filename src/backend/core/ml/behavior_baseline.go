package ml

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// BehaviorBaseline implements behavioral baseline modeling
type BehaviorBaseline struct {
	profiles       map[string]*UserProfile
	globalBaseline *GlobalBaseline
	mu             sync.RWMutex
	logger         *logrus.Logger
	learningPeriod time.Duration
	updateInterval time.Duration
}

// UserProfile represents a user's behavioral baseline
type UserProfile struct {
	UserID         string
	CreatedAt      time.Time
	LastUpdated    time.Time
	
	// Request patterns
	AvgRequestRate float64
	StdDevRequestRate float64
	PeakRequestRate float64
	
	// Temporal patterns
	ActiveHours    []int // Hours of day when user is typically active
	ActiveDays     []int // Days of week when user is typically active
	AvgSessionDuration float64
	
	// Access patterns
	CommonEndpoints map[string]int
	CommonUserAgents map[string]int
	CommonIPs       map[string]int
	
	// Geographic patterns
	CommonCountries map[string]int
	CommonCities    map[string]int
	
	// Behavioral metrics
	AvgResponseTime float64
	ErrorRate       float64
	SuccessRate     float64
	
	// Statistical data
	TotalRequests   int64
	TotalSessions   int64
	
	// Anomaly scores
	AnomalyScores   []float64
	LastAnomalyTime time.Time
}

// GlobalBaseline represents global system baseline
type GlobalBaseline struct {
	AvgRequestsPerSecond float64
	AvgActiveUsers       float64
	AvgErrorRate         float64
	PeakHours            []int
	CommonProtocols      map[string]float64
	LastUpdated          time.Time
}

// AnomalyDetection represents an anomaly detection result
type AnomalyDetection struct {
	UserID         string
	AnomalyScore   float64
	IsAnomaly      bool
	Deviations     []Deviation
	Severity       string
	DetectedAt     time.Time
}

// Deviation represents a specific deviation from baseline
type Deviation struct {
	Metric      string
	Expected    float64
	Actual      float64
	Deviation   float64
	Severity    string
}

// NewBehaviorBaseline creates a new behavior baseline system
func NewBehaviorBaseline(logger *logrus.Logger) *BehaviorBaseline {
	if logger == nil {
		logger = logrus.New()
	}

	return &BehaviorBaseline{
		profiles:       make(map[string]*UserProfile),
		globalBaseline: &GlobalBaseline{
			CommonProtocols: make(map[string]float64),
			PeakHours:       make([]int, 0),
		},
		logger:         logger,
		learningPeriod: 7 * 24 * time.Hour, // 7 days learning period
		updateInterval: 1 * time.Hour,
	}
}

// LearnUserBehavior learns a user's behavioral pattern
func (bb *BehaviorBaseline) LearnUserBehavior(ctx context.Context, userID string, metrics *UserMetrics) error {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	profile, exists := bb.profiles[userID]
	if !exists {
		profile = &UserProfile{
			UserID:           userID,
			CreatedAt:        time.Now(),
			CommonEndpoints:  make(map[string]int),
			CommonUserAgents: make(map[string]int),
			CommonIPs:        make(map[string]int),
			CommonCountries:  make(map[string]int),
			CommonCities:     make(map[string]int),
			AnomalyScores:    make([]float64, 0),
		}
		bb.profiles[userID] = profile
	}

	// 更新統計數據
	bb.updateProfile(profile, metrics)

	// 計算基線指標
	bb.computeBaseline(profile)

	profile.LastUpdated = time.Now()

	bb.logger.Debugf("Updated behavior baseline for user %s", userID)
	return nil
}

// DetectAnomaly detects anomalies in user behavior
func (bb *BehaviorBaseline) DetectAnomaly(ctx context.Context, userID string, metrics *UserMetrics) (*AnomalyDetection, error) {
	bb.mu.RLock()
	profile, exists := bb.profiles[userID]
	bb.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no baseline profile found for user %s", userID)
	}

	// 檢查學習期是否完成
	if time.Since(profile.CreatedAt) < bb.learningPeriod {
		return &AnomalyDetection{
			UserID:     userID,
			IsAnomaly:  false,
			DetectedAt: time.Now(),
		}, nil
	}

	detection := &AnomalyDetection{
		UserID:     userID,
		Deviations: make([]Deviation, 0),
		DetectedAt: time.Now(),
	}

	// 檢測各種偏差
	bb.checkRequestRateDeviation(profile, metrics, detection)
	bb.checkTemporalDeviation(profile, metrics, detection)
	bb.checkAccessPatternDeviation(profile, metrics, detection)
	bb.checkGeographicDeviation(profile, metrics, detection)
	bb.checkErrorRateDeviation(profile, metrics, detection)

	// 計算總體異常分數
	detection.AnomalyScore = bb.calculateAnomalyScore(detection.Deviations)

	// 判斷是否為異常
	detection.IsAnomaly = detection.AnomalyScore > 0.7

	// 確定嚴重程度
	detection.Severity = bb.determineSeverity(detection.AnomalyScore)

	if detection.IsAnomaly {
		bb.logger.Warnf("Anomaly detected for user %s: score=%.2f, severity=%s", 
			userID, detection.AnomalyScore, detection.Severity)
	}

	return detection, nil
}

// updateProfile updates user profile with new metrics
func (bb *BehaviorBaseline) updateProfile(profile *UserProfile, metrics *UserMetrics) {
	profile.TotalRequests += metrics.RequestCount
	profile.TotalSessions += metrics.SessionCount

	// 更新請求率（移動平均）
	alpha := 0.1 // 平滑因子
	profile.AvgRequestRate = alpha*metrics.RequestRate + (1-alpha)*profile.AvgRequestRate

	// 更新端點訪問
	for endpoint, count := range metrics.Endpoints {
		profile.CommonEndpoints[endpoint] += count
	}

	// 更新 User-Agent
	for ua, count := range metrics.UserAgents {
		profile.CommonUserAgents[ua] += count
	}

	// 更新 IP 地址
	for ip, count := range metrics.IPs {
		profile.CommonIPs[ip] += count
	}

	// 更新地理位置
	for country, count := range metrics.Countries {
		profile.CommonCountries[country] += count
	}

	// 更新錯誤率
	if metrics.TotalRequests > 0 {
		newErrorRate := float64(metrics.ErrorCount) / float64(metrics.TotalRequests)
		profile.ErrorRate = alpha*newErrorRate + (1-alpha)*profile.ErrorRate
		profile.SuccessRate = 1.0 - profile.ErrorRate
	}
}

// computeBaseline computes baseline metrics
func (bb *BehaviorBaseline) computeBaseline(profile *UserProfile) {
	// 計算請求率標準差
	if len(profile.AnomalyScores) > 10 {
		mean := profile.AvgRequestRate
		variance := 0.0
		
		for _, score := range profile.AnomalyScores {
			diff := score - mean
			variance += diff * diff
		}
		
		profile.StdDevRequestRate = math.Sqrt(variance / float64(len(profile.AnomalyScores)))
	}

	// 識別活躍時段
	profile.ActiveHours = bb.identifyActiveHours(profile)
	profile.ActiveDays = bb.identifyActiveDays(profile)
}

// checkRequestRateDeviation checks for request rate anomalies
func (bb *BehaviorBaseline) checkRequestRateDeviation(profile *UserProfile, metrics *UserMetrics, detection *AnomalyDetection) {
	if profile.AvgRequestRate == 0 {
		return
	}

	deviation := math.Abs(metrics.RequestRate-profile.AvgRequestRate) / profile.AvgRequestRate

	if deviation > 2.0 { // 超過 200% 偏差
		detection.Deviations = append(detection.Deviations, Deviation{
			Metric:    "request_rate",
			Expected:  profile.AvgRequestRate,
			Actual:    metrics.RequestRate,
			Deviation: deviation,
			Severity:  bb.getDeviationSeverity(deviation),
		})
	}
}

// checkTemporalDeviation checks for temporal anomalies
func (bb *BehaviorBaseline) checkTemporalDeviation(profile *UserProfile, metrics *UserMetrics, detection *AnomalyDetection) {
	currentHour := time.Now().Hour()
	isActiveHour := false

	for _, hour := range profile.ActiveHours {
		if hour == currentHour {
			isActiveHour = true
			break
		}
	}

	// 如果在非活躍時段有大量活動
	if !isActiveHour && metrics.RequestRate > profile.AvgRequestRate*0.5 {
		detection.Deviations = append(detection.Deviations, Deviation{
			Metric:    "temporal_pattern",
			Expected:  0.0,
			Actual:    metrics.RequestRate,
			Deviation: 1.0,
			Severity:  "medium",
		})
	}
}

// checkAccessPatternDeviation checks for access pattern anomalies
func (bb *BehaviorBaseline) checkAccessPatternDeviation(profile *UserProfile, metrics *UserMetrics, detection *AnomalyDetection) {
	// 檢查是否訪問了不常見的端點
	uncommonEndpoints := 0
	for endpoint := range metrics.Endpoints {
		if _, exists := profile.CommonEndpoints[endpoint]; !exists {
			uncommonEndpoints++
		}
	}

	if uncommonEndpoints > 5 {
		detection.Deviations = append(detection.Deviations, Deviation{
			Metric:    "access_pattern",
			Expected:  0.0,
			Actual:    float64(uncommonEndpoints),
			Deviation: float64(uncommonEndpoints) / 5.0,
			Severity:  "medium",
		})
	}
}

// checkGeographicDeviation checks for geographic anomalies
func (bb *BehaviorBaseline) checkGeographicDeviation(profile *UserProfile, metrics *UserMetrics, detection *AnomalyDetection) {
	// 檢查是否從不常見的國家訪問
	for country := range metrics.Countries {
		if _, exists := profile.CommonCountries[country]; !exists {
			detection.Deviations = append(detection.Deviations, Deviation{
				Metric:    "geographic_location",
				Expected:  0.0,
				Actual:    1.0,
				Deviation: 1.0,
				Severity:  "high",
			})
			break
		}
	}
}

// checkErrorRateDeviation checks for error rate anomalies
func (bb *BehaviorBaseline) checkErrorRateDeviation(profile *UserProfile, metrics *UserMetrics, detection *AnomalyDetection) {
	if metrics.TotalRequests == 0 {
		return
	}

	currentErrorRate := float64(metrics.ErrorCount) / float64(metrics.TotalRequests)
	
	if currentErrorRate > profile.ErrorRate*2 {
		detection.Deviations = append(detection.Deviations, Deviation{
			Metric:    "error_rate",
			Expected:  profile.ErrorRate,
			Actual:    currentErrorRate,
			Deviation: currentErrorRate / profile.ErrorRate,
			Severity:  "high",
		})
	}
}

// calculateAnomalyScore calculates overall anomaly score
func (bb *BehaviorBaseline) calculateAnomalyScore(deviations []Deviation) float64 {
	if len(deviations) == 0 {
		return 0.0
	}

	totalScore := 0.0
	weights := map[string]float64{
		"request_rate":        0.25,
		"temporal_pattern":    0.15,
		"access_pattern":      0.20,
		"geographic_location": 0.25,
		"error_rate":          0.15,
	}

	for _, dev := range deviations {
		weight := weights[dev.Metric]
		severityMultiplier := bb.getSeverityMultiplier(dev.Severity)
		score := math.Min(dev.Deviation, 3.0) / 3.0 * weight * severityMultiplier
		totalScore += score
	}

	return math.Min(totalScore, 1.0)
}

// Helper functions

func (bb *BehaviorBaseline) identifyActiveHours(profile *UserProfile) []int {
	// 簡化版本：假設所有時段都活躍
	// 實際實現需要分析歷史數據
	return []int{9, 10, 11, 12, 13, 14, 15, 16, 17}
}

func (bb *BehaviorBaseline) identifyActiveDays(profile *UserProfile) []int {
	// 簡化版本：工作日
	return []int{1, 2, 3, 4, 5}
}

func (bb *BehaviorBaseline) getDeviationSeverity(deviation float64) string {
	if deviation > 5.0 {
		return "critical"
	} else if deviation > 3.0 {
		return "high"
	} else if deviation > 2.0 {
		return "medium"
	}
	return "low"
}

func (bb *BehaviorBaseline) getSeverityMultiplier(severity string) float64 {
	switch severity {
	case "critical":
		return 2.0
	case "high":
		return 1.5
	case "medium":
		return 1.0
	case "low":
		return 0.5
	default:
		return 1.0
	}
}

func (bb *BehaviorBaseline) determineSeverity(score float64) string {
	if score > 0.9 {
		return "critical"
	} else if score > 0.8 {
		return "high"
	} else if score > 0.7 {
		return "medium"
	}
	return "low"
}

// UserMetrics represents current user metrics
type UserMetrics struct {
	RequestRate    float64
	RequestCount   int64
	SessionCount   int64
	TotalRequests  int64
	ErrorCount     int64
	Endpoints      map[string]int
	UserAgents     map[string]int
	IPs            map[string]int
	Countries      map[string]int
	Cities         map[string]int
}

// GetProfile returns a user's profile
func (bb *BehaviorBaseline) GetProfile(userID string) (*UserProfile, error) {
	bb.mu.RLock()
	defer bb.mu.RUnlock()

	profile, exists := bb.profiles[userID]
	if !exists {
		return nil, fmt.Errorf("profile not found for user %s", userID)
	}

	return profile, nil
}

// GetGlobalBaseline returns the global baseline
func (bb *BehaviorBaseline) GetGlobalBaseline() *GlobalBaseline {
	bb.mu.RLock()
	defer bb.mu.RUnlock()

	return bb.globalBaseline
}

// UpdateGlobalBaseline updates the global baseline
func (bb *BehaviorBaseline) UpdateGlobalBaseline(ctx context.Context) error {
	bb.mu.Lock()
	defer bb.mu.Unlock()

	// 計算全局統計
	totalUsers := len(bb.profiles)
	if totalUsers == 0 {
		return nil
	}

	totalRequestRate := 0.0
	totalErrorRate := 0.0

	for _, profile := range bb.profiles {
		totalRequestRate += profile.AvgRequestRate
		totalErrorRate += profile.ErrorRate
	}

	bb.globalBaseline.AvgRequestsPerSecond = totalRequestRate
	bb.globalBaseline.AvgActiveUsers = float64(totalUsers)
	bb.globalBaseline.AvgErrorRate = totalErrorRate / float64(totalUsers)
	bb.globalBaseline.LastUpdated = time.Now()

	bb.logger.Debugf("Updated global baseline: users=%d, avg_rps=%.2f", 
		totalUsers, bb.globalBaseline.AvgRequestsPerSecond)

	return nil
}

