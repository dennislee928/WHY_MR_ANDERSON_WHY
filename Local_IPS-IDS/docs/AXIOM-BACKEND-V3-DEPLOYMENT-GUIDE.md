# Axiom Backend V3 部署指南

> **版本**: 3.0.0  
> **日期**: 2025-10-16

---

## 📋 目錄

1. [系統要求](#系統要求)
2. [環境準備](#環境準備)
3. [Docker 部署](#docker-部署)
4. [手動部署](#手動部署)
5. [配置說明](#配置說明)
6. [資料庫遷移](#資料庫遷移)
7. [驗證部署](#驗證部署)
8. [故障排除](#故障排除)

---

## 系統要求

### 硬體要求
- **CPU**: 4 核心以上
- **記憶體**: 8GB 以上
- **磁碟**: 50GB 可用空間

### 軟體要求
- **Go**: 1.21 或更高版本
- **PostgreSQL**: 15 或更高版本
- **Redis**: 7 或更高版本
- **Docker**: 24.0 或更高版本 (可選)
- **Docker Compose**: 2.20 或更高版本 (可選)

### 依賴服務
- Prometheus
- Grafana
- Loki
- Cyber-AI-Quantum
- Nginx
- RabbitMQ

---

## 環境準備

### 1. 安裝 PostgreSQL

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql-15

# 創建資料庫
sudo -u postgres createdb pandora_db
sudo -u postgres createuser pandora
sudo -u postgres psql -c "ALTER USER pandora WITH PASSWORD 'pandora123';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE pandora_db TO pandora;"
```

### 2. 安裝 Redis

```bash
# Ubuntu/Debian
sudo apt install redis-server

# 設置密碼
sudo redis-cli
> CONFIG SET requirepass "pandora123"
> exit
```

### 3. 克隆代碼

```bash
git clone https://github.com/your-org/Local_IPS-IDS.git
cd Local_IPS-IDS/Application/be
```

---

## Docker 部署

### 使用 Docker Compose

**最簡單的部署方式**

```bash
cd Application
docker-compose up -d axiom-be
```

`docker-compose.yml` 配置：

```yaml
services:
  axiom-be:
    build:
      context: ..
      dockerfile: Application/docker/axiom-be-v3.dockerfile
    container_name: axiom-be-v3
    restart: unless-stopped
    ports:
      - "3001:3001"
    environment:
      - PORT=3001
      - POSTGRES_HOST=postgres
      - POSTGRES_PORT=5432
      - POSTGRES_USER=pandora
      - POSTGRES_PASSWORD=pandora123
      - POSTGRES_DB=pandora_db
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=pandora123
      - REDIS_DB=0
      - PROMETHEUS_URL=http://prometheus:9090
      - GRAFANA_URL=http://grafana:3000
      - LOKI_URL=http://loki:3100
      - QUANTUM_URL=http://cyber-ai-quantum:8000
      - NGINX_URL=http://nginx:80
      - NGINX_CONFIG_PATH=/etc/nginx/nginx.conf
    depends_on:
      - postgres
      - redis
      - prometheus
      - loki
      - cyber-ai-quantum
    networks:
      - pandora-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3001/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### 查看日誌

```bash
docker logs -f axiom-be-v3
```

---

## 手動部署

### 1. 安裝依賴

```bash
cd Application/be
go mod download
```

### 2. 配置環境變數

```bash
cp .env.example .env
# 編輯 .env 文件，修改配置
nano .env
```

### 3. 運行資料庫遷移

```bash
make migrate
```

或手動執行：

```bash
go run cmd/server/main.go migrate
```

### 4. 構建應用

```bash
make build
```

### 5. 運行應用

```bash
./bin/axiom-backend
```

或使用 Makefile：

```bash
make run
```

---

## 配置說明

### 環境變數完整列表

| 變數名 | 默認值 | 說明 |
|--------|--------|------|
| `PORT` | 3001 | 服務端口 |
| `POSTGRES_HOST` | localhost | PostgreSQL 主機 |
| `POSTGRES_PORT` | 5432 | PostgreSQL 端口 |
| `POSTGRES_USER` | pandora | PostgreSQL 用戶 |
| `POSTGRES_PASSWORD` | pandora123 | PostgreSQL 密碼 |
| `POSTGRES_DB` | pandora_db | PostgreSQL 資料庫名 |
| `REDIS_HOST` | localhost | Redis 主機 |
| `REDIS_PORT` | 6379 | Redis 端口 |
| `REDIS_PASSWORD` | pandora123 | Redis 密碼 |
| `REDIS_DB` | 0 | Redis 資料庫編號 |
| `PROMETHEUS_URL` | http://localhost:9090 | Prometheus URL |
| `GRAFANA_URL` | http://localhost:3000 | Grafana URL |
| `LOKI_URL` | http://localhost:3100 | Loki URL |
| `QUANTUM_URL` | http://localhost:8000 | Quantum Service URL |
| `NGINX_URL` | http://localhost:80 | Nginx URL |
| `NGINX_CONFIG_PATH` | /etc/nginx/nginx.conf | Nginx 配置文件路徑 |

---

## 資料庫遷移

### 自動遷移

應用啟動時會自動執行 `AutoMigrate()`，創建所有必要的表。

### 手動遷移腳本

如需更精細的控制，使用遷移腳本：

```bash
cd database/migrations
# 執行所有遷移
./run-migrations.sh
```

### Migration 文件

位於 `database/migrations/`:
- `001_create_services_table.sql`
- `002_create_quantum_jobs_table.sql`
- `003_create_windows_logs_table.sql`
- 等...

---

## 驗證部署

### 1. 健康檢查

```bash
curl http://localhost:3001/health
```

預期響應：

```json
{
  "status": "healthy",
  "service": "axiom-backend-v3",
  "version": "3.0.0",
  "time": "..."
}
```

### 2. 測試 Prometheus API

```bash
curl -X POST http://localhost:3001/api/v2/prometheus/query \
  -H "Content-Type: application/json" \
  -d '{"query": "up"}'
```

### 3. 測試 Quantum API

```bash
curl -X POST http://localhost:3001/api/v2/quantum/qkd/generate \
  -H "Content-Type: application/json" \
  -d '{"key_length": 256, "backend": "simulator", "shots": 1024}'
```

### 4. 查看資料庫

```bash
psql -U pandora -d pandora_db -h localhost
\dt
```

應該看到 9 個表。

---

## 故障排除

### 問題 1: 無法連接資料庫

**症狀**: `Failed to connect to PostgreSQL`

**解決方案**:
1. 檢查 PostgreSQL 是否運行：`sudo systemctl status postgresql`
2. 檢查連接參數是否正確
3. 檢查防火牆設置

### 問題 2: Redis 連接失敗

**症狀**: `Failed to connect to Redis`

**解決方案**:
1. 檢查 Redis 是否運行：`sudo systemctl status redis`
2. 檢查密碼是否正確
3. 測試連接：`redis-cli -a pandora123 ping`

### 問題 3: Prometheus 查詢失敗

**症狀**: `prometheus health check failed`

**解決方案**:
1. 確認 Prometheus 服務運行：`curl http://localhost:9090/-/healthy`
2. 檢查 `PROMETHEUS_URL` 環境變數
3. 檢查網路連接

### 問題 4: 量子服務不可用

**症狀**: `quantum service health check failed`

**解決方案**:
1. 確認 cyber-ai-quantum 服務運行
2. 檢查 `QUANTUM_URL` 環境變數
3. 查看 quantum 服務日誌

---

## 性能調優

### 資料庫優化

```sql
-- 創建額外索引
CREATE INDEX idx_windows_logs_time ON windows_logs(time_created DESC);
CREATE INDEX idx_quantum_jobs_type_status ON quantum_jobs(type, status);
```

### Redis 優化

```conf
# redis.conf
maxmemory 2gb
maxmemory-policy allkeys-lru
```

### 連接池調優

在 `internal/database/db.go` 中：

```go
sqlDB.SetMaxIdleConns(20)    // 增加空閒連接
sqlDB.SetMaxOpenConns(200)   // 增加最大連接
sqlDB.SetConnMaxLifetime(time.Hour)
```

---

## 監控

### Prometheus Metrics

Axiom Backend 導出以下指標：

```
# HTTP 請求
axiom_http_requests_total{method, endpoint, status}
axiom_http_request_duration_seconds{method, endpoint}

# 資料庫
axiom_db_connections_active
axiom_db_query_duration_seconds

# 快取
axiom_cache_hits_total
axiom_cache_misses_total

# 量子作業
axiom_quantum_jobs_total{type, status}
axiom_quantum_job_duration_seconds{type}
```

### 健康檢查端點

- `/health` - 基本健康檢查
- `/api/v2/prometheus/health` - Prometheus 健康檢查
- `/api/v2/loki/health` - Loki 健康檢查
- `/api/v2/quantum/health` - Quantum 健康檢查

---

## 升級指南

### 從 V2 升級到 V3

1. **備份資料庫**

```bash
pg_dump -U pandora pandora_db > backup_$(date +%Y%m%d).sql
```

2. **停止舊版服務**

```bash
docker-compose stop axiom-be
```

3. **更新代碼**

```bash
git pull origin main
```

4. **運行遷移**

```bash
cd Application/be
go run cmd/server/main.go migrate
```

5. **啟動新版服務**

```bash
docker-compose up -d axiom-be-v3
```

6. **驗證**

```bash
curl http://localhost:3001/health
```

---

## 安全建議

1. **修改默認密碼**：修改 PostgreSQL、Redis 的默認密碼
2. **啟用 TLS**：在生產環境使用 HTTPS
3. **啟用認證**：配置 JWT 或 API Key 認證
4. **防火牆**：限制端口訪問
5. **定期備份**：設置自動備份計劃

---

## 支援

- **文檔**: [Complete Documentation](./README.md)
- **Issues**: [GitHub Issues](https://github.com/your-org/Local_IPS-IDS/issues)
- **Email**: support@example.com

---

**文檔版本**: 3.0.0  
**最後更新**: 2025-10-16



> **版本**: 3.0.0  
> **日期**: 2025-10-16

---

## 📋 目錄

1. [系統要求](#系統要求)
2. [環境準備](#環境準備)
3. [Docker 部署](#docker-部署)
4. [手動部署](#手動部署)
5. [配置說明](#配置說明)
6. [資料庫遷移](#資料庫遷移)
7. [驗證部署](#驗證部署)
8. [故障排除](#故障排除)

---

## 系統要求

### 硬體要求
- **CPU**: 4 核心以上
- **記憶體**: 8GB 以上
- **磁碟**: 50GB 可用空間

### 軟體要求
- **Go**: 1.21 或更高版本
- **PostgreSQL**: 15 或更高版本
- **Redis**: 7 或更高版本
- **Docker**: 24.0 或更高版本 (可選)
- **Docker Compose**: 2.20 或更高版本 (可選)

### 依賴服務
- Prometheus
- Grafana
- Loki
- Cyber-AI-Quantum
- Nginx
- RabbitMQ

---

## 環境準備

### 1. 安裝 PostgreSQL

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql-15

# 創建資料庫
sudo -u postgres createdb pandora_db
sudo -u postgres createuser pandora
sudo -u postgres psql -c "ALTER USER pandora WITH PASSWORD 'pandora123';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE pandora_db TO pandora;"
```

### 2. 安裝 Redis

```bash
# Ubuntu/Debian
sudo apt install redis-server

# 設置密碼
sudo redis-cli
> CONFIG SET requirepass "pandora123"
> exit
```

### 3. 克隆代碼

```bash
git clone https://github.com/your-org/Local_IPS-IDS.git
cd Local_IPS-IDS/Application/be
```

---

## Docker 部署

### 使用 Docker Compose

**最簡單的部署方式**

```bash
cd Application
docker-compose up -d axiom-be
```

`docker-compose.yml` 配置：

```yaml
services:
  axiom-be:
    build:
      context: ..
      dockerfile: Application/docker/axiom-be-v3.dockerfile
    container_name: axiom-be-v3
    restart: unless-stopped
    ports:
      - "3001:3001"
    environment:
      - PORT=3001
      - POSTGRES_HOST=postgres
      - POSTGRES_PORT=5432
      - POSTGRES_USER=pandora
      - POSTGRES_PASSWORD=pandora123
      - POSTGRES_DB=pandora_db
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=pandora123
      - REDIS_DB=0
      - PROMETHEUS_URL=http://prometheus:9090
      - GRAFANA_URL=http://grafana:3000
      - LOKI_URL=http://loki:3100
      - QUANTUM_URL=http://cyber-ai-quantum:8000
      - NGINX_URL=http://nginx:80
      - NGINX_CONFIG_PATH=/etc/nginx/nginx.conf
    depends_on:
      - postgres
      - redis
      - prometheus
      - loki
      - cyber-ai-quantum
    networks:
      - pandora-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3001/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### 查看日誌

```bash
docker logs -f axiom-be-v3
```

---

## 手動部署

### 1. 安裝依賴

```bash
cd Application/be
go mod download
```

### 2. 配置環境變數

```bash
cp .env.example .env
# 編輯 .env 文件，修改配置
nano .env
```

### 3. 運行資料庫遷移

```bash
make migrate
```

或手動執行：

```bash
go run cmd/server/main.go migrate
```

### 4. 構建應用

```bash
make build
```

### 5. 運行應用

```bash
./bin/axiom-backend
```

或使用 Makefile：

```bash
make run
```

---

## 配置說明

### 環境變數完整列表

| 變數名 | 默認值 | 說明 |
|--------|--------|------|
| `PORT` | 3001 | 服務端口 |
| `POSTGRES_HOST` | localhost | PostgreSQL 主機 |
| `POSTGRES_PORT` | 5432 | PostgreSQL 端口 |
| `POSTGRES_USER` | pandora | PostgreSQL 用戶 |
| `POSTGRES_PASSWORD` | pandora123 | PostgreSQL 密碼 |
| `POSTGRES_DB` | pandora_db | PostgreSQL 資料庫名 |
| `REDIS_HOST` | localhost | Redis 主機 |
| `REDIS_PORT` | 6379 | Redis 端口 |
| `REDIS_PASSWORD` | pandora123 | Redis 密碼 |
| `REDIS_DB` | 0 | Redis 資料庫編號 |
| `PROMETHEUS_URL` | http://localhost:9090 | Prometheus URL |
| `GRAFANA_URL` | http://localhost:3000 | Grafana URL |
| `LOKI_URL` | http://localhost:3100 | Loki URL |
| `QUANTUM_URL` | http://localhost:8000 | Quantum Service URL |
| `NGINX_URL` | http://localhost:80 | Nginx URL |
| `NGINX_CONFIG_PATH` | /etc/nginx/nginx.conf | Nginx 配置文件路徑 |

---

## 資料庫遷移

### 自動遷移

應用啟動時會自動執行 `AutoMigrate()`，創建所有必要的表。

### 手動遷移腳本

如需更精細的控制，使用遷移腳本：

```bash
cd database/migrations
# 執行所有遷移
./run-migrations.sh
```

### Migration 文件

位於 `database/migrations/`:
- `001_create_services_table.sql`
- `002_create_quantum_jobs_table.sql`
- `003_create_windows_logs_table.sql`
- 等...

---

## 驗證部署

### 1. 健康檢查

```bash
curl http://localhost:3001/health
```

預期響應：

```json
{
  "status": "healthy",
  "service": "axiom-backend-v3",
  "version": "3.0.0",
  "time": "..."
}
```

### 2. 測試 Prometheus API

```bash
curl -X POST http://localhost:3001/api/v2/prometheus/query \
  -H "Content-Type: application/json" \
  -d '{"query": "up"}'
```

### 3. 測試 Quantum API

```bash
curl -X POST http://localhost:3001/api/v2/quantum/qkd/generate \
  -H "Content-Type: application/json" \
  -d '{"key_length": 256, "backend": "simulator", "shots": 1024}'
```

### 4. 查看資料庫

```bash
psql -U pandora -d pandora_db -h localhost
\dt
```

應該看到 9 個表。

---

## 故障排除

### 問題 1: 無法連接資料庫

**症狀**: `Failed to connect to PostgreSQL`

**解決方案**:
1. 檢查 PostgreSQL 是否運行：`sudo systemctl status postgresql`
2. 檢查連接參數是否正確
3. 檢查防火牆設置

### 問題 2: Redis 連接失敗

**症狀**: `Failed to connect to Redis`

**解決方案**:
1. 檢查 Redis 是否運行：`sudo systemctl status redis`
2. 檢查密碼是否正確
3. 測試連接：`redis-cli -a pandora123 ping`

### 問題 3: Prometheus 查詢失敗

**症狀**: `prometheus health check failed`

**解決方案**:
1. 確認 Prometheus 服務運行：`curl http://localhost:9090/-/healthy`
2. 檢查 `PROMETHEUS_URL` 環境變數
3. 檢查網路連接

### 問題 4: 量子服務不可用

**症狀**: `quantum service health check failed`

**解決方案**:
1. 確認 cyber-ai-quantum 服務運行
2. 檢查 `QUANTUM_URL` 環境變數
3. 查看 quantum 服務日誌

---

## 性能調優

### 資料庫優化

```sql
-- 創建額外索引
CREATE INDEX idx_windows_logs_time ON windows_logs(time_created DESC);
CREATE INDEX idx_quantum_jobs_type_status ON quantum_jobs(type, status);
```

### Redis 優化

```conf
# redis.conf
maxmemory 2gb
maxmemory-policy allkeys-lru
```

### 連接池調優

在 `internal/database/db.go` 中：

```go
sqlDB.SetMaxIdleConns(20)    // 增加空閒連接
sqlDB.SetMaxOpenConns(200)   // 增加最大連接
sqlDB.SetConnMaxLifetime(time.Hour)
```

---

## 監控

### Prometheus Metrics

Axiom Backend 導出以下指標：

```
# HTTP 請求
axiom_http_requests_total{method, endpoint, status}
axiom_http_request_duration_seconds{method, endpoint}

# 資料庫
axiom_db_connections_active
axiom_db_query_duration_seconds

# 快取
axiom_cache_hits_total
axiom_cache_misses_total

# 量子作業
axiom_quantum_jobs_total{type, status}
axiom_quantum_job_duration_seconds{type}
```

### 健康檢查端點

- `/health` - 基本健康檢查
- `/api/v2/prometheus/health` - Prometheus 健康檢查
- `/api/v2/loki/health` - Loki 健康檢查
- `/api/v2/quantum/health` - Quantum 健康檢查

---

## 升級指南

### 從 V2 升級到 V3

1. **備份資料庫**

```bash
pg_dump -U pandora pandora_db > backup_$(date +%Y%m%d).sql
```

2. **停止舊版服務**

```bash
docker-compose stop axiom-be
```

3. **更新代碼**

```bash
git pull origin main
```

4. **運行遷移**

```bash
cd Application/be
go run cmd/server/main.go migrate
```

5. **啟動新版服務**

```bash
docker-compose up -d axiom-be-v3
```

6. **驗證**

```bash
curl http://localhost:3001/health
```

---

## 安全建議

1. **修改默認密碼**：修改 PostgreSQL、Redis 的默認密碼
2. **啟用 TLS**：在生產環境使用 HTTPS
3. **啟用認證**：配置 JWT 或 API Key 認證
4. **防火牆**：限制端口訪問
5. **定期備份**：設置自動備份計劃

---

## 支援

- **文檔**: [Complete Documentation](./README.md)
- **Issues**: [GitHub Issues](https://github.com/your-org/Local_IPS-IDS/issues)
- **Email**: support@example.com

---

**文檔版本**: 3.0.0  
**最後更新**: 2025-10-16


> **版本**: 3.0.0  
> **日期**: 2025-10-16

---

## 📋 目錄

1. [系統要求](#系統要求)
2. [環境準備](#環境準備)
3. [Docker 部署](#docker-部署)
4. [手動部署](#手動部署)
5. [配置說明](#配置說明)
6. [資料庫遷移](#資料庫遷移)
7. [驗證部署](#驗證部署)
8. [故障排除](#故障排除)

---

## 系統要求

### 硬體要求
- **CPU**: 4 核心以上
- **記憶體**: 8GB 以上
- **磁碟**: 50GB 可用空間

### 軟體要求
- **Go**: 1.21 或更高版本
- **PostgreSQL**: 15 或更高版本
- **Redis**: 7 或更高版本
- **Docker**: 24.0 或更高版本 (可選)
- **Docker Compose**: 2.20 或更高版本 (可選)

### 依賴服務
- Prometheus
- Grafana
- Loki
- Cyber-AI-Quantum
- Nginx
- RabbitMQ

---

## 環境準備

### 1. 安裝 PostgreSQL

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql-15

# 創建資料庫
sudo -u postgres createdb pandora_db
sudo -u postgres createuser pandora
sudo -u postgres psql -c "ALTER USER pandora WITH PASSWORD 'pandora123';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE pandora_db TO pandora;"
```

### 2. 安裝 Redis

```bash
# Ubuntu/Debian
sudo apt install redis-server

# 設置密碼
sudo redis-cli
> CONFIG SET requirepass "pandora123"
> exit
```

### 3. 克隆代碼

```bash
git clone https://github.com/your-org/Local_IPS-IDS.git
cd Local_IPS-IDS/Application/be
```

---

## Docker 部署

### 使用 Docker Compose

**最簡單的部署方式**

```bash
cd Application
docker-compose up -d axiom-be
```

`docker-compose.yml` 配置：

```yaml
services:
  axiom-be:
    build:
      context: ..
      dockerfile: Application/docker/axiom-be-v3.dockerfile
    container_name: axiom-be-v3
    restart: unless-stopped
    ports:
      - "3001:3001"
    environment:
      - PORT=3001
      - POSTGRES_HOST=postgres
      - POSTGRES_PORT=5432
      - POSTGRES_USER=pandora
      - POSTGRES_PASSWORD=pandora123
      - POSTGRES_DB=pandora_db
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=pandora123
      - REDIS_DB=0
      - PROMETHEUS_URL=http://prometheus:9090
      - GRAFANA_URL=http://grafana:3000
      - LOKI_URL=http://loki:3100
      - QUANTUM_URL=http://cyber-ai-quantum:8000
      - NGINX_URL=http://nginx:80
      - NGINX_CONFIG_PATH=/etc/nginx/nginx.conf
    depends_on:
      - postgres
      - redis
      - prometheus
      - loki
      - cyber-ai-quantum
    networks:
      - pandora-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3001/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### 查看日誌

```bash
docker logs -f axiom-be-v3
```

---

## 手動部署

### 1. 安裝依賴

```bash
cd Application/be
go mod download
```

### 2. 配置環境變數

```bash
cp .env.example .env
# 編輯 .env 文件，修改配置
nano .env
```

### 3. 運行資料庫遷移

```bash
make migrate
```

或手動執行：

```bash
go run cmd/server/main.go migrate
```

### 4. 構建應用

```bash
make build
```

### 5. 運行應用

```bash
./bin/axiom-backend
```

或使用 Makefile：

```bash
make run
```

---

## 配置說明

### 環境變數完整列表

| 變數名 | 默認值 | 說明 |
|--------|--------|------|
| `PORT` | 3001 | 服務端口 |
| `POSTGRES_HOST` | localhost | PostgreSQL 主機 |
| `POSTGRES_PORT` | 5432 | PostgreSQL 端口 |
| `POSTGRES_USER` | pandora | PostgreSQL 用戶 |
| `POSTGRES_PASSWORD` | pandora123 | PostgreSQL 密碼 |
| `POSTGRES_DB` | pandora_db | PostgreSQL 資料庫名 |
| `REDIS_HOST` | localhost | Redis 主機 |
| `REDIS_PORT` | 6379 | Redis 端口 |
| `REDIS_PASSWORD` | pandora123 | Redis 密碼 |
| `REDIS_DB` | 0 | Redis 資料庫編號 |
| `PROMETHEUS_URL` | http://localhost:9090 | Prometheus URL |
| `GRAFANA_URL` | http://localhost:3000 | Grafana URL |
| `LOKI_URL` | http://localhost:3100 | Loki URL |
| `QUANTUM_URL` | http://localhost:8000 | Quantum Service URL |
| `NGINX_URL` | http://localhost:80 | Nginx URL |
| `NGINX_CONFIG_PATH` | /etc/nginx/nginx.conf | Nginx 配置文件路徑 |

---

## 資料庫遷移

### 自動遷移

應用啟動時會自動執行 `AutoMigrate()`，創建所有必要的表。

### 手動遷移腳本

如需更精細的控制，使用遷移腳本：

```bash
cd database/migrations
# 執行所有遷移
./run-migrations.sh
```

### Migration 文件

位於 `database/migrations/`:
- `001_create_services_table.sql`
- `002_create_quantum_jobs_table.sql`
- `003_create_windows_logs_table.sql`
- 等...

---

## 驗證部署

### 1. 健康檢查

```bash
curl http://localhost:3001/health
```

預期響應：

```json
{
  "status": "healthy",
  "service": "axiom-backend-v3",
  "version": "3.0.0",
  "time": "..."
}
```

### 2. 測試 Prometheus API

```bash
curl -X POST http://localhost:3001/api/v2/prometheus/query \
  -H "Content-Type: application/json" \
  -d '{"query": "up"}'
```

### 3. 測試 Quantum API

```bash
curl -X POST http://localhost:3001/api/v2/quantum/qkd/generate \
  -H "Content-Type: application/json" \
  -d '{"key_length": 256, "backend": "simulator", "shots": 1024}'
```

### 4. 查看資料庫

```bash
psql -U pandora -d pandora_db -h localhost
\dt
```

應該看到 9 個表。

---

## 故障排除

### 問題 1: 無法連接資料庫

**症狀**: `Failed to connect to PostgreSQL`

**解決方案**:
1. 檢查 PostgreSQL 是否運行：`sudo systemctl status postgresql`
2. 檢查連接參數是否正確
3. 檢查防火牆設置

### 問題 2: Redis 連接失敗

**症狀**: `Failed to connect to Redis`

**解決方案**:
1. 檢查 Redis 是否運行：`sudo systemctl status redis`
2. 檢查密碼是否正確
3. 測試連接：`redis-cli -a pandora123 ping`

### 問題 3: Prometheus 查詢失敗

**症狀**: `prometheus health check failed`

**解決方案**:
1. 確認 Prometheus 服務運行：`curl http://localhost:9090/-/healthy`
2. 檢查 `PROMETHEUS_URL` 環境變數
3. 檢查網路連接

### 問題 4: 量子服務不可用

**症狀**: `quantum service health check failed`

**解決方案**:
1. 確認 cyber-ai-quantum 服務運行
2. 檢查 `QUANTUM_URL` 環境變數
3. 查看 quantum 服務日誌

---

## 性能調優

### 資料庫優化

```sql
-- 創建額外索引
CREATE INDEX idx_windows_logs_time ON windows_logs(time_created DESC);
CREATE INDEX idx_quantum_jobs_type_status ON quantum_jobs(type, status);
```

### Redis 優化

```conf
# redis.conf
maxmemory 2gb
maxmemory-policy allkeys-lru
```

### 連接池調優

在 `internal/database/db.go` 中：

```go
sqlDB.SetMaxIdleConns(20)    // 增加空閒連接
sqlDB.SetMaxOpenConns(200)   // 增加最大連接
sqlDB.SetConnMaxLifetime(time.Hour)
```

---

## 監控

### Prometheus Metrics

Axiom Backend 導出以下指標：

```
# HTTP 請求
axiom_http_requests_total{method, endpoint, status}
axiom_http_request_duration_seconds{method, endpoint}

# 資料庫
axiom_db_connections_active
axiom_db_query_duration_seconds

# 快取
axiom_cache_hits_total
axiom_cache_misses_total

# 量子作業
axiom_quantum_jobs_total{type, status}
axiom_quantum_job_duration_seconds{type}
```

### 健康檢查端點

- `/health` - 基本健康檢查
- `/api/v2/prometheus/health` - Prometheus 健康檢查
- `/api/v2/loki/health` - Loki 健康檢查
- `/api/v2/quantum/health` - Quantum 健康檢查

---

## 升級指南

### 從 V2 升級到 V3

1. **備份資料庫**

```bash
pg_dump -U pandora pandora_db > backup_$(date +%Y%m%d).sql
```

2. **停止舊版服務**

```bash
docker-compose stop axiom-be
```

3. **更新代碼**

```bash
git pull origin main
```

4. **運行遷移**

```bash
cd Application/be
go run cmd/server/main.go migrate
```

5. **啟動新版服務**

```bash
docker-compose up -d axiom-be-v3
```

6. **驗證**

```bash
curl http://localhost:3001/health
```

---

## 安全建議

1. **修改默認密碼**：修改 PostgreSQL、Redis 的默認密碼
2. **啟用 TLS**：在生產環境使用 HTTPS
3. **啟用認證**：配置 JWT 或 API Key 認證
4. **防火牆**：限制端口訪問
5. **定期備份**：設置自動備份計劃

---

## 支援

- **文檔**: [Complete Documentation](./README.md)
- **Issues**: [GitHub Issues](https://github.com/your-org/Local_IPS-IDS/issues)
- **Email**: support@example.com

---

**文檔版本**: 3.0.0  
**最後更新**: 2025-10-16



> **版本**: 3.0.0  
> **日期**: 2025-10-16

---

## 📋 目錄

1. [系統要求](#系統要求)
2. [環境準備](#環境準備)
3. [Docker 部署](#docker-部署)
4. [手動部署](#手動部署)
5. [配置說明](#配置說明)
6. [資料庫遷移](#資料庫遷移)
7. [驗證部署](#驗證部署)
8. [故障排除](#故障排除)

---

## 系統要求

### 硬體要求
- **CPU**: 4 核心以上
- **記憶體**: 8GB 以上
- **磁碟**: 50GB 可用空間

### 軟體要求
- **Go**: 1.21 或更高版本
- **PostgreSQL**: 15 或更高版本
- **Redis**: 7 或更高版本
- **Docker**: 24.0 或更高版本 (可選)
- **Docker Compose**: 2.20 或更高版本 (可選)

### 依賴服務
- Prometheus
- Grafana
- Loki
- Cyber-AI-Quantum
- Nginx
- RabbitMQ

---

## 環境準備

### 1. 安裝 PostgreSQL

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql-15

# 創建資料庫
sudo -u postgres createdb pandora_db
sudo -u postgres createuser pandora
sudo -u postgres psql -c "ALTER USER pandora WITH PASSWORD 'pandora123';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE pandora_db TO pandora;"
```

### 2. 安裝 Redis

```bash
# Ubuntu/Debian
sudo apt install redis-server

# 設置密碼
sudo redis-cli
> CONFIG SET requirepass "pandora123"
> exit
```

### 3. 克隆代碼

```bash
git clone https://github.com/your-org/Local_IPS-IDS.git
cd Local_IPS-IDS/Application/be
```

---

## Docker 部署

### 使用 Docker Compose

**最簡單的部署方式**

```bash
cd Application
docker-compose up -d axiom-be
```

`docker-compose.yml` 配置：

```yaml
services:
  axiom-be:
    build:
      context: ..
      dockerfile: Application/docker/axiom-be-v3.dockerfile
    container_name: axiom-be-v3
    restart: unless-stopped
    ports:
      - "3001:3001"
    environment:
      - PORT=3001
      - POSTGRES_HOST=postgres
      - POSTGRES_PORT=5432
      - POSTGRES_USER=pandora
      - POSTGRES_PASSWORD=pandora123
      - POSTGRES_DB=pandora_db
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=pandora123
      - REDIS_DB=0
      - PROMETHEUS_URL=http://prometheus:9090
      - GRAFANA_URL=http://grafana:3000
      - LOKI_URL=http://loki:3100
      - QUANTUM_URL=http://cyber-ai-quantum:8000
      - NGINX_URL=http://nginx:80
      - NGINX_CONFIG_PATH=/etc/nginx/nginx.conf
    depends_on:
      - postgres
      - redis
      - prometheus
      - loki
      - cyber-ai-quantum
    networks:
      - pandora-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:3001/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### 查看日誌

```bash
docker logs -f axiom-be-v3
```

---

## 手動部署

### 1. 安裝依賴

```bash
cd Application/be
go mod download
```

### 2. 配置環境變數

```bash
cp .env.example .env
# 編輯 .env 文件，修改配置
nano .env
```

### 3. 運行資料庫遷移

```bash
make migrate
```

或手動執行：

```bash
go run cmd/server/main.go migrate
```

### 4. 構建應用

```bash
make build
```

### 5. 運行應用

```bash
./bin/axiom-backend
```

或使用 Makefile：

```bash
make run
```

---

## 配置說明

### 環境變數完整列表

| 變數名 | 默認值 | 說明 |
|--------|--------|------|
| `PORT` | 3001 | 服務端口 |
| `POSTGRES_HOST` | localhost | PostgreSQL 主機 |
| `POSTGRES_PORT` | 5432 | PostgreSQL 端口 |
| `POSTGRES_USER` | pandora | PostgreSQL 用戶 |
| `POSTGRES_PASSWORD` | pandora123 | PostgreSQL 密碼 |
| `POSTGRES_DB` | pandora_db | PostgreSQL 資料庫名 |
| `REDIS_HOST` | localhost | Redis 主機 |
| `REDIS_PORT` | 6379 | Redis 端口 |
| `REDIS_PASSWORD` | pandora123 | Redis 密碼 |
| `REDIS_DB` | 0 | Redis 資料庫編號 |
| `PROMETHEUS_URL` | http://localhost:9090 | Prometheus URL |
| `GRAFANA_URL` | http://localhost:3000 | Grafana URL |
| `LOKI_URL` | http://localhost:3100 | Loki URL |
| `QUANTUM_URL` | http://localhost:8000 | Quantum Service URL |
| `NGINX_URL` | http://localhost:80 | Nginx URL |
| `NGINX_CONFIG_PATH` | /etc/nginx/nginx.conf | Nginx 配置文件路徑 |

---

## 資料庫遷移

### 自動遷移

應用啟動時會自動執行 `AutoMigrate()`，創建所有必要的表。

### 手動遷移腳本

如需更精細的控制，使用遷移腳本：

```bash
cd database/migrations
# 執行所有遷移
./run-migrations.sh
```

### Migration 文件

位於 `database/migrations/`:
- `001_create_services_table.sql`
- `002_create_quantum_jobs_table.sql`
- `003_create_windows_logs_table.sql`
- 等...

---

## 驗證部署

### 1. 健康檢查

```bash
curl http://localhost:3001/health
```

預期響應：

```json
{
  "status": "healthy",
  "service": "axiom-backend-v3",
  "version": "3.0.0",
  "time": "..."
}
```

### 2. 測試 Prometheus API

```bash
curl -X POST http://localhost:3001/api/v2/prometheus/query \
  -H "Content-Type: application/json" \
  -d '{"query": "up"}'
```

### 3. 測試 Quantum API

```bash
curl -X POST http://localhost:3001/api/v2/quantum/qkd/generate \
  -H "Content-Type: application/json" \
  -d '{"key_length": 256, "backend": "simulator", "shots": 1024}'
```

### 4. 查看資料庫

```bash
psql -U pandora -d pandora_db -h localhost
\dt
```

應該看到 9 個表。

---

## 故障排除

### 問題 1: 無法連接資料庫

**症狀**: `Failed to connect to PostgreSQL`

**解決方案**:
1. 檢查 PostgreSQL 是否運行：`sudo systemctl status postgresql`
2. 檢查連接參數是否正確
3. 檢查防火牆設置

### 問題 2: Redis 連接失敗

**症狀**: `Failed to connect to Redis`

**解決方案**:
1. 檢查 Redis 是否運行：`sudo systemctl status redis`
2. 檢查密碼是否正確
3. 測試連接：`redis-cli -a pandora123 ping`

### 問題 3: Prometheus 查詢失敗

**症狀**: `prometheus health check failed`

**解決方案**:
1. 確認 Prometheus 服務運行：`curl http://localhost:9090/-/healthy`
2. 檢查 `PROMETHEUS_URL` 環境變數
3. 檢查網路連接

### 問題 4: 量子服務不可用

**症狀**: `quantum service health check failed`

**解決方案**:
1. 確認 cyber-ai-quantum 服務運行
2. 檢查 `QUANTUM_URL` 環境變數
3. 查看 quantum 服務日誌

---

## 性能調優

### 資料庫優化

```sql
-- 創建額外索引
CREATE INDEX idx_windows_logs_time ON windows_logs(time_created DESC);
CREATE INDEX idx_quantum_jobs_type_status ON quantum_jobs(type, status);
```

### Redis 優化

```conf
# redis.conf
maxmemory 2gb
maxmemory-policy allkeys-lru
```

### 連接池調優

在 `internal/database/db.go` 中：

```go
sqlDB.SetMaxIdleConns(20)    // 增加空閒連接
sqlDB.SetMaxOpenConns(200)   // 增加最大連接
sqlDB.SetConnMaxLifetime(time.Hour)
```

---

## 監控

### Prometheus Metrics

Axiom Backend 導出以下指標：

```
# HTTP 請求
axiom_http_requests_total{method, endpoint, status}
axiom_http_request_duration_seconds{method, endpoint}

# 資料庫
axiom_db_connections_active
axiom_db_query_duration_seconds

# 快取
axiom_cache_hits_total
axiom_cache_misses_total

# 量子作業
axiom_quantum_jobs_total{type, status}
axiom_quantum_job_duration_seconds{type}
```

### 健康檢查端點

- `/health` - 基本健康檢查
- `/api/v2/prometheus/health` - Prometheus 健康檢查
- `/api/v2/loki/health` - Loki 健康檢查
- `/api/v2/quantum/health` - Quantum 健康檢查

---

## 升級指南

### 從 V2 升級到 V3

1. **備份資料庫**

```bash
pg_dump -U pandora pandora_db > backup_$(date +%Y%m%d).sql
```

2. **停止舊版服務**

```bash
docker-compose stop axiom-be
```

3. **更新代碼**

```bash
git pull origin main
```

4. **運行遷移**

```bash
cd Application/be
go run cmd/server/main.go migrate
```

5. **啟動新版服務**

```bash
docker-compose up -d axiom-be-v3
```

6. **驗證**

```bash
curl http://localhost:3001/health
```

---

## 安全建議

1. **修改默認密碼**：修改 PostgreSQL、Redis 的默認密碼
2. **啟用 TLS**：在生產環境使用 HTTPS
3. **啟用認證**：配置 JWT 或 API Key 認證
4. **防火牆**：限制端口訪問
5. **定期備份**：設置自動備份計劃

---

## 支援

- **文檔**: [Complete Documentation](./README.md)
- **Issues**: [GitHub Issues](https://github.com/your-org/Local_IPS-IDS/issues)
- **Email**: support@example.com

---

**文檔版本**: 3.0.0  
**最後更新**: 2025-10-16

