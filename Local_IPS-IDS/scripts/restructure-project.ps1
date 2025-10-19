# 專案結構重整腳本
# 請在執行前備份專案！

param(
    [switch]$DryRun = $false
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "   Pandora Box 專案結構重整腳本" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

if ($DryRun) {
    Write-Host "⚠️  DRY RUN 模式 - 不會實際移動文件`n" -ForegroundColor Yellow
}

# 記錄操作
$script:operations = @()

function Log-Operation {
    param($From, $To, $Type = "MOVE")
    $script:operations += [PSCustomObject]@{
        Type = $Type
        From = $From
        To   = $To
    }
    Write-Host "  [$Type] $From -> $To" -ForegroundColor Gray
}

function Move-FileIfExists {
    param($Source, $Destination)
    
    if (Test-Path $Source) {
        Log-Operation -From $Source -To $Destination -Type "MOVE"
        if (-not $DryRun) {
            $destDir = Split-Path -Parent $Destination
            if (-not (Test-Path $destDir)) {
                New-Item -ItemType Directory -Force -Path $destDir | Out-Null
            }
            Move-Item -Force $Source $Destination
        }
        return $true
    }
    return $false
}

# 階段 1: 移動編譯產物
Write-Host "`n📦 階段 1: 移動編譯產物" -ForegroundColor Green
Write-Host "─────────────────────────────" -ForegroundColor Green
Move-FileIfExists "agent.exe" "bin/agent.exe"
Move-FileIfExists "console.exe" "bin/console.exe"

# 階段 2: 移動 Dockerfile
Write-Host "`n🐳 階段 2: 移動 Dockerfile" -ForegroundColor Green
Write-Host "─────────────────────────────" -ForegroundColor Green
Move-FileIfExists "Dockerfile.agent" "build/docker/agent.dockerfile"
Move-FileIfExists "Dockerfile.agent.koyeb" "build/docker/agent.koyeb.dockerfile"
Move-FileIfExists "Dockerfile.monitoring" "build/docker/monitoring.dockerfile"
Move-FileIfExists "Dockerfile.nginx" "build/docker/nginx.dockerfile"
Move-FileIfExists "Dockerfile.server-be" "build/docker/server-be.dockerfile"
Move-FileIfExists "Dockerfile.server-fe" "build/docker/server-fe.dockerfile"
Move-FileIfExists "Dockerfile.test" "build/docker/test.dockerfile"
Move-FileIfExists "Dockerfile.ui.patr" "build/docker/ui.patr.dockerfile"

# 階段 3: 移動文檔
Write-Host "`n📚 階段 3: 移動文檔" -ForegroundColor Green
Write-Host "─────────────────────────────" -ForegroundColor Green

# 部署文檔
Move-FileIfExists "DEPLOYMENT.md" "docs/deployment/README.md"
Move-FileIfExists "DEPLOYMENT-GCP.md" "docs/deployment/gcp.md"
Move-FileIfExists "DEPLOYMENT-SUMMARY.md" "docs/deployment/summary.md"
Move-FileIfExists "DEPLOY-SPEC.MD" "docs/deployment/spec.md"
Move-FileIfExists "README-DEPLOYMENT.md" "docs/deployment/quickstart.md"
Move-FileIfExists "README-PAAS-DEPLOYMENT.md" "docs/deployment/paas.md"

# Fly.io 文檔
Move-FileIfExists "FLYIO-TROUBLESHOOTING.md" "docs/deployment/flyio/troubleshooting.md"
Move-FileIfExists "FLYIO-GRAFANA-FIX.md" "docs/deployment/flyio/grafana-fix.md"
Move-FileIfExists "FLYIO-NEXTJS-CONFLICT-FIX.md" "docs/deployment/flyio/nextjs-conflict-fix.md"
Move-FileIfExists "FLYIO-NEXTJS-TEMPORARY-FIX.md" "docs/deployment/flyio/nextjs-temporary-fix.md"
Move-FileIfExists "FLYIO-VOLUME-FIX.md" "docs/deployment/flyio/volume-fix.md"

# Koyeb 文檔
Move-FileIfExists "KOYEB-DEPLOYMENT-GUIDE.md" "docs/deployment/koyeb/deployment-guide.md"
Move-FileIfExists "KOYEB-FIX-SUMMARY.md" "docs/deployment/koyeb/fix-summary.md"
Move-FileIfExists "KOYEB-QUICK-START.md" "docs/deployment/koyeb/quickstart.md"

# 其他文檔
Move-FileIfExists "FINAL-STATUS.md" "docs/operations/final-status.md"
Move-FileIfExists "FIXES-SUMMARY.md" "docs/operations/fixes-summary.md"
Move-FileIfExists "IMPLEMENTATION-SUMMARY.md" "docs/development/implementation-summary.md"
Move-FileIfExists "TERRAFORM-IMPLEMENTATION-SUMMARY.md" "docs/deployment/terraform-summary.md"
Move-FileIfExists "MQTT-PUBSUB-RATELIMIT-LOADBALANCER.md" "docs/architecture/modules.md"
Move-FileIfExists "WINDOWS-SETUP-COMPLETE.md" "docs/development/windows-setup.md"

# 階段 4: 移動 Docker Compose
Write-Host "`n🐋 階段 4: 移動 Docker Compose" -ForegroundColor Green
Write-Host "─────────────────────────────" -ForegroundColor Green
Move-FileIfExists "docker-compose.yml" "deployments/docker-compose/docker-compose.yml"
Move-FileIfExists "docker-compose.test.yml" "deployments/docker-compose/docker-compose.test.yml"

# 階段 5: 移動 K8s 配置
Write-Host "`n☸️  階段 5: 移動 K8s 配置" -ForegroundColor Green
Write-Host "─────────────────────────────" -ForegroundColor Green

if (Test-Path "k8s") {
    Log-Operation -From "k8s/*" -To "deployments/kubernetes/base/" -Type "COPY"
    if (-not $DryRun) {
        Copy-Item -Recurse -Force "k8s/*" "deployments/kubernetes/base/"
    }
}

if (Test-Path "k8s-gcp") {
    Log-Operation -From "k8s-gcp/*" -To "deployments/kubernetes/gcp/" -Type "COPY"
    if (-not $DryRun) {
        Copy-Item -Recurse -Force "k8s-gcp/*" "deployments/kubernetes/gcp/"
    }
}

# 階段 6: 移動 Terraform
Write-Host "`n🏗️  階段 6: 移動 Terraform" -ForegroundColor Green
Write-Host "─────────────────────────────" -ForegroundColor Green

if (Test-Path "terraform") {
    Log-Operation -From "terraform/*" -To "deployments/terraform/" -Type "COPY"
    if (-not $DryRun) {
        Copy-Item -Recurse -Force "terraform/*" "deployments/terraform/"
    }
}

# 階段 7: 移動 PaaS 配置
Write-Host "`n☁️  階段 7: 移動 PaaS 配置" -ForegroundColor Green
Write-Host "─────────────────────────────" -ForegroundColor Green
Move-FileIfExists "fly.toml" "deployments/paas/flyio/fly.toml"
Move-FileIfExists "fly-monitoring.toml" "deployments/paas/flyio/fly-monitoring.toml"
Move-FileIfExists "koyeb.yaml" "deployments/paas/koyeb/koyeb.yaml"
Move-FileIfExists ".koyeb/config.yaml" "deployments/paas/koyeb/config.yaml"
Move-FileIfExists "railway.json" "deployments/paas/railway/railway.json"
Move-FileIfExists "railway.toml" "deployments/paas/railway/railway.toml"
Move-FileIfExists "render.yaml" "deployments/paas/render/render.yaml"
Move-FileIfExists "patr.yaml" "deployments/paas/patr/patr.yaml"
Move-FileIfExists "vercel.json" "deployments/paas/vercel/vercel.json"

# 階段 8: 移動備份文件到 docs/archive
Write-Host "`n📦 階段 8: 移動備份文件" -ForegroundColor Green
Write-Host "─────────────────────────────" -ForegroundColor Green

if (-not $DryRun) {
    New-Item -ItemType Directory -Force -Path "docs/archive" | Out-Null
}

Move-FileIfExists "next.config.js.backup" "docs/archive/next.config.js.backup"
Move-FileIfExists "package.json.backup" "docs/archive/package.json.backup"
Move-FileIfExists "tailwind.config.js.backup" "docs/archive/tailwind.config.js.backup"
Move-FileIfExists "tsconfig.json.backup" "docs/archive/tsconfig.json.backup"
Move-FileIfExists "vercel.json.backup" "docs/archive/vercel.json.backup"

# 總結
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "   重整完成！" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

Write-Host "總共執行了 $($script:operations.Count) 個操作`n" -ForegroundColor Green

if ($DryRun) {
    Write-Host "這是 DRY RUN 模式的結果" -ForegroundColor Yellow
    Write-Host "請檢查上述操作，確認無誤後執行:`n" -ForegroundColor Yellow
    Write-Host "  .\scripts\restructure-project.ps1`n" -ForegroundColor Cyan
}

# 導出操作日誌
$logFile = "docs/restructure-operations.csv"
$script:operations | Export-Csv -Path $logFile -NoTypeInformation -Encoding UTF8

Write-Host "操作日誌已保存到: $logFile" -ForegroundColor Gray

Write-Host "`n⚠️  下一步操作:" -ForegroundColor Yellow
Write-Host "  1. 更新 .gitignore" -ForegroundColor White
Write-Host "  2. 更新 Makefile" -ForegroundColor White
Write-Host "  3. 更新 CI/CD workflows" -ForegroundColor White
Write-Host "  4. 測試建置流程`n" -ForegroundColor White

