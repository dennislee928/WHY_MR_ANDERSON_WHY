# Koyeb 快速部署參考卡 - Pandora Agent

## 🚀 5 分鐘快速部署

### 步驟 1: 準備

```bash
# 確認檔案存在
ls Dockerfile.agent.koyeb
```

### 步驟 2: Koyeb Dashboard 設定

前往 https://app.koyeb.com

### 步驟 3: 關鍵配置（容易出錯的地方）

#### ✅ 正確配置

| 欄位 | 正確值 | ❌ 錯誤範例 |
|------|--------|-----------|
| **Dockerfile path** | `Dockerfile.agent.koyeb` | ~~`./Dockerfile.agent.koyeb`~~ |
| **Build context** | `/` | ~~`.`~~ 或 ~~留空~~ |
| **Port** | `8080` | ~~`80`~~ |
| **Region** | `fra` | 其他區域（可能有延遲） |

### 步驟 4: 必要環境變數

```env
# 基礎設定
LOG_LEVEL=info
GIN_MODE=release
PORT=8080

# 資料庫與快取（從 Railway 和 Render 取得）
DATABASE_URL=postgresql://...
REDIS_URL=redis://...

# 監控系統（從 Fly.io 取得）
PROMETHEUS_URL=https://...
LOKI_URL=https://...
GRAFANA_URL=https://...

# 安全設定（自己生成）
JWT_SECRET=<openssl rand -base64 48>
ENCRYPTION_KEY=<openssl rand -hex 32>
```

### 步驟 5: 驗證部署

```bash
# 替換為您的 Koyeb URL
KOYEB_URL="https://pandora-agent-xxx.koyeb.app"

# 健康檢查
curl $KOYEB_URL/health

# API 狀態
curl $KOYEB_URL/api/v1/status

# Metrics
curl $KOYEB_URL/metrics
```

## 🐛 常見錯誤與解決

### 錯誤 1: "no such file or directory"

```
error: failed to solve: failed to read dockerfile: 
open ./Dockerfile.agent.koyeb: no such file or directory
```

**原因**: Dockerfile path 欄位填寫錯誤

**解決**:
- ✅ 使用: `Dockerfile.agent.koyeb`
- ❌ 不要用: `./Dockerfile.agent.koyeb`

### 錯誤 2: "connection refused"

```
健康檢查失敗: connection refused
```

**原因**: PORT 環境變數與實際監聽端口不一致

**解決**:
```env
PORT=8080  # 必須與 Dockerfile 中的 EXPOSE 一致
```

### 錯誤 3: "database connection failed"

```
Error: database connection failed
```

**原因**: DATABASE_URL 格式錯誤或未設定

**解決**:
```bash
# 從 Railway 複製完整的 DATABASE_URL
DATABASE_URL=postgresql://postgres:password@host:5432/database
```

## 📊 資源限制

Koyeb 免費方案 (Nano):

- **CPU**: 0.1 vCPU
- **Memory**: 512 MB
- **Disk**: Ephemeral (不持久化)
- **Instances**: 2 個（永不休眠）
- **Bandwidth**: 無限制

**提示**: 
- Agent 預估記憶體使用: 150-200 MB
- 建議監控記憶體使用，避免 OOM

## 🔄 更新部署

### 方法 1: 自動部署（推薦）

```bash
# Git push 自動觸發重新部署
git add .
git commit -m "Update agent"
git push origin main
```

### 方法 2: 手動觸發

Koyeb Dashboard → Services → pandora-agent → "Redeploy"

### 方法 3: CLI

```bash
koyeb service redeploy pandora-agent/pandora-agent
```

## 💡 最佳實踐

1. **使用 Secret 儲存敏感資料**
   - Dashboard → Secrets → Add Secret
   - 在環境變數中引用 Secret

2. **設定告警**
   - 監控 Memory > 80%
   - 監控 CPU > 90%
   - 健康檢查失敗通知

3. **定期查看日誌**
   ```bash
   koyeb service logs pandora-agent/pandora-agent --follow
   ```

4. **使用多個 Instances**
   - 免費方案提供 2 個 Nano
   - 實現基本的高可用性

## 📝 檢查清單

部署前確認：

- [ ] `Dockerfile.agent.koyeb` 存在於 Repository 根目錄
- [ ] Dockerfile path 正確填寫（不含 `./` 前綴）
- [ ] Build context 設定為 `/`
- [ ] Port 設定為 `8080`
- [ ] 所有必要環境變數已設定
- [ ] 健康檢查路徑為 `/health`
- [ ] Region 選擇 `fra` (Frankfurt)

## 🆘 需要幫助？

1. **查看詳細指南**: [KOYEB-DEPLOYMENT-GUIDE.md](KOYEB-DEPLOYMENT-GUIDE.md)
2. **完整文件**: [README-PAAS-DEPLOYMENT.md](README-PAAS-DEPLOYMENT.md)
3. **Koyeb 官方文件**: https://www.koyeb.com/docs
4. **Koyeb Discord**: https://discord.gg/koyeb

---

**最後更新**: 2024-12-19
**適用版本**: Pandora Box Console IDS-IPS v1.0.0

