# Pandora Box Console - 本地構建腳本（Windows）
# 用於在本地環境構建地端部署版本

param(
    [string]$Target = "all",
    [string]$Version = "dev",
    [switch]$SkipFrontend = $false,
    [switch]$SkipBackend = $false,
    [switch]$Clean = $false
)

$ErrorActionPreference = "Stop"

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "  Pandora Box Console 本地構建工具   " -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

$SCRIPT_DIR = $PSScriptRoot
$ROOT_DIR = Split-Path -Parent $SCRIPT_DIR
$BACKEND_DIR = Join-Path $SCRIPT_DIR "be"
$FRONTEND_DIR = Join-Path $SCRIPT_DIR "Fe"
$DIST_DIR = Join-Path $SCRIPT_DIR "dist"
$VERSION = $Version
$BUILD_DATE = Get-Date -Format "yyyy-MM-dd_HH:mm:ss"
$GIT_COMMIT = (git rev-parse --short HEAD 2>$null) ?? "unknown"

Write-Host "📋 構建資訊:" -ForegroundColor Yellow
Write-Host "   版本: $VERSION"
Write-Host "   構建日期: $BUILD_DATE"
Write-Host "   Git Commit: $GIT_COMMIT"
Write-Host "   目標: $Target"
Write-Host ""

# 清理
if ($Clean) {
    Write-Host "🧹 清理舊的構建產物..." -ForegroundColor Yellow
    if (Test-Path $DIST_DIR) {
        Remove-Item -Path $DIST_DIR -Recurse -Force
    }
    Write-Host "✅ 清理完成" -ForegroundColor Green
    Write-Host ""
}

# 創建輸出目錄
New-Item -ItemType Directory -Force -Path $DIST_DIR | Out-Null
New-Item -ItemType Directory -Force -Path "$DIST_DIR\backend" | Out-Null
New-Item -ItemType Directory -Force -Path "$DIST_DIR\frontend" | Out-Null

# 構建後端
if (-not $SkipBackend) {
    Write-Host "🔨 構建後端..." -ForegroundColor Yellow
    
    # 檢查 Go 是否安裝
    try {
        $goVersion = go version
        Write-Host "   使用 Go: $goVersion"
    } catch {
        Write-Host "❌ 錯誤: 未找到 Go。請安裝 Go 1.24 或更高版本。" -ForegroundColor Red
        exit 1
    }
    
    # 設定環境變數
    $env:CGO_ENABLED = "0"
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    
    $LDFLAGS = "-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_DATE -X main.GitCommit=$GIT_COMMIT"
    
    # 構建 Agent
    Write-Host "   正在構建 Agent..."
    Push-Location $ROOT_DIR
    
    # 下載依賴
    Write-Host "   正在下載 Go 依賴..."
    go mod download
    
    # 構建 Agent
    go build -ldflags "$LDFLAGS" -o "$DIST_DIR\backend\pandora-agent.exe" .\cmd\agent\main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Agent 構建失敗" -ForegroundColor Red
        Pop-Location
        exit 1
    }
    
    # 構建 Console
    Write-Host "   正在構建 Console..."
    go build -ldflags "$LDFLAGS" -o "$DIST_DIR\backend\pandora-console.exe" .\cmd\console\main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Console 構建失敗" -ForegroundColor Red
        Pop-Location
        exit 1
    }
    
    # 構建 UI Server
    Write-Host "   正在構建 UI Server..."
    go build -ldflags "$LDFLAGS" -o "$DIST_DIR\backend\axiom-ui.exe" .\cmd\ui\main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ UI Server 構建失敗" -ForegroundColor Red
        Pop-Location
        exit 1
    }
    
    Pop-Location
    
    # 複製配置檔案
    Write-Host "   正在複製配置檔案..."
    Copy-Item -Path "$ROOT_DIR\configs" -Destination "$DIST_DIR\backend\configs" -Recurse -Force
    
    Write-Host "✅ 後端構建完成" -ForegroundColor Green
    Write-Host ""
}

