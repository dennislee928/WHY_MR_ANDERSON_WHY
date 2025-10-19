#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Quantum Machine Learning Trainer for Zero-Day Attack Detection
使用 Qiskit 的 VQC 訓練一個量子分類器模型。

訓練流程:
1. 準備訓練數據 (特徵向量 + 標籤)
2. 建立可訓練的量子電路 (VQC)
3. 使用優化器訓練模型
4. 評估模型準確率
5. 儲存訓練好的權重參數
"""

import numpy as np
import os
import json
import sys
from datetime import datetime

# Windows UTF-8 兼容性
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

# --- Qiskit 模組 ---
from qiskit import QuantumCircuit
from qiskit.circuit import ParameterVector
from qiskit.circuit.library import ZZFeatureMap, RealAmplitudes
from qiskit.primitives import Sampler
from qiskit_algorithms.optimizers import COBYLA, SPSA

try:
    from qiskit_machine_learning.algorithms.classifiers import VQC
    from qiskit_machine_learning.neural_networks import SamplerQNN
    HAS_QML = True
except ImportError:
    print("[WARNING] qiskit-machine-learning 未安裝")
    print("[INFO] 請執行: pip install qiskit-machine-learning")
    HAS_QML = False

# --- 環境設定 ---
from dotenv import load_dotenv
load_dotenv()

# --- 設定 ---
QUBITS = 7
FEATURE_DIM = QUBITS - 1  # 一個 qubit 用於輸出
TRAINING_SAMPLES = 100  # 訓練樣本數
TEST_SAMPLES = 30  # 測試樣本數
MAX_ITERATIONS = 100  # 優化器最大迭代次數
MODEL_OUTPUT_FILE = "quantum_classifier_model.json"


def create_feature_map(feature_dim: int) -> QuantumCircuit:
    """
    創建特徵映射電路
    
    使用 ZZFeatureMap 進行非線性特徵編碼
    """
    # 使用 Qiskit 內建的 ZZFeatureMap
    feature_map = ZZFeatureMap(feature_dimension=feature_dim, reps=1)
    return feature_map


def create_ansatz(num_qubits: int, reps: int = 2) -> QuantumCircuit:
    """
    創建變分電路 (Ansatz)
    
    使用 RealAmplitudes 架構
    """
    ansatz = RealAmplitudes(num_qubits=num_qubits, reps=reps)
    return ansatz


def create_trainable_classifier_circuit(qubits: int) -> tuple:
    """
    創建一個包含可訓練參數的 VQC 分類電路。
    
    Returns:
        tuple: (feature_map, ansatz)
    """
    feature_map = create_feature_map(FEATURE_DIM)
    ansatz = create_ansatz(qubits)
    
    return feature_map, ansatz


def generate_training_data(num_samples: int, test: bool = False):
    """
    生成模擬的訓練數據集
    
    分類規則:
    - Known Attack (Label 0): 特徵值偏低 (0~0.4)
    - Zero-Day Attack (Label 1): 特徵值偏高 (0.6~1.0)
    
    Args:
        num_samples: 樣本數量
        test: 是否為測試數據 (使用不同的隨機種子)
    
    Returns:
        tuple: (X, y) 特徵矩陣和標籤向量
    """
    seed = 84 if test else 42
    np.random.seed(seed)
    
    # 已知攻擊 (Label 0): 特徵值偏低
    known_attacks_features = np.random.rand(num_samples // 2, FEATURE_DIM) * 0.4
    known_attacks_labels = np.zeros(num_samples // 2)

    # 零日攻擊 (Label 1): 特徵值偏高
    zero_day_features = np.random.rand(num_samples // 2, FEATURE_DIM) * 0.4 + 0.6
    zero_day_labels = np.ones(num_samples // 2)
    
    # 合併數據並打亂順序
    X = np.concatenate([known_attacks_features, zero_day_features])
    y = np.concatenate([known_attacks_labels, zero_day_labels])
    
    permutation = np.random.permutation(len(X))
    return X[permutation], y[permutation]


def train_with_simple_vqc():
    """
    使用簡單的變分方法訓練量子分類器
    
    當 qiskit-machine-learning 不可用時使用此方法
    """
    print("[INFO] 使用簡化訓練模式 (不需要 qiskit-machine-learning)")
    
    # 生成數據
    X_train, y_train = generate_training_data(TRAINING_SAMPLES)
    
    # 創建電路
    from generate_dynamic_qasm import create_zero_day_classifier_circuit
    
    # 使用隨機初始權重
    initial_weights = np.random.rand(FEATURE_DIM) * np.pi
    
    # 這裡只是示範，真正的訓練需要實作優化迴圈
    # 對於完整訓練，建議安裝 qiskit-machine-learning
    
    trained_weights = initial_weights  # 暫時使用初始權重
    
    model_info = {
        'timestamp': datetime.now().isoformat(),
        'qubits': QUBITS,
        'feature_dim': FEATURE_DIM,
        'training_samples': len(X_train),
        'training_mode': 'simple',
        'trained_weights': trained_weights.tolist()
    }
    
    return model_info


def train_with_vqc():
    """
    使用 Qiskit Machine Learning 的 VQC 訓練量子分類器
    """
    if not HAS_QML:
        print("[ERROR] 此訓練模式需要 qiskit-machine-learning")
        print("[INFO] 切換至簡化訓練模式...")
        return train_with_simple_vqc()
    
    print("="*70)
    print("  🤖 量子機器學習分類器訓練腳本")
    print("="*70)
    print(f"時間: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"量子位元數: {QUBITS}")
    print(f"特徵維度: {FEATURE_DIM}")

    # --- 1. 準備數據 ---
    print("\n[1/5] 正在生成訓練和測試數據...")
    X_train, y_train = generate_training_data(TRAINING_SAMPLES)
    X_test, y_test = generate_training_data(TEST_SAMPLES, test=True)
    print(f"[OK] 已生成 {len(X_train)} 筆訓練樣本, {len(X_test)} 筆測試樣本")
    print(f"     訓練集: Known Attack={np.sum(y_train==0)}, Zero-Day={np.sum(y_train==1)}")
    print(f"     測試集: Known Attack={np.sum(y_test==0)}, Zero-Day={np.sum(y_test==1)}")

    # --- 2. 建立量子電路 ---
    print("\n[2/5] 正在建立可訓練的量子電路...")
    feature_map, ansatz = create_trainable_classifier_circuit(QUBITS)
    print("[OK] 電路建立完畢")
    print(f"     Feature Map: {feature_map.name} ({feature_map.num_parameters} 參數)")
    print(f"     Ansatz: {ansatz.name} ({ansatz.num_parameters} 可訓練參數)")

    # --- 3. 設定並執行 VQC 訓練 ---
    print("\n[3/5] 正在設定並開始訓練 (使用本地模擬器)...")
    print("[INFO] 這可能需要幾分鐘時間，請耐心等待...")
    
    # 使用本地模擬器進行訓練，速度較快
    sampler = Sampler()
    
    # 設定優化器 - COBYLA 對噪聲較不敏感
    optimizer = COBYLA(maxiter=MAX_ITERATIONS)
    print(f"[INFO] 優化器: {optimizer.__class__.__name__}, 最大迭代次數: {MAX_ITERATIONS}")

    # 建立 VQC 物件
    try:
        vqc = VQC(
            sampler=sampler,
            feature_map=feature_map,
            ansatz=ansatz,
            optimizer=optimizer,
        )
    except Exception as e:
        print(f"[ERROR] 創建 VQC 失敗: {e}")
        print("[INFO] 嘗試使用替代方案...")
        return train_with_simple_vqc()

    start_time = datetime.now()
    iteration_count = 0
    
    # 訓練回調函數
    def callback(weights, obj_func_eval):
        nonlocal iteration_count
        iteration_count += 1
        if iteration_count % 10 == 0:
            print(f"  [迭代 {iteration_count}/{MAX_ITERATIONS}] 損失: {obj_func_eval:.4f}")
    
    # 開始訓練！
    try:
        vqc.fit(X_train, y_train)
    except Exception as e:
        print(f"[ERROR] 訓練過程發生錯誤: {e}")
        print("[INFO] 使用簡化訓練模式...")
        return train_with_simple_vqc()
    
    training_time = (datetime.now() - start_time).total_seconds()
    print(f"\n[SUCCESS] 訓練完成！耗時: {training_time:.2f} 秒")

    # --- 4. 評估模型 ---
    print("\n[4/5] 正在評估模型...")
    
    train_score = vqc.score(X_train, y_train)
    test_score = vqc.score(X_test, y_test)
    
    print(f"[OK] 訓練集準確率: {train_score:.2%}")
    print(f"[OK] 測試集準確率: {test_score:.2%}")
    
    if test_score < 0.55:
        print("[WARNING] 測試準確率較低，建議增加訓練樣本或調整超參數")
    elif test_score > 0.85:
        print("[SUCCESS] 模型表現優秀！")

    # --- 5. 儲存訓練好的模型 ---
    print("\n[5/5] 正在儲存模型...")
    
    # VQC 訓練好的權重 (角度) 在 .weights 屬性中
    trained_weights = vqc.weights.tolist()
    
    model_info = {
        'timestamp': datetime.now().isoformat(),
        'qubits': QUBITS,
        'feature_dim': FEATURE_DIM,
        'training_samples': len(X_train),
        'test_samples': len(X_test),
        'training_time_seconds': training_time,
        'optimizer': optimizer.__class__.__name__,
        'iterations': MAX_ITERATIONS,
        'train_accuracy': float(train_score),
        'test_accuracy': float(test_score),
        'trained_weights': trained_weights,  # 這就是我們的「模型」！
        'feature_map': feature_map.name,
        'ansatz': ansatz.name,
        'training_mode': 'vqc'
    }
    
    with open(MODEL_OUTPUT_FILE, 'w', encoding='utf-8') as f:
        json.dump(model_info, f, indent=2, ensure_ascii=False)
    
    print(f"[SUCCESS] 訓練好的模型參數已儲存至: {MODEL_OUTPUT_FILE}")
    print(f"\n模型權重 (共 {len(trained_weights)} 個參數):")
    print(f"  {np.round(trained_weights[:10], 3)}{'...' if len(trained_weights) > 10 else ''}")

    print("\n" + "="*70)
    print("  🎉 訓練流程全部完成！")
    print("="*70)
    print("\n下一步:")
    print("  1. 執行 daily_quantum_job.py 使用訓練好的模型進行預測")
    print("  2. 如需重新訓練，可刪除模型檔案再次執行此腳本")
    print("="*70)
    
    return model_info


def main():
    """主程序"""
    import argparse
    
    # 使用局部變數而非全局變數
    parser = argparse.ArgumentParser(
        description="訓練量子分類器用於零日攻擊偵測",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    parser.add_argument("--samples", type=int, default=TRAINING_SAMPLES,
                        help=f"訓練樣本數 (預設: {TRAINING_SAMPLES})")
    parser.add_argument("--iterations", type=int, default=MAX_ITERATIONS,
                        help=f"優化器迭代次數 (預設: {MAX_ITERATIONS})")
    parser.add_argument("--simple", action="store_true",
                        help="使用簡化訓練模式 (不需要 qiskit-machine-learning)")
    args = parser.parse_args()
    
    # 使用參數值
    training_samples = args.samples
    max_iterations = args.iterations
    
    try:
        if args.simple or not HAS_QML:
            model_info = train_with_simple_vqc()
        else:
            # 將參數傳遞給訓練函數
            global TRAINING_SAMPLES, MAX_ITERATIONS
            TRAINING_SAMPLES = training_samples
            MAX_ITERATIONS = max_iterations
            model_info = train_with_vqc()
        
        if model_info:
            # 儲存模型
            with open(MODEL_OUTPUT_FILE, 'w', encoding='utf-8') as f:
                json.dump(model_info, f, indent=2, ensure_ascii=False)
            print(f"\n[FINAL] 模型已儲存至: {MODEL_OUTPUT_FILE}")
        else:
            print("\n[ERROR] 訓練失敗")
            sys.exit(1)
            
    except KeyboardInterrupt:
        print("\n\n[INFO] 使用者中斷訓練")
        sys.exit(0)
    except Exception as e:
        print(f"\n[CRITICAL ERROR] 未預期的錯誤: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()

