package multitenant

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// TenantManager manages multi-tenant operations
type TenantManager struct {
	tenants    map[string]*Tenant
	mu         sync.RWMutex
	logger     *logrus.Logger
	isolation  IsolationLevel
}

// Tenant represents a tenant in the system
type Tenant struct {
	ID          string
	Name        string
	Domain      string
	Status      TenantStatus
	Plan        SubscriptionPlan
	CreatedAt   time.Time
	UpdatedAt   time.Time
	
	// Resource limits
	Limits      *ResourceLimits
	
	// Current usage
	Usage       *ResourceUsage
	
	// Configuration
	Config      *TenantConfig
	
	// Metadata
	Metadata    map[string]string
}

// TenantStatus represents tenant status
type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusDeleted   TenantStatus = "deleted"
)

// SubscriptionPlan represents subscription plans
type SubscriptionPlan string

const (
	PlanFree       SubscriptionPlan = "free"
	PlanBasic      SubscriptionPlan = "basic"
	PlanProfessional SubscriptionPlan = "professional"
	PlanEnterprise SubscriptionPlan = "enterprise"
)

// IsolationLevel represents data isolation level
type IsolationLevel string

const (
	IsolationSharedDB      IsolationLevel = "shared_db"       // 共享資料庫，不同 schema
	IsolationSeparateDB    IsolationLevel = "separate_db"     // 獨立資料庫
	IsolationDedicatedHost IsolationLevel = "dedicated_host"  // 專用主機
)

// ResourceLimits defines resource limits for a tenant
type ResourceLimits struct {
	MaxUsers           int
	MaxDevices         int
	MaxEvents          int64
	MaxStorageGB       int
	MaxAPICallsPerDay  int64
	MaxConcurrentConns int
}

// ResourceUsage tracks current resource usage
type ResourceUsage struct {
	CurrentUsers       int
	CurrentDevices     int
	EventsToday        int64
	StorageUsedGB      float64
	APICallsToday      int64
	ConcurrentConns    int
	LastUpdated        time.Time
}

// TenantConfig contains tenant-specific configuration
type TenantConfig struct {
	EnableAdvancedFeatures bool
	EnableMLDetection      bool
	EnableAutoResponse     bool
	RetentionDays          int
	CustomDomain           string
	LogLevel               string
	Timezone               string
}

// NewTenantManager creates a new tenant manager
func NewTenantManager(isolation IsolationLevel, logger *logrus.Logger) *TenantManager {
	if logger == nil {
		logger = logrus.New()
	}

	return &TenantManager{
		tenants:   make(map[string]*Tenant),
		logger:    logger,
		isolation: isolation,
	}
}

// CreateTenant creates a new tenant
func (tm *TenantManager) CreateTenant(ctx context.Context, tenant *Tenant) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.tenants[tenant.ID]; exists {
		return fmt.Errorf("tenant %s already exists", tenant.ID)
	}

	// 設置預設值
	tenant.Status = TenantStatusActive
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()

	// 根據訂閱計劃設置資源限制
	tenant.Limits = tm.getDefaultLimits(tenant.Plan)

	// 初始化使用量
	tenant.Usage = &ResourceUsage{
		LastUpdated: time.Now(),
	}

	// 初始化配置
	if tenant.Config == nil {
		tenant.Config = tm.getDefaultConfig(tenant.Plan)
	}

	tm.tenants[tenant.ID] = tenant
	tm.logger.Infof("Created tenant: %s (%s)", tenant.ID, tenant.Name)

	return nil
}

// GetTenant retrieves a tenant by ID
func (tm *TenantManager) GetTenant(ctx context.Context, tenantID string) (*Tenant, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tenant, exists := tm.tenants[tenantID]
	if !exists {
		return nil, fmt.Errorf("tenant %s not found", tenantID)
	}

	return tenant, nil
}

