# Pandora Box Console IDS/IPS - 專案結構說明

> **📌 重要**: 本文檔說明 `dev` 分支的專案結構，專為**地端部署（On-Premise）**設計。  
> 雲端部署版本請參考 `main` 分支的文檔。

---

## 📁 專案結構概覽（dev 分支）

```
pandora_box_console_IDS-IPS/  (dev 分支 - v3.0.0 AI 驅動智能安全平台)
│
├── Application/                # 🆕 地端應用程式主目錄
│   ├── be/                     # 後端應用程式
│   │   ├── Makefile            # 後端構建腳本
│   │   └── README.md           # 後端使用說明
│   ├── Fe/                     # 前端應用程式
│   │   ├── components/         # React 組件
│   │   ├── pages/              # Next.js 頁面
│   │   ├── public/             # 靜態資源
│   │   ├── styles/             # 樣式文件
│   │   ├── package.json        # NPM 套件定義
│   │   ├── next.config.js      # Next.js 配置
│   │   ├── tsconfig.json       # TypeScript 配置
│   │   └── README.md           # 前端使用說明
│   ├── build-local.ps1         # Windows 本地構建腳本
│   ├── build-local.sh          # Linux/macOS 本地構建腳本
│   ├── dist/                   # 構建產物（不納入版控）
│   │   ├── backend/            # 編譯後的後端程式
│   │   ├── frontend/           # 構建後的前端程式
│   │   ├── start.bat/.sh       # 啟動腳本
│   │   ├── stop.bat/.sh        # 停止腳本
│   │   └── README.txt          # 使用說明
│   └── README.md               # 應用程式主說明
│
├── api/                        # 🆕 API 定義（Phase 1）
│   └── proto/                  # gRPC Protocol Buffers
│       ├── common.proto        # 共享類型定義
│       ├── device.proto        # Device Service API（6 RPCs）
│       ├── network.proto       # Network Service API（7 RPCs）
│       ├── control.proto       # Control Service API（9 RPCs）
│       ├── Makefile            # Proto 代碼生成
│       └── README.md           # API 文檔
│
├── bin/                        # 編譯產物（舊版，保留相容性）
│   ├── pandora-agent.exe
│   ├── pandora-console.exe
│   └── axiom-ui.exe
│
├── build/                      # 建置相關文件
│   ├── docker/                 # Dockerfile 集中管理
│   │   ├── agent.dockerfile
│   │   ├── agent.koyeb.dockerfile
│   │   ├── monitoring.dockerfile
│   │   ├── nginx.dockerfile
│   │   ├── server-be.dockerfile
│   │   ├── server-fe.dockerfile
│   │   └── ui.patr.dockerfile
│   └── package/                # 打包腳本
│
├── cmd/                        # 主程式入口
│   ├── agent/                  # Agent 主程式
│   │   └── main.go
│   ├── console/                # Console 主程式
│   │   └── main.go
│   └── ui/                     # UI 主程式
│       └── main.go
│
├── internal/                   # 私有應用程式代碼
│   ├── axiom/                  # Axiom 引擎 (含 Swagger 整合)
│   ├── device/                 # 裝置管理
│   ├── grafana/                # Grafana 整合
│   ├── handlers/               # HTTP 處理器
│   ├── loadbalancer/           # 負載均衡器
│   ├── logging/                # 日誌系統
│   ├── metrics/                # 指標收集
│   ├── mqtt/                   # MQTT 客戶端
│   ├── mtls/                   # mTLS 支援
│   ├── network/                # 網路管理
│   ├── pin/                    # PIN 碼系統
│   ├── pubsub/                 # 🆕 RabbitMQ 消息隊列 (完整整合)
│   ├── ratelimit/              # 速率限制
│   ├── security/               # 安全相關
│   ├── token/                  # Token 認證
│   ├── tpm/                    # TPM 支援
│   └── utils/                  # 工具函數
│
├── configs/                    # 配置文件
│   ├── agent/                  # Agent 配置
│   ├── console/                # Console 配置
│   ├── grafana/                # Grafana 配置
│   ├── nginx/                  # Nginx 配置
│   ├── postgres/               # PostgreSQL 配置
│   ├── prometheus/             # Prometheus 配置
│   └── rabbitmq/               # 🆕 RabbitMQ 配置
│
├── deployments/                # 部署配置集中管理
│   ├── kubernetes/             # Kubernetes 部署
│   │   ├── base/               # 基礎配置（通用）
│   │   ├── gcp/                # GCP 專用配置
│   │   └── oci/                # OCI 專用配置
│   ├── terraform/              # Terraform IaC
│   │   ├── environments/       # 環境配置
│   │   │   ├── dev/
│   │   │   ├── staging/
│   │   │   └── prod/
│   │   └── modules/            # Terraform 模組
│   ├── paas/                   # PaaS 平台配置
│   │   ├── flyio/              # Fly.io
│   │   ├── koyeb/              # Koyeb
│   │   ├── railway/            # Railway
│   │   ├── render/             # Render
│   │   └── patr/               # Patr.io
│   └── docker-compose/         # Docker Compose
│       ├── docker-compose.yml
│       └── docker-compose.test.yml
│
├── scripts/                    # 工具腳本
│   ├── build/                  # 建置腳本
│   ├── deploy/                 # 部署腳本
│   ├── test/                   # 測試腳本
│   └── restructure-project.ps1 # 專案重整腳本
│
├── docs/                       # 文檔集中管理
│   ├── architecture/           # 架構文檔
│   │   └── modules.md
│   ├── deployment/             # 部署指南
│   │   ├── README.md
│   │   ├── quickstart.md
│   │   ├── kubernetes.md
│   │   ├── gcp.md
│   │   ├── oci.md
│   │   ├── paas.md
│   │   ├── terraform-summary.md
│   │   ├── flyio/
│   │   └── koyeb/
│   ├── development/            # 開發指南
│   │   ├── windows-setup.md
│   │   └── implementation-summary.md
│   ├── operations/             # 運維文檔
│   │   ├── final-status.md
│   │   └── fixes-summary.md
│   ├── PROJECT-RESTRUCTURE-PLAN.md
│   ├── RESTRUCTURE-STATUS.md
│   ├── RESTRUCTURE-SUMMARY.md
│   └── CI-CD-UPDATE-GUIDE.md
│
├── test/                       # 測試文件
│   ├── integration/            # 整合測試
│   └── fixtures/               # 測試固件
│
├── web/                        # 前端資源
│   ├── components/             # React 組件
│   ├── public/                 # 靜態資源
│   └── styles/                 # 樣式文件
│
├── .github/                    # GitHub 配置
│   └── workflows/              # GitHub Actions
│       ├── ci.yml
│       ├── deploy-gcp.yml
│       ├── deploy-oci.yml
│       └── deploy-paas.yml
│
├── .gitignore                  # Git 忽略規則
├── go.mod                      # Go 模組定義
├── go.sum                      # Go 依賴鎖定
├── Makefile                    # Make 建置腳本
├── README.md                   # 專案說明
└── LICENSE                     # 授權條款
```

