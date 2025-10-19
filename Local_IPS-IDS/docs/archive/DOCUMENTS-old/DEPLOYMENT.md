# Pandora Box Console IDS-IPS 部署指南

## 概述

本文件說明如何將 Pandora Box Console IDS-IPS 系統部署到 Oracle Cloud Infrastructure (OCI) 和 Vercel。

### 系統架構

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Vercel (FE)   │    │   OCI K8s (BE)   │    │   OCI Storage   │
│                 │    │                  │    │                 │
│  - Next.js UI   │◄──►│  - Pandora Agent │    │  - PostgreSQL   │
│  - Static Files │    │  - Pandora Console│    │  - Redis        │
│  - API Proxy    │    │  - Prometheus     │    │  - Grafana      │
│                 │    │  - Loki          │    │  - Loki         │
│                 │    │  - Grafana       │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## 前置需求

### 必要工具

- [GitHub CLI](https://cli.github.com/) (gh)
- [OCI CLI](https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Docker](https://www.docker.com/)
- [Node.js](https://nodejs.org/) (v18+)

### OCI 資源需求

1. **OCI Kubernetes Cluster (OKE)**
   - 最少 3 個 worker nodes
   - 每個 node 最少 2 OCPU, 8GB RAM
   - 使用 OCI Container Registry

2. **儲存需求**
   - PostgreSQL: 10GB
   - Prometheus: 20GB
   - Loki: 10GB
   - Grafana: 5GB
   - Redis: 5GB

3. **網路需求**
   - Load Balancer
   - SSL 憑證 (Let's Encrypt 或 OCI SSL)
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
# OCI 配置
OCI_USER=ocid1.user.oc1..your-user-ocid
OCI_TENANCY=ocid1.tenancy.oc1..your-tenancy-ocid
OCI_REGION=us-ashburn-1
OCI_NAMESPACE=your-oci-namespace
CLUSTER_OCID=ocid1.cluster.oc1..your-cluster-ocid

# 網域配置
DOMAIN=pandora.yourdomain.com
```

### 步驟 2: 密鑰和憑證設定

#### 2.1 執行密鑰設定腳本

```bash
chmod +x scripts/setup-secrets.sh
./scripts/setup-secrets.sh
```

#### 2.2 手動設定 GitHub Actions 密鑰

如果腳本無法執行，請手動在 GitHub Repository Settings 中設定以下密鑰：

**OCI 配置密鑰:**
- `OCI_USER`
- `OCI_TENANCY`
- `OCI_REGION`
- `OCI_FINGERPRINT`
- `OCI_NAMESPACE`
- `CLUSTER_OCID`
- `OCI_USERNAME`
- `OCI_PASSWORD`

**Vercel 配置密鑰:**
- `VERCEL_TOKEN`
- `VERCEL_ORG_ID`
- `VERCEL_PROJECT_ID`
- `VERCEL_API_BASE_URL`
- `VERCEL_GRAFANA_URL`
- `VERCEL_PROMETHEUS_URL`

### 步驟 3: OCI 部署

#### 3.1 更新 Kubernetes manifests

更新 `k8s/` 目錄中的映像 URL：

```bash
# 替換 YOUR_NAMESPACE 為實際的 OCI namespace
find k8s/ -name "*.yaml" -exec sed -i "s/YOUR_NAMESPACE/your-actual-namespace/g" {} \;
```

#### 3.2 執行部署腳本

```bash
chmod +x scripts/deploy-oci.sh
./scripts/deploy-oci.sh
```

#### 3.3 手動部署 (可選)

如果自動腳本失敗，可以手動執行：

```bash
# 1. 建置並推送 Docker 映像
docker build -f Dockerfile.agent -t iad.ocir.io/your-namespace/pandora-agent:latest .
docker build -f Dockerfile -t iad.ocir.io/your-namespace/pandora-console:latest .

# 2. 推送到 OCI Container Registry
docker push iad.ocir.io/your-namespace/pandora-agent:latest
docker push iad.ocir.io/your-namespace/pandora-console:latest

# 3. 部署到 Kubernetes
kubectl apply -k k8s/

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
NEXT_PUBLIC_API_BASE_URL=https://pandora.yourdomain.com/api
NEXT_PUBLIC_GRAFANA_URL=https://pandora.yourdomain.com/grafana
NEXT_PUBLIC_PROMETHEUS_URL=https://pandora.yourdomain.com/prometheus
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

將您的網域指向 OCI Load Balancer IP：

```
A    pandora.yourdomain.com    →    OCI Load Balancer IP
CNAME www.pandora.yourdomain.com →   pandora.yourdomain.com
```

#### 5.2 SSL 憑證

使用 cert-manager 自動管理 SSL 憑證：

```bash
# 安裝 cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# 創建 Let's Encrypt ClusterIssuer
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
EOF
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

# 檢查日誌
kubectl logs -f deployment/pandora-console -n pandora-box
```

### 訪問服務

部署完成後，您可以訪問以下服務：

- **前端 UI:** `https://pandora.yourdomain.com`
- **Grafana Dashboard:** `https://pandora.yourdomain.com/grafana`
  - 用戶名: `admin`
  - 密碼: `pandora123`
- **Prometheus:** `https://pandora.yourdomain.com/prometheus`
- **API 健康檢查:** `https://pandora.yourdomain.com/api/v1/health`

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
kubectl create secret docker-registry oci-registry-secret \
  --docker-server=iad.ocir.io \
  --docker-username=your-namespace/your-user \
  --docker-password=your-auth-token \
  --namespace=pandora-box
```

### 效能調優

#### 1. 資源限制

根據實際使用情況調整 `k8s/` 目錄中各服務的資源限制：

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

- 使用 OCI 的內建負載均衡器
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
kubectl set image deployment/pandora-console pandora-console=iad.ocir.io/your-namespace/pandora-console:new-tag -n pandora-box

# 檢查滾動更新狀態
kubectl rollout status deployment/pandora-console -n pandora-box
```

### 系統升級

1. 備份現有資料
2. 更新 Kubernetes manifests
3. 執行滾動更新
4. 驗證服務正常運作

## 支援和聯絡

如果遇到問題，請：

1. 檢查本文檔的故障排除章節
2. 查看 GitHub Issues
3. 聯繫技術支援團隊

---

**注意:** 請確保在生產環境中使用強密碼和適當的安全配置。
