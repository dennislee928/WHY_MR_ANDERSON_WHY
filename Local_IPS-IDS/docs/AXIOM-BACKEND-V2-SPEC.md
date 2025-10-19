# Axiom Backend V2 規格文檔

> **版本**: 2.0.0  
> **日期**: 2025-10-16  
> **狀態**: 設計中

---

## 📋 概述

Axiom Backend V2 是一個統一的 API Gateway 和控制中心，提供對整個 Pandora Box Console IDS-IPS 系統中所有服務的集中管理和控制。

### 核心目標

1. **統一 API Gateway**: 所有服務透過 Axiom Backend 統一訪問
2. **服務編排**: 協調多個服務間的互動
3. **量子功能整合**: 統一接口觸發量子相關功能
4. **配置管理**: 集中管理所有服務配置
5. **日誌聚合**: 收集和處理來自所有源的日誌
6. **資料持久化**: 使用 GORM 管理 PostgreSQL 和 Redis

---

## 🏗️ 架構設計

### 服務清單與功能映射

#### 1. **pandora-agent** (Agent 服務)
**功能 APIs:**
- 查詢 Agent 狀態和健康檢查
- 控制 Agent 行為 (啟動/停止監控)
- 收集 Windows Event Log
- 配置 Agent 參數
- 查詢收集的日誌和指標

#### 2. **prometheus** (指標收集)
**功能 APIs:**
- 查詢即時指標數據
- 查詢歷史指標數據
- 執行 PromQL 查詢
- 管理 Alert Rules
- 管理 Scrape Targets

#### 3. **grafana** (視覺化)
**功能 APIs:**
- 管理 Dashboard (CRUD)
- 管理 Data Sources
- 管理 Alert Notifications
- 快照和分享 Dashboard
- 查詢面板數據

#### 4. **loki** (日誌聚合)
**功能 APIs:**
- 查詢日誌 (LogQL)
- 查詢日誌標籤
- 查詢日誌統計
- 管理日誌保留策略

#### 5. **alertmanager** (告警管理)
**功能 APIs:**
- 查詢活躍告警
- 管理告警靜默
- 管理告警分組
- 配置告警路由
- 查詢告警歷史

#### 6. **redis** (快取)
**功能 APIs:**
- 查詢 Redis 狀態和統計
- 管理快取鍵值
- 查詢快取使用情況
- 清理快取

#### 7. **rabbitmq** (消息隊列)
**功能 APIs:**
- 查詢隊列狀態
- 管理 Exchanges 和 Queues
- 發送測試消息
- 查詢消息統計
- 管理連接和通道

#### 8. **postgres** (資料庫)
**功能 APIs:**
- 查詢資料庫狀態
- 執行資料庫備份
- 查詢連接統計
- 管理資料庫用戶和權限
- 執行 Vacuum 和維護操作

#### 9. **cyber-ai-quantum** (AI/量子服務)
**功能 APIs:**
- 觸發量子密鑰分發 (QKD)
- 執行量子加密/解密
- 提交量子作業到 IBM Quantum
- 查詢量子作業狀態
- 執行 Zero Trust 預測
- 觸發 ML 威脅檢測
- 量子 QSVM/QAOA/QWalk 算法
- 管理量子電路和結果

#### 10. **nginx** (反向代理)
**功能 APIs:**
- 查詢 Nginx 狀態
- 讀取當前配置
- 更新配置
- 重載配置
- 查詢訪問日誌統計
- 管理上游服務

#### 11. **portainer** (容器管理)
**功能 APIs:**
- 查詢容器狀態
- 啟動/停止/重啟容器
- 查詢容器日誌
- 查詢容器資源使用
- 管理容器配置

#### 12. **n8n** (工作流自動化)
**功能 APIs:**
- 管理工作流 (CRUD)
- 觸發工作流執行
- 查詢執行歷史
- 管理 Credentials
- 查詢工作流統計

#### 13. **node-exporter** (系統指標)
**功能 APIs:**
- 查詢系統 CPU/Memory/Disk 指標
- 查詢網路流量統計
- 查詢檔案系統狀態

---

## 🗄️ 資料庫設計

### PostgreSQL Schema (GORM Models)

