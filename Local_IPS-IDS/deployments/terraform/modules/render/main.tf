# ================================
# Render Redis + Nginx 模組
# 使用 Render API 管理
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

# Render Redis Service
resource "null_resource" "render_redis" {
  triggers = {
    service_name = var.redis_name
    plan         = var.redis_plan
  }

  provisioner "local-exec" {
    command = <<-EOT
      echo "Creating Render Redis service via API..."
      
      curl -X POST https://api.render.com/v1/services \
        -H "Authorization: Bearer ${var.api_key}" \
        -H "Content-Type: application/json" \
        -d '{
          "type": "redis",
          "name": "${var.redis_name}",
          "plan": "${var.redis_plan}",
          "region": "${var.region}",
          "maxmemoryPolicy": "allkeys-lru"
        }'
    EOT
  }
}

# Render Nginx Service
resource "null_resource" "render_nginx" {
  triggers = {
    service_name = var.nginx_name
    dockerfile   = var.nginx_dockerfile
  }

  provisioner "local-exec" {
    command = <<-EOT
      echo "Creating Render Nginx service via API..."
      
      curl -X POST https://api.render.com/v1/services \
        -H "Authorization: Bearer ${var.api_key}" \
        -H "Content-Type: application/json" \
        -d '{
          "type": "web_service",
          "name": "${var.nginx_name}",
          "plan": "${var.nginx_plan}",
          "region": "${var.region}",
          "env": "docker",
          "dockerfilePath": "${var.nginx_dockerfile}",
          "repo": "${var.repository_url}",
          "branch": "${var.branch}",
          "healthCheckPath": "/health",
          "envVars": [
            {"key": "AXIOM_UI_URL", "value": "${var.axiom_ui_url}"},
            {"key": "GRAFANA_URL", "value": "${var.grafana_url}"},
            {"key": "PROMETHEUS_URL", "value": "${var.prometheus_url}"},
            {"key": "LOKI_URL", "value": "${var.loki_url}"},
            {"key": "ALERTMANAGER_URL", "value": "${var.alertmanager_url}"}
          ]
        }'
    EOT
  }

  depends_on = [null_resource.render_redis]
}

