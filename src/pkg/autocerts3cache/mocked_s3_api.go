package autocerts3cache

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/mock"
)

type MockedS3Api struct {
	mock.Mock
}

func NewMockedS3Api() *MockedS3Api {
	return &MockedS3Api{}
}

func (m *MockedS3Api) GetObject(ctx context.Context, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

func (m *MockedS3Api) PutObject(ctx context.Context, params *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

func (m *MockedS3Api) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*s3.DeleteObjectOutput), args.Error(1)
}
