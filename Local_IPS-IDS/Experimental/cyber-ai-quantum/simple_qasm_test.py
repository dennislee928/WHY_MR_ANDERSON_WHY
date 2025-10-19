#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
簡單的 QASM 生成測試（不需要連接 IBM）
"""

import sys

# 設置 UTF-8 輸出
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

try:
    from qiskit import QuantumCircuit, qasm2, qasm3
    
    print("=" * 60)
    print("  QASM Generation Test (Offline)")
    print("=" * 60)
    
    # ========== 電路 1: Bell State ==========
    print("\n[CIRCUIT 1] Bell State (Quantum Entanglement)")
    bell = QuantumCircuit(2, 2)
    bell.h(0)
    bell.cx(0, 1)
    bell.measure([0, 1], [0, 1])
    
    print("\nQASM 2.0:")
    print("-" * 60)
    print(qasm2.dumps(bell))
    print("-" * 60)
    
    print("\nQASM 3.0:")
    print("-" * 60)
    print(qasm3.dumps(bell))
    print("-" * 60)
    
    # ========== 電路 2: Superposition ==========
    print("\n[CIRCUIT 2] Superposition (Equal States)")
    superpos = QuantumCircuit(3, 3)
    superpos.h(0)
    superpos.h(1)
    superpos.h(2)
    superpos.measure([0, 1, 2], [0, 1, 2])
    
    print("\nQASM 2.0:")
    print("-" * 60)
    print(qasm2.dumps(superpos))
    print("-" * 60)
    
    # ========== 電路 3: Phase Kickback ==========
    print("\n[CIRCUIT 3] Phase Kickback")
    phase = QuantumCircuit(2, 2)
    phase.h(0)
    phase.x(1)
    phase.h(1)
    phase.cz(0, 1)
    phase.h(0)
    phase.h(1)
    phase.measure([0, 1], [0, 1])
    
    print("\nQASM 2.0:")
    print("-" * 60)
    print(qasm2.dumps(phase))
    print("-" * 60)
    
    # 保存 QASM 文件
    print("\n[INFO] Saving QASM files...")
    
    with open("qasm_output/bell_state.qasm", 'w') as f:
        f.write(qasm2.dumps(bell))
    print("[OK] Saved: qasm_output/bell_state.qasm")
    
    with open("qasm_output/superposition.qasm", 'w') as f:
        f.write(qasm2.dumps(superpos))
    print("[OK] Saved: qasm_output/superposition.qasm")
    
    with open("qasm_output/phase_kickback.qasm", 'w') as f:
        f.write(qasm2.dumps(phase))
    print("[OK] Saved: qasm_output/phase_kickback.qasm")
    
    # 生成 OpenQASM 3.0 版本
    with open("qasm_output/bell_state_v3.qasm", 'w') as f:
        f.write(qasm3.dumps(bell))
    print("[OK] Saved: qasm_output/bell_state_v3.qasm")
    
    print("\n[SUCCESS] All QASM files generated!")
    print("\n[INFO] You can now:")
    print("  1. Upload these .qasm files to IBM Quantum Composer")
    print("  2. Run them on simulators or real hardware")
    print("  3. Use them in your quantum applications")
    
    print("\n[INFO] IBM Quantum Composer:")
    print("  https://quantum.ibm.com/composer")

except ImportError as e:
    print(f"[ERROR] Missing module: {e}")
    print("\nInstall: pip install qiskit")
    sys.exit(1)

except Exception as e:
    print(f"[ERROR] {e}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

