variable "api_token" {
  description = "Patr.io API token"
  type        = string
  sensitive   = true
}

variable "app_name" {
  description = "Application name"
  type        = string
}

variable "docker_image" {
  description = "Docker image"
  type        = string
}

variable "api_url" {
  description = "Backend API URL"
  type        = string
}

variable "grafana_url" {
  description = "Grafana URL"
  type        = string
}

variable "prometheus_url" {
  description = "Prometheus URL"
  type        = string
}

