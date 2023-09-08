package stats

import (
	"github.com/fededp/dbg-go/pkg/root"
	s3utils "github.com/fededp/dbg-go/pkg/utils/s3"
	"log/slog"
	"strings"
)

type s3Statter struct {
	*s3utils.Client
}

func NewS3Statter() (Statter, error) {
	client, err := s3utils.NewClient(true)
	if err != nil {
		return nil, err
	}
	return &s3Statter{Client: client}, nil
}

func (f *s3Statter) Info() string {
	return "gathering stats for remote drivers"
}

func (s *s3Statter) GetDriverStats(opts root.Options) (driverStatsByDriverVersion, error) {
	slog.SetDefault(slog.With("bucket", s3utils.S3Bucket))

	driverStatsByVersion := make(driverStatsByDriverVersion)
	err := s.LoopFiltered(opts, "computing stats", "key", func(driverVersion, key string) error {
		dStats := driverStatsByVersion[driverVersion]
		if strings.HasSuffix(key, ".ko") {
			dStats.NumModules++
		} else if strings.HasSuffix(key, ".o") {
			dStats.NumProbes++
		}
		driverStatsByVersion[driverVersion] = dStats
		return nil
	})
	return driverStatsByVersion, err
}
