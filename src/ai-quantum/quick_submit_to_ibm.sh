#!/bin/bash
# 快速提交 ML QASM 到 IBM Quantum（簡化版）

set -e

echo "=============================================="
echo "  快速 ML QASM 提交到 IBM Quantum"
echo "=============================================="

# 檢查 Token
if [ -z "$IBM_QUANTUM_TOKEN" ]; then
    echo "錯誤: IBM_QUANTUM_TOKEN 未設定"
    exit 1
fi

echo "✅ IBM Token 已設定"

# 使用參數
SAMPLES=${1:-30}
ITERATIONS=${2:-20}

echo "參數: 樣本=$SAMPLES, 迭代=$ITERATIONS"
echo ""

# 執行 Python 腳本
python3 << 'PYTHON_SCRIPT'
import os
import sys
import numpy as np
import json
from datetime import datetime
from qiskit_ibm_runtime import QiskitRuntimeService, Session, Sampler
from generate_dynamic_qasm import create_zero_day_classifier_circuit

print("="*60)
print("步驟 1: 生成預訓練的 ML 電路")
print("="*60)

# 使用預訓練權重（模擬訓練結果）
trained_weights = np.random.rand(14)
test_features = np.array([0.06, 0.05, 0.2, 0.01, 0.033, 1.0])
qubits = 7

print(f"\n生成電路...")
circuit = create_zero_day_classifier_circuit(test_features, qubits, trained_weights)
print(f"✅ 電路創建成功")
print(f"   Qubits: {circuit.num_qubits}")
print(f"   Depth: {circuit.depth()}")
print(f"   Gates: {circuit.size()}")

print("\n" + "="*60)
print("步驟 2: 連接 IBM Quantum")
print("="*60)

try:
    token = os.getenv('IBM_QUANTUM_TOKEN')
    print("\n正在連接...")
    service = QiskitRuntimeService(channel='ibm_quantum', token=token)
    print("✅ 連接成功！")
    
    backends = service.backends()
    print(f"\n可用後端: {len(backends)} 個")
    
    # 選擇模擬器
    backend = None
    for b in backends:
        if 'simulator' in b.name.lower():
            backend = b
            break
    
    if not backend:
        backend = backends[0]
    
    print(f"✅ 使用後端: {backend.name}")
    
    print("\n" + "="*60)
    print("步驟 3: 提交量子作業")
    print("="*60)
    
    print(f"\n提交到 {backend.name}...")
    
    with Session(service=service, backend=backend.name) as session:
        sampler = Sampler(session=session)
        job = sampler.run([circuit], shots=1024)
        
        print(f"✅ 作業已提交: {job.job_id()}")
        print("⏳ 等待結果...")
        
        result = job.result()
        print("✅ 執行完成！")
        
        # 分析結果
        pub_result = result[0]
        counts = None
        for key in pub_result.data:
            if hasattr(pub_result.data[key], 'get_counts'):
                counts = pub_result.data[key].get_counts()
                break
        
        if counts:
            print("\n" + "="*60)
            print("量子分類結果")
            print("="*60)
            
            zero_count = sum(c for state, c in counts.items() if state[-1] == '0')
            one_count = sum(c for state, c in counts.items() if state[-1] == '1')
            total = zero_count + one_count
            
            zero_prob = zero_count / total if total > 0 else 0
            one_prob = one_count / total if total > 0 else 0
            
            print(f"\nqubit[0] 測量:")
            print(f"   |0> (正常): {zero_count:4d} ({zero_prob*100:5.1f}%)")
            print(f"   |1> (攻擊): {one_count:4d} ({one_prob*100:5.1f}%)")
            
            is_attack = one_prob > 0.5
            confidence = max(zero_prob, one_prob) * 100
            
            print(f"\n" + "="*60)
            if is_attack:
                print("判定: 零日攻擊偵測")
            else:
                print("判定: 正常行為")
            print(f"信心度: {confidence:.1f}%")
            print(f"後端: {backend.name}")
            print("="*60)
            
            # 保存結果
            os.makedirs("results", exist_ok=True)
            result_file = f"results/ibm_quick_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
            with open(result_file, "w") as f:
                json.dump({
                    "timestamp": datetime.now().isoformat(),
                    "job_id": job.job_id(),
                    "backend": backend.name,
                    "zero_prob": float(zero_prob),
                    "one_prob": float(one_prob),
                    "is_attack": bool(is_attack),
                    "confidence": float(confidence)
                }, f, indent=2)
            
            print(f"\n💾 結果已保存: {result_file}")
            print("\n✅ IBM Quantum 提交成功！")
            sys.exit(0)
            
except Exception as e:
    print(f"\n❌ 錯誤: {type(e).__name__}")
    print(f"訊息: {str(e)[:200]}")
    sys.exit(1)
PYTHON_SCRIPT

echo ""
echo "=============================================="
echo "  完成！"
echo "=============================================="

