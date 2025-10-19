# Pandora Box Console - Terraform 部署指南

完整的基礎設施即程式碼 (Infrastructure as Code) 解決方案，用於部署 Pandora Box Console IDS-IPS 系統到多個 PaaS 平台。

## 📋 目錄

- [架構概覽](#架構概覽)
- [前置需求](#前置需求)
- [快速開始](#快速開始)
- [模組說明](#模組說明)
- [環境管理](#環境管理)
- [CI/CD 整合](#cicd-整合)
- [故障排除](#故障排除)

## 🏗️ 架構概覽

```
┌─────────────────────────────────────────────────────────────┐
│                   Pandora Box Console                        │
│                Multi-Platform PaaS Deployment                │
└─────────────────────────────────────────────────────────────┘

Railway.app                Render.com              Koyeb
┌──────────────┐          ┌──────────────┐       ┌──────────────┐
│ PostgreSQL   │◄─────────│ Redis        │◄──────│ Pandora Agent│
│ Database     │          │ + Nginx      │       │ + Promtail   │
└──────────────┘          └──────────────┘       └──────────────┘
                                                         ▲
                                                         │
Patr.io                   Fly.io                        │
┌──────────────┐          ┌──────────────────────────────┐
│ Axiom UI     │◄─────────│ Monitoring Stack         │───┘
│ (Next.js)    │          │ • Prometheus             │
└──────────────┘          │ • Loki                   │
                          │ • Grafana                │
                          │ • AlertManager           │
                          └──────────────────────────────┘
```

## 🎯 部署平台分配

| 平台 | 服務 | 原因 |
|------|------|------|
| **Railway** | PostgreSQL | 託管資料庫，自動備份 |
| **Render** | Redis + Nginx | 快取和反向代理 |
| **Koyeb** | Agent + Promtail | 低延遲，適合實時處理 |
| **Patr.io** | Axiom UI | Next.js 部署最佳化 |
| **Fly.io** | 監控系統 | 全球分佈，持久化儲存 |

## ✅ 前置需求

### 1. 安裝必要工具

```bash
# Terraform
wget https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip
unzip terraform_1.6.0_linux_amd64.zip
sudo mv terraform /usr/local/bin/

# 驗證安裝
terraform version
```

### 2. 準備 API Tokens

您需要從以下平台獲取 API tokens：

1. **Railway**: https://railway.app/account/tokens
2. **Render**: https://dashboard.render.com/account/api-keys
3. **Koyeb**: https://app.koyeb.com/account/api
4. **Patr.io**: https://patr.cloud/dashboard/settings/tokens
5. **Fly.io**: `fly auth token`

### 3. 準備 Git Repository

確保您的 repository 包含所有必要的 Dockerfile：
- `Dockerfile.monitoring`
- `Dockerfile.agent.koyeb`
- `Dockerfile.ui.patr`
- `Dockerfile.nginx`

## 🚀 快速開始

### 1. 初始化配置

```bash
# 進入 terraform 目錄
cd terraform

# 複製變數範本
cp terraform.tfvars.example terraform.tfvars

# 編輯 terraform.tfvars，填入您的 API tokens
vim terraform.tfvars
```

### 2. 初始化 Terraform

```bash
# 初始化 Terraform（下載 providers）
terraform init

# 驗證配置
terraform validate

# 格式化配置檔案
terraform fmt -recursive
```

### 3. 查看計劃

```bash
# 查看將要建立的資源
terraform plan

# 將計劃儲存到檔案
terraform plan -out=tfplan
```

### 4. 應用配置

```bash
# 應用配置（需要確認）
terraform apply

# 或自動確認
terraform apply -auto-approve

# 或使用儲存的計劃
terraform apply tfplan
```

### 5. 查看輸出

```bash
# 查看所有輸出
terraform output

# 查看特定輸出
terraform output service_urls

# 以 JSON 格式輸出
terraform output -json
```

## 📦 模組說明

### Fly.io 監控系統模組

```hcl
module "flyio_monitoring" {
  source = "./modules/flyio"

  project_name           = "pandora"
  environment            = "prod"
  organization           = "personal"
  region                 = "nrt"
  volume_size            = 10
  grafana_admin_password = "your-password"
}
```

**資源**:
- Fly.io Application
- Persistent Volume (10GB)
- Machine with health checks
- Secrets management

### Railway PostgreSQL 模組

```hcl
module "railway_postgres" {
  source = "./modules/railway"

  project_id    = "your-project-id"
  api_token     = "your-api-token"
  database_name = "pandora"
}
```

**資源**:
- PostgreSQL 15 database
- Automatic backups
- Connection string

### Render Services 模組

```hcl
module "render_services" {
  source = "./modules/render"

  api_key        = "your-api-key"
  redis_name     = "pandora-redis"
  nginx_name     = "pandora-nginx"
  repository_url = "https://github.com/your-org/repo"
}
```

**資源**:
- Redis instance
- Nginx reverse proxy
- GitHub integration

### Koyeb Agent 模組

```hcl
module "koyeb_agent" {
  source = "./modules/koyeb"

  api_token    = "your-api-token"
  app_name     = "pandora-agent"
  docker_image = "your-registry/agent:latest"
  database_url = module.railway_postgres.database_url
}
```

**資源**:
- Koyeb service
- Auto-scaling configuration
- Health checks

### Patr.io UI 模組

```hcl
module "patr_ui" {
  source = "./modules/patr"

  api_token    = "your-api-token"
  app_name     = "axiom-ui"
  docker_image = "your-registry/ui:latest"
  api_url      = module.koyeb_agent.api_url
}
```

**資源**:
- Next.js application
- Environment variables
- Health checks

## 🌍 環境管理

### 開發環境

```bash
# 使用開發環境變數
terraform apply -var-file="environments/dev/terraform.tfvars"

# 或設定 workspace
terraform workspace new dev
terraform workspace select dev
terraform apply
```

### 生產環境

```bash
# 使用生產環境變數
terraform apply -var-file="environments/prod/terraform.tfvars"

# 或使用 workspace
terraform workspace new prod
terraform workspace select prod
terraform apply
```

### Workspace 命令

```bash
# 列出所有 workspaces
terraform workspace list

# 建立新 workspace
terraform workspace new staging

# 切換 workspace
terraform workspace select prod

# 顯示當前 workspace
terraform workspace show

# 刪除 workspace
terraform workspace delete staging
```

## 🔄 CI/CD 整合

### GitHub Actions

已包含 `.github/workflows/terraform-deploy.yml`，自動執行：

1. **Pull Request**: 
   - Terraform format check
   - Terraform validate
   - Terraform plan
   - 將計劃結果評論到 PR

2. **Push to main/dev**:
   - 自動 apply 變更
   - 上傳 state 檔案
   - 輸出部署摘要

### 設定 GitHub Secrets

在 GitHub Repository Settings → Secrets 中新增：

```
RAILWAY_PROJECT_ID
RAILWAY_API_TOKEN
POSTGRES_PASSWORD
RENDER_API_KEY
KOYEB_API_TOKEN
KOYEB_ORG_ID
PATR_API_TOKEN
FLY_API_TOKEN
GRAFANA_PASSWORD
```

## 🔧 常用命令

```bash
# 查看資源
terraform show

# 列出所有資源
terraform state list

# 查看特定資源
terraform state show module.flyio_monitoring.fly_app.monitoring

# 導入現有資源
terraform import module.flyio_monitoring.fly_app.monitoring pandora-monitoring

# 移除資源（不刪除實際資源）
terraform state rm module.flyio_monitoring.fly_app.monitoring

# 重新整理 state
terraform refresh

# 驗證並格式化
terraform fmt -recursive
terraform validate
```

## 🧹 清理資源

```bash
# 銷毀所有資源（需要確認）
terraform destroy

# 銷毀特定模組
terraform destroy -target=module.flyio_monitoring

# 自動確認銷毀
terraform destroy -auto-approve
```

## ❗ 故障排除

### 1. Provider 下載失敗

```bash
# 清理 lock 檔案
rm .terraform.lock.hcl

# 重新初始化
terraform init -upgrade
```

### 2. State Lock

```bash
# 強制解鎖（謹慎使用）
terraform force-unlock LOCK_ID
```

### 3. API Token 問題

```bash
# 驗證環境變數
echo $FLY_API_TOKEN

# 重新設定
export FLY_API_TOKEN="your-new-token"
```

### 4. 模組錯誤

```bash
# 更新模組
terraform get -update

# 清理並重新初始化
rm -rf .terraform
terraform init
```

## 📊 監控部署

### 查看 Fly.io 狀態

```bash
fly status --app pandora-monitoring-prod
fly logs --app pandora-monitoring-prod
```

### 驗證服務

```bash
# Grafana
curl https://pandora-monitoring-prod.fly.dev/api/health

# Prometheus
curl https://pandora-monitoring-prod.fly.dev/prometheus/-/healthy

# Loki
curl https://pandora-monitoring-prod.fly.dev/loki/ready
```

## 🔐 最佳實踐

1. **永遠不要提交 `terraform.tfvars`** - 包含敏感資訊
2. **使用 Remote State** - 團隊協作時使用 S3 或 Terraform Cloud
3. **定期備份 State** - `terraform.tfstate` 很重要
4. **使用 Workspaces** - 分離不同環境
5. **Code Review** - 所有變更都應該經過 PR review
6. **測試變更** - 先在 dev 環境測試
7. **標記版本** - 使用 Git tags 標記穩定版本

## 📚 進階主題

### Remote State Backend

```hcl
# terraform/versions.tf
terraform {
  backend "s3" {
    bucket = "pandora-terraform-state"
    key    = "prod/terraform.tfstate"
    region = "us-west-2"
  }
}
```

### 狀態遷移

```bash
# 遷移到 remote backend
terraform init -migrate-state
```

### Import 現有資源

```bash
# 導入 Fly.io app
terraform import module.flyio_monitoring.fly_app.monitoring pandora-monitoring
```

## 🆘 獲取幫助

- **Terraform 文件**: https://www.terraform.io/docs
- **Fly.io Provider**: https://registry.terraform.io/providers/fly-apps/fly
- **GitHub Issues**: 在專案中回報問題

---

**維護者**: Pandora DevOps Team  
**最後更新**: 2024-12-19  
**Terraform 版本**: >= 1.0

