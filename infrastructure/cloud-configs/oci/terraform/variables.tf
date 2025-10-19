# OCI Terraform Variables

variable "tenancy_ocid" {
  description = "OCI Tenancy OCID"
  type        = string
}

variable "user_ocid" {
  description = "OCI User OCID"
  type        = string
}

variable "fingerprint" {
  description = "API Key Fingerprint"
  type        = string
}

variable "private_key_path" {
  description = "Path to OCI API private key"
  type        = string
  default     = "~/.oci/oci_api_key.pem"
}

variable "region" {
  description = "OCI Region"
  type        = string
  default     = "us-ashburn-1"
}

variable "compartment_id" {
  description = "Compartment OCID"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

# Compute Instance Configuration (Always Free Tier)
variable "instance_shape" {
  description = "Instance shape (VM.Standard.A1.Flex for Arm, VM.Standard.E2.1.Micro for AMD)"
  type        = string
  default     = "VM.Standard.A1.Flex" # Arm-based, Always Free
}

variable "instance_ocpus" {
  description = "Number of OCPUs (max 4 for Always Free Arm)"
  type        = number
  default     = 2
}

variable "instance_memory_gb" {
  description = "Memory in GB (max 24 for Always Free Arm total)"
  type        = number
  default     = 12
}

variable "data_volume_size_gb" {
  description = "Block volume size in GB (max 200 for Always Free)"
  type        = number
  default     = 100
}

variable "ssh_public_key_path" {
  description = "Path to SSH public key"
  type        = string
  default     = "~/.ssh/id_rsa.pub"
}

