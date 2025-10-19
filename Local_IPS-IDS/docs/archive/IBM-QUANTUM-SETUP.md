# IBM Quantum 設置指南

> **版本**: 1.0.0  
> **更新日期**: 2025-01-14  
> **目標**: 整合真實量子硬體到 Pandora IDS-IPS

---

## 📋 概述

本指南將引導您完成從模擬量子計算到**真實 IBM Quantum 硬體**的完整設置過程。

---

## 🎯 Phase 0: 準備工作

### 步驟 1: 創建 IBM Quantum 帳號

1. **訪問 IBM Quantum**
   - 前往：https://quantum.ibm.com/
   - 點擊 "Sign in" 或 "Create an IBMid"

2. **註冊免費帳號**
   - 使用電子郵件註冊
   - 驗證郵箱
   - 完成個人資料設置

3. **獲取 API Token**
   ```
   登入後：
   1. 點擊右上角頭像
   2. 選擇 "Account settings"
   3. 找到 "API token" 部分
   4. 複製您的 token（40-50 字元的字符串）
   ```

**重要**: 妥善保管您的 Token，不要提交到 Git!

---

### 步驟 2: 配置環境變數

1. **創建環境配置文件**
   ```bash
   cd Experimental/cyber-ai-quantum
   cp env.example .env
   ```

2. **編輯 .env 文件**
   ```bash
   nano .env  # 或使用您喜歡的編輯器
   ```

3. **填入您的 Token**
   ```env
   # IBM Quantum 配置
   IBM_QUANTUM_TOKEN=your_actual_token_here_40_chars_long
   
   # 選擇後端
   QUANTUM_BACKEND=ibmq_qasm_simulator  # 開始使用模擬器
   QUANTUM_REAL_HARDWARE_ENABLED=false  # 稍後改為 true
   ```

---

### 步驟 3: 安裝 Qiskit 依賴

```bash
# 安裝完整 Qiskit 套件
pip install qiskit==0.45.0
pip install qiskit-aer==0.13.0
pip install qiskit-ibm-runtime==0.15.0
pip install qiskit-machine-learning==0.7.0

# 驗證安裝
python -c "import qiskit; print(f'Qiskit version: {qiskit.__version__}')"
```

---

## 🚀 Phase 1: 本地測試（模擬器）

### 測試 PoC 量子分類器

```bash
cd Experimental/cyber-ai-quantum
python poc_quantum_classifier.py
```

**預期輸出**:
```
=== Pandora Real Quantum Classifier PoC ===

--- Circuit Visualization Test ---
電路資訊:
  num_qubits: 4
  num_parameters: 12
  circuit_depth: 15
  ...

--- Quantum vs Classical Benchmark ---
Classical predictions: [0.48, 0.52, ...]
Classical total time: 0.085s
Classical avg time: 8.5ms

Quantum predictions: [0.51, 0.49, ...]
Quantum total time: 2.341s
Quantum avg time: 234.1ms

✅ PoC 完成！
```

**分析**:
- 量子模擬比 NumPy 慢 25-30倍（正常，因為模擬開銷）
- 功能性已驗證 ✅
- 準備好進入雲端測試

---

## ☁️ Phase 2: 雲端模擬器測試

### 步驟 1: 配置雲端後端

編輯 `.env`:
```env
QUANTUM_BACKEND=ibmq_qasm_simulator
QUANTUM_REAL_HARDWARE_ENABLED=false
```

### 步驟 2: 測試連接

創建測試腳本 `test_ibm_connection.py`:

```python
#!/usr/bin/env python3
import os
from dotenv import load_dotenv
from qiskit_ibm_runtime import QiskitRuntimeService

load_dotenv()

token = os.getenv('IBM_QUANTUM_TOKEN')

try:
    service = QiskitRuntimeService(channel='ibm_quantum', token=token)
    print("✅ IBM Quantum 連接成功!")
    
    # 列出可用後端
    backends = service.backends()
    print(f"\n可用後端 ({len(backends)}):")
    for backend in backends[:5]:
        print(f"  - {backend.name}")
    
except Exception as e:
    print(f"❌ 連接失敗: {e}")
```

運行：
```bash
python test_ibm_connection.py
```

---

## 🔧 Phase 3: 真實量子硬體

### 可用的免費量子處理器

IBM Quantum 提供以下免費設備（需排隊）：

