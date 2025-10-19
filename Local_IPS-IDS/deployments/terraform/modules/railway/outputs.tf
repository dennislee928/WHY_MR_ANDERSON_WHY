output "database_url" {
  description = "PostgreSQL connection string"
  value       = try(data.external.railway_connection.result.database_url, "")
  sensitive   = true
}

output "service_id" {
  description = "Railway service ID"
  value       = null_resource.railway_postgres.id
}

