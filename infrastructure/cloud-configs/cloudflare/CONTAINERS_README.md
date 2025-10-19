# Cloudflare Containers - Security Platform

## 🚀 **概述**

本專案使用 Cloudflare Containers 來增強 Security Platform Workers，提供容器化的微服務架構，包括：

- **Backend API Container** - Go 後端 API 服務
- **AI/Quantum Container** - Python AI 和量子計算服務
- **Security Tools Container** - 安全掃描工具集合
- **Database Container** - PostgreSQL 資料庫
- **Monitoring Container** - Prometheus 監控系統

## 🏗️ **架構圖**

```
┌─────────────────────────────────────────────────────────────┐
│                    Cloudflare Workers                       │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │              Container Orchestrator                     │ │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐        │ │
│  │  │   Backend   │ │    AI/      │ │  Security   │        │ │
│  │  │     API     │ │  Quantum    │ │   Tools     │        │ │
│  │  └─────────────┘ └─────────────┘ └─────────────┘        │ │
│  │  ┌─────────────┐ ┌─────────────┐                        │ │
│  │  │  Database   │ │ Monitoring  │                        │ │
│  │  └─────────────┘ └─────────────┘                        │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 📁 **檔案結構**

```
infrastructure/cloud-configs/cloudflare/
├── wrangler-containers.toml          # Cloudflare Containers 配置
├── src/
│   └── containers-worker.js          # 容器編排 Worker
├── containers/
│   ├── backend-api/
│   │   ├── Dockerfile               # Backend API 容器
│   │   ├── package.json
│   │   └── src/
│   ├── ai-quantum/
│   │   ├── Dockerfile               # AI/Quantum 容器
│   │   ├── requirements.txt
│   │   └── main.py
│   ├── security-tools/
│   │   ├── Dockerfile               # 安全工具容器
│   │   └── security_api.py
│   ├── database/
│   │   ├── Dockerfile               # 資料庫容器
│   │   └── init-scripts/
│   └── monitoring/
│       ├── Dockerfile               # 監控容器
│       ├── prometheus.yml
│       └── rules/
├── docker-compose.yml                # 本地開發用
├── deploy-containers.sh              # Linux/Mac 部署腳本
├── deploy-containers.ps1             # Windows 部署腳本
└── README.md                         # 本檔案
```

## 🛠️ **安裝與設定**

### **1. 安裝依賴**

#### **Docker**
```bash
# 安裝 Docker Desktop
# Windows: https://docs.docker.com/desktop/windows/install/
# Mac: https://docs.docker.com/desktop/mac/install/
# Linux: https://docs.docker.com/engine/install/
```

#### **Wrangler CLI**
```bash
# 安裝 Wrangler
npm install -g wrangler

# 驗證安裝
wrangler --version
```

#### **GitHub Container Registry 存取**
```bash
# 設定 GitHub Token
export GITHUB_TOKEN="your_github_token"
export GITHUB_USERNAME="your_username"

# 登入 GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u $GITHUB_USERNAME --password-stdin
```

### **2. 配置 Cloudflare**

#### **設定 Wrangler**
```bash
# 登入 Cloudflare
wrangler login

# 設定專案
wrangler init security-platform-containers
```

#### **更新 wrangler-containers.toml**
```toml
# 更新您的帳戶 ID
account_id = "your_account_id"

# 更新路由
routes = [
  { pattern = "api.yourdomain.com/*", zone_name = "yourdomain.com" }
]
```

## 🚀 **部署**

### **方法 1: 使用部署腳本**

#### **Linux/Mac**
```bash
# 執行完整部署
./deploy-containers.sh

# 指定版本和環境
./deploy-containers.sh v1.0.0 production all

# 只構建容器
./deploy-containers.sh latest production build

# 只部署到 Cloudflare
./deploy-containers.sh latest production deploy
```

#### **Windows PowerShell**
```powershell
# 執行完整部署
.\deploy-containers.ps1

# 指定版本和環境
.\deploy-containers.ps1 v1.0.0 production all

# 只構建容器
.\deploy-containers.ps1 latest production build

# 只部署到 Cloudflare
.\deploy-containers.ps1 latest production deploy
```

### **方法 2: 手動部署**

#### **1. 構建容器**
```bash
# 構建所有容器
docker build -t ghcr.io/yourusername/security-platform-backend-api:latest ./containers/backend-api/
docker build -t ghcr.io/yourusername/security-platform-ai-quantum:latest ./containers/ai-quantum/
docker build -t ghcr.io/yourusername/security-platform-security-tools:latest ./containers/security-tools/
docker build -t ghcr.io/yourusername/security-platform-database:latest ./containers/database/
docker build -t ghcr.io/yourusername/security-platform-monitoring:latest ./containers/monitoring/
```

#### **2. 推送到 Registry**
```bash
# 推送所有容器
docker push ghcr.io/yourusername/security-platform-backend-api:latest
docker push ghcr.io/yourusername/security-platform-ai-quantum:latest
docker push ghcr.io/yourusername/security-platform-security-tools:latest
docker push ghcr.io/yourusername/security-platform-database:latest
docker push ghcr.io/yourusername/security-platform-monitoring:latest
```

#### **3. 部署到 Cloudflare**
```bash
# 部署 Worker
wrangler deploy --config wrangler-containers.toml
```

## 🧪 **測試**

### **本地測試**
```bash
# 使用 Docker Compose 啟動所有服務
docker-compose up -d

# 檢查服務狀態
docker-compose ps

# 測試健康檢查
curl http://localhost:3000/health
curl http://localhost:8000/health
curl http://localhost:8080/health
curl http://localhost:9090/-/healthy
```

### **Cloudflare Workers 測試**
```bash
# 測試容器健康檢查
curl https://your-worker.your-subdomain.workers.dev/api/v1/containers/health

