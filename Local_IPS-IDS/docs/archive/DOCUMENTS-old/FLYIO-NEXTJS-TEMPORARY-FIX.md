# Fly.io Next.js 偵測問題 - 臨時解決方案

## 🚨 問題狀況

Fly.io 持續偵測到 Next.js 應用，即使使用 `.flyignore` 也無法避免。錯誤訊息：

```
Detected a Next.js app
Error: launch manifest was created for a app, but this is a Next.js app
```

## ✅ 臨時解決方案

### 已執行的操作

為了讓 Fly.io 成功部署監控系統，我們暫時重新命名了 Next.js 相關檔案：

```bash
# 重新命名 Next.js 檔案（避免偵測）
git mv package.json package.json.backup
git mv next.config.js next.config.js.backup
git mv tailwind.config.js tailwind.config.js.backup
git mv tsconfig.json tsconfig.json.backup
git mv vercel.json vercel.json.backup
```

### 檔案狀態

| 原始檔案 | 新名稱 | 狀態 |
|---------|--------|------|
| `package.json` | `package.json.backup` | 已重新命名 |
| `next.config.js` | `next.config.js.backup` | 已重新命名 |
| `tailwind.config.js` | `tailwind.config.js.backup` | 已重新命名 |
| `tsconfig.json` | `tsconfig.json.backup` | 已重新命名 |
| `vercel.json` | `vercel.json.backup` | 已重新命名 |

## 🎯 目的

這個臨時解決方案的目的是：

1. **讓 Fly.io 成功部署監控系統** (Prometheus + Loki + Grafana + AlertManager)
2. **避免 Next.js 自動偵測衝突**
3. **保持檔案完整性** (只是重新命名，不是刪除)

## 📋 部署後恢復步驟

### 1. 提交變更並推送到 main 分支

```bash
git add .
git commit -m "Temporary fix: Rename Next.js files to avoid Fly.io detection

- Rename package.json to package.json.backup
- Rename next.config.js to next.config.js.backup  
- Rename tailwind.config.js to tailwind.config.js.backup
- Rename tsconfig.json to tsconfig.json.backup
- Rename vercel.json to vercel.json.backup

This allows Fly.io to deploy monitoring system without Next.js conflicts."
git push origin main
```

### 2. 等待 Fly.io 部署成功

### 3. 恢復 Next.js 檔案（在 dev 分支）

```bash
# 切換到 dev 分支
git checkout dev

# 恢復檔案名稱
git mv package.json.backup package.json
git mv next.config.js.backup next.config.js
git mv tailwind.config.js.backup tailwind.config.js
git mv tsconfig.json.backup tsconfig.json
git mv vercel.json.backup vercel.json

# 提交恢復
git add .
git commit -m "Restore Next.js files after Fly.io deployment"
git push origin dev
```

## 🔄 長期解決方案

### 方案 A: 分離 Repository

考慮將專案分離為兩個 Repository：

1. **pandora-backend**: Go 後端 + 監控系統
2. **pandora-frontend**: Next.js 前端

### 方案 B: 使用不同的 Fly.io 應用

為不同服務建立不同的 Fly.io 應用：

```bash
# 監控系統
fly apps create pandora-monitoring

# 前端 (如果需要)
fly apps create pandora-frontend
```

### 方案 C: 使用 Docker 映像部署

先建置 Docker 映像，再部署到 Fly.io：

```bash
# 建置監控系統映像
docker build -f Dockerfile.monitoring -t pandora-monitoring:latest .

# 推送到 Docker Hub
docker push pandora-monitoring:latest

# 在 Fly.io 中使用預建映像
fly deploy --image pandora-monitoring:latest
```

## 📊 當前部署狀態

### 已完成的配置

- ✅ **Railway**: PostgreSQL 資料庫
- ✅ **Render**: Redis + Nginx 反向代理
- ✅ **Koyeb**: Pandora Agent + Promtail
- ✅ **Patr.io**: Axiom UI 前端
- 🔄 **Fly.io**: 監控系統 (進行中)

### 預期結果

Fly.io 部署成功後，您將獲得：

- **Grafana**: https://pandora-monitoring.fly.dev
- **Prometheus**: https://pandora-monitoring.fly.dev/prometheus
- **Loki**: https://pandora-monitoring.fly.dev/loki
- **AlertManager**: https://pandora-monitoring.fly.dev/alertmanager

## 🚨 注意事項

1. **這是臨時解決方案**: 檔案只是重新命名，不是刪除
2. **不影響其他平台**: Railway, Render, Koyeb, Patr.io 不受影響
3. **可完全恢復**: 所有檔案都可以恢復到原始狀態
4. **不影響開發**: dev 分支保持完整

## 📚 相關文件

- [Fly.io 部署指南](README-PAAS-DEPLOYMENT.md)
- [Fly.io 故障排除](FLYIO-TROUBLESHOOTING.md)
- [Next.js 衝突修復](FLYIO-NEXTJS-CONFLICT-FIX.md)

## 🆘 需要幫助？

如果遇到問題：

1. 檢查 Fly.io Dashboard 的建置日誌
2. 確認所有 Next.js 檔案已重新命名
3. 驗證 `Dockerfile.monitoring` 存在且正確
4. 檢查 `fly.toml` 配置

---

**執行日期**: 2024-12-19
**狀態**: 🔄 臨時解決方案已執行
**下一步**: 等待 Fly.io 部署成功後恢復檔案

