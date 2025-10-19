# Qiskit Integration Guide - 真實量子計算整合指南

## 📋 概述

本指南說明如何將 **真實的量子計算** 整合到 Pandora Box Console IDS-IPS 系統中，使用 **IBM Quantum** 和 **Qiskit Runtime**。

---

## 🎯 架構概覽

```
┌─────────────────────────────────────────────────────────────┐
│                  Pandora Zero Trust System                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│          Zero Trust Quantum Predictor (FastAPI)             │
│  • Context Aggregation                                       │
│  • Hybrid Quantum-Classical ML                               │
│  • Quantum Policy Optimization                               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              Quantum Executor Service                        │
│  • Job Submission (async)                                    │
│  • Status Monitoring                                         │
│  • Result Retrieval                                          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Qiskit Runtime                              │
│  • Local Simulator (Aer)                                     │
│  • Cloud Simulator (ibmq_qasm_simulator)                     │
│  • Real Quantum Hardware (ibm_* devices)                     │
└─────────────────────────────────────────────────────────────┘
```

---

## 🚀 快速開始

### 1. 獲取 IBM Quantum Token

1. 訪問 https://quantum.ibm.com/
2. 登入或註冊帳號
3. 前往 **Account Settings** → **API Token**
4. 複製您的 API Token

### 2. 設置環境變數

創建 `.env` 文件：

```bash
# IBM Quantum Configuration
IBM_QUANTUM_TOKEN=your_token_here
QUANTUM_BACKEND=ibmq_qasm_simulator  # 或 ibm_kyoto, ibm_osaka, 等

# Zero Trust Configuration
ZERO_TRUST_ENABLED=true
QUANTUM_HYBRID_MODE=true
AUTO_RESPONSE_ENABLED=false
PREDICTION_THRESHOLD=0.7
RL_LEARNING_RATE=0.001
QUANTUM_DIMENSIONS=4
```

### 3. 測試連接

```bash
cd Experimental/cyber-ai-quantum
python test_ibm_connection.py
```

**預期輸出：**
```
=== IBM Quantum Connection Test ===

✅ Token loaded (40 characters)

連接到 IBM Quantum...
✅ IBM Quantum 連接成功!

可用後端總數: 15

模擬器:
  ✓ ibmq_qasm_simulator
  ✓ simulator_statevector
  ✓ simulator_mps
  ✓ simulator_extended_stabilizer
  ✓ simulator_stabilizer

真實量子處理器:
  ✓ ibm_kyoto (127 qubits, 排隊: 3)
  ✓ ibm_osaka (127 qubits, 排隊: 5)
  ✓ ibm_brisbane (127 qubits, 排隊: 2)
  ✓ ibm_sherbrooke (127 qubits, 排隊: 4)
  ✓ ibm_torino (133 qubits, 排隊: 1)

推薦後端（最不忙碌）:
  🎯 ibm_torino
     Qubits: 133
     排隊作業: 1
     狀態: active

✅ 連接測試完成！
```

---

## 📊 性能基準測試

### 運行完整基準測試套件

```bash
cd Experimental/cyber-ai-quantum
python benchmark_quantum_performance.py
```

### 運行特定測試

```bash
# 僅測試本地模擬器
python benchmark_quantum_performance.py local

# 僅測試雲端模擬器
python benchmark_quantum_performance.py cloud

# 僅測試電路優化
python benchmark_quantum_performance.py optimization

# 僅測試錯誤緩解
python benchmark_quantum_performance.py mitigation
```

### 基準測試結果範例