```go
// 服務狀態表
type Service struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"uniqueIndex;not null"`
    Type        string    // prometheus, grafana, etc.
    Status      string    // healthy, unhealthy, unknown
    URL         string    
    Version     string
    LastCheck   time.Time
    Config      datatypes.JSON
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// 配置歷史表
type ConfigHistory struct {
    ID         uint   `gorm:"primaryKey"`
    ServiceID  uint   `gorm:"index"`
    Service    Service
    ConfigType string // nginx, agent, etc.
    Content    string `gorm:"type:text"`
    AppliedBy  string
    AppliedAt  time.Time
    Status     string // pending, applied, failed
}

// 量子作業表
type QuantumJob struct {
    ID            uint   `gorm:"primaryKey"`
    JobID         string `gorm:"uniqueIndex"`
    Type          string // qkd, qsvm, qaoa, etc.
    Status        string // pending, running, completed, failed
    Backend       string // ibm_quantum, simulator
    Circuit       string `gorm:"type:text"`
    Result        datatypes.JSON
    SubmittedAt   time.Time
    CompletedAt   *time.Time
    Error         string
}

// Windows 日誌表
type WindowsLog struct {
    ID          uint   `gorm:"primaryKey"`
    AgentID     string `gorm:"index"`
    LogType     string // System, Security, Application
    Source      string
    EventID     int
    Level       string
    Message     string `gorm:"type:text"`
    TimeCreated time.Time `gorm:"index"`
    ReceivedAt  time.Time
    Metadata    datatypes.JSON
}

// 告警表
type Alert struct {
    ID          uint   `gorm:"primaryKey"`
    AlertName   string `gorm:"index"`
    Severity    string
    Source      string
    Message     string `gorm:"type:text"`
    Status      string // active, resolved, acknowledged
    CreatedAt   time.Time `gorm:"index"`
    UpdatedAt   time.Time
    ResolvedAt  *time.Time
    ResolvedBy  string
    Labels      datatypes.JSON
}

// API 請求日誌表
type APILog struct {
    ID           uint   `gorm:"primaryKey"`
    Method       string
    Path         string
    Status       int
    Duration     int64 // microseconds
    ClientIP     string
    UserAgent    string
    RequestBody  string `gorm:"type:text"`
    ResponseBody string `gorm:"type:text"`
    Error        string
    CreatedAt    time.Time `gorm:"index"`
}
```

### Redis Schema

```
# 服務健康狀態快取 (TTL: 30s)
service:health:{service_name} -> {"status": "healthy", "last_check": "2025-10-16T10:00:00Z"}

# 即時指標快取 (TTL: 10s)
metrics:realtime:{service_name} -> {"cpu": 15.5, "memory": 42.3, ...}

# 量子作業狀態快取 (TTL: 5min)
quantum:job:{job_id} -> {"status": "running", "progress": 75}

# API 速率限制 (TTL: 1min)
ratelimit:api:{client_ip}:{endpoint} -> 100

# 會話管理 (TTL: 24h)
session:{session_id} -> {"user": "admin", "expires": "2025-10-17T10:00:00Z"}

