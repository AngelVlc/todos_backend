package autocerts3cache

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Api interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}
