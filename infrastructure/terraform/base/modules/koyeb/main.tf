# ================================
# Koyeb Agent + Promtail 模組
# 使用 Koyeb API 管理
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

# Koyeb Application
resource "null_resource" "koyeb_agent" {
  triggers = {
    app_name      = var.app_name
    docker_image  = var.docker_image
    instance_type = var.instance_type
  }

  provisioner "local-exec" {
    command = <<-EOT
      echo "Creating Koyeb Agent service via API..."
      
      curl -X POST https://app.koyeb.com/v1/services \
        -H "Authorization: Bearer ${var.api_token}" \
        -H "Content-Type: application/json" \
        -d '{
          "name": "${var.app_name}",
          "definition": {
            "name": "${var.app_name}",
            "regions": ["${var.region}"],
            "instance_type": "${var.instance_type}",
            "docker": {
              "image": "${var.docker_image}"
            },
            "env": [
              {"key": "LOG_LEVEL", "value": "${var.log_level}"},
              {"key": "DATABASE_URL", "value": "${var.database_url}"},
              {"key": "REDIS_URL", "value": "${var.redis_url}"},
              {"key": "PROMETHEUS_URL", "value": "${var.prometheus_url}"},
              {"key": "LOKI_URL", "value": "${var.loki_url}"}
            ],
            "ports": [
              {"port": 8080, "protocol": "http"}
            ],
            "health_checks": [{
              "http": {
                "path": "/health",
                "port": 8080
              }
            }],
            "scaling": {
              "min": 1,
              "max": 1
            }
          }
        }'
    EOT
  }
}

