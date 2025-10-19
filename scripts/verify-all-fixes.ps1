# ============================================================================
# 驗證所有修復 - 一鍵檢查腳本
# ============================================================================
# 用途: 驗證 SAST 修復、量子 ML 實作、Docker 部署
# 執行方式: .\verify-all-fixes.ps1
# ============================================================================

param(
    [Parameter(Mandatory=$false)]
    [switch]$SkipDocker = $false,
    
    [Parameter(Mandatory=$false)]
    [switch]$Verbose = $false
)

$ErrorActionPreference = "Continue"

Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host "  驗證所有修復 - Pandora v3.4.1" -ForegroundColor Cyan
Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host ""

$passed = 0
$failed = 0

# 切換到專案根目錄
$projectRoot = "C:\Users\dennis.lee\Documents\GitHub\Local_IPS-IDS"
Set-Location $projectRoot

# ============================================================================
# 1. 驗證 Go 依賴版本
# ============================================================================
Write-Host "[1/8] 檢查 Go 依賴版本..." -ForegroundColor Yellow

$requiredVersions = @{
    "golang.org/x/crypto" = "v0.43.0"
    "golang.org/x/net" = "v0.46.0"
    "golang.org/x/oauth2" = "v0.30.0"
}

$goList = go list -m all | Out-String

foreach ($pkg in $requiredVersions.Keys) {
    $required = $requiredVersions[$pkg]
    if ($goList -match "$pkg\s+($required|v0\.\d+\.\d+)") {
        $actualVersion = $matches[1]
        if ($actualVersion -ge $required) {
            Write-Host "  ✅ $pkg $actualVersion" -ForegroundColor Green
            $passed++
        } else {
            Write-Host "  ❌ $pkg $actualVersion (需要 $required)" -ForegroundColor Red
            $failed++
        }
    } else {
        Write-Host "  ❌ $pkg 未找到" -ForegroundColor Red
        $failed++
    }
}

# ============================================================================
# 2. 驗證建構
# ============================================================================
Write-Host "`n[2/8] 測試 Go 建構..." -ForegroundColor Yellow

try {
    $buildOutput = go build -o bin/test-verify.exe ./cmd/main.go 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✅ 建構成功" -ForegroundColor Green
        $passed++
        Remove-Item bin/test-verify.exe -ErrorAction SilentlyContinue
    } else {
        Write-Host "  ❌ 建構失敗" -ForegroundColor Red
        if ($Verbose) { Write-Host $buildOutput }
        $failed++
    }
} catch {
    Write-Host "  ❌ 建構失敗: $_" -ForegroundColor Red
    $failed++
}

# ============================================================================
# 3. 檢查量子 ML 檔案
# ============================================================================
Write-Host "`n[3/8] 檢查量子 ML 檔案..." -ForegroundColor Yellow

$quantumFiles = @(
    "Experimental/cyber-ai-quantum/feature_extractor.py",
    "Experimental/cyber-ai-quantum/generate_dynamic_qasm.py",
    "Experimental/cyber-ai-quantum/train_quantum_classifier.py",
    "Experimental/cyber-ai-quantum/daily_quantum_job.py",
    "Experimental/cyber-ai-quantum/analyze_results.py"
)

foreach ($file in $quantumFiles) {
    if (Test-Path $file) {
        Write-Host "  ✅ $($file.Split('/')[-1])" -ForegroundColor Green
        $passed++
    } else {
        Write-Host "  ❌ $file 未找到" -ForegroundColor Red
        $failed++
    }
}

# ============================================================================
# 4. 檢查 Docker 檔案
# ============================================================================
Write-Host "`n[4/8] 檢查 Docker 配置..." -ForegroundColor Yellow

$dockerFiles = @(
    "Experimental/cyber-ai-quantum/Dockerfile",
    "Application/docker-compose.yml",
    "Application/rebuild-quantum.ps1",
    "Application/rebuild-quantum.sh"
)

foreach ($file in $dockerFiles) {
    if (Test-Path $file) {
        Write-Host "  ✅ $($file.Split('/')[-1])" -ForegroundColor Green
        $passed++
    } else {
        Write-Host "  ❌ $file 未找到" -ForegroundColor Red
        $failed++
    }
}

# ============================================================================
# 5. 檢查文檔
# ============================================================================
Write-Host "`n[5/8] 檢查文檔..." -ForegroundColor Yellow

$docFiles = @(
    "Experimental/cyber-ai-quantum/README-QUANTUM-TESTING.md",
    "Experimental/QUANTUM-ML-IMPLEMENTATION-COMPLETE.md",
    "Experimental/cyber-ai-quantum/FIXES-APPLIED.md",
    "SAST/2025-10-15-FIXES.md",
    "Experimental/ALL-FIXES-COMPLETE-v3.4.1.md"
)

foreach ($file in $docFiles) {
    if (Test-Path $file) {
        Write-Host "  ✅ $($file.Split('/')[-1])" -ForegroundColor Green
        $passed++
    } else {
        Write-Host "  ❌ $file 未找到" -ForegroundColor Red
        $failed++
    }
}

