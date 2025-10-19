variable "api_key" {
  description = "Render API key"
  type        = string
  sensitive   = true
}

variable "region" {
  description = "Render region"
  type        = string
  default     = "oregon"
}

# Redis 配置
variable "redis_name" {
  description = "Redis service name"
  type        = string
}

variable "redis_plan" {
  description = "Redis plan"
  type        = string
  default     = "free"
}

# Nginx 配置
variable "nginx_name" {
  description = "Nginx service name"
  type        = string
}

variable "nginx_plan" {
  description = "Nginx plan"
  type        = string
  default     = "free"
}

variable "nginx_dockerfile" {
  description = "Path to Nginx Dockerfile"
  type        = string
  default     = "./Dockerfile.nginx"
}

variable "repository_url" {
  description = "Git repository URL"
  type        = string
}

variable "branch" {
  description = "Git branch"
  type        = string
  default     = "main"
}

# 環境變數 - 其他服務 URLs
variable "axiom_ui_url" {
  description = "Axiom UI URL"
  type        = string
  default     = ""
}

variable "grafana_url" {
  description = "Grafana URL"
  type        = string
  default     = ""
}

variable "prometheus_url" {
  description = "Prometheus URL"
  type        = string
  default     = ""
}

variable "loki_url" {
  description = "Loki URL"
  type        = string
  default     = ""
}

variable "alertmanager_url" {
  description = "AlertManager URL"
  type        = string
  default     = ""
}

