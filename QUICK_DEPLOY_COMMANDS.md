# ⚡ 快速部署命令 - Cloudflare Containers

## 🎯 **立即執行 (3 步驟)**

### **步驟 1: 獲取 API Token**
前往: https://dash.cloudflare.com/profile/api-tokens
點擊 "Create Token" → 選擇 "Edit Cloudflare Workers" → 創建並複製 Token

### **步驟 2: 設定 Token**
```powershell
# 在 PowerShell 中執行 (替換 YOUR_TOKEN 為實際 Token)
$env:CLOUDFLARE_API_TOKEN = "YOUR_TOKEN_HERE"

# 驗證設定
echo $env:CLOUDFLARE_API_TOKEN
```

### **步驟 3: 部署**
```powershell
# 進入專案目錄
cd C:\Users\USER\Desktop\WHY_MR_ANDERSON_WHY\spring-queen-719d

# 執行部署
npm run deploy
```

---

## 🔍 **如果需要 Account ID**

1. 前往: https://dash.cloudflare.com/
2. 右側邊欄會顯示 "Account ID"
3. 複製並設定:

```powershell
$env:CLOUDFLARE_ACCOUNT_ID = "YOUR_ACCOUNT_ID_HERE"
```

---

## ✅ **部署成功後**

您會看到:
```
✨ Successfully published your Worker to:
   https://spring-queen-719d.YOUR-SUBDOMAIN.workers.dev
```

測試:
```powershell
curl https://spring-queen-719d.YOUR-SUBDOMAIN.workers.dev
```

---

## 🚨 **如果遇到錯誤**

### 錯誤: "Authentication error"
→ 重新創建 API Token，確保使用 "Edit Cloudflare Workers" 模板

### 錯誤: "Missing account_id"  
→ 設定 Account ID:
```powershell
$env:CLOUDFLARE_ACCOUNT_ID = "your_account_id"
```

### 錯誤: "Worker name already exists"
→ 使用不同名稱或刪除現有 Worker

---

## 📞 **需要幫助？**

告訴我：
1. ✅ Token 已獲取並設定
2. ❌ 遇到特定錯誤 (告訴我錯誤訊息)
3. ❓ 不確定如何獲取 Token/Account ID
