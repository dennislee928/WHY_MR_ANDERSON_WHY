-- Migration 002: Agent 和合規性架構
-- 版本: 3.1.0
-- 日期: 2025-10-16

-- ============================================
-- Agent 管理表
-- ============================================

CREATE TABLE IF NOT EXISTS agents (
    id SERIAL PRIMARY KEY,
    agent_id VARCHAR(64) UNIQUE NOT NULL,
    mode VARCHAR(16) NOT NULL CHECK (mode IN ('external', 'internal')),
    hostname VARCHAR(256),
    ip_address INET,
    capabilities JSONB,
    status VARCHAR(32) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'offline')),
    api_key_hash VARCHAR(64),
    client_cert TEXT,
    last_heartbeat TIMESTAMPTZ,
    registered_at TIMESTAMPTZ DEFAULT NOW(),
    deregistered_at TIMESTAMPTZ,
    config JSONB,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_agents_agent_id ON agents(agent_id);
CREATE INDEX idx_agents_mode ON agents(mode);
CREATE INDEX idx_agents_status ON agents(status);
CREATE INDEX idx_agents_last_heartbeat ON agents(last_heartbeat);

-- ============================================
-- 保留策略表
-- ============================================

CREATE TABLE IF NOT EXISTS retention_policies (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(64) NOT NULL,
    agent_mode VARCHAR(16) CHECK (agent_mode IN ('external', 'internal', NULL)),
    retention_days INTEGER NOT NULL CHECK (retention_days > 0),
    legal_hold BOOLEAN DEFAULT FALSE,
    regulation VARCHAR(32) NOT NULL CHECK (regulation IN ('GDPR', 'PCI-DSS', 'HIPAA', 'SOX', 'ISO27001')),
    auto_delete BOOLEAN DEFAULT TRUE,
    archive_required BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_retention_event_type ON retention_policies(event_type);
CREATE INDEX idx_retention_regulation ON retention_policies(regulation);

-- 插入默認策略
INSERT INTO retention_policies (event_type, agent_mode, retention_days, regulation, auto_delete, archive_required) VALUES
    ('windows_event_log', 'external', 90, 'PCI-DSS', TRUE, TRUE),
    ('security_alert', 'external', 365, 'GDPR', TRUE, TRUE),
    ('access_log', 'internal', 180, 'SOX', TRUE, FALSE),
    ('compliance_scan', 'external', 2555, 'HIPAA', FALSE, TRUE)
ON CONFLICT DO NOTHING;

-- ============================================
-- GDPR 刪除請求表
-- ============================================

CREATE TABLE IF NOT EXISTS gdpr_deletion_requests (
    id SERIAL PRIMARY KEY,
    request_id UUID DEFAULT gen_random_uuid() UNIQUE,
    subject_identifier VARCHAR(256) NOT NULL,
    requested_by VARCHAR(128),
    requested_at TIMESTAMPTZ DEFAULT NOW(),
    approved_by VARCHAR(128),
    approved_at TIMESTAMPTZ,
    status VARCHAR(32) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'completed', 'rejected')),
    completion_date TIMESTAMPTZ,
    verification_hash VARCHAR(64),
    deleted_count INTEGER,
    notes TEXT
);

CREATE INDEX idx_gdpr_request_id ON gdpr_deletion_requests(request_id);
CREATE INDEX idx_gdpr_status ON gdpr_deletion_requests(status);
CREATE INDEX idx_gdpr_requested_at ON gdpr_deletion_requests(requested_at);

-- ============================================
-- 審計訪問日誌表
-- ============================================

CREATE TABLE IF NOT EXISTS audit_access_log (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    user_id VARCHAR(64) NOT NULL,
    action VARCHAR(64) NOT NULL CHECK (action IN ('query', 'export', 'delete', 'update', 'create')),
    resource_type VARCHAR(64),
    resource_id VARCHAR(256),
    query_text TEXT,
    record_count INTEGER,
    ip_address INET,
    user_agent TEXT,
    justification TEXT,
    approved_by VARCHAR(64),
    session_id VARCHAR(128)
);

