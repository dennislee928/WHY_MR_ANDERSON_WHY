package loadbalancer

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

// LoadBalancer 負載均衡器
type LoadBalancer struct {
	config   *Config
	backends []*Backend
	current  uint64 // Round-robin 計數器
	logger   *logrus.Logger
	mu       sync.RWMutex
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// Backend 後端服務器
type Backend struct {
	URL       string
	Healthy   bool
	LastCheck time.Time
	Failures  int
	mu        sync.RWMutex
}

// Config 負載均衡器配置
type Config struct {
	Backends            []string      `yaml:"backends" json:"backends"`                           // 後端服務器列表
	Strategy            string        `yaml:"strategy" json:"strategy"`                           // 負載均衡策略: "round-robin", "least-conn", "random"
	HealthCheckEnabled  bool          `yaml:"health_check_enabled" json:"health_check_enabled"`   // 是否啟用健康檢查
	HealthCheckInterval time.Duration `yaml:"health_check_interval" json:"health_check_interval"` // 健康檢查間隔
	HealthCheckTimeout  time.Duration `yaml:"health_check_timeout" json:"health_check_timeout"`   // 健康檢查超時
	HealthCheckPath     string        `yaml:"health_check_path" json:"health_check_path"`         // 健康檢查路徑
	MaxRetries          int           `yaml:"max_retries" json:"max_retries"`                     // 最大重試次數
	RetryDelay          time.Duration `yaml:"retry_delay" json:"retry_delay"`                     // 重試延遲
}

// NewLoadBalancer 創建新的負載均衡器
func NewLoadBalancer(config *Config, logger *logrus.Logger) (*LoadBalancer, error) {
	if logger == nil {
		logger = logrus.New()
	}

	if len(config.Backends) == 0 {
		return nil, fmt.Errorf("沒有配置後端服務器")
	}

	// 設定預設值
	if config.Strategy == "" {
		config.Strategy = "round-robin"
	}
	if config.HealthCheckInterval == 0 {
		config.HealthCheckInterval = 30 * time.Second
	}
	if config.HealthCheckTimeout == 0 {
		config.HealthCheckTimeout = 5 * time.Second
	}
	if config.HealthCheckPath == "" {
		config.HealthCheckPath = "/health"
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 1 * time.Second
	}

	// 初始化後端服務器
	backends := make([]*Backend, len(config.Backends))
	for i, url := range config.Backends {
		backends[i] = &Backend{
			URL:     url,
			Healthy: true, // 初始假設健康
		}
	}

	lb := &LoadBalancer{
		config:   config,
		backends: backends,
		logger:   logger,
		stopCh:   make(chan struct{}),
	}

	// 啟動健康檢查
	if config.HealthCheckEnabled {
		lb.wg.Add(1)
		go lb.healthCheckRoutine()
	}

	return lb, nil
}

// GetBackend 獲取一個後端服務器
func (lb *LoadBalancer) GetBackend() (*Backend, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	// 過濾健康的後端
	healthyBackends := make([]*Backend, 0)
	for _, backend := range lb.backends {
		backend.mu.RLock()
		if backend.Healthy {
			healthyBackends = append(healthyBackends, backend)
		}
		backend.mu.RUnlock()
	}

	if len(healthyBackends) == 0 {
		return nil, fmt.Errorf("沒有可用的後端服務器")
	}

	// 根據策略選擇後端
	switch lb.config.Strategy {
	case "round-robin":
		return lb.roundRobin(healthyBackends), nil
	case "random":
		return lb.random(healthyBackends), nil
	case "least-conn":
		return lb.leastConnection(healthyBackends), nil
	default:
		return lb.roundRobin(healthyBackends), nil
	}
}

// roundRobin Round-robin 負載均衡
func (lb *LoadBalancer) roundRobin(backends []*Backend) *Backend {
	n := len(backends)
	if n == 0 {
		return nil
	}

	idx := atomic.AddUint64(&lb.current, 1) % uint64(n)
	return backends[idx]
}

// random 隨機負載均衡
func (lb *LoadBalancer) random(backends []*Backend) *Backend {
	if len(backends) == 0 {
		return nil
	}

	// 簡化版：使用當前時間作為隨機種子
	idx := time.Now().UnixNano() % int64(len(backends))
	return backends[idx]
}

// leastConnection 最少連接負載均衡
func (lb *LoadBalancer) leastConnection(backends []*Backend) *Backend {
	if len(backends) == 0 {
		return nil
	}

	// 簡化版：返回第一個健康的後端
	// 實際實作需要追蹤每個後端的連接數
	return backends[0]
}

// healthCheckRoutine 健康檢查協程
func (lb *LoadBalancer) healthCheckRoutine() {
	defer lb.wg.Done()

	ticker := time.NewTicker(lb.config.HealthCheckInterval)
	defer ticker.Stop()

	// 立即執行一次健康檢查
	lb.performHealthChecks()

	for {
		select {
		case <-ticker.C:
			lb.performHealthChecks()
		case <-lb.stopCh:
			return
		}
	}
}

// performHealthChecks 執行健康檢查
func (lb *LoadBalancer) performHealthChecks() {
	lb.mu.RLock()
	backends := lb.backends
	lb.mu.RUnlock()

	for _, backend := range backends {
		go lb.checkBackend(backend)
	}
}

// checkBackend 檢查單個後端
func (lb *LoadBalancer) checkBackend(backend *Backend) {
	client := &http.Client{
		Timeout: lb.config.HealthCheckTimeout,
	}

	url := backend.URL + lb.config.HealthCheckPath
	resp, err := client.Get(url)

	backend.mu.Lock()
	defer backend.mu.Unlock()

	backend.LastCheck = time.Now()

	if err != nil || resp.StatusCode != http.StatusOK {
		backend.Failures++
		if backend.Healthy && backend.Failures >= lb.config.MaxRetries {
			backend.Healthy = false
			lb.logger.Errorf("後端服務器不健康: %s (錯誤: %v)", backend.URL, err)
		}
	} else {
		if !backend.Healthy {
			lb.logger.Infof("後端服務器恢復健康: %s", backend.URL)
		}
		backend.Healthy = true
		backend.Failures = 0
	}

	if resp != nil {
		resp.Body.Close()
	}
}

// GetHealthyBackends 獲取所有健康的後端
func (lb *LoadBalancer) GetHealthyBackends() []*Backend {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	healthy := make([]*Backend, 0)
	for _, backend := range lb.backends {
		backend.mu.RLock()
		if backend.Healthy {
			healthy = append(healthy, backend)
		}
		backend.mu.RUnlock()
	}

	return healthy
}

// GetStatus 獲取負載均衡器狀態
func (lb *LoadBalancer) GetStatus() map[string]interface{} {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	backendStatus := make([]map[string]interface{}, len(lb.backends))
	healthyCount := 0

	for i, backend := range lb.backends {
		backend.mu.RLock()
		backendStatus[i] = map[string]interface{}{
			"url":        backend.URL,
			"healthy":    backend.Healthy,
			"last_check": backend.LastCheck,
			"failures":   backend.Failures,
		}
		if backend.Healthy {
			healthyCount++
		}
		backend.mu.RUnlock()
	}

	return map[string]interface{}{
		"strategy":  lb.config.Strategy,
		"total":     len(lb.backends),
		"healthy":   healthyCount,
		"unhealthy": len(lb.backends) - healthyCount,
		"backends":  backendStatus,
	}
}

// Stop 停止負載均衡器
func (lb *LoadBalancer) Stop() {
	lb.logger.Info("停止負載均衡器...")
	close(lb.stopCh)
	lb.wg.Wait()
	lb.logger.Info("負載均衡器已停止")
}
