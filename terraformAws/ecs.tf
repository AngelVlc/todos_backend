resource "aws_ecs_cluster" "main" {
  name = "${local.name}-cluster"
  tags = local.tags
}

resource "aws_ecs_task_definition" "main" {
  family                   = "${local.name}-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn
  tags                     = local.tags
  container_definitions = jsonencode([
    {
      name      = "todos-backend"
      image     = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com/${aws_ecr_repository.todos_repository.name}:latest@${data.aws_ecr_image.app.image_digest}"
      essential = true
      environment = [
        { name : "PORT", value : tostring(local.port) },
        { name : "MYSQL_HOST", value : module.aurora_mysql.cluster_endpoint },
        { name : "MYSQL_PORT", value : tostring(module.aurora_mysql.cluster_port) },
        { name : "MYSQL_USER", value : module.aurora_mysql.cluster_master_username },
        { name : "MYSQL_PASSWORD", value : module.aurora_mysql.cluster_master_password },
        { name : "MYSQL_DATABASE", value : local.database_name },
        { name : "JWT_SECRET", value : var.jwt_secret },
        { name : "CORS_ALLOWED_ORIGINS", value : var.cors_allowed_origins },
        { name : "NEW_RELIC_LICENSE_KEY", value : var.new_relic_license_key },
        { name : "HONEYBADGER_API_KEY", value : var.honeybadger_api_key },
        { name : "ENVIRONMENT", value : var.environment },
        { name : "CLEARDB_DATABASE_URL", value : "mysql://${module.aurora_mysql.cluster_master_username}:${module.aurora_mysql.cluster_master_password}@${module.aurora_mysql.cluster_endpoint}/${local.database_name}" },
        { name : "DELETE_EXPIRED_REFRESH_TOKEN_INTERVAL", value : "24h" },
        { name : "DOMAIN", value : var.domain },
        { name : "BUCKET_NAME", value: aws_s3_bucket.bucket.id }
      ]
      portMappings = [{
        protocol      = "tcp"
        containerPort = local.port
        hostPort      = local.port
      }]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-region        = data.aws_region.current.name
          awslogs-group         = local.ecs_log_group_name
          awslogs-stream-prefix = "api"
        }
      }
    },
    {
      name      = "ecs-sidecar"
      image     = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com/ecs-sidecar"
      essential = false
      environment = [
        { name : "CLUSTER_NAME", value : aws_ecs_cluster.main.name },
        { name : "DOMAIN", value : var.domain },
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-region        = data.aws_region.current.name
          awslogs-group         = local.ecs_log_group_name
          awslogs-stream-prefix = "ecs-sidecar"
        }
      }
    }
  ])
}

resource "aws_iam_role" "ecs_task_execution_role" {
  name = "${local.name}-ecsTaskExecutionRole"

  assume_role_policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": "sts:AssumeRole",
     "Principal": {
       "Service": "ecs-tasks.amazonaws.com"
     },
     "Effect": "Allow",
     "Sid": ""
   }
 ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-task-execution-role-policy-attachment" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role" "ecs_task_role" {
  name = "${local.name}-ecsTaskRole"

  assume_role_policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": "sts:AssumeRole",
     "Principal": {
       "Service": "ecs-tasks.amazonaws.com"
     },
     "Effect": "Allow",
     "Sid": ""
   }
 ]
}
EOF
}

resource "aws_iam_policy" "ecs" {
  name        = "${local.name}-task-policy-ecs"
  description = "Policy for ECS"

  policy = <<EOF
{
   "Version": "2012-10-17",
   "Statement": [
       {
           "Effect": "Allow",
           "Action": [
               "ecr:*",
               "ecs:ListTasks",
               "ecs:DescribeTasks",
               "ec2:DescribeNetworkInterfaces",
               "route53:ListHostedZones",
               "route53:ChangeResourceRecordSets",
               "s3:*"
           ],
           "Resource": "*"
       }
   ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "ecs-task-role-policy-attachment" {
  role       = aws_iam_role.ecs_task_role.name
  policy_arn = aws_iam_policy.ecs.arn
}

resource "aws_ecs_service" "main" {
  name                               = "${local.name}-service"
  tags                               = local.tags
  cluster                            = aws_ecs_cluster.main.id
  task_definition                    = aws_ecs_task_definition.main.arn
  desired_count                      = 1
  deployment_minimum_healthy_percent = 100
  deployment_maximum_percent         = 200
  launch_type                        = "FARGATE"
  scheduling_strategy                = "REPLICA"
  force_new_deployment               = true

  network_configuration {
    security_groups  = [aws_security_group.ecs_security_group.id]
    subnets          = data.aws_subnets.default.ids
    assign_public_ip = true
  }

  lifecycle {
    ignore_changes = [desired_count]
  }
}

resource "aws_security_group" "ecs_security_group" {
  name        = "${local.name}-ecs-service-sg"
  description = "Allow traffic to ${local.name}"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description      = "API port"
    from_port        = local.port
    to_port          = local.port
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  ingress {
    description      = "API TLS port"
    from_port        = 443
    to_port          = 443
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    description      = "API port from VPC"
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = local.tags
}
