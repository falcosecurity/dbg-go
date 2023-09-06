package s3utils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/fededp/dbg-go/pkg/root"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
)

const (
	S3Bucket             = "falco-distribution"
	s3Region             = "eu-west-1"
	s3DriverNameRegexFmt = `^%s_(?P<Distro>[a-zA-Z-0-9.0-9]*)_(?P<KernelRelease>.*)_(?P<KernelVersion>.*)(\.o|\.ko)`
)

func (cl *Client) LoopDriversFiltered(opts root.Options,
	message, tag string,
	keyProcessor func(driverVersion, key string) error,
) error {
	s3DriverNameRegex := regexp.MustCompile(fmt.Sprintf(s3DriverNameRegexFmt, opts.DriverName))
	for _, driverVersion := range opts.DriverVersion {
		prefix := filepath.Join("driver", driverVersion, opts.Architecture.ToNonDeb())
		params := &s3.ListObjectsV2Input{
			Bucket: aws.String(S3Bucket),
			Prefix: aws.String(prefix),
		}
		maxKeys := 1000
		p := s3.NewListObjectsV2Paginator(cl, params, func(o *s3.ListObjectsV2PaginatorOptions) {
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
				slog.Info(message, tag, key)
				if opts.DryRun {
					slog.Info("skipping because of dry-run.")
					return nil
				}
				err = keyProcessor(driverVersion, filepath.Join(prefix, key))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (cl *Client) HeadDriver(opts root.Options, driverVersion, key string) bool {
	prefix := filepath.Join("driver", driverVersion, opts.Architecture.ToNonDeb())
	fullKey := filepath.Join(prefix, key)
	object, _ := cl.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(S3Bucket),
		Key:    aws.String(fullKey),
	})
	return object != nil
}

func (cl *Client) PutDriver(opts root.Options, driverVersion, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	err = cl.putObject(opts, driverVersion, filepath.Base(path), f)
	_ = f.Close()
	return err
}

func (cl *Client) putObject(opts root.Options, driverVersion, key string, reader io.Reader) error {
	prefix := filepath.Join("driver", driverVersion, opts.Architecture.ToNonDeb())
	fullKey := filepath.Join(prefix, key)
	_, err := cl.Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:               aws.String(S3Bucket),
		Key:                  aws.String(fullKey),
		ACL:                  types.ObjectCannedACLPublicRead,
		Body:                 reader,
		ContentType:          aws.String("binary/octet-stream"),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	})
	return err
}
