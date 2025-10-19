# Axiom Backend V3 資料庫遷移指南

> **版本**: 3.0.0  
> **日期**: 2025-10-16

---

## 📋 概述

本指南說明如何安全地執行 Axiom Backend V3 的資料庫遷移，包括初始化部署和後續升級。

---

## 遷移策略

### 方式 1: 自動遷移 (推薦用於開發/測試)

應用啟動時自動執行 GORM AutoMigrate。

**優點**:
- 簡單快速
- 自動處理表結構變更
- 適合開發環境

**缺點**:
- 缺乏精細控制
- 無法回滾
- 不適合生產環境

### 方式 2: SQL 遷移腳本 (推薦用於生產)

使用手動編寫的 SQL 腳本。

**優點**:
- 完全控制
- 可審查和測試
- 支援回滾
- 適合生產環境

**缺點**:
- 需要手動維護
- 需要更多工作

---

## 初始部署 (全新安裝)

### 步驟 1: 準備資料庫

```bash
# 創建資料庫和用戶
sudo -u postgres createdb pandora_db
sudo -u postgres createuser pandora
sudo -u postgres psql -c "ALTER USER pandora WITH PASSWORD 'pandora123';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE pandora_db TO pandora;"
```

### 步驟 2: 執行初始遷移

#### 方式 A: 使用 Auto Migrate

```bash
cd Application/be
go run cmd/server/main.go
```

應用會自動創建所有表。

#### 方式 B: 使用 SQL 腳本 (推薦)

```bash
psql -U pandora -d pandora_db -h localhost -f database/migrations/001_initial_schema.sql
```

### 步驟 3: 驗證

```sql
-- 連接資料庫
psql -U pandora -d pandora_db -h localhost

-- 列出所有表
\dt

-- 應該看到以下 9 個表：
-- services
-- config_histories
-- quantum_jobs
-- windows_logs
-- alerts
-- api_logs
-- metric_snapshots
-- users
-- sessions

-- 檢查表結構
\d services
\d quantum_jobs
\d windows_logs
```

---

## 版本升級

### 從 V2 升級到 V3

#### 前置檢查

```bash
# 1. 備份資料庫
pg_dump -U pandora -h localhost pandora_db > backup_v2_$(date +%Y%m%d_%H%M%S).sql

# 2. 檢查磁碟空間
df -h

# 3. 檢查當前版本
curl http://localhost:3001/health
```

#### 遷移步驟

**步驟 1: 停止服務**

```bash
cd Application
docker-compose stop axiom-be
```

**步驟 2: 更新代碼**

```bash
git pull origin main
cd Application/be
go mod download
```

**步驟 3: 執行遷移**

```bash
# 如果使用 SQL 腳本
psql -U pandora -d pandora_db -h localhost -f database/migrations/002_v3_upgrade.sql

# 或使用 GORM Auto Migrate
go run cmd/server/main.go migrate
```

**步驟 4: 啟動新版本**

```bash
cd Application
docker-compose up -d axiom-be-v3
```

**步驟 5: 驗證**

```bash
# 檢查健康狀態
curl http://localhost:3001/health

# 應該返回版本 3.0.0
```

#### 回滾方案

如果升級失敗：

```bash
# 1. 停止新版本
docker-compose stop axiom-be-v3

# 2. 恢復資料庫
psql -U pandora -d pandora_db -h localhost < backup_v2_20251016_100000.sql

# 3. 啟動舊版本
docker-compose up -d axiom-be
```

---

## 遷移檔案列表

### database/migrations/

| 檔案 | 版本 | 說明 |
|------|------|------|
| `001_initial_schema.sql` | 3.0.0 | 初始 Schema |
| `002_add_indexes.sql` | 3.0.1 | 添加性能索引 |
| `003_add_audit_fields.sql` | 3.1.0 | 添加審計欄位 |

---

## 資料遷移

### 遷移舊版本數據

如果從舊系統遷移數據：

```sql
-- 遷移舊的量子作業數據
INSERT INTO quantum_jobs (job_id, type, status, submitted_at, ...)
SELECT 
    old_job_id,
    job_type,
    job_status,
    submit_time,
    ...
FROM old_quantum_jobs_table
WHERE submit_time >= '2025-01-01';

-- 遷移 Windows 日誌
INSERT INTO windows_logs (agent_id, log_type, event_id, message, time_created, received_at)
SELECT
    agent_identifier,
    log_category,
    event_number,
    log_message,
    event_time,
    CURRENT_TIMESTAMP
FROM old_event_logs
WHERE event_time >= NOW() - INTERVAL '30 days';
```

