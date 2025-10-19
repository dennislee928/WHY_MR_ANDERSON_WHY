#!/bin/bash
# ============================================================================
# 提交 Machine Learning 訓練的 QASM 到 IBM Quantum
# ============================================================================
# 功能: 
#   1. 訓練量子機器學習模型
#   2. 生成訓練好的 QASM 電路
#   3. 提交到 IBM Quantum 真實硬體或雲端模擬器
# ============================================================================

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 函數定義
print_header() {
    echo -e "${BLUE}============================================================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}============================================================================${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# ============================================================================
# 檢查環境
# ============================================================================
print_header "檢查環境配置"

# 檢查 IBM Token
if [ -z "$IBM_QUANTUM_TOKEN" ]; then
    print_error "IBM_QUANTUM_TOKEN 環境變數未設定"
    echo ""
    echo "請執行:"
    echo "  export IBM_QUANTUM_TOKEN='your_token_here'"
    echo ""
    echo "或使用參數:"
    echo "  $0 --token YOUR_TOKEN"
    exit 1
fi

print_success "IBM Token 已設定 (長度: ${#IBM_QUANTUM_TOKEN} 字元)"

# 解析參數
TRAINING_SAMPLES=50
MAX_ITERATIONS=30
USE_SIMULATOR=${USE_SIMULATOR:-false}
BACKEND_NAME=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --token)
            export IBM_QUANTUM_TOKEN="$2"
            shift 2
            ;;
        --samples)
            TRAINING_SAMPLES="$2"
            shift 2
            ;;
        --iterations)
            MAX_ITERATIONS="$2"
            shift 2
            ;;
        --simulator)
            USE_SIMULATOR=true
            shift
            ;;
        --backend)
            BACKEND_NAME="$2"
            shift 2
            ;;
        --help)
            echo "用法: $0 [選項]"
            echo ""
            echo "選項:"
            echo "  --token TOKEN        IBM Quantum Token"
            echo "  --samples N          訓練樣本數 (預設: 50)"
            echo "  --iterations N       訓練迭代數 (預設: 30)"
            echo "  --simulator          使用雲端模擬器"
            echo "  --backend NAME       指定後端名稱"
            echo "  --help               顯示此幫助訊息"
            echo ""
            echo "範例:"
            echo "  $0 --token YOUR_TOKEN --samples 100 --iterations 50"
            echo "  $0 --simulator --samples 30"
            exit 0
            ;;
        *)
            print_error "未知選項: $1"
            echo "使用 --help 查看用法"
            exit 1
            ;;
    esac
done

print_info "訓練參數: 樣本數=$TRAINING_SAMPLES, 迭代數=$MAX_ITERATIONS"
print_info "使用模擬器: $USE_SIMULATOR"

# ============================================================================
# 步驟 1: 訓練量子機器學習模型
# ============================================================================
print_header "步驟 1: 訓練量子分類器"

echo "正在訓練 VQC (Variational Quantum Classifier)..."
echo ""

python3 << END_PYTHON
import sys
import numpy as np
import json
from datetime import datetime
from train_quantum_classifier import (
    generate_training_data,
    create_trainable_classifier_circuit,
    train_with_vqc
)

print("="*60)
print("Quantum ML Training")
print("="*60)

# 生成訓練數據
print(f"\n生成訓練數據 (樣本數: $TRAINING_SAMPLES)...")
X_train, y_train = generate_training_data($TRAINING_SAMPLES)
print(f"✅ 訓練數據生成完成")
print(f"   特徵維度: {X_train.shape}")
print(f"   標籤分布: 正常={sum(y_train==0)}, 攻擊={sum(y_train==1)}")

# 創建可訓練電路
print(f"\n創建 VQC 電路...")
qubits = 7
ansatz, initial_weights = create_trainable_classifier_circuit(qubits)
print(f"✅ VQC 電路創建成功")
print(f"   量子位元: {qubits}")
print(f"   可訓練參數: {len(initial_weights)}")

# 訓練模型
print(f"\n開始訓練 (最大迭代: $MAX_ITERATIONS)...")
print("─"*60)

trained_weights, training_loss = train_with_vqc(
    X_train, y_train,
    ansatz, initial_weights,
    max_iterations=$MAX_ITERATIONS
)

print("─"*60)
print(f"\n✅ 訓練完成！")
print(f"   最終 Loss: {training_loss:.6f}")
print(f"   訓練權重: {len(trained_weights)} 個參數")

# 保存訓練好的權重
weights_file = "models/trained_weights.json"
metadata = {
    "timestamp": datetime.now().isoformat(),
    "training_samples": $TRAINING_SAMPLES,
    "max_iterations": $MAX_ITERATIONS,
    "final_loss": float(training_loss),
    "num_qubits": qubits,
    "num_weights": len(trained_weights),
    "weights": trained_weights.tolist()
}

