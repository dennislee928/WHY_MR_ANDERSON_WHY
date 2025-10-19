variable "project_name" {
  description = "Project name prefix"
  type        = string
}

variable "environment" {
  description = "Environment (dev/staging/prod)"
  type        = string
}

variable "organization" {
  description = "Fly.io organization"
  type        = string
}

variable "region" {
  description = "Fly.io region"
  type        = string
  default     = "nrt"
}

variable "volume_size" {
  description = "Volume size in GB"
  type        = number
  default     = 10
}

variable "docker_image" {
  description = "Docker image for monitoring system"
  type        = string
}

variable "log_level" {
  description = "Log level"
  type        = string
  default     = "info"
}

variable "grafana_admin_password" {
  description = "Grafana admin password"
  type        = string
  sensitive   = true
}

variable "fly_api_token" {
  description = "Fly.io API token"
  type        = string
  sensitive   = true
}

