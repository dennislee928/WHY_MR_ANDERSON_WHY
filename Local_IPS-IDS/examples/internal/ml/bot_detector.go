package ml

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// BotDetector implements ML-based bot detection
type BotDetector struct {
	model          *BotDetectionModel
	featureCache   map[string]*UserFeatures
	cacheMu        sync.RWMutex
	logger         *logrus.Logger
	threshold      float64
	cacheExpiry    time.Duration
}

// BotDetectionModel represents the trained model
type BotDetectionModel struct {
	weights       map[string]float64
	bias          float64
	version       string
	trainedAt     time.Time
}

// UserFeatures represents behavioral features for bot detection
type UserFeatures struct {
	// 請求模式特徵
	RequestRate         float64   // 請求速率 (req/min)
	RequestVariance     float64   // 請求間隔變異數
	BurstScore          float64   // 突發請求分數
	
	// 行為特徵
	UserAgentEntropy    float64   // User-Agent 熵值
	PathEntropy         float64   // 訪問路徑熵值
	SessionDuration     float64   // 會話持續時間
	
	// 網路特徵
	IPReputation        float64   // IP 信譽分數
	GeoAnomalyScore     float64   // 地理位置異常分數
	
	// 瀏覽器特徵
	JavaScriptEnabled   bool      // JavaScript 是否啟用
	CookiesEnabled      bool      // Cookies 是否啟用
	ScreenResolution    string    // 螢幕解析度
	
	// 時間特徵
	TimeOfDay           float64   // 訪問時段
	DayOfWeek           float64   // 星期幾
	
	// 統計特徵
	TotalRequests       int64     // 總請求數
	ErrorRate           float64   // 錯誤率
	UniqueEndpoints     int       // 訪問的唯一端點數
	
	LastUpdate          time.Time
}

// BotPrediction represents the detection result
type BotPrediction struct {
	IsBot          bool
	Confidence     float64
	Score          float64
	Features       *UserFeatures
	Reasons        []string
	DetectedAt     time.Time
}

// NewBotDetector creates a new bot detector
func NewBotDetector(logger *logrus.Logger) *BotDetector {
	if logger == nil {
		logger = logrus.New()
	}

	// 初始化預訓練模型權重
	model := &BotDetectionModel{
		weights: map[string]float64{
			"request_rate":         0.25,
			"request_variance":     -0.15,
			"burst_score":          0.30,
			"user_agent_entropy":   -0.20,
			"path_entropy":         -0.10,
			"session_duration":     -0.15,
			"ip_reputation":        -0.25,
			"geo_anomaly":          0.20,
			"js_enabled":           -0.10,
			"cookies_enabled":      -0.10,
			"error_rate":           0.15,
			"unique_endpoints":     -0.05,
		},
		bias:      0.0,
		version:   "1.0.0",
		trainedAt: time.Now(),
	}

	return &BotDetector{
		model:        model,
		featureCache: make(map[string]*UserFeatures),
		logger:       logger,
		threshold:    0.7, // Bot 檢測閾值
		cacheExpiry:  30 * time.Minute,
	}
}

// Predict predicts if a user is a bot
func (bd *BotDetector) Predict(ctx context.Context, userID string, features *UserFeatures) (*BotPrediction, error) {
	// 更新特徵緩存
	bd.updateFeatureCache(userID, features)

	// 計算特徵向量
	featureVector := bd.extractFeatureVector(features)

	// 計算預測分數（邏輯回歸）
	score := bd.calculateScore(featureVector)

	// 應用 sigmoid 函數
	probability := sigmoid(score)

	// 判斷是否為 Bot
	isBot := probability >= bd.threshold

	// 生成原因
	reasons := bd.generateReasons(features, featureVector, probability)

	prediction := &BotPrediction{
		IsBot:      isBot,
		Confidence: probability,
		Score:      score,
		Features:   features,
		Reasons:    reasons,
		DetectedAt: time.Now(),
	}

	if isBot {
		bd.logger.Warnf("Bot detected: user=%s, confidence=%.2f, reasons=%v", 
			userID, probability, reasons)
	}

	return prediction, nil
}

