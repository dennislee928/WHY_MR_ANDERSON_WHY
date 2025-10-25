# 🔑 Cloudflare Wrangler 登入指南 - API Token 方式

## 🚨 **問題**: `wrangler login` 無法打開瀏覽器

由於無法自動打開瀏覽器，我們使用 **API Token** 方式登入。

---

## 📋 **步驟 1: 獲取 Cloudflare API Token**

### **方法 A: 透過 Dashboard 創建 (推薦)**

1. **登入 Cloudflare Dashboard**
   - 前往: https://dash.cloudflare.com/
   - 登入您的帳戶

2. **進入 API Tokens 頁面**
   - 點擊右上角的個人頭像
   - 選擇 **"My Profile"**
   - 點擊左側選單的 **"API Tokens"**
   - 或直接前往: https://dash.cloudflare.com/profile/api-tokens

3. **創建新 Token**
   - 點擊 **"Create Token"** 按鈕
   - 找到 **"Edit Cloudflare Workers"** 模板
   - 點擊 **"Use template"**

4. **配置 Token 權限**
   ```
   Token name: Wrangler Deploy Token
   
   Permissions:
   ✅ Account - Workers Scripts - Edit
   ✅ Account - Workers KV Storage - Edit (可選)
   ✅ Account - Workers Tail - Read (可選)
   ✅ Zone - Workers Routes - Edit (可選)
   
   Account Resources:
   ✅ Include - Your Account
   
   Zone Resources:
   ✅ Include - All zones (或選擇特定域名)
   
   Client IP Address Filtering:
   留空 (允許所有 IP)
   
   TTL:
   選擇 "Custom" → 設定有效期限 (建議 1 年)
   ```

5. **生成並複製 Token**
   - 點擊 **"Continue to summary"**
   - 點擊 **"Create Token"**
   - **重要**: 複製顯示的 Token (只會顯示一次！)
   - Token 格式類似: `xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`

---

## 📋 **步驟 2: 使用 API Token 登入**

### **選項 A: 設定環境變數 (推薦)**

在 PowerShell 中執行:

```powershell
# 設定 Cloudflare API Token
$env:CLOUDFLARE_API_TOKEN = "your_token_here"

# 驗證設定
echo $env:CLOUDFLARE_API_TOKEN

# 現在可以部署了
cd spring-queen-719d
npm run deploy
```

### **選項 B: 在專案中設定 (永久)**

1. **創建 `.env` 檔案** (在 `spring-queen-719d` 目錄)
   ```env
   CLOUDFLARE_API_TOKEN=your_token_here
   CLOUDFLARE_ACCOUNT_ID=your_account_id_here
   ```

2. **獲取 Account ID**
   - 前往: https://dash.cloudflare.com/
   - 在右側邊欄找到 **"Account ID"**
   - 或在任何 Worker 頁面的 URL 中找到

3. **部署**
   ```powershell
   cd spring-queen-719d
   npm run deploy
   ```

### **選項 C: 使用 wrangler config (舊方法)**

```powershell
# 手動配置
npx wrangler config

# 按提示輸入 API Token
```

---

## 📋 **步驟 3: 完成部署**

```powershell
# 確保在正確目錄
cd C:\Users\USER\Desktop\WHY_MR_ANDERSON_WHY\spring-queen-719d

# 方法 1: 使用環境變數
$env:CLOUDFLARE_API_TOKEN = "your_token_here"
npm run deploy

# 或方法 2: 如果已設定 .env
npm run deploy

# 或方法 3: 直接使用 wrangler
npx wrangler deploy --env production
```

---

## ✅ **驗證部署成功**

部署成功後，您會看到類似輸出:

```
✨ Successfully published your Worker to:
   https://spring-queen-719d.YOUR-SUBDOMAIN.workers.dev
   
✨ Container deployed successfully!
```

**測試 Worker**:
```powershell
# 測試健康檢查
curl https://spring-queen-719d.YOUR-SUBDOMAIN.workers.dev/health

# 或在瀏覽器打開
start https://spring-queen-719d.YOUR-SUBDOMAIN.workers.dev
```

---

## 🔧 **快速命令參考**

```powershell
# === 設定 API Token ===
$env:CLOUDFLARE_API_TOKEN = "your_token_here"

# === 部署 ===
cd spring-queen-719d
npm run deploy

# === 查看部署 ===
npx wrangler deployments list

# === 查看即時日誌 ===
npx wrangler tail

# === 刪除 Worker (如果需要) ===
npx wrangler delete

# === 查看 Worker 資訊 ===
npx wrangler whoami
```

---

## 🚨 **常見問題**

### **問題 1: Token 權限不足**
```
Error: Authentication error
```

**解決方案**:
- 重新創建 Token，確保勾選所有必要權限
- 使用 "Edit Cloudflare Workers" 模板

### **問題 2: Account ID 找不到**
```
Error: Missing account_id
```

**解決方案**:
1. 前往 https://dash.cloudflare.com/
2. 右側邊欄會顯示 "Account ID"
3. 或在 `.env` 中設定:
   ```env
   CLOUDFLARE_ACCOUNT_ID=your_account_id
   ```

### **問題 3: Token 過期**
```
Error: Invalid token
```

**解決方案**:
- 在 Dashboard 創建新 Token
- 更新環境變數

---

## 📝 **完整部署流程**

```powershell
# 1. 設定 Token (選擇其中一個方法)
# 方法 A: 環境變數
$env:CLOUDFLARE_API_TOKEN = "your_token_here"

# 方法 B: 創建 .env 檔案
cd spring-queen-719d
echo "CLOUDFLARE_API_TOKEN=your_token_here" > .env
echo "CLOUDFLARE_ACCOUNT_ID=your_account_id" >> .env

# 2. 部署
npm run deploy

# 3. 測試
curl https://spring-queen-719d.YOUR-SUBDOMAIN.workers.dev

# 4. 查看日誌
npx wrangler tail
```

---

## 🎯 **下一步**

1. **立即獲取 API Token**:
   - 前往: https://dash.cloudflare.com/profile/api-tokens
   - 創建 "Edit Cloudflare Workers" Token

2. **設定環境變數**:
   ```powershell
   $env:CLOUDFLARE_API_TOKEN = "your_token_here"
   ```

3. **執行部署**:
   ```powershell
   cd spring-queen-719d
   npm run deploy
   ```

---

**準備好了嗎？請先獲取 API Token，然後我們繼續部署！** 🚀

**需要協助**:
- 📝 如何找到 Account ID
- 🔑 如何創建 API Token
- 🚀 部署過程中的任何問題
