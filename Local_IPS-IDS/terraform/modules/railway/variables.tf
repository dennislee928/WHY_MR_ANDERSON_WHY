variable "project_id" {
  description = "Railway project ID"
  type        = string
  sensitive   = true
}

variable "api_token" {
  description = "Railway API token"
  type        = string
  sensitive   = true
}

variable "service_name" {
  description = "PostgreSQL service name"
  type        = string
  default     = "postgres"
}

variable "postgres_version" {
  description = "PostgreSQL version"
  type        = string
  default     = "15"
}

variable "database_name" {
  description = "Database name"
  type        = string
  default     = "pandora"
}

