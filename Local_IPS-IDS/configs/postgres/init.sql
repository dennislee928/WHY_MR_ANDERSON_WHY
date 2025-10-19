-- Pandora Box Console IDS-IPS 資料庫初始化腳本
-- 此腳本會在 PostgreSQL 首次啟動時自動執行

-- 建立資料庫 (Railway 通常已建立)
-- CREATE DATABASE pandora;

-- 連接到 pandora 資料庫
\c pandora;

-- 建立擴充功能
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 建立使用者表
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE
);

-- 建立設備表
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id VARCHAR(255) UNIQUE NOT NULL,
    device_name VARCHAR(255) NOT NULL,
    device_type VARCHAR(100) NOT NULL,
    port VARCHAR(50),
    status VARCHAR(50) NOT NULL DEFAULT 'offline',
    last_seen TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 建立安全事件表
CREATE TABLE IF NOT EXISTS security_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    source_ip INET,
    destination_ip INET,
    source_port INTEGER,
    destination_port INTEGER,
    protocol VARCHAR(50),
    description TEXT,
    payload JSONB,
    device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    resolved BOOLEAN DEFAULT FALSE,
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolved_by UUID REFERENCES users(id) ON DELETE SET NULL
);

-- 建立網路阻斷記錄表
CREATE TABLE IF NOT EXISTS network_blocks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ip_address INET NOT NULL,
    reason TEXT NOT NULL,
    block_type VARCHAR(50) NOT NULL,
    blocked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    blocked_by UUID REFERENCES users(id) ON DELETE SET NULL,
    unblocked_at TIMESTAMP WITH TIME ZONE,
    unblocked_by UUID REFERENCES users(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT TRUE,
    auto_unblock_at TIMESTAMP WITH TIME ZONE
);

-- 建立系統日誌表
CREATE TABLE IF NOT EXISTS system_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    level VARCHAR(50) NOT NULL,
    service VARCHAR(100) NOT NULL,
    message TEXT NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 建立 PIN 碼驗證記錄表
CREATE TABLE IF NOT EXISTS pin_verifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    pin_hash TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    verified_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 建立 Token 認證記錄表
CREATE TABLE IF NOT EXISTS token_auth (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    device_info TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    last_used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 建立告警規則表
CREATE TABLE IF NOT EXISTS alert_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rule_name VARCHAR(255) UNIQUE NOT NULL,
    rule_type VARCHAR(100) NOT NULL,
    condition TEXT NOT NULL,
    threshold NUMERIC,
    severity VARCHAR(50) NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    notification_channels JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 建立告警歷史表
CREATE TABLE IF NOT EXISTS alert_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rule_id UUID REFERENCES alert_rules(id) ON DELETE SET NULL,
    alert_type VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    metadata JSONB,
    triggered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    acknowledged BOOLEAN DEFAULT FALSE,
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    acknowledged_by UUID REFERENCES users(id) ON DELETE SET NULL,
    resolved BOOLEAN DEFAULT FALSE,
    resolved_at TIMESTAMP WITH TIME ZONE
);

-- 建立索引以提升查詢效能
CREATE INDEX IF NOT EXISTS idx_security_events_created_at ON security_events(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_security_events_severity ON security_events(severity);
CREATE INDEX IF NOT EXISTS idx_security_events_source_ip ON security_events(source_ip);
CREATE INDEX IF NOT EXISTS idx_network_blocks_ip_address ON network_blocks(ip_address);
CREATE INDEX IF NOT EXISTS idx_network_blocks_is_active ON network_blocks(is_active);
CREATE INDEX IF NOT EXISTS idx_system_logs_created_at ON system_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_system_logs_level ON system_logs(level);
CREATE INDEX IF NOT EXISTS idx_alert_history_triggered_at ON alert_history(triggered_at DESC);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 建立預設管理員帳號 (密碼: pandora123)
-- 使用 pgcrypto 的 crypt 函數加密密碼
INSERT INTO users (username, email, password_hash, role)
VALUES (
    'admin',
    'admin@pandora-ids.local',
    crypt('pandora123', gen_salt('bf')),
    'admin'
) ON CONFLICT (username) DO NOTHING;

-- 建立更新時間觸發器函數
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 為需要自動更新 updated_at 的表建立觸發器
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_devices_updated_at
    BEFORE UPDATE ON devices
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_alert_rules_updated_at
    BEFORE UPDATE ON alert_rules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 建立自動清理舊日誌的函數
CREATE OR REPLACE FUNCTION cleanup_old_logs()
RETURNS void AS $$
BEGIN
    -- 刪除 30 天前的系統日誌
    DELETE FROM system_logs WHERE created_at < NOW() - INTERVAL '30 days';
    
    -- 刪除 90 天前的已解決安全事件
    DELETE FROM security_events 
    WHERE resolved = TRUE 
    AND resolved_at < NOW() - INTERVAL '90 days';
    
    -- 刪除 90 天前的告警歷史
    DELETE FROM alert_history 
    WHERE resolved = TRUE 
    AND resolved_at < NOW() - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;

-- 授權 (Railway 會自動處理權限)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO pandora;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO pandora;

-- 完成訊息
DO $$
BEGIN
    RAISE NOTICE '✅ Pandora Box Console IDS-IPS 資料庫初始化完成!';
    RAISE NOTICE '📊 已建立 10 個資料表';
    RAISE NOTICE '🔐 預設管理員帳號: admin / pandora123';
    RAISE NOTICE '⚡ 已建立效能索引';
END $$;

