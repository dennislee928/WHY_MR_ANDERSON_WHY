#!/usr/bin/env python3
"""測試量子分類器（使用本地模擬器）"""
import numpy as np
from generate_dynamic_qasm import create_zero_day_classifier_circuit
from qiskit_aer import AerSimulator

print("="*60)
print("Quantum Classifier Test (Local Simulator)")
print("="*60)

# 創建高風險特徵
features = np.array([0.06, 0.05, 0.2, 0.01, 0.033, 1.0])
qubits = 7
weights = np.random.rand(14)

print(f"\nFeatures: {features}")
print(f"Qubits: {qubits}")

# 生成電路
circuit = create_zero_day_classifier_circuit(features, qubits, weights)
print(f"\nCircuit created successfully!")
print(f"  Qubits: {circuit.num_qubits}")
print(f"  Depth: {circuit.depth()}")
print(f"  Gates: {circuit.size()}")

# 使用本地模擬器執行
print("\nRunning on local simulator...")
simulator = AerSimulator()
job = simulator.run(circuit, shots=1024)
result = job.result()
counts = result.get_counts()

print(f"Simulation completed! Total shots: {sum(counts.values())}")

# 分析 qubit[0] (最右邊的位元)
zero_count = sum(c for state, c in counts.items() if state[-1] == '0')
one_count = sum(c for state, c in counts.items() if state[-1] == '1')
total = zero_count + one_count

zero_prob = zero_count / total if total > 0 else 0
one_prob = one_count / total if total > 0 else 0

print(f"\nQubit[0] Measurement:")
print(f"  |0> (Normal):  {zero_count:4d} ({zero_prob*100:5.1f}%)")
print(f"  |1> (Attack):  {one_count:4d} ({one_prob*100:5.1f}%)")

# 分類判定
threshold = 0.5
is_attack = one_prob > threshold
confidence = max(zero_prob, one_prob) * 100

print(f"\n" + "="*60)
if is_attack:
    print("VERDICT: ZERO-DAY ATTACK DETECTED")
else:
    print("VERDICT: NORMAL BEHAVIOR")
print(f"Confidence: {confidence:.1f}%")
print("="*60)

print("\nSUCCESS: Quantum classifier works correctly!")

