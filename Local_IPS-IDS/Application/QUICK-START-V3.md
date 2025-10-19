# Axiom Backend V3 - 快速啟動指南

> **5 分鐘快速測試**

---

## 🚀 方法 1: 使用腳本 (最簡單)

### Windows
```powershell
cd Application
.\build-axiom-be-v3.ps1
```

### Linux/Mac
```bash
cd Application
chmod +x build-axiom-be-v3.sh
./build-axiom-be-v3.sh
```

---

## 🐳 方法 2: Docker Compose (推薦)

### 完整啟動 (所有服務)

```bash
cd Application

# 1. 啟動所有服務
docker-compose up -d

# 2. 等待 30 秒讓服務啟動
# (Windows)
timeout /t 30
# (Linux/Mac)
sleep 30

# 3. 執行資料庫 Migration
docker-compose exec -T postgres psql -U pandora -d pandora_db < ../database/migrations/001_initial_schema.sql
docker-compose exec -T postgres psql -U pandora -d pandora_db < ../database/migrations/002_agent_and_compliance_schema.sql

# 4. 測試 API
curl http://localhost:3001/health
```

### 只啟動 Axiom Backend 及依賴

```bash
cd Application

# 啟動必要服務
docker-compose up -d postgres redis prometheus loki cyber-ai-quantum axiom-be

# 查看日誌
docker-compose logs -f axiom-be
```

---

## 📝 方法 3: 手動 Docker 構建

```bash
cd C:\Users\dennis.lee\Documents\GitHub\Local_IPS-IDS

# 構建鏡像
docker build -f Application/docker/axiom-be.dockerfile -t axiom-backend:v3 .

# 啟動容器（確保 PostgreSQL 和 Redis 已運行）
docker run -d \
  --name axiom-be-test \
  -p 3001:3001 \
  -e POSTGRES_HOST=host.docker.internal \
  -e POSTGRES_PORT=5432 \
  -e POSTGRES_USER=pandora \
  -e POSTGRES_PASSWORD=pandora123 \
  -e POSTGRES_DB=pandora_db \
  -e REDIS_HOST=host.docker.internal \
  -e REDIS_PORT=6379 \
  -e REDIS_PASSWORD=pandora123 \
  axiom-backend:v3

# 查看日誌
docker logs -f axiom-be-test
```

---

## 🧪 快速測試

### 測試 API (PowerShell)

```powershell
# 健康檢查
Invoke-RestMethod http://localhost:3001/health

# Agent 註冊
$body = @{
    mode = "internal"
    hostname = "test-server"
    ip_address = "127.0.0.1"
    capabilities = @("windows_logs")
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:3001/api/v2/agent/register -Method Post -Body $body -ContentType "application/json"

# PII 檢測
$piiBody = @{
    text = "Contact: test@example.com"
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:3001/api/v2/compliance/pii/detect -Method Post -Body $piiBody -ContentType "application/json"
```

### 測試 API (Curl)

```bash
# 健康檢查
curl http://localhost:3001/health

# Agent 註冊
curl -X POST http://localhost:3001/api/v2/agent/register \
  -H "Content-Type: application/json" \
  -d '{"mode":"internal","hostname":"test","ip_address":"127.0.0.1","capabilities":["windows_logs"]}'

# PII 檢測
curl -X POST http://localhost:3001/api/v2/compliance/pii/detect \
  -H "Content-Type: application/json" \
  -d '{"text":"Email: test@example.com, Card: 4532-1234-5678-9010"}'

# 儲存統計
curl http://localhost:3001/api/v2/storage/tiers/stats
```

---

## 🔍 常用指令

### Docker Compose

```bash
# 查看狀態
docker-compose ps

# 查看日誌
docker-compose logs axiom-be
docker-compose logs -f axiom-be  # 實時

# 重啟服務
docker-compose restart axiom-be

# 停止服務
docker-compose stop axiom-be

# 刪除並重建
docker-compose up --build --force-recreate axiom-be
```

### Docker

```bash
# 查看容器
docker ps | grep axiom

# 查看日誌
docker logs axiom-be-v3
docker logs -f axiom-be-v3  # 實時

# 進入容器
docker exec -it axiom-be-v3 sh

# 停止容器
docker stop axiom-be-v3

# 刪除容器
docker rm axiom-be-v3
```

---

## 📊 服務端點

- **健康檢查**: `http://localhost:3001/health`
- **API v2**: `http://localhost:3001/api/v2/`
- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000`
- **Loki**: `http://localhost:3100`

---

## 🛠️ 故障排除

### 容器無法啟動

```bash
# 查看詳細日誌
docker logs axiom-be-v3

# 檢查網路
docker network inspect pandora-network

# 檢查環境變量
docker exec axiom-be-v3 env
```

### API 無響應

```bash
# 檢查服務是否運行
docker ps | grep axiom-be

# 測試連接
curl -v http://localhost:3001/health

# 查看端口映射
docker port axiom-be-v3
```

### 資料庫連接失敗

```bash
# 測試 PostgreSQL
docker-compose exec postgres psql -U pandora -d pandora_db -c "SELECT version();"

# 測試 Redis
docker-compose exec redis redis-cli -a pandora123 ping
```

---

**完整文檔**: `Application/DOCKER-BUILD-TEST-GUIDE.md`

