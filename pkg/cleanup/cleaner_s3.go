package cleanup

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3utils "github.com/fededp/dbg-go/pkg/utils/s3"
	"log/slog"
)

type s3Cleaner struct {
	*s3utils.Client
}

func NewS3Cleaner(awsProfile string) (Cleaner, error) {
	client, err := s3utils.NewClient(false, awsProfile)
	if err != nil {
		return nil, err
	}
	return &s3Cleaner{Client: client}, nil
}

func (s *s3Cleaner) Info() string {
	return "cleaning up remote driver files"
}

func (s *s3Cleaner) Cleanup(opts Options, driverVersion string) error {
	err := s.LoopDriversFiltered(opts.Options, driverVersion, func(key string) error {
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

func (s *s3Cleaner) removeKey(key string) error {
	_, err := s.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(s3utils.S3Bucket),
		Key:    aws.String(key),
	})
	return err
}
