# Fly.io Volume 快速調整命令

您已經成功登入 Fly.io！現在可以直接執行以下命令來調整 volume：

## 🚀 方式 1: 重新部署（推薦）

在您當前的 PowerShell 終端機中執行：

```powershell
# 1. 檢查當前 volumes
flyctl volumes list --app pandora-monitoring

# 2. 重新部署應用（會使用新的 3GB volume 配置）
flyctl deploy --app pandora-monitoring --config deployments/paas/flyio/fly-monitoring.toml --dockerfile build/docker/monitoring.dockerfile

# 3. 檢查新的 volumes
flyctl volumes list --app pandora-monitoring

# 4. 刪除舊的大 volumes（重要！這樣才能停止計費）
# 記下舊 volume 的 ID，然後執行：
flyctl volumes delete <OLD_VOLUME_ID_1>
flyctl volumes delete <OLD_VOLUME_ID_2>
flyctl volumes delete <OLD_VOLUME_ID_3>
flyctl volumes delete <OLD_VOLUME_ID_4>
```

## 🔍 檢查應用狀態

```powershell
# 查看應用狀態
flyctl status --app pandora-monitoring

# 查看日誌
flyctl logs --app pandora-monitoring

# 查看 volumes
flyctl volumes list --app pandora-monitoring
```

## 💰 預期結果

- **之前**: 4 個 volumes，總計 18GB
- **之後**: 1 個 volume，3GB
- **節省**: 約 $15/月

## ⚠️ 重要提醒

1. 重新部署會有短暫停機時間（約 2-5 分鐘）
2. **必須手動刪除舊 volumes** 才能停止計費
3. Prometheus 和 Loki 的歷史數據會丟失（這通常是可接受的）
4. Grafana 儀表板設定會丟失，建議先匯出

## 🆘 如果遇到問題

### 問題 1: 部署失敗
```powershell
# 檢查錯誤日誌
flyctl logs --app pandora-monitoring

# 重試部署
flyctl deploy --app pandora-monitoring --config deployments/paas/flyio/fly-monitoring.toml --dockerfile build/docker/monitoring.dockerfile --remote-only
```

### 問題 2: Volume 刪除失敗
```powershell
# 先停止應用
flyctl apps stop pandora-monitoring

# 等待 30 秒後刪除
flyctl volumes delete <VOLUME_ID>

# 重新啟動應用
flyctl apps restart pandora-monitoring
```

### 問題 3: 需要回滾
```powershell
# 查看部署歷史
flyctl releases --app pandora-monitoring

# 回滾到上一個版本
flyctl releases rollback --app pandora-monitoring
```

## 📝 執行步驟檢查清單

- [ ] 已登入 Fly.io (`flyctl auth login`)
- [ ] 已檢查當前 volumes
- [ ] 已執行重新部署
- [ ] 已確認新 volume 創建成功
- [ ] 已刪除所有舊 volumes
- [ ] 已驗證應用正常運行
- [ ] 已確認費用降低

---

**立即在您的終端機中執行第一個命令開始！** 🚀
