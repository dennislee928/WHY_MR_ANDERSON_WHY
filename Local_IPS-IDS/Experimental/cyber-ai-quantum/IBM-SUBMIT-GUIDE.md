# IBM Quantum ML QASM 提交指南

**目的**: 訓練量子機器學習模型並提交 QASM 到 IBM Quantum 真實硬體

---

## 📋 功能說明

腳本 `submit_ml_qasm_to_ibm.sh` 會自動完成以下步驟：

1. ✅ 訓練 VQC (Variational Quantum Classifier) 模型
2. ✅ 保存訓練好的權重
3. ✅ 生成使用訓練權重的 QASM 電路
4. ✅ 提交到 IBM Quantum 真實硬體或雲端模擬器
5. ✅ 獲取並分析結果
6. ✅ 保存完整報告

---

## 🚀 使用方式

### 方式 1: 在容器內執行（推薦）

```bash
# 進入容器
docker exec -it cyber-ai-quantum bash

# 執行腳本（使用預設參數）
./submit_ml_qasm_to_ibm.sh

# 或使用自訂參數
./submit_ml_qasm_to_ibm.sh --samples 100 --iterations 50

# 使用雲端模擬器（推薦先測試）
./submit_ml_qasm_to_ibm.sh --simulator

# 指定特定後端
./submit_ml_qasm_to_ibm.sh --backend ibm_brisbane
```

### 方式 2: 從 Host 執行

```bash
# Windows PowerShell
docker exec cyber-ai-quantum bash -c "cd /app && ./submit_ml_qasm_to_ibm.sh --simulator"

# Linux/macOS
docker exec cyber-ai-quantum bash -c "cd /app && ./submit_ml_qasm_to_ibm.sh --simulator"
```

---

## 📝 參數說明

| 參數 | 說明 | 預設值 |
|------|------|--------|
| `--token TOKEN` | IBM Quantum Token | 從環境變數 |
| `--samples N` | 訓練樣本數 | 50 |
| `--iterations N` | 訓練迭代數 | 30 |
| `--simulator` | 使用雲端模擬器 | false |
| `--backend NAME` | 指定後端名稱 | 自動選擇 |
| `--help` | 顯示幫助 | - |

---

## 📊 執行範例

### 範例 1: 快速測試（雲端模擬器）

```bash
docker exec cyber-ai-quantum bash -c "cd /app && ./submit_ml_qasm_to_ibm.sh --simulator --samples 30 --iterations 20"
```

**預期輸出**:
```
============================================================================
  檢查環境配置
============================================================================
✅ IBM Token 已設定 (長度: 44 字元)
ℹ️  訓練參數: 樣本數=30, 迭代數=20
ℹ️  使用模擬器: true

============================================================================
  步驟 1: 訓練量子分類器
============================================================================
...
✅ 訓練完成！
   最終 Loss: 0.234567
   訓練權重: 14 個參數

============================================================================
  步驟 2: 生成 ML QASM 電路
============================================================================
✅ 電路生成成功
   深度: 13
   閘門數: 18

============================================================================
  步驟 3: 提交到 IBM Quantum
============================================================================
✅ 連接成功！
✅ 使用雲端模擬器: ibm_qasm_simulator
...
✅ 量子執行完成！

============================================================
量子分類結果
============================================================
qubit[0] 測量:
   |0⟩ (正常):  456 ( 44.5%)
   |1⟩ (攻擊):  568 ( 55.5%)

============================================================
判定: 🚨 零日攻擊偵測
信心度: 55.5%
後端: ibm_qasm_simulator
============================================================

💾 結果已保存: results/ibm_result_20251015_103045.json
✅ IBM Quantum 提交完成！
```

### 範例 2: 提交到真實量子硬體

```bash
# 查看可用後端
docker exec cyber-ai-quantum python -c "
from qiskit_ibm_runtime import QiskitRuntimeService
import os
service = QiskitRuntimeService(channel='ibm_quantum', token=os.getenv('IBM_QUANTUM_TOKEN'))
backends = service.backends()
print('可用後端:')
for i, b in enumerate(backends[:10]):
    print(f'  {i+1}. {b.name}')
"

# 提交到特定後端
docker exec cyber-ai-quantum bash -c "cd /app && ./submit_ml_qasm_to_ibm.sh --backend ibm_brisbane --samples 50"
```

### 範例 3: 大規模訓練

```bash
# 使用更多樣本和迭代次數
docker exec cyber-ai-quantum bash -c "cd /app && ./submit_ml_qasm_to_ibm.sh --samples 200 --iterations 100 --simulator"
```

---

## 📁 生成的檔案

執行後會生成以下檔案：

```
models/
  └── trained_weights.json          # 訓練好的模型權重

qasm_output/
  └── ml_trained_circuit.qasm       # QASM 2.0 格式電路

results/
  └── ibm_result_YYYYMMDD_HHMMSS.json  # IBM Quantum 執行結果
```

