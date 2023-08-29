package root

import (
	"fmt"
	"github.com/spf13/viper"
	logger "log/slog"
	"strings"
)

type Target struct {
	Distro        string
	KernelRelease string
	KernelVersion string
}

func (t Target) IsSet() bool {
	return t.Distro != "" && t.KernelRelease != "" && t.KernelVersion != ""
}

func (t Target) toGlob() string {
	// Empty filters fallback at ".*" since we are using a regex match below
	if t.Distro == "" {
		t.Distro = "*"
	} else {
		dkDistro, found := SupportedDistros[t.Distro]
		if found {
			// Filenames use driverkit lowercase target, instead of the kernel-crawler naming.
			t.Distro = dkDistro
		} else {
			// Perhaps a regex? ToLower and pray
			t.Distro = strings.ToLower(t.Distro)
		}
	}
	if t.KernelRelease == "" {
		t.KernelRelease = "*"
	}
	if t.KernelVersion == "" {
		t.KernelVersion = "*"
	}
	return fmt.Sprintf("%s_%s_%s.yaml", t.Distro, t.KernelRelease, t.KernelVersion)
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
