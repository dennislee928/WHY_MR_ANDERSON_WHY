package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"axiom-backend/internal/model"
)

// Database 資料庫管理器
type Database struct {
	PG    *gorm.DB
	Redis *redis.Client
}

// Config 資料庫配置
type Config struct {
	// PostgreSQL 配置
	PGHost     string
	PGPort     int
	PGUser     string
	PGPassword string
	PGDatabase string
	PGSSL      string

	// Redis 配置
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int
}

// NewDatabase 創建新的資料庫管理器
func NewDatabase(cfg *Config) (*Database, error) {
	// 初始化 PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.PGHost, cfg.PGPort, cfg.PGUser, cfg.PGPassword, cfg.PGDatabase, cfg.PGSSL,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// 配置連接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 初始化 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
		PoolSize: 20,
	})

	// 測試 Redis 連接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Database{
		PG:    db,
		Redis: rdb,
	}, nil
}

// AutoMigrate 自動遷移資料庫表
func (d *Database) AutoMigrate() error {
	// 遷移基礎 Models
	err := d.PG.AutoMigrate(
		&model.Service{},
		&model.ConfigHistory{},
		&model.QuantumJob{},
		&model.WindowsLog{},
		&model.Alert{},
		&model.APILog{},
		&model.MetricSnapshot{},
		&model.User{},
		&model.Session{},
	)
	if err != nil {
		return err
	}
	
	// 遷移合規性 Models (Phase 13)
	err = model.AutoMigrateComplianceTables(d.PG)
	if err != nil {
		return err
	}
	
	// 初始化默認策略
	if err := model.SeedDefaultPolicies(d.PG); err != nil {
		return err
	}
	
	if err := model.SeedDefaultPIIPatterns(d.PG); err != nil {
		return err
	}
	
	return nil
}

// Close 關閉資料庫連接
func (d *Database) Close() error {
	// 關閉 PostgreSQL
	sqlDB, err := d.PG.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}

	// 關閉 Redis
	if err := d.Redis.Close(); err != nil {
		return err
	}

	return nil
}

// HealthCheck 健康檢查
func (d *Database) HealthCheck(ctx context.Context) error {
	// 檢查 PostgreSQL
	sqlDB, err := d.PG.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("PostgreSQL health check failed: %w", err)
	}

	// 檢查 Redis
	if err := d.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	return nil
}

