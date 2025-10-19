#!/usr/bin/env python3
"""
在 Host 環境測試 IBM Quantum 連接（模擬昨天的成功方式）
"""
import os
import sys
import numpy as np
from datetime import datetime

print("="*60)
print("Host Environment IBM Quantum Test")
print("="*60)
print(f"Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")

# 檢查 Token
token = os.getenv('IBM_QUANTUM_TOKEN')
if not token:
    print("ERROR: IBM_QUANTUM_TOKEN not set")
    print("Run: export IBM_QUANTUM_TOKEN='your_token'")
    sys.exit(1)

print(f"Token configured: {len(token)} characters")

try:
    print("\nImporting Qiskit...")
    from qiskit_ibm_runtime import QiskitRuntimeService, SamplerV2 as Sampler
    from qiskit.transpiler.preset_passmanagers import generate_preset_pass_manager
    from generate_dynamic_qasm import create_zero_day_classifier_circuit
    print("✅ Imports successful")
    
    print("\n" + "="*60)
    print("Step 1: Generate ML Circuit")
    print("="*60)
    
    features = np.array([0.06, 0.05, 0.2, 0.01, 0.033, 1.0])
    qubits = 7
    weights = np.random.rand(14)
    
    circuit = create_zero_day_classifier_circuit(features, qubits, weights)
    print(f"\n✅ Circuit created")
    print(f"   Qubits: {circuit.num_qubits}")
    print(f"   Depth: {circuit.depth()}")
    print(f"   Gates: {circuit.size()}")
    
    print("\n" + "="*60)
    print("Step 2: Connect to IBM Quantum (ibm_cloud)")
    print("="*60)
    
    print("\nConnecting using ibm_cloud channel...")
    service = QiskitRuntimeService(channel='ibm_cloud', token=token)
    print("✅ Connected successfully!")
    
    backends = service.backends()
    print(f"\nAvailable backends: {len(backends)}")
    
    # 選擇後端
    backend = None
    for b in backends[:10]:
        print(f"  - {b.name}")
        if 'simulator' in b.name.lower():
            backend = b
    
    if not backend:
        backend = backends[0]
    
    print(f"\n✅ Selected backend: {backend.name}")
    
    print("\n" + "="*60)
    print("Step 3: Transpile Circuit")
    print("="*60)
    
    pm = generate_preset_pass_manager(backend=backend, optimization_level=1)
    transpiled = pm.run(circuit)
    print(f"✅ Transpiled: {transpiled.depth()} depth, {transpiled.size()} gates")
    
    print("\n" + "="*60)
    print("Step 4: Submit Job")
    print("="*60)
    
    print(f"\nSubmitting to {backend.name}...")
    sampler = Sampler(mode=backend)
    job = sampler.run([transpiled], shots=1024)
    
    print(f"✅ Job submitted: {job.job_id()}")
    print("Waiting for result...")
    
    result = job.result()
    print("✅ Job completed!")
    
    # 分析
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
        
        print("\n" + "="*60)
        print("Classification Result")
        print("="*60)
        print(f"|0> (Normal): {zero_prob*100:.1f}%")
        print(f"|1> (Attack): {one_prob*100:.1f}%")
        
        is_attack = one_prob > 0.5
        print(f"\nVerdict: {'ATTACK' if is_attack else 'NORMAL'}")
        print(f"Backend: {backend.name}")
        print("="*60)
        
        print("\n✅ SUCCESS! ML QASM submitted to IBM Quantum!")
        
except Exception as e:
    print(f"\n❌ ERROR: {type(e).__name__}")
    print(f"Message: {str(e)[:300]}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

