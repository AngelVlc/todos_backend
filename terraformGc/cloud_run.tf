resource "google_cloud_run_service" "run_service" {
  name     = var.app_name
  location = local.region

  template {
    spec {
      containers {
        image = var.container_image
        env {
          name  = "ENVIRONMENT"
          value = "production"
        }
        env {
          name  = "MYSQL_TLS"
          value = var.mysql_tls
        }
        env {
          name  = "MYSQL_HOST"
          value = var.mysql_host
        }
        env {
          name  = "MYSQL_PORT"
          value = var.mysql_port
        }
        env {
          name  = "MYSQL_DATABASE"
          value = var.mysql_database
        }
        env {
          name  = "MYSQL_USER"
          value = var.mysql_user
        }
        env {
          name  = "MYSQL_PASSWORD"
          value = var.mysql_password
        }
        env {
          name  = "JWT_SECRET"
          value = var.jwt_secret
        }
        env {
          name  = "CORS_ALLOWED_ORIGINS"
          value = var.cors_allowed_origins
        }
        env {
          name  = "NEW_RELIC_LICENSE_KEY"
          value = var.new_relic_license_key
        }
        env {
          name  = "HONEYBADGER_API_KEY"
          value = var.honeybadger_api_key
        }
        env {
          name  = "DELETE_EXPIRED_REFRESH_TOKEN_INTERVAL"
          value = var.delete_expired_refresh_token_interval
        }
        env {
          name  = "ALGOLIA_APP_ID"
          value = var.algolia_app_id
        }
        env {
          name  = "ALGOLIA_API_KEY"
          value = var.algolia_api_key
        }
        env {
          name  = "ALGOLIA_SEARCH_ONLY_KEY"
          value = var.algolia_search_only_key
        }
      }
      service_account_name = "todos-service-account"
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.run]
}

resource "google_cloud_run_service_iam_member" "allUsers" {
  service  = google_cloud_run_service.run_service.name
  location = google_cloud_run_service.run_service.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}
