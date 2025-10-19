# Pandora Box Console IDS-IPS - PaaS 多平台部署指南

[![Deploy Status](https://img.shields.io/badge/Deploy-PaaS-success)](https://github.com)
[![Railway](https://img.shields.io/badge/Railway-PostgreSQL-brightgreen)](https://railway.app)
[![Render](https://img.shields.io/badge/Render-Redis%20%2B%20Nginx-blue)](https://render.com)
[![Koyeb](https://img.shields.io/badge/Koyeb-Agent-orange)](https://koyeb.com)
[![Patr.io](https://img.shields.io/badge/Patr.io-UI-purple)](https://patr.cloud)
[![Fly.io](https://img.shields.io/badge/Fly.io-Monitoring-lightblue)](https://fly.io)

本文件說明如何將 Pandora Box Console IDS-IPS 系統部署到多個免費 PaaS 平台，實現零成本的雲端微服務架構。

## 📊 架構概覽

```
┌─────────────────────────────────────────────────────────────────────┐
│                    Pandora Box Console IDS-IPS                      │
│                      PaaS 多平台架構                                 │
└─────────────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              │               │               │
              ▼               ▼               ▼
       ┌──────────┐    ┌──────────┐    ┌──────────┐
       │ Railway  │    │  Render  │    │  Koyeb   │
       │PostgreSQL│    │Redis+Nginx│   │  Agent   │
       └──────────┘    └──────────┘    └──────────┘
              │               │               │
              └───────────────┼───────────────┘
                              │
              ┌───────────────┼───────────────┐
              │               │               │
              ▼               ▼               ▼
       ┌──────────┐    ┌──────────────────────────┐
       │ Patr.io  │    │       Fly.io             │
       │   UI     │    │ Prometheus+Loki+Grafana  │
       └──────────┘    └──────────────────────────┘
```

## 🎯 部署策略

| 微服務 | 部署平台 | 部署方式 | 理由 |product|
|--------|---------|---------|------|------|
| **PostgreSQL** | Railway.app | Managed Database | 免費託管 PostgreSQL，穩定可靠 |ballast.proxy.rlwy.net:53524|
| **Redis** | Render | Managed Redis | 免費託管 Redis，低延遲快取 |https://redis-7-2-11-alpine3-21.onrender.com|
| **Nginx** | Render | Web Service | 反向代理，統一流量入口 |https://nginx-stable-perl-boqt.onrender.com|
| **pandora-agent** | Koyeb | Docker Container | 2 個永不休眠容器，24/7 運行 ||
| **promtail** | Koyeb | Sidecar | 與 Agent 共同部署，收集日誌 ||
| **axiom-ui** | Patr.io | Container | UI 伺服器，用戶介面流量分離 ||
| **Prometheus** | Fly.io | Docker App | 永久儲存空間，監控指標收集 |https://pandora-monitoring.fly.dev/prometheus|
| **Loki** | Fly.io | Docker App | 永久儲存空間，日誌聚合 |https://pandora-monitoring.fly.dev/loki|
| **Grafana** | Fly.io | Docker App | 視覺化儀表板，隨時可用 | https://pandora-monitoring.fly.dev|
| **AlertManager** | Fly.io | Docker App | 告警系統，全時運行 | https://pandora-monitoring.fly.dev/alertmanager|
| **node-exporter** | - | 不部署 | PaaS 環境無法監控底層主機 ||

## 🚀 快速開始

### 前置需求

1. **帳號註冊**
   - [Railway.app](https://railway.app)
   - [Render](https://render.com)
   - [Koyeb](https://koyeb.com)
   - [Patr.io](https://patr.cloud)
   - [Fly.io](https://fly.io)

2. **安裝 CLI 工具**
   ```bash
   # Railway CLI
   npm install -g @railway/cli
   
   # Render CLI (可選)
   # 主要使用 Web Dashboard
   
   # Koyeb CLI
   curl -fsSL https://cli.koyeb.com/install.sh | bash
   
   # Patr.io CLI (可選)
   # 主要使用 Web Dashboard
   
   # Fly.io CLI
   curl -L https://fly.io/install.sh | sh
   ```

3. **其他工具**
   - Docker 20.10+
   - Git
   - curl
   - jq (可選)

### 一鍵部署

```bash
# 1. 複製環境變數範本
cp env.paas.example .env.paas

# 2. 編輯 .env.paas，填入各平台的 API Token 和 URL

# 3. 執行部署腳本
chmod +x scripts/deploy-paas.sh
./scripts/deploy-paas.sh

# 4. 驗證部署
chmod +x scripts/verify-paas-deployment.sh
./scripts/verify-paas-deployment.sh
```

## 📝 詳細部署步驟

### 1️⃣ Railway - PostgreSQL 資料庫

#### 使用 CLI 部署

```bash
# 登入 Railway
railway login

# 建立新專案
railway init --name pandora-postgresql

# 添加 PostgreSQL 資料庫
railway add --database postgresql

# 上傳初始化 SQL
railway up configs/postgres/init.sql
```

#### 使用 Web Dashboard

1. 前往 [Railway Dashboard](https://railway.app/dashboard)
2. 點擊 "New Project"
3. 選擇 "Provision PostgreSQL"
4. 專案建立後，進入 PostgreSQL 服務
5. 在 "Variables" 標籤中複製 `DATABASE_URL`
6. 在 "Data" 標籤中執行 `configs/postgres/init.sql`

#### 環境變數設定

```bash
# 複製以下變數到 .env.paas
RAILWAY_DATABASE_URL=postgresql://postgres:xxx@containers-us-west-xxx.railway.app:5432/railway
POSTGRES_HOST=containers-us-west-xxx.railway.app
POSTGRES_PORT=5432
POSTGRES_DB=railway
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your-password
```

### 2️⃣ Render - Redis + Nginx

#### 部署 Redis

1. 前往 [Render Dashboard](https://dashboard.render.com)
2. 點擊 "New +" → "Redis"
3. 選擇免費方案 (Free)
4. 建立後複製 "Internal Redis URL" 和 "External Redis URL"

#### 部署 Nginx

1. 在 Render Dashboard 點擊 "New +" → "Web Service"
2. 連接 GitHub Repository
3. 設定如下：
   - **Name**: `pandora-nginx`
   - **Environment**: `Docker`
   - **Dockerfile Path**: `Dockerfile.nginx`
   - **Plan**: `Free`

4. 設定環境變數：
   ```
   AXIOM_UI_URL=https://axiom-ui-xxx.patr.cloud
   GRAFANA_URL=https://pandora-monitoring.fly.dev
   PROMETHEUS_URL=https://pandora-monitoring.fly.dev/prometheus
   LOKI_URL=https://pandora-monitoring.fly.dev/loki
   ALERTMANAGER_URL=https://pandora-monitoring.fly.dev/alertmanager
   ```

5. 部署完成後，複製服務 URL

### 3️⃣ Koyeb - Pandora Agent + Promtail

#### 使用 CLI 部署

```bash
# 登入 Koyeb
koyeb login

# 建置 Docker 映像
docker build -f Dockerfile.agent.koyeb -t your-docker-username/pandora-agent:latest .
docker push your-docker-username/pandora-agent:latest

# 建立應用
koyeb app create pandora-agent

# 部署服務
koyeb service create pandora-agent \
  --app pandora-agent \
  --docker your-docker-username/pandora-agent:latest \
  --ports 8080:http \
  --routes /:8080 \
  --regions fra \
  --instance-type nano \
  --replicas 2 \
  --env DATABASE_URL=$RAILWAY_DATABASE_URL \
  --env REDIS_URL=$RENDER_REDIS_URL \
  --env PROMETHEUS_URL=$PROMETHEUS_URL \
  --env LOKI_URL=$LOKI_URL \
  --env GRAFANA_URL=$GRAFANA_URL
```

#### 使用 Web Dashboard（推薦）

1. **前往 [Koyeb Dashboard](https://app.koyeb.com)**
2. **點擊 "Create App"**
3. **選擇部署來源**
   - 選擇 "GitHub"
   - 選擇 Repository: `pandora_box_console_IDS-IPS`
   - 選擇 Branch: `main`

4. **設定建置配置** ⚠️ 重要
   - **Builder**: Docker
   - **Dockerfile path**: `Dockerfile.agent.koyeb` （必須明確指定！）
   - **Build context**: `/` (root)

5. **設定部署配置**
   - **App name**: `pandora-agent`
   - **Service name**: `pandora-agent`
   - **Region**: Europe (Frankfurt) - `fra`
   - **Instance type**: Nano
   - **Instances**: 2
   - **Port**: 8080

6. **新增環境變數**（見上述 CLI 範例）

7. **設定健康檢查**
   - **Type**: HTTP
   - **Path**: `/health`
   - **Port**: 8080

8. **點擊 "Deploy"**

> **注意**: 如果遇到 Dockerfile 找不到的錯誤，請確認在 "Dockerfile path" 欄位中正確填入 `Dockerfile.agent.koyeb`（不要包含 `./` 前綴）。詳細故障排除請參閱 [KOYEB-DEPLOYMENT-GUIDE.md](KOYEB-DEPLOYMENT-GUIDE.md)

### 4️⃣ Patr.io - Axiom UI

#### 使用 Web Dashboard

1. 前往 [Patr.io Dashboard](https://patr.cloud)
2. 點擊 "New Deployment"
3. 連接 GitHub Repository
4. 設定：
   - **Name**: `axiom-ui`
   - **Dockerfile**: `Dockerfile.ui.patr`
   - **Port**: 3001
   - **Resources**: 512Mi RAM, 0.5 CPU

5. 設定環境變數：
   ```
   PORT=3001
   LOG_LEVEL=info
   GIN_MODE=release
   DATABASE_URL=$RAILWAY_DATABASE_URL
   REDIS_URL=$RENDER_REDIS_URL
   AGENT_URL=$KOYEB_AGENT_URL
   PROMETHEUS_URL=$PROMETHEUS_URL
   GRAFANA_URL=$GRAFANA_URL
   JWT_SECRET=your-jwt-secret-key
   ```

6. 部署完成後，複製服務 URL

### 5️⃣ Fly.io - 監控系統

#### 初始化專案

```bash
# 登入 Fly.io
fly auth login

# 建立應用
fly apps create pandora-monitoring --org personal

# 建立持久化儲存
fly volumes create prometheus_data --size 3 --region nrt
fly volumes create loki_data --size 3 --region nrt
fly volumes create grafana_data --size 1 --region nrt
fly volumes create alertmanager_data --size 1 --region nrt
```

#### 設定環境變數

```bash
fly secrets set \
  GRAFANA_ADMIN_PASSWORD=pandora123 \
  LOG_LEVEL=info \
  TZ=Asia/Taipei
```

#### 部署應用

```bash
fly deploy --config fly.toml --dockerfile Dockerfile.monitoring
```

#### 驗證部署

```bash
# 檢查應用狀態
fly status

# 檢查日誌
fly logs

# 開啟應用
fly open
```

## 🔒 安全設定

### 環境變數管理

1. **絕不提交 `.env.paas` 到版本控制**
   ```bash
   # 確保 .gitignore 包含
   .env.paas
   .env.paas.local
   ```

2. **使用強密碼和隨機密鑰**
   ```bash
   # 生成加密密鑰
   openssl rand -hex 32
   
   # 生成 JWT Secret
   openssl rand -base64 48
   ```

3. **定期輪換密鑰**
   - 每 90 天更換一次資料庫密碼
   - 每 30 天更換一次 API Token
   - 每 60 天更換一次 JWT Secret

### Secrets 設定

#### GitHub Secrets

在 GitHub Repository Settings → Secrets and variables → Actions 中新增：

```
# Railway
RAILWAY_TOKEN=xxx
RAILWAY_PROJECT_ID=xxx
RAILWAY_DATABASE_URL=xxx

# Render
RENDER_API_KEY=xxx
RENDER_SERVICE_ID=xxx
RENDER_REDIS_URL=xxx
RENDER_NGINX_URL=xxx

# Koyeb
KOYEB_TOKEN=xxx
KOYEB_AGENT_URL=xxx
DOCKER_USERNAME=xxx
DOCKER_PASSWORD=xxx

# Patr.io
PATR_API_TOKEN=xxx
PATR_DEPLOYMENT_ID=xxx
PATR_UI_URL=xxx

# Fly.io
FLY_API_TOKEN=xxx
PROMETHEUS_URL=xxx
LOKI_URL=xxx
GRAFANA_URL=xxx
ALERTMANAGER_URL=xxx

# 應用程式
ENCRYPTION_KEY=xxx
JWT_SECRET=xxx
GRAFANA_ADMIN_PASSWORD=xxx
GRAFANA_API_KEY=xxx

# 告警
SLACK_WEBHOOK_URL=xxx
EMAIL_PASSWORD=xxx
```

## 📊 監控與維護

### 健康檢查

```bash
# Agent 健康檢查
curl https://pandora-agent-xxx.koyeb.app/health

# UI 健康檢查
curl https://axiom-ui-xxx.patr.cloud/api/v1/status

# Prometheus 健康檢查
curl https://pandora-monitoring.fly.dev/prometheus/-/healthy

# Loki 健康檢查
curl https://pandora-monitoring.fly.dev/loki/ready

# Grafana 健康檢查
curl https://pandora-monitoring.fly.dev/api/health

# AlertManager 健康檢查
curl https://pandora-monitoring.fly.dev/alertmanager/-/healthy
```

### 查看日誌

```bash
# Railway 日誌
railway logs

# Render 日誌 (Web Dashboard)
https://dashboard.render.com/services/xxx/logs

# Koyeb 日誌
koyeb service logs pandora-agent/pandora-agent --follow

# Patr.io 日誌 (Web Dashboard)
https://patr.cloud/deployments/xxx/logs

# Fly.io 日誌
fly logs
```

### 資源監控

訪問各平台的 Dashboard：

- **Railway**: https://railway.app/project/xxx
- **Render**: https://dashboard.render.com/services/xxx
- **Koyeb**: https://app.koyeb.com/apps/pandora-agent
- **Patr.io**: https://patr.cloud/deployments/xxx
- **Fly.io**: https://fly.io/dashboard/pandora-monitoring

## 🔄 更新部署

### 手動更新

```bash
# 更新 Agent (Koyeb)
docker build -f Dockerfile.agent.koyeb -t your-docker-username/pandora-agent:latest .
docker push your-docker-username/pandora-agent:latest
koyeb service redeploy pandora-agent/pandora-agent

# 更新 UI (Patr.io)
# 在 Patr.io Dashboard 點擊 "Redeploy"

# 更新監控系統 (Fly.io)
fly deploy --config fly.toml --dockerfile Dockerfile.monitoring
```

### 自動化更新 (GitHub Actions)

推送到 `main` 分支會自動觸發部署：

```bash
git add .
git commit -m "Update services"
git push origin main
```

手動觸發特定平台部署：

1. 前往 GitHub Actions
2. 選擇 "Deploy to PaaS Platforms"
3. 點擊 "Run workflow"
4. 選擇要部署的平台
5. 點擊 "Run workflow"

## 🐛 故障排除

### 常見問題

#### 1. Agent 無法連接資料庫

```bash
# 檢查資料庫 URL
echo $RAILWAY_DATABASE_URL

# 測試資料庫連接
psql $RAILWAY_DATABASE_URL -c "SELECT 1;"

# 檢查 Koyeb 環境變數
koyeb service env pandora-agent/pandora-agent
```

#### 2. Redis 連接失敗

```bash
# 檢查 Redis URL
echo $RENDER_REDIS_URL

# 測試 Redis 連接
redis-cli -u $RENDER_REDIS_URL ping

# 確認 Internal URL 用於內部服務通訊
```

#### 3. Prometheus 無法抓取指標

```bash
# 檢查 Prometheus 目標
curl https://pandora-monitoring.fly.dev/prometheus/api/v1/targets

# 確認 Agent 暴露 metrics 端點
curl https://pandora-agent-xxx.koyeb.app/metrics
```

#### 4. Fly.io 儲存空間不足

```bash
# 檢查儲存空間使用
fly volumes list

# 擴展儲存空間
fly volumes extend prometheus_data --size 5

# 清理舊資料
fly ssh console
cd /prometheus
rm -rf chunks_head/*.tmp
```

#### 5. 部署失敗

```bash
# 檢查建置日誌
fly logs

# 驗證 Dockerfile
docker build -f Dockerfile.monitoring -t test .

# 檢查配置檔案
fly config validate
```

## 💰 成本估算

所有服務使用免費方案：

| 平台 | 服務 | 免費額度 | 預估成本 |
|------|------|----------|---------|
| Railway | PostgreSQL | 5GB 儲存, 500hrs/月 | $0 |
| Render | Redis | 25MB RAM | $0 |
| Render | Nginx | 750hrs/月 | $0 |
| Koyeb | Agent | 2 Nano 容器 | $0 |
| Patr.io | UI | 512Mi RAM | $0 |
| Fly.io | Monitoring | 3GB 儲存, 160GB 流量/月 | $0 |
| **總計** | | | **$0/月** |

**注意**: 免費額度可能隨時變更，請參考各平台最新定價。

## 📚 相關文件

- [部署腳本文件](scripts/deploy-paas.sh)
- [驗證腳本文件](scripts/verify-paas-deployment.sh)
- [環境變數範本](env.paas.example)
- [GitHub Actions Workflow](.github/workflows/deploy-paas.yml)
- [主要 README](README.md)
- [OCI 部署指南](DEPLOYMENT.md)
- [GCP 部署指南](DEPLOYMENT-GCP.md)

## 🤝 支援與聯絡

如有問題請：

1. 查看本文件的故障排除章節
2. 檢查 [GitHub Issues](https://github.com/your-org/pandora_box_console_IDS-IPS/issues)
3. 加入討論 [GitHub Discussions](https://github.com/your-org/pandora_box_console_IDS-IPS/discussions)

---

**© 2024 Pandora Security Team | MIT License**