---

## 索引優化

### 創建額外索引

根據查詢模式創建額外索引：

```sql
-- Windows 日誌常用查詢
CREATE INDEX idx_windows_logs_agent_type_time 
ON windows_logs(agent_id, log_type, time_created DESC);

-- 量子作業常用查詢
CREATE INDEX idx_quantum_jobs_type_time 
ON quantum_jobs(type, submitted_at DESC);

-- API 日誌分析
CREATE INDEX idx_api_logs_path_status_time 
ON api_logs(path, status, created_at DESC);
```

### 分區表 (大數據量)

如果單表超過 1000 萬行，考慮分區：

```sql
-- Windows 日誌按月分區
CREATE TABLE windows_logs_y2025m10 PARTITION OF windows_logs
FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

CREATE TABLE windows_logs_y2025m11 PARTITION OF windows_logs
FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
```

---

## 數據清理

### 清理舊數據

```sql
-- 刪除 30 天前的 Windows 日誌
DELETE FROM windows_logs 
WHERE time_created < NOW() - INTERVAL '30 days';

-- 刪除 7 天前的 API 日誌
DELETE FROM api_logs 
WHERE created_at < NOW() - INTERVAL '7 days';

-- 刪除已完成的舊量子作業 (90 天前)
DELETE FROM quantum_jobs 
WHERE status = 'completed' 
AND completed_at < NOW() - INTERVAL '90 days';
```

### 設置自動清理

```sql
-- 創建清理函數
CREATE OR REPLACE FUNCTION cleanup_old_logs()
RETURNS void AS $$
BEGIN
    DELETE FROM windows_logs WHERE time_created < NOW() - INTERVAL '30 days';
    DELETE FROM api_logs WHERE created_at < NOW() - INTERVAL '7 days';
    DELETE FROM quantum_jobs WHERE status = 'completed' AND completed_at < NOW() - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;

-- 使用 pg_cron 每天執行
SELECT cron.schedule('cleanup-old-logs', '0 2 * * *', 'SELECT cleanup_old_logs();');
```

---

## 性能監控

### 查詢慢查詢

```sql
-- 啟用慢查詢日誌
ALTER SYSTEM SET log_min_duration_statement = 1000; -- 1秒

-- 查看慢查詢
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### 表大小監控

```sql
-- 查看表大小
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_total_relation_size(schemaname||'.'||tablename) AS size_bytes
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY size_bytes DESC;
```

---

## 備份和恢復

### 每日備份

```bash
#!/bin/bash
# backup-daily.sh

BACKUP_DIR="/var/backups/pandora"
DATE=$(date +%Y%m%d_%H%M%S)

# 創建備份
pg_dump -U pandora -h localhost pandora_db > "$BACKUP_DIR/pandora_db_$DATE.sql"

# 壓縮備份
gzip "$BACKUP_DIR/pandora_db_$DATE.sql"

# 刪除 7 天前的備份
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +7 -delete

echo "Backup completed: pandora_db_$DATE.sql.gz"
```

### 恢復備份

```bash
# 1. 停止服務
docker-compose stop axiom-be-v3

# 2. 恢復資料庫
gunzip -c /var/backups/pandora/pandora_db_20251016_020000.sql.gz | \
psql -U pandora -d pandora_db -h localhost

