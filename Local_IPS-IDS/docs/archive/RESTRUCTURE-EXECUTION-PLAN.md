# 專案重構執行計劃

> **建立**: 2025-10-09  
> **方法**: 規劃 → 驗證 → 執行 → 記錄  
> **狀態**: 📋 執行中

---

## 🔍 Phase 1: 深度分析（已完成）

### 發現的問題

#### ❌ 問題 1: build/docker/ 有重複的 Dockerfiles

**現狀**: 16 個檔案（8個 `*.dockerfile` + 8個 `Dockerfile.*`）  
**應該**: 8 個檔案（只保留 `*.dockerfile` 命名）

**檔案清單**:
```
重複檔案（需刪除）:
- Dockerfile.agent           → 保留 agent.dockerfile
- Dockerfile.agent.koyeb     → 保留 agent.koyeb.dockerfile
- Dockerfile.monitoring      → 保留 monitoring.dockerfile
- Dockerfile.nginx           → 保留 nginx.dockerfile
- Dockerfile.server-be       → 保留 server-be.dockerfile
- Dockerfile.server-fe       → 保留 server-fe.dockerfile
- Dockerfile.test            → 保留 test.dockerfile
- Dockerfile.ui.patr         → 保留 ui.patr.dockerfile
```

**執行計劃**: 刪除所有 `Dockerfile.*` 檔案，保留 `*.dockerfile`

---

#### ❌ 問題 2: .terraform 目錄存在

**現狀**: 
- `terraform/.terraform/` 存在
- `deployments/terraform/.terraform/` 存在

**應該**: 這些應該被 .gitignore 排除，不應該在版控中

**執行計劃**: 
1. 移除這些目錄
2. 確保 .gitignore 包含 `.terraform/`

---

#### ⚠️ 問題 3: 部分文檔檔案被刪除

**現狀**: 用戶刪除了一些我創建的報告檔案  
**應該**: 保留必要的文檔，整理到 docs/ 目錄

**執行計劃**: 在 docs/ 目錄重新創建必要文檔

---

## 📋 執行待辦清單

### Todo Group 1: 清理重複檔案

- [ ] T1.1: 刪除 `build/docker/Dockerfile.agent`
- [ ] T1.2: 刪除 `build/docker/Dockerfile.agent.koyeb`
- [ ] T1.3: 刪除 `build/docker/Dockerfile.monitoring`
- [ ] T1.4: 刪除 `build/docker/Dockerfile.nginx`
- [ ] T1.5: 刪除 `build/docker/Dockerfile.server-be`
- [ ] T1.6: 刪除 `build/docker/Dockerfile.server-fe`
- [ ] T1.7: 刪除 `build/docker/Dockerfile.test`
- [ ] T1.8: 刪除 `build/docker/Dockerfile.ui.patr`

**驗證命令**:
```powershell
Get-ChildItem -Path "build\docker" | Measure-Object
# Count 應該是 8
```

---

### Todo Group 2: 清理 .terraform 目錄

- [ ] T2.1: 刪除 `terraform/.terraform/`
- [ ] T2.2: 刪除 `terraform/.terraform.lock.hcl`
- [ ] T2.3: 刪除 `deployments/terraform/.terraform/`
- [ ] T2.4: 刪除 `deployments/terraform/.terraform.lock.hcl`
- [ ] T2.5: 更新 .gitignore 確保包含 `.terraform/`

**驗證命令**:
```bash
find . -name ".terraform" -type d
# 應該返回空
```

---

### Todo Group 3: 整理文檔結構

- [ ] T3.1: 在 docs/ 創建所有必要的文檔
- [ ] T3.2: 移除根目錄的臨時報告
- [ ] T3.3: 創建 docs/onpremise/ 子目錄
- [ ] T3.4: 組織所有部署相關文檔

**目標結構**:
```
docs/
├── onpremise/              # 地端部署文檔
│   ├── QUICK-START.md
│   ├── DEPLOYMENT-GUIDE.md
│   └── TESTING-CHECKLIST.md
├── development/            # 開發文檔
│   ├── FRONTEND-GUIDE.md
│   └── BACKEND-GUIDE.md
├── archive/                # 存檔
└── RESTRUCTURE-FINAL-REPORT.md
```

---

### Todo Group 4: 驗證 Application/ 結構

- [ ] T4.1: 驗證 Application/Fe/ 所有檔案存在
- [ ] T4.2: 驗證 Application/be/ 所有檔案存在
- [ ] T4.3: 測試前端構建（npm install, npm run build）
- [ ] T4.4: 測試後端構建（make all或build腳本）
- [ ] T4.5: 驗證 build-local.* 腳本可執行

---

### Todo Group 5: 更新和測試 CI/CD

- [ ] T5.1: 驗證 ci.yml 語法
- [ ] T5.2: 驗證 build-onpremise-installers.yml 語法
- [ ] T5.3: 測試所有路徑引用
- [ ] T5.4: 創建 workflow 測試文檔

---

### Todo Group 6: 最終文檔

- [ ] T6.1: 更新 README.md
- [ ] T6.2: 更新 README-PROJECT-STRUCTURE.md
- [ ] T6.3: 創建 COMMIT-MESSAGE.md
- [ ] T6.4: 創建 CHANGELOG.md
- [ ] T6.5: 創建最終驗證清單

---

## ⚡ 立即執行計劃

### Step 1: 清理 build/docker/ 重複檔案

```powershell
cd build\docker
Remove-Item Dockerfile.* -Force
Get-ChildItem | Measure-Object  # 應該是 8
```

### Step 2: 清理 .terraform 目錄

```powershell
Remove-Item -Path "terraform\.terraform" -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item -Path "terraform\.terraform.lock.hcl" -Force -ErrorAction SilentlyContinue
Remove-Item -Path "deployments\terraform\.terraform" -Recurse -Force -ErrorAction SilentlyContinue
```

### Step 3: 整理文檔

```powershell
New-Item -ItemType Directory -Path "docs\onpremise", "docs\development" -Force
# 移動和創建文檔
```

### Step 4: 測試構建

```powershell
# 測試前端
cd Application\Fe
npm install
npm run build

# 測試後端
cd ..\be
# 如果有 make: make all
# 否則: ..\build.ps1
```

---

## 📊 執行進度追蹤

| Group | 任務數 | 完成 | 進度 |
|-------|--------|------|------|
| Group 1 | 8 | 0 | 0% |
| Group 2 | 5 | 0 | 0% |
| Group 3 | 4 | 0 | 0% |
| Group 4 | 5 | 0 | 0% |
| Group 5 | 4 | 0 | 0% |
| Group 6 | 5 | 0 | 0% |
| **總計** | **31** | **0** | **0%** |

---

## 🎯 執行順序

1. **Group 1**: 清理 Dockerfiles（5分鐘）
2. **Group 2**: 清理 .terraform（2分鐘）
3. **Group 3**: 整理文檔（10分鐘）
4. **Group 4**: 驗證 Application（20分鐘）
5. **Group 5**: 驗證 CI/CD（10分鐘）
6. **Group 6**: 最終文檔（15分鐘）

**總預計時間**: ~60分鐘

---

**下一步**: 等待確認後開始執行 Group 1

