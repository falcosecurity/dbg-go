package validate

import (
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
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
			err = validateConfig(test.confName, test.opts.Architecture, test.opts.DriverName, test.opts.DriverVersion[0])
			if test.errorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type DkConfigNamed struct {
	DriverkitYaml
	confPath string
}

func TestValidateConfigFiltered(t *testing.T) {
	configPath := fmt.Sprintf(root.ConfigPathFmt,
		"./test/",
		"1.0.0+driver",
		"x86_64",
		"")

	dkConfigs := []DkConfigNamed{
		{
			DriverkitYaml: DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.10.0",
				Target:        "centos",
				Architecture:  "amd64",
				Output: DriverkitYamlOutputs{
					Module: fmt.Sprintf(OutputPathFmt+".ko",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"centos",
						"5.10.0",
						"1"),
					Probe: fmt.Sprintf(OutputPathFmt+".o",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"centos",
						"5.10.0",
						"1"),
				},
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			confPath: configPath + "centos_5.10.0_1.yaml",
		},
		{
			DriverkitYaml: DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.15.0",
				Target:        "centos",
				Architecture:  "amd64",
				Output: DriverkitYamlOutputs{
					Module: fmt.Sprintf(OutputPathFmt+".ko",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"centos",
						"5.15.0",
						"1"),
					Probe: fmt.Sprintf(OutputPathFmt+".o",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"centos",
						"5.15.0",
						"1"),
				},
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			confPath: configPath + "centos_5.15.0_1.yaml",
		},
		{
			DriverkitYaml: DriverkitYaml{
				KernelVersion: "13",
				KernelRelease: "5.15.0",
				Target:        "ubuntu",
				Architecture:  "amd64",
				Output: DriverkitYamlOutputs{
					Module: fmt.Sprintf(OutputPathFmt+".ko",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"ubuntu",
						"5.15.0",
						"13"),
					Probe: fmt.Sprintf(OutputPathFmt+".o",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"ubuntu",
						"5.15.0",
						"13"),
				},
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			confPath: configPath + "ubuntu_5.15.0_13.yaml",
		},
		{
			DriverkitYaml: DriverkitYaml{
				KernelVersion: "1",
				KernelRelease: "5.15.25",
				Target:        "bottlerocket",
				Architecture:  "amd64",
				Output: DriverkitYamlOutputs{
					Module: fmt.Sprintf(OutputPathFmt+".ko",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"bottlerocket",
						"5.15.25",
						"1"),
					Probe: fmt.Sprintf(OutputPathFmt+".o",
						"1.0.0+driver",
						"x86_64",
						"falco",
						"bottlerocket",
						"5.15.25",
						"1"),
				},
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			confPath: configPath + "bottlerocket_5.15.25_1.yaml",
		},
	}

	tests := map[string]struct {
		opts                   Options
		expectedOutputContains []string
	}{
		"validate all": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "",
						KernelRelease: "",
						KernelVersion: "",
					},
				},
				DriverName: "falco",
			},
			expectedOutputContains: []string{"centos_5.10", "centos_5.15", "ubuntu_5.15", "bottlerocket_5.15"},
		},
		"validate centos only": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "CentOS",
						KernelRelease: "",
						KernelVersion: "",
					},
				},
				DriverName: "falco",
			},
			expectedOutputContains: []string{"centos_5.10", "centos_5.15"},
		},
		"validate centos regex only": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "Cent*",
						KernelRelease: "",
						KernelVersion: "",
					},
				},
				DriverName: "falco",
			},
			expectedOutputContains: []string{"centos_5.10", "centos_5.15"},
		},
		"validate bottlerocket only": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "BottleRocket",
						KernelRelease: "",
						KernelVersion: "",
					},
				},
				DriverName: "falco",
			},
			expectedOutputContains: []string{"bottlerocket_5.15"},
		},
		"validate 5.15 kernelrelease only": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "",
						KernelRelease: "5.15.*",
						KernelVersion: "",
					},
				},
				DriverName: "falco",
			},
			expectedOutputContains: []string{"centos_5.15", "ubuntu_5.15", "bottlerocket_5.15"},
		},
		"validate kernelversion 1 only": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "",
						KernelRelease: "",
						KernelVersion: "1",
					},
				},
				DriverName: "falco",
			},
			expectedOutputContains: []string{"centos_5.10", "centos_5.15", "bottlerocket_5.15"},
		},
	}

	err := os.MkdirAll(configPath, 0700)
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll("./test/")
	})

	for _, dkConfig := range dkConfigs {
		cleanup, err := generateConfigFile(dkConfig.DriverkitYaml, dkConfig.confPath)
		assert.NoError(t, err)
		t.Cleanup(cleanup)
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			type MessageJSON struct {
				Config string `json:"config,omitempty"`
			}
			var messageJSON MessageJSON
			found := 0
			lines := 0
			utils.RunTestParsingLogs(t,
				func() error {
					return Run(test.opts)
				},
				&messageJSON,
				func() bool {
					if messageJSON.Config == "" {
						return true // go on
					}
					lines++
					for _, expectedOutput := range test.expectedOutputContains {
						if strings.Contains(messageJSON.Config, expectedOutput) {
							found++
							break
						}
					}
					return true
				})
			if found != lines {
				t.Fail()
			}
		})
	}
}
