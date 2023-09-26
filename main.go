package main

import (
	"github.com/falcosecurity/dbg-go/cmd"
	"log/slog"
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.With("err", err).Error("error executing dbg-go")
	}
}
