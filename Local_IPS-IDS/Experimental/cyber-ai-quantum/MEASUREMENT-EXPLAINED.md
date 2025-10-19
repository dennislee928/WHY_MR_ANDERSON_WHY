# Measurement 機制說明

## 📊 提交到 ibm_brisbane 或 ibm_torino 會自動 measurement 嗎？

### ✅ 是的，但需要理解兩種方式

---

## 方式 1: 電路內包含 Measurement（我們使用的）✅

### 當前實作（`generate_dynamic_qasm.py`）

```python
def create_zero_day_classifier_circuit(features, qubits, weights):
    # 創建電路：n 個量子位元 + 1 個古典位元
    qc = QuantumCircuit(qubits, 1, name="zero_day_classifier")
    
    # ... 特徵編碼、糾纏、變分層 ...
    
    # ✅ 測量層 (Measurement)
    output_qubit = 0  # qubit[0] 是輸出位元
    qc.measure(output_qubit, 0)  # 測量 qubit[0] → 古典位元 c[0]
    
    return qc
```

**生成的 QASM 2.0 格式**:
```qasm
OPENQASM 2.0;
include "qelib1.inc";

qreg q[7];        // 7 個量子位元
creg c[1];        // 1 個古典位元（用於儲存測量結果）

// ... 量子閘操作 ...

measure q[0] -> c[0];  // ✅ 測量 qubit[0]
```

### ✅ 提交到 IBM 後的流程

```
1. 提交電路 → IBM Quantum
   ↓
2. IBM 硬體執行所有量子閘
   ↓
3. 執行到 measure 指令時，自動進行測量
   ↓
4. 測量 qubit[0]，1024 次 (shots)
   ↓
5. 回傳結果：計數分布
   {
     '0': 628,  // qubit[0] 測量到 |0⟩ 628 次
     '1': 396   // qubit[0] 測量到 |1⟩ 396 次
   }
```

---

## 方式 2: 使用 Sampler 自動添加 Measurement

### Qiskit Runtime SamplerV2 的行為

**如果電路沒有 measurement**:
```python
# 電路沒有 measure
qc = QuantumCircuit(7)
qc.h(0)
# 沒有 qc.measure()

# 提交時 SamplerV2 會自動添加測量所有量子位元
sampler = SamplerV2(mode=backend)
job = sampler.run([qc], shots=1024)
# → 自動測量所有 7 個 qubits
```

**如果電路有 measurement**（我們的情況）:
```python
# 電路已包含 measure
qc = QuantumCircuit(7, 1)
qc.h(0)
qc.measure(0, 0)  # ✅ 明確指定只測量 qubit[0]

# 提交時使用我們指定的 measurement
sampler = SamplerV2(mode=backend)
job = sampler.run([qc], shots=1024)
# → 只測量 qubit[0]，結果存在 c[0]
```

---

## 🎯 我們的實作：明確控制 Measurement

### 為什麼只測量 qubit[0]？

根據 `new_spec.md` 的需求：

```
✅ 需求 6: 分析 qubit[0] measurement
   - 測量 qubit[0] 的機率分布
   - P(|0⟩) = 正常行為機率
   - P(|1⟩) = 零日攻擊機率
```

### 實作細節

```python
# 電路結構
qc = QuantumCircuit(7, 1)  # 7 qubits, 1 classical bit

# qubit 用途分配:
# - qubit[1] ~ qubit[6]: 特徵編碼（6 個特徵）
# - qubit[0]: 輸出位元（用於分類）

# 只測量輸出位元
qc.measure(0, 0)  # qubit[0] → classical_bit[0]
```

### 提交到 IBM 的 QASM

```qasm
OPENQASM 2.0;
include "qelib1.inc";

qreg q[7];
creg c[1];   // ✅ 只有 1 個古典位元，只存 qubit[0] 的結果

// ... 量子閘操作 ...

h q[1];
rx(0.1884955592153876) q[2];
rx(0.15707963267948966) q[3];
// ...

measure q[0] -> c[0];  // ✅ 只測量 qubit[0]
```

---

## 📊 實際執行結果

### 剛才的成功提交（Job ID: d3nhnq83qtks738ed9t0）

```python
# 原始電路
Circuit: 7 qubits, 13 depth, 18 gates
  - 包含 measure q[0] -> c[0]

# 轉譯到 ibm_brisbane
Transpiled: 131 depth, 229 gates
  - IBM 硬體自動優化和適配
  - measurement 指令保留

# 執行結果（1024 shots）
Results:
  |0>: 628 次 (61.3%)  ← qubit[0] 測量到 |0⟩
  |1>: 396 次 (38.7%)  ← qubit[0] 測量到 |1⟩
  
# 分類
P(|1⟩) = 38.7% < 50% → ✅ 正常行為
```

---

## 🔬 驗證 Measurement 是否存在

### 檢查 QASM 文件

```bash
# 查看生成的 QASM
python generate_dynamic_qasm.py --qubits 7 --output test.qasm

# 查看文件內容
cat test.qasm | grep measure
# 輸出: measure q[0] -> c[0];  ✅
```

### 檢查 QuantumCircuit 物件

```python
from generate_dynamic_qasm import create_zero_day_classifier_circuit
import numpy as np

circuit = create_zero_day_classifier_circuit(np.random.rand(6), 7)

print(f"量子位元: {circuit.num_qubits}")     # 7
print(f"古典位元: {circuit.num_clbits}")     # 1 ✅
print(f"測量操作: {circuit.count_ops().get('measure', 0)}")  # 1 ✅

# 查看電路圖
print(circuit.draw(output='text'))
```

**輸出範例**:
```
     ┌──────────┐     ┌──────────┐
q_0: ┤ CRY(π/4) ├─────┤ CRY(π/4) ├─────M  ← 測量
     └─┬──────┬─┘┌────┴────────┬─┘     ║
q_1: ──┤ RX(θ) ├──┤ RX(features) ├──────╫
       └───────┘  └─────────────┘      ║
                                       c: 1/══════════════════════════════╩
```

---

## ✅ 總結

### 問：提交到 ibm_brisbane 或 ibm_torino 會自動做 measurement 嗎？

**答：是的！**

| 項目 | 狀態 | 說明 |
|------|------|------|
| **電路包含 measure** | ✅ | `qc.measure(0, 0)` |
| **只測量 qubit[0]** | ✅ | 根據需求設計 |
| **IBM 自動執行** | ✅ | 執行 1024 shots |
| **回傳計數** | ✅ | `{'0': 628, '1': 396}` |
| **自動分類** | ✅ | 根據 qubit[0] 結果 |

### 優點

1. ✅ **明確控制**: 我們指定只測量 qubit[0]
2. ✅ **節省資源**: 不浪費在測量不需要的 qubits
3. ✅ **清晰語義**: qubit[0] = 分類輸出
4. ✅ **高效傳輸**: 只傳回 1 bit 的結果

### 如果想測量所有 qubits

```python
# 修改 generate_dynamic_qasm.py
qc = QuantumCircuit(7, 7)  # 7 個古典位元

# 測量所有
for i in range(7):
    qc.measure(i, i)

# 結果會是完整的量子態
# 例如: |0101010⟩, |1100011⟩ 等
```

---

**結論**: ✅ 當前實作已正確包含 measurement，提交到 IBM 真實硬體會自動執行並回傳 qubit[0] 的測量結果！

