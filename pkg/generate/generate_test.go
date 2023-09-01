package generate

import (
	"github.com/fededp/dbg-go/pkg/root"
	testutils "github.com/fededp/dbg-go/pkg/utils/test"
	"github.com/fededp/dbg-go/pkg/validate"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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

	b.StopTimer()

	for n := 0; n < b.N; n++ {
		b.StartTimer()
		err := Run(opts)
		assert.NoError(b, err)
		b.StopTimer()
		_ = os.RemoveAll("./test/")
	}
}

func TestGenerate(t *testing.T) {
	testCacheData = true // enable json data caching for subsequent tests
	tests := map[string]struct {
		opts        Options
		expectError bool
	}{
		"run in auto mode with loaded from kernel-crawler distro on multiple driver versions": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "amd64",
					DriverVersion: []string{"1.0.0+driver", "2.0.0+driver"},
					Target: root.Target{
						Distro: "load", // Should load it from lastDistro kernel crawler file
					},
					DriverName: "falco",
				},
				Auto: true,
			},
			expectError: false,
		},
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
			expectError: false,
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
			expectError: false,
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
			expectError: false, // we do not expect any error; no configs will be generated though
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
			expectError: false,
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
			expectError: false,
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
			expectError: true,
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
			expectError: true,
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
			expectError: true,
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
			expectError: true,
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
			expectError: false,
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
				// Validate all generated files
				validateOpts := validate.Options{Options: test.opts.Options}
				if validateOpts.Distro == "load" {
					// When we load from last kernel-crawler distro,
					// force validate all
					validateOpts.Distro = ""
				}
				err = validate.Run(validateOpts)
				assert.NoError(t, err)
			}
		})
	}
}
