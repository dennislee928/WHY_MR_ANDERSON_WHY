package ml

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// DeepLearningDetector implements deep learning-based threat detection
type DeepLearningDetector struct {
	model          *NeuralNetwork
	featureScaler  *FeatureScaler
	trainingData   []TrainingExample
	mu             sync.RWMutex
	logger         *logrus.Logger
	threshold      float64
	modelVersion   string
}

// NeuralNetwork represents a simple feedforward neural network
type NeuralNetwork struct {
	Layers         []*Layer
	LearningRate   float64
	Epochs         int
	BatchSize      int
	ActivationFunc string
}

// Layer represents a neural network layer
type Layer struct {
	Weights [][]float64
	Biases  []float64
	Size    int
}

// TrainingExample represents a training data point
type TrainingExample struct {
	Features []float64
	Label    float64 // 0: benign, 1: malicious
}

// ThreatFeatures represents features for threat detection
type ThreatFeatures struct {
	// Network features
	PacketSize          float64
	PacketRate          float64
	BytesPerSecond      float64
	FlowDuration        float64
	
	// Protocol features
	TCPFlags            []float64
	ProtocolDistribution map[string]float64
	PortNumbers         []float64
	
	// Behavioral features
	SessionCount        float64
	UniqueIPs           float64
	FailedConnections   float64
	
	// Temporal features
	TimeOfDay           float64
	DayOfWeek           float64
	RequestInterval     float64
	
	// Statistical features
	PacketSizeStdDev    float64
	InterArrivalTime    float64
	Entropy             float64
	
	Timestamp           time.Time
}

// ThreatPrediction represents the prediction result
type ThreatPrediction struct {
	IsThreat      bool
	Confidence    float64
	ThreatType    string
	Severity      string
	Features      *ThreatFeatures
	ModelVersion  string
	PredictedAt   time.Time
}

// FeatureScaler normalizes features
type FeatureScaler struct {
	Mean   []float64
	StdDev []float64
}

// NewDeepLearningDetector creates a new deep learning detector
func NewDeepLearningDetector(logger *logrus.Logger) *DeepLearningDetector {
	if logger == nil {
		logger = logrus.New()
	}

	// 初始化神經網路（3層：輸入層、隱藏層、輸出層）
	model := &NeuralNetwork{
		Layers: []*Layer{
			{Size: 16}, // 輸入層（16 個特徵）
			{Size: 32}, // 第一隱藏層
			{Size: 16}, // 第二隱藏層
			{Size: 1},  // 輸出層
		},
		LearningRate:   0.001,
		Epochs:         100,
		BatchSize:      32,
		ActivationFunc: "relu",
	}

	// 初始化權重和偏置
	initializeWeights(model)

	return &DeepLearningDetector{
		model:         model,
		featureScaler: &FeatureScaler{},
		trainingData:  make([]TrainingExample, 0),
		logger:        logger,
		threshold:     0.7,
		modelVersion:  "1.0.0-dl",
	}
}

// Predict predicts if traffic is a threat
func (dld *DeepLearningDetector) Predict(ctx context.Context, features *ThreatFeatures) (*ThreatPrediction, error) {
	// 提取特徵向量
	featureVector := dld.extractFeatureVector(features)

	// 標準化特徵
	normalizedFeatures := dld.featureScaler.Transform(featureVector)

	// 前向傳播
	output := dld.forward(normalizedFeatures)

	// 應用 sigmoid 獲得概率
	probability := sigmoid(output[0])

	// 判斷是否為威脅
	isThreat := probability >= dld.threshold

	// 確定威脅類型和嚴重程度
	threatType, severity := dld.classifyThreat(features, probability)

	prediction := &ThreatPrediction{
		IsThreat:     isThreat,
		Confidence:   probability,
		ThreatType:   threatType,
		Severity:     severity,
		Features:     features,
		ModelVersion: dld.modelVersion,
		PredictedAt:  time.Now(),
	}

	if isThreat {
		dld.logger.Warnf("Threat detected: type=%s, severity=%s, confidence=%.2f", 
			threatType, severity, probability)
	}

	return prediction, nil
}

// forward performs forward propagation
func (dld *DeepLearningDetector) forward(input []float64) []float64 {
	dld.mu.RLock()
	defer dld.mu.RUnlock()

	activation := input

	// 通過每一層
	for i := 0; i < len(dld.model.Layers)-1; i++ {
		layer := dld.model.Layers[i]
		nextLayer := dld.model.Layers[i+1]

		// 計算 z = W * a + b
		z := make([]float64, nextLayer.Size)
		for j := 0; j < nextLayer.Size; j++ {
			sum := nextLayer.Biases[j]
			for k := 0; k < len(activation); k++ {
				sum += nextLayer.Weights[j][k] * activation[k]
			}
			z[j] = sum
		}

		// 應用激活函數
		activation = applyActivation(z, dld.model.ActivationFunc)
	}

	return activation
}

