package cleanup

import (
	"log/slog"
)

func Run(opts Options, cleaner Cleaner) error {
	var err error
	slog.Info(cleaner.Info())
	for _, driverVersion := range opts.DriverVersion {
		if opts.KernelVersion == "" && opts.KernelRelease == "" && opts.Distro == "" {
			// No filters, remove entire folder/all bucket keys
			err = cleaner.CleanupAll(opts, driverVersion)
		} else {
			err = cleaner.Cleanup(opts, driverVersion)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
