# Axiom Backend V3 - Docker 本地構建和測試指南

> **版本**: 3.1.0  
> **日期**: 2025-10-16

---

## 🚀 快速開始

### 方法 1: 使用自動化腳本 (推薦)

#### Windows (PowerShell)
```powershell
cd Application
.\build-axiom-be-v3.ps1
```

#### Linux/Mac (Bash)
```bash
cd Application
chmod +x build-axiom-be-v3.sh
./build-axiom-be-v3.sh
```

---

### 方法 2: 手動 Docker 指令

#### Step 1: 構建鏡像

```bash
# 進入專案根目錄
cd C:\Users\dennis.lee\Documents\GitHub\Local_IPS-IDS

# 構建 Axiom Backend V3 鏡像
docker build \
  -f Application/docker/axiom-be.dockerfile \
  -t axiom-backend:v3.1.0 \
  -t axiom-backend:latest \
  .
```

#### Step 2: 獨立測試容器

```bash
# 啟動測試容器（假設 PostgreSQL 和 Redis 已在本地運行）
docker run -d \
  --name axiom-be-v3-test \
  -p 3001:3001 \
  -e POSTGRES_HOST=host.docker.internal \
  -e POSTGRES_PORT=5432 \
  -e POSTGRES_USER=pandora \
  -e POSTGRES_PASSWORD=pandora123 \
  -e POSTGRES_DB=pandora_db \
  -e REDIS_HOST=host.docker.internal \
  -e REDIS_PORT=6379 \
  -e REDIS_PASSWORD=pandora123 \
  -e PROMETHEUS_URL=http://host.docker.internal:9090 \
  -e LOKI_URL=http://host.docker.internal:3100 \
  -e QUANTUM_URL=http://host.docker.internal:8000 \
  axiom-backend:v3.1.0
```

#### Step 3: 查看日誌

```bash
# 查看容器日誌
docker logs axiom-be-v3-test

# 實時查看日誌
docker logs -f axiom-be-v3-test
```

#### Step 4: 測試 API

```bash
# 健康檢查
curl http://localhost:3001/health

# Agent 註冊
curl -X POST http://localhost:3001/api/v2/agent/register \
  -H "Content-Type: application/json" \
  -d '{
    "mode": "internal",
    "hostname": "test-server",
    "ip_address": "127.0.0.1",
    "capabilities": ["windows_logs"]
  }'

# PII 檢測
curl -X POST http://localhost:3001/api/v2/compliance/pii/detect \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Contact: test@example.com"
  }'

# 儲存統計
curl http://localhost:3001/api/v2/storage/tiers/stats
```

---

### 方法 3: 使用 Docker Compose (完整環境)

#### Step 1: 啟動完整環境

```bash
# 進入 Application 目錄
cd Application

# 啟動所有服務（包括 PostgreSQL, Redis, Prometheus, Loki 等）
docker-compose up -d

# 或只啟動 Axiom Backend 及其依賴
docker-compose up -d postgres redis prometheus loki cyber-ai-quantum axiom-be
```

#### Step 2: 查看服務狀態

```bash
# 查看所有服務狀態
docker-compose ps

# 查看 Axiom BE 日誌
docker-compose logs -f axiom-be
```

#### Step 3: 執行資料庫 Migration

```bash
# 進入 postgres 容器
docker-compose exec postgres psql -U pandora -d pandora_db

# 或從外部執行
docker-compose exec -T postgres psql -U pandora -d pandora_db < ../database/migrations/001_initial_schema.sql
docker-compose exec -T postgres psql -U pandora -d pandora_db < ../database/migrations/002_agent_and_compliance_schema.sql
```

#### Step 4: 測試服務

```bash
# 等待服務啟動（約 30 秒）
sleep 30

# 測試健康檢查
curl http://localhost:3001/health

# 測試 Agent 列表
curl http://localhost:3001/api/v2/agent/list
```

#### Step 5: 停止服務

```bash
# 停止所有服務
docker-compose down

# 停止並刪除 volumes（清理所有數據）
docker-compose down -v
```

---

## 🧪 測試腳本

### Windows PowerShell 測試