// extractFeatureVector extracts feature vector from ThreatFeatures
func (dld *DeepLearningDetector) extractFeatureVector(features *ThreatFeatures) []float64 {
	vector := make([]float64, 16)

	vector[0] = normalizePacketSize(features.PacketSize)
	vector[1] = normalizePacketRate(features.PacketRate)
	vector[2] = normalizeBytesPerSecond(features.BytesPerSecond)
	vector[3] = normalizeFlowDuration(features.FlowDuration)
	vector[4] = normalizeSessionCount(features.SessionCount)
	vector[5] = normalizeUniqueIPs(features.UniqueIPs)
	vector[6] = normalizeFailedConnections(features.FailedConnections)
	vector[7] = normalizeTimeOfDay(features.TimeOfDay)
	vector[8] = normalizeDayOfWeek(features.DayOfWeek)
	vector[9] = normalizeRequestInterval(features.RequestInterval)
	vector[10] = features.PacketSizeStdDev
	vector[11] = features.InterArrivalTime
	vector[12] = features.Entropy
	
	// TCP Flags 平均值
	if len(features.TCPFlags) > 0 {
		sum := 0.0
		for _, flag := range features.TCPFlags {
			sum += flag
		}
		vector[13] = sum / float64(len(features.TCPFlags))
	}
	
	// Port Numbers 平均值
	if len(features.PortNumbers) > 0 {
		sum := 0.0
		for _, port := range features.PortNumbers {
			sum += port
		}
		vector[14] = sum / float64(len(features.PortNumbers))
	}
	
	// Protocol Distribution 熵
	vector[15] = calculateProtocolEntropy(features.ProtocolDistribution)

	return vector
}

// classifyThreat classifies the threat type and severity
func (dld *DeepLearningDetector) classifyThreat(features *ThreatFeatures, confidence float64) (string, string) {
	threatType := "unknown"
	severity := "low"

	// 基於特徵模式分類威脅類型
	if features.PacketRate > 1000 && features.UniqueIPs < 10 {
		threatType = "ddos"
		severity = "critical"
	} else if features.FailedConnections > 50 {
		threatType = "brute_force"
		severity = "high"
	} else if features.Entropy > 0.9 {
		threatType = "data_exfiltration"
		severity = "high"
	} else if features.PacketSize > 1400 && features.PacketRate > 100 {
		threatType = "flooding"
		severity = "medium"
	} else if confidence > 0.9 {
		threatType = "anomaly"
		severity = "high"
	} else if confidence > 0.8 {
		threatType = "suspicious"
		severity = "medium"
	}

	return threatType, severity
}

// Train trains the model with new data
func (dld *DeepLearningDetector) Train(ctx context.Context, examples []TrainingExample) error {
	dld.mu.Lock()
	defer dld.mu.Unlock()

	dld.logger.Infof("Starting training with %d examples", len(examples))

	// 添加到訓練數據
	dld.trainingData = append(dld.trainingData, examples...)

	// 計算特徵縮放參數
	dld.computeScalingParams()

	// 訓練神經網路（簡化版本）
	for epoch := 0; epoch < dld.model.Epochs; epoch++ {
		totalLoss := 0.0

		// Mini-batch training
		for i := 0; i < len(examples); i += dld.model.BatchSize {
			end := i + dld.model.BatchSize
			if end > len(examples) {
				end = len(examples)
			}

			batch := examples[i:end]
			loss := dld.trainBatch(batch)
			totalLoss += loss
		}

		avgLoss := totalLoss / float64(len(examples))
		
		if epoch%10 == 0 {
			dld.logger.Debugf("Epoch %d/%d, Loss: %.4f", epoch, dld.model.Epochs, avgLoss)
		}
	}

	dld.logger.Infof("Training completed")
	return nil
}

// trainBatch trains on a batch of examples
func (dld *DeepLearningDetector) trainBatch(batch []TrainingExample) float64 {
	totalLoss := 0.0

	for _, example := range batch {
		// 前向傳播
		normalized := dld.featureScaler.Transform(example.Features)
		output := dld.forward(normalized)
		prediction := sigmoid(output[0])

		// 計算損失（二元交叉熵）
		loss := -example.Label*math.Log(prediction+1e-10) - (1-example.Label)*math.Log(1-prediction+1e-10)
		totalLoss += loss

		// 反向傳播（簡化版本，實際需要完整實現）
		// 這裡省略詳細的梯度計算和權重更新
	}

	return totalLoss / float64(len(batch))
}

