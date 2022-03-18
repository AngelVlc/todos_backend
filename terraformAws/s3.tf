resource "aws_s3_bucket" "bucket" {
  bucket = local.name
  tags = local.tags
}

resource "aws_s3_bucket_acl" "bucket_acl" {
  bucket = aws_s3_bucket.bucket.id
  acl    = "authenticated-read"
}