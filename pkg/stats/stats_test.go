// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package stats

import (
	"os"
	"testing"

	"github.com/falcosecurity/dbg-go/pkg/root"
	testutils "github.com/falcosecurity/dbg-go/pkg/utils/test"
	"github.com/falcosecurity/dbg-go/pkg/validate"
	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type DkConfigNamed struct {
	validate.DriverkitYaml
	hasProbe  bool
	hasModule bool
}

func TestStats(t *testing.T) {
	configPath := root.BuildConfigPath(root.Options{
		RepoRoot:     "./test/",
		Architecture: "amd64",
	}, "1.0.0+driver", "")

	dkConfigs := []DkConfigNamed{
		{
			DriverkitYaml: validate.DriverkitYaml{
				KernelVersion:    "1",
				KernelRelease:    "5.10.0",
				Target:           "centos",
				Architecture:     "amd64",
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			hasProbe:  true,
			hasModule: true,
		},
		{
			DriverkitYaml: validate.DriverkitYaml{
				KernelVersion:    "1",
				KernelRelease:    "5.15.0",
				Target:           "centos",
				Architecture:     "amd64",
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			hasProbe:  true,
			hasModule: true,
		},
		{
			DriverkitYaml: validate.DriverkitYaml{
				KernelVersion:    "13",
				KernelRelease:    "5.15.0",
				Target:           "ubuntu",
				Architecture:     "amd64",
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			hasModule: true,
		},
		{
			DriverkitYaml: validate.DriverkitYaml{
				KernelVersion:    "1",
				KernelRelease:    "5.15.25",
				Target:           "bottlerocket",
				Architecture:     "amd64",
				KernelConfigData: "aaaa", // just to avoid failing validation
			},
			hasProbe:  true,
			hasModule: true,
		},
	}

	err := os.MkdirAll(configPath, 0700)
	t.Cleanup(func() {
		_ = os.RemoveAll("./test")
	})
	assert.NoError(t, err)

	// Create all configs needed by the test
	for _, dkConf := range dkConfigs {
		file, err := os.OpenFile(configPath+dkConf.ToConfigName(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		assert.NoError(t, err)

		dkConf.FillOutputs("1.0.0+driver",
			root.Options{
				DriverName:   "falco",
				Architecture: kernelrelease.Architecture(dkConf.Architecture),
			})
		// Remove when test requires it
		if !dkConf.hasModule {
			dkConf.DriverkitYaml.Output.Module = ""
		}
		if !dkConf.hasProbe {
			dkConf.DriverkitYaml.Output.Probe = ""
		}
		enc := yaml.NewEncoder(file)
		err = enc.Encode(dkConf.DriverkitYaml)
		_ = file.Close()
		assert.NoError(t, err)
	}

	tests := map[string]struct {
		opts          Options
		expectedStats driverStats
	}{
		"stats 1.0.0+driver x86_64": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				DriverName:    "falco",
			}},
			expectedStats: driverStats{
				NumProbes:  3,
				NumModules: 4,
			},
		},
		"stats 2.0.0+driver x86_64": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"2.0.0+driver"}, // not present
				DriverName:    "falco",
			}},
			expectedStats: driverStats{
				NumProbes:  0,
				NumModules: 0,
			},
		},
		"stats 1.0.0+driver aarch64": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "arm64", // not present
				DriverVersion: []string{"1.0.0+driver"},
				DriverName:    "falco",
			}},
			expectedStats: driverStats{
				NumProbes:  0,
				NumModules: 0,
			},
		},
		"stats 1.0.0+driver x86_64 filtered by distro": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				DriverName:    "falco",
				Target: root.Target{
					Distro: "centos",
				},
			}},
			expectedStats: driverStats{
				NumProbes:  2,
				NumModules: 2,
			},
		},
		"stats 1.0.0+driver x86_64 filtered by distro regex": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				DriverName:    "falco",
				Target: root.Target{
					Distro: "cent*",
				},
			}},
			expectedStats: driverStats{
				NumProbes:  2,
				NumModules: 2,
			},
		},
		"stats 1.0.0+driver x86_64 filtered by kernel release": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				DriverName:    "falco",
				Target: root.Target{
					KernelRelease: "5.10.*",
				},
			}},
			expectedStats: driverStats{
				NumProbes:  1,
				NumModules: 1,
			},
		},
		"stats 1.0.0+driver x86_64 filtered by kernel version": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				DriverName:    "falco",
				Target: root.Target{
					KernelVersion: "1",
				},
			}},
			expectedStats: driverStats{
				NumProbes:  3,
				NumModules: 3,
			},
		},
	}

	statter := NewFileStatter()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			driverStatsByVersion, err := statter.GetDriverStats(test.opts.Options)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedStats, driverStatsByVersion["1.0.0+driver"])
		})
	}
}

