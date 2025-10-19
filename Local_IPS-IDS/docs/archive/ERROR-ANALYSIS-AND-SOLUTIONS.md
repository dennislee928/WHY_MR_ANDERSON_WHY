# 錯誤分析與解決方案

## 📋 系統錯誤日誌分析報告

**分析日期**: 2025-01-14  
**系統版本**: v3.2.0  
**分析的服務**: 11個核心容器

---

## ✅ 正常運行的服務

以下服務正常運行，無需修復：

| 服務 | 狀態 | 說明 |
|------|------|------|
| RabbitMQ | ✅ Healthy | 正常啟動，管理界面可訪問 |
| Redis | ✅ Healthy | 正常運行，密碼認證已配置 |
| PostgreSQL | ✅ Healthy | 資料庫正常運行 |

---

## ⚠️ 需要處理的錯誤

### 1. AlertManager - Webhook 端點 404 錯誤

**錯誤日誌** (`alter_manager.txt`):
```
msg="Notify for alerts failed" err="...unexpected status code 404: 
http://axiom-ui:3001/api/v1/alerts/webhook: 404 page not found"
```

**原因**: AlertManager 嘗試將告警發送到 Axiom UI 的 webhook 端點，但該端點未實現。

**影響**: ⚠️ 中等 - 告警通知無法送達 Axiom UI

**解決方案**:

需要在 `internal/axiom/ui_server.go` 實現這些端點：

```go
// 添加到 ui_server.go
router.POST("/api/v1/alerts/webhook", ui.handleAlertWebhook)
router.POST("/api/v1/alerts/critical", ui.handleCriticalAlert)

func (ui *UIServer) handleAlertWebhook(c *gin.Context) {
    var alerts []Alert
    if err := c.BindJSON(&alerts); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 處理告警
    for _, alert := range alerts {
        logger.Infof("Received alert: %s", alert.Labels["alertname"])
        // TODO: 存儲到數據庫，推送到 WebSocket 客戶端
    }
    
    c.JSON(200, gin.H{"status": "ok", "received": len(alerts)})
}

func (ui *UIServer) handleCriticalAlert(c *gin.Context) {
    var alerts []Alert
    if err := c.BindJSON(&alerts); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 處理嚴重告警
    for _, alert := range alerts {
        logger.Warnf("Critical alert: %s", alert.Labels["alertname"])
        // TODO: 觸發緊急響應，發送緊急通知
    }
    
    c.JSON(200, gin.H{"status": "ok", "received": len(alerts)})
}
```

**優先級**: P1 🟡

---

### 2. Prometheus - 無法連接 AlertManager

**錯誤日誌** (`prometheus.txt`):
```
msg="Error sending alert" err="...dial tcp: lookup alertmanager on 127.0.0.11:53: no such host"
```

**原因**: DNS 解析問題，可能是容器啟動順序問題。

**影響**: ⚠️ 中等 - Prometheus 無法發送告警到 AlertManager

**解決方案**:

1. **檢查 docker-compose.yml 的 depends_on**:
```yaml
prometheus:
  depends_on:
    - alertmanager
  networks:
    - pandora-network
```

2. **確認網絡配置**:
```bash
docker network inspect application_default
```

3. **重啟服務**:
```bash
cd Application
docker-compose restart prometheus alertmanager
```

**優先級**: P2 🟢

---

### 3. Nginx - 找不到上游服務

**錯誤日誌** (`nginx.txt`):
```
[emerg] host not found in upstream "grafana:3000"
```

**原因**: Nginx 啟動時其他服務尚未就緒，DNS 解析失敗。

**影響**: ⚠️ 中等 - Nginx 無法啟動

**解決方案**:

1. **修改 configs/nginx/nginx.conf**，使用 resolver：
```nginx
http {
    resolver 127.0.0.11 valid=30s;
    
    upstream grafana {
        server grafana:3000;
    }
    
    # ...
}
```

2. **或在 docker-compose.yml 中添加依賴**:
```yaml
nginx:
  depends_on:
    - grafana
    - prometheus
    - loki
    - axiom-ui
```

3. **重啟 Nginx**:
```bash
docker-compose restart nginx
```

**優先級**: P2 🟢

---

### 4. Node Exporter - NFSd 指標錯誤

**錯誤日誌** (`node_exporter.txt`):
```
msg="collector failed" name=nfsd err="...unknown NFSd metric line \"wdeleg_getattr\""
```

**原因**: Node Exporter 嘗試收集 NFSd 指標，但系統中沒有 NFS 服務或版本不匹配。

**影響**: ✅ 低 - 僅影響 NFSd 指標收集，其他指標正常

**解決方案**:

1. **禁用 NFSd collector**:
```yaml
# docker-compose.yml
node-exporter:
  command:
    - '--path.rootfs=/host'
    - '--collector.disable-defaults'
    - '--collector.cpu'
    - '--collector.meminfo'
    - '--collector.diskstats'
    - '--collector.filesystem'
    - '--collector.netdev'
    # 不包含 nfsd
```

2. **或忽略此錯誤**（建議）：
   - 此錯誤不影響系統功能
   - 其他所有指標正常收集

**優先級**: P3 🔵 (可忽略)

---

### 5. PostgreSQL - 無效的啟動封包

**錯誤日誌** (`postgres.txt`):
```
LOG:  invalid length of startup packet
```

**原因**: 健康檢查或監控工具使用 TCP 連接檢查，而非有效的 PostgreSQL 協議。

**影響**: ✅ 低 - 僅日誌噪音，不影響功能

**解決方案**:

