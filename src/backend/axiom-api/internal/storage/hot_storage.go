package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// HotStorage Hot Storage (Redis Streams) - 1小時實時日誌
type HotStorage struct {
	redis *redis.Client
}

// LogEntry 日誌條目
type LogEntry struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	AgentID   string                 `json:"agent_id"`
	AgentMode string                 `json:"agent_mode"`
	EventType string                 `json:"event_type"`
	Source    string                 `json:"source"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
}

// NewHotStorage 創建 Hot Storage
func NewHotStorage(redisClient *redis.Client) *HotStorage {
	return &HotStorage{
		redis: redisClient,
	}
}

// Write 寫入日誌到 Redis Streams
func (h *HotStorage) Write(ctx context.Context, agentID string, logs []LogEntry) error {
	streamKey := h.getStreamKey(agentID, time.Now())
	
	pipe := h.redis.Pipeline()
	
	for _, log := range logs {
		// 序列化日誌數據
		logData, err := json.Marshal(log)
		if err != nil {
			return fmt.Errorf("failed to marshal log: %w", err)
		}
		
		// 添加到 Redis Stream
		pipe.XAdd(ctx, &redis.XAddArgs{
			Stream: streamKey,
			MaxLen: 100000, // 每個 stream 最多保留 10萬條
			Values: map[string]interface{}{
				"data": string(logData),
			},
		})
	}
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to write to Redis Stream: %w", err)
	}
	
	// 設置過期時間 (1小時)
	h.redis.Expire(ctx, streamKey, 1*time.Hour)
	
	return nil
}

// Query 查詢日誌
func (h *HotStorage) Query(ctx context.Context, agentID string, start, end time.Time) ([]LogEntry, error) {
	streamKey := h.getStreamKey(agentID, start)
	
	// 從 Redis Stream 讀取
	messages, err := h.redis.XRange(ctx, streamKey, "-", "+").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to read from Redis Stream: %w", err)
	}
	
	var logs []LogEntry
	for _, msg := range messages {
		logDataStr, ok := msg.Values["data"].(string)
		if !ok {
			continue
		}
		
		var log LogEntry
		if err := json.Unmarshal([]byte(logDataStr), &log); err != nil {
			continue
		}
		
		// 時間範圍過濾
		if log.Timestamp.After(start) && log.Timestamp.Before(end) {
			logs = append(logs, log)
		}
	}
	
	return logs, nil
}

// GetEntriesOlderThan 獲取舊日誌（用於轉移）
func (h *HotStorage) GetEntriesOlderThan(ctx context.Context, duration time.Duration) ([]LogEntry, error) {
	cutoffTime := time.Now().Add(-duration)
	
	// 掃描所有 streams
	var logs []LogEntry
	
	// 使用 SCAN 遍歷所有 logs: 開頭的 key
	iter := h.redis.Scan(ctx, 0, "logs:agent:*", 0).Iterator()
	for iter.Next(ctx) {
		streamKey := iter.Val()
		
		messages, err := h.redis.XRange(ctx, streamKey, "-", "+").Result()
		if err != nil {
			continue
		}
		
		for _, msg := range messages {
			logDataStr, ok := msg.Values["data"].(string)
			if !ok {
				continue
			}
			
			var log LogEntry
			if err := json.Unmarshal([]byte(logDataStr), &log); err != nil {
				continue
			}
			
			if log.Timestamp.Before(cutoffTime) {
				logs = append(logs, log)
			}
		}
	}
	
	return logs, nil
}

// Delete 刪除日誌
func (h *HotStorage) Delete(ctx context.Context, entries []LogEntry) error {
	// 按 stream 分組
	streamGroups := make(map[string][]string)
	
	for _, entry := range entries {
		streamKey := h.getStreamKey(entry.AgentID, entry.Timestamp)
		streamGroups[streamKey] = append(streamGroups[streamKey], entry.ID)
	}
	
	// 刪除每個 stream 中的條目
	pipe := h.redis.Pipeline()
	for streamKey, ids := range streamGroups {
		for _, id := range ids {
			pipe.XDel(ctx, streamKey, id)
		}
	}
	
	_, err := pipe.Exec(ctx)
	return err
}

// GetStats 獲取統計信息
func (h *HotStorage) GetStats(ctx context.Context) (map[string]interface{}, error) {
	// 統計所有 streams
	var totalEntries int64
	var streamCount int
	
	iter := h.redis.Scan(ctx, 0, "logs:agent:*", 0).Iterator()
	for iter.Next(ctx) {
		streamKey := iter.Val()
		streamCount++
		
		length, err := h.redis.XLen(ctx, streamKey).Result()
		if err != nil {
			continue
		}
		
		totalEntries += length
	}
	
	return map[string]interface{}{
		"total_entries":  totalEntries,
		"stream_count":   streamCount,
		"retention":      "1 hour",
		"max_per_stream": 100000,
	}, nil
}

// CreateConsumerGroup 創建消費者組
func (h *HotStorage) CreateConsumerGroup(ctx context.Context, agentID, groupName string) error {
	streamKey := h.getStreamKey(agentID, time.Now())
	
	err := h.redis.XGroupCreateMkStream(ctx, streamKey, groupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	
	return nil
}

// ReadFromGroup 從消費者組讀取
func (h *HotStorage) ReadFromGroup(ctx context.Context, agentID, groupName, consumerName string, count int64) ([]LogEntry, error) {
	streamKey := h.getStreamKey(agentID, time.Now())
	
	streams, err := h.redis.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{streamKey, ">"},
		Count:    count,
		Block:    0,
	}).Result()
	
	if err != nil {
		return nil, fmt.Errorf("failed to read from group: %w", err)
	}
	
	var logs []LogEntry
	for _, stream := range streams {
		for _, msg := range stream.Messages {
			logDataStr, ok := msg.Values["data"].(string)
			if !ok {
				continue
			}
			
			var log LogEntry
			if err := json.Unmarshal([]byte(logDataStr), &log); err != nil {
				continue
			}
			
			log.ID = msg.ID
			logs = append(logs, log)
		}
	}
	
	return logs, nil
}

// getStreamKey 生成 stream key
func (h *HotStorage) getStreamKey(agentID string, timestamp time.Time) string {
	return fmt.Sprintf("logs:agent:%s:%s", agentID, timestamp.Format("2006-01-02"))
}

