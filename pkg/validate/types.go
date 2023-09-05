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

type WrongConfigNameErr struct {
	configName         string
	expectedConfigName string
}

func (w *WrongConfigNameErr) Error() string {
	return fmt.Sprintf("config filename is wrong (%s); should be %s", w.configName, w.expectedConfigName)
}

type WrongArchInConfigErr struct {
	configPath string
	arch       string
}

func (w *WrongArchInConfigErr) Error() string {
	return fmt.Sprintf("wrong architecture in config file %s: %s", w.configPath, w.arch)
}

type WrongOutputProbeNameErr struct {
	outputProbeName         string
	expectedOutputProbeName string
}

func (w *WrongOutputProbeNameErr) Error() string {
	return fmt.Sprintf("output probe filename is wrong (%s); expected: %s.o", w.outputProbeName, w.expectedOutputProbeName)
}

type WrongOutputProbeArchErr struct {
	probe string
	arch  string
}

func (w *WrongOutputProbeArchErr) Error() string {
	return fmt.Sprintf("output probe filename has wrong architecture in its path (%s); expected %s", w.probe, w.arch)
}

type WrongOutputModuleNameErr struct {
	outputModuleName         string
	expectedOutputModuleName string
}

func (w *WrongOutputModuleNameErr) Error() string {
	return fmt.Sprintf("output module filename is wrong (%s); expected: %s.o", w.outputModuleName, w.expectedOutputModuleName)
}

type WrongOutputModuleArchErr struct {
	module string
	arch   string
}

func (w *WrongOutputModuleArchErr) Error() string {
	return fmt.Sprintf("output module filename has wrong architecture in its path (%s); expected %s", w.module, w.arch)
}

type KernelConfigDataNotBase64Err struct{}

func (k *KernelConfigDataNotBase64Err) Error() string {
	return fmt.Sprintf("kernelconfigdata must be a base64 encoded string")
}
