resource "aws_ecr_repository" "todos_repository" {
  name                 = local.name
  image_tag_mutability = "MUTABLE"
  tags                 = local.tags
}