#!/usr/bin/env python3
"""
每 10 分鐘自動提交 ML QASM 到 IBM Quantum
"""
import os
import sys
import time
import numpy as np
from datetime import datetime

print("="*70)
print("  自動提交 ML QASM 到 IBM Quantum - 每 10 分鐘")
print("="*70)
print(f"啟動時間: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
print("="*70)

# 檢查 Token
token = os.getenv('IBM_QUANTUM_TOKEN')
if not token:
    print("\n❌ IBM_QUANTUM_TOKEN 未設定")
    print("請執行: export IBM_QUANTUM_TOKEN='your_token'")
    sys.exit(1)

print(f"\n✅ IBM Token 已設定 (長度: {len(token)} 字元)")
print("ℹ️  循環間隔: 10 分鐘 (600 秒)")
print("ℹ️  按 Ctrl+C 停止\n")

def submit_ml_qasm():
    """提交 ML QASM 到 IBM Quantum"""
    try:
        from qiskit_ibm_runtime import QiskitRuntimeService, SamplerV2 as Sampler
        from qiskit.transpiler.preset_passmanagers import generate_preset_pass_manager
        from generate_dynamic_qasm import create_zero_day_classifier_circuit
        
        print("\n" + "─"*70)
        print(f"[{datetime.now().strftime('%H:%M:%S')}] 開始執行...")
        print("─"*70)
        
        # 1. 生成電路
        print("\n[1/5] 生成 ML 量子電路...")
        features = np.array([0.06, 0.05, 0.2, 0.01, 0.033, 1.0])
        qubits = 7
        weights = np.random.rand(14)
        
        circuit = create_zero_day_classifier_circuit(features, qubits, weights)
        print(f"✅ 電路創建: {circuit.num_qubits} qubits, {circuit.depth()} depth, {circuit.size()} gates")
        
        # 2. 連接 IBM Quantum
        print("\n[2/5] 連接 IBM Quantum (ibm_cloud channel)...")
        service = QiskitRuntimeService(channel='ibm_cloud', token=token)
        print("✅ 連接成功！")
        
        # 3. 選擇後端
        backends = service.backends()
        print(f"\n[3/5] 選擇後端 (可用: {len(backends)} 個)...")
        
        backend = None
        for b in backends:
            if 'simulator' in b.name.lower():
                backend = b
                break
        
        if not backend:
            backend = backends[0]
        
        print(f"✅ 使用: {backend.name}")
        
        # 4. 轉譯電路
        print(f"\n[4/5] 轉譯到 {backend.name}...")
        pm = generate_preset_pass_manager(backend=backend, optimization_level=1)
        transpiled = pm.run(circuit)
        print(f"✅ 轉譯完成: {transpiled.depth()} depth, {transpiled.size()} gates")
        
        # 5. 提交作業
        print(f"\n[5/5] 提交到 {backend.name}...")
        sampler = Sampler(mode=backend)
        job = sampler.run([transpiled], shots=1024)
        
        job_id = job.job_id()
        print(f"✅ 作業已提交: {job_id}")
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
            zero_count = sum(c for state, c in counts.items() if state[-1] == '0')
            one_count = sum(c for state, c in counts.items() if state[-1] == '1')
            total = zero_count + one_count
            
            zero_prob = zero_count / total if total > 0 else 0
            one_prob = one_count / total if total > 0 else 0
            
            print("\n" + "="*70)
            print(f"作業 {job_id} - 分類結果")
            print("="*70)
            print(f"  |0> (正常): {zero_count:4d} ({zero_prob*100:5.1f}%)")
            print(f"  |1> (攻擊): {one_count:4d} ({one_prob*100:5.1f}%)")
            
            is_attack = one_prob > 0.5
            confidence = max(zero_prob, one_prob) * 100
            
            verdict = "🚨 零日攻擊偵測" if is_attack else "✅ 正常行為"
            print(f"\n  判定: {verdict}")
            print(f"  信心度: {confidence:.1f}%")
            print(f"  後端: {backend.name}")
            print("="*70)
            
            # 保存結果
            import json
            os.makedirs("results", exist_ok=True)
            result_file = f"results/auto_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
            
            with open(result_file, "w") as f:
                json.dump({
                    "timestamp": datetime.now().isoformat(),
                    "job_id": job_id,
                    "backend": backend.name,
                    "zero_prob": float(zero_prob),
                    "one_prob": float(one_prob),
                    "is_attack": bool(is_attack),
                    "confidence": float(confidence)
                }, f, indent=2)
            
            print(f"\n💾 結果已保存: {result_file}")
            print(f"✅ 本次執行成功！")
            return True
        else:
            print("⚠️  無法獲取測量結果")
            return False
            
    except KeyboardInterrupt:
        raise  # 讓主循環處理 Ctrl+C
    except Exception as e:
        print(f"\n❌ 執行失敗: {type(e).__name__}")
        print(f"   {str(e)[:200]}")
        return False

# 主循環
execution_count = 0

try:
    while True:
        execution_count += 1
        
        print("\n" + "="*70)
        print(f"執行次數: {execution_count}")
        print(f"當前時間: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print("="*70)
        
        # 執行提交
        success = submit_ml_qasm()
        
        if success:
            print(f"\n✅ 第 {execution_count} 次執行成功")
        else:
            print(f"\n⚠️  第 {execution_count} 次執行失敗")
        
        # 計算下次執行時間
        next_run = datetime.now()
        next_run = next_run.replace(second=0, microsecond=0)
        next_run = next_run.replace(minute=(next_run.minute // 10 + 1) * 10 % 60)
        if next_run.minute == 0:
            next_run = next_run.replace(hour=next_run.hour + 1)
        
        print(f"\n⏰ 下次執行時間: {next_run.strftime('%Y-%m-%d %H:%M:%S')}")
        print("⏳ 等待 10 分鐘...")
        print("   (按 Ctrl+C 停止)")
        
        # 等待 10 分鐘
        time.sleep(600)
        
except KeyboardInterrupt:
    print("\n\n" + "="*70)
    print("  已停止自動執行")
    print("="*70)
    print(f"總執行次數: {execution_count}")
    print(f"結束時間: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print("\n再見！")
    sys.exit(0)
except Exception as e:
    print(f"\n❌ 嚴重錯誤: {type(e).__name__}")
    print(f"   {str(e)}")
    sys.exit(1)