// UpdateTenant updates a tenant
func (tm *TenantManager) UpdateTenant(ctx context.Context, tenant *Tenant) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.tenants[tenant.ID]; !exists {
		return fmt.Errorf("tenant %s not found", tenant.ID)
	}

	tenant.UpdatedAt = time.Now()
	tm.tenants[tenant.ID] = tenant

	tm.logger.Infof("Updated tenant: %s", tenant.ID)
	return nil
}

// DeleteTenant deletes a tenant
func (tm *TenantManager) DeleteTenant(ctx context.Context, tenantID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tenant, exists := tm.tenants[tenantID]
	if !exists {
		return fmt.Errorf("tenant %s not found", tenantID)
	}

	tenant.Status = TenantStatusDeleted
	tenant.UpdatedAt = time.Now()

	tm.logger.Infof("Deleted tenant: %s", tenantID)
	return nil
}

// CheckResourceLimit checks if tenant has exceeded resource limits
func (tm *TenantManager) CheckResourceLimit(ctx context.Context, tenantID string, resource string) (bool, error) {
	tenant, err := tm.GetTenant(ctx, tenantID)
	if err != nil {
		return false, err
	}

	switch resource {
	case "users":
		return tenant.Usage.CurrentUsers < tenant.Limits.MaxUsers, nil
	case "devices":
		return tenant.Usage.CurrentDevices < tenant.Limits.MaxDevices, nil
	case "events":
		return tenant.Usage.EventsToday < tenant.Limits.MaxEvents, nil
	case "storage":
		return tenant.Usage.StorageUsedGB < float64(tenant.Limits.MaxStorageGB), nil
	case "api_calls":
		return tenant.Usage.APICallsToday < tenant.Limits.MaxAPICallsPerDay, nil
	case "connections":
		return tenant.Usage.ConcurrentConns < tenant.Limits.MaxConcurrentConns, nil
	default:
		return false, fmt.Errorf("unknown resource: %s", resource)
	}
}

// IncrementUsage increments resource usage
func (tm *TenantManager) IncrementUsage(ctx context.Context, tenantID string, resource string, amount int64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tenant, exists := tm.tenants[tenantID]
	if !exists {
		return fmt.Errorf("tenant %s not found", tenantID)
	}

	switch resource {
	case "events":
		tenant.Usage.EventsToday += amount
	case "api_calls":
		tenant.Usage.APICallsToday += amount
	case "storage":
		tenant.Usage.StorageUsedGB += float64(amount)
	}

	tenant.Usage.LastUpdated = time.Now()
	return nil
}

// GetDatabaseName returns the database name for a tenant
func (tm *TenantManager) GetDatabaseName(tenantID string) string {
	switch tm.isolation {
	case IsolationSharedDB:
		return "pandora_shared"
	case IsolationSeparateDB:
		return fmt.Sprintf("pandora_tenant_%s", tenantID)
	case IsolationDedicatedHost:
		return fmt.Sprintf("pandora_dedicated_%s", tenantID)
	default:
		return "pandora_shared"
	}
}

// GetSchemaName returns the schema name for a tenant
func (tm *TenantManager) GetSchemaName(tenantID string) string {
	switch tm.isolation {
	case IsolationSharedDB:
		return fmt.Sprintf("tenant_%s", tenantID)
	default:
		return "public"
	}
}

// getDefaultLimits returns default resource limits based on plan
func (tm *TenantManager) getDefaultLimits(plan SubscriptionPlan) *ResourceLimits {
	switch plan {
	case PlanFree:
		return &ResourceLimits{
			MaxUsers:           5,
			MaxDevices:         10,
			MaxEvents:          10000,
			MaxStorageGB:       1,
			MaxAPICallsPerDay:  1000,
			MaxConcurrentConns: 10,
		}
	case PlanBasic:
		return &ResourceLimits{
			MaxUsers:           20,
			MaxDevices:         50,
			MaxEvents:          100000,
			MaxStorageGB:       10,
			MaxAPICallsPerDay:  10000,
			MaxConcurrentConns: 50,
		}
	case PlanProfessional:
		return &ResourceLimits{
			MaxUsers:           100,
			MaxDevices:         500,
			MaxEvents:          1000000,
			MaxStorageGB:       100,
			MaxAPICallsPerDay:  100000,
			MaxConcurrentConns: 200,
		}
	case PlanEnterprise:
		return &ResourceLimits{
			MaxUsers:           -1, // unlimited
			MaxDevices:         -1,
			MaxEvents:          -1,
			MaxStorageGB:       -1,
			MaxAPICallsPerDay:  -1,
			MaxConcurrentConns: -1,
		}
	default:
		return tm.getDefaultLimits(PlanFree)
	}
}

