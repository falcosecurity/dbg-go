package cleanup

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fededp/dbg-go/pkg/root"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func preCreateFolders(repoRoot, architecture string, driverVersionsToBeCreated []string) error {
	for _, driverVersion := range driverVersionsToBeCreated {
		configPath := fmt.Sprintf(root.ConfigPathFmt,
			repoRoot,
			driverVersion,
			architecture,
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
			err := preCreateFolders(test.opts.RepoRoot, test.opts.Architecture, test.driverVersionsToBeCreated)
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

func TestCleanupFiltered(t *testing.T) {
	tobeCreated := []string{
		"./test/driverkit/config/1.0.0+driver/x86_64/ubuntu_5.15.0_1.yaml",
		"./test/driverkit/config/1.0.0+driver/x86_64/ubuntu_5.19.2_1.yaml",
		"./test/driverkit/config/1.0.0+driver/x86_64/fedora_5.15.0_24.yaml",
		"./test/driverkit/config/1.0.0+driver/x86_64/talos_6.0.0_1.yaml",
		"./test/driverkit/config/1.0.0+driver/x86_64/amazonlinux_6.0.0_23.yaml",
	}

	err := preCreateFolders("./test", "x86_64", []string{"1.0.0+driver"})
	t.Cleanup(func() {
		os.RemoveAll("./test")
	})
	assert.NoError(t, err)

	// Create all empty files needed by the test
	for _, filepath := range tobeCreated {
		emptyFile, err := os.Create(filepath)
		assert.NoError(t, err)
		_ = emptyFile.Close()
	}

	tests := map[string]struct {
		opts                   Options
		expectedOutputContains []string
	}{
		"delete ubuntu configs": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "x86_64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					Distro:        "Ubun*",
					KernelRelease: "",
					KernelVersion: "",
				},
			}},
			expectedOutputContains: []string{"ubuntu_5.15", "ubuntu_5.19"},
		},
		"delete 24 kernelversion configs": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "x86_64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					Distro:        "",
					KernelRelease: "",
					KernelVersion: "24",
				},
			}},
			expectedOutputContains: []string{"fedora_5.15.0_24"},
		},
		"delete 6.0.0 configs": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "x86_64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					Distro:        "",
					KernelRelease: "6.0.0",
					KernelVersion: "",
				},
			}},
			expectedOutputContains: []string{"amazonlinux_6.0", "talos_6.0"},
		},
	}

	// Store logged data, will be used by test
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(io.Writer(&buf), nil))
	slog.SetDefault(logger)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err = Run(test.opts)
			assert.NoError(t, err)
			// Use logged output to ensure we really cleanup only correct configs:
			// parse every logged line to a structured json (we print "path:" for each config path being cleaned up)
			// then, for each parsed logged line, check if it contains one of the requested string by the test.
			// Count all "containing" lines; they must match total lines logged (that have a "config:" key).
			type MessageJSON struct {
				Path string `json:"config,omitempty"`
			}
			var messageJSON MessageJSON
			scanner := bufio.NewScanner(&buf)
			found := 0
			lines := 0
			for scanner.Scan() {
				err = json.Unmarshal(scanner.Bytes(), &messageJSON)
				assert.NoError(t, err)
				if messageJSON.Path == "" {
					continue
				}
				lines++
				for _, expectedOutput := range test.expectedOutputContains {
					if strings.Contains(messageJSON.Path, expectedOutput) {
						found++
						break
					}
				}
			}
			if found != lines {
				t.Errorf("wrong number of printed lines; expected %d, found %d", lines, found)
			}
			buf.Reset()
		})
	}

	// Check that we removed everything in the folder
	entries, err := os.ReadDir("./test/driverkit/config/1.0.0+driver/x86_64/")
	assert.NoError(t, err)
	assert.Len(t, entries, 0)
}
