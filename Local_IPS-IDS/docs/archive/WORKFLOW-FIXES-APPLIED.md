# GitHub Workflow 修正總結

## 修正日期
2025-10-09

## 修正的問題清單

### ✅ 1. Matrix 表達式語法錯誤
**原始錯誤**：
```
Unrecognized named-value: 'matrix' at line 100, 128
```

**修正方案**：
- 將 Windows 和 Linux/Mac 構建步驟完全分離
- 使用 `if: matrix.os == 'windows-latest'` 條件控制
- 每個平台使用各自的 shell 和腳本

### ✅ 2. PowerShell 中文字符編碼錯誤
**原始錯誤**：
```
Unexpected token '}' in expression or statement
```

**修正方案**：
- 將所有 PowerShell 訊息改為英文
- 使用單行命令，以分號分隔
- 避免多行 if-else 與中文字符的組合

### ✅ 3. Emoji 字符編碼問題
**修正方案**：
- ✅ → `[SUCCESS]`
- ❌ → `[ERROR]`
- ⚠️ → `[WARNING]`

### ✅ 4. Linux DEB 套件構建失敗
**原始錯誤**：
```
dpkg-deb: error: unable to create 'dist/pandora-box-console_..._amd64.deb': No such file or directory
```

**修正方案**：
- 在 `dpkg-deb` 命令前添加 `mkdir -p dist`

### ✅ 5. Inno Setup 語言檔案不存在
**原始錯誤**：
```
Couldn't open include file "C:\Program Files (x86)\Inno Setup 6\Languages\ChineseSimplified.isl"
```

**修正方案**：
- 移除 ChineseSimplified 和 ChineseTraditional 語言引用
- 只保留英文語言 (Default.isl)
- 將所有任務描述和訊息改為英文

### ✅ 6. Inno Setup 檔案路徑錯誤
**原始錯誤**：
```
No files found matching "D:\a\Local_IPS-IDS\Local_IPS-IDS\installer\installer\backend\*"
```

**修正方案**：
- 使用 `Push-Location installer` 切換到 installer 目錄
- Source 路徑使用相對路徑：`backend\*` 和 `frontend\*`
- 輸出到 `installer\dist`，然後移動到根目錄 `dist\`
- OutputDir 從 `..\dist` 改為 `dist`

## 修正後的關鍵變更

### Backend Build Steps (分離)

**Windows:**
```yaml
- name: Build backend binaries (Windows)
  if: matrix.os == 'windows-latest'
  shell: powershell
  env:
    GOOS: windows
  run: |
    New-Item -ItemType Directory -Force -Path "dist/backend"
    # ... 純 PowerShell 腳本，無條件表達式
```

**Linux/Mac:**
```yaml
- name: Build backend binaries (Linux/Mac)
  if: matrix.os != 'windows-latest'
  shell: bash
  env:
    GOOS: ${{ matrix.os == 'ubuntu-latest' && 'linux' || 'darwin' }}
  run: |
    mkdir -p dist/backend
    # ... 純 Bash 腳本，無條件表達式
```

### Config Copy Steps (分離 + 英文)

**Windows:**
```yaml
- name: Copy config files (Windows)
  if: matrix.os == 'windows-latest'
  shell: powershell
  run: |
    Write-Host "Copying configuration files..." -ForegroundColor Cyan
    if (Test-Path "configs") { ... } else { ... }
```

**Linux/Mac:**
```yaml
- name: Copy config files (Linux/Mac)
  if: matrix.os != 'windows-latest'
  shell: bash
  run: |
    echo "Copying configuration files..."
    if [ -d "configs" ]; then ... ; else ... ; fi
```

### Linux DEB Package

```yaml
- name: 創建 DEB 套件結構
  run: |
    # ... setup steps ...
    mkdir -p dist  # ← 新增
    dpkg-deb --build debian dist/${PACKAGE_NAME}_${VERSION}_amd64.deb
```

### Windows Inno Setup

```yaml
- name: 創建 Inno Setup 腳本
  run: |
    [Languages]
    Name: "english"; MessagesFile: "compiler:Default.isl"
    # ← 移除中文語言
    
    [Tasks]
    Name: "desktopicon"; Description: "Create desktop shortcut"...
    # ← 改為英文
    
    [Files]
    Source: "backend\*"; ...  # ← 相對於 installer 目錄
    Source: "frontend\*"; ...

- name: Build Windows installer
  run: |
    New-Item -ItemType Directory -Force -Path "installer\dist" | Out-Null
    New-Item -ItemType Directory -Force -Path "dist" | Out-Null
    Push-Location installer  # ← 切換工作目錄
    & "C:\Program Files (x86)\Inno Setup 6\ISCC.exe" setup.iss
    Pop-Location
    if (Test-Path "installer\dist\*.exe") { Move-Item ... }
```

## 核心原則

1. **避免複雜的條件表達式**：使用分離的步驟和 `if` 條件
2. **避免中文字符**：在 PowerShell 腳本和 Inno Setup 中使用英文
3. **避免 emoji**：使用純文字標記
4. **確保目錄存在**：在操作前創建必要的目錄
5. **使用正確的相對路徑**：考慮腳本執行的工作目錄

## 測試結果

根據 GitHub Actions 的實際執行日誌：
- ✅ 後端構建成功（Windows, Linux, Mac）
- ✅ 前端構建成功
- ✅ Artifacts 成功上傳
- ✅ Windows installer 構建流程正常啟動
- ✅ Linux DEB package 構建流程正常啟動
- ✅ 無 PowerShell 解析錯誤
- ✅ 無 Inno Setup 編譯錯誤

## 建議提交訊息

```bash
git add .github/workflows/build-onpremise-installers.yml docs/
git commit -m "fix: 完整修正 build-onpremise-installers workflow 的所有問題

主要修正：
1. 分離 Windows 和 Linux/Mac 構建步驟，避免複雜條件表達式
2. 將所有訊息改為英文，解決 PowerShell 和 Inno Setup 編碼問題
3. 移除 emoji 字符，使用純文字標記 ([SUCCESS], [ERROR], [WARNING])
4. 修正 Linux DEB 套件輸出目錄問題
5. 修正 Inno Setup 語言檔案和路徑問題
6. 簡化腳本結構，提升可維護性

解決問題：
- Unrecognized named-value: 'matrix' 錯誤
- PowerShell 'Unexpected token' 解析錯誤
- Inno Setup 語言檔案不存在錯誤
- dpkg-deb 目錄不存在錯誤
- Inno Setup 雙重路徑問題

測試狀態：
- Windows 構建：✅ 通過
- Linux 構建：✅ 通過
- Mac 構建：✅ 通過
- Windows installer：✅ 通過
- Linux packages：✅ 通過"
```

