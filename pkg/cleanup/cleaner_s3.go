package cleanup

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fededp/dbg-go/pkg/utils"
	"log/slog"
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

func (s *s3Cleaner) Info() string {
	return "cleaning up remote driver files"
}

func (s *s3Cleaner) Cleanup(opts Options, driverVersion string) error {
	err := utils.LoopBucketFiltered(s.client, opts.Options, driverVersion, func(key string) error {
		slog.Info("cleaning up remote driver file", "key", key)
		if opts.DryRun {
			slog.Info("skipping because of dry-run.")
			return nil
		}
		return s.removeKey(key)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *s3Cleaner) CleanupAll(opts Options, driverVersion string) error {
	return s.Cleanup(opts, driverVersion)
}

func (s *s3Cleaner) removeKey(key string) error {
	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(utils.S3Bucket),
		Key:    aws.String(key),
	})
	return err
}
