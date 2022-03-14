module "aurora_mysql" {
  source  = "terraform-aws-modules/rds-aurora/aws"
  version = "6.2.0"

  name              = "${local.name}-aurora-db-main"
  engine            = "aurora-mysql"
  engine_mode       = "serverless"
  storage_encrypted = true
  tags              = local.tags
  database_name     = local.database_name

  subnets                 = data.aws_subnet_ids.default.ids
  create_security_group   = true
  allowed_security_groups = [aws_security_group.ecs_security_group.id]


  monitoring_interval = 60

  apply_immediately   = true
  skip_final_snapshot = true

  db_parameter_group_name         = aws_db_parameter_group.main.id
  db_cluster_parameter_group_name = aws_rds_cluster_parameter_group.main.id
  # enabled_cloudwatch_logs_exports = # NOT SUPPORTED

  deletion_protection = true

  scaling_configuration = {
    auto_pause               = true
    min_capacity             = 1
    max_capacity             = 1
    seconds_until_auto_pause = 300
    timeout_action           = "ForceApplyCapacityChange"
  }
}

resource "aws_db_parameter_group" "main" {
  name        = "${local.name}-aurora-db-parameter-group"
  family      = "aurora-mysql5.7"
  description = "${local.name}-aurora-db-parameter-group"
  tags        = local.tags
}

resource "aws_rds_cluster_parameter_group" "main" {
  name        = "${local.name}-aurora-db-cluster-parameter-group"
  family      = "aurora-mysql5.7"
  description = "${local.name}-aurora-db-cluster-parameter-group"
  tags        = local.tags
}
