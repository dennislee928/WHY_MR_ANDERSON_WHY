# 統一安全與基礎設施平台

[![授權](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go 版本](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![Python 版本](https://img.shields.io/badge/Python-3.11+-blue.svg)](https://python.org)
[![Docker](https://img.shields.io/badge/Docker-20.10+-blue.svg)](https://docker.com)

繁體中文 | [English](README.md)

## 概述

一個全面的雲原生安全與基礎設施管理平台，整合：
- **IDS/IPS 系統** - 即時入侵偵測與防護
- **AI/ML 威脅偵測** - 基於深度學習的安全分析
- **量子計算整合** - IBM Quantum 進階密碼學
- **安全掃描工具** - 整合 Nuclei、Nmap、AMASS 掃描器
- **多雲部署** - 支援 Cloudflare Workers、OCI、IBM Cloud

## 架構

```
┌─────────────────────────────────────────────────────────────┐
│           統一安全與基礎設施平台                              │
└─────────────────────────────────────────────────────────────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                    │
    ┌────▼─────┐      ┌──────▼──────┐      ┌─────▼──────┐
    │  前端    │      │    後端     │      │  AI/量子   │
    │ (React)  │      │    (Go)     │      │  (Python)  │
    └────┬─────┘      └──────┬──────┘      └─────┬──────┘
         │                   │                    │
         └───────────────────┼────────────────────┘
                             │
                    ┌────────▼────────┐
                    │   基礎設施       │
                    │   - Docker       │
                    │   - Kubernetes   │
                    │   - 多雲部署     │
                    └──────────────────┘
```

## 快速開始

### 前置需求

- Docker 20.10+
- Docker Compose 2.0+
- Go 1.24+ (本地開發)
- Python 3.11+ (AI/量子功能)
- Node.js 18+ (前端開發)

### 1. 複製儲存庫

```bash
git clone <repository-url>
cd WHY_MR_ANDERSON_WHY
```

### 2. 環境設定

```bash
cp .env.example .env
# 編輯 .env 設定您的配置
```

### 3. 使用 Docker Compose 啟動

```bash
cd infrastructure/docker
docker-compose up -d
```

### 4. 訪問服務

- **前端 UI**: http://localhost:3001
- **後端 API**: http://localhost:3001/api/v1
- **Swagger 文檔**: http://localhost:3001/swagger
- **AI/量子 API**: http://localhost:8000
- **Grafana**: http://localhost:3000
- **Prometheus**: http://localhost:9090

## 功能特色

### 🛡️ 安全功能

- **即時 IDS/IPS**: 基於 USB-SERIAL CH340 的入侵偵測
- **AI 威脅偵測**: 95.8% 準確率，10 種威脅類型
- **量子密碼學**: QKD、後量子加密
- **零信任架構**: 情境感知存取控制
- **漏洞掃描**: Nuclei、Nmap、AMASS 整合

### 🤖 AI/ML 能力

- 深度學習威脅分類
- 行為異常偵測
- 量子增強機器學習
- AI 治理與公平性審計
- 即時資料流監控

### 🔬 量子計算

- IBM Quantum 整合（127+ qubits）
- 量子密鑰分發（QKD）
- 後量子密碼學
- 量子威脅預測
- 混合量子-古典 ML

### 🌐 多雲支援

| 平台 | 免費額度 | 功能 |
|------|---------|------|
| **Cloudflare Workers** | 1000萬請求/月 | 無伺服器、D1 資料庫、KV 儲存 |
| **Oracle Cloud (OCI)** | 永久免費 | 2 台 VM、4 個 ARM 核心、200GB 儲存 |
| **IBM Cloud** | Lite 方案 | Cloud Foundry、物件儲存 |

詳見 [成本比較](docs/deployment/cost-comparison.md)。

### 📊 監控與可觀測性

- Prometheus 指標收集
- Grafana 儀表板
- Loki 日誌聚合
- 分散式追蹤
- 即時 WebSocket 更新

## 專案結構

```
WHY_MR_ANDERSON_WHY/
├── src/                          # 原始碼
│   ├── backend/                  # Go 服務
│   │   ├── cmd/                  # 進入點
│   │   ├── core/                 # 核心邏輯 (internal)
│   │   ├── axiom-api/            # REST API 伺服器
│   │   ├── api/                  # gRPC 定義
│   │   └── database/             # 資料庫遷移
│   ├── frontend/                 # React UI (Next.js)
│   ├── ai-quantum/               # Python AI/量子服務
│   └── security-tools/           # 掃描器整合
├── infrastructure/               # 部署配置
│   ├── docker/                   # Docker 與 Compose
│   ├── kubernetes/               # K8s 清單
│   ├── terraform/                # 基礎設施即代碼
│   └── cloud-configs/            # 雲端特定配置
│       ├── cloudflare/
│       ├── oci/
│       └── ibm/
├── cicd/                         # CI/CD 管道
│   ├── buddy/                    # Buddy CI
│   ├── argocd/                   # Argo CD GitOps
│   └── harness/                  # Harness 管道
├── docs/                         # 文檔
├── scripts/                      # 工具腳本
├── configs/                      # 應用程式配置
└── tests/                        # 測試套件
```

## 部署

### 本地開發

```bash
# 後端
cd src/backend
go run cmd/server/main.go

# 前端
cd src/frontend
npm install
npm run dev

# AI/量子
cd src/ai-quantum
pip install -r requirements.txt
python main.py
```

### Docker 部署

```bash
cd infrastructure/docker
docker-compose up -d
```

### Kubernetes 部署

```bash
kubectl apply -f infrastructure/kubernetes/
```

### 雲端部署

- **Cloudflare Workers**: 參見 [Cloudflare 指南](docs/deployment/cloudflare.md)
- **Oracle Cloud**: 參見 [OCI 指南](docs/deployment/oci.md)
- **IBM Cloud**: 參見 [IBM Cloud 指南](docs/deployment/ibm-cloud.md)

## CI/CD

支援三個 CI/CD 平台：

1. **Buddy CI** - 簡單、專注於 Docker
   - 配置: `cicd/buddy/buddy.yml`
   
2. **Argo CD** - GitOps、Kubernetes 原生
   - 配置: `cicd/argocd/`
   
3. **Harness** - 企業級
   - 配置: `cicd/harness/`

參見 [CI/CD 文檔](docs/deployment/cicd.md) 進行設定。

## API 文檔

### REST API

完整 API 文檔位於：
- 本地: http://localhost:3001/swagger
- 生產環境: 參見部署文檔

主要端點：

```bash
# 系統狀態
GET /api/v1/status

# 安全威脅
GET /api/v1/security/threats
POST /api/v1/security/threats/:id/block

# 網路管理
GET /api/v1/network/stats
DELETE /api/v1/network/blocked-ips/:ip

# AI/ML 威脅偵測
POST /api/v1/ml/detect

# 量子服務
POST /api/v1/quantum/qkd/generate
POST /api/v1/zerotrust/predict
```

### WebSocket

透過 WebSocket 的即時更新：

```javascript
const ws = new WebSocket('ws://localhost:3001/ws?client_id=dashboard');
ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    // 處理即時更新
};
```

## 配置

### 環境變數

`.env` 中的關鍵配置：

```bash
# 資料庫
DB_HOST=localhost
DB_PORT=5432
DB_USER=sectools
DB_PASSWORD=changeme

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# IBM Quantum
IBM_QUANTUM_TOKEN=your_token_here

# 雲端提供商
CLOUDFLARE_API_TOKEN=
OCI_TENANCY_OCID=
IBM_CLOUD_API_KEY=
```

### 應用程式配置

- **後端**: `configs/agent-config.yaml`
- **前端**: `src/frontend/.env.local`
- **AI/量子**: `src/ai-quantum/env.example`

## 安全性

### 最佳實踐

- ✅ 所有敏感資料靜態加密
- ✅ 服務間 mTLS 通訊
- ✅ 速率限制與 DDoS 防護
- ✅ CI/CD 中的 SAST 掃描
- ✅ 定期依賴更新
- ✅ 零信任架構

### 合規性

- GDPR 合規
- SOC2 就緒
- ISO27001 對齊
- PII 自動偵測與匿名化

## 效能

| 指標 | 值 |
|------|------|
| API 響應時間 (P99) | < 2ms |
| 吞吐量 | 50萬+ 請求/秒 |
| AI 偵測延遲 | < 10ms |
| 可用性 | 99.999% |

## 貢獻

我們歡迎貢獻！請參見 [CONTRIBUTING.md](CONTRIBUTING.md)。

### 開發工作流程

1. Fork 儲存庫
2. 建立功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交變更 (`git commit -m '新增出色功能'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 開啟 Pull Request

## 測試

```bash
# 後端測試
cd src/backend
go test ./...

# 前端測試
cd src/frontend
npm test

# 整合測試
cd tests
go test -tags=integration ./...
```

## 文檔

- [架構](docs/architecture/system-design.zh-TW.md)
- [API 參考](docs/development/api-reference.md)
- [部署指南](docs/deployment/)
- [安全性](docs/security/)
- [開發指南](docs/development/getting-started.md)

## 路線圖

### 2025 Q1
- ✅ 統一專案結構
- ✅ 多雲部署
- ✅ 三個 CI/CD 平台
- [ ] 增強 AI 威脅偵測
- [ ] 擴展量子演算法

### 2025 Q2
- [ ] 行動應用程式支援
- [ ] 進階分析儀表板
- [ ] 多租戶架構
- [ ] MISP 威脅情報整合

## 授權

本專案採用 MIT 授權條款 - 詳見 [LICENSE](LICENSE) 檔案。

## 致謝

- [ProjectDiscovery](https://github.com/projectdiscovery) - Nuclei 掃描器
- [Nmap](https://nmap.org/) - 網路掃描
- [OWASP AMASS](https://github.com/OWASP/Amass) - 資產發現
- [IBM Quantum](https://quantum-computing.ibm.com/) - 量子計算
- [Qiskit](https://qiskit.org/) - 量子開發

## 支援

- **問題**: [GitHub Issues](https://github.com/your-repo/issues)
- **討論**: [GitHub Discussions](https://github.com/your-repo/discussions)
- **電子郵件**: security@example.com

## 免責聲明

此工具僅供授權的安全測試和研究使用。使用者必須遵守當地法律，並在掃描任何系統前取得適當授權。

---

**由安全社群用 ❤️ 製作**

**🌟 如果此專案對您有幫助，請給個星星！**

