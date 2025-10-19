# Oracle Cloud Infrastructure (OCI) Terraform Configuration
# Deploys Security Platform to OCI Always Free Tier
#
# Free Tier Resources:
# - 2x AMD Compute VMs (1/8 OCPU, 1 GB RAM each)
# - 4x Arm-based Ampere A1 cores (24 GB RAM total)
# - 200 GB Block Volume storage
# - 10 GB Object Storage

terraform {
  required_version = ">= 1.0"
  
  required_providers {
    oci = {
      source  = "oracle/oci"
      version = "~> 5.0"
    }
  }
}

# Provider configuration
provider "oci" {
  tenancy_ocid     = var.tenancy_ocid
  user_ocid        = var.user_ocid
  fingerprint      = var.fingerprint
  private_key_path = var.private_key_path
  region           = var.region
}

# Data source for availability domain
data "oci_identity_availability_domains" "ads" {
  compartment_id = var.tenancy_ocid
}

# Create VCN (Virtual Cloud Network)
resource "oci_core_vcn" "security_platform_vcn" {
  compartment_id = var.compartment_id
  display_name   = "security-platform-vcn"
  cidr_blocks    = ["10.0.0.0/16"]
  dns_label      = "secplatform"
}

# Internet Gateway
resource "oci_core_internet_gateway" "security_platform_igw" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.security_platform_vcn.id
  display_name   = "security-platform-igw"
  enabled        = true
}

# Route Table
resource "oci_core_route_table" "security_platform_rt" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.security_platform_vcn.id
  display_name   = "security-platform-rt"

  route_rules {
    destination       = "0.0.0.0/0"
    network_entity_id = oci_core_internet_gateway.security_platform_igw.id
  }
}

# Security List
resource "oci_core_security_list" "security_platform_sl" {
  compartment_id = var.compartment_id
  vcn_id         = oci_core_vcn.security_platform_vcn.id
  display_name   = "security-platform-sl"

  # Egress: Allow all outbound
  egress_security_rules {
    destination = "0.0.0.0/0"
    protocol    = "all"
  }

  # Ingress: SSH
  ingress_security_rules {
    protocol = "6" # TCP
    source   = "0.0.0.0/0"
    tcp_options {
      min = 22
      max = 22
    }
  }

  # Ingress: HTTP
  ingress_security_rules {
    protocol = "6"
    source   = "0.0.0.0/0"
    tcp_options {
      min = 80
      max = 80
    }
  }

  # Ingress: HTTPS
  ingress_security_rules {
    protocol = "6"
    source   = "0.0.0.0/0"
    tcp_options {
      min = 443
      max = 443
    }
  }

  # Ingress: Backend API (3001)
  ingress_security_rules {
    protocol = "6"
    source   = "0.0.0.0/0"
    tcp_options {
      min = 3001
      max = 3001
    }
  }

  # Ingress: AI/Quantum API (8000)
  ingress_security_rules {
    protocol = "6"
    source   = "0.0.0.0/0"
    tcp_options {
      min = 8000
      max = 8000
    }
  }

  # Ingress: Grafana (3000)
  ingress_security_rules {
    protocol = "6"
    source   = "0.0.0.0/0"
    tcp_options {
      min = 3000
      max = 3000
    }
  }

  # Ingress: Prometheus (9090)
  ingress_security_rules {
    protocol = "6"
    source   = "0.0.0.0/0"
    tcp_options {
      min = 9090
      max = 9090
    }
  }
}

# Subnet
resource "oci_core_subnet" "security_platform_subnet" {
  compartment_id      = var.compartment_id
  vcn_id              = oci_core_vcn.security_platform_vcn.id
  cidr_block          = "10.0.1.0/24"
  display_name        = "security-platform-subnet"
  dns_label           = "secplatformsub"
  route_table_id      = oci_core_route_table.security_platform_rt.id
  security_list_ids   = [oci_core_security_list.security_platform_sl.id]
  dhcp_options_id     = oci_core_vcn.security_platform_vcn.default_dhcp_options_id
}

# Get Ubuntu image
data "oci_core_images" "ubuntu_images" {
  compartment_id           = var.compartment_id
  operating_system         = "Canonical Ubuntu"
  operating_system_version = "22.04"
  shape                    = var.instance_shape
  sort_by                  = "TIMECREATED"
  sort_order               = "DESC"
}