import os
os.makedirs("models", exist_ok=True)
with open(weights_file, "w") as f:
    json.dump(metadata, f, indent=2)

print(f"\n💾 權重已保存: {weights_file}")

sys.exit(0)
END_PYTHON

if [ $? -ne 0 ]; then
    print_error "訓練失敗"
    exit 1
fi

print_success "訓練完成！"

# ============================================================================
# 步驟 2: 生成訓練好的 QASM 電路
# ============================================================================
print_header "步驟 2: 生成 ML QASM 電路"

echo "正在生成使用訓練權重的 QASM 電路..."
echo ""

python3 << END_PYTHON
import sys
import json
import numpy as np
from generate_dynamic_qasm import create_zero_day_classifier_circuit
from qiskit import qasm2

# 載入訓練好的權重
with open("models/trained_weights.json", "r") as f:
    metadata = json.load(f)

trained_weights = np.array(metadata["weights"])
print(f"✅ 載入訓練權重: {len(trained_weights)} 個參數")
print(f"   訓練時間: {metadata['timestamp']}")
print(f"   最終 Loss: {metadata['final_loss']:.6f}")

# 生成測試特徵（高風險情境）
test_features = np.array([0.06, 0.05, 0.2, 0.01, 0.033, 1.0])
qubits = metadata["num_qubits"]

print(f"\n生成 QASM 電路...")
print(f"   測試特徵: {test_features}")
print(f"   量子位元: {qubits}")

# 創建電路
circuit = create_zero_day_classifier_circuit(
    test_features,
    qubits,
    trained_weights
)

print(f"✅ 電路生成成功")
print(f"   深度: {circuit.depth()}")
print(f"   閘門數: {circuit.size()}")

# 保存 QASM 2.0 格式
qasm_str = qasm2.dumps(circuit)
qasm_file = "qasm_output/ml_trained_circuit.qasm"

import os
os.makedirs("qasm_output", exist_ok=True)
with open(qasm_file, "w") as f:
    f.write(qasm_str)

print(f"\n💾 QASM 已保存: {qasm_file}")
print(f"\nQASM 預覽:")
print("─"*60)
print(qasm_str[:500])
if len(qasm_str) > 500:
    print("...")
print("─"*60)

sys.exit(0)
END_PYTHON

if [ $? -ne 0 ]; then
    print_error "QASM 生成失敗"
    exit 1
fi

print_success "QASM 電路生成完成！"

# ============================================================================
# 步驟 3: 提交到 IBM Quantum
# ============================================================================
print_header "步驟 3: 提交到 IBM Quantum"

echo "正在連接到 IBM Quantum..."
echo ""

python3 << END_PYTHON
import sys
import os
import json
import numpy as np
from datetime import datetime
from qiskit_ibm_runtime import QiskitRuntimeService, Session, Sampler
from qiskit import QuantumCircuit, qasm2
from generate_dynamic_qasm import create_zero_day_classifier_circuit

# IBM Token
token = os.getenv('IBM_QUANTUM_TOKEN')
use_simulator = "$USE_SIMULATOR" == "true"
backend_name = "$BACKEND_NAME"

print("="*60)
print("IBM Quantum 提交")
print("="*60)

