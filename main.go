package main

import (
	"github.com/fededp/dbg-go/cmd"
	"log/slog"
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.With("err", err).Error("error executing dbg-go")
	}
}
