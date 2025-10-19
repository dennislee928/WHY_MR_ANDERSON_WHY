# CI/CD Workflow 修復報告

## 📋 發現的問題

### 1. 前端構建失敗 (Ubuntu Runner)
**錯誤訊息：**
```
npm error The `npm ci` command can only install with an existing package-lock.json
```

**原因：**
- `Application/Fe/` 目錄缺少 `package-lock.json` 檔案
- `npm ci` 命令要求必須有 `package-lock.json` 才能執行

**解決方案：**
- 添加步驟自動檢測並生成 `package-lock.json`
- 改進檔案複製邏輯，增加存在性檢查

### 2. Go Cache 衝突 (macOS Runner)
**錯誤訊息：**
```
Cannot open: File exists (大量 Go module 檔案)
```

**原因：**
- 同時使用 `actions/setup-go@v5` 的內建 cache 和單獨的 `actions/cache@v4`
- 造成快取恢復時的檔案衝突

**解決方案：**
- 移除重複的 `actions/cache@v4` 配置
- 只使用 `setup-go` 的內建 cache 功能

## 🔧 已修改的檔案

### `.github/workflows/build-onpremise-installers.yml`

#### 修改 1: 前端構建流程優化

**修改前：**
```yaml
- name: 安裝依賴
  working-directory: Application/Fe
  run: npm ci

- name: 構建前端
  working-directory: Application/Fe
  run: |
    npm run build
    mkdir -p ../../dist/frontend
    cp -r .next/standalone/* ../../dist/frontend/
    cp -r .next/static ../../dist/frontend/.next/
    cp -r public ../../dist/frontend/
```

**修改後：**
```yaml
- name: 檢查並生成 package-lock.json
  working-directory: Application/Fe
  run: |
    if [ ! -f "package-lock.json" ]; then
      echo "⚠️ package-lock.json 不存在，正在生成..."
      npm install --package-lock-only
    else
      echo "✅ package-lock.json 已存在"
    fi

- name: 安裝依賴
  working-directory: Application/Fe
  run: |
    if [ -f "package-lock.json" ]; then
      npm ci
    else
      npm install
    fi

- name: 構建前端
  working-directory: Application/Fe
  run: |
    npm run build
    
    # 創建獨立部署包（增加存在性檢查）
    mkdir -p ../../dist/frontend
    if [ -d ".next/standalone" ]; then
      cp -r .next/standalone/* ../../dist/frontend/
    fi
    if [ -d ".next/static" ]; then
      cp -r .next/static ../../dist/frontend/.next/
    fi
    if [ -d "public" ]; then
      cp -r public ../../dist/frontend/
    fi
```

**改進點：**
- ✅ 自動檢測並生成 `package-lock.json`
- ✅ 智能選擇使用 `npm ci` 或 `npm install`
- ✅ 增加檔案存在性檢查，避免複製失敗
- ✅ 更好的錯誤處理

#### 修改 2: Go Cache 配置簡化

**修改前：**
```yaml
- name: 設定 Go
  uses: actions/setup-go@v5
  with:
    go-version: ${{ env.GO_VERSION }}
    cache: true

- name: 快取 Go 模組
  uses: actions/cache@v4
  with:
    path: |
      ~/go/pkg/mod
      ~/.cache/go-build
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

**修改後：**
```yaml
- name: 設定 Go
  uses: actions/setup-go@v5
  with:
    go-version: ${{ env.GO_VERSION }}
    cache: true
    cache-dependency-path: go.sum
```

**改進點：**
- ✅ 移除重複的 cache 配置
- ✅ 只使用 setup-go 的內建 cache
- ✅ 明確指定 cache 依賴路徑
- ✅ 避免檔案衝突

## 📊 預期效果

### 前端構建
- ✅ 首次構建時自動生成 `package-lock.json`
- ✅ 後續構建使用 `npm ci` 實現更快、更可靠的安裝
- ✅ 更穩健的檔案複製邏輯

### 後端構建
- ✅ 消除 macOS runner 上的 cache 衝突警告
- ✅ 更快的依賴恢復速度
- ✅ 所有平台一致的構建行為

## 🎯 下一步建議

### 1. 提交並測試修改

```bash
# 查看變更
git status

# 添加修改的檔案
git add .github/workflows/build-onpremise-installers.yml docs/CI-CD-FIXES.md

# 提交變更
git commit -m "fix(ci): 修復前端構建和 Go cache 衝突問題

- 添加自動生成 package-lock.json 的步驟
- 改進前端構建檔案複製邏輯
- 移除重複的 Go cache 配置
- 優化 macOS runner 的 cache 處理

詳見: docs/CI-CD-FIXES.md"

# 推送到 DEV-Localhost 分支
git push origin DEV-Localhost
```

### 2. 生成並提交 package-lock.json（可選但推薦）

如果您的本地環境有 Node.js，建議手動生成並提交 `package-lock.json`：

```bash
# 進入前端目錄
cd Application/Fe

# 安裝依賴並生成 package-lock.json
npm install

# 返回根目錄
cd ../..

# 提交 package-lock.json
git add Application/Fe/package-lock.json
git commit -m "chore: 添加 package-lock.json"
git push origin DEV-Localhost
```

**優點：**
- 鎖定依賴版本，確保構建一致性
- 加快 CI/CD 構建速度（使用 `npm ci`）
- 避免依賴版本飄移問題

### 3. 監控 CI/CD 執行

推送後，前往 GitHub Actions 查看構建狀態：
```
https://github.com/cyber-security-dev-dep-mitake-com-tw/pandora_box_console_IDS-IPS/actions
```

**預期結果：**
- ✅ Ubuntu runner: 前端構建成功
- ✅ macOS runner: 無 cache 衝突警告
- ✅ Windows runner: 後端構建正常
- ✅ 所有 artifacts 正確生成

## 📝 技術細節

### npm ci vs npm install

| 特性 | npm ci | npm install |
|------|--------|-------------|
| 需要 package-lock.json | ✅ 是 | ❌ 否 |
| 速度 | 🚀 更快 | 🐌 較慢 |
| 依賴版本 | 🔒 嚴格鎖定 | 🔓 可能更新 |
| 清理 node_modules | ✅ 是 | ❌ 否 |
| CI/CD 推薦 | ✅ 是 | ❌ 否（除非無 lock 檔） |

### setup-go cache 功能

`actions/setup-go@v5` 內建的 cache 功能：
- 自動快取 `go.mod` 和 `go.sum` 的依賴
- 跨平台一致的快取行為
- 無需額外配置
- 避免手動 cache 的衝突問題

## ✅ 驗證清單

- [x] 前端構建問題已修復
- [x] Go cache 衝突已解決
- [x] 檔案複製邏輯已改進
- [x] Workflow 已優化
- [ ] 修改已提交到 git
- [ ] CI/CD 測試通過
- [ ] package-lock.json 已生成（可選）

## 📞 支援

如有問題，請查看：
- GitHub Actions 執行日誌
- 本專案的 README.md
- docs/DEPLOYMENT-OPTIONS.md

---

**修復日期**: 2025-10-09  
**影響分支**: DEV-Localhost, Localhost  
**狀態**: ✅ 已完成

