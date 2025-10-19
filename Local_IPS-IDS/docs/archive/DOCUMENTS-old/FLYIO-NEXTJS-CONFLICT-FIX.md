# Fly.io Next.js 偵測衝突修復指南

## 🐛 問題描述

Fly.io 部署失敗，錯誤訊息：

```
Detected a Next.js app
Error: launch manifest was created for a app, but this is a Next.js app
unsuccessful command 'flyctl launch plan generate /tmp/manifest.json'
```

## 🔍 問題分析

### 根本原因

Fly.io 自動偵測到專案中的 Next.js 檔案：
- `package.json` (包含 Next.js 依賴)
- `next.config.js`
- `tailwind.config.js`
- `tsconfig.json`

但我們的 `fly.toml` 配置是為監控系統設計的，不是 Next.js 應用，導致配置衝突。

### 專案結構說明

這個專案包含：
- **Go 後端**: `cmd/`, `internal/` (要部署到 Koyeb)
- **Next.js 前端**: `package.json`, `next.config.js` (要部署到 Patr.io)
- **監控系統**: Prometheus + Loki + Grafana + AlertManager (要部署到 Fly.io)

## ✅ 解決方案

### 方案 1: 使用專用配置檔案（推薦）

使用 `fly-monitoring.toml` 而不是 `fly.toml`：

```bash
# 部署時指定配置檔案
fly deploy --config fly-monitoring.toml --dockerfile Dockerfile.monitoring
```

### 方案 2: 修改現有 fly.toml

已更新 `fly.toml` 加入明確的 builder 指定：

```toml
[build]
  builder = "dockerfile"  # 明確指定使用 Docker
  dockerfile = "Dockerfile.monitoring"
```

### 方案 3: 使用 .flyignore

已建立 `.flyignore` 檔案排除 Next.js 相關檔案：

```
# 排除 Next.js 檔案
package.json
next.config.js
tailwind.config.js
tsconfig.json
vercel.json
node_modules/
.next/
dist/
web/
```

## 🚀 部署步驟

### 方法 1: 使用 Fly.io Dashboard

1. **前往 Fly.io Dashboard**
   - 登入 https://fly.io/dashboard
   - 選擇或建立 `pandora-monitoring` 應用

2. **設定建置配置**
   - **Builder**: Docker
   - **Dockerfile**: `Dockerfile.monitoring`
   - **Build context**: `/`

3. **設定環境變數**
   ```
   LOG_LEVEL=info
   TZ=Asia/Taipei
   GRAFANA_ADMIN_PASSWORD=pandora123
   ```

4. **建立 Volumes**
   ```bash
   fly volumes create prometheus_data --size 3 --region nrt
   fly volumes create loki_data --size 3 --region nrt
   fly volumes create grafana_data --size 1 --region nrt
   fly volumes create alertmanager_data --size 1 --region nrt
   ```

5. **部署**
   - 點擊 "Deploy" 按鈕

### 方法 2: 使用 Fly CLI

```bash
# 1. 安裝 Fly CLI
curl -L https://fly.io/install.sh | sh

# 2. 登入
fly auth login

# 3. 建立應用
fly apps create pandora-monitoring --org personal

# 4. 建立 Volumes
fly volumes create prometheus_data --size 3 --region nrt
fly volumes create loki_data --size 3 --region nrt
fly volumes create grafana_data --size 1 --region nrt
fly volumes create alertmanager_data --size 1 --region nrt

# 5. 設定環境變數
fly secrets set GRAFANA_ADMIN_PASSWORD=pandora123

# 6. 部署（使用專用配置）
fly deploy --config fly-monitoring.toml --dockerfile Dockerfile.monitoring
```

### 方法 3: 使用 GitHub Actions

更新 `.github/workflows/deploy-paas.yml` 中的 Fly.io 部署步驟：

```yaml
- name: Deploy to Fly.io
  env:
    FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
  run: |
    flyctl deploy --config fly-monitoring.toml --dockerfile Dockerfile.monitoring --remote-only
```

## 🔧 驗證部署

部署成功後，驗證各服務：

```bash
# 設定 Fly.io 應用 URL
FLY_URL="https://pandora-monitoring.fly.dev"

# 1. 檢查整體健康狀態
curl $FLY_URL/health

# 2. 檢查 Prometheus
curl $FLY_URL/prometheus/-/healthy

# 3. 檢查 Loki
curl $FLY_URL/loki/ready

# 4. 檢查 Grafana
curl $FLY_URL/api/health

# 5. 檢查 AlertManager
curl $FLY_URL/alertmanager/-/healthy
```

## 📊 預期結果

部署成功後，您將看到：

- ✅ **Grafana**: https://pandora-monitoring.fly.dev
  - 用戶: admin
  - 密碼: pandora123

- ✅ **Prometheus**: https://pandora-monitoring.fly.dev/prometheus

- ✅ **Loki**: https://pandora-monitoring.fly.dev/loki

- ✅ **AlertManager**: https://pandora-monitoring.fly.dev/alertmanager

## 🚨 常見問題

### 問題 1: 仍然偵測到 Next.js

**解決方案**:
```bash
# 使用 --no-buildpacks 強制使用 Docker
fly deploy --config fly-monitoring.toml --dockerfile Dockerfile.monitoring --no-buildpacks
```

### 問題 2: Volume 不存在

**解決方案**:
```bash
# 建立所有必要的 volumes
fly volumes create prometheus_data --size 3 --region nrt
fly volumes create loki_data --size 3 --region nrt
fly volumes create grafana_data --size 1 --region nrt
fly volumes create alertmanager_data --size 1 --region nrt
```

### 問題 3: 記憶體不足

**解決方案**:
```toml
# 在 fly-monitoring.toml 中調整資源
[vm]
  memory_mb = 2048  # 增加到 2GB
```

## 💡 最佳實踐

1. **使用專用配置檔案**: `fly-monitoring.toml` 避免與 Next.js 衝突
2. **明確指定 Builder**: 使用 `builder = "dockerfile"`
3. **使用 .flyignore**: 排除不需要的檔案
4. **分離部署**: 不同服務使用不同的 Fly.io 應用
5. **監控資源**: 定期檢查 CPU 和記憶體使用

## 📚 相關文件

- [Fly.io Docker 部署指南](https://fly.io/docs/languages-and-frameworks/dockerfile/)
- [Fly.io 配置參考](https://fly.io/docs/reference/configuration/)
- [Fly.io Volumes 文件](https://fly.io/docs/reference/volumes/)

## 🆘 需要幫助？

1. **Fly.io 官方文件**: https://fly.io/docs
2. **Fly.io Discord**: https://fly.io/discord
3. **GitHub Issues**: 在專案中回報問題

---

**修復日期**: 2024-12-19
**狀態**: ✅ 提供多種解決方案
**建議**: 使用 `fly-monitoring.toml` 配置檔案