# Compute Instance 1 (Main Application Server)
resource "oci_core_instance" "app_server" {
  compartment_id      = var.compartment_id
  availability_domain = data.oci_identity_availability_domains.ads.availability_domains[0].name
  display_name        = "security-platform-app"
  shape               = var.instance_shape

  shape_config {
    ocpus         = var.instance_ocpus
    memory_in_gbs = var.instance_memory_gb
  }

  create_vnic_details {
    subnet_id        = oci_core_subnet.security_platform_subnet.id
    assign_public_ip = true
    display_name     = "app-server-vnic"
  }

  source_details {
    source_type = "image"
    source_id   = data.oci_core_images.ubuntu_images.images[0].id
    boot_volume_size_in_gbs = 50
  }

  metadata = {
    ssh_authorized_keys = file(var.ssh_public_key_path)
    user_data = base64encode(templatefile("${path.module}/cloud-init-app.yaml", {
      docker_compose_version = "2.23.0"
    }))
  }

  freeform_tags = {
    "Environment" = var.environment
    "Project"     = "security-platform"
    "Role"        = "application-server"
  }
}

# Compute Instance 2 (Database & Monitoring Server)
resource "oci_core_instance" "db_server" {
  compartment_id      = var.compartment_id
  availability_domain = data.oci_identity_availability_domains.ads.availability_domains[0].name
  display_name        = "security-platform-db"
  shape               = var.instance_shape

  shape_config {
    ocpus         = var.instance_ocpus
    memory_in_gbs = var.instance_memory_gb
  }

  create_vnic_details {
    subnet_id        = oci_core_subnet.security_platform_subnet.id
    assign_public_ip = true
    display_name     = "db-server-vnic"
  }

  source_details {
    source_type = "image"
    source_id   = data.oci_core_images.ubuntu_images.images[0].id
    boot_volume_size_in_gbs = 50
  }

  metadata = {
    ssh_authorized_keys = file(var.ssh_public_key_path)
    user_data = base64encode(templatefile("${path.module}/cloud-init-db.yaml", {
      docker_compose_version = "2.23.0"
    }))
  }

  freeform_tags = {
    "Environment" = var.environment
    "Project"     = "security-platform"
    "Role"        = "database-server"
  }
}

# Block Volume for persistent data
resource "oci_core_volume" "data_volume" {
  compartment_id      = var.compartment_id
  availability_domain = data.oci_identity_availability_domains.ads.availability_domains[0].name
  display_name        = "security-platform-data"
  size_in_gbs         = var.data_volume_size_gb

  freeform_tags = {
    "Environment" = var.environment
    "Project"     = "security-platform"
  }
}

# Attach block volume to DB server
resource "oci_core_volume_attachment" "data_volume_attachment" {
  attachment_type = "paravirtualized"
  instance_id     = oci_core_instance.db_server.id
  volume_id       = oci_core_volume.data_volume.id
  display_name    = "data-volume-attachment"
}

# Object Storage Bucket for backups
resource "oci_objectstorage_bucket" "backup_bucket" {
  compartment_id = var.compartment_id
  namespace      = data.oci_objectstorage_namespace.ns.namespace
  name           = "security-platform-backups-${var.environment}"
  access_type    = "NoPublicAccess"

  freeform_tags = {
    "Environment" = var.environment
    "Project"     = "security-platform"
  }
}

data "oci_objectstorage_namespace" "ns" {
  compartment_id = var.compartment_id
}

# Outputs
output "app_server_public_ip" {
  description = "Public IP of application server"
  value       = oci_core_instance.app_server.public_ip
}

output "db_server_public_ip" {
  description = "Public IP of database server"
  value       = oci_core_instance.db_server.public_ip
}

output "vcn_id" {
  description = "VCN ID"
  value       = oci_core_vcn.security_platform_vcn.id
}

output "backup_bucket_name" {
  description = "Backup bucket name"
  value       = oci_objectstorage_bucket.backup_bucket.name
}

output "ssh_connection_app" {
  description = "SSH connection command for app server"
  value       = "ssh ubuntu@${oci_core_instance.app_server.public_ip}"
}

output "ssh_connection_db" {
  description = "SSH connection command for DB server"
  value       = "ssh ubuntu@${oci_core_instance.db_server.public_ip}"
}

