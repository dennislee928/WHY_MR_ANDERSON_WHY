package waitingroom

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// WaitingRoom implements a virtual waiting room for traffic peak handling
// 虛擬等待室實現，用於處理流量峰值
type WaitingRoom struct {
	redis      *redis.Client
	queueKey   string
	maxActive  int64
	timeout    time.Duration
	logger     *logrus.Logger
	mu         sync.RWMutex
	enabled    bool
}

// Config contains waiting room configuration
type Config struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	QueueKey      string
	MaxActive     int64
	Timeout       time.Duration
	AutoEnable    bool
	Threshold     int64 // 自動啟用閾值
}

// NewWaitingRoom creates a new waiting room
func NewWaitingRoom(config *Config, logger *logrus.Logger) (*WaitingRoom, error) {
	if logger == nil {
		logger = logrus.New()
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// 測試連接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &WaitingRoom{
		redis:     redisClient,
		queueKey:  config.QueueKey,
		maxActive: config.MaxActive,
		timeout:   config.Timeout,
		logger:    logger,
		enabled:   config.AutoEnable,
	}, nil
}

// Enqueue adds a user to the waiting room queue
func (wr *WaitingRoom) Enqueue(ctx context.Context, userID string) (*QueuePosition, error) {
	wr.mu.RLock()
	if !wr.enabled {
		wr.mu.RUnlock()
		return &QueuePosition{
			UserID:   userID,
			Position: 0,
			Admitted: true,
		}, nil
	}
	wr.mu.RUnlock()

	// 添加到 Redis sorted set（使用時間戳作為分數）
	score := float64(time.Now().UnixNano())
	if err := wr.redis.ZAdd(ctx, wr.queueKey, redis.Z{
		Score:  score,
		Member: userID,
	}).Err(); err != nil {
		return nil, fmt.Errorf("failed to enqueue: %w", err)
	}

	// 獲取排隊位置
	rank, err := wr.redis.ZRank(ctx, wr.queueKey, userID).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get rank: %w", err)
	}

	// 檢查是否可以進入
	admitted := rank < wr.maxActive

	position := &QueuePosition{
		UserID:         userID,
		Position:       rank + 1,
		QueueLength:    wr.getQueueLength(ctx),
		EstimatedWait:  wr.estimateWaitTime(rank),
		Admitted:       admitted,
		EnqueuedAt:     time.Now(),
	}

	wr.logger.Infof("User %s enqueued at position %d (admitted: %v)", userID, position.Position, admitted)

	return position, nil
}

// CheckPosition checks a user's position in the queue
func (wr *WaitingRoom) CheckPosition(ctx context.Context, userID string) (*QueuePosition, error) {
	// 獲取排隊位置
	rank, err := wr.redis.ZRank(ctx, wr.queueKey, userID).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("user not in queue")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get rank: %w", err)
	}

	admitted := rank < wr.maxActive

	return &QueuePosition{
		UserID:        userID,
		Position:      rank + 1,
		QueueLength:   wr.getQueueLength(ctx),
		EstimatedWait: wr.estimateWaitTime(rank),
		Admitted:      admitted,
	}, nil
}

// Dequeue removes a user from the queue
func (wr *WaitingRoom) Dequeue(ctx context.Context, userID string) error {
	if err := wr.redis.ZRem(ctx, wr.queueKey, userID).Err(); err != nil {
		return fmt.Errorf("failed to dequeue: %w", err)
	}

	wr.logger.Debugf("User %s dequeued", userID)
	return nil
}

// Enable enables the waiting room
func (wr *WaitingRoom) Enable() {
	wr.mu.Lock()
	defer wr.mu.Unlock()

	wr.enabled = true
	wr.logger.Info("Waiting room enabled")
}

// Disable disables the waiting room
func (wr *WaitingRoom) Disable() {
	wr.mu.Lock()
	defer wr.mu.Unlock()

	wr.enabled = false
	wr.logger.Info("Waiting room disabled")
}

// IsEnabled checks if the waiting room is enabled
func (wr *WaitingRoom) IsEnabled() bool {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	return wr.enabled
}

// getQueueLength returns the total queue length
func (wr *WaitingRoom) getQueueLength(ctx context.Context) int64 {
	count, err := wr.redis.ZCard(ctx, wr.queueKey).Result()
	if err != nil {
		wr.logger.Errorf("Failed to get queue length: %v", err)
		return 0
	}
	return count
}

// estimateWaitTime estimates wait time based on position
func (wr *WaitingRoom) estimateWaitTime(rank int64) time.Duration {
	if rank < wr.maxActive {
		return 0
	}

	// 假設平均每個用戶停留 30 秒
	avgSessionTime := 30 * time.Second
	usersAhead := rank - wr.maxActive

	return time.Duration(usersAhead) * avgSessionTime / time.Duration(wr.maxActive)
}

// CleanupExpired removes expired entries from the queue
func (wr *WaitingRoom) CleanupExpired(ctx context.Context) error {
	// 移除超過 timeout 的條目
	minScore := float64(time.Now().Add(-wr.timeout).UnixNano())

	removed, err := wr.redis.ZRemRangeByScore(ctx, wr.queueKey, "-inf", fmt.Sprintf("%f", minScore)).Result()
	if err != nil {
		return fmt.Errorf("failed to cleanup: %w", err)
	}

	if removed > 0 {
		wr.logger.Infof("Cleaned up %d expired queue entries", removed)
	}

	return nil
}

// StartCleanupWorker starts a background worker to cleanup expired entries
func (wr *WaitingRoom) StartCleanupWorker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := wr.CleanupExpired(ctx); err != nil {
				wr.logger.Errorf("Cleanup error: %v", err)
			}
		}
	}
}

// QueuePosition represents a user's position in the queue
type QueuePosition struct {
	UserID        string
	Position      int64
	QueueLength   int64
	EstimatedWait time.Duration
	Admitted      bool
	EnqueuedAt    time.Time
}

// Close closes the waiting room
func (wr *WaitingRoom) Close() error {
	if wr.redis != nil {
		return wr.redis.Close()
	}
	return nil
}

