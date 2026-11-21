resource "google_project_iam_member" "cloudsql_client" {
  project = var.gc_project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${local.service_account_name}@${var.gc_project_id}.iam.gserviceaccount.com"
}