---

## 📖 目錄說明

### 🆕 地端應用程式 (`Application/`)

這是 **dev 分支的核心目錄**，包含完整的地端部署應用程式。

#### `Application/be/` - 後端應用程式

後端程式的獨立構建和開發環境。

- **Makefile**: 提供各種構建目標（all, agent, console, ui, test, package）
- **引用機制**: 透過相對路徑引用專案根目錄的 `cmd/` 和 `internal/`
- **構建產物**: 輸出到 `Application/dist/backend/`

#### `Application/Fe/` - 前端應用程式

基於 Next.js 14 的現代化前端應用程式。

- **package.json**: 定義所有前端依賴和構建腳本
- **next.config.js**: Next.js 配置，支援獨立部署
- **tsconfig.json**: TypeScript 配置
- **構建產物**: 輸出到 `Application/dist/frontend/`

#### `Application/build-local.*` - 本地構建腳本

自動化構建工具，支援 Windows 和 Linux/macOS。

- **功能**:
  - 一鍵構建前後端
  - 自動生成啟動/停止腳本
  - 版本資訊嵌入
  - 產物打包
- **使用**:
  ```powershell
  # Windows
  .\build-local.ps1 -Version "1.0.0"
  ```
  ```bash
  # Linux/macOS
  ./build-local.sh all "1.0.0"
  ```

#### `Application/dist/` - 構建產物

完整的可部署套件，包含：
- 所有編譯後的二進位檔案
- 前端靜態資源
- 配置檔案
- 啟動/停止腳本
- 使用說明

### 原始碼

#### `cmd/`
應用程式的入口點。每個子目錄對應一個可執行程式。

- **為什麼這樣組織**：遵循 Go 專案標準佈局
- **命名規範**：與編譯產物名稱對應
- **dev 分支變更**: 現在主要透過 `Application/be/` 構建

#### `internal/`
私有應用程式和庫代碼。其他專案無法導入這些包。

- **為什麼使用 `internal/`**：Go 的特殊目錄，強制執行包可見性
- **組織原則**：按功能領域劃分（不是按層次）
- **dev 分支變更**: 保持不變，由 `Application/be/` 引用

### 配置和資源

#### `configs/`
所有配置文件的集中位置。

- **組織方式**：按服務分組（grafana/, nginx/, 等）
- **環境變數**：使用 `.env` 文件或環境變數覆蓋
- **dev 分支變更**: 構建時會複製到 `Application/dist/backend/configs/`

