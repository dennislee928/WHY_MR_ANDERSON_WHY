# Commit Message (建議使用)

```
feat: 完成專案重構 - 地端部署版本 v3.0.0

## 重大更新

將專案完全重構為地端部署版本，新增 Application/ 完整應用程式結構。

## 新增功能

### Application/ 目錄結構
- Application/Fe/ - 完整的 Next.js 14 前端應用程式
  - 28個檔案：components, pages, services, hooks, types
  - UI組件系統：Card, Button, Badge, Loading, Alert
  - 響應式佈局：MainLayout
  - API服務層：完整的 REST 和 WebSocket 支援
  - TailwindCSS 自定義主題

- Application/be/ - 後端構建系統
  - Makefile：17個目標
  - build.ps1：Windows構建腳本
  - build.sh：Linux/macOS構建腳本

- Application/build-local.* - 主構建腳本
  - 一鍵構建前後端
  - 自動生成啟動/停止腳本
  - 版本資訊嵌入

### CI/CD 自動化
- build-onpremise-installers.yml - 安裝檔構建workflow
  - 支援多平台：Windows/Linux/macOS
  - 支援多架構：amd64/arm64
  - 生成多格式：.exe/.deb/.rpm/.iso/.ova
  - 自動發布到 GitHub Releases

### 安裝檔生成系統
- build/installer/ - 完整的安裝檔資源
  - Windows：Inno Setup腳本
  - Linux：Debian/RPM套件腳本
  - ISO：自動安裝腳本
  - Systemd：服務配置

### 文檔系統
- docs/onpremise/ - 地端部署文檔
- docs/development/ - 開發指南
- docs/cicd/ - CI/CD文檔
- docs/archive/ - 存檔文檔

## 變更

### 專案結構重組
- 移動 web/ → Application/Fe/legacy/web-old/
- 移動 DOCUMENTS/ → docs/archive/DOCUMENTS-old/
- 移動 k8s/ → deployments/kubernetes/legacy/
- 清理根目錄 Dockerfiles → build/docker/

### CI/CD Workflows
- 更新 ci.yml：支援dev分支，使用Application/Fe/
- 停用雲端部署workflows：僅手動觸發
- 修正所有路徑引用

### 配置和規範
- 更新 .gitignore：新增Application/相關規則
- 新增 .editorconfig：統一編碼風格
- 更新 README.md：完整的地端部署說明
- 更新 README-PROJECT-STRUCTURE.md：新結構說明

## 刪除

### 清理重複檔案
- build/docker/Dockerfile.* (8個) → 保留 *.dockerfile

### 清理Terraform產物
- terraform/.terraform/
- deployments/terraform/.terraform/
- 所有 .terraform.lock.hcl

## 修復

- 修正 Makefile 路徑引用
- 修正 CI workflow 前端路徑
- 修正安裝檔構建命令
- 修正 Dashboard 組件 import 路徑

## 安全性

- 新增 .env.local 到 .gitignore
- 新增 systemd 安全性設定
- 新增安裝程式權限檢查

---

## 統計數據

- **新建檔案**: 61+
- **修改檔案**: 15+
- **刪除檔案**: 20+
- **程式碼行數**: ~5000+
- **文檔頁面**: 18+

---

## 測試狀態

- [x] 結構驗證通過
- [ ] 本地構建測試（待執行）
- [ ] CI/CD測試（待觸發）

---

## 相關Issue/PR

- Closes #XXX
- Related to #YYY

---

## Breaking Changes

此版本為地端部署專用版本，與v2.0.0（雲端版本）架構不同。

如需雲端部署，請使用 main 分支。

---

## 後續工作

- [ ] 測試本地構建
- [ ] 測試CI/CD自動構建
- [ ] 創建第一個正式Release (v3.0.0)
- [ ] 更新部署文檔

---

Co-authored-by: Pandora Security Team <support@pandora-ids.com>
```

---

## 簡短版本（如果上面太長）

```
feat: 完成專案重構 - 地端部署版本 v3.0.0

- 新增 Application/ 完整應用程式結構
  - Fe/: Next.js 14 前端（28檔案）
  - be/: Go 後端構建系統
  - build-local.*: 一鍵構建腳本

- 新增 CI/CD 安裝檔生成
  - 支援 .exe/.deb/.rpm/.iso/.ova
  - 多平台多架構

- 清理專案結構
  - 移除重複 Dockerfiles (8個)
  - 清理 .terraform 目錄
  - 整理文檔到 docs/

- 完整文檔系統 (18+ 文檔)

總計：新增 61+ 檔案，清理 20+ 檔案

詳見：docs/RESTRUCTURE-FINAL-REPORT.md
```
