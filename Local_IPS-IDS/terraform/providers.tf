# Provider 配置
# 注意：所有 PaaS 平台（Railway, Render, Koyeb, Patr, Fly.io）
# 都沒有穩定的官方 Terraform provider
# 使用標準 HashiCorp providers 和 null_resource + CLI/API 管理

# HTTP Provider - 用於 API 呼叫
provider "http" {}

# Null Provider - 用於執行 CLI 命令
provider "null" {}

# External Provider - 用於執行外部腳本
provider "external" {}

