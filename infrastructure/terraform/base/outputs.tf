# ================================
# 服務 URLs
# ================================

output "service_urls" {
  description = "All service URLs"
  value = {
    grafana      = module.flyio_monitoring.grafana_url
    prometheus   = module.flyio_monitoring.prometheus_url
    loki         = module.flyio_monitoring.loki_url
    alertmanager = module.flyio_monitoring.alertmanager_url
    ui           = module.patr_ui.ui_url
    api          = module.koyeb_agent.api_url
    nginx        = module.render_services.nginx_url
  }
}

# ================================
# 資料庫連線資訊
# ================================

output "database_url" {
  description = "PostgreSQL connection string"
  value       = module.railway_postgres.database_url
  sensitive   = true
}

output "redis_url" {
  description = "Redis connection string"
  value       = module.render_services.redis_url
  sensitive   = true
}

# ================================
# Fly.io 資訊
# ================================

output "fly_app_name" {
  description = "Fly.io app name"
  value       = module.flyio_monitoring.app_name
}

output "fly_volume_id" {
  description = "Fly.io volume ID"
  value       = module.flyio_monitoring.volume_id
}

# ================================
# 部署摘要
# ================================

output "deployment_summary" {
  description = "Complete deployment summary"
  value = {
    environment = var.environment
    project     = var.project_name
    services = {
      database   = "Railway PostgreSQL"
      cache      = "Render Redis"
      proxy      = "Render Nginx"
      backend    = "Koyeb Agent + Promtail"
      frontend   = "Patr.io Axiom UI"
      monitoring = "Fly.io (Prometheus + Loki + Grafana + AlertManager)"
    }
    deployment_time = timestamp()
  }
}

