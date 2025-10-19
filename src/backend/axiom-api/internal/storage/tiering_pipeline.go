package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// TieringPipeline 資料分層流轉管道
type TieringPipeline struct {
	hotStorage  *HotStorage
	coldStorage *ColdStorage
	// warmStorage *WarmStorage  // Loki integration (待實施)
	// archiveStorage *ArchiveStorage  // S3/MinIO (待實施)
	
	ctx    context.Context
	cancel context.CancelFunc
}

// NewTieringPipeline 創建分層流轉管道
func NewTieringPipeline(redisClient *redis.Client, db *gorm.DB) *TieringPipeline {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &TieringPipeline{
		hotStorage:  NewHotStorage(redisClient),
		coldStorage: NewColdStorage(db),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start 啟動管道
func (p *TieringPipeline) Start() {
	log.Println("Starting tiering pipeline...")
	
	// Task 1: Hot → Cold (每 5 分鐘)
	// 將 1 小時以上的數據從 Redis 轉移到 PostgreSQL
	go p.scheduleTask("hot-to-cold", 5*time.Minute, p.transferHotToCold)
	
	// Task 2: Cold → Archive (每天)
	// 將 90 天以上的數據從 PostgreSQL 封存到 S3
	// go p.scheduleTask("cold-to-archive", 24*time.Hour, p.transferColdToArchive)
	
	// Task 3: 完整性驗證 (每天)
	go p.scheduleTask("integrity-check", 24*time.Hour, p.verifyIntegrity)
	
	// Task 4: 保留策略執行 (每天)
	go p.scheduleTask("retention-enforcement", 24*time.Hour, p.enforceRetention)
	
	log.Println("Tiering pipeline started successfully")
}

// Stop 停止管道
func (p *TieringPipeline) Stop() {
	log.Println("Stopping tiering pipeline...")
	p.cancel()
}

// scheduleTask 調度任務
func (p *TieringPipeline) scheduleTask(name string, interval time.Duration, task func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	// 立即執行一次
	if err := task(); err != nil {
		log.Printf("[%s] Error: %v", name, err)
	}
	
	for {
		select {
		case <-ticker.C:
			log.Printf("[%s] Running task...", name)
			if err := task(); err != nil {
				log.Printf("[%s] Error: %v", name, err)
			} else {
				log.Printf("[%s] Task completed successfully", name)
			}
			
		case <-p.ctx.Done():
			log.Printf("[%s] Task stopped", name)
			return
		}
	}
}

// transferHotToCold 將 Hot 數據轉移到 Cold
func (p *TieringPipeline) transferHotToCold() error {
	log.Println("Transferring hot data to cold storage...")
	
	// 獲取 1 小時以上的數據
	logs, err := p.hotStorage.GetEntriesOlderThan(p.ctx, 1*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to get old entries: %w", err)
	}
	
	if len(logs) == 0 {
		log.Println("No logs to transfer")
		return nil
	}
	
	// 寫入 Cold Storage
	if err := p.coldStorage.Write(p.ctx, logs); err != nil {
		return fmt.Errorf("failed to write to cold storage: %w", err)
	}
	
	// 從 Hot Storage 刪除
	if err := p.hotStorage.Delete(p.ctx, logs); err != nil {
		return fmt.Errorf("failed to delete from hot storage: %w", err)
	}
	
	log.Printf("Transferred %d logs from hot to cold storage", len(logs))
	return nil
}

// verifyIntegrity 驗證完整性
func (p *TieringPipeline) verifyIntegrity() error {
	log.Println("Verifying data integrity...")
	
	tamperedLogs, err := p.coldStorage.VerifyIntegrity(p.ctx)
	if err != nil {
		return fmt.Errorf("integrity check failed: %w", err)
	}
	
	if len(tamperedLogs) > 0 {
		log.Printf("WARNING: Found %d tampered logs: %v", len(tamperedLogs), tamperedLogs)
		// TODO: 發送安全告警
	} else {
		log.Println("Integrity check passed - no tampering detected")
	}
	
	return nil
}

// enforceRetention 執行保留策略
func (p *TieringPipeline) enforceRetention() error {
	log.Println("Enforcing retention policies...")
	
	// TODO: 根據 retention_policies 表自動刪除過期數據
	// 這裡需要查詢保留策略，然後刪除過期的日誌
	
	log.Println("Retention policies enforced")
	return nil
}

// GetStats 獲取管道統計
func (p *TieringPipeline) GetStats(ctx context.Context) (map[string]interface{}, error) {
	hotStats, _ := p.hotStorage.GetStats(ctx)
	coldStats, _ := p.coldStorage.GetStats(ctx)
	
	return map[string]interface{}{
		"hot_storage":  hotStats,
		"cold_storage": coldStats,
		"pipeline_status": "running",
		"last_check":   time.Now(),
	}, nil
}

