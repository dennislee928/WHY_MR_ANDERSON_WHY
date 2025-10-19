#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
測試 QASM 代碼上傳到 IBM Quantum
生成 QASM 並直接提交執行
"""

import os
import sys
from datetime import datetime

# 設置 UTF-8 輸出（Windows 兼容性）
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

from dotenv import load_dotenv

load_dotenv()

print("=" * 60)
print("  QASM Upload to IBM Quantum Test")
print("=" * 60)
print(f"\nTest Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")

token = os.getenv('IBM_QUANTUM_TOKEN')

if not token:
    print("[ERROR] IBM_QUANTUM_TOKEN not found")
    sys.exit(1)

try:
    from qiskit import QuantumCircuit, qasm2, qasm3
    from qiskit_ibm_runtime import QiskitRuntimeService, SamplerV2 as Sampler
    
    print("[OK] Modules imported\n")
    
    # ========== 創建 3 個測試電路 ==========
    
    # 電路 1: Bell State（量子糾纏）
    print("[CIRCUIT 1] Bell State (Quantum Entanglement)")
    bell_circuit = QuantumCircuit(2, 2, name="bell_state")
    bell_circuit.h(0)
    bell_circuit.cx(0, 1)
    bell_circuit.measure([0, 1], [0, 1])
    
    print("  Purpose: Test quantum entanglement")
    print("  Qubits: 2")
    print("  Expected: ~50% |00>, ~50% |11>")
    
    # 電路 2: GHZ State（3-qubit 糾纏）
    print("\n[CIRCUIT 2] GHZ State (3-qubit Entanglement)")
    ghz_circuit = QuantumCircuit(3, 3, name="ghz_state")
    ghz_circuit.h(0)
    ghz_circuit.cx(0, 1)
    ghz_circuit.cx(0, 2)
    ghz_circuit.measure([0, 1, 2], [0, 1, 2])
    
    print("  Purpose: Test 3-way entanglement")
    print("  Qubits: 3")
    print("  Expected: ~50% |000>, ~50% |111>")
    
    # 電路 3: Quantum Fourier Transform（量子傅立葉變換）
    print("\n[CIRCUIT 3] Quantum Fourier Transform (QFT)")
    qft_circuit = QuantumCircuit(3, 3, name="qft_3qubit")
    
    # QFT 實現
    qft_circuit.h(2)
    qft_circuit.cp(3.14159/2, 1, 2)
    qft_circuit.cp(3.14159/4, 0, 2)
    qft_circuit.h(1)
    qft_circuit.cp(3.14159/2, 0, 1)
    qft_circuit.h(0)
    
    # Swap qubits
    qft_circuit.swap(0, 2)
    qft_circuit.measure([0, 1, 2], [0, 1, 2])
    
    print("  Purpose: Test quantum phase estimation")
    print("  Qubits: 3")
    print("  Algorithm: QFT")
    
    # ========== 生成 QASM 代碼 ==========
    print("\n" + "=" * 60)
    print("  Generated QASM Code")
    print("=" * 60)
    
    circuits = [
        ("Bell State", bell_circuit),
        ("GHZ State", ghz_circuit),
        ("QFT", qft_circuit)
    ]
    
    for name, circuit in circuits:
        print(f"\n[{name}] QASM 2.0:")
        print("-" * 50)
        qasm_code = qasm2.dumps(circuit)
        for i, line in enumerate(qasm_code.split('\n')[:15], 1):
            print(f"  {i:2d} | {line}")
        if len(qasm_code.split('\n')) > 15:
            print(f"  ... ({len(qasm_code.split('\n')) - 15} more lines)")
        print("-" * 50)
    
    # ========== 連接並提交 ==========
    print("\n[STEP 3/6] Connecting to IBM Quantum...")
    service = QiskitRuntimeService(channel='ibm_cloud', token=token)
    print("[OK] Connected")
    
    print("\n[STEP 4/6] Selecting backend...")
    backend = service.least_busy(operational=True, simulator=False)
    print(f"[OK] Selected: {backend.name} ({backend.num_qubits} qubits)")
    
    # 選擇要執行的電路
    print("\n[STEP 5/6] Preparing to submit job...")
    print("[INFO] Submitting Bell State circuit (simplest test)")
    
    # 轉譯電路
    from qiskit.transpiler.preset_passmanagers import generate_preset_pass_manager
    pm = generate_preset_pass_manager(backend=backend, optimization_level=1)
    transpiled = pm.run(bell_circuit)
    
    print(f"[OK] Circuit transpiled")
    print(f"     Original gates: {bell_circuit.size()}")
    print(f"     Transpiled gates: {transpiled.size()}")
    
    # 提交作業
    print("\n[STEP 6/6] Submitting job to quantum hardware...")
    sampler = Sampler(backend)
    job = sampler.run([transpiled], shots=1024)
    
    job_id = job.job_id()
    
    print("\n" + "=" * 60)
    print("  JOB SUBMITTED SUCCESSFULLY!")
    print("=" * 60)
    print(f"\n[SUCCESS] Job ID: {job_id}")
    print(f"[INFO] Backend: {backend.name}")
    print(f"[INFO] Circuit: Bell State")
    print(f"[INFO] Shots: 1024")
    print(f"[INFO] Status: {job.status()}")
    
    # 保存作業信息
    job_file = f"results/quantum_job_{job_id}_info.txt"
    os.makedirs("results", exist_ok=True)
    
    with open(job_file, 'w', encoding='utf-8') as f:
        f.write(f"IBM Quantum Job Information\n")
        f.write(f"=" * 60 + "\n\n")
        f.write(f"Job ID: {job_id}\n")
        f.write(f"Backend: {backend.name}\n")
        f.write(f"Qubits: {backend.num_qubits}\n")
        f.write(f"Circuit: Bell State\n")
        f.write(f"Shots: 1024\n")
        f.write(f"Submitted: {datetime.now()}\n")
        f.write(f"Status: {job.status()}\n\n")
        f.write(f"QASM Code:\n")
        f.write("-" * 60 + "\n")
        f.write(qasm2.dumps(bell_circuit))
        f.write("\n" + "-" * 60 + "\n")
    
    print(f"\n[OK] Job info saved: {job_file}")
    
    # 監控選項
    print("\n[OPTIONS] What would you like to do?")
    print("  1. Wait for results (may take 1-30 minutes)")
    print("  2. Exit and check later")
    print("\n[INFO] To check job status later:")
    print(f"       python check_job_status.py {job_id}")
    print("\n[INFO] Or use:")
    print(f"       from qiskit_ibm_runtime import QiskitRuntimeService")
    print(f"       service = QiskitRuntimeService(channel='ibm_cloud', token=token)")
    print(f"       job = service.job('{job_id}')")
    print(f"       print(job.status())")
    print(f"       result = job.result()  # if DONE")
    
    # 簡單等待（可選）
    print("\n[INFO] Waiting 30 seconds for initial status...")
    time.sleep(30)
    
    final_status = job.status()
    print(f"[INFO] Current status: {final_status}")
    
    if final_status == 'DONE':
        print("[SUCCESS] Job completed quickly!")
        result = job.result()
        pub_result = result[0]
        counts = pub_result.data.meas.get_counts()
        print(f"\n[RESULTS] {counts}")
    else:
        print(f"[INFO] Job still {final_status}")
        print("[INFO] Check back later or run monitoring script")
    
    print("\n[SUCCESS] Test completed!")
    print(f"\n[SUMMARY]")
    print(f"  - QASM circuits generated: 3")
    print(f"  - Job submitted: YES")
    print(f"  - Job ID: {job_id}")
    print(f"  - Backend: {backend.name}")
    print(f"  - Status: {final_status}")

except ImportError as e:
    print(f"[ERROR] Missing module: {e}")
    print("\nInstall:")
    print("  pip install qiskit qiskit-ibm-runtime")
    sys.exit(1)

except Exception as e:
    print(f"\n[ERROR] {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

