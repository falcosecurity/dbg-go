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
	"testing"
)

// NOTE: this test might be flaking because it tries to build some configs against a driverversion
// When it fails, just update configs to be built.
func TestBuild(t *testing.T) {
	// MUST RUN IN STRICT LOGICAL ORDER; USE A SLICE.
	tests := []struct {
		opts            Options
		name            string
		expectedObjects []string
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
				Publish:      false,
				IgnoreErrors: false,
			},
			expectedObjects: []string{
				"driver/5.0.1+driver/x86_64/falco_centos_5.14.0-284.11.1.el9_2.x86_64_1.ko",
				"driver/5.0.1+driver/x86_64/falco_centos_5.10.96-90.460.amzn2022.x86_64_1.o",
			},
			name: "build 5.0.1+driver centos 5.14.0-361.el9.x86_64",
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

	// Now, for each test, build the drivers then check s3 bucket objects
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// First, generate needed configs
			err = generate.Run(generate.Options{
				Options: test.opts.Options,
				Auto:    false,
			})
			assert.NoError(t, err)

			// Build the configs!
			err = Run(test.opts)
			assert.NoError(t, err)

			// Check the remaining objects in the bucket
			objects, err := testClient.ListObjects(context.Background(), &s3.ListObjectsInput{
				Bucket: aws.String(s3utils.S3Bucket),
			})
			assert.NoError(t, err)
			assert.Len(t, objects.Contents, len(test.expectedObjects))
			for _, obj := range objects.Contents {
				key := *obj.Key
				assert.Contains(t, test.expectedObjects, key)
			}
		})
	}
}
