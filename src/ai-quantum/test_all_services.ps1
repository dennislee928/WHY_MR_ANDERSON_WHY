# Pandora Cyber AI/Quantum - 完整服務測試腳本 (PowerShell)

Write-Host "=== Pandora Cyber AI/Quantum 服務測試 ===" -ForegroundColor Cyan
Write-Host ""

# 測試計數器
$script:Total = 0
$script:Passed = 0
$script:Failed = 0

# 測試函數
function Test-Endpoint {
    param(
        [string]$Name,
        [string]$Url,
        [string]$Method = "GET",
        [string]$Body = $null
    )
    
    $script:Total++
    Write-Host "測試 $Name... " -NoNewline
    
    try {
        $params = @{
            Uri = $Url
            Method = $Method
            UseBasicParsing = $true
            ErrorAction = 'Stop'
            TimeoutSec = 10
        }
        
        if ($Method -eq "POST" -and $Body) {
            $params.ContentType = "application/json"
            $params.Body = $Body
        }
        
        $response = Invoke-WebRequest @params
        
        if ($response.StatusCode -eq 200) {
            Write-Host "✅ 通過 (HTTP $($response.StatusCode))" -ForegroundColor Green
            $script:Passed++
        } else {
            Write-Host "❌ 失敗 (HTTP $($response.StatusCode))" -ForegroundColor Red
            $script:Failed++
        }
    }
    catch {
        Write-Host "❌ 失敗 ($($_.Exception.Message))" -ForegroundColor Red
        $script:Failed++
    }
}

Write-Host "--- 健康檢查 ---" -ForegroundColor Yellow
Test-Endpoint "Health Check" "http://localhost:8000/health"
Test-Endpoint "Root Endpoint" "http://localhost:8000/"
Write-Host ""

Write-Host "--- ML 威脅檢測 ---" -ForegroundColor Yellow
Test-Endpoint "ML Detect" "http://localhost:8000/api/v1/ml/detect" "POST" `
    '{"source_ip":"192.168.1.100","packets_per_second":1000,"syn_count":50}'
Test-Endpoint "ML Model Status" "http://localhost:8000/api/v1/ml/model/status"
Write-Host ""

Write-Host "--- 量子密碼學 ---" -ForegroundColor Yellow
Test-Endpoint "Quantum QKD" "http://localhost:8000/api/v1/quantum/qkd/generate" "POST" `
    '{"key_length":256}'
Test-Endpoint "Quantum Encrypt" "http://localhost:8000/api/v1/quantum/encrypt" "POST" `
    '{"message":"Test Message"}'
Test-Endpoint "Quantum Predict" "http://localhost:8000/api/v1/quantum/predict" "POST" `
    '{"historical_threats":[{"severity":0.8,"frequency":0.6,"impact":0.7}]}'
Write-Host ""

Write-Host "--- AI 治理 ---" -ForegroundColor Yellow
Test-Endpoint "Governance Integrity" "http://localhost:8000/api/v1/governance/integrity"
Test-Endpoint "Adversarial Detect" "http://localhost:8000/api/v1/governance/adversarial/detect" "POST" `
    '{"source_ip":"192.168.1.100","packets_per_second":100}'
Test-Endpoint "Governance Report" "http://localhost:8000/api/v1/governance/report"
Write-Host ""

Write-Host "--- 資料流監控 ---" -ForegroundColor Yellow
Test-Endpoint "DataFlow Stats" "http://localhost:8000/api/v1/dataflow/stats"
Test-Endpoint "DataFlow Anomalies" "http://localhost:8000/api/v1/dataflow/anomalies"
Test-Endpoint "DataFlow Baseline" "http://localhost:8000/api/v1/dataflow/baseline"
Write-Host ""

Write-Host "--- 系統狀態 ---" -ForegroundColor Yellow
Test-Endpoint "System Status" "http://localhost:8000/api/v1/status"
Write-Host ""

Write-Host "--- 相關服務檢查 ---" -ForegroundColor Yellow
Test-Endpoint "Axiom UI" "http://localhost:3001/"
Test-Endpoint "Axiom API" "http://localhost:3001/api/v1/status"
Test-Endpoint "RabbitMQ Mgmt" "http://localhost:15672/"
Test-Endpoint "Grafana" "http://localhost:3000/api/health"
Test-Endpoint "Prometheus" "http://localhost:9090/-/healthy"
Write-Host ""

# 總結
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "測試總結："
Write-Host "  總計: $Total"
Write-Host "  通過: $Passed" -ForegroundColor Green
Write-Host "  失敗: $Failed" -ForegroundColor Red
Write-Host ""

if ($Failed -eq 0) {
    Write-Host "🎉 所有測試通過！" -ForegroundColor Green
    exit 0
} else {
    $successRate = [math]::Round(($Passed / $Total) * 100, 1)
    Write-Host "⚠️  成功率: $successRate%" -ForegroundColor Yellow
    exit 1
}