# 即時統計計數器
counter:api:requests:total
counter:threats:detected:total
counter:quantum:jobs:completed
```

---

## 📡 API 端點設計

### Base URL
```
http://localhost:3001/api/v2
```

### 1. 服務管理 APIs

#### 1.1 查詢所有服務狀態
```http
GET /api/v2/services
Response: {
    "services": [
        {
            "name": "prometheus",
            "type": "monitoring",
            "status": "healthy",
            "url": "http://prometheus:9090",
            "last_check": "2025-10-16T10:00:00Z",
            "metrics": {
                "uptime": "72h",
                "memory_usage": "256MB"
            }
        }
    ],
    "total": 13,
    "healthy": 12,
    "unhealthy": 1
}
```

#### 1.2 查詢單個服務狀態
```http
GET /api/v2/services/{service_name}
```

#### 1.3 健康檢查單個服務
```http
POST /api/v2/services/{service_name}/health-check
```

#### 1.4 重啟服務 (透過 Portainer)
```http
POST /api/v2/services/{service_name}/restart
```

### 2. Prometheus 管理 APIs

#### 2.1 執行 PromQL 查詢
```http
POST /api/v2/prometheus/query
Body: {
    "query": "up{job='node-exporter'}",
    "time": "2025-10-16T10:00:00Z"
}
```

#### 2.2 執行範圍查詢
```http
POST /api/v2/prometheus/query-range
Body: {
    "query": "rate(http_requests_total[5m])",
    "start": "2025-10-16T09:00:00Z",
    "end": "2025-10-16T10:00:00Z",
    "step": "15s"
}
```

#### 2.3 管理 Alert Rules
```http
GET /api/v2/prometheus/rules
POST /api/v2/prometheus/rules
PUT /api/v2/prometheus/rules/{rule_name}
DELETE /api/v2/prometheus/rules/{rule_name}
```

### 3. Grafana 管理 APIs

#### 3.1 查詢 Dashboards
```http
GET /api/v2/grafana/dashboards
```

#### 3.2 創建 Dashboard
```http
POST /api/v2/grafana/dashboards
Body: {
    "dashboard": {...},
    "folderUid": "..."
}
```

#### 3.3 查詢 Dashboard 數據
```http
GET /api/v2/grafana/dashboards/{uid}/data
```

### 4. Loki 日誌查詢 APIs

#### 4.1 查詢日誌
```http
POST /api/v2/loki/query
Body: {
    "query": "{job=\"varlogs\"} |= \"error\"",
    "limit": 100,
    "start": "2025-10-16T09:00:00Z",
    "end": "2025-10-16T10:00:00Z"
}
```

#### 4.2 查詢日誌標籤
```http
GET /api/v2/loki/labels
```

### 5. 量子功能 APIs

#### 5.1 提交量子密鑰分發作業
```http
POST /api/v2/quantum/qkd/generate
Body: {
    "key_length": 256,
    "backend": "ibm_quantum" // or "simulator"
}
Response: {
    "job_id": "qkd-job-123",
    "status": "pending",
    "estimated_time": "30s"
}
```

#### 5.2 查詢量子作業狀態
```http
GET /api/v2/quantum/jobs/{job_id}
Response: {
    "job_id": "qkd-job-123",
    "status": "completed",
    "result": {...},
    "submitted_at": "2025-10-16T10:00:00Z",
    "completed_at": "2025-10-16T10:00:30Z"
}
```

#### 5.3 執行 Zero Trust 預測
```http
POST /api/v2/quantum/zerotrust/predict
Body: {
    "user_id": "user123",
    "ip_address": "192.168.1.100",
    "features": {...}
}
```

#### 5.4 執行 QSVM 分類
```http
POST /api/v2/quantum/qsvm/classify
Body: {
    "features": [0.1, 0.2, ...],
    "backend": "ibm_quantum"
}
```

### 6. Nginx 配置管理 APIs

#### 6.1 查詢當前配置
```http
GET /api/v2/nginx/config
Response: {
    "config": "user nginx; worker_processes auto; ...",
    "last_modified": "2025-10-16T09:00:00Z"
}
```

#### 6.2 更新配置
```http
PUT /api/v2/nginx/config
Body: {
    "config": "user nginx; ...",
    "validate": true // 先驗證配置
}
```

#### 6.3 重載配置
```http
POST /api/v2/nginx/reload
```

#### 6.4 查詢 Nginx 狀態
```http
GET /api/v2/nginx/status
Response: {
    "active_connections": 42,
    "accepts": 1523,
    "handled": 1523,
    "requests": 3046,
    "reading": 0,
    "writing": 3,
    "waiting": 39
}
```

### 7. Windows 日誌 APIs

#### 7.1 接收 Agent 日誌
```http
POST /api/v2/logs/windows
Body: {
    "agent_id": "agent-001",
    "logs": [
        {
            "log_type": "Security",
            "source": "Microsoft-Windows-Security-Auditing",
            "event_id": 4624,
            "level": "Information",
            "message": "An account was successfully logged on",
            "time_created": "2025-10-16T10:00:00Z",
            "metadata": {...}
        }
    ]
}
```

#### 7.2 查詢 Windows 日誌
```http
GET /api/v2/logs/windows
Query params:
- agent_id: 過濾 Agent
- log_type: System/Security/Application
- level: Information/Warning/Error
- start_time: 開始時間
- end_time: 結束時間
- limit: 數量限制
```

#### 7.3 搜索 Windows 日誌
```http
POST /api/v2/logs/windows/search
Body: {
    "query": "failed login",
    "filters": {
        "log_type": "Security",
        "event_id": [4625, 4771]
    },
    "time_range": "24h"
}
```

### 8. RabbitMQ 管理 APIs

#### 8.1 查詢隊列狀態
```http
GET /api/v2/rabbitmq/queues
```

#### 8.2 發送測試消息
```http
POST /api/v2/rabbitmq/publish
Body: {
    "exchange": "pandora-events",
    "routing_key": "test",
    "message": {...}
}
```

### 9. Portainer 容器管理 APIs

#### 9.1 查詢所有容器
```http
GET /api/v2/containers
```

#### 9.2 啟動/停止容器
```http
POST /api/v2/containers/{container_id}/start
POST /api/v2/containers/{container_id}/stop
POST /api/v2/containers/{container_id}/restart
```

#### 9.3 查詢容器日誌
```http
GET /api/v2/containers/{container_id}/logs
Query params:
- tail: 最後 N 行
- since: 起始時間
```

---

## 🔧 技術實現

### 後端技術棧
- **語言**: Go 1.21+
- **框架**: Gin
- **ORM**: GORM
- **資料庫**: PostgreSQL 15, Redis 7
- **消息隊列**: RabbitMQ 3.12
- **文檔**: Swagger/OpenAPI 3.0

### 前端技術棧
- **框架**: Next.js 14
- **語言**: TypeScript
- **UI庫**: Tailwind CSS, shadcn/ui
- **狀態管理**: React Hooks
- **API客戶端**: Fetch API, Auto-generated from Swagger

---

## 📊 監控和日誌

### Prometheus 指標

```
# Axiom Backend 指標
axiom_api_requests_total{method, endpoint, status}
axiom_api_duration_seconds{method, endpoint}
axiom_service_health{service_name, status}
axiom_quantum_jobs_total{type, status}
axiom_windows_logs_received_total{agent_id, log_type}

