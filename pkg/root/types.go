package root

import (
	"github.com/spf13/viper"
	logger "log/slog"
)

type Target struct {
	Distro        string
	KernelRelease string
	KernelVersion string
}

func (t Target) IsSet() bool {
	return t.Distro != "" && t.KernelRelease != "" && t.KernelVersion != ""
}

type Options struct {
	DryRun        bool
	RepoRoot      string
	Architecture  string
	DriverVersion []string
	Target
}

func LoadRootOptions() Options {
	opts := Options{
		DryRun:        viper.GetBool("dry-run"),
		RepoRoot:      viper.GetString("repo-root"),
		Architecture:  viper.GetString("architecture"),
		DriverVersion: viper.GetStringSlice("driver-version"),
		Target: Target{
			Distro:        viper.GetString("target-distro"),
			KernelRelease: viper.GetString("target-kernelrelease"),
			KernelVersion: viper.GetString("target-kernelversion"),
		},
	}
	logger.Debug("loaded root options", "opts", opts)
	return opts
}
