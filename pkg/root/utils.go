package root

import (
	"fmt"
	"log/slog"
	"path/filepath"
)

type ConfigLooper func(driverVersion, configPath string) error

func LoopConfigsFiltered(opts Options, message string, worker ConfigLooper) error {
	configNameGlob := opts.Target.ToGlob()
	for _, driverVersion := range opts.DriverVersion {
		configPath := fmt.Sprintf(ConfigPathFmt,
			opts.RepoRoot,
			driverVersion,
			opts.Architecture,
			configNameGlob)
		configs, err := filepath.Glob(configPath)
		if err != nil {
			return err
		}
		for _, config := range configs {
			slog.Info(message, "config", config)
			if opts.DryRun {
				slog.Info("skipping because of dry-run.")
				return nil
			}
			err = worker(driverVersion, config)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
