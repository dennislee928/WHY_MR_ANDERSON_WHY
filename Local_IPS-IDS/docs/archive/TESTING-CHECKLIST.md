# 專案重構測試清單

> **更新時間**: 2025-10-09  
> **狀態**: ✅ 準備測試

---

## 📋 結構驗證

### Application/ 目錄

- [x] Application/Fe/ 存在
- [x] Application/be/ 存在
- [x] Application/build-local.ps1 存在
- [x] Application/build-local.sh 存在
- [x] Application/README.md 存在

### 前端檔案（Application/Fe/）

#### 配置檔案
- [x] package.json
- [x] tsconfig.json
- [x] next.config.js
- [x] tailwind.config.js
- [x] postcss.config.js
- [x] .eslintrc.json
- [x] .gitignore
- [x] .env.example

#### 程式碼檔案
- [x] pages/index.tsx
- [x] pages/_app.tsx
- [x] pages/_document.tsx
- [x] components/ui/card.tsx
- [x] components/ui/button.tsx
- [x] components/ui/badge.tsx
- [x] components/ui/loading.tsx
- [x] components/ui/alert.tsx
- [x] components/layout/MainLayout.tsx
- [x] components/dashboard/Dashboard.tsx
- [x] services/api.ts
- [x] hooks/useSystemStatus.ts
- [x] hooks/useWebSocket.ts
- [x] types/index.ts
- [x] lib/utils.ts
- [x] styles/globals.css

### 後端檔案（Application/be/）

- [x] Makefile
- [x] build.ps1
- [x] build.sh
- [x] go.mod
- [x] README.md

### CI/CD Workflows

- [x] .github/workflows/ci.yml（已更新）
- [x] .github/workflows/build-onpremise-installers.yml（已創建）
- [x] .github/workflows/deploy-gcp.yml（已停用）
- [x] .github/workflows/deploy-oci.yml（已停用）
- [x] .github/workflows/deploy-paas.yml（已停用）
- [x] .github/workflows/terraform-deploy.yml（已停用）

### 安裝檔資源

- [x] build/installer/windows/setup-template.iss
- [x] build/installer/linux/postinst.sh
- [x] build/installer/linux/prerm.sh
- [x] build/installer/linux/systemd/pandora-agent.service
- [x] build/installer/iso/install.sh
- [x] build/installer/README.md

### 文檔

- [x] README.md（已更新）
- [x] README-PROJECT-STRUCTURE.md（已更新）
- [x] ONPREMISE-DEPLOYMENT-GUIDE.md
- [x] PROJECT-RESTRUCTURE-PROGRESS.md
- [x] DEV-BRANCH-ONPREMISE-SUMMARY.md
- [x] PHASE2-COMPLETE.md
- [x] PHASE3-COMPLETE.md
- [x] Application/README.md
- [x] Application/Fe/README.md
- [x] Application/be/README.md
- [x] configs/README.md
- [x] deployments/onpremise/README.md
- [x] build/installer/README.md

---

## 🧪 功能測試

### 本地構建測試

#### Windows

```powershell
# 測試 1: 主構建腳本
cd Application
.\build-local.ps1 -Version "test-1.0.0"
# 預期: dist/ 目錄生成，包含 backend/ 和 frontend/

# 測試 2: 只構建後端
.\build-local.ps1 -SkipFrontend
# 預期: 只有 dist/backend/ 生成

# 測試 3: 只構建前端
.\build-local.ps1 -SkipBackend
# 預期: 只有 dist/frontend/ 生成

# 測試 4: 清理後重建
.\build-local.ps1 -Clean
# 預期: 先刪除 dist/，再重新構建
```

#### Linux/macOS

```bash
# 測試 1: 主構建腳本
cd Application
./build-local.sh
# 預期: dist/ 目錄生成

# 測試 2: 環境變數控制
SKIP_FRONTEND=true ./build-local.sh
# 預期: 只構建後端

SKIP_BACKEND=true ./build-local.sh
# 預期: 只構建前端
```

### 後端構建測試

```bash
cd Application/be

# 測試 1: Make info
make info
# 預期: 顯示所有構建配置

# 測試 2: Make all
make all
# 預期: 構建 3 個二進位檔案到 bin/

# 測試 3: 分別構建
make agent
make console  
make ui
# 預期: 每個命令生成對應的二進位檔案

# 測試 4: 清理
make clean
# 預期: 刪除 bin/ 目錄
```

### 前端構建測試