// getDefaultConfig returns default configuration based on plan
func (tm *TenantManager) getDefaultConfig(plan SubscriptionPlan) *TenantConfig {
	config := &TenantConfig{
		RetentionDays: 30,
		LogLevel:      "info",
		Timezone:      "UTC",
	}

	switch plan {
	case PlanProfessional, PlanEnterprise:
		config.EnableAdvancedFeatures = true
		config.EnableMLDetection = true
		config.EnableAutoResponse = true
		config.RetentionDays = 90
	case PlanBasic:
		config.EnableAdvancedFeatures = true
		config.RetentionDays = 60
	}

	return config
}

// ListTenants lists all tenants
func (tm *TenantManager) ListTenants(ctx context.Context) ([]*Tenant, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tenants := make([]*Tenant, 0, len(tm.tenants))
	for _, tenant := range tm.tenants {
		if tenant.Status != TenantStatusDeleted {
			tenants = append(tenants, tenant)
		}
	}

	return tenants, nil
}

// GetTenantByDomain retrieves a tenant by domain
func (tm *TenantManager) GetTenantByDomain(ctx context.Context, domain string) (*Tenant, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	for _, tenant := range tm.tenants {
		if tenant.Domain == domain && tenant.Status == TenantStatusActive {
			return tenant, nil
		}
	}

	return nil, fmt.Errorf("tenant not found for domain: %s", domain)
}

// SuspendTenant suspends a tenant
func (tm *TenantManager) SuspendTenant(ctx context.Context, tenantID string, reason string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tenant, exists := tm.tenants[tenantID]
	if !exists {
		return fmt.Errorf("tenant %s not found", tenantID)
	}

	tenant.Status = TenantStatusSuspended
	tenant.UpdatedAt = time.Now()
	if tenant.Metadata == nil {
		tenant.Metadata = make(map[string]string)
	}
	tenant.Metadata["suspension_reason"] = reason

	tm.logger.Warnf("Suspended tenant %s: %s", tenantID, reason)
	return nil
}

// ReactivateTenant reactivates a suspended tenant
func (tm *TenantManager) ReactivateTenant(ctx context.Context, tenantID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tenant, exists := tm.tenants[tenantID]
	if !exists {
		return fmt.Errorf("tenant %s not found", tenantID)
	}

	tenant.Status = TenantStatusActive
	tenant.UpdatedAt = time.Now()

	tm.logger.Infof("Reactivated tenant: %s", tenantID)
	return nil
}

// GetTenantStats returns statistics for all tenants
func (tm *TenantManager) GetTenantStats() map[string]interface{} {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	stats := map[string]interface{}{
		"total_tenants":     len(tm.tenants),
		"active_tenants":    0,
		"suspended_tenants": 0,
		"deleted_tenants":   0,
		"by_plan":           make(map[SubscriptionPlan]int),
	}

	for _, tenant := range tm.tenants {
		switch tenant.Status {
		case TenantStatusActive:
			stats["active_tenants"] = stats["active_tenants"].(int) + 1
		case TenantStatusSuspended:
			stats["suspended_tenants"] = stats["suspended_tenants"].(int) + 1
		case TenantStatusDeleted:
			stats["deleted_tenants"] = stats["deleted_tenants"].(int) + 1
		}

		byPlan := stats["by_plan"].(map[SubscriptionPlan]int)
		byPlan[tenant.Plan]++
	}

	return stats
}

