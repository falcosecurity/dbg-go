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

package generate

import (
	"os"
	"testing"

	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/falcosecurity/dbg-go/pkg/stats"
	testutils "github.com/falcosecurity/dbg-go/pkg/utils/test"
	"github.com/falcosecurity/dbg-go/pkg/validate"
	"github.com/stretchr/testify/assert"
)

func BenchmarkAutogenerate(b *testing.B) {
	testCacheData = true // enable json data caching for subsequent tests
	opts := Options{
		Options: root.Options{
			RepoRoot:      "./test/",
			Architecture:  "amd64",
			DriverVersion: []string{"5.0.1+driver"},
			DriverName:    "falco",
		},
		Auto: true,
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err := Run(opts)
		assert.NoError(b, err)
		_ = os.RemoveAll("./test/")
	}
}

func TestGenerate(t *testing.T) {
	testCacheData = true // enable json data caching for subsequent tests
	tests := map[string]struct {
		opts               Options
		expectError        bool
		expectedMinConfigs int
	}{
		"run in auto mode with any target distro filter on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "amd64",
					DriverVersion: []string{"1.0.0+driver"},
					DriverName:    "falco",
				},
				Auto: true,
			},
			expectError:        false,
			expectedMinConfigs: 1,
		},
		"run in auto mode with target distro filter on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "amd64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro: "centos",
					},
					DriverName: "falco",
				},
				Auto: true,
			},
			expectError:        false,
			expectedMinConfigs: 1,
		},
		"run in auto mode with non existent target distro on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "amd64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro: "WRONG_DISTRO",
					},
					DriverName: "falco",
				},
				Auto: true,
			},
			expectError:        false, // we do not expect any error; no configs will be generated though
			expectedMinConfigs: 0,
		},
		"run in auto mode with single target distro on single driver version with custom driver name": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "amd64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro: "centos",
					},
					DriverName: "CUSTOM",
				},
				Auto: true,
			},
			expectError:        false,
			expectedMinConfigs: 1,
		},
		"run in auto mode with single target distro with regex kernel version on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "amd64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "centos",
						KernelVersion: "^1$",
					},
					DriverName: "falco",
				},
				Auto: true,
			},
			expectError:        false,
			expectedMinConfigs: 1,
		},
		"run in auto mode with single target distro with regex kernel release on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "amd64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "centos",
						KernelRelease: `^5\..+$`,
						KernelVersion: "1",
					},
					DriverName: "falco",
				},
				Auto: true,
			},
			expectError:        false,
			expectedMinConfigs: 1,
		},
		"run in auto mode with regex target distro on single driver version with custom driver name": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "amd64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro: "^cent.*$",
					},
					DriverName: "CUSTOM",
				},
				Auto: true,
			},
			expectError:        false,
			expectedMinConfigs: 1,
		},
		"run with empty target kernel release on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "amd64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "centos",
						KernelVersion: "1",
					},
					DriverName: "falco",
				},
			},
			expectError:        true,
			expectedMinConfigs: 1,
		},
		"run with empty target kernel version on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "amd64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "centos",
						KernelRelease: "5.10.0",
					},
					DriverName: "falco",
				},
			},
			expectError:        true,
			expectedMinConfigs: 1,
		},
		"run with empty target distro on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "amd64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						KernelRelease: "5.10.0",
						KernelVersion: "1",
					},
					DriverName: "falco",
				},
			},
			expectError:        true,
			expectedMinConfigs: 1,
		},
		"run with unsupported target distro on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "amd64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "WRONG_DISTRO",
						KernelRelease: "5.10.0",
						KernelVersion: "1",
					},
					DriverName: "falco",
				},
			},
			expectError:        true,
			expectedMinConfigs: 0,
		},
		// NOTE: the below test is flaky: if debian pulls down the headers, we will break.
		// in case, just update to a newer version.
		"run with target values on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "amd64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "debian",
						KernelRelease: "6.1.38-2-amd64",
						KernelVersion: "1",
					},
					DriverName: "falco",
				},
			},
			expectError:        false,
			expectedMinConfigs: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := testutils.PreCreateFolders(test.opts.Options, test.opts.DriverVersion)
			t.Cleanup(func() {
				_ = os.RemoveAll(test.opts.RepoRoot)
			})
			assert.NoError(t, err)
			err = Run(test.opts)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Ensure configs are generated.
				statter := stats.NewFileStatter()

				// Get stats on generated configurations for all kernels and distros.
				statsOpts := test.opts.Options
				statsOpts.KernelVersion = ""
				statsOpts.KernelRelease = ""
				statsOpts.Distro = ""

				driverStats, err := statter.GetDriverStats(statsOpts)
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, len(driverStats), test.expectedMinConfigs)

				// Validate all generated files
				validateOpts := validate.Options{Options: test.opts.Options}
				err = validate.Run(validateOpts)
				assert.NoError(t, err)
			}
		})
	}
}
