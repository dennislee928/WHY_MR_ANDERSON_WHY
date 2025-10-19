# Fly.io Volume 手動調整步驟

## 🎯 目標
- **之前**: 4 個獨立 volumes (總計 18GB)
- **現在**: 1 個 3GB volume
- **預估節省**: 大幅降低儲存費用（約節省 $15/月）

## 📋 前置要求

1. 確保已安裝 Fly.io CLI：
```powershell
# 方法 1: 使用 PowerShell 安裝腳本
iwr https://fly.io/install.ps1 -useb | iex

# 方法 2: 使用我們的簡化腳本
.\scripts\install-flyctl-simple.ps1
```

2. 重新啟動 PowerShell 終端機

3. 登入 Fly.io：
```powershell
flyctl auth login
```

## 🔧 手動調整步驟

### 步驟 1: 檢查當前 volumes

```powershell
flyctl volumes list --app pandora-monitoring
```

預期輸出類似：
```
ID                      NAME            SIZE    REGION  ZONE    ATTACHED VM     STATUS
vol_xxxxx1             grafana_data     5GB     hkg     xxxx    xxxxxxxxx       created
vol_xxxxx2             prometheus_data  5GB     hkg     xxxx    xxxxxxxxx       created
vol_xxxxx3             loki_data        5GB     hkg     xxxx    xxxxxxxxx       created
vol_xxxxx4             alertmanager     3GB     hkg     xxxx    xxxxxxxxx       created
```

### 步驟 2: 備份重要數據（可選但建議）

```powershell
# 使用 flyctl ssh 連接到應用
flyctl ssh console --app pandora-monitoring

# 在容器內備份
tar -czf /tmp/backup.tar.gz /data/grafana /data/prometheus /data/loki /data/alertmanager
exit
```

### 步驟 3: 停止應用

```powershell
flyctl apps stop pandora-monitoring
```

### 步驟 4: 刪除舊 volumes

```powershell
# 列出所有 volumes 並記錄 ID
flyctl volumes list --app pandora-monitoring

# 逐一刪除（替換 vol_xxxxx 為實際 ID）
flyctl volumes delete vol_xxxxx1 --yes
flyctl volumes delete vol_xxxxx2 --yes
flyctl volumes delete vol_xxxxx3 --yes
flyctl volumes delete vol_xxxxx4 --yes
```

### 步驟 5: 創建新的 3GB volume

```powershell
flyctl volumes create data --size 3 --region hkg --app pandora-monitoring
```

### 步驟 6: 更新 fly.toml 配置

確保 `deployments/paas/flyio/fly.toml` 包含：

```toml
[[mounts]]
  source = "data"
  destination = "/data"
  initial_size = "3gb"
```

### 步驟 7: 重新部署應用

```powershell
flyctl deploy --config deployments/paas/flyio/fly.toml --dockerfile build/docker/monitoring.dockerfile --remote-only --app pandora-monitoring
```

### 步驟 8: 驗證部署

```powershell
# 檢查應用狀態
flyctl status --app pandora-monitoring

# 檢查 volumes
flyctl volumes list --app pandora-monitoring

# 檢查日誌
flyctl logs --app pandora-monitoring
```

### 步驟 9: 測試服務

訪問以下 URL 確認服務正常：
- Grafana: https://pandora-monitoring.fly.dev:3000
- Prometheus: https://pandora-monitoring.fly.dev:9090
- Loki: https://pandora-monitoring.fly.dev:3100
- AlertManager: https://pandora-monitoring.fly.dev:9093

## 🔍 故障排除

### 問題 1: Volume 刪除失敗

```powershell
# 如果 volume 仍然附加到 VM，需要先分離
flyctl apps stop pandora-monitoring
# 等待 30 秒
flyctl volumes delete vol_xxxxx --yes
```

### 問題 2: 新 volume 未掛載

```powershell
# 檢查 fly.toml 配置
cat deployments/paas/flyio/fly.toml

# 確保只有一個 [[mounts]] 區塊
# 重新部署
flyctl deploy --config deployments/paas/flyio/fly.toml --remote-only
```

### 問題 3: 數據丟失

```powershell
# 如果需要恢復備份
flyctl ssh console --app pandora-monitoring
# 在容器內
tar -xzf /tmp/backup.tar.gz -C /
```

## 💰 費用節省計算

| 項目 | 之前 | 現在 | 節省 |
|------|------|------|------|
| Volume 數量 | 4 個 | 1 個 | -75% |
| 總容量 | 18 GB | 3 GB | -83% |
| 月費用 (估算) | ~$18 | ~$3 | ~$15/月 |

## ⚠️ 重要提醒

1. **數據丟失風險**: 刪除 volume 會永久刪除數據，請確保已備份重要數據
2. **停機時間**: 調整過程需要約 5-10 分鐘的停機時間
3. **監控數據**: Prometheus 和 Loki 的歷史數據會丟失，但這通常是可接受的
4. **Grafana 設定**: Grafana 的儀表板和設定會丟失，建議先匯出

## 📝 自動化腳本

如果您想使用自動化腳本，請在新的 PowerShell 終端機中執行：

```powershell
# 確保 flyctl 已安裝並在 PATH 中
flyctl version

# 執行調整腳本
.\scripts\flyio-volume-resize.ps1
```

## ✅ 完成檢查清單

- [ ] 已安裝並登入 Fly.io CLI
- [ ] 已備份重要數據（Grafana 儀表板）
- [ ] 已停止應用
- [ ] 已刪除所有舊 volumes
- [ ] 已創建新的 3GB volume
- [ ] 已更新 fly.toml 配置
- [ ] 已重新部署應用
- [ ] 已驗證所有服務正常運行
- [ ] 已確認費用降低

## 🔗 相關文檔

- [Fly.io Volumes 文檔](https://fly.io/docs/reference/volumes/)
- [Fly.io 定價](https://fly.io/docs/about/pricing/)
- [專案 Volume 修復指南](./FLYIO-VOLUME-FIX.md)