### 結果檔案格式

```json
{
  "timestamp": "2025-10-15T10:30:45.123456",
  "job_id": "ch6jab6cgf...",
  "backend": "ibm_qasm_simulator",
  "circuit_info": {
    "qubits": 7,
    "depth": 13,
    "gates": 18
  },
  "measurements": {
    "zero_count": 456,
    "one_count": 568,
    "zero_prob": 0.445,
    "one_prob": 0.555
  },
  "classification": {
    "is_attack": true,
    "confidence": 55.5,
    "threshold": 0.5
  },
  "training_info": {
    "training_samples": 50,
    "max_iterations": 30,
    "final_loss": 0.234567
  }
}
```

---

## 🔧 故障排除

### 問題 1: 網路連接失敗

**錯誤**:
```
Failed to resolve 'auth.quantum-computing.ibm.com'
```

**解決方案**:
```bash
# 1. 檢查 DNS
docker exec cyber-ai-quantum cat /etc/resolv.conf

# 2. 更新 docker-compose.yml
services:
  cyber-ai-quantum:
    dns:
      - 8.8.8.8
      - 8.8.4.4

# 3. 重啟容器
docker-compose restart cyber-ai-quantum

# 4. 或使用本地模擬器
python test_local_simulator.py
```

### 問題 2: Token 無效

**錯誤**:
```
Authentication failed
```

**解決方案**:
```bash
# 檢查 Token
echo $IBM_QUANTUM_TOKEN

# 重新設定
export IBM_QUANTUM_TOKEN="your_new_token"

# 或在 docker-compose.yml 更新
environment:
  - IBM_QUANTUM_TOKEN=your_new_token
```

### 問題 3: 後端佇列滿

**訊息**:
```
pending_jobs: 50
```

**解決方案**:
```bash
# 使用雲端模擬器
./submit_ml_qasm_to_ibm.sh --simulator

# 或選擇佇列較少的後端
./submit_ml_qasm_to_ibm.sh --backend ibm_cairo
```

---

## 🎯 建議的測試流程

### 第一次使用

1. **測試本地模擬器** ✅（確保代碼正常）
   ```bash
   docker exec cyber-ai-quantum python test_local_simulator.py
   ```

2. **測試雲端模擬器** ✅（確保 IBM 連接正常）
   ```bash
   docker exec cyber-ai-quantum bash -c "cd /app && ./submit_ml_qasm_to_ibm.sh --simulator --samples 30"
   ```

3. **提交到真實硬體** 🎯（實際量子計算）
   ```bash
   docker exec cyber-ai-quantum bash -c "cd /app && ./submit_ml_qasm_to_ibm.sh --samples 50 --iterations 30"
   ```

---

## 📊 性能建議

| 場景 | 樣本數 | 迭代數 | 預計時間 |
|------|--------|--------|----------|
| 快速測試 | 30 | 20 | ~30 秒 |
| 標準訓練 | 50 | 30 | ~1 分鐘 |
| 高品質訓練 | 100 | 50 | ~3 分鐘 |
| 最佳訓練 | 200 | 100 | ~10 分鐘 |

**注意**: IBM Quantum 真實硬體可能需要等待佇列（幾分鐘到幾小時）

---

## 🔄 整合到自動化工作流程

### 每日自動訓練並提交

**方式 1: Crontab (Linux/macOS)**
```bash
# 每天凌晨 2:00 執行
0 2 * * * docker exec cyber-ai-quantum bash -c "cd /app && ./submit_ml_qasm_to_ibm.sh --simulator" >> /var/log/quantum-ml.log 2>&1
```

**方式 2: Windows Task Scheduler**
```powershell
# 創建排程任務
$Action = New-ScheduledTaskAction -Execute "docker" -Argument "exec cyber-ai-quantum bash -c 'cd /app && ./submit_ml_qasm_to_ibm.sh --simulator'"
$Trigger = New-ScheduledTaskTrigger -Daily -At "02:00"
Register-ScheduledTask -TaskName "QuantumML-Daily" -Action $Action -Trigger $Trigger
```

---

## ✅ 檢查清單

提交前檢查：

- [ ] IBM Token 已設定且有效
- [ ] Docker 容器運行中
- [ ] 網路連接正常
- [ ] 已測試本地模擬器
- [ ] 已測試雲端模擬器（可選）
- [ ] 了解預期等待時間

---

## 📚 相關文件

- **IBM Quantum 文檔**: https://quantum.ibm.com/docs
- **Qiskit Runtime API**: https://docs.quantum.ibm.com/api/qiskit-ibm-runtime
- **本地測試**: `test_local_simulator.py`
- **API 端點**: `http://localhost:8000/docs`

---

**最後更新**: 2025-10-15  
**狀態**: ✅ 生產就緒

