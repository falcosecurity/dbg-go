package validate

import (
	"fmt"
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

func (dy *DriverkitYaml) FillOutputs(driverVersion string, opts root.Options) {
	outputPath := root.BuildOutputPath(opts, driverVersion, dy.ToName())
	// Tricky because driverkit configs Outputs assume
	// that the tool is called from the `driverkit` folder of test-infra.
	paths := strings.Split(outputPath, "/")
	configOutputPath := filepath.Join(paths[3:]...)
	dy.Output.Module = configOutputPath + ".ko"
	dy.Output.Probe = configOutputPath + ".o"
}
