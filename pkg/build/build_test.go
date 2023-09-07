package build

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fededp/dbg-go/pkg/generate"
	"github.com/fededp/dbg-go/pkg/root"
	s3utils "github.com/fededp/dbg-go/pkg/utils/s3"
	testutils "github.com/fededp/dbg-go/pkg/utils/test"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// NOTE: this test might be flaking because it tries to build some configs against a driver version.
// When it fails, just update configs to be built.
func TestBuild(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("only supported on amd64")
	}

	// MUST RUN IN STRICT LOGICAL ORDER; USE A SLICE.
	tests := []struct {
		opts                  Options
		name                  string
		expectedLocalObjects  []string
		expectedBucketObjects []string
		shouldCreate          bool
	}{
		{
			opts: Options{
				Options: root.Options{
					Architecture:  "amd64",
					DriverVersion: []string{"5.0.1+driver"},
					DriverName:    "falco",
					RepoRoot:      "./test",
					Target: root.Target{
						Distro:        "centos",
						KernelRelease: "5.14.0-361.el9.x86_64",
						KernelVersion: "1",
					},
				},
				SkipExisting: true,
				Publish:      true,
				IgnoreErrors: false,
			},
			expectedLocalObjects: []string{
				"falco_centos_5.14.0-361.el9.x86_64_1.ko",
				"falco_centos_5.14.0-361.el9.x86_64_1.o",
			},
			expectedBucketObjects: []string{
				"falco_centos_5.14.0-361.el9.x86_64_1.ko",
				"falco_centos_5.14.0-361.el9.x86_64_1.o",
			},
			shouldCreate: true,
			name:         "build 5.0.1+driver centos 5.14.0-361.el9.x86_64",
		},
		{
			opts: Options{
				Options: root.Options{
					Architecture:  "amd64",
					DriverVersion: []string{"5.0.1+driver"},
					DriverName:    "falco",
					RepoRoot:      "./test",
					Target: root.Target{
						Distro:        "centos",
						KernelRelease: "5.14.0-354.el9.x86_64",
						KernelVersion: "1",
					},
				},
				SkipExisting: true,
				Publish:      false,
				IgnoreErrors: false,
			},
			expectedLocalObjects: []string{
				"falco_centos_5.14.0-354.el9.x86_64_1.ko",
				"falco_centos_5.14.0-354.el9.x86_64_1.o",
				"falco_centos_5.14.0-361.el9.x86_64_1.ko",
				"falco_centos_5.14.0-361.el9.x86_64_1.o",
			},
			expectedBucketObjects: []string{
				"falco_centos_5.14.0-361.el9.x86_64_1.ko",
				"falco_centos_5.14.0-361.el9.x86_64_1.o",
			},
			shouldCreate: false, // since it is not publishing
			name:         "build 5.0.1+driver centos 5.14.0-354.el9.x86_64",
		},
		{
			opts: Options{
				Options: root.Options{
					Architecture:  "amd64",
					DriverVersion: []string{"5.0.1+driver"},
					DriverName:    "falco",
					RepoRoot:      "./test",
					Target: root.Target{
						Distro:        "centos",
						KernelRelease: "5.14.0-361.el9.x86_64", // try to rebuild same object.
						KernelVersion: "1",
					},
				},
				SkipExisting: true,
				Publish:      true,
				IgnoreErrors: false,
			},
			expectedLocalObjects: []string{
				"falco_centos_5.14.0-354.el9.x86_64_1.ko",
				"falco_centos_5.14.0-354.el9.x86_64_1.o",
				"falco_centos_5.14.0-361.el9.x86_64_1.ko",
				"falco_centos_5.14.0-361.el9.x86_64_1.o",
			},
			expectedBucketObjects: []string{
				"falco_centos_5.14.0-361.el9.x86_64_1.ko",
				"falco_centos_5.14.0-361.el9.x86_64_1.o",
			},
			shouldCreate: false, // since objects are already present, nothing should be created
			name:         "rebuild 5.0.1+driver centos 5.14.0-361.el9.x86_64",
		},
		{
			opts: Options{
				Options: root.Options{
					Architecture:  "amd64",
					DriverVersion: []string{"5.0.1+driver"},
					DriverName:    "falco",
					RepoRoot:      "./test",
					Target: root.Target{
						Distro:        "centos",
						KernelRelease: "5.14.0-361.el9.x86_64", // try to rebuild same object.
						KernelVersion: "1",
					},
				},
				SkipExisting: false, // this time, force-republish
				Publish:      true,
				IgnoreErrors: false,
			},
			expectedLocalObjects: []string{
				"falco_centos_5.14.0-354.el9.x86_64_1.ko",
				"falco_centos_5.14.0-354.el9.x86_64_1.o",
				"falco_centos_5.14.0-361.el9.x86_64_1.ko",
				"falco_centos_5.14.0-361.el9.x86_64_1.o",
			},
			expectedBucketObjects: []string{
				"falco_centos_5.14.0-361.el9.x86_64_1.ko",
				"falco_centos_5.14.0-361.el9.x86_64_1.o",
			},
			shouldCreate: true, // since objects are already present, nothing should be created
			name:         "rebuild 5.0.1+driver centos 5.14.0-361.el9.x86_64",
		},
	}

	// This client will be used by the Run action
	testClient = testutils.S3CreateTestBucket(t, nil)

	configPath := root.BuildConfigPath(root.Options{
		RepoRoot:     "./test",
		Architecture: "amd64",
	}, "5.0.1+driver", "")
	err := os.MkdirAll(configPath, 0700)
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll("./test/")
	})

	outputPath := root.BuildOutputPath(root.Options{
		RepoRoot:     "./test",
		Architecture: "amd64",
	}, "5.0.1+driver", "")
	err = os.MkdirAll(outputPath, 0700)
	assert.NoError(t, err)

	// Now, for each test, build the drivers then check s3 bucket objects
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			now := time.Now()
			// First, generate needed configs
			err = generate.Run(generate.Options{
				Options: test.opts.Options,
				Auto:    false,
			})
			assert.NoError(t, err)

			// Build the configs!
			err = Run(test.opts)
			assert.NoError(t, err)

			// Check that the files were created
			f, err := os.Open(outputPath)
			assert.NoError(t, err)
			t.Cleanup(func() {
				_ = f.Close()
			})
			entries, err := f.Readdirnames(0)
			assert.NoError(t, err)
			for _, e := range entries {
				assert.Contains(t, test.expectedLocalObjects, e)
			}

			if test.opts.Publish {
				// Check the remaining objects in the bucket
				objects, err := testClient.ListObjects(context.Background(), &s3.ListObjectsInput{
					Bucket: aws.String(s3utils.S3Bucket),
					Prefix: aws.String("driver/5.0.1+driver/x86_64/"),
				})
				assert.NoError(t, err)
				assert.Len(t, objects.Contents, len(test.expectedBucketObjects))
				for _, obj := range objects.Contents {
					key := filepath.Base(*obj.Key)
					lastMod := *obj.LastModified
					if test.shouldCreate {
						assert.True(t, lastMod.After(now))
					} else {
						assert.True(t, lastMod.Before(now))
					}
					assert.Contains(t, test.expectedBucketObjects, key)
				}
			}
		})
	}
}
