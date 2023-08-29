package cleanup

import (
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	logger "log/slog"
	"os"
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
	opts.DriverVersion = []string{driverVersion} // locally overwrite driverVersions to only match current driverVersion
	return root.LoopConfigsFiltered(opts.Options, "removing file", func(driverVersion, configPath string) error {
		return os.Remove(configPath)
	})
}
