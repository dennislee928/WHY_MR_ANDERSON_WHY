# Fly.io 部署故障排除指南

## 🐛 問題描述

Fly.io 部署失敗，錯誤訊息：

```
Error: failed loading app config from /usr/src/app/fly.toml: 
toml: table mounts already exists
```

## 🔍 問題分析

### 根本原因

`fly.toml` 配置檔案中使用了錯誤的 TOML 語法：

```toml
# ❌ 錯誤語法 - 多個 [mounts] 區段
[mounts]
  source = "prometheus_data"
  destination = "/prometheus"

[mounts]  # ← 這裡重複定義了 mounts 區段
  source = "loki_data"
  destination = "/loki"
```

### TOML 語法規則

在 TOML 中：
- `[section]` 定義單一區段
- `[[section]]` 定義陣列區段（可有多個）

## ✅ 解決方案

### 修復後的 fly.toml

```toml
# ✅ 正確語法 - 使用 [[mounts]] 陣列
[[mounts]]
  source = "prometheus_data"
  destination = "/prometheus"

[[mounts]]
  source = "loki_data"
  destination = "/loki"

[[mounts]]
  source = "grafana_data"
  destination = "/var/lib/grafana"

[[mounts]]
  source = "alertmanager_data"
  destination = "/alertmanager"
```

## 🚀 重新部署步驟

### 方法 1: 使用 Fly.io Dashboard

1. **前往 Fly.io Dashboard**
   - 登入 https://fly.io/dashboard
   - 選擇 `pandora-monitoring` 應用

2. **觸發重新部署**
   - 點擊 "Deploy" 按鈕
   - 或點擊 "Redeploy" 按鈕

3. **監控部署狀態**
   - 查看 Build Logs
   - 確認沒有 TOML 錯誤

### 方法 2: 使用 Fly CLI

```bash
# 1. 安裝 Fly CLI (如果尚未安裝)
curl -L https://fly.io/install.sh | sh

# 2. 登入
fly auth login

# 3. 驗證配置
fly config validate

# 4. 重新部署
fly deploy --config fly.toml --dockerfile Dockerfile.monitoring
```

### 方法 3: 透過 GitHub Actions

如果使用 GitHub Actions 自動部署：

```bash
# 推送修復到 main 分支
git add fly.toml
git commit -m "Fix fly.toml TOML syntax - use [[mounts]] instead of [mounts]"
git push origin main
```

## 🔧 其他常見 TOML 錯誤

### 1. 重複的區段名稱

```toml
# ❌ 錯誤
[env]
  LOG_LEVEL = "info"

[env]  # ← 重複定義
  TZ = "Asia/Taipei"

# ✅ 正確
[env]
  LOG_LEVEL = "info"
  TZ = "Asia/Taipei"
```

### 2. 錯誤的陣列語法

```toml
# ❌ 錯誤
[services]
  ports = [80, 443]

# ✅ 正確
[[services]]
  protocol = "tcp"
  internal_port = 80

  [[services.ports]]
    port = 80
    handlers = ["http"]

  [[services.ports]]
    port = 443
    handlers = ["tls", "http"]
```

### 3. 字串引號問題

```toml
# ❌ 錯誤
app = pandora-monitoring  # 需要引號

# ✅ 正確
app = "pandora-monitoring"
```

## 📋 完整的 fly.toml 範例

```toml
# Fly.io 監控系統配置檔案
app = "pandora-monitoring"
primary_region = "nrt"

[build]
  dockerfile = "Dockerfile.monitoring"

[env]
  LOG_LEVEL = "info"
  TZ = "Asia/Taipei"

[http_service]
  internal_port = 80
  force_https = true
  auto_stop_machines = false
  auto_start_machines = true
  min_machines_running = 1

  [http_service.concurrency]
    type = "connections"
    hard_limit = 1000
    soft_limit = 500

# 持久化儲存 - 使用 [[mounts]] 陣列語法
[[mounts]]
  source = "prometheus_data"
  destination = "/prometheus"

[[mounts]]
  source = "loki_data"
  destination = "/loki"

[[mounts]]
  source = "grafana_data"
  destination = "/var/lib/grafana"

[[mounts]]
  source = "alertmanager_data"
  destination = "/alertmanager"
```

## 🔍 驗證部署

部署成功後，驗證各服務：

```bash
# 設定 Fly.io 應用 URL
FLY_URL="https://pandora-monitoring.fly.dev"

# 1. 檢查 Prometheus
curl $FLY_URL/prometheus/-/healthy

# 2. 檢查 Loki
curl $FLY_URL/loki/ready

# 3. 檢查 Grafana
curl $FLY_URL/api/health

# 4. 檢查 AlertManager
curl $FLY_URL/alertmanager/-/healthy

# 5. 檢查整體健康狀態
curl $FLY_URL/health
```

## 🚨 其他可能的問題

### 1. Volume 不存在

如果遇到 volume 相關錯誤：

```bash
# 建立必要的 volumes
fly volumes create prometheus_data --size 3 --region nrt
fly volumes create loki_data --size 3 --region nrt
fly volumes create grafana_data --size 1 --region nrt
fly volumes create alertmanager_data --size 1 --region nrt
```

### 2. Dockerfile 問題

確認 `Dockerfile.monitoring` 存在且正確：

```bash
# 檢查 Dockerfile
ls -la Dockerfile.monitoring

# 本地測試建置
docker build -f Dockerfile.monitoring -t test-monitoring .
```

### 3. 環境變數問題

確保必要的環境變數已設定：

```bash
# 設定 Grafana 管理員密碼
fly secrets set GRAFANA_ADMIN_PASSWORD=pandora123

# 檢查 secrets
fly secrets list
```

## 📊 監控部署狀態

在 Fly.io Dashboard 中：

1. **Overview**: 查看應用整體狀態
2. **Logs**: 監控即時日誌
3. **Metrics**: 查看資源使用情況
4. **Volumes**: 確認持久化儲存狀態

## 💡 最佳實踐

1. **使用 fly config validate** 驗證配置
2. **本地測試 Dockerfile** 確保建置成功
3. **分階段部署** 先部署基本版本，再添加複雜功能
4. **監控日誌** 及時發現問題
5. **備份重要資料** 定期備份 volume 資料

## 🆘 需要幫助？

1. **Fly.io 官方文件**: https://fly.io/docs
2. **TOML 語法參考**: https://toml.io/en/
3. **Fly.io Discord**: https://fly.io/discord
4. **GitHub Issues**: 在專案中回報問題

---

**修復日期**: 2024-12-19
**狀態**: ✅ TOML 語法已修復
**下一步**: 重新部署到 Fly.io

