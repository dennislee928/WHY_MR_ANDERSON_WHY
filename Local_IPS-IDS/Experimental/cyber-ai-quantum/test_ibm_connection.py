#!/usr/bin/env python3
"""
測試 IBM Quantum 連接和提交功能
"""
import os
import sys
import numpy as np
from datetime import datetime

try:
    from qiskit_ibm_runtime import QiskitRuntimeService, Session, Sampler
    from qiskit import QuantumCircuit
    print("✅ Qiskit Runtime 導入成功")
except ImportError as e:
    print(f"❌ 導入失敗: {e}")
    sys.exit(1)

def test_ibm_connection():
    """測試 IBM Quantum 連接"""
    print("\n" + "="*60)
    print("測試 IBM Quantum 連接")
    print("="*60)
    
    token = os.getenv('IBM_QUANTUM_TOKEN')
    if not token:
        print("❌ IBM_QUANTUM_TOKEN 環境變數未設定")
        return False
    
    print(f"✅ IBM Token 已設定 (長度: {len(token)})")
    
    try:
        print("\n正在連接到 IBM Quantum...")
        service = QiskitRuntimeService(
            channel='ibm_quantum',
            token=token
        )
        
        print("✅ 成功連接到 IBM Quantum！")
        
        # 列出可用後端
        backends = service.backends()
        print(f"\n可用後端數量: {len(backends)}")
        
        print("\n前 10 個可用後端:")
        for i, backend in enumerate(backends[:10]):
            status = backend.status()
            queue = status.pending_jobs if hasattr(status, 'pending_jobs') else 'N/A'
            print(f"  {i+1}. {backend.name:<30} | 佇列: {queue}")
        
        # 找到模擬器
        simulator = None
        for backend in backends:
            if 'simulator' in backend.name.lower():
                simulator = backend
                break
        
        if simulator:
            print(f"\n✅ 找到模擬器: {simulator.name}")
        else:
            print("\n⚠️  未找到模擬器，將使用第一個可用後端")
            simulator = backends[0] if backends else None
        
        return True, service, simulator
        
    except Exception as e:
        print(f"\n❌ 連接失敗: {type(e).__name__}")
        print(f"   錯誤訊息: {str(e)[:200]}")
        return False, None, None

def test_simple_circuit_submission(service, backend):
    """測試提交簡單電路到 IBM Quantum"""
    print("\n" + "="*60)
    print("測試提交量子電路")
    print("="*60)
    
    if not service or not backend:
        print("❌ 服務或後端未初始化")
        return False
    
    try:
        # 創建簡單的測試電路
        print("\n創建測試電路...")
        qc = QuantumCircuit(2, 2)
        qc.h(0)
        qc.cx(0, 1)
        qc.measure([0, 1], [0, 1])
        
        print(f"✅ 電路創建成功: {qc.num_qubits} qubits, {qc.depth()} depth")
        print(f"   電路內容:")
        print(qc.draw(output='text', initial_state=False))
        
        # 使用 Session 和 Sampler 提交作業
        print(f"\n正在提交到後端: {backend.name}...")
        
        with Session(service=service, backend=backend.name) as session:
            sampler = Sampler(session=session)
            
            print("⏳ 正在執行量子電路...")
            job = sampler.run([qc], shots=1000)
            
            print(f"✅ 作業已提交！")
            print(f"   作業 ID: {job.job_id()}")
            print(f"   狀態: {job.status()}")
            
            # 等待結果（最多 30 秒）
            print("\n⏳ 等待結果...")
            result = job.result()
            
            print(f"✅ 作業完成！")
            
            # 獲取結果
            pub_result = result[0]
            
            # 嘗試獲取計數（兼容 V2 API）
            counts = None
            for key in pub_result.data:
                if hasattr(pub_result.data[key], 'get_counts'):
                    counts = pub_result.data[key].get_counts()
                    break
            
            if counts:
                print(f"\n量子測量結果:")
                for state, count in sorted(counts.items(), key=lambda x: x[1], reverse=True):
                    prob = count / 1000 * 100
                    print(f"   |{state}⟩: {count:4d} ({prob:5.1f}%)")
                
                return True
            else:
                print("⚠️  無法獲取測量結果")
                return False
                
    except Exception as e:
        print(f"\n❌ 提交失敗: {type(e).__name__}")
        print(f"   錯誤訊息: {str(e)[:300]}")
        import traceback
        print("\n完整錯誤追蹤:")
        traceback.print_exc()
        return False