# 3. 啟動服務
docker-compose up -d axiom-be-v3
```

---

## 注意事項

### ⚠️ 重要提醒

1. **生產環境**：永遠先在測試環境驗證遷移
2. **備份**：執行任何遷移前務必備份
3. **維護窗口**：在低峰時段執行遷移
4. **回滾計劃**：準備回滾腳本
5. **通知**：提前通知用戶

### ✅ 遷移檢查清單

- [ ] 備份資料庫完成
- [ ] 測試環境驗證通過
- [ ] 回滾腳本準備完成
- [ ] 維護通知已發送
- [ ] 監控系統就緒
- [ ] 回滾計劃已準備
- [ ] 執行遷移
- [ ] 驗證數據完整性
- [ ] 驗證應用功能
- [ ] 監控系統狀態

---

## 聯繫支援

如遇到遷移問題，請聯繫技術支援：
- Email: support@pandora.local
- Issue: [GitHub Issues](https://github.com/your-org/Local_IPS-IDS/issues)

---

**文檔版本**: 3.0.0  
**最後更新**: 2025-10-16



> **版本**: 3.0.0  
> **日期**: 2025-10-16

---

## 📋 概述

本指南說明如何安全地執行 Axiom Backend V3 的資料庫遷移，包括初始化部署和後續升級。

---

## 遷移策略

### 方式 1: 自動遷移 (推薦用於開發/測試)

應用啟動時自動執行 GORM AutoMigrate。

**優點**:
- 簡單快速
- 自動處理表結構變更
- 適合開發環境

**缺點**:
- 缺乏精細控制
- 無法回滾
- 不適合生產環境

### 方式 2: SQL 遷移腳本 (推薦用於生產)

使用手動編寫的 SQL 腳本。

**優點**:
- 完全控制
- 可審查和測試
- 支援回滾
- 適合生產環境

**缺點**:
- 需要手動維護
- 需要更多工作

---

## 初始部署 (全新安裝)

### 步驟 1: 準備資料庫

```bash
# 創建資料庫和用戶
sudo -u postgres createdb pandora_db
sudo -u postgres createuser pandora
sudo -u postgres psql -c "ALTER USER pandora WITH PASSWORD 'pandora123';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE pandora_db TO pandora;"
```

### 步驟 2: 執行初始遷移

#### 方式 A: 使用 Auto Migrate

```bash
cd Application/be
go run cmd/server/main.go
```

應用會自動創建所有表。

#### 方式 B: 使用 SQL 腳本 (推薦)

```bash
psql -U pandora -d pandora_db -h localhost -f database/migrations/001_initial_schema.sql
```

### 步驟 3: 驗證

```sql
-- 連接資料庫
psql -U pandora -d pandora_db -h localhost

-- 列出所有表
\dt

-- 應該看到以下 9 個表：
-- services
-- config_histories
-- quantum_jobs
-- windows_logs
-- alerts
-- api_logs
-- metric_snapshots
-- users
-- sessions

-- 檢查表結構
\d services
\d quantum_jobs
\d windows_logs
```

---

## 版本升級

### 從 V2 升級到 V3

#### 前置檢查

```bash
# 1. 備份資料庫
pg_dump -U pandora -h localhost pandora_db > backup_v2_$(date +%Y%m%d_%H%M%S).sql

# 2. 檢查磁碟空間
df -h

# 3. 檢查當前版本
curl http://localhost:3001/health
```

#### 遷移步驟

**步驟 1: 停止服務**

```bash
cd Application
docker-compose stop axiom-be
```

**步驟 2: 更新代碼**

```bash
git pull origin main
cd Application/be
go mod download
```

**步驟 3: 執行遷移**

```bash
# 如果使用 SQL 腳本
psql -U pandora -d pandora_db -h localhost -f database/migrations/002_v3_upgrade.sql

# 或使用 GORM Auto Migrate
go run cmd/server/main.go migrate
```

**步驟 4: 啟動新版本**

```bash
cd Application
docker-compose up -d axiom-be-v3
```

**步驟 5: 驗證**

```bash
# 檢查健康狀態
curl http://localhost:3001/health

# 應該返回版本 3.0.0
```

#### 回滾方案

如果升級失敗：

```bash
# 1. 停止新版本
docker-compose stop axiom-be-v3

# 2. 恢復資料庫
psql -U pandora -d pandora_db -h localhost < backup_v2_20251016_100000.sql

# 3. 啟動舊版本
docker-compose up -d axiom-be
```

---

## 遷移檔案列表

### database/migrations/

| 檔案 | 版本 | 說明 |
|------|------|------|
| `001_initial_schema.sql` | 3.0.0 | 初始 Schema |
| `002_add_indexes.sql` | 3.0.1 | 添加性能索引 |
| `003_add_audit_fields.sql` | 3.1.0 | 添加審計欄位 |

---

## 資料遷移

### 遷移舊版本數據

如果從舊系統遷移數據：

```sql
-- 遷移舊的量子作業數據
INSERT INTO quantum_jobs (job_id, type, status, submitted_at, ...)
SELECT 
    old_job_id,
    job_type,
    job_status,
    submit_time,
    ...
FROM old_quantum_jobs_table
WHERE submit_time >= '2025-01-01';

-- 遷移 Windows 日誌
INSERT INTO windows_logs (agent_id, log_type, event_id, message, time_created, received_at)
SELECT
    agent_identifier,
    log_category,
    event_number,
    log_message,
    event_time,
    CURRENT_TIMESTAMP
