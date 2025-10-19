# Terraform 部署檢查清單

## 🎯 部署前檢查

### ✅ API Tokens 準備

- [ ] Railway API Token 和 Project ID
- [ ] Render API Key
- [ ] Koyeb API Token 和 Organization ID
- [ ] Patr.io API Token
- [ ] Fly.io API Token
- [ ] Grafana Admin Password 設定
- [ ] PostgreSQL Password 設定

### ✅ Docker Images 準備

- [ ] `pandora-agent` 映像已建置並推送
- [ ] `axiom-ui` 映像已建置並推送
- [ ] `monitoring` 映像已建置並推送到 Fly.io registry
- [ ] `nginx` Dockerfile 存在且正確

### ✅ 配置檔案檢查

- [ ] `terraform.tfvars` 已建立（從 `.example` 複製）
- [ ] 所有敏感資訊已填入
- [ ] 環境變數已設定（dev/staging/prod）
- [ ] Git repository URL 正確
- [ ] Branch 名稱正確

### ✅ 前置作業

- [ ] Terraform 已安裝（>= 1.0）
- [ ] Git 已配置
- [ ] Docker 已安裝（用於本地測試）
- [ ] 網路連線穩定

## 🚀 部署步驟

### 第一步：初始化

```bash
cd terraform
terraform init
```

**檢查**:
- [ ] Providers 下載成功
- [ ] 沒有錯誤訊息
- [ ] `.terraform` 目錄已建立

### 第二步：驗證

```bash
terraform validate
terraform fmt -check
```

**檢查**:
- [ ] 配置有效
- [ ] 格式正確

### 第三步：計劃

```bash
terraform plan -out=tfplan
```

**檢查**:
- [ ] 計劃輸出清楚
- [ ] 要建立的資源正確
- [ ] 沒有意外的銷毀操作
- [ ] 所有模組都被引用

### 第四步：應用

```bash
terraform apply tfplan
```

**檢查**:
- [ ] 應用成功
- [ ] 所有資源都已建立
- [ ] 沒有錯誤訊息

### 第五步：驗證

```bash
terraform output
```

**檢查**:
- [ ] 所有 URL 都有效
- [ ] 服務都可訪問
- [ ] 健康檢查通過

## 🔍 部署後驗證

### Fly.io 監控系統

```bash
# 檢查應用狀態
fly status --app $(terraform output -raw fly_app_name)

# 查看日誌
fly logs --app $(terraform output -raw fly_app_name)

# 測試服務
curl $(terraform output -raw grafana_url)/api/health
curl $(terraform output -raw prometheus_url)/-/healthy
curl $(terraform output -raw loki_url)/ready
```

**檢查**:
- [ ] Grafana 可訪問
- [ ] Prometheus 健康
- [ ] Loki 健康
- [ ] AlertManager 健康
- [ ] Volume 已掛載

### Railway PostgreSQL

```bash
# 使用 psql 連線測試
psql $(terraform output -raw database_url)
```

**檢查**:
- [ ] 資料庫可連線
- [ ] 資料庫已建立
- [ ] 使用者權限正確

### Render Services

**檢查**:
- [ ] Redis 運行中
- [ ] Nginx 運行中
- [ ] 健康檢查通過

### Koyeb Agent

**檢查**:
- [ ] 服務運行中
- [ ] 日誌正常
- [ ] API 端點可訪問

### Patr.io UI

**檢查**:
- [ ] UI 可訪問
- [ ] 可以登入
- [ ] 可以連接到後端 API

## 📊 監控設定

### Grafana

1. [ ] 登入 Grafana (admin / your-password)
2. [ ] 驗證 Datasources
   - [ ] Prometheus connected
   - [ ] Loki connected
3. [ ] 檢查 Dashboards
4. [ ] 設定告警規則

### Prometheus

1. [ ] 訪問 Prometheus UI
2. [ ] 檢查 Targets
   - [ ] 所有 targets 都是 UP
3. [ ] 執行測試查詢
4. [ ] 驗證資料收集

### Loki

1. [ ] 訪問 Loki API
2. [ ] 查詢測試日誌
3. [ ] 驗證日誌收集

## 🔐 安全檢查

- [ ] 所有密碼都是強密碼
- [ ] API tokens 已安全儲存
- [ ] `terraform.tfvars` 不在 Git 中
- [ ] GitHub Secrets 已設定
- [ ] Fly.io secrets 已設定
- [ ] 沒有硬編碼的敏感資訊
- [ ] HTTPS 已啟用
- [ ] 防火牆規則已設定（如適用）

## 🧪 測試清單

### 功能測試

- [ ] 用戶可以登入 UI
- [ ] Agent 可以連接到資料庫
- [ ] Agent 可以連接到 Redis
- [ ] 日誌被正確收集
- [ ] 指標被正確收集
- [ ] 告警可以正常發送
- [ ] Nginx 反向代理正常工作

### 效能測試

- [ ] 頁面載入時間 < 3秒
- [ ] API 回應時間 < 500ms
- [ ] 資料庫查詢效能正常
- [ ] 沒有記憶體洩漏

### 整合測試

- [ ] 所有服務可以互相通訊
- [ ] 資料流正常
- [ ] 監控告警正常
- [ ] 日誌聚合正常

## 📝 文件檢查

- [ ] 所有配置都已記錄
- [ ] API endpoints 已記錄
- [ ] 故障排除步驟已記錄
- [ ] 變更日誌已更新
- [ ] README 已更新

## 🔄 回滾計劃

如果部署失敗：

1. [ ] 保存錯誤日誌
2. [ ] 執行 `terraform destroy` （如需要）
3. [ ] 檢查 state 檔案
4. [ ] 恢復到上一個已知良好狀態
5. [ ] 分析失敗原因
6. [ ] 修復問題
7. [ ] 重新部署

## ✅ 完成

- [ ] 所有檢查項目都已完成
- [ ] 系統運行穩定
- [ ] 監控正常運作
- [ ] 團隊已通知
- [ ] 文件已更新
- [ ] 部署記錄已儲存

---

**部署日期**: _______________  
**部署者**: _______________  
**環境**: _______________  
**版本**: _______________  
**狀態**: _______________

