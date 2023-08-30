package generate

import (
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/utils"
	"github.com/fededp/dbg-go/pkg/validate"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAutogenerate(t *testing.T) {
	testCacheData = true // enable json data caching for subsequent tests
	tests := map[string]struct {
		opts        Options
		expectError bool
	}{
		"run in auto mode without target distro on multiple driver versions": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver", "2.0.0+driver"},
					Target: root.Target{
						Distro:        "", // Should load it from lastDistro kernel crawler file
						KernelRelease: "",
						KernelVersion: "",
					},
				},
				DriverName: "falco",
				Auto:       true,
			},
			expectError: false,
		},
		"run in auto mode with any target distro filter on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "^.*$", // Avoid loading from lastDistro kernel-crawler file, instead force-set any distro
						KernelRelease: "",
						KernelVersion: "",
					},
				},
				DriverName: "falco",
				Auto:       true,
			},
			expectError: false,
		},
		"run in auto mode with target distro filter on single driver version": {
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
				Auto:       true,
			},
			expectError: false,
		},
		"run in auto mode with non existent target distro on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "WRONG_DISTRO",
						KernelRelease: "",
						KernelVersion: "",
					},
				},
				DriverName: "falco",
				Auto:       true,
			},
			expectError: true,
		},
		"run in auto mode with single target distro on single driver version with custom driver name": {
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
				DriverName: "CUSTOM",
				Auto:       true,
			},
			expectError: false,
		},
		"run in auto mode with regex target distro on single driver version with custom driver name": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "^Cent.*$",
						KernelRelease: "",
						KernelVersion: "",
					},
				},
				DriverName: "CUSTOM",
				Auto:       true,
			},
			expectError: false,
		},
		"run with empty target kernel release on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "CentOS",
						KernelRelease: "",
						KernelVersion: "1",
					},
				},
				DriverName: "falco",
			},
			expectError: true,
		},
		"run with empty target kernel version on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "CentOS",
						KernelRelease: "5.10.0",
						KernelVersion: "",
					},
				},
				DriverName: "falco",
			},
			expectError: true,
		},
		"run with empty target distro on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "",
						KernelRelease: "5.10.0",
						KernelVersion: "1",
					},
				},
				DriverName: "falco",
			},
			expectError: true,
		},
		"run with unsupported target distro on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "WRONG_DISTRO",
						KernelRelease: "5.10.0",
						KernelVersion: "1",
					},
				},
				DriverName: "falco",
			},
			expectError: true,
		},
		"run with target values on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
					Target: root.Target{
						Distro:        "CentOS",
						KernelRelease: "5.10.0",
						KernelVersion: "1",
					},
				},
				DriverName: "falco",
			},
			expectError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := utils.PreCreateFolders(test.opts.RepoRoot, test.opts.Architecture, test.opts.DriverVersion)
			t.Cleanup(func() {
				_ = os.RemoveAll(test.opts.RepoRoot)
			})
			assert.NoError(t, err)
			err = Run(test.opts)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Validate all generated files
				validateOpts := validate.Options{Options: test.opts.Options, DriverName: test.opts.DriverName}
				err = validate.Run(validateOpts)
				assert.NoError(t, err)
			}
		})
	}
}
