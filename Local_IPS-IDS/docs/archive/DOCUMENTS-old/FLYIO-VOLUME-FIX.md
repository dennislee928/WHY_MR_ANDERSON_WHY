# Fly.io Volume 限制修復指南

## 🚨 問題

```
✖ Failed: error creating a new machine: failed to launch VM: 
invalid config.mounts, only 1 volume supported
```

**Fly.io 免費方案限制：每個應用只能有 1 個 Volume！**

## ✅ 解決方案

### 已完成的修改

1. **修改 fly.toml** - 使用單一 Volume
   ```toml
   [[mounts]]
     source = "monitoring_data"
     destination = "/data"
   ```

2. **修改 Dockerfile.monitoring** - 使用統一數據目錄
   - 所有數據儲存在 `/data/` 下
   - 使用符號連結保持相容性

### 需要執行的步驟

#### 步驟 1: 刪除現有的多個 Volumes

```bash
# 列出現有 volumes
fly volumes list --app pandora-monitoring

# 刪除所有舊 volumes
fly volumes delete prometheus_data --app pandora-monitoring --yes
fly volumes delete loki_data --app pandora-monitoring --yes
fly volumes delete grafana_data --app pandora-monitoring --yes
fly volumes delete alertmanager_data --app pandora-monitoring --yes
```

#### 步驟 2: 建立單一大容量 Volume

```bash
# 建立 10GB 的統一數據 volume
fly volumes create monitoring_data --size 10 --region nrt --app pandora-monitoring
```

#### 步驟 3: 提交變更

```bash
# 提交修改
git add fly.toml Dockerfile.monitoring
git commit -m "Fix Fly.io volume limitation - use single volume for all data"
git push origin main
```

#### 步驟 4: 重新部署

```bash
# 重新部署
fly deploy --app pandora-monitoring
```

## 📊 新的數據結構

```
/data/                      # 單一 Volume 掛載點
├── prometheus/            # Prometheus 數據
├── loki/                  # Loki 日誌
├── grafana/               # Grafana 儀表板和設定
└── alertmanager/          # AlertManager 數據

# 符號連結保持相容性
/prometheus -> /data/prometheus
/loki -> /data/loki
/var/lib/grafana -> /data/grafana
/alertmanager -> /data/alertmanager
```

## 💡 優勢

1. **符合免費方案限制** - 只使用 1 個 Volume
2. **更大容量** - 10GB vs 分散的小容量
3. **簡化管理** - 單一 Volume 更容易備份和管理
4. **成本效益** - 免費方案可用

## 🔍 驗證部署

部署成功後：

```bash
# 檢查應用狀態
fly status --app pandora-monitoring

# 檢查 Volume
fly volumes list --app pandora-monitoring

# 查看日誌
fly logs --app pandora-monitoring

# SSH 進入容器驗證
fly ssh console --app pandora-monitoring
ls -la /data/
```

## ⚠️ 注意事項

1. **數據會遺失** - 刪除舊 Volumes 會遺失現有數據（但目前還沒有重要數據）
2. **Volume 大小** - 可以根據需求調整（10GB, 20GB 等）
3. **區域選擇** - 使用 `nrt` (Tokyo) 以獲得更好的延遲

## 📋 快速執行腳本

```bash
#!/bin/bash
# 快速修復腳本

echo "🗑️  刪除舊 Volumes..."
fly volumes delete prometheus_data --app pandora-monitoring --yes 2>/dev/null
fly volumes delete loki_data --app pandora-monitoring --yes 2>/dev/null
fly volumes delete grafana_data --app pandora-monitoring --yes 2>/dev/null
fly volumes delete alertmanager_data --app pandora-monitoring --yes 2>/dev/null

echo "📦 建立新 Volume..."
fly volumes create monitoring_data --size 10 --region nrt --app pandora-monitoring

echo "🚀 重新部署..."
fly deploy --app pandora-monitoring

echo "✅ 完成！"
```

---

**修復日期**: 2024-12-19  
**狀態**: 配置已更新，等待重新部署  
**預期結果**: 部署成功，所有監控服務正常運行

