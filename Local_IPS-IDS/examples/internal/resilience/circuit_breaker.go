package resilience

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreakerConfig contains circuit breaker configuration
type CircuitBreakerConfig struct {
	MaxFailures     int
	Timeout         time.Duration
	ResetTimeout    time.Duration
	HalfOpenMaxCalls int
}

// DefaultCircuitBreakerConfig returns default circuit breaker configuration
func DefaultCircuitBreakerConfig() *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		MaxFailures:      5,
		Timeout:          60 * time.Second,
		ResetTimeout:     30 * time.Second,
		HalfOpenMaxCalls: 1,
	}
}

// CircuitBreaker implements the circuit breaker pattern
// 斷路器模式實現
type CircuitBreaker struct {
	config       *CircuitBreakerConfig
	state        CircuitState
	failures     int
	lastFailTime time.Time
	halfOpenCalls int
	logger       *logrus.Logger
	mu           sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config *CircuitBreakerConfig, logger *logrus.Logger) *CircuitBreaker {
	if config == nil {
		config = DefaultCircuitBreakerConfig()
	}

	if logger == nil {
		logger = logrus.New()
	}

	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
		logger: logger,
	}
}

// Execute executes a function through the circuit breaker
// 通過斷路器執行函數
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	// 檢查斷路器狀態
	if err := cb.beforeRequest(); err != nil {
		return err
	}

	// 執行函數
	err := fn()

	// 記錄結果
	cb.afterRequest(err)

	return err
}

// ExecuteWithResult executes a function with result through the circuit breaker
// 帶返回值的斷路器執行
func ExecuteWithResult[T any](cb *CircuitBreaker, ctx context.Context, fn func() (T, error)) (T, error) {
	var result T

	if err := cb.beforeRequest(); err != nil {
		return result, err
	}

	result, err := fn()
	cb.afterRequest(err)

	return result, err
}

// beforeRequest checks if request should be allowed
func (cb *CircuitBreaker) beforeRequest() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		// 正常狀態，允許請求
		return nil

	case StateOpen:
		// 斷路器打開，檢查是否可以進入半開狀態
		if time.Since(cb.lastFailTime) > cb.config.Timeout {
			cb.state = StateHalfOpen
			cb.halfOpenCalls = 0
			cb.logger.Info("Circuit breaker entering half-open state")
			return nil
		}
		return fmt.Errorf("circuit breaker is open")

	case StateHalfOpen:
		// 半開狀態，限制請求數量
		if cb.halfOpenCalls >= cb.config.HalfOpenMaxCalls {
			return fmt.Errorf("circuit breaker is half-open, max calls reached")
		}
		cb.halfOpenCalls++
		return nil

	default:
		return fmt.Errorf("unknown circuit breaker state")
	}
}

// afterRequest records the result of a request
func (cb *CircuitBreaker) afterRequest(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		// 請求失敗
		cb.failures++
		cb.lastFailTime = time.Now()

		if cb.state == StateHalfOpen {
			// 半開狀態失敗，重新打開斷路器
			cb.state = StateOpen
			cb.logger.Warn("Circuit breaker reopened after failure in half-open state")
		} else if cb.failures >= cb.config.MaxFailures {
			// 失敗次數達到閾值，打開斷路器
			cb.state = StateOpen
			cb.logger.Warnf("Circuit breaker opened after %d failures", cb.failures)
		}
	} else {
		// 請求成功
		if cb.state == StateHalfOpen {
			// 半開狀態成功，關閉斷路器
			cb.state = StateClosed
			cb.failures = 0
			cb.halfOpenCalls = 0
			cb.logger.Info("Circuit breaker closed after successful call in half-open state")
		} else if cb.state == StateClosed {
			// 正常狀態成功，重置失敗計數
			cb.failures = 0
		}
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetFailures returns the current number of failures
func (cb *CircuitBreaker) GetFailures() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = StateClosed
	cb.failures = 0
	cb.halfOpenCalls = 0
	cb.logger.Info("Circuit breaker manually reset")
}

// String returns the string representation of circuit state
func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

