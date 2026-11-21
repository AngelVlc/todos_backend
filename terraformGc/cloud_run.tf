resource "google_cloud_run_service" "run_service" {
  name     = var.app_name
  location = local.region

  template {
    metadata {
      annotations = {
        "run.googleapis.com/container-dependencies" = jsonencode({ (var.app_name) = ["cloudsql-proxy"] })
      }
    }

    spec {
      containers {
        name  = var.app_name
        image = var.container_image

        ports {
          container_port = 8080
        }

        volume_mounts {
          name       = "cloudsql"
          mount_path = "/cloudsql"
        }

        env {
          name  = "ENVIRONMENT"
          value = "production"
        }
        env {
          name  = "MYSQL_UNIX_SOCKET"
          value = "/cloudsql/${var.cloud_sql_instance_connection_name}"
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
          name  = "TOKEN_EXPIRATION_TIME"
          value = var.token_expiration_time
        }
        env {
          name  = "REFRESH_TOKEN_EXPIRATION_TIME"
          value = var.refresh_token_expiration_time
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

        startup_probe {
          tcp_socket {
            port = 8080
          }
        }
      }

      containers {
        name  = "cloudsql-proxy"
        image = "gcr.io/cloud-sql-connectors/cloud-sql-proxy:2.28.0"
        args = [
          "--structured-logs",
          "--unix-socket=/cloudsql/${var.cloud_sql_instance_connection_name}",
          "--health-check",
          "--http-address=0.0.0.0",
        ]

        volume_mounts {
          name       = "cloudsql"
          mount_path = "/cloudsql"
        }

        startup_probe {
          http_get {
            path = "/startup"
            port = 9090
          }
        }
      }

      volumes {
        name = "cloudsql"
        empty_dir {
          medium = "Memory"
        }
      }

      service_account_name = local.service_account_name
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [
    google_project_service.run,
    google_project_service.sqladmin,
  ]
}

resource "google_cloud_run_service_iam_member" "allUsers" {
  service  = google_cloud_run_service.run_service.name
  location = google_cloud_run_service.run_service.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}