| 設備名稱 | 量子位元 | 拓撲 | 平均排隊時間 |
|---------|---------|------|-------------|
| ibm_brisbane | 127 | Heavy-hex | ~2-5 分鐘 |
| ibm_kyoto | 127 | Heavy-hex | ~3-8 分鐘 |
| ibm_osaka | 127 | Heavy-hex | ~2-6 分鐘 |
| ibm_sherbrooke | 127 | Heavy-hex | ~5-15 分鐘 |

**推薦**: 使用 `service.least_busy()` 自動選擇

### 步驟 1: 啟用真實硬體

編輯 `.env`:
```env
QUANTUM_BACKEND=auto  # 自動選擇最少忙碌的
QUANTUM_REAL_HARDWARE_ENABLED=true
```

### 步驟 2: 提交第一個量子作業

```python
from poc_quantum_classifier import QuantumThreatClassifier

# 創建真實量子分類器
classifier = QuantumThreatClassifier(use_real_quantum=True)

# 預測（會提交到量子計算機排隊）
result = await classifier.predict(test_features)

print(f"量子預測結果: {result}")
```

### 步驟 3: 監控作業

```bash
# 查看作業狀態
curl http://localhost:8000/api/v1/quantum/jobs

# 檢查特定作業
curl http://localhost:8000/api/v1/quantum/jobs/{job_id}/status

# 獲取結果
curl http://localhost:8000/api/v1/quantum/jobs/{job_id}/result
```

---

## 📊 性能優化

### 電路優化

```python
from qiskit import transpile

# 轉譯電路以適應硬體拓撲
transpiled_circuit = transpile(
    circuit,
    backend,
    optimization_level=3,  # 0-3，3最優化
    seed_transpiler=42
)

print(f"原始深度: {circuit.depth()}")
print(f"優化後深度: {transpiled_circuit.depth()}")
```

### 錯誤緩解

```python
from qiskit_ibm_runtime import Estimator, Options

# 配置錯誤緩解
options = Options()
options.resilience_level = 1  # 0-3
options.optimization_level = 3

estimator = Estimator(backend=backend, options=options)
```

---

## 🔄 混合執行策略

Pandora 使用智能混合策略：

```
輸入威脅
    ↓
快速古典預測
    ↓
風險 < 70%？
    ├─ Yes → 返回古典結果 (< 10ms)
    └─ No  → 提交量子作業
                ↓
           排隊 (1-10分鐘)
                ↓
           量子執行 (~30秒)
                ↓
           返回精確結果
```

**優勢**:
- 95% 的請求在 <10ms 內完成（古典）
- 5% 的高風險請求獲得量子級精確度
- 最佳化量子資源使用

---

## 📅 定期量子分析

### 每日分析 (凌晨 2:00)

```bash
python scheduled_quantum_analysis.py daily
```

**任務**:
- 重新評估過去24小時的高風險事件
- 使用真實量子計算進行深度分析
- 識別誤報和遺漏

### 每週訓練 (週日 凌晨 3:00)

```bash
python scheduled_quantum_analysis.py weekly
```

**任務**:
- 使用過去一週數據重訓練 VQC 參數
- 更新量子模型權重
- 評估模型性能

### 每月批次 (每月1號 凌晨 4:00)

```bash
python scheduled_quantum_analysis.py monthly
```

**任務**:
- 分析過去30天的所有事件
- 識別長期威脅趨勢
- 生成詳細報告

---

## 🎯 配置排程（Windows）

```powershell
# 運行排程腳本
cd Experimental\cyber-ai-quantum
.\schedule_quantum_tasks.ps1
```

這將創建 3 個 Windows 排程任務：
- `Pandora_Daily_Quantum_Analysis`
- `Pandora_Weekly_Quantum_Training`
- `Pandora_Monthly_Quantum_Batch`

### 管理排程

```powershell
# 查看任務
Get-ScheduledTask | Where-Object {$_.TaskName -like "Pandora_*"}

# 手動運行
Start-ScheduledTask -TaskName "Pandora_Daily_Quantum_Analysis"

# 停用
Disable-ScheduledTask -TaskName "Pandora_Daily_Quantum_Analysis"
```

---

## 📈 監控與日誌

### Prometheus 指標

```bash
curl http://localhost:8000/api/v1/quantum/executor/statistics
```

