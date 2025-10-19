# Pandora Box Console - 安全修復應用腳本 (PowerShell)
# 自動應用 SAST 掃描發現的安全修復

$ErrorActionPreference = "Stop"

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "  🔒 Pandora 安全修復應用工具" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

# 1. 更新 Go 依賴
Write-Host "`n📦 步驟 1/5: 更新 Go 依賴..." -ForegroundColor Yellow
go mod tidy
go mod download
Write-Host "✅ Go 依賴已更新" -ForegroundColor Green

# 2. 更新 Python 依賴
Write-Host "`n📦 步驟 2/5: 更新 Python 依賴..." -ForegroundColor Yellow
Set-Location Experimental/cyber-ai-quantum
pip install -r requirements.txt --upgrade --quiet
Set-Location ../..
Write-Host "✅ Python 依賴已更新" -ForegroundColor Green

# 3. 驗證 Dockerfile USER 指令
Write-Host "`n🔍 步驟 3/5: 驗證 Dockerfile 安全性..." -ForegroundColor Yellow
$dockerfiles = @(
    "Application/docker/agent.koyeb.dockerfile",
    "Application/docker/monitoring.dockerfile",
    "Application/docker/nginx.dockerfile",
    "Application/docker/test.dockerfile",
    "Application/docker/axiom-be.dockerfile"
)

foreach ($dockerfile in $dockerfiles) {
    if (Select-String -Path $dockerfile -Pattern "^USER " -Quiet) {
        Write-Host "  ✅ $dockerfile - USER 指令已存在" -ForegroundColor Green
    } else {
        Write-Host "  ❌ $dockerfile - 缺少 USER 指令" -ForegroundColor Red
    }
}

# 4. 檢查 Alpine 版本
Write-Host "`n🔍 步驟 4/5: 檢查 Alpine 基礎映像版本..." -ForegroundColor Yellow
Get-ChildItem Application/docker/*.dockerfile | ForEach-Object {
    $content = Get-Content $_.FullName -Raw
    if ($content -match "FROM alpine:3\.21" -or $content -match "FROM alpine:3\.22") {
        Write-Host "  ✅ $($_.Name) - Alpine 版本安全" -ForegroundColor Green
    } elseif ($content -match "FROM alpine:") {
        Write-Host "  ⚠️  $($_.Name) - 建議更新到 Alpine 3.21+" -ForegroundColor Yellow
    }
}

# 5. 重新構建關鍵服務
Write-Host "`n🔨 步驟 5/5: 重新構建 Docker 映像..." -ForegroundColor Yellow
Set-Location Application

Write-Host "  構建 axiom-be..." -ForegroundColor Yellow
docker-compose build --no-cache axiom-be

Write-Host "  構建 cyber-ai-quantum..." -ForegroundColor Yellow
docker-compose build --no-cache cyber-ai-quantum

Set-Location ..
Write-Host "✅ Docker 映像已重新構建" -ForegroundColor Green

# 完成
Write-Host "`n=========================================" -ForegroundColor Cyan
Write-Host "  ✅ 安全修復應用完成！" -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "`n📋 下一步:" -ForegroundColor Cyan
Write-Host "  1. 查看詳細報告: " -NoNewline
Write-Host "docs/SAST-SECURITY-FIXES.md" -ForegroundColor Yellow
Write-Host "  2. 重啟服務: " -NoNewline
Write-Host "cd Application && docker-compose up -d" -ForegroundColor Yellow
Write-Host "  3. 驗證服務: " -NoNewline
Write-Host "docker-compose ps" -ForegroundColor Yellow
Write-Host "`n⚠️  需要手動處理的項目:" -ForegroundColor Yellow
Write-Host "  - 配置 gRPC TLS 證書"
Write-Host "  - 修復 exec.Command 輸入驗證"
Write-Host "  - 修復 RWMutex 死鎖風險"
Write-Host "  - 更新 GitHub Actions 配置`n"

