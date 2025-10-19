variable "api_token" {
  description = "Koyeb API token"
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

variable "region" {
  description = "Koyeb region"
  type        = string
  default     = "fra"
}

variable "instance_type" {
  description = "Instance type"
  type        = string
  default     = "nano"
}

variable "log_level" {
  description = "Log level"
  type        = string
  default     = "info"
}

variable "database_url" {
  description = "Database connection URL"
  type        = string
  sensitive   = true
}

variable "redis_url" {
  description = "Redis connection URL"
  type        = string
  sensitive   = true
}

variable "prometheus_url" {
  description = "Prometheus URL"
  type        = string
}

variable "loki_url" {
  description = "Loki URL"
  type        = string
}

