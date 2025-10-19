# Fly.io Grafana 崩潰問題修復

## 🐛 問題描述

Grafana 不斷崩潰重啟，日誌顯示：

```
INFO success: grafana entered RUNNING state
WARN exited: grafana (exit status 1; not expected)
Error: ✗ failed to connect to database: mkdir /var/lib/grafana: file exists
```

## 🔍 根本原因

**符號連結衝突**：

1. Dockerfile 中建立了符號連結：
   ```dockerfile
   RUN ln -s /data/grafana /var/lib/grafana
   ```

2. Grafana 啟動時嘗試建立 `/var/lib/grafana` 目錄

3. 因為符號連結已存在（指向 `/data/grafana`），mkdir 失敗

4. Grafana 崩潰並不斷重啟

## ✅ 解決方案

### 修改 1: Dockerfile.monitoring

**移除 Grafana 的符號連結**：

```dockerfile
# 舊版本（有問題）
RUN ln -s /data/prometheus /prometheus && \
    ln -s /data/loki /loki && \
    ln -s /data/grafana /var/lib/grafana && \  # ← 問題所在
    ln -s /data/alertmanager /alertmanager

# 新版本（修復）
RUN ln -s /data/prometheus /prometheus && \
    ln -s /data/loki /loki && \
    ln -s /data/alertmanager /alertmanager
    # 移除 Grafana 符號連結，直接使用環境變數指定路徑
```

### 修改 2: configs/supervisord-flyio.conf

**使用環境變數指定 Grafana 資料路徑**：

```ini
# 舊版本
environment=GF_PATHS_DATA="/var/lib/grafana"

# 新版本
environment=GF_PATHS_DATA="/data/grafana"
```

同時確保 `/var/log/grafana` 目錄存在。

## 📊 修復後的目錄結構

```
/data/                          # Volume 掛載點
├── prometheus/                # Prometheus 數據
├── loki/                      # Loki 日誌
├── grafana/                   # Grafana 數據（直接使用，不用符號連結）
│   ├── grafana.db            # Grafana 資料庫
│   ├── plugins/              # 外掛
│   └── dashboards/           # 儀表板
└── alertmanager/             # AlertManager 數據

# 符號連結（僅用於其他服務）
/prometheus -> /data/prometheus
/loki -> /data/loki
/alertmanager -> /data/alertmanager

# Grafana 直接使用 /data/grafana（透過環境變數）
# 不使用符號連結
```

## 🚀 重新部署

修復已提交並推送到 main 分支：

```bash
git add Dockerfile.monitoring configs/supervisord-flyio.conf
git commit -m "Fix Grafana crash - resolve symlink conflict"
git push origin main
```

Fly.io 會自動檢測到更新並觸發重新部署。

## 🔍 驗證修復

等待重新部署完成後，檢查：

```bash
# 查看日誌
fly logs -a pandora-monitoring

# 應該看到所有服務正常運行：
# ✅ prometheus entered RUNNING state
# ✅ loki entered RUNNING state
# ✅ grafana entered RUNNING state (不再崩潰)
# ✅ alertmanager entered RUNNING state
# ✅ nginx entered RUNNING state
```

### 健康檢查

```bash
# Grafana
curl https://pandora-monitoring.fly.dev/api/health

# Prometheus
curl https://pandora-monitoring.fly.dev/prometheus/-/healthy

# Loki
curl https://pandora-monitoring.fly.dev/loki/ready

# AlertManager
curl https://pandora-monitoring.fly.dev/alertmanager/-/healthy
```

## 📚 學到的教訓

1. **符號連結與應用程式衝突**
   - 某些應用程式（如 Grafana）期望自己建立目錄
   - 符號連結可能會導致 "file exists" 錯誤
   - 解決：直接使用環境變數指定路徑

2. **環境變數的重要性**
   - Grafana 支援 `GF_PATHS_DATA` 環境變數
   - 比符號連結更靈活且不易出錯

3. **Supervisor 配置**
   - 環境變數可以在 Supervisor 配置中設定
   - 適合容器化部署

## 🎯 預期結果

修復後，Grafana 應該：
- ✅ 正常啟動
- ✅ 使用 `/data/grafana` 儲存數據
- ✅ 不再崩潰重啟
- ✅ 可以透過瀏覽器訪問
- ✅ 登入: admin / pandora123

---

**修復日期**: 2024-12-19  
**狀態**: ✅ 已修復並推送  
**等待**: Fly.io 自動重新部署

