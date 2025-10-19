# 🚀 Pandora Box Console IDS-IPS 部署指南

## 📋 快速開始

本專案提供完整的 CI/CD 部署方案，將後端服務部署到 **Oracle Cloud Infrastructure (OCI)**，前端部署到 **Vercel**。

### 🎯 部署架構

```
┌─────────────────────────────────────────────────────────────┐
│                    Pandora Box Console                     │
│                      IDS/IPS System                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   🌐 Vercel     │    │   ☁️ OCI K8s     │    │   💾 Storage    │
│                 │    │                  │    │                 │
│  • Next.js UI   │◄──►│  • Pandora Agent │    │  • PostgreSQL   │
│  • Static Files │    │  • Console API   │    │  • Redis        │
│  • API Proxy    │    │  • Prometheus    │    │  • Grafana      │
│                 │    │  • Loki          │    │  • Loki         │
│                 │    │  • Grafana       │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## 🛠️ 前置需求

### 必要工具
- [GitHub CLI](https://cli.github.com/) - 管理 GitHub 密鑰
- [OCI CLI](https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm) - Oracle Cloud 管理
- [kubectl](https://kubernetes.io/docs/tasks/tools/) - Kubernetes 管理
- [Docker](https://www.docker.com/) - 容器建置
- [Node.js](https://nodejs.org/) v18+ - 前端開發

### 雲端資源
- **OCI Kubernetes Cluster (OKE)** - 最少 3 nodes, 2 OCPU, 8GB RAM
- **OCI Container Registry** - 存放 Docker 映像
- **Vercel 帳號** - 前端部署平台
- **網域和 SSL 憑證** - 生產環境必備

## ⚡ 一鍵部署

### 1️⃣ 設定環境變數

```bash
# 複製環境變數範本
cp env.example .env

# 編輯並填入您的配置
nano .env
```

### 2️⃣ 設定密鑰和憑證

```bash
# 執行密鑰設定腳本
chmod +x scripts/setup-secrets.sh
./scripts/setup-secrets.sh
```

### 3️⃣ 執行部署

```bash
# 執行 OCI 部署腳本
chmod +x scripts/deploy-oci.sh
./scripts/deploy-oci.sh
```

### 4️⃣ 設定 Vercel

1. 在 [Vercel Dashboard](https://vercel.com/dashboard) 創建專案
2. 連接 GitHub Repository
3. 設定環境變數
4. 部署完成！

## 📁 檔案結構

```
pandora_box_console_IDS-IPS/
├── 📁 k8s/                          # Kubernetes 部署配置
│   ├── namespace.yaml               # 命名空間
│   ├── configmap.yaml              # 配置對應
│   ├── secrets.yaml                # 密鑰管理
│   ├── postgres.yaml               # PostgreSQL 部署
│   ├── redis.yaml                  # Redis 部署
│   ├── prometheus.yaml             # Prometheus 部署
│   ├── loki.yaml                   # Loki 部署
│   ├── grafana.yaml                # Grafana 部署
│   ├── pandora-backend.yaml        # 後端服務部署
│   ├── ingress.yaml                # 入口配置
│   └── kustomization.yaml          # Kustomize 配置
├── 📁 scripts/                      # 部署腳本
│   ├── deploy-oci.sh               # OCI 部署腳本
│   └── setup-secrets.sh            # 密鑰設定腳本
├── 📁 .github/workflows/            # GitHub Actions
│   └── ci.yml                      # CI/CD Pipeline
├── 📄 vercel.json                   # Vercel 配置
├── 📄 package.json                  # 前端依賴
├── 📄 next.config.js               # Next.js 配置
├── 📄 tailwind.config.js           # Tailwind CSS 配置
├── 📄 tsconfig.json                # TypeScript 配置
├── 📄 env.example                  # 環境變數範本
├── 📄 DEPLOYMENT.md                # 詳細部署文件
└── 📄 README-DEPLOYMENT.md         # 快速部署指南
```

## 🔧 GitHub Actions CI/CD

本專案包含完整的 CI/CD Pipeline：

### Pipeline 流程

1. **🔍 基本檢查** - Go 程式碼檢查、格式化、測試
2. **🎨 前端檢查** - TypeScript 檢查、Linting、測試
3. **🐳 Docker 建置** - 建置並測試 Docker 映像
4. **🔒 安全掃描** - Trivy 漏洞掃描
5. **📦 映像推送** - 推送映像到 OCI Container Registry
6. **🚀 OCI 部署** - 部署到 OCI Kubernetes Cluster
7. **🌐 Vercel 部署** - 部署前端到 Vercel

### 觸發條件

- **Push 到 main/develop** - 執行完整 CI/CD
- **Pull Request** - 執行檢查和測試
- **Tag 推送** - 執行生產部署
- **手動觸發** - 支援 workflow_dispatch

## 🔐 密鑰管理

### GitHub Secrets

需要在 GitHub Repository Settings 中設定：

**OCI 配置:**
- `OCI_USER` - OCI 使用者 OCID
- `OCI_TENANCY` - OCI Tenancy OCID
- `OCI_REGION` - OCI 區域
- `OCI_FINGERPRINT` - API 金鑰指紋
- `OCI_NAMESPACE` - OCI Container Registry 命名空間
- `CLUSTER_OCID` - Kubernetes Cluster OCID
- `OCI_USERNAME` - Registry 使用者名稱
- `OCI_PASSWORD` - Registry 認證密碼

**Vercel 配置:**
- `VERCEL_TOKEN` - Vercel API Token
- `VERCEL_ORG_ID` - Vercel 組織 ID
- `VERCEL_PROJECT_ID` - Vercel 專案 ID
- `VERCEL_API_BASE_URL` - API 基礎 URL
- `VERCEL_GRAFANA_URL` - Grafana URL
- `VERCEL_PROMETHEUS_URL` - Prometheus URL

### Kubernetes Secrets

自動創建的密鑰：
- `pandora-secrets` - 應用程式密鑰
- `oci-registry-secret` - Docker Registry 認證
- `pandora-mtls-certs` - mTLS 伺服器憑證
- `pandora-client-certs` - mTLS 客戶端憑證

## 🌐 服務端點

部署完成後可訪問：

| 服務 | URL | 說明 |
|------|-----|------|
| 🏠 前端 UI | `https://pandora.yourdomain.com` | 主要使用者介面 |
| 📊 Grafana | `https://pandora.yourdomain.com/grafana` | 監控儀表板 |
| 📈 Prometheus | `https://pandora.yourdomain.com/prometheus` | 指標查詢 |
| 🔌 API | `https://pandora.yourdomain.com/api/v1/health` | API 健康檢查 |
| 🤖 Agent | `https://pandora.yourdomain.com/agent/health` | Agent 健康檢查 |