try:
    # 連接到 IBM Quantum
    print("\n正在連接到 IBM Quantum...")
    service = QiskitRuntimeService(
        channel='ibm_quantum',
        token=token
    )
    print("✅ 連接成功！")
    
    # 列出可用後端
    backends = service.backends()
    print(f"\n可用後端: {len(backends)} 個")
    
    # 選擇後端
    if backend_name:
        backend = service.backend(backend_name)
        print(f"✅ 使用指定後端: {backend.name}")
    elif use_simulator:
        # 尋找模擬器
        simulator = None
        for b in backends:
            if 'simulator' in b.name.lower():
                simulator = b
                break
        if simulator:
            backend = simulator
            print(f"✅ 使用雲端模擬器: {backend.name}")
        else:
            print("⚠️  未找到模擬器，使用第一個可用後端")
            backend = backends[0]
            print(f"✅ 使用後端: {backend.name}")
    else:
        # 使用第一個可用的真實量子後端
        backend = backends[0]
        print(f"✅ 使用後端: {backend.name}")
    
    # 顯示後端資訊
    status = backend.status()
    print(f"\n後端狀態:")
    print(f"   名稱: {backend.name}")
    if hasattr(status, 'pending_jobs'):
        print(f"   佇列: {status.pending_jobs} 個作業")
    if hasattr(backend, 'num_qubits'):
        print(f"   量子位元: {backend.num_qubits}")
    
    # 載入訓練好的電路
    print(f"\n載入 ML 訓練電路...")
    with open("models/trained_weights.json", "r") as f:
        metadata = json.load(f)
    
    trained_weights = np.array(metadata["weights"])
    test_features = np.array([0.06, 0.05, 0.2, 0.01, 0.033, 1.0])
    qubits = metadata["num_qubits"]
    
    circuit = create_zero_day_classifier_circuit(
        test_features,
        qubits,
        trained_weights
    )
    
    print(f"✅ 電路載入成功")
    print(f"   量子位元: {circuit.num_qubits}")
    print(f"   深度: {circuit.depth()}")
    print(f"   閘門: {circuit.size()}")
    
    # 提交作業
    print(f"\n正在提交作業到 {backend.name}...")
    print("⏳ 請稍候...")
    
    with Session(service=service, backend=backend.name) as session:
        sampler = Sampler(session=session)
        
        job = sampler.run([circuit], shots=1024)
        
        job_id = job.job_id()
        print(f"\n✅ 作業已提交！")
        print(f"   作業 ID: {job_id}")
        print(f"   狀態: {job.status()}")
        
        # 等待結果
        print(f"\n⏳ 等待量子執行完成...")
        result = job.result()
        
        print(f"✅ 量子執行完成！")
        
        # 分析結果
        pub_result = result[0]
        
        # 獲取計數（兼容 V2 API）
        counts = None
        for key in pub_result.data:
            if hasattr(pub_result.data[key], 'get_counts'):
                counts = pub_result.data[key].get_counts()
                break
        
        if counts:
            print(f"\n" + "="*60)
            print("量子分類結果")
            print("="*60)
            
            # 分析 qubit[0]
            zero_count = sum(c for state, c in counts.items() if state[-1] == '0')
            one_count = sum(c for state, c in counts.items() if state[-1] == '1')
            total = zero_count + one_count
            
            zero_prob = zero_count / total if total > 0 else 0
            one_prob = one_count / total if total > 0 else 0
            
            print(f"\nqubit[0] 測量:")
            print(f"   |0⟩ (正常): {zero_count:4d} ({zero_prob*100:5.1f}%)")
            print(f"   |1⟩ (攻擊): {one_count:4d} ({one_prob*100:5.1f}%)")
            
            # 判定
            threshold = 0.5
            is_attack = one_prob > threshold
            confidence = max(zero_prob, one_prob) * 100
            
            print(f"\n" + "="*60)
            if is_attack:
                print("判定: 🚨 零日攻擊偵測")
            else:
                print("判定: ✅ 正常行為")
            print(f"信心度: {confidence:.1f}%")
            print(f"後端: {backend.name}")
            print("="*60)
            
            # 保存結果
            result_file = f"results/ibm_result_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
            result_data = {
                "timestamp": datetime.now().isoformat(),
                "job_id": job_id,
                "backend": backend.name,
                "circuit_info": {
                    "qubits": circuit.num_qubits,
                    "depth": circuit.depth(),
                    "gates": circuit.size()
                },
                "measurements": {
                    "zero_count": int(zero_count),
                    "one_count": int(one_count),
                    "zero_prob": float(zero_prob),
                    "one_prob": float(one_prob)
                },
                "classification": {
                    "is_attack": bool(is_attack),
                    "confidence": float(confidence),
                    "threshold": float(threshold)
                },
                "training_info": metadata
            }
            
            os.makedirs("results", exist_ok=True)
            with open(result_file, "w") as f:
                json.dump(result_data, f, indent=2)
            
            print(f"\n💾 結果已保存: {result_file}")
            print(f"\n✅ IBM Quantum 提交完成！")
            sys.exit(0)
        else:
            print("\n⚠️  無法獲取測量結果")
            sys.exit(1)
            
except Exception as e:
    print(f"\n❌ 錯誤: {type(e).__name__}")
    print(f"   訊息: {str(e)[:200]}")
    
    import traceback
    print("\n完整錯誤:")
    traceback.print_exc()
    
    sys.exit(1)
END_PYTHON

RESULT=$?

if [ $RESULT -eq 0 ]; then
    print_success "所有步驟完成！"
    echo ""
    print_info "生成的檔案:"
    echo "  - models/trained_weights.json (訓練權重)"
    echo "  - qasm_output/ml_trained_circuit.qasm (QASM 電路)"
    echo "  - results/ibm_result_*.json (執行結果)"
else
    print_error "提交失敗"
    echo ""
    print_warning "可能的原因:"
    echo "  1. 網路連接問題"
    echo "  2. IBM Token 無效"
    echo "  3. 後端暫時不可用"
    echo ""
    print_info "嘗試使用本地模擬器:"
    echo "  python3 test_local_simulator.py"
fi

exit $RESULT

