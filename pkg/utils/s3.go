package utils

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	S3Bucket = "falco-distribution"
	S3Region = "eu-west-1"
)

func NewS3Client(readOnly bool, awsProfile string) (*s3.Client, error) {
	var (
		cfg aws.Config
		err error
	)
	if !readOnly {
		cfg, err = config.LoadDefaultConfig(context.Background(), config.WithRegion(S3Bucket), config.WithSharedConfigProfile(awsProfile))
		if err != nil {
			return nil, err
		}
	} else {
		cfg = aws.Config{
			Region:      S3Region,
			Credentials: aws.AnonymousCredentials{},
		}
	}
	return s3.NewFromConfig(cfg), nil
}
