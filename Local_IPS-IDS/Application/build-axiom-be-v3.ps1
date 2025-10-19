# Axiom Backend V3 本地構建和測試腳本 (PowerShell)
# 版本: 3.1.0
# 使用方式：
#   .\build-axiom-be-v3.ps1           # 正常構建
#   .\build-axiom-be-v3.ps1 -NoCache  # 不使用快取構建

param(
    [switch]$NoCache
)

# 切換到腳本所在目錄
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $ScriptDir

Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  Axiom Backend V3 本地構建和測試" -ForegroundColor Cyan
Write-Host "  版本: 3.1.0" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "工作目錄: $ScriptDir" -ForegroundColor Gray
Write-Host ""

# 設置變量
$IMAGE_NAME = "axiom-backend"
$IMAGE_TAG = "v3.1.0"
$CONTAINER_NAME = "axiom-be-v3-test"
$PORT = 3001

# Step 1: 清理舊容器和鏡像
Write-Host "[1/6] 清理舊容器..." -ForegroundColor Yellow

# 清理目標測試容器
$existingContainer = docker ps -a --format "{{.Names}}" | Where-Object { $_ -eq $CONTAINER_NAME }
if ($existingContainer) {
    Write-Host "  停止並刪除容器: $CONTAINER_NAME" -ForegroundColor Gray
    docker stop $CONTAINER_NAME 2>$null | Out-Null
    docker rm $CONTAINER_NAME 2>$null | Out-Null
}

# 清理所有佔用 3001 端口的容器
$portContainers = docker ps -a --filter "publish=$PORT" --format "{{.Names}}" 2>$null
if ($portContainers) {
    Write-Host "  發現佔用端口 $PORT 的容器:" -ForegroundColor Gray
    foreach ($container in $portContainers) {
        Write-Host "    - $container" -ForegroundColor Gray
        docker stop $container 2>$null | Out-Null
        docker rm $container 2>$null | Out-Null
    }
}

Write-Host "完成" -ForegroundColor Green
Write-Host ""

# Step 2: 構建 Docker 鏡像
Write-Host "[2/6] 構建 Docker 鏡像..." -ForegroundColor Yellow
Write-Host "鏡像: ${IMAGE_NAME}:${IMAGE_TAG}" -ForegroundColor Gray

$BUILD_START = Get-Date

# 構建參數
$buildArgs = @(
    "-f", "docker/axiom-be.dockerfile",
    "-t", "${IMAGE_NAME}:${IMAGE_TAG}",
    "-t", "${IMAGE_NAME}:latest"
)

if ($NoCache) {
    $buildArgs += "--no-cache"
    Write-Host "使用 --no-cache 選項" -ForegroundColor Gray
}

$buildArgs += "be"

docker build @buildArgs

if ($LASTEXITCODE -ne 0) {
    Write-Host "構建失敗！" -ForegroundColor Red
    exit 1
}

$BUILD_END = Get-Date
$BUILD_TIME = ($BUILD_END - $BUILD_START).TotalSeconds
Write-Host "構建成功！耗時: $([math]::Round($BUILD_TIME, 2))秒" -ForegroundColor Green
Write-Host ""

# Step 3: 檢查鏡像大小
Write-Host "[3/6] 檢查鏡像..." -ForegroundColor Yellow
docker images | Select-String -Pattern $IMAGE_NAME | Select-Object -First 1
Write-Host ""

# Step 4: 啟動容器（單機測試）
Write-Host "[4/6] 啟動測試容器..." -ForegroundColor Yellow
docker run -d `
    --name $CONTAINER_NAME `
    -p ${PORT}:3001 `
    -e POSTGRES_HOST=host.docker.internal `
    -e POSTGRES_PORT=5432 `
    -e POSTGRES_USER=pandora `
    -e POSTGRES_PASSWORD=pandora123 `
    -e POSTGRES_DB=pandora_db `
    -e REDIS_HOST=host.docker.internal `
    -e REDIS_PORT=6379 `
    -e REDIS_PASSWORD=pandora123 `
    -e PROMETHEUS_URL=http://host.docker.internal:9090 `
    -e LOKI_URL=http://host.docker.internal:3100 `
    -e QUANTUM_URL=http://host.docker.internal:8000 `
    ${IMAGE_NAME}:${IMAGE_TAG}

if ($LASTEXITCODE -ne 0) {
    Write-Host "容器啟動失敗！" -ForegroundColor Red
    exit 1
}
Write-Host "容器已啟動" -ForegroundColor Green
Write-Host ""

