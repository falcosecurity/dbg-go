package cleanup

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/falcosecurity/dbg-go/pkg/root"
	s3utils "github.com/falcosecurity/dbg-go/pkg/utils/s3"
	"github.com/falcosecurity/dbg-go/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCleanup(t *testing.T) {
	tobeCreated := []string{
		"./test/driverkit/config/1.0.0+driver/x86_64/ubuntu_5.15.0_1.yaml",
		"./test/driverkit/config/1.0.0+driver/x86_64/ubuntu_5.19.2_1.yaml",
		"./test/driverkit/config/1.0.0+driver/x86_64/fedora_5.15.0_24.yaml",
		"./test/driverkit/config/1.0.0+driver/x86_64/talos_6.0.0_1.yaml",
		"./test/driverkit/config/1.0.0+driver/x86_64/amazonlinux_6.0.0_23.yaml",
	}

	err := testutils.PreCreateFolders(root.Options{
		RepoRoot:     "./test",
		Architecture: "amd64",
	}, []string{"1.0.0+driver"})
	t.Cleanup(func() {
		_ = os.RemoveAll("./test")
	})
	assert.NoError(t, err)

	// Create all empty files needed by the test
	for _, filepath := range tobeCreated {
		emptyFile, err := os.Create(filepath)
		assert.NoError(t, err)
		_ = emptyFile.Close()
	}

	tests := map[string]struct {
		opts          Options
		errorExpected bool
	}{
		"delete all": {
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				DriverName:    "falco",
			}},
			errorExpected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err = Run(test.opts, NewFileCleaner())
			if test.errorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Check that any file that was asked for removal, is no more present
			for _, driverVersion := range test.opts.DriverVersion {
				f, err := os.Open("./test/driverkit/config/" + driverVersion + "/x86_64/")
				assert.NoError(t, err)
				t.Cleanup(func() {
					_ = f.Close()
				})
				_, err = f.Readdirnames(1)
				assert.Equal(t, io.EOF, err)
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

	err := testutils.PreCreateFolders(root.Options{
		RepoRoot:     "./test",
		Architecture: "amd64",
	}, []string{"1.0.0+driver"})
	t.Cleanup(func() {
		_ = os.RemoveAll("./test")
	})
	assert.NoError(t, err)

	// Create all empty files needed by the test
	for _, filepath := range tobeCreated {
		emptyFile, err := os.Create(filepath)
		assert.NoError(t, err)
		_ = emptyFile.Close()
	}

	// MUST RUN IN STRICT LOGICAL ORDER; USE A SLICE.
	tests := []struct {
		opts                   Options
		expectedOutputContains []string
		name                   string
	}{
		{
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					Distro: "ubun*",
				},
				DriverName: "falco",
			}},
			expectedOutputContains: []string{"ubuntu_5.15", "ubuntu_5.19"},
			name:                   "delete ubuntu configs",
		},
		{
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					KernelVersion: "24",
				},
				DriverName: "falco",
			}},
			expectedOutputContains: []string{"fedora_5.15.0_24"},
			name:                   "delete 24 kernelversion configs",
		},
		{
			opts: Options{Options: root.Options{
				RepoRoot:      "./test/",
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					KernelRelease: "6.0.0",
				},
				DriverName: "falco",
			}},
			expectedOutputContains: []string{"amazonlinux_6.0", "talos_6.0"},
			name:                   "delete 6.0.0 configs",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			type MessageJSON struct {
				Path string `json:"config,omitempty"`
			}
			var messageJSON MessageJSON
			found := 0
			lines := 0
			testutils.RunTestParsingLogs(t, func() error {
				return Run(test.opts, NewFileCleaner())
			}, &messageJSON, func() bool {
				if messageJSON.Path == "" {
					return true // skip and go on
				}
				lines++
				for _, expectedOutput := range test.expectedOutputContains {
					if strings.Contains(messageJSON.Path, expectedOutput) {
						found++
						break
					}
				}
				return true
			})
			if found != lines {
				t.Errorf("wrong number of printed lines; expected %d, found %d", lines, found)
			}
		})
	}

	// Check that we removed everything in the folder
	entries, err := os.ReadDir("./test/driverkit/config/1.0.0+driver/x86_64/")
	assert.NoError(t, err)
	assert.Len(t, entries, 0)
}