FROM old_event_logs
WHERE event_time >= NOW() - INTERVAL '30 days';
```

---

## 索引優化

### 創建額外索引

根據查詢模式創建額外索引：

```sql
-- Windows 日誌常用查詢
CREATE INDEX idx_windows_logs_agent_type_time 
ON windows_logs(agent_id, log_type, time_created DESC);

-- 量子作業常用查詢
CREATE INDEX idx_quantum_jobs_type_time 
ON quantum_jobs(type, submitted_at DESC);

-- API 日誌分析
CREATE INDEX idx_api_logs_path_status_time 
ON api_logs(path, status, created_at DESC);
```

### 分區表 (大數據量)

如果單表超過 1000 萬行，考慮分區：

```sql
-- Windows 日誌按月分區
CREATE TABLE windows_logs_y2025m10 PARTITION OF windows_logs
FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

CREATE TABLE windows_logs_y2025m11 PARTITION OF windows_logs
FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
```

---

## 數據清理

### 清理舊數據

```sql
-- 刪除 30 天前的 Windows 日誌
DELETE FROM windows_logs 
WHERE time_created < NOW() - INTERVAL '30 days';

-- 刪除 7 天前的 API 日誌
DELETE FROM api_logs 
WHERE created_at < NOW() - INTERVAL '7 days';

-- 刪除已完成的舊量子作業 (90 天前)
DELETE FROM quantum_jobs 
WHERE status = 'completed' 
AND completed_at < NOW() - INTERVAL '90 days';
```

### 設置自動清理

```sql
-- 創建清理函數
CREATE OR REPLACE FUNCTION cleanup_old_logs()
RETURNS void AS $$
BEGIN
    DELETE FROM windows_logs WHERE time_created < NOW() - INTERVAL '30 days';
    DELETE FROM api_logs WHERE created_at < NOW() - INTERVAL '7 days';
    DELETE FROM quantum_jobs WHERE status = 'completed' AND completed_at < NOW() - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;

-- 使用 pg_cron 每天執行
SELECT cron.schedule('cleanup-old-logs', '0 2 * * *', 'SELECT cleanup_old_logs();');
```

---

## 性能監控

### 查詢慢查詢

```sql
-- 啟用慢查詢日誌
ALTER SYSTEM SET log_min_duration_statement = 1000; -- 1秒

-- 查看慢查詢
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### 表大小監控

```sql
-- 查看表大小
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_total_relation_size(schemaname||'.'||tablename) AS size_bytes
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY size_bytes DESC;
```

---

## 備份和恢復

### 每日備份

```bash
#!/bin/bash
# backup-daily.sh

BACKUP_DIR="/var/backups/pandora"
DATE=$(date +%Y%m%d_%H%M%S)

# 創建備份
pg_dump -U pandora -h localhost pandora_db > "$BACKUP_DIR/pandora_db_$DATE.sql"

# 壓縮備份
gzip "$BACKUP_DIR/pandora_db_$DATE.sql"

# 刪除 7 天前的備份
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +7 -delete

echo "Backup completed: pandora_db_$DATE.sql.gz"
```

### 恢復備份

```bash
# 1. 停止服務
docker-compose stop axiom-be-v3

# 2. 恢復資料庫
gunzip -c /var/backups/pandora/pandora_db_20251016_020000.sql.gz | \
psql -U pandora -d pandora_db -h localhost

# 3. 啟動服務
docker-compose up -d axiom-be-v3
```

---

## 注意事項

### ⚠️ 重要提醒

1. **生產環境**：永遠先在測試環境驗證遷移
2. **備份**：執行任何遷移前務必備份
3. **維護窗口**：在低峰時段執行遷移
4. **回滾計劃**：準備回滾腳本
5. **通知**：提前通知用戶

### ✅ 遷移檢查清單

- [ ] 備份資料庫完成
- [ ] 測試環境驗證通過
- [ ] 回滾腳本準備完成
- [ ] 維護通知已發送
- [ ] 監控系統就緒
- [ ] 回滾計劃已準備
- [ ] 執行遷移
- [ ] 驗證數據完整性
- [ ] 驗證應用功能
- [ ] 監控系統狀態

---

## 聯繫支援