# Step 5: 等待服務啟動
Write-Host "[5/6] 等待服務就緒..." -ForegroundColor Yellow
$RETRY = 0
$MAX_RETRY = 30

while ($RETRY -lt $MAX_RETRY) {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:${PORT}/health" -UseBasicParsing -TimeoutSec 2
        if ($response.StatusCode -eq 200) {
            Write-Host "服務已就緒！" -ForegroundColor Green
            break
        }
    } catch {
        $RETRY++
        Write-Host "  等待中... ($RETRY/$MAX_RETRY)" -ForegroundColor Gray
        Start-Sleep -Seconds 2
    }
}

if ($RETRY -eq $MAX_RETRY) {
    Write-Host "服務啟動超時！" -ForegroundColor Red
    Write-Host "查看日誌:" -ForegroundColor Yellow
    docker logs $CONTAINER_NAME
    exit 1
}
Write-Host ""

# Step 6: 測試 API
Write-Host "[6/6] 測試 API 端點..." -ForegroundColor Yellow

# 測試健康檢查
Write-Host "  測試健康檢查..." -ForegroundColor Gray
try {
    $health = Invoke-RestMethod -Uri "http://localhost:${PORT}/health"
    if ($health.status -eq "healthy") {
        Write-Host "  健康檢查通過: $($health.service) v$($health.version)" -ForegroundColor Green
    }
} catch {
    Write-Host "  健康檢查失敗: $_" -ForegroundColor Red
}

# 測試 Agent 註冊
Write-Host "  測試 Agent 註冊..." -ForegroundColor Gray
try {
    $agentBody = @{
        mode = "internal"
        hostname = "test-server"
        ip_address = "127.0.0.1"
        capabilities = @("windows_logs")
    } | ConvertTo-Json

    $agentResponse = Invoke-RestMethod -Uri "http://localhost:${PORT}/api/v2/agent/register" -Method Post -Body $agentBody -ContentType "application/json"
    if ($agentResponse.success) {
        Write-Host "  Agent 註冊成功: $($agentResponse.data.agent_id)" -ForegroundColor Green
    }
} catch {
    Write-Host "  Agent 註冊失敗: $_" -ForegroundColor Red
}

# 測試 PII 檢測
Write-Host "  測試 PII 檢測..." -ForegroundColor Gray
try {
    $piiBody = @{
        text = "Contact: test@example.com, Card: 4532-1234-5678-9010"
    } | ConvertTo-Json

    $piiResponse = Invoke-RestMethod -Uri "http://localhost:${PORT}/api/v2/compliance/pii/detect" -Method Post -Body $piiBody -ContentType "application/json"
    if ($piiResponse.success -and $piiResponse.data.pii_found) {
        Write-Host "  PII 檢測成功: 發現 $($piiResponse.data.matches.Count) 個 PII" -ForegroundColor Green
    }
} catch {
    Write-Host "  PII 檢測失敗: $_" -ForegroundColor Red
}

# 測試儲存統計
Write-Host "  測試儲存統計..." -ForegroundColor Gray
try {
    $storageResponse = Invoke-RestMethod -Uri "http://localhost:${PORT}/api/v2/storage/tiers/stats"
    if ($storageResponse.success) {
        Write-Host "  儲存統計成功" -ForegroundColor Green
    }
} catch {
    Write-Host "  儲存統計失敗: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "  構建和測試完成！" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "容器信息:" -ForegroundColor Yellow
Write-Host "  名稱: $CONTAINER_NAME" -ForegroundColor White
Write-Host "  端口: http://localhost:$PORT" -ForegroundColor White
Write-Host ""
Write-Host "常用指令:" -ForegroundColor Yellow
Write-Host "  查看日誌: docker logs $CONTAINER_NAME" -ForegroundColor White
Write-Host "  查看日誌(實時): docker logs -f $CONTAINER_NAME" -ForegroundColor White
Write-Host "  停止容器: docker stop $CONTAINER_NAME" -ForegroundColor White
Write-Host "  刪除容器: docker rm $CONTAINER_NAME" -ForegroundColor White
Write-Host ""
Write-Host "API 端點:" -ForegroundColor Yellow
Write-Host "  健康檢查: http://localhost:$PORT/health" -ForegroundColor White
Write-Host "  API 文檔: 查看 docs/AXIOM-BACKEND-V3-API-DOCUMENTATION.md" -ForegroundColor White
Write-Host ""
