output "app_name" {
  description = "Fly.io app name"
  value       = local.app_name
}

output "app_url" {
  description = "Fly.io app URL"
  value       = "https://${local.app_name}.fly.dev"
}

output "grafana_url" {
  description = "Grafana URL"
  value       = "https://${local.app_name}.fly.dev"
}

output "prometheus_url" {
  description = "Prometheus URL"
  value       = "https://${local.app_name}.fly.dev/prometheus"
}

output "loki_url" {
  description = "Loki URL"
  value       = "https://${local.app_name}.fly.dev/loki"
}

output "alertmanager_url" {
  description = "AlertManager URL"
  value       = "https://${local.app_name}.fly.dev/alertmanager"
}

output "volume_id" {
  description = "Volume ID"
  value       = null_resource.fly_volume.id
}

output "deploy_id" {
  description = "Deployment ID"
  value       = null_resource.fly_deploy.id
}