### 預設認證

- **Grafana:** 用戶名 `admin`, 密碼 `pandora123`
- **API:** 使用 JWT Token 認證

## 🔍 監控和維護

### 日誌查看

```bash
# 查看所有服務日誌
kubectl logs -f deployment/pandora-console -n pandora-box
kubectl logs -f deployment/pandora-agent -n pandora-box

# 使用 Loki 查詢日誌
kubectl port-forward service/loki 3100:3100 -n pandora-box
```

### 指標監控

```bash
# 查看 Prometheus 指標
kubectl port-forward service/prometheus 9090:9090 -n pandora-box
```

### 備份

```bash
# 備份資料庫
kubectl exec -it deployment/postgres -n pandora-box -- pg_dump -U pandora pandora > backup.sql
```

## 🚨 故障排除

### 常見問題

1. **Pod 無法啟動**
   ```bash
   kubectl describe pod <pod-name> -n pandora-box
   kubectl get events -n pandora-box --sort-by='.lastTimestamp'
   ```

2. **服務無法訪問**
   ```bash
   kubectl get endpoints -n pandora-box
   kubectl describe ingress pandora-ingress -n pandora-box
   ```

3. **映像拉取失敗**
   ```bash
   kubectl get secrets -n pandora-box
   # 重新創建 Registry 密鑰
   ```

## 🔄 更新流程

### 自動更新

推送程式碼到 `main` 分支會自動觸發部署：

```bash
git add .
git commit -m "feat: 新功能"
git push origin main
```

### 手動更新

```bash
# 更新特定服務
kubectl set image deployment/pandora-console pandora-console=iad.ocir.io/your-namespace/pandora-console:new-tag -n pandora-box

# 檢查更新狀態
kubectl rollout status deployment/pandora-console -n pandora-box
```

## 📞 支援

- 📖 [詳細部署文件](DEPLOYMENT.md)
- 🐛 [GitHub Issues](https://github.com/your-repo/issues)
- 📧 技術支援: support@yourdomain.com

## 📄 授權

本專案採用 MIT 授權條款。

---

**🎉 部署完成後，您將擁有一個完整的 IDS/IPS 監控系統！**