如遇到遷移問題，請聯繫技術支援：
- Email: support@pandora.local
- Issue: [GitHub Issues](https://github.com/your-org/Local_IPS-IDS/issues)

---

**文檔版本**: 3.0.0  
**最後更新**: 2025-10-16


> **版本**: 3.0.0  
> **日期**: 2025-10-16

---

## 📋 概述

本指南說明如何安全地執行 Axiom Backend V3 的資料庫遷移，包括初始化部署和後續升級。

---

## 遷移策略

### 方式 1: 自動遷移 (推薦用於開發/測試)

應用啟動時自動執行 GORM AutoMigrate。

**優點**:
- 簡單快速
- 自動處理表結構變更
- 適合開發環境

**缺點**:
- 缺乏精細控制
- 無法回滾
- 不適合生產環境

### 方式 2: SQL 遷移腳本 (推薦用於生產)

使用手動編寫的 SQL 腳本。

**優點**:
- 完全控制
- 可審查和測試
- 支援回滾
- 適合生產環境

**缺點**:
- 需要手動維護
- 需要更多工作

---

## 初始部署 (全新安裝)

### 步驟 1: 準備資料庫

```bash
# 創建資料庫和用戶
sudo -u postgres createdb pandora_db
sudo -u postgres createuser pandora
sudo -u postgres psql -c "ALTER USER pandora WITH PASSWORD 'pandora123';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE pandora_db TO pandora;"
```

### 步驟 2: 執行初始遷移

#### 方式 A: 使用 Auto Migrate

```bash
cd Application/be
go run cmd/server/main.go
```

應用會自動創建所有表。

#### 方式 B: 使用 SQL 腳本 (推薦)

```bash
psql -U pandora -d pandora_db -h localhost -f database/migrations/001_initial_schema.sql
```

### 步驟 3: 驗證

```sql
-- 連接資料庫
psql -U pandora -d pandora_db -h localhost

-- 列出所有表
\dt

-- 應該看到以下 9 個表：
-- services
-- config_histories
-- quantum_jobs
-- windows_logs
-- alerts
-- api_logs
-- metric_snapshots
-- users
-- sessions

-- 檢查表結構
\d services
\d quantum_jobs
\d windows_logs
```

---

## 版本升級

### 從 V2 升級到 V3

#### 前置檢查

```bash
# 1. 備份資料庫
pg_dump -U pandora -h localhost pandora_db > backup_v2_$(date +%Y%m%d_%H%M%S).sql

# 2. 檢查磁碟空間
df -h

# 3. 檢查當前版本
curl http://localhost:3001/health
```

#### 遷移步驟

**步驟 1: 停止服務**

```bash
cd Application
docker-compose stop axiom-be
```

**步驟 2: 更新代碼**

```bash
git pull origin main
cd Application/be
go mod download
```

**步驟 3: 執行遷移**

```bash
# 如果使用 SQL 腳本
psql -U pandora -d pandora_db -h localhost -f database/migrations/002_v3_upgrade.sql

# 或使用 GORM Auto Migrate
go run cmd/server/main.go migrate
```

**步驟 4: 啟動新版本**

```bash
cd Application
docker-compose up -d axiom-be-v3
```

**步驟 5: 驗證**

```bash
# 檢查健康狀態
curl http://localhost:3001/health

# 應該返回版本 3.0.0
```

#### 回滾方案

如果升級失敗：

```bash
# 1. 停止新版本
docker-compose stop axiom-be-v3

# 2. 恢復資料庫
psql -U pandora -d pandora_db -h localhost < backup_v2_20251016_100000.sql

# 3. 啟動舊版本
docker-compose up -d axiom-be
```

---

## 遷移檔案列表

### database/migrations/

| 檔案 | 版本 | 說明 |
|------|------|------|
| `001_initial_schema.sql` | 3.0.0 | 初始 Schema |
| `002_add_indexes.sql` | 3.0.1 | 添加性能索引 |
| `003_add_audit_fields.sql` | 3.1.0 | 添加審計欄位 |

---

## 資料遷移

### 遷移舊版本數據

如果從舊系統遷移數據：

```sql
-- 遷移舊的量子作業數據
INSERT INTO quantum_jobs (job_id, type, status, submitted_at, ...)
SELECT 
    old_job_id,
    job_type,
    job_status,
    submit_time,
    ...
FROM old_quantum_jobs_table
WHERE submit_time >= '2025-01-01';

-- 遷移 Windows 日誌
INSERT INTO windows_logs (agent_id, log_type, event_id, message, time_created, received_at)
SELECT
    agent_identifier,
    log_category,
    event_number,
    log_message,
    event_time,
    CURRENT_TIMESTAMP
FROM old_event_logs
WHERE event_time >= NOW() - INTERVAL '30 days';
```

---

## 索引優化

### 創建額外索引

根據查詢模式創建額外索引：

```sql
-- Windows 日誌常用查詢
CREATE INDEX idx_windows_logs_agent_type_time 
ON windows_logs(agent_id, log_type, time_created DESC);

-- 量子作業常用查詢
CREATE INDEX idx_quantum_jobs_type_time 
ON quantum_jobs(type, submitted_at DESC);

-- API 日誌分析
CREATE INDEX idx_api_logs_path_status_time 
ON api_logs(path, status, created_at DESC);
```

### 分區表 (大數據量)

如果單表超過 1000 萬行，考慮分區：

```sql
-- Windows 日誌按月分區
CREATE TABLE windows_logs_y2025m10 PARTITION OF windows_logs
FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

CREATE TABLE windows_logs_y2025m11 PARTITION OF windows_logs
FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
```

---

## 數據清理

### 清理舊數據

```sql
-- 刪除 30 天前的 Windows 日誌
DELETE FROM windows_logs 
WHERE time_created < NOW() - INTERVAL '30 days';

-- 刪除 7 天前的 API 日誌
DELETE FROM api_logs 
WHERE created_at < NOW() - INTERVAL '7 days';

-- 刪除已完成的舊量子作業 (90 天前)
DELETE FROM quantum_jobs 
WHERE status = 'completed' 
AND completed_at < NOW() - INTERVAL '90 days';
```

### 設置自動清理

```sql
-- 創建清理函數
CREATE OR REPLACE FUNCTION cleanup_old_logs()
RETURNS void AS $$
BEGIN
    DELETE FROM windows_logs WHERE time_created < NOW() - INTERVAL '30 days';
    DELETE FROM api_logs WHERE created_at < NOW() - INTERVAL '7 days';
    DELETE FROM quantum_jobs WHERE status = 'completed' AND completed_at < NOW() - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;

-- 使用 pg_cron 每天執行
SELECT cron.schedule('cleanup-old-logs', '0 2 * * *', 'SELECT cleanup_old_logs();');
```

---

## 性能監控

### 查詢慢查詢

```sql
-- 啟用慢查詢日誌
ALTER SYSTEM SET log_min_duration_statement = 1000; -- 1秒

-- 查看慢查詢
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### 表大小監控

```sql
-- 查看表大小
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_total_relation_size(schemaname||'.'||tablename) AS size_bytes
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY size_bytes DESC;
```

---

## 備份和恢復

### 每日備份

```bash
#!/bin/bash
# backup-daily.sh

BACKUP_DIR="/var/backups/pandora"
DATE=$(date +%Y%m%d_%H%M%S)

# 創建備份
pg_dump -U pandora -h localhost pandora_db > "$BACKUP_DIR/pandora_db_$DATE.sql"

# 壓縮備份
gzip "$BACKUP_DIR/pandora_db_$DATE.sql"

# 刪除 7 天前的備份
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +7 -delete

echo "Backup completed: pandora_db_$DATE.sql.gz"
```

### 恢復備份

```bash
# 1. 停止服務
docker-compose stop axiom-be-v3

# 2. 恢復資料庫
gunzip -c /var/backups/pandora/pandora_db_20251016_020000.sql.gz | \
psql -U pandora -d pandora_db -h localhost

# 3. 啟動服務
docker-compose up -d axiom-be-v3
```

---

## 注意事項

### ⚠️ 重要提醒

1. **生產環境**：永遠先在測試環境驗證遷移
2. **備份**：執行任何遷移前務必備份
3. **維護窗口**：在低峰時段執行遷移
4. **回滾計劃**：準備回滾腳本
5. **通知**：提前通知用戶

### ✅ 遷移檢查清單

- [ ] 備份資料庫完成
- [ ] 測試環境驗證通過
- [ ] 回滾腳本準備完成
- [ ] 維護通知已發送
- [ ] 監控系統就緒
- [ ] 回滾計劃已準備
- [ ] 執行遷移
- [ ] 驗證數據完整性
- [ ] 驗證應用功能
- [ ] 監控系統狀態

---

## 聯繫支援

如遇到遷移問題，請聯繫技術支援：
- Email: support@pandora.local
- Issue: [GitHub Issues](https://github.com/your-org/Local_IPS-IDS/issues)

---

**文檔版本**: 3.0.0  
**最後更新**: 2025-10-16



> **版本**: 3.0.0  
> **日期**: 2025-10-16

---

## 📋 概述

本指南說明如何安全地執行 Axiom Backend V3 的資料庫遷移，包括初始化部署和後續升級。

---

## 遷移策略

### 方式 1: 自動遷移 (推薦用於開發/測試)

應用啟動時自動執行 GORM AutoMigrate。

**優點**:
- 簡單快速
- 自動處理表結構變更
- 適合開發環境

**缺點**:
- 缺乏精細控制
- 無法回滾
- 不適合生產環境

### 方式 2: SQL 遷移腳本 (推薦用於生產)

使用手動編寫的 SQL 腳本。

**優點**:
- 完全控制
- 可審查和測試
- 支援回滾
- 適合生產環境

**缺點**:
- 需要手動維護
- 需要更多工作

---

## 初始部署 (全新安裝)

### 步驟 1: 準備資料庫

```bash
# 創建資料庫和用戶
sudo -u postgres createdb pandora_db
sudo -u postgres createuser pandora
sudo -u postgres psql -c "ALTER USER pandora WITH PASSWORD 'pandora123';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE pandora_db TO pandora;"
```

### 步驟 2: 執行初始遷移

#### 方式 A: 使用 Auto Migrate

```bash
cd Application/be
go run cmd/server/main.go
```

應用會自動創建所有表。

#### 方式 B: 使用 SQL 腳本 (推薦)

```bash
psql -U pandora -d pandora_db -h localhost -f database/migrations/001_initial_schema.sql
```

### 步驟 3: 驗證

```sql
-- 連接資料庫
psql -U pandora -d pandora_db -h localhost

-- 列出所有表
\dt

-- 應該看到以下 9 個表：
-- services
-- config_histories
-- quantum_jobs
-- windows_logs
-- alerts
-- api_logs
-- metric_snapshots
-- users
-- sessions

-- 檢查表結構
\d services
\d quantum_jobs
\d windows_logs
```

---

## 版本升級

### 從 V2 升級到 V3

#### 前置檢查

```bash
# 1. 備份資料庫
pg_dump -U pandora -h localhost pandora_db > backup_v2_$(date +%Y%m%d_%H%M%S).sql

# 2. 檢查磁碟空間
df -h

# 3. 檢查當前版本
curl http://localhost:3001/health
```

#### 遷移步驟

**步驟 1: 停止服務**

```bash
cd Application
docker-compose stop axiom-be
```

**步驟 2: 更新代碼**

```bash
git pull origin main
cd Application/be
go mod download
```

**步驟 3: 執行遷移**

```bash
# 如果使用 SQL 腳本
psql -U pandora -d pandora_db -h localhost -f database/migrations/002_v3_upgrade.sql

# 或使用 GORM Auto Migrate
go run cmd/server/main.go migrate
```

**步驟 4: 啟動新版本**

```bash
cd Application
docker-compose up -d axiom-be-v3
```

**步驟 5: 驗證**

```bash
# 檢查健康狀態
curl http://localhost:3001/health

# 應該返回版本 3.0.0
```

#### 回滾方案

如果升級失敗：

```bash
# 1. 停止新版本
docker-compose stop axiom-be-v3

# 2. 恢復資料庫
psql -U pandora -d pandora_db -h localhost < backup_v2_20251016_100000.sql

# 3. 啟動舊版本
docker-compose up -d axiom-be
```

---

## 遷移檔案列表

### database/migrations/

| 檔案 | 版本 | 說明 |
|------|------|------|
| `001_initial_schema.sql` | 3.0.0 | 初始 Schema |
| `002_add_indexes.sql` | 3.0.1 | 添加性能索引 |
| `003_add_audit_fields.sql` | 3.1.0 | 添加審計欄位 |

---

## 資料遷移

### 遷移舊版本數據

如果從舊系統遷移數據：

```sql
-- 遷移舊的量子作業數據
INSERT INTO quantum_jobs (job_id, type, status, submitted_at, ...)
SELECT 
    old_job_id,
    job_type,
    job_status,
    submit_time,
    ...
FROM old_quantum_jobs_table
WHERE submit_time >= '2025-01-01';

-- 遷移 Windows 日誌
INSERT INTO windows_logs (agent_id, log_type, event_id, message, time_created, received_at)
SELECT
    agent_identifier,
    log_category,
    event_number,
    log_message,
    event_time,
    CURRENT_TIMESTAMP
FROM old_event_logs
WHERE event_time >= NOW() - INTERVAL '30 days';
```

---

## 索引優化

### 創建額外索引

根據查詢模式創建額外索引：

```sql
-- Windows 日誌常用查詢
CREATE INDEX idx_windows_logs_agent_type_time 
ON windows_logs(agent_id, log_type, time_created DESC);

-- 量子作業常用查詢
CREATE INDEX idx_quantum_jobs_type_time 
ON quantum_jobs(type, submitted_at DESC);

-- API 日誌分析
CREATE INDEX idx_api_logs_path_status_time 
ON api_logs(path, status, created_at DESC);
```

### 分區表 (大數據量)

如果單表超過 1000 萬行，考慮分區：

```sql
-- Windows 日誌按月分區
CREATE TABLE windows_logs_y2025m10 PARTITION OF windows_logs
FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');

CREATE TABLE windows_logs_y2025m11 PARTITION OF windows_logs
FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
```

---

## 數據清理

### 清理舊數據

```sql
-- 刪除 30 天前的 Windows 日誌
DELETE FROM windows_logs 
WHERE time_created < NOW() - INTERVAL '30 days';

-- 刪除 7 天前的 API 日誌
DELETE FROM api_logs 
WHERE created_at < NOW() - INTERVAL '7 days';

-- 刪除已完成的舊量子作業 (90 天前)
DELETE FROM quantum_jobs 
WHERE status = 'completed' 
AND completed_at < NOW() - INTERVAL '90 days';
```

### 設置自動清理

```sql
-- 創建清理函數
CREATE OR REPLACE FUNCTION cleanup_old_logs()
RETURNS void AS $$
BEGIN
    DELETE FROM windows_logs WHERE time_created < NOW() - INTERVAL '30 days';
    DELETE FROM api_logs WHERE created_at < NOW() - INTERVAL '7 days';
    DELETE FROM quantum_jobs WHERE status = 'completed' AND completed_at < NOW() - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;

-- 使用 pg_cron 每天執行
SELECT cron.schedule('cleanup-old-logs', '0 2 * * *', 'SELECT cleanup_old_logs();');
```

---

## 性能監控

### 查詢慢查詢

```sql
-- 啟用慢查詢日誌
ALTER SYSTEM SET log_min_duration_statement = 1000; -- 1秒

-- 查看慢查詢
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### 表大小監控

```sql
-- 查看表大小
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_total_relation_size(schemaname||'.'||tablename) AS size_bytes
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY size_bytes DESC;
```

---

## 備份和恢復

### 每日備份

```bash
#!/bin/bash
# backup-daily.sh

BACKUP_DIR="/var/backups/pandora"
DATE=$(date +%Y%m%d_%H%M%S)

# 創建備份
pg_dump -U pandora -h localhost pandora_db > "$BACKUP_DIR/pandora_db_$DATE.sql"

# 壓縮備份
gzip "$BACKUP_DIR/pandora_db_$DATE.sql"

# 刪除 7 天前的備份
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +7 -delete

echo "Backup completed: pandora_db_$DATE.sql.gz"
```

### 恢復備份

```bash
# 1. 停止服務
docker-compose stop axiom-be-v3

# 2. 恢復資料庫
gunzip -c /var/backups/pandora/pandora_db_20251016_020000.sql.gz | \
psql -U pandora -d pandora_db -h localhost

# 3. 啟動服務
docker-compose up -d axiom-be-v3
```

---

## 注意事項

### ⚠️ 重要提醒

1. **生產環境**：永遠先在測試環境驗證遷移
2. **備份**：執行任何遷移前務必備份
3. **維護窗口**：在低峰時段執行遷移
4. **回滾計劃**：準備回滾腳本
5. **通知**：提前通知用戶

### ✅ 遷移檢查清單

- [ ] 備份資料庫完成
- [ ] 測試環境驗證通過
- [ ] 回滾腳本準備完成
- [ ] 維護通知已發送
- [ ] 監控系統就緒
- [ ] 回滾計劃已準備
- [ ] 執行遷移
- [ ] 驗證數據完整性
- [ ] 驗證應用功能
- [ ] 監控系統狀態

---

## 聯繫支援

如遇到遷移問題，請聯繫技術支援：
- Email: support@pandora.local
- Issue: [GitHub Issues](https://github.com/your-org/Local_IPS-IDS/issues)

---

**文檔版本**: 3.0.0  
**最後更新**: 2025-10-16

