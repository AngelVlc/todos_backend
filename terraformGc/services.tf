resource "google_project_service" "resource_manager" {
  service = "cloudresourcemanager.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "container_registry" {
  service = "containerregistry.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "run" {
  service = "run.googleapis.com"
}