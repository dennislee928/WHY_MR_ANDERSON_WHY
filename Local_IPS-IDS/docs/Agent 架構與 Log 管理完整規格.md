# Axiom Backend V2 - Agent 架構與 Log 管理完整規格

> **版本** : 2.5.0
>
> **日期** : 2025-10-16
>
> **焦點** : Agent 雙模式連接 + 持續日誌上傳 + 合規性

---

## 📋 目錄

1. [Agent 架構設計](https://claude.ai/chat/1dcab708-3d68-4bb7-901e-cbbece89e27a#1-agent-%E6%9E%B6%E6%A7%8B%E8%A8%AD%E8%A8%88)
2. [Log 儲存策略](https://claude.ai/chat/1dcab708-3d68-4bb7-901e-cbbece89e27a#2-log-%E5%84%B2%E5%AD%98%E7%AD%96%E7%95%A5)
3. [合規性實現](https://claude.ai/chat/1dcab708-3d68-4bb7-901e-cbbece89e27a#3-%E5%90%88%E8%A6%8F%E6%80%A7%E5%AF%A6%E7%8F%BE)
4. [API 規格](https://claude.ai/chat/1dcab708-3d68-4bb7-901e-cbbece89e27a#4-api-%E8%A6%8F%E6%A0%BC)
5. [安全與完整性](https://claude.ai/chat/1dcab708-3d68-4bb7-901e-cbbece89e27a#5-%E5%AE%89%E5%85%A8%E8%88%87%E5%AE%8C%E6%95%B4%E6%80%A7)
6. [負載平衡](https://claude.ai/chat/1dcab708-3d68-4bb7-901e-cbbece89e27a#6-%E8%B2%A0%E8%BC%89%E5%B9%B3%E8%A1%A1)

---

## 1. Agent 架構設計

### 1.1 雙模式連接架構

```
┌─────────────────────────────────────────────────────────────┐
│                    Axiom Backend V2                          │
│  ┌────────────┐  ┌──────────────┐  ┌───────────────────┐   │
│  │ Agent API  │  │ Log Ingestion│  │ Compliance Engine │   │
│  │ Gateway    │  │ Pipeline     │  │                   │   │
│  └─────┬──────┘  └──────┬───────┘  └─────────┬─────────┘   │
│        │                │                     │              │
│        │  ┌─────────────▼──────────────┐     │              │
│        │  │   Log Storage Layer         │◄────┘              │
│        │  │  - Hot: Redis Streams       │                    │
│        │  │  - Warm: Loki               │                    │
│        │  │  - Cold: PostgreSQL         │                    │
│        │  │  - Archive: S3/MinIO        │                    │
│        │  └─────────────────────────────┘                    │
└────────┼──────────────────────────────────────────────────────┘
         │
    ┌────┴──────────────────────────────────────┐
    │                                            │
    │  ┌─────────────────┐                      │
    │  │     Nginx        │                      │
    │  │  (Reverse Proxy) │                      │
    │  │  - Rate Limiting │                      │
    │  │  - mTLS          │                      │
    │  │  - WAF           │                      │
    │  └────────┬─────────┘                      │
    │           │                                 │
┌───▼───────────▼─────┐              ┌───────────▼─────────┐
│  External Agents     │              │  Internal Agents     │
│  (通過 Nginx)        │              │  (直連 Backend)      │
│                      │              │                      │
│  • 遠端辦公室        │              │  • 數據中心內部      │
│  • 分支機構          │              │  • Kubernetes Pod    │
│  • 雲端實例          │              │  • Docker 容器       │
│  • 客戶端點          │              │  • 本地服務器        │
└──────────────────────┘              └──────────────────────┘
```

### 1.2 Agent 連接模式規格

#### **External Agent (外部連接)**

 **連接路徑** : `Agent → Internet → Nginx → Axiom Backend`

**配置檔案** (`agent-config-external.yaml`):

```yaml
agent:
  mode: external
  id: "agent-external-{unique-id}"
  
connection:
  endpoint: "https://axiom.yourdomain.com/api/v2/agent"
  protocol: https
  port: 443
  through_nginx: true
  
authentication:
  method: mtls  # Mutual TLS
  client_cert: "/etc/pandora/certs/client.crt"
  client_key: "/etc/pandora/certs/client.key"
  ca_cert: "/etc/pandora/certs/ca.crt"
  api_key: "${AGENT_API_KEY}"  # Backup auth method
  
upload:
  method: streaming  # continuous upload
  batch_size: 100    # events per batch
  flush_interval: 10s
  max_retry: 5
  retry_backoff: exponential
  compression: gzip
  
buffer:
  type: persistent
  path: "/var/lib/pandora/buffer"
  max_size: 1GB
  overflow_strategy: drop_oldest  # or block, drop_newest
  
security:
  encrypt_in_transit: true
  encrypt_at_rest: true
  key_rotation: 24h
```

 **特性** :

* ✅ 通過 Nginx WAF 防護
* ✅ Rate Limiting 保護
* ✅ mTLS 雙向認證
* ✅ 支援動態 IP
* ✅ 本地緩衝防止資料遺失
* ✅ 壓縮傳輸節省頻寬

---

#### **Internal Agent (內部連接)**

 **連接路徑** : `Agent → Internal Network → Axiom Backend (直連)`

**配置檔案** (`agent-config-internal.yaml`):

```yaml
agent:
  mode: internal
  id: "agent-internal-{unique-id}"
  
connection:
  endpoint: "http://axiom-backend.internal:8080/api/v2/agent"
  protocol: http  # or https with internal CA
  port: 8080
  through_nginx: false
  direct_connect: true
  
authentication:
  method: api_key  # Simpler auth for internal
  api_key: "${AGENT_API_KEY}"
  
upload:
  method: streaming
  batch_size: 500    # Larger batches in internal network
  flush_interval: 5s # Faster flush
  max_retry: 3
  compression: none  # No compression needed in internal network
  
buffer:
  type: memory       # Memory-only for internal (faster)
  max_size: 256MB
  overflow_strategy: block  # Block until buffer space available
  
security:
  encrypt_in_transit: false  # Optional for internal
  encrypt_at_rest: false     # Data encrypted at storage layer
```

 **特性** :

* ✅ 低延遲直連
* ✅ 無需 mTLS 開銷
* ✅ 更大批次大小
* ✅ 記憶體緩衝（更快）
* ✅ 可選加密（降低 CPU 負載）
* ✅ Service Discovery 整合

---

### 1.3 Agent 註冊與生命週期管理

#### Agent 註冊流程

```
┌─────────┐         ┌────────────┐         ┌────────────┐
│  Agent  │         │   Nginx    │         │   Axiom    │
└────┬────┘         └─────┬──────┘         └─────┬──────┘
     │                    │                       │
     │  1. Register       │                       │
     ├────────────────────┼──────────────────────►│
     │    (External)      │                       │
     │  or Direct (Internal)                      │
     │                    │                       │
     │                    │   2. Validate & Store │
     │                    │   ◄───────────────────┤
     │                    │                       │
     │  3. Return Credentials                     │
     │◄───────────────────┼───────────────────────┤
     │   - Agent ID                               │
     │   - API Key                                │
     │   - Cert (if external)                     │
     │                    │                       │
     │  4. Heartbeat (every 30s)                  │
     ├────────────────────┼──────────────────────►│
     │                    │                       │
     │  5. Log Stream                             │
     ├═══════════════════►│═══════════════════════►│
     │    (Continuous)    │                       │
```

---

## 2. Log 儲存策略

### 2.1 四層儲存架構 (Hot-Warm-Cold-Archive)

```
┌──────────────────────────────────────────────────────────────┐
│                    Log Storage Tiers                         │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  HOT (實時)   │  │ WARM (近期)  │  │ COLD (歷史)  │      │
│  │              │  │              │  │              │      │
│  │ Redis Streams│→ │    Loki      │→ │  PostgreSQL  │      │
│  │              │  │              │  │              │      │
│  │ 最近 1 小時   │  │  最近 7 天   │  │  8-90 天     │      │
│  │ 100% 數據    │  │  100% 數據   │  │  結構化索引  │      │
│  │ 超快查詢     │  │  快速查詢    │  │  中速查詢    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│         │                  │                  │             │
│         │                  │                  │             │
│         └──────────────────┴──────────────────┘             │
│                            │                                │
│                            ▼                                │
│                  ┌──────────────────┐                       │
│                  │ ARCHIVE (封存)    │                       │
│                  │                  │                       │
│                  │  S3/MinIO/WORM   │                       │
│                  │                  │                       │
│                  │  91+ 天          │                       │
│                  │  壓縮 + 加密     │                       │
│                  │  不可變儲存      │                       │
│                  │  法規遵循        │                       │
│                  └──────────────────┘                       │
└──────────────────────────────────────────────────────────────┘
```

### 2.2 詳細儲存規格

#### **Tier 1: Hot Storage (Redis Streams)**

 **用途** : 實時日誌接收與即時查詢

 **技術規格** :

```yaml
redis:
  streams:
    name_pattern: "logs:agent:{agentId}:{date}"
    max_length: 100000  # per stream
    retention: 1h
  
  features:
    - consumer_groups      # Multiple consumers
    - exactly_once_delivery
    - pending_entry_list   # Retry mechanism
    - time_series_support
  
  performance:
    expected_throughput: 100k events/sec
    query_latency: <10ms
  
  high_availability:
    mode: redis_cluster
    replicas: 3
    sentinel: true
```

 **資料結構** :

```json
{
  "stream": "logs:agent:ext-001:2025-10-16",
  "entry": {
    "id": "1729065600000-0",
    "data": {
      "timestamp": "2025-10-16T10:00:00Z",
      "agent_id": "ext-001",
      "agent_mode": "external",
      "event_type": "windows_event_log",
      "source": "Security",
      "event_id": 4624,
      "level": "Information",
      "computer": "WS2019-SERVER",
      "message": "An account was successfully logged on",
      "metadata": {
        "upload_timestamp": "2025-10-16T10:00:01Z",
        "batch_id": "batch-12345",
        "sequence": 1
      }
    }
  }
}
```

 **自動轉移規則** :

```go
// Pseudo-code
if entry.age > 1hour {
    transferToLoki(entry)
    deleteFromRedis(entry)
}
```

---

#### **Tier 2: Warm Storage (Loki)**

 **用途** : 近期日誌查詢與分析

 **技術規格** :

```yaml
loki:
  storage:
    chunk_store: filesystem  # or S3, GCS
    retention: 7d
  
  index:
    period: 24h
    prefix: "loki_index_"
  
  ingester:
    chunk_idle_period: 30m
    chunk_retain_period: 1m
    max_chunk_age: 1h
  
  query:
    max_concurrent: 100
    timeout: 5m
  
  labels:
    - agent_id
    - agent_mode
    - event_type
    - source
    - level
    - computer
```

 **LogQL 查詢範例** :

```logql
{agent_mode="external", event_type="windows_event_log"} 
  |= "error" 
  | json 
  | event_id="4625" 
  | line_format "Failed login from {{.computer}}"
```

 **資料壓縮** :

* Compression: Snappy
* Average compression ratio: 10:1
* Storage cost: ~$0.023/GB/month

---

#### **Tier 3: Cold Storage (PostgreSQL)**

 **用途** : 長期儲存、合規查詢、關聯分析

 **資料庫架構** :

```sql
-- Main event log table (partitioned by date)
CREATE TABLE event_logs (
    id BIGSERIAL,
    timestamp TIMESTAMPTZ NOT NULL,
    agent_id VARCHAR(64) NOT NULL,
    agent_mode VARCHAR(16) NOT NULL,  -- 'external' or 'internal'
    event_type VARCHAR(64) NOT NULL,
    source VARCHAR(128),
    event_id INTEGER,
    level VARCHAR(32),
    computer VARCHAR(256),
    message TEXT,
    raw_data JSONB,
  
    -- Compliance fields
    retention_until TIMESTAMPTZ,
    archived BOOLEAN DEFAULT false,
    integrity_hash VARCHAR(64),  -- SHA-256
  
    -- Metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
  
    PRIMARY KEY (id, timestamp)
) PARTITION BY RANGE (timestamp);

-- Create monthly partitions
CREATE TABLE event_logs_2025_10 PARTITION OF event_logs
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

-- Indexes for fast querying
CREATE INDEX idx_agent_id ON event_logs(agent_id, timestamp DESC);
CREATE INDEX idx_event_type ON event_logs(event_type, timestamp DESC);
CREATE INDEX idx_computer ON event_logs(computer, timestamp DESC);
CREATE INDEX idx_compliance ON event_logs(retention_until, archived);
CREATE INDEX idx_integrity ON event_logs(integrity_hash);

-- Full-text search
CREATE INDEX idx_message_fts ON event_logs USING gin(to_tsvector('english', message));

-- JSONB index for raw_data queries
CREATE INDEX idx_raw_data ON event_logs USING gin(raw_data);
```

 **表分割策略** :

```sql
-- Automatic partition management
CREATE OR REPLACE FUNCTION create_monthly_partition()
RETURNS void AS $$
DECLARE
    partition_date DATE;
    partition_name TEXT;
BEGIN
    partition_date := DATE_TRUNC('month', NOW() + INTERVAL '1 month');
    partition_name := 'event_logs_' || TO_CHAR(partition_date, 'YYYY_MM');
  
    EXECUTE format('
        CREATE TABLE IF NOT EXISTS %I PARTITION OF event_logs
        FOR VALUES FROM (%L) TO (%L)',
        partition_name,
        partition_date,
        partition_date + INTERVAL '1 month'
    );
END;
$$ LANGUAGE plpgsql;

-- Run daily via cron
SELECT cron.schedule('create-partition', '0 0 * * *', 'SELECT create_monthly_partition()');
```

 **資料完整性** :

```sql
-- Trigger to calculate integrity hash
CREATE OR REPLACE FUNCTION calculate_integrity_hash()
RETURNS TRIGGER AS $$
BEGIN
    NEW.integrity_hash := encode(
        digest(
            NEW.timestamp::text || 
            NEW.agent_id || 
            NEW.event_type || 
            NEW.message,
            'sha256'
        ),
        'hex'
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_integrity_hash
    BEFORE INSERT ON event_logs
    FOR EACH ROW
    EXECUTE FUNCTION calculate_integrity_hash();
```

 **合規查詢表** :

```sql
-- Audit trail for compliance
CREATE TABLE audit_access_log (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    user_id VARCHAR(64),
    action VARCHAR(64),  -- 'query', 'export', 'delete'
    query_text TEXT,
    record_count INTEGER,
    ip_address INET,
    justification TEXT,  -- Required for GDPR
    approved_by VARCHAR(64)
);

-- Data retention policy table
CREATE TABLE retention_policies (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(64),
    retention_days INTEGER,
    legal_hold BOOLEAN DEFAULT false,
    regulation VARCHAR(32),  -- 'GDPR', 'PCI-DSS', 'HIPAA'
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Sample policies
INSERT INTO retention_policies (event_type, retention_days, regulation) VALUES
    ('windows_event_log', 90, 'PCI-DSS'),
    ('security_alert', 365, 'GDPR'),
    ('access_log', 180, 'SOX');
```

---

#### **Tier 4: Archive Storage (S3/MinIO/WORM)**

 **用途** : 長期封存、法規遵循、不可變儲存

 **技術規格** :

```yaml
archive:
  storage:
    type: s3_compatible  # AWS S3, MinIO, Wasabi
    bucket: "axiom-logs-archive"
    region: "us-east-1"
    storage_class: "GLACIER"  # or DEEP_ARCHIVE
  
  compliance:
    worm: true  # Write-Once-Read-Many
    object_lock: true
    retention_mode: "COMPLIANCE"  # Cannot be deleted even by admin
    retention_years: 7  # Default for most regulations
  
  encryption:
    method: "AES-256-GCM"
    key_management: "AWS KMS"  # or HashiCorp Vault
    key_rotation: true
  
  lifecycle:
    archive_after_days: 90
    transition_to_deep_archive: 365
    delete_after_years: 7  # Configurable by policy
  
  indexing:
    metadata_only: true
    stored_in: postgresql
    manifest_files: true
```

 **封存檔案結構** :

```
s3://axiom-logs-archive/
├── 2025/
│   ├── 10/
│   │   ├── 16/
│   │   │   ├── agent-ext-001/
│   │   │   │   ├── logs-000001.parquet.gz.encrypted
│   │   │   │   ├── logs-000001.manifest.json
│   │   │   │   └── logs-000001.checksum.sha256
│   │   │   └── agent-int-001/
│   │   │       └── ...
│   │   └── manifest-2025-10-16.json
│   └── manifest-2025-10.json
└── retention-policies.json
```

 **Manifest 檔案範例** :

```json
{
  "file": "logs-000001.parquet.gz.encrypted",
  "date": "2025-10-16",
  "agent_id": "ext-001",
  "record_count": 1000000,
  "size_bytes": 52428800,
  "compression_ratio": 12.5,
  "encryption": "AES-256-GCM",
  "integrity": {
    "algorithm": "SHA-256",
    "checksum": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  },
  "compliance": {
    "regulations": ["PCI-DSS", "GDPR"],
    "retention_until": "2032-10-16",
    "worm_enabled": true,
    "legal_hold": false
  },
  "time_range": {
    "start": "2025-10-16T00:00:00Z",
    "end": "2025-10-16T23:59:59Z"
  }
}
```

---

### 2.3 資料流轉管道

```go
// Pseudo-code for data tiering pipeline
package pipeline

type LogTieringPipeline struct {
    redis      *RedisClient
    loki       *LokiClient
    postgres   *PostgresClient
    s3         *S3Client
    config     *TieringConfig
}

func (p *LogTieringPipeline) Start() {
    // Stage 1: Redis → Loki (every 5 minutes)
    go p.scheduleTask("redis-to-loki", 5*time.Minute, func() {
        entries := p.redis.GetEntriesOlderThan(1 * time.Hour)
        p.loki.BatchInsert(entries)
        p.redis.Delete(entries)
    })
  
    // Stage 2: Loki → PostgreSQL (every day)
    go p.scheduleTask("loki-to-postgres", 24*time.Hour, func() {
        logs := p.loki.GetLogsOlderThan(7 * 24 * time.Hour)
        p.postgres.BulkInsert(logs)
        p.loki.Delete(logs)
    })
  
    // Stage 3: PostgreSQL → S3 Archive (every week)
    go p.scheduleTask("postgres-to-archive", 7*24*time.Hour, func() {
        logs := p.postgres.GetLogsOlderThan(90 * 24 * time.Hour)
      
        // Compress and encrypt
        archive := p.compressAndEncrypt(logs)
      
        // Generate manifest
        manifest := p.generateManifest(archive)
      
        // Upload with WORM
        p.s3.UploadWithObjectLock(archive, manifest)
      
        // Mark as archived in PostgreSQL
        p.postgres.MarkAsArchived(logs)
    })
  
    // Stage 4: Integrity verification (daily)
    go p.scheduleTask("integrity-check", 24*time.Hour, func() {
        p.verifyIntegrity()
    })
}

func (p *LogTieringPipeline) verifyIntegrity() {
    // Check PostgreSQL integrity hashes
    tampered := p.postgres.FindTamperedRecords()
    if len(tampered) > 0 {
        p.alertSecurityTeam(tampered)
    }
  
    // Verify S3 archive checksums
    archives := p.s3.ListArchives()
    for _, archive := range archives {
        if !p.s3.VerifyChecksum(archive) {
            p.alertSecurityTeam(archive)
        }
    }
}
```

---

## 3. 合規性實現

### 3.1 多法規支援矩陣

| 需求                   | GDPR | PCI-DSS | HIPAA | SOX | ISO 27001 | 實現方式                |
| ---------------------- | ---- | ------- | ----- | --- | --------- | ----------------------- |
| **資料加密**     | ✅   | ✅      | ✅    | ✅  | ✅        | AES-256-GCM (傳輸+靜態) |
| **訪問控制**     | ✅   | ✅      | ✅    | ✅  | ✅        | RBAC + Attribute-based  |
| **審計日誌**     | ✅   | ✅      | ✅    | ✅  | ✅        | 不可變審計表            |
| **資料保留**     | ✅   | ✅      | ✅    | ✅  | ✅        | 自動化保留策略          |
| **資料刪除**     | ✅   | ❌      | ✅    | ❌  | ✅        | 安全刪除 + 驗證         |
| **個人資料識別** | ✅   | ✅      | ✅    | ❌  | ✅        | PII 檢測引擎            |
| **資料匿名化**   | ✅   | ❌      | ✅    | ❌  | ✅        | 自動脫敏                |
| **完整性驗證**   | ✅   | ✅      | ✅    | ✅  | ✅        | SHA-256 hash chain      |
| **不可抵賴性**   | ✅   | ✅      | ❌    | ✅  | ✅        | 數位簽章                |
| **資料可攜性**   | ✅   | ❌      | ❌    | ❌  | ❌        | 標準格式匯出            |

### 3.2 GDPR 合規實現

#### **個人資料識別 (PII Detection)**

```sql
-- PII detection patterns
CREATE TABLE pii_patterns (
    id SERIAL PRIMARY KEY,
    pattern_type VARCHAR(64),  -- 'email', 'ssn', 'credit_card', etc.
    regex VARCHAR(512),
    description TEXT
);

INSERT INTO pii_patterns (pattern_type, regex, description) VALUES
    ('email', '[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}', 'Email addresses'),
    ('credit_card', '\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b', 'Credit card numbers'),
    ('ssn', '\b\d{3}-\d{2}-\d{4}\b', 'Social Security Numbers'),
    ('ip_address', '\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b', 'IP addresses');

-- Function to detect PII
CREATE OR REPLACE FUNCTION detect_pii(text_content TEXT)
RETURNS JSONB AS $$
DECLARE
    found_pii JSONB := '[]'::jsonb;
    pattern RECORD;
BEGIN
    FOR pattern IN SELECT * FROM pii_patterns LOOP
        IF text_content ~ pattern.regex THEN
            found_pii := found_pii || jsonb_build_object(
                'type', pattern.pattern_type,
                'description', pattern.description
            );
        END IF;
    END LOOP;
    RETURN found_pii;
END;
$$ LANGUAGE plpgsql;

-- PII tracking table
CREATE TABLE pii_occurrences (
    id BIGSERIAL PRIMARY KEY,
    log_id BIGINT REFERENCES event_logs(id),
    pii_type VARCHAR(64),
    field_name VARCHAR(128),
    detected_at TIMESTAMPTZ DEFAULT NOW(),
    masked BOOLEAN DEFAULT false
);
```

#### **資料匿名化**

```go
package anonymization

import (
    "crypto/sha256"
    "encoding/hex"
    "regexp"
)

type Anonymizer struct {
    patterns map[string]*regexp.Regexp
    salt     string
}

func NewAnonymizer(salt string) *Anonymizer {
    return &Anonymizer{
        patterns: map[string]*regexp.Regexp{
            "email":       regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
            "credit_card": regexp.MustCompile(`\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`),
            "ssn":         regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),
            "ip":          regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`),
        },
        salt: salt,
    }
}

func (a *Anonymizer) Mask(text string, piiType string) string {
    pattern := a.patterns[piiType]
    return pattern.ReplaceAllStringFunc(text, func(match string) string {
        return a.hash(match)
    })
}

func (a *Anonymizer) hash(value string) string {
    h := sha256.New()
    h.Write([]byte(value + a.salt))
    return "REDACTED_" + hex.EncodeToString(h.Sum(nil))[:16]
}

// Example usage
func (a *Anonymizer) AnonymizeLog(log *EventLog) {
    log.Message = a.Mask(log.Message, "email")
    log.Message = a.Mask(log.Message, "credit_card")
    log.Message = a.Mask(log.Message, "ssn")
    log.Message = a.Mask(log.Message, "ip")
}
```

#### **刪除權 (Right to Erasure)**

```sql
-- GDPR deletion request table
CREATE TABLE gdpr_deletion_requests (
    id SERIAL PRIMARY KEY,
    request_id UUID DEFAULT gen_random_uuid(),
    subject_identifier VARCHAR(256),  -- Email, user ID, etc.
    requested_by VARCHAR(128),
    requested_at TIMESTAMPTZ DEFAULT NOW(),
    approved_by VARCHAR(128),
    approved_at TIMESTAMPTZ,
    status VARCHAR(32),  -- 'pending', 'approved', 'completed', 'rejected'
    completion_date TIMESTAMPTZ,
    verification_hash VARCHAR(64),
    notes TEXT
);

-- Function to securely delete data
CREATE OR REPLACE FUNCTION gdpr_delete_subject_data(
    subject_id VARCHAR(256),
    verification_hash VARCHAR(64)
)
RETURNS TABLE(deleted_count INTEGER, affected_tables TEXT[]) AS $$
```
