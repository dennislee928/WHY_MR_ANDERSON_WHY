# 專案結構重整計劃

## 📋 目標

整理專案結構，使其符合 Go 專案最佳實踐和企業標準。

## 🗂️ 新的目錄結構

```
pandora_box_console_IDS-IPS/
├── bin/                        # 編譯產物（不納入版控）
│   ├── pandora-agent
│   ├── pandora-console
│   └── axiom-ui
├── build/                      # 建置相關文件
│   ├── docker/                 # Dockerfile 集中管理
│   │   ├── agent.dockerfile
│   │   ├── agent.koyeb.dockerfile
│   │   ├── console.dockerfile
│   │   ├── monitoring.dockerfile
│   │   ├── nginx.dockerfile
│   │   └── ui.patr.dockerfile
│   └── package/                # 打包腳本
├── cmd/                        # 主程式入口
│   ├── agent/
│   ├── console/
│   └── ui/
├── internal/                   # 私有應用程式代碼
│   ├── handlers/
│   ├── ratelimit/
│   ├── pubsub/
│   ├── mqtt/
│   ├── loadbalancer/
│   └── ...
├── pkg/                        # 公開庫代碼（可選）
├── configs/                    # 配置文件
│   ├── agent/
│   ├── console/
│   ├── grafana/
│   ├── nginx/
│   ├── postgres/
│   └── prometheus/
├── deployments/                # 部署配置集中管理
│   ├── kubernetes/
│   │   ├── base/               # 基礎配置
│   │   ├── overlays/
│   │   │   ├── development/
│   │   │   ├── staging/
│   │   │   └── production/
│   │   ├── gcp/                # GCP 專用
│   │   └── oci/                # OCI 專用
│   ├── terraform/              # Terraform 配置
│   │   ├── environments/
│   │   └── modules/
│   ├── paas/                   # PaaS 平台配置
│   │   ├── flyio/
│   │   ├── koyeb/
│   │   ├── railway/
│   │   ├── render/
│   │   └── patr/
│   └── docker-compose/         # Docker Compose
│       ├── docker-compose.yml
│       └── docker-compose.test.yml
├── scripts/                    # 工具腳本
│   ├── build/
│   ├── deploy/
│   └── test/
├── docs/                       # 文檔集中管理
│   ├── architecture/           # 架構文檔
│   ├── deployment/             # 部署指南
│   │   ├── kubernetes.md
│   │   ├── gcp.md
│   │   ├── oci.md
│   │   └── paas.md
│   ├── development/            # 開發指南
│   ├── operations/             # 運維文檔
│   └── api/                    # API 文檔
├── test/                       # 測試文件
│   ├── integration/
│   ├── e2e/
│   └── fixtures/
├── web/                        # 前端資源
│   ├── components/
│   ├── public/
│   └── styles/
├── .github/                    # GitHub 配置
│   ├── workflows/
│   └── ISSUE_TEMPLATE/
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── LICENSE
```

## 🔄 遷移對照表

### 根目錄清理

| 原位置 | 新位置 | 說明 |
|--------|--------|------|
| `*.exe` | `bin/` | 編譯產物 |
| `Dockerfile.*` | `build/docker/` | Docker 建置文件 |
| `docker-compose.yml` | `deployments/docker-compose/` | Compose 配置 |
| `*-DEPLOYMENT.md` | `docs/deployment/` | 部署文檔 |
| `FLYIO-*.md` | `docs/deployment/flyio/` | Fly.io 文檔 |
| `KOYEB-*.md` | `docs/deployment/koyeb/` | Koyeb 文檔 |
| `*.toml` | `deployments/paas/` | PaaS 配置 |
| `*.yaml` (paas) | `deployments/paas/` | PaaS 配置 |

### K8s 配置重組

| 原位置 | 新位置 |
|--------|--------|
| `k8s/` | `deployments/kubernetes/base/` |
| `k8s-gcp/` | `deployments/kubernetes/gcp/` |
| 新增 | `deployments/kubernetes/oci/` |

### Terraform 重組

| 原位置 | 新位置 |
|--------|--------|
| `terraform/` | `deployments/terraform/` |

### 文檔重組

| 原位置 | 新位置 |
|--------|--------|
| `DOCUMENTS/` | `docs/` |
| `README-*.md` | `docs/` |
| `DEPLOYMENT*.md` | `docs/deployment/` |

## 📝 需要更新的文件

### CI/CD Workflows

1. **.github/workflows/ci.yml**
   - Docker build context 路徑
   - Dockerfile 路徑更新

2. **.github/workflows/deploy-gcp.yml**
   - Dockerfile 路徑
   - K8s manifests 路徑

3. **.github/workflows/deploy-oci.yml**
   - Dockerfile 路徑
   - K8s manifests 路徑

4. **.github/workflows/deploy-paas.yml**
   - Dockerfile 路徑
   - 配置文件路徑

### 建置配置

5. **Makefile**
   - 輸出目錄更新為 `bin/`
   - Docker build 路徑

6. **.gitignore**
   - 添加 `bin/`
   - 更新忽略規則

### 部署配置

7. **Kustomization 文件**
   - 所有 K8s kustomization 路徑

8. **PaaS 配置文件**
   - fly.toml, koyeb.yaml 等

## 🚀 執行步驟

### 階段 1: 創建新目錄結構
- [ ] 創建 `bin/`
- [ ] 創建 `build/docker/`
- [ ] 創建 `docs/`
- [ ] 創建 `deployments/`

### 階段 2: 遷移文件
- [ ] 遷移 Dockerfile
- [ ] 遷移文檔
- [ ] 遷移部署配置
- [ ] 遷移 K8s 配置

### 階段 3: 更新配置
- [ ] 更新 CI/CD workflows
- [ ] 更新 Makefile
- [ ] 更新 .gitignore
- [ ] 更新 README

### 階段 4: 清理
- [ ] 刪除舊文件
- [ ] 驗證所有路徑
- [ ] 測試建置流程

### 階段 5: 驗證
- [ ] 本地建置測試
- [ ] CI/CD 測試
- [ ] 部署測試

## ⚠️ 注意事項

1. **向後兼容**: 保留舊路徑的符號連結（如需要）
2. **文檔同步**: 所有 README 更新路徑引用
3. **團隊通知**: 通知團隊成員路徑變更
4. **分支策略**: 在單獨的分支進行重構
5. **漸進遷移**: 可以分多個 PR 完成

## 📊 影響評估

### 高影響
- CI/CD workflows（必須更新）
- Makefile（必須更新）
- 部署腳本（必須更新）

### 中影響
- 開發工作流程（需要適應新路徑）
- 文檔連結（需要更新）

### 低影響
- 內部代碼（不需要改動）
- Git 歷史（保持完整）

## ✅ 完成標準

- [ ] 所有編譯產物在 `bin/`
- [ ] 所有 Dockerfile 在 `build/docker/`
- [ ] 所有文檔在 `docs/`
- [ ] 所有部署配置在 `deployments/`
- [ ] CI/CD 全部通過
- [ ] 本地建置成功
- [ ] 文檔更新完成
- [ ] README 更新完成

---

**開始日期**: 2025-10-08  
**預計完成**: 2025-10-08  
**負責人**: AI Assistant + User

