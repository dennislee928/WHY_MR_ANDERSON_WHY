# ================================
# Patr.io Axiom UI 模組
# 使用 Patr.io API 管理
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

# Patr.io Application
resource "null_resource" "patr_ui" {
  triggers = {
    app_name     = var.app_name
    docker_image = var.docker_image
  }

  provisioner "local-exec" {
    command = <<-EOT
      echo "Creating Patr.io UI service via API..."
      
      curl -X POST https://api.patr.cloud/api/v1/deployments \
        -H "Authorization: Bearer ${var.api_token}" \
        -H "Content-Type: application/json" \
        -d '{
          "name": "${var.app_name}",
          "image": "${var.docker_image}",
          "env": [
            {"key": "NEXT_PUBLIC_API_BASE_URL", "value": "${var.api_url}"},
            {"key": "NEXT_PUBLIC_GRAFANA_URL", "value": "${var.grafana_url}"},
            {"key": "NEXT_PUBLIC_PROMETHEUS_URL", "value": "${var.prometheus_url}"}
          ],
          "resources": {
            "cpu": "0.5",
            "memory": "512Mi"
          },
          "ports": [{
            "containerPort": 3000,
            "protocol": "HTTP"
          }],
          "healthCheck": {
            "path": "/api/health",
            "port": 3000
          }
        }'
    EOT
  }
}