```powershell
# 1. 健康檢查
$health = Invoke-RestMethod -Uri "http://localhost:3001/health"
Write-Host "健康狀態: $($health.status)"

# 2. Agent 註冊
$agentBody = @{
    mode = "internal"
    hostname = "test-server"
    ip_address = "127.0.0.1"
    capabilities = @("windows_logs", "compliance_scan")
} | ConvertTo-Json

$agentResponse = Invoke-RestMethod -Uri "http://localhost:3001/api/v2/agent/register" -Method Post -Body $agentBody -ContentType "application/json"
Write-Host "Agent ID: $($agentResponse.data.agent_id)"

# 3. PII 檢測
$piiBody = @{
    text = "Email: john@example.com, Card: 4532-1234-5678-9010, SSN: 123-45-6789"
} | ConvertTo-Json

$piiResponse = Invoke-RestMethod -Uri "http://localhost:3001/api/v2/compliance/pii/detect" -Method Post -Body $piiBody -ContentType "application/json"
Write-Host "PII 發現: $($piiResponse.data.matches.Count) 個"
$piiResponse.data.matches | ForEach-Object { Write-Host "  - $($_.type): $($_.masked)" }

# 4. 資料匿名化
$anonBody = @{
    text = "User: john@example.com, IP: 192.168.1.100"
    method = "hash"
} | ConvertTo-Json

$anonResponse = Invoke-RestMethod -Uri "http://localhost:3001/api/v2/compliance/pii/anonymize" -Method Post -Body $anonBody -ContentType "application/json"
Write-Host "匿名化結果: $($anonResponse.data.anonymized_text)"

# 5. 儲存統計
$storageResponse = Invoke-RestMethod -Uri "http://localhost:3001/api/v2/storage/tiers/stats"
Write-Host "儲存統計: $($storageResponse.data | ConvertTo-Json)"
```

### Linux/Mac Bash 測試

```bash
#!/bin/bash

# 1. 健康檢查
echo "測試健康檢查..."
curl -s http://localhost:3001/health | jq .

# 2. Agent 註冊
echo "測試 Agent 註冊..."
curl -s -X POST http://localhost:3001/api/v2/agent/register \
  -H "Content-Type: application/json" \
  -d '{
    "mode": "internal",
    "hostname": "test-server",
    "ip_address": "127.0.0.1",
    "capabilities": ["windows_logs"]
  }' | jq .

# 3. PII 檢測
echo "測試 PII 檢測..."
curl -s -X POST http://localhost:3001/api/v2/compliance/pii/detect \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Contact: john@example.com, Card: 4532-1234-5678-9010"
  }' | jq .

# 4. GDPR 刪除請求
echo "測試 GDPR 刪除請求..."
curl -s -X POST http://localhost:3001/api/v2/compliance/gdpr/deletion-request \
  -H "Content-Type: application/json" \
  -d '{
    "subject_id": "john@example.com",
    "requested_by": "dpo@company.com",
    "notes": "User requested deletion"
  }' | jq .

# 5. 儲存統計
echo "測試儲存統計..."
curl -s http://localhost:3001/api/v2/storage/tiers/stats | jq .
```

---

## 📋 本地構建步驟詳解

### 前置條件

1. **確保服務正在運行**:
```bash
# 檢查 PostgreSQL
docker ps | grep postgres

# 檢查 Redis
docker ps | grep redis

# 如果沒有運行，啟動它們
docker-compose up -d postgres redis prometheus loki cyber-ai-quantum
```

2. **執行資料庫 Migration**:
```bash
# 方法 1: 使用 docker-compose exec
docker-compose exec -T postgres psql -U pandora -d pandora_db < ../database/migrations/001_initial_schema.sql
docker-compose exec -T postgres psql -U pandora -d pandora_db < ../database/migrations/002_agent_and_compliance_schema.sql

# 方法 2: 使用本地 psql (如果已安裝)
psql -h localhost -U pandora -d pandora_db -f ../database/migrations/001_initial_schema.sql
psql -h localhost -U pandora -d pandora_db -f ../database/migrations/002_agent_and_compliance_schema.sql
```

### 構建選項

#### 選項 1: 只構建 Axiom Backend

```bash
cd Application
docker build -f docker/axiom-be.dockerfile -t axiom-backend:v3.1.0 ..
```

#### 選項 2: 構建並立即運行

```bash
cd Application
docker-compose up --build axiom-be
```

#### 選項 3: 重建（不使用快取）

```bash
cd Application
docker build --no-cache -f docker/axiom-be.dockerfile -t axiom-backend:v3.1.0 ..
```