**輸出**:
```json
{
  "total_jobs": 145,
  "status_distribution": {
    "DONE": 120,
    "RUNNING": 3,
    "QUEUED": 2,
    "ERROR": 5
  },
  "average_execution_time_seconds": 45.2,
  "backend_type": "real_hardware"
}
```

### 日誌位置

- 每日分析：`analysis_results/daily_*.json`
- 每週訓練：`analysis_results/weekly_training_*.json`
- 每月批次：`analysis_results/monthly_batch_*.json`

---

## ⚠️ 常見問題

### Q1: Token 驗證失敗
```
Error: Invalid IBM Quantum token
```

**解決**:
1. 檢查 `.env` 中的 token 是否正確
2. 確保沒有多餘空格
3. 重新從 IBM Quantum 網站複製 token

### Q2: 作業長時間排隊
```
Job status: QUEUED for 30 minutes
```

**解決**:
1. 使用 `service.least_busy()` 自動選擇
2. 改用雲端模擬器進行開發
3. 考慮付費帳號獲得優先級

### Q3: 電路太大無法執行
```
Error: Circuit has 127 qubits but backend only supports 5
```

**解決**:
1. 減少 `num_qubits` 參數
2. 使用 `transpile()` 優化電路
3. 選擇更大的量子處理器

### Q4: 結果質量差
```
Accuracy: 52% (barely better than random)
```

**解決**:
1. 啟用錯誤緩解 (`resilience_level=1`)
2. 增加電路重複次數提高信噪比
3. 使用更多訓練迭代

---

## 🔐 安全注意事項

### API Token 安全

```bash
# ✅ 正確：使用環境變數
export IBM_QUANTUM_TOKEN="your_token"

# ❌ 錯誤：硬編碼在代碼中
token = "abc123..."  # 永遠不要這樣做！

# ✅ 正確：.env 文件（並加入 .gitignore）
echo ".env" >> .gitignore
```

### Docker 部署

```dockerfile
# Dockerfile
ENV IBM_QUANTUM_TOKEN=${IBM_QUANTUM_TOKEN}

# docker-compose.yml
environment:
  - IBM_QUANTUM_TOKEN=${IBM_QUANTUM_TOKEN}

# 運行時傳入
docker-compose up -d
```

---

## 📊 成本與配額

### 免費帳號限制

| 項目 | 限制 |
|------|------|
| 每月作業數 | 未明確限制 |
| 並發作業 | 1-3 個 |
| 作業優先級 | 低 |
| 電路深度 | 建議 < 100 |
| 執行時間 | < 3 小時 |

### 使用建議

1. **開發階段**: 使用本地 `AerSimulator` ✅
2. **測試階段**: 使用雲端 `ibmq_qasm_simulator` ✅
3. **生產環境**: 混合執行（古典 + 量子）✅
4. **研究目的**: 真實硬體批次分析 ✅

---

## 🎓 學習資源

### IBM Quantum 文檔
- 官方文檔：https://docs.quantum.ibm.com/
- Qiskit 教程：https://qiskit.org/textbook
- Runtime API：https://docs.quantum.ibm.com/api/qiskit-ibm-runtime

### Pandora 相關文檔
- `QISKIT-INTEGRATION-GUIDE.md` - 技術整合指南
- `poc_quantum_classifier.py` - 完整 PoC 實現
- `services/quantum_executor.py` - 執行器服務

---

## 🚀 快速開始

### 1分鐘測試

```bash
# 1. 設置 Token
export IBM_QUANTUM_TOKEN="your_token"

# 2. 測試連接
cd Experimental/cyber-ai-quantum
python -c "
from services.quantum_executor import get_quantum_executor
executor = get_quantum_executor()
print(executor.get_statistics())
"

# 3. 運行 PoC
python poc_quantum_classifier.py
```

### 整合到 FastAPI

```bash
# 啟動服務
cd Application
docker-compose up -d cyber-ai-quantum

# 測試 API
curl http://localhost:8000/api/v1/quantum/executor/statistics
```

---

## 📞 支援

### IBM Quantum 支援
- Slack: https://qiskit.slack.com/
- GitHub: https://github.com/Qiskit

### Pandora 團隊
- 技術文檔：`docs/`
- Issue Tracker: GitHub Issues
- Email: support@pandora-ids.com

---

**準備好開始您的量子之旅了嗎？** 🚀

從本地模擬器開始，逐步過渡到真實量子硬體，打造世界級的量子增強安全系統！

