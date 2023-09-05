package validate

import (
	"fmt"
	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
	"github.com/fededp/dbg-go/pkg/root"
	"path/filepath"
	"strings"
)

type Options struct {
	root.Options
}

type DriverkitYamlOutputs struct {
	Module string `yaml:"module"`
	Probe  string `yaml:"probe"`
}

// DriverkitYaml is the driverkit config schema
type DriverkitYaml struct {
	KernelVersion    string               `yaml:"kernelversion" json:"kernelversion"`
	KernelRelease    string               `yaml:"kernelrelease" json:"kernelrelease"`
	Target           string               `yaml:"target" json:"target"`
	Architecture     string               `yaml:"architecture"`
	Output           DriverkitYamlOutputs `yaml:"output"`
	KernelUrls       []string             `yaml:"kernelurls,omitempty" json:"headers"`
	KernelConfigData string               `yaml:"kernelconfigdata,omitempty" json:"kernelconfigdata"`
}

func (dy *DriverkitYaml) ToName() string {
	return fmt.Sprintf("%s_%s_%s", dy.Target, dy.KernelRelease, dy.KernelVersion)
}

func (dy *DriverkitYaml) ToConfigName() string {
	return fmt.Sprintf("%s.yaml", dy.ToName())
}

func (dy *DriverkitYaml) FillOutputs(driverVersion string, opts root.Options) {
	outputPath := root.BuildOutputPath(opts, driverVersion, dy.ToName())
	// Tricky because driverkit configs Outputs assume
	// that the tool is called from the `driverkit` folder of test-infra repo.
	// Only keep last 4 parts, ie: from "output/" onwards
	paths := strings.Split(outputPath, "/")
	configOutputPath := filepath.Join(paths[len(paths)-4:]...)

	kr := kernelrelease.FromString(dy.KernelRelease)
	kr.Architecture = opts.Architecture
	if kr.SupportsModule() {
		dy.Output.Module = configOutputPath + ".ko"
	}
	if kr.SupportsProbe() {
		dy.Output.Probe = configOutputPath + ".o"
	}
}
