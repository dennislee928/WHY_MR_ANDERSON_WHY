# 🚀 IBM Quantum 自動執行指南

## 執行方式

### Git Bash 中執行（推薦）

```bash
# 1. 切換到目錄
cd ~/Documents/GitHub/Local_IPS-IDS/Experimental/cyber-ai-quantum

# 2. 設定 Token（使用 Git Bash 語法）
export IBM_QUANTUM_TOKEN="7PzS0AdaFB83OVn7WghkdLupzthCe1fTl3yI4Qrp7G6o"

# 3. 執行 10 分鐘循環
python auto_submit_every_10min.py

# 或單次執行
python test_host_ibm.py
```

### Windows PowerShell 中執行

```powershell
# 1. 切換到目錄
cd C:\Users\dennis.lee\Documents\GitHub\Local_IPS-IDS\Experimental\cyber-ai-quantum

# 2. 設定 Token（使用 PowerShell 語法）
$env:IBM_QUANTUM_TOKEN="7PzS0AdaFB83OVn7WghkdLupzthCe1fTl3yI4Qrp7G6o"

# 3. 執行 10 分鐘循環
python auto_submit_every_10min.py

# 或單次執行
python test_host_ibm.py
```

---

## 📋 腳本說明

### `auto_submit_every_10min.py` - 自動循環執行

**功能**:
- ✅ 每 10 分鐘自動執行一次
- ✅ 提交 ML QASM 到 IBM 真實硬體
- ✅ 自動保存結果到 `results/`
- ✅ 按 Ctrl+C 停止

**輸出範例**:
```
======================================================================
執行次數: 1
當前時間: 2025-10-15 12:10:00
======================================================================

[1/5] 生成 ML 量子電路...
✅ 電路創建: 7 qubits, 13 depth, 18 gates

[2/5] 連接 IBM Quantum (ibm_cloud channel)...
✅ 連接成功！

[3/5] 選擇後端 (可用: 2 個)...
✅ 使用: ibm_brisbane

[4/5] 轉譯到 ibm_brisbane...
✅ 轉譯完成: 131 depth, 229 gates

[5/5] 提交到 ibm_brisbane...
✅ 作業已提交: d3nhnq83qtks738ed9t0
⏳ 等待結果...
✅ 執行完成！

======================================================================
作業 d3nhnq83qtks738ed9t0 - 分類結果
======================================================================
  |0> (正常):  628 ( 61.3%)
  |1> (攻擊):  396 ( 38.7%)

  判定: ✅ 正常行為
  信心度: 61.3%
  後端: ibm_brisbane
======================================================================

💾 結果已保存: results/auto_20251015_121045.json
✅ 本次執行成功！

⏰ 下次執行時間: 2025-10-15 12:20:00
⏳ 等待 10 分鐘...
   (按 Ctrl+C 停止)
```

### `test_host_ibm.py` - 單次執行

**功能**:
- ✅ 執行一次後結束
- ✅ 適合手動測試
- ✅ 提交到 IBM 真實硬體

---

## 🔧 常見問題

### 問題 1: Token 未設定

**Git Bash**:
```bash
export IBM_QUANTUM_TOKEN="你的Token"
```

**PowerShell**:
```powershell
$env:IBM_QUANTUM_TOKEN="你的Token"
```

### 問題 2: 模組未安裝

```bash
pip install -r requirements.txt
```

### 問題 3: 停止循環執行

按 `Ctrl+C` 即可安全停止

---

## 📊 預期行為

### 每 10 分鐘執行

```
12:00 → 提交作業 1
12:10 → 提交作業 2
12:20 → 提交作業 3
...
```

### 保存的結果

```
results/
├── auto_20251015_120000.json
├── auto_20251015_121000.json
├── auto_20251015_122000.json
└── ...
```

每個檔案包含：
- 作業 ID
- 後端名稱
- 分類結果
- 時間戳記

---

## ✅ 快速啟動

**在 Git Bash 中複製貼上**:
```bash
cd ~/Documents/GitHub/Local_IPS-IDS/Experimental/cyber-ai-quantum
export IBM_QUANTUM_TOKEN="7PzS0AdaFB83OVn7WghkdLupzthCe1fTl3yI4Qrp7G6o"
python auto_submit_every_10min.py
```

**在 PowerShell 中複製貼上**:
```powershell
cd C:\Users\dennis.lee\Documents\GitHub\Local_IPS-IDS\Experimental\cyber-ai-quantum
$env:IBM_QUANTUM_TOKEN="7PzS0AdaFB83OVn7WghkdLupzthCe1fTl3yI4Qrp7G6o"
python auto_submit_every_10min.py
```

---

**最後更新**: 2025-10-15  
**狀態**: ✅ 就緒

