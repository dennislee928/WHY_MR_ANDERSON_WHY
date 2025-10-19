#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
每日自動化量子作業
整合了 QASM 生成、提交、監控和結果獲取。

執行流程:
1. 載入訓練好的模型參數 (如果存在)
2. 生成動態 QASM 電路 (基於特徵和權重)
3. 連接 IBM Quantum
4. 轉譯並提交作業
5. 等待作業完成
6. 獲取並儲存結果
7. 自動分析結果並生成報告
"""

import os
import sys
import time
import json
from datetime import datetime
from dotenv import load_dotenv
import numpy as np

# Windows UTF-8 兼容性
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

from qiskit_ibm_runtime import QiskitRuntimeService, SamplerV2 as Sampler
from qiskit.transpiler.preset_passmanagers import generate_preset_pass_manager

# 導入我們建立的動態 QASM 生成器
from generate_dynamic_qasm import create_zero_day_classifier_circuit
from analyze_results import analyze_classification_results, save_analysis_report

# --- 載入環境 ---
load_dotenv()
token = os.getenv('IBM_QUANTUM_TOKEN')
if not token:
    print("[ERROR] IBM_QUANTUM_TOKEN not found.")
    print("[INFO] 請在 .env 檔案中設定 IBM_QUANTUM_TOKEN")
    sys.exit(1)

# --- 設定 ---
QUBITS = 7
SHOTS = 2048  # 建議使用更高的 shots 以獲得更穩定的統計結果
USE_SIMULATOR = os.getenv('USE_SIMULATOR', 'false').lower() == 'true'
MODEL_FILE = "quantum_classifier_model.json"
RESULTS_DIR = "results"
CLASSIFICATION_THRESHOLD = 0.5


def load_trained_model():
    """載入訓練好的模型參數"""
    try:
        with open(MODEL_FILE, 'r', encoding='utf-8') as f:
            model_info = json.load(f)
        trained_weights = np.array(model_info['trained_weights'])
        print(f"[OK] 成功載入訓練好的模型參數 (來自 {MODEL_FILE})")
        print(f"     模型訓練時間: {model_info.get('timestamp', 'N/A')}")
        print(f"     訓練準確率: {model_info.get('accuracy', 'N/A')}")
        return trained_weights, model_info
    except FileNotFoundError:
        print(f"[WARNING] 找不到模型檔案 '{MODEL_FILE}'")
        print("[INFO] 將使用預設角度進行預測")
        print("[HINT] 執行 train_quantum_classifier.py 來訓練模型")
        return None, None
    except Exception as e:
        print(f"[ERROR] 載入模型時發生錯誤: {e}")
        return None, None


def get_features_from_logs():
    """
    從 Windows Agent Logs 獲取特徵向量
    
    TODO: 實作真實的特徵提取邏輯
    目前使用模擬特徵
    """
    # !!! 重要 !!!
    # 在真實應用中，這裡應該:
    # 1. 從資料庫或檔案系統讀取最新的 Windows Log
    # 2. 呼叫 feature_extractor.py 進行特徵提取
    # 3. 返回標準化的特徵向量
    
    print("[INFO] 正在獲取日誌特徵...")
    print("[WARNING] 目前使用模擬特徵 (TODO: 整合真實 Windows Agent Logs)")
    
    # 模擬特徵
    features = np.random.rand(QUBITS - 1)
    print(f"[INFO] 模擬特徵向量: {np.round(features, 3)}")
    
    return features


def run_daily_job():
    """執行每日量子作業"""
    print("="*70)
    print("  🚀 開始每日零日攻擊偵測量子作業")
    print("="*70)
    print(f"時間: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"量子位元數: {QUBITS}")
    print(f"測量次數 (Shots): {SHOTS}")
    
    # 確保結果目錄存在
    if not os.path.exists(RESULTS_DIR):
        os.makedirs(RESULTS_DIR)
        print(f"[INFO] 創建結果目錄: {RESULTS_DIR}")
    
    # --- 1. 載入模型並獲取特徵 ---
    print(f"\n[1/7] 載入模型和獲取特徵...")
    trained_weights, model_info = load_trained_model()
    features = get_features_from_logs()
    
    # --- 2. 生成動態 QASM 電路 ---
    print(f"\n[2/7] 生成 {QUBITS}-qubit 分類電路...")
    try:
        circuit = create_zero_day_classifier_circuit(features, QUBITS, trained_weights)
        print(f"[OK] 電路生成完畢")
        print(f"     閘門數: {len(circuit.data)}")
        print(f"     電路深度: {circuit.depth()}")
    except Exception as e:
        print(f"[ERROR] 電路生成失敗: {e}")
        return None
    
    # --- 3. 連接 IBM Quantum ---
    print("\n[3/7] 連接 IBM Quantum...")
    try:
        service = QiskitRuntimeService(channel='ibm_cloud', token=token)
        
        if USE_SIMULATOR:
            print("[INFO] 使用模擬器進行測試...")
            backend = service.backend('ibmq_qasm_simulator')
        else:
            print("[INFO] 正在選擇可用的真實量子後端...")
            backend = service.least_busy(operational=True, simulator=False)
        
        print(f"[OK] 連接成功！選擇後端: {backend.name}")
        print(f"     後端狀態: {backend.status().status_msg}")
    except Exception as e:
        print(f"[ERROR] 連接失敗: {e}")
        print("[INFO] 請檢查網路連線和 IBM Quantum Token")
        return None

    # --- 4. 轉譯並提交作業 ---
    print("\n[4/7] 轉譯並提交作業...")
    try:
        pm = generate_preset_pass_manager(backend=backend, optimization_level=1)
        transpiled_qc = pm.run(circuit)
        
        print(f"[OK] 電路轉譯完成")
        print(f"     原始閘門數: {len(circuit.data)}")
        print(f"     轉譯後閘門數: {len(transpiled_qc.data)}")
        
        sampler = Sampler(backend)
        job = sampler.run([transpiled_qc], shots=SHOTS)
        job_id = job.job_id()
        print(f"[SUCCESS] 作業已提交！ Job ID: {job_id}")
        
        # 儲存作業資訊
        job_info_file = f"{RESULTS_DIR}/job_{job_id}_info.txt"
        with open(job_info_file, 'w', encoding='utf-8') as f:
            f.write(f"Job ID: {job_id}\n")
            f.write(f"Backend: {backend.name}\n")
            f.write(f"Submitted: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
            f.write(f"Shots: {SHOTS}\n")
            f.write(f"Qubits: {QUBITS}\n")
            f.write(f"Features: {features.tolist()}\n")
            if trained_weights is not None:
                f.write(f"Weights: {trained_weights.tolist()}\n")
        print(f"[INFO] 作業資訊已儲存至: {job_info_file}")
        
    except Exception as e:
        print(f"[ERROR] 提交失敗: {e}")
        import traceback
        traceback.print_exc()
        return None

    # --- 5. 等待作業完成 ---
    print("\n[5/7] 等待作業結果...")
    print("[INFO] 這可能需要幾分鐘到幾小時，取決於佇列狀況")
    
    max_wait_time = 3600  # 最多等待 1 小時
    check_interval = 30  # 每 30 秒檢查一次
    elapsed_time = 0
    
    try:
        while elapsed_time < max_wait_time:
            status = job.status()
            print(f"  [{datetime.now().strftime('%H:%M:%S')}] 狀態: {status} (已等待 {elapsed_time}s)")
            
            if status == 'DONE':
                print("\n[SUCCESS] 作業完成！")
                break
            elif status in ['ERROR', 'CANCELLED']:
                print(f"\n[ERROR] 作業失敗: {status}")
                return None
            
            time.sleep(check_interval)
            elapsed_time += check_interval
        else:
            print(f"\n[WARNING] 等待超時 ({max_wait_time}s)")
            print(f"[INFO] 作業仍在執行中，Job ID: {job_id}")
            print(f"[INFO] 稍後可使用 check_job_status.py 檢查結果")
            return {'job_id': job_id, 'status': 'PENDING'}
        
        result = job.result()
        
    except Exception as e:
        print(f"[ERROR] 等待結果時發生錯誤: {e}")
        return None

    # --- 6. 獲取並儲存結果 ---
    print("\n[6/7] 處理並儲存結果...")
    try:
        pub_result = result[0]
        
        # 使用修復後的結果獲取方式
        counts = {}
        for key in pub_result.data:
            if hasattr(pub_result.data[key], 'get_counts'):
                counts = pub_result.data[key].get_counts()
                break
        
        if not counts:
            raise RuntimeError("在結果中找不到計數數據。")

        # 將結果儲存為 JSON 檔案
        result_filename = f"{RESULTS_DIR}/result_{job_id}.json"
        result_data = {
            'job_id': job_id,
            'backend': backend.name,
            'shots': SHOTS,
            'qubits': QUBITS,
            'timestamp': datetime.now().isoformat(),
            'features': features.tolist(),
            'counts': {k: int(v) for k, v in counts.items()}  # 確保 value 是標準 int
        }
        
        if trained_weights is not None:
            result_data['weights'] = trained_weights.tolist()
            result_data['model_info'] = model_info
        
        with open(result_filename, 'w', encoding='utf-8') as f:
            json.dump(result_data, f, indent=2, ensure_ascii=False)

        print(f"[SUCCESS] 結果已儲存至: {result_filename}")
        print("\n測量結果預覽:")
        for bitstring, count in sorted(counts.items(), key=lambda x: x[1], reverse=True)[:5]:
            print(f"  {bitstring}: {count} 次 ({count/SHOTS*100:.1f}%)")
        
    except Exception as e:
        print(f"[ERROR] 處理結果失敗: {e}")
        import traceback
        traceback.print_exc()
        return None

    # --- 7. 自動分析結果 ---
    print("\n[7/7] 分析結果並生成報告...")
    try:
        analysis_result = analyze_classification_results(
            result_filename, 
            threshold=CLASSIFICATION_THRESHOLD,
            output_report=True
        )
        
        # 儲存分析報告
        report_filename = f"{RESULTS_DIR}/analysis_{job_id}.json"
        save_analysis_report(analysis_result, report_filename)
        
        print("\n" + "="*70)
        print("  🎉 每日作業全部完成！")
        print("="*70)
        
        return {
            'job_id': job_id,
            'status': 'SUCCESS',
            'result_file': result_filename,
            'analysis_file': report_filename,
            'is_zero_day': analysis_result.get('is_zero_day', False),
            'confidence': analysis_result.get('confidence', 0)
        }
        
    except Exception as e:
        print(f"[ERROR] 分析結果失敗: {e}")
        import traceback
        traceback.print_exc()
        return {
            'job_id': job_id,
            'status': 'COMPLETED_WITHOUT_ANALYSIS',
            'result_file': result_filename
        }


if __name__ == "__main__":
    try:
        result = run_daily_job()
        if result:
            print(f"\n[FINAL STATUS] {result['status']}")
            if result.get('is_zero_day'):
                print("⚠️  警告: 偵測到疑似零日攻擊！請立即檢查。")
        else:
            print("\n[FINAL STATUS] FAILED")
            sys.exit(1)
    except KeyboardInterrupt:
        print("\n\n[INFO] 使用者中斷執行")
        sys.exit(0)
    except Exception as e:
        print(f"\n[CRITICAL ERROR] 未預期的錯誤: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

