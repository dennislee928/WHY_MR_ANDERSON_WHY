# Terraform Windows Setup Script
# PowerShell 腳本用於 Windows 環境設定

Write-Host "🚀 Pandora Box Console - Terraform Setup" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green
Write-Host ""

# 檢查是否在 terraform 目錄
if (!(Test-Path "terraform.tfvars.example")) {
    Write-Host "❌ 錯誤: 請在 terraform 目錄中執行此腳本" -ForegroundColor Red
    Write-Host "   cd terraform" -ForegroundColor Yellow
    exit 1
}

# 步驟 1: 複製變數範本
Write-Host "📋 步驟 1: 建立 terraform.tfvars..." -ForegroundColor Cyan
if (Test-Path "terraform.tfvars") {
    Write-Host "   ⚠️  terraform.tfvars 已存在" -ForegroundColor Yellow
    $overwrite = Read-Host "   要覆蓋嗎? (y/N)"
    if ($overwrite -ne "y") {
        Write-Host "   ⏭️  跳過建立 terraform.tfvars" -ForegroundColor Yellow
    } else {
        Copy-Item "terraform.tfvars.example" "terraform.tfvars"
        Write-Host "   ✅ 已建立 terraform.tfvars" -ForegroundColor Green
    }
} else {
    Copy-Item "terraform.tfvars.example" "terraform.tfvars"
    Write-Host "   ✅ 已建立 terraform.tfvars" -ForegroundColor Green
}

Write-Host ""

# 步驟 2: 檢查 Terraform 安裝
Write-Host "🔍 步驟 2: 檢查 Terraform..." -ForegroundColor Cyan
try {
    $terraformVersion = terraform version
    Write-Host "   ✅ Terraform 已安裝" -ForegroundColor Green
    Write-Host "   版本: $($terraformVersion[0])" -ForegroundColor Gray
} catch {
    Write-Host "   ❌ Terraform 未安裝" -ForegroundColor Red
    Write-Host ""
    Write-Host "   請安裝 Terraform:" -ForegroundColor Yellow
    Write-Host "   1. 下載: https://www.terraform.io/downloads" -ForegroundColor Yellow
    Write-Host "   2. 或使用 Chocolatey: choco install terraform" -ForegroundColor Yellow
    Write-Host "   3. 或使用 Scoop: scoop install terraform" -ForegroundColor Yellow
    exit 1
}

Write-Host ""

# 步驟 3: 初始化 Terraform
Write-Host "⚙️  步驟 3: 初始化 Terraform?" -ForegroundColor Cyan
$init = Read-Host "   要執行 terraform init 嗎? (Y/n)"
if ($init -ne "n") {
    Write-Host "   執行 terraform init..." -ForegroundColor Gray
    terraform init
    if ($LASTEXITCODE -eq 0) {
        Write-Host "   ✅ Terraform 初始化成功" -ForegroundColor Green
    } else {
        Write-Host "   ❌ Terraform 初始化失敗" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "   ⏭️  跳過初始化" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Green
Write-Host "✅ 設定完成!" -ForegroundColor Green
Write-Host ""
Write-Host "📝 下一步:" -ForegroundColor Cyan
Write-Host "   1. 編輯 terraform.tfvars 填入您的 API tokens" -ForegroundColor White
Write-Host "      code terraform.tfvars" -ForegroundColor Gray
Write-Host "      或" -ForegroundColor Gray
Write-Host "      notepad terraform.tfvars" -ForegroundColor Gray
Write-Host ""
Write-Host "   2. 驗證配置" -ForegroundColor White
Write-Host "      terraform validate" -ForegroundColor Gray
Write-Host ""
Write-Host "   3. 查看計劃" -ForegroundColor White
Write-Host "      terraform plan" -ForegroundColor Gray
Write-Host ""
Write-Host "   4. 應用配置" -ForegroundColor White
Write-Host "      terraform apply" -ForegroundColor Gray
Write-Host ""
Write-Host "📚 查看完整文件:" -ForegroundColor Cyan
Write-Host "   README.md" -ForegroundColor Gray
Write-Host "   DEPLOYMENT-CHECKLIST.md" -ForegroundColor Gray
Write-Host ""