---

## 🔍 故障排除

### 問題 1: 構建失敗

```bash
# 檢查 Go 版本
go version

# 清理 Docker 快取
docker system prune -a

# 重新構建
docker build --no-cache -f docker/axiom-be.dockerfile -t axiom-backend:v3.1.0 ..
```

### 問題 2: 容器無法啟動

```bash
# 查看詳細日誌
docker logs axiom-be-v3

# 進入容器調試
docker exec -it axiom-be-v3 sh

# 檢查環境變量
docker exec axiom-be-v3 env
```

### 問題 3: 無法連接資料庫

```bash
# 檢查網路
docker network ls
docker network inspect pandora-network

# 測試資料庫連接
docker-compose exec postgres psql -U pandora -d pandora_db -c "SELECT version();"

# 檢查 Redis 連接
docker-compose exec redis redis-cli -a pandora123 ping
```

### 問題 4: API 無響應

```bash
# 檢查端口映射
docker port axiom-be-v3

# 檢查服務是否在監聽
docker exec axiom-be-v3 netstat -tuln

# 重啟服務
docker-compose restart axiom-be
```

---

## 📊 性能測試

### 基準測試

```bash
# 使用 Apache Bench
ab -n 1000 -c 10 http://localhost:3001/health

# 使用 wrk
wrk -t4 -c100 -d30s http://localhost:3001/health
```

### 負載測試

```bash
# 批量 Agent 註冊
for i in {1..100}; do
  curl -s -X POST http://localhost:3001/api/v2/agent/register \
    -H "Content-Type: application/json" \
    -d "{\"mode\":\"internal\",\"hostname\":\"test-$i\",\"ip_address\":\"127.0.0.$i\",\"capabilities\":[\"windows_logs\"]}" &
done
wait

# 查看註冊的 Agents
curl http://localhost:3001/api/v2/agent/list | jq '.data | length'
```

---

## 🎯 完整測試清單

### 基礎功能測試

- [ ] 健康檢查 `GET /health`
- [ ] Prometheus 查詢 `POST /api/v2/prometheus/query`
- [ ] Loki 日誌查詢 `GET /api/v2/loki/query`
- [ ] Quantum 任務 `POST /api/v2/quantum/qkd/generate`

### Phase 11: Agent 測試

- [ ] Agent 註冊 (Internal) `POST /api/v2/agent/register`
- [ ] Agent 註冊 (External) `POST /api/v2/agent/register`
- [ ] Agent 心跳 `POST /api/v2/agent/heartbeat`
- [ ] Agent 列表 `GET /api/v2/agent/list`
- [ ] 資產發現 `POST /api/v2/agent/practical/discover-assets`
- [ ] 合規檢查 `POST /api/v2/agent/practical/check-compliance`
- [ ] 遠端執行 `POST /api/v2/agent/practical/execute-command`

### Phase 12: Storage 測試

- [ ] 儲存統計 `GET /api/v2/storage/tiers/stats`
- [ ] 手動轉移 `POST /api/v2/storage/tier/transfer`

### Phase 13: Compliance 測試

- [ ] PII 檢測 `POST /api/v2/compliance/pii/detect`
- [ ] 資料匿名化 `POST /api/v2/compliance/pii/anonymize`
- [ ] 反假名化 `POST /api/v2/compliance/pii/depseudonymize`
- [ ] GDPR 刪除請求 `POST /api/v2/compliance/gdpr/deletion-request`
- [ ] GDPR 審批 `POST /api/v2/compliance/gdpr/deletion-request/{id}/approve`
- [ ] GDPR 執行 `POST /api/v2/compliance/gdpr/deletion-request/{id}/execute`
- [ ] GDPR 驗證 `GET /api/v2/compliance/gdpr/deletion-request/{id}/verify`
- [ ] 資料匯出 `POST /api/v2/compliance/gdpr/data-export`

### 創新功能測試

- [ ] 時間旅行快照 `POST /api/v2/time-travel/snapshot/create`
- [ ] 風險評分 `POST /api/v2/adaptive-security/risk/calculate`
- [ ] 自癒修復 `POST /api/v2/combined/self-healing/remediate`
- [ ] API 健康評分 `GET /api/v2/governance/api-health/{apiPath}`
- [ ] 技術債務掃描 `POST /api/v2/tech-debt/scan`

---

## 🐳 Docker Compose 完整指令

