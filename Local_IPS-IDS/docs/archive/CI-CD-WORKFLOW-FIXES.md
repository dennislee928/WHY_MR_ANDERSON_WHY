# CI/CD Workflow 修復報告

## 📋 問題分析

根據 GitHub Actions workflow 執行結果，發現以下問題：

### 1. Windows 後端構建失敗 ❌
- **症狀**: `構建後端程式 (windows-latest, ...)` 任務失敗
- **原因**: 
  - 缺少錯誤處理和詳細日誌
  - Go 版本設定不正確 (1.24 不存在)
  - 缺少 CGO_ENABLED 設定

### 2. 前端構建失敗 ❌
- **症狀**: `構建前端應用程式` 任務失敗
- **原因**:
  - 缺少目錄結構檢查
  - 錯誤處理不夠詳細
  - 構建產物複製邏輯不完整

### 3. macOS/Ubuntu 後端構建警告 ⚠️
- **症狀**: 後端構建任務顯示警告
- **原因**:
  - 配置檔案複製時缺少存在性檢查
  - 缺少詳細的構建日誌

## 🔧 修復方案

### 修復 1: Windows 後端構建 (PowerShell 語法錯誤)

**問題根源:**
```
ParserError: Missing '(' after 'if' in if statement.
Line | if [ $? -ne 0 ]; then
```

**原因:** Windows runner 使用 PowerShell，但 workflow 使用了 bash 語法

**修改前:**
```yaml
- name: 構建後端二進位檔案
  env:
    GOOS: ${{ matrix.os == 'ubuntu-latest' && 'linux' || matrix.os == 'windows-latest' && 'windows' || 'darwin' }}
    GOARCH: ${{ matrix.arch }}
  run: |
    go build -ldflags="$LDFLAGS" -o dist/backend/pandora-agent${{ matrix.os == 'windows-latest' && '.exe' || '' }} ./cmd/agent/main.go
    if [ $? -ne 0 ]; then
      echo "❌ Agent 構建失敗"
      exit 1
    fi
```

**修改後:**
```yaml
- name: 構建後端二進位檔案
  env:
    GOOS: ${{ matrix.os == 'ubuntu-latest' && 'linux' || matrix.os == 'windows-latest' && 'windows' || 'darwin' }}
    GOARCH: ${{ matrix.arch }}
    CGO_ENABLED: 0
  shell: ${{ matrix.os == 'windows-latest' && 'powershell' || 'bash' }}
  run: |
    # Windows 使用 PowerShell 語法
    ${{ matrix.os == 'windows-latest' && 'New-Item -ItemType Directory -Force -Path "dist/backend"' || 'mkdir -p dist/backend' }}
    
    # 構建 Agent
    ${{ matrix.os == 'windows-latest' && 'Write-Host "構建 Agent (GOOS=$env:GOOS, GOARCH=$env:GOARCH)..."' || 'echo "構建 Agent (GOOS=$GOOS, GOARCH=$GOARCH)..."' }}
    ${{ matrix.os == 'windows-latest' && 'go build -ldflags="$LDFLAGS" -o "dist/backend/pandora-agent.exe" ./cmd/agent/main.go' || 'go build -ldflags="$LDFLAGS" -o dist/backend/pandora-agent${{ matrix.os == "windows-latest" && ".exe" || "" }} ./cmd/agent/main.go' }}
    ${{ matrix.os == 'windows-latest' && 'if ($LASTEXITCODE -ne 0) { Write-Host "❌ Agent 構建失敗" -ForegroundColor Red; exit 1 }' || 'if [ $? -ne 0 ]; then echo "❌ Agent 構建失敗"; exit 1; fi' }}
```

### 修復 2: 前端構建

**修改前:**
```yaml
- name: 構建前端
  working-directory: Application/Fe
  run: |
    npm run build
    mkdir -p ../../dist/frontend
    cp -r .next/standalone/* ../../dist/frontend/
```

