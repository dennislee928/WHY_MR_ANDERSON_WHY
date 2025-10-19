terraform {
  required_version = ">= 1.0"

  required_providers {
    http = {
      source  = "hashicorp/http"
      version = "~> 3.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
    external = {
      source  = "hashicorp/external"
      version = "~> 2.0"
    }
    # 注意：所有 PaaS 平台（Railway, Render, Koyeb, Patr, Fly.io）
    # 都沒有穩定的官方 Terraform provider
    # 使用 null_resource + CLI/API 呼叫方式管理
  }

  # Backend 配置 - 用於儲存 Terraform state
  # 生產環境建議使用遠端 backend (S3, Terraform Cloud, 等)
  backend "local" {
    path = "terraform.tfstate"
  }
}

