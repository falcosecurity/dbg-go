package cleanup

import (
	"github.com/fededp/dbg-go/pkg/root"
	"log/slog"
	"os"
)

type fileCleaner struct{}

func NewFileCleaner() Cleaner {
	return &fileCleaner{}
}

func (f *fileCleaner) Info() string {
	return "cleaning up local config files"
}

func (f *fileCleaner) CleanupAll(opts Options, driverVersion string) error {
	configPath := root.BuildConfigPath(opts.Options, driverVersion, "")
	slog.Info("removing folder", "config", configPath)
	if opts.DryRun {
		slog.Info("skipping because of dry-run.")
		return nil
	}
	err := os.RemoveAll(configPath)
	return err
}

func (f *fileCleaner) Cleanup(opts Options, driverVersion string) error {
	opts.DriverVersion = []string{driverVersion} // locally overwrite driverVersions to only match current driverVersion
	return root.LoopConfigsFiltered(opts.Options, "removing file", func(driverVersion, configPath string) error {
		return os.Remove(configPath)
	})
}
