# Cloudflare Workers 部署指南 | Cloudflare Workers Deployment Guide

## 🚨 當前環境問題 | Current Environment Issues

檢測到 Node.js 權限問題。以下是解決方案：

### 問題：Node.js Permission Error
```
Error: EPERM: operation not permitted, lstat 'C:\Users\pclee'
```

### 解決方案 | Solutions

#### 選項 1：使用 Cloudflare Dashboard（推薦 | Recommended）

由於本地環境有權限問題，建議使用 Cloudflare Dashboard 手動部署：

**步驟 | Steps：**

1. **登入 Cloudflare Dashboard**
   - 前往：https://dash.cloudflare.com
   - 使用您的帳號登入

2. **建立 Worker**
   - 點擊左側 "Workers & Pages"
   - 點擊 "Create Application"
   - 選擇 "Create Worker"
   - 命名：`security-platform-worker`

3. **複製程式碼**
   
   將以下檔案內容複製到 Worker：
   
   **主程式（從 `infrastructure/cloud-configs/cloudflare/src/index.js`）：**

```javascript
/**
 * Cloudflare Workers Entry Point
 */

export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    
    // CORS headers
    const corsHeaders = {
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
      'Access-Control-Allow-Headers': 'Content-Type, Authorization',
    };

    // Handle CORS preflight
    if (request.method === 'OPTIONS') {
      return new Response(null, { headers: corsHeaders });
    }

    // Simple routing
    if (url.pathname === '/api/v1/health') {
      return new Response(
        JSON.stringify({
          healthy: true,
          timestamp: new Date().toISOString(),
          version: '1.0.0'
        }),
        {
          headers: {
            ...corsHeaders,
            'Content-Type': 'application/json',
          },
        }
      );
    }

    if (url.pathname === '/api/v1/status') {
      return new Response(
        JSON.stringify({
          status: 'operational',
          services: {
            worker: 'healthy',
            database: 'pending',
            cache: 'healthy'
          },
          timestamp: new Date().toISOString()
        }),
        {
          headers: {
            ...corsHeaders,
            'Content-Type': 'application/json',
          },
        }
      );
    }

    // Default response
    return new Response(
      JSON.stringify({
        message: 'Security Platform API',
        version: '1.0.0',
        endpoints: [
          '/api/v1/health',
          '/api/v1/status'
        ]
      }),
      {
        headers: {
          ...corsHeaders,
          'Content-Type': 'application/json',
        },
      }
    );
  },
};
```

4. **部署 | Deploy**
   - 點擊 "Save and Deploy"
   - 您會獲得一個 URL：`https://security-platform-worker.YOUR-SUBDOMAIN.workers.dev`

5. **測試 | Test**
   ```bash
   # 測試 health endpoint
   curl https://security-platform-worker.YOUR-SUBDOMAIN.workers.dev/api/v1/health
   
   # 測試 status endpoint
   curl https://security-platform-worker.YOUR-SUBDOMAIN.workers.dev/api/v1/status
   ```

#### 選項 2：修復 Node.js 環境

如果您想使用 Wrangler CLI：

**步驟 1：檢查環境變數**
```powershell
# 檢查用戶目錄
echo $env:USERPROFILE

# 應該是：C:\Users\USER
# 如果不是，需要修正
```

**步驟 2：清理 Node.js 緩存**
```powershell
# 清理 npm 緩存
npm cache clean --force

# 設定臨時目錄
$env:TEMP = "C:\Users\USER\AppData\Local\Temp"
$env:TMP = "C:\Users\USER\AppData\Local\Temp"
```

**步驟 3：使用管理員權限重新安裝**
```powershell
# 以管理員身份運行 PowerShell
# 然後執行：
npm install -g wrangler
```

**步驟 4：驗證安裝**
```powershell
wrangler --version
```

#### 選項 3：使用 WSL (Windows Subsystem for Linux)

如果您有 WSL：