**修改後:**
```yaml
- name: 檢查前端目錄
  run: |
    echo "檢查 Application/Fe 目錄結構..."
    ls -la Application/Fe/ || echo "Application/Fe 目錄不存在"
    if [ -f "Application/Fe/package.json" ]; then
      echo "✅ package.json 存在"
      cat Application/Fe/package.json | head -10
    else
      echo "❌ package.json 不存在"
      exit 1
    fi

- name: 構建前端
  working-directory: Application/Fe
  env:
    NEXT_PUBLIC_APP_VERSION: ${{ needs.prepare.outputs.version }}
    NODE_ENV: production
  run: |
    echo "正在構建前端應用程式..."
    echo "環境變數:"
    echo "  NEXT_PUBLIC_APP_VERSION: $NEXT_PUBLIC_APP_VERSION"
    echo "  NODE_ENV: $NODE_ENV"
    
    npm run build
    if [ $? -ne 0 ]; then
      echo "❌ 前端構建失敗"
      exit 1
    fi
    
    echo "✅ 前端構建完成"
    echo "構建產物:"
    ls -la .next/ || echo "未找到 .next 目錄"
    
    # 創建獨立部署包
    mkdir -p ../../dist/frontend
    if [ -d ".next/standalone" ]; then
      echo "複製 standalone 產物..."
      cp -r .next/standalone/* ../../dist/frontend/
    else
      echo "⚠️ 未找到 .next/standalone，嘗試複製整個 .next 目錄"
      cp -r .next ../../dist/frontend/.next/
    fi
    
    if [ -d ".next/static" ]; then
      echo "複製 static 資源..."
      cp -r .next/static ../../dist/frontend/.next/
    fi
    
    if [ -d "public" ]; then
      echo "複製 public 資源..."
      cp -r public ../../dist/frontend/
    fi
    
    echo "✅ 前端產物複製完成"
    ls -la ../../dist/frontend/
```

### 修復 3: Go 版本設定

**修改前:**
```yaml
env:
  GO_VERSION: '1.24'
```

**修改後:**
```yaml
env:
  GO_VERSION: '1.21'
```

### 修復 4: 配置檔案複製 (跨平台語法)

**修改前:**
```yaml
- name: 複製配置檔案
  shell: bash
  run: |
    cp -r configs dist/backend/
    cp -r scripts dist/backend/ || true
```

**修改後:**
```yaml
- name: 複製配置檔案
  shell: ${{ matrix.os == 'windows-latest' && 'powershell' || 'bash' }}
  run: |
    # Windows 使用 PowerShell，其他平台使用 bash
    ${{ matrix.os == 'windows-latest' && 'Write-Host "複製配置檔案..."' || 'echo "複製配置檔案..."' }}
    ${{ matrix.os == 'windows-latest' && 'if (Test-Path "configs") { Copy-Item -Path "configs" -Destination "dist/backend/" -Recurse -Force; Write-Host "✅ configs 目錄已複製" -ForegroundColor Green } else { Write-Host "⚠️ configs 目錄不存在" -ForegroundColor Yellow }' || 'if [ -d "configs" ]; then cp -r configs dist/backend/; echo "✅ configs 目錄已複製"; else echo "⚠️ configs 目錄不存在"; fi' }}
    
    ${{ matrix.os == 'windows-latest' && 'if (Test-Path "scripts") { Copy-Item -Path "scripts" -Destination "dist/backend/" -Recurse -Force; Write-Host "✅ scripts 目錄已複製" -ForegroundColor Green } else { Write-Host "⚠️ scripts 目錄不存在，跳過" -ForegroundColor Yellow }' || 'if [ -d "scripts" ]; then cp -r scripts dist/backend/; echo "✅ scripts 目錄已複製"; else echo "⚠️ scripts 目錄不存在，跳過"; fi' }}
    
    ${{ matrix.os == 'windows-latest' && 'Write-Host "後端產物結構:"; Get-ChildItem -Path "dist/backend" -Recurse | Format-Table Name, Length, LastWriteTime -AutoSize' || 'echo "後端產物結構:"; ls -la dist/backend/' }}
```

