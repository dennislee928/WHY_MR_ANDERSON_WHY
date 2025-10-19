# Application/be/build.ps1
# Windows 後端構建腳本

param(
    [string]$Target = "all",
    [string]$Version = "dev"
)

$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  後端構建腳本 (Windows)" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$ROOT_DIR = (Get-Location).Path + "\..\..\"
$ROOT_DIR = (Resolve-Path $ROOT_DIR).Path

Write-Host "專案根目錄: $ROOT_DIR" -ForegroundColor Gray
Write-Host "構建目標: $Target" -ForegroundColor Gray
Write-Host "版本: $Version" -ForegroundColor Gray
Write-Host ""

# 檢查 Go 是否安裝
try {
    $goVersion = go version
    Write-Host "✓ Go 環境: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "✗ 錯誤: 未找到 Go" -ForegroundColor Red
    Write-Host "  請安裝 Go 1.24+ 後再試" -ForegroundColor Yellow
    exit 1
}

# 檢查根目錄是否正確
if (-not (Test-Path "$ROOT_DIR\go.mod")) {
    Write-Host "✗ 錯誤: 找不到 go.mod" -ForegroundColor Red
    Write-Host "  ROOT_DIR: $ROOT_DIR" -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "開始構建..." -ForegroundColor Yellow
Write-Host ""

# 使用 make
$makeCommand = "make"
if ($Target -eq "all") {
    & $makeCommand all
} else {
    & $makeCommand $Target
}

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  ✓ 構建成功！" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "二進位檔案位於: .\bin\" -ForegroundColor Cyan
    Write-Host ""
    if (Test-Path ".\bin\") {
        Get-ChildItem ".\bin\" | Format-Table Name, Length, LastWriteTime
    }
} else {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Red
    Write-Host "  ✗ 構建失敗" -ForegroundColor Red
    Write-Host "========================================" -ForegroundColor Red
    exit 1
}

