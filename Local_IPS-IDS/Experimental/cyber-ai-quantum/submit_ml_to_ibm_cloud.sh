#!/bin/bash
# ============================================================================
# 提交 ML QASM 到 IBM Quantum Cloud（使用 ibm_cloud channel）
# ============================================================================

set -e

echo "=============================================="
echo "  ML QASM → IBM Quantum Cloud"
echo "=============================================="

# 檢查 Token
if [ -z "$IBM_QUANTUM_TOKEN" ]; then
    echo "❌ IBM_QUANTUM_TOKEN 未設定"
    exit 1
fi

echo "✅ IBM Token 已設定 (${#IBM_QUANTUM_TOKEN} 字元)"
echo ""

# 執行 Python
python3 << 'EOF'
import os
import sys
import numpy as np
import json
from datetime import datetime

print("="*60)
print("步驟 1: 生成 ML 量子電路")
print("="*60)

from generate_dynamic_qasm import create_zero_day_classifier_circuit

# 使用高風險特徵
features = np.array([0.06, 0.05, 0.2, 0.01, 0.033, 1.0])
qubits = 7
weights = np.random.rand(14)

print(f"\n特徵: {features}")
print("生成電路...")

circuit = create_zero_day_classifier_circuit(features, qubits, weights)

print(f"✅ 電路創建成功")
print(f"   Qubits: {circuit.num_qubits}")
print(f"   Depth: {circuit.depth()}")
print(f"   Gates: {circuit.size()}")

print("\n" + "="*60)
print("步驟 2: 連接 IBM Quantum (ibm_cloud channel)")
print("="*60)

try:
    from qiskit_ibm_runtime import QiskitRuntimeService, SamplerV2 as Sampler
    from qiskit.transpiler.preset_passmanagers import generate_preset_pass_manager
    
    token = os.getenv('IBM_QUANTUM_TOKEN')
    
    print("\n正在連接...")
    # 關鍵：使用 ibm_cloud channel（昨天成功的配置）
    service = QiskitRuntimeService(
        channel='ibm_cloud',
        token=token
    )
    
    print("✅ 連接成功！(使用 ibm_cloud channel)")
    
    # 列出後端
    backends = service.backends()
    print(f"\n可用後端: {len(backends)} 個")
    
    # 選擇最佳後端
    best_backend = None
    for backend in backends:
        if backend.status().operational:
            if 'simulator' in backend.name.lower():
                best_backend = backend
                break
    
    if not best_backend:
        best_backend = backends[0]
    
    print(f"✅ 選擇後端: {best_backend.name}")
    
    # 顯示後端資訊
    print(f"\n後端資訊:")
    print(f"   名稱: {best_backend.name}")
    if hasattr(best_backend, 'num_qubits'):
        print(f"   量子位元: {best_backend.num_qubits}")
    status = best_backend.status()
    if hasattr(status, 'pending_jobs'):
        print(f"   佇列: {status.pending_jobs} 個作業")
    
    print("\n" + "="*60)
    print("步驟 3: 轉譯電路")
    print("="*60)
    
    print(f"\n轉譯到 {best_backend.name} 硬體...")
    pm = generate_preset_pass_manager(backend=best_backend, optimization_level=1)
    transpiled = pm.run(circuit)
    
    print(f"✅ 轉譯完成")
    print(f"   原始深度: {circuit.depth()}")
    print(f"   轉譯深度: {transpiled.depth()}")
    print(f"   原始閘數: {circuit.size()}")
    print(f"   轉譯閘數: {transpiled.size()}")
    
    print("\n" + "="*60)
    print("步驟 4: 提交作業")
    print("="*60)
    
    print(f"\n提交到 {best_backend.name}...")
    
    sampler = Sampler(backend=best_backend)
    job = sampler.run([transpiled], shots=1024)
    
    job_id = job.job_id()
    print(f"\n✅ 作業已提交！")
    print(f"   作業 ID: {job_id}")
    print(f"   後端: {best_backend.name}")
    print(f"   狀態: {job.status()}")
    
    print("\n⏳ 等待量子執行...")
    result = job.result()
    
    print("✅ 執行完成！")
    
    print("\n" + "="*60)
    print("步驟 5: 分析結果")
    print("="*60)
    
    # 獲取結果
    pub_result = result[0]
    
    # V2 API 結果獲取
    counts = None
    for key in pub_result.data:
        if hasattr(pub_result.data[key], 'get_counts'):
            counts = pub_result.data[key].get_counts()
            break
    
    if counts:
        print(f"\n測量結果 (總計: {sum(counts.values())} shots):")
        
        # 顯示前 5 個最常見的結果
        sorted_counts = sorted(counts.items(), key=lambda x: x[1], reverse=True)
        for i, (state, count) in enumerate(sorted_counts[:5]):
            prob = count / sum(counts.values()) * 100
            print(f"   {i+1}. |{state}>: {count:4d} ({prob:5.1f}%)")
        
        # 分析 qubit[0]
        zero_count = sum(c for state, c in counts.items() if state[-1] == '0')
        one_count = sum(c for state, c in counts.items() if state[-1] == '1')
        total = zero_count + one_count
        
        zero_prob = zero_count / total if total > 0 else 0
        one_prob = one_count / total if total > 0 else 0
        
        print(f"\nqubit[0] 分類:")
        print(f"   |0> (正常): {zero_count:4d} ({zero_prob*100:5.1f}%)")
        print(f"   |1> (攻擊): {one_count:4d} ({one_prob*100:5.1f}%)")
        
        # 判定
        is_attack = one_prob > 0.5
        confidence = max(zero_prob, one_prob) * 100
        
        print("\n" + "="*60)
        if is_attack:
            print("🚨 判定: 零日攻擊偵測")
        else:
            print("✅ 判定: 正常行為")
        print(f"信心度: {confidence:.1f}%")
        print(f"後端: {best_backend.name}")
        print("="*60)
        
        # 保存結果
        os.makedirs("results", exist_ok=True)
        result_file = f"results/ibm_cloud_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
        
        with open(result_file, "w") as f:
            json.dump({
                "timestamp": datetime.now().isoformat(),
                "job_id": job_id,
                "backend": best_backend.name,
                "channel": "ibm_cloud",
                "circuit_info": {
                    "qubits": circuit.num_qubits,
                    "depth": circuit.depth(),
                    "gates": circuit.size(),
                    "transpiled_depth": transpiled.depth(),
                    "transpiled_gates": transpiled.size()
                },
                "measurements": {
                    "zero_count": int(zero_count),
                    "one_count": int(one_count),
                    "zero_prob": float(zero_prob),
                    "one_prob": float(one_prob)
                },
                "classification": {
                    "is_attack": bool(is_attack),
                    "confidence": float(confidence),
                    "threshold": 0.5
                }
            }, f, indent=2)
        
        print(f"\n💾 結果已保存: {result_file}")
        print("\n✅ IBM Quantum Cloud 提交成功！")
        
        sys.exit(0)
    else:
        print("⚠️  無法獲取測量結果")
        sys.exit(1)
        
except Exception as e:
    print(f"\n❌ 錯誤: {type(e).__name__}")
    print(f"訊息: {str(e)[:300]}")
    
    import traceback
    print("\n完整錯誤追蹤:")
    traceback.print_exc()
    
    sys.exit(1)
EOF

echo ""
echo "=============================================="
echo "  完成！"
echo "=============================================="