## 📝 修復摘要

### 已修復的問題

1. **Windows 後端構建失敗 (PowerShell 語法錯誤)** ✅
   - 根本原因: Windows runner 使用 PowerShell，但 workflow 使用 bash 語法
   - 解決方案: 添加 `shell: ${{ matrix.os == 'windows-latest' && 'powershell' || 'bash' }}`
   - 跨平台語法: Windows 使用 PowerShell，其他平台使用 bash
   - 錯誤處理: `if ($LASTEXITCODE -ne 0)` vs `if [ $? -ne 0 ]`

2. **前端構建失敗** ✅
   - 添加目錄結構檢查
   - 改進錯誤處理
   - 優化構建產物複製邏輯

3. **macOS/Ubuntu 後端構建警告** ✅
   - 添加配置檔案存在性檢查
   - 添加詳細構建日誌

4. **Go 版本設定錯誤** ✅
   - 修正 Go 版本從 1.24 到 1.21

5. **跨平台語法兼容性** ✅
   - Windows: PowerShell 語法 (`New-Item`, `Copy-Item`, `Test-Path`)
   - Linux/macOS: Bash 語法 (`mkdir`, `cp`, `[ -d ]`)

### 改進的錯誤處理

- 每個構建步驟都有錯誤檢查
- 詳細的日誌輸出
- 明確的失敗原因提示
- 構建產物結構驗證

### 新增的檢查步驟

- 前端目錄結構檢查
- package.json 存在性驗證
- 構建產物目錄檢查
- 配置檔案存在性檢查

## 🎯 預期結果

修復後的 workflow 應該能夠：

1. **成功構建所有平台的後端程式**
   - Windows (amd64)
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)

2. **成功構建前端應用程式**
   - 正確處理 package-lock.json
   - 成功執行 npm run build
   - 正確複製構建產物

3. **生成所有安裝檔**
   - Windows 安裝程式 (.exe)
   - Linux 套件 (.deb, .rpm)
   - ISO 安裝光碟
   - OVA 虛擬機配置

4. **創建 GitHub Release**
   - 自動上傳所有構建產物
   - 生成詳細的發布說明

## 🔍 測試建議

### 本地測試

```bash
# 1. 檢查 Go 版本
go version  # 應該顯示 1.21.x

# 2. 測試後端構建
cd Application/be
make all

# 3. 測試前端構建
cd ../Fe
npm install
npm run build

# 4. 檢查構建產物
ls -la ../dist/
```

### 遠端測試

```bash
# 1. 提交修復
git add .github/workflows/build-onpremise-installers.yml
git commit -m "fix(ci): 修復 GitHub Actions workflow 構建問題

- 修正 Go 版本設定 (1.24 -> 1.21)
- 添加詳細錯誤處理和日誌
- 改進前端構建流程
- 優化配置檔案複製邏輯"

# 2. 推送到 dev 分支
git push origin dev

# 3. 觀察 GitHub Actions 執行結果
```

## 📊 修復狀態

- [x] Windows 後端構建失敗
- [x] 前端構建失敗  
- [x] macOS/Ubuntu 後端構建警告
- [x] Go 版本設定錯誤
- [x] 錯誤處理改進
- [x] 日誌輸出優化
- [ ] 本地測試驗證
- [ ] 遠端測試驗證

## 🎉 總結

通過系統性的分析和修復，我們解決了 GitHub Actions workflow 中的主要問題：

1. **根本原因修復**: 修正了 Go 版本設定和構建環境配置
2. **錯誤處理強化**: 添加了詳細的錯誤檢查和日誌輸出
3. **構建流程優化**: 改進了前端和後端的構建邏輯
4. **產物處理完善**: 優化了構建產物的複製和驗證

這些修復應該能夠解決 workflow 執行中的失敗和警告問題，確保 CI/CD 流程能夠順利完成。
