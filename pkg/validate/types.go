package validate

import (
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
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

func (dy *DriverkitYaml) ToConfigName() string {
	return fmt.Sprintf("%s_%s_%s.yaml", dy.Target, dy.KernelRelease, dy.KernelVersion)
}

func (dy *DriverkitYaml) ToOutputPath(driverVersion string, opts root.Options) string {
	return fmt.Sprintf(outputPathFmt,
		driverVersion,
		opts.Architecture.ToNonDeb(),
		opts.DriverName,
		dy.Target,
		dy.KernelRelease,
		dy.KernelVersion)
}
