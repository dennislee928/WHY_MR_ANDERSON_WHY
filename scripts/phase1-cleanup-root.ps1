# Phase 1: 清理根目錄
# 移除和整理不需要的檔案

param(
    [switch]$DryRun = $false
)

$ErrorActionPreference = "Continue"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  階段 1: 清理根目錄" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

if ($DryRun) {
    Write-Host "⚠️  DRY RUN MODE - 不會實際刪除檔案" -ForegroundColor Yellow
    Write-Host ""
}

# 1. 移除根目錄的舊 Dockerfile
Write-Host "1. 清理根目錄的 Dockerfile..." -ForegroundColor Yellow
$oldDockerfiles = @(
    "Dockerfile.agent",
    "Dockerfile.agent.koyeb",
    "Dockerfile.monitoring",
    "Dockerfile.nginx",
    "Dockerfile.server-be",
    "Dockerfile.server-fe",
    "Dockerfile.test",
    "Dockerfile.ui.patr"
)

foreach ($file in $oldDockerfiles) {
    if (Test-Path $file) {
        Write-Host "  移除: $file" -ForegroundColor Gray
        if (-not $DryRun) {
            Remove-Item $file -Force
        }
    }
}

# 2. 整理 DOCUMENTS 目錄到 docs/archive
Write-Host ""
Write-Host "2. 整理 DOCUMENTS 目錄..." -ForegroundColor Yellow
if (Test-Path "DOCUMENTS") {
    Write-Host "  移動 DOCUMENTS 到 docs/archive/" -ForegroundColor Gray
    if (-not $DryRun) {
        if (-not (Test-Path "docs/archive")) {
            New-Item -ItemType Directory -Path "docs/archive" -Force | Out-Null
        }
        Move-Item "DOCUMENTS" "docs/archive/DOCUMENTS-old" -Force
    }
}

# 3. 清理舊的 k8s 目錄（保留 deployments/kubernetes）
Write-Host ""
Write-Host "3. 清理舊的 k8s 目錄..." -ForegroundColor Yellow
$oldK8sDirs = @("k8s", "k8s-gcp")
foreach ($dir in $oldK8sDirs) {
    if (Test-Path $dir) {
        Write-Host "  移動: $dir 到 deployments/kubernetes/legacy/" -ForegroundColor Gray
        if (-not $DryRun) {
            if (-not (Test-Path "deployments/kubernetes/legacy")) {
                New-Item -ItemType Directory -Path "deployments/kubernetes/legacy" -Force | Out-Null
            }
            Move-Item $dir "deployments/kubernetes/legacy/$dir" -Force
        }
    }
}

# 4. 清理臨時和報告檔案
Write-Host ""
Write-Host "4. 清理臨時檔案..." -ForegroundColor Yellow
$tempFiles = @(
    "report.md",
    "RESTRUCTURE-SUCCESS.md",
    "PROJECT-RESTRUCTURE-FINAL-REPORT.md"
)

foreach ($file in $tempFiles) {
    if (Test-Path $file) {
        Write-Host "  移動: $file 到 docs/archive/" -ForegroundColor Gray
        if (-not $DryRun) {
            Move-Item $file "docs/archive/" -Force
        }
    }
}

# 5. 清理 bin 目錄（構建產物）
Write-Host ""
Write-Host "5. 清理 bin 目錄..." -ForegroundColor Yellow
if (Test-Path "bin") {
    Write-Host "  bin 目錄應該由 .gitignore 排除" -ForegroundColor Gray
    # 不刪除，但確保在 .gitignore 中
}

# 6. 整理 web 目錄（舊版前端）
Write-Host ""
Write-Host "6. 整理舊版 web 目錄..." -ForegroundColor Yellow
if (Test-Path "web") {
    Write-Host "  移動 web 到 Application/Fe/legacy/" -ForegroundColor Gray
    if (-not $DryRun) {
        if (-not (Test-Path "Application/Fe/legacy")) {
            New-Item -ItemType Directory -Path "Application/Fe/legacy" -Force | Out-Null
        }
        Move-Item "web" "Application/Fe/legacy/web-old" -Force
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  ✅ 階段 1 完成" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""

if ($DryRun) {
    Write-Host "這是 DRY RUN，實際執行請運行：" -ForegroundColor Yellow
    Write-Host "  .\scripts\phase1-cleanup-root.ps1" -ForegroundColor White
}

