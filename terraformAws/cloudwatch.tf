resource "aws_cloudwatch_log_group" "ecs_log_group" {
  name              = local.ecs_log_group_name
  tags              = local.tags
  retention_in_days = 30
}
