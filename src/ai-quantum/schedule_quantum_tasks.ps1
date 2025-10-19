# Pandora Quantum Analysis - Windows Task Scheduler Script
# 使用 Windows 工作排程器配置定期量子分析

Write-Host "=== Pandora Quantum Analysis Task Scheduler ===" -ForegroundColor Cyan
Write-Host ""

# 腳本路徑
$scriptPath = "C:\Users\dennis.lee\Documents\GitHub\Local_IPS-IDS\Experimental\cyber-ai-quantum"
$pythonExe = "python"

# ========== 創建每日量子分析任務 ==========
Write-Host "創建每日量子分析任務..." -ForegroundColor Yellow

$dailyAction = New-ScheduledTaskAction -Execute $pythonExe `
    -Argument "$scriptPath\scheduled_quantum_analysis.py daily" `
    -WorkingDirectory $scriptPath

$dailyTrigger = New-ScheduledTaskTrigger -Daily -At 2:00AM

$dailySettings = New-ScheduledTaskSettingsSet -ExecutionTimeLimit (New-TimeSpan -Hours 2)

Register-ScheduledTask `
    -TaskName "Pandora_Daily_Quantum_Analysis" `
    -Action $dailyAction `
    -Trigger $dailyTrigger `
    -Settings $dailySettings `
    -Description "每日量子威脅分析 - 重新評估高風險事件" `
    -Force

Write-Host "✅ 每日任務已創建 (每天 02:00)" -ForegroundColor Green

# ========== 創建每週量子訓練任務 ==========
Write-Host ""
Write-Host "創建每週量子訓練任務..." -ForegroundColor Yellow

$weeklyAction = New-ScheduledTaskAction -Execute $pythonExe `
    -Argument "$scriptPath\scheduled_quantum_analysis.py weekly" `
    -WorkingDirectory $scriptPath

$weeklyTrigger = New-ScheduledTaskTrigger -Weekly -DaysOfWeek Sunday -At 3:00AM

$weeklySettings = New-ScheduledTaskSettingsSet -ExecutionTimeLimit (New-TimeSpan -Hours 4)

Register-ScheduledTask `
    -TaskName "Pandora_Weekly_Quantum_Training" `
    -Action $weeklyAction `
    -Trigger $weeklyTrigger `
    -Settings $weeklySettings `
    -Description "每週量子模型重訓練" `
    -Force

Write-Host "✅ 每週任務已創建 (每週日 03:00)" -ForegroundColor Green

# ========== 創建每月批次分析任務 ==========
Write-Host ""
Write-Host "創建每月批次分析任務..." -ForegroundColor Yellow

$monthlyAction = New-ScheduledTaskAction -Execute $pythonExe `
    -Argument "$scriptPath\scheduled_quantum_analysis.py monthly" `
    -WorkingDirectory $scriptPath

# 每月1號
$monthlyTrigger = New-ScheduledTaskTrigger -Daily -At 4:00AM
$monthlyTrigger.DaysInterval = 30

$monthlySettings = New-ScheduledTaskSettingsSet -ExecutionTimeLimit (New-TimeSpan -Hours 8)

Register-ScheduledTask `
    -TaskName "Pandora_Monthly_Quantum_Batch" `
    -Action $monthlyAction `
    -Trigger $monthlyTrigger `
    -Settings $monthlySettings `
    -Description "每月量子批次深度分析" `
    -Force

Write-Host "✅ 每月任務已創建 (每月1號 04:00)" -ForegroundColor Green

# ========== 顯示已創建的任務 ==========
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "已創建的排程任務:" -ForegroundColor Cyan
Write-Host ""

Get-ScheduledTask | Where-Object {$_.TaskName -like "Pandora_*Quantum*"} | ForEach-Object {
    Write-Host "📅 $($_.TaskName)" -ForegroundColor Green
    Write-Host "   狀態: $($_.State)"
    Write-Host "   描述: $($_.Description)"
    Write-Host ""
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "管理命令:" -ForegroundColor Yellow
Write-Host "  查看任務: Get-ScheduledTask | Where-Object {`$_.TaskName -like 'Pandora_*'}"
Write-Host "  啟用任務: Enable-ScheduledTask -TaskName 'Pandora_Daily_Quantum_Analysis'"
Write-Host "  停用任務: Disable-ScheduledTask -TaskName 'Pandora_Daily_Quantum_Analysis'"
Write-Host "  刪除任務: Unregister-ScheduledTask -TaskName 'Pandora_Daily_Quantum_Analysis' -Confirm:`$false"
Write-Host ""
Write-Host "日誌位置: $scriptPath\analysis_results\"
Write-Host ""

