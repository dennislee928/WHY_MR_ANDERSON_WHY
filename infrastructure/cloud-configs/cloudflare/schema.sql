-- D1 Database Schema for Security Platform
-- 
-- Create this database with:
-- wrangler d1 create security_platform_db
-- wrangler d1 execute security_platform_db --file=schema.sql

-- Threats table
CREATE TABLE IF NOT EXISTS threats (
    id TEXT PRIMARY KEY,
    source_ip TEXT NOT NULL,
    threat_type TEXT NOT NULL,
    severity TEXT NOT NULL,
    status TEXT DEFAULT 'active',
    description TEXT,
    discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_threats_discovered ON threats(discovered_at);
CREATE INDEX idx_threats_status ON threats(status);
CREATE INDEX idx_threats_severity ON threats(severity);

-- Blocked IPs table
CREATE TABLE IF NOT EXISTS blocked_ips (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip TEXT NOT NULL UNIQUE,
    reason TEXT,
    blocked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME
);

CREATE INDEX idx_blocked_ips_ip ON blocked_ips(ip);

-- Devices table
CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    ip_address TEXT,
    mac_address TEXT,
    status TEXT DEFAULT 'online',
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_last_seen ON devices(last_seen);

-- Network statistics table
CREATE TABLE IF NOT EXISTS network_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    connections INTEGER DEFAULT 0,
    bytes_in INTEGER DEFAULT 0,
    bytes_out INTEGER DEFAULT 0,
    packets_dropped INTEGER DEFAULT 0
);

CREATE INDEX idx_network_stats_timestamp ON network_stats(timestamp);

-- ML detections table
CREATE TABLE IF NOT EXISTS ml_detections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_ip TEXT NOT NULL,
    threat_type TEXT NOT NULL,
    confidence REAL NOT NULL,
    features TEXT, -- JSON string
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ml_detections_created ON ml_detections(created_at);
CREATE INDEX idx_ml_detections_type ON ml_detections(threat_type);

-- Alerts table
CREATE TABLE IF NOT EXISTS alerts (
    id TEXT PRIMARY KEY,
    level TEXT NOT NULL,
    message TEXT NOT NULL,
    source TEXT,
    resolved BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    resolved_at DATETIME
);

CREATE INDEX idx_alerts_level ON alerts(level);
CREATE INDEX idx_alerts_resolved ON alerts(resolved);

-- API logs table (for monitoring and debugging)
CREATE TABLE IF NOT EXISTS api_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    method TEXT NOT NULL,
    path TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    response_time_ms INTEGER,
    ip_address TEXT,
    user_agent TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_api_logs_timestamp ON api_logs(timestamp);
CREATE INDEX idx_api_logs_status ON api_logs(status_code);

-- Insert sample data for testing
INSERT INTO threats (id, source_ip, threat_type, severity, description) VALUES
    ('threat_001', '192.168.1.100', 'port_scan', 'medium', 'Suspicious port scanning detected'),
    ('threat_002', '10.0.0.50', 'brute_force', 'high', 'Multiple failed login attempts'),
    ('threat_003', '172.16.0.25', 'ddos', 'critical', 'DDoS attack pattern detected');

INSERT INTO devices (id, name, type, ip_address, mac_address, status) VALUES
    ('device_001', 'Main Server', 'server', '192.168.1.10', '00:1B:44:11:3A:B7', 'online'),
    ('device_002', 'Firewall', 'firewall', '192.168.1.1', '00:1B:44:11:3A:B8', 'online'),
    ('device_003', 'IDS Sensor', 'sensor', '192.168.1.20', '00:1B:44:11:3A:B9', 'online');