# 測試服務發現
curl https://your-worker.your-subdomain.workers.dev/api/v1/services

# 測試後端 API
curl https://your-worker.your-subdomain.workers.dev/api/v1/backend/health

# 測試 AI 服務
curl https://your-worker.your-subdomain.workers.dev/api/v1/ai/analyze
```

## 📊 **監控與日誌**

### **容器健康檢查**
```bash
# 檢查所有容器健康狀態
curl https://your-worker.your-subdomain.workers.dev/api/v1/containers/health

# 回應範例
{
  "overall": "healthy",
  "containers": {
    "backend": { "healthy": true, "lastCheck": 1640995200000 },
    "ai": { "healthy": true, "lastCheck": 1640995200000 },
    "security": { "healthy": true, "lastCheck": 1640995200000 },
    "database": { "healthy": true, "lastCheck": 1640995200000 },
    "monitoring": { "healthy": true, "lastCheck": 1640995200000 }
  },
  "timestamp": "2024-01-01T00:00:00.000Z"
}
```

### **服務發現**
```bash
# 獲取所有服務資訊
curl https://your-worker.your-subdomain.workers.dev/api/v1/services

# 回應範例
{
  "services": {
    "backend": {
      "binding": "BACKEND_API",
      "healthy": true,
      "endpoints": ["/backend/api/v1/*", "/backend/health"]
    },
    "ai": {
      "binding": "AI_QUANTUM",
      "healthy": true,
      "endpoints": ["/ai/api/v1/*", "/ai/health"]
    }
  }
}
```

### **指標聚合**
```bash
# 獲取所有服務指標
curl https://your-worker.your-subdomain.workers.dev/api/v1/metrics
```

## 🔧 **配置選項**

### **容器資源限制**
```toml
[containers.resources]
cpu_limit = "1000m"      # CPU 限制
memory_limit = "512Mi"    # 記憶體限制
```

### **環境變數**
```toml
[containers.env]
NODE_ENV = "production"
LOG_LEVEL = "info"
API_VERSION = "v1"
```

### **健康檢查配置**
```toml
[containers.health_check]
path = "/health"          # 健康檢查路徑
interval = 30            # 檢查間隔 (秒)
timeout = 10             # 超時時間 (秒)
retries = 3              # 重試次數
```

## 🚨 **故障排除**

### **常見問題**

#### **1. 容器無法啟動**
```bash
# 檢查容器日誌
docker logs container_name

# 檢查容器狀態
docker ps -a

# 重新構建容器
docker build --no-cache -t image_name .
```

#### **2. Cloudflare Workers 部署失敗**
```bash
# 檢查 Wrangler 配置
wrangler whoami

# 檢查帳戶權限
wrangler accounts list

# 查看部署日誌
wrangler tail
```

#### **3. 容器間通訊問題**
```bash
# 檢查網路配置
docker network ls
docker network inspect network_name

# 測試容器間連線
docker exec container1 ping container2
```

#### **4. 健康檢查失敗**
```bash
# 手動測試健康檢查端點
curl -f http://container_ip:port/health

# 檢查容器資源使用
docker stats container_name
```

### **除錯模式**
```bash
# 啟用詳細日誌
wrangler tail --format=pretty

# 本地除錯
wrangler dev --config wrangler-containers.toml
```

## 📈 **效能優化**

### **容器優化**
- 使用多階段構建減少映像大小
- 設定適當的資源限制
- 啟用容器快取
- 使用健康檢查確保服務可用性

### **網路優化**
- 使用 Cloudflare 的全球網路
- 啟用 HTTP/2 和 HTTP/3
- 設定適當的快取策略
- 使用 CDN 加速靜態資源

### **監控優化**
- 設定 Prometheus 指標收集
- 配置 Grafana 儀表板
- 設定告警規則
- 監控容器資源使用

## 🔄 **CI/CD 整合**

### **GitHub Actions**
```yaml
name: Deploy Containers
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Build and Deploy
        run: |
          ./deploy-containers.sh ${{ github.sha }} production all
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### **自動化部署**
- 推送到 main 分支自動觸發部署
- 使用 Git SHA 作為容器版本
- 自動執行測試和驗證
- 部署失敗時自動回滾

## 📚 **API 參考**

### **容器管理 API**

#### **健康檢查**
```http
GET /api/v1/containers/health
```

#### **服務發現**
```http
GET /api/v1/services
```

#### **容器擴展**
```http
POST /api/v1/containers/{serviceName}/scale
Content-Type: application/json

{
  "replicas": 3
}
```

#### **容器日誌**
```http
GET /api/v1/containers/{serviceName}/logs
```

#### **指標聚合**
```http
GET /api/v1/metrics
```

### **服務 API**

#### **後端 API**
```http
GET /api/v1/backend/health
GET /api/v1/backend/api/v1/threats
POST /api/v1/backend/api/v1/threats
```

#### **AI/Quantum**
```http
GET /api/v1/ai/health
POST /api/v1/ai/api/v1/analyze
POST /api/v1/ai/api/v1/quantum-process
```

#### **安全工具**
```http
GET /api/v1/security/health
POST /api/v1/security/api/v1/scan
GET /api/v1/security/api/v1/reports
```

## 🤝 **貢獻指南**

1. Fork 專案
2. 建立功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交變更 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 開啟 Pull Request

## 📄 **授權**

本專案採用 MIT 授權 - 詳見 [LICENSE](LICENSE) 檔案

---

**注意**: 請確保在部署前正確配置所有環境變數和憑證，並測試所有容器功能。
