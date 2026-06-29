variable "project_id" {
  description = "GCP project ID"
  type        = string
  default     = "biomech-217fe"
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "vpc_id" {
  description = "VPC network self-link (projects/.../networks/...)"
  type        = string
  default     = null
}

variable "vpc_name" {
  description = "VPC network name"
  type        = string
  default     = "default"
}
