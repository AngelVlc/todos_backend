resource "google_artifact_registry_repository" "todos_backend_repo" {
  location      = "us-central1"
  repository_id = "todos-backend"
  description   = "todos backend repository"
  format        = "DOCKER"
}