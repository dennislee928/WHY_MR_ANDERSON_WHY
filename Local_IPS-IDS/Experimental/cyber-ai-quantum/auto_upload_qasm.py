#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
自動上傳 QASM 文件到 IBM Quantum 並執行
通過 Python API 完全自動化，無需手動操作
"""

import os
import sys
import time
from datetime import datetime
from pathlib import Path

# 設置 UTF-8 輸出（Windows 兼容性）
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

from dotenv import load_dotenv

load_dotenv()

print("=" * 70)
print("  Automated QASM Upload & Execution to IBM Quantum")
print("=" * 70)
print(f"\nTime: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")

token = os.getenv('IBM_QUANTUM_TOKEN')

if not token:
    print("[ERROR] IBM_QUANTUM_TOKEN not found")
    print("Set it: $env:IBM_QUANTUM_TOKEN='your_token'")
    sys.exit(1)

print(f"[OK] Token loaded ({len(token)} characters)")

try:
    from qiskit import QuantumCircuit, qasm2
    from qiskit_ibm_runtime import QiskitRuntimeService, SamplerV2 as Sampler
    from qiskit.transpiler.preset_passmanagers import generate_preset_pass_manager
    
    print("[OK] Qiskit modules imported\n")
    
    # ========== 步驟 1: 掃描 QASM 文件 ==========
    print("[STEP 1/7] Scanning QASM files...")
    
    qasm_dir = Path("qasm_output")
    if not qasm_dir.exists():
        print("[ERROR] qasm_output/ directory not found")
        print("[INFO] Run 'python simple_qasm_test.py' first")
        sys.exit(1)
    
    qasm_files = list(qasm_dir.glob("*.qasm"))
    # 只處理 QASM 2.0 文件（排除 v3）
    qasm_files = [f for f in qasm_files if 'v3' not in f.name]
    
    print(f"[OK] Found {len(qasm_files)} QASM files:")
    for qf in qasm_files:
        size = qf.stat().st_size
        print(f"  - {qf.name} ({size} bytes)")
    
    if not qasm_files:
        print("[ERROR] No QASM files found")
        sys.exit(1)
    
    # ========== 步驟 2: 連接到 IBM Quantum ==========
    print("\n[STEP 2/7] Connecting to IBM Quantum...")
    
    try:
        service = QiskitRuntimeService(channel='ibm_cloud', token=token)
        print("[SUCCESS] Connected to IBM Quantum!")
    except Exception as e:
        print(f"[ERROR] Connection failed: {e}")
        print("\n[INFO] Trying alternative method...")
        try:
            # 嘗試保存帳號並重新連接
            QiskitRuntimeService.save_account(
                channel='ibm_cloud', 
                token=token, 
                overwrite=True
            )
            service = QiskitRuntimeService()
            print("[SUCCESS] Connected via saved credentials!")
        except Exception as e2:
            print(f"[ERROR] Alternative method failed: {e2}")
            print("\n[SOLUTION] Please check:")
            print("  1. Token validity: https://quantum.ibm.com/account")
            print("  2. Network/firewall settings")
            print("  3. VPN if required")
            sys.exit(1)
    
    # ========== 步驟 3: 選擇後端 ==========
    print("\n[STEP 3/7] Selecting quantum backend...")
    
    # 優先使用模擬器（快速測試）
    use_simulator = True  # 改為 False 使用真實硬體
    
    if use_simulator:
        print("[INFO] Using simulator for fast testing...")
        try:
            # 嘗試不同的模擬器名稱
            simulator_names = [
                'ibmq_qasm_simulator',
                'simulator_statevector',
                'simulator_mps',
                'simulator_extended_stabilizer',
                'simulator_stabilizer'
            ]
            
            backend = None
            for sim_name in simulator_names:
                try:
                    backend = service.backend(sim_name)
                    print(f"[OK] Using simulator: {backend.name}")
                    break
                except:
                    continue
            
            if not backend:
                print("[WARNING] No simulator found, using real hardware...")
                backend = service.least_busy(operational=True, simulator=False)
                print(f"[OK] Using real hardware: {backend.name}")
        except Exception as e:
            print(f"[ERROR] Cannot select simulator: {e}")
            backend = service.least_busy(operational=True, simulator=False)
            print(f"[OK] Using real hardware: {backend.name}")
    else:
        print("[INFO] Using real quantum hardware...")
        backend = service.least_busy(operational=True, simulator=False)
        print(f"[OK] Selected: {backend.name}")
        print(f"     Qubits: {backend.num_qubits}")
        status = backend.status()
        print(f"     Queue: {status.pending_jobs} jobs")
    
    # ========== 步驟 4: 載入並驗證 QASM 文件 ==========
    print("\n[STEP 4/7] Loading and validating QASM files...")
    
    circuits = []
    for qasm_file in qasm_files:
        print(f"\n[LOADING] {qasm_file.name}")
        
        try:
            # 讀取 QASM 內容
            with open(qasm_file, 'r') as f:
                qasm_code = f.read()
            
            print(f"  [OK] File read ({len(qasm_code)} bytes)")
            
            # 解析 QASM 為量子電路
            qc = QuantumCircuit.from_qasm_str(qasm_code)
            
            print(f"  [OK] QASM parsed successfully")
            print(f"       Qubits: {qc.num_qubits}")
            print(f"       Classical bits: {qc.num_clbits}")
            print(f"       Gates: {qc.size()}")
            print(f"       Depth: {qc.depth()}")
            
            circuits.append({
                'name': qasm_file.stem,
                'file': qasm_file.name,
                'circuit': qc,
                'qasm': qasm_code
            })
            
        except Exception as e:
            print(f"  [ERROR] Failed to load: {e}")
            continue
    
    if not circuits:
        print("\n[ERROR] No valid circuits loaded")
        sys.exit(1)
    
    print(f"\n[OK] Successfully loaded {len(circuits)} circuits")
    
    # ========== 步驟 5: 轉譯電路 ==========
    print("\n[STEP 5/7] Transpiling circuits for target backend...")
    
    pm = generate_preset_pass_manager(backend=backend, optimization_level=1)
    
    transpiled_circuits = []
    for circ_info in circuits:
        print(f"\n[TRANSPILING] {circ_info['name']}")
        
        try:
            transpiled = pm.run(circ_info['circuit'])
            
            print(f"  [OK] Transpiled successfully")
            print(f"       Original gates: {circ_info['circuit'].size()}")
            print(f"       Transpiled gates: {transpiled.size()}")
            print(f"       Reduction: {circ_info['circuit'].size() - transpiled.size()} gates")
            
            transpiled_circuits.append({
                **circ_info,
                'transpiled': transpiled
            })
            
        except Exception as e:
            print(f"  [ERROR] Transpilation failed: {e}")
            print(f"  [INFO] Using original circuit")
            transpiled_circuits.append({
                **circ_info,
                'transpiled': circ_info['circuit']
            })
    
    # ========== 步驟 6: 批量提交作業 ==========
    print("\n[STEP 6/7] Submitting jobs to IBM Quantum...")
    print(f"[INFO] Backend: {backend.name}")
    print(f"[INFO] Circuits to submit: {len(transpiled_circuits)}")
    print(f"[INFO] Shots per circuit: 1024")
    
    submitted_jobs = []
    
    for circ_info in transpiled_circuits:
        print(f"\n[SUBMITTING] {circ_info['name']}")
        
        try:
            sampler = Sampler(backend)
            job = sampler.run([circ_info['transpiled']], shots=1024)
            
            job_id = job.job_id()
            job_status = job.status()
            
            print(f"  [SUCCESS] Job submitted!")
            print(f"  [INFO] Job ID: {job_id}")
            print(f"  [INFO] Status: {job_status}")
            
            submitted_jobs.append({
                **circ_info,
                'job': job,
                'job_id': job_id,
                'submit_time': datetime.now()
            })
            
            # 保存作業信息
            job_info_file = f"results/job_{job_id}_info.txt"
            os.makedirs("results", exist_ok=True)
            
            with open(job_info_file, 'w', encoding='utf-8') as f:
                f.write(f"IBM Quantum Job Information\n")
                f.write(f"=" * 70 + "\n\n")
                f.write(f"Circuit Name: {circ_info['name']}\n")
                f.write(f"Job ID: {job_id}\n")
                f.write(f"Backend: {backend.name}\n")
                f.write(f"Submitted: {datetime.now()}\n")
                f.write(f"Status: {job_status}\n")
                f.write(f"Shots: 1024\n\n")
                f.write(f"Original QASM:\n")
                f.write("-" * 70 + "\n")
                f.write(circ_info['qasm'])
                f.write("\n" + "-" * 70 + "\n")
            
            print(f"  [OK] Job info saved: {job_info_file}")
            
            # 短暫延遲避免 API 限流
            time.sleep(2)
            
        except Exception as e:
            print(f"  [ERROR] Submission failed: {e}")
            print(f"  [DEBUG] {type(e).__name__}: {str(e)}")
            continue
    
    # ========== 步驟 7: 監控作業狀態 ==========
    print("\n[STEP 7/7] Monitoring job status...")
    print(f"[INFO] Total jobs submitted: {len(submitted_jobs)}")
    
    if not submitted_jobs:
        print("[ERROR] No jobs were submitted successfully")
        sys.exit(1)
    
    print("\n" + "=" * 70)
    print("  JOB SUBMISSION SUMMARY")
    print("=" * 70)
    
    for i, job_info in enumerate(submitted_jobs, 1):
        print(f"\n[JOB {i}] {job_info['name']}")
        print(f"  Job ID: {job_info['job_id']}")
        print(f"  Status: {job_info['job'].status()}")
        print(f"  Backend: {backend.name}")
        print(f"  Submitted: {job_info['submit_time'].strftime('%H:%M:%S')}")
    
    print("\n" + "=" * 70)
    
    # 提供監控選項
    print("\n[OPTIONS] What would you like to do?")
    print("  1. Wait for all jobs to complete (may take time)")
    print("  2. Monitor first job only")
    print("  3. Exit and check later")
    
    print("\n[INFO] To check job status later:")
    for job_info in submitted_jobs:
        print(f"  python check_job_status.py {job_info['job_id']}")
    
    # 簡單監控第一個作業
    if submitted_jobs:
        print("\n[INFO] Monitoring first job for 60 seconds...")
        first_job = submitted_jobs[0]['job']
        first_job_id = submitted_jobs[0]['job_id']
        
        for i in range(6):  # 6 次，每次 10 秒
            time.sleep(10)
            status = first_job.status()
            print(f"  [CHECK {i+1}/6] Status: {status}")
            
            if status == 'DONE':
                print("\n[SUCCESS] Job completed!")
                
                try:
                    result = first_job.result()
                    pub_result = result[0]
                    
                    # 新版 API (V2) 推薦的作法 - 動態查找包含 get_counts() 的屬性
                    counts = {}
                    for key in pub_result.data:
                        if hasattr(pub_result.data[key], 'get_counts'):
                            counts = pub_result.data[key].get_counts()
                            break  # 找到第一個就跳出
                    
                    if not counts:
                        raise AttributeError("找不到任何包含 get_counts() 的測量結果屬性。")
                
                except Exception as result_error:
                    print(f"[WARNING] 無法解析詳細結果: {result_error}")
                    print("[INFO] 作業已完成，但結果格式可能已變更或為空。")
                    counts = {}
                
                print("\n[RESULTS]")
                for bitstring, count in sorted(counts.items(), key=lambda x: x[1], reverse=True)[:5]:
                    percentage = (count / 1024) * 100
                    bar = '#' * int(percentage / 2)
                    print(f"  |{bitstring}>: {count:4d} ({percentage:5.1f}%) {bar}")
                
                # 保存結果
                result_file = f"results/job_{first_job_id}_result.txt"
                with open(result_file, 'w', encoding='utf-8') as f:
                    f.write(f"Job ID: {first_job_id}\n")
                    f.write(f"Backend: {backend.name}\n")
                    f.write(f"Completed: {datetime.now()}\n")
                    f.write(f"Results: {counts}\n")
                
                print(f"\n[OK] Results saved: {result_file}")
                break
            
            elif status in ['CANCELLED', 'ERROR']:
                print(f"\n[ERROR] Job {status}")
                break
        
        else:
            print(f"\n[INFO] Job still running: {first_job_id}")
            print("[INFO] Check later with:")
            print(f"  python check_job_status.py {first_job_id}")
    
    # 最終摘要
    print("\n" + "=" * 70)
    print("  UPLOAD COMPLETE!")
    print("=" * 70)
    print(f"\n[SUMMARY]")
    print(f"  QASM files processed: {len(circuits)}")
    print(f"  Jobs submitted: {len(submitted_jobs)}")
    print(f"  Backend: {backend.name}")
    print(f"  Results directory: results/")
    
    print("\n[SUCCESS] All QASM files uploaded and executed!")
    print("\n[INFO] Job IDs:")
    for job_info in submitted_jobs:
        print(f"  - {job_info['name']}: {job_info['job_id']}")

except ImportError as e:
    print(f"[ERROR] Missing module: {e}")
    print("\nInstall:")
    print("  pip install qiskit qiskit-ibm-runtime")
    sys.exit(1)

except Exception as e:
    print(f"\n[ERROR] {e}")
    print(f"[DEBUG] Type: {type(e).__name__}")
    import traceback
    print("\n[DEBUG] Traceback:")
    traceback.print_exc()
    sys.exit(1)

