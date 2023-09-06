package cleanup

import (
	"log/slog"
)

func Run(opts Options, cleaner Cleaner) error {
	slog.Info(cleaner.Info())
	return cleaner.Cleanup(opts)
}
