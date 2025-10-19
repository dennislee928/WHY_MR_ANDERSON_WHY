#!/usr/bin/env python3
"""簡單的 IBM Quantum 連接測試"""
import os
import sys

print("="*60)
print("IBM Quantum Connection Test")
print("="*60)

# 檢查 Token
token = os.getenv('IBM_QUANTUM_TOKEN')
if not token:
    print("ERROR: IBM_QUANTUM_TOKEN not set")
    sys.exit(1)

print(f"Token configured: {len(token)} characters")

try:
    from qiskit_ibm_runtime import QiskitRuntimeService
    print("Qiskit Runtime imported successfully")
    
    print("\nConnecting to IBM Quantum...")
    service = QiskitRuntimeService(channel='ibm_quantum', token=token)
    
    print("SUCCESS: Connected to IBM Quantum!")
    
    backends = service.backends()
    print(f"\nAvailable backends: {len(backends)}")
    
    print("\nTop 10 backends:")
    for i, backend in enumerate(backends[:10]):
        print(f"  {i+1}. {backend.name}")
    
    # 測試提交簡單電路
    print("\n" + "="*60)
    print("Testing Simple Circuit Submission")
    print("="*60)
    
    from qiskit import QuantumCircuit
    from qiskit_ibm_runtime import Session, Sampler
    
    # 創建電路
    qc = QuantumCircuit(2, 2)
    qc.h(0)
    qc.cx(0, 1)
    qc.measure([0, 1], [0, 1])
    
    print(f"Circuit created: {qc.num_qubits} qubits, {qc.depth()} depth")
    
    # 選擇模擬器
    backend_name = backends[0].name
    for b in backends:
        if 'simulator' in b.name.lower():
            backend_name = b.name
            break
    
    print(f"\nUsing backend: {backend_name}")
    print("Submitting job...")
    
    with Session(service=service, backend=backend_name) as session:
        sampler = Sampler(session=session)
        job = sampler.run([qc], shots=100)
        
        print(f"Job submitted: {job.job_id()}")
        print(f"Status: {job.status()}")
        
        print("Waiting for result...")
        result = job.result()
        
        print("Job completed!")
        
        # 獲取結果
        pub_result = result[0]
        counts = None
        for key in pub_result.data:
            if hasattr(pub_result.data[key], 'get_counts'):
                counts = pub_result.data[key].get_counts()
                break
        
        if counts:
            print("\nMeasurement results:")
            for state, count in sorted(counts.items(), key=lambda x: x[1], reverse=True):
                print(f"  |{state}>: {count}")
            print("\nSUCCESS: IBM Quantum submission works!")
        else:
            print("WARNING: Could not get counts")
            
except Exception as e:
    print(f"\nERROR: {type(e).__name__}")
    print(f"Message: {str(e)[:200]}")
    import traceback
    traceback.print_exc()
    sys.exit(1)

print("\n" + "="*60)
print("Test Complete!")
print("="*60)


