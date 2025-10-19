# Cloudflare Bindings 設定指南

## 🚀 **完整的 Cloudflare Bindings 配置**

本指南說明如何為安全平台專案設定所有可用的 Cloudflare bindings，以實現完整的功能。

## 📋 **Bindings 分類與用途**

### **1. 核心資料儲存** 📊
- **D1 Database**: 主要關聯式資料庫
- **KV Namespace (CACHE)**: 快取層
- **KV Namespace (SESSIONS)**: 會話管理

### **2. 檔案儲存與日誌** 📁
- **R2 Bucket (LOGS)**: 安全日誌儲存
- **R2 Bucket (THREAT_INTEL)**: 威脅情報檔案

### **3. AI 與機器學習** 🤖
- **Workers AI**: AI 威脅檢測
- **Vectorize Index**: 威脅情報向量搜尋

### **4. 訊息與佇列** 🔄
- **Queue (THREAT_PROCESSING)**: 威脅事件處理
- **Queue (LOG_PROCESSING)**: 日誌處理

### **5. 安全性與速率限制** 🛡️
- **Rate Limiter**: API 速率限制
- **Secrets Store**: 安全憑證儲存

### **6. 即時通訊** 🌐
- **Durable Objects (WEBSOCKET_MANAGER)**: WebSocket 管理
- **Durable Objects (THREAT_BROADCASTER)**: 威脅事件廣播

### **7. 進階功能** ⚡
- **Analytics Engine**: 安全事件分析
- **Hyperdrive**: 外部資料庫加速

## 🛠️ **設定步驟**

### **步驟 1: 建立 D1 Database**
```bash
# 建立資料庫
wrangler d1 create security_platform_db

# 執行 schema
wrangler d1 execute security_platform_db --file=schema.sql
```

### **步驟 2: 建立 KV Namespaces**
```bash
# 建立快取 namespace
wrangler kv:namespace create "CACHE"

# 建立會話 namespace
wrangler kv:namespace create "SESSIONS"
```

### **步驟 3: 建立 R2 Buckets**
```bash
# 建立日誌 bucket
wrangler r2 bucket create security-logs

# 建立威脅情報 bucket
wrangler r2 bucket create threat-intelligence
```

### **步驟 4: 建立 Queues**
```bash
# 建立威脅處理佇列
wrangler queues create threat-analysis-queue

# 建立日誌處理佇列
wrangler queues create log-processing-queue
```

### **步驟 5: 建立 Vectorize Index**
```bash
# 建立威脅情報向量索引
wrangler vectorize create threat-intelligence --dimensions=768
```

### **步驟 6: 建立 Hyperdrive**
```bash
# 建立 Hyperdrive 連線
wrangler hyperdrive create external-db --connection-string="postgresql://..."
```

## 📝 **更新 wrangler.toml**

將建立的資源 ID 填入 `wrangler.toml`：

```toml
# D1 Database
[[d1_databases]]
binding = "DB"
database_name = "security_platform_db"
database_id = "YOUR_D1_DATABASE_ID"

# KV Namespaces
[[kv_namespaces]]
binding = "CACHE"
id = "YOUR_CACHE_NAMESPACE_ID"

[[kv_namespaces]]
binding = "SESSIONS"
id = "YOUR_SESSIONS_NAMESPACE_ID"

# R2 Buckets
[[r2_buckets]]
binding = "LOGS"
bucket_name = "security-logs"

[[r2_buckets]]
binding = "THREAT_INTEL"
bucket_name = "threat-intelligence"

# Queues
[[queues]]
binding = "THREAT_PROCESSING"
queue_name = "threat-analysis-queue"

[[queues]]
binding = "LOG_PROCESSING"
queue_name = "log-processing-queue"

# Vectorize
[[vectorize]]
binding = "THREAT_VECTORS"
index_name = "threat-intelligence"

# Hyperdrive
[[hyperdrive]]
binding = "HYPERDRIVE"
hyperdrive_id = "YOUR_HYPERDRIVE_ID"
```

## 💰 **成本考量**

### **免費額度**
- **D1**: 5GB 儲存，25M 讀取，5M 寫入
- **KV**: 1GB 儲存，100K 讀取，100K 寫入
- **R2**: 10GB 儲存，1M 請求
- **Workers AI**: 10K 請求
- **Vectorize**: 30M 向量
- **Queue**: 100K 訊息

