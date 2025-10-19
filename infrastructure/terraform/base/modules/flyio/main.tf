# ================================
# Fly.io 監控系統模組
# Prometheus + Loki + Grafana + AlertManager
# 使用 Fly CLI 管理（沒有穩定的 Terraform Provider）
# ================================

terraform {
  required_providers {
    null = {
      source  = "hashicorp/null"
      version = "~> 3.0"
    }
  }
}

locals {
  app_name = "${var.project_name}-monitoring-${var.environment}"
}

# 建立 Fly.io 應用
resource "null_resource" "fly_app" {
  triggers = {
    app_name     = local.app_name
    organization = var.organization
    region       = var.region
  }

  provisioner "local-exec" {
    command = <<-EOT
      echo "Creating Fly.io app: ${local.app_name}"
      fly apps create ${local.app_name} --org ${var.organization} || echo "App already exists"
    EOT

    environment = {
      FLY_API_TOKEN = var.fly_api_token
    }
  }
}

# 建立 Volume
resource "null_resource" "fly_volume" {
  triggers = {
    app_name    = local.app_name
    volume_name = "monitoring_data"
    size        = var.volume_size
    region      = var.region
  }

  provisioner "local-exec" {
    command = <<-EOT
      echo "Creating Fly.io volume..."
      fly volumes create monitoring_data \
        --size ${var.volume_size} \
        --region ${var.region} \
        --app ${local.app_name} \
        --yes || echo "Volume already exists"
    EOT

    environment = {
      FLY_API_TOKEN = var.fly_api_token
    }
  }

  depends_on = [null_resource.fly_app]
}

# 設定 Secrets
resource "null_resource" "fly_secrets" {
  triggers = {
    grafana_password = md5(var.grafana_admin_password)
  }

  provisioner "local-exec" {
    command = <<-EOT
      echo "Setting Fly.io secrets..."
      fly secrets set \
        GRAFANA_ADMIN_PASSWORD="${var.grafana_admin_password}" \
        LOG_LEVEL="${var.log_level}" \
        TZ="Asia/Taipei" \
        --app ${local.app_name}
    EOT

    environment = {
      FLY_API_TOKEN = var.fly_api_token
    }
  }

  depends_on = [null_resource.fly_app]
}

# 部署應用
resource "null_resource" "fly_deploy" {
  triggers = {
    app_name     = local.app_name
    docker_image = var.docker_image
    timestamp    = timestamp()
  }

  provisioner "local-exec" {
    command = <<-EOT
      echo "Deploying to Fly.io..."
      cd ..
      fly deploy \
        --config fly.toml \
        --dockerfile Dockerfile.monitoring \
        --app ${local.app_name} \
        --remote-only \
        --yes
    EOT

    environment = {
      FLY_API_TOKEN = var.fly_api_token
    }
  }

  depends_on = [
    null_resource.fly_app,
    null_resource.fly_volume,
    null_resource.fly_secrets
  ]
}