```
╔════════════════════════════════════════════════════════════╗
║   Pandora Quantum Performance Benchmark Suite              ║
╚════════════════════════════════════════════════════════════╝

============================================================
Benchmark 1: Local Aer Simulator
============================================================

✅ Local Simulator Results:
  Total time: 2.45s
  Avg per prediction: 245.0ms
  Throughput: 4.1 pred/s

============================================================
Benchmark 3: Circuit Transpilation & Optimization
============================================================

  Testing optimization level 0...
    Depth: 128 → 128 (0 gates saved)
    Time: 52.3ms

  Testing optimization level 3...
    Depth: 128 → 64 (64 gates saved)
    Time: 89.1ms

✅ Optimization Benchmark Complete
  Best optimization: Level 3
  Depth reduction: 64 gates

============================================================
PERFORMANCE COMPARISON SUMMARY
============================================================

┌─────────────────────────┬──────────────┬───────────────┬─────────────┐
│ Backend                 │ Avg Time     │ Throughput    │ Use Case    │
├─────────────────────────┼──────────────┼───────────────┼─────────────┤
│ Local Simulator         │    245.0ms │      4.1 p/s │ Development │
│ Cloud Simulator         │   1550.0ms │     ~0.6 p/s │ Testing     │
│ Real Hardware (est.)    │ ~90000ms │     ~0.01 p/s │ Research    │
└─────────────────────────┴──────────────┴───────────────┴─────────────┘

📋 建議:
  ✓ 開發/測試: 使用本地模擬器 (快速迭代)
  ✓ 生產環境: 混合執行（古典 + 量子，低風險用古典）
  ✓ 批次分析: 夜間提交到真實硬體（每日/每週/每月）
  ✓ 優化: 始終使用 optimization_level=3
  ✓ 錯誤緩解: 真實硬體使用 T-REx，關鍵任務使用 Combined
```

---

## 🔧 API 使用範例

### 1. Zero Trust 攻擊預測

```bash
curl -X POST http://localhost:8000/api/v1/zerotrust/predict \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_123",
    "device_trust": 0.8,
    "location_anomaly": 0.3,
    "behavior_score": 0.7,
    "network_features": [0.5, 0.2, 0.8, 0.1, 0.9, 0.4, 0.6, 0.3, 0.7, 0.5],
    "force_quantum": false
  }'
```

**響應：**
```json
{
  "attack_probability": 0.23,
  "trust_score": 0.77,
  "risk_level": "LOW",
  "threat_types": ["none"],
  "recommended_actions": ["allow"],
  "quantum_advantage": 0.05,
  "computation_time_ms": 245,
  "backend_used": "local_simulator",
  "timestamp": "2025-01-14T10:30:45Z"
}
```

### 2. 查詢量子作業狀態

```bash
# 提交作業後會返回 job_id
JOB_ID="c123456789abcdef"

curl http://localhost:8000/api/v1/quantum/jobs/$JOB_ID/status
```

**響應：**
```json
{
  "job_id": "c123456789abcdef",
  "status": "COMPLETED",
  "backend": "ibmq_qasm_simulator",
  "queue_position": null,
  "estimated_wait_time_seconds": null,
  "created_at": "2025-01-14T10:25:00Z",
  "completed_at": "2025-01-14T10:30:45Z"
}
```

### 3. 獲取作業結果

```bash
curl http://localhost:8000/api/v1/quantum/jobs/$JOB_ID/result
```

**響應：**
```json
{
  "job_id": "c123456789abcdef",
  "status": "COMPLETED",
  "result": {
    "counts": {"0": 512, "1": 512},
    "probabilities": {"0": 0.5, "1": 0.5}
  },
  "execution_time_seconds": 345,
  "shots": 1024,
  "backend": "ibmq_qasm_simulator"
}
```

---

## 🎯 混合量子-古典策略

系統使用智能混合執行策略：

### 1. 快速古典預測

```python
# 低風險場景 (< 70% 攻擊機率)
use_classical = True
response_time = < 10ms
```

### 2. 量子增強預測

```python
# 高風險場景 (>= 70% 攻擊機率)
use_quantum = True
response_time = 200-500ms (local) or 5-60s (cloud)
quantum_advantage = +5-15% accuracy
```

### 3. 批次量子分析

```python
# 定期深度分析
schedule = "daily 03:00 AM"
backend = "ibm_torino"  # 真實硬體
analysis_depth = "comprehensive"
```

---

## 📅 定期量子分析

### 設置 Cron 作業 (Linux)

```bash
cd Experimental/cyber-ai-quantum
chmod +x crontab_quantum.sh
./crontab_quantum.sh
```

### 設置排程任務 (Windows)

```powershell
cd Experimental\cyber-ai-quantum
.\schedule_quantum_tasks.ps1
```

### 手動運行

```bash
python scheduled_quantum_analysis.py --mode daily
python scheduled_quantum_analysis.py --mode weekly
python scheduled_quantum_analysis.py --mode monthly
```

