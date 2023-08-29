package root

import (
	"fmt"
	logger "log/slog"
	"path/filepath"
)

type ConfigLooper func(driverVersion, configPath string) error

func LoopConfigsFiltered(opts Options, message string, worker ConfigLooper) error {
	configNameGlob := opts.Target.toGlob()
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
			logger.Info(message, "config", config)
			if opts.DryRun {
				logger.Info("skipping because of dry-run.")
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
