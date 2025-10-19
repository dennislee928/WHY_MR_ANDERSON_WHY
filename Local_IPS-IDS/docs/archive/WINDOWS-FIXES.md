# Windows 環境修復指南

**日期**: 2025-01-14  
**問題**: Windows 環境的兼容性問題

---

## 🐛 問題 1: qiskit-ibm-runtime 版本不兼容

### 錯誤信息
```
ImportError: cannot import name 'BackendV1' from 'qiskit.providers.backend'
```

### 原因
- `qiskit-ibm-runtime 0.15.0` 與 `qiskit 2.2.1` 不兼容
- 需要使用兼容的版本組合

### 修復方案

#### 選項 1: 更新到最新版本（推薦）
```bash
cd Experimental/cyber-ai-quantum
pip install --upgrade qiskit qiskit-aer qiskit-ibm-runtime
```

#### 選項 2: 使用兼容版本
```bash
pip install qiskit==1.3.1 qiskit-aer==0.15.1 qiskit-ibm-runtime==0.30.0
```

#### 選項 3: 降級到穩定版本
```bash
pip install qiskit==0.45.3 qiskit-aer==0.13.3 qiskit-ibm-runtime==0.15.0
```

### 驗證
```bash
python -c "from qiskit_ibm_runtime import QiskitRuntimeService; print('OK')"
```

---

## 🐛 問題 2: Windows 缺少 OpenSSL

### 錯誤信息
```
openssl : 無法辨識 'openssl' 詞彙是否為 Cmdlet
```

### 原因
- Windows 預設不包含 OpenSSL
- 需要手動安裝

### 修復方案

#### 選項 1: 使用 Git for Windows 的 OpenSSL（推薦）
```powershell
# 添加到 PATH
$env:PATH += ";C:\Program Files\Git\usr\bin"

# 驗證
openssl version
```

#### 選項 2: 安裝 OpenSSL for Windows
1. 下載: https://slproweb.com/products/Win32OpenSSL.html
2. 安裝 "Win64 OpenSSL v3.x.x Light"
3. 添加到 PATH: `C:\Program Files\OpenSSL-Win64\bin`

#### 選項 3: 使用 Chocolatey
```powershell
choco install openssl
```

#### 選項 4: 使用 Docker 生成證書（最簡單）
```powershell
docker run --rm -v ${PWD}/certs:/certs alpine/openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout /certs/ca-key.pem -out /certs/ca-cert.pem -subj "/C=TW/ST=Taipei/L=Taipei/O=Pandora/CN=CA"
```

---

## 🔧 快速修復腳本

創建 `scripts/windows-setup.ps1`:

```powershell
# Windows 環境設置腳本

Write-Host "=== Windows 環境設置 ===" -ForegroundColor Cyan

# 1. 檢查 OpenSSL
if (Get-Command openssl -ErrorAction SilentlyContinue) {
    Write-Host "[OK] OpenSSL already installed" -ForegroundColor Green
} else {
    Write-Host "[INFO] Adding Git OpenSSL to PATH..." -ForegroundColor Yellow
    $env:PATH += ";C:\Program Files\Git\usr\bin"
    
    if (Get-Command openssl -ErrorAction SilentlyContinue) {
        Write-Host "[OK] OpenSSL found in Git" -ForegroundColor Green
    } else {
        Write-Host "[ERROR] OpenSSL not found" -ForegroundColor Red
        Write-Host "Please install OpenSSL or use Docker method" -ForegroundColor Yellow
    }
}

# 2. 更新 Qiskit
Write-Host "`n[INFO] Updating Qiskit packages..." -ForegroundColor Yellow
cd Experimental/cyber-ai-quantum
pip install --upgrade qiskit qiskit-aer qiskit-ibm-runtime --quiet

# 3. 驗證
Write-Host "`n[INFO] Verifying installation..." -ForegroundColor Yellow
python -c "from qiskit_ibm_runtime import QiskitRuntimeService; print('[OK] qiskit-ibm-runtime works')"

Write-Host "`n[SUCCESS] Windows environment setup completed!" -ForegroundColor Green
```

---

## 🚀 立即修復

### 步驟 1: 設置 OpenSSL
```powershell
# 臨時添加到 PATH
$env:PATH += ";C:\Program Files\Git\usr\bin"

# 永久添加（可選）
[Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";C:\Program Files\Git\usr\bin", "User")
```

### 步驟 2: 更新 Qiskit
```powershell
cd Experimental/cyber-ai-quantum
pip install --upgrade qiskit==1.3.1 qiskit-aer==0.15.1 qiskit-ibm-runtime==0.30.0
```

### 步驟 3: 測試
```powershell
# 設置 Token
$env:IBM_QUANTUM_TOKEN="7PzS0AdaFB83OVn7WghkdLupzthCe1fTl3yI4Qrp7G6o"

# 測試連接
python test_ibm_connection.py
```

### 步驟 4: 生成證書
```powershell
cd ../..
.\scripts\generate-grpc-certs.ps1
```

---

## ✅ 驗證清單

- [ ] OpenSSL 已安裝並在 PATH 中
- [ ] Qiskit 版本已更新
- [ ] IBM Quantum 連接測試成功
- [ ] gRPC 證書已生成
- [ ] 所有服務正常運行

---

**維護者**: Pandora Security Team  
**支援**: support@pandora-ids.com

