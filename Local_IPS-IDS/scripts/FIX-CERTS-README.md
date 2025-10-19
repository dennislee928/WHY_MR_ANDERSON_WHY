# mTLS 憑證生成問題修復

## 🔍 問題

在 Git Bash 中執行 `generate-certs.sh` 時遇到路徑轉換錯誤：

```bash
'/C=TW/ST=Taipei/...' → 'C:/Program Files/Git/C=TW/ST=Taipei/...'
```

---

## ✅ 解決方案

### 已應用修復

**所有 `-subj` 參數已修復**:
```bash
# 修復前（Git Bash 會轉換）
-subj "/C=TW/ST=Taipei/O=Pandora/CN=..."

# 修復後（Git Bash 相容）
-subj "//C=TW//ST=Taipei//O=Pandora//CN=..."
```

**修復數量**: 5 處

---

## 🚀 使用方式

### 方式 1: Git Bash（已修復）

```bash
cd ~/Documents/GitHub/Local_IPS-IDS

# 直接執行（已修復路徑問題）
./scripts/generate-certs.sh
```

### 方式 2: PowerShell（建議）

```powershell
cd C:\Users\dennis.lee\Documents\GitHub\Local_IPS-IDS

# 在 Git Bash 中執行
bash scripts/generate-certs.sh
```

### 方式 3: 使用環境變數

```bash
# 禁用路徑轉換
export MSYS_NO_PATHCONV=1

# 執行腳本
./scripts/generate-certs.sh
```

---

## 📁 生成的憑證

執行成功後會生成：

```
deployments/onpremise/certs/
├── ca/
│   ├── ca.key            (CA 私鑰)
│   └── ca.crt            (CA 證書)
├── device/
│   ├── server.key        (Device Service 私鑰)
│   ├── server.csr        (CSR)
│   └── server.crt        (證書)
├── network/
│   ├── server.key
│   └── server.crt
├── control/
│   ├── server.key
│   └── server.crt
└── client/
    ├── client.key
    └── client.crt
```

---

## ✅ 驗證

```bash
# 檢查 CA 證書
openssl x509 -in deployments/onpremise/certs/ca/ca.crt -noout -text

# 驗證證書鏈
openssl verify -CAfile deployments/onpremise/certs/ca/ca.crt \
  deployments/onpremise/certs/device/server.crt
```

---

**修復狀態**: ✅ 完成  
**可立即使用**: 是

