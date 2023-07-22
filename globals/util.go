package globals

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func S3PutRequest(bucket, key, body string) (*s3.PutObjectOutput, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(body),
	}

	res, err := S3Client.PutObject(context.TODO(), input)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func S3HeadRequest(bucket, key string) (*s3.HeadObjectOutput, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	res, err := S3Client.HeadObject(context.TODO(), input)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func S3ObjectSize(bucket, key string) (uint64, error) {
	res, err := S3HeadRequest(bucket, key)

	if err != nil {
		return 0, err
	}

	return uint64(res.ContentLength), nil
}
