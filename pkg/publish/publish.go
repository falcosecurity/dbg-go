package publish

import (
	"github.com/fededp/dbg-go/pkg/root"
	s3utils "github.com/fededp/dbg-go/pkg/utils/s3"
	"log/slog"
)

// Used by tests
var testClient *s3utils.Client

func Run(opts Options) error {
	slog.Info("publishing drivers")
	var (
		client *s3utils.Client
		err    error
	)
	if testClient == nil {
		client, err = s3utils.NewClient(false, opts.AwsProfile)
		if err != nil {
			return err
		}
	} else {
		client = testClient
	}
	return root.LoopPathFiltered(opts.Options, root.BuildOutputPath, "publishing", "driver", func(driverVersion, path string) error {
		return client.PutDriver(opts.Options, driverVersion, path)
	})
}
