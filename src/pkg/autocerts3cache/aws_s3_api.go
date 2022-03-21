package autocerts3cache

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsS3Api struct {
	s3Client *s3.Client
}

func NewAwsS3Api(cfg aws.Config) *AwsS3Api {
	client := s3.NewFromConfig(cfg)

	return &AwsS3Api{s3Client: client}
}

func (a *AwsS3Api) GetObject(ctx context.Context, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return a.s3Client.GetObject(ctx, params)
}

func (a *AwsS3Api) PutObject(ctx context.Context, params *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return a.s3Client.PutObject(ctx, params)
}

func (a *AwsS3Api) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return a.s3Client.DeleteObject(ctx, params)
}
