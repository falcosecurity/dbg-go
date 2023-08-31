package utils

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fededp/dbg-go/pkg/root"
	"log/slog"
	"path/filepath"
	"regexp"
)

const (
	S3Bucket = "falco-distribution"
	S3Region = "eu-west-1"
)

var s3DriverNameRegex = regexp.MustCompile(`^falco_(?P<Distro>[a-zA-Z-0-9.0-9]*)_(?P<KernelRelease>.*)_(?P<KernelVersion>.*)(\.o|\.ko)`)

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

func LoopBucketFiltered(client *s3.Client,
	opts root.Options,
	driverVersion string,
	keyProcessor func(key string) error,
) error {
	kDistro := root.KernelCrawlerDistro(opts.Distro)
	dkDistro := kDistro.ToDriverkitDistro()
	opts.Distro = string(dkDistro)

	prefix := filepath.Join("driver", driverVersion, opts.Architecture)
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(S3Bucket),
		Prefix: aws.String(prefix),
	}
	maxKeys := 1000
	p := s3.NewListObjectsV2Paginator(client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
		if v := int32(maxKeys); v != 0 {
			o.Limit = v
		}
	})
	for p.HasMorePages() {
		slog.Debug("fetched a page of objects", "prefix", prefix)
		page, err := p.NextPage(context.TODO())
		if err != nil {
			return err
		}
	keyLoop:
		for _, object := range page.Contents {
			if object.Key == nil {
				continue
			}
			key := filepath.Base(*object.Key)
			matches := s3DriverNameRegex.FindStringSubmatch(key)
			if len(matches) == 0 {
				slog.Warn("skipping key, malformed", "key", key)
				continue
			}
			for i, name := range s3DriverNameRegex.SubexpNames() {
				if i > 0 && i <= len(matches) {
					switch name {
					case "Distro":
						if !opts.DistroFilter(matches[i]) {
							continue keyLoop
						}
					case "KernelRelease":
						if !opts.KernelReleaseFilter(matches[i]) {
							continue keyLoop
						}
					case "KernelVersion":
						if !opts.KernelVersionFilter(matches[i]) {
							continue keyLoop
						}
					}
				}
			}
			err = keyProcessor(filepath.Join(prefix, key))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
