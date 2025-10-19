output "redis_url" {
  description = "Redis connection URL"
  value       = "redis://pandora-redis.render.com:6379"
  sensitive   = true
}

output "nginx_url" {
  description = "Nginx service URL"
  value       = "https://${var.nginx_name}.onrender.com"
}

output "redis_service_id" {
  description = "Redis service ID"
  value       = null_resource.render_redis.id
}

output "nginx_service_id" {
  description = "Nginx service ID"
  value       = null_resource.render_nginx.id
}

