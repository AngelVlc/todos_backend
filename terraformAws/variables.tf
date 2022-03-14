variable "app_name" {
  description = "Name of the Heroku app to be provisioned"
}

variable "jwt_secret" {
  description = "JWT secret"
}

variable "cors_allowed_origins" {
  description = "Comma separated CORS allowed domains"
}

variable "environment" {
  description = "Name of the environment to be provisioned"
}

variable "new_relic_license_key" {
  description = "NewRelic license key"
}

variable "honeybadger_api_key" {
  description = "HoneyBadger api key"
}

variable "subdomain" {
  description = "Subdomain"
}