def test_zero_day_circuit(service, backend):
    """測試零日攻擊分類電路"""
    print("\n" + "="*60)
    print("測試零日攻擊分類電路")
    print("="*60)
    
    if not service or not backend:
        print("❌ 服務或後端未初始化")
        return False
    
    try:
        from generate_dynamic_qasm import create_zero_day_classifier_circuit
        
        # 生成模擬特徵（高風險情境）
        features = np.array([0.06, 0.05, 0.2, 0.01, 0.033, 1.0])  # 高風險
        qubits = 7
        weights = np.random.rand(14)
        
        print(f"\n特徵向量: {features}")
        print(f"量子位元數: {qubits}")
        
        # 創建電路
        circuit = create_zero_day_classifier_circuit(features, qubits, weights)
        print(f"✅ 分類電路創建成功: {circuit.num_qubits} qubits, {circuit.depth()} depth")
        
        # 提交到 IBM Quantum
        print(f"\n正在提交分類電路到: {backend.name}...")
        
        with Session(service=service, backend=backend.name) as session:
            sampler = Sampler(session=session)
            
            print("⏳ 正在執行量子分類...")
            job = sampler.run([circuit], shots=1024)
            
            print(f"✅ 分類作業已提交！")
            print(f"   作業 ID: {job.job_id()}")
            
            # 等待結果
            result = job.result()
            pub_result = result[0]
            
            # 獲取 qubit[0] 的測量結果
            counts = None
            for key in pub_result.data:
                if hasattr(pub_result.data[key], 'get_counts'):
                    counts = pub_result.data[key].get_counts()
                    break
            
            if counts:
                print(f"\n✅ 量子分類完成！")
                
                # 分析 qubit[0] (最右邊的位元)
                zero_count = sum(c for state, c in counts.items() if state[-1] == '0')
                one_count = sum(c for state, c in counts.items() if state[-1] == '1')
                total = zero_count + one_count
                
                zero_prob = zero_count / total if total > 0 else 0
                one_prob = one_count / total if total > 0 else 0
                
                print(f"\nqubit[0] 測量結果:")
                print(f"   |0⟩ (正常): {zero_count:4d} ({zero_prob*100:5.1f}%)")
                print(f"   |1⟩ (異常): {one_count:4d} ({one_prob*100:5.1f}%)")
                
                # 分類判定
                threshold = 0.5
                is_attack = one_prob > threshold
                
                print(f"\n分類結果: {'🚨 零日攻擊' if is_attack else '✅ 正常行為'}")
                print(f"信心度: {max(zero_prob, one_prob)*100:.1f}%")
                
                return True
            else:
                print("⚠️  無法獲取分類結果")
                return False

except Exception as e:
        print(f"\n❌ 分類測試失敗: {type(e).__name__}")
        print(f"   錯誤訊息: {str(e)[:300]}")
    import traceback
    traceback.print_exc()
        return False

def main():
    """主測試函數"""
    print("\n" + "="*60)
    print("IBM Quantum 功能完整測試")
    print("="*60)
    print(f"測試時間: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    
    # 測試 1: 連接測試
    success, service, backend = test_ibm_connection()
    if not success:
        print("\n❌ 連接測試失敗，無法繼續")
    sys.exit(1)
    
    # 測試 2: 簡單電路提交
    print("\n" + "─"*60)
    input("\n按 Enter 繼續測試簡單電路提交...")
    test_simple_circuit_submission(service, backend)
    
    # 測試 3: 零日攻擊分類電路
    print("\n" + "─"*60)
    input("\n按 Enter 繼續測試零日攻擊分類電路...")
    test_zero_day_circuit(service, backend)
    
    print("\n" + "="*60)
    print("測試完成！")
    print("="*60)

if __name__ == "__main__":
    main()