# Redis 快取指標
axiom_redis_hits_total
axiom_redis_misses_total
axiom_redis_keys_total

# Database 指標
axiom_db_connections_active
axiom_db_query_duration_seconds
```

### 日誌格式

```json
{
  "timestamp": "2025-10-16T10:00:00Z",
  "level": "info",
  "service": "axiom-backend",
  "endpoint": "/api/v2/quantum/jobs",
  "method": "POST",
  "client_ip": "192.168.1.100",
  "duration_ms": 125,
  "status": 201,
  "message": "Quantum job created successfully",
  "job_id": "qkd-job-123"
}
```

---

## 🔒 安全性

### 認證和授權
- JWT Token 認證
- Role-Based Access Control (RBAC)
- API Key 管理

### 資料安全
- 所有敏感資料加密儲存
- 使用 TLS/mTLS 通訊
- 密鑰定期輪換

### 速率限制
- 基於 IP 和用戶的速率限制
- 使用 Redis 實現分散式限流

---

## 📦 部署

### Docker Compose 配置
```yaml
axiom-be:
  build:
    context: ..
    dockerfile: Application/docker/axiom-be-v2.dockerfile
  container_name: axiom-be-v2
  ports:
    - "3001:3001"
  environment:
    - DB_HOST=postgres
    - DB_PORT=5432
    - DB_NAME=pandora_db
    - REDIS_HOST=redis
    - REDIS_PORT=6379
  depends_on:
    - postgres
    - redis
    - rabbitmq
```

---

## 🚀 開發路線圖

### Phase 1: 基礎架構 (1-2 天)
- [x] 設計資料庫 Schema
- [ ] 實現 GORM Models
- [ ] 設置 Redis 快取
- [ ] 設置基本 API 框架

### Phase 2: 服務整合 (2-3 天)
- [ ] Prometheus 集成
- [ ] Grafana 集成
- [ ] Loki 集成
- [ ] 其他服務集成

### Phase 3: 量子功能 (1-2 天)
- [ ] 量子服務代理 APIs
- [ ] 作業管理系統

### Phase 4: Nginx & Windows Logs (1-2 天)
- [ ] Nginx 配置管理
- [ ] Windows 日誌收集和查詢

### Phase 5: Frontend (2-3 天)
- [ ] Next.js UI 更新
- [ ] 所有新功能 UI

### Phase 6: 測試和文檔 (1-2 天)
- [ ] 單元測試
- [ ] 集成測試
- [ ] Swagger 文檔
- [ ] 用戶手冊

---

## 📝 總結

Axiom Backend V2 將成為整個 Pandora Box Console IDS-IPS 系統的核心控制中心，提供統一、強大、易用的 API 接口，支持所有服務的管理和編排。

