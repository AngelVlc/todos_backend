terraform {
  required_version = ">= 0.14.9"

  required_providers {
    google = ">= 3.3"
  }

  backend "gcs" {
    bucket = "terraform-todos-backend"
    prefix = "state"
  }
}

provider "google" {
  project = var.gc_project_id
  region  = local.region
}