CREATE INDEX idx_audit_timestamp ON audit_access_log(timestamp DESC);
CREATE INDEX idx_audit_user_id ON audit_access_log(user_id, timestamp DESC);
CREATE INDEX idx_audit_action ON audit_access_log(action, timestamp DESC);
CREATE INDEX idx_audit_session ON audit_access_log(session_id);

-- ============================================
-- PII 模式表
-- ============================================

CREATE TABLE IF NOT EXISTS pii_patterns (
    id SERIAL PRIMARY KEY,
    pattern_type VARCHAR(64) UNIQUE NOT NULL,
    regex VARCHAR(512) NOT NULL,
    description TEXT,
    severity VARCHAR(20) CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_pii_patterns_type ON pii_patterns(pattern_type);
CREATE INDEX idx_pii_patterns_enabled ON pii_patterns(enabled);

-- 插入默認 PII 模式
INSERT INTO pii_patterns (pattern_type, regex, description, severity, enabled) VALUES
    ('email', '[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}', 'Email addresses', 'medium', TRUE),
    ('credit_card', '\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b', 'Credit card numbers', 'critical', TRUE),
    ('ssn', '\b\d{3}-\d{2}-\d{4}\b', 'Social Security Numbers', 'critical', TRUE),
    ('ip_address', '\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b', 'IP addresses', 'low', TRUE),
    ('phone', '\b\d{3}[-.]?\d{3}[-.]?\d{4}\b', 'Phone numbers', 'medium', TRUE)
ON CONFLICT (pattern_type) DO NOTHING;

-- ============================================
-- PII 出現記錄表
-- ============================================

CREATE TABLE IF NOT EXISTS pii_occurrences (
    id BIGSERIAL PRIMARY KEY,
    log_id BIGINT,
    pii_type VARCHAR(64),
    field_name VARCHAR(128),
    detected_at TIMESTAMPTZ DEFAULT NOW(),
    masked BOOLEAN DEFAULT FALSE,
    hash VARCHAR(64)
);

CREATE INDEX idx_pii_log_id ON pii_occurrences(log_id);
CREATE INDEX idx_pii_type ON pii_occurrences(pii_type, detected_at DESC);
CREATE INDEX idx_pii_detected_at ON pii_occurrences(detected_at DESC);

-- ============================================
-- 封存清單表
-- ============================================

CREATE TABLE IF NOT EXISTS archive_manifests (
    id BIGSERIAL PRIMARY KEY,
    file_path VARCHAR(512) UNIQUE NOT NULL,
    date DATE NOT NULL,
    agent_id VARCHAR(64),
    record_count INTEGER,
    size_bytes BIGINT,
    compression_ratio FLOAT,
    checksum VARCHAR(64),
    encryption_key_id VARCHAR(64),
    compliance_tags JSONB,
    retention_until DATE,
    worm_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_archive_date ON archive_manifests(date DESC);
CREATE INDEX idx_archive_agent_id ON archive_manifests(agent_id);
CREATE INDEX idx_archive_retention ON archive_manifests(retention_until);

-- ============================================
-- 為現有 event_logs 表添加分區支持
-- ============================================

-- 注意：如果表已存在，這些修改需要謹慎處理
-- 建議在生產環境中使用 ALTER TABLE 逐步遷移

-- 添加缺失的列（如果不存在）
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='event_logs' AND column_name='agent_mode') THEN
        ALTER TABLE event_logs ADD COLUMN agent_mode VARCHAR(16);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='event_logs' AND column_name='retention_until') THEN
        ALTER TABLE event_logs ADD COLUMN retention_until TIMESTAMPTZ;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='event_logs' AND column_name='archived') THEN
        ALTER TABLE event_logs ADD COLUMN archived BOOLEAN DEFAULT FALSE;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='event_logs' AND column_name='integrity_hash') THEN
        ALTER TABLE event_logs ADD COLUMN integrity_hash VARCHAR(64);
    END IF;
