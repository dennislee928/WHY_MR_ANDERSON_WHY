output "api_url" {
  description = "Koyeb application URL"
  value       = "https://${var.app_name}.koyeb.app"
}

output "service_id" {
  description = "Koyeb service ID"
  value       = null_resource.koyeb_agent.id
}

