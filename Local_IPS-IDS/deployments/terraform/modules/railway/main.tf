# ================================
# Railway PostgreSQL 模組
# 使用 HTTP API 管理（沒有官方 Terraform Provider）
# ================================

terraform {
  required_providers {
    http = {
      source  = "hashicorp/http"
      version = "~> 3.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
  }
}

# 使用 null_resource 和 local-exec 來管理 Railway 資源
# 這是因為 Railway 沒有官方 Terraform provider

resource "null_resource" "railway_postgres" {
  triggers = {
    project_id       = var.project_id
    service_name     = var.service_name
    postgres_version = var.postgres_version
  }

  provisioner "local-exec" {
    command = <<-EOT
      # 使用 Railway CLI 建立 PostgreSQL 服務
      # 需要先安裝: npm install -g @railway/cli
      
      echo "Creating Railway PostgreSQL service..."
      
      # 設定 Railway API token
      export RAILWAY_TOKEN="${var.api_token}"
      
      # 建立或更新 PostgreSQL 插件
      railway link ${var.project_id}
      railway add --plugin postgres
      
      echo "Railway PostgreSQL service created/updated"
    EOT

    environment = {
      RAILWAY_TOKEN = var.api_token
    }
  }

  provisioner "local-exec" {
    when    = destroy
    command = <<-EOT
      echo "Destroying Railway PostgreSQL service..."
      # 注意：實際刪除需要手動確認或使用 Railway API
    EOT
  }
}

# 使用 Data Source 獲取連線資訊
# 注意：這需要 Railway API 支援
data "external" "railway_connection" {
  program = ["bash", "-c", <<-EOT
    # 獲取 Railway 服務連線資訊
    export RAILWAY_TOKEN="${var.api_token}"
    
    # 使用 Railway CLI 獲取環境變數
    railway variables | grep DATABASE_URL | awk '{print $2}' | jq -R '{database_url: .}'
  EOT
  ]

  depends_on = [null_resource.railway_postgres]
}