END $$;

-- 添加索引
CREATE INDEX IF NOT EXISTS idx_event_logs_agent_mode ON event_logs(agent_mode);
CREATE INDEX IF NOT EXISTS idx_event_logs_retention ON event_logs(retention_until, archived);
CREATE INDEX IF NOT EXISTS idx_event_logs_integrity ON event_logs(integrity_hash);

-- ============================================
-- 完整性 Hash 計算函數
-- ============================================

CREATE OR REPLACE FUNCTION calculate_integrity_hash()
RETURNS TRIGGER AS $$
BEGIN
    NEW.integrity_hash := encode(
        digest(
            NEW.timestamp::text || 
            COALESCE(NEW.agent_id, '') || 
            COALESCE(NEW.event_type, '') || 
            COALESCE(NEW.message, ''),
            'sha256'
        ),
        'hex'
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 創建觸發器（如果不存在）
DROP TRIGGER IF EXISTS trg_event_logs_integrity_hash ON event_logs;
CREATE TRIGGER trg_event_logs_integrity_hash
    BEFORE INSERT ON event_logs
    FOR EACH ROW
    EXECUTE FUNCTION calculate_integrity_hash();

-- ============================================
-- 月度分區創建函數
-- ============================================

CREATE OR REPLACE FUNCTION create_monthly_partition()
RETURNS void AS $$
DECLARE
    partition_date DATE;
    partition_name TEXT;
    start_date TEXT;
    end_date TEXT;
BEGIN
    -- 創建下個月的分區
    partition_date := DATE_TRUNC('month', NOW() + INTERVAL '1 month')::DATE;
    partition_name := 'event_logs_' || TO_CHAR(partition_date, 'YYYY_MM');
    start_date := partition_date::TEXT;
    end_date := (partition_date + INTERVAL '1 month')::DATE::TEXT;
    
    -- 檢查分區是否已存在
    IF NOT EXISTS (
        SELECT 1 FROM pg_class 
        WHERE relname = partition_name
    ) THEN
        EXECUTE format('
            CREATE TABLE %I PARTITION OF event_logs
            FOR VALUES FROM (%L) TO (%L)',
            partition_name,
            start_date,
            end_date
        );
        
        RAISE NOTICE 'Created partition: %', partition_name;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- 創建當前和下個月的分區
SELECT create_monthly_partition();

-- ============================================
-- 自動化任務（需要 pg_cron 擴展）
-- ============================================

-- 注意：需要先安裝 pg_cron 擴展
-- CREATE EXTENSION IF NOT EXISTS pg_cron;

-- 每天創建分區（如果 pg_cron 可用）
-- SELECT cron.schedule('create-monthly-partition', '0 0 * * *', 'SELECT create_monthly_partition()');

-- ============================================
-- 保留策略執行函數
-- ============================================

CREATE OR REPLACE FUNCTION enforce_retention_policies()
RETURNS TABLE(deleted_count INTEGER, event_type TEXT) AS $$
DECLARE
    policy RECORD;
    cutoff_date TIMESTAMPTZ;
    rows_deleted INTEGER;
BEGIN
    FOR policy IN SELECT * FROM retention_policies WHERE auto_delete = TRUE LOOP
        cutoff_date := NOW() - (policy.retention_days || ' days')::INTERVAL;
        
        -- 刪除過期數據（未設置 legal_hold）
        DELETE FROM event_logs
        WHERE event_logs.event_type = policy.event_type
          AND event_logs.timestamp < cutoff_date
          AND event_logs.archived = FALSE
          AND NOT EXISTS (
              SELECT 1 FROM retention_policies rp
              WHERE rp.event_type = event_logs.event_type
                AND rp.legal_hold = TRUE
          );
        
        GET DIAGNOSTICS rows_deleted = ROW_COUNT;
        
        deleted_count := rows_deleted;
        event_type := policy.event_type;
        RETURN NEXT;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- 數據統計視圖
-- ============================================

CREATE OR REPLACE VIEW storage_stats AS
SELECT
    'hot' AS tier,
    NULL AS event_count,
    '1 hour' AS retention,
    'Redis Streams' AS technology
UNION ALL
SELECT
    'warm' AS tier,
    NULL AS event_count,
    '7 days' AS retention,
    'Loki' AS technology
UNION ALL
SELECT
    'cold' AS tier,
    COUNT(*) AS event_count,
    '90 days' AS retention,
    'PostgreSQL' AS technology
FROM event_logs
WHERE archived = FALSE
UNION ALL
SELECT
    'archive' AS tier,
    COUNT(*) AS event_count,
    '7+ years' AS retention,
    'S3/MinIO WORM' AS technology
FROM event_logs
WHERE archived = TRUE;

-- ============================================
-- 合規性報告視圖
-- ============================================

CREATE OR REPLACE VIEW compliance_summary AS
SELECT
    rp.regulation,
    rp.event_type,
    COUNT(el.id) AS total_logs,
    SUM(CASE WHEN el.archived THEN 1 ELSE 0 END) AS archived_logs,
    SUM(CASE WHEN el.timestamp < NOW() - (rp.retention_days || ' days')::INTERVAL THEN 1 ELSE 0 END) AS overdue_logs
FROM retention_policies rp
LEFT JOIN event_logs el ON el.event_type = rp.event_type
GROUP BY rp.regulation, rp.event_type;

-- ============================================
-- PII 統計視圖
-- ============================================

CREATE OR REPLACE VIEW pii_detection_stats AS
SELECT
    pii_type,
    COUNT(*) AS occurrences,
    SUM(CASE WHEN masked THEN 1 ELSE 0 END) AS masked_count,
    SUM(CASE WHEN NOT masked THEN 1 ELSE 0 END) AS unmasked_count,
    DATE_TRUNC('day', detected_at) AS detection_date
FROM pii_occurrences
GROUP BY pii_type, DATE_TRUNC('day', detected_at)
ORDER BY detection_date DESC;

-- ============================================
-- 審計報告視圖
-- ============================================

CREATE OR REPLACE VIEW audit_summary AS
SELECT
    user_id,
    action,
    COUNT(*) AS action_count,
    SUM(record_count) AS total_records_accessed,
    MAX(timestamp) AS last_access,
    DATE_TRUNC('day', timestamp) AS access_date
FROM audit_access_log
GROUP BY user_id, action, DATE_TRUNC('day', timestamp)
ORDER BY access_date DESC;

-- ============================================
-- 註釋
-- ============================================

COMMENT ON TABLE agents IS 'Agent 註冊和管理表，支援雙模式連接';
COMMENT ON TABLE retention_policies IS '多法規資料保留策略表';
COMMENT ON TABLE gdpr_deletion_requests IS 'GDPR 刪除請求追蹤表';
COMMENT ON TABLE audit_access_log IS '所有資料訪問的審計日誌';
COMMENT ON TABLE pii_patterns IS 'PII 檢測模式庫';
COMMENT ON TABLE pii_occurrences IS 'PII 出現記錄和追蹤';
COMMENT ON TABLE archive_manifests IS 'S3/MinIO 封存檔案清單';

COMMENT ON COLUMN agents.mode IS 'external: 通過 Nginx + mTLS, internal: 直連';
COMMENT ON COLUMN retention_policies.legal_hold IS 'Legal Hold: 禁止自動刪除';
COMMENT ON COLUMN event_logs.integrity_hash IS 'SHA-256 完整性哈希值，用於防篡改';

-- ============================================
-- 完成
-- ============================================

-- 記錄 Migration 版本
INSERT INTO schema_migrations (version, description, applied_at) VALUES
    ('002', 'Agent and Compliance Schema', NOW())
ON CONFLICT (version) DO NOTHING;