# ============================================================================
# 6. 驗證 Python 環境
# ============================================================================
Write-Host "`n[6/8] 檢查 Python 環境..." -ForegroundColor Yellow

try {
    $pythonVersion = python --version 2>&1
    Write-Host "  ✅ Python: $pythonVersion" -ForegroundColor Green
    $passed++
} catch {
    Write-Host "  ❌ Python 未安裝" -ForegroundColor Red
    $failed++
}

# ============================================================================
# 7. 測試 Docker（可選）
# ============================================================================
if (-not $SkipDocker) {
    Write-Host "`n[7/8] 測試 Docker..." -ForegroundColor Yellow
    
    try {
        $dockerVersion = docker version --format '{{.Server.Version}}' 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✅ Docker: $dockerVersion" -ForegroundColor Green
            $passed++
            
            # 檢查容器是否運行
            $containerStatus = docker ps --filter "name=cyber-ai-quantum" --format "{{.Status}}"
            if ($containerStatus -match "Up") {
                Write-Host "  ✅ cyber-ai-quantum 容器運行中" -ForegroundColor Green
                $passed++
                
                # 測試健康檢查
                try {
                    $health = Invoke-WebRequest -Uri "http://localhost:8000/health" -UseBasicParsing -TimeoutSec 2 -ErrorAction Stop
                    if ($health.StatusCode -eq 200) {
                        Write-Host "  ✅ 健康檢查通過" -ForegroundColor Green
                        $passed++
                    }
                } catch {
                    Write-Host "  ⚠️  健康檢查失敗（容器可能未完全啟動）" -ForegroundColor Yellow
                }
            } else {
                Write-Host "  ⚠️  cyber-ai-quantum 容器未運行" -ForegroundColor Yellow
                Write-Host "     提示: 執行 'docker-compose up -d cyber-ai-quantum'" -ForegroundColor Gray
            }
        }
    } catch {
        Write-Host "  ❌ Docker 未運行或未安裝" -ForegroundColor Red
        $failed++
    }
} else {
    Write-Host "`n[7/8] 跳過 Docker 測試 (-SkipDocker)" -ForegroundColor Yellow
}

# ============================================================================
# 8. 驗證 git 狀態
# ============================================================================
Write-Host "`n[8/8] 檢查 Git 狀態..." -ForegroundColor Yellow

$gitStatus = git status --short
$modifiedFiles = ($gitStatus | Measure-Object).Count

if ($modifiedFiles -gt 0) {
    Write-Host "  📝 有 $modifiedFiles 個檔案已修改" -ForegroundColor Cyan
    if ($Verbose) {
        Write-Host "`n修改的檔案:" -ForegroundColor Gray
        git status --short
    }
} else {
    Write-Host "  ✅ 工作目錄乾淨" -ForegroundColor Green
}

# ============================================================================
# 總結
# ============================================================================
Write-Host "`n============================================================================" -ForegroundColor Cyan
Write-Host "  驗證完成！" -ForegroundColor Cyan
Write-Host "============================================================================" -ForegroundColor Cyan

$total = $passed + $failed
$percentage = if ($total -gt 0) { [math]::Round(($passed / $total) * 100, 1) } else { 0 }

Write-Host "`n測試結果:" -ForegroundColor Yellow
Write-Host "  ✅ 通過: $passed" -ForegroundColor Green
Write-Host "  ❌ 失敗: $failed" -ForegroundColor Red
Write-Host "  📊 成功率: $percentage%" -ForegroundColor $(if ($percentage -ge 90) { "Green" } elseif ($percentage -ge 70) { "Yellow" } else { "Red" })

if ($percentage -ge 90) {
    Write-Host "`n🎉 狀態: 優秀！所有主要功能正常運作" -ForegroundColor Green
} elseif ($percentage -ge 70) {
    Write-Host "`n⚠️  狀態: 良好，但仍有部分需要改進" -ForegroundColor Yellow
} else {
    Write-Host "`n❌ 狀態: 需要修復失敗的項目" -ForegroundColor Red
}

Write-Host "`n詳細報告:" -ForegroundColor Yellow
Write-Host "  - 量子 ML 實作: Experimental/QUANTUM-ML-IMPLEMENTATION-COMPLETE.md" -ForegroundColor White
Write-Host "  - SAST 修復: SAST/2025-10-15-FIXES.md" -ForegroundColor White
Write-Host "  - 完整總覽: Experimental/ALL-FIXES-COMPLETE-v3.4.1.md" -ForegroundColor White

Write-Host "`n下一步建議:" -ForegroundColor Yellow
if ($failed -eq 0) {
    Write-Host "  1. 執行完整測試: go test ./..." -ForegroundColor White
    Write-Host "  2. 提交變更: git add . && git commit -m 'feat: quantum ML + SAST fixes v3.4.1'" -ForegroundColor White
    Write-Host "  3. 重新掃描: snyk test" -ForegroundColor White
} else {
    Write-Host "  1. 查看失敗項目並修復" -ForegroundColor White
    Write-Host "  2. 重新執行驗證: .\verify-all-fixes.ps1" -ForegroundColor White
}

Write-Host ""

