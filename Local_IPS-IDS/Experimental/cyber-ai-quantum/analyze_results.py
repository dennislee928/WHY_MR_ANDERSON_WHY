#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
量子作業結果分析器
讀取作業結果，根據 qubit[0] 的測量統計進行分類

分類邏輯:
- qubit[0] = 1 → Zero-Day Attack (Potential)
- qubit[0] = 0 → Known Attack / Benign
"""

import json
import argparse
import sys
import os
from datetime import datetime
from typing import Dict, Tuple

# Windows UTF-8 兼容性
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')


def analyze_classification_results(result_filename: str, threshold: float = 0.5, 
                                   output_report: bool = True) -> Dict:
    """
    分析分類結果並輸出報告。

    Args:
        result_filename (str): 包含作業結果的 JSON 檔案路徑。
        threshold (float): 判定為 Zero-Day Attack 的機率閾值 (預設: 0.5)。
        output_report (bool): 是否輸出詳細報告到控制台 (預設: True)。

    Returns:
        Dict: 包含分析結果的字典
    """
    try:
        with open(result_filename, 'r', encoding='utf-8') as f:
            data = json.load(f)
    except FileNotFoundError:
        print(f"[ERROR] 找不到檔案: {result_filename}")
        sys.exit(1)
    except json.JSONDecodeError:
        print(f"[ERROR] 檔案格式錯誤，無法解析 JSON: {result_filename}")
        sys.exit(1)

    counts = data.get('counts', {})
    if not counts:
        print("[WARNING] 結果中沒有 'counts' 數據。")
        return {}

    # --- 核心分析邏輯 ---
    zero_day_counts = 0  # qubit[0] 測量為 '1'
    known_attack_counts = 0  # qubit[0] 測量為 '0'
    bitstring_details = []

    if output_report:
        print("="*70)
        print("  📊 零日攻擊分類分析報告")
        print("="*70)
        print(f"Job ID: {data.get('job_id', 'N/A')}")
        print(f"Backend: {data.get('backend', 'N/A')}")
        print(f"總測量次數 (Shots): {data.get('shots', 'N/A')}")
        print(f"分析時間: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print("-"*70)
        print("\n詳細測量結果分析:")
    
    # bitstring 在 Qiskit 中是反向的，所以 c[0] 對應的是最右邊的位元
    for bitstring, count in sorted(counts.items(), key=lambda x: x[1], reverse=True):
        qubit0_state = bitstring[-1]  # 獲取 qubit[0] 的狀態
        
        if qubit0_state == '1':
            zero_day_counts += count
            attack_type = "Zero-Day Attack (Potential)"
            risk_level = "🔴 HIGH"
        else:
            known_attack_counts += count
            attack_type = "Known Attack / Benign"
            risk_level = "🟢 LOW"
        
        bitstring_details.append({
            'bitstring': bitstring,
            'qubit0': qubit0_state,
            'type': attack_type,
            'count': count
        })
        
        if output_report:
            print(f"  [{risk_level}] Bitstring: '{bitstring}' → qubit[0]='{qubit0_state}' → {attack_type:<28} | 次數: {count}")

    total_shots = zero_day_counts + known_attack_counts
    if total_shots == 0:
        print("\n[ERROR] 總測量次數為 0，無法計算機率。")
        return {}

    # --- 統計與推論 ---
    prob_zero_day = zero_day_counts / total_shots
    prob_known_attack = known_attack_counts / total_shots

    if output_report:
        print("\n" + "-"*70)
        print("統計摘要:")
        print(f"  - 總計 'Zero-Day' (qubit[0]=1) 次數: {zero_day_counts}")
        print(f"  - 總計 'Known Attack' (qubit[0]=0) 次數: {known_attack_counts}")
        print(f"  - P(|1⟩) 機率 (判定為 Zero-Day): {prob_zero_day:.2%}")
        print(f"  - P(|0⟩) 機率 (判定為 Known Attack): {prob_known_attack:.2%}")
        print("-"*70)

    # --- 最終結論 ---
    is_zero_day = prob_zero_day > threshold
    confidence = prob_zero_day if is_zero_day else prob_known_attack
    
    if output_report:
        print("\n最終推論:")
        if is_zero_day:
            final_conclusion = "高度可能為 Zero-Day Attack"
            print(f"  [🔴 CRITICAL] {final_conclusion}")
            print(f"     原因: P(|1⟩) = {prob_zero_day:.2%}，超過了 {threshold:.2%} 的閾值。")
            print(f"     建議: 立即啟動事件回應程序，隔離可疑主機，進行深度分析。")
        else:
            final_conclusion = "較可能為已知攻擊或正常流量"
            print(f"  [🟢 INFO] {final_conclusion}")
            print(f"     原因: P(|1⟩) = {prob_zero_day:.2%}，未達到 {threshold:.2%} 的閾值。")
            print(f"     建議: 持續監控，記錄日誌供後續分析。")
        
        print("="*70)

    # 返回分析結果
    analysis_result = {
        'job_id': data.get('job_id', 'N/A'),
        'backend': data.get('backend', 'N/A'),
        'shots': data.get('shots', 0),
        'timestamp': datetime.now().isoformat(),
        'zero_day_count': zero_day_counts,
        'known_attack_count': known_attack_counts,
        'prob_zero_day': prob_zero_day,
        'prob_known_attack': prob_known_attack,
        'is_zero_day': is_zero_day,
        'confidence': confidence,
        'threshold': threshold,
        'conclusion': final_conclusion,
        'bitstring_details': bitstring_details
    }

    return analysis_result


def save_analysis_report(analysis_result: Dict, output_file: str = None):
    """
    儲存分析報告為 JSON 檔案。

    Args:
        analysis_result (Dict): 分析結果字典
        output_file (str): 輸出檔案路徑，若為 None 則自動生成
    """
    if output_file is None:
        timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
        output_file = f"results/analysis_report_{timestamp}.json"
    
    # 確保目錄存在
    output_dir = os.path.dirname(output_file)
    if output_dir and not os.path.exists(output_dir):
        os.makedirs(output_dir)
    
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(analysis_result, f, indent=2, ensure_ascii=False)
    
    print(f"\n[INFO] 分析報告已儲存至: {output_file}")


def main():
    parser = argparse.ArgumentParser(
        description="分析量子作業的分類結果。",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
範例:
  python analyze_results.py results/result_<job_id>.json
  python analyze_results.py results/result_<job_id>.json --threshold 0.6
  python analyze_results.py results/result_<job_id>.json --save
        """
    )
    parser.add_argument("result_file", type=str, 
                        help="包含 'counts' 的 JSON 結果檔案路徑。")
    parser.add_argument("--threshold", type=float, default=0.5, 
                        help="判定為 Zero-Day Attack 的機率閾值 (預設: 0.5)。")
    parser.add_argument("--save", action="store_true", 
                        help="儲存分析報告為 JSON 檔案。")
    parser.add_argument("--output", type=str, default=None,
                        help="指定輸出報告檔案路徑 (需配合 --save 使用)。")
    args = parser.parse_args()

    # 執行分析
    analysis_result = analyze_classification_results(args.result_file, args.threshold)

    # 儲存報告
    if args.save and analysis_result:
        save_analysis_report(analysis_result, args.output)


if __name__ == "__main__":
    main()

