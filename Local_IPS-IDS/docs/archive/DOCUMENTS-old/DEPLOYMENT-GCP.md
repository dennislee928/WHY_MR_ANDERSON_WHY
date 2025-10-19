# Pandora Box Console IDS-IPS GCP 部署指南

## 概述

本文件說明如何將 Pandora Box Console IDS-IPS 系統部署到 Google Cloud Platform (GCP) 和 Vercel。

### 系統架構

```
┌─────────────────────────────────────────────────────────────┐
│                    Pandora Box Console                     │
│                      IDS/IPS System                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   🌐 Vercel     │    │   ☁️ GCP GKE     │    │   💾 Storage    │
│                 │    │                  │    │                 │
│  • Next.js UI   │◄──►│  • Pandora Agent │    │  • Cloud SQL    │
│  • Static Files │    │  • Console API   │    │  • Memorystore  │
│  • API Proxy    │    │  • Prometheus    │    │  • Grafana      │
│                 │    │  • Loki          │    │  • Loki         │
│                 │    │  • Grafana       │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## 前置需求

### 必要工具

- [GitHub CLI](https://cli.github.com/) (gh)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Docker](https://www.docker.com/)
- [Node.js](https://nodejs.org/) (v18+)

### GCP 資源需求

1. **Google Kubernetes Engine (GKE)**
   - 最少 3 個 worker nodes
   - 每個 node 最少 2 vCPU, 8GB RAM
   - 使用 Google Container Registry (GCR)

2. **儲存需求**
   - PostgreSQL: 10GB (Standard persistent disk)
   - Prometheus: 20GB (Standard persistent disk)
   - Loki: 10GB (Standard persistent disk)
   - Grafana: 5GB (Standard persistent disk)
   - Redis: 5GB (Standard persistent disk)

3. **網路需求**
   - Load Balancer
   - SSL 憑證 (Google Managed SSL)
   - 網域設定

### Vercel 需求

1. **Vercel 帳號**
2. **專案設定**
3. **環境變數配置**

## 部署步驟

### 步驟 1: 環境準備

#### 1.1 複製環境變數範本

```bash
cp env.example .env
```

#### 1.2 編輯環境變數

更新 `.env` 檔案中的以下變數：

```bash
# GCP 配置
GCP_PROJECT_ID=your-gcp-project-id
GCP_REGION=us-central1
GCP_ZONE=us-central1-a
GCP_CLUSTER_NAME=pandora-cluster
GCP_SA_KEY=/path/to/your/service-account-key.json

# 網域配置
DOMAIN_GCP=pandora-gcp.yourdomain.com
```

### 步驟 2: 密鑰和憑證設定

#### 2.1 執行密鑰設定腳本

```bash
chmod +x scripts/setup-secrets.sh
./scripts/setup-secrets.sh
```

#### 2.2 手動設定 GitHub Actions 密鑰

如果腳本無法執行，請手動在 GitHub Repository Settings 中設定以下密鑰：

**GCP 配置密鑰:**
- `GCP_PROJECT_ID`
- `GCP_CLUSTER_NAME`
- `GCP_SA_KEY`

**Vercel 配置密鑰:**
- `VERCEL_TOKEN`
- `VERCEL_ORG_ID`
- `VERCEL_PROJECT_ID_GCP`
- `VERCEL_API_BASE_URL_GCP`
- `VERCEL_GRAFANA_URL_GCP`
- `VERCEL_PROMETHEUS_URL_GCP`

### 步驟 3: GCP 部署

#### 3.1 更新 Kubernetes manifests

更新 `k8s-gcp/` 目錄中的映像 URL：

```bash
# 替換 YOUR_PROJECT_ID 為實際的 GCP Project ID
find k8s-gcp/ -name "*.yaml" -exec sed -i "s/YOUR_PROJECT_ID/your-actual-project-id/g" {} \;
```

#### 3.2 執行部署腳本

```bash
chmod +x scripts/deploy-gcp.sh
./scripts/deploy-gcp.sh
```

#### 3.3 手動部署 (可選)

如果自動腳本失敗，可以手動執行：

```bash
# 1. 建置並推送 Docker 映像
docker build -f Dockerfile.agent -t gcr.io/your-project-id/pandora-agent:latest .
docker build -f Dockerfile -t gcr.io/your-project-id/pandora-console:latest .

