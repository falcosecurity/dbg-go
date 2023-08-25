package main

import (
	"github.com/fededp/dbg-go/cmd"
	logger "log/slog"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logger.With("err", err).Error("error executing dbg-go")
	}
}
