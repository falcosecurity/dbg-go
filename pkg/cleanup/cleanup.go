package cleanup

import (
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/validate"
	logger "log/slog"
	"os"
	"path/filepath"
)

func Run(opts Options) error {
	logger.Info("cleaning up existing config files")
	var err error
	for _, driverVersion := range opts.DriverVersion {
		if opts.KernelVersion == "" && opts.KernelRelease == "" && opts.Distro == "" {
			err = cleanupFolder(opts, driverVersion)
		} else {
			err = cleanupMatchingConfigs(opts, driverVersion)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanupFolder(opts Options, driverVersion string) error {
	configPath := fmt.Sprintf(root.ConfigPathFmt,
		opts.RepoRoot,
		driverVersion,
		opts.Architecture,
		"")
	logger.Info("removing folder", "path", configPath)
	if opts.DryRun {
		logger.Info("skipping because of dry-run.")
		return nil
	}
	err := os.RemoveAll(configPath)
	return err
}

func cleanupMatchingConfigs(opts Options, driverVersion string) error {
	configNameGlob := validate.ConfGlobFromDistro(opts.Target)

	configPath := fmt.Sprintf(root.ConfigPathFmt,
		opts.RepoRoot,
		driverVersion,
		opts.Architecture,
		configNameGlob)

	files, err := filepath.Glob(configPath)
	if err != nil {
		return err
	}
	for _, f := range files {
		logger.Info("removing file", "path", f)
		if opts.DryRun {
			logger.Info("skipping because of dry-run.")
			continue
		}
		if err = os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}