### 啟動服務

```bash
# 啟動所有服務
docker-compose up -d

# 啟動特定服務組
docker-compose up -d postgres redis axiom-be

# 啟動並查看日誌
docker-compose up axiom-be
```

### 管理服務

```bash
# 查看狀態
docker-compose ps

# 重啟服務
docker-compose restart axiom-be

# 停止服務
docker-compose stop axiom-be

# 刪除服務
docker-compose down

# 刪除服務和數據
docker-compose down -v
```

### 查看日誌

```bash
# 查看所有日誌
docker-compose logs

# 查看特定服務
docker-compose logs axiom-be

# 實時日誌
docker-compose logs -f axiom-be

# 最近 100 行
docker-compose logs --tail=100 axiom-be
```

### 執行命令

```bash
# 進入容器
docker-compose exec axiom-be sh

# 執行命令
docker-compose exec axiom-be curl http://localhost:3001/health
```

---

## 🔧 開發模式

### 本地開發（不使用 Docker）

```bash
# 1. 確保 PostgreSQL 和 Redis 正在運行
docker-compose up -d postgres redis prometheus loki cyber-ai-quantum

# 2. 設置環境變量
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export REDIS_HOST=localhost
export PROMETHEUS_URL=http://localhost:9090
export LOKI_URL=http://localhost:3100
export QUANTUM_URL=http://localhost:8000

# 3. 執行 Migration
psql -h localhost -U pandora -d pandora_db -f database/migrations/001_initial_schema.sql
psql -h localhost -U pandora -d pandora_db -f database/migrations/002_agent_and_compliance_schema.sql

# 4. 運行應用
cd Application/be
go run cmd/server/main.go
```

### 熱重載開發

```bash
# 安裝 air
go install github.com/cosmtrek/air@latest

# 運行熱重載
cd Application/be
air
```

---

## 📈 監控

### 容器資源使用

```bash
# 查看資源使用
docker stats axiom-be-v3

# 查看 top processes
docker top axiom-be-v3
```

### 應用監控

```bash
# Prometheus 指標（如果已實現）
curl http://localhost:3001/metrics

# 健康檢查詳情
curl http://localhost:3001/health | jq .
```

---

## 🚀 生產環境構建

### 構建優化鏡像

```bash
# 多階段構建 + 優化
docker build \
  -f Application/docker/axiom-be.dockerfile \
  --target builder \
  --build-arg GO_VERSION=1.21 \
  -t axiom-backend:v3.1.0-prod \
  .
```

### 推送到 Registry

```bash
# 標記鏡像
docker tag axiom-backend:v3.1.0 your-registry.com/axiom-backend:v3.1.0

# 推送
docker push your-registry.com/axiom-backend:v3.1.0
```

---

## 📝 注意事項

1. **資料庫連接**: 確保 PostgreSQL 和 Redis 已啟動並可訪問
2. **Migration**: 首次啟動前必須執行資料庫 Migration
3. **環境變量**: 檢查所有必需的環境變量是否已設置
4. **端口衝突**: 確保 3001 端口未被佔用
5. **網路**: Docker Compose 會自動創建 `pandora-network`

---

## 🎯 快速測試流程

### 完整測試流程（5分鐘）

```bash
# 1. 啟動依賴服務 (2分鐘)
cd Application
docker-compose up -d postgres redis prometheus loki cyber-ai-quantum

# 2. 執行 Migration (30秒)
sleep 30
docker-compose exec -T postgres psql -U pandora -d pandora_db < ../database/migrations/001_initial_schema.sql
docker-compose exec -T postgres psql -U pandora -d pandora_db < ../database/migrations/002_agent_and_compliance_schema.sql

# 3. 構建並啟動 Axiom BE (2分鐘)
docker-compose up --build -d axiom-be

# 4. 等待啟動 (30秒)
sleep 30

# 5. 測試 API (1分鐘)
curl http://localhost:3001/health
curl -X POST http://localhost:3001/api/v2/agent/register \
  -H "Content-Type: application/json" \
  -d '{"mode":"internal","hostname":"test","ip_address":"127.0.0.1","capabilities":["windows_logs"]}'
```

---

**文檔版本**: 3.1.0  
**最後更新**: 2025-10-16  
**參考**: [最終報告](../docs/AXIOM-BACKEND-V3-FINAL-REPORT.md)

