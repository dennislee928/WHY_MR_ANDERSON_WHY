# ============================================================================
# Pandora Cyber AI/Quantum - 重新建構與部署腳本
# ============================================================================
# 用途: 重新建構 cyber-ai-quantum 服務並重新啟動
# 執行方式: .\rebuild-quantum.ps1
# ============================================================================

param(
    [Parameter(Mandatory=$false)]
    [switch]$NoBuild = $false,
    
    [Parameter(Mandatory=$false)]
    [switch]$Clean = $false,
    
    [Parameter(Mandatory=$false)]
    [string]$IBMToken = $env:IBM_QUANTUM_TOKEN
)

Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host "  Pandora Cyber AI/Quantum - 重新建構與部署" -ForegroundColor Cyan
Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host ""

# 切換到正確的目錄
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
Push-Location $scriptPath

# 檢查 Docker 是否運行
Write-Host "[1/6] 檢查 Docker 環境..." -ForegroundColor Yellow
try {
    $dockerVersion = docker version --format '{{.Server.Version}}' 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "[ERROR] Docker 未運行或未安裝" -ForegroundColor Red
        Write-Host "[INFO] 請啟動 Docker Desktop" -ForegroundColor Yellow
        Pop-Location
        exit 1
    }
    Write-Host "[OK] Docker 版本: $dockerVersion" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] 無法檢測 Docker: $_" -ForegroundColor Red
    Pop-Location
    exit 1
}

# 設定環境變數
if ($IBMToken) {
    Write-Host "`n[2/6] 設定 IBM Quantum Token..." -ForegroundColor Yellow
    $env:IBM_QUANTUM_TOKEN = $IBMToken
    Write-Host "[OK] Token 已設定 (${IBMToken.Substring(0, 8)}...)" -ForegroundColor Green
} else {
    Write-Host "`n[2/6] 警告: 未設定 IBM Quantum Token" -ForegroundColor Yellow
    Write-Host "[INFO] 可使用 -IBMToken 參數或設定環境變數 IBM_QUANTUM_TOKEN" -ForegroundColor Gray
}

# 清理舊容器和映像（如果需要）
if ($Clean) {
    Write-Host "`n[3/6] 清理舊容器和映像..." -ForegroundColor Yellow
    
    # 停止並刪除容器
    Write-Host "  停止容器..." -ForegroundColor Gray
    docker-compose stop cyber-ai-quantum 2>$null
    docker-compose rm -f cyber-ai-quantum 2>$null
    
    # 刪除映像
    Write-Host "  刪除舊映像..." -ForegroundColor Gray
    docker rmi application-cyber-ai-quantum 2>$null
    
    # 清理無用映像
    Write-Host "  清理 dangling 映像..." -ForegroundColor Gray
    docker image prune -f
    
    Write-Host "[OK] 清理完成" -ForegroundColor Green
} else {
    Write-Host "`n[3/6] 跳過清理 (使用 -Clean 參數強制清理)" -ForegroundColor Yellow
}

# 重新建構映像
if (-not $NoBuild) {
    Write-Host "`n[4/6] 重新建構 cyber-ai-quantum 映像..." -ForegroundColor Yellow
    Write-Host "[INFO] 這可能需要幾分鐘時間..." -ForegroundColor Gray
    
    docker-compose build --no-cache cyber-ai-quantum
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "[ERROR] 映像建構失敗" -ForegroundColor Red
        Pop-Location
        exit 1
    }
    Write-Host "[OK] 映像建構成功" -ForegroundColor Green
} else {
    Write-Host "`n[4/6] 跳過建構 (使用現有映像)" -ForegroundColor Yellow
}

# 停止舊容器
Write-Host "`n[5/6] 停止舊容器..." -ForegroundColor Yellow
docker-compose stop cyber-ai-quantum
Write-Host "[OK] 容器已停止" -ForegroundColor Green

# 啟動新容器
Write-Host "`n[6/6] 啟動新容器..." -ForegroundColor Yellow
docker-compose up -d cyber-ai-quantum

if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] 容器啟動失敗" -ForegroundColor Red
    Write-Host "[INFO] 查看日誌: docker-compose logs cyber-ai-quantum" -ForegroundColor Yellow
    Pop-Location
    exit 1
}

Write-Host "[OK] 容器已啟動" -ForegroundColor Green

# 等待服務就緒
Write-Host "`n[等待] 等待服務就緒..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# 檢查健康狀態
Write-Host "`n[健康檢查] 檢查服務狀態..." -ForegroundColor Yellow
$maxAttempts = 10
$attempt = 0
$healthy = $false

while ($attempt -lt $maxAttempts -and -not $healthy) {
    $attempt++
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8000/health" -UseBasicParsing -TimeoutSec 2 -ErrorAction SilentlyContinue
        if ($response.StatusCode -eq 200) {
            $healthy = $true
            Write-Host "[OK] 服務健康檢查通過！" -ForegroundColor Green
        }
    } catch {
        Write-Host "  嘗試 $attempt/$maxAttempts : 服務尚未就緒..." -ForegroundColor Gray
        Start-Sleep -Seconds 3
    }
}

if (-not $healthy) {
    Write-Host "[WARNING] 健康檢查失敗，但容器可能仍在啟動中" -ForegroundColor Yellow
    Write-Host "[INFO] 請稍後手動檢查: http://localhost:8000/health" -ForegroundColor Yellow
}

# 顯示容器狀態
Write-Host "`n[狀態] 容器資訊:" -ForegroundColor Yellow
docker ps --filter "name=cyber-ai-quantum" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# 顯示日誌
Write-Host "`n[日誌] 最近的日誌 (最後 20 行):" -ForegroundColor Yellow
Write-Host "============================================================================" -ForegroundColor Cyan
docker-compose logs --tail=20 cyber-ai-quantum
Write-Host "============================================================================" -ForegroundColor Cyan

# 管理指令提示
Write-Host "`n============================================================================" -ForegroundColor Cyan
Write-Host "  部署完成！" -ForegroundColor Green
Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host "`n常用指令:" -ForegroundColor Yellow
Write-Host "  查看日誌:       docker-compose logs -f cyber-ai-quantum" -ForegroundColor White
Write-Host "  進入容器:       docker exec -it cyber-ai-quantum /bin/bash" -ForegroundColor White
Write-Host "  重新啟動:       docker-compose restart cyber-ai-quantum" -ForegroundColor White
Write-Host "  停止服務:       docker-compose stop cyber-ai-quantum" -ForegroundColor White
Write-Host "  查看狀態:       docker-compose ps cyber-ai-quantum" -ForegroundColor White
Write-Host ""
Write-Host "API 端點:" -ForegroundColor Yellow
Write-Host "  健康檢查:       http://localhost:8000/health" -ForegroundColor White
Write-Host "  API 文檔:       http://localhost:8000/docs" -ForegroundColor White
Write-Host "  接收日誌:       POST http://localhost:8000/api/v1/agent/log" -ForegroundColor White
Write-Host "  Zero Trust:     POST http://localhost:8000/api/v1/zerotrust/predict" -ForegroundColor White
Write-Host ""

Pop-Location
Write-Host "[完成] 所有操作完成！" -ForegroundColor Green
Write-Host ""

