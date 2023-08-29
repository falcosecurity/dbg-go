package utils

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	S3Bucket = "falco-distribution"
	S3Region = "eu-west-1"
)

func NewS3Client() *s3.Client {
	return s3.New(s3.Options{
		Credentials:     aws.AnonymousCredentials{},
		EndpointOptions: s3.EndpointResolverOptions{},
		Region:          S3Region,
	})
}
