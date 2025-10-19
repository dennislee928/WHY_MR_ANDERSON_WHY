# ============================================================================
# Pandora Cyber AI/Quantum - 每日量子作業排程腳本 (Windows)
# ============================================================================
# 用途: 設定 Windows 工作排程器，每日 00:00 (TPE 時區) 執行量子分類作業
# 執行方式: 以管理員身分執行此 PowerShell 腳本
# ============================================================================

param(
    [Parameter(Mandatory=$false)]
    [string]$PythonPath = "python",
    
    [Parameter(Mandatory=$false)]
    [string]$ScriptPath = "$PSScriptRoot\daily_quantum_job.py",
    
    [Parameter(Mandatory=$false)]
    [string]$TaskName = "PandoraQuantumDailyJob",
    
    [Parameter(Mandatory=$false)]
    [string]$ExecutionTime = "00:00",
    
    [Parameter(Mandatory=$false)]
    [string]$LogPath = "$PSScriptRoot\logs\daily_job.log"
)

Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host "  Pandora Cyber AI/Quantum - 每日量子作業排程設定" -ForegroundColor Cyan
Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host ""

# 檢查管理員權限
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host "[ERROR] 此腳本需要管理員權限才能設定工作排程器" -ForegroundColor Red
    Write-Host "[INFO] 請以管理員身分重新執行此腳本" -ForegroundColor Yellow
    Read-Host "按 Enter 鍵退出"
    exit 1
}

# 檢查 Python 是否可用
Write-Host "[1/5] 檢查 Python 環境..." -ForegroundColor Yellow
try {
    $pythonVersion = & $PythonPath --version 2>&1
    Write-Host "[OK] Python 版本: $pythonVersion" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] 找不到 Python 執行檔: $PythonPath" -ForegroundColor Red
    Write-Host "[INFO] 請確認 Python 已安裝並加入 PATH" -ForegroundColor Yellow
    Read-Host "按 Enter 鍵退出"
    exit 1
}

# 檢查腳本是否存在
Write-Host "`n[2/5] 檢查量子作業腳本..." -ForegroundColor Yellow
if (-not (Test-Path $ScriptPath)) {
    Write-Host "[ERROR] 找不到腳本: $ScriptPath" -ForegroundColor Red
    Read-Host "按 Enter 鍵退出"
    exit 1
}
Write-Host "[OK] 腳本路徑: $ScriptPath" -ForegroundColor Green

# 建立日誌目錄
$logDir = Split-Path -Parent $LogPath
if (-not (Test-Path $logDir)) {
    New-Item -ItemType Directory -Path $logDir -Force | Out-Null
    Write-Host "[OK] 建立日誌目錄: $logDir" -ForegroundColor Green
}

# 檢查是否已存在同名任務
Write-Host "`n[3/5] 檢查現有排程任務..." -ForegroundColor Yellow
$existingTask = Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue

if ($existingTask) {
    Write-Host "[WARNING] 發現現有任務: $TaskName" -ForegroundColor Yellow
    $response = Read-Host "是否要刪除並重新建立? (Y/N)"
    if ($response -eq 'Y' -or $response -eq 'y') {
        Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false
        Write-Host "[OK] 已刪除舊任務" -ForegroundColor Green
    } else {
        Write-Host "[INFO] 保留現有任務，退出腳本" -ForegroundColor Yellow
        Read-Host "按 Enter 鍵退出"
        exit 0
    }
}

# 建立排程任務
Write-Host "`n[4/5] 建立新的排程任務..." -ForegroundColor Yellow

# 設定觸發器 (每日 00:00 執行)
$trigger = New-ScheduledTaskTrigger -Daily -At $ExecutionTime

# 設定動作 (執行 Python 腳本)
$workingDir = Split-Path -Parent $ScriptPath
$action = New-ScheduledTaskAction -Execute $PythonPath `
    -Argument "`"$ScriptPath`" >> `"$LogPath`" 2>&1" `
    -WorkingDirectory $workingDir

# 設定主體 (使用當前使用者，不論是否登入都執行)
$principal = New-ScheduledTaskPrincipal -UserId $env:USERNAME -LogonType S4U -RunLevel Highest

# 設定其他選項
$settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable

# 註冊任務
try {
    Register-ScheduledTask -TaskName $TaskName -Trigger $trigger -Action $action -Principal $principal -Settings $settings -Description "Pandora Cyber AI/Quantum 每日零日攻擊偵測量子作業" | Out-Null
    Write-Host "[SUCCESS] 排程任務建立成功！" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] 建立任務失敗: $_" -ForegroundColor Red
    Read-Host "按 Enter 鍵退出"
    exit 1
}

# 顯示任務資訊
Write-Host "`n[5/5] 任務資訊:" -ForegroundColor Yellow
Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host "  任務名稱:     $TaskName" -ForegroundColor White
Write-Host "  執行時間:     每日 $ExecutionTime" -ForegroundColor White
Write-Host "  執行腳本:     $ScriptPath" -ForegroundColor White
Write-Host "  Python 路徑:  $PythonPath" -ForegroundColor White
Write-Host "  工作目錄:     $workingDir" -ForegroundColor White
Write-Host "  日誌檔案:     $LogPath" -ForegroundColor White
Write-Host "============================================================================" -ForegroundColor Cyan

# 詢問是否要立即測試執行
Write-Host "`n[選項] 是否要立即測試執行一次? (Y/N)" -ForegroundColor Yellow
$testRun = Read-Host

if ($testRun -eq 'Y' -or $testRun -eq 'y') {
    Write-Host "`n[測試] 開始執行量子作業..." -ForegroundColor Yellow
    Write-Host "============================================================================" -ForegroundColor Cyan
    
    try {
        # 切換到腳本目錄
        Push-Location $workingDir
        
        # 執行腳本
        & $PythonPath $ScriptPath
        
        Write-Host "============================================================================" -ForegroundColor Cyan
        Write-Host "[OK] 測試執行完成" -ForegroundColor Green
    } catch {
        Write-Host "[ERROR] 測試執行失敗: $_" -ForegroundColor Red
    } finally {
        Pop-Location
    }
}

# 顯示管理指令
Write-Host "`n============================================================================" -ForegroundColor Cyan
Write-Host "  排程任務設定完成！" -ForegroundColor Green
Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host "`n管理指令:" -ForegroundColor Yellow
Write-Host "  查看任務狀態:   Get-ScheduledTask -TaskName '$TaskName'" -ForegroundColor White
Write-Host "  手動執行一次:   Start-ScheduledTask -TaskName '$TaskName'" -ForegroundColor White
Write-Host "  停用任務:       Disable-ScheduledTask -TaskName '$TaskName'" -ForegroundColor White
Write-Host "  啟用任務:       Enable-ScheduledTask -TaskName '$TaskName'" -ForegroundColor White
Write-Host "  刪除任務:       Unregister-ScheduledTask -TaskName '$TaskName' -Confirm:`$false" -ForegroundColor White
Write-Host "  查看日誌:       Get-Content '$LogPath' -Tail 50" -ForegroundColor White
Write-Host ""

# 開啟工作排程器
Write-Host "[INFO] 按 Enter 鍵開啟工作排程器查看任務..." -ForegroundColor Yellow
Read-Host
Start-Process "taskschd.msc" "/s"

Write-Host "`n[完成] 排程設定全部完成！" -ForegroundColor Green
Write-Host ""

