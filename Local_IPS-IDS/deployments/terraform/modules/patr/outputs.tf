output "ui_url" {
  description = "Patr.io UI URL"
  value       = "https://${var.app_name}.patr.cloud"
}

output "service_id" {
  description = "Patr.io service ID"
  value       = null_resource.patr_ui.id
}