### **付費升級建議**
- 高流量生產環境建議升級到付費方案
- 監控使用量避免超出免費額度

## 🔧 **程式碼範例**

### **D1 Database 操作**
```javascript
// 儲存威脅事件
await env.DB.prepare(`
  INSERT INTO threats (id, severity, type, timestamp, data)
  VALUES (?, ?, ?, ?, ?)
`).bind(threatId, severity, type, timestamp, JSON.stringify(data)).run();

// 查詢威脅事件
const threats = await env.DB.prepare(`
  SELECT * FROM threats 
  WHERE severity = ? AND timestamp > ?
`).bind('high', Date.now() - 86400000).all();
```

### **KV Cache 操作**
```javascript
// 快取威脅情報
await env.CACHE.put(`threat:${threatId}`, JSON.stringify(threatData), {
  expirationTtl: 3600 // 1小時過期
});

// 取得快取資料
const cached = await env.CACHE.get(`threat:${threatId}`);
```

### **R2 儲存操作**
```javascript
// 儲存日誌檔案
await env.LOGS.put(`logs/${date}/security.log`, logData);

// 取得威脅情報檔案
const threatIntel = await env.THREAT_INTEL.get('ioc-list.json');
```

### **Workers AI 操作**
```javascript
// AI 威脅分析
const analysis = await env.AI.run("@cf/meta/llama-2-7b-chat-int8", {
  prompt: `分析以下安全事件：${JSON.stringify(eventData)}`
});
```

### **Vectorize 搜尋**
```javascript
// 向量化威脅情報
await env.THREAT_VECTORS.insert([
  {
    id: threatId,
    values: threatVector,
    metadata: { type: 'malware', severity: 'high' }
  }
]);

// 相似性搜尋
const matches = await env.THREAT_VECTORS.query(threatVector, {
  topK: 10,
  filter: { severity: 'high' }
});
```

### **Queue 操作**
```javascript
// 發送威脅事件到佇列
await env.THREAT_PROCESSING.send({
  threatId: '123',
  severity: 'high',
  timestamp: Date.now(),
  data: threatData
});

// 處理佇列訊息
export class ThreatProcessor {
  async consume(batch, env) {
    for (const message of batch.messages) {
      await this.processThreat(message.body);
    }
  }
}
```

### **Durable Objects WebSocket**
```javascript
export class WebSocketManager {
  async fetch(request) {
    const webSocketPair = new WebSocketPair();
    const [client, server] = Object.values(webSocketPair);
    
    server.accept();
    
    // 廣播威脅事件
    server.addEventListener('message', async (event) => {
      const data = JSON.parse(event.data);
      if (data.type === 'threat_alert') {
        await this.broadcastThreat(data);
      }
    });
    
    return new Response(null, {
      status: 101,
      webSocket: client,
    });
  }
}
```

## 🚀 **部署與測試**

### **部署**
```bash
# 部署到 Cloudflare
wrangler deploy

# 部署到特定環境
wrangler deploy --env production
```

### **測試 API**
```bash
# 測試健康檢查
curl https://security-platform-worker.workers.dev/api/v1/health

# 測試威脅檢測
curl -X POST https://security-platform-worker.workers.dev/api/v1/ml/detect \
  -H "Content-Type: application/json" \
  -d '{"event": "suspicious_login", "data": {...}}'
```

## 📊 **監控與分析**

### **Analytics Engine**
```javascript
// 記錄安全事件
await env.ANALYTICS.writeDataPoint({
  blobs: [eventType, severity],
  doubles: [timestamp, riskScore],
  indexes: [userId, deviceId]
});
```

### **Observability**
- 啟用 `[observability]` 設定
- 使用 Cloudflare Dashboard 監控
- 設定告警規則

## 🔒 **安全性最佳實踐**

1. **使用 Secrets Store** 儲存敏感資料
2. **實施 Rate Limiting** 防止濫用
3. **加密敏感資料** 在儲存前
4. **定期輪換 API 金鑰**
5. **監控異常活動**

## 📚 **相關資源**

- [Cloudflare Workers 文件](https://developers.cloudflare.com/workers/)
- [D1 Database 指南](https://developers.cloudflare.com/d1/)
- [Workers AI 範例](https://developers.cloudflare.com/workers-ai/)
- [Vectorize 文件](https://developers.cloudflare.com/vectorize/)

---

**注意**: 某些 bindings 可能需要付費方案才能使用。請檢查 Cloudflare 定價頁面了解詳細資訊。