---

## 🔍 監控量子作業

### Prometheus 指標

系統自動導出量子作業指標到 Prometheus：

```promql
# 總作業數
pandora_quantum_jobs_total{status="completed"}

# 平均執行時間
pandora_quantum_job_duration_seconds{backend="ibmq_qasm_simulator"}

# 錯誤率
rate(pandora_quantum_jobs_total{status="error"}[5m])

# 排隊時間
pandora_quantum_queue_wait_seconds{backend="ibm_torino"}
```

### 查看統計

```bash
curl http://localhost:8000/api/v1/quantum/executor/statistics
```

**響應：**
```json
{
  "total_jobs": 150,
  "completed_jobs": 142,
  "failed_jobs": 3,
  "pending_jobs": 5,
  "avg_execution_time_seconds": 12.5,
  "backends_used": {
    "local_simulator": 100,
    "ibmq_qasm_simulator": 40,
    "ibm_torino": 10
  }
}
```

---

## ⚙️ 進階配置

### 1. 電路優化

```python
# services/quantum_executor.py
transpile_options = {
    'optimization_level': 3,  # 0-3, 3 為最佳
    'seed_transpiler': 42,
    'scheduling_method': 'alap'
}
```

### 2. 錯誤緩解

```python
# 啟用 T-REx (Readout Error Mitigation)
resilience_options = {
    'resilience_level': 1,  # 0=off, 1=basic, 2=advanced
    'optimization_level': 3
}
```

### 3. 自訂後端選擇

```python
# 根據工作負載選擇後端
def select_backend(urgency: str, accuracy_required: float):
    if urgency == "immediate":
        return "local_simulator"
    elif accuracy_required > 0.95:
        return service.least_busy(simulator=False)  # 真實硬體
    else:
        return "ibmq_qasm_simulator"
```

---

## 🐛 故障排除

### 常見問題

#### 1. Token 無效

**錯誤：**
```
IBMAccountError: Invalid authentication credentials
```

**解決方案：**
- 確認 Token 正確（40 字元）
- 檢查 Token 未過期
- 重新生成 Token

#### 2. 後端不可用

**錯誤：**
```
IBMBackendError: Backend 'ibm_torino' is not available
```

**解決方案：**
- 檢查後端狀態：https://quantum.ibm.com/services/resources
- 使用 `least_busy()` 自動選擇
- 降級到模擬器

#### 3. 作業超時

**錯誤：**
```
JobTimeoutError: Job exceeded maximum wait time
```

**解決方案：**
- 增加 `max_time` 參數
- 使用異步 API，稍後查詢結果
- 考慮使用模擬器

---

## 📊 性能優化建議

### 1. 開發階段
- ✅ 使用本地模擬器 (Aer)
- ✅ 小規模電路 (< 10 qubits)
- ✅ 快速迭代

### 2. 測試階段
- ✅ 使用雲端模擬器
- ✅ 中等規模電路 (10-20 qubits)
- ✅ 驗證正確性

### 3. 生產階段
- ✅ 混合執行策略
- ✅ 批次量子作業（夜間）
- ✅ 啟用錯誤緩解
- ✅ 監控和告警

---

## 🎓 學習資源

- 📚 [Qiskit 官方文檔](https://qiskit.org/documentation/)
- 📚 [IBM Quantum 用戶指南](https://quantum.ibm.com/docs/)
- 📚 [Qiskit Runtime 教程](https://qiskit.org/documentation/partners/qiskit_ibm_runtime/)
- 📚 [Variational Quantum Algorithms](https://qiskit.org/textbook/ch-applications/vqe-molecules.html)

---

## 📝 總結

本指南涵蓋了從零到完整整合真實量子計算的所有步驟：

- ✅ IBM Quantum 帳號設置
- ✅ 環境配置
- ✅ 連接測試
- ✅ 性能基準測試
- ✅ API 使用
- ✅ 混合執行策略
- ✅ 定期分析
- ✅ 監控和故障排除

**下一步**：開始使用真實量子硬體進行 Zero Trust 攻擊預測！

---

**維護者**: Pandora Security Team  
**版本**: 1.0.0  
**最後更新**: 2025-01-14

