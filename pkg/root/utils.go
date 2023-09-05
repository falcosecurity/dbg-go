package root

import (
	"fmt"
	"log/slog"
	"path/filepath"
)

type PathWorker func(driverVersion, path string) error

type PathBuilder func(opts Options, driverVersion, configName string) string

func LoopPathFiltered(opts Options, pathBuilder PathBuilder, message, tag string, worker PathWorker) error {
	configNameGlob := opts.Target.toGlob()
	for _, driverVersion := range opts.DriverVersion {
		path := pathBuilder(opts, driverVersion, configNameGlob)
		files, err := filepath.Glob(path)
		if err != nil {
			return err
		}
		for _, file := range files {
			slog.Info(message, tag, file)
			if opts.DryRun {
				slog.Info("skipping because of dry-run.")
				return nil
			}
			err = worker(driverVersion, file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func BuildConfigPath(opts Options, driverVersion, configName string) string {
	return fmt.Sprintf(configPathFmt,
		opts.RepoRoot,
		driverVersion,
		opts.Architecture.ToNonDeb(),
		configName)
}

func BuildOutputPath(opts Options, driverVersion, outputName string) string {
	fullName := ""
	if outputName != "" {
		// only add "drivername_" prefix when outputName is not empty,
		// ie: when we are not generating a folder path.
		fullName = opts.DriverName + "_" + outputName
	}

	return fmt.Sprintf(outputPathFmt,
		opts.RepoRoot,
		driverVersion,
		opts.Architecture.ToNonDeb(),
		fullName)
}