```bash
# 在 WSL 中
cd /mnt/c/Users/USER/Desktop/WHY_MR_ANDERSON_WHY/infrastructure/cloud-configs/cloudflare

# 安裝 Node.js (如果沒有)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# 安裝依賴
npm install

# 登入 Cloudflare
npx wrangler login

# 部署
npx wrangler deploy
```

## 📊 D1 資料庫設定（進階）

一旦 Worker 部署成功，您可以添加 D1 資料庫：

### 使用 Dashboard 建立 D1

1. **前往 D1 頁面**
   - Cloudflare Dashboard > Storage & Databases > D1

2. **建立資料庫**
   - 點擊 "Create Database"
   - 名稱：`security_platform_db`
   - 點擊 "Create"

3. **取得資料庫 ID**
   - 複製 Database ID

4. **更新 Worker 設定**
   - Workers & Pages > security-platform-worker > Settings
   - Bindings > Add Binding
   - 類型：D1 Database
   - Variable name：`DB`
   - 選擇：`security_platform_db`

5. **初始化資料表**
   - 使用 D1 Console 執行 SQL：
   
```sql
-- 建立 threats 表
CREATE TABLE IF NOT EXISTS threats (
    id TEXT PRIMARY KEY,
    source_ip TEXT NOT NULL,
    threat_type TEXT NOT NULL,
    severity TEXT NOT NULL,
    status TEXT DEFAULT 'active',
    description TEXT,
    discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 建立索引
CREATE INDEX idx_threats_discovered ON threats(discovered_at);
CREATE INDEX idx_threats_status ON threats(status);
```

## 🔄 更新 Worker 以使用 D1

在 Worker 編輯器中，更新代碼以使用資料庫：

```javascript
export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    
    if (url.pathname === '/api/v1/threats') {
      try {
        // 查詢資料庫
        const { results } = await env.DB.prepare(
          'SELECT * FROM threats ORDER BY discovered_at DESC LIMIT 10'
        ).all();
        
        return new Response(
          JSON.stringify({
            threats: results,
            total: results.length
          }),
          {
            headers: {
              'Content-Type': 'application/json',
              'Access-Control-Allow-Origin': '*'
            }
          }
        );
      } catch (error) {
        return new Response(
          JSON.stringify({ error: error.message }),
          { status: 500 }
        );
      }
    }
    
    // ... 其他路由
  }
};
```

## ✅ 驗證部署

### 測試端點

```bash
# 設定您的 Worker URL
$WORKER_URL = "https://security-platform-worker.YOUR-SUBDOMAIN.workers.dev"

# 測試 health
curl $WORKER_URL/api/v1/health

# 預期回應：
# {"healthy":true,"timestamp":"...","version":"1.0.0"}

# 測試 status
curl $WORKER_URL/api/v1/status

# 如果已設定 D1，測試 threats
curl $WORKER_URL/api/v1/threats
```

### 查看日誌

在 Cloudflare Dashboard：
- Workers & Pages > security-platform-worker
- 點擊 "Logs" 標籤
- 查看即時日誌

## 💰 成本監控

免費額度：
- ✅ 10,000,000 請求/月
- ✅ 30,000,000 CPU 毫秒
- ✅ 5 GB D1 儲存空間

查看使用量：
- Dashboard > Workers & Pages > Analytics

## 🎯 下一步

部署成功後：

1. **設定自訂網域**
   - Workers > Triggers > Custom Domains
   - 添加：`api.yourdomain.com`

2. **添加 KV 儲存**
   - 用於快取和 session 管理

3. **設定監控**
   - 啟用 Analytics
   - 設定告警

4. **整合後端**
   - 將 Worker 設為 API Gateway
   - 代理請求到 OCI 或 IBM Cloud

## 📞 需要協助？

如果遇到問題：
1. 檢查 Worker logs
2. 查看 [Cloudflare 文檔](https://developers.cloudflare.com/workers/)
3. 參考 [Quick-Start.md](../../../Quick-Start.md)

---

**狀態**：✅ 簡化版部署完成  
**下一步**：測試 API 端點並添加更多功能