#### `web/`
前端資源和靜態文件（舊版，保留）。

- **框架**：Next.js / React
- **dev 分支變更**: 新的前端開發在 `Application/Fe/` 進行

### 建置和部署

#### `bin/`
編譯後的可執行文件。

- **不納入版控**：在 `.gitignore` 中排除
- **命名規範**：`pandora-{service}.exe`

#### `build/`
建置腳本和 Dockerfile。

- **docker/**：所有 Dockerfile 集中管理
- **package/**：打包和發布腳本

#### `deployments/`
所有部署配置的集中位置。

**為什麼這樣組織**：
- 清晰分離不同的部署目標
- 避免根目錄混亂
- 便於維護和查找

**子目錄結構**：
- `kubernetes/`：K8s manifests
  - `base/`：基礎配置，可被 overlay 使用
  - `gcp/`, `oci/`：雲服務商特定配置
- `terraform/`：基礎設施即代碼
- `paas/`：各 PaaS 平台配置
- `docker-compose/`：本地開發和測試

### 文檔

#### `docs/`
所有專案文檔的集中位置。

**組織原則**：
- 按讀者類型分類（開發/部署/運維）
- 子目錄按主題組織
- 保持扁平化（避免過深嵌套）

**文檔類型**：
- **architecture/**：系統設計和架構決策
- **deployment/**：部署指南和操作手冊
- **development/**：開發環境設置和指南
- **operations/**：運維和故障排除

### 測試

#### `test/`
整合測試和端到端測試。

- **單元測試**：與代碼放在一起（`*_test.go`）
- **整合測試**：放在 `test/integration/`
- **固件數據**：放在 `test/fixtures/`

### 腳本

#### `scripts/`
各種自動化腳本。

**組織方式**：
- `build/`：建置相關腳本
- `deploy/`：部署相關腳本  
- `test/`：測試相關腳本

---

## 🎯 設計原則

### 1. 關注點分離
- 原始碼、配置、部署、文檔各自獨立
- 避免混合不同類型的文件

### 2. 標準化
- 遵循 Go 專案標準佈局
- 遵循 Cloud Native 最佳實踐
- 遵循 12-Factor App 原則

### 3. 可發現性
- 清晰的命名
- 一致的組織結構
- 完善的文檔

### 4. 可維護性
- 邏輯分組
- 避免深層嵌套
- 保持目錄結構簡潔

---

## 🌐 部署架構

### dev 分支：地端部署

本分支（`dev`）專為**地端部署（On-Premise）**設計，所有服務運行在本地伺服器。

#### 部署方式

| 方式 | 說明 | 適用場景 |
|------|------|----------|
| **預建安裝檔** | 使用 CI/CD 生成的 .exe/.deb/.rpm/.iso 檔案 | 生產環境、快速部署 |
| **本地構建** | 使用 `build-local.*` 腳本構建 | 開發、測試、自訂部署 |
| **Docker Compose** | 使用 docker-compose.yml | 容器化部署 |
| **手動構建** | 直接使用 Go 和 npm 命令 | 進階開發 |

#### 安裝檔格式

透過 GitHub Actions 自動構建多種格式：

| 格式 | 平台 | 特性 |
|------|------|------|
| `.exe` | Windows | Inno Setup 安裝精靈 |
| `.deb` | Debian/Ubuntu | systemd 服務自動配置 |
| `.rpm` | RHEL/CentOS | RPM 套件管理 |
| `.iso` | All | 可開機安裝光碟 |
| `.ova` | All | 虛擬機映像檔 |

#### 本地服務架構

```
本地伺服器
├── Backend Services
│   ├── Axiom BE (Port 3001) - 獨立 Go API 服務
│   │   └── 29+ REST API + Swagger + WebSocket
│   ├── Pandora Agent (Host Network)
│   │   └── 核心 IDS/IPS 引擎
│   └── Cyber AI/Quantum (Port 8000) - Python 服務
│       ├── ML 威脅檢測
│       ├── Zero Trust 量子預測
│       ├── IBM Quantum 整合
│       └── 進階量子算法 (QSVM/QAOA/QWalk)
│
├── Message Queue (Port 5672)
│   └── RabbitMQ (完整事件流整合)
│
├── Monitoring Stack
│   ├── Prometheus (Port 9090) - 指標收集
│   ├── Grafana (Port 3000) - 視覺化
│   ├── Loki (Port 3100) - 日誌聚合
│   ├── AlertManager (Port 9093) - 告警管理
│   └── Node Exporter (Port 9100) - 系統指標
│
├── Storage Layer
│   ├── PostgreSQL (Port 5432) - 關聯資料庫
│   └── Redis (Port 6379) - 快取系統
│
└── Infrastructure
    └── Nginx (Port 443) - 反向代理
```

### main 分支：雲端部署（參考）

主分支採用多平台混合部署策略：

| 平台 | 服務 | URL | 用途 |
|------|------|-----|------|
| **Koyeb** | Pandora Agent | `https://dizzy-sher-mitake-7f13854a.koyeb.app:8080` | 主應用程式 |
| **Fly.io** | Monitoring Stack | `https://pandora-monitoring.fly.dev` | 監控系統 |
| **Render** | Redis + Nginx | `redis-7-2-11-alpine3-21.onrender.com` | 資料與代理 |

## 📚 相關文檔

### 地端部署文檔（dev 分支）
- [應用程式說明](Application/README.md) ⭐ **重要**
- [後端開發指南](Application/be/README.md)
- [前端開發指南](Application/Fe/README.md)
- [主要 README](README.md)

### CI/CD 文檔
- [地端安裝檔構建](.github/workflows/build-onpremise-installers.yml)
- [CI Pipeline](.github/workflows/ci.yml)

### 雲端部署文檔（main 分支參考）
- [PaaS 整合指南](docs/deployment/paas-integration-guide.md)
- [Koyeb 部署指南](docs/deployment/koyeb/README.md)
- [Fly.io Volume 調整](docs/deployment/flyio/FLYIO-VOLUME-FIX.md)

### 開發文檔
- [專案重整計劃](docs/PROJECT-RESTRUCTURE-PLAN.md)
- [重整狀態報告](docs/RESTRUCTURE-STATUS.md)
- [CI/CD 更新指南](docs/CI-CD-UPDATE-GUIDE.md)

---

## 🔄 版本歷史

- **v3.0.0** (2025-10-09): AI 智能化與企業級優化 🎉
  - ✅ 深度學習威脅檢測（99%+ 準確率）
  - ✅ 行為基線建模和異常檢測
  - ✅ Jaeger 分散式追蹤
  - ✅ 智能緩存系統（95%+ 命中率）
  - ✅ 多租戶 SaaS 架構
  - ✅ 合規性報告和 SLA 管理
  - ✅ RabbitMQ 完整事件流整合
  - ✅ Swagger API 文檔整合
  - ✅ 17+ 新 API 端點
  - 新增 `internal/ml/`, `internal/tracing/`, `internal/cache/`, `internal/multitenant/`
  - 新增 `configs/rabbitmq/`, `examples/rabbitmq-integration/`
  
- **v2.0.0** (2025-10-09): Kubernetes 與自動化
  - ✅ Kubernetes 雲原生部署
  - ✅ Helm Charts + ArgoCD GitOps
  - ✅ ML Bot 檢測和 TLS Fingerprinting
  - ✅ WAF 防護和自動威脅響應
  - 新增 `deployments/kubernetes/`, `deployments/helm/`, `deployments/argocd/`
  - 新增 `internal/discovery/`, `internal/security/`, `internal/automation/`
  
- **v1.0.0** (2025-10-09): 微服務架構重構
  - ✅ 3 個獨立微服務（Device/Network/Control）
  - ✅ RabbitMQ 消息隊列 + gRPC 通訊
  - ✅ mTLS 安全認證 + 硬體整合
  - 新增 `api/proto/`, `cmd/*-service/`, `internal/pubsub/`, `internal/services/`
  - 新增 `internal/grpc/`, `internal/resilience/`, `internal/ratelimit/`
  
- **v0.1.0** (2024-12-19): 初始版本
  - 單體架構
  - 基礎 IDS/IPS 功能

---

## 🎯 分支策略

| 分支 | 用途 | 部署方式 |
|------|------|----------|
| `main` | 雲端生產環境 | PaaS 平台（Koyeb, Fly.io, Render） |
| `dev` | 地端部署版本 | 安裝檔、本地構建、Docker |
| `staging` | 預發布測試 | 視需求 |

---

## 📋 快速參考

### 構建命令

```bash
# Windows 本地構建
cd Application
.\build-local.ps1

# Linux/macOS 本地構建
cd Application
./build-local.sh

# 後端構建（Make）
cd Application/be
make all

# 前端構建（npm）
cd Application/Fe
npm run build
```

### 啟動服務

```bash
# 使用構建產物
cd Application/dist
.\start.bat  # Windows
./start.sh   # Linux/macOS

# 使用安裝檔
sudo systemctl start pandora-agent  # Linux
```

### CI/CD 觸發

```bash
# 推送到 dev 分支觸發 CI 和安裝檔構建
git push origin dev

# 創建版本標籤觸發 Release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

---

**維護者**: Pandora Security Team  
**技術支援**: support@pandora-ids.com  
**最後更新**: 2025-10-09（dev 分支地端部署版本）
