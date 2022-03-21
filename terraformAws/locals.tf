locals {
  app_name           = "todosbackend"
  name               = "${local.app_name}-${var.environment}"
  port               = 80
  ecs_log_group_name = "/aws/ecs/${local.name}"
  database_name      = "todos"
  tags = {
    Environment = var.environment
    Project     = local.app_name
  }
}