# 2. 推送到 Google Container Registry
docker push gcr.io/your-project-id/pandora-agent:latest
docker push gcr.io/your-project-id/pandora-console:latest

# 3. 部署到 Kubernetes
kubectl apply -k k8s-gcp/

# 4. 檢查部署狀態
kubectl get pods -n pandora-box
kubectl get services -n pandora-box
kubectl get ingress -n pandora-box
```

### 步驟 4: Vercel 前端部署

#### 4.1 設定 Vercel 專案

1. 在 [Vercel Dashboard](https://vercel.com/dashboard) 創建新專案
2. 連接 GitHub Repository
3. 設定專案配置：
   - **Framework Preset:** Next.js
   - **Build Command:** `npm run build`
   - **Output Directory:** `dist`
   - **Install Command:** `npm ci`

#### 4.2 設定環境變數

在 Vercel 專案設定中添加以下環境變數：

```
NEXT_PUBLIC_API_BASE_URL=https://pandora-gcp.yourdomain.com/api
NEXT_PUBLIC_GRAFANA_URL=https://pandora-gcp.yourdomain.com/grafana
NEXT_PUBLIC_PROMETHEUS_URL=https://pandora-gcp.yourdomain.com/prometheus
```

#### 4.3 部署前端

```bash
# 安裝依賴
npm install

# 本地建置測試
npm run build

# 部署到 Vercel (通過 GitHub Actions 自動部署)
git push origin main
```

### 步驟 5: 網域和 SSL 設定

#### 5.1 設定 DNS

將您的網域指向 GCP Load Balancer IP：

```
A    pandora-gcp.yourdomain.com    →    GCP Load Balancer IP
CNAME www.pandora-gcp.yourdomain.com →   pandora-gcp.yourdomain.com
```

#### 5.2 SSL 憑證

使用 Google Managed SSL 憑證：

```bash
# 創建 ManagedCertificate
kubectl apply -f - <<EOF
apiVersion: networking.gke.io/v1
kind: ManagedCertificate
metadata:
  name: pandora-gcp-ssl-cert
  namespace: pandora-box
spec:
  domains:
    - pandora-gcp.yourdomain.com
EOF

# 檢查憑證狀態
kubectl describe managedcertificate pandora-gcp-ssl-cert -n pandora-box
```

## 驗證部署

### 檢查服務狀態

```bash
# 檢查 Pod 狀態
kubectl get pods -n pandora-box

# 檢查服務狀態
kubectl get services -n pandora-box

# 檢查 Ingress 狀態
kubectl get ingress -n pandora-box

# 檢查 ManagedCertificate 狀態
kubectl get managedcertificate -n pandora-box

# 檢查日誌
kubectl logs -f deployment/pandora-console -n pandora-box
```

### 訪問服務

部署完成後，您可以訪問以下服務：

- **前端 UI:** `https://pandora-gcp.yourdomain.com`
- **Grafana Dashboard:** `https://pandora-gcp.yourdomain.com/grafana`
  - 用戶名: `admin`
  - 密碼: `pandora123`
- **Prometheus:** `https://pandora-gcp.yourdomain.com/prometheus`
- **API 健康檢查:** `https://pandora-gcp.yourdomain.com/api/v1/health`

## 監控和維護

### 日誌監控

```bash
# 查看所有服務日誌
kubectl logs -f deployment/pandora-console -n pandora-box
kubectl logs -f deployment/pandora-agent -n pandora-box
kubectl logs -f deployment/grafana -n pandora-box

# 使用 Loki 查詢日誌
kubectl port-forward service/loki 3100:3100 -n pandora-box
# 然後訪問 http://localhost:3100
```

### 指標監控

```bash
# 查看 Prometheus 指標
kubectl port-forward service/prometheus 9090:9090 -n pandora-box
# 然後訪問 http://localhost:9090
```

### 備份