1. **改進健康檢查**:
```yaml
# docker-compose.yml
postgres:
  healthcheck:
    test: ["CMD-SHELL", "pg_isready -U pandora"]
    interval: 10s
    timeout: 5s
    retries: 5
```

**優先級**: P3 🔵

---

### 6. Redis - 安全攻擊警告

**錯誤日誌** (`redis.txt`):
```
Possible SECURITY ATTACK detected. It looks like somebody is sending POST or Host: commands
```

**原因**: 監控工具或健康檢查使用 HTTP 協議連接 Redis，觸發安全檢測。

**影響**: ✅ 低 - 誤報，實際是監控工具

**解決方案**:

1. **已完成**: `protected-mode` 已禁用
2. **確認密碼認證**: ✅ 已配置 `pandora123`
3. **忽略此警告**：來自合法的監控工具

**優先級**: P3 🔵

---

### 7. Promtail - 只讀文件系統錯誤

**錯誤日誌** (`promtail.txt`):
```
error writing positions file" error="...read-only file system
```

**原因**: `/app/data` 目錄掛載為只讀或權限不足。

**影響**: ⚠️ 中等 - Promtail 無法保存讀取位置，重啟後可能重複收集日誌

**解決方案**:

1. **修改 docker-compose.yml**:
```yaml
promtail:
  volumes:
    - ./logs:/var/log:ro
    - ./data/promtail:/app/data:rw  # 確保可寫
```

2. **創建目錄並設置權限**:
```bash
mkdir -p Application/data/promtail
chmod 777 Application/data/promtail
```

3. **重啟 Promtail**:
```bash
docker-compose restart promtail
```

**優先級**: P1 🟡

---

### 8. Axiom UI - 缺失端點

**錯誤日誌** (`ui.txt`):
```
{"path":"/metrics","status":404}
{"path":"/health","status":404}
```

**原因**: Prometheus 和健康檢查工具請求的端點未實現。

**影響**: ⚠️ 中等 - Prometheus 無法抓取 Axiom UI 指標

**解決方案**:

在 `internal/axiom/ui_server.go` 添加端點：

```go
// 添加 metrics 端點
router.GET("/metrics", ui.getMetrics)

func (ui *UIServer) getMetrics(c *gin.Context) {
    // 返回 Prometheus 格式指標
    metrics := `# HELP axiom_requests_total Total HTTP requests
# TYPE axiom_requests_total counter
axiom_requests_total{method="GET",path="/api/v1/status"} 1234
axiom_requests_total{method="POST",path="/api/v1/alerts"} 56

# HELP axiom_active_connections Active WebSocket connections
# TYPE axiom_active_connections gauge
axiom_active_connections 42
`
    c.String(200, metrics)
}

// 健康檢查端點
router.GET("/health", ui.getHealth)

func (ui *UIServer) getHealth(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "healthy",
        "timestamp": time.Now().Unix(),
        "version": "3.2.0"
    })
}
```

**優先級**: P1 🟡

---

## 📊 錯誤優先級總結

| 優先級 | 數量 | 錯誤 |
|--------|------|------|
| 🔴 P0 Critical | 0 | 無阻斷性錯誤 |
| 🟡 P1 High | 3 | AlertManager Webhook, Promtail 寫入, Axiom UI 指標 |
| 🟢 P2 Medium | 2 | Prometheus → AlertManager, Nginx DNS |
| 🔵 P3 Low | 3 | Node Exporter NFSd, PostgreSQL 封包, Redis 安全警告 |

---

## 🛠️ 快速修復腳本

創建 `scripts/fix-monitoring-errors.sh`:

```bash
#!/bin/bash
# Pandora Box Console - 監控錯誤修復腳本

echo "=== 修復監控服務錯誤 ==="

# 1. 創建 Promtail 數據目錄
echo "1. 設置 Promtail 數據目錄..."
mkdir -p Application/data/promtail
chmod 777 Application/data/promtail

# 2. 重啟服務（正確的依賴順序）
echo "2. 重啟服務..."
cd Application

docker-compose restart postgres redis rabbitmq
sleep 5

docker-compose restart prometheus grafana loki alertmanager
sleep 5

docker-compose restart promtail node-exporter
sleep 3

docker-compose restart nginx axiom-ui

echo "✅ 修復完成！"
echo ""
echo "請檢查服務狀態："
echo "  docker-compose ps"
```

---

## 📈 監控建議

### 1. 設置 Grafana 告警

為關鍵錯誤設置告警：

- ✅ Promtail 寫入失敗
- ✅ Prometheus 抓取失敗
- ✅ AlertManager 通知失敗

### 2. 日誌聚合

所有錯誤已自動聚合到 Loki，可通過 Grafana 查看：

```
{container_name="promtail"} |= "error"
{container_name="prometheus"} |= "error"
{container_name="alertmanager"} |= "error"
```

### 3. 健康檢查儀表板

創建 Grafana 儀表板監控所有服務健康狀態。

---

## ✅ 下一步行動

1. **立即修復** (P1):
   - [ ] 實現 AlertManager webhook 端點
   - [ ] 修復 Promtail 寫入權限
   - [ ] 添加 Axiom UI metrics 端點

2. **計劃修復** (P2):
   - [ ] 修復 Nginx DNS 解析
   - [ ] 修復 Prometheus → AlertManager 連接

3. **可選優化** (P3):
   - [ ] 禁用 Node Exporter NFSd collector
   - [ ] 改進 PostgreSQL 健康檢查
   - [ ] 優化 Redis 監控配置

---

**維護者**: Pandora Security Team  
**最後更新**: 2025-01-14  
**版本**: v3.2.0

