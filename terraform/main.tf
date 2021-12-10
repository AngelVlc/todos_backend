terraform {
  backend "pg" {
    schema_name = "d29h8hsdmnk12h"
  }
}

variable "heroku_username" {
  description = "Heroku user name"
}

variable "heroku_api_key" {
  description = "Heroku api key"
}

variable "app_name" {
  description = "Name of the Heroku app to be provisioned"
}

variable "jwt_secret" {
  description = "JWT secret"
}

variable "admin_password" {
  description = "Password for admin"
}

variable "cors_allowed_origins" {
  description = "Comma separated CORS allowed domains"
}

variable "environment" {
  description = "Name of the environment to be provisioned"
}

provider "heroku" {
  email   = var.heroku_username
  api_key = var.heroku_api_key
}

resource "heroku_app" "default" {
  name   = var.app_name
  region = "eu"
  stack  = "container"
}

resource "heroku_addon" "database" {
  app    = heroku_app.default.name
  plan   = "cleardb:ignite"
  lifecycle {
    prevent_destroy = true
  }
}

resource "heroku_addon" "errors" {
  app    = heroku_app.default.name
  plan   = "honeybadger:free"
}

resource "heroku_addon" "apm" {
  app    = heroku_app.default.name
  plan   = "newrelic:wayne"
}

resource "heroku_addon" "log" {
  app    = heroku_app.default.name
  plan   = "sumologic:free"
}

resource "heroku_config" "default" {
  sensitive_vars = {
    JWT_SECRET = var.jwt_secret
    ADMIN_PASSWORD = var.admin_password
    CORS_ALLOWED_ORIGINS = var.cors_allowed_origins
    ENVIRONMENT = var.environment
    TOKEN_EXPIRATION_DURATION = "5m"
    REFRESH_TOKEN_EXPIRATION_DURATION = "24h"
    DELETE_EXPIRED_REFRESH_TOKEN_INTERVAL = "30s"
  }
}

resource "heroku_app_config_association" "default" {
  app_id = heroku_app.default.id

  sensitive_vars = heroku_config.default.sensitive_vars
}
