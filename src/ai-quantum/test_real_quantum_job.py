#!/usr/bin/env python3
# -*- coding: utf-8 -*-
#  $env:IBM_QUANTUM_TOKEN="7PzS0AdaFB83OVn7WghkdLupzthCe1fTl3yI4Qrp7G6o"
"""
測試真實量子硬體執行
生成簡單的 QASM 電路並提交到 IBM Quantum
"""

import os
import sys
import time
from datetime import datetime

# 設置 UTF-8 輸出（Windows 兼容性）
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

from dotenv import load_dotenv

# 載入環境變數
load_dotenv()

print("=" * 60)
print("  IBM Quantum Real Hardware Test")
print("=" * 60)
print(f"\nTest Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")

token = os.getenv('IBM_QUANTUM_TOKEN')

if not token:
    print("[ERROR] IBM_QUANTUM_TOKEN not found")
    print("Please set: $env:IBM_QUANTUM_TOKEN='your_token'")
    sys.exit(1)

print(f"[OK] Token loaded ({len(token)} characters)")

try:
    from qiskit import QuantumCircuit
    from qiskit_ibm_runtime import QiskitRuntimeService, SamplerV2 as Sampler
    from qiskit.transpiler.preset_passmanagers import generate_preset_pass_manager
    
    print("[OK] Qiskit modules imported successfully")
    
    # ========== 步驟 1: 連接到 IBM Quantum ==========
    print("\n[STEP 1/6] Connecting to IBM Quantum...")
    try:
        service = QiskitRuntimeService(channel='ibm_cloud', token=token)
        print("[SUCCESS] Connected to IBM Quantum!")
    except Exception as e:
        print(f"[ERROR] Connection failed: {e}")
        sys.exit(1)
    
    # ========== 步驟 2: 選擇量子後端 ==========
    print("\n[STEP 2/6] Selecting quantum backend...")
    try:
        # 獲取最不忙碌的真實量子處理器
        backend = service.least_busy(operational=True, simulator=False)
        print(f"[OK] Selected backend: {backend.name}")
        print(f"     Qubits: {backend.num_qubits}")
        
        status = backend.status()
        print(f"     Queue: {status.pending_jobs} jobs")
        print(f"     Status: {status.status_msg}")
    except Exception as e:
        print(f"[ERROR] Cannot select backend: {e}")
        print("[INFO] Falling back to simulator...")
        backend = service.backend("ibmq_qasm_simulator")
        print(f"[OK] Using simulator: {backend.name}")
    
    # ========== 步驟 3: 創建量子電路 ==========
    print("\n[STEP 3/6] Creating quantum circuit...")
    
    # 創建一個簡單的 Bell State 電路（量子糾纏）
    qc = QuantumCircuit(2, 2)
    qc.h(0)           # Hadamard gate on qubit 0
    qc.cx(0, 1)       # CNOT gate (qubit 0 controls qubit 1)
    qc.measure([0, 1], [0, 1])  # Measure both qubits
    
    print("[OK] Quantum circuit created:")
    print(f"     Circuit: Bell State (Quantum Entanglement)")
    print(f"     Qubits: 2")
    print(f"     Gates: Hadamard + CNOT")
    print(f"     Measurements: 2")
    
    # 顯示 QASM 代碼
    print("\n[INFO] QASM Code:")
    print("-" * 50)
    qasm_str = qc.qasm()
    for line in qasm_str.split('\n')[:10]:  # 只顯示前 10 行
        print(f"  {line}")
    print("-" * 50)
    
    # ========== 步驟 4: 轉譯電路 ==========
    print("\n[STEP 4/6] Transpiling circuit for target backend...")
    try:
        pm = generate_preset_pass_manager(backend=backend, optimization_level=3)
        transpiled_qc = pm.run(qc)
        
        print(f"[OK] Circuit transpiled successfully")
        print(f"     Original depth: {qc.depth()}")
        print(f"     Transpiled depth: {transpiled_qc.depth()}")
        print(f"     Optimization: Level 3 (maximum)")
    except Exception as e:
        print(f"[WARNING] Transpilation failed: {e}")
        print("[INFO] Using original circuit")
        transpiled_qc = qc
    
    # ========== 步驟 5: 提交作業到量子硬體 ==========
    print("\n[STEP 5/6] Submitting job to IBM Quantum...")
    print(f"[INFO] Backend: {backend.name}")
    print(f"[INFO] Shots: 1024")
    print("[INFO] This will create a real quantum job!")
    print("[INFO] Job may take 1-30 minutes depending on queue...")
    
    try:
        # 創建 Sampler 並執行
        sampler = Sampler(backend)
        job = sampler.run([transpiled_qc], shots=1024)
        
        job_id = job.job_id()
        print(f"\n[SUCCESS] Job submitted!")
        print(f"[INFO] Job ID: {job_id}")
        print(f"[INFO] Status: {job.status()}")
        print(f"[INFO] Queue position: Checking...")
        
        # ========== 步驟 6: 等待結果 ==========
        print("\n[STEP 6/6] Waiting for results...")
        print("[INFO] Polling job status every 10 seconds...")
        print("[INFO] Press Ctrl+C to stop monitoring (job will continue)")
        
        start_time = time.time()
        poll_count = 0
        
        while True:
            status = job.status()
            elapsed = time.time() - start_time
            poll_count += 1
            
            print(f"\n[POLL {poll_count}] Time: {elapsed:.1f}s | Status: {status}")
            
            if status == 'DONE':
                print("\n[SUCCESS] Job completed!")
                
                # 獲取結果
                result = job.result()
                pub_result = result[0]
                counts = pub_result.data.meas.get_counts()
                
                print("\n[RESULTS] Measurement outcomes:")
                print("-" * 50)
                for bitstring, count in sorted(counts.items(), key=lambda x: x[1], reverse=True):
                    percentage = (count / 1024) * 100
                    bar = '#' * int(percentage / 2)
                    print(f"  |{bitstring}>: {count:4d} ({percentage:5.1f}%) {bar}")
                print("-" * 50)
                
                # 驗證 Bell State 特性
                print("\n[ANALYSIS] Bell State Verification:")
                count_00 = counts.get('00', 0)
                count_11 = counts.get('11', 0)
                count_01 = counts.get('01', 0)
                count_10 = counts.get('10', 0)
                
                entanglement_ratio = (count_00 + count_11) / 1024
                print(f"  |00> + |11>: {count_00 + count_11} ({entanglement_ratio*100:.1f}%)")
                print(f"  |01> + |10>: {count_01 + count_10} ({(1-entanglement_ratio)*100:.1f}%)")
                
                if entanglement_ratio > 0.85:
                    print("  [OK] Strong entanglement detected!")
                elif entanglement_ratio > 0.70:
                    print("  [OK] Moderate entanglement (noise present)")
                else:
                    print("  [WARNING] Weak entanglement (high noise)")
                
                print(f"\n[INFO] Total execution time: {elapsed:.1f} seconds")
                print(f"[INFO] Backend used: {backend.name}")
                print(f"[INFO] Job ID: {job_id}")
                
                # 保存結果
                result_file = f"results/quantum_job_{job_id}.txt"
                os.makedirs("results", exist_ok=True)
                with open(result_file, 'w') as f:
                    f.write(f"Job ID: {job_id}\n")
                    f.write(f"Backend: {backend.name}\n")
                    f.write(f"Time: {datetime.now()}\n")
                    f.write(f"Counts: {counts}\n")
                
                print(f"\n[OK] Results saved to: {result_file}")
                break
                
            elif status in ['CANCELLED', 'ERROR']:
                print(f"\n[ERROR] Job {status}")
                if hasattr(job, 'error_message'):
                    print(f"[ERROR] Message: {job.error_message()}")
                sys.exit(1)
                
            elif status == 'QUEUED':
                queue_info = job.queue_info()
                if queue_info:
                    print(f"       Queue position: {queue_info.position}")
                    print(f"       Estimated start: {queue_info.estimated_start_time}")
                
            # 每 10 秒輪詢一次
            time.sleep(10)
            
            # 超時保護（30 分鐘）
            if elapsed > 1800:
                print("\n[WARNING] Timeout after 30 minutes")
                print(f"[INFO] Job is still running: {job_id}")
                print("[INFO] Check status later with:")
                print(f"       job = service.job('{job_id}')")
                break
    
    except KeyboardInterrupt:
        print("\n\n[INFO] Monitoring stopped by user")
        print(f"[INFO] Job is still running: {job_id}")
        print("[INFO] Check status later with:")
        print(f"       from qiskit_ibm_runtime import QiskitRuntimeService")
        print(f"       service = QiskitRuntimeService(channel='ibm_cloud', token=token)")
        print(f"       job = service.job('{job_id}')")
        print(f"       print(job.status())")
        sys.exit(0)
        
    except Exception as e:
        print(f"\n[ERROR] Job execution failed: {e}")
        print(f"[DEBUG] Error type: {type(e).__name__}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

except ImportError as e:
    print(f"[ERROR] Missing module: {e}")
    print("\nInstall required packages:")
    print("  pip install qiskit qiskit-ibm-runtime")
    sys.exit(1)

except Exception as e:
    print(f"[ERROR] Unexpected error: {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

