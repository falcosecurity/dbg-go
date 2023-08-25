package root

import (
	"github.com/spf13/viper"
	logger "log/slog"
)

type Options struct {
	DryRun        bool
	RepoRoot      string
	Architecture  string
	DriverVersion []string
}

func LoadRootOptions() Options {
	opts := Options{
		DryRun:        viper.GetBool("dry-run"),
		RepoRoot:      viper.GetString("repo-root"),
		Architecture:  viper.GetString("architecture"),
		DriverVersion: viper.GetStringSlice("driver-version"),
	}
	logger.Debug("loaded root options", "opts", opts)
	return opts
}
