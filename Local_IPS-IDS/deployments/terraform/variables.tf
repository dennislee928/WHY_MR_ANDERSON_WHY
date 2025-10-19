# ================================
# 環境配置
# ================================

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "dev"
  
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be dev, staging, or prod."
  }
}

variable "project_name" {
  description = "Project name prefix"
  type        = string
  default     = "pandora-box-console"
}

# ================================
# Railway PostgreSQL 配置
# ================================

variable "railway_project_id" {
  description = "Railway project ID"
  type        = string
  sensitive   = true
}

variable "railway_api_token" {
  description = "Railway API token"
  type        = string
  sensitive   = true
}

variable "postgres_password" {
  description = "PostgreSQL password"
  type        = string
  sensitive   = true
}

# ================================
# Render 配置
# ================================

variable "render_api_key" {
  description = "Render API key"
  type        = string
  sensitive   = true
}

# ================================
# Koyeb 配置
# ================================

variable "koyeb_api_token" {
  description = "Koyeb API token"
  type        = string
  sensitive   = true
}

variable "koyeb_organization_id" {
  description = "Koyeb organization ID"
  type        = string
  sensitive   = true
}

# ================================
# Patr.io 配置
# ================================

variable "patr_api_token" {
  description = "Patr.io API token"
  type        = string
  sensitive   = true
}

# ================================
# Fly.io 配置
# ================================

variable "fly_api_token" {
  description = "Fly.io API token"
  type        = string
  sensitive   = true
}

variable "fly_organization" {
  description = "Fly.io organization"
  type        = string
  default     = "personal"
}

variable "fly_region" {
  description = "Fly.io primary region"
  type        = string
  default     = "nrt"
}

# ================================
# 應用配置
# ================================

variable "grafana_admin_password" {
  description = "Grafana admin password"
  type        = string
  sensitive   = true
}

variable "monitoring_volume_size" {
  description = "Fly.io monitoring volume size in GB"
  type        = number
  default     = 10
}

variable "log_level" {
  description = "Application log level"
  type        = string
  default     = "info"
  
  validation {
    condition     = contains(["debug", "info", "warn", "error"], var.log_level)
    error_message = "Log level must be debug, info, warn, or error."
  }
}

# ================================
# Docker Image 配置
# ================================

variable "agent_image" {
  description = "Pandora Agent Docker image"
  type        = string
  default     = "ghcr.io/your-org/pandora-agent:latest"
}

variable "ui_image" {
  description = "Axiom UI Docker image"
  type        = string
  default     = "ghcr.io/your-org/axiom-ui:latest"
}

variable "monitoring_image" {
  description = "Monitoring system Docker image"
  type        = string
  default     = "registry.fly.io/pandora-monitoring:latest"
}

