package root

import (
	"fmt"
	"github.com/falcosecurity/driverkit/pkg/driverbuilder/builder"
	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
	"github.com/spf13/viper"
	"log/slog"
	"regexp"
)

type Target struct {
	Distro        builder.Type
	KernelRelease string
	KernelVersion string
}

func (t Target) IsSet() bool {
	return t.Distro != "" && t.KernelRelease != "" && t.KernelVersion != ""
}

func (t Target) ToGlob() string {
	// Empty filters fallback at ".*" since we are using a regex match below
	if t.Distro == "" {
		t.Distro = "*"
	}
	if t.KernelRelease == "" {
		t.KernelRelease = "*"
	}
	if t.KernelVersion == "" {
		t.KernelVersion = "*"
	}
	return fmt.Sprintf("%s_%s_%s.yaml", t.Distro, t.KernelRelease, t.KernelVersion)
}

func (t Target) DistroFilter(distro string) bool {
	matched, _ := regexp.MatchString(t.Distro.String(), distro)
	// check if key is actually supported
	if matched {
		_, ok := SupportedDistros[builder.Type(distro)]
		return ok
	}
	return matched
}

func (t Target) KernelReleaseFilter(kernelrelease string) bool {
	matched, _ := regexp.MatchString(t.KernelRelease, kernelrelease)
	return matched
}

func (t Target) KernelVersionFilter(kernelversion string) bool {
	matched, _ := regexp.MatchString(t.KernelVersion, kernelversion)
	return matched
}

type Options struct {
	DryRun        bool
	RepoRoot      string
	Architecture  kernelrelease.Architecture
	DriverVersion []string
	Target
}

func LoadRootOptions() Options {
	opts := Options{
		DryRun:        viper.GetBool("dry-run"),
		RepoRoot:      viper.GetString("repo-root"),
		Architecture:  kernelrelease.Architecture(viper.GetString("architecture")),
		DriverVersion: viper.GetStringSlice("driver-version"),
		Target: Target{
			Distro:        builder.Type(viper.GetString("target-distro")),
			KernelRelease: viper.GetString("target-kernelrelease"),
			KernelVersion: viper.GetString("target-kernelversion"),
		},
	}
	slog.Debug("loaded root options", "opts", opts)
	return opts
}
