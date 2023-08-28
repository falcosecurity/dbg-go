package cleanup

import (
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func preCreateFolders(opts Options, driverVersionsToBeCreated []string) error {
	for _, driverVersion := range driverVersionsToBeCreated {
		configPath := fmt.Sprintf(root.ConfigPathFmt,
			opts.RepoRoot,
			driverVersion,
			opts.Architecture,
			"")
		err := os.MkdirAll(configPath, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

// difference returns the elements in `a` that aren't in `b`.
func sliceDifference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func TestCleanup(t *testing.T) {
	tests := map[string]struct {
		opts                          Options
		driverVersionsToBeCreated     []string
		errorExpected                 bool
		driverFolderRemainingExpected []string
	}{
		"delete all": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "x86_64",
				DriverVersion: []string{"1.0.0+driver", "2.0.0+driver"},
			}},
			driverVersionsToBeCreated:     []string{"1.0.0+driver", "2.0.0+driver"},
			errorExpected:                 false,
			driverFolderRemainingExpected: nil,
		},
		"delete only one": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "x86_64",
				DriverVersion: []string{"1.0.0+driver"},
			}},
			driverVersionsToBeCreated:     []string{"1.0.0+driver", "2.0.0+driver"},
			errorExpected:                 false,
			driverFolderRemainingExpected: []string{"2.0.0+driver"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := preCreateFolders(test.opts, test.driverVersionsToBeCreated)
			t.Cleanup(func() {
				os.RemoveAll(test.opts.RepoRoot)
			})
			assert.NoError(t, err)
			err = Run(test.opts)
			if test.errorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Check that any folder that was asked for removal, is no more present
			for _, driverVersion := range sliceDifference(test.driverVersionsToBeCreated, test.driverFolderRemainingExpected) {
				configPath := fmt.Sprintf(root.ConfigPathFmt,
					test.opts.RepoRoot,
					driverVersion,
					test.opts.Architecture,
					"")

				_, err = os.Stat(configPath)
				fmt.Println(err)
				assert.True(t, os.IsNotExist(err))
			}

			// Check that any folder that was NOT asked for removal, is still present
			for _, driverVersion := range test.driverFolderRemainingExpected {
				configPath := fmt.Sprintf(root.ConfigPathFmt,
					test.opts.RepoRoot,
					driverVersion,
					test.opts.Architecture,
					"")

				_, err = os.Stat(configPath)
				assert.NoError(t, err)
			}
		})
	}
}
