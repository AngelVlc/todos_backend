variable "jwt_secret" {
  description = "JWT secret"
}

variable "cors_allowed_origins" {
  description = "Comma separated CORS allowed domains"
}

variable "new_relic_license_key" {
  description = "NewRelic license key"
}

variable "honeybadger_api_key" {
  description = "HoneyBadger api key"
}

variable "app_name" {
  description = "Name of the app to be provisioned"
}

variable "gc_project_id" {
  description = "Google Cloud Project Id"
}

variable "container_image" {
  description = "Container imaged used in the Cloud Run service"
  default = "image"
}

variable "delete_expired_refresh_token_interval" {
  description = "Interval for the delete expired tokens process. Valid time units are s, m or h"
}

variable "mysql_tls" {}
variable "mysql_host" {}
variable "mysql_user" {}
variable "mysql_password" {}
variable "mysql_port" {
  default = "3306"
}
variable "mysql_database" {}

variable "algolia_app_id" {}
variable "algolia_api_key" {}
variable "algolia_search_only_key" {}
