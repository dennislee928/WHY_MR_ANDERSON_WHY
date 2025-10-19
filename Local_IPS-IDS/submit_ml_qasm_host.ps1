# ============================================================================
# 在 Host 環境提交 ML QASM 到 IBM Quantum（模擬昨天的成功方式）
# ============================================================================

param(
    [Parameter(Mandatory=$false)]
    [string]$Token = $env:IBM_QUANTUM_TOKEN,
    
    [Parameter(Mandatory=$false)]
    [int]$Samples = 30,
    
    [Parameter(Mandatory=$false)]
    [switch]$UseSimulator = $false
)

Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host "  在 Host 環境提交 ML QASM 到 IBM Quantum" -ForegroundColor Cyan
Write-Host "============================================================================" -ForegroundColor Cyan
Write-Host ""

# 檢查 Token
if (-not $Token) {
    Write-Host "❌ IBM_QUANTUM_TOKEN 未設定" -ForegroundColor Red
    Write-Host ""
    Write-Host "請執行:" -ForegroundColor Yellow
    Write-Host '  $env:IBM_QUANTUM_TOKEN="your_token"' -ForegroundColor White
    Write-Host "  .\submit_ml_qasm_host.ps1" -ForegroundColor White
    exit 1
}

Write-Host "✅ IBM Token 已設定 (長度: $($Token.Length) 字元)" -ForegroundColor Green
Write-Host "ℹ️  訓練樣本: $Samples" -ForegroundColor Blue
Write-Host "ℹ️  使用模擬器: $UseSimulator" -ForegroundColor Blue
Write-Host ""

# 設定環境變數
$env:IBM_QUANTUM_TOKEN = $Token

# 切換到 cyber-ai-quantum 目錄
$quantumDir = "C:\Users\dennis.lee\Documents\GitHub\Local_IPS-IDS\Experimental\cyber-ai-quantum"
Set-Location $quantumDir

# 檢查 Python 和套件
Write-Host "[1/4] 檢查 Python 環境..." -ForegroundColor Yellow

try {
    $pythonVersion = python --version 2>&1
    Write-Host "  ✅ $pythonVersion" -ForegroundColor Green
} catch {
    Write-Host "  ❌ Python 未安裝" -ForegroundColor Red
    exit 1
}

# 檢查必要檔案
Write-Host ""
Write-Host "[2/4] 檢查必要檔案..." -ForegroundColor Yellow

$requiredFiles = @(
    "generate_dynamic_qasm.py",
    "requirements.txt"
)

foreach ($file in $requiredFiles) {
    if (Test-Path $file) {
        Write-Host "  ✅ $file" -ForegroundColor Green
    } else {
        Write-Host "  ❌ $file 未找到" -ForegroundColor Red
        exit 1
    }
}

# 執行提交
Write-Host ""
Write-Host "[3/4] 生成並提交 ML QASM..." -ForegroundColor Yellow
Write-Host ""

$pythonScript = @'
import os
import sys
import numpy as np
import json
from datetime import datetime

try:
    from qiskit_ibm_runtime import QiskitRuntimeService, SamplerV2 as Sampler
    from qiskit.transpiler.preset_passmanagers import generate_preset_pass_manager
    from generate_dynamic_qasm import create_zero_day_classifier_circuit
    
    print("="*60)
    print("生成 ML 量子電路")
    print("="*60)
    
    # 生成電路
    features = np.array([0.06, 0.05, 0.2, 0.01, 0.033, 1.0])
    qubits = 7
    weights = np.random.rand(14)
    
    circuit = create_zero_day_classifier_circuit(features, qubits, weights)
    print(f"\n✅ 電路創建: {circuit.num_qubits} qubits, {circuit.depth()} depth")
    
    print("\n" + "="*60)
    print("連接 IBM Quantum")
    print("="*60)
    
    token = os.getenv('IBM_QUANTUM_TOKEN')
    print(f"\nToken: {len(token)} 字元")
    print("正在連接 (ibm_cloud channel)...")
    
    service = QiskitRuntimeService(channel='ibm_cloud', token=token)
    print("✅ 連接成功！")
    
    backends = service.backends()
    print(f"可用後端: {len(backends)} 個")
    
    # 選擇後端
    backend = backends[0]
    for b in backends:
        if 'simulator' in b.name.lower():
            backend = b
            break
    
    print(f"\n使用後端: {backend.name}")
    
    # 轉譯
    print("\n轉譯電路...")
    pm = generate_preset_pass_manager(backend=backend, optimization_level=1)
    transpiled = pm.run(circuit)
    print(f"✅ 轉譯完成: {transpiled.depth()} depth, {transpiled.size()} gates")
    
    # 提交
    print("\n" + "="*60)
    print("提交量子作業")
    print("="*60)
    
    sampler = Sampler(backend=backend)
    job = sampler.run([transpiled], shots=1024)
    
    print(f"\n✅ 作業已提交: {job.job_id()}")
    print("⏳ 等待結果...")
    
    result = job.result()
    print("✅ 執行完成！")
    
    # 分析
    pub_result = result[0]
    counts = None
    for key in pub_result.data:
        if hasattr(pub_result.data[key], 'get_counts'):
            counts = pub_result.data[key].get_counts()
            break
    
    if counts:
        zero_count = sum(c for state, c in counts.items() if state[-1] == '0')
        one_count = sum(c for state, c in counts.items() if state[-1] == '1')
        total = zero_count + one_count
        
        zero_prob = zero_count / total if total > 0 else 0
        one_prob = one_count / total if total > 0 else 0
        
        print(f"\n" + "="*60)
        print("分類結果")
        print("="*60)
        print(f"|0> (正常): {zero_prob*100:.1f}%")
        print(f"|1> (攻擊): {one_prob*100:.1f}%")
        
        is_attack = one_prob > 0.5
        print(f"\n判定: {'🚨 攻擊' if is_attack else '✅ 正常'}")
        print(f"後端: {backend.name}")
        print("="*60)
        
        print("\n✅ 成功！")
        
except Exception as e:
    print(f"\n❌ 錯誤: {type(e).__name__}")
    print(f"{str(e)[:300]}")
    sys.exit(1)
'@

$pythonScript | python -

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "[4/4] ✅ 完成！" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "[4/4] ❌ 失敗" -ForegroundColor Red
    Write-Host ""
    Write-Host "可能原因:" -ForegroundColor Yellow
    Write-Host "  1. 需要安裝 Python 套件: pip install -r requirements.txt" -ForegroundColor White
    Write-Host "  2. 網路連接問題（防火牆/代理）" -ForegroundColor White
    Write-Host "  3. IBM Quantum Token 無效" -ForegroundColor White
}

Write-Host ""

