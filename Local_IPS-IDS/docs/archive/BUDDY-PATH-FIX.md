# Buddy Works 路徑修復指南
## 解決 "YAML PATH: WRONG LOCATION" 錯誤

> 📅 **修復日期**: 2025-10-09  
> 🐛 **問題**: YAML PATH: WRONG LOCATION  
> ✅ **狀態**: 提供多種解決方案

---

## 🐛 問題描述

### 錯誤訊息

```
YAML PATH: WRONG LOCATION
./buddy/pipelines/01-build-installers
```

### 可能原因

1. **路徑格式問題** - Buddy 對路徑格式有特定要求
2. **檔案副檔名問題** - 可能需要或不需要 `.yml` 副檔名
3. **分支問題** - 檔案可能不在指定分支中
4. **權限問題** - Buddy 可能無法訪問該路徑

---

## ✅ 解決方案

### 方案 1: 嘗試不同的路徑格式

按以下順序嘗試：

1. **完整路徑（推薦）**:
   ```
   .buddy/pipelines/01-build-installers.yml
   ```

2. **無前導點**:
   ```
   buddy/pipelines/01-build-installers.yml
   ```

3. **無副檔名**:
   ```
   .buddy/pipelines/01-build-installers
   ```

4. **無前導點和副檔名**:
   ```
   buddy/pipelines/01-build-installers
   ```

### 方案 2: 檢查分支

確保檔案存在於指定分支中：

1. **切換到 main 分支**:
   - 檢查 `.buddy/pipelines/` 目錄是否存在
   - 確認 `01-build-installers.yml` 檔案存在

2. **切換到 dev 分支**:
   - 重複上述檢查

### 方案 3: 使用根目錄檔案

如果路徑問題持續，可以將管道檔案移到根目錄：

1. **創建單一管道檔案**:
   ```bash
   # 複製到根目錄
   cp .buddy/pipelines/01-build-installers.yml build-installers.yml
   ```

2. **在 Buddy 中使用**:
   ```
   build-installers.yml
   ```

---

## 🔧 詳細步驟

### Step 1: 驗證檔案存在

1. 在 GitHub 倉庫中檢查：
   - 訪問 `https://github.com/your-org/Local_IPS-IDS`
   - 確認分支（main 或 dev）
   - 導航到 `.buddy/pipelines/`
   - 確認 `01-build-installers.yml` 存在

### Step 2: 嘗試導入

1. 在 Buddy Works 中：
   - 點擊 "Pipelines" → "Add new"
   - 選擇 "Import YAML" → "From Git"
   - 設置 PROJECT: `Local_IPS-IDS (This project)`
   - 設置 BRANCH: `main` 或 `dev`

2. 嘗試不同路徑：
   ```
   嘗試 1: .buddy/pipelines/01-build-installers.yml
   嘗試 2: buddy/pipelines/01-build-installers.yml
   嘗試 3: .buddy/pipelines/01-build-installers
   嘗試 4: buddy/pipelines/01-build-installers
   ```

### Step 3: 檢查錯誤訊息

如果仍然失敗，檢查：
- 檔案是否真的存在於該分支
- 檔案內容是否有效 YAML
- 是否有權限問題

---

## 🚀 替代方案

### 方案 A: 使用 Inline YAML

如果路徑問題持續，使用 "Inline YAML" 選項：

1. 在 Buddy 中選擇 "Inline YAML"
2. 複製 `01-build-installers.yml` 的內容
3. 直接貼上到 Buddy 編輯器

### 方案 B: 創建根目錄檔案

將管道檔案移到根目錄：

```bash
# 在倉庫根目錄創建
cp .buddy/pipelines/01-build-installers.yml build-installers.yml
git add build-installers.yml
git commit -m "Add build-installers.yml to root for Buddy import"
git push
```

然後在 Buddy 中使用路徑：
```
build-installers.yml
```

### 方案 C: 使用 GitHub Actions 作為參考

如果 Buddy 導入持續有問題，可以：

1. 先使用現有的 GitHub Actions workflows
2. 稍後再嘗試 Buddy 導入
3. 參考 `.github/workflows/build-onpremise-installers.yml`

---

## 📊 路徑測試表

| 路徑格式 | 狀態 | 備註 |
|----------|------|------|
| `.buddy/pipelines/01-build-installers.yml` | ✅ 推薦 | 完整路徑 |
| `buddy/pipelines/01-build-installers.yml` | 🔄 嘗試 | 無前導點 |
| `.buddy/pipelines/01-build-installers` | 🔄 嘗試 | 無副檔名 |
| `buddy/pipelines/01-build-installers` | 🔄 嘗試 | 無前導點和副檔名 |
| `build-installers.yml` | 🔄 備選 | 根目錄檔案 |

---

## 🎯 故障排除

### 檢查清單

- [ ] 檔案存在於指定分支
- [ ] 路徑格式正確
- [ ] YAML 語法有效
- [ ] Buddy 有倉庫訪問權限
- [ ] 分支名稱正確

### 常見錯誤

1. **檔案不存在**:
   - 檢查分支
   - 確認檔案已推送

2. **權限問題**:
   - 檢查 GitHub 倉庫權限
   - 確認 Buddy 整合設置

3. **YAML 語法錯誤**:
   - 驗證 YAML 格式
   - 檢查縮排

---

## 🎊 成功指標

導入成功後，您應該看到：

- ✅ 管道出現在 Buddy 項目中
- ✅ 管道配置正確顯示
- ✅ 觸發條件設置正確
- ✅ Actions 列表完整

---

## 📚 相關文檔

- [管道導入指南](BUDDY-PIPELINE-IMPORT-GUIDE.md)
- [Buddy Works 設置](BUDDY-WORKS-SETUP.md)
- [YAML 修復說明](BUDDY-YAML-FIX.md)

---

**狀態**: 🔄 提供多種解決方案  
**下一步**: 嘗試不同路徑格式  
**預計時間**: 5-10 分鐘

**🎯 選擇最適合的方案，成功導入管道！**
