package autocerts3cache

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"golang.org/x/crypto/acme/autocert"
)

type S3Cache struct {
	bucketName string
	s3Api      S3Api
}

var _ autocert.Cache = (*S3Cache)(nil)

func NewS3Cache(bucketName string, s3Api S3Api) *S3Cache {
	return &S3Cache{
		bucketName: bucketName,
		s3Api:      s3Api,
	}
}

func (c *S3Cache) Get(ctx context.Context, key string) ([]byte, error) {
	log.Printf("Autocert S3 Cache Get '%v'\n", key)

	input := &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}

	output, err := c.s3Api.GetObject(ctx, input)
	if err != nil {
		var noSuckKeyError *types.NoSuchKey
		if errors.As(err, &noSuckKeyError) {
			return nil, autocert.ErrCacheMiss
		}

		return nil, err
	}
	defer output.Body.Close()

	return ioutil.ReadAll(output.Body)
}

func (c *S3Cache) Put(ctx context.Context, key string, data []byte) error {
	log.Printf("Autocert S3 Cache Put '%v'\n", key)

	input := &s3.PutObjectInput{
		Bucket:               aws.String(c.bucketName),
		Key:                  aws.String(key),
		Body:                 bytes.NewReader(data),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	}

	_, err := c.s3Api.PutObject(ctx, input)

	return err
}

func (c *S3Cache) Delete(ctx context.Context, key string) error {
	log.Printf("Autocert S3 Cache Delete '%v'\n", key)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}

	_, err := c.s3Api.DeleteObject(ctx, input)

	return err
}
