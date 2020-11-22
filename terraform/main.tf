variable "gcp_project" {}
variable "gcp_region" {}
variable "gcp_config_file" {}

provider "google" {
  credentials = file(var.gcp_config_file)
  project     = var.gcp_project
  region      = var.gcp_region
  version     = "v3.48.0"
}

terraform {
  required_version = "v0.12.26"
}
