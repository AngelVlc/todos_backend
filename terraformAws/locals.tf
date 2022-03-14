locals {
  name               = "${var.app_name}-${var.environment}"
  port               = 80
  ecs_log_group_name = "/aws/ecs/${local.name}"
  database_name      = "todos"
  tags = {
    Environment = var.environment
    Project     = var.app_name
  }
}
