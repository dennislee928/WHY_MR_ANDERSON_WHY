package resilience

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// RetryConfig contains retry configuration
type RetryConfig struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	MaxJitter       time.Duration
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 1 * time.Second,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		MaxJitter:       1 * time.Second,
	}
}

// Retry executes a function with exponential backoff retry
// 使用指數退避重試執行函數
func Retry(ctx context.Context, config *RetryConfig, fn func() error, logger *logrus.Logger) error {
	if config == nil {
		config = DefaultRetryConfig()
	}

	var lastErr error
	interval := config.InitialInterval

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		// 執行函數
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// 最後一次嘗試，不再重試
		if attempt == config.MaxAttempts {
			break
		}

		// 檢查 context 是否已取消
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		default:
		}

		// 記錄重試
		if logger != nil {
			logger.Warnf("Attempt %d/%d failed: %v, retrying in %v",
				attempt, config.MaxAttempts, err, interval)
		}

		// 等待後重試
		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		}

		// 計算下次重試間隔（指數退避）
		interval = time.Duration(float64(interval) * config.Multiplier)
		if interval > config.MaxInterval {
			interval = config.MaxInterval
		}

		// 添加 jitter 避免雷鳴群效應
		if config.MaxJitter > 0 {
			jitter := time.Duration(float64(config.MaxJitter) * (0.5 + 0.5*time.Now().UnixNano()%1000/1000.0))
			interval += jitter
		}
	}

	return fmt.Errorf("max retry attempts reached: %w", lastErr)
}

// RetryWithResult executes a function with retry and returns result
// 帶返回值的重試函數
func RetryWithResult[T any](ctx context.Context, config *RetryConfig, fn func() (T, error), logger *logrus.Logger) (T, error) {
	var result T
	var lastErr error
	interval := config.InitialInterval

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		var err error
		result, err = fn()
		if err == nil {
			return result, nil
		}

		lastErr = err

		if attempt == config.MaxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return result, fmt.Errorf("retry cancelled: %w", ctx.Err())
		default:
		}

		if logger != nil {
			logger.Warnf("Attempt %d/%d failed: %v, retrying in %v",
				attempt, config.MaxAttempts, err, interval)
		}

		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return result, fmt.Errorf("retry cancelled: %w", ctx.Err())
		}

		interval = time.Duration(float64(interval) * config.Multiplier)
		if interval > config.MaxInterval {
			interval = config.MaxInterval
		}
	}

	return result, fmt.Errorf("max retry attempts reached: %w", lastErr)
}

