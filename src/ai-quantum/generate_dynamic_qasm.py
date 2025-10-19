#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
動態 QASM 電路生成器
用於從輸入特徵生成零日攻擊分類電路

根據 VQC (Variational Quantum Circuit) 架構設計
使用 qubit[0] 作為輸出分類位元
"""

import numpy as np
from qiskit import QuantumCircuit, qasm2
import argparse
from datetime import datetime
import sys
import os

# Windows UTF-8 兼容性
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')


def create_zero_day_classifier_circuit(features: np.ndarray, qubits: int, weights: np.ndarray = None) -> QuantumCircuit:
    """
    根據輸入特徵創建一個 VQC (Variational Quantum Circuit) 分類電路。

    Args:
        features (np.ndarray): 從日誌數據提取的特徵向量，值應在 [0, 1] 之間。
        qubits (int): 電路中使用的量子位元數量。
        weights (np.ndarray, optional): 訓練好的權重參數。如果為 None，使用預設角度。

    Returns:
        QuantumCircuit: 用於分類的量子電路。
    """
    # 確保特徵數量與量子位元數量匹配（特徵不足時補 0）
    if len(features) < qubits:
        features = np.pad(features, (0, qubits - len(features)))
    elif len(features) > qubits:
        features = features[:qubits]
    
    # 我們將使用 (qubits-1) 個量子位元來編碼特徵，最後 1 個 (qubit[0]) 作為輸出
    feature_qubits = qubits - 1
    output_qubit = 0  # 根據設計，qubit[0] 為輸出
    
    # 如果沒有提供權重，使用預設值
    if weights is None:
        weights = np.full(feature_qubits, np.pi / 4)
    
    # 創建量子電路，包含 n 個量子位元和 1 個古典位元 (只測量 qubit[0])
    qc = QuantumCircuit(qubits, 1, name="zero_day_classifier")

    # --- 1. 特徵編碼層 (Feature Encoding) ---
    # 使用 RX 門將古典特徵數據編碼到量子位元上
    # 這裡我們將特徵映射到 qubit[1] 到 qubit[n-1]
    for i in range(feature_qubits):
        # features[i] * np.pi 將特徵值 (0~1) 映射到旋轉角度 (0~pi)
        qc.rx(features[i] * np.pi, i + 1)
    
    qc.barrier()  # 分隔層，方便視覺化

    # --- 2. 糾纏層 (Entanglement Layer) ---
    # 使用 CNOT 門在相鄰的特徵量子位元之間創建糾纏
    for i in range(feature_qubits - 1):
        qc.cx(i + 1, i + 2)
    
    qc.barrier()

    # --- 3. 可訓練/決策層 (Variational/Decision Layer) ---
    # 讓特徵量子位元的資訊影響輸出的 qubit[0]
    # 這裡使用受控的 RY 門
    for i in range(feature_qubits):
        # 使用訓練好的角度或預設角度
        trainable_angle = weights[i] if i < len(weights) else np.pi / 4
        qc.cry(trainable_angle, i + 1, output_qubit)

    qc.barrier()
    
    # --- 4. 測量層 (Measurement) ---
    # 只測量輸出的 qubit[0] 到古典位元 c[0]
    qc.measure(output_qubit, 0)
    
    return qc


def main():
    parser = argparse.ArgumentParser(
        description="生成用於零日攻擊偵測的動態 QASM 電路。",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
範例:
  python generate_dynamic_qasm.py --qubits 7
  python generate_dynamic_qasm.py --qubits 8 --output q8.qasm
  python generate_dynamic_qasm.py --features 0.2,0.5,0.8,0.1,0.9,0.3
        """
    )
    parser.add_argument("--qubits", type=int, default=7, 
                        help="要使用的 qubit 總數量 (預設: 7)。")
    parser.add_argument("--output", type=str, 
                        default=f"qasm_output/daily_log_{datetime.now().strftime('%Y%m%d_%H%M%S')}.qasm", 
                        help="輸出的 .qasm 檔案名稱。")
    parser.add_argument("--features", type=str, default=None,
                        help="特徵向量，用逗號分隔 (例如: 0.2,0.5,0.8)。未提供時使用隨機特徵。")
    parser.add_argument("--weights", type=str, default=None,
                        help="訓練好的權重參數，用逗號分隔。未提供時使用預設角度。")
    args = parser.parse_args()

    print("="*70)
    print("  動態 QASM 電路生成器 - 零日攻擊偵測")
    print("="*70)
    print(f"時間: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"量子位元數: {args.qubits}")
    
    # 確保輸出目錄存在
    output_dir = os.path.dirname(args.output)
    if output_dir and not os.path.exists(output_dir):
        os.makedirs(output_dir)
        print(f"[INFO] 創建輸出目錄: {output_dir}")

    # 準備特徵
    num_features = args.qubits - 1
    if args.features:
        try:
            features = np.array([float(x.strip()) for x in args.features.split(',')])
            if len(features) < num_features:
                print(f"[WARNING] 提供的特徵數量 ({len(features)}) 少於需要的數量 ({num_features})，將補 0")
                features = np.pad(features, (0, num_features - len(features)))
            elif len(features) > num_features:
                print(f"[WARNING] 提供的特徵數量 ({len(features)}) 多於需要的數量 ({num_features})，將截斷")
                features = features[:num_features]
            print(f"[OK] 使用提供的特徵: {np.round(features, 3)}")
        except ValueError as e:
            print(f"[ERROR] 無法解析特徵向量: {e}")
            print("[INFO] 改用隨機特徵")
            features = np.random.rand(num_features)
    else:
        # !!! 重要 !!!
        # 在真實應用中，這裡的 features 應該來自你的 Windows Agent Log 處理結果
        # 目前我們先用隨機數據模擬
        features = np.random.rand(num_features)
        print(f"[INFO] 使用模擬特徵 (隨機): {np.round(features, 3)}")

    # 準備權重
    weights = None
    if args.weights:
        try:
            weights = np.array([float(x.strip()) for x in args.weights.split(',')])
            print(f"[OK] 使用提供的訓練權重: {np.round(weights, 3)}")
        except ValueError as e:
            print(f"[ERROR] 無法解析權重: {e}")
            print("[INFO] 使用預設角度")

    # 創建電路
    print(f"\n[*] 正在生成 {args.qubits}-qubit 的量子分類電路...")
    classifier_circuit = create_zero_day_classifier_circuit(features, args.qubits, weights)

    # 生成 QASM 2.0 程式碼
    qasm_code = qasm2.dumps(classifier_circuit)

    # 儲存到檔案
    with open(args.output, 'w') as f:
        f.write(qasm_code)
    
    print(f"\n[SUCCESS] 電路已成功生成並儲存至: {args.output}")
    print(f"[INFO] 檔案大小: {os.path.getsize(args.output)} bytes")
    print("\n電路預覽 (QASM 2.0):")
    print("-" * 70)
    print(qasm_code)
    print("-" * 70)
    
    # 顯示電路統計
    print(f"\n電路統計:")
    print(f"  - 量子位元數: {classifier_circuit.num_qubits}")
    print(f"  - 古典位元數: {classifier_circuit.num_clbits}")
    print(f"  - 閘門總數: {len(classifier_circuit.data)}")
    print(f"  - 電路深度: {classifier_circuit.depth()}")
    print("\n" + "="*70)


if __name__ == "__main__":
    main()

