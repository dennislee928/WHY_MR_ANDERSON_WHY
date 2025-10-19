# Fly.io Volume 大小調整指南

## 問題
Fly.io 的 pandora-monitoring 應用使用了 18GB 的 volume，產生額外費用。

## 解決方案

### 步驟 1: 安裝 Fly.io CLI

**選項 A: 使用 winget (推薦)**
```powershell
winget install Fly.Flyctl
```

**選項 B: 使用安裝腳本**
```powershell
.\scripts\install-flyctl.ps1
```

**選項 C: 手動下載**
- 下載：https://github.com/superfly/flyctl/releases
- 解壓並加入 PATH

### 步驟 2: 登入 Fly.io
```bash
flyctl auth login
```

### 步驟 3: 檢查當前 volumes
```bash
flyctl volumes list --app pandora-monitoring
```

### 步驟 4: 執行 volume 調整

**方式 1 (推薦): 重新部署**
```powershell
# 執行調整腳本
.\scripts\flyio-volume-resize.ps1

# 選擇 "1" 重新部署
```

**方式 2: 手動調整**
```bash
# 1. 創建新的 3GB volume
flyctl volumes create monitoring_data_new --app pandora-monitoring --region nrt --size 3

# 2. 重新部署使用新配置
flyctl deploy --app pandora-monitoring --config deployments/paas/flyio/fly-monitoring.toml

# 3. 刪除舊的大 volume
flyctl volumes list --app pandora-monitoring
flyctl volumes delete <OLD_VOLUME_ID>
```

## 已修改的配置

### fly-monitoring.toml
- **之前**: 4 個獨立 volumes (可能 18GB 總計)
- **現在**: 1 個合併的 3GB volume

```toml
[[mounts]]
  source = "monitoring_data"
  destination = "/data"
  size_gb = 3
```

## 費用節省
- **之前**: ~18GB × $0.15/GB/月 = ~$2.7/月
- **現在**: 3GB × $0.15/GB/月 = ~$0.45/月
- **節省**: ~$2.25/月 (~85% 減少)

## 注意事項
1. **重要**: 記得刪除舊的大 volume 才能停止計費
2. 調整過程可能需要短暫停機
3. 確保重要數據已備份