```bash
cd Application/Fe

# 測試 1: 安裝依賴
npm install
# 預期: node_modules/ 生成，無錯誤

# 測試 2: 開發模式
npm run dev
# 預期: 啟動在 http://localhost:3001

# 測試 3: 生產構建
npm run build
# 預期: .next/ 目錄生成，無錯誤

# 測試 4: 類型檢查
npm run type-check
# 預期: 無 TypeScript 錯誤

# 測試 5: Linting
npm run lint
# 預期: 無 ESLint 錯誤
```

---

## 🔄 CI/CD 測試

### 觸發條件測試

#### CI Workflow

```bash
# 測試: 推送到 dev 分支
git checkout dev
git add .
git commit -m "test: trigger CI"
git push origin dev

# 預期:
# - ci.yml workflow 被觸發
# - basic-check 執行
# - frontend-check 執行（使用 Application/Fe/）
# - docker-build-test 執行
# - security-scan 執行
```

#### 安裝檔構建 Workflow

```bash
# 測試 1: 推送到 dev 分支
git push origin dev
# 預期: build-onpremise-installers.yml 被觸發

# 測試 2: 創建標籤
git tag -a v3.0.0-test -m "Test release"
git push origin v3.0.0-test
# 預期:
# - build-onpremise-installers.yml 被觸發
# - 生成所有安裝檔
# - 創建 GitHub Release
# - 上傳所有 artifacts
```

#### 手動觸發測試

1. 進入 GitHub Actions 頁面
2. 選擇 "Build On-Premise Installers"
3. 點擊 "Run workflow"
4. 選擇 `dev` 分支
5. 輸入版本號（如：test-1.0.0）
6. 點擊 "Run workflow"

預期:
- Workflow 開始執行
- 所有 jobs 成功完成
- Artifacts 可下載

---

## ✅ 驗收標準

### 結構完整性

- [x] 所有必要目錄存在
- [x] 所有必要檔案存在
- [x] 無孤立或冗餘檔案
- [x] .gitignore 正確配置

### 構建系統

- [ ] Windows 構建腳本可執行
- [ ] Linux 構建腳本可執行
- [ ] Makefile 所有目標可用
- [ ] 可生成所有二進位檔案
- [ ] 版本資訊正確嵌入

### 前端應用

- [ ] npm install 成功
- [ ] npm run dev 成功
- [ ] npm run build 成功
- [ ] 無 TypeScript 錯誤
- [ ] 無 ESLint 錯誤
- [ ] TailwindCSS 正常工作
- [ ] 所有組件可正常使用

### 後端應用

- [ ] Make 編譯成功
- [ ] 所有二進位檔案可執行
- [ ] 配置檔案可正常載入
- [ ] 無 Go lint 錯誤
- [ ] 測試全部通過

### CI/CD

- [ ] ci.yml 可正常執行
- [ ] build-onpremise-installers.yml 可正常執行
- [ ] 所有安裝檔可成功生成
- [ ] Artifacts 可正常下載
- [ ] Release 可自動創建

### 文檔

- [ ] 所有 README 內容正確
- [ ] 所有連結有效
- [ ] 使用說明完整
- [ ] 範例程式碼可執行

---

## 🚨 已知問題

### PowerShell 編碼問題

- **問題**: 中文字符在 PowerShell 腳本中可能顯示異常
- **影響**: 測試腳本執行失敗
- **解決**: 使用 UTF-8 編碼，或使用英文訊息

### Make on Windows

- **問題**: Windows 預設沒有 make 命令
- **影響**: 無法直接使用 Makefile
- **解決**: 
  - 使用 build.ps1 腳本
  - 或安裝 make: `choco install make`

---

## 📝 測試報告模板

```markdown
## 測試報告

**測試日期**: YYYY-MM-DD
**測試者**: XXX
**環境**: Windows/Linux/macOS

### 構建測試結果

- [ ] Windows 構建: PASS/FAIL
- [ ] Linux 構建: PASS/FAIL
- [ ] 前端構建: PASS/FAIL
- [ ] 後端構建: PASS/FAIL

### CI/CD 測試結果

- [ ] CI Workflow: PASS/FAIL
- [ ] 安裝檔構建: PASS/FAIL

### 發現的問題

1. ...
2. ...

### 建議改進

1. ...
2. ...
```

---

**狀態**: ✅ 測試清單已準備  
**下一步**: 執行測試並記錄結果  
**最後更新**: 2025-10-09