// extractFeatureVector converts features to vector
func (bd *BotDetector) extractFeatureVector(features *UserFeatures) map[string]float64 {
	vector := make(map[string]float64)

	// 標準化特徵值到 [0, 1] 範圍
	vector["request_rate"] = normalizeRequestRate(features.RequestRate)
	vector["request_variance"] = normalizeVariance(features.RequestVariance)
	vector["burst_score"] = features.BurstScore
	vector["user_agent_entropy"] = features.UserAgentEntropy
	vector["path_entropy"] = features.PathEntropy
	vector["session_duration"] = normalizeSessionDuration(features.SessionDuration)
	vector["ip_reputation"] = features.IPReputation
	vector["geo_anomaly"] = features.GeoAnomalyScore
	vector["js_enabled"] = boolToFloat(features.JavaScriptEnabled)
	vector["cookies_enabled"] = boolToFloat(features.CookiesEnabled)
	vector["error_rate"] = features.ErrorRate
	vector["unique_endpoints"] = normalizeEndpoints(features.UniqueEndpoints)

	return vector
}

// calculateScore calculates the prediction score
func (bd *BotDetector) calculateScore(features map[string]float64) float64 {
	score := bd.model.bias

	for feature, value := range features {
		if weight, exists := bd.model.weights[feature]; exists {
			score += weight * value
		}
	}

	return score
}

// generateReasons generates human-readable reasons
func (bd *BotDetector) generateReasons(features *UserFeatures, vector map[string]float64, probability float64) []string {
	reasons := []string{}

	if features.RequestRate > 100 {
		reasons = append(reasons, fmt.Sprintf("High request rate: %.1f req/min", features.RequestRate))
	}

	if features.BurstScore > 0.8 {
		reasons = append(reasons, fmt.Sprintf("Burst pattern detected: %.2f", features.BurstScore))
	}

	if features.RequestVariance < 0.1 {
		reasons = append(reasons, "Low request variance (robotic pattern)")
	}

	if features.UserAgentEntropy < 0.3 {
		reasons = append(reasons, "Suspicious User-Agent")
	}

	if !features.JavaScriptEnabled {
		reasons = append(reasons, "JavaScript disabled")
	}

	if !features.CookiesEnabled {
		reasons = append(reasons, "Cookies disabled")
	}

	if features.ErrorRate > 0.5 {
		reasons = append(reasons, fmt.Sprintf("High error rate: %.1f%%", features.ErrorRate*100))
	}

	if features.IPReputation < 0.3 {
		reasons = append(reasons, "Low IP reputation")
	}

	if features.GeoAnomalyScore > 0.7 {
		reasons = append(reasons, "Geographical anomaly detected")
	}

	return reasons
}

// updateFeatureCache updates the feature cache
func (bd *BotDetector) updateFeatureCache(userID string, features *UserFeatures) {
	bd.cacheMu.Lock()
	defer bd.cacheMu.Unlock()

	features.LastUpdate = time.Now()
	bd.featureCache[userID] = features
}

// GetCachedFeatures retrieves cached features
func (bd *BotDetector) GetCachedFeatures(userID string) (*UserFeatures, bool) {
	bd.cacheMu.RLock()
	defer bd.cacheMu.RUnlock()

	features, exists := bd.featureCache[userID]
	if !exists {
		return nil, false
	}

	// 檢查是否過期
	if time.Since(features.LastUpdate) > bd.cacheExpiry {
		return nil, false
	}

	return features, true
}

// CleanupExpiredCache removes expired cache entries
func (bd *BotDetector) CleanupExpiredCache() {
	bd.cacheMu.Lock()
	defer bd.cacheMu.Unlock()

	now := time.Now()
	for userID, features := range bd.featureCache {
		if now.Sub(features.LastUpdate) > bd.cacheExpiry {
			delete(bd.featureCache, userID)
		}
	}

	bd.logger.Debugf("Cache cleanup completed, remaining entries: %d", len(bd.featureCache))
}

// StartCleanupWorker starts background cache cleanup
func (bd *BotDetector) StartCleanupWorker(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			bd.CleanupExpiredCache()
		}
	}
}

// UpdateModel updates the detection model
func (bd *BotDetector) UpdateModel(weights map[string]float64, bias float64) {
	bd.model.weights = weights
	bd.model.bias = bias
	bd.model.trainedAt = time.Now()
	bd.logger.Infof("Model updated: version=%s", bd.model.version)
}

// Helper functions

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func normalizeRequestRate(rate float64) float64 {
	// 正常用戶: 0-10 req/min, Bot: 50+ req/min
	return math.Min(rate/100.0, 1.0)
}

func normalizeVariance(variance float64) float64 {
	// 正常用戶: 高變異, Bot: 低變異
	return math.Min(variance, 1.0)
}

func normalizeSessionDuration(duration float64) float64 {
	// 標準化到 [0, 1]，假設正常會話 < 60 分鐘
	return math.Min(duration/60.0, 1.0)
}

func normalizeEndpoints(count int) float64 {
	// 標準化訪問的端點數量
	return math.Min(float64(count)/50.0, 1.0)
}

func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

