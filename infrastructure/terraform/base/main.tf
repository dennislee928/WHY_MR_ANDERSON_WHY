# ================================
# Pandora Box Console IDS-IPS
# Multi-Platform PaaS Deployment
# ================================

# ================================
# Railway PostgreSQL
# ================================

module "railway_postgres" {
  source = "./modules/railway"

  project_id       = var.railway_project_id
  api_token        = var.railway_api_token
  service_name     = "${var.project_name}-postgres-${var.environment}"
  postgres_version = "15"
  database_name    = "pandora"
}

# ================================
# Render Services (Redis + Nginx)
# ================================

module "render_services" {
  source = "./modules/render"

  api_key         = var.render_api_key
  region          = "oregon"
  redis_name      = "${var.project_name}-redis-${var.environment}"
  nginx_name      = "${var.project_name}-nginx-${var.environment}"
  repository_url  = "https://github.com/cyber-security-dev-dep-mitake-com-tw/pandora_box_console_IDS-IPS"
  branch          = var.environment == "prod" ? "main" : "dev"

  # 使用計算的 URL 而不是模組輸出（避免循環依賴）
  axiom_ui_url     = "https://${var.project_name}-ui-${var.environment}.patr.cloud"
  grafana_url      = "https://${var.project_name}-monitoring-${var.environment}.fly.dev"
  prometheus_url   = "https://${var.project_name}-monitoring-${var.environment}.fly.dev/prometheus"
  loki_url         = "https://${var.project_name}-monitoring-${var.environment}.fly.dev/loki"
  alertmanager_url = "https://${var.project_name}-monitoring-${var.environment}.fly.dev/alertmanager"
}

# ================================
# Fly.io Monitoring System
# ================================

module "flyio_monitoring" {
  source = "./modules/flyio"

  project_name   = var.project_name
  environment    = var.environment
  organization   = var.fly_organization
  region         = var.fly_region
  volume_size    = var.monitoring_volume_size
  docker_image   = var.monitoring_image
  log_level      = var.log_level
  fly_api_token  = var.fly_api_token

  grafana_admin_password = var.grafana_admin_password
}

# ================================
# Koyeb Agent + Promtail
# ================================

module "koyeb_agent" {
  source = "./modules/koyeb"

  api_token      = var.koyeb_api_token
  app_name       = "${var.project_name}-agent-${var.environment}"
  docker_image   = var.agent_image
  region         = "fra"
  instance_type  = "nano"
  log_level      = var.log_level

  database_url   = module.railway_postgres.database_url
  redis_url      = module.render_services.redis_url
  prometheus_url = module.flyio_monitoring.prometheus_url
  loki_url       = module.flyio_monitoring.loki_url

  depends_on = [
    module.railway_postgres,
    module.render_services,
    module.flyio_monitoring
  ]
}

# ================================
# Patr.io UI
# ================================

module "patr_ui" {
  source = "./modules/patr"

  api_token      = var.patr_api_token
  app_name       = "${var.project_name}-ui-${var.environment}"
  docker_image   = var.ui_image

  api_url        = module.koyeb_agent.api_url
  grafana_url    = module.flyio_monitoring.grafana_url
  prometheus_url = module.flyio_monitoring.prometheus_url

  depends_on = [
    module.koyeb_agent,
    module.flyio_monitoring
  ]
}

