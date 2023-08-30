package cleanup

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fededp/dbg-go/pkg/utils"
)

type s3Cleaner struct {
	client *s3.Client
}

func NewS3Cleaner(awsProfile string) (Cleaner, error) {
	client, err := utils.NewS3Client(false, awsProfile)
	if err != nil {
		return nil, err
	}
	return &s3Cleaner{client: client}, nil
}

func (s s3Cleaner) Remove(key string) error {
	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(utils.S3Bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s s3Cleaner) RemoveAll(bucket string) error {
	_, err := s.client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	return err
}
