# Buddy Works CI/CD 設置指南
## 完整管道配置

> 📅 **創建日期**: 2025-10-09  
> 🎯 **目標**: 整合所有 GitHub Actions 到 Buddy Works  
> 📊 **管道數量**: 12 個  
> ✅ **狀態**: 配置完成

---

## 🚀 快速開始

### 1. 連接 Buddy Works

1. 訪問 [app.buddy.works](https://app.buddy.works)
2. 點擊 "New Project"
3. 選擇 "GitHub" 作為 Git 提供商
4. 授權 Buddy 訪問您的倉庫
5. 選擇 `Local_IPS-IDS` 倉庫

### 2. 導入管道配置

```bash
# 方法 1: 使用 buddy.yml（推薦）
# Buddy 會自動檢測倉庫根目錄的 buddy.yml 並導入所有管道

# 方法 2: 手動創建
# 在 Buddy UI 中逐個創建管道
```

### 3. 配置環境變數

在 Buddy 項目設置中添加以下變數：

| 變數名稱 | 描述 | 範例值 |
|----------|------|--------|
| `GITHUB_TOKEN` | GitHub Personal Access Token | ghp_xxxx |
| `DOCKER_USERNAME` | Docker Hub 用戶名 | your-username |
| `DOCKER_PASSWORD` | Docker Hub 密碼 | your-password |
| `SLACK_WEBHOOK_URL` | Slack Webhook URL | https://hooks.slack.com/... |

---

## 📋 管道總覽

### 12 個 Buddy 管道

| # | 管道名稱 | 觸發方式 | 優先級 | 狀態 | 對應 GitHub Action |
|---|----------|----------|--------|------|-------------------|
| 1 | Build On-Premise Installers | Push (main/dev) | HIGH | ✅ | build-onpremise-installers.yml |
| 2 | CI Pipeline | Push (main/dev) | NORMAL | ✅ | ci.yml |
| 3 | Kubernetes Deployment | Manual | HIGH | ✅ | 新增 |
| 4 | ArgoCD GitOps Sync | Manual | NORMAL | ✅ | 新增 |
| 5 | Performance Testing | Manual | NORMAL | ✅ | 新增（Phase 4） |
| 6 | Security Audit | Manual | HIGH | ✅ | 新增（Phase 4） |
| 7 | Chaos Engineering | Manual | NORMAL | ✅ | 新增（Phase 4） |
| 8 | Deploy to GCP | Manual | LOW | ❌ Disabled | deploy-gcp.yml |
| 9 | Deploy to OCI | Manual | LOW | ❌ Disabled | deploy-oci.yml |
| 10 | Deploy to PaaS | Manual | LOW | ❌ Disabled | deploy-paas.yml |
| 11 | Terraform Deploy | Manual | LOW | ❌ Disabled | terraform-deploy.yml |
| 12 | ML Model Validation | Manual | NORMAL | ✅ | 新增（Phase 4） |

### 額外的自動化管道

| # | 管道名稱 | 觸發方式 | 頻率 | 用途 |
|---|----------|----------|------|------|
| 13 | Monitoring & Alerts | Schedule | 每 6 小時 | 健康檢查 |
| 14 | Backup & DR | Schedule | 每天 2:00 AM | 備份 |
| 15 | Documentation Build | Push | - | 文檔生成 |
| 16 | Notification Pipeline | Schedule | 每週一 9:00 AM | 週報 |
| 17 | Dependency Updates | Schedule | 每週一 0:00 AM | 依賴更新 |

---

## 🎯 管道詳細說明

### 1. Build On-Premise Installers

**觸發**: Push to main/dev  
**優先級**: HIGH  
**用途**: 構建所有平台的安裝檔

**階段**:
1. 準備構建環境（取得版本資訊）
2. 構建後端（Linux/Windows/macOS）
3. 構建前端（Next.js）
4. 構建 Linux 套件（.deb）
5. 構建 ISO 映像
6. 創建 GitHub Release（使用 GitHub API）

**產物**:
- Windows: `.exe` 安裝程式
- Linux: `.deb`, `.rpm` 套件
- macOS: `.tar.gz` 壓縮包
- ISO: `.iso` 安裝光碟

**注意**: GitHub Release 使用 GitHub REST API 創建，需要配置 `GITHUB_TOKEN` 環境變數

### 2. CI Pipeline

**觸發**: Push to main/dev  
**優先級**: NORMAL  
**用途**: 持續集成檢查

**階段**:
1. Go 基本檢查（vet, fmt, test, build）
2. 前端檢查（type-check, lint, test, build）
3. Docker 建置測試（3 個鏡像）
4. 安全掃描（Trivy）

**質量門檻**:
- ✅ 所有測試通過
- ✅ 代碼格式正確
- ✅ 無高危漏洞

### 3. Kubernetes Deployment

**觸發**: Manual  
**優先級**: HIGH  
**用途**: 部署到 Kubernetes 集群

**階段**:
1. 設置 kubectl
2. 部署微服務（Device/Network/Control）
3. 使用 Helm 部署（可選）

**驗證**:
- Pod 狀態
- Service 狀態
- HPA 狀態

### 4. ArgoCD GitOps Sync

**觸發**: Manual  
**優先級**: NORMAL  
**用途**: GitOps 自動化部署

**階段**:
1. 應用 ArgoCD Application 定義
2. 等待同步完成
3. 檢查應用狀態

### 5. Performance Testing

**觸發**: Manual  
**優先級**: NORMAL  
**用途**: 驗證性能聲明

**階段**:
1. 負載測試（k6）
2. 基準測試（Go Benchmark）

**驗證指標**:
- 吞吐量: 500K req/s
- 延遲: < 2ms P99
- 錯誤率: < 1%

### 6. Security Audit

**觸發**: Manual  
**優先級**: HIGH  
**用途**: 全面安全檢查

**階段**:
1. Trivy 漏洞掃描
2. GoSec 代碼掃描
3. OWASP ZAP Web 掃描

**報告**:
- trivy-report.json
- gosec-report.json
- zap-report.html

### 7. Chaos Engineering

**觸發**: Manual  
**優先級**: NORMAL  
**用途**: 測試系統彈性

**階段**:
1. 安裝 Chaos Mesh
2. 運行混沌測試（Pod 故障、網路延遲）

**測試場景**:
- Pod 隨機終止
- 網路延遲 100ms
- CPU 壓力測試

### 8-11. 雲端部署（已停用）

**狀態**: Disabled  
**原因**: dev 分支專注於地端部署

這些管道保留但停用，不會執行：
- Deploy to GCP
- Deploy to OCI
- Deploy to PaaS
- Terraform Deploy

### 12. ML Model Validation

**觸發**: Manual  
**優先級**: NORMAL  
**用途**: 驗證 AI 模型準確率

**階段**:
1. 運行 ML 驗證測試（Python pytest）
2. 生成準確率報告

**驗證**:
- 深度學習: 99%+ 準確率
- Bot 檢測: 95%+ 準確率
- TLS FP: 98%+ 識別率

### 13. Monitoring & Alerts

**觸發**: Schedule (每 6 小時)  
**優先級**: LOW  
**用途**: 自動健康檢查

**檢查項目**:
- 所有微服務健康狀態
- RabbitMQ 狀態
- Prometheus 狀態
- Grafana 狀態

**通知**: Slack #pandora-monitoring

### 14. Backup & DR

**觸發**: Schedule (每天 2:00 AM)  
**優先級**: NORMAL  
**用途**: 自動備份

**備份內容**:
- PostgreSQL 資料庫
- 配置檔案
- mTLS 證書

### 15. Documentation Build

**觸發**: Push to main/dev  
**優先級**: LOW  
**用途**: 自動生成文檔

**生成**:
- API 文檔（Swagger）
- 文檔網站（Docusaurus）

### 16. Notification Pipeline

**觸發**: Schedule (每週一 9:00 AM)  
**優先級**: LOW  
**用途**: 週報通知

**內容**:
- 項目狀態
- 本週重點
- 最新文檔

### 17. Dependency Updates

**觸發**: Schedule (每週一 0:00 AM)  
**優先級**: LOW  
**用途**: 自動更新依賴

**更新**:
- Go 模組
- NPM 套件
- Docker 鏡像

---

## 🔧 配置步驟

### Step 1: 創建項目

```bash
# 在 Buddy Works 中
1. 點擊 "New Project"
2. 選擇 "GitHub"
3. 選擇倉庫: Local_IPS-IDS
4. 分支: dev (主要) 和 main
```

### Step 2: 導入 buddy.yml

```bash
# Buddy 會自動檢測 buddy.yml
# 或手動導入：
1. 點擊 "Pipelines"
2. 點擊 "Import from YAML"
3. 選擇 buddy.yml
4. 點擊 "Import"
```

### Step 3: 配置整合

#### GitHub 整合

```
Settings → Integrations → GitHub
- 授權 Buddy 訪問倉庫
- 啟用 Webhook
- 配置 Status Checks
```

#### Docker Registry 整合

```
Settings → Integrations → Docker Registry
- Registry: ghcr.io
- Username: $GITHUB_USERNAME
- Password: $GITHUB_TOKEN
```

#### Slack 整合

```
Settings → Integrations → Slack
- Workspace: your-workspace
- Channel: #pandora-monitoring, #pandora-updates
- Webhook URL: $SLACK_WEBHOOK_URL
```

#### Kubernetes 整合

```
Settings → Integrations → Kubernetes
- Cluster URL: https://your-k8s-cluster
- Certificate: (從 kubeconfig 複製)
- Token: (從 kubeconfig 複製)
```

### Step 4: 配置環境變數

```
Settings → Variables
- VERSION: (自動從 Git 取得)
- BUILD_DATE: (自動生成)
- GIT_COMMIT: (自動從 Git 取得)
- GITHUB_TOKEN: (Secrets)
- DOCKER_USERNAME: (Secrets)
- DOCKER_PASSWORD: (Secrets)
- SLACK_WEBHOOK_URL: (Secrets)
```

### Step 5: 配置通知

```
Settings → Notifications
- Email: admin@pandora-ids.com
- Slack: #pandora-updates
- 觸發條件:
  - Pipeline failed
  - Pipeline succeeded (僅 main 分支)
```

---

## 🎯 管道執行順序

### 自動觸發（Push）

```
Push to main/dev
  ↓
[1] CI Pipeline (自動)
  ├── Go 檢查
  ├── 前端檢查
  ├── Docker 建置
  └── 安全掃描
  ↓
[2] Build On-Premise Installers (自動)
  ├── 準備環境
  ├── 構建後端（3 平台）
  ├── 構建前端
  ├── 構建套件
  ├── 構建 ISO
  └── GitHub Release
  ↓
[3] Documentation Build (自動)
  └── 生成文檔
```

### 手動觸發（按需）

```
開發者點擊 "Run"
  ↓
選擇管道:
  • Kubernetes Deployment
  • ArgoCD GitOps Sync
  • Performance Testing
  • Security Audit
  • Chaos Engineering
  • ML Model Validation
```

### 定時觸發（Schedule）

```
每 6 小時:
  • Monitoring & Alerts

每天 2:00 AM:
  • Backup & DR

每週一 0:00 AM:
  • Dependency Updates

每週一 9:00 AM:
  • Notification Pipeline (週報)
```

---

## 📊 管道對照表

### GitHub Actions vs Buddy Works

| GitHub Action | Buddy Pipeline | 觸發方式 | 狀態 |
|---------------|----------------|----------|------|
| build-onpremise-installers.yml | Build On-Premise Installers | Push | ✅ 活躍 |
| ci.yml | CI Pipeline | Push | ✅ 活躍 |
| deploy-gcp.yml | Deploy to GCP | Manual | ❌ 停用 |
| deploy-oci.yml | Deploy to OCI | Manual | ❌ 停用 |
| deploy-paas.yml | Deploy to PaaS | Manual | ❌ 停用 |
| terraform-deploy.yml | Terraform Deploy | Manual | ❌ 停用 |
| - | Kubernetes Deployment | Manual | ✅ 新增 |
| - | ArgoCD GitOps Sync | Manual | ✅ 新增 |
| - | Performance Testing | Manual | ✅ 新增 |
| - | Security Audit | Manual | ✅ 新增 |
| - | Chaos Engineering | Manual | ✅ 新增 |
| - | ML Model Validation | Manual | ✅ 新增 |
| - | Monitoring & Alerts | Schedule | ✅ 新增 |
| - | Backup & DR | Schedule | ✅ 新增 |
| - | Documentation Build | Push | ✅ 新增 |
| - | Notification Pipeline | Schedule | ✅ 新增 |
| - | Dependency Updates | Schedule | ✅ 新增 |

---

## 🎨 Buddy Works 優勢

### vs GitHub Actions

| 特性 | GitHub Actions | Buddy Works | 優勢 |
|------|----------------|-------------|------|
| **視覺化** | YAML 編輯 | 拖放 GUI | ✅ 更直觀 |
| **執行速度** | 中等 | 快速 | ✅ 2-3x 更快 |
| **緩存** | 基本 | 智能 | ✅ 更高效 |
| **並行** | 有限 | 無限 | ✅ 更快構建 |
| **監控** | 基本 | 詳細 | ✅ 更好洞察 |
| **通知** | 基本 | 豐富 | ✅ 多渠道 |
| **定時任務** | Cron | 視覺化 | ✅ 更易配置 |
| **成本** | 免費（有限） | 付費 | ❌ 需要訂閱 |

---

## 🔐 安全最佳實踐

### Secrets 管理

```
Buddy Settings → Variables → Secrets
- 啟用加密
- 限制訪問權限
- 定期輪換
```

### 權限控制

```
Project Settings → Members
- Admin: 完全訪問
- Developer: 運行管道，查看日誌
- Viewer: 僅查看
```

### 審計日誌

```
Settings → Audit Log
- 追蹤所有管道執行
- 追蹤配置變更
- 追蹤用戶操作
```

---

## 📈 監控和報告

### 管道執行統計

```
Dashboard → Statistics
- 成功率
- 平均執行時間
- 失敗原因
- 資源使用
```

### 自訂報告

```
Reports → Custom Reports
- 每週構建報告
- 部署頻率
- 測試覆蓋率趨勢
- 安全掃描結果
```

### Slack 集成

```
每次管道執行後自動發送通知：
- 成功: #pandora-updates
- 失敗: #pandora-alerts
- 安全問題: #pandora-security
```

---

## 🚀 進階功能

### 1. 管道模板

創建可重用的管道模板：

```yaml
# 範例：微服務部署模板
- template: "Microservice Deployment"
  parameters:
    - service_name
    - port
    - replicas
  actions:
    - "Build Docker Image"
    - "Push to Registry"
    - "Deploy to K8s"
    - "Run Health Check"
```

### 2. 條件執行

```yaml
- action: "Deploy to Production"
  type: "BUILD"
  trigger_condition: "$BUDDY_EXECUTION_BRANCH == 'main' && $BUDDY_EXECUTION_TAG != ''"
  execute_commands:
    - "echo 'Deploying to production...'"
```

### 3. 並行執行

```yaml
# 同時構建多個平台
- parallel:
  - "Build Linux"
  - "Build Windows"
  - "Build macOS"
```

### 4. 手動批准

```yaml
- action: "Manual Approval"
  type: "WAIT_FOR_APPROVAL"
  required_approvers: 2
  timeout: 3600
```

---

## 🔧 故障排除

### 管道失敗

```bash
# 1. 查看日誌
Pipelines → [Pipeline Name] → Executions → [Failed Execution] → Logs

# 2. 檢查環境變數
Settings → Variables → 確認所有變數已設置

# 3. 重新運行
點擊 "Retry" 按鈕
```

### Docker 建置失敗

```bash
# 1. 清除緩存
Pipeline Settings → Cache → Clear Cache

# 2. 檢查 Dockerfile
確認路徑正確：build/docker/*.dockerfile

# 3. 檢查 Docker Registry 認證
Settings → Integrations → Docker Registry
```

### Kubernetes 部署失敗

```bash
# 1. 檢查 kubeconfig
Settings → Integrations → Kubernetes

# 2. 驗證 YAML 文件
kubectl apply --dry-run=client -f deployments/kubernetes/

# 3. 檢查集群資源
kubectl get nodes
kubectl top nodes
```

---

## 📋 檢查清單

### 初始設置

- [ ] 創建 Buddy Works 帳號
- [ ] 連接 GitHub 倉庫
- [ ] 導入 buddy.yml
- [ ] 配置環境變數
- [ ] 設置 Docker Registry 整合
- [ ] 設置 Slack 整合
- [ ] 設置 Kubernetes 整合
- [ ] 測試第一個管道

### 日常使用

- [ ] 監控管道執行狀態
- [ ] 查看 Slack 通知
- [ ] 審查安全掃描報告
- [ ] 檢查備份狀態
- [ ] 更新依賴

### 定期維護

- [ ] 每週審查失敗管道
- [ ] 每月審查性能報告
- [ ] 每季度審查安全報告
- [ ] 每半年審查管道配置

---

## 🎯 最佳實踐

### 1. 管道命名

```
✅ 好的命名:
- "Build On-Premise Installers"
- "CI Pipeline"
- "Performance Testing"

❌ 不好的命名:
- "Pipeline 1"
- "Test"
- "Deploy"
```

### 2. 錯誤處理

```yaml
execute_commands:
  - "command1 || echo 'Command1 failed but continuing'"
  - "command2"
  - "command3 || exit 1"  # 失敗時停止
```

### 3. 緩存使用

```yaml
cached_dirs:
  - "/go/pkg/mod"        # Go modules
  - "node_modules"       # NPM packages
  - ".next/cache"        # Next.js cache
```

### 4. 並行化

```yaml
# 同時運行獨立任務
- parallel:
  - "Unit Tests"
  - "Lint Check"
  - "Security Scan"
```

---

## 📚 參考資源

- [Buddy Works 官方文檔](https://buddy.works/docs)
- [YAML 配置參考](https://buddy.works/docs/yaml/yaml-gui)
- [Docker 整合](https://buddy.works/docs/integrations/docker)
- [Kubernetes 整合](https://buddy.works/docs/integrations/kubernetes)
- [Slack 整合](https://buddy.works/docs/integrations/slack)

---

## 🎉 總結

**Buddy Works 配置已完成！**

我們創建了：
- ✅ 12 個主要管道
- ✅ 5 個自動化管道
- ✅ 完整的 buddy.yml 配置
- ✅ 詳細的設置文檔

**優勢**:
- 🚀 更快的構建速度
- 🎨 視覺化管道編輯
- 📊 更好的監控和報告
- 🔔 豐富的通知選項
- 🔄 智能緩存和並行

**下一步**:
1. 在 Buddy Works 創建帳號
2. 連接 GitHub 倉庫
3. 導入 buddy.yml
4. 配置整合和變數
5. 運行第一個管道！

---

**狀態**: ✅ 配置完成  
**檔案**: buddy.yml  
**管道數**: 17 個  
**準備好**: 開始使用 Buddy Works！

**🚀 享受更快、更直觀的 CI/CD 體驗！**

