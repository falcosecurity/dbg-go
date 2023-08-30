package stats

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/utils"
	"log/slog"
	"strings"
)

type s3Statter struct {
	client *s3.Client
}

func NewS3Statter() (Statter, error) {
	client, err := utils.NewS3Client(true, "UNNEEDED")
	if err != nil {
		return nil, err
	}
	return &s3Statter{client: client}, nil
}

func (s *s3Statter) GetDriverStats(opts root.Options) (driverStatsByDriverVersion, error) {
	slog.SetDefault(slog.With("bucket", utils.S3Bucket))

	slog.Info("fetching stats for remote drivers")
	driverStatsByVersion := make(driverStatsByDriverVersion)

	for _, driverVersion := range opts.DriverVersion {
		dStats := driverStatsByVersion[driverVersion]
		err := utils.LoopBucketFiltered(s.client, opts, driverVersion, func(key string) error {
			slog.Info("computing stats", "key", key)
			if opts.DryRun {
				slog.Info("skipping because of dry-run.")
				return nil
			}
			if strings.HasSuffix(key, ".ko") {
				dStats.NumModules++
			} else if strings.HasSuffix(key, ".o") {
				dStats.NumProbes++
			}
			return nil
		})
		if err != nil {
			return driverStatsByVersion, err
		}
		driverStatsByVersion[driverVersion] = dStats
	}
	return driverStatsByVersion, nil
}