# 構建前端
if (-not $SkipFrontend) {
    Write-Host "🎨 構建前端..." -ForegroundColor Yellow
    
    # 檢查 Node.js 是否安裝
    try {
        $nodeVersion = node --version
        Write-Host "   使用 Node.js: $nodeVersion"
    } catch {
        Write-Host "❌ 錯誤: 未找到 Node.js。請安裝 Node.js 18 或更高版本。" -ForegroundColor Red
        exit 1
    }
    
    Push-Location $FRONTEND_DIR
    
    # 安裝依賴
    if (-not (Test-Path "node_modules")) {
        Write-Host "   正在安裝依賴..."
        npm install
        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ 依賴安裝失敗" -ForegroundColor Red
            Pop-Location
            exit 1
        }
    }
    
    # 構建前端
    Write-Host "   正在構建前端應用程式..."
    $env:NEXT_PUBLIC_APP_VERSION = $VERSION
    $env:NODE_ENV = "production"
    npm run build
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ 前端構建失敗" -ForegroundColor Red
        Pop-Location
        exit 1
    }
    
    # 複製構建產物
    Write-Host "   正在複製構建產物..."
    if (Test-Path ".next\standalone") {
        Copy-Item -Path ".next\standalone\*" -Destination "$DIST_DIR\frontend\" -Recurse -Force
    }
    if (Test-Path ".next\static") {
        Copy-Item -Path ".next\static" -Destination "$DIST_DIR\frontend\.next\" -Recurse -Force -ErrorAction SilentlyContinue
    }
    if (Test-Path "public") {
        Copy-Item -Path "public" -Destination "$DIST_DIR\frontend\" -Recurse -Force -ErrorAction SilentlyContinue
    }
    
    Pop-Location
    
    Write-Host "✅ 前端構建完成" -ForegroundColor Green
    Write-Host ""
}

# 創建啟動腳本
Write-Host "📝 創建啟動腳本..." -ForegroundColor Yellow

$startScript = @"
@echo off
echo =====================================
echo   Pandora Box Console IDS-IPS
echo   版本: $VERSION
echo =====================================
echo.

REM 設定環境變數
set LOG_LEVEL=info
set DEVICE_PORT=COM3
set CONFIG_DIR=%~dp0backend\configs

echo 正在啟動服務...
echo.

REM 啟動後端服務
start "Pandora Agent" /D "%~dp0backend" pandora-agent.exe --config "%CONFIG_DIR%\agent-config.yaml"
timeout /t 2 /nobreak >nul

start "Pandora Console" /D "%~dp0backend" pandora-console.exe --config "%CONFIG_DIR%\console-config.yaml"
timeout /t 2 /nobreak >nul

start "Axiom UI" /D "%~dp0backend" axiom-ui.exe --config "%CONFIG_DIR%\ui-config.yaml"
timeout /t 2 /nobreak >nul

echo.
echo =====================================
echo   所有服務已啟動！
echo =====================================
echo.
echo 訪問 Web 介面: http://localhost:3001
echo 訪問 Grafana: http://localhost:3000
echo 訪問 Prometheus: http://localhost:9090
echo.
echo 按任意鍵關閉此視窗...
pause >nul
"@

$startScript | Out-File -FilePath "$DIST_DIR\start.bat" -Encoding ASCII

$stopScript = @"
@echo off
echo 正在停止 Pandora Box Console 服務...
taskkill /F /IM pandora-agent.exe 2>nul
taskkill /F /IM pandora-console.exe 2>nul
taskkill /F /IM axiom-ui.exe 2>nul
echo 所有服務已停止。
pause
"@

$stopScript | Out-File -FilePath "$DIST_DIR\stop.bat" -Encoding ASCII

Write-Host "✅ 啟動腳本已創建" -ForegroundColor Green
Write-Host ""

# 創建 README
$readmeContent = @"
Pandora Box Console IDS-IPS v$VERSION
=====================================

構建資訊
--------
版本: $VERSION
構建日期: $BUILD_DATE
Git Commit: $GIT_COMMIT

快速開始
--------

1. 確保已安裝必要的依賴：
   - PostgreSQL 14+
   - Redis 7+

2. 編輯配置檔案（位於 backend\configs\）

3. 執行 start.bat 啟動所有服務

4. 訪問 http://localhost:3001 使用 Web 介面

停止服務
--------
執行 stop.bat 停止所有服務

服務端口
--------
- Axiom UI: 3001
- Grafana: 3000
- Prometheus: 9090
- Agent API: 8080

技術支援
--------
問題回報: https://github.com/your-org/pandora_box_console_IDS-IPS/issues
電子郵件: support@pandora-ids.com

授權條款
--------
MIT License - 詳見 LICENSE 檔案
"@

$readmeContent | Out-File -FilePath "$DIST_DIR\README.txt" -Encoding UTF8

Write-Host ""
Write-Host "=====================================" -ForegroundColor Green
Write-Host "  ✅ 構建完成！" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green
Write-Host ""
Write-Host "構建產物位於: $DIST_DIR" -ForegroundColor Cyan
Write-Host ""
Write-Host "目錄結構:" -ForegroundColor Yellow
Write-Host "  backend\          - 後端程式"
Write-Host "  frontend\         - 前端程式"
Write-Host "  start.bat         - 啟動所有服務"
Write-Host "  stop.bat          - 停止所有服務"
Write-Host "  README.txt        - 說明文件"
Write-Host ""
Write-Host "下一步:" -ForegroundColor Yellow
Write-Host "  1. cd $DIST_DIR"
Write-Host "  2. 編輯 backend\configs\ 中的配置檔案"
Write-Host "  3. 執行 start.bat 啟動服務"
Write-Host ""

