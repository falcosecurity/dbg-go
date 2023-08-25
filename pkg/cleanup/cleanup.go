package cleanup

import (
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	logger "log/slog"
	"os"
)

func Run(opts Options) error {
	logger.Info("cleaning up existing config files")
	for _, driverVersion := range opts.DriverVersion {
		configPath := fmt.Sprintf(root.ConfigPathFmt,
			opts.RepoRoot,
			driverVersion,
			opts.Architecture,
			"")
		logger.Info("removing folder", "path", configPath)
		if opts.DryRun {
			logger.Info("skipping because of dry-run.")
			continue
		}
		err := os.RemoveAll(configPath)
		if err != nil {
			return err
		}
	}
	return nil
}