```bash
# 備份 PostgreSQL
kubectl exec -it deployment/postgres -n pandora-box -- pg_dump -U pandora pandora > backup.sql

# 備份 Prometheus 資料
kubectl exec -it deployment/prometheus -n pandora-box -- tar -czf /tmp/prometheus-backup.tar.gz /prometheus
kubectl cp pandora-box/prometheus-xxx:/tmp/prometheus-backup.tar.gz ./prometheus-backup.tar.gz
```

## 故障排除

### 常見問題

#### 1. Pod 無法啟動

```bash
# 檢查 Pod 狀態
kubectl describe pod <pod-name> -n pandora-box

# 檢查事件
kubectl get events -n pandora-box --sort-by='.lastTimestamp'
```

#### 2. 服務無法訪問

```bash
# 檢查服務端點
kubectl get endpoints -n pandora-box

# 檢查 Ingress 狀態
kubectl describe ingress pandora-ingress -n pandora-box

# 檢查 ManagedCertificate 狀態
kubectl describe managedcertificate pandora-gcp-ssl-cert -n pandora-box
```

#### 3. 資料庫連接問題

```bash
# 檢查 PostgreSQL 連接
kubectl exec -it deployment/postgres -n pandora-box -- psql -U pandora -d pandora

# 檢查 Redis 連接
kubectl exec -it deployment/redis -n pandora-box -- redis-cli ping
```

#### 4. 映像拉取失敗

```bash
# 檢查映像拉取密鑰
kubectl get secrets -n pandora-box

# 重新創建映像拉取密鑰
kubectl create secret docker-registry gcr-registry-secret \
  --docker-server=gcr.io \
  --docker-username=_json_key \
  --docker-password="$(gcloud auth print-access-token)" \
  --namespace=pandora-box
```

#### 5. SSL 憑證問題

```bash
# 檢查 ManagedCertificate 狀態
kubectl describe managedcertificate pandora-gcp-ssl-cert -n pandora-box

# 檢查 DNS 設定
nslookup pandora-gcp.yourdomain.com

# 檢查 Load Balancer 狀態
gcloud compute forwarding-rules list
```

### 效能調優

#### 1. 資源限制

根據實際使用情況調整 `k8s-gcp/` 目錄中各服務的資源限制：

```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

#### 2. 儲存優化

- 使用 SSD 儲存類別提高 I/O 效能
- 定期清理舊的日誌和指標資料
- 設定適當的保留期

#### 3. 網路優化

- 使用 GCP 的內建負載均衡器
- 設定適當的連接池大小
- 啟用 HTTP/2 和 gzip 壓縮

## 安全考量

### 1. 網路安全

- 使用 TLS 加密所有通訊
- 設定適當的防火牆規則
- 限制管理端口的訪問

### 2. 資料安全

- 定期備份資料庫
- 加密敏感資料
- 使用強密碼和定期輪換

### 3. 存取控制

- 使用 RBAC 控制 Kubernetes 存取
- 設定適當的 Grafana 使用者權限
- 啟用 API 認證

## 更新和升級

### 應用程式更新

```bash
# 更新映像標籤
kubectl set image deployment/pandora-console pandora-console=gcr.io/your-project-id/pandora-console:new-tag -n pandora-box

# 檢查滾動更新狀態
kubectl rollout status deployment/pandora-console -n pandora-box
```

### 系統升級

1. 備份現有資料
2. 更新 Kubernetes manifests
3. 執行滾動更新
4. 驗證服務正常運作

## 成本優化

### 1. 資源優化

- 使用 Preemptible VMs 降低成本
- 設定適當的節點大小
- 啟用自動擴縮容

### 2. 儲存優化

- 使用適當的儲存類別
- 定期清理不需要的資料
- 使用快照備份

### 3. 網路優化

- 使用區域負載均衡器
- 優化 CDN 設定
- 減少跨區域流量

## 支援和聯絡

如果遇到問題，請：

1. 檢查本文檔的故障排除章節
2. 查看 GitHub Issues
3. 聯繫技術支援團隊

---

**注意:** 請確保在生產環境中使用強密碼和適當的安全配置。
