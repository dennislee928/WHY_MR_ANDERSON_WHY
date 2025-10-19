# ============================================================================
# 推送 Pandora 映像到 Docker Hub (PowerShell 版本)
# ============================================================================

param(
    [Parameter(Mandatory=$false)]
    [string]$Username = $env:DOCKERHUB_USERNAME,
    
    [Parameter(Mandatory=$false)]
    [string]$Version = "v3.4.1"
)

Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host "  推送 Pandora 映像到 Docker Hub" -ForegroundColor Cyan
Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host ""

# 檢查 Docker Hub 帳號
if (-not $Username) {
    $Username = Read-Host "請輸入您的 Docker Hub 帳號"
}

if (-not $Username) {
    Write-Host "❌ Docker Hub 帳號未設定" -ForegroundColor Red
    exit 1
}

Write-Host "✅ Docker Hub 帳號: $Username" -ForegroundColor Green
Write-Host "ℹ️  版本標籤: $Version" -ForegroundColor Blue
Write-Host ""

# 要推送的映像列表
$images = @(
    "application-axiom-be",
    "application-axiom-ui",
    "application-pandora-agent",
    "application-cyber-ai-quantum"
)

# 步驟 1: 登入 Docker Hub
Write-Host "[1/3] 登入 Docker Hub..." -ForegroundColor Yellow
Write-Host ""

try {
    docker login
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ 登入成功" -ForegroundColor Green
    } else {
        Write-Host "❌ 登入失敗" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "❌ 登入錯誤: $_" -ForegroundColor Red
    exit 1
}

Write-Host ""

# 步驟 2: 標記映像
Write-Host "[2/3] 標記映像..." -ForegroundColor Yellow
Write-Host ""

foreach ($image in $images) {
    $localImage = "${image}:latest"
    $shortName = $image -replace '^application-', ''
    $remoteImage = "$Username/${shortName}:$Version"
    $remoteLatest = "$Username/${shortName}:latest"
    
    Write-Host "標記: $localImage" -ForegroundColor Blue
    Write-Host "  → $remoteImage" -ForegroundColor White
    Write-Host "  → $remoteLatest" -ForegroundColor White
    
    # 標記版本
    docker tag $localImage $remoteImage
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✅ 版本標籤成功" -ForegroundColor Green
    } else {
        Write-Host "  ❌ 標記失敗" -ForegroundColor Red
        continue
    }
    
    # 標記 latest
    docker tag $localImage $remoteLatest
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✅ latest 標籤成功" -ForegroundColor Green
    } else {
        Write-Host "  ❌ 標記失敗" -ForegroundColor Red
    }
    
    Write-Host ""
}

# 步驟 3: 推送映像
Write-Host "[3/3] 推送映像到 Docker Hub..." -ForegroundColor Yellow
Write-Host ""

$pushedCount = 0
$failedCount = 0

foreach ($image in $images) {
    $shortName = $image -replace '^application-', ''
    $remoteImage = "$Username/${shortName}:$Version"
    $remoteLatest = "$Username/${shortName}:latest"
    
    Write-Host "═══════════════════════════════════════" -ForegroundColor Cyan
    Write-Host "推送: $shortName" -ForegroundColor Cyan
    Write-Host "═══════════════════════════════════════" -ForegroundColor Cyan
    Write-Host ""
    
    # 推送版本標籤
    Write-Host "推送 $Version 標籤..." -ForegroundColor Yellow
    docker push $remoteImage
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ $remoteImage 推送成功" -ForegroundColor Green
    } else {
        Write-Host "❌ $remoteImage 推送失敗" -ForegroundColor Red
        $failedCount++
        continue
    }
    
    # 推送 latest 標籤
    Write-Host ""
    Write-Host "推送 latest 標籤..." -ForegroundColor Yellow
    docker push $remoteLatest
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ $remoteLatest 推送成功" -ForegroundColor Green
        $pushedCount++
    } else {
        Write-Host "❌ $remoteLatest 推送失敗" -ForegroundColor Red
        $failedCount++
    }
    
    Write-Host ""
}

# 總結
Write-Host ""
Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host "  推送完成！" -ForegroundColor Cyan
Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "✅ 成功推送: $pushedCount 個映像" -ForegroundColor Green

if ($failedCount -gt 0) {
    Write-Host "❌ 失敗: $failedCount 個映像" -ForegroundColor Red
}

Write-Host ""
Write-Host "您的映像現已可用於：" -ForegroundColor Yellow
Write-Host ""

foreach ($image in $images) {
    $shortName = $image -replace '^application-', ''
    Write-Host "  docker pull $Username/${shortName}:$Version" -ForegroundColor White
    Write-Host "  docker pull $Username/${shortName}:latest" -ForegroundColor White
    Write-Host ""
}

Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host "完成！映像已推送到 Docker Hub" -ForegroundColor Green
Write-Host "Docker Hub: https://hub.docker.com/u/$Username" -ForegroundColor Blue
Write-Host "============================================================================" -ForegroundColor Cyan

