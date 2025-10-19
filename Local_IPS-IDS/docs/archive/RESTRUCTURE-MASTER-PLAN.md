# 專案重構主計劃

> **日期**: 2025-10-09  
> **狀態**: 🔄 規劃中  
> **方法**: 規劃 → 驗證 → 執行 → 記錄

---

## 📊 當前狀態分析

### ❌ 發現的問題

根據 `tree` 輸出，發現以下問題：

#### 1. 目錄位置問題
- ❌ `terraform/.terraform/` 存在（應該被 .gitignore 排除）
- ❌ `Application/Fe/legacy/web-old/` 存在但未整合
- ❌ `deployments/terraform/.terraform/` 重複
- ⚠️ `build/docker/` 可能缺少某些 Dockerfiles

#### 2. 缺少的檔案
- ❌ Application/Fe/ 缺少實際的 pages 內容
- ❌ Application/be/ 可能缺少某些連結或配置
- ❌ 某些 Dockerfiles 可能遺失

#### 3. 結構不一致
- ❌ 部分舊檔案未清理
- ❌ 部分新結構未完全建立

---

## 📋 重構主計劃（10個階段）

### Phase 1: 深度分析和規劃 ✅
**目標**: 徹底分析現有結構，列出所有需要移動/創建的檔案

**任務**:
- [x] 分析 tree 輸出
- [x] 列出所有問題
- [ ] 建立完整的檔案清單
- [ ] 建立移動計劃
- [ ] 建立創建計劃
- [ ] 建立刪除計劃

**驗證**: 檔案清單完整，計劃可行

---

### Phase 2: 清理 Terraform 和臨時檔案
**目標**: 清理所有不應存在的目錄

**任務**:
- [ ] 移除 `terraform/.terraform/`
- [ ] 移除 `deployments/terraform/.terraform/`
- [ ] 更新 .gitignore 確保這些不再出現
- [ ] 清理所有 .terraform.lock.hcl

**驗證**: 
```bash
find . -name ".terraform" -type d
# 應該返回空
```

---

### Phase 3: 整理 build/ 目錄
**目標**: 確保所有 Dockerfiles 和構建資源在正確位置

**任務**:
- [ ] 檢查 build/docker/ 中的所有 Dockerfiles
- [ ] 從根目錄移動遺漏的 Dockerfiles
- [ ] 驗證每個 Dockerfile 可用
- [ ] 創建 build/docker/README.md

**驗證**:
```bash
ls -la build/docker/*.dockerfile
# 應該看到所有必要的 Dockerfiles
```

---

### Phase 4: 完整建立 Application/Fe/
**目標**: 建立完整可運行的前端應用

**任務**:
- [ ] 整合 legacy/web-old/ 的內容
- [ ] 確保所有 components 完整
- [ ] 確保所有 pages 完整
- [ ] 確保所有配置檔案完整
- [ ] 測試 npm install
- [ ] 測試 npm run build

**驗證**:
```bash
cd Application/Fe
npm install
npm run build
npm run dev
# 應該能正常啟動
```

---

### Phase 5: 完整建立 Application/be/
**目標**: 建立完整可編譯的後端結構

**任務**:
- [ ] 驗證 Makefile 所有路徑
- [ ] 測試編譯每個程式
- [ ] 確保 configs 引用正確
- [ ] 測試所有 make 目標

**驗證**:
```bash
cd Application/be
make info
make all
# 應該成功編譯
```

---

### Phase 6: 完善安裝檔生成資源
**目標**: 確保所有安裝檔可以正確生成

**任務**:
- [ ] 驗證 build/installer/ 所有檔案
- [ ] 測試 Inno Setup 腳本
- [ ] 測試 Debian 打包腳本
- [ ] 測試 ISO 生成腳本
- [ ] 創建測試說明

**驗證**: 每種安裝檔都有完整的生成資源

---

### Phase 7: 修正所有 CI/CD Workflows
**目標**: 確保所有 workflows 可以正常執行

**任務**:
- [ ] 驗證 ci.yml 路徑
- [ ] 驗證 build-onpremise-installers.yml 所有步驟
- [ ] 修正所有路徑引用
- [ ] 測試 workflow 語法
- [ ] 創建 workflow 測試計劃

**驗證**:
```bash
# 使用 actionlint 驗證
actionlint .github/workflows/*.yml
```

---

### Phase 8: 完善主構建腳本
**目標**: build-local.* 腳本完全可用

**任務**:
- [ ] 測試 build-local.ps1
- [ ] 測試 build-local.sh
- [ ] 修正所有路徑問題
- [ ] 添加更多錯誤處理
- [ ] 測試各種參數組合

**驗證**: 在 Windows 和 Linux 上實際運行成功

---

### Phase 9: 整理和清理
**目標**: 移除所有不需要的檔案，整理目錄結構

**任務**:
- [ ] 清理所有臨時檔案
- [ ] 整理 docs/ 目錄
- [ ] 清理舊的測試檔案
- [ ] 更新 .gitignore
- [ ] 驗證 git status

**驗證**:
```bash
git status
# 應該只顯示真正要提交的變更
```

---

### Phase 10: 最終驗證和文檔
**目標**: 完整測試，完善文檔，準備發布

**任務**:
- [ ] 執行完整構建測試
- [ ] 執行 CI/CD 測試
- [ ] 更新所有 README
- [ ] 創建 CHANGELOG
- [ ] 創建 commit message
- [ ] 準備 Release Notes

**驗證**: 所有測試通過，文檔完整

---

## 📝 檔案清單（待建立）

### 需要檢查的目錄

1. **build/docker/**
   - [ ] agent.dockerfile
   - [ ] agent.koyeb.dockerfile
   - [ ] server-be.dockerfile
   - [ ] server-fe.dockerfile
   - [ ] monitoring.dockerfile
   - [ ] nginx.dockerfile
   - [ ] ui.patr.dockerfile
   - [ ] test.dockerfile

2. **Application/Fe/**
   - [ ] 所有 components
   - [ ] 所有 pages
   - [ ] 所有 services
   - [ ] 所有配置

3. **Application/be/**
   - [ ] Makefile
   - [ ] 構建腳本
   - [ ] 配置引用

---

## 🎯 下一步

**立即執行**: Phase 1 深度分析

我將：
1. 列出所有現有檔案
2. 列出所有缺少的檔案
3. 建立完整的移動/創建/刪除計劃
4. 等待您確認後執行

---

**狀態**: 📋 規劃階段  
**下一步**: 深度分析  
**預計完成**: 需要系統性執行所有階段

