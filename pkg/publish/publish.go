package publish

import (
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	s3utils "github.com/fededp/dbg-go/pkg/utils/s3"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
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
	return LoopDriversFiltered(client, opts.Options)
}

func LoopDriversFiltered(client *s3utils.Client, opts root.Options) error {
	outputNameGlob := opts.Target.ToGlob()
	outputNameGlob = strings.ReplaceAll(outputNameGlob, ".yaml", ".*")
	for _, driverVersion := range opts.DriverVersion {
		driversPath := fmt.Sprintf("output/%s/%s/falco_%s", driverVersion, opts.Architecture.ToNonDeb(), outputNameGlob)
		drivers, err := filepath.Glob(driversPath)
		if err != nil {
			return err
		}
		for _, driver := range drivers {
			slog.Info("publishing", "driver", driver)
			if opts.DryRun {
				slog.Info("skipping because of dry-run.")
				return nil
			}
			f, err := os.Open(driver)
			if err != nil {
				return err
			}
			err = client.PutObject(opts, driverVersion, filepath.Base(driver), f)
			_ = f.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
