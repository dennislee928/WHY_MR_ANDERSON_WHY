#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
批量上傳 QASM 到 IBM Quantum（使用 Batch 模式）
更高效的批量執行方式
"""

import os
import sys
import json
from datetime import datetime
from pathlib import Path

# 設置 UTF-8 輸出
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

from dotenv import load_dotenv

load_dotenv()

print("=" * 70)
print("  Batch QASM Upload to IBM Quantum")
print("=" * 70)
print(f"\nTime: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")

token = os.getenv('IBM_QUANTUM_TOKEN')

if not token:
    print("[ERROR] IBM_QUANTUM_TOKEN not found")
    sys.exit(1)

try:
    from qiskit import QuantumCircuit
    from qiskit_ibm_runtime import QiskitRuntimeService, Batch, SamplerV2 as Sampler
    from qiskit.transpiler.preset_passmanagers import generate_preset_pass_manager
    
    print("[OK] Modules imported\n")
    
    # ========== 載入 QASM 文件 ==========
    print("[STEP 1/5] Loading QASM files...")
    
    qasm_dir = Path("qasm_output")
    qasm_files = [f for f in qasm_dir.glob("*.qasm") if 'v3' not in f.name]
    
    print(f"[OK] Found {len(qasm_files)} files\n")
    
    circuits = []
    for qf in qasm_files:
        with open(qf, 'r') as f:
            qasm_code = f.read()
        
        qc = QuantumCircuit.from_qasm_str(qasm_code)
        circuits.append({
            'name': qf.stem,
            'circuit': qc
        })
        print(f"  [OK] {qf.name} - {qc.num_qubits} qubits, {qc.size()} gates")
    
    # ========== 連接 ==========
    print("\n[STEP 2/5] Connecting...")
    
    try:
        service = QiskitRuntimeService(channel='ibm_cloud', token=token)
        print("[OK] Connected")
    except:
        QiskitRuntimeService.save_account(channel='ibm_cloud', token=token, overwrite=True)
        service = QiskitRuntimeService()
        print("[OK] Connected via saved credentials")
    
    # ========== 選擇後端 ==========
    print("\n[STEP 3/5] Selecting backend...")
    
    # 使用模擬器進行快速測試
    try:
        backends = service.backends(simulator=True)
        if backends:
            backend = backends[0]
            print(f"[OK] Using simulator: {backend.name}")
        else:
            backend = service.least_busy(operational=True)
            print(f"[OK] Using: {backend.name}")
    except:
        backend = service.least_busy(operational=True)
        print(f"[OK] Using: {backend.name}")
    
    # ========== 轉譯 ==========
    print("\n[STEP 4/5] Transpiling circuits...")
    
    pm = generate_preset_pass_manager(backend=backend, optimization_level=1)
    transpiled = [pm.run(c['circuit']) for c in circuits]
    
    print(f"[OK] {len(transpiled)} circuits transpiled")
    
    # ========== 批量提交 ==========
    print("\n[STEP 5/5] Submitting batch job...")
    
    try:
        # 使用 Batch 模式（更高效）
        with Batch(backend=backend) as batch:
            sampler = Sampler()
            job = sampler.run(transpiled, shots=1024)
            
            job_id = job.job_id()
            
            print(f"\n[SUCCESS] Batch job submitted!")
            print(f"[INFO] Job ID: {job_id}")
            print(f"[INFO] Circuits: {len(transpiled)}")
            print(f"[INFO] Backend: {backend.name}")
            print(f"[INFO] Status: {job.status()}")
            
            # 保存批量作業信息
            batch_info = {
                'job_id': job_id,
                'backend': backend.name,
                'circuits': [c['name'] for c in circuits],
                'submitted': datetime.now().isoformat(),
                'shots': 1024
            }
            
            batch_file = f"results/batch_{job_id}.json"
            with open(batch_file, 'w') as f:
                json.dump(batch_info, f, indent=2)
            
            print(f"\n[OK] Batch info saved: {batch_file}")
            
            print("\n[INFO] To check results:")
            print(f"  python check_job_status.py {job_id}")
    
    except Exception as e:
        print(f"[ERROR] Batch submission failed: {e}")
        print("\n[INFO] Falling back to individual submissions...")
        
        # 逐個提交
        for circ_info in transpiled_circuits:
            try:
                sampler = Sampler(backend)
                job = sampler.run([circ_info['transpiled']], shots=1024)
                print(f"  [OK] {circ_info['name']}: {job.job_id()}")
            except Exception as e2:
                print(f"  [ERROR] {circ_info['name']}: {e2}")

except Exception as e:
    print(f"\n[ERROR] {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

