package autogenerate

import (
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/fededp/dbg-go/pkg/utils"
	"github.com/fededp/dbg-go/pkg/validate"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func preCreateFolders(opts Options) (func(), error) {
	toBeRemoved := make([]string, 0)
	f := func() {
		for _, path := range toBeRemoved {
			_ = os.RemoveAll(path)
		}
	}
	for _, driverVersion := range opts.DriverVersion {
		configPath := fmt.Sprintf(root.ConfigPathFmt,
			opts.RepoRoot,
			driverVersion,
			opts.Architecture,
			"")
		err := os.MkdirAll(configPath, 0700)
		if err != nil {
			return f, err
		}
		toBeRemoved = append(toBeRemoved, configPath)
	}
	return f, nil
}

func TestAutogenerate(t *testing.T) {
	tests := map[string]struct {
		opts        Options
		expectError bool
	}{
		"run without target on multiple driver versions": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver", "2.0.0+driver"},
				},
				DriverName: "falco",
				Target:     "",
			},
			expectError: false,
		},
		"run with * target on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
				},
				DriverName: "falco",
				Target:     "*",
			},
			expectError: false,
		},
		"run with single target on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
				},
				DriverName: "falco",
				Target:     "CentOS",
			},
			expectError: false,
		},
		"run with non existent target on single driver version": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
				},
				DriverName: "falco",
				Target:     "WRONG_TARGET",
			},
			expectError: true,
		},
		"run with single target on single driver version with custom driver name": {
			opts: Options{
				Options: root.Options{
					RepoRoot:      "./test/",
					Architecture:  "x86_64",
					DriverVersion: []string{"1.0.0+driver"},
				},
				DriverName: "CUSTOM",
				Target:     "CentOS",
			},
			expectError: false,
		},
	}

	// Download most recent json to be used during the test
	url := fmt.Sprintf(urlArchFmt, "x86_64")
	jsonData, err := utils.GetURL(url)
	assert.NoError(t, err)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cleanup, err := preCreateFolders(test.opts)
			t.Cleanup(cleanup)
			assert.NoError(t, err)
			err = generateConfigs(test.opts, jsonData)
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
