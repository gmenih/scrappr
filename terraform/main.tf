variable "gcp_project" {}
variable "gcp_region" {}

provider "google" {
  credentials = file("../config/google-key.json")
  project     = var.gcp_project
  region      = var.gcp_region
  version     = "v3.48.0"
}

terraform {
  required_version = "v0.12.26"
}
