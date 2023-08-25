package validate

import (
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func generateConfigFile(dkConf DriverkitYaml, confName string) (func(), error) {
	data, err := yaml.Marshal(dkConf)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(confName, data, 0644)
	if err != nil {
		return nil, err
	}
	return func() {
		_ = os.Remove(confName)
	}, nil
}

func TestValidateConfig(t *testing.T) {
	// Normal options
	opts := Options{
		Options: root.Options{
			Architecture:  "x86_64",
			DriverVersion: []string{"1.0.0+driver"},
		},
		DriverName: "falco",
	}
	namedDriverOpts := Options{
		Options: root.Options{
			Architecture:  "x86_64",
			DriverVersion: []string{"2.0.0+driver"},
		},
		DriverName: "TEST",
	}

	tests := map[string]struct {
		opts          Options
		dkConf        DriverkitYaml
		confName      string
		errorExpected bool
	}{
		"correct config": {
			opts: opts,
			dkConf: DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.10.0",
				Target:        "centos",
				Architecture:  "amd64",
				Output: DriverkitYamlOutputs{
					Module: fmt.Sprintf(OutputPathFmt+".ko",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
					Probe: fmt.Sprintf(OutputPathFmt+".o",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
				},
				KernelUrls:       nil,
				KernelConfigData: "test",
			},
			confName:      "centos_5.10.0_1.yaml",
			errorExpected: false,
		},
		"correct config with custom driver name": {
			opts: namedDriverOpts,
			dkConf: DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.10.0",
				Target:        "centos",
				Architecture:  "amd64",
				Output: DriverkitYamlOutputs{
					Module: fmt.Sprintf(OutputPathFmt+".ko",
						namedDriverOpts.DriverVersion[0],
						namedDriverOpts.Architecture,
						namedDriverOpts.DriverName,
						"centos",
						"5.10.0",
						"1"),
					Probe: fmt.Sprintf(OutputPathFmt+".o",
						namedDriverOpts.DriverVersion[0],
						namedDriverOpts.Architecture,
						namedDriverOpts.DriverName,
						"centos",
						"5.10.0",
						"1"),
				},
				KernelUrls:       nil,
				KernelConfigData: "test",
			},
			confName:      "centos_5.10.0_1.yaml",
			errorExpected: false,
		},
		"wrong arch config": {
			opts: opts,
			dkConf: DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.10.0",
				Target:        "centos",
				Architecture:  "arm64", // arm64 config running in x86_64 mode
				Output: DriverkitYamlOutputs{
					Module: fmt.Sprintf(OutputPathFmt+".ko",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
					Probe: fmt.Sprintf(OutputPathFmt+".o",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
				},
				KernelUrls:       nil,
				KernelConfigData: "test",
			},
			confName:      "centos_5.10.0_1.yaml",
			errorExpected: true,
		},
		"wrong name config": {
			opts: opts,
			dkConf: DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.10.0",
				Target:        "centos",
				Architecture:  "amd64",
				Output: DriverkitYamlOutputs{
					Module: fmt.Sprintf(OutputPathFmt+".ko",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
					Probe: fmt.Sprintf(OutputPathFmt+".o",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
				},
				KernelUrls:       nil,
				KernelConfigData: "test",
			},
			confName:      "centos_WRONG_5.10.0_1.yaml",
			errorExpected: true,
		},
		"wrong arch in config outputs": {
			opts: opts,
			dkConf: DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.10.0",
				Target:        "centos",
				Architecture:  "amd64",
				Output: DriverkitYamlOutputs{
					Module: fmt.Sprintf(OutputPathFmt+".ko",
						opts.DriverVersion[0],
						"WRONGARCH",
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
					Probe: fmt.Sprintf(OutputPathFmt+".o",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
				},
				KernelUrls:       nil,
				KernelConfigData: "test",
			},
			confName:      "centos_5.10.0_1.yaml",
			errorExpected: true,
		},
		"wrong path in config outputs": {
			opts: opts,
			dkConf: DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.10.0",
				Target:        "centos",
				Architecture:  "amd64",
				Output: DriverkitYamlOutputs{
					Module: fmt.Sprintf(OutputPathFmt+".ko",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"WRONGTARGET",
						"5.10.0",
						"1"),
					Probe: fmt.Sprintf(OutputPathFmt+".o",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
				},
				KernelUrls:       nil,
				KernelConfigData: "test",
			},
			confName:      "centos_5.10.0_1.yaml",
			errorExpected: true,
		},
		"wrong suffix in config output module": {
			opts: opts,
			dkConf: DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.10.0",
				Target:        "centos",
				Architecture:  "amd64",
				Output: DriverkitYamlOutputs{
					Module: fmt.Sprintf(OutputPathFmt+".kooo",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
					Probe: fmt.Sprintf(OutputPathFmt+".o",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
				},
				KernelUrls:       nil,
				KernelConfigData: "test",
			},
			confName:      "centos_5.10.0_1.yaml",
			errorExpected: true,
		},
		"wrong suffix in config output probe": {
			opts: opts,
			dkConf: DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.10.0",
				Target:        "centos",
				Architecture:  "amd64",
				Output: DriverkitYamlOutputs{
					Module: fmt.Sprintf(OutputPathFmt+".ooo",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
					Probe: fmt.Sprintf(OutputPathFmt+".o",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
				},
				KernelUrls:       nil,
				KernelConfigData: "test",
			},
			confName:      "centos_5.10.0_1.yaml",
			errorExpected: true,
		},
		"no kernelurls nor kernelconfigdata set in config": {
			opts: opts,
			dkConf: DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.10.0",
				Target:        "centos",
				Architecture:  "amd64",
				Output: DriverkitYamlOutputs{
					Module: fmt.Sprintf(OutputPathFmt+".ooo",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
					Probe: fmt.Sprintf(OutputPathFmt+".o",
						opts.DriverVersion[0],
						opts.Architecture,
						opts.DriverName,
						"centos",
						"5.10.0",
						"1"),
				},
				KernelUrls:       nil,
				KernelConfigData: "",
			},
			confName:      "centos_5.10.0_1.yaml",
			errorExpected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cleanup, err := generateConfigFile(test.dkConf, test.confName)
			assert.NoError(t, err)
			t.Cleanup(cleanup)
			err = validateConfig(test.opts, test.opts.DriverVersion[0], test.confName)
			if test.errorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