func TestCleanupS3(t *testing.T) {
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
	cleaner := &s3Cleaner{Client: client}

	// MUST RUN IN STRICT LOGICAL ORDER; USE A SLICE.
	tests := []struct {
		opts             Options
		remainingObjects []string
		name             string
	}{
		{
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"2.0.0+driver"},
				DriverName:    "falco",
			}},
			remainingObjects: []string{
				"driver/1.0.0+driver/x86_64/falco_almalinux_5.14.0-284.11.1.el9_2.x86_64_1.ko",
				"driver/1.0.0+driver/x86_64/falco_amazonlinux2022_5.10.96-90.460.amzn2022.x86_64_1.o",
				"driver/1.0.0+driver/x86_64/falco_debian_6.3.11-1-amd64_1.o",
				"driver/1.0.0+driver/x86_64/falco_debian_6.3.11-1-amd64_1.ko",
				"driver/2.0.0+driver/aarch64/falco_almalinux_4.18.0-477.10.1.el8_8.aarch64_1.ko",
				"driver/2.0.0+driver/aarch64/falco_bottlerocket_5.10.165_1_1.13.1-aws.o",
			},
			name: "cleanup 2.0.0+driver x86_64",
		},
		{
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					Distro: "debian",
				},
				DriverName: "falco",
			}},
			remainingObjects: []string{
				"driver/1.0.0+driver/x86_64/falco_almalinux_5.14.0-284.11.1.el9_2.x86_64_1.ko",
				"driver/1.0.0+driver/x86_64/falco_amazonlinux2022_5.10.96-90.460.amzn2022.x86_64_1.o",
				"driver/2.0.0+driver/aarch64/falco_almalinux_4.18.0-477.10.1.el8_8.aarch64_1.ko",
				"driver/2.0.0+driver/aarch64/falco_bottlerocket_5.10.165_1_1.13.1-aws.o",
			},
			name: "cleanup 1.0.0+driver x86_64 debian",
		},
		{
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver"},
				Target: root.Target{
					Distro: "amazonlin*",
				},
				DriverName: "falco",
			}},
			remainingObjects: []string{
				"driver/1.0.0+driver/x86_64/falco_almalinux_5.14.0-284.11.1.el9_2.x86_64_1.ko",
				"driver/2.0.0+driver/aarch64/falco_almalinux_4.18.0-477.10.1.el8_8.aarch64_1.ko",
				"driver/2.0.0+driver/aarch64/falco_bottlerocket_5.10.165_1_1.13.1-aws.o",
			},
			name: "cleanup 1.0.0+driver x86_64 amazonlinux regex",
		},
		{
			opts: Options{Options: root.Options{
				Architecture:  "amd64",
				DriverVersion: []string{"1.0.0+driver", "2.0.0+driver"},
				Target: root.Target{
					KernelRelease: "5.*",
				},
				DriverName: "falco",
			}},
			remainingObjects: []string{
				"driver/2.0.0+driver/aarch64/falco_almalinux_4.18.0-477.10.1.el8_8.aarch64_1.ko",
				"driver/2.0.0+driver/aarch64/falco_bottlerocket_5.10.165_1_1.13.1-aws.o",
			},
			name: "cleanup 1.0.0+driver,2.0.0+driver x86_64 kernel release regex",
		},
		{
			opts: Options{Options: root.Options{
				Architecture:  "arm64",
				DriverVersion: []string{"2.0.0+driver"},
				DriverName:    "falco",
			}},
			remainingObjects: []string{},
			name:             "cleanup 2.0.0+driver aarch64 drivers",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Run(test.opts, cleaner)
			assert.NoError(t, err)

			// Check the remaining objects in the bucket
			objects, err := client.ListObjects(context.Background(), &s3.ListObjectsInput{
				Bucket: aws.String(s3utils.S3Bucket),
			})
			assert.NoError(t, err)
			for _, obj := range objects.Contents {
				key := *obj.Key
				assert.Contains(t, test.remainingObjects, key)
			}
		})
	}
}
