# 變更日誌

所有重要的專案變更都會記錄在此檔案。

格式基於 [Keep a Changelog](https://keepachangelog.com/zh-TW/1.0.0/)，  
版本號遵循 [語義化版本](https://semver.org/lang/zh-TW/)。

---

## [3.0.0] - 2025-10-09

### 🎉 重大更新 - 地端部署版本

此版本將專案完全重構為地端部署版本，新增完整的應用程式結構和自動化構建系統。

### 新增 (Added)

#### Application/ 目錄結構
- **Application/Fe/** - 完整的 Next.js 14 前端應用程式
  - 25+ 個新檔案
  - 7 個 UI 組件（Card, Button, Badge, Loading, Alert, MainLayout, Dashboard）
  - 完整的 API 服務層
  - 2 個自定義 Hooks（useSystemStatus, useWebSocket）
  - TypeScript 類型系統（10+ 介面）
  - 工具函數庫（9+ 函數）
  - TailwindCSS 自定義主題
  
- **Application/be/** - 後端構建系統
  - Makefile（17 個目標）
  - Windows 構建腳本（build.ps1）
  - Linux/macOS 構建腳本（build.sh）
  - 完整的引用結構
  
- **Application/build-local.*** - 主構建腳本
  - Windows PowerShell 版本
  - Linux/macOS Bash 版本
  - 自動生成啟動/停止腳本
  - 版本資訊嵌入

#### CI/CD 自動化
- **build-onpremise-installers.yml** - 安裝檔構建 workflow
  - 支援 Windows/Linux/macOS
  - 支援 amd64/arm64
  - 生成 .exe/.deb/.rpm/.iso/.ova
  - 自動發布到 GitHub Releases

#### 安裝檔生成系統
- **build/installer/windows/** - Windows 安裝程式資源
  - Inno Setup 腳本範本
  - 系統服務配置
  
- **build/installer/linux/** - Linux 套件資源
  - Debian postinst/prerm 腳本
  - Systemd 服務檔案
  
- **build/installer/iso/** - ISO 光碟資源
  - 自動安裝腳本

#### 文檔系統
- `RESTRUCTURE-FINAL-REPORT.md` - 完整重構報告
- `QUICK-START-GUIDE.md` - 3分鐘快速入門
- `ONPREMISE-DEPLOYMENT-GUIDE.md` - 詳細部署指南
- `TESTING-CHECKLIST.md` - 完整測試清單
- `PROJECT-RESTRUCTURE-PROGRESS.md` - 進度追蹤
- `CHANGELOG.md` - 本檔案
- 各子目錄 README.md（6+個）

### 變更 (Changed)

#### 專案結構
- 移動舊 `web/` → `Application/Fe/legacy/web-old/`
- 移動舊 `DOCUMENTS/` → `docs/archive/DOCUMENTS-old/`
- 移動舊 `k8s/` → `deployments/kubernetes/legacy/`
- 移動根目錄 Dockerfiles → `build/docker/`

#### CI/CD Workflows
- 更新 `ci.yml` - 支援 dev 分支，使用 Application/Fe/
- 停用雲端部署 workflows（deploy-gcp, deploy-oci, deploy-paas, terraform-deploy）

#### 配置檔案
- 更新 `.gitignore` - 新增 Application/ 相關規則
- 新增 `.editorconfig` - 統一編碼風格
- 更新 `README.md` - 完整的地端部署說明
- 更新 `README-PROJECT-STRUCTURE.md` - 新結構說明

#### 構建系統
- 更新 `Application/be/Makefile`
  - 修正路徑引用
  - 新增 copy-to-dist 目標
  - 新增 info 目標
- 改進 `Application/build-local.ps1`
- 改進 `Application/build-local.sh`

### 棄用 (Deprecated)

- 舊版 `web/` 目錄（已移至 legacy）
- 根目錄的 Dockerfiles（已移至 build/docker/）
- 舊版 k8s 配置（已移至 legacy）

### 移除 (Removed)

無。所有檔案都已歸檔而非刪除，確保可追溯性。

### 修復 (Fixed)

- 修正 Makefile 路徑引用錯誤
- 修正 CI workflow 前端路徑
- 修正 build-onpremise-installers.yml 構建命令
- 修正 Dashboard 組件 import 路徑

### 安全性 (Security)

- 新增 .env.local 到 .gitignore
- 新增 systemd 安全性設定
- 新增安裝程式權限檢查

---

## [2.0.0] - 2025-10-08

### 變更
- 重整專案結構（雲端部署版本）
- 詳見 PROJECT-RESTRUCTURE-FINAL-REPORT.md

---

## [1.0.0] - 初始版本

### 新增
- 基礎 IDS-IPS 功能
- Agent 和 Console 程式
- 基礎監控整合

---

## 版本對照

| 版本 | 日期 | 分支 | 部署方式 | 說明 |
|------|------|------|----------|------|
| v3.0.0 | 2025-10-09 | dev | 地端部署 | 完整重構，Application/ 結構 |
| v2.0.0 | 2025-10-08 | main | 雲端部署 | PaaS 平台整合 |
| v1.0.0 | 初始 | main | 混合 | 初始版本 |

---

## 升級指南

### 從 v2.0.0 升級到 v3.0.0

**重要**: v3.0.0 是全新的地端部署版本，與 v2.0.0 (雲端版本) 並行發展。

如果您：
- **使用雲端部署**: 保持在 `main` 分支和 v2.x.x
- **要轉換為地端**: 切換到 `dev` 分支和 v3.x.x

**遷移步驟**:

1. 備份現有資料
2. 切換分支: `git checkout dev`
3. 執行本地構建或使用安裝檔
4. 遷移配置和資料
5. 測試驗證

---

## 貢獻者

- **專案負責人**: Pandora Security Team
- **重構執行**: AI Assistant (Claude Sonnet 4.5)
- **審核者**: Dennis Lee

---

**維護**: 持續更新  
**支援**: support@pandora-ids.com  
**文檔**: [完整文檔](README.md)

