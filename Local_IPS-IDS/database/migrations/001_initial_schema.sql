-- Axiom Backend V3 Initial Schema
-- Version: 3.0.0
-- Date: 2025-10-16

-- 服務狀態表
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(50),
    status VARCHAR(20),
    url VARCHAR(255),
    version VARCHAR(50),
    last_check TIMESTAMP WITH TIME ZONE,
    config JSONB,
    metrics JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_services_name ON services(name);
CREATE INDEX idx_services_type ON services(type);
CREATE INDEX idx_services_status ON services(status);
CREATE INDEX idx_services_last_check ON services(last_check);

-- 配置歷史表
CREATE TABLE IF NOT EXISTS config_histories (
    id SERIAL PRIMARY KEY,
    service_id INTEGER NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    config_type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    applied_by VARCHAR(100),
    applied_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL,
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_config_histories_service_id ON config_histories(service_id);
CREATE INDEX idx_config_histories_applied_at ON config_histories(applied_at);
CREATE INDEX idx_config_histories_status ON config_histories(status);

-- 量子作業表
CREATE TABLE IF NOT EXISTS quantum_jobs (
    id SERIAL PRIMARY KEY,
    job_id VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    backend VARCHAR(50),
    backend_name VARCHAR(100),
    circuit TEXT,
    input_data JSONB,
    result JSONB,
    shots INTEGER DEFAULT 1024,
    qubits INTEGER,
    depth INTEGER,
    estimated_time INTEGER,
    actual_time INTEGER,
    submitted_by VARCHAR(100),
    submitted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    error TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_quantum_jobs_job_id ON quantum_jobs(job_id);
CREATE INDEX idx_quantum_jobs_type ON quantum_jobs(type);
CREATE INDEX idx_quantum_jobs_status ON quantum_jobs(status);
CREATE INDEX idx_quantum_jobs_submitted_at ON quantum_jobs(submitted_at);
CREATE INDEX idx_quantum_jobs_completed_at ON quantum_jobs(completed_at);
CREATE INDEX idx_quantum_jobs_type_status ON quantum_jobs(type, status);

-- Windows 日誌表
CREATE TABLE IF NOT EXISTS windows_logs (
    id SERIAL PRIMARY KEY,
    agent_id VARCHAR(100) NOT NULL,
    log_type VARCHAR(50) NOT NULL,
    source VARCHAR(255),
    event_id INTEGER,
    level VARCHAR(20),
    message TEXT,
    time_created TIMESTAMP WITH TIME ZONE NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    computer VARCHAR(255),
    user_id VARCHAR(100),
    process_id INTEGER,
    thread_id INTEGER,
    keywords VARCHAR(255),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_windows_logs_agent_id ON windows_logs(agent_id);
CREATE INDEX idx_windows_logs_log_type ON windows_logs(log_type);
CREATE INDEX idx_windows_logs_level ON windows_logs(level);
CREATE INDEX idx_windows_logs_event_id ON windows_logs(event_id);
CREATE INDEX idx_windows_logs_time_created ON windows_logs(time_created DESC);
CREATE INDEX idx_windows_logs_received_at ON windows_logs(received_at DESC);
CREATE INDEX idx_windows_logs_source ON windows_logs(source);

-- 告警表
CREATE TABLE IF NOT EXISTS alerts (
    id SERIAL PRIMARY KEY,
    alert_name VARCHAR(255) NOT NULL,
    fingerprint VARCHAR(64) UNIQUE,
    severity VARCHAR(20) NOT NULL,
    source VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    message TEXT NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL,
    priority INTEGER DEFAULT 0,
    count INTEGER DEFAULT 1,
    labels JSONB,
    annotations JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolved_by VARCHAR(100),
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    acknowledged_by VARCHAR(100),
    last_occurred_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_alerts_alert_name ON alerts(alert_name);
CREATE INDEX idx_alerts_severity ON alerts(severity);
CREATE INDEX idx_alerts_source ON alerts(source);
CREATE INDEX idx_alerts_category ON alerts(category);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_created_at ON alerts(created_at DESC);
CREATE INDEX idx_alerts_resolved_at ON alerts(resolved_at);

-- API 請求日誌表
CREATE TABLE IF NOT EXISTS api_logs (
    id SERIAL PRIMARY KEY,
    method VARCHAR(10) NOT NULL,
    path VARCHAR(500) NOT NULL,
    status INTEGER NOT NULL,
    duration BIGINT NOT NULL,
    client_ip VARCHAR(45),
    user_agent VARCHAR(500),
    request_body TEXT,
    response_body TEXT,
    error TEXT,
    user_id VARCHAR(100),
    api_key VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_api_logs_method ON api_logs(method);
CREATE INDEX idx_api_logs_path ON api_logs(path);
CREATE INDEX idx_api_logs_status ON api_logs(status);
CREATE INDEX idx_api_logs_client_ip ON api_logs(client_ip);
CREATE INDEX idx_api_logs_created_at ON api_logs(created_at DESC);
CREATE INDEX idx_api_logs_user_id ON api_logs(user_id);
CREATE INDEX idx_api_logs_created_at_status ON api_logs(created_at, status);

-- 指標快照表
CREATE TABLE IF NOT EXISTS metric_snapshots (
    id SERIAL PRIMARY KEY,
    service_id INTEGER NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    metric_name VARCHAR(255) NOT NULL,
    metric_type VARCHAR(50),
    value DOUBLE PRECISION NOT NULL,
    unit VARCHAR(20),
    labels JSONB,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_metric_snapshots_service_id ON metric_snapshots(service_id);
CREATE INDEX idx_metric_snapshots_metric_name ON metric_snapshots(metric_name);
CREATE INDEX idx_metric_snapshots_timestamp ON metric_snapshots(timestamp DESC);
CREATE INDEX idx_metric_snapshots_service_metric_time ON metric_snapshots(service_id, metric_name, timestamp);

-- 使用者表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'viewer',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    api_key VARCHAR(64) UNIQUE,
    permissions JSONB,
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_login_ip VARCHAR(45),
    login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_api_key ON users(api_key);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- 會話表
CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(128) UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(512) UNIQUE NOT NULL,
    client_ip VARCHAR(45),
    user_agent VARCHAR(500),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_session_id ON sessions(session_id);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- 創建默認管理員用戶
-- 密碼: admin123 (請在生產環境中修改)
INSERT INTO users (username, email, password_hash, role, status)
VALUES (
    'admin',
    'admin@pandora.local',
    '$2a$10$YourBcryptHashHere', -- 需要生成實際的 bcrypt hash
    'admin',
    'active'
) ON CONFLICT (username) DO NOTHING;

-- 完成
SELECT 'Migration 001 completed successfully' AS message;


-- Version: 3.0.0
-- Date: 2025-10-16

-- 服務狀態表
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(50),
    status VARCHAR(20),
    url VARCHAR(255),
    version VARCHAR(50),
    last_check TIMESTAMP WITH TIME ZONE,
    config JSONB,
    metrics JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_services_name ON services(name);
CREATE INDEX idx_services_type ON services(type);
CREATE INDEX idx_services_status ON services(status);
CREATE INDEX idx_services_last_check ON services(last_check);

-- 配置歷史表
CREATE TABLE IF NOT EXISTS config_histories (
    id SERIAL PRIMARY KEY,
    service_id INTEGER NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    config_type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    applied_by VARCHAR(100),
    applied_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL,
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_config_histories_service_id ON config_histories(service_id);
CREATE INDEX idx_config_histories_applied_at ON config_histories(applied_at);
CREATE INDEX idx_config_histories_status ON config_histories(status);

-- 量子作業表
CREATE TABLE IF NOT EXISTS quantum_jobs (
    id SERIAL PRIMARY KEY,
    job_id VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    backend VARCHAR(50),
    backend_name VARCHAR(100),
    circuit TEXT,
    input_data JSONB,
    result JSONB,
    shots INTEGER DEFAULT 1024,
    qubits INTEGER,
    depth INTEGER,
    estimated_time INTEGER,
    actual_time INTEGER,
    submitted_by VARCHAR(100),
    submitted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    error TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_quantum_jobs_job_id ON quantum_jobs(job_id);
CREATE INDEX idx_quantum_jobs_type ON quantum_jobs(type);
CREATE INDEX idx_quantum_jobs_status ON quantum_jobs(status);
CREATE INDEX idx_quantum_jobs_submitted_at ON quantum_jobs(submitted_at);
CREATE INDEX idx_quantum_jobs_completed_at ON quantum_jobs(completed_at);
CREATE INDEX idx_quantum_jobs_type_status ON quantum_jobs(type, status);

-- Windows 日誌表
CREATE TABLE IF NOT EXISTS windows_logs (
    id SERIAL PRIMARY KEY,
    agent_id VARCHAR(100) NOT NULL,
    log_type VARCHAR(50) NOT NULL,
    source VARCHAR(255),
    event_id INTEGER,
    level VARCHAR(20),
    message TEXT,
    time_created TIMESTAMP WITH TIME ZONE NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    computer VARCHAR(255),
    user_id VARCHAR(100),
    process_id INTEGER,
    thread_id INTEGER,
    keywords VARCHAR(255),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_windows_logs_agent_id ON windows_logs(agent_id);
CREATE INDEX idx_windows_logs_log_type ON windows_logs(log_type);
CREATE INDEX idx_windows_logs_level ON windows_logs(level);
CREATE INDEX idx_windows_logs_event_id ON windows_logs(event_id);
CREATE INDEX idx_windows_logs_time_created ON windows_logs(time_created DESC);
CREATE INDEX idx_windows_logs_received_at ON windows_logs(received_at DESC);
CREATE INDEX idx_windows_logs_source ON windows_logs(source);

-- 告警表
CREATE TABLE IF NOT EXISTS alerts (
    id SERIAL PRIMARY KEY,
    alert_name VARCHAR(255) NOT NULL,
    fingerprint VARCHAR(64) UNIQUE,
    severity VARCHAR(20) NOT NULL,
    source VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    message TEXT NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL,
    priority INTEGER DEFAULT 0,
    count INTEGER DEFAULT 1,
    labels JSONB,
    annotations JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolved_by VARCHAR(100),
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    acknowledged_by VARCHAR(100),
    last_occurred_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_alerts_alert_name ON alerts(alert_name);
CREATE INDEX idx_alerts_severity ON alerts(severity);
CREATE INDEX idx_alerts_source ON alerts(source);
CREATE INDEX idx_alerts_category ON alerts(category);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_created_at ON alerts(created_at DESC);
CREATE INDEX idx_alerts_resolved_at ON alerts(resolved_at);

-- API 請求日誌表
CREATE TABLE IF NOT EXISTS api_logs (
    id SERIAL PRIMARY KEY,
    method VARCHAR(10) NOT NULL,
    path VARCHAR(500) NOT NULL,
    status INTEGER NOT NULL,
    duration BIGINT NOT NULL,
    client_ip VARCHAR(45),
    user_agent VARCHAR(500),
    request_body TEXT,
    response_body TEXT,
    error TEXT,
    user_id VARCHAR(100),
    api_key VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_api_logs_method ON api_logs(method);
CREATE INDEX idx_api_logs_path ON api_logs(path);
CREATE INDEX idx_api_logs_status ON api_logs(status);
CREATE INDEX idx_api_logs_client_ip ON api_logs(client_ip);
CREATE INDEX idx_api_logs_created_at ON api_logs(created_at DESC);
CREATE INDEX idx_api_logs_user_id ON api_logs(user_id);
CREATE INDEX idx_api_logs_created_at_status ON api_logs(created_at, status);

-- 指標快照表
CREATE TABLE IF NOT EXISTS metric_snapshots (
    id SERIAL PRIMARY KEY,
    service_id INTEGER NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    metric_name VARCHAR(255) NOT NULL,
    metric_type VARCHAR(50),
    value DOUBLE PRECISION NOT NULL,
    unit VARCHAR(20),
    labels JSONB,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_metric_snapshots_service_id ON metric_snapshots(service_id);
CREATE INDEX idx_metric_snapshots_metric_name ON metric_snapshots(metric_name);
CREATE INDEX idx_metric_snapshots_timestamp ON metric_snapshots(timestamp DESC);
CREATE INDEX idx_metric_snapshots_service_metric_time ON metric_snapshots(service_id, metric_name, timestamp);

-- 使用者表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'viewer',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    api_key VARCHAR(64) UNIQUE,
    permissions JSONB,
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_login_ip VARCHAR(45),
    login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_api_key ON users(api_key);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- 會話表
CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(128) UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(512) UNIQUE NOT NULL,
    client_ip VARCHAR(45),
    user_agent VARCHAR(500),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_session_id ON sessions(session_id);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- 創建默認管理員用戶
-- 密碼: admin123 (請在生產環境中修改)
INSERT INTO users (username, email, password_hash, role, status)
VALUES (
    'admin',
    'admin@pandora.local',
    '$2a$10$YourBcryptHashHere', -- 需要生成實際的 bcrypt hash
    'admin',
    'active'
) ON CONFLICT (username) DO NOTHING;

-- 完成
SELECT 'Migration 001 completed successfully' AS message;

-- Version: 3.0.0
-- Date: 2025-10-16

-- 服務狀態表
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(50),
    status VARCHAR(20),
    url VARCHAR(255),
    version VARCHAR(50),
    last_check TIMESTAMP WITH TIME ZONE,
    config JSONB,
    metrics JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_services_name ON services(name);
CREATE INDEX idx_services_type ON services(type);
CREATE INDEX idx_services_status ON services(status);
CREATE INDEX idx_services_last_check ON services(last_check);

-- 配置歷史表
CREATE TABLE IF NOT EXISTS config_histories (
    id SERIAL PRIMARY KEY,
    service_id INTEGER NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    config_type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    applied_by VARCHAR(100),
    applied_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL,
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_config_histories_service_id ON config_histories(service_id);
CREATE INDEX idx_config_histories_applied_at ON config_histories(applied_at);
CREATE INDEX idx_config_histories_status ON config_histories(status);

-- 量子作業表
CREATE TABLE IF NOT EXISTS quantum_jobs (
    id SERIAL PRIMARY KEY,
    job_id VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    backend VARCHAR(50),
    backend_name VARCHAR(100),
    circuit TEXT,
    input_data JSONB,
    result JSONB,
    shots INTEGER DEFAULT 1024,
    qubits INTEGER,
    depth INTEGER,
    estimated_time INTEGER,
    actual_time INTEGER,
    submitted_by VARCHAR(100),
    submitted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    error TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_quantum_jobs_job_id ON quantum_jobs(job_id);
CREATE INDEX idx_quantum_jobs_type ON quantum_jobs(type);
CREATE INDEX idx_quantum_jobs_status ON quantum_jobs(status);
CREATE INDEX idx_quantum_jobs_submitted_at ON quantum_jobs(submitted_at);
CREATE INDEX idx_quantum_jobs_completed_at ON quantum_jobs(completed_at);
CREATE INDEX idx_quantum_jobs_type_status ON quantum_jobs(type, status);

-- Windows 日誌表
CREATE TABLE IF NOT EXISTS windows_logs (
    id SERIAL PRIMARY KEY,
    agent_id VARCHAR(100) NOT NULL,
    log_type VARCHAR(50) NOT NULL,
    source VARCHAR(255),
    event_id INTEGER,
    level VARCHAR(20),
    message TEXT,
    time_created TIMESTAMP WITH TIME ZONE NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    computer VARCHAR(255),
    user_id VARCHAR(100),
    process_id INTEGER,
    thread_id INTEGER,
    keywords VARCHAR(255),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_windows_logs_agent_id ON windows_logs(agent_id);
CREATE INDEX idx_windows_logs_log_type ON windows_logs(log_type);
CREATE INDEX idx_windows_logs_level ON windows_logs(level);
CREATE INDEX idx_windows_logs_event_id ON windows_logs(event_id);
CREATE INDEX idx_windows_logs_time_created ON windows_logs(time_created DESC);
CREATE INDEX idx_windows_logs_received_at ON windows_logs(received_at DESC);
CREATE INDEX idx_windows_logs_source ON windows_logs(source);

-- 告警表
CREATE TABLE IF NOT EXISTS alerts (
    id SERIAL PRIMARY KEY,
    alert_name VARCHAR(255) NOT NULL,
    fingerprint VARCHAR(64) UNIQUE,
    severity VARCHAR(20) NOT NULL,
    source VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    message TEXT NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL,
    priority INTEGER DEFAULT 0,
    count INTEGER DEFAULT 1,
    labels JSONB,
    annotations JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolved_by VARCHAR(100),
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    acknowledged_by VARCHAR(100),
    last_occurred_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_alerts_alert_name ON alerts(alert_name);
CREATE INDEX idx_alerts_severity ON alerts(severity);
CREATE INDEX idx_alerts_source ON alerts(source);
CREATE INDEX idx_alerts_category ON alerts(category);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_created_at ON alerts(created_at DESC);
CREATE INDEX idx_alerts_resolved_at ON alerts(resolved_at);

-- API 請求日誌表
CREATE TABLE IF NOT EXISTS api_logs (
    id SERIAL PRIMARY KEY,
    method VARCHAR(10) NOT NULL,
    path VARCHAR(500) NOT NULL,
    status INTEGER NOT NULL,
    duration BIGINT NOT NULL,
    client_ip VARCHAR(45),
    user_agent VARCHAR(500),
    request_body TEXT,
    response_body TEXT,
    error TEXT,
    user_id VARCHAR(100),
    api_key VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_api_logs_method ON api_logs(method);
CREATE INDEX idx_api_logs_path ON api_logs(path);
CREATE INDEX idx_api_logs_status ON api_logs(status);
CREATE INDEX idx_api_logs_client_ip ON api_logs(client_ip);
CREATE INDEX idx_api_logs_created_at ON api_logs(created_at DESC);
CREATE INDEX idx_api_logs_user_id ON api_logs(user_id);
CREATE INDEX idx_api_logs_created_at_status ON api_logs(created_at, status);

-- 指標快照表
CREATE TABLE IF NOT EXISTS metric_snapshots (
    id SERIAL PRIMARY KEY,
    service_id INTEGER NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    metric_name VARCHAR(255) NOT NULL,
    metric_type VARCHAR(50),
    value DOUBLE PRECISION NOT NULL,
    unit VARCHAR(20),
    labels JSONB,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_metric_snapshots_service_id ON metric_snapshots(service_id);
CREATE INDEX idx_metric_snapshots_metric_name ON metric_snapshots(metric_name);
CREATE INDEX idx_metric_snapshots_timestamp ON metric_snapshots(timestamp DESC);
CREATE INDEX idx_metric_snapshots_service_metric_time ON metric_snapshots(service_id, metric_name, timestamp);

-- 使用者表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'viewer',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    api_key VARCHAR(64) UNIQUE,
    permissions JSONB,
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_login_ip VARCHAR(45),
    login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_api_key ON users(api_key);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- 會話表
CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(128) UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(512) UNIQUE NOT NULL,
    client_ip VARCHAR(45),
    user_agent VARCHAR(500),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_session_id ON sessions(session_id);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- 創建默認管理員用戶
-- 密碼: admin123 (請在生產環境中修改)
INSERT INTO users (username, email, password_hash, role, status)
VALUES (
    'admin',
    'admin@pandora.local',
    '$2a$10$YourBcryptHashHere', -- 需要生成實際的 bcrypt hash
    'admin',
    'active'
) ON CONFLICT (username) DO NOTHING;

-- 完成
SELECT 'Migration 001 completed successfully' AS message;


-- Version: 3.0.0
-- Date: 2025-10-16

-- 服務狀態表
CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(50),
    status VARCHAR(20),
    url VARCHAR(255),
    version VARCHAR(50),
    last_check TIMESTAMP WITH TIME ZONE,
    config JSONB,
    metrics JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_services_name ON services(name);
CREATE INDEX idx_services_type ON services(type);
CREATE INDEX idx_services_status ON services(status);
CREATE INDEX idx_services_last_check ON services(last_check);

-- 配置歷史表
CREATE TABLE IF NOT EXISTS config_histories (
    id SERIAL PRIMARY KEY,
    service_id INTEGER NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    config_type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    applied_by VARCHAR(100),
    applied_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL,
    error TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_config_histories_service_id ON config_histories(service_id);
CREATE INDEX idx_config_histories_applied_at ON config_histories(applied_at);
CREATE INDEX idx_config_histories_status ON config_histories(status);

-- 量子作業表
CREATE TABLE IF NOT EXISTS quantum_jobs (
    id SERIAL PRIMARY KEY,
    job_id VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    backend VARCHAR(50),
    backend_name VARCHAR(100),
    circuit TEXT,
    input_data JSONB,
    result JSONB,
    shots INTEGER DEFAULT 1024,
    qubits INTEGER,
    depth INTEGER,
    estimated_time INTEGER,
    actual_time INTEGER,
    submitted_by VARCHAR(100),
    submitted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    error TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_quantum_jobs_job_id ON quantum_jobs(job_id);
CREATE INDEX idx_quantum_jobs_type ON quantum_jobs(type);
CREATE INDEX idx_quantum_jobs_status ON quantum_jobs(status);
CREATE INDEX idx_quantum_jobs_submitted_at ON quantum_jobs(submitted_at);
CREATE INDEX idx_quantum_jobs_completed_at ON quantum_jobs(completed_at);
CREATE INDEX idx_quantum_jobs_type_status ON quantum_jobs(type, status);

-- Windows 日誌表
CREATE TABLE IF NOT EXISTS windows_logs (
    id SERIAL PRIMARY KEY,
    agent_id VARCHAR(100) NOT NULL,
    log_type VARCHAR(50) NOT NULL,
    source VARCHAR(255),
    event_id INTEGER,
    level VARCHAR(20),
    message TEXT,
    time_created TIMESTAMP WITH TIME ZONE NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE NOT NULL,
    computer VARCHAR(255),
    user_id VARCHAR(100),
    process_id INTEGER,
    thread_id INTEGER,
    keywords VARCHAR(255),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_windows_logs_agent_id ON windows_logs(agent_id);
CREATE INDEX idx_windows_logs_log_type ON windows_logs(log_type);
CREATE INDEX idx_windows_logs_level ON windows_logs(level);
CREATE INDEX idx_windows_logs_event_id ON windows_logs(event_id);
CREATE INDEX idx_windows_logs_time_created ON windows_logs(time_created DESC);
CREATE INDEX idx_windows_logs_received_at ON windows_logs(received_at DESC);
CREATE INDEX idx_windows_logs_source ON windows_logs(source);

-- 告警表
CREATE TABLE IF NOT EXISTS alerts (
    id SERIAL PRIMARY KEY,
    alert_name VARCHAR(255) NOT NULL,
    fingerprint VARCHAR(64) UNIQUE,
    severity VARCHAR(20) NOT NULL,
    source VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    message TEXT NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL,
    priority INTEGER DEFAULT 0,
    count INTEGER DEFAULT 1,
    labels JSONB,
    annotations JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolved_by VARCHAR(100),
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    acknowledged_by VARCHAR(100),
    last_occurred_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_alerts_alert_name ON alerts(alert_name);
CREATE INDEX idx_alerts_severity ON alerts(severity);
CREATE INDEX idx_alerts_source ON alerts(source);
CREATE INDEX idx_alerts_category ON alerts(category);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_created_at ON alerts(created_at DESC);
CREATE INDEX idx_alerts_resolved_at ON alerts(resolved_at);

-- API 請求日誌表
CREATE TABLE IF NOT EXISTS api_logs (
    id SERIAL PRIMARY KEY,
    method VARCHAR(10) NOT NULL,
    path VARCHAR(500) NOT NULL,
    status INTEGER NOT NULL,
    duration BIGINT NOT NULL,
    client_ip VARCHAR(45),
    user_agent VARCHAR(500),
    request_body TEXT,
    response_body TEXT,
    error TEXT,
    user_id VARCHAR(100),
    api_key VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_api_logs_method ON api_logs(method);
CREATE INDEX idx_api_logs_path ON api_logs(path);
CREATE INDEX idx_api_logs_status ON api_logs(status);
CREATE INDEX idx_api_logs_client_ip ON api_logs(client_ip);
CREATE INDEX idx_api_logs_created_at ON api_logs(created_at DESC);
CREATE INDEX idx_api_logs_user_id ON api_logs(user_id);
CREATE INDEX idx_api_logs_created_at_status ON api_logs(created_at, status);

-- 指標快照表
CREATE TABLE IF NOT EXISTS metric_snapshots (
    id SERIAL PRIMARY KEY,
    service_id INTEGER NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    metric_name VARCHAR(255) NOT NULL,
    metric_type VARCHAR(50),
    value DOUBLE PRECISION NOT NULL,
    unit VARCHAR(20),
    labels JSONB,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_metric_snapshots_service_id ON metric_snapshots(service_id);
CREATE INDEX idx_metric_snapshots_metric_name ON metric_snapshots(metric_name);
CREATE INDEX idx_metric_snapshots_timestamp ON metric_snapshots(timestamp DESC);
CREATE INDEX idx_metric_snapshots_service_metric_time ON metric_snapshots(service_id, metric_name, timestamp);

-- 使用者表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'viewer',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    api_key VARCHAR(64) UNIQUE,
    permissions JSONB,
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_login_ip VARCHAR(45),
    login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_api_key ON users(api_key);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- 會話表
CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(128) UNIQUE NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(512) UNIQUE NOT NULL,
    client_ip VARCHAR(45),
    user_agent VARCHAR(500),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_session_id ON sessions(session_id);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- 創建默認管理員用戶
-- 密碼: admin123 (請在生產環境中修改)
INSERT INTO users (username, email, password_hash, role, status)
VALUES (
    'admin',
    'admin@pandora.local',
    '$2a$10$YourBcryptHashHere', -- 需要生成實際的 bcrypt hash
    'admin',
    'active'
) ON CONFLICT (username) DO NOTHING;

-- 完成
SELECT 'Migration 001 completed successfully' AS message;