func TestStatsS3(t *testing.T) {
	keysToBeCreated := []string{
		"driver/1.0.0+driver/x86_64/falco_almalinux_5.14.0-284.11.1.el9_2.x86_64_1.ko",
		"driver/1.0.0+driver/x86_64/falco_amazonlinux2022_5.10.96-90.460.amzn2022.x86_64_1.o",
		"driver/1.0.0+driver/x86_64/falco_debian_6.3.11-1-amd64_1.o",
		"driver/1.0.0+driver/x86_64/falco_debian_6.3.11-1-amd64_1.ko",
		"driver/2.0.0+driver/x86_64/falco_almalinux_5.14.0-284.11.1.el9_2.x86_64_1.ko",
		"driver/2.0.0+driver/aarch64/falco_almalinux_4.18.0-477.10.1.el8_8.aarch64_1.ko",
		"driver/2.0.0+driver/aarch64/falco_bottlerocket_5.10.165_1_1.13.1-aws.o",
	}
	client := testutils.S3CreateTestBucket(t, keysToBeCreated)
	statter := s3Statter{Client: client}

	tests := map[string]struct {
		opts          Options
		expectedStats driverStatsByDriverVersion
	}{
		"stats 1.0.0+driver, 2.0.0+driver x86_64": {
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver", "2.0.0+driver"},
				DriverName:    "falco",
			}},
			expectedStats: driverStatsByDriverVersion{
				"1.0.0+driver": {
					NumProbes:  2,
					NumModules: 2,
				},
				"2.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
			},
		},
		"stats 2.0.0+driver aarch64": {
			opts: Options{Options: root.Options{
				Architecture:  "arm64",
				DriverVersion: []string{"2.0.0+driver"},
				DriverName:    "falco",
			}},
			expectedStats: driverStatsByDriverVersion{
				"2.0.0+driver": {
					NumProbes:  1,
					NumModules: 1,
				},
			},
		},
		"stats 2.0.0+driver x86_64": {
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"2.0.0+driver"},
				DriverName:    "falco",
			}},
			expectedStats: driverStatsByDriverVersion{
				"2.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
			},
		},
		"stats 1.0.0+driver, 2.0.0+driver x86_64 filtered by distro": {
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver", "2.0.0+driver"},
				DriverName:    "falco",
				Target: root.Target{
					Distro: "almalinux",
				},
			}},
			expectedStats: driverStatsByDriverVersion{
				"1.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
				"2.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
			},
		},
		"stats 1.0.0+driver, 2.0.0+driver x86_64 filtered by distro regex": {
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver", "2.0.0+driver"},
				DriverName:    "falco",
				Target: root.Target{
					Distro: "almali*",
				},
			}},
			expectedStats: driverStatsByDriverVersion{
				"1.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
				"2.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
			},
		},
		"stats 1.0.0+driver, 2.0.0+driver x86_64 filtered by kernelrelease": {
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver", "2.0.0+driver"},
				DriverName:    "falco",
				Target: root.Target{
					KernelRelease: "5.*",
				},
			}},
			expectedStats: driverStatsByDriverVersion{
				"1.0.0+driver": {
					NumProbes:  1,
					NumModules: 1,
				},
				"2.0.0+driver": {
					NumProbes:  0,
					NumModules: 1,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dStats, err := statter.GetDriverStats(test.opts.Options)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedStats, dStats)
		})
	}
}
