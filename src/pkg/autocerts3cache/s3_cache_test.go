package autocerts3cache

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
)

func Test_Get_WithError(t *testing.T) {
	mockedS3Api := NewMockedS3Api()

	cache := NewS3Cache("bucket", mockedS3Api)

	ctx := context.TODO()

	input := &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	}

	mockedS3Api.On("GetObject", ctx, input).Return(nil, fmt.Errorf("some error"))

	res, err := cache.Get(ctx, "key")

	assert.Nil(t, res)
	assert.EqualError(t, err, "some error")

	mockedS3Api.AssertExpectations(t)
}

func Test_Get_Ok(t *testing.T) {
	mockedS3Api := NewMockedS3Api()

	cache := NewS3Cache("bucket", mockedS3Api)

	ctx := context.TODO()

	input := &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	}

	output := &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("content"))),
	}

	mockedS3Api.On("GetObject", ctx, input).Return(output, nil)

	res, err := cache.Get(ctx, "key")

	assert.Equal(t, []byte("content"), res)
	assert.Nil(t, err)

	mockedS3Api.AssertExpectations(t)
}

func Test_Put_WithError(t *testing.T) {
	mockedS3Api := NewMockedS3Api()

	cache := NewS3Cache("bucket", mockedS3Api)

	ctx := context.TODO()

	input := &s3.PutObjectInput{
		Bucket:               aws.String("bucket"),
		Key:                  aws.String("key"),
		Body:                 bytes.NewReader([]byte("content")),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	}

	mockedS3Api.On("PutObject", ctx, input).Return(nil, fmt.Errorf("some error"))

	err := cache.Put(ctx, "key", []byte("content"))

	assert.EqualError(t, err, "some error")

	mockedS3Api.AssertExpectations(t)
}

func Test_Put_Ok(t *testing.T) {
	mockedS3Api := NewMockedS3Api()

	cache := NewS3Cache("bucket", mockedS3Api)

	ctx := context.TODO()

	input := &s3.PutObjectInput{
		Bucket:               aws.String("bucket"),
		Key:                  aws.String("key"),
		Body:                 bytes.NewReader([]byte("content")),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	}

	mockedS3Api.On("PutObject", ctx, input).Return(nil, nil)

	err := cache.Put(ctx, "key", []byte("content"))

	assert.Nil(t, err)

	mockedS3Api.AssertExpectations(t)
}

func Test_Delete_WithError(t *testing.T) {
	mockedS3Api := NewMockedS3Api()

	cache := NewS3Cache("bucket", mockedS3Api)

	ctx := context.TODO()

	input := &s3.DeleteObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	}

	mockedS3Api.On("DeleteObject", ctx, input).Return(nil, fmt.Errorf("some error"))

	err := cache.Delete(ctx, "key")

	assert.EqualError(t, err, "some error")

	mockedS3Api.AssertExpectations(t)
}

func Test_Delete_Ok(t *testing.T) {
	mockedS3Api := NewMockedS3Api()

	cache := NewS3Cache("bucket", mockedS3Api)

	ctx := context.TODO()

	input := &s3.DeleteObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	}

	mockedS3Api.On("DeleteObject", ctx, input).Return(nil, nil)

	err := cache.Delete(ctx, "key")

	assert.Nil(t, err)

	mockedS3Api.AssertExpectations(t)
}
