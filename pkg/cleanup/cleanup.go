package cleanup

import (
	"log/slog"
)

func Run(opts Options, cleaner Cleaner) error {
	var err error
	slog.Info(cleaner.Info())
	for _, driverVersion := range opts.DriverVersion {
		err = cleaner.Cleanup(opts, driverVersion)
		if err != nil {
			return err
		}
	}
	return nil
}
