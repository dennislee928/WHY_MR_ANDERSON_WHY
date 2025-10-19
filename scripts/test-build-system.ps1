# 測試構建系統
# 驗證所有構建腳本是否正常工作

param(
    [switch]$Quick = $false
)

$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  構建系統測試" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$ROOT_DIR = Split-Path -Parent $PSScriptRoot
$APP_DIR = Join-Path $ROOT_DIR "Application"
$testsPassed = 0
$testsFailed = 0

function Test-Step {
    param(
        [string]$Name,
        [scriptblock]$Test
    )
    
    Write-Host "測試: $Name..." -ForegroundColor Yellow
    try {
        & $Test
        Write-Host "  ✓ 通過" -ForegroundColor Green
        $script:testsPassed++
        return $true
    } catch {
        Write-Host "  ✗ 失敗: $_" -ForegroundColor Red
        $script:testsFailed++
        return $false
    }
}

# 測試 1: 檢查目錄結構
Test-Step "檢查 Application 目錄結構" {
    if (-not (Test-Path "$APP_DIR\Fe")) { throw "Application/Fe 不存在" }
    if (-not (Test-Path "$APP_DIR\be")) { throw "Application/be 不存在" }
}

# 測試 2: 檢查前端檔案
Test-Step "檢查前端關鍵檔案" {
    $files = @(
        "$APP_DIR\Fe\package.json",
        "$APP_DIR\Fe\next.config.js",
        "$APP_DIR\Fe\tsconfig.json",
        "$APP_DIR\Fe\tailwind.config.js"
    )
    foreach ($file in $files) {
        if (-not (Test-Path $file)) {
            throw "缺少檔案: $file"
        }
    }
}

# 測試 3: 檢查後端檔案
Test-Step "檢查後端關鍵檔案" {
    $files = @(
        "$APP_DIR\be\Makefile",
        "$APP_DIR\be\build.ps1",
        "$APP_DIR\be\build.sh"
    )
    foreach ($file in $files) {
        if (-not (Test-Path $file)) {
            throw "缺少檔案: $file"
        }
    }
}

# 測試 4: 檢查構建腳本
Test-Step "檢查主構建腳本" {
    $files = @(
        "$APP_DIR\build-local.ps1",
        "$APP_DIR\build-local.sh"
    )
    foreach ($file in $files) {
        if (-not (Test-Path $file)) {
            throw "缺少檔案: $file"
        }
    }
}

# 測試 5: 檢查 CI/CD workflows
Test-Step "檢查 CI/CD workflows" {
    $files = @(
        ".github\workflows\ci.yml",
        ".github\workflows\build-onpremise-installers.yml"
    )
    foreach ($file in $files) {
        if (-not (Test-Path $file)) {
            throw "缺少檔案: $file"
        }
    }
}

# 測試 6: 檢查根目錄檔案
Test-Step "檢查專案根目錄" {
    if (-not (Test-Path "go.mod")) { throw "go.mod 不存在" }
    if (-not (Test-Path "README.md")) { throw "README.md 不存在" }
    if (-not (Test-Path ".gitignore")) { throw ".gitignore 不存在" }
}

# 測試 7: 驗證 Go 環境
Test-Step "驗證 Go 環境" {
    $null = go version 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw "Go 未安裝或不在 PATH 中"
    }
}

# 測試 8: 驗證 Node.js 環境（如果不是快速模式）
if (-not $Quick) {
    Test-Step "驗證 Node.js 環境" {
        $null = node --version 2>&1
        if ($LASTEXITCODE -ne 0) {
            throw "Node.js 未安裝或不在 PATH 中"
        }
    }
}

# 測試 9: 檢查安裝檔資源
Test-Step "檢查安裝檔資源" {
    $dirs = @(
        "build\installer\windows",
        "build\installer\linux",
        "build\installer\iso"
    )
    foreach ($dir in $dirs) {
        if (-not (Test-Path $dir)) {
            throw "缺少目錄: $dir"
        }
    }
}

# 測試 10: 檢查文檔
Test-Step "檢查文檔完整性" {
    $files = @(
        "README.md",
        "README-PROJECT-STRUCTURE.md",
        "Application\README.md",
        "Application\Fe\README.md",
        "Application\be\README.md"
    )
    foreach ($file in $files) {
        if (-not (Test-Path $file)) {
            throw "缺少檔案: $file"
        }
    }
}

# 顯示結果
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  測試結果" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "通過: $testsPassed" -ForegroundColor Green
Write-Host "失敗: $testsFailed" -ForegroundColor Red
Write-Host ""

if ($testsFailed -eq 0) {
    Write-Host "✓ 所有測試通過！" -ForegroundColor Green
    Write-Host ""
    Write-Host "專案結構驗證成功，可以進行構建。" -ForegroundColor Cyan
    exit 0
} else {
    Write-Host "✗ 有 $testsFailed 個測試失敗" -ForegroundColor Red
    Write-Host ""
    Write-Host "請修正上述問題後重新測試。" -ForegroundColor Yellow
    exit 1
}

