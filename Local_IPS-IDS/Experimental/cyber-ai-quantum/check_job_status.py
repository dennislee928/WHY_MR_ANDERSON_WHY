#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
檢查 IBM Quantum 作業狀態
"""

import os
import sys

# 設置 UTF-8 輸出（Windows 兼容性）
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

from dotenv import load_dotenv

load_dotenv()

if len(sys.argv) < 2:
    print("[ERROR] Usage: python check_job_status.py <job_id>")
    sys.exit(1)

job_id = sys.argv[1]
token = os.getenv('IBM_QUANTUM_TOKEN')

if not token:
    print("[ERROR] IBM_QUANTUM_TOKEN not found")
    sys.exit(1)

try:
    from qiskit_ibm_runtime import QiskitRuntimeService
    
    print(f"[INFO] Checking job: {job_id}")
    
    service = QiskitRuntimeService(channel='ibm_cloud', token=token)
    job = service.job(job_id)
    
    status = job.status()
    print(f"\n[STATUS] {status}")
    
    if status == 'DONE':
        print("\n[SUCCESS] Job completed!")
        try:
            result = job.result()
            pub_result = result[0]
            
            # 新版 API (V2) 推薦的作法 - 動態查找包含 get_counts() 的屬性
            counts = {}
            for key in pub_result.data:
                if hasattr(pub_result.data[key], 'get_counts'):
                    counts = pub_result.data[key].get_counts()
                    break  # 找到第一個就跳出
            
            if not counts:
                raise AttributeError("找不到任何包含 get_counts() 的測量結果屬性。")
            
            # 備用方案（保留向後兼容）
            if not counts:
                try:
                    counts = pub_result.data.meas.get_counts()
            except AttributeError:
                try:
                    counts = pub_result.data.c.get_counts()
                except:
                    # 獲取所有可用的數據鍵
                    data_keys = dir(pub_result.data)
                    print(f"[DEBUG] Available data keys: {[k for k in data_keys if not k.startswith('_')]}")
                    
                    # 嘗試獲取第一個可用的測量結果
                    for key in data_keys:
                        if not key.startswith('_') and hasattr(getattr(pub_result.data, key), 'get_counts'):
                            counts = getattr(pub_result.data, key).get_counts()
                            break
                    else:
                        counts = {}
        except Exception as e:
            print(f"[WARNING] Cannot parse results: {e}")
            counts = {}
        
        print("\n[RESULTS]")
        for bitstring, count in sorted(counts.items(), key=lambda x: x[1], reverse=True):
            percentage = (count / sum(counts.values())) * 100
            print(f"  |{bitstring}>: {count:4d} ({percentage:5.1f}%)")
    
    elif status == 'QUEUED':
        queue_info = job.queue_info()
        if queue_info:
            print(f"[INFO] Queue position: {queue_info.position}")
    
    elif status == 'RUNNING':
        print("[INFO] Job is currently running on quantum hardware")
    
    else:
        print(f"[INFO] Job status: {status}")

except Exception as e:
    print(f"[ERROR] {e}")
    sys.exit(1)

