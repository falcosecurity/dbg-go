package cleanup

import (
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	"log/slog"
)

func Run(opts Options, cleaner Cleaner) error {
	slog.Info("cleaning up existing config files")
	var err error
	for _, driverVersion := range opts.DriverVersion {
		if opts.KernelVersion == "" && opts.KernelRelease == "" && opts.Distro == "" {
			err = cleanupFolder(opts, driverVersion, cleaner)
		} else {
			err = cleanupMatchingConfigs(opts, driverVersion, cleaner)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanupFolder(opts Options, driverVersion string, cleaner Cleaner) error {
	configPath := fmt.Sprintf(root.ConfigPathFmt,
		opts.RepoRoot,
		driverVersion,
		opts.Architecture,
		"")
	slog.Info("removing folder", "config", configPath)
	if opts.DryRun {
		slog.Info("skipping because of dry-run.")
		return nil
	}
	err := cleaner.RemoveAll(configPath)
	return err
}

func cleanupMatchingConfigs(opts Options, driverVersion string, cleaner Cleaner) error {
	opts.DriverVersion = []string{driverVersion} // locally overwrite driverVersions to only match current driverVersion
	return root.LoopConfigsFiltered(opts.Options, "removing file", func(driverVersion, configPath string) error {
		return cleaner.Remove(configPath)
	})
}
