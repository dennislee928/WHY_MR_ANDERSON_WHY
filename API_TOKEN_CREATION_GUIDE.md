# 🔧 Cloudflare API Token 創建指南 - 完整步驟

## 🚨 **問題診斷**

您的 Token 返回 401 錯誤，可能原因：
1. ❌ Token 權限不足
2. ❌ Token 未正確創建
3. ❌ 使用了錯誤的 Token 類型

---

## ✅ **正確創建 API Token 的步驟**

### **步驟 1: 前往 API Tokens 頁面**

1. 打開瀏覽器
2. 前往: **https://dash.cloudflare.com/profile/api-tokens**
3. 確保您已登入

---

### **步驟 2: 創建新 Token**

1. 點擊右上角藍色按鈕 **"Create Token"**

2. **不要使用自訂模板！** 找到這個預設模板：
   ```
   📝 Edit Cloudflare Workers
   ```
   點擊右側的 **"Use template"** 按鈕

---

### **步驟 3: 確認權限設定**

創建頁面應該顯示：

```
Token name: Edit Cloudflare Workers

Permissions:
┌─────────────────────────────────────────────┐
│ Account                                     │
│ ├─ Workers Scripts ............... Edit    │
│ ├─ Workers KV Storage ............ Edit    │
│ └─ Workers Tail .................. Read    │
│                                             │
│ Zone                                        │
│ └─ Workers Routes ................ Edit    │
└─────────────────────────────────────────────┘

Account Resources:
✅ Include: [您的帳戶名稱]

Zone Resources:
✅ Include: All zones
```

**重要**: 
- ✅ 確保 **"Workers Scripts"** 有 **"Edit"** 權限
- ✅ 確保 **Account Resources** 已選擇您的帳戶
- ✅ **不要** 修改任何預設設定

---

### **步驟 4: 創建 Token**

1. 向下滾動到底部
2. 點擊 **"Continue to summary"**
3. 檢查摘要頁面
4. 點擊 **"Create Token"**

---

### **步驟 5: 複製 Token**

🚨 **超級重要！**

Token 會顯示在螢幕上，**只會顯示這一次**！

```
┌─────────────────────────────────────────────┐
│ ✅ Token created successfully               │
│                                             │
│ [Copy] xxxx-xxxxxxxxxxxxxxxxxxxxxxxx-xxxx  │
│                                             │
│ ⚠️  Make sure to copy your API token now.  │
│    You won't be able to see it again!      │
└─────────────────────────────────────────────┘
```

**立即複製這個 Token！**

---

### **步驟 6: 驗證 Token (可選)**

複製 Token 後，在 PowerShell 測試：

```powershell
# 設定新 Token
$token = "your_new_token_here"

# 驗證 Token
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

Invoke-RestMethod -Uri "https://api.cloudflare.com/client/v4/user/tokens/verify" -Headers $headers -Method GET
```

**成功的回應**:
```json
{
  "success": true,
  "result": {
    "id": "...",
    "status": "active"
  }
}
```

---

### **步驟 7: 使用新 Token 部署**

```powershell
# 設定新 Token
$env:CLOUDFLARE_API_TOKEN = "your_new_token_here"

# 進入專案目錄
cd C:\Users\USER\Desktop\WHY_MR_ANDERSON_WHY\spring-queen-719d

# 部署
npm run deploy
```

---

## 🎯 **快速檢查清單**

創建 Token 時，請確認：

- [ ] ✅ 使用了 **"Edit Cloudflare Workers"** 預設模板
- [ ] ✅ **沒有** 修改任何權限設定
- [ ] ✅ Account Resources 選擇了您的帳戶
- [ ] ✅ 成功點擊 "Create Token"
- [ ] ✅ 已複製完整的 Token（包含所有字元）
- [ ] ✅ Token 沒有空格或換行符

---

## 🔍 **如果找不到 "Edit Cloudflare Workers" 模板**

如果看不到這個模板，手動創建：

### **權限設定**:
```
Account Permissions:
├─ Workers Scripts ................. Edit
├─ Workers KV Storage .............. Edit
└─ Account Settings ................ Read

Zone Permissions:
└─ Workers Routes .................. Edit
```

### **Account Resources**:
```
✅ Include: [選擇您的帳戶]
```

### **Zone Resources**:
```
✅ Include: All zones
或
✅ Include: [選擇特定域名]
```

---

## ❓ **常見問題**

### **Q: Token 複製後有換行怎麼辦？**
A: 移除所有空格和換行，Token 應該是一串連續的字元

### **Q: 可以重新查看 Token 嗎？**
A: 不行，Token 只顯示一次。如果遺失，需要重新創建

### **Q: 一個帳戶可以有多個 Token 嗎？**
A: 可以，您可以創建多個 Token 用於不同用途

### **Q: Token 會過期嗎？**
A: 預設不會過期，除非您設定了 TTL

---

## 🚀 **完成創建後請執行**

```powershell
# 1. 設定新 Token
$env:CLOUDFLARE_API_TOKEN = "your_new_complete_token"

# 2. 驗證 Token（可選）
$headers = @{"Authorization" = "Bearer $env:CLOUDFLARE_API_TOKEN"}
Invoke-RestMethod -Uri "https://api.cloudflare.com/client/v4/user/tokens/verify" -Headers $headers

# 3. 如果驗證成功，執行部署
cd spring-queen-719d
npm run deploy
```

---

## 📸 **視覺化步驟**

```
Cloudflare Dashboard
    ↓
Profile (右上角頭像)
    ↓
API Tokens
    ↓
Create Token (藍色按鈕)
    ↓
Edit Cloudflare Workers (Use template)
    ↓
Continue to summary
    ↓
Create Token
    ↓
複製 Token ← 這裡！只有一次機會！
```

---

## 💡 **提示**

1. **不要分享 Token**: Token 等同於您的密碼
2. **安全儲存**: 將 Token 儲存在安全的地方（如密碼管理器）
3. **定期輪換**: 建議每 3-6 個月更換一次 Token
4. **刪除不用的**: 在 API Tokens 頁面可以刪除舊的 Token

---

**準備好了嗎？**

1. 前往: https://dash.cloudflare.com/profile/api-tokens
2. 創建 "Edit Cloudflare Workers" Token
3. 複製完整的 Token
4. 回來告訴我，我們繼續部署！

或者，如果您已經獲得新的 Token，請執行：
```powershell
$env:CLOUDFLARE_API_TOKEN = "your_new_token"
cd spring-queen-719d
npm run deploy
```