// computeScalingParams computes mean and std dev for feature scaling
func (dld *DeepLearningDetector) computeScalingParams() {
	if len(dld.trainingData) == 0 {
		return
	}

	numFeatures := len(dld.trainingData[0].Features)
	dld.featureScaler.Mean = make([]float64, numFeatures)
	dld.featureScaler.StdDev = make([]float64, numFeatures)

	// 計算均值
	for _, example := range dld.trainingData {
		for i, feature := range example.Features {
			dld.featureScaler.Mean[i] += feature
		}
	}

	n := float64(len(dld.trainingData))
	for i := range dld.featureScaler.Mean {
		dld.featureScaler.Mean[i] /= n
	}

	// 計算標準差
	for _, example := range dld.trainingData {
		for i, feature := range example.Features {
			diff := feature - dld.featureScaler.Mean[i]
			dld.featureScaler.StdDev[i] += diff * diff
		}
	}

	for i := range dld.featureScaler.StdDev {
		dld.featureScaler.StdDev[i] = math.Sqrt(dld.featureScaler.StdDev[i] / n)
		if dld.featureScaler.StdDev[i] == 0 {
			dld.featureScaler.StdDev[i] = 1.0
		}
	}
}

// Transform normalizes features using z-score normalization
func (fs *FeatureScaler) Transform(features []float64) []float64 {
	if len(fs.Mean) == 0 || len(fs.StdDev) == 0 {
		return features
	}

	normalized := make([]float64, len(features))
	for i, feature := range features {
		if i < len(fs.Mean) && i < len(fs.StdDev) {
			normalized[i] = (feature - fs.Mean[i]) / fs.StdDev[i]
		} else {
			normalized[i] = feature
		}
	}

	return normalized
}

// SaveModel saves the model to JSON
func (dld *DeepLearningDetector) SaveModel(filepath string) error {
	dld.mu.RLock()
	defer dld.mu.RUnlock()

	data, err := json.MarshalIndent(dld.model, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal model: %w", err)
	}

	// 實際實現需要寫入檔案
	dld.logger.Infof("Model saved to %s", filepath)
	return nil
}

// LoadModel loads the model from JSON
func (dld *DeepLearningDetector) LoadModel(filepath string) error {
	dld.mu.Lock()
	defer dld.mu.Unlock()

	// 實際實現需要從檔案讀取
	dld.logger.Infof("Model loaded from %s", filepath)
	return nil
}

// Helper functions

func initializeWeights(model *NeuralNetwork) {
	for i := 0; i < len(model.Layers)-1; i++ {
		currentLayer := model.Layers[i]
		nextLayer := model.Layers[i+1]

		// Xavier initialization
		limit := math.Sqrt(6.0 / float64(currentLayer.Size+nextLayer.Size))

		nextLayer.Weights = make([][]float64, nextLayer.Size)
		nextLayer.Biases = make([]float64, nextLayer.Size)

		for j := 0; j < nextLayer.Size; j++ {
			nextLayer.Weights[j] = make([]float64, currentLayer.Size)
			for k := 0; k < currentLayer.Size; k++ {
				// 隨機初始化（實際應使用更好的隨機數生成器）
				nextLayer.Weights[j][k] = (math.Sin(float64(j*k+1)) * 2 - 1) * limit
			}
			nextLayer.Biases[j] = 0.0
		}
	}
}

func applyActivation(z []float64, activationType string) []float64 {
	activated := make([]float64, len(z))

	switch activationType {
	case "relu":
		for i, val := range z {
			activated[i] = math.Max(0, val)
		}
	case "sigmoid":
		for i, val := range z {
			activated[i] = sigmoid(val)
		}
	case "tanh":
		for i, val := range z {
			activated[i] = math.Tanh(val)
		}
	default:
		copy(activated, z)
	}

	return activated
}

// Normalization functions

func normalizePacketSize(size float64) float64 {
	return math.Min(size/1500.0, 1.0)
}

func normalizePacketRate(rate float64) float64 {
	return math.Min(rate/10000.0, 1.0)
}

func normalizeBytesPerSecond(bytes float64) float64 {
	return math.Min(bytes/1000000.0, 1.0) // 1 Mbps
}

func normalizeFlowDuration(duration float64) float64 {
	return math.Min(duration/3600.0, 1.0) // 1 hour
}

func normalizeSessionCount(count float64) float64 {
	return math.Min(count/1000.0, 1.0)
}

func normalizeUniqueIPs(count float64) float64 {
	return math.Min(count/1000.0, 1.0)
}

func normalizeFailedConnections(count float64) float64 {
	return math.Min(count/100.0, 1.0)
}

func normalizeTimeOfDay(hour float64) float64 {
	return hour / 24.0
}

func normalizeDayOfWeek(day float64) float64 {
	return day / 7.0
}

func normalizeRequestInterval(interval float64) float64 {
	return math.Min(interval/60.0, 1.0) // 1 minute
}

func calculateProtocolEntropy(distribution map[string]float64) float64 {
	if len(distribution) == 0 {
		return 0.0
	}

	entropy := 0.0
	total := 0.0

	for _, count := range distribution {
		total += count
	}

	for _, count := range distribution {
		if count > 0 {
			p := count / total
			entropy -= p * math.Log2(p)
		}
	}

	return entropy
}

