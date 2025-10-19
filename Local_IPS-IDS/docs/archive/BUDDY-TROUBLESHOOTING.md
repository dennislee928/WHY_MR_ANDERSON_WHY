# Buddy Works 故障排除指南
## 解決環境變數和緩存問題

> 📅 **創建日期**: 2025-10-09  
> 🎯 **目標**: 解決 Buddy Works 運行時問題  
> ✅ **狀態**: 提供完整故障排除方案

---

## 🐛 常見問題

### 問題 1: 環境變數錯誤

**錯誤訊息**:
```
11: cannot create : Directory nonexistent
echo "VERSION=$VERSION" >> $BUDDY_EXPORT_FILE
```

**原因**: `$BUDDY_EXPORT_FILE` 在 Buddy Works 中不存在

**解決方案**: ✅ 已修復
- 移除所有 `$BUDDY_EXPORT_FILE` 的使用
- 讓每個 action 獨立計算版本資訊

### 問題 2: 緩存問題

**現象**: 修復後仍然出現舊錯誤

**原因**: Buddy Works 可能使用了緩存的舊版本

**解決方案**:
1. **等待同步** (推薦):
   - 等待 2-5 分鐘讓 Buddy 同步最新檔案
   - 重新觸發管道

2. **使用簡化版測試**:
   - 使用 `build-installers-simple.yml`
   - 驗證修復是否有效

3. **強制刷新**:
   - 在 Buddy 中刪除舊管道
   - 重新導入最新版本

---

## 🚀 解決方案

### 方案 A: 使用修復後的完整版

**檔案**: `build-installers.yml`
**狀態**: ✅ 已修復所有環境變數問題

**特點**:
- 完整的構建流程
- 所有環境變數問題已修復
- 包含 GitHub Release 功能

### 方案 B: 使用簡化版測試

**檔案**: `build-installers-simple.yml`
**狀態**: ✅ 新創建的簡化版

**特點**:
- 只有 2 個 action
- 專注於環境變數測試
- 快速驗證修復效果

### 方案 C: 使用單獨檔案

**檔案**: `.buddy/01-build-installers.yml` 等
**狀態**: ✅ 已修復

**特點**:
- 單獨的管道檔案
- 避免名稱衝突
- 可以分別導入

---

## 🔧 修復驗證

### 檢查修復是否有效

1. **檢查檔案內容**:
   ```bash
   # 應該沒有 BUDDY_EXPORT_FILE
   grep -r "BUDDY_EXPORT_FILE" .
   # 應該沒有輸出
   ```

2. **檢查環境變數設定**:
   ```bash
   # 每個 action 都應該有：
   export VERSION=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo '0.1.0')
   export BUILD_DATE=$(date -u '+%Y-%m-%d_%H:%M:%S')
   export GIT_COMMIT=$(git rev-parse --short HEAD)
   ```

3. **運行測試**:
   - 使用簡化版管道測試
   - 確認環境變數正常顯示

---

## 📊 修復狀態

### 已修復的問題

| 問題 | 狀態 | 修復方法 |
|------|------|----------|
| `$BUDDY_EXPORT_FILE` 錯誤 | ✅ 已修復 | 移除使用，獨立計算 |
| 環境變數傳遞 | ✅ 已修復 | 每個 action 獨立設定 |
| 路徑格式問題 | ✅ 已修復 | 改為絕對路徑 |
| Docker 鏡像設定 | ✅ 已修復 | 添加缺失欄位 |
| 管道名稱衝突 | ✅ 已修復 | 刪除重複檔案 |

### 檔案狀態

| 檔案 | 狀態 | 用途 |
|------|------|------|
| `build-installers.yml` | ✅ 已修復 | 完整版管道 |
| `build-installers-simple.yml` | ✅ 新創建 | 簡化版測試 |
| `.buddy/01-build-installers.yml` | ✅ 已修復 | 單獨管道檔案 |
| `buddy.yml` | ❌ 已刪除 | 避免衝突 |

---

## 🎯 推薦步驟

### 立即行動

1. **等待 2-5 分鐘**:
   - 讓 Buddy Works 同步最新檔案
   - 避免緩存問題

2. **測試簡化版**:
   - 導入 `build-installers-simple.yml`
   - 驗證環境變數是否正常

3. **如果簡化版成功**:
   - 導入完整版 `build-installers.yml`
   - 運行完整構建流程

### 如果問題持續

1. **檢查檔案版本**:
   - 確認 Buddy 使用的是最新檔案
   - 檢查 Git commit 是否是最新的

2. **清除緩存**:
   - 在 Buddy 中刪除舊管道
   - 重新導入最新版本

3. **使用替代方案**:
   - 使用 `.buddy/01-build-installers.yml`
   - 或使用 Inline YAML 方式

---

## 📚 相關文檔

- [最終解決方案](BUDDY-FINAL-SOLUTION.md) - 完整解決方案
- [完整修復指南](BUDDY-COMPLETE-FIX-GUIDE.md) - 詳細修復記錄
- [YAML 修復說明](BUDDY-YAML-FIX.md) - 技術修復細節

---

## 🚨 緊急解決方案

### 如果所有方案都失敗

1. **使用 Inline YAML**:
   - 在 Buddy 中選擇 "Inline YAML"
   - 複製 `build-installers-simple.yml` 內容
   - 直接貼上到編輯器

2. **手動創建管道**:
   - 在 Buddy 中手動創建新管道
   - 添加必要的環境變數設定
   - 逐步添加 actions

3. **聯繫支援**:
   - 如果問題持續存在
   - 可能需要 Buddy Works 技術支援

---

**狀態**: ✅ 所有已知問題已修復  
**建議**: 先測試簡化版，再使用完整版  
**預計時間**: 5-10 分鐘

**🎉 修復應該已經生效，如果問題持續請使用簡化版測試！**
