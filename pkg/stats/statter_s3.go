package stats

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/utils"
	"log/slog"
	"path/filepath"
	"regexp"
	"strings"
)

var s3DriverNameRegex = regexp.MustCompile(`^falco_(?P<Distro>[a-zA-Z-0-9.0-9]*)_(?P<KernelRelease>.*)_(?P<KernelVersion>.*)(\.o|\.ko)`)

type s3Statter struct {
	client *s3.Client
}

func NewS3Statter() Statter {
	return &s3Statter{client: utils.NewS3Client()}
}

func (s *s3Statter) GetDriverStats(opts root.Options) (driverStatsByDriverVersion, error) {
	slog.SetDefault(slog.With("bucket", utils.S3Bucket))

	slog.Info("computing stats")
	driverStatsByVersion := make(driverStatsByDriverVersion)
	distroFilter := func(distro string) bool {
		matched, _ := regexp.MatchString(opts.Distro, distro)
		return matched
	}

	kernelreleaseFilter := func(kernelrelease string) bool {
		matched, _ := regexp.MatchString(opts.KernelRelease, kernelrelease)
		return matched
	}

	kernelversionFilter := func(kernelversion string) bool {
		matched, _ := regexp.MatchString(opts.KernelVersion, kernelversion)
		return matched
	}

	for _, driverVersion := range opts.DriverVersion {
		dStats := driverStatsByVersion[driverVersion]

		prefix := filepath.Join("driver", driverVersion, opts.Architecture)

		params := &s3.ListObjectsV2Input{
			Bucket: aws.String(utils.S3Bucket),
			Prefix: aws.String(prefix),
		}
		maxKeys := 1000
		p := s3.NewListObjectsV2Paginator(s.client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
			if v := int32(maxKeys); v != 0 {
				o.Limit = v
			}
		})
		for p.HasMorePages() {
			slog.Debug("fetched a page of objects", "prefix", prefix)
			page, err := p.NextPage(context.TODO())
			if err != nil {
				return driverStatsByVersion, err
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
							if !distroFilter(matches[i]) {
								continue keyLoop
							}
						case "KernelRelease":
							if !kernelreleaseFilter(matches[i]) {
								continue keyLoop
							}
						case "KernelVersion":
							if !kernelversionFilter(matches[i]) {
								continue keyLoop
							}
						}
					}
				}

				slog.Info("computing stats", "key", key)
				if strings.HasSuffix(key, ".ko") {
					dStats.NumModules++
				} else if strings.HasSuffix(key, ".o") {
					dStats.NumProbes++
				}
			}
		}
		driverStatsByVersion[driverVersion] = dStats
	}
	return driverStatsByVersion, nil
}